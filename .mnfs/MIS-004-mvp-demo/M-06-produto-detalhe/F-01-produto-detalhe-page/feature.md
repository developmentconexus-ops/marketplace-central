# F-01-produto-detalhe-page

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-06-produto-detalhe.

## Brief

Página `/catalogo/produtos/:productId` real (substitui placeholder W1 via `routes/produto.tsx`): header de identidade, box VEREDICTO, tab Anúncios vinculados, tab Estoque. Composição client-side pura de APIs existentes — ZERO endpoint/SDK novo. Cena central da demo.

## Inputs

- Design handoff (tela Produto) + leitura R-02 `research/design-screens-2026-07-17.md`.
- Clients SDK existentes: catalog (identidade/custo — M-01), `market.ts` (verdicts/aggregates — M-02), `productLinks.ts` (vínculos — M-04), `listings.ts` (by-product + sinais — M-05).
- IC-01 (chips identidade/completude), IC-03 (estados de veredicto + obrigação de citar evidência), IC-05 (tokens/primitivas/seams).

## Expected Output

- Header: descrição, CODPROD, EAN (ou chip "sem EAN" + refforn), marca, NCM, custo ERP (valor + `imported_at` como freshness), completude de identidade (chips por campo).
- Box VEREDICTO: label (`saudavel|viavel_preco_mercado|apertado|nao_vale`), faixa de preço sugerida, decomposição resumida, evidência citada (N sellers, fontes, fetched_at), botão "coletar agora" (POST /market/collections + refetch).
- Tab Anúncios vinculados: tabela dos anúncios do produto (by-product) com sinal por anúncio; vazio ⇒ EmptyState com link /vinculos.
- Tab Estoque: KPIs FÍSICO / RESERVADO / DISPONÍVEL (Reader via API catalog; disponível desconhecido (reservado ausente no snapshot — IC-02) ⇒ card "—") + banner de cobertura (dias, se vendas disponíveis via orders; sem dado ⇒ banner oculto).
- EARS: While produto tem veredicto com evidência, when página carrega, the box shall mostrar label + faixa + fontes com timestamps. While botão coletar clicado, when POST síncrono retorna 200 com sumário, the box shall refetch veredicto/sinais imediatamente e renderizar o resultado (estado "coletando" dura só o request em voo — sem polling). While produto inexistente, when rota carrega, the página shall mostrar ErrorState 404 com link ao catálogo.

## Negative Scenarios (matriz de estados compostos)

- Sem EAN ⇒ header chip "sem EAN — matching limitado a REVIEW" + veredicto NO_PRICE_EVIDENCE com call-to-action /vinculos.
- EAN ok, sem vínculo, sem coleta ⇒ veredicto NO_PRICE_EVIDENCE + botão coletar em destaque.
- Coleta feita, <5 sellers ⇒ INSUFFICIENT_MARKET com contagem exata ("3 vendedores — mínimo 5").
- Custo desconhecido ⇒ faixa de mercado sem label de margem (`blocking_state: SEM_CUSTO` do IC-03), custo "—".
- Snapshot ERP ausente ⇒ Estoque tab inteira em estado "importar planilha" (link /estoque ou import).
- Falha parcial (market fora, catalog ok) ⇒ header renderiza, box em ErrorState isolado com retry.

## Deliberate Exclusions

Tabela de movimentações, KPI Full ML, tabs Concorrência/Pedidos/Histórico/Dados ⇒ MIS-005 (sem fonte de dados no MVP). Não renderizar tabs vazias.

## State / Interaction Model

- Tab ativa na URL (`?tab=`); deep-link + F5 restauram.
- Queries independentes por widget (header/veredicto/anúncios/estoque) — falha de uma não derruba as outras (Suspense/erro por seção).
- "Coletar agora": mutation síncrona (POST /market/collections, 200 com sumário) ⇒ invalidate+refetch; botão disabled enquanto request em voo; erro/timeout do request ⇒ toast com causa (sem polling, sem query de collection).
- Keys: `['catalog','product',id]`, `['market','verdict',id]`, `['listings','by-product',id]`, invalidação pós-coleta nos dois últimos.

## Constraints

- ZERO backend novo — qualquer gap de API descoberto = ESCALATION ao hub (provável dono: milestone da API), nunca endpoint improvisado.
- Estados honestos obrigatórios (ADR-17): nenhum "0", nenhum placeholder numérico.
- Só SDK clients; tokens M-03.

## Ownership

- Owned paths: `apps/web/src/pages/produto/**` (novo), `apps/web/src/routes/produto.tsx` (conteúdo).
- Forbidden paths: `apps/web/src/app/**`, `packages/ui/**`, outros routes/pages, `sdk-runtime/**`, qualquer backend.
- Parallel-safe with: M-09 F-01 (disjoint: rotas/páginas distintas) — dentro do M-06, feature única.

## Validation Expectations

- Screenshots dos 4 estados de veredicto (saudável/…/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET com contagem) + estado sem-EAN.
- Transcript "coletar agora": clique ⇒ 200 síncrono com sumário ⇒ box atualiza com evidência nova (fetched_at maior).
- Deep-link `/catalogo/produtos/X?tab=estoque` + F5 ⇒ tab correta ativa.
- KPI estoque batendo com valor do snapshot importado (mesmo número do transcript M-01).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-06, wave C pós-merge wave B).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots da matriz de estados.
- Blockers or open decisions: none.
