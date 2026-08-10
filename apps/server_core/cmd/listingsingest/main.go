// Command listingsingest is the operator entry point that drives the Mercado
// Livre listing feed into the listings context — the listings-slice
// counterpart to cmd/catalogingest. It is read-only against ML (GETs only:
// scan + multiget, plus at most one /users/me lookup to resolve the seller
// id) and writes only to the listings schema. The access token is resolved
// the same way cmd/mlprobe resolves it (Postgres + AES-GCM local key) and is
// never printed, logged, or embedded in an error.
//
// Which tenant to ingest for and how big a scan page to request are
// parameters of the run, so they arrive as flags:
//
//	go run ./cmd/listingsingest -tenant tenant_default [-page-size 50]
//
// The environment carries only deployment configuration — the database URL
// and the encryption key — read once by pgdb.LoadConfig.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

const (
	commandName     = "listingsingest"
	defaultPageSize = 50
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx, os.Args[1:]); err != nil {
		// -h is a request, not a failure: the usage text is already on stderr
		// and this process has nothing to report.
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		// Never a fabricated success: any failure here is printed verbatim to
		// stderr and the process exits non-zero.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// options are the parameters of one run — what this invocation is for, as
// opposed to how the deployment is configured.
type options struct {
	tenant   string
	pageSize int
}

// parseOptions reads the run's parameters from args and nothing else. There
// is deliberately no environment fallback for either of them: a required flag
// has no default that a forgotten export could be mistaken for, so the
// fail-closed behaviour a live write path needs is a property of the channel
// rather than a guard that has to be remembered (global constraint 9 —
// unknown never becomes a plausible default).
//
// usageOut receives the usage text; the FlagSet's own printer is silenced so
// that every failure leaves this process through main's single error path
// instead of being written twice.
func parseOptions(args []string, usageOut io.Writer) (options, error) {
	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	tenantFlag := fs.String("tenant", "", "tenant id this run ingests for (required)")
	pageSizeFlag := fs.Int("page-size", defaultPageSize, "listing ids requested per scan page")

	if err := fs.Parse(args); err != nil {
		fs.SetOutput(usageOut)
		fs.Usage()
		return options{}, err
	}
	if fs.NArg() > 0 {
		return options{}, fmt.Errorf("unexpected argument %q (this command takes flags only)", fs.Arg(0))
	}
	tenantID := strings.TrimSpace(*tenantFlag)
	if tenantID == "" {
		return options{}, errors.New("-tenant is required: the tenant to ingest for is a parameter of this run, not a deployment setting")
	}
	if *pageSizeFlag <= 0 {
		return options{}, fmt.Errorf("-page-size must be a positive integer, got %d", *pageSizeFlag)
	}
	return options{tenant: tenantID, pageSize: *pageSizeFlag}, nil
}

func run(ctx context.Context, args []string) error {
	opts, err := parseOptions(args, os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return err
		}
		return fmt.Errorf("listingsingest: %w", err)
	}
	tenantID, err := tenant.Parse(opts.tenant)
	if err != nil {
		return fmt.Errorf("listingsingest: tenant: %w", err)
	}

	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		return fmt.Errorf("listingsingest: database config: %w", err)
	}

	pool, err := pgdb.NewPool(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("listingsingest: open postgres pool: %w", err)
	}
	defer pool.Close()

	// dbCfg.EncryptionKey, not a second direct environment read of the same
	// encryption-key variable: pgdb.LoadConfig already fails closed on it
	// above, and re-reading it here could only ever agree with what it saw
	// or never run — re-validating an already-validated value would be dead
	// weight, not an extra safeguard.
	token, accountID, credTenant, err := resolveMLCredential(ctx, pool, dbCfg.EncryptionKey)
	if err != nil {
		return fmt.Errorf("listingsingest: resolve marketplace credential: %w", err)
	}
	// A credential from any tenant but the one this process is configured for
	// must never drive a write: cross-tenant use is not a bug to log, it is a
	// request this command refuses outright.
	if credTenant != tenantID.String() {
		return fmt.Errorf("listingsingest: connected marketplace credential belongs to tenant %q, refusing to run for configured tenant %q", credTenant, tenantID.String())
	}

	wiring, err := composition.WireListings(pool, composition.MLBaseURL, accountID, accountID, func(context.Context) (string, error) {
		return token, nil
	})
	if err != nil {
		return fmt.Errorf("listingsingest: wire listings: %w", err)
	}

	report, err := composition.RunListingsIngest(ctx, wiring.Module, wiring.Feed.ListingFeed, tenantID, opts.pageSize)
	if err != nil {
		return fmt.Errorf("listingsingest: run ingest: %w", err)
	}

	fmt.Printf("listings ingest report: pages=%d observed=%d created=%d changed=%d idempotent=%d\n",
		report.Pages, report.Observed, report.Created, report.Changed, report.Idempotent)
	return nil
}

// resolveMLCredential resolves the live access token and the seller/account
// id for the connected mercado_livre installation, plus the tenant that
// installation belongs to (so the caller can refuse a cross-tenant
// credential). The SQL join, the provider/status/is_active/revoked_at
// filter, the "most recent version wins" ordering, and the
// crypto.NewLocalKeyService(key, "local-key-v1").DecryptJSON call are the
// measured precedent from cmd/mlprobe/main.go's resolveToken (~lines
// 164-190) — same shape, same keys, same fail-closed reads.
//
// mlprobe's decrypted payload was checked directly (not assumed): it carries
// "access_token" but the seller id is NOT stored there under "user_id" —
// integrations/application/auth_flow_service.go persists it as
// "provider_account_id", and it is also mirrored onto
// integration_installations.external_account_id at connect time (0016, 0017:
// non-empty is enforced for any installation with status NOT IN
// ('disconnected','failed')). That column is read here alongside the
// credential row, in the same query, so the account id is resolved from
// already-stored state whenever possible; the payload's "provider_account_id"
// is the second-line fallback; only if both are empty does this command make
// the one live GET /users/me call mlprobe makes, using the "id" field of its
// response exactly as mlprobe does.
func resolveMLCredential(ctx context.Context, pool *pgxpool.Pool, key string) (token, accountID, credTenant string, err error) {
	// provider_code is bound as a parameter, not inlined: the vendor name is
	// composition.MLProviderCode, not a literal repeated here (RuleVendorToken,
	// internal/arch/scan.go — see the doc comment on that constant).
	row := pool.QueryRow(ctx, `
		SELECT i.installation_id, i.tenant_id, i.external_account_id, c.encrypted_payload
		FROM integration_installations i
		JOIN integration_credentials c
		  ON c.tenant_id = i.tenant_id AND c.installation_id = i.installation_id
		 AND c.is_active = true AND c.revoked_at IS NULL
		WHERE i.provider_code = $1 AND i.status = 'connected'
		ORDER BY c.version DESC LIMIT 1`, composition.MLProviderCode)
	var installationID, externalAccountID string
	var payload []byte
	if err := row.Scan(&installationID, &credTenant, &externalAccountID, &payload); err != nil {
		return "", "", "", fmt.Errorf("no connected marketplace credential: %w", err)
	}

	svc, err := crypto.NewLocalKeyService(key, "local-key-v1")
	if err != nil {
		return "", "", "", fmt.Errorf("encryption key: %w", err)
	}
	decoded, _, err := svc.DecryptJSON(payload)
	if err != nil {
		return "", "", "", fmt.Errorf("decrypt credential: %w", err)
	}
	tok, _ := decoded["access_token"].(string)
	if strings.TrimSpace(tok) == "" {
		return "", "", "", errors.New("access_token missing in decrypted credential payload")
	}

	accountID = strings.TrimSpace(externalAccountID)
	if accountID == "" {
		accountID = credentialAccountIDFromPayload(decoded)
	}
	if accountID == "" {
		accountID, err = fetchMLSellerID(ctx, tok)
		if err != nil {
			return "", "", "", fmt.Errorf("resolve seller id: %w", err)
		}
	}
	return tok, accountID, credTenant, nil
}

// credentialAccountIDFromPayload reads the seller/account id out of the
// decrypted credential payload map, trying both keys observed across this
// codebase's ML integration: "provider_account_id" (the key
// auth_flow_service.go actually persists) and "user_id" (the raw field name
// ML's own OAuth token response uses, in case an older or differently-shaped
// payload carries it verbatim). Either may arrive as a string or a JSON
// number, since encoding/json decodes an untyped payload's numeric ids as
// float64.
func credentialAccountIDFromPayload(payload map[string]any) string {
	for _, key := range []string{"provider_account_id", "user_id"} {
		switch v := payload[key].(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				return s
			}
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		}
	}
	return ""
}

// fetchMLSellerID performs the one-shot GET /users/me lookup mlprobe uses to
// derive the connected account's own seller id, for the case where neither
// the installation row nor the decrypted payload already carried it. This
// lives in cmd/, not in the mercadolivre adapter: it is a one-time bootstrap
// concern of resolving "who is this credential for", not a listing feed
// operation the adapter's port needs to expose.
func fetchMLSellerID(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, composition.MLBaseURL+"/users/me", nil)
	if err != nil {
		return "", fmt.Errorf("build /users/me request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET /users/me: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read /users/me: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		// The status and body say what went wrong; the token never appears —
		// it is only ever sent in the Authorization header, never echoed.
		return "", fmt.Errorf("GET /users/me: status %d", resp.StatusCode)
	}
	var me struct {
		ID json.Number `json:"id"`
	}
	if err := json.Unmarshal(body, &me); err != nil {
		return "", fmt.Errorf("decode /users/me: %w", err)
	}
	id := strings.TrimSpace(me.ID.String())
	if id == "" {
		return "", errors.New("/users/me response carried no id")
	}
	return id, nil
}
