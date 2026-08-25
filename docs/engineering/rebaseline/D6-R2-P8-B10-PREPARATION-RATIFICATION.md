# D6-R2 P8 — B10 Preparação Ratification

> **CURRENT P8:** MAIN STRUCTURE OPERATOR-RATIFIED / LOCKED; CORRESPONDENCE REGION BOUNDED REOPEN / CANDIDATE
> **Block:** B10 — Preparação / R10
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked evidence:** `qualification/d6-r2-wireframes/b10-preparation.html`
> **Bounded candidate:** [`D6-R2-P8-B10-CORRESPONDENCE-REVALIDATION.md`](D6-R2-P8-B10-CORRESPONDENCE-REVALIDATION.md)
> **Operator disposition:** `LOCK` on 2026-08-25
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Why B10 was reopened

The prior B10 LOCK projected `source_evidence.state = known` as `Atendido`. P9 showed that source evidence did not authorize a per-requirement satisfaction conclusion.

The operator required a Global Maximum re-evaluation rather than a cosmetic label patch. Current marketplace-hub/provider evidence supported the smaller sustainable model:

```text
marketplace requirement
+ source value when available
+ downstream listing authoring for missing/ambiguous values
+ provider-specific validation/feedback
```

The proposed generic `source_sufficiency` layer was rejected as accidental complexity. No new Product operation, Permission, Principal kind or wire field was required.

## 2. Operator walkthrough and LOCK

The simplified browser-operable candidate was supplied to the operator for the fresh P8 walkthrough required by the Frontend Product Experience Planning Method.

The candidate exposed:

```text
Campo do marketplace
Exigência
Valor encontrado
Na configuração do anúncio
```

and kept the search/list → exact-subject detail → correspondence → ListingIntent boundary executable.

After receiving the exact candidate, the operator returned **`Aprovado`** as the disposition. Under the previously stated P8 choice set (`LOCK / REVISE / UPSTREAM FINDING`), that approval is recorded as:

```text
final disposition: LOCK
material revision requested: NONE
upstream finding raised: NONE
```

## 3. Locked B10 structure

The current LOCK protects:

- `OFERTA > Preparação` as the B10 placement;
- search/list → exact source-qualified subject detail;
- exact Organization + Marketplace Installation context;
- no ambient/default SourceInstance;
- known-empty distinct from unknown/unavailable;
- complete provider/context-specific marketplace requirement census;
- primary projection `Campo do marketplace / Exigência / Valor encontrado / Na configuração do anúncio`;
- no `Situação`, `Atendido`, readiness score or per-field sufficiency label;
- `known / missing / conflicting / unknown / unavailable / unsupported` remaining distinguishable beneath the human projection;
- missing/conflicting source data not becoming publication-impossible;
- Product↔channel correspondence as a separate ProductChannelReadiness concern;
- explicit Resolve/Clear correspondence followed by authoritative reread;
- `Continuar para configurar o anúncio` as a downstream Offering/ListingIntent boundary, not a publication effect;
- technical identifiers/details as secondary support disclosure rather than primary operator vocabulary;
- responsive search→detail/mobile-stack behavior without changing authority meaning.

The LOCK does **not** freeze final visual design, production components, final typography/palette, provider payload validation UI, or the unopened ListingIntent editor.

## 4. A01 disposition

`A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE` remains unchanged.

A01 concerns the search/list → exact-subject detail structure and relative task-frequency/density evidence. The bounded requirement/value simplification did not falsify that structure.

## 5. Authority and reopen triggers

No D0–D5 Product authority, Product operation, Permission, Principal kind, D7/D8 authority or implementation authorization changes with this LOCK.

Reopen B10 only on material evidence such as:

- a real provider requirement cannot be represented without losing material meaning;
- a real operator task proves the search/detail structure materially insufficient;
- source/provider evidence cannot remain honest in the four-column projection;
- downstream ListingIntent authoring proves a required B10 handoff identity/state is missing;
- P9 finds a contradiction against the locked experience or accepted Product authority.

## 6. Current bounded method gate

The main B10 structure remains locked. PR #70 triggered only the previously declared correspondence-region reopen condition. The bounded candidate must be operated by the operator before that region can be re-LOCKED.

```text
main B10 structure OPERATOR-RATIFIED / LOCKED
+ correspondence-region CANDIDATE
→ operator walkthrough
→ LOCK / REVISE / UPSTREAM FINDING
→ only after re-LOCK: rerun bounded P9 trace
```

P10/P11, B20, Pre-D9/D9 and Product implementation may not bypass the bounded re-LOCK and P9 rerun.
