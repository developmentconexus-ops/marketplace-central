# Marketplace Central — Rebaseline Evidence Register

> **Status:** SUPPORTING EVIDENCE  
> **Baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> Evidence informs D1–D9; it does not become target architecture by itself.

## Current repository facts

- `internal/modules/` has 21 legacy module directories.
- `internal/contexts/` has 2 new contexts: `catalog` and `listings`.
- `internal/adapters/` already has `erp` and `marketplace` families.
- `cmd/` has 7 entrypoint directories.
- Frontend ownership is split between app-local routes/pages and feature packages.
- Legacy route redirects remain in `AppRouter.tsx` with no production-user compatibility requirement.
- OpenAPI, a hand-written SDK and routing knowledge currently coexist.
- The migration chain contains both legacy and new-context state.
- `scripts/gate.ps1` and `scripts/harness.ps1` are active verification/runtime-test mechanisms.
- Current governance rules still describe parts of the legacy tree; they are changed with the D-stage that owns each code surface, not during documentation cleanup.

## Operator constraints

- No production users require backward compatibility with the current application.
- Hard cutover is allowed when target design requires it.
- Git history is the archive; no active `old/` source tree is desired.
- Technical design must answer context, identity, database, communication, events, external integrations, API, frontend and runtime questions before implementation planning.
- YAGNI removes accidental complexity, not correctness.

## Product/source direction

Historical material consistently describes MPC as an internal marketplace operations/intelligence cockpit:

- Mercado Livre first;
- Sankhya/Oracle supplies internal business facts;
- MPC adds linkage, stock reconciliation, pricing/profitability, order/reconciliation workflows and safe operational controls;
- external protocol details belong behind adapters/ports;
- unknown operational facts must not silently become plausible defaults.

D1–D4 decide the exact boundaries and reverify external/source semantics.

## Domain evidence to carry, not blindly ratify

### Product/listing relationship

Historical operator/domain work indicates the business may require many-to-many product↔listing relationships, variation-level identity and separate kit/BOM semantics. D1/D2 must validate the actual model; a convenient 1:1 current table is not proof.

### Product identity conflict

Old drafts treated `CODPROD` as internal ProductID. Later design work proposed opaque MPC identity plus source keys. These conclusions conflict. D2 must explicitly decide internal identity vs `CODPROD`, EAN, REFFORN and source-instance identifiers.

### Provider identity

Historical code/research distinguishes listing item ID, variation ID, catalog-product relation and seller SKU/custom-field values. D2/D4 define which are provider references versus internal identities and reverify current provider behavior.

### Pricing source mismatch

Historical review found a real class of defect where ERP taxonomy/category could be confused with Mercado Livre listing category when calculating commission/margin. D1/D4 must assign ownership and inputs; D8 must prove the golden flow uses the right facts.

## Sankhya evidence to reverify in D4

A prior read-only pass in the Metal Nobre environment recorded these candidates:

- product code: `TGFPRO.CODPROD`;
- EAN observed at `TGFPRO.REFERENCIA` in that environment;
- supplier reference: `TGFPRO.REFFORN`;
- cost candidate: `TGFCUS.CUSSEMICM` with product/company/as-of semantics;
- tax rule inputs associated with `TGFICM` and realized item tax facts in `TGFITE`.

These are historical measurements, not permanent contracts. D4 rechecks current schema/query semantics before target ports/contracts are ratified.

## Durable uncertainty/safety properties

Evidence repeatedly supports:

- unknown is not zero/default;
- absence from a partial source pull is not automatically terminal state;
- provider wire payload is not domain truth merely because fields are useful;
- external identifier is not automatically canonical identity;
- blind retry of a potentially accepted write is unsafe;
- provider success response is not automatically convergence;
- raw observation/reprocessing can be useful where mappings evolve, but must be justified per capability rather than universalized.

Exact mechanisms belong to D2–D4/D7.

## Verification lessons already absorbed into the canonical method

- presence is not execution;
- zero executed checks is not proof;
- a measurement needs a stated universe;
- current directory/schema/ADR shape is not sufficient justification for target structure;
- negative fixtures/counterexamples are valuable proof that a control can fail;
- stale binaries/worktrees can invalidate evidence;
- converging authorities is preferable to hand-syncing duplicate authorities with more guards.

Full historical reports remain in Git history.

## Open evidence questions

- **D1:** final contexts and dissolution of legacy modules.
- **D2:** identity model, table ownership, exact value/knowledge semantics, recoverability/reset.
- **D3:** sync/event/projection map and event/outbox semantics.
- **D4:** current Mercado Livre and Sankhya capability contracts.
- **D5:** one API contract authority and generation/runtime validation.
- **D6:** frontend feature/package and data-consumption topology.
- **D7:** process, scheduler, transaction, outbox and deployment topology.
- **D8:** end-to-end golden-flow proof.

These are design gates, not tasks for an implementation worker to decide locally.