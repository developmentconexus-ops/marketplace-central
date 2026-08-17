# D3 — Communication / Events

> **Status:** OPEN / IN PROGRESS — D3 stage opened; no D3 batch is canonical yet; D3-B1 direction is operator-approved for independent review  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-16

## 1. Purpose and boundary

D3 answers:

> **How do the owners and boundaries already accepted in D1/D2 communicate without transferring, duplicating or hiding authority?**

D3 decides, proportionately:

- when an accepted semantic dependency requires a synchronous owner query/capability interaction;
- when a committed producer-owned fact justifies an event;
- when multi-authority composition belongs in a projection/read model;
- how accepted D1 workflow cycles can be realized without private-code or shared-write-authority cycles;
- the minimum event/communication semantics needed for delivery, ordering, duplication, replay and recovery contracts;
- how Organization scope, Principal/actor attribution, domain-local Business Intent identity, provenance and material historical correctness survive communication.

D3 does **not** choose:

- Mercado Livre/Sankhya/payment concrete contracts, DTOs, auth, pagination, provider webhook/polling semantics or source completeness — **D4**;
- HTTP/OpenAPI endpoint/operation/error contracts — **D5**;
- frontend screens, work inboxes, approval UX or concrete UI projection topology — **D6**;
- workers, queues, brokers, outbox implementation, transaction boundaries, locks, retries, RLS enforcement, deployment/process topology or execution-safety runtime realization — **D7**;
- product implementation.

Implementation remains blocked until D9 is accepted.

## 2. Imported parent invariants

These are already authoritative from D0–D2; D3 does not reopen them by choosing communication:

1. **Every material business meaning has one semantic authority.** Communication may transport or reference meaning; it does not create a second write authority.
2. **D3 may communicate only inside the accepted D1 semantic edge set.** If a necessary semantic dependency is missing, stop and reopen only the implicated D1 decision instead of hiding the dependency in an event, API, queue, database or projection.
3. **A consumer uses producer-owned meaning only through an explicit public semantic boundary owned by the producer.** No private implementation import, cross-context SQL/private-table access or shared mutable business entity.
4. **Organization scope is explicit.** Communication must not reconstruct tenant scope from provider account, Installation, Selling Entity, source key or process-global defaults.
5. **Actor attribution remains semantically honest.** Human, delegated automation and system action remain distinguishable where material.
6. **Material domain-owned Business Intents retain domain-local identity.** D3 must not introduce a generic Action/Mutation/Command/BusinessIntent owner merely to route communication.
7. **Historical snapshots/references/projections never become current producer authority.** A read model may compose authorities; it may not become write authority.
8. **External provider/business-system facts remain external/source-qualified.** Provider notifications or payloads are not automatically MPC domain events; D4 owns concrete acquisition/translation contracts.
9. **Unknown/absence, provenance and material time meanings remain fail-honest.** Communication must not turn stale, partial, duplicated or delayed information into plausible current truth.
10. **Mechanism ≠ Authority.** Shared transport/runtime mechanics may later centralize accidental complexity without acquiring business meaning.

## 3. D3 decision surface

D3 must close two material decision groups before final stage review.

### D3-B1 — Communication Topology & Edge Matrix

Determine the smallest correct semantic communication form for every accepted D1 edge and the D2 cross-cutting identity/access dependency:

- current owner query;
- explicit owner capability request;
- committed domain event;
- projection/read model;
- combinations such as event-triggered reaction followed by current owner revalidation when current truth is material.

B1 must explicitly resolve the highest-risk cycles/edges:

- Business-System Materialization ⇄ Fulfillment Lifecycle;
- Controlled Action Governance ⇄ action-owning domains;
- Operational Work ⇄ originating domains;
- Post-Sale Resolution ⇄ Sales / Materialization / Fulfillment / Economics;
- Commercial Economics ⇄ Marketplace Offering Operations.

### D3-B2 — Communication Contract & Failure Semantics

After B1 identifies which communications actually exist, define only the semantic contract properties required before D4/D7, including as materially necessary:

- committed-fact/event identity and producer ownership;
- Organization scope;
- Principal/actor attribution;
- domain-local Intent reference;
- causation/correlation;
- provenance and material time meanings;
- duplication and ordering assumptions;
- replay semantics;
- delayed/missed-delivery recovery expectations;
- projection rebuildability and authority fence.

D7 will decide concrete transport/runtime implementation.

## 4. Review and authority protocol

D3 uses the same accelerated protocol proven in D2:

1. GPT prepares a coherent candidate batch from repository authority/evidence.
2. The operator approves the **candidate direction** for independent challenge.
3. A disposable `D3-B<n>-REVIEW-CANDIDATE.md` may be committed, clearly marked **NOT ARCHITECTURE AUTHORITY**.
4. The operator invokes **Fable** separately. Fable reconstructs the authority path independently and writes material findings through the GitHub review channel (`AI-DIALOG.md`).
5. Reviewer findings are evidence, never authority. GPT independently adjudicates each material finding against the current repository authority/evidence.
6. Reviewer disagreement is named and returned for another round or escalated to the operator; GPT does not simulate Fable and Fable does not silently decide operator authority.
7. Only the operator-approved converged batch is consolidated into this canonical D3 artifact.
8. Review candidates are disposable and should be removed after consolidation.
9. After material batches converge, D3 receives a final **Global Coherence + YAGNI / Overengineering / Future-Cost review**.
10. D3 closes only after explicit operator ratification as a whole.

`AI-DIALOG.md`, review candidates, chat summaries and reviewer statements are **not** part of the architecture authority path.

## 5. Current D3 state / exact next action

D3 is **OPEN / IN PROGRESS**.

The operator approved the direction of **D3-B1 — Communication Topology & Edge Matrix** for independent review. That approval authorizes review preparation; it does **not** make the B1 candidate canonical D3 architecture.

Exact next action:

1. review `docs/engineering/rebaseline/D3-B1-REVIEW-CANDIDATE.md` through an independent Fable round;
2. GPT independently adjudicates the returned material findings against repository authority/evidence;
3. continue reviewer rounds only where a material dispute remains;
4. return the converged B1 batch to the operator for explicit batch acceptance;
5. only then consolidate accepted B1 meaning into this artifact and proceed to B2.

Do not begin D4 or product implementation while D3 remains open.