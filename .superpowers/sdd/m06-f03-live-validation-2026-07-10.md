# M-06 F-03 live integration validation — 2026-07-10

## Target

- Docker Compose environment: PostgreSQL, backend, frontend all healthy.
- Mercado Livre installation: `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`.
- Oracle/Sankhya target: configured live reader; credentials intentionally not recorded here.

## Real-environment evidence

| Surface | Evidence | Result |
| --- | --- | --- |
| Docker + PostgreSQL | `docker compose ps`; backend `/healthz`; database-backed adapter tests | Healthy; real persistence tests passed. |
| Mercado Livre | Account probe, listings probe (20), orders probe (20); import with limit 50 | Real provider returned account, listings, and orders; 30 orders imported. |
| Oracle/Sankhya | `MPC_ORACLE_LIVE_TEST=1 go test ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -count=1 -v` | Passed product lookup, sellable stock, current price, as-of CUSSEMICM cost, sales history, and tax inputs. |
| Margin pipeline | Live ML order import -> margin-input import -> snapshot calculation -> API readback | 270 inputs and 60 snapshots calculated before response limiting. |
| Unknown-data safety | Margin-input and snapshot API readback plus browser state | Missing freight/commission/tax/link data remained `null` / `—`; no unknown-to-zero or unknown-to-realized conversion. |
| Browser | Built-in browser at `http://127.0.0.1:5174/orders`, Mercado Livre selected | Desktop and mobile rendered; no browser console errors; mobile document had no page-level horizontal overflow. |

## Live margin input readback

- Revenue: 30 valued.
- Sale fee: 30 valued.
- Oracle CUSSEMICM cost: 29 valued, 1 null.
- Oracle PIS: 5 valued, 25 null.
- Oracle COFINS: 5 valued, 25 null.
- Freight and commission: 30 null each (explicitly not yet provided by live source).

## Deterministic and real-PostgreSQL verification

- `TestProfitSnapshotRealizationPersistence`: PASS against Docker PostgreSQL.
- `TestManualAdjustmentsAppendOnlyReadbackAndConstraints`: PASS against Docker PostgreSQL.
- `go test ./internal/modules/profitability/... ./internal/composition -count=1`: PASS.
- SDK: 35/35 PASS.
- feature-orders: 13/13 PASS (including the idempotency retry lifecycle regression).
- web: 184/184 PASS.
- `npm run build --workspace @marketplace-central/web`: PASS.

## Observations and remaining gate conditions

- Oracle driver emitted a timezone warning (`SESSIONTIMEZONE +00:00` vs `SYSTIMESTAMP -03:00`), but all live queries passed. This is an operational follow-up, not a test failure.
- Candidate generation found exact-EAN candidates, but no new product link was approved during this validation because approval requires a truthful, explicit audit actor.
- This evidence validates live connectivity and the current safe data-quality behavior. It does **not** mark M-06 passed: the independent cold milestone gate and any required approved resolved-link scenario remain separate acceptance conditions.

## Visual evidence

- `evidence/m06-f03-total-validation-desktop.png`
- `evidence/m06-f03-total-validation-mobile.png`

## Superseding correction — unapproved candidate leakage

The earlier live readback reporting 29 valued Oracle costs is **withdrawn as
resolved-link evidence**. Cold-gate preflight proved those values were reached
through exact-match candidates whose order quality was still `unresolved`.
Candidate generation is not operator approval.

The source boundary and profitability defense were corrected under TDD and
received independent SPEC and QUALITY approval. After restarting the Docker
backend onto the corrected binary, the orchestrator ran a fresh live Mercado
Livre import and recalculation:

- Mercado Livre order probe: PASS.
- First imports occasionally returned the provider's transient error; a
  bounded idempotent retry succeeded on attempt 2.
- Orders imported: 30.
- Margin inputs imported: 270.
- Snapshots calculated: 60.
- Non-resolved order items with an internal product ID: 0.
- Known cost/tax inputs attached to non-resolved items: 0.
- Target `MLB4834373620`: `link_quality=unresolved`,
  `internal_product_id=NULL`.
- Its cost, ICMS, IPI, PIS, and COFINS inputs are all `NULL`, quality
  `unresolved_link`, reason `internal product link is not resolved`.
- Live installation snapshots: 48 incomplete and 12 not-realized; 0 complete,
  0 with contribution/margin.

This corrected evidence proves safe unresolved-link behavior. It does not
prove the separate required resolved-link scenario and does not mark M-06
passed.

## Fresh Docker baseline and live revalidation — 2026-07-10

The Docker volume was recreated after the earlier evidence. The stack was
rebuilt from the current worktree and revalidated independently; this section
is the authoritative evidence for that fresh database state.

- `postgres`, `backend`, `frontend`, and the optional OAuth tunnel were up;
  PostgreSQL and backend health checks passed.
- The already-persisted Mercado Livre installation was `connected` and
  `healthy` with an active credential. No credential material is recorded.
- The operator drove the real `/orders` UI: order import returned 20 Mercado
  Livre orders; margin-input import then completed against the configured
  live internal reader; recalculation persisted 40 snapshots (20 item and 20
  order scopes).
- Snapshot grouping in real PostgreSQL was: 15 item and 15 order snapshots
  `incomplete` / `realized`, plus 5 item and 5 order snapshots
  `not_realized` / `not_realized`. All 40 had `NULL` cost, tax, and
  contribution where their product link remained unresolved. The `realized`
  realization state is the explicit paid-order state, while `incomplete`
  quality kept contribution and margin unavailable; it is not realized-margin
  math.
- Browser QA at desktop and 390x844 mobile showed the imported records,
  `Missing` quality flags, `Order cancelled` for cancelled orders, and em
  dashes for unavailable contribution/margin. Console error/warning capture
  was empty. The mobile document width did not overflow its viewport.
- Fresh targeted checks passed in Docker: profitability/composition Go suite,
  SDK (35 tests), Orders router/client/proxy tests (8 tests), and web build.

The broad server command `go test ./apps/server_core/... -count=1` did not
pass. Its F-03 profitability packages passed, but unrelated integration and
legacy tests failed because the live OAuth redirect environment changes their
implicit default, marketplace fixtures omit the now-required marketplace
code, one migration test resolves the wrong relative directory, and an
inventory expectation diverges. Those failures are recorded as a cold-gate
blocker; they were not changed in the shared dirty worktree.

## Candidate-generation follow-up — 2026-07-10

After the fresh revalidation, a full-server test command created three
synthetic integration installations in the shared Docker database. They were
identified by their test-only `inst-cred-*` identifiers and creation times,
then deleted in a single transaction. The authenticated Mercado Livre
installation and real orders, inputs, snapshots, listing identities, and
candidates were preserved.

The operator then used the real installation to import 20 live listing
snapshots and generate 20 exact-EAN candidates through the live Oracle
reader. The initial candidate request returned `503
PRODUCT_LINKS_INTERNAL_READ_UNAVAILABLE` because the backend had started
while Oracle ping was temporarily unavailable. The official live Oracle smoke
then passed in Docker; restarting the backend reconstructed the reader and
the same candidate generation succeeded. No candidate was approved.

Two generated candidates overlap imported order items:

| Provider order | Order status | Listing | Candidate internal product | Match |
| --- | --- | --- | --- | --- |
| `2000012659424976` | `cancelled` | `MLB4834373620` | `20303` / `PLACA ADES.PUXE/EMPURRE` | exact EAN `7898016503065` |
| `2000012747964244` | `cancelled` | `MLB4834408384` | `20322` / `PLACA ADES.CALANDRADA PROIBIDO FUMAR` | exact EAN `7898016503621` |

Reimporting up to 50 real orders produced 30 persisted orders, but no paid
order overlapped an available generated candidate. Approving either listed
candidate would correctly prove a resolved *not-realized* cancellation path,
not the required resolved paid-order margin scenario. A truthful actor
approval plus either an approved paid-order resolution or a later qualifying
provider order remains necessary. M-06 remains blocked.
