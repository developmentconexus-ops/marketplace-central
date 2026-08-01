# P7 seat 5 — adversarial double-pass ★2 + ★7 — round r02 (verbatim)

```yaml
round: r02
seat: 5/5
scope: "★2 Consistency + ★7 Security Posture (independent adversarial double-pass)"
model: claude (cold mission-reviewer crew)
manifest: planning-reviews/p7-input-r02.sha256
manifest_digest: e238f447823ab01eaa1780deb84d0db3280488ab4eda864033d51134f8a4ebaa
transport: task-notification (result in-context, complete)
verdict: "★2 FAIL (★2-A PK listing_variations; ★2-B slug webhook); ★7 PASS"
persisted: 2026-08-01
```

Resultado verbatim do assento (abaixo, sem edição):

---

Review complete. All 65 manifested entries covered; findings below.

Mission root: `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync` — P7 round r02, seat 5/5 (adversarial double-pass), scope ★2 + ★7.

## ★2 Consistency — FAIL

Red-team grep of every enum, PK/shape, route token, origem vocabulary, cursor key, status enum, migration range, and list ordering across all 65 manifested files found two genuine cross-artifact shape divergences. Everything else swept clean: origem CHECK vocabulary (`api_listing_prices|api_shipping_options(reserve)|api_order|api_shipment|config`) is identical across IC-01, M-05/M-06/M-07 briefs and `validation-contract.md:164`; inbox status enum (`received/processing/done/malformed/unmapped/dropped`) identical across IC-04, M-08 F-01/F-02 and M08-C2/C4; phase vocabulary `{backfill, incremental, sweep}` identical across IC-05/IC-06/M-02 F-03/M-04 F-03/M-06 F-01; subject_id/entity_id composite formats identical across IC-01/IC-02/M-05 F-02/M-06 F-02; migration ranges disjoint and matching the ownership matrix (0086-0089/0090-0092/0093/0094-0095); every list operation declares order (IC-02 `detected_at desc`, IC-04 FIFO `received_at asc`, IC-05 name asc, IC-03/IC-07 by-reference to measured baselines in `research/p5-prerequisites.md:95`); the dedupe key `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL` is identical in mission.md ADR-11, IC-04 and M-08.

**Finding ★2-A — `listing_variations` PK stated two incompatible ways (5 columns vs 4).**
- Excerpt (spine): `mission.md:224-225` — "**ADR-13 `listing_variations` aditiva; PK de `listings` NÃO muda.** Child table `(tenant_id, installation_id, provider, provider_listing_id, variation_id)`, mesmo writer"
- Excerpt (binding IC): `research/listings-sync-interface-contract.md:65` — "PK `(tenant_id, provider, provider_listing_id, variation_id)`."
- Excerpt (DDL brief): `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31` — "0091: `listing_variations` PK (tenant_id, provider, provider_listing_id, variation_id);"
- Defect locus: `research/listings-sync-interface-contract.md:65` ∥ `mission.md:225`; offending token: `installation_id` present in the ADR-13 tuple, absent from IC-07/F-01. No manifested amendment records the drop (grep of `p5-reconciliation-r01..r08` for `installation_id|ADR-13|variations` = zero hits on this seam), and A-13 was adopted at `planning-reviews/p3-reconciliation-r01.md:37` with the spine reading. Note the parent's real PK is `(tenant_id, installation_id, provider_listing_id, variation_id)` (`research/codebase-ingest-side.md:65`), so the 4-column child key cannot address parent rows per installation — the divergence is load-bearing for 0091 DDL, and a worker reading the spine vs the IC produces incompatible migrations.
- Yes-if: pin ONE tuple in all three loci — either amend `mission.md:225` (ADR-13) to the 4-column PK with a recorded reason for dropping `installation_id` (e.g., single-installation-per-provider invariant, stated), or amend `research/listings-sync-interface-contract.md:65` + `M-04.../F-01-listings-ddl/feature.md:31` to include `installation_id` per the ratified A-13.

**Finding ★2-B — webhook route slug `mercadolivre` claimed equal to a provider_code vocabulary that is actually `mercado_livre`.**
- Excerpt: `research/webhook-inbox-interface-contract.md:31` — "Rota `POST /webhooks/{provider}` (`provider` = slug `mercadolivre`, mesmo vocabulário provider_code das 4 superfícies — CHIP-PED-FILA)"; repeated at `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md:39` — "slug provider_code das 4 superfícies (`mercadolivre`)".
- Contradicting manifested fact: `research/p5-prerequisites.md:113` — "`client.createIntegrationInstallation({installation_id, provider_code: "mercado_livre", ...})`". Repo confirms the vocabulary (read-only ★-check): `apps/server_core/internal/composition/market_adapters.go:239` — `const mercadoLivreProviderCode = "mercado_livre"`; `apps/web/src/pages/integracoes/IntegracoesPage.tsx:514` — `i.provider_code === "mercado_livre"`; `apps/web/src/pages/pedidos/PedidosTable.tsx:99` — "provider_code is the channel slug (\"mercado_livre\")". The token `mercadolivre` is the Go *package* name, not the provider_code (`M-01-ml-client-hardening/milestone.md:48` itself distinguishes "dir `mercado_livre`, package `mercadolivre`").
- Defect locus: `research/webhook-inbox-interface-contract.md:31`; offending token: `mercadolivre`. Blast radius: the `provider` column value in `notifications_inbox` (and its dedupe UNIQUE), the health-block provider join (IC-05), and the public callback URL registered in the ML app (M08-C6) all key on this token; two features resolving it differently produce a silent dedupe/join miss.
- Yes-if: change the slug token to `mercado_livre` at `webhook-inbox-interface-contract.md:31` and `M-08.../F-01-inbox-endpoint/feature.md:39` (matching the four-surface provider_code vocabulary the IC itself invokes), or — if a distinct URL slug is genuinely intended — delete the claim "mesmo vocabulário provider_code das 4 superfícies" and pin an explicit slug→provider_code mapping row in IC-04.

## ★7 Security Posture — PASS

Enumerated every auth boundary and PII/attacker-controlled surface; each has a mitigation plus a Security-typed criterion, or an explicit declined-with-reason. No silent omission found.

- Mission-level posture present: `mission.md:354` — "Q2 Security | webhook = hint não-confiável; dado real sempre refetch na API autenticada; dedupe no inbox; IP origem gravado + comparado à allowlist oficial (log-only) | ... | notificação forjada não injeta dado (só fetch idempotente); IP não-oficial aparece em log". Crypto verification explicitly declined with operator reason: `mission.md:365` — "Verificação criptográfica de webhook (assinatura/token) | decisão P1 do operador: ... NUNCA gate de aceitação".
- Public unauthenticated write surface (`POST /webhooks/{provider}`): attacker-controlled `resource` is pinned and never composed into authenticated URLs — `research/webhook-inbox-interface-contract.md:54`: "`resource` é attacker-controlled — NUNCA concatenado em URL de chamada ML autenticada sem [casar o pin `^/orders/[0-9]+$`]"; tested with traversal payloads (`/orders/../users/me`) and forge in Security-typed `M08-C1` (`M-08-webhook-ingest/validation-contract.md:38-60`), which also carries non-constant positive/negative XFF controls, so the spoofable IP derivation is never a criterion of acceptance (`webhook-inbox-interface-contract.md:66` — "INFORMATIVO/log-only ... NUNCA gate de aceitação"). Envelope hardening (always-200, 64KB truncation with marker, interactive 15s deadline proved by trickle outside `registerBatchRoutes`) is criterioned in M08-C2.
- Live ML write gate: Security-typed `M08-C6` requires recorded explicit operator authorization before callback registration; without it M08-C5 is `could-not-run` — no stub pass.
- PII containment: raw jsonb confined to `listings` (B-7); `order_shipments` typed-only (`research/orders-persistence-interface-contract.md`, "SEM coluna raw"); billing_info raw banned with assert (`M-03.../F-01-ml-ingest-readers/feature.md:45-46,88`); schema-level grep + synthetic-PII forge criterioned in Security-typed `MIS07-C5` (`validation-contract.md:116-133`) and `M02-C7`; the evidence-dump scrub pendency is named (`research/external-ml-api-facts.md`, `validation-contract.md:210-211`).
- New read surfaces (`GET /sync/health`, /listings filter, /orders additive DTO) sit behind existing API auth (`research/sync-health-interface-contract.md`, "Auth/session existentes da API"); no multi-role surface exists in scope.

## Advisories

1. (auto-fixable) `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md:58` — "refresh enfileira/rejeita com 409 semântico EXISTENTE do batch": hedged disjunction ("enfileira/rejeita"), and the 409 case has no row in IC-07's Error Matrix. It is a preserved existing behavior (`research/codebase-read-side.md:87` — "POST /listings/refresh, 202+409 `refresh_in_progress`"), so not a ★2 divergence; de-hedge and add the row by reference.
2. (auto-fixable) `M-07-pricing-fee-read/validation-contract.md:67` — "camada 1 → origem `api_listing_prices`": IC-01 defines no layer-1 producer this mission (layer 1 = MIS-008); the fixture should state that the test plants the camada-1 row, so the expected origem is fixture-defined rather than contract-derived.
3. (auto-fixable) `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:87` — "entities todas null-honestas" for a tenant with zero `sync_state` rows reads two ways; M09-C1 pins "`entities: []`" and IC-05 pins "uma row por entidade REGISTRADA". Reword the negative scenario to "`entities: []`".
4. Flood without `_id` bypasses dedupe (ADR-11 emendado) and grows `notifications_inbox` unboundedly; no retention/purge policy is declared anywhere in IC-04/M-08. Rows are inert (M08-C3), but storage exhaustion on the public surface is unaddressed.
5. Amplification residual: each plausible forged notification costs the worker ≥1 authenticated ML GET (refetch posture), drawing down the per-installation token bucket (M-01) — a forge flood can starve real sync freshness. M-08 names flood mitigations (cap/dedupe/deadline) but not this budget-drain vector.
6. (auto-fixable) `validation-contract.md:124` — evidence command "`grep -rn "raw" apps/server_core/migrations/0086*..0095*`": `0086*..0095*` is not valid glob/range syntax; the command as written will not run.

## Coverage

Read line-by-line (all under the mission root unless noted): `mission.md`; `validation-contract.md`; `research/channel-fees-interface-contract.md`; `research/divergences-interface-contract.md`; `research/orders-persistence-interface-contract.md`; `research/webhook-inbox-interface-contract.md`; `research/sync-health-interface-contract.md`; `research/sync-ingest-ports-interface-contract.md`; `research/listings-sync-interface-contract.md`; `research/external-ml-api-facts.md`; `research/p5-prerequisites.md`; `research/codebase-ingest-side.md`; `research/codebase-read-side.md`; `planning-reviews/p3-reconciliation-r01.md`; `planning-reviews/p7-input-r02.sha256`; all 9 `M-0*/milestone.md`; all 23 `M-0*/F-*/feature.md`; all 9 `M-0*/validation-contract.md`; the rubric at `C:\Users\leandro.theodoro\.claude\plugins\cache\mnfs-harness\mnfs-workflow\0.2.0\skills\mission-planning\references\readiness-review-rubric.md`.

Swept by grep only (disclosed, not read line-by-line): `planning-reviews/p5-reconciliation-r01..r08.md` (three targeted sweeps: `installation_id|ADR-13|variations`; `B-7|B-8|B-9|slug|webhook|PII|raw`; token sweeps for `mercadolivre|mercado_livre` and `listing_variations|variation_id`). Mission-root-wide greps incidentally surfaced match lines from non-manifested files (`p3-claude-candidate-r01.md`, `p3-opus-counterproposal-r01.md`, `p5-claude-decomposition-audit-r0*.md`, `p5-passes-r01.md`, and two `p7-*` review outputs); I did not open any of those files and no finding relies on their content — every FAIL locus and citation is grounded in manifested files or repo source.

Repo source verified read-only (permitted for claim-testing): `apps/server_core/internal/composition/market_adapters.go:239`, `root.go:741/743`, `apps/server_core/internal/modules/integrations/transport/http_handler_test.go:103,130`, `apps/web/src` grep for `provider_code` (incl. `IntegracoesPage.tsx:514`, `PedidosTable.tsx:99`).

Per dispatch: no seven-★ verdict computed; the dispatching session folds these per-criterion findings.
