## VERDICT

REFUTED

Two BLOCKING defects are reachable in the already-merged endpoint. They now require corrective commits on `main`: valid zero-padded CODPRODs can be omitted from `vinculados`, and malformed UUID path parameters return 500.

## PART A — C1..C11

| Criterion | Verdict | Evidence (file:line + quoted line) |
|---|---|---|
| C1 | PASS | `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:40`: `var knownIdentityAnchors = []IdentityAnchor{` contains only seller SKU, EAN, title, and marca through line 44. ERP data remains at `apps/server_core/migrations/0046_create_erp_import_products.sql:10`: `refforn TEXT,`. An exhaustive `git grep` found no `IdentityAnchorRefforn` in `apps`, `contracts`, or `packages`. |
| C2 | PASS | The four classifications are implemented at `generation_service.go:704`: `return domain.LinkCandidateReasonDirectionUnavailable, "", ...`; `:723`: `return domain.LinkCandidateReasonDirectionIncomparable, domain.LinkCandidateReasonSideBoth, ...`; `:726`: `return ...SideProvider...`; `:728`: `return ...SideERP...`. Presence in the result is asserted at `generation_service_test.go:356`: `reason, ok := findReason(result.Items[0].Reasons, "title", tc.direction)`. |
| C3 | PASS | `apps/server_core/internal/modules/product_links/domain/link_candidate.go:72`: `Side LinkCandidateReasonSide json:"side,omitempty"`. Production assignments occur only in incomparable classification paths. The serialization guard checks absence at `generation_service_test.go:540`: `if strings.Contains(string(payload), '"side"') {`. |
| C4 | PASS | D-121 remains structurally isolated from reasons: concordant CODPROD+EAN assigns ACCEPT at `generation_service.go:505`: `candidate.MatchStatus = domain.LinkCandidateMatchStatusAccept`, then reasons are appended at `:507`: `candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, comparison)`. Single-anchor paths remain CONFIRM at `:529` and `:539`. The automatic resolver checks only ACCEPT at `resolution_service.go:231`: `if candidate.MatchStatus != domain.LinkCandidateMatchStatusAccept {`. No confidence, band, or status calculation reads `INCOMPARABLE`. |
| C5 | NOT-PROVEN | The implementation is statically sound: `generation_service.go:861`: `slices.SortStableFunc(pairs, func(a, b dimensionPair) int {`, and `:873`: `last.display += " ≡ " + pair.display`. Input order is deterministic regex order, not map iteration. However C5 explicitly requires independent `-count=10` and mutated-`SortFunc` failure artifacts. The referenced `scratchpad/c5a.txt` and `c5b.txt` were not independently available, and running tests was prohibited. Those two logs tied to `dbdcdfb1` would prove it. |
| C6 | PASS | Every relevant CTE is tenant-scoped at `query_repository.go:76`, `:83`, `:90`, and `:100`; both fan-out sets use DISTINCT at `:86`: `SELECT DISTINCT products.codprod AS codprod` and `:94`: the same for the queue. The fixture exercises two installations and duplicate linkage, with its exact assertion at `chain_query_repository_integration_test.go:80`: `got.Importados != 4 \|\| got.Vinculados != 2 \|\| got.Enfileirados != 3`. This criterion’s narrow fixture passes by inspection, although finding G1 shows an uncovered valid CODPROD representation. |
| C7 | PASS | `contracts/api/marketplace-central.openapi.yaml:8091`: `Current market queue at queue_read_at, never an import-history total:` and `:8092`: `the number falls as the queue is consumed.` `queue_read_at` is required at `:8080` and formatted as date-time at `:8094-8096`. |
| C8 | PASS | Missing rows map through `query_repository.go:121`: `if errors.Is(err, pgx.ErrNoRows) {` to `ErrImportNotFound`, then HTTP 404 at `http_handler.go:126-127`: `if errors.Is(err, ports.ErrImportNotFound) { writeError(...StatusNotFound...) }`. Empty subqueries return numeric `count(*) = 0`, and `http_handler.go:215` has no `omitempty`: `Enfileirados int json:"enfileirados"`. The raw response assertion is at `chain_query_repository_integration_test.go:175`: `!strings.Contains(empty.Body.String(), '"enfileirados":0')`. |
| C9 | PASS | The public endpoint behavior, OpenAPI, and SDK landed together in `37d6b7cc`, verified with `git show --stat`. The three surfaces are `http_handler.go:44`: `registerInteractiveRoute(mux, "/erp/imports/{id}/chain", ...)`; OpenAPI `:3257`: `/erp/imports/{id}/chain:`; SDK `packages/sdk-runtime/src/index.ts:1901`: `getErpImportChain: (id: string) => ...`. F-02 likewise placed Go, OpenAPI, and SDK in `9c030154`. |
| C10 | PASS | `git diff --name-only e98d8193 dbdcdfb1 -- apps/web apps/server_core/migrations` returned no paths. This matches the binding prohibition at `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/chip.md:212`: `FE \| nenhum. apps/web/ está fora de escopo...`. |
| C11 | PASS | The only emissions of `UNAVAILABLE` are the provider-not-supplied path at `generation_service.go:704` and the explicitly excluded both-nonempty branch at `:642`: `reason.Direction = domain.LinkCandidateReasonDirectionUnavailable`. The excluded branch is pinned at `generation_service_test.go:381`: `findReason(candidate.Reasons, "ean", ...DirectionUnavailable)`. A2-R2 is respected at `generation_service.go:719-720`: `if listingValue != "" && productValue != "" { return "", "", "", false }`; no unverified `FOR` is emitted. |

## PART B — FINDINGS

### G1. Zero-padded CODPRODs disappear from `vinculados` — BLOCKING

- **Locus**: `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:89`: `ON links.internal_product_id::text = products.codprod`
- **Why it is wrong**: ERP import validation preserves any nonempty CODPROD after trimming—`domain/validation.go:24`: `row.Codprod = strings.TrimSpace(row.Codprod)`—and does not reject leading zeroes. Linking then canonicalizes that value numerically at `adapters/internalread/reader.go:162`: `id, parseErr := strconv.ParseInt(strings.TrimSpace(row.CodigoProduto), 10, 0)`. The resolved link stores an integer (`migrations/0025_product_link_workflows.sql:9`: `internal_product_id integer,`). Consequently ERP product `"00101"` can produce resolved link ID `101`, but the chain comparison is `"101" = "00101"`, which is false. The endpoint silently underreports a product that is resolved.
- **Reachability**: Import a valid workbook row with `codprod = "00101"` and an EAN. Generate a concordant candidate and resolve it; the link stores `internal_product_id = 101`. Calling `/erp/imports/{id}/chain` returns `vinculados: 0` for that product. `IsValidCodprod` explicitly accepts leading zeros: `internal_read/domain/seller_sku.go:25-30` accepts every digit and merely requires one nonzero digit.
- **Yes-if**: This stops being a finding if import persistence enforces canonical, non-zero-padded decimal CODPRODs, or the chain join safely compares the same canonical identity representation. A regression fixture must use product `"00101"`, resolved link `101`, and expect it in `vinculados`.

### G2. Malformed import UUIDs return 500 — BLOCKING

- **Locus**: `apps/server_core/internal/modules/erp_import/transport/http_handler.go:124`: `chain, err := h.queries.GetImportChain(r.Context(), domain.ImportID(r.PathValue("id")))`
- **Why it is wrong**: The path value is cast to a named string without UUID validation. It is then bound to `eip.id = $2` at `query_repository.go:77`, while the column is UUID (`migrations/0045_create_erp_import_protocols.sql:2`: `id UUID PRIMARY KEY,`). PostgreSQL rejects a value such as `not-a-uuid`; this is not `pgx.ErrNoRows`, so the handler falls through to `http_handler.go:130`: `writeError(w, http.StatusInternalServerError, "internal_error", "")`. A client input error is reported as a server failure.
- **Reachability**: `GET /erp/imports/not-a-uuid/chain`. The SDK also accepts arbitrary strings at `packages/sdk-runtime/src/index.ts:1901`, so this does not require bypassing the supported client API.
- **Yes-if**: Validate `{id}` before repository access and return a documented 400, or deliberately map malformed identifiers to the existing 404 contract. Add a real-repository HTTP test for a malformed UUID.

### G3. A legal non-array `pending` value makes the endpoint fail — NON-BLOCKING

- **Locus**: `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:96-98`: `jsonb_array_elements_text(COALESCE(state.cursor -> 'pending', '[]'::jsonb))`
- **Why it is wrong**: `COALESCE` handles SQL NULL but does not verify JSON type. `{"pending":"101"}` or `{"pending":{}}` is non-null, so `jsonb_array_elements_text` raises a PostgreSQL error. The schema permits this because `migrations/0075_sync_sync_state.sql:23` declares only `cursor jsonb,`, and the generic success writer accepts opaque `json.RawMessage`.
- **Reachability**: A legal database row `(entity='market', cursor='{"pending":"101"}')` followed by the chain GET reaches the error. However, the current market enqueuer writes arrays and no market scheduler job is registered, so a current external producer of malformed market state was not shown reachable; severity is therefore NON-BLOCKING.
- **Yes-if**: Enforce the market cursor shape at every writer/database boundary, or type-check with `jsonb_typeof` and define explicit behavior for corrupt cursor state.

### G4. The resolved-link count has no supporting state/identity index — NON-BLOCKING

- **Locus**: `apps/server_core/migrations/0025_product_link_workflows.sql:17-18`: `CREATE INDEX IF NOT EXISTS product_links_installation_idx ON product_links (tenant_id, installation_id, updated_at DESC);`
- **Why it is wrong**: The new query filters `tenant_id`, `state = 'resolved'`, and joins through `internal_product_id::text`. The only product-links index begins with tenant but then installation and timestamp; it cannot directly satisfy the state/identity lookup, leaving all link rows for a tenant subject to filtering. The cast also prevents an ordinary integer identity lookup from matching the text side directly.
- **Reachability**: Not shown reachable at a tenant size sufficient to exceed the interactive 15-second deadline; therefore NON-BLOCKING.
- **Yes-if**: A production-scale `EXPLAIN (ANALYZE, BUFFERS)` establishes the current plan meets the route budget, or an index/canonical join supports `(tenant_id, state, internal_product_id)`.

## WHAT I READ

Binding artifacts, in the required order:

- `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/validation-contract.md`
- `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/hub-rulings.md`
- `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/chip.md`
- `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-2/EVIDENCE.md`

All 20 changed code/contract paths:

- `apps/server_core/internal/composition/root_test.go`
- `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`
- `apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go`
- `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go`
- `apps/server_core/internal/modules/erp_import/application/query_service.go`
- `apps/server_core/internal/modules/erp_import/application/query_service_test.go`
- `apps/server_core/internal/modules/erp_import/domain/import.go`
- `apps/server_core/internal/modules/erp_import/ports/repository.go`
- `apps/server_core/internal/modules/erp_import/transport/http_handler.go`
- `apps/server_core/internal/modules/erp_import/transport/http_handler_test.go`
- `apps/server_core/internal/modules/product_links/adapters/connectors/identity_anchor_adapter_test.go`
- `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go`
- `apps/server_core/internal/modules/product_links/application/generation_service.go`
- `apps/server_core/internal/modules/product_links/application/generation_service_test.go`
- `apps/server_core/internal/modules/product_links/domain/link_candidate.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/erpImport.test.ts`
- `packages/sdk-runtime/src/erpImport.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `packages/sdk-runtime/src/index.ts`

Supporting blast-radius paths:

- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/platform/httpx/router.go`
- `apps/server_core/internal/platform/httpx/route_deadline.go`
- `apps/server_core/internal/modules/product_links/application/resolution_service.go`
- `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go`
- `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go`
- `apps/server_core/internal/modules/erp_import/adapters/sync/enqueuer.go`
- `apps/server_core/internal/modules/erp_import/domain/validation.go`
- `apps/server_core/internal/modules/internal_read/domain/seller_sku.go`
- `apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go`
- `apps/server_core/internal/modules/sync/application/scheduler.go`
- `apps/server_core/internal/modules/sync/composition/products_job.go`
- `apps/server_core/internal/modules/sync/domain/sync_state.go`
- `apps/server_core/migrations/0025_product_link_workflows.sql`
- `apps/server_core/migrations/0045_create_erp_import_protocols.sql`
- `apps/server_core/migrations/0046_create_erp_import_products.sql`
- `apps/server_core/migrations/0075_sync_sync_state.sql`
- `apps/web/src/pages/vinculos/QueueRow.tsx`
- `apps/web/src/pages/vinculos/VinculoDrawer.tsx`
- `apps/web/src/pages/vinculos/QueueTab.test.tsx`

I regenerated the reviewed 20-file diff with `git diff e98d8193 dbdcdfb1 -- apps contracts packages`; its `+1220/-124` count matched the frozen input description. I did not open the external frozen patch itself. I could not independently read the prior-session `scratchpad/c5a.txt`/`c5b.txt` artifacts, which is why C5 is NOT-PROVEN.