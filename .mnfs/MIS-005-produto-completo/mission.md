# MIS-005-produto-completo

```yaml
id: MIS-005
type: mission
status: draft
owner: Mission Strategist
parent: none
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: mission
planning_phase: scope
```

## Objective

Produto completo pós-MVP: tudo que o MIS-004 deliberadamente excluiu. Planejado AGORA em grão de milestone (aprovado no STOP P3 de 2026-07-17); P4–P7 completos obrigatórios ANTES da execução deste mission.

## Outcome

Plataforma operável em produção: writes ML governados, auth/multi-tenant real, convergência por webhooks, histórico de mercado agendado, radar de mercado, pós-venda completo, fiscal completo, repasses/conciliação e hardening de produção.

## Scope — Milestone headlines (aprovadas P3 r01)

| ID | Headline | Nota |
| --- | --- | --- |
| M-01 mercado-radar | Reprecificação/Oportunidades/Monitorados + scheduler snapshots diários + benchmarks oficiais | ligar coleta CEDO (sem histórico retroativo) |
| M-02 repasses | MP release report assíncrono/CSV ingestão agendada + conciliação ERP + calendário | API-MAP §Repasses |
| M-03 difal-config-completo | Agendamento/lembretes/marcar-pago DIFAL + exceções + tela Configurações completa (3 camadas) | herda tabela IC-04 do MIS-004 |
| M-04 pedidos-pos-venda | Cancelados/Devoluções/claims reais, reputação, logística reversa, contestação, NF-e produção | junto/depois M-06 (eventos) |
| M-05 produto-detalhe-completo | Abas Concorrência, Pedidos, Histórico 90d, auditoria, Dados, visits | Histórico exige M-01 (snapshots acumulados) |
| M-06 webhooks-eventos | Topics orders_v2/shipments/claims/items/payments/item competition/fbm stock + reconciliação GET + proteção duplicata/out-of-order | substitui polling do MIS-004 |
| M-07 provedor-externo | Qualificação gated: contrato → canário-5 congelado → batch-50; sem passar = fecha no-provider | research §9–10; NUNCA antes dos gates |
| M-08 full-visits-benchmarks | Estoque Full (inventories), fulfillment ops/eventos, visits, items_to_win/price references, etiquetas em massa, profundidade ERP | absorve inventory/fulfillment depth |
| M-09 auth-multitenancy | Auth middleware, tenant real por request, RBAC mínimo, CORS produção, administração de tenants | baseline: hoje single-tenant sem auth (R-01) |
| M-10 writes-producao | Habilitação governada de writes ML via envelope M-03: linkage resolvido, policy/source time, idempotência, auditoria, reconciliação, rollout controlado pelo operador | 7 gates do lane provider-write |
| M-11 producao-hardening | Deploy/runtime produção, observabilidade, retenção, recovery, security hardening, validação live integrada | |

## Non-Scope

NO-GOs permanentes do research §3 (scraping, busca pública, proxies, SearchAPI/Gecko/Pricefy). Multi-canal (Amazon/Shopee/Magalu) — fora de ambas as missões.

## Current State

Depende do fechamento do MIS-004. Base será main pós-MIS-004; contratos IC-01..IC-06 herdados e extensíveis (compat rules nos ICs).

## Handoff

- Current status: draft (grão milestone, escopo aprovado).
- Current owner: Mission Strategist.
- Next owner: Mission Strategist (sessão nova de planning P4–P7) após MIS-004 fechar.
- Next action: P4 architecture + ICs próprios; re-validar headlines contra estado pós-MIS-004.
- Required artifact paths: este mission.md; futuros M-*/milestone.md.
- Blocked decisions: None — execução gated no fechamento do MIS-004.
