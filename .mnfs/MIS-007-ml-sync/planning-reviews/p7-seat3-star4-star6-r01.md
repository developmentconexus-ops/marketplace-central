# P7 r01 — Claude crew seat 3/5 verbatim (★4 Verifiability + ★6 Evidence Honesty)

```yaml
round: r01
seat: 3/5
scope: "★4 Verifiability + ★6 Evidence Honesty"
model: opus (cold, read-only)
manifest: planning-reviews/p7-input-r01.sha256
manifest_digest: 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e
task_id: aca532ffe4e19f46a
transport: task-notification result (in-context), persisted same turn
verdict: "★4 FAIL (4-A MIS07-C2 p95-vs-mediana n=3; 4-B M01-U2 baseline pré-merge ausente); ★6 PASS; 5 advisories"
persisted: 2026-08-01
```

Verbatim result below (unedited):

---

# P7 Readiness Review — MIS-007-ml-sync, round r01
## Claude crew seat 3 of 5 — SCOPED pass: ★4 Verifiability + ★6 Evidence Honesty

**Mission root:** `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync`
**Rubric applied:** `C:\Users\leandro.theodoro\.claude\plugins\cache\mnfs-harness\mnfs-workflow\0.2.0\skills\mission-planning\references\readiness-review-rubric.md` (read first)
**Manifest:** `planning-reviews/p7-input-r01.sha256`, 65 entries.
**Digest line read (line 66, verbatim):**
`# top-level digest (sha256 of the 65 sorted entries above): 3429ceb0a3c91e09c2eef18af526dd14bb44a46ddc9c94303bbc146238cd9f3e`
— matches the dispatched expected digest. **Attestation caveat: I read the recorded digest line; I did not recompute the hash (no shell used). Digest is attested from the artifact, not independently verified by this seat.**

Coverage: **65 of 65 manifest paths read.** No unread-path residue. No overall seven-★ verdict computed (scoped seat).

---

## ★4 Verifiability — **FAIL** (2 loci)

Procedure run: every criterion in the mission VC and the 9 milestone VCs checked for command-or-interaction + observable Expected + named blocking failure + concrete evidence path; grep for generic-verb stand-ins (`funciona`, `correto`, `adequado`, `válido`, `continua`, `suporta`) and for `Expected:` naming a location instead of an observable; statistic/threshold agreement between `Expected` and `Blocking failure`.

The corpus is overwhelmingly strong — must-fail pairs are named (`M-04/validation-contract.md:61-64` requires the RED test to name itself via `failure_token=test=`), vacuous passes are pre-banned by name (`M-01/validation-contract.md:48-50`: *"'eventually succeeds' sem asserção de tempo = passe vácuo, REPROVA o critério"*), and fixture-shape traps are pinned (`M-04/validation-contract.md:135`: *"fixture de 1 página (passe vácuo de paginação — lição CHIP-MERCADO...)"*). Two criteria nonetheless fail the binary-decidability test.

### FAIL 4-A — mission criterion MIS07-C2 states two different statistics; not binary-decidable at the stated sample size

- **Rubric criterion:** ★4 Verifiability — binary verdict; an `Expected` that cannot be decided (or that the blocking failure contradicts) is a ★4 defect.
- **Cited excerpt — `validation-contract.md:67-72` (verbatim):**
  ```
  - Command: browser QA — /pedidos e /anuncios com a conta real sincronizada; medir load
    (DevTools network, DOMContentLoaded→dados renderizados) 3 amostras cada
  - Expected: p95 das amostras <2s por tela; nenhuma request ML no waterfall
  - Actual:
  - Artifact:
  Blocking failure: qualquer amostra >2s causada por chamada viva ML, ou mediana >2s
  ```
- **Defect locus:** `validation-contract.md:69` (in conjunction with `:68` and `:72`).
- **Offending token/value:** `p95 das amostras` — computed over `3 amostras cada` (`:68`), while the blocking failure at `:72` uses a *different* statistic, `mediana >2s`.
- **Why it fails the procedure:** (a) p95 is not defined at n=3 by any stated method — the QA validator must invent one (max? interpolation? nearest-rank ⇒ p95 collapses to the max), so two honest validators can return opposite verdicts on identical measurements; (b) the two statistics can disagree on the same run: samples 1.0s / 1.5s / 2.5s with **no** live ML call ⇒ median 1.5s passes the blocking rule while p95 (nearest-rank = 2.5s) fails the Expected. The criterion has no single decidable outcome. Note the mission's own Q1 is otherwise well-instrumented (MIS07-C1 at `:55-56` names the live-GET blocking failure precisely) — the defect is isolated to the statistic.
- **Yes-if (minimal, inside approved scope):** state ONE statistic at the stated sample size in both places — e.g. `Expected: as 3 amostras <2s por tela; nenhuma request ML no waterfall` and `Blocking failure: qualquer amostra >2s` — or raise the sample count to a number at which p95 is defined and name the computation method, then use that same statistic verbatim in the blocking failure. No new scope: Q1 (`<2s`, zero live ML) is already ratified in `mission.md` Quality Attributes.

### FAIL 4-B — M01-U2 requires a pre-merge baseline that the contract never requires anyone to capture

- **Rubric criterion:** ★4 Verifiability — concrete evidence path; a criterion whose comparison basis is unrecoverable at validation time cannot fail honestly (vacuous pass).
- **Cited excerpt — `M-01-ml-client-hardening/validation-contract.md:174` (verbatim):**
  ```
  | M01-U2 | Refresh/import existente que atravessa o adapter ML continua funcional (caminho sem erro inalterado) | browser drive de 1 operação real que toca o adapter + resposta idêntica à pré-merge |
  ```
- **Defect locus:** `M-01-ml-client-hardening/validation-contract.md:174`, against that same file's `## Evidence Requirements` at `:142-143`, which require only the saved test output and `git diff --stat` of the chip — no baseline capture.
- **Offending token/value:** `resposta idêntica à pré-merge` (compounded by the generic verb `continua funcional`).
- **Why it fails the procedure:** the drive runs after the merge; nothing in the contract obliges anyone to record the pre-merge response, and after the merge that response is unrecoverable. The validator can only assert "looked fine", which is the vacuous-pass shape this mission bans elsewhere. This is not a stylistic gap — the sibling contract already carries the exact fix, which makes the omission a real divergence rather than a house convention: `M-07-pricing-fee-read/validation-contract.md:156` — *"Baseline before/after salvo ANTES do merge (irrecuperável depois)."*
- **Yes-if (minimal, grounded in an existing contract):** add to `M-01/validation-contract.md` `## Evidence Requirements` the M-07 clause verbatim-equivalent — baseline before/after saved BEFORE the merge (unrecoverable afterwards) — naming the operation driven and the captured response body/status, and restate M01-U2's proof column as the comparison of those two saved captures.

---

## ★6 Evidence Honesty — **PASS**

Procedure run: every version-sensitive external claim traced to `research/external-ml-api-facts.md`; each load-bearing claim tested for `verified` (source + date) vs explicit `assumed`/`verify-at-install`; `ran`/`assumed`/`could-not-run` vocabulary checked for declaration AND operationalization; fabricated-default sweep; provenance-token sweep (every `F-r0X-N`, `N-x`, `P-x`, `fato #N` handle cited in a contract resolved against its source file).

- **Anchor — `research/external-ml-api-facts.md:36-37` (verbatim):**
  > Nada aqui é `verify-at-install`: tudo verificado por doc oficial datada ou medição live própria, exceto #5, #10, #11 marcados explicitamente.
- **The three non-verified rows are labeled at the point of claim, not buried** — e.g. `research/external-ml-api-facts.md:21`:
  > `| 11 | Rate limit ~1500 req/min por seller | assumed (fonte parcial); doc oficial recomenda backoff exp + jitter | design §10; T12-429-behavior.json |`
  and the `assumed` label is repeated at every consumption site rather than laundered into fact: `mission.md:142` (*"Limite configurável (fato ~1500 req/min é `assumed`)"*), `M-01/milestone.md:63-65` (*"limiter CONFIGURÁVEL, nunca constante compilada; default conservador"*), `M-01/F-01-resilience-decorator/feature.md:47-48`.
- **Verified-and-load-bearing claims carry source + date + a regression hook.** Fact #1 (`sale_fee` per unit) is the one the fee ledger genuinely rests on, and it is `verified (medição live 2026-07-20; doc NÃO declara)` with the artifact `evidence/ml-api/T2-order-sale-fee.json`, plus the R-4 regression pin at `research/channel-fees-interface-contract.md:107-108` and a live must-fail at `validation-contract.md:88-89`. The plan discriminates correctly between what it measured and what it assumed.
- **`could-not-run` is operationalized, not decorative** — `M-08-webhook-ingest/validation-contract.md:140-141`:
  > Expected: autorização EXPLÍCITA do operador, com data, ANTES do registro do callback no app ML; sem ela, M08-C5 fica `could-not-run` e o milestone NÃO passa por stub
- **Deferrals are dated and owned rather than silently dropped** — `M-01/validation-contract.md:18-20` (*"Deferral registrado: 2026-08-01, dono = live-drives das lanes B/C"*) and `M-09/validation-contract.md:148-149` (*"critério DIFERIDO — não bloqueia close do M-09; bloqueia close da MISSÃO via MIS07-C8"*).
- **No fabricated defaults; the ban is contractual, not aspirational** — `M-09/validation-contract.md:48` (*"entidade nunca-rodada = campos null (NUNCA timestamp fabricado)"*), `:107` (*"cinza 'nunca' (nulls — NUNCA '0 min atrás')"*), `M-06/validation-contract.md:101-103` (margem NULL persisted, `—` on screen, `incompleto[]` naming the missing term).
- **Stale premises were actively refuted rather than inherited** — `mission.md:35-37` self-corrects the design's 16%/22% premise; `research/codebase-read-side.md:155` declares the design's `fee_sync.go:29` reference OBSOLETE; `planning-reviews/p5-reconciliation-r02.md:47-52` records that r01's own claim "EARS-3 is now reachable" was **WRONG** and supersedes it in writing.
- **Un-remeasured numbers are labeled as such** — `research/codebase-read-side.md:81` and `:171` both mark the 10.8s figure as *"valor vem de memória de missão (M-08), não re-medido aqui"*.
- **Provenance handles resolve.** Every cited token I sampled exists at its named source with matching semantics: `F-r04-1` (`p5-reconciliation-r04.md:28-39`) ⇒ `M09-C1:49`/`M09-C2`; `F-r04-2` (`r04.md:43-49`) ⇒ `M09-C3:86`; `F-r05-1` (`r05.md:30-40`) ⇒ `M09-C1:47`; `F-r07-3` (`r07.md:30-40`) ⇒ `M09-C4:107`; `N-4` (`r02.md:53-65`) ⇒ `M08-C3:89-90` and `M-08/milestone.md:78-80`; `F-r06-5` (`r06.md:68-80`) ⇒ `research/codebase-ingest-side.md:50`; `F-r07-4` (`r07.md:55-57`) ⇒ `codebase-ingest-side.md:96-97`. No dangling handle found.
- **Process failures are disclosed against self-interest** — `p5-reconciliation-r01.md:15-21` records that the r01 audit was nearly lost and was recovered post-compaction; `r03.md:19-24` records the compaction gap as a **deviation**; `r08.md:33-38` states plainly that folding advisories after a PASS does not reopen the subgate and pins which manifest the PASS binds to. `mission.md:6` still reads `status: draft` with `planning_phase: readiness` while `mission.md:389-390` keeps the Sol P3/P5/P7 retroactive gates recorded as mandatory from 2026-08-05 — the artifact does not claim a readiness it has not earned.
- **Known dirty state is surfaced, not hidden** — `research/external-ml-api-facts.md:7`: *"ATENÇÃO: dumps untracked com PII scrub pendente ... não commitar sem scrub."*

**Load-bearing test on the single `assumed` claim (#11), explicitly:** it is absent from `mission.md:101-111` `Clarified Decisions → Accepted assumptions:` (which lists 4 other items). I tested whether that omission is a ★6 FAIL and concluded it is not: ADR-02 is constructed so the number is a **configurable default**, never a design premise; the `assumed` label travels with it to all four consumption sites; R-1 carries mitigation/trigger/owner; and `M01-C2` makes a hardcoded limit a blocking failure. If #11 is wrong, the blast radius is one config value — no ADR, IC, or milestone boundary moves. Not load-bearing ⇒ the rubric's "verified OR recorded as an accepted assumption" obligation does not bind. Logged as advisory 6-i below.

---

## Advisory findings (do NOT flip either verdict)

1. **`M-01/validation-contract.md:47`** — `Command: ... -run <TestRetryAfter> (nome exato no spec)` uses a placeholder test name. The observable (clock assertion) is pinned hard at `:48-50`, so the criterion is not vacuous; but the QA validator cannot copy-paste the command. Cheap fix at spec time.
2. **6-i — fact #11 traceability (see above).** Mirroring `#11 rate limit ~1500 req/min = assumed; default conservador 900; falso ⇒ ajuste de config, sem mudança de desenho` into `mission.md` `Accepted assumptions:` would make the single surviving assumption visible to a reader who never opens `research/`. Quality, not a defect.
3. **`validation-contract.md:203-210`** — `## Evidence Requirements` names the destination directory `docs/design/evidence/` while the observables live in the per-criterion `Expected` blocks. Readable as-is; naming the file-per-criterion would remove the last indirection.
4. **`M-05/validation-contract.md:16-17`** — cites the session memory `ml-catalog-offers-pricing-api` for the `?context=channel_marketplace` fact rather than a manifested research row. Mitigated by design: `M-01/F-02-items-multiget-raw-dto/feature.md:55-60` makes the execution VERIFY that behavior and RECORD the outcome, with a named fallback deliverable if it does not hold — so the plan does not depend on the memory being right.
5. **`planning-reviews/p5-reconciliation-r07.md:33-34`** cites `IC-05 §InboxHealth`, a section `r08.md:60-66` later proved does not exist. Archive-only; `r08` fixed the live citations and `r08.md:75-80` explicitly scopes the residual sweep to exclude verbatim archives. Correct handling — recorded here only so a later reader does not re-file it.

---

## Manifest paths read (65 / 65)

All paths relative to `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\`.

**Mission level (2):** `mission.md`, `validation-contract.md`

**M-01 (4):** `M-01-ml-client-hardening/milestone.md`, `/validation-contract.md`, `/F-01-resilience-decorator/feature.md`, `/F-02-items-multiget-raw-dto/feature.md`

**M-02 (6):** `M-02-sync-core-seam/milestone.md`, `/validation-contract.md`, `/F-01-core-ddl/feature.md`, `/F-02-fee-divergence-ports/feature.md`, `/F-03-scheduler-incremental-cursor/feature.md`, `/F-04-read-guard-allowlist/feature.md`

**M-03 (5):** `M-03-orders-shipment-persist/milestone.md`, `/validation-contract.md`, `/F-01-ml-ingest-readers/feature.md`, `/F-02-ingest-order-v1/feature.md`, `/F-03-read-path-switch/feature.md`

**M-04 (6):** `M-04-listings-backfill-ingest/milestone.md`, `/validation-contract.md`, `/F-01-listings-ddl/feature.md`, `/F-02-mass-closure-replacement/feature.md`, `/F-03-backfill-cursor-ingest/feature.md`, `/F-04-scheduler-refresh-wiring/feature.md`

**M-05 (5):** `M-05-listings-fees-divergence/milestone.md`, `/validation-contract.md`, `/F-01-camada2-fee-ingest/feature.md`, `/F-02-stock-divergence/feature.md`, `/F-03-anuncios-fe-contract/feature.md`

**M-06 (5):** `M-06-orders-backfill-decomposition/milestone.md`, `/validation-contract.md`, `/F-01-backfill-incremental/feature.md`, `/F-02-decomposition-camada3/feature.md`, `/F-03-audit-fe-pedidos/feature.md`

**M-07 (4):** `M-07-pricing-fee-read/milestone.md`, `/validation-contract.md`, `/F-01-fee-read-resolver/feature.md`, `/F-02-precos-provenance-fe/feature.md`

**M-08 (4):** `M-08-webhook-ingest/milestone.md`, `/validation-contract.md`, `/F-01-inbox-endpoint/feature.md`, `/F-02-worker-callback/feature.md`

**M-09 (4):** `M-09-sync-observability/milestone.md`, `/validation-contract.md`, `/F-01-sync-health-endpoint/feature.md`, `/F-02-integracoes-health-section/feature.md`

**Research (11):** `research/external-ml-api-facts.md`, `research/p5-prerequisites.md`, `research/codebase-read-side.md`, `research/codebase-ingest-side.md`, `research/channel-fees-interface-contract.md`, `research/divergences-interface-contract.md`, `research/orders-persistence-interface-contract.md`, `research/webhook-inbox-interface-contract.md`, `research/sync-health-interface-contract.md`, `research/sync-ingest-ports-interface-contract.md`, `research/listings-sync-interface-contract.md`

**Planning reviews (9):** `planning-reviews/p3-reconciliation-r01.md`, `planning-reviews/p5-reconciliation-r01.md` … `r08.md`

**Plus (non-manifest, read as instructed):** the rubric; `planning-reviews/p7-input-r01.sha256` (the manifest itself).

**Unread-path sweep: empty.** Every manifested path was read; no criterion in this seat's scope rests on unread content. No repo file was needed to adjudicate ★4 or ★6 beyond one targeted grep of `mission.md` YAML front-matter (`status: draft`, `planning_phase: readiness`) used in the ★6 self-claim check.

**Seat scope reminder:** ★4 = FAIL (2 loci), ★6 = PASS. No seven-★ verdict computed. No file written or edited by this seat.
