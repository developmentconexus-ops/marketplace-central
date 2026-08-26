# D6-R2 P9 — B24 Preços Screen Contract

> **Status:** DERIVED / PASS — P8 LOCKED 2026-08-26; BACKEND SUFFICIENT after the bounded MKT-01 projection repair; UPSTREAM FINDING NONE
> **Block:** B24 — Preços / R24
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked P8 evidence:** `qualification/d6-r2-wireframes/b24-price-intents.html` (revision 5, LOCKED 2026-08-26)
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml` (108/31/H-A-S after MKT-01)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 result

P9 ran against the operator-LOCKED revision-5 candidate. The locked human job is:

```text
open the pricing workbench for the organization
→ read, per row, the current price, the Economics margin fact and the Market position fact
→ narrow to what needs attention (below policy / not at the lowest delivered price) — server facts only
→ type a candidate price; the owners re-answer (margin, and position in the comparable range)
→ open the cascade and the reference anchors when the trade-off needs proving
→ create exactly one explicit PriceIntent for that row, superseding any active one
→ read the honest external outcome; verify (never blindly retry) on ambiguity
```

**P9 verdict: PASS / BACKEND SUFFICIENT / UPSTREAM FINDING NONE** — after the bounded D5 projection repair recorded in §6.

## 2. Route, identity and client-state ownership

Route family: `/org/:organizationId/publicacoes/precos`, two fixed views (`Decidir` workbench / `Intenções` ledger) — not a saved-view platform.

| State class | B24 ownership |
| --- | --- |
| `GLOBAL_WORKSPACE_CONTEXT` | `organization_id` |
| `URL_NAVIGATION_STATE` | active view, account filter, attention filter |
| `SERVER_STATE` | listings, price intents, expected economics, competitive positions, scenario evaluations, effect state |
| `LOCAL_EPHEMERAL` | the unsent typed price, row disclosure state, mobile navigation |

The typed price is local and inert until it is either sent to an owner for evaluation or committed as an explicit intent. No pricing number displayed on the screen is produced by the screen (`data-client-computation="none"`).

## 3. Product operation / access binding

| Screen need | Operation | Semantic owner | Permission | Principal kinds |
| --- | --- | --- | --- | --- |
| workbench rows (listings) | `ListMarketplaceListings` | Offering | `offering.read` | H/A/S |
| active intent per row + ledger | `ListPriceIntents` | Offering | `offering.read` | H/A/S |
| read one intent / verify outcome | `GetPriceIntent` | Offering | `offering.read` | H/A/S |
| create intent (supersede-only) | `CreatePriceIntent` + Idempotency-Key | Offering | `price.manage` | H/A |
| current margin per row | `ListExpectedEconomics` | CommercialEconomics | `economics.read` | H/A/S |
| margin at the typed price + waterfall + anchor evaluation | `EvaluatePriceScenario` | CommercialEconomics | `economics.read` | H/A/S |
| current market position per row | `ListCompetitivePositions` | MarketIntelligence | `market.read` | H/A/S |
| comparable evidence behind the range | `ListComparableOffers` | MarketIntelligence | `market.read` | H/A/S |
| position of the **typed candidate** price | `EvaluateCompetitivePositionScenario` (MKT-01) | MarketIntelligence | `market.read` | H/A/S |

The workbench composes these cursor collections **page-level**: one page of rows, one page of each owner collection, no N+1 and no screen-shaped aggregate endpoint.

## 4. What traces cleanly

- **Write path.** Exactly one consequential write (`CreatePriceIntent`), explicit per row, Idempotency-Key carried, supersede-instead-of-edit, no bulk apply, no repricing engine.
- **Honest outcome.** `pending / applied / rejected / ambiguous / superseded` plus `convergence`; rejection surfaces provider feedback verbatim with a substitute path; ambiguity is verification-only, never blindly retried.
- **Owner separation.** Economics answers margin and policy judgment; Market answers position; Offering owns the write. The screen sets no threshold and computes no money.
- **Honest population.** known-empty ≠ unknown ≠ unavailable, for the workbench, the ledger, economics conclusions and market evidence alike.
- **Attention filter.** Projects server-issued facts (`below_policy`, position not first), never a client score.
- **Anchors.** Owner-evaluated; the market anchor is gated by `evidence_sufficiency` and its absence is explained; clicking only fills the input.

## 5. Backward trace

No admitted operation above is orphaned, and no locked control lacks an operation — with the single exception recorded below.

## 6. MKT-01 — bounded D5 projection repair (not a domain reopen)

The first P9 pass raised the locked Mercado column (delivered range, our rank, rank at the typed price) as an upstream finding. **The operator revoked that framing and was right.** Re-reading the smallest owner set proves the meaning was already accepted upstream; only its API projection was missing.

| Already-accepted authority | What it says |
| --- | --- |
| [D1 Domains/Boundaries](D1-DOMAINS-BOUNDARIES.md) §domain table | Market Intelligence owns *comparability, **competitive position/change**, market-evidence sufficiency*; pricing authority stays with Economics and the write with Offering |
| [D4 External Integrations](D4-EXTERNAL-INTEGRATIONS.md) §6.2 (D4-B4, ACCEPTED / CANONICAL) | The Mercado Livre lane proved catalog offer population, own/winner offer price, buyer-facing shipping and `price_to_win` as materially useful Market evidence; *"Market Intelligence may derive delivered-price/competitive explanation"* |

So the delivered range and our position within the comparable population are **already inside an accepted owner's meaning**, proven against a real provider. No D0/D1/D4 reopen was warranted. What was missing was the **D5 projection**: the OAD exposed only `relation` + `delivered_price_gap`, which is strictly less than the accepted meaning and is exactly the ahead/behind wording the locked screen removed.

Bounded repair executed (MKT-01), inside existing owners and Permissions:

- `MarketDeliveredPriceRange` (low/high) and `MarketRank` (`position`, `comparable_count`, closed `basis: observed_comparable_population`);
- both added as **optional** projections on `CompetitivePosition`, `CompetitivePositionListItem` and the new scenario evaluation, while `evidence_sufficiency` becomes **required** on the collection item — insufficient or unavailable evidence therefore has no rank to state and can never read as a confident position;
- `EvaluateCompetitivePositionScenario` — stateless class C, `market.read`, H/A/S, owner MarketIntelligence — answering the position of a **candidate** price, mirroring what `EvaluatePriceScenario` does for Economics. Product surface **107 → 108**; Permissions unchanged at 31; Principal kinds unchanged.

D4 §6.2's knowledge-state controls are carried into the projection rather than assumed away: *"catalog offer paging is bounded provider population, never general-market completeness."* The rank is therefore explicitly a rank **within the observed comparable population**, both in the closed `basis` vocabulary and in the locked screen's own copy (*"população observada, não o mercado inteiro"*), and `coverage` travels with every answer. This was applied as a bounded copy/binding adjustment to the LOCKED B24 evidence; the locked structure — regions, placement, density, navigation, state placement — is unchanged.

Rejected alternatives remain rejected: deriving range/rank in the client (breaks the block LOCK and discards coverage/sufficiency); folding the market answer into `EvaluatePriceScenario` (collapses the D1.8 Economics/Market separation); dropping rank from the locked screen (backend-shaped UX, §3.10A).

## 7. Adversarial checks

P9 rejects: client-computed range, rank or margin; a rank stated without sufficient evidence; a rank implying whole-market completeness; label-carrying writes; bulk apply or an automatic repricing engine; in-place price edit; blind retry after an ambiguous external effect; hiding provider rejection feedback; collapsing known-empty into unknown or unavailable. All are excluded by the locked evidence plus `verify-d6-r-b24-price-intents-wireframe.mjs` (**30/30** negative controls) and the projection proof (**23/23**).

## 8. P9 closure and P10 note

```text
P8 OPERATOR-RATIFIED / LOCKED (2026-08-26, revision 5)
→ exact route/state/identity binding
→ exact owner/operation/Permission binding
→ frontend → backend trace: one gap found
→ gap re-scoped against accepted D1/D4-B4 authority: D5 projection only
→ bounded MKT-01 repair + proof (108/31/H-A-S, historical non-regression PASS)
→ backend → frontend trace PASS
→ adversarial shortcuts rejected
→ BACKEND SUFFICIENT
→ UPSTREAM FINDING NONE
```

**P9: PASS / CLOSED for B24.**

P10: B24 reuses the established laws (owner-issued facts vs screen rendering; known-empty ≠ unknown ≠ unavailable; navigation ≠ mutation; server-gated consequential writes; ambiguous-verify). The pricing workbench row — owner facts, a debounced scenario evaluation and one explicit supersede-only write — is B24-local until a second locked block proves the same shape; **no new shared component/pattern authority is claimed.** P11, Pre-D9/D9 and Product implementation remain outside this closure.
