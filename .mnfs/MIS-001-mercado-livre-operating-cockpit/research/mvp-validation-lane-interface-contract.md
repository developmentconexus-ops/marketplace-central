# MVP Validation Lane Interface Contract

```yaml
id: IC-004
type: interface-contract
status: ready
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Proportional Windows PowerShell validation for the MIS-001 MVP: deterministic
tests, real read-only sources, PostgreSQL state, safe evidence, and one browser
journey at a reviewed SHA.

## Principles

- QA records the reviewed SHA before running evidence and stops if HEAD changes.
- Use existing repository/harness commands when they already prove the criterion.
- Go commands run from `apps/server_core` with repository-local absolute `GOCACHE`.
- Real credentials and selector values remain local; evidence contains no secret,
  buyer PII, raw provider payload, or full Oracle row.
- Oracle and Mercado Livre operations are read-only. Any provider mutation stops QA.
- Mocks prove deterministic behavior only; they never prove a real integration.
- Missing real data is `externally_blocked`, not a fixture-backed pass.

## Local Scenario

The operator selects one installation, positive CODPROD, Mercado Livre listing,
related order, and nullable variation using a Git-ignored local file or process-local
environment. The values are never copied into MNFS evidence. A concise evidence note
records source, observed time, read-only status, command or browser action, and the
artifact path. Hash manifests and a custom evidence schema are not required.

## Validation Ladder

| ID | Purpose | Command or interaction | Evidence |
| --- | --- | --- | --- |
| MVP-QA-01 | Server behavior | From `apps/server_core`, set absolute `.gocache` and run targeted packages; run `go test ./... -count=1` when the milestone impact warrants it | command output + SHA |
| MVP-QA-02 | SDK/web behavior | From repository root run impacted SDK/web tests and web build | command output + SHA |
| MVP-QA-03 | PostgreSQL | Run the existing integration lane; repeat the selected listing/order import and verify stable natural-key counts | concise output, no connection string |
| MVP-QA-04 | Oracle read | Run the governed live Oracle read/smoke for a positive CODPROD | source, timestamp, `read_only=true`, sanitized result |
| MVP-QA-05 | Mercado Livre read | Run the governed live listing/order GET lane | method/path class, timestamp, `read_only=true`, sanitized result |
| MVP-QA-06 | Browser journey | Complete MVP-BD-01 at desktop and 390x844 | interaction note + screenshots |
| MVP-QA-07 | Safety fold | Inspect evidence and reachable MVP UI/network behavior for secrets/PII and provider mutation | explicit pass/fail note |

Exact commands may be supplied by the applicable milestone contract or existing
harness registry. QA records the commands actually run; it does not need to build a
new wrapper merely to serialize them.

## Browser Drive MVP-BD-01

1. Open `/` and confirm Visão geral, Produtos, Anúncios, Vendas, Operações, and the
   selected installation.
2. Open a stock attention item and verify the exact listing plus source, quality,
   and observation time.
3. Open Product 360 and verify positive CODPROD is distinct from EAN, reference,
   seller SKU, and provider identity.
4. Return to the listing, open the related sale, and verify revenue, nullable margin
   inputs, quality, source, and observation time.
5. Return to the listing and review the stock simulation. Confirm current/proposed
   values, reason, preview payload, `executed=false`, and no execute control.
6. Reload and confirm installation/entity context is retained. Capture desktop and
   390x844 evidence without buyer PII.
7. Exercise at least one applicable empty/error/unknown/stale/conflict state and
   confirm unknown business values are not rendered as zero.

## Evidence Paths

QA may write under the milestone `_fixed-sha-qa/` directory:

- `commands.md` or sanitized command outputs;
- `real-reads.md`;
- `postgres.md`;
- `browser-drive.md` and named screenshots;
- `safety.md`;
- `validation-result.md` at the milestone root.

The final files name the reviewed SHA. Temporary capture mechanics are an
implementation detail, not a product acceptance criterion.

## Environment Ownership

Environment key names remain governed by
`contracts/governance/execution-lanes.json` and
`research/m08-governance-runtime-inventory.md`. Values never enter evidence.
No new alias, dependency install, cache purge, cold clone, or WSL lane is allowed.

## Terminal Outcomes

- `passed` and `failed` are QA verdicts.
- `externally_blocked` means a credential, real source, or coherent real sample is
  unavailable. Persist the compact harness checkpoint and callback Portfolio.
- QA alone writes the final milestone verdict.

## Handoff

- Current status: Ready for proportional use by M-14.
- Next owner: M-14 Milestone Orchestrator, then `mpc-verifier` for QA.
- Required proof: ladder rows applicable to the claims plus MVP-BD-01.
- Blocking failures: provider/Oracle write, secret/PII leak, fake-as-real evidence,
  SHA drift, or an unknown numeric represented as zero.
