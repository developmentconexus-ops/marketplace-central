# M-09 Execution Guide

## Objective

Replace the MVP catalog's ambiguous legacy product identity with positive Sankhya
CODPROD, preserve separate nullable external identifiers, cut the active catalog
read path to Oracle/internal_read, and keep unknown facts honest.

## Required Initial Reads

- `AGENTS.md`
- `.agents/skills/mpc-goal-harness/SKILL.md`
- mission `mission.md`
- this milestone's `milestone.md` and `validation-contract.md`
- `research/mnos-sankhya-read-interface-contract.md` (IC-002)
- `research/mvp-operator-workspace-interface-contract.md` (IC-003), only ProductRef,
  SourceFact, identity/compatibility, transport, and must-preserve sections
- the current feature brief before each dispatch
- compiled selectors from knowledge route `portfolio-core`

Paths above are relative to the repository or mission root as evident from context.
Do not load transcripts or broad historical M-06 evidence.

## Feature Order

1. F-01 canonical product identity contract.
2. F-02 Oracle catalog cutover.
3. F-03 catalog compatibility and legacy-runtime cutover.

The Milestone may combine adjacent briefs only when one coherent writer/commit still
has a reviewable boundary. Otherwise dispatch sequentially. Use only the named
custom agent `mpc-implementer`; its project profile fixes Luna/high.

## Ownership And Seams

- One writer at a time in the shared checkout.
- F-01 owns catalog/internal_read identity types and the atomic OpenAPI/SDK seam.
- F-02 owns catalog application/ports plus Oracle/internal_read adapter composition.
- F-03 owns deterministic migration/compatibility, active MSDB removal, and the one
  migration-number seam if a forward migration is required.
- Feature plans must name exact allowed/forbidden paths before editing.
- M-06 evidence, web workspace redesign, provider mutation, auth/RBAC, and broad
  runtime consolidation are forbidden.

## Proof And Commit Rule

Each feature writes `validation.md`, runs targeted proof with absolute repository
`.gocache`, and returns one intentional commit. Public API shape changes update
OpenAPI and `packages/sdk-runtime` in the same commit. Mocks prove contracts only.

After integration, freeze the SHA and request `mpc-verifier` in
`fixed_sha_review` mode, then `proportional_qa` mode. The verifier project profile
fixes Luna/high. QA follows `validation-contract.md`; only QA may pass M-09.

## Stop Conditions

- Legacy identity cannot map deterministically to CODPROD.
- Unknown cost/price/stock would need a zero/default.
- Oracle/provider payload or adapter type would cross into domain/application.
- Tenant scope, OpenAPI/SDK lockstep, shared-seam ownership, or read-only Oracle
  guarantees conflict with the planned change.
- Required live Oracle read is unavailable: terminal `externally_blocked`, not Pass.

On terminal, persist the compact checkpoint and send its path/verdict to the exact
Portfolio `hub_task_id` from the `/goal` packet.
