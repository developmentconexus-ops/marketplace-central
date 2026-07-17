# REPLAN BRIEF — MVP demo (segunda 2026-07-20) + produto completo

> Autor: hub MIS-003, 2026-07-17. Consumidor: **sessão nova de mission planning (do zero)**.
> Status: proposta de escopo alinhada com operador; planning session produz o plano; operador aprova antes de qualquer implementação.

## 1. Contexto e gatilho

O operador entregou o pacote de design **high-fidelity final** (não wireframe): 9 telas + mapa
UI×API + decisões de produto ratificadas. Isso supersede o planejamento de telas do MIS-003
(W2 M-04, W3 M-05-FE/M-06 cancelados). Há uma **apresentação de MVP a cliente na segunda-feira
2026-07-20**; a feature central pedida pelo cliente é: *simulador + produtos em estoque
(cadastrados no ERP) que podem ser vendidos por bons preços comparado ao mercado ML*.

O trabalho vira (pelo menos) **duas missões**: MIS-004 (MVP demo) e MIS-005 (produto completo).
Planejar as duas; detalhar MIS-004; MIS-005 pode ficar em grão de milestone.

## 2. Inputs obrigatórios (ler nesta ordem)

1. `docs/design/handoff-2026-07/README.md` — overview + tokens + fidelity.
2. `docs/design/handoff-2026-07/API-MAP.md` — **binding**: cada bloco de cada tela classificado ✅ API direta / ⚠️ derivado / ❌ domínio nosso. Nada na UI foi inventado fora desse mapa.
3. `docs/design/handoff-2026-07/HANDOFF.md` — decisões de produto/design (não regredir: sem edição inline, kanban read-only, config 3 camadas global→produto→anúncio, fila de sync com preview+protocolo para TODA mutação).
4. As 9 telas `docs/design/handoff-2026-07/*.dc.html` — referência de design; **recriar no stack do repo, não copiar HTML**. Markup em `<x-dc>`, dados de exemplo no `<script data-dc-script>`.
5. `docs/research/2026-07-17-pricing-intelligence-implementation-handoff.md` — **binding para tudo de preço/mercado** (+ os 3 docs de autoridade que ele referencia).
6. Doutrina: `docs/HARNESS-CORE.md` + `docs/HARNESS-PROFILE.md`, `ARCHITECTURE.md`/ADRs (ADR-17!), OpenAPI + `sdk-runtime`, `contracts/governance/`.
7. Estado atual: MIS-003 W1 (M-02 plataforma FE + Anúncios; M-03 mutation envelope; SAT contratos) — em fechamento, merge iminente. **Construir em cima, não refazer.**

## 3. Régua de honestidade de preço (fechada pela pesquisa — não reabrir)

- **Pós-anúncio = GO**: item próprio + `sale_price` + `price_to_win` (10/10 provado). Sinal primário para anúncios ativos.
- **Pré-anúncio = GO só para identidade + estados honestos**: descoberta por catálogo oficial (`products/search` por EAN → `products/{id}`); `buy_box_winner*` consumir só se preenchido (veio null 22/22); resposta pode e deve ser `ACCEPT / REVIEW / REJECT / NO_CANDIDATE / NO_PRICE_EVIDENCE / INSUFFICIENT_MARKET`.
- **`products/{id}/items` (distribuição de ofertas)**: funciona hoje (236/256), rota marcada p/ desligamento na doc EN — implementar SÓ atrás de feature flag + paginação completa + telemetria + fallback explícito. Nunca fundação única.
- **Proibido**: scraping, busca pública, proxies, SearchAPI/Gecko/Pricefy (NO-GO fechado), provedor externo sem contrato+canário (§9-10 do research). `/sites/MLB/search` deu 403 com nossa app — não depender dele.
- Matching: EAN não é verdade suficiente (colisões reais provadas). Gate determinístico §5 do research: 2 âncoras independentes, contradição vence EAN, hard negatives (kit, cor, medida, voltagem...).
- Identidade ERP: `CODPROD` (SKU canônico) ≠ `REFERENCIA` (EAN) ≠ `REFFORN` (ref. fabricante) — corrigir contrato antes de ligar `seller_sku`.
- Modelo de dados de preço: persistir origem+evidência (`ProductIdentity`, `MarketPriceSnapshot`, `CompetitiveSignal`, `ValidatedOffer`, `MarketPriceAggregate` — §6). Nunca sobrescrever snapshot válido com zero/falha (ADR-17).
- Copy de venda: "monitoramos posicionamento competitivo e preço-alvo dos anúncios ativos, com sinais oficiais e evidência de coleta" — NÃO prometer "menor preço"/"preço de mercado automático" pré-anúncio.

## 4. Proposta de corte MVP (MIS-004) vs completo (MIS-005)

Critério: história da demo = *"do estoque ERP ao veredicto 'vale a pena vender?' com preço vs mercado honesto, e simulador de margem real"* + frontend novo + Pedidos funcional. Sem levar em conta tempo humano — paralelização IA; mas cortar pelo risco (dependências externas assíncronas, contratos, features sem fonte de dado provada).

### MIS-004 — MVP demo (dentro: tudo que a história precisa)

| # | Item | Base |
|---|---|---|
| 1 | **Contrato de identidade** CODPROD/EAN/REFFORN + Identity Resolver determinístico + testes de colisão | research §4-5, §11.1 |
| 2 | **Modelo de dados de preço** (snapshot/oferta validada/agregado/sinal competitivo/auditoria) | research §6 |
| 3 | **Adapter ML oficial**: itens próprios + sale_price + price_to_win; descoberta catálogo por EAN; ofertas de catálogo flag-gated | research §8 |
| 4 | **Shell novo design** (tokens papel+verde, Instrument Sans/Plex Mono, nav pills, light/dark) — retheme da plataforma M-02 | README tokens |
| 5 | **Vínculos e Importação** (tela) — matching com chips de confiança, drawer manual, lote+protocolo; vincular não altera ML | tela 7; é pré-requisito da história |
| 6 | **Anúncios** (tela, retheme + sinais competitivos: nosso preço, vencedor, alvo) | tela 2; M-02 base |
| 7 | **Produto Detalhe parcial**: header + box **"Vale a pena vender?"** (a feature pedida pelo cliente) + aba Anúncios vinculados + aba Estoque (físico/reservado/disponível ERP) | tela 3 |
| 8 | **Simulador** completo: matriz + painel preço↔margem bidirecional + decomposição real via `listing_prices`/`shipping_options` (nunca hardcodar) + drawer parâmetros + cenários; "aplicar" via fila de sync (M-03) | tela 5 |
| 9 | **Pedidos**: KPIs + Fila/Lista/Kanban read-only + drawer detalhe (decomposição, timeline, rastreio, comprador). SEM claims/devoluções, SEM DIFAL completo (chip/coluna pode aparecer com dado da tabela seed se barato) | tela 6 |
| 10 | **Dashboard** (agregações dos dados já trazidos) — última prioridade, corta se apertar | tela 1 |

Estados honestos em TODA a UI de preço: chip de match, `n_offers/n_sellers`, idade da coleta, `NO_PRICE_EVIDENCE`/`REVIEW` visíveis. Demo com dados reais do seller (leitura); **zero writes live no ML sem autorização explícita do operador**.

### MIS-005 — produto completo (fora do MVP)

- **Mercado** (radar completo: reprecificação, oportunidades, monitorados + scheduler de snapshots diários — ligar coleta CEDO, sem retroativo).
- **Repasses** (release report MP assíncrono/CSV — ingestão agendada, conciliação ERP, calendário).
- **DIFAL completo** (tabela 27 UFs seed "padrão 2026", exceções, agendamento, lembretes, marcar pago) + **Configurações** completa.
- **Pedidos**: Cancelados/Devoluções reais (claims, reputação `has_incentive`, reversa, contestar), NF-e/faturamento.
- **Produto Detalhe**: abas Concorrência, Pedidos, Histórico 90d (exige snapshots), auditoria, Dados.
- **Webhooks completos** (`orders_v2`, `shipments`, `claims`, `items`, `payments`, `item competition`, `fbm stock operations`) como fonte de eventos + GET reconciliação. (No MVP, polling/GET basta — planning decide o mínimo.)
- **Provedor externo** de preço (DataForSEO canário → Precifica contrato → batch50 re-run) — só após gates §9-10.
- Estoque Full (`inventories`), visits, benchmarks (`items_to_win`, price references), etiquetas em massa.

**Decisões em aberto para o planning** (não decidir sozinho, propor): 1) webhooks mínimos vs polling no MVP; 2) DIFAL chip-only no MVP ou 100% fora; 3) Dashboard dentro/fora; 4) retheme M-02 como primeiro milestone vs shell paralelo; 5) grão de milestones/chips e matriz de colisão.

## 5. Restrições de arquitetura (repo)

- **MPC-native, não Mercado Livre**: payloads de provider morrem nos adapters; domínio fala anúncio/vínculo/produto/pedido MPC. API-MAP referencia endpoints ML — isso é fonte do ADAPTER, não shape do domínio.
- Camadas domain/application/ports/adapters/transport; tenant scoping `tenant_id`; ADR-17 (unknown ≠ 0/default) — vale para todo ⚠️ derivado.
- API change = OpenAPI + `sdk-runtime` juntos. Mutações = fila de sync com preview+protocolo (M-03 envelope). YAGNI: nada além do que as telas + API-MAP pedem.
- Provider writes: linkage resolvido, policy/source time explícitos, proteção de duplicata, audit.

## 6. Estado W1 (não bloquear planning; merge iminente)

- M-02 (shell FE, query state, Anúncios v1): corrective M02-COR-1 em voo → novo CLOSED → dual gate → merge. Base do retheme.
- M-03 (mutation envelope/fila de sync): F-04 gated no merge do M-02 → merge. Core do padrão de mutação do produto.
- CHIP-SAT (contratos backend M-05 F-01/M-06 F-02): corrective C1 → close.
- Planning assume main pós-W1 como base; se plano ficar pronto antes do merge, marcar dependência explícita.
