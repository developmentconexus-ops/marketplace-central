# D6-R2 P8 — B23 Intenções de Anúncio Candidate

> **Status:** OPERATOR-RATIFIED / LOCKED 2026-08-26 (revision 3)
> **Block:** B23 — ListingIntent authoring / R22–R23
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Candidate evidence:** `qualification/d6-r2-wireframes/b23-listing-intents.html`
> **Prerequisites:** PR #70/#75 read projections + PR #69 B20 closure — ACCEPTED / INTEGRATED
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Design adjudication

The bounded B23 P8 structural design was adjudicated in chat on 2026-08-26 and **approved by the operator**:

- **R22 `/publicacoes/intencoes`** — ListingIntent collection: source product (name/SKU via `source_product_presentation`), target (new listing / existing listing), lifecycle (draft/submitted/discarded) and server-owned dispatchability with blockers;
- **R23** — revision-aware ListingIntent editor: header with source product, target, lifecycle and `requirements_revision` (stale warning when the marketplace requirements moved); per-requirement resolution choice restricted to exactly `follow_source` vs `explicit_override`; media selection from source media and authored media; explicit submit gated by server `dispatchability`; explicit discard;
- **post-submit honesty**: `external_effect_state` not_attempted/pending/accepted/rejected/ambiguous; rejected surfaces the provider feedback verbatim with a fix-and-resubmit path; **ambiguous is never blindly resubmitted** — the only offered action is authoritative verification; `convergence` shows whether the observed listing reflects the request;
- Offering is the sole ListingIntent write owner (`listing.manage`, H/A); creates carry `Idempotency-Key`; writes carry canonical keys only; no live marketplace effect (prototype simulates everything; production live writes require explicit operator authorization);
- no generic form engine, provider field bag, screen-shaped API, bulk selection or saved views.

## 2. Authority consumed

| Screen need | Operation | Permission | Principal kinds |
| --- | --- | --- | --- |
| R22 collection | `ListListingIntents` | `offering.read` | H/A/S |
| create draft | `CreateListingIntentDraft` (+Idempotency-Key) | `listing.manage` | H/A |
| read one intent | `GetListingIntent` | `offering.read` | H/A/S |
| edit draft | `UpdateListingIntentDraft` | `listing.manage` | H/A |
| discard draft | `DiscardListingIntentDraft` | `listing.manage` | H/A |
| submit | `SubmitListingIntent` | `listing.manage` | H/A |
| authored media | `CreateListingIntentMedia` (+Idempotency-Key) | `listing.manage` | H/A |

All operations are owned by Offering. `ListingIntent` carries `source_product_presentation`, `target`, `desired` (lifecycle + requirement_resolutions + media_selection), `resolved_requirements`, `dispatchability` + `dispatch_blockers`, `external_effect_state`, `convergence`, `requirements_revision` and `effect_history` — the candidate binds each to exactly one rendering home.

## 3. Proof

- deterministic verifier `scripts/verify-d6-r-b23-listing-intents-wireframe.mjs`: **10/10 negative controls** (submit gate, ambiguous blind-retry prohibition, resolution-kind closure, canonical-key writes, idempotency, population honesty, stale revision, simulated-effect declaration), wired diff-aware into `npm run gate`;
- browser-operated flow (Chromium): list with 3 lifecycles and honest population states → blocked draft with server blockers and disabled submit → filling "Material" + save re-enables submit via simulated server reread → submit → pending → rejected shows verbatim provider feedback + fix path → ambiguous shows verification-only action and the no-blind-retry law → verification resolves to accepted → stale-revision blocks submit with warning;
- 390px viewport: no horizontal document overflow; mobile drawer navigation law active;
- browser console warnings/errors: **0**.

## 4. Revision 2 — P8 resumed after the variations repair (2026-08-26)

The ratified variations upstream finding was repaired (see the [adjudicated design](../../superpowers/specs/2026-08-26-listing-variations-global-maximum-design.md) and the D4-R1 §16 amendment); B23 P8 resumed applying every PRESENT-IN-AUTHORITY disposition from the [P6 study](D6-R2-P6-B23-LISTING-INTENT-REFERENCE-STUDY.md):

- **publication context region** (category/product type with human descriptors; applying context re-evaluates the census);
- **grouped typed technical sheet**: Obrigatórios with live pendency count → Recomendados → progressive "Exibir todos os campos" for Opcionais/Condicionais; each field rendered by its `value_spec` (text with 60-char title counter, option selects, number+unit, boolean with `not_applicable`, text list, exact decimal);
- **Variações region**: provider axis vocabulary (`variation_axes`: Cor obrigatório com opções, Voltagem opcional), coordinate-keyed options, per-option SKU-level source product, per-option photos and per-variation-scoped fields; per-option dispatch blockers ("Opção Cor: Inox sem foto") resolve interactively; price/quantity explicitly excluded (owners compose at dispatch);
- **primary photo emphasis** in the listing media region;
- R22 list gains a variation summary column.

Revised proof: verifier **14/14 negative controls** (adds: variation write widened past coordinate keys, per-variation scope collapse, price/quantity pulled into authoring, title limit removed); browser operation re-passed end to end (two blockers resolving independently → submit enable → pending; ambiguous law intact) with **0 console issues** and no 390px overflow.

## 5. Revision 3 — category discovery in the context region (2026-08-26)

The operator's category/taxonomy question opened the [category discovery/equivalence evaluation](../../superpowers/specs/2026-08-26-category-discovery-equivalence-design.md); C1/C2 were ratified and repaired (`SearchPublicationContexts`, census 107/31). The context region now renders real discovery: text search, provider prediction with confidence, organization-history suggestion with usage provenance, honest unavailable state, **no preselection** and apply-only-after-explicit-choice. Verifier extended to **16/16** negative controls; browser proof re-passed with 0 console issues.

## 6. P8 operator gate — resolved

The operator walked through revision 3 and returned **`Aprovado`** on 2026-08-26. Under the stated choice set (`LOCK / REVISE / UPSTREAM FINDING`):

```text
final disposition: LOCK on 2026-08-26 (revision 3)
walkthrough rounds: 2 upstream-driven REVISE cycles (variations; category discovery), then LOCK
material revision requested: NONE outstanding
upstream finding raised: NONE outstanding (both ratified findings repaired in-increment)
```

**P8 is OPERATOR-RATIFIED / LOCKED for B23.** This authorized the B23 P9 run, recorded in [`D6-R2-P9-B23-LISTING-INTENTS-SCREEN-CONTRACT.md`](D6-R2-P9-B23-LISTING-INTENTS-SCREEN-CONTRACT.md).

## 7. Reopen triggers

Reopen only if operation of this candidate proves that:

- the two admitted resolution kinds cannot express a real authoring decision;
- a required authoring fact is missing from the ListingIntent read;
- the server dispatchability gate cannot preserve a safe submit decision;
- ambiguous-outcome verification cannot be decided from current truth;
- responsive or accessible operation cannot preserve the authoring job.
