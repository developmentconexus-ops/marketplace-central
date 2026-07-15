# Mission Validation Contract

```yaml
id: MIS-003
type: mission-validation-contract
status: draft
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: mission
```

## Mission ID

MIS-003-operator-cockpit-wireframe-replan

## QA Level

QA-2 — integration-proven contracts (ephemeral-postgres + stub provider adapter); live-provider-read lane where a criterion names it; live provider WRITE only under governed lane with explicit operator authorization (RK-03). Dual gate per milestone: Codex `mpc-verifier` + independent Claude review at fixed SHA; only QA passes.

## Required Final State

Wireframe deck-2 cockpit operational: 8 sidebar workspaces + protocolo detail live on the M-02 platform; listings read spine (IC-02), mutation envelope (IC-03), market contract-only module (IC-04), frontend platform (IC-05) all proven per their milestone contracts; legacy pages replaced with redirects; MIS-001 M-13/M-14/M-07 marked superseded.

## Criteria

## Criterion: All six milestones passed by QA
ID: MIS-003-C01
Level: Mission
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: `Get-ChildItem .mnfs/MIS-003-operator-cockpit-wireframe-replan/M-0*/validation-result.md`
- Expected: six files, each with verdict `Pass` issued by QA Validator at a named fixed SHA; each milestone's dual-gate review (Codex mpc-verifier + Claude) recorded in its rollup
- Actual:
- Artifact: `<milestone-root>/validation-result.md` × 6
Blocking failure: any milestone verdict ≠ Pass, or a rollup missing dual-gate records
Blocking failure observed: No
Owner: QA Validator

## Criterion: Wireframe route map live with reload fidelity
ID: MIS-003-C02
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser walkthrough (dev build) over the IC-05 route table + all 6 legacy redirects
- Expected: each IC-05 route renders its workspace; `/anuncios?installation=X&tab=pendencia&filter.exception=sync_error` after F5 shows identical tab/filter/installation; each legacy URL lands on its new route with query preserved
- Actual:
- Artifact: `validation-result.md` screenshot set at mission root evidence section
Blocking failure: any route 404s, any redirect drops query, any deep link loses state on reload
Blocking failure observed: No
Owner: QA Validator

## Criterion: Unknown-never-zero end to end
ID: MIS-003-C03
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: seeded walkthrough: listing without cost on `/anuncios` + `/precos`; dashboard with one degraded source on `/`; virgin market data on `/precos` competitive columns
- Expected: margem cells show "—" with hint "sem custo no ERP → não simulado"; degraded counter card shows "—" with tooltip "fonte indisponível: {source}"; competitive columns show `no_price_evidence` copy; JSON evidence shows `null` (not 0) at each source endpoint
- Actual:
- Artifact: mission `validation-result.md` unknown-audit section (JSON excerpts + screenshots)
Blocking failure: any unknown rendered as 0, R$ 0,00, or fabricated value
Blocking failure observed: No
Owner: QA Validator

## Criterion: Envelope governs every provider write
ID: MIS-003-C04
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `Grep` provider write capability call sites in `apps/server_core/internal/modules/`; review of write paths from UI surfaces
- Expected: all provider write invocations reachable only from the mutations poller apply path; UI mutation surfaces call only `/mutations*` endpoints; StockActionService entry points produce envelope protocols
- Actual:
- Artifact: grep transcript + call-graph note in mission `validation-result.md`
Blocking failure: any provider write path bypassing the envelope (except pre-existing non-superseded flows explicitly listed in M-03 rollup)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Governance lanes green at final SHA
ID: MIS-003-C05
Level: Mission
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: repo governance lane suite (incl. `scripts/run-live-oracle-docker.ps1` read-only lanes) at the mission-final SHA
- Expected: all lanes green; `GOV_API_SDK_SPLIT` check passes (OpenAPI and sdk-runtime consistent); Go tests green with `GOCACHE` absolute
- Actual:
- Artifact: lane output logs referenced from mission `validation-result.md`
Blocking failure: any red lane at final SHA
Blocking failure observed: No
Owner: QA Validator

## Criterion: Supersession recorded
ID: MIS-003-C06
Level: Mission
Type: Documentation
Required: Yes
Status: Pending
Evidence:
- Command: `Grep -n "superseded_by: MIS-003" .mnfs/MIS-001*/M-13*/milestone.md, M-14*, M-07*`
- Expected: three matches, one per superseded milestone frontmatter; no other MIS-001 artifact modified (git diff scope proof)
- Actual:
- Artifact: grep output + git diff excerpt in mission `validation-result.md`
Blocking failure: missing note or out-of-scope MIS-001 edits
Blocking failure observed: No
Owner: QA Validator

## Criterion: Market module honestly empty
ID: MIS-003-C07
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: `GET /market/observations?installation_id=<inst>&listing_ids=<any>` (IC-04 params) on production-shaped DB; grep for collector implementations and scraping dependencies
- Expected: 200 with all items `evidence_state: "no_price_evidence"`; zero production CollectorPort implementations; zero scraping libraries in go.mod; no market seed rows outside `_test`
- Actual:
- Artifact: response JSON + grep transcript in mission `validation-result.md`
Blocking failure: any production collector, scraped data, or fabricated market value (ML ToS 7.6 / G1 posture)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Milestone rollups: `<milestone-root>/validation-result.md` (six).
- Mission rollup: `<mission-root>/validation-result.md` with fixed SHA, lane logs, screenshot set, JSON excerpts.
- Feature quick proofs: `<feature-root>/validation.md` (referenced, not re-validated here).

## Blocking Failures

Any criterion's named blocking failure; additionally: secret/token/PII found in any evidence artifact (Q2) blocks the mission regardless of criterion status.

## Retry Policy

Milestone-level corrections per milestone contracts (max 2 correction attempts each). Mission criteria re-checked only after affected milestone re-passes.

## Handoff

- Current status: draft (planning).
- Next owner: QA Validator (after all milestones integrate).
- Next action: none until M-01..M-06 rollups exist.
- Required files/evidence: as listed under Evidence Requirements.
- Blockers or open decisions: live provider-write lane authorization decided by operator at M-03 QA time (RK-03).
