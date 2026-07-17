# Research Note

```yaml
id: R-03
type: research
status: draft
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Addendum pós-merge W1: deltas entre o baseline R-01 (main @ cd74b401, chips não mergeados) e main atual após merges M-02 (79d6787f) e M-03 (f4612be3). R-01 permanece intacto (manifest P3 congelado sobre ele); este note é a evidência corrente para P4/P5.

## Findings

### Estado W1: RESOLVIDO

- `chip/m-02-frontend-platform-anuncios` mergeado @ 79d6787f (2026-07-17 11:32).
- `chip/m-03-mutation-envelope-writes` mergeado @ f4612be3 (com correctives COR-1..3 + QA-1 poller no server main).
- Precondição ADR-01 (base = main pós-merge W1) SATISFEITA. Chips MIS-004 partem de main ≥ f4612be3.

### FE atual (apps/web) — substitui seção "Frontend" do R-01

- **Rotas PT-BR canônicas JÁ registradas** (`AppRouter.tsx`): `/` (DashboardPage), `/anuncios` (AnunciosPage real), `/catalogo` (CatalogPage), `/catalogo/produtos/:productId` (PLACEHOLDER), `/vinculos` (PLACEHOLDER), `/estoque` (StockSeguroPage), `/precos` (PricingSimulatorPage), `/pedidos` (PLACEHOLDER), `/integracoes` (PLACEHOLDER), `/protocolos/:protocolId` (ProtocoloPage — UI do envelope M-03), `/classifications`, `/marketplaces` (PLACEHOLDER) + `LegacyRedirect` das rotas antigas (/orders→/pedidos etc.).
- **Anúncios workspace EXISTE**: `apps/web/src/pages/` AnunciosPage (323 l) + AnunciosTable + ListingDetailPanel + ListingsRefreshControl + ListingsSummary + anunciosQueries/QueryState. Padrão novo: telas moram em `apps/web/src/pages/`, NÃO em packages feature-*.
- **4 packages legacy REMOVIDOS** (COR-2, diretiva do operador): feature-integrations (IntegrationsHubPage), feature-marketplaces, feature-orders (OrdersPage), feature-product-links (ProductLinksPage) — páginas deletadas. Sobrevivem: feature-classifications, feature-simulator, feature-inventory, feature-products, ui, web-query.
- **InstallationContext global** (`app/InstallationContext.tsx`): provider + InstallationGate + seletor de installation no header. Telas consomem; nenhum milestone recria.
- **Primitivas ui novas**: UnknownValue, FreshnessIndicator, EmptyState/ErrorState/LoadingState, ConflictTag (alinhadas com ADR-17/observabilidade P1c).
- **web-query**: invalidation + failureCopy adicionados.
- **Tema: AINDA NÃO EXISTE.** Sidebar dark `#0F172A` hardcoded, Tailwind slate/blue, sem data-theme/tokens/fontes do design. Retheme M-03 (MIS-004) continua necessário integral.
- **Nav atual ≠ nav canônica do HANDOFF**: sidebar com 8 itens (Visão geral · Catálogo · Anúncios · Vínculos & Import. · Estoque · Preços & Simulador · Pedidos · Integrações & Sync). Canônica (HANDOFF, ratificada P1): pills no header Visão geral · Anúncios · Mercado("em breve") · Simulador · Pedidos · Repasses("em breve") · ⚙; Vínculos FORA da nav global. Reconciliação = escopo do M-03 shell-retheme.
- **UI de conexão OAuth REMOVIDA** com feature-integrations; `/integracoes` é placeholder. Endpoints OAuth do módulo integrations continuam no server. Demo depende de installation JÁ conectada no DB (assunção P1) ou conexão via API — entra no runbook.

### Backend/migrations — atualiza R-01

- Módulos governance: 15 (mutations REGISTRADO via merge M-03).
- Migrações: 41 arquivos; M-03 preencheu gaps 0038/0039 (mutation_protocols, inventory_stock_actions_mutation_protocol); topo permanece `0044_market_references.sql`; gaps restantes 0040-0042; fixture `runner_test.go:25,64` = 41. **Blocos 0045+ do plano MIS-004 continuam livres.**
- Envelope M-03 completo em main: /mutations OpenAPI + poller iniciado no server main (QA-1) + ProtocoloPage.

## Impact On Mission (deltas de plano — absorvidos na reconciliação P3 §5)

1. ADR-01: precondição satisfeita — vira registro de base SHA por chip.
2. M-03 shell-retheme: escopo REDUZ (rotas/PT-BR/placeholders já existem) e ganha: tokens papel+verde + data-theme + fontes + sidebar→pills canônicas + reconciliação de nav + indireção de rotas por área (para telas trocarem placeholder sem editar seam AppRouter).
3. M-04/M-08: páginas legacy não existem — telas novas nascem em `apps/web/src/pages/<área>/` (padrão M-02). Sem "retheme de página existente" para Vínculos/Pedidos.
4. M-05: estende AnunciosPage existente (não cria).
5. M-06: rota placeholder `/catalogo/produtos/:productId` já registrada.
6. M-07: PricingSimulatorPage vive em packages/feature-simulator (padrão legacy sobrevivente) — brief decide evoluir vs reconstruir em pages/; superfície do milestone = /precos + feature-simulator.
7. Runbook/risco: conexão de conta ML sem UI — garantir installation conectada antes da demo (DB existente ou API).
