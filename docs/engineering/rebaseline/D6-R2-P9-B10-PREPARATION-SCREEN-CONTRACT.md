# D6-R2 P9 — B10 Preparação Screen Contract

> **Status:** PAUSED — P8 REOPENED
> **Block:** B10 — Preparação / R10
> **Method profile:** `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD + FRONTEND-METHOD`
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Why P9 is paused

The first P9 pass found:

**`F-P9-B10-01` — `known` source evidence does not equal requirement satisfied.**

The initial local correction proposed replacing `Atendido` with `Informação disponível`. The operator challenged that local maximum and required a full Global Maximum evaluation against the DevelopmentConexus Method and current marketplace-hub/provider evidence.

Result: **GLOBAL MAXIMUM REVALIDATED.**

The correct simplification is not another status label and not another backend state. B10 should project:

**requirements + source values + downstream authoring/provider validation**.

The unnecessary intermediate layer is **REJECTED — `source_sufficiency`**.

**NO NEW UPSTREAM WIRE FIELD** is required. The current Product contracts already preserve the essential distinctions; the accidental complexity was in the frontend projection.

Because P8 was materially reopened and the new candidate is not yet operator-LOCKED, this P9 cannot be represented as closed. **rerun P9 after operator re-LOCK**.

## 2. Stable authority that survives the reopen

The bounded rebaseline does not change the accepted route or Product operation ownership.

Canonical workspace:

`/preparacao`

Stable identity/navigation carriers remain:

```text
organization_id
marketplace_installation_id
q
source_instance_id                 optional search narrowing only
selected_source_instance_id
selected_native_product_key
```

The source-qualified selected identity remains exact. Organization changes invalidate Installation and selected-subject context; Installation changes invalidate selected-subject state. URL state is navigation, never business authority.

## 3. Stable Product operation trace

B10 continues to consume:

- `ListMarketplaceInstallations` when the exact page-local account selector needs population;
- `SearchSourceProductsForMarketplace` — owner ProductChannelReadiness, `readiness.read`;
- `GetProductChannelReadiness` — exact subject/correspondence truth, `readiness.read`;
- `GetPublicationRequirements` — provider/context requirement census + source evidence, `readiness.read`;
- `ResolveProductChannelCorrespondence` — `readiness.manage`;
- `ClearProductChannelCorrespondence` — `readiness.manage`.

The downstream operation remains `CreateListingIntentDraft`, owned by Offering with `listing.manage` and `Idempotency-Key` semantics.

**B10 does not call `CreateListingIntentDraft`.** Its continue control is a navigation handoff because B10 does not own the final desired listing representation.

## 4. Simplified screen contract candidate

The new P8 projection deliberately avoids a per-field `satisfied`, `ready`, `met` or sufficiency conclusion.

| Frontend region | Product truth | Human projection |
| --- | --- | --- |
| marketplace account | MarketplacePortfolio / exact Installation | Conta do marketplace |
| source search | `SearchSourceProductsForMarketplace` | Produto + cadastro de origem |
| requirement field | `GetPublicationRequirements.requirement_class` | Campo do marketplace + Exigência |
| source evidence | `GetPublicationRequirements.source_evidence` | Valor encontrado or honest absence/ambiguity |
| downstream handoff | Offering boundary | Na configuração do anúncio |
| correspondence | `GetProductChannelReadiness` + Resolve/Clear | Vínculo com o marketplace |

The provider declaration and source value remain distinct even when the UI places them side by side.

A `known` source value may be carried forward; a `missing`, `unsupported`, `unknown`, `unavailable` or `conflicting` value remains honest evidence. B10 does not need to label any of them as a satisfied requirement.

Missing/conflicting source values do not become publication-impossible by themselves. ListingIntent may later choose `FOLLOW_SOURCE` or `EXPLICIT_OVERRIDE` under Offering authority.

## 5. Correspondence safety remains unchanged

Resolve/Clear continues to use the exact subject + current `correspondence_etag`.

```text
explicit operator action
→ Resolve/Clear
→ no blind retry after ambiguous possible acceptance
→ authoritative GetProductChannelReadiness reread
→ requirement reread when the current context/revision can change
```

This safety property is unrelated to the removed requirement-status UI and remains binding.

## 6. Bidirectional trace to re-prove after LOCK

### frontend → backend

```text
Organization/account context
→ Product scope/navigation

search
→ SearchSourceProductsForMarketplace

Campo do marketplace / Exigência / Valor encontrado
→ GetPublicationRequirements

Vínculo com o marketplace
→ GetProductChannelReadiness
→ ResolveProductChannelCorrespondence / ClearProductChannelCorrespondence

Continuar para configurar o anúncio
→ downstream navigation boundary
→ NOT a B10 CreateListingIntentDraft call
```

### backend → frontend

```text
SearchSourceProductsForMarketplace
→ source-qualified result list

GetProductChannelReadiness
→ exact subject + correspondence + reread basis

GetPublicationRequirements
→ marketplace fields + obligation/type/value constraints + source values/evidence

Resolve/Clear correspondence
→ explicit correspondence controls + authoritative reread

CreateListingIntentDraft
→ no B10 mutation home; downstream Offering authoring only
```

The final P9 must re-check these mappings against the exact operator-LOCKED P8 artifact. This document records the stable candidate trace only; it does not substitute for that rerun.

## 7. Forbidden shortcuts

B10 must not:

- create a Product/PIM master;
- introduce `source_sufficiency` or another generic per-field readiness state;
- compute `satisfied` in React from requirement/source data;
- add a provider field bag;
- treat missing source data as publication impossible;
- turn provider validation into MPC-owned universal business truth;
- blindly retry ambiguous correspondence writes;
- manufacture ListingIntent desired state merely to leave B10.

## 8. Current disposition

```text
F-P9-B10-01
→ Global Maximum revalidation
→ source_sufficiency rejected
→ simplified P8 operator projection
→ P8 REOPENED / CANDIDATE
→ operator walkthrough + re-LOCK required
→ rerun P9 after operator re-LOCK
→ P10 only after P9 closure
```

**P9: PAUSED — P8 REOPENED.**
