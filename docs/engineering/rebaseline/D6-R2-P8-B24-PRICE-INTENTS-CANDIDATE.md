# D6-R2 P8 — B24 Preços Ratification

> **Status:** P8 OPERATOR-LOCKED 2026-08-26 (revision 5)
> **Block:** B24 — Preços / R24
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Candidate evidence:** `qualification/d6-r2-wireframes/b24-price-intents.html`
> **Prerequisites:** PR #76 (B23) integrated; 107/31/H-A-S surface
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Design adjudication

The bounded B24 P8 structural design was adjudicated in chat on 2026-08-26 and **approved by the operator**:

- **R24 `/publicacoes/precos`** — PriceIntent collection + explicit create; detail outcomes inline;
- **supersede-instead-of-edit law**: changing a price creates a new intent that supersedes the previous one (kept in history as `superseded`); no in-place price edit anywhere;
- typed targets with typed presentation: existing listing or **pre-creation** target (price decided before the listing exists — feeds the B23 flow);
- honest state truth (`pending / applied / rejected / ambiguous / superseded`) + `convergence`; rejected surfaces provider feedback verbatim with a substitute-intent path; **ambiguous is verification-only** (no blind retry);
- decision context as inline **read-only** owner-separated facts (Market: comparable range + position at the typed price; Economics: predicted margin at the typed price), with navigation continuations — no price is computed or applied automatically;
- Offering owns price writes (`price.manage`, H/A, Idempotency-Key on create);
- **rejected by default**: automatic repricing engine, follow-competitor rules, bulk repricing, saved views, screen-shaped API.

## 2. Authority consumed

| Screen need | Operation | Permission | Principal kinds |
| --- | --- | --- | --- |
| R24 collection | `ListPriceIntents` | `offering.read` | H/A/S |
| create intent | `CreatePriceIntent` (+Idempotency-Key) | `price.manage` | H/A |
| read one intent | `GetPriceIntent` | `offering.read` | H/A/S |
| Market context facts | `GetCompetitivePosition` | market read | H/A/S |
| Economics context facts | `GetExpectedEconomics` | economics read | H/A/S |

`PriceIntent` carries typed `target` + `target_presentation` (existing/pre-creation pair-enforced), `desired_price` (Money), `state`, `supersedes_price_intent_id` and `convergence` — each bound to exactly one rendering home.

## 3. Proof

- deterministic verifier `scripts/verify-d6-r-b24-price-intents-wireframe.mjs`: **10/10 negative controls** (edit-in-place, repricing engine, premature create, ambiguous retry, context-facts writes, population collapse, idempotency drop, hidden feedback, listing-mutation smuggling), wired diff-aware into `npm run gate`;
- browser-operated flow (Chromium): create gate (target+price explicit) → supersede notice on an already-priced target → live market-position/margin facts at the typed price → create → pending → verify resolves to applied and the list updates → ambiguous shows the verification-only law → rejected shows verbatim feedback + substitute path → population states distinct;
- 390px viewport: no horizontal document overflow; mobile drawer law active; browser console: **0** warnings/errors.

## 4. Revision 2 — pricing workbench (operator-aligned, 2026-08-26)

The operator's walkthrough of revision 1 found the intent-first shape non-dynamic for the real pricing job ("one target at a time"). The alignment conversation approved a listing-centric restructure, verified as PRESENT-IN-AUTHORITY (the workbench composes the existing cursor collections — listings, price intents, expected economics, competitive positions — page-level, with no N+1 and no new endpoint):

- **two fixed views** (not a saved-view platform): **Decidir** (default) — the pricing workbench: one row per listing with current price, Economics margin fact, Market range/position fact, active intent, and an inline "new price" input with an explicit per-row confirm (`data-row-write="one-intent-per-row"`); attention filter projects server facts (low margin / out of range); **Intenções** — the intent ledger (states, supersessions, outcomes) plus the pre-creation special case;
- the per-row confirm creates a normal PriceIntent (Idempotency-Key, supersede law, honest outcome flow) — bulk apply remains rejected;
- verifier extended to **13/13 negative controls** (bulk apply, premature row confirm, client-scored attention filter, listing-mutation smuggling); browser proof re-passed end to end (default view, attention filter, inline confirm → pending → verify → applied reflected in row and ledger, workbench population honesty) with **0 console issues** and no 390px overflow.

## 5. Revision 3 — owner-evaluated pricing indicators + single create home (2026-08-26)

The operator's second alignment produced two changes, both verified as PRESENT-IN-AUTHORITY:

**(a) Single create home.** Having a create form in the ledger view was ambiguous. The workbench now also carries **pre-creation target rows** (products with a draft listing intent, no listing yet), so every price decision happens in **Decidir**; **Intenções** is a pure ledger (`data-intents-view="ledger-only"`, `data-create-home="workbench-only"`) that creates and alters nothing.

**(b) Pricing indicators from the owner.** The typed price is no longer a bare box. `EvaluatePriceScenario` (an already-admitted stateless C operation, `economics.read`) is triggered **explicitly** (never per keystroke) and the row expands with:

- the **price waterfall** from owner components: price − marketplace fee − seller-borne shipping − tax − promotion − product cost = **contribution**;
- **contribution (R$) and contribution margin (%)** side by side;
- the **owner's policy judgment** (`acceptable` / `below_policy`) — the screen never sets the threshold;
- **position at the delivered price** (product + shipping) versus comparables, with the gap;
- **reference anchors** already evaluated — current price, active intent, **policy floor**, and *igualar o mais barato*, the last one **gated by `evidence_sufficiency`**: where a marketplace gives no sufficient comparables the anchor is absent and the screen explains why (`data-anchor-missing="market"`). Clicking an anchor only fills the input (`data-anchor-apply="fills-input-only"`); confirmation stays explicit.

Honest degradation is rendered throughout: economics `insufficient_evidence` / `unavailable` and market `insufficient` / `unavailable` never masquerade as a number.

Revised proof: verifier **20/20 negative controls** (adds: evaluation detached from the owner operation, per-keystroke evaluation, screen-owned policy judgment, ungated market anchor, auto-applying anchor, stripped waterfall, ledger regaining a create form, evaluate enabled without input); browser operation re-passed end to end (7-line waterfall, 4 indicators, policy floor anchor filling and re-evaluating at R$ 420,59, market anchor absent with explanation on the insufficient-evidence row, confirm → pending → verify → applied, ledger with zero create controls) with **0 console issues** and no 390px overflow.

## 6. Revision 4 — live owner evaluation while typing (2026-08-26)

Third operator alignment. The typed-price box was still inert until an explicit click, the Mercado column had lost the price range, and the waterfall was always expanded. Resolved without breaking the "the screen never computes" law:

- **Live evaluation, still owner-owned.** Typing schedules a **debounced** `EvaluatePriceScenario` call (`data-evaluation-trigger="debounced-owner-call"`, `data-evaluation-debounce-ms="400"`, `data-client-computation="none"`): the owner still produces every number, but the operator sees it update as they work instead of clicking. Never per keystroke.
- **New margin next to the current one.** The row keeps its current contribution/margin column (`data-current-margin-preserved="true"`) and the price cell shows the new contribution + margin with an explicit **percentage-point delta** (`vs atual 9,9% · +4,2 p.p.`), plus the owner policy verdict for the typed price.
- **Market range restored** in the Mercado column (`Faixa entregue: R$ 237,00 – R$ 289,00`) and a positional bar under the live evaluation showing where the typed delivered price falls inside that range.
- **Waterfall behind a disclosure** (`▸ / ▾`, `aria-expanded`, `data-waterfall-disclosure="collapsible"`), collapsed by default and holding the reference anchors, so the fast read stays in the row and the detail is one click away.
- Fixtures were re-based so the demo reads truthfully: an acceptable row, an attention row (below policy and behind), and one with insufficient market evidence — the policy-floor anchor lands exactly on the 12% threshold and *igualar o mais barato* honestly shows 6,4% (below policy), i.e. the real trade-off.

Revised proof: verifier **24/24 negative controls** (adds: debounce removed, screen computing locally, current margin dropped, market range removed, waterfall forced always-open); browser operation re-passed (live update at two different prices, delta in p.p., policy verdict flip at R$ 265, range bar, collapse/expand cycle, anchor refilling and re-evaluating, insufficient-market honesty) with **0 console issues** and no 390px overflow.

## 7. Revision 5 — compact row, competitive rank in the Mercado column (2026-08-26)

Operator finding on revision 4: the live evaluation block was unlabeled ("R$ 17,60 · 5,3%" carried no meaning), and it made the row tall enough that pricing several listings in one pass became unworkable. The requested structure moves the dynamism into the column that already existed instead of growing the row.

Bounded change:

| Concern | Revision 5 structure |
| --- | --- |
| Live indicator label | Explicit `Margem nova (contribuição · % sobre o preço)`; the preserved current margin reads `margem atual R$ … · % · ± p.p.` |
| Row density | Compact: the dashed live box and the inline range bar left the price cell; the cell is input + `✓ Criar intenção` + a `▸/▾` arrow |
| Waterfall placement | Behind that arrow, in the adjacent expandable row (unchanged owner components + reference anchors) |
| Mercado column | Owner-issued delivered range + **our competitive rank** (`Sua posição hoje: 3º de 11`) + positional bar; the typed price updates it live (`Com R$ 265,00: 6º de 11`) |
| Ahead/behind wording | Removed — the rank and the range carry the position without a judgment word |
| Attention basis | `below_policy` OR rank > 1, both server-issued facts; still `data-filter-basis="server-facts"` |

Preserved laws: rank and range come from the Market owner's comparable set (`data-market-cell="range-plus-rank"`), never computed on the screen; economics stays a debounced `EvaluatePriceScenario` call; no bulk apply; supersede-only; anchors fill the input only.

Proof: `node scripts/verify-d6-r-b24-price-intents-wireframe.mjs` → **28/28** negative controls (added: market rank removed, market cell frozen against the typed price, live indicator lost its label, row density expanded back inline). Browser: 0 console issues, no 390px overflow, priced row height 156px.

## 8. P8 operator gate — LOCKED

The operator operated revision 5 in the browser and returned **LOCK on 2026-08-26**. The locked structure is the compact pricing workbench: one row per pricing subject; Economics margin (current, and new-with-label at the typed price) and Market (owner delivered range + our competitive rank, live against the typed price) as read-only owner-separated facts; the waterfall and reference anchors behind a row-level expand/collapse arrow; one explicit supersede-only write per row.

Only the operator can change this LOCK. Assistant, verifier and reviewer output cannot.

## 9. Reopen triggers

Reopen only if operation of this candidate proves that: the supersede model cannot express a real pricing decision; a required pricing fact is missing from the reads; per-variation pricing is materially required (would be a new upstream finding against PriceIntent targets); or responsive/accessible operation cannot preserve the pricing job.
