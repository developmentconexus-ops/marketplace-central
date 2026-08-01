# P7 seat 3 — ★4 Verifiability + ★6 Evidence Honesty — round r02 (verbatim)

```yaml
round: r02
seat: 3/5
scope: "★4 Verifiability + ★6 Evidence Honesty"
model: claude (cold mission-reviewer crew)
manifest: planning-reviews/p7-input-r02.sha256
manifest_digest: e238f447823ab01eaa1780deb84d0db3280488ab4eda864033d51134f8a4ebaa
transport: task-notification (result in-context, complete; entidades HTML da notificação
  decodificadas p/ literais na transcrição)
verdict: "★4 PASS; ★6 PASS"
persisted: 2026-08-01
```

Resultado verbatim do assento (abaixo, sem edição):

---

## ★4 Verifiability — PASS

Procedure ran over mission VC + 9 milestone VCs (68 `Expected:` blocks) + Validation Expectations of all 24 feature briefs. Every Required criterion has command + observable expected value + named blocking failure + evidence-path anchor. Findings that drove verdict:

- Single-statistic thresholds pinned (r01 B-5 repaired): `validation-contract.md:69` — "Expected: as 3 amostras <2s por tela (estatística ÚNICA — P7 r01 B-5); nenhuma request ML no waterfall". Blocking failure at `validation-contract.md:73` names the same statistic ("qualquer das 3 amostras ≥2s") — no dual-statistic divergence.
- Irrecoverable-baseline shape repaired (r01 B-6): `M-01-ml-client-hardening/validation-contract.md:144-146` requires baseline "capturada e salva ANTES do merge … irrecuperável depois"; same rule at `M-07-pricing-fee-read/validation-contract.md:155` — "Baseline before/after salvo ANTES do merge (irrecuperável depois)." Comparison targets exist at check time.
- Must-fails not satisfiable by defect-free runs: `M-02-sync-core-seam/validation-contract.md` M02-C4 requires red output with named message ("terminal cursor must be non-nil") before green; `M-03-orders-shipment-persist/F-03-read-path-switch/feature.md:86` — "import vivo reintroduzido → allowlist test VERMELHO nomeando o site." Red-names-target idiom throughout.
- Vacuous-pass guards explicit: `M-01-ml-client-hardening/validation-contract.md` M01-C1 pins elapsed ≥2s clock assertion and states "eventually succeeds = passe vácuo REPROVA"; `M-08-webhook-ingest/validation-contract.md:51` pins resource guard to regex `^/orders/[0-9]+$` with negative case producing `malformed` row and "NENHUMA URL ML construída"; `:52-54` pins XFF positive+negative controls ("critério não-constante"). `M-09-sync-observability/validation-contract.md` M09-C3 pins byte-equal canonical default `{"last_notification_at":null,"pending":0,"dropped_24h":0}` with JSONB structure caveat.
- Unpinned references absent: every VC carries `base_sha: dd89d4b3` (e.g. `validation-contract.md:13`); fee origin vocabulary enumerated at `validation-contract.md:163-164` ("origem ∈ vocabulário IC-01 (api_listing_prices | api_order | api_shipment | config)"); tolerances pinned in ICs (estoque 0, tarifa R$0.01).
- Generic-verb sweep over all VCs: only benign hits (`correction_attempts` boilerplate; `M-01` M01-U2 "continua funcional" — grounded by diff of the two saved captures, decidable).
- Stub/deferral honesty: `M-08-webhook-ingest/validation-contract.md:146` — without operator authorization "M08-C5 fica `could-not-run` e o milestone NÃO passa por stub"; M09-C6 deferral explicitly discharged via MIS07-C8 (`validation-contract.md:180-184`).

No criterion found that is non-binary, vacuous, baseline-unrecoverable, or evidence-label-only. PASS.

## ★6 Evidence Honesty — PASS

Claim registry audited (`research/external-ml-api-facts.md`, 20 facts) plus both codebase research files and p5-prerequisites, all against the version-sensitive-claim rule:

- No silent acceptance: every fact carries status + source + date — live measurement 2026-07-20 or doc via context7 dated 2026-07-31. Closing declaration `research/external-ml-api-facts.md:36-37` — "Nada aqui é `verify-at-install`."
- The one `assumed` fact is structurally severed from correctness: `research/external-ml-api-facts.md:21` — fact #11 rate limit "~1500 req/min … `assumed` (fonte parcial)". Design does not depend on the value: ADR-02 makes the limit configurable, `M-01-ml-client-hardening/validation-contract.md` M01-C2 requires "limite CONFIGURÁVEL, nunca constante compilada", risk R-1 (mission.md) has mitigation/trigger/owner, and `planning-reviews/p3-reconciliation-r01.md:19` records the ratified reasoning verbatim: "limite configurável pois fato #11 é `assumed`". Not load-bearing → not a ★6 defect.
- Explicit unknowns stay unknown: facts #5 and #10 marked "NÃO confirmado" with design independence stated in the registry itself — no criterion or IC consumes them as verified.
- Codebase claims verified at pin: `research/codebase-read-side.md:8-9` — "base: main (dd89d4b3) / method: leitura estática de arquivos"; `research/codebase-ingest-side.md:3` — "Base: main @ dd89d4b3 … file:line exatos"; `research/p5-prerequisites.md` 12 items verbatim with file:line at same SHA — matching `base_sha` of every VC, so criteria citing them compare against a recoverable reference. Zero `assumed`/`verify-at-install` markers in either codebase file (grep = 0 hits).
- Evidence pendencies disclosed, not hidden: `validation-contract.md:209-211` — "PII scrub antes de commit — pendência conhecida em `docs/design/evidence/ml-api/`"; `:213-214` — "NENHUM write ML (registro de callback incluso) sem autorização explícita do operador registrada."

No version-sensitive claim silently accepted; no load-bearing claim outside verified/accepted-assumption channels. PASS.

## Advisories

(auto-fixable; none flips a criterion)

1. `M-07-pricing-fee-read/validation-contract.md:118` — "N simulações reais": N unpinned. Verdict stays binary (every delta must be explained; `:162` "Delta órfão no before/after = blocking"), but pinning N (e.g. ≥3, incl. 1 com camada-2 e 1 sem ledger, per `:145-147` Drive steps) removes sampling discretion.
2. `M-09-sync-observability/validation-contract.md:104` — "verde (sucesso recente)": "recente" not numerically pinned in the VC (fixture-determined in context; suggest citing the exact threshold from IC-05/FE mapping).
3. `M-09-sync-observability/validation-contract.md:68` — same word "recente" in timestamp rendering criterion; anti-criterion is concrete ("NUNCA 'há N dias'/cinza") so decidable, but a numeric bound would harden it.
4. `research/external-ml-api-facts.md:21` fact #11 (`assumed`) is not echoed under mission.md `Clarified Decisions -> Accepted assumptions:` (4 entries there omit it). Design severs dependency (see ★6), so advisory only; one echoed line would make the registry and mission self-consistent.

## Coverage

Read line-by-line (all under `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\` unless noted):
- Rubric: `C:\Users\leandro.theodoro\.claude\plugins\cache\mnfs-harness\mnfs-workflow\0.2.0\skills\mission-planning\references\readiness-review-rubric.md`
- `planning-reviews/p7-input-r02.sha256` (manifest, 65 entries)
- `mission.md`; mission `validation-contract.md`
- All 9 `M-0*/validation-contract.md`; all 9 `M-0*/milestone.md`; all 24 `M-0*/F-*/feature.md`
- `research/external-ml-api-facts.md`; `research/p5-prerequisites.md`
- All 7 interface contracts: `research/channel-fees-interface-contract.md`, `divergences-…`, `listings-sync-…`, `orders-persistence-…`, `sync-health-…`, `sync-ingest-ports-…`, `webhook-inbox-interface-contract.md`
- `research/codebase-read-side.md` (lines 1-60) and `research/codebase-ingest-side.md` (lines 1-30) — partial reads, headers + method + site inventory

Swept by grep only (disclosed):
- Claim-marker sweep (`assumed|NÃO confirmado|verify-at-install|unverified`) over `research/` and `planning-reviews/p3-*` — one hit outside manifest (`p3-opus-counterproposal-r01.md`; not used as evidence), one manifest hit (`p3-reconciliation-r01.md:19`, cited above)
- Generic-verb sweep (`corretamente|funciona|works|correct|proper` etc.) over all `**/validation-contract.md`
- `Expected:` count sweep over all VCs (68 blocks / 10 files)
- Targeted sweeps: `simula|before/after|recente` in M-07 VC; `recente|verde` in M-09 VC

Not read: `planning-reviews/p5-reconciliation-r01..r08.md` and remainder of `p3-reconciliation-r01.md` (manifest-listed process records; only grep-swept), remainders of `codebase-read-side.md`/`codebase-ingest-side.md`. `planning-reviews/p7-*` review outputs: not opened (per dispatch prohibition).
