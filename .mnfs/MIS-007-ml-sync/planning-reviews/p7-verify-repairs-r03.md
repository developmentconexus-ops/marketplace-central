# P7 r03 — verificação mecânica dos reparos r02 (assento de conferência)

```yaml
round: r03
seat: verify (não-rubrica — conferência textual dos 17 reparos r02)
model: claude sonnet (caveman:cavecrew-investigator, read-only)
manifest: NENHUM (desvio registrado — shell wedge B-8, ver p7-claude-readiness-r03.md)
verdict: "17/17 PRESENT"
persisted: 2026-08-01
operator_ruling: "Deixe assim e mande sonnet subagents conferir e somente um revisar
  gastou muitos creditos" (2026-08-01 — rodada enxuta, custo)
```

Resultado verbatim do assento (abaixo, sem edição):

---

All 17 items verified against files under `.mnfs/MIS-007-ml-sync`.

| # | Status | Evidence |
|---|--------|----------|
| 1 | PRESENT | `research/listings-sync-interface-contract.md:65-66` — `- PK (tenant_id, installation_id, provider, provider_listing_id, variation_id) — tuple` / `verbatim do ADR-13 (P7 r02 ★2-A: installation_id é OBRIGATÓRIO no tuple...` |
| 2 | PRESENT | `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31-32` — `- 0091: listing_variations PK (tenant_id, installation_id, provider, provider_listing_id,` / `variation_id) — tuple verbatim ADR-13/IC-07 (P7 r02 ★2-A);` |
| 3 | PRESENT | `mission.md:229-230` — ADR-13 child table `(tenant_id, installation_id, provider, provider_listing_id, variation_id)` — 5 columns, unchanged |
| 4 | PRESENT | `research/webhook-inbox-interface-contract.md:31-33` — slug `mercado_livre`, "P7 r02 ★2-B: provider_code real é mercado_livre — market_adapters.go:239, p5-prerequisites.md:113" |
| 5 | PRESENT | `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md:39-40` — slug `mercado_livre` + "não confundir com package Go mercadolivre" |
| 6 | PRESENT | `M-08-webhook-ingest/validation-contract.md:175-176` — "M-09 fecha na lane A antes do despacho da lane D — P7 r02 A-1"; "já fechado na prática" ausente |
| 7 | PRESENT | `mission.md:118-121` (A-2 eco fato #11) e `mission.md:116-117` (A-10 exceção agregados globais) |
| 8 | PRESENT | `M-07-pricing-fee-read/validation-contract.md:68-69` (A-4 fixture-defined) e `:120-121` (A-7 N≥3) |
| 9 | PRESENT | `research/sync-health-interface-contract.md:49-51` (A-5 `:853-855`/`:858`) e `:154-155` (A-10 agregados GLOBAIS) |
| 10 | PRESENT | `M-09-sync-observability/F-01-sync-health-endpoint/feature.md:87-88` (A-13 `entities: []`) e `:95-96` (A-6 `health_*`) |
| 11 | PRESENT | `M-09-sync-observability/validation-contract.md:69` e `:106` — discriminador por estado `consecutive_failures == 0` (A-8, ambos critérios) |
| 12 | PRESENT | `M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md:57-60` — 409 refresh_in_progress, medido codebase-read-side.md:87 (A-12) |
| 13 | PRESENT | `research/listings-sync-interface-contract.md:99` — row Error Matrix 409 (A-12) |
| 14 | PRESENT | `validation-contract.md:124-125` — grep `00{86..95}*.sql` brace expansion (A-14) |
| 15 | PRESENT | `M-01-ml-client-hardening/validation-contract.md:147-149` — scrub + headers de auth NUNCA capturados (A-11) |
| 16 | PRESENT (0 hits fora de planning-reviews) | tuple 4-col só em `p7-claude-readiness-r02.md:41`, `p7-seat2-star2-star3-r01.md:108`, `p7-seat5-doublepass-r02.md:29,30` (registros históricos) |
| 17 | PRESENT | todos os hits de `mercadolivre` fora de planning-reviews são referência de package/path Go, heading de seção, URL context7, padrão de grep no VC, ou registro retrospectivo do achado no Handoff — NENHUM afirma slug/provider_code |

**Totals: 17/17 PRESENT.**
