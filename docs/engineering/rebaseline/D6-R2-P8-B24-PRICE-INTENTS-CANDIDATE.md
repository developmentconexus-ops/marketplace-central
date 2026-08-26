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

## 4. P8 operator gate

Current disposition:

```text
CANDIDATE
operator walkthrough required
LOCK / REVISE / UPSTREAM FINDING
```

The assistant, verifier and reviewer cannot set `LOCKED`. P9 must not run for B24 until the operator explicitly LOCKs this candidate.

## 5. Reopen triggers

Reopen only if operation of this candidate proves that: the supersede model cannot express a real pricing decision; a required pricing fact is missing from the reads; per-variation pricing is materially required (would be a new upstream finding against PriceIntent targets); or responsive/accessible operation cannot preserve the pricing job.
