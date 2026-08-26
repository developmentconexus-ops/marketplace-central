# D6-R2 P9 — B23 Intenções de Anúncio Screen Contract

> **Status:** DERIVED / PASS — P8 LOCKED 2026-08-26; BACKEND SUFFICIENT; UPSTREAM FINDING NONE
> **Block:** B23 — ListingIntent authoring / R22–R23
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked P8 evidence:** `qualification/d6-r2-wireframes/b23-listing-intents.html`
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml` (107/31/H-A-S)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 result

P9 ran only after the operator LOCKed the revision-3 B23 candidate. Both in-increment upstream findings (listing variations; publication-context discovery/equivalence) were ratified, repaired and consumed by the locked candidate before this run.

The locked human job is:

```text
browse/create listing intents for the organization
→ open one intent (draft or submitted)
→ choose the marketplace category with discovery/suggestion, explicitly
→ resolve the typed, grouped requirement census (follow_source vs explicit_override, not_applicable where allowed)
→ author variations from the provider axis vocabulary (per-option SKU/photos/fields)
→ select listing media (source + authored, primary emphasized)
→ submit only when server dispatchability allows; discard explicitly
→ read the honest external outcome; verify (never blindly retry) on ambiguity; fix and resubmit on rejection
```

**P9 verdict: PASS / BACKEND SUFFICIENT / UPSTREAM FINDING NONE.**

## 2. Route, identity and client-state ownership

Route family: `/org/:organizationId/publicacoes/intencoes` (R22) and `/publicacoes/intencoes/:listingIntentId` (R23).

| State class | B23 ownership |
| --- | --- |
| `GLOBAL_WORKSPACE_CONTEXT` | `organization_id` |
| `URL_NAVIGATION_STATE` | selected `listing_intent_id`, list filter |
| `SERVER_STATE` | intent collection/detail, publication requirements + variation axes, context candidates, dispatchability, effect state |
| `LOCAL_EPHEMERAL` | unsent field edits, disclosure toggles, mobile navigation |

Prototype scenario selectors are Evidence controls. TanStack Query remains the server-state owner; unsent edits become truth only through explicit draft writes.

## 3. Product operation / access binding

| Screen need | Operation | Semantic owner | Permission | Principal kinds |
| --- | --- | --- | --- | --- |
| R22 collection | `ListListingIntents` | Offering | `offering.read` | H/A/S |
| create draft | `CreateListingIntentDraft` + Idempotency-Key | Offering | `listing.manage` | H/A |
| read one intent | `GetListingIntent` | Offering | `offering.read` | H/A/S |
| census + variation axes | `GetPublicationRequirements` | ProductChannelReadiness | `readiness.read` | H/A/S |
| category discovery/suggestion | `SearchPublicationContexts` | ProductChannelReadiness | `readiness.read` | H/A/S |
| edit draft (fields, context, variations, media selection) | `UpdateListingIntentDraft` | Offering | `listing.manage` | H/A |
| authored media | `CreateListingIntentMedia` + Idempotency-Key | Offering | `listing.manage` | H/A |
| submit | `SubmitListingIntent` | Offering | `listing.manage` | H/A |
| discard | `DiscardListingIntentDraft` | Offering | `listing.manage` | H/A |

All writes carry canonical keys only (requirement/axis/option/category keys, candidate keys); no presentation label is ever a write carrier. Live provider effects require explicit operator authorization; the prototype simulates all effects.

## 4. Material screen contract (summary of locked regions)

- **R22 collection:** source product presentation (name/SKU), target (new/existing listing), variation summary, lifecycle (draft/submitted/discarded), server dispatchability with first blocker; honest known-empty ≠ unknown ≠ unavailable; cursor grammar; no bulk selection.
- **Categoria no marketplace:** `SearchPublicationContexts` candidates with closed `suggestion_basis` (provider prediction + confidence / text search / organization history + usage provenance); no preselection; apply only after explicit choice; unavailable ≠ nonexistent; provider capability variance honest per the hub-neutrality clause.
- **Ficha técnica:** grouped by `requirement_class` with pendency counts and progressive show-all; typed inputs per the seven value specs; per-field `follow_source` vs `explicit_override`; `not_applicable` where allowed; title constrained with counter.
- **Variações:** provider `variation_axes` vocabulary; coordinate-keyed options with optional SKU-level source refs, per-option media and per-variation-scoped fields; per-option dispatch blockers; price/quantity explicitly excluded (owners compose at dispatch).
- **Envio/resultado:** submit gated by server `dispatchability` with explicit blockers; `external_effect_state` pending/accepted/rejected(+verbatim provider feedback, fix path)/ambiguous(verification only — never blind retry); `convergence` projection; stale `requirements_revision` warning blocks submission until review.

## 5. Bidirectional trace — PASS

Every admitted operation above lands on exactly one locked surface region; no locked control lacks an operation; no admitted operation is orphaned. Reads that shape safety (dispatchability, effect state, revision) are server truth rendered, never client-computed.

## 6. Adversarial checks

P9 rejects: client-computed dispatchability; label-carrying writes; auto-applied category or history suggestion; a third resolution kind; per-option price/stock authoring; blind retry after ambiguity; hiding provider rejection feedback; census collapse of unavailable into empty. All are excluded by the locked evidence and `verify-d6-r-b23-listing-intents-wireframe.mjs` (16/16 negative controls) plus the projection proof (18/18).

## 7. P9 closure and P10 note

```text
P8 OPERATOR-RATIFIED / LOCKED (2026-08-26, revision 3)
→ exact route/state/identity binding
→ exact owner/operation/Permission binding
→ frontend → backend trace PASS
→ backend → frontend trace PASS
→ adversarial shortcuts rejected
→ BACKEND SUFFICIENT
→ UPSTREAM FINDING NONE
```

**P9: PASS / CLOSED for B23.**

P10: B23 reuses the established laws (typed presentation ≠ canonical ref; known-empty ≠ unknown ≠ unavailable; navigation ≠ mutation; server-gated consequential writes; ambiguous-verify). The grouped typed requirement editor and the suggestion-candidate list are B23-local until a second locked block proves the same shape; **no new shared component/pattern authority is claimed.** P11, Pre-D9/D9 and Product implementation remain outside this closure.
