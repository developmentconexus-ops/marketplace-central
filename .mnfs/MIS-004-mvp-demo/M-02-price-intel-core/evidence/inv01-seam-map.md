## 1. ML adapter HTTP + auth

`apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:20-23`

```go
const (
    defaultBaseURL = "https://api.mercadolibre.com"
    defaultSiteID  = "MLB"
)
```

`.../capability_adapter.go:43-65`

```go
func NewCapabilityAdapter(cfg CapabilityAdapterConfig) *CapabilityAdapter {
    baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
    if baseURL == "" { baseURL = defaultBaseURL }
    httpClient := cfg.HTTPClient
    if httpClient == nil { httpClient = &http.Client{Timeout: 15 * time.Second} }
```

`.../capability_adapter.go:84-106`

```go
func (a *CapabilityAdapter) ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (...) {
    accountRef, err := normalizeAccountRef(ref)
    token, err := a.accessToken(ctx, accountRef)
    var response mlAccountResponse
    if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/me", nil, &response); err != nil {
```

`.../capability_adapter.go:452-464`

```go
func (a *CapabilityAdapter) accessToken(ctx context.Context, accountRef domain.ProviderAccountRef) (string, error) {
    if a.accessTokenResolver == nil { ... }
    token, err := a.accessTokenResolver(ctx, accountRef)
    token = strings.TrimSpace(token)
    if token == "" { ... }
```

`.../capability_adapter.go:522-542`

```go
req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, body)
req.Header.Set("Authorization", "Bearer "+token)
req.Header.Set("X-Tenant-ID", accountRef.TenantID)
req.Header.Set("X-Installation-ID", accountRef.InstallationID)
resp, err := a.httpClient.Do(req)
defer resp.Body.Close()
```

`.../capability_adapter.go:147-160`, `467-472`

```go
token, err := a.accessToken(ctx, accountRef)
item, err := a.readItem(ctx, accountRef, token, itemID)

func (a *CapabilityAdapter) readItem(...) (...) {
    if err := a.doJSON(ctx, accountRef, token, http.MethodGet,
        "/items/"+url.PathEscape(itemID), nil, &item); err != nil { ... }
```

`apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go:13-14,27-54`

```go
// FeeSyncer seeds Mercado Livre fee schedules with static defaults only.
// Slice 1 runtime fee evidence comes from live ReadFeeQuote, not this seeder.
func (f *FeeSyncer) Sync(...) (...) {
    ...
    return len(schedules), nil
}
```

No GET occurs in `fee_sync.go`.

## 2. Write path — untouched

`apps/server_core/internal/modules/connectors/adapters/mercado_livre/listing_writer.go:24-44`

```go
func (a *CapabilityAdapter) UpdateListing(...) (...) {
    ...
    token, err := a.accessToken(ctx, account)
    ...
    resp, raw, err := a.doRawWithIdempotency(ctx, account, token,
        http.MethodPut, "/items/"+url.PathEscape(request.ListingID), ...)
```

`apps/server_core/internal/modules/connectors/adapters/mercado_livre/price_writer.go:37-61`

```go
func (a *CapabilityAdapter) UpdatePrice(...) (...) {
    ...
    token, err := a.accessToken(ctx, accountRef)
    ...
    resp, rawBody, err := a.doRawWithIdempotency(ctx, accountRef, token,
        http.MethodPut, "/items/"+url.PathEscape(request.ListingID), ...)
```

Both reuse `CapabilityAdapter.httpClient`, `baseURL`, `accessToken`, and `doRawWithIdempotency`: `capability_adapter.go:35-40,452-464,518-548`.

## 3. Connectors ports pattern

`apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:22-29`

```go
type AccountProber interface {
    ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error)
}

type ListingReader interface {
    ListListings(ctx context.Context, input domain.ListListingsInput) ([]domain.ListingSnapshot, error)
    ReadListing(ctx context.Context, ref domain.ProviderListingRef) (domain.ListingSnapshot, error)
}
```

`apps/server_core/internal/modules/connectors/domain/capability.go:145-153`

```go
type AccountSnapshot struct {
    ProviderCode        string    `json:"provider_code"`
    ProviderAccountID   string    `json:"provider_account_id"`
    ProviderAccountName string    `json:"provider_account_name"`
    SiteID              string    `json:"site_id"`
    Status              string    `json:"status"`
    FetchedAt           time.Time `json:"fetched_at"`
    RawProviderRef      any       `json:"raw_provider_ref,omitempty"`
}
```

`apps/server_core/internal/modules/connectors/ports/me_auth.go:5-9`

```go
type MEAuthPort interface {
    HandleStart(http.ResponseWriter, *http.Request)
    HandleCallback(http.ResponseWriter, *http.Request)
    HandleStatus(http.ResponseWriter, *http.Request)
}
```

## 4. CredentialResolver / installation credentials

`apps/server_core/internal/modules/integrations/application/credential_resolver.go:11-15,25-30,45`

```go
type CredentialResolutionRef struct {
    TenantID, InstallationID, ProviderCode, ProviderAccountID string
}
type credentialLookup interface {
    GetActiveCredential(ctx context.Context, installationID string) (domain.Credential, bool, error)
}
type payloadDecryptor interface {
    DecryptJSON(encoded []byte) (map[string]any, string, error)
}
func (r *CredentialResolver) ResolveAccessToken(ctx context.Context,
    ref CredentialResolutionRef) (ResolvedCredential, error)
```

`.../credential_resolver.go:50-78`

```go
if ref.TenantID == "" || ref.InstallationID == "" ||
   ref.ProviderCode == "" || ref.ProviderAccountID == "" { ... }
credential, found, err := r.credentials.GetActiveCredential(ctx, ref.InstallationID)
payload, _, err := r.decryptor.DecryptJSON(credential.EncryptedPayload)
accessToken, ok := credentialPayloadString(payload, "access_token")
```

Implementation: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go:548-550`.

Tenant-scoped DB implementation: `apps/server_core/internal/modules/integrations/adapters/postgres/credential_repo.go:35-48`.

`MPC_ENCRYPTION_KEY` loading: `apps/server_core/internal/platform/pgdb/config.go:14-28`.

AES-GCM service: `apps/server_core/internal/modules/integrations/adapters/crypto/local_key.go:21-33,74-80,107-118`.

Adapter composition: `apps/server_core/internal/composition/root.go:330-343`.

## 5. Current `market` structure

Files and roles:

- `domain/market.go` — `Money`, observations, references, validation.
- `ports/store.go` — `ObservationStore`, `ReferenceStore`.
- `ports/collector.go` — collector port.
- `application/collection_service.go` — append orchestration.
- `application/read_service.go` — latest-read and synthetic empty evidence.
- `adapters/postgres/repository.go` — shared pool/tenant repository struct.
- `adapters/postgres/observation_repository.go` — observation append/read.
- `adapters/postgres/reference_repository.go` — reference append/read.
- `transport/query.go` — query parsing.
- `transport/http_handler.go` — `/market/observations`, `/market/references`.
- Remaining files are tests for the corresponding layers.

`apps/server_core/internal/modules/market/adapters/postgres/repository.go:5-12`

```go
type Repository struct {
    pool     *pgxpool.Pool
    tenantID string
}
func NewRepository(pool *pgxpool.Pool, tenantID string) *Repository {
    return &Repository{pool: pool, tenantID: tenantID}
}
```

`.../observation_repository.go:57-72`

```go
func (r *Repository) LatestObservations(ctx context.Context, listingIDs []string) (...) {
    rows, err := r.pool.Query(ctx, `
        SELECT DISTINCT ON (listing_id) ...
        FROM market_observations
        WHERE tenant_id = $1 AND listing_id = ANY($2)
        ORDER BY listing_id, captured_at DESC
    `, r.tenantID, listingIDs)
```

`.../reference_repository.go:45-56`

```go
rows, err := r.pool.Query(ctx, `
    SELECT DISTINCT ON (product_id) ...
    FROM market_references
    WHERE tenant_id = $1 AND product_id = ANY($2)
    ORDER BY product_id, captured_at DESC
`, r.tenantID, productIDs)
```

`apps/server_core/migrations/0043_market_observations.sql:1-18` and `0044_market_references.sql:1-18` define append-only tenant-keyed tables. `0044` primary key: `(tenant_id, product_id, captured_at)`.

## 6. Migration mechanics

`apps/server_core/migrations/source.go:8-16`

```go
// embedded contains the canonical, forward-only migration set.
//
//go:embed *.sql
var embedded embed.FS

func Source() fs.FS {
    return embedded
}
```

`apps/server_core/internal/platform/migrate/runner.go:12-26`

```go
func Filenames(source fs.FS) ([]string, error) {
    entries, err := fs.ReadDir(source, ".")
    ...
    if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
        filenames = append(filenames, entry.Name())
    }
    sort.Strings(filenames)
```

`.../runner.go:69-83`: applies pending files in lexicographic order; no down-migration path.

`apps/server_core/internal/platform/migrate/runner_test.go:16-26`

```go
wantMatches, err := filepath.Glob(filepath.Join("..", "..", "..", "migrations", "*.sql"))
...
for _, path := range wantMatches {
    want = append(want, filepath.Base(path))
}
sort.Strings(want)
if len(want) != 41 {
```

Exact count assertion: `apps/server_core/internal/platform/migrate/runner_test.go:25`

```go
if len(want) != 41 {
```

## 7. Composition root

Market module: `apps/server_core/internal/composition/root.go:531-533`

```go
marketModuleRepo := marketpostgres.NewRepository(pool, cfg.DefaultTenantID)
marketReadSvc := marketapp.NewReadService(marketModuleRepo, marketModuleRepo, time.Now)
markettransport.NewHandler(marketReadSvc).Register(mux)
```

Connectors capability construction: `.../root.go:330-347`.

```go
credentialResolver := integrationsapp.NewCredentialResolver(credentialSvc, encryptionSvc)
mercadoLivreCapabilities := mercadolivreconnector.NewCapabilityAdapter(...)
marketplaceCapabilities := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{
    mercadoLivreCapabilities.ProviderCapabilitySet(),
})
```

Connector transport registration: `.../root.go:656-657`.

```go
// Connectors (Melhor Envio auth + fee seeding foundations)
connectorstransport.NewHandler(meOAuth).Register(mux)
```

## 8. Telemetry / metrics

No Prometheus counter, metrics package, or counter helper found.

Existing telemetry is structured `slog` timing: `apps/server_core/internal/modules/internal_read/observability/timing.go:124-136`.

```go
started := time.Now()
err := call()
duration := time.Since(started)
attrs := []any{"method", method, "duration_ms", duration.Milliseconds()}
if err != nil {
    r.logger.Error("oracle_read", attrs...)
    return err
}
r.logger.Info("oracle_read", attrs...)
```

## 9. Tenant, decimal, time conventions

Tenant predicate: `apps/server_core/internal/modules/market/adapters/postgres/observation_repository.go:69-72`.

```sql
FROM market_observations
WHERE tenant_id = $1 AND listing_id = ANY($2)
```

Money uses decimal strings, not cents or `shopspring/decimal`: `apps/server_core/internal/modules/market/domain/market.go:11-20,268-279`.

```go
type Money struct {
    Amount   string `json:"amount"`
    Currency string `json:"currency"`
}
var decimalPattern = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)
```

Microsecond truncation appears in market integration tests: `apps/server_core/internal/modules/market/adapters/postgres/observation_repository_integration_test.go:96`.

```go
when := time.Now().UTC().Truncate(time.Microsecond)
```

## GAPS

- No existing counter/metrics helper with route/status/page labels.
- `fee_sync.go` has no HTTP GET.
- `credentialLookup` is unexported; public resolver access is through `CredentialResolver.ResolveAccessToken`.
- No separate Mercado Livre HTTP client file; HTTP is embedded in `CapabilityAdapter`.