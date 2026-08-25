# D6-R2 P6 — B10 Preparation Reference Study

> **Status:** DERIVED / GLOBAL MAXIMUM REVALIDATED — P8 REOPENED / CANDIDATE
> **Block:** B10 Preparation / R10
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Human job and accepted boundary

B10 supports a deliberately small preparation job:

```text
exact Organization + Marketplace Installation
→ search admitted source products
→ select exact source-qualified product
→ inspect marketplace-required fields and source values
→ resolve/clear Product↔channel correspondence when needed
→ continue to ListingIntent configuration
```

B10 does **not** edit a source Product master, become a PIM, decide final ListingIntent dispatchability, validate a complete provider payload, publish to the provider, or become a second marketplace rule engine.

The Product distinction remains:

```text
provider requirement
!= source evidence/value
!= ListingIntent desired value
!= provider acceptance/convergence
```

And specifically:

```text
missing source != publication impossible
```

## 2. Earlier structure and PR #68 result

The previously selected navigation structure remains sound:

```text
search/list → selected exact-subject detail
```

**CURRENT STRUCTURE CONFIRMED.**

The earlier alternatives — persistent split view and progressive inline expansion — remain rejected because no evidence justifies their extra responsive and identity complexity. **P7 layout hypotheses remain NOT TRIGGERED.**

PR #68 is integrated into `main` at:

`ed3d164b0574b7950c2c7467d150c89576bba1ec`

It correctly preserves provider/context-qualified publication requirements, independent `requirement_class` and `applicability`, six source-evidence states, seven bounded value-spec families, source-candidate identity and source-media separation. That upstream wire repair remains accepted and is not reopened here.

## 3. P9 trigger

The first locked human-first P8 projected source evidence into a separate per-requirement `Situação` column and used labels such as `Atendido`.

P9 correctly proved one semantic overclaim:

```text
source_evidence.state = known
!=
provider/Readiness fact that the requirement is satisfied
```

The first proposed correction was only to rename `Atendido` to `Informação disponível`. Operator challenge then exposed a broader question: whether B10 needs a per-requirement satisfaction layer at all.

That question is material because a convenience layer could duplicate marketplace validation, turn Readiness into a generic rules engine, or force the frontend to explain states that provide no additional operator value.

## 4. Fresh market/reference evidence

The Global Maximum was re-evaluated against current marketplace-hub patterns and the Mercado Livre contract itself.

### AnyMarket

Official support material models channel preparation around marketplace attributes that are mandatory/optional, attribute/value binding and explicit pendencies. Missing or invalid required data is corrected from marketplace feedback; the hub does not need a separate universal per-field sufficiency ontology.

Reference:
- <https://suporte.anymarket.com.br/hc/pt-br/articles/51175297537555-Como-funciona-a-Central-de-v%C3%ADnculos>

### Hub2b

Official category-attribute guidance exposes the marketplace attribute, whether it is mandatory, its type and allowed values. Values are then supplied from the source system or Hub2b authoring context.

Reference:
- <https://suporte.hub2b.com.br/hc/pt-br/articles/360030486354-Atributos-por-categoria>

### Magis5

Official publication guidance treats omitted mandatory information as a publication/pre-publication error and lets channel-specific review remain channel-specific. The useful operator model is required information + supplied value + marketplace feedback, not another generic readiness state machine.

Reference:
- <https://ajuda.magis5.com.br/como-publicar-anuncios-em-massa-na-magis5>

### Mercado Livre

Mercado Livre itself supplies category/context attribute requirements, including required and conditional-required behavior, and owns validation of the effective item payload. Conditional requirements may depend on the item context; provider validation remains the strongest authority for the actual provider payload.

References:
- <https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/atributos>
- <https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/validador-de-publicacoes>

### Evidence conclusion

The common sustainable pattern is:

```text
marketplace requirements
+ source/catalog values when available
+ explicit mapping/authoring for what is missing or ambiguous
+ downstream marketplace-specific validation/feedback
```

The evidence does **not** justify a generic intermediate `source_sufficiency` concept.

## 5. Decision Core — GLOBAL MAXIMUM REVALIDATED

### Root cause

B10 had started to convert technical honesty into accidental Product complexity:

```text
requirement + source evidence
→ invented operator satisfaction state
→ pressure for another backend/client semantic layer
```

The actual operator job is simpler: see what the channel asks for, see what value exists now, and know what will need to be completed or chosen while configuring the listing.

### Target invariant

> **B10 must faithfully expose marketplace requirements and available source values, preserve missing/conflicting/unknown evidence honestly, and hand unresolved listing values to Offering authoring without creating a second validation authority. Final draft validity and provider acceptance remain with their accepted owners.**

### Alternatives

| Alternative | Decision | Reason |
| --- | --- | --- |
| add per-requirement `source_sufficiency` | **REJECTED — `source_sufficiency`** | duplicates meaning, adds states with little operator value, and risks becoming a generic rules engine |
| frontend computes `satisfied` from `source_evidence + value_spec` | REJECTED | frontend would become parallel business authority |
| send blindly and rely only on provider rejection | REJECTED | throws away cheap declared requirement/type/value constraints already available before dispatch |
| **requirements + source values + downstream authoring/provider validation** | **CURRENT STRUCTURE CONFIRMED / GLOBAL MAXIMUM** | smallest model that keeps provider truth, source truth, authoring and final validation distinct |

**NO NEW UPSTREAM WIRE FIELD** is required for this decision. `GetPublicationRequirements` already carries the requirement/value constraints and source evidence needed by B10; Offering/ListingIntent already owns downstream desired values and dispatchability.

This is removal of accidental complexity, not removal of essential provider complexity.

## 6. Revised P8 operator projection

The primary detail table becomes:

```text
Campo do marketplace
Exigência
Valor encontrado
Na configuração do anúncio
```

Examples:

```text
Marca       | Obrigatório | Deca          | usar o valor encontrado
Material    | Obrigatório | —             | preencher na configuração do anúncio
Acabamento  | Recomendado | Cromado / inox| escolher na configuração do anúncio
```

No `Situação`, `Atendido`, readiness score or per-field sufficiency status is needed.

Technical evidence remains encoded beneath the UI so that:

- `requirement_class != applicability`;
- `known / missing / conflicting / unknown / unavailable / unsupported` remain distinct;
- all seven `value_spec` families remain representable;
- opaque source candidates remain distinguishable;
- `not_applicable` remains override-only;
- source media stays separate;
- no raw `provider_fields` bag appears.

Missing/conflicting source values do not by themselves block entry into ListingIntent authoring. B10 may still fail closed when its own exact identity/correspondence/read authority cannot be established.

## 7. Correspondence and downstream boundary

Correspondence stays inside ProductChannelReadiness:

```text
Resolve/Clear correspondence
→ no blind retry after ambiguous possible acceptance
→ authoritative reread
```

The B10 continue control remains a navigation handoff. B10 does not manufacture Offering-owned desired state and does not call `CreateListingIntentDraft` merely to leave the screen.

ListingIntent later chooses `FOLLOW_SOURCE` or `EXPLICIT_OVERRIDE` and completes missing/ambiguous values as required. Provider-specific validation belongs after the effective draft/payload exists, not in B10.

## 8. P8 reopen and LOCK impact

The earlier human-first P8 was operated and later operator-ratified. Its A01 assumption was explicitly accepted:

`A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE`.

The operator has now authorized this bounded Global Maximum rebaseline. Because the visible requirement model changes materially, the prior B10 LOCK is reopened only for B10.

**P8 REOPENED / CANDIDATE.**

A01 is **not** reopened: the current change does not alter the search/list → exact-subject structure that A01 governs.

LOCK impact sweep:

```text
B00 | UNAFFECTED
B01 | UNAFFECTED
B00-R2 | UNAFFECTED
B11 | UNAFFECTED
B12 | UNAFFECTED
B110 | UNAFFECTED
B10 | REOPENED — bounded requirement/value projection only
```

No D0–D5 Product authority, Product operation, Permission, Principal kind, D7/D8 authority or implementation authorization is changed.

## 9. Current gate

The revised low-fi remains disposable P8 evidence and must be operated again before any new LOCK.

```text
proof revised P8 candidate
→ next gate: operator walkthrough
→ operator LOCK / REVISE / UPSTREAM FINDING
→ rerun P9 only after operator re-LOCK
```

Do not begin P10, Pre-D9/D9 or Product implementation from this candidate state.
