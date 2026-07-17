# P3 Claude Candidate — r01 (2026-07-17)

Missões: **MIS-004 mvp-demo** (demo cliente 2026-07-20) + **MIS-005 produto-completo** (grão milestone).
Evidência congelada: ver `p3-input-r01.sha256`. Decisões P1: `research/p1-clarified-decisions-2026-07-17.md`.

## Outcome (MIS-004)

Demo segunda: história "do estoque ERP (planilha .xlsx do cliente) ao veredicto honesto 'vale a pena vender?' com preço vs mercado ML, e simulador de margem real" — no design novo (papel+verde), com Pedidos funcional, dados live read-only, zero writes ML.

## Architecture spine (ADR-lite, MIS-004)

### ADR-01: Base de execução = main pós-merge W1
- Decisão: todos os chips MIS-004 partem de main APÓS merge de chip/m-02 (FE platform), chip/m-03 (mutation envelope), chip/sat. Planning artifacts aterrissam antes; execução gated no merge.
- Prevents: retrabalho contra shell antigo; duplicar envelope de mutação.
- Must preserve: M-03 envelope como ÚNICO write path (UI → /mutations → poller → adapter).
- Trade-off: execução espera merges W1 (horas, hub já em fechamento).
- Validation impact: base SHA 40-hex registrado por chip.

### ADR-02: Fonte ERP dual — Reader port ganha adapter de snapshot importado (.xlsx)
- Decisão: novo módulo `erp_import`: endpoint de import .xlsx + tabelas de snapshot ERP (produtos/estoque/custo) + adapter implementando o subset do port `internal_read/ports.Reader` usado pelo MVP; Oracle/Sankhya adapter permanece; seleção por installation/config. Import é lote com protocolo + relatório de validação (linhas aceitas/rejeitadas/avisos).
- Prevents: demo refém de VPN Sankhya; stub em live path (proibido, profile §7).
- Must preserve: contrato de colunas IC-02 (obrigatórias CODPROD, DESCRPROD, CUSTO, ESTOQUE_FISICO; opcionais ESTOQUE_RESERVADO, EAN, REFFORN, MARCA, NCM); opcional ausente ⇒ estado honesto (sem EAN ⇒ matching nunca auto-ACCEPT).
- Trade-off: módulo novo + entrada governance (via merge do chip).
- Validation impact: import de planilha exemplo real; linha inválida rejeitada com motivo; QA drive tela Vínculos pós-import.

### ADR-03: Correção do contrato de identidade (research §4)
- Decisão: CODPROD = SKU canônico; REFERENCIA = EAN/GTIN (validar checksum, unicidade); REFFORN = referência fabricante; marca+atributos = âncoras. Corrigir semântica em internal_read/catalog ANTES de ligar matching novo (hoje reader trata REFERENCIA como referência de fabricante — defeito estrutural confirmado no baseline).
- Prevents: colisões de EAN viram vínculos errados; seller_sku futuro ligado ao campo errado.
- Must preserve: gate determinístico research §5 (2 âncoras independentes, contradição vence EAN, hard negatives kit/cor/medida/voltagem); decisões ACCEPT/REVIEW/REJECT/NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET.
- Trade-off: nenhum.
- Validation impact: testes de colisão reais (Doka/Menegotti, Doka/VW) no resolver.

### ADR-04: Modelo de dados de preço no módulo `market` (research §6)
- Decisão: `market` deixa de ser contract-only e ganha persistência: MarketPriceSnapshot, ValidatedOffer, MarketPriceAggregate, CompetitiveSignal (+ vínculo a identidade catalog). Agregados: só BRL, condição new, identidade ACCEPT, dedupe por seller (menor oferta válida), expor n_offers/n_sellers/idade.
- Prevents: guardar só "market_price" sem origem/evidência.
- Must preserve: ADR-17 — snapshot válido nunca sobrescrito por zero/falha; fetched_at/erro/expiração/fonte registrados.
- Trade-off: mais tabelas que o mínimo da demo; é o modelo já fechado pela pesquisa.
- Validation impact: teste negativo obrigatório (falha de coleta não zera snapshot).

### ADR-05: Extensões do adapter ML oficial — um dono
- Decisão: connectors/mercado_livre ganha: /items/{id}/sale_price, /items/{id}/price_to_win?version=v2, /products/search (por EAN), /products/{id} (buy_box só se preenchido), /products/{id}/items FLAG-GATED (paginação completa + telemetria + fallback explícito), /shipments/* (SLA/custos), /users/{id}/shipping_options/free. Proibido: /sites/MLB/search como dependência, scraping, provedor não homologado.
- Prevents: fundação em rota marcada p/ desligamento; 403 conhecido do /sites/MLB/search.
- Must preserve: payloads de provider morrem no adapter (MPC-native).
- Trade-off: rota de ofertas útil hoje pode sumir — por isso flag+fallback.
- Validation impact: telemetria da rota flag-gated visível; lane live-provider-read.

### ADR-06: Pré-anúncio honesto = identidade + estados (sem promessa de preço)
- Decisão: veredicto "vale a pena vender?" emite ACCEPT/REVIEW/REJECT/NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET; pós-anúncio usa sale_price+price_to_win como sinal primário; copy nunca promete "menor preço"/"preço de mercado automático" pré-anúncio.
- Prevents: reabrir NO-GOs fechados (SearchAPI/Gecko/Pricefy/scraping).
- Must preserve: evidência de coleta visível em toda UI de preço (target P1c).
- Trade-off: demo pode mostrar NO_PRICE_EVIDENCE em produtos do cliente (buy_box veio null 22/22) — runbook de demo mistura produtos com anúncio próprio ativo (sinal forte) e pré-anúncio honesto.
- Validation impact: QA verifica estados negativos renderizados, não escondidos.

### ADR-07: FE — retheme/shell primeiro; nav canônica
- Decisão: milestone FE inicial reconstrói shell (tokens papel+verde, Instrument Sans/Plex Mono, data-theme light/dark, PT-BR, nav pills Visão geral·Anúncios·Mercado·Simulador·Pedidos·Repasses·⚙; Mercado/Repasses desabilitadas "em breve"; Vínculos fora da nav global). Telas novas nascem no shell novo; mocks são referência, defeitos conhecidos não se reproduzem (DIFAL SP-hardcoded, Dashboard stale, contagens fake).
- Prevents: dois visuais coexistindo; colisão múltipla no seam AppRouter/Layout.
- Must preserve: decisões de design HANDOFF (sem edição inline, kanban read-only, config 3 camadas, drawers).
- Trade-off: telas esperam shell (~metade de wave A).
- Validation impact: QA visual light/dark por tela.

### ADR-08: Zero writes ML no MIS-004; fila com execução desabilitada
- Decisão: toda mutação (aplicar preço, vincular em massa, pausar) entra na fila M-03 com preview+protocolo; dispatcher para ML DESABILITADO por flag na demo/MIS-004. Vínculo é tabela nossa (nunca toca ML).
- Prevents: write acidental em anúncio de produção na frente do cliente.
- Must preserve: gates provider-write do lane (7 gates) intactos p/ quando ligar.
- Trade-off: "aplicar" não conclui de verdade na demo — protocolo mostra estado "na fila".
- Validation impact: QA prova que NENHUMA chamada de escrita saiu (telemetria/log do adapter).

### ADR-09: Polling/GET only no MVP
- Decisão: sem webhooks no MIS-004; refresh por GET/polling on-demand. Webhooks completos = MIS-005 M-06.
- Prevents: endpoint público/tunnel + async em 3 dias.
- Trade-off: dados "ao vivo" exigem refresh manual/agendado local.
- Validation impact: refresh visível na UI (idade da coleta).

### ADR-10: DIFAL fonte única mínima
- Decisão: tabela 27 UFs seed "padrão 2026" + overrides (mapa esparso) persiste no módulo `pricing` (dono dos parâmetros de cálculo no MVP); Simulador aplica DIFAL pelo destino REAL selecionado; Pedidos mostra chip informativo por pedido (UF do shipment × tabela). SEM agendamento/lembrete/marcar-pago.
- Prevents: 3 implementações divergentes do mock; SP hardcoded.
- Must preserve: interna−interestadual=DIFAL efetivo; origem SC; 12% Sul/Sudeste, 7% resto.
- Trade-off: Configurações completa fica MIS-005; edição via drawer do Simulador só.
- Validation impact: teste com UF de exceção ajustada refletindo em Simulador e chip Pedidos.

### ADR-11: Coleta on-demand; runtime demo local
- Decisão: sinais/snapshots coletados on-demand (botão/CLI na prep da demo); sem scheduler no MVP (scheduler diário = MIS-005 M-01, ligar CEDO). Demo roda docker dev stack local.
- Prevents: infra nova no prazo.
- Trade-off: variação 7d/histórico indisponíveis (honesto: "sem histórico ainda").
- Validation impact: rehearsal completo da demo no stack local = QA da missão.

### ADR-12: sdk-runtime é client manual — sync same-commit
- Decisão: registrar que packages/sdk-runtime NÃO tem codegen; toda mudança OpenAPI vem com edição manual do SDK no mesmo commit (regra profile mantida); adotar codegen fica fora das duas missões (backlog).
- Prevents: "regenerar" fantasma nos briefs.
- Trade-off: sync manual propenso a drift — mitigado por parity check L2.
- Validation impact: parity OpenAPI↔SDK↔handler no smoke L2.

## MIS-004 — Milestone headlines (ordem + dependências)

| ID | Headline | Módulos/seams (dono) | Depende |
|---|---|---|---|
| M-01 erp-xlsx-identity | Contrato de identidade corrigido (ADR-03) + módulo erp_import (.xlsx → snapshot ERP + Reader adapter + protocolo) | internal_read, catalog(identity), erp_import(novo), OpenAPI seção erp-import, migrações bloco A | merge W1 |
| M-02 price-intel-core | Modelo de preço (ADR-04) + Identity Resolver determinístico + extensões adapter ML (ADR-05) + agregador + motor de veredicto (ADR-06) | market, connectors/mercado_livre, OpenAPI seção market, migrações bloco B | merge W1; edge de feature: resolver depende do IC-01 (M-01) |
| M-03 shell-retheme | Shell novo (ADR-07): tokens, tema, nav, Layout/AppRouter PT-BR | apps/web shell seam (AppRouter/Layout/nav/theme) | merge W1 (M-02 chip) |
| M-04 vinculos-import-ui | Tela Vínculos & Importação completa + gaps product_links (confiança, lote+protocolo, undo) | product_links, /product-links, FE feature-product-links + rota | M-01, M-03 |
| M-05 anuncios-sinais | Anúncios retheme + sinais competitivos (nosso preço, vencedor, alvo) + exceções | listings, /listings, FE feature-listings | M-02, M-03 |
| M-06 produto-detalhe | Header + box "Vale a pena vender?" + abas Anúncios vinculados e Estoque | FE feature-product-detail(nova) + rota; composição client-side de APIs existentes | M-01, M-02, M-04, M-03 |
| M-07 simulador | Simulador completo: decomposição real, params globais + DIFAL seed (ADR-10), cenários, aplicar→fila (execução off, ADR-08) | pricing, /pricing, FE feature-simulator, migrações bloco C | M-02, M-03; envelope M-03(W1) na base |
| M-08 pedidos | Pedidos: KPIs + Fila/Lista/Kanban + drawer (decomposição, timeline, rastreio, comprador) + chip DIFAL | orders, /orders, FE feature-orders; connectors: shipments via ADDITIVE LOCK se M-02 aberto | M-03; edges: M-02(adapter), M-07(difal read) |
| M-09 dashboard-demo | Dashboard MPC (agregações reais) + demo hardening (runbook, coleta on-demand, rehearsal) — CORTÁVEL | dashboard, /dashboard, FE Dashboard | resto; última |

Waves: **A** = M-01 ∥ M-02 ∥ M-03 · **B** = M-04 ∥ M-05 ∥ M-07 ∥ M-08 · **C** = M-06 ∥ M-09.
Colisões pré-resolvidas: OpenAPI+sdk-runtime = contract lock com seções disjuntas por milestone (hub); migrações = blocos disjuntos (A: 0045-0049, B: 0050-0054, C: 0055-0059); connectors = M-02 dono, M-08 additive lock p/ shipments; FE = shell só M-03, depois rotas disjuntas por tela.

## MIS-005 — Milestone headlines (grão milestone; P4-P7 antes da execução)

| ID | Headline | Nota |
|---|---|---|
| M-01 mercado-radar | Reprecificação/Oportunidades/Monitorados + scheduler snapshots diários | ligar coleta CEDO (sem retroativo) |
| M-02 repasses | MP release report assíncrono/CSV ingestão agendada + conciliação ERP + calendário | |
| M-03 difal-config-completo | Agendamento/lembretes/marcar pago + exceções + tela Configurações completa | migra ownership de params se preciso |
| M-04 pedidos-pos-venda | Claims/devoluções reais, reputação has_incentive, reversa, contestar, NF-e | depois/junto M-06 (eventos) |
| M-05 produto-detalhe-completo | Abas Concorrência, Pedidos, Histórico 90d, auditoria, Dados | Histórico exige M-01 (snapshots) |
| M-06 webhooks-eventos | Topics orders_v2/shipments/claims/items/payments/item competition/fbm stock + GET reconciliação | substitui polling |
| M-07 provedor-externo | DataForSEO canário → Precifica contrato → batch50 re-run | gated research §9-10; nunca antes |
| M-08 full-visits-benchmarks | Estoque Full (inventories), visits, items_to_win/price references, etiquetas em massa | |
| M-09 auth-multitenancy | Auth middleware, tenant real por request, RBAC mínimo, CORS prod | baseline: hoje single-tenant sem auth |

## Top risks

| Risco | L | I | Mitigação |
|---|---|---|---|
| Merge W1 atrasa → execução bloqueada | M | A | hub prioriza; planning independe; chips prontos p/ dispatch |
| 3 dias vs 9 milestones full-doctrine | M | A | waves paralelas; M-09 cortável; M-06 abas mínimas; gates são horas |
| Pré-anúncio devolve NO_PRICE_EVIDENCE em massa (buy_box null 22/22) | A | M | runbook demo mistura anúncios próprios ativos (sinal forte) + pré-anúncio honesto; copy ADR-06 |
| xlsx real do cliente com colunas imprevistas | M | M | IC-02 + template exemplo + validação com relatório de rejeição |
| Colisão connectors (M-02×M-08) | M | B | additive lock nomeado no plano |
| Rota products/{id}/items desliga | B | M | flag+fallback+telemetria (ADR-05) |
