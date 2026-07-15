# M-01 Milestone-Session Decisions

Ratified by the M-01 chip (milestone session) to resolve the non-blocking open risks from F-01/plan.md.
Authority basis: IC-02 "provider-mapped at adapter" + ADR-17 (unknown→null) + IC-02 extension rule ("add keys/cases, never change existing semantics"). Contract-fixed items (identity, existing enums, existing error codes, filter grammar, sort, group envelope) are NOT touched here. The two hub-gated blockers (migration number, connectors capability seam + provider-auth error code) are excluded — awaiting hub.

## D-5 · Modality (listing_type) allowlist + label registry
IC-02 pins `listing_type` nullable, provider-mapped, example `{code:"gold_pro",label:"Premium"}`, but does not enumerate codes. Ratified recognized ML `listing_type_id` → pt-BR label registry:

| code | label |
|---|---|
| gold_pro | Premium |
| gold_special | Clássico |
| gold_premium | Ouro Premium (legado) |
| gold | Ouro (legado) |
| silver | Prata (legado) |
| bronze | Bronze (legado) |
| free | Grátis |

Rules: recognized code → persist `listing_type_code`, label applied at read (F-02) from THIS registry. Empty/absent OR unrecognized provider modality → `listing_type_code` SQL NULL and the read-time `listing_type` object is `null` (ADR-17 unmappable→null; never guess a label). Registry lives once in the listings module, consumed by F-01 (recognition) + F-02 ({code,label} assembly). Extensible by adding rows, never remapping.

## D-6 · Active run id placement in 409 refresh_in_progress
IC-02 requires the 409 body to carry the active `operation_run_id` but not its location. Ratified: `error.details.operation_run_id` (string). Uses the existing shared `ErrorResponse.error.details` free-object (openapi.yaml:4292-4321) — least-disruptive, no schema-shape change. Body: `{"error":{"code":"refresh_in_progress","message":"...","details":{"operation_run_id":"<id>"}}}`.

## D-7 · Known-but-not-connected installation
IC-02 error matrix covers missing (400 installation_required) / unknown (404 installation_not_found) / concurrent (409). It is silent on an installation that EXISTS but is not in a runnable/connected state (draft/disconnected/suspended/requires-reauth). Ratified NEW error case (extension, not a change): existing installation not in a connected/runnable state → **409 `installation_not_connected`**, message directs operator to reconnect via Integrações. RK-06 honored: NO StartAuthorize/reauth triggered. M-01 functional + live lanes exercise only connected installations; this case is proven by a table-test row. Do NOT collapse into installation_not_found.

## D-8 · price_currency nullability
IC-02 marks the price object nullable; DB shape marks `price_amount NUMERIC NULL` but does not annotate `price_currency`. Ratified: `price_currency` is nullable and is NULL whenever the provider price fact is absent (no fabricated "BRL"). Amount+currency move together: both present or both NULL. Matches ADR-17.

## D-9 · Listings→installation FK / delete policy
IC-02 does not specify an FK or ON DELETE. Ratified: **no cross-module FK** from `listings` to the integrations installations table. Rationale: keeps module boundary clean (listings does not depend on integrations' physical schema; it depends on the published interface), and preserves listing history if an installation row is removed. Installation existence is enforced at the application layer (refresh resolves the installation via the integrations published port before starting a run). Tenant scoping (`tenant_id` predicate) is unconditional on every query regardless.

## D-10 · Pagination completion semantics
The only `ListingReader` (ML) uses an offset cursor and returns a slice with no total/next-cursor (capability_adapter.go:117-141). Ratified full-pull termination: start cursor `""`(→offset 0), advance offset by the count of provider items returned, terminate on a short page (returned < requested limit). Absent-row closing runs ONLY after a fully-completed pull. If a future reader cannot honor offset/short-page semantics, the published interface must gain explicit completion metadata before that reader is onboarded (documented constraint, not F-01 work).

## D-11 · Duplicate canonical key across provider pages
IC-02 fixes the key but is silent on duplicates. Ratified **fail-honest**: if two snapshots in one pull map to the same `(tenant_id, installation_id, provider_listing_id, variation_id)`, the refresh FAILS before any persistence (operation run → failed, honest code) rather than last-write-wins. Rationale: last-write-wins hides provider inconsistency and makes closed-marking unsafe. Existing rows unchanged.

## D-12 · Stale queued/running run recovery — DEFERRED
Crash-recovery for a run left queued/running after process termination is out of M-01 scope. M-01 accepts manual cleanup. F-01 must NOT invent a timeout that could permit overlapping refreshes. A stale-run recovery rule is a later-feature/mission concern. (Recorded as a milestone defer.)

## D-13 · Seed link states vs F-01 scope
IC-02 seed (MLBTEST0001..0006) includes resolved/unlinked examples, but F-01 must not touch product_links. Ratified boundary: F-01 seeds the six fixed listing IDs + listings-owned facts ONLY; the link fixtures (product_links rows) and joined resolved/unlinked assertions belong to F-02. Not authorization to duplicate link columns in `listings`.

## D-14 · GOCACHE absolute path
Go requires an absolute GOCACHE. Ratified: resolve the repo `.gocache` dir to an absolute path and assign before Go commands (Windows/pwsh): e.g. `$env:GOCACHE = (Resolve-Path .gocache).Path`. Preserves the intended cache location while satisfying Go.

## D-15 · No SDK generator — contract-first hand-maintained parity
`packages/sdk-runtime` has no OpenAPI codegen command. Ratified process for the refresh endpoint slice: edit OpenAPI first, then hand-maintain the sdk-runtime types + method + tests in the SAME commit. Validation language: "contract-first SDK parity" — never claim regeneration that did not occur. Satisfies GOV_API_SDK_SPLIT (both files in one commit).

## D-1 · Migration number — RESOLVED by hub
Hub reallocated (0033–0035 confirmed committed at base). Granted block **0036–0037**. Listings migration = **`0036_listings.sql`**; 0037 reserved headroom. Nothing beyond 0037 without a new REQUEST.

## D-3 · Connectors capability seam — RESOLVED by hub (temporary contract-lock, option a)
M-01 holds a temporary contract-lock on the connectors capability seam until CLOSED. Scope, strictly:
- **Additive only**: new OPTIONAL fields on the ListListings snapshot — `price {amount (string→NUMERIC-safe decimal), currency}`, `listing_type code` — + ML adapter mapping. NO renames, NO removals, NO behavior change for existing consumers; existing connectors tests stay green untouched.
- **Provider-auth error code** (matches domain convention capability.go:39-44): `ErrCodeProviderAuth ErrorCode = "CONNECTORS_PROVIDER_AUTH"`. ML adapter maps HTTP 401/403 to it, splitting them out of provider_validation. Listings surfaces it as provider_auth-class on the operation run (RK-06). No StartAuthorize/reauth.
- Unmappable price/currency/listing_type from provider → NULL/absent, never zeroed (ADR-17).
- **Guardrail**: if the additive change requires touching non-ML adapters beyond compiling stubs / zero-value defaults → STOP and re-REQUEST to hub.
- Connectors diff travels IN the milestone diff (per-slice review + dual gate cover it); CLOSED payload must call out the connectors files + why.

## D-16 · below_margin formula — OPERATOR-RATIFIED (contract-literal)
F-02 planner surfaced a formula ambiguity: IC-02 says "latest cost + min-margin"; the F-02 brief says "cost + price"; `pricing/application/service.go:34-51` additionally deducts commission + fixed fee + shipping. Operator ratified **contract-literal** (2026-07-15):
- `below_margin = price_amount < cost_amount × (1 + min_margin_percent/100)`.
- Inputs: latest cost fact (internal_read cost port) + `min_margin_percent` (marketplaces policy for the installation). NO commission/fixed-fee/shipping deduction — that economic margin is the pricing module's concern, out of the read-spine scope.
- **Null-safe (ADR-17):** any of {price, cost, policy} absent → `below_margin = null` (never false).
- Summary `exceptions.below_margin` counter counts only rows where below_margin is non-null AND true; excludes null-cost rows (C07).

## D-17 · Exception precedence + group severity (R3) — milestone-ratified
IC-02 fixes the exception filter values but not precedence/severity. Ratified ONE severity table, reused by `has_exception`, item `exception` mapping, `group_state`, and summary counters:
- Per-listing active-issue precedence (worst first): **sync_error > stale > below_margin > unlinked**. The item's `exception` field is the highest-precedence active issue; `has_exception` = any active.
- `attribute_required` is a sync_error sub-reason (surfaced in message_pt), not a separate top-level exception.
- `group_state` = worst child: **error** if any child sync_error; **attention** if any child in {stale, below_margin, unlinked}; else **ok**. Order-independent.

## D-18 · Keyset stability semantics (R5) — milestone-ratified
C09 "stable cursor" = **deterministic keyset ordering** (title ASC, listing_id ASC), opaque base64, `limit+1` look-ahead, no offset. Snapshot isolation across a concurrent title-changing refresh is **NOT** required for M-01: refresh is atomic (D-10/D-11), and a title mutation between pages may reposition a row — accepted and documented. No as-of snapshot / row-version / session storage invented. Consistent with D-10.

## D-19 · Timeline event source (R4) — milestone-ratified, FOLDED INTO F-01
F-02 needs `GET /listings/{id}` timeline (last 10 sync events); F-01's current row + installation-wide operation runs cannot supply it. Ratified (listings-owned, within lane):
- New table (NOT new listing columns) **`listing_sync_events`** in migration `0036_listings.sql`, keyed `(tenant_id, installation_id, provider_listing_id, variation_id, event_id)`, columns `at`, `kind`, `message_pt`.
- Event-kind vocabulary: `synced`, `sync_error`, `closed`, `paused`, `refreshed`.
- F-01 writes events in the **same completed-pull transaction** as the affected listing state; F-02 only reads last 10 (`at DESC, event_id DESC`). Empty-fallback / installation-wide history forbidden.
- Index `(tenant_id, installation_id, provider_listing_id, variation_id, at DESC, event_id DESC)`.
- **Scope impact:** this EXPANDS F-01. Fold into F-01 repository/migration design before F-02 Slice 3.

## D-20 · Cost source vs C09 single-query (R1) — RESOLVED by hub (option a)
Latest cost is authoritative in **Oracle**; listings/link/policy are **PostgreSQL**; no local latest-cost projection. Hub GRANTED option (a): **C09 "one query" is reinterpreted** as ONE PostgreSQL conditional-aggregate over listings/link/policy **+ ONE bounded Oracle batch cost read** (single `GetCostFactsByIDs` for the whole page/id-set — never per-row). Intent of C09 = no-N+1 + honest counters, not a literal single statement. Option (b) local cost projection = mid-milestone architecture invention → refused.
Constraints (binding):
- **Unknown-cost rows are NEVER silently classified as not-below-margin.** They must be surfaced observably. Ratified shape (additive per IC-02 extension rule): summary `exceptions.below_margin` counts only rows where below_margin is **known true**; add observable counter `exceptions.margin_unknown` = rows where below_margin is null due to missing cost/price/policy. Item-level below_margin stays null for those (D-16). Unknown is thus derivable and never hidden.
- Document this ratified interpretation in `F-02/validation.md` AND flag it in the CLOSED payload as a **criterion reinterpretation** so QA validates M01-C09 against the amended reading.
- Unblocks Slice 5. Below_margin list filter/counter use the single batch Oracle cost read.

## D-21 · Marketplaces pricing-policy seam — RESOLVED by hub (additive contract-lock)
Hub GRANTED M-01 a temporary additive contract-lock on the marketplaces module for exactly one new published method: **`GetPricingPolicyForInstallation(ctx, installationID) (Policy, bool, error)`**, tenant-scoped via `marketplace_accounts.integration_installation_id`. Terms (same as D-3 connectors lock):
- **Additive only**: no renames/removals/behavior change; existing marketplaces tests stay green untouched.
- Unknown installation → typed not-found (bool false / typed error), **never a nil/zero-policy default** (ADR-17).
- Multiple matching policies → fail honestly unless a selection rule is later ratified (no silent first-row).
- Lock ends at M-01 CLOSED; the marketplaces diff is called out explicitly in the CLOSED payload.
- Touching anything in marketplaces beyond this one method + its tests = new REQUEST to hub.

## Escalation (out-of-scope, not blocking)
- ADR-12 / ADR-17 have no formal record under docs/architecture/decisions (behavior unambiguous in mission.md). Architecture owner repairs; F-01 proceeds on the mission-fixed behavior.
