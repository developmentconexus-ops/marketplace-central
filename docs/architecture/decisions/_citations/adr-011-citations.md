# ADR-11 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 36
**Spellings found:** `ADR-11` (MIS-004/MIS-007) and `ADR-011` (MIS-001/MIS-003).

## Assertion A1 — Webhook body is a pointer, never becomes domain data; handler always returns 200; dedupe key was amended (narrowed to `_id`-only) from its original full tuple
MIS-007's dominant, most-cited meaning, including a recorded amendment.
- Citations: 22
- Verbatim: "Webhook: ponteiro, nunca dado; always-200." / dedupe-amendment verbatim: "dedupe de inbox só com `_id` — ADR-11 emendado; ver M-08 Done Means"
- Anchors:
  - `.mnfs/MIS-007-ml-sync/mission.md:224`
  - `.mnfs/MIS-007-ml-sync/mission.md:373`
  - `.mnfs/MIS-007-ml-sync/mission.md:402`
  - `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:25`
  - `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:104`
  - `.mnfs/MIS-007-ml-sync/M-08-webhook-ingest/validation-contract.md:54`
  - `.mnfs/MIS-007-ml-sync/M-08-webhook-ingest/validation-contract.md:94`
  - `.mnfs/MIS-007-ml-sync/M-08-webhook-ingest/milestone.md:78`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r02.md:99`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-reconciliation-r02.md:53`

## Assertion A2 — (MIS-004 mission, unrelated subject) Market-reference collection is on-demand only; runtime is a local docker dev stack, with no retroactive/historical backfill claimed
- Citations: 5
- Verbatim: "ADR-11 Coleta on-demand + runtime docker local"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:95`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:81`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:97`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:18`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:19`

## Assertion A3 — (MIS-001/MIS-003 lineage) Historical gate preservation: the M-06 failed QA gate stays fixed at its SHA and is never rewritten to manufacture a pass
- Citations: 4
- Verbatim: "ADR-011 Historical gate preservation ... Fixed SHA and failed QA remain intact"
- Anchors:
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md:199`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md:277`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:262`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/market-intelligence-digest.md:35`

## Contradictions
- **Number collision across missions, three ways.** ADR-11 means webhook-pointer-never-data (MIS-007), on-demand-collection/local-docker-runtime (MIS-004), and M-06-gate-frozen (MIS-001/MIS-003) — three unrelated rules under one label, none referencing the others.
- **Amendment recorded but reasserted verbatim elsewhere.** `p5-claude-decomposition-audit-r03.md:66` and `:68` show the dedupe tuple narrowing (A1) was explicitly negotiated and logged as an amendment ("emenda registrada em `mission.md` ADR-11, P5 r02 N-4"), yet several later citations (e.g. `p7-seat1-star1-star5-r02.md:27`) still cite "ADR-11" bare, without the "emendado" qualifier, risking readers picking up the pre-amendment (full-tuple) rule.
- **Open unaddressed gap flagged against the amendment.** `p7-seat5-doublepass-r02.md:55` notes the narrowed dedupe key means "flood without `_id` bypasses dedupe ... grows `notifications_inbox` unboundedly; no retention/purge policy is declared anywhere" — a live risk the amendment introduced but never closed.

## Exceptions / carve-outs
- The dedupe-tuple narrowing itself is the carve-out: original ratified tuple was `(provider, topic, resource, notification_id)`; amended to require only `_id`/`notification_id` presence, with the no-`_id` path explicitly declared as *not* covered by transport dedupe — only by `IngestOrder` domain idempotence (`M-08-webhook-ingest/milestone.md:77-80`).
