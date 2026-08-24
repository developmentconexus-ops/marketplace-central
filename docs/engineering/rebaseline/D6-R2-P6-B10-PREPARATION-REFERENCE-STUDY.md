# D6-R2 P6 — B10 Preparation Reference Study

> **Status:** DERIVED — bounded reference evidence for B10; P7 NOT TRIGGERED
> **Parent ledger:** [D6-R2 P8 Structural Wireframe Block Ledger](D6-R2-P8-BLOCK-LEDGER.md)
> **Block:** B10 Preparation / R10
> **Authority boundary:** references inform task structure only; accepted MPC Product/IA remains authoritative
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Question

B10 must support this accepted job without hidden source authority:

```text
exact Organization + Marketplace Installation
→ search admitted source products
→ select exact SourceInstance + native product key
→ inspect readiness + publication requirements
→ resolve/clear correspondence when needed and authorized
→ re-read readiness
→ continue to ListingIntent only when current state permits
```

P5 left three plausible presentation patterns:

```text
A  persistent split-view / master-detail
B  progressive inline detail inside the search workspace
C  search/list → selected-subject detail state
```

P6 studies mature products solving analogous catalog-diagnostics/readiness tasks. It does not copy their visual design or import their business authority.

## 2. References

### R1 — Google Merchant Center — Needs attention

Primary evidence:

- <https://support.google.com/merchants/answer/12476548>
- <https://support.google.com/merchants/answer/12153802>

Observed task pattern:

```text
Products / Needs attention
→ issue prioritization + filters
→ affected-product collection
→ compact issue explanation in collection
→ product / issue detail for material diagnosis and correction
```

Useful lessons:

- triage context remains in the collection;
- product rows expose enough issue/status information to choose the next subject;
- complex explanation/correction moves to a dedicated product or issue detail;
- issue state and business impact are not reduced to a single generic readiness color;
- filters reduce noise without replacing exact product identity.

Mismatch/risk for MPC:

- Merchant Center owns Google-side issue correction mechanisms that MPC does not own;
- MPC must never convert marketplace/provider mutation patterns into local Product writes;
- Google prioritization/click-potential models are not authority for an MPC score.

### R2 — Akeneo — Product Grid + Product Readiness

Primary evidence:

- <https://help.akeneo.com/en_US/product-readiness/product-readiness-overview>
- <https://help.akeneo.com/en_US/serenity-take-the-power-over-your-products/serenity-get-familiar-with-the-product-grid>

Observed task pattern:

```text
product grid
→ readiness/completeness summary per product
→ select exact product
→ product edit/detail form
→ readiness panel exposes unmet requirements
→ correct material product information in the durable subject context
```

Useful lessons:

- a collection can expose a compact readiness summary without becoming the correction workspace;
- the detailed unmet-requirement list belongs beside the exact product identity;
- channel/context qualification matters to readiness;
- durable detail is suitable when several requirement groups may need inspection.

Mismatch/risk for MPC:

- Akeneo is itself a product-information authority/editor; MPC source products are externally owned;
- B10 may resolve only admitted correspondence, not edit arbitrary source-product attributes;
- completeness percentages are not an MPC readiness model unless Product authority explicitly supplies them.

### R3 — Amazon Seller Central — Manage All Inventory / Fix Your Products

Reference evidence:

- <https://sellercentral.amazon.com/seller-forums/discussions/t/a8234e1f-f065-44d8-a4c4-f23cf92b3f4b>
- <https://sellercentral.amazon.com/seller-forums/discussions/t/03276fe9-b943-4308-b2b9-02a76b915dfd>

Observed task pattern from Amazon guidance:

```text
Manage All Inventory
→ filter/search inactive, suppressed or incomplete listing
→ inspect status/reason
→ Fix Your Products / exact listing context
→ material edit/detail when required
```

Useful lessons:

- search/filter first is appropriate when the operator starts from a product identity or known problem;
- reason/status belongs in the collection so the operator can choose the correct subject;
- material corrections are anchored to the exact listing/product context rather than a global inline editor.

Mismatch/risk for MPC:

- Amazon Seller Central mixes marketplace-owned listing mutation with diagnostics; MPC must preserve its Offering/Readiness owner boundary;
- forum/help evidence is weaker than a formal interaction contract and is used only as corroboration.

### R4 — Mirakl / Lowe's Seller Portal — catalog validation feedback

Reference evidence:

- <https://seller.lowes.com/mirakl-seller-documentation/>

Observed task pattern:

```text
catalog submission/import
→ validation feedback/errors
→ row/product-specific issue information
→ correction
→ resubmit/revalidate
```

Useful lessons:

- validation feedback must stay tied to the exact source/import/product identity;
- errors should explain what failed rather than expose only a red/green status;
- revalidation after correction is a distinct state.

Mismatch/risk for MPC:

- the Mirakl example is batch/file-oriented while B10 is an interactive source-product search flow;
- MPC does not introduce upload/import UX or batch correction from this reference.

## 3. Evidence matrix

| Criterion | Google Merchant Center | Akeneo | Amazon Seller Central | Mirakl/Lowe's | MPC implication |
| --- | --- | --- | --- | --- | --- |
| search/filter before subject | strong | strong | strong | batch-oriented | **YES** — B10 remains search-first |
| compact readiness/issue signal in collection | strong | strong | strong | strong | **YES**, but no invented score |
| full material correction inline in collection | limited | no | limited/mixed | no | **NO baseline** |
| dedicated exact-subject detail | strong | strong | strong | error-detail/report analogue | **YES** |
| durable exact identity during correction | strong | strong | strong | strong | **YES** — SourceInstance + native key stay explicit |
| permanent split-view required | no | no | no | no | **NO evidence** |
| re-read/revalidation after correction | yes | yes | eventual refresh | yes | **YES** — readiness re-read stays explicit |

## 4. Derived B10 structure

Reference evidence plus accepted MPC authority supports one leading structure:

```text
B10 Preparation — one search-first workspace

SEARCH MODE
├─ exact Marketplace Installation in page-local context
├─ search input
├─ optional explicit SourceInstance narrowing
└─ structured result list
   ├─ source-product presentation
   ├─ SourceInstance + native key qualification
   ├─ compact readiness state
   ├─ requirements/correspondence signal when known
   └─ open preparation

SELECTED-SUBJECT MODE
├─ exact source identity header
├─ readiness summary / knowledge state
├─ publication requirements
├─ correspondence state
│  └─ resolve / clear only when Product authority admits it
├─ re-read readiness after consequential resolution
└─ continuation to ListingIntent when current state permits
```

`SEARCH MODE` and `SELECTED-SUBJECT MODE` remain inside accepted R10/Preparação. P6 does not choose the exact TanStack route/search-param carrier; P9 binds canonical identity/navigation inputs later.

The search context should remain recoverable when returning from the selected subject, but this is navigation state, not Product truth.

## 5. Alternative disposition

### A — Persistent split-view / master-detail — REJECTED for baseline

Why it remains technically possible but is not materially justified now:

- no accepted evidence yet shows high-frequency rapid cycling across many products; A01 remains OPEN;
- readiness + requirements + correspondence can become vertically deep, making a permanent narrow detail pane fragile;
- responsive/mobile realization would need a mode switch anyway;
- no studied mature reference requires permanent split-view for this class of correction.

A later P12 frequency finding may justify a productivity enhancement without changing Product authority, but P8 does not assume it now.

### B — Full progressive inline expansion — REJECTED for material detail

A compact issue/readiness preview in a result row is useful. Full inline correction is not selected because:

- exact SourceInstance/native identity must stay unmistakable;
- requirements and correspondence resolution can be materially deep;
- recovery/revalidation after resolution deserves durable selected-subject context;
- long expanded rows degrade scanning and responsive behavior.

### C — Search/list → selected-subject detail state — SELECTED CANDIDATE BASIS

This pattern preserves search-first findability while giving the material subject enough space for identity, requirements, correspondence and revalidation.

## 6. P7 trigger disposition

**P7 NOT TRIGGERED.** After P6, only one structure remains materially justified by both references and current MPC evidence.

```text
persistent split-view       plausible implementation technique, not evidence-backed baseline
full inline expansion       insufficient for material detail/recovery
list → selected detail      evidence-backed candidate basis
```

Creating three rendered hypotheses would therefore manufacture ambiguity for ceremony rather than resolve a real product question.

## 7. B10 negative controls entering P8

The B10 HTML candidate must make these visibly invalid:

```text
hidden/default SourceInstance
first search result silently selected
Marketplace Installation used as ambient authority
source product edited as if MPC owned its master data
provider write exposed from Readiness
readiness percentage/health score invented by frontend
known-empty confused with unknown/unavailable/unsupported
correspondence resolution represented as generic product edit
bulk preparation/actions introduced without authority/evidence
permanent split-view treated as mandatory product structure
```

## 8. P6 exit

**DERIVED.** B10 should render one low-fidelity HTML candidate using a search/list → selected exact-subject detail pattern inside the accepted Preparação workspace. P7 is not triggered. No Product operation, Permission, owner or D0–D8 authority is reopened.

## 9. Post-methodology adoption revalidation

**Method profile:** `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD.md + FRONTEND-METHOD.md`.

**Outcome:** `CURRENT STRUCTURE CONFIRMED`.

The accepted methodology adds stronger P7 feasibility and P8 operating/assumption proof, but it does not provide Evidence that falsifies B10's search/list → exact-subject detail structure. **P7 layout hypotheses remain NOT TRIGGERED**: rendering rejected split-view/inline alternatives would manufacture ambiguity.

The pre-adoption HTML does require bounded P8 correction before operator adjudication: its search control is a no-op; it does not inspect a known-empty search result; its Organization selector does not prove invalidation of Installation context; its shell still projects an older IA snapshot; its ListingIntent continuation is a no-op rather than an explicit unopened-block boundary; and its mobile drawer does not fully expose keyboard/focus state. These are prototype-proof defects, not Product capability gaps.

### 9.1 P7 feasibility disposition under the accepted method

| Feasibility question | Current authority | Disposition |
| --- | --- | --- |
| Required fields / summaries | `SearchSourceProductsForMarketplace`, `GetProductChannelReadiness`, `GetPublicationRequirements` expose the source-qualified search/readiness/requirements truth needed by R10. | **PRESENT-IN-AUTHORITY** |
| Identity sources | Organization + exact Marketplace Installation + explicit `SourceInstance` + native product key are accepted identities; no hidden/default source is admitted. | **PRESENT-IN-AUTHORITY** |
| Pagination / scale | Search admits `limit + cursor`; no `total_count` or complete-corpus inference is required. P8 may use deterministic local fixtures without claiming global totals. | **PRESENT-IN-AUTHORITY** |
| Sort / filter | Search query + optional exact `source_instance_id` narrowing are admitted. No material user need or Product contract currently justifies an additional sort surface. | **PRESENT-IN-AUTHORITY** for admitted filters; additional sort **REJECTED — YAGNI / no evidenced need** |
| Preview / content truth | Compact readiness/requirements/correspondence signals are projections of `ProductChannelReadiness`, never frontend scores or source-master edits. | **PRESENT-IN-AUTHORITY** |
| Material writes | B10 may expose only `ResolveProductChannelCorrespondence` / `ClearProductChannelCorrespondence` when authorized. `CreateListingIntentDraft` is the explicit downstream B10 boundary, not a B10 editor implementation. | **PRESENT-IN-AUTHORITY** |

`UPSTREAM FINDING: NONE`.

### 9.2 Lock-time assumption and walkthrough gate

A01 materially influenced rejection of a permanent split-view because real high-frequency rapid cycling is not evidenced. Under the accepted Method its lock-time disposition remains **PENDING OPERATOR**. Before B10 can LOCK, the operator must choose exactly one:

```text
ACCEPT_FOR_LOCK_WITH_LATER_PROBE
BLOCK_LOCK
```

A03 detailed terminology remains an OPEN P12 comprehension probe, but current Evidence does not make it a structural B10 blocker unless the operator walkthrough exposes material confusion. A05 continues to support no invented bulk baseline; absence of Product bulk authority is not converted into a convenience UI.

`Operator walkthrough: PENDING`.

`P8 status: CANDIDATE / NOT LOCKED`.
