# P5 Reconciliation — round 03 (MIS-007-ml-sync)

```yaml
type: planning-review-reconciliation
phase: P5
round: 03
created: 2026-07-31
auditor_artifact: p5-claude-decomposition-audit-r03.md
auditor: cold Claude Opus crew (task a37f6975f3368ccfc; operator-ratified waiver — Sol P5
  retroactive mandatory before status: planned)
input_manifest: p5-input-r03.sha256 (digest 94faf38b37e907073434b423fd331faa6828b001ef6a3bcc86c8b7ef0d1e7407)
verdict_received: NEEDS-REVISION (N-1..N-10 all CONFIRMED CLOSED; both r02 arm-adaptations
  judged SOUND under adversarial reading; 2 new blocking P-1..P-2, 3 advisory P-3..P-5)
disposition: ALL findings ACCEPTED as valid; fold applied in full; re-audit r04 required
```

## Process note

Task output-file came back 0 bytes — THIRD instance of the same transport fragility.
Compaction fired between the notification's arrival and persistence; the verbatim was
recovered from the session transcript JSONL (line 1164, task-notification block for
a37f6975f3368ccfc), HTML entities decoded, and persisted with a provenance header as the
FIRST action of the resumed turn. The same-turn persistence rule was honored in spirit
(no fold work preceded persistence); the compaction gap is recorded here as a deviation.

## Per-finding disposition (all ACCEPTED, no downgrades)

### Blocking

- **P-1 (phase mandatory over ALL cursors contradicts ratified ADR-07, the live products
  job, and M-09's tolerant-parse pin)** — FIXED via the auditor's FIRST arm (narrow to the
  ratified scope): IC-06 Compatibility Rules now scope `phase` obrigatório to the cursors
  of jobs NOVOS (M-04/M-06 — ADR-07's ratified wording), name `products_job.go:22-25`
  (`ProductsCursor`, no phase field) as INTOCADO, and pin the scheduler's parse as
  TOLERANTE — phase ausente/desconhecida ⇒ `incremental=false` (today's behavior), never
  an error, verbatim-matching M-09 F-01's constraint. M-02 F-03 Brief (items 1 and 2),
  Expected Output, and Negative Scenarios all narrowed to the same scope; the contract
  test's domain is stated as jobs de M-04/M-06; a NEW negative scenario pins the legacy
  cursor path (products sem phase → incremental=false, SEM erro). The products regression
  is now falsifiable against the REAL job: F-03 Validation Expectations and M-02
  milestone Done-Means both assert `incremental=false` recorded for the real
  `ProductsCursor` (no longer provable only via "job fake"). `products_job.go` remains in
  Forbidden paths — arm (b) rejected because editing the only live job was exactly what
  F-03's ownership block and r02 N-1's closure forbid; arm (a) removes the contradiction
  without touching any concrete job. Lane-A totality conflict resolved: M-02 F-03 and
  M-09 F-01 now state the SAME tolerant semantics for the same field.
- **P-2 (ADR-14 "≤1 milestone FE-contract em voo" contradicts ratified lane C
  M-05∥M-06∥M-07)** — FIXED via the FIRST arm (amend ADR-14 to the rule actually
  planned): mission.md ADR-14 now reads COMMITS de contrato FE serializados — ≤1 COMMIT
  de contrato (YAML+SDK) em voo por vez, hub arbitra a ordem, código paraleliza — with a
  recorded EMENDA noting the original form contradicted the ratified lane C and naming
  the hand-written client literal as a hub-arbitrated seam (same form as M-07's pricing
  region of root.go). The transversal edge line and the matrix OpenAPI/SDK rule restate
  the amended form; R-7's trigger re-worded to "2 COMMITS de contrato simultâneos ou não
  arbitrados pelo hub" (no longer fires on the ratified plan by construction). The three
  IC restatements aligned: IC-03 Route Namespace, IC-07, IC-05 (M-09's contract commit
  enters the hub serialization). Lane-C milestone files already carried the
  commit-serialization phrasing (M-05:61-62, M-06:61-62, M-07:61-62) — the amendment
  makes the ADR match them instead of contradicting them. Arm (b) (serialize lane C)
  rejected: it would redraw the operator-ratified lane structure to preserve a drifted
  ADR sentence; the ratified artifact is the lane plan, the ADR sentence is the defect.
  P-4 folded into the same edit (`index.ts:2113-2330` → `index.ts:2113-2446`).

### Advisory (all applied)

- **P-3** — mission.md Q3 criterion narrowed: "reprocessar notificação → zero duplicata de
  DOMÍNIO (IngestOrder idempotente; dedupe de inbox só com `_id` — ADR-11 emendado; ver
  M-08 Done Means)" — mission-level criterion now matches the N-4 amendment; P6 cannot
  author a row-count criterion the partial key cannot satisfy.
- **P-4** — stale anchor fixed inside the ADR-14 amendment: `index.ts:2113-2446` (client
  literal closes at :2446; verified by the auditor). `codebase-read-side.md:105` keeps its
  open-ended "2113-2330+" — a P2 research observation, not a range claim; not flagged.
- **P-5** — SECOND arm: M-05 F-01 and F-02 Expected Output now pin that the fee/divergence
  writer wiring lives INSIDE the listings composition package under M-04's recorded
  additive lock and that `root.go` is NOT touched by M-05 (matrix cell `—` is now true by
  construction). First arm (add a root.go cell) rejected: M-05 has no reason to enter
  root.go — the composition package is already the additive-lock surface, and keeping the
  cell `—` preserves the lane-C root.go surface exclusively for M-06/M-07.

## Effect on prior artifacts

- `p5-input-r03.sha256` INVALIDATED by the fold edits. New manifest: `p5-input-r04.sha256`
  (same 46-file set, post-fold content).
- mission.md ADR-14 now carries a recorded amendment (P-2); ADR-07 unchanged (it was
  already the ratified narrow form — the drift was in IC-06/F-03, now corrected).
- r02 closures untouched: no N-1..N-10 locus was reopened; P-1's fix narrows the N-1
  mechanism without reintroducing any of N-1's ambiguity (JobFunc byte-identical, no
  concrete job edited — both clauses preserved verbatim).

## Next

1. Freeze `p5-input-r04.sha256` over the same 46-file set (post-fold content).
2. Re-dispatch cold decomposition auditor (r04) on the frozen manifest; persist verbatim
   IMMEDIATELY on return.
3. Advance to P6 only on r04 PASS. Sol P5 retroactive touchpoint remains MANDATORY before
   `status: planned` (≥ 2026-08-05).
