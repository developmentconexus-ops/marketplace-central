# D6-R2 P8 — B10 Preparação Ratification / Reopen Record

> **PRIOR P8: OPERATOR-RATIFIED / LOCKED**
> **CURRENT P8: REOPENED / CANDIDATE**
> **Block:** B10 — Preparação / R10
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Candidate evidence:** `qualification/d6-r2-wireframes/b10-preparation.html`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Historical operator LOCK

The prior human-first B10 candidate was operated and explicitly LOCKED by the operator after the earlier backend-language revision.

Historical walkthrough:

```text
OPERATED
actual task attempted: inspect the revised human-first Preparation flow as a normal operator
material issues found: none after the human-language revision
final disposition: LOCK
```

That decision remains historical evidence; it is not erased by the current bounded reopen.

## 2. A01 remains accepted

The prior lock-time assumption remains:

`A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE`.

A01 concerns the choice of search/list → exact-subject detail rather than a persistent split view. The current simplification does not falsify that structure, so A01 is not reopened.

## 3. Why the LOCK reopened

P9 exposed that `source_evidence.state = known` had been projected as `Atendido`. The operator then required a Global Maximum evaluation rather than a cosmetic label change.

Fresh marketplace-hub/provider research showed that the simpler sustainable model is:

```text
marketplace requirement
+ source value when available
+ downstream listing authoring for missing/ambiguous values
+ provider-specific validation/feedback
```

The proposed `source_sufficiency` layer was rejected as overengineering. No new Product wire field is required.

The operator explicitly approved this **operator-authorized bounded rebaseline**.

Because the visible requirement model changes from five columns with a satisfaction/status projection to four columns focused on requirement/value/handoff, the old B10 LOCK cannot silently carry forward.

## 4. Current candidate boundary

The reopened P8 preserves:

- search/list → exact-subject detail;
- exact Organization + Marketplace Installation;
- no ambient/default SourceInstance;
- known-empty distinct from unknown/unavailable;
- complete provider/context-specific requirement census;
- PR #68 class/applicability/value/source distinctions beneath the UI;
- missing source truth not becoming publication-impossible;
- correspondence Resolve/Clear under ProductChannelReadiness with authoritative reread;
- ListingIntent as downstream authoring boundary;
- technical identifiers/details as secondary disclosure;
- B00/B01/B00-R2/B11/B12/B110 as UNAFFECTED.

The primary operator projection now removes per-requirement satisfaction semantics and presents:

```text
Campo do marketplace
Exigência
Valor encontrado
Na configuração do anúncio
```

## 5. Current gate

**CURRENT P8: REOPENED / CANDIDATE.**

The revised browser-operable candidate must pass its structural proof and then a fresh **operator walkthrough required** gate.

Only the operator may choose the next P8 disposition:

```text
LOCK
REVISE
UPSTREAM FINDING
```

P9 remains paused while P8 is reopened. Do not begin P10, Pre-D9/D9 or Product implementation before a fresh operator LOCK and P9 rerun.
