// Command listingsingest is the operator entry point that drives the Mercado
// Livre listing feed into the listings context — the listings-slice
// counterpart to cmd/catalogingest. It is read-only against ML (GETs only:
// scan + multiget, plus at most one /users/me lookup to resolve the seller
// id) and writes only to the listings schema. The access token is resolved
// the same way cmd/mlprobe resolves it (Postgres + AES-GCM local key) and is
// never printed, logged, or embedded in an error.
package main

import (
	"context"
	"encoding/json"
	"errors"
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

const defaultPageSize = 50

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		// Never a fabricated success: any failure here is printed verbatim to
		// stderr and the process exits non-zero.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dbCfg, err := pgdb.LoadConfig()
	if err != nil {
		return fmt.Errorf("listingsingest: database config: %w", err)
	}
	if err := requireTenantConfigured(os.Getenv); err != nil {
		return err
	}
	tenantID, err := tenant.Parse(dbCfg.DefaultTenantID)
	if err != nil {
		return fmt.Errorf("listingsingest: tenant: %w", err)
	}
	pageSize, err := loadPageSize(os.Getenv)
	if err != nil {
		return fmt.Errorf("listingsingest: %w", err)
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

	report, err := composition.RunListingsIngest(ctx, wiring.Module, wiring.Feed.ListingFeed, tenantID, pageSize)
	if err != nil {
		return fmt.Errorf("listingsingest: run ingest: %w", err)
	}

	fmt.Printf("listings ingest report: pages=%d observed=%d created=%d changed=%d idempotent=%d\n",
		report.Pages, report.Observed, report.Created, report.Changed, report.Idempotent)
	return nil
}

// requireTenantConfigured fails closed when MC_DEFAULT_TENANT_ID is unset or
// empty. pgdb.LoadConfig (internal/platform/pgdb/config.go) silently
// substitutes "tenant_default" for every caller when the variable is absent
// — correct for its many read-oriented callers, but wrong here: this command
// performs live writes across an entire seller's listings (D-39,
// .mnfs/HARNESS-DEBTS.md). An operator who forgets to export the variable
// must see a failure naming it, not have thousands of rows land under an
// invented tenant with a silent exit code 0 (global constraint 9: unknown
// never becomes a plausible default).
//
// getenv is read directly here, exactly mirroring how pgdb.LoadConfig reads
// the same variable (os.Getenv("MC_DEFAULT_TENANT_ID"), no trimming), so
// this guard can never disagree with what LoadConfig saw. It keys on the
// variable being absent/empty, never on the resolved value: a deployment
// that legitimately names its tenant "tenant_default" still passes, and a
// typo'd variable name still fails.
//
// pgdb.LoadConfig's defaulting itself is not changed — it is shared
// infrastructure with other callers (cmd/server) that may have distinct
// reasons to tolerate the default today.
func requireTenantConfigured(getenv func(string) string) error {
	if getenv("MC_DEFAULT_TENANT_ID") == "" {
		return errors.New("listingsingest: MC_DEFAULT_TENANT_ID is required (refusing to silently default to \"tenant_default\" for a live write path)")
	}
	return nil
}

// loadPageSize reads the operator-tunable page size. Unlike a domain fact, a
// paging batch size has no "unknown" state to preserve — it is an operational
// default, not a fabricated business value — so an unset variable falls back
// to defaultPageSize rather than failing closed.
func loadPageSize(getenv func(string) string) (int, error) {
	raw := strings.TrimSpace(getenv("MPC_LISTINGS_INGEST_PAGE_SIZE"))
	if raw == "" {
		return defaultPageSize, nil
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("MPC_LISTINGS_INGEST_PAGE_SIZE must be a positive integer, got %q", raw)
	}
	return parsed, nil
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
