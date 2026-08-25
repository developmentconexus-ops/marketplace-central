# Repository Governance & Context Health — Execution Audit

> **Status:** TASK 1 / NO-DELETE CHECKPOINT  
> **Baseline:** `main@bdbbef43ed3a5e9d912e67ddac5173024352eaa3`  
> **Health branch:** `chore/repository-health-context-governance`  
> **Rule:** no path classified `REHOME THEN RETIRE` may be removed before its surviving meaning is consolidated into the named current owner. This file is temporary execution evidence and must retire with the health spec/plan before PR #71 is integration-ready.

## 1. Material finding from the audit

The dominant repository-health defect is **canonical consolidation drift**, not merely file count.

The later NOTIF-01 / AuthorizationRequest work is accepted and reflected in the current Product OAD (`106 operations / 31 Permissions`), but major semantic parents are still pre-amendment snapshots:

- canonical D0 does not yet carry the accepted Personal Notifications Product scope;
- canonical D1 still presents the pre-NOTIF boundary catalog and edge set;
- canonical D2 does not yet carry Notification / later AuthorizationRequest identity-data semantics;
- canonical D3 remains the pre-NOTIF communication matrix;
- canonical W4 still states the historical `95 operations / 29 Permissions` vocabulary;
- P5 still states the historical `99 / 99` surface basis and predates the final Notifications / AuthorizationRequest frontend rebaseline.

Therefore the safe cleanup order is:

```text
accepted intermediate meaning
→ compact canonical rehome in D0/D1/D2/D3/D5/D6/D7/D8
→ update current P5/P8/P9 references where needed
→ retire intermediate chain
```

The Product OAD itself is not changed by this health work.

## 2. D6-R2 / NOTIF-01 / AuthorizationRequest audit

| Path | Class | Surviving meaning | Replacement owner / live consumer | Action |
| --- | --- | --- | --- | --- |
| `docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md` | RETIRE | none; route-pack-only reachability | `docs/index.md` direct current-owner routing | rm |
| `docs/engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md` | REHOME THEN RETIRE | D6-R2 tailoring, user-needs/flows and current frontend-realization laws | `D6-FRONTEND.md` + `D6-R2-P5-SCREEN-SURFACE-INVENTORY.md` + current block evidence | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md` | RETIRE | accepted IA result already expressed by P5 and locked P8 evidence | `D6-FRONTEND.md` + P5 + current P8 registry/evidence | rm |
| `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md` | KEEP CURRENT EVIDENCE | stable operator LOCK registry, including B00/B01 where no smaller ratification record owns the lock | active D6-R2 P8 navigation; compact to stable lock/evidence registry and remove mutable status chronology | keep |
| `docs/engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md` | RETIRE | reference study completed; no current authority | B10 ratification + B10 P9 + locked HTML | rm |
| `docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md` | KEEP CURRENT EVIDENCE | operator LOCK and accepted B10 structural meaning | active bounded B10 reopen / future re-LOCK comparison | keep |
| `docs/engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md` | KEEP CURRENT EVIDENCE | current B10 P9 baseline that must be rerun after upstream repair | active bounded B10 reopen | keep |
| `docs/engineering/rebaseline/D6-R2-P5-B110-AUTHORIZATION-REQUEST-SUPERSESSION.md` | REHOME THEN RETIRE | current P5 approvals/actionable-request surface correction | `D6-R2-P5-SCREEN-SURFACE-INVENTORY.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-P8-B110-APPROVALS-CANDIDATE.md` | RETIRE | pre-ratification candidate only | B110 ratification + locked HTML | rm |
| `docs/engineering/rebaseline/D6-R2-P8-B110-APPROVALS-RATIFICATION.md` | KEEP CURRENT EVIDENCE | operator LOCK and accepted approvals structure | active D6-R2 acceptance evidence | keep |
| `docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md` | KEEP CURRENT EVIDENCE | final current Notifications/Approvals P9 trace and frontend/backend binding | active D6-R2 acceptance evidence | keep |
| `docs/engineering/rebaseline/D6-R2-AUTHORIZATION-REQUEST-FABLE-ADJUDICATION.md` | REHOME THEN RETIRE | accepted corrections to AuthorizationRequest/Product/P9 composition not all consolidated into canonical parents | D2/D3/W1/W2/W4 + retained P9 contract | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-AUTHORIZATION-REQUEST-FABLE-RATIFICATION.md` | RETIRE | ratification wrapper after adjudication | canonical rehome + retained P9 | rm |
| `docs/engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md` | REHOME THEN RETIRE | accepted typed AuthorizationRequest ETag carrier and bounded B110 revalidation | `D5-B2-WIRE-CONTRACT.md` + W2/W4 as applicable + current OAD/P9 | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-FABLE-R2-OAD-SOURCE-REACHABILITY-HYGIENE.md` | REHOME THEN RETIRE | current source-reachability/orphan hygiene rule | `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md` + `verify-oad-source-reachability.mjs` + allowlist | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT-PLAN.md` | RETIRE | execution plan only | accepted D0-R→D8-R results | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md` | RETIRE | early umbrella snapshot superseded by later D0-R/D1-R/D2-R... amendments | final accepted substantive NOTIF docs before canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-REFERENCE-STUDY.md` | RETIRE | research/reference-study evidence; decisions now accepted | D0/D1 canonical rehome + Git history | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md` | RETIRE | proposal/design superseded by accepted stage-specific decisions | D0→D8 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-FAMILY-SEMANTIC-CONTRACTS.md` | REHOME THEN RETIRE | 14 accepted family human meanings, birth transitions, audiences, deep links, repeat/suppression laws | `D1-DOMAINS-BOUNDARIES.md` compact notification-family boundary section | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D0-R-TRIGGER-SCOPE-CORRECTION.md` | REHOME THEN RETIRE | Personal Notifications Product 1.0 scope, 14 awareness families, attention-transition and exact-human audience laws | `D0-PRODUCT-SYSTEM-DEFINITION.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D1-R-PRODUCER-ROUTING-BOUNDARY-CORRECTION.md` | REHOME THEN RETIRE | Personal Notifications supporting owner, producer edges and routing boundary | `D1-DOMAINS-BOUNDARIES.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D1-R-SEMANTIC-CLOSURE.md` | REHOME THEN RETIRE | ratified ten-edge/14-family boundary, materialization audience, bounded Work and Sale→Post-Sale suppression | `D1-DOMAINS-BOUNDARIES.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md` | REHOME THEN RETIRE | AuthorizationRequest semantic boundary correction without new generic workflow owner | `D1-DOMAINS-BOUNDARIES.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R-IDENTITY-DATA-OWNERSHIP.md` | REHOME THEN RETIRE | Notification identity, Organization/recipient ownership, awareness lifecycle/source-correlation model | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R-RATIFICATION.md` | RETIRE | ratification wrapper | D2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R2-TEMPORAL-ROUTING-SUPERSESSION.md` | REHOME THEN RETIRE | routing changes affect future occurrences; historical recipients remain stable | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R2-RATIFICATION.md` | RETIRE | ratification wrapper | D2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-ROUTE-UNCONFIGURE-CUTOVER.md` | REHOME THEN RETIRE | typed configured/unconfigured route lifecycle and history semantics | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md` | RETIRE | ratification wrapper | D2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-PRESENTATION-SNAPSHOT.md` | REHOME THEN RETIRE | immutable/non-authoritative presentation snapshot law for historical human explanation | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md` | RETIRE | ratification wrapper | D2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-TYPED-RESULT-CONTINUATION.md` | REHOME THEN RETIRE | typed source/target continuation semantics and current-authorization fence | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` + D5 W2 where wire-specific | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md` | RETIRE | ratification wrapper | D2/W2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R6-AUTHORIZATION-REQUEST-IDENTITY-LIFECYCLE.md` | REHOME THEN RETIRE | AuthorizationRequest identity/lifecycle, request-target/review lineage and zero-decider semantics | `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R6-RATIFICATION.md` | RETIRE | ratification wrapper | D2 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-F1-TEMPORAL-ORDERING-FINDING.md` | RETIRE | finding resolved by accepted D3-R3 communication/recovery result | D3 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R-COMMUNICATION-PROPAGATION.md` | REHOME THEN RETIRE | producer→PersonalNotifications propagation, recovery/dedup/current reread semantics | `D3-COMMUNICATION-EVENTS.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md` | RETIRE | ratification wrapper | D3 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md` | REHOME THEN RETIRE | presentation snapshot feed-forward does not transfer current source authority | `D3-COMMUNICATION-EVENTS.md` + D2 current snapshot law | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md` | REHOME THEN RETIRE | typed result continuation feed-forward and current-owner reread/recovery | `D3-COMMUNICATION-EVENTS.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R3-AUTHORIZATION-REQUEST-COMMUNICATION-RECOVERY.md` | REHOME THEN RETIRE | recoverable AuthorizationRequest/decision/notification propagation and zero-decider Work behavior | `D3-COMMUNICATION-EVENTS.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R3-RATIFICATION.md` | RETIRE | ratification wrapper | D3 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-F1-OPERATION-SURFACE-ROUTE-REVERSIBILITY.md` | REHOME THEN RETIRE | accepted bounded route/operation surface correction | `D5-API.md` + `D5-B2-WIRE-CONTRACT.md` | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-F2-INBOX-PRESENTATION-CONTEXT-FINDING.md` | RETIRE | finding resolved by D2 presentation snapshot + D5/D6 final projection | D2/W2/P9 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md` | REHOME THEN RETIRE | final notification operation/Permission/client-class decisions | `D5-B2-OPERATION-ADMISSION-MATRIX.md` + W4 | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md` | REHOME THEN RETIRE | bounded recipient-candidate discovery surface and disclosure/authorization fences | W1 + W2 + W3 + W4 | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-OPERATION-ADMISSION-TABLE.md` | REHOME THEN RETIRE | admitted notification operations and exact access/client-class surface | operation matrix + W3 + W4 | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-RATIFICATION.md` | RETIRE | ratification wrapper | D5 canonical rehome | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md` | RETIRE | historical OAD proof snapshot | current Product OAD + `verify-notification-oad.mjs` + Product verifier | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R5-AUTHORIZATION-REQUEST-PRODUCT-SURFACE.md` | REHOME THEN RETIRE | final AuthorizationRequest read/decision/review-basis Product surface | D5 W1/W2/W3/W4 + operation matrix | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md` | RETIRE | historical OAD proof snapshot | current Product OAD + `verify-authorization-request-oad.mjs` + Product verifier | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md` | REHOME THEN RETIRE | Notifications shell/Inbox/settings route/surface and client-state laws | `D6-FRONTEND.md` + P5 + retained Notifications P8/P9 evidence | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md` | KEEP CURRENT EVIDENCE | operator LOCK for B00-R2/B11/B12 and B12 structural laws | active D6-R2 frontend acceptance evidence | keep |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-ACTIONABLE-GOVERNANCE-CONTEXT.md` | RETIRE | finding resolved by AuthorizationRequest Product surface + P5 supersession + final P9 | D5 canonical rehome + P5 + retained P9 | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md` | RETIRE | supersession ratification wrapper | P5/B110/P9 current evidence | rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md` | REHOME THEN RETIRE | accepted runtime/recovery realization constraints for AuthorizationRequest/Notifications | `D7-RUNTIME-JOBS-TRANSACTIONS.md` + accepted D7-C/D7-E owner as applicable | rehome+rm |
| `docs/engineering/rebaseline/D6-R2-NOTIF-01-D8-R-AUTHORIZATION-REQUEST-GOLDEN-FLOW-REVALIDATION.md` | REHOME THEN RETIRE | accepted D8 composed/golden-flow revalidation for Notifications/AuthorizationRequest | `D8-GOLDEN-FLOWS.md` | rehome+rm |

## 3. ADR / citation audit

| Path | Class | Surviving meaning | Replacement owner / live consumer | Action |
| --- | --- | --- | --- | --- |
| `docs/architecture/decisions/README.md` | KEEP AUTHORITY | sole ADR file-status/disposition registry | current ADR routing | keep |
| `docs/architecture/decisions/008-production-deploy-topology.md` | RETIRE | D7-conditioned residue; condition is now satisfied | D7-E / accepted deployment-operability authority | rm |
| `docs/architecture/decisions/010-mercado-livre-polling-visible-refresh.md` | RETIRE | D7-conditioned acquisition cadence residue | D7 runtime + D4 acquisition authority | rm |
| `docs/architecture/decisions/017-unknown-is-never-zero.md` | KEEP CURRENT EVIDENCE | retained Fact/domain-judgment evidence until its separate rehome condition closes | ADR 034 / D2 retained Fact condition | keep |
| `docs/architecture/decisions/018-mutation-envelope-table-and-poller.md` | RETIRE | D7-conditioned execution-safety residue | D7-C + accepted D3/D4 external-effect safety | rm |
| `docs/architecture/decisions/026-scheduler-phase-vocabulary.md` | RETIRE | D7-conditioned scheduler/runtime residue | D7 runtime/jobs authority | rm |
| `docs/architecture/decisions/030-scheduler-second-instance-per-installation.md` | RETIRE | D7-conditioned process/scheduler topology residue | D7 runtime/jobs authority | rm |
| `docs/architecture/decisions/034-fact-substitui-adr-017.md` | KEEP CURRENT EVIDENCE | retained Fact replacement/evidence anchor pending separate condition | D2 / ADR017 retirement gate | keep |
| `docs/architecture/decisions/035-architecture-rebaseline-governs-target-design.md` | KEEP AUTHORITY | D0–D9 transition/implementation-block authority | active until D0–D9 closes | keep |
| `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md` | KEEP CURRENT EVIDENCE | provenance still consumed by retained reconstructed ADRs | ADR017/034 while retained | keep |
| `docs/architecture/decisions/_citations/adr-009-citations.md` | RETIRE | last retained consumer is ADR010 | retire with ADR010 | rm |
| `docs/architecture/decisions/_citations/adr-013-citations.md` | RETIRE | last retained consumer is ADR018 | retire with ADR018 | rm |
| `docs/architecture/decisions/_citations/adr-017-citations.md` | KEEP CURRENT EVIDENCE | still consumed by ADR017/034 | ADR017/034 | keep |
| `docs/architecture/decisions/_citations/adr-07-twodigit-citations.md` | RETIRE | last retained consumer is ADR026 | retire with ADR026 | rm |
| `docs/architecture/decisions/_citations/adr-08-twodigit-citations.md` | RETIRE | last retained consumer is ADR030 | retire with ADR030 | rm |

## 4. Completed plan audit

| Path | Class | Surviving meaning | Replacement owner / live consumer | Action |
| --- | --- | --- | --- | --- |
| `docs/plans/2026-08-22-op-read-01-repair.md` | RETIRE | completed execution plan only | `D5-R2-OPERATIONAL-READ-PROJECTION-REPAIR.md` + current OAD/proof | rm |
| `docs/plans/2026-08-23-authorization-request-d5-wire.md` | RETIRE | completed execution plan only | D5 AuthorizationRequest rehome + current OAD/proof | rm |
| `docs/plans/2026-08-23-notif-d6-r-b00-r2.md` | RETIRE | completed frontend implementation/qualification plan | Notifications P8 ratification + locked HTML | rm |
| `docs/plans/2026-08-23-notif-d6-r-b11.md` | RETIRE | completed frontend implementation/qualification plan | Notifications P8 ratification + locked HTML | rm |

## 5. Verification-script audit

| Path | Class | Surviving meaning | Replacement owner / live consumer | Action |
| --- | --- | --- | --- | --- |
| `scripts/verify-product-oad.mjs` | KEEP CURRENT EVIDENCE | required current aggregate Product proof | `npm run gate` | keep |
| `scripts/verify-product-oad-baseline.mjs` | KEEP CURRENT EVIDENCE | current verifier dependency carrying still-valid baseline assertions | `verify-product-oad.mjs` replay chain | keep |
| `scripts/verify-product-oad-current99.mjs` | KEEP CURRENT EVIDENCE | current verifier dependency | `verify-product-oad.mjs` replay chain | keep |
| `scripts/verify-product-oad-pre-auth.mjs` | KEEP CURRENT EVIDENCE | current verifier dependency | `verify-product-oad.mjs` replay chain | keep |
| `scripts/verify-notification-oad.mjs` | KEEP CURRENT EVIDENCE | current 106/31 notification semantic proof | targeted Product proof | keep |
| `scripts/verify-authorization-request-oad.mjs` | KEEP CURRENT EVIDENCE | current 106/31 AuthorizationRequest semantic proof | targeted Product proof | keep |
| `scripts/verify-oad-source-reachability.mjs` | KEEP CURRENT EVIDENCE | current OAD source-reachability/orphan policy proof | targeted Product proof | keep |
| `scripts/verify-operational-read-contract.mjs` | KEEP CURRENT EVIDENCE | current owner-local operational-read proof | D5-R2 / targeted proof | keep |
| `scripts/verify-performance-evidence-knowledge.mjs` | KEEP CURRENT EVIDENCE | current performance evidence/knowledge proof | D6-R1 / targeted proof | keep |
| `scripts/verify-d6-r-b00-r2-wireframe.mjs` | KEEP CURRENT EVIDENCE | reproducible proof for currently locked B00-R2 HTML | Notifications P8 evidence | keep |
| `scripts/verify-d6-r-b11-inbox-wireframe.mjs` | KEEP CURRENT EVIDENCE | reproducible proof for currently locked B11 HTML | Notifications P8 evidence | keep |
| `scripts/verify-d6-r-b110-approvals-wireframe.mjs` | KEEP CURRENT EVIDENCE | reproducible proof for currently locked B110 HTML | B110 acceptance evidence | keep |
| `scripts/verify-d6-r-b12-routing-settings-wireframe.mjs` | KEEP CURRENT EVIDENCE | reproducible proof for currently locked B12 HTML | Notifications P8 evidence | keep |
| `scripts/verify-d6-r-p9-screen-contracts.mjs` | RETIRE | historical P9 meta-proof tied to an intermediate Fable repair document; not required by current gate | retained final P9 contract + W1 canonical rehome + current OAD proofs | rm |
| `scripts/verify-authorization-request-w1-carrier.mjs` | RETIRE | repair-specific W1 proof; accepted carrier meaning must move to canonical W1 | W1 rehome + current AuthorizationRequest OAD proof | rm |
| `scripts/verify-authorization-request-fable-fixes.mjs` | RETIRE | repair/reviewer-specific proof after final corrections | canonical rehome + `verify-authorization-request-oad.mjs` | rm |
| `scripts/verify-authorization-request-d7r.mjs` | RETIRE | stage-repair proof after accepted D7-R result | D7 canonical rehome + current Product/AuthorizationRequest proofs | rm |
| `scripts/verify-authorization-request-d8r.mjs` | RETIRE | stage-repair proof after accepted D8-R result | D8 canonical rehome + current Product/AuthorizationRequest proofs | rm |

## 6. Current owners that must receive consolidation but are not retirement candidates

These files stay in the active tree and must be updated during Task 3 because the audit proved they are stale/incomplete relative to already accepted later work:

| Current owner | Required consolidation |
| --- | --- |
| `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md` | accepted Personal Notifications Product scope and attention-only boundary |
| `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md` | Personal Notifications supporting owner, producer edges, family/audience boundary, AuthorizationRequest boundary correction |
| `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` | Notification identity/data/routing/snapshot/continuation and AuthorizationRequest identity/lifecycle |
| `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md` | notification/AuthorizationRequest propagation, recovery, reread and zero-decider semantics |
| `docs/engineering/rebaseline/D5-API.md` | current 106/31 post-amendment API closure status/laws |
| `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` | final Notification/AuthorizationRequest routes and typed request-ETag carrier law |
| `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md` | current 106-operation admitted surface |
| `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` | current Notification/AuthorizationRequest schema meanings |
| `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` | current notification/actionable-request collections/query semantics |
| `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` | current 31-Permission / 106-operation access surface, including `notifications.manage` and authenticated-self reads |
| `docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md` | source-reachability/orphan hygiene rule already mechanically enforced |
| `docs/engineering/rebaseline/D6-FRONTEND.md` | accepted Notifications/Approvals frontend architecture and authority boundaries |
| `docs/engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md` | current Notifications/Approvals surface inventory and current 106-operation basis |
| `docs/engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md` + exact D7-C/D7-E owner | accepted AuthorizationRequest/notification runtime/recovery corrections |
| `docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md` | accepted AuthorizationRequest/notification composed revalidation |

## 7. Task 1 verdict

**PASS — safe to proceed to rehome/retirement execution, with no deletions performed in Task 1.**

The audit falsified one unsafe cleanup shortcut: the NOTIF-01 chain cannot be treated as pure history yet because its accepted semantics were never fully consolidated into D0–D8/W1–W4/P5. Task 3 must perform that consolidation before the corresponding `rehome+rm` paths can leave the active tree.
