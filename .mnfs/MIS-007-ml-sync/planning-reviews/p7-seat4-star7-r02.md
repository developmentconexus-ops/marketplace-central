# P7 seat 4 — ★7 Security Posture (adversarial) — round r02 (verbatim)

```yaml
round: r02
seat: 4/5
scope: "★7 Security Posture (adversarial)"
model: claude (cold mission-reviewer crew)
manifest: planning-reviews/p7-input-r02.sha256
manifest_digest: e238f447823ab01eaa1780deb84d0db3280488ab4eda864033d51134f8a4ebaa
transport: task-notification (result in-context, complete)
verdict: "★7 PASS"
persisted: 2026-08-01
```

Resultado verbatim do assento (abaixo, sem edição):

---

## ★7 Security Posture — PASS

**Procedure run:** enumerated every auth/PII/secret/multi-role surface in the to-be-built scope, then verified each is targeted (mitigation + Security-typed validation criterion) or declined-with-reason — adversarial stance, assuming the defect exists.

**Surfaces found and their guards (all cited verbatim):**

1. **Mission-level anchor (rubric's named check).** Security appears in `## Quality Attributes` with mitigation and a non-vacuous validation criterion — `mission.md:354`: "Q2 Security | webhook = hint não-confiável; dado real sempre refetch na API autenticada; dedupe no inbox; IP origem gravado + comparado à allowlist oficial (log-only) | contrato do webhook (P4) | notificação forjada não injeta dado (só fetch idempotente); IP não-oficial aparece em log". Declined dimension carries an explicit reason — `mission.md:365`: "Verificação criptográfica de webhook (assinatura/token) | decisão P1 do operador: webhook = hint não-confiável + refetch autenticado; `ip_official` é INFORMATIVO/log-only, NUNCA gate de aceitação". Four `Type: Security` criteria exist across the tree (`validation-contract.md:119` MIS07-C5, `M-02-sync-core-seam/validation-contract.md:148` M02-C7, `M-08-webhook-ingest/validation-contract.md:41,140` M08-C1/M08-C6).

2. **Unauthenticated public write surface (`POST /webhooks/{provider}`), abuse envelope.** Bounded body before parse — `research/webhook-inbox-interface-contract.md:158`: "Cap de body 64KB ANTES do parse (LimitReader)". Attacker-controlled `resource` flowing into an authenticated URL is pinned to a named owner — `webhook-inbox-interface-contract.md:52-55`: "`resource` só é usável como ponteiro quando casa `^/orders/[0-9]+$` (ancorado)... O `resource` é attacker-controlled — NUNCA concatenado em URL de chamada ML autenticada sem casar o pin; o guard mora no package webhook (dono M-08)". The criterion tests it including traversal and is non-constant on both IP branches — `M-08-webhook-ingest/validation-contract.md:46-48`: "POST com `resource` fora do pin (ex.: `/users/123` e `/orders/../users/me` — traversal...) + par de POSTs p/ derivação de IP... controles positivo/negativo". Forgery-to-domain-write is the blocking failure — `M-08.../validation-contract.md:57`: "forja alcançando tabela de domínio, ou corpo de webhook persistido como dado". Spoofable XFF is acknowledged, not silently trusted — `webhook-inbox-interface-contract.md:65-66`: "Header é forjável fora do túnel ⇒ `ip_official` é INFORMATIVO/log-only".

3. **PII (buyer fiscal + shipment delivery data).** Raw persistence banned at ADR, schema, reader, and criterion levels: `mission.md:154-155` "**Raw de `billing_info` NUNCA persiste**"; `research/orders-persistence-interface-contract.md:55-56` "SEM coluna `raw` (P7 r01 B-7, ADR-03 emendado: payload de shipment carrega PII de entrega...)"; `M-03.../F-01-ml-ingest-readers/feature.md:88` "Assert: raw não contém substring `billing_info`"; Security-typed schema gate `M-02.../validation-contract.md:152-154`: "colunas TIPADAS só...; nenhuma coluna `raw`/jsonb de billing_info"; mission gate `validation-contract.md:123-128` greps migrations 0086–0095. Evidence-side PII: `M-03.../validation-contract.md:18-19` "dumps `docs/design/evidence/ml-api/`, PII scrub antes de qualquer commit".

4. **Live ML write (callback registration).** Gated by a Security-typed, non-stubable criterion — `M-08.../validation-contract.md:145-147`: "autorização EXPLÍCITA do operador, com data, ANTES do registro do callback no app ML; sem ela, M08-C5 fica `could-not-run` e o milestone NÃO passa por stub"; blocking at `:149`.

5. **Tenant scoping in new tables.** `research/divergences-interface-contract.md:44` "`tenant_id` NOT NULL"; `research/channel-fees-interface-contract.md:46` "`tenant_id` NOT NULL (tenant-scoped, sempre no WHERE)"; `order_shipments`/`listing_variations` PKs lead with `tenant_id` (`M-02.../F-01-core-ddl/feature.md:30`, `M-04.../F-01-listings-ddl/feature.md:31`). Inbox `tenant_id` NULL is deliberate (unmapped notifications) under the recorded assumption `mission.md:115-116`.

6. **New read endpoint `/sync/health`.** No new auth boundary — `research/sync-health-interface-contract.md:175`: "Auth/session existentes da API; nada cruza origem nova"; read-only, zero ML.

No auth/PII/secret/multi-role surface is silently omitted. **★7 PASS.**

## Advisories

1. **Inbox flood: no retention/prune policy and no per-source throttle on the public endpoint.** IC-04 mitigates per-request abuse (64KB cap, interactive 15s class, dedupe when `_id` present) but declares no retention for terminal rows (`done`/`malformed`/`dropped`) and no cap on insert rate; a sustained flood grows `notifications_inbox` unboundedly, and each plausible forgery also spends one authenticated ML refetch from the installation's token-bucket (quota starvation of legitimate sync). Locus: `research/webhook-inbox-interface-contract.md:97-108` (Persistence Expectations — no retention row). Suggest: one line in IC-04 naming a retention/prune rule and acknowledging the refetch-amplification bound. Not auto-fixable (needs a decided value).
2. **Inbox health aggregation has no tenant predicate while the port is tenant-parameterized.** `research/sync-health-interface-contract.md:39` (`WebhookStats(ctx, tenantID)`) vs `:149-150` ("inbox agregados (`status='received'` count; `dropped` nas últimas 24h; max `received_at`)") with `webhook-inbox-interface-contract.md:49` (`tenant_id` NULL). Covered today by the recorded single-tenant assumption (`mission.md:115-116`), but that assumption's clause "tudo já tenant-scoped" is not literally true of this aggregate — pin the intended semantics (global vs tenant-filtered) in IC-05 so a future multi-tenant change doesn't silently leak counts. Auto-fixable (one sentence).
3. **M-01 pre-merge baseline capture lacks a scrub/no-secret rule.** `M-01.../validation-contract.md:143-146`: "resposta da operação real que atravessa o adapter ML (payload/headers relevantes)" — if the driven operation is an order read, the captured payload can carry buyer PII, and "headers" could include `Authorization: Bearer` if captured at the adapter (`research/codebase-ingest-side.md:23`). M-03's VC pins scrub-before-commit; M-01's does not. Auto-fixable (copy the M-03 scrub clause + "nunca headers de auth").

## Coverage

Read line-by-line (absolute paths under `C:\Users\leandro.theodoro\Documents\marketplace-central\.mnfs\MIS-007-ml-sync\`):
- rubric (`C:\Users\leandro.theodoro\.claude\plugins\cache\mnfs-harness\mnfs-workflow\0.2.0\skills\mission-planning\references\readiness-review-rubric.md`), `planning-reviews\p7-input-r02.sha256`
- `mission.md`, `validation-contract.md`
- `research\webhook-inbox-interface-contract.md`, `research\orders-persistence-interface-contract.md`, `research\sync-health-interface-contract.md`
- `M-08-webhook-ingest\milestone.md`, `M-08-webhook-ingest\validation-contract.md`, `M-08-webhook-ingest\F-01-inbox-endpoint\feature.md`, `M-08-webhook-ingest\F-02-worker-callback\feature.md`
- `M-01-ml-client-hardening\validation-contract.md`, `M-02-sync-core-seam\F-01-core-ddl\feature.md`, `M-02-sync-core-seam\validation-contract.md` (lines 130–190), `M-03-orders-shipment-persist\F-02-ingest-order-v1\feature.md`, `M-03-orders-shipment-persist\validation-contract.md`, `M-04-listings-backfill-ingest\F-01-listings-ddl\feature.md`

Swept by grep (whole mission tree): `tenant_id`; `Authorization|Bearer|access_token|client_secret|APP_SECRET|ngrok` (zero secrets embedded in artifacts); `Type: Security` (all validation contracts); `PII|billing_info|documento|scrub` (surfaced the relevant lines of `M-03...\F-01-ml-ingest-readers\feature.md`, `research\channel-fees-interface-contract.md:46`, `research\divergences-interface-contract.md:44`, `research\external-ml-api-facts.md:7` — judged from those excerpts, not full reads). Not read: M-05/M-06/M-07/M-09 milestone/feature bodies, M-04 F-02..F-04, remaining research files, `planning-reviews\p3-*`/`p5-*` (grep coverage only — none of these adds an auth/PII surface beyond those enumerated: they are read-only FE/resolver/observability work on existing authenticated screens). Grep sweeps incidentally matched `planning-reviews\p7-*` round-1 review files; their content was disregarded per dispatch — all findings above were verified independently at source loci.
