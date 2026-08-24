# D6-R2 Fable R-2 — OAD Source Reachability Hygiene

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Operator adjudication:** 2026-08-24 — Fable R-2 accepted; continue bounded remediation
> **Trigger:** independent methodology + whole-repository Fable review `review/d6r2-methodology-whole-repo-fable @ 975a9176cb4019960166e436d63e33349819c046`
> **Canonical Product wire:** [`contracts/api/product/openapi.yaml`](../../../contracts/api/product/openapi.yaml)
> **Historical debt manifest:** [`source-orphan-allowlist.json`](../../../contracts/api/product/source-orphan-allowlist.json)
> **Verifier:** `scripts/verify-oad-source-reachability.mjs`
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · H/A/S only
> **Implementation:** BLOCKED UNTIL accepted D9

```text
R2_POLICY:EXACT_FROZEN_HISTORICAL_ALLOWLIST
R2_ALLOWED_ORPHANS:49
R2_NEW_ORPHANS:0
R2_PATHITEMS:16
R2_SCHEMAS:33
R2_ORPHAN_SNAPSHOT_SHA256:2ccb6d5716659cd46cc652358d55b125c96e154724163be7d6850bb501b8aed7
R2_PRODUCT_BUNDLE_SHA256_AT_CLOSE:7ad4581686879517d0585f54896e618b3a5b1cc57bd7576c629c1015bf330a02
R2_NEW_PRODUCT_SEMANTICS:0
```

## 1. Finding

Fable found that supersession had left Product-source definitions unreachable from the current `openapi.yaml` graph. A future editor or source-text verifier could touch an old copy believing it current even though the generated canonical bundle prunes it.

A RED source-graph proof measured:

```text
16 unreachable source pathItems
33 unreachable source schemas
49 total unreachable enforced definitions
```

The initial hypothesis was to prune every unreachable definition.

## 2. Diagnostic falsification of blind pruning

The RED exposed an important distinction: several unreachable definitions are still deliberate inputs to accepted historical non-regression proofs.

The 95/29 and 99/30 proof harnesses reconstruct earlier Product generations by rewinding current root references to superseded source symbols, including prior Access, AuthorizationDecision and Work definitions. Deleting all 49 from the current source tree would therefore either break accepted historical proof or force an unrelated proof-fixture refactor.

The current Product bundle and historical proof fixture are different authority roles:

```text
openapi.yaml reachable graph / generated bundle
→ CURRENT Product wire authority

exact allowlisted unreachable definitions
→ HISTORICAL proof source debt only
→ never current Product authority
```

## 3. Alternatives

### A — delete every unreachable source definition now

**REJECTED.** It confuses current-root reachability with historical-proof usage and would create avoidable proof-infrastructure churn.

### B — immediately extract every historical generation into a separate fixture tree, then require zero current-tree orphans

**DEFERRED.** Architecturally clean, but current evidence does not justify the migration cost now. It adds no Product capability or correctness beyond a smaller bounded control.

Reopen this extraction when retained historical source debt materially impedes maintainability/tooling or a legitimate Product change needs one of the retained symbols.

### C — exact frozen historical allowlist + no growth/content drift

**SELECTED / GLOBAL MAXIMUM UNDER CURRENT YAGNI PRESSURE.** It closes the actual failure mode without deleting evidence or creating a proof-infrastructure project.

## 4. Binding source-hygiene law

`contracts/api/product/source-orphan-allowlist.json` contains the exact retained set.

The gate requires all of the following:

```text
current unreachable enforced definitions == exact manifest set
pathItems == 16
schemas   == 33
total     == 49
orphan source-content digest == manifest digest
new orphan definitions == 0
```

Therefore:

- a 50th orphan fails;
- silent mutation of one retained historical definition fails;
- silent removal/migration of one retained definition fails until the manifest is deliberately updated;
- current source code must not infer Product authority merely because a definition exists in a YAML file;
- the canonical root/bundle remains the current wire authority.

A future intentional extraction may reduce or eliminate the manifest, but it must preserve the historical non-regression contract deliberately rather than by accidental source retention.

## 5. RED → GREEN proof

### RED

The first detector reached all existing Product/P8/P9/Fable/D7-R/D8-R/D5-R7 proofs and failed on the measured 49-definition orphan set.

It recorded:

```text
bundle SHA-256
7ad4581686879517d0585f54896e618b3a5b1cc57bd7576c629c1015bf330a02

orphan snapshot SHA-256
2ccb6d5716659cd46cc652358d55b125c96e154724163be7d6850bb501b8aed7
```

### GREEN

At `a8b1e7d0a5b0245577da872789217f6643b24b4c`, CI #639 and pr-title #715 succeeded.

Observed proof:

```text
95/29 historical Product proof       PASS
99/30 historical Product proof       PASS
106/31 current Product proof          PASS
D5-R7 W1 carrier proof                PASS
canonical bundle SHA-256              7ad4581686879517d0585f54896e618b3a5b1cc57bd7576c629c1015bf330a02
allowed orphan pathItems              16
allowed orphan schemas                33
allowed orphan total                  49/49
new orphans                           0
orphan snapshot SHA-256               2ccb6d5716659cd46cc652358d55b125c96e154724163be7d6850bb501b8aed7
source reachability negative controls 2/2
source reachability                   PASS
```

The bundle hash is identical to the diagnostic RED baseline, proving this hygiene closure did not alter the canonical Product bundle.

## 6. Result

R-2 is closed without Product semantic change, OAD census change, new Permission, source pruning by guess, or weakening historical proof.

The independent whole-repository review's repo-specific MATERIAL/IMPORTANT findings R-1 and R-2 are now closed. Global frontend-method evolution remains a separate operator-owned workstream and must be reconciled into repository authority before B10 proceeds.
