# Lane: delivery

Repo `1473e863`. Scope: `apps/server_core` HTTP transport layer (20 modules under
`internal/modules/*/transport`, platform `internal/platform/{httpx,apierror}`,
`internal/composition/root.go`), `contracts/api/marketplace-central.openapi.yaml`,
`packages/sdk-runtime/src`. `internal/contexts/` has no `transport` package yet
(`find … -iname transport` under `internal/contexts` → 0 hits) — the whole delivery
surface still lives in the legacy `internal/modules/` tree.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| F1 | gap | No HTTP-layer authentication/authorization extraction exists anywhere. No auth middleware, no principal type, no `Authorization` header parsing outside provider OAuth adapters. | `grep -rln "Authorization\b"` in non-test, non-vendor Go → only `router.go` (CORS allow-list), `market_adapters.go`, `mercadolivre/auth_adapter.go`, `cmd/mlprobe` — all provider-outbound, none inbound-request. No `RequireAuth`/`AuthMiddleware`/`Authenticate` function exists repo-wide. | 0 implementations across 20 transport packages, 104+ route registrations |
| F2 | gap | No per-request tenant extraction. Tenant is a single value (`cfg.DefaultTenantID`) wired at composition time, not derived from the request. | `apps/server_core/internal/platform/pgdb/tenant.go:5-9` `DefaultTenantID` returns the config fallback or the literal `"tenant_default"`; `internal/composition/root.go` passes `cfg.DefaultTenantID` into every service constructor (e.g. lines 862, 865, 297-298, 521-522, 527) | 0 request-scoped tenant readers; ≥8 composition call sites hardwire the same config value |
| F3 | duplication+hazard | Status-code selection is reimplemented independently per module via string-prefix matching on `err.Error()`, never `errors.Is`/`errors.As`. | `mapOrdersError` `apps/server_core/internal/modules/orders/transport/http_handler.go:739-749` (`strings.HasPrefix(msg, "ORDERS_")` etc.); same pattern in `mapClassificationsError`, `mapIntegrationError` (`integrations/transport/http_handler.go:47`), `mapInventoryError`, `mapMarketplacesError`, `mapMutationError`, `mapFaturadoError`, `mapCalcError`, `mapPricingError`, `mapProductLinksError`, `mapError` (profitability) | 12 independent `map*Error` functions; `grep -c "errors.Is(err\|errors.As(err"` in `*/transport/*.go` → 0 |
| F4 | drift | Go 1.22+ method-prefixed routes (`"GET /x"`) get the stdlib's automatic 405 on a method mismatch — plain text, bypassing `apierror.Write` entirely — while bare-path routes (`"/x"` + manual `r.Method` check) always answer through the JSON envelope. Two different wire shapes for the same condition, and the difference is invisible unless you know which registration style a route uses. | stdlib: `C:\Program Files\Go\src\net\http\server.go:2707` — `Error(w, StatusText(StatusMethodNotAllowed), StatusMethodNotAllowed)` (Go 1.25.1, `go.mod:1`). Registration split verified by grep (see Dialect count). | 51 method-prefixed route registrations get this behavior; 47 bare-path registrations don't |
| F5 | duplication+drift | `parsePositiveInt(raw string, fallback int) int` is copy-pasted into 5 separate transport packages, and one copy has silently diverged in behavior. | Definitions: `integrations/transport/auth_handler.go:452`, `inventory/transport/http_handler.go:157`, `orders/transport/http_handler.go:751`, `product_links/transport/http_handler.go:514`, `profitability/transport/http_handler.go:199`. 4 of 5 use `strconv.Atoi`; `product_links/transport/http_handler.go:514-521` uses `fmt.Sscanf(value, "%d", &parsed)`, which (unlike `Atoi`) accepts trailing garbage (`"20xyz"` parses as `20` instead of falling back) | 5 copies, 1 behaviorally different from the other 4 |
| F6 | duplication | Cursor/pagination envelope shape (`items` + `next_cursor` + `page_size` [+ `as_of`]) is hand-copied as a private struct in each module rather than shared. | `catalogPageEnvelope` (`catalog/transport/http_handler.go:311`), `runPageEnvelope` (`integrations/transport/run_read_handler.go:37`), `listingPageEnvelope`/`listingGroupPageEnvelope` (`listings/transport/http_handler.go:150,156`), `mutationPageEnvelope` (`mutations/transport/query_handler.go:43`), `orderPageEnvelope`/`enrichedOrderPageEnvelope` (`orders/transport/http_handler.go:176,181`) | 7 independent struct definitions sharing the identical 3-field JSON shape by convention, no shared type |
| F7 | drift | Pagination is not universal: `erp_import` returns an unpaginated `items`-only list (`erp_import/transport/http_handler.go:223`), `orders/sankhya_linkage_handler.go:210` returns a flat `candidates` array with no cursor, and `market/collection_handler.go:64` returns a wholly custom envelope (`status`/`decisões`/`contagens`/`causas`) with a **Portuguese** JSON key (`json:"decisões"`) — the only non-English wire field name found in the transport layer. | `market/transport/collection_handler.go:64-68` | 3 divergent shapes alongside the 7-struct convention in F6 |
| F8 | drift | Money/decimal wire encoding has two incompatible conventions inside the same API: exact decimal-as-string (deliberate, ADR-driven) vs bare `float64`. | String: `pricing/domain/decimal.go:12-15` `Money{Amount string}` routed through `FormatRatHalfUp` (`decimal.go:46`); `catalog/transport/http_handler.go:339` `Amount *string`. Float64: `orders/domain/read_model.go:14` `Total *float64`; `orders/transport/sankhya_linkage_handler.go:221,230` `Amount *float64`; `profitability/domain/input.go:66` and `profitability/transport/http_handler.go:105` `Amount float64` | 2 dialects for the same concept (currency amount) inside one published API |
| F9 | gap | No request ID / correlation ID is assigned or propagated at the HTTP boundary. The only `request_id` in the codebase is a market-module *domain* field for ML price-collection evidence, unrelated to HTTP tracing. | `grep -rln "RequestID\|request_id\|correlation"` non-test/non-vendor → only `market/adapters/postgres/snapshot_repository.go`, `market/application/collection_pipeline_service.go`, `market/domain/evidence.go` (all business-domain, not transport) | 0 HTTP correlation-ID implementations |
| F10 | gap | Most JSON body decodes have no size limit; `http.MaxBytesReader` is used in only 2 of 20+ handler files. | `grep -rln "MaxBytesReader"` → `erp_import/transport/http_handler.go`, `orders/transport/sankhya_linkage_handler.go` only | 2/20 transport packages bound request body size |
| F11 | hazard | CORS is wide open (`Access-Control-Allow-Origin: *`) for every route including state-changing `POST`/`PUT`/`DELETE`, admitted in the code comment as dev-only but shipped unconditionally — no env gate. | `apps/server_core/internal/platform/httpx/router.go:20` | 1 site, applies to every one of the ~104 registered routes (global middleware) |
| F12 | drift | Body-decode-failure error **codes** are module-invented rather than shared, even though the message and mechanism (`apierror.Write`, `http.StatusBadRequest`, `"malformed request body"`) are consistent: `CATALOG_ENRICHMENT_INVALID`, `CLASSIFICATIONS_CREATE_INVALID`, `INTEGRATIONS_INSTALLATION_INVALID`, `INVENTORY_INVALID_REQUEST`, `MARKETPLACES_ACCOUNT_INVALID`, `MARKETPLACES_POLICY_INVALID`, `ORDERS_INVALID_REQUEST`, `PRICING_REQUEST_INVALID`, `PRODUCT_LINKS_INVALID_REQUEST`, `PROFITABILITY_INVALID_REQUEST`, plus `listings`' `invalid_request` and `mutations`' `invalid_body` — two of which don't even follow the `MODULE_THING_INVALID` naming convention the rest use. | grep transcript in Commands run | 12 distinct codes for one semantic condition (malformed JSON body) |
| F13 | hand-sync | The unified error envelope shape (`{"error":{"code","message","details"}}`) is right, but 2 of its 3 producers are hardcoded string literals that must be kept byte-identical to `apierror.envelope`'s Go struct by hand, not by sharing the type. | `httpx.WriteJSON`'s `encodeFailureBody` (`internal/platform/httpx/json.go:15`) and `route_deadline.go`'s `deadlineExceededBody` (`internal/platform/httpx/route_deadline.go:129`) — both acknowledge in their own comments they can't import `apierror` (cycle: `apierror -> httpx`) so they duplicate the JSON literal by hand instead | 2 hand-maintained copies of one struct's wire shape |
| F14 | gap | Contract: 10 operations declared in the OpenAPI spec **and** live in the router have zero SDK method reaching them — the frontend has no sanctioned way to call them. | See "Contract agreement" below for the full path list and verification that each is genuinely router-registered | 10 of 95 unique OpenAPI paths (≈11 of 111 operations, `/pricing/tariff-defaults` has 2 methods) |
| F15 | duplication | The SDK hand-rolls 5 near-identical query-string builder functions instead of one, plus several handwritten `encodeURIComponent` concatenations for the same job. | `catalogQuery`, `listingQuery`, `syncRunQuery`, `orderQuery`, `mutationQuery` — `packages/sdk-runtime/src/index.ts:2010,2021,2042,2052,2064`; ad hoc concatenation e.g. `index.ts:2227-2240` (`probeIntegrationFeeQuote`, `probeIntegrationStock`) | 5 builder functions + ≥4 inline hand-built query strings |

## The five heaviest, with detail

**1. F1/F2 — no authentication and no per-request tenant extraction (gap).**
Every one of the ~104 registered routes is reachable with no credential, and every
request is served as if it belongs to `cfg.DefaultTenantID` — a value fixed once at
process boot (`internal/composition/root.go`), never read from the request. This is
the delivery-lane half of a foundational gap: `AGENTS.md` states "tenant queries
scope `tenant_id`" as a binding rule, and the persistence layer (D-44 in the debt
ledger; `mpc_app` role is `NOLOGIN` and the app DSN bypasses RLS) already shows the
enforcement is contingent. At the HTTP boundary the story is stronger, not weaker: an
authenticated multi-tenant concept doesn't exist here at all yet — there is no code
that could enforce it even if RLS worked. Any remediation that only fixes RLS while
this layer stays as-is still ships a single-tenant-in-practice system wearing
multi-tenant persistence code.

**2. F3 — 12 independent status-code-selection functions, all via string matching.**
`mapOrdersError`, `mapClassificationsError`, `mapIntegrationError`,
`mapInventoryError`, `mapMarketplacesError`, `mapMutationError`, `mapFaturadoError`,
`mapCalcError`, `mapPricingError`, `mapProductLinksError`, `mapError`
(profitability) each parse `err.Error()` text (`strings.HasPrefix`/`HasSuffix`) to
decide an HTTP status and re-derive a code string. Zero use of `errors.Is`/`errors.As`
anywhere in the transport layer (`grep -c` → 0). A domain error's message wording is
now load-bearing for the HTTP status it produces — rename an error string in
`application` or `domain` and the transport layer silently starts answering `500`
instead of `400`/`404`/`409` for the same condition, with no compiler signal. This is
exactly the "wrong noun used correctly everywhere" failure mode PHASE-0.md's operator
section warns about, twelve times over, independently invented.

**3. F4 — the JSON error envelope has an accidental 4th, uncoded producer.**
`apierror.go`'s own doc comment claims "the wire shape cannot drift between modules"
(package doc, `apierror.go:1-3`) because every handler funnels through `Write`. That
claim is false for method-mismatch on the 51 routes registered with Go's
method-prefixed pattern syntax (`"GET /catalog/products"` etc): net/http's `ServeMux`
itself answers a bare-text `405 Method Not Allowed` (verified in the installed Go
1.25.1 toolchain source, `server.go:2707`) before any of this repo's code — including
`apierror.Recover` and `apierror.Write` — ever runs. A client parsing every error
response as `{"error":{"code","message","details"}}` will fail to parse this one.
Nobody wrote this bug; it is a side effect of two routing idioms (method-prefixed vs
bare-path-plus-manual-check) coexisting at near-even split (51/47) with different
implicit contracts.

**4. F14 — SDK is missing 10 real, implemented, spec-declared operations.**
`/admin/fee-schedules/{seed,sync}`, `/connectors/melhor-envio/auth/{start,callback}`,
`/integrations/auth/callback`, `/market/{aggregates,collections,signals,verdicts}`,
`/pricing/tariff-defaults` (GET+PUT) are all present in
`contracts/api/marketplace-central.openapi.yaml`, all confirmed live in
`internal/composition/root.go`'s router wiring (verified by grep against each
handler's own `mux.HandleFunc` call — see Commands run), and none of them has any
`getJson`/`postJson`/`putJson` call anywhere in `packages/sdk-runtime/src`. Because
`GOV_FRONTEND_FETCH` (`contracts/governance/invariants.json:13`,
`scripts/harness/Policy.psm1:467-474`) forbids raw `fetch()` in `apps/web/src` outside
`sdk-runtime` — and that gate is currently clean, 0 violations — these 10 operations
are not just "undocumented in the SDK," they are **unreachable from the web app under
the repo's own governance rule**, unless someone adds a governance exception. Market
signals/verdicts/aggregates and the fee-schedule admin actions are not small
surfaces.

**5. F5/F6/F7/F15 — pagination is a convention, not a mechanism, on both sides of the wire.**
Backend: `parsePositiveInt` copy-pasted 5 times, one copy (`product_links`) using
`fmt.Sscanf` instead of `strconv.Atoi` and therefore accepting trailing garbage the
other 4 reject — a hand-sync defect that has already drifted. The `items`/
`next_cursor`/`page_size` envelope is repeated as a private struct 7 times with no
shared type, and 3 more handlers (`erp_import`, `orders/sankhya_linkage_handler`,
`market/collection_handler`) don't follow it at all — one of them (`market`) uses
Portuguese JSON field names (`"decisões"`) found nowhere else in the transport layer.
Frontend: 5 separate query-builder functions in the SDK reimplement the same
`URLSearchParams` assembly, plus ad hoc `encodeURIComponent` string-building for
routes that didn't get a builder. None of this is wrong per-request; all of it is the
same fact (how a paginated request is encoded, how a paginated response is shaped)
kept synchronized by nine-plus independent authors across two languages, by hand.

## Dialect count

| concern | distinct implementations | sites |
|---|---|---|
| error response body shape | 1 real mechanism (`apierror.Write`) + 2 hand-copied literal duplicates of its shape + 1 accidental stdlib plain-text producer (F4) | 192 `apierror.Write` call sites; 2 literal copies (`json.go:15`, `route_deadline.go:129`); 51 routes exposed to the accidental 4th |
| status-code selection | 12 independent `map*Error` functions, all string-matching, 0 `errors.Is`/`errors.As` | 12 functions across 11 modules |
| request validation | 0 shared library; manual `json.Decode` + manual field checks | 19 `json.NewDecoder(r.Body)` sites; 44 manual `== ""` checks in `*/transport/*.go` |
| pagination | 2 backend styles (cursor-envelope convention vs ad hoc/none) × 5 duplicated `parsePositiveInt` (1 behaviorally different) × 7 duplicated envelope structs; 5 SDK query builders + ad hoc concatenation | see F5-F7, F15 |
| authentication/authorization extraction | 0 | 0 of ~104 routes |
| tenant extraction | 0 (single composition-time constant) | 0 request-scoped readers |
| request IDs / correlation | 0 | 0 |
| content negotiation / encoding | 1 (always `application/json`, no `Accept` handling — consistent, not a dialect) | n/a |
| time serialization | 1 (stdlib default RFC3339Nano via `encoding/json`, 0 custom `MarshalJSON`) | 0 custom marshalers found |
| decimal/money serialization | 2 (`string`+`big.Rat` exact vs bare `float64`) | see F8 |
| route registration idiom | 2 (Go 1.22 method-prefixed pattern vs bare-path + manual `r.Method` check) | 51 method-prefixed / 47 bare-path (104 total registration call sites) |
| middleware chain | 1 (`CORS → apierror.Recover → RouteClassMux[deadline-per-route] → handler`), applied uniformly — no bypass path found | 1 chain, `root.go:994`, threaded through all 20 `Register(mux)` implementors |

## Contract agreement

- **OpenAPI operations:** 111 unique `(method, path)` pairs, 95 unique paths, 111 unique `operationId` values (1:1 with method+path — no duplicate or missing operationIds). Counted from `contracts/api/marketplace-central.openapi.yaml` (8574 lines) by parsing `paths:` block structure directly (method line indentation), not by trusting prose.
- **Router:** 20 `Register(mux httpx.RouteRegistrar)` implementations (one per `internal/modules/*/transport` package that exposes HTTP, plus `connectors/adapters/melhorenvio/oauth.go`), ~104-105 individual `mux.Handle`/`HandleFunc`/`registerInteractiveRoute` call sites. Route classes (`interactive`/`batch`) are declared centrally in `registerBatchRoutes` (`root.go:268-281`, 8 batch patterns) and default to `interactive` otherwise.
- **SDK:** 104 methods in the object returned by `createMarketplaceCentralClient` (`packages/sdk-runtime/src/index.ts:2137-2471`), 172 exported interfaces/types (matches established fact #5), 2595 total lines.
- **Join key used:** *not* `operationId` vs SDK method name (established fact warns this join is wrong — verified true: e.g. spec operationId `listCatalogProductFacts` and SDK method `listCatalogProductFacts` happen to match by convention in some cases but nothing enforces it). Joined instead on the literal HTTP path template, normalizing both `{param}` (OpenAPI) and `${var}` (SDK template literals) to a wildcard, and checking static-segment containment — the actual wire-compatible fact.
- **Deltas found by path-based join:**
  - **In OpenAPI + router, absent from SDK (10 paths, F14):** `/admin/fee-schedules/seed`, `/admin/fee-schedules/sync`, `/connectors/melhor-envio/auth/callback`, `/connectors/melhor-envio/auth/start`, `/integrations/auth/callback`, `/market/aggregates`, `/market/collections`, `/market/signals`, `/market/verdicts`, `/pricing/tariff-defaults`. All 10 verified present in router wiring by direct grep against each handler's registration call (see Commands run) — these are not stale spec entries, they are live, callable, undocumented-to-the-SDK endpoints.
  - **In SDK, absent from OpenAPI:** 0 found (containment check: every static path prefix referenced by an SDK `getJson`/`postJson`/`putJson` call in `index.ts` matches some OpenAPI path prefix). The SDK does not reach past what the spec declares.
  - `GOV_API_SDK_SPLIT` (established fact 6, `invariants.json:12`, `api-sdk-atomicity` check) only requires the spec and `packages/sdk-runtime` to change in the same commit — it has no mechanism to catch a spec-and-router operation that never got an SDK method added at all, because from the governance check's point of view "the SDK changed" is satisfied by *any* SDK-tree edit in that commit, not by adding a matching method. This 10-operation gap is exactly what that rule cannot see by construction.
  - **Response-shape agreement, handlers vs spec:** not exhaustively diffed field-by-field against the YAML schemas (that is a distinct, larger measurement — 111 operations × full schema comparison) — flagged under Unverified. What was verified directly against handler code (not annotations) is the pagination-envelope divergence (F6/F7) and the money-type divergence (F8), both of which are real shape facts a schema diff would also have to reconcile.

## What is actually fine

- **`apierror.Write` is a genuine single producer for the cases that reach it.** 192 call sites, 0 raw `http.Error`/`json.NewEncoder` bypasses found anywhere in first-party Go code (only in `.gomodcache` vendor trees). CHIP-ERROR-UNIFY (referenced in `.mnfs/HARNESS-DEBTS.md` section E and D-7/D-8/D-9) evidently landed: this is not the "2 families + 4 writeError" state the debt ledger describes from 2026-07-31 — today it's 1 mechanism with 19 thin per-module wrapper functions that all delegate to it (verified by reading all 19: `writeCatalogPageError`, `writeDashboardError`, `writeListError`, `writeMarketError`, `writeMutationError`, etc. — every one bottoms out in `apierror.Write`). Do not re-litigate "unify the error writer"; it's done. What remains (F3, F4, F12, F13) is one layer deeper than the ledger's snapshot.
- **`GOV_FRONTEND_FETCH` is a real, currently-clean, mechanically-enforced gate**, and it is precisely the kind of control PHASE-0.md's operator section asks for ("a mechanism that makes a bad change hard to land," not a discipline rule). 0 raw `fetch()` calls found in `apps/web/src` or `packages/*` outside `sdk-runtime` itself. Worth citing as a positive existence proof, not just a finding to fix elsewhere.
- **Time serialization is uniform.** 0 custom `MarshalJSON` implementations found in the transport layer; every `time.Time` field goes through `encoding/json`'s default RFC3339Nano. Not a dialect.
- **`pricing/domain.Money` (string + `big.Rat`, `FormatRatHalfUp`) is deliberately engineered, not accidental** — ADR-17-aware doc comments, an explicit decimal-string grammar (`decimalStringPattern`), and half-up rounding at a stated precision. This is the correct model to converge *other* money fields toward, not something to touch itself.
- **Deadline middleware has no bypass path.** Every module's `Register(mux)` receives the same `*httpx.RouteClassMux`, and `RouteClassMux.Handle` unconditionally wraps every registration in `deadlineMiddleware` — there is no second code path into the mux that skips it. The interactive/batch split is declared centrally (`registerBatchRoutes`) rather than scattered.
- **SDK never calls outside the spec.** The delta at the contract boundary is one-directional (SDK undershoots what's declared+built, F14) — there is no evidence of the SDK reaching an undocumented/shadow endpoint.

## Unverified / needs judgment

- **Full per-operation response-schema diff (handler struct vs OpenAPI `schema:` block), all 111 operations.** Only spot-checked (catalog page envelope, orders envelope, error envelope). A complete diff needs either `oapi-codegen`-generated types compared against hand-written response structs, or a field-by-field manual pass — both out of this lane's budget. Established facts 5-7 flag this as the reason a naive "does it compile" check can't be trusted; this lane confirms the risk is real (F6-F8) without sizing it exhaustively.
- **The 47 vs 51 route-registration-style split (F4) is a call-site count, not a verified-distinct-route count.** A few call sites register the same logical path under both styles in different modules (e.g., `/connectors/melhor-envio/auth/start` is registered bare-path in two different files — `oauth.go:69` and `connectors/transport/http_handler.go:38` — which is itself worth a second look but wasn't chased further here: is one dead code, or do both run and the second one silently wins/loses under `http.ServeMux`'s pattern-conflict rules?). Flagging as a lead for the `duplication` or `layering` lane rather than re-deriving here.
- **Whether the 10 SDK-missing operations (F14) are intentionally admin/internal-only** (e.g. `/admin/fee-schedules/*` reads as an operator tool, not something the web app should expose) is a judgment call this lane does not make. The fact reported is narrower and load-bearing regardless of intent: they are in the *public* OpenAPI spec with no `x-internal`-style marker found, and `GOV_FRONTEND_FETCH` would currently block anyone from wiring them up via raw `fetch` if a future PR needed them from the web app.
- **Whether `market/collection_handler.go`'s Portuguese field name and bespoke envelope (F7) is deliberate** (matches a Portuguese-speaking internal consumer) or an oversight — not established either way; reported as a fact (the only non-English wire key found), not a verdict.

## Commands run

```
git status --porcelain --untracked-files=all
find apps/server_core -type d -iname "transport" | sort
find apps/server_core/internal/contexts -type d -iname "transport"
find apps/server_core/internal/platform -maxdepth 2 -type d
find apps/server_core -iname "*router*.go" | grep -v _test
find . -iname "*openapi*" -not -path "*/node_modules/*" -not -path "*/scripts/.runs/*"
find . -path "*sdk-runtime*" -iname "*.ts" | grep -v node_modules
grep -rn "apierror.Write(" apps/server_core --include=*.go | grep -v _test | wc -l          # 192
grep -rn "httpx.WriteJSON(" apps/server_core --include=*.go | grep -v _test | wc -l          # 112
grep -rln "http.Error(" apps/server_core --include=*.go | grep -v _test                      # 0 (only .gomodcache)
grep -rn "^func write[A-Za-z]*Error\|^func writeError" apps/server_core --include=*.go | grep -v _test   # 19 wrappers
grep -rn "^func map[A-Za-z]*Error" apps/server_core --include=*.go | grep -v _test           # 12 in transport
grep -rn "strings.HasPrefix(msg\|errors.Is(err\|errors.As(err" apps/server_core --include=*/transport/*.go | grep -v _test | wc -l   # 0 errors.Is/As
grep -rln "json.NewDecoder(r.Body)" apps/server_core --include=*.go | grep -v _test | wc -l  # 19
grep -rln "== \"\"" apps/server_core/internal/modules/*/transport | wc -l                    # 44
grep -rn "^func parsePositiveInt" apps/server_core --include=*.go | grep -v _test            # 5 definitions
grep -rn "type.*[Pp]ageEnvelope\|type.*[Ll]istResponse\|type.*[Cc]ollectionResponse" apps/server_core --include=*.go | grep -v _test
grep -rln "Authorization\b" apps/server_core --include=*.go | grep -v _test | grep -v .gomodcache
grep -rn "tenant.ID(\|WithTenant(" apps/server_core --include=*.go | grep -v _test           # 0 request-scoped
grep -rln "X-Request-Id\|RequestID\|request_id\|correlation" apps/server_core --include=*.go | grep -v _test
grep -rln "MaxBytesReader" apps/server_core --include=*.go | grep -v _test                   # 2 files
GOROOT=$(go env GOROOT) && grep -n "StatusMethodNotAllowed" "$GOROOT/src/net/http/server.go"  # server.go:2707
grep -m1 "^go " apps/server_core/go.mod                                                       # go 1.25.1
grep -rlE '"(GET|POST|PUT|DELETE|PATCH) /' apps/server_core --include=*.go | grep -v _test | grep -v .gomodcache   # 17 files, method-prefixed style
grep -rn 'mux\.\(Handle\|HandleFunc\)("/\|registerInteractiveRoute(mux, "/' apps/server_core --include=*.go | grep -v _test | grep -v .gomodcache | wc -l   # 47 bare-path
grep -rn 'mux\.\(Handle\|HandleFunc\)("\(GET\|POST\|PUT\|DELETE\|PATCH\) ' apps/server_core --include=*.go | grep -v _test | grep -v .gomodcache | wc -l   # 51 method-prefixed
python3 <<'PY'  # OpenAPI path/operation extraction from contracts/api/marketplace-central.openapi.yaml
# parses `  /path:` and `    method:` lines directly -> 95 unique paths, 111 unique (method,path)
PY
python3 <<'PY'  # SDK path-literal extraction from packages/sdk-runtime/src/index.ts return block (lines 2137-2471)
# regex over getJson<T>(`...`) / postJson<T>(`...`) etc -> 104 methods, 80 inline path literals
PY
python3 <<'PY'  # static-segment containment diff, OpenAPI paths vs SDK source text
# -> 10 OpenAPI paths with zero string presence in index.ts; 0 SDK-only paths
PY
grep -rn "/admin/fee-schedules/seed\|/admin/fee-schedules/sync\|/connectors/melhor-envio/auth/callback\|/connectors/melhor-envio/auth/start\|/integrations/auth/callback\|/market/aggregates\|/market/collections\|/market/signals\|/market/verdicts\|/pricing/tariff-defaults" apps/server_core --include=*.go | grep -v _test | grep -v .gomodcache | grep -iE 'mux\.|HandleFunc'   # confirms all 10 are router-live
grep -n "GOV_FRONTEND_FETCH\|frontend-fetch" contracts/governance/invariants.json scripts/harness/Policy.psm1
grep -rln "[^A-Za-z0-9_]fetch(" apps/web/src packages --include=*.ts --include=*.tsx --include=*.js --include=*.mjs | grep -v "sdk-runtime/"   # 0 hits
grep -rln "go-playground/validator\|validator.New()" apps/server_core --include=*.go | grep -v _test
grep -rln "func.*MarshalJSON\|time.Format(" apps/server_core --include=*/transport/*.go | grep -v _test   # 0
wc -l contracts/api/marketplace-central.openapi.yaml    # 8574
wc -l packages/sdk-runtime/src/index.ts                  # 2595
grep -c "^export interface\|^interface " packages/sdk-runtime/src/index.ts   # 172
```
