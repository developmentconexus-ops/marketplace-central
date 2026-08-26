# D6-R2 P8 — B24 Preços Candidate

> **Status:** CANDIDATE / OPERATOR WALKTHROUGH REQUIRED — no LOCK claimed
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

## 6. P8 operator gate

Current disposition:

```text
CANDIDATE
operator walkthrough required
LOCK / REVISE / UPSTREAM FINDING
```

The assistant, verifier and reviewer cannot set `LOCKED`. P9 must not run for B24 until the operator explicitly LOCKs this candidate.

## 7. Reopen triggers

Reopen only if operation of this candidate proves that: the supersede model cannot express a real pricing decision; a required pricing fact is missing from the reads; per-variation pricing is materially required (would be a new upstream finding against PriceIntent targets); or responsive/accessible operation cannot preserve the pricing job.
