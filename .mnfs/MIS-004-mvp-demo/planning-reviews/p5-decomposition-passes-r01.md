# P5 Decomposition Passes — r01

```yaml
id: P5-PASSES-R01
type: planning-review
status: complete
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
lifecycle_scope: support
```

Três passes obrigatórios do P5 (HARNESS/mission-planning), executados ANTES do audit Sol.
Autor: sessão de planning (Claude). Evidência de código citada com caminho:linha.

## 1. Feature-level Write-DAG

20 features. Write-sets declarados em cada brief (§Ownership). Sobreposições e resolução:

| Par | Sobreposição | Resolução |
| --- | --- | --- |
| M-01 F-02 × F-03 | módulo `erp_import` (F-02 tudo exceto `adapter/`; F-03 só `adapter/` + root + modules.json); fixture `runner_test.go` (ambos bump) | edge serial F-02→F-03 já no DAG interno (F-03 depende das tabelas) — sem lock extra |
| M-01 F-01 × F-02 | nenhuma (módulos distintos; F-02 consome validador GTIN via import — assinatura IC-01) | paralelo confirmado |
| M-02 F-02 × F-03 × F-04 | módulo `market` (subtrees distintos, mas mesmo módulo) | cadeia serial F-02→F-03→F-04 já no DAG interno |
| M-03 F-01 × F-02 | `apps/web/src/app/**` (F-01 só `app/theme/**` + index.css; F-02 Layout/AppRouter/routes) | subtrees disjuntos + edge F-01→F-02 já existente |
| M-0{4,5,7,8} F-01 × F-02 | nenhuma (backend/OpenAPI/SDK vs pages/routes) | edge serial F-01→F-02 por dependência de contrato (OpenAPI committed), não por write-overlap |
| Cross-milestone wave B (M-04∥M-05∥M-07∥M-08) | ZERO em módulos, seções OpenAPI, arquivos SDK, blocos de migration, rotas FE, páginas, tabelas | disjunção seis-eixos confirmada (matriz de ownership da missão) |
| Cross-milestone wave C (M-06∥M-09) | zero (pages/rotas/módulos distintos; M-06 não tem backend) | paralelo confirmado |

Seams compartilhados (hub adjudica no merge, já registrados na missão §Hub merge seams):
`sdk-runtime/src/index.ts` (barrel, 1 linha export por milestone), composition root (registro
aditivo), `contracts/governance/modules.json` (entry por chip via merge), fixture de contagem
de migrations `runner_test.go` (bump numérico — conflito trivialmente resolvível), arquivo
OpenAPI raiz (seções distribuídas por milestone; colisão estrutural = hub).

Additive-lock grants (pré-autorizados na matriz da missão):
- **IC-04 `DifalForUF` + `Decompose`**: M-07 publica, M-08 consome read-only. Lock: M-07 NÃO
  altera assinatura pós-publicação (congelada no IC-04); liberado no close do M-08.
- **Envelope mutations**: M-04 (tipo `link_apply`) e M-07 (tipo `price_update`) consomem via
  API pública `/mutations*`; `modules/mutations/**` é forbidden path para AMBOS. Nenhum tipo
  novo necessário (ver §3).

Veredicto: nenhum write-overlap sem edge serial ou grant. PASS.

## 2. Contract-Satisfiability

Promessas dos briefs × conjunto ADR/IC ratificado:

| Promessa | Constraints em jogo | Veredicto |
| --- | --- | --- |
| M-04 batch approve via envelope local, zero ML | enum fechado `ProtocolType` + zero-writes-ML + forbidden `modules/mutations` | SATISFEITO por reuso: `link_apply` existe, habilitado, e seu executor é LOCAL (`applyLink` → linkage writer, ações `approve_candidate/manual_resolve/reject_listing` — `mutations/application/writer.go:116-128`). Brief M-04 F-01 corrigido nesta rodada p/ citar o tipo existente. |
| M-07 "aplicar preço" sem write ML | executor de `price_update` ESCREVE ML (`adapters/connectors/price_writer.go`) × zero-writes | SATISFEITO por gate de estado: protocolo para em `previewed` (transição `previewed→approved` nunca disparada pela UI; runbook proíbe approve na demo). Brief M-07 F-02 corrigido p/ fixar teto `previewed`. Residual risk aceito e registrado (R10 abaixo). |
| Verdict sem custo ERP | IC-03 verdict_label × ADR-17 unknown | SATISFEITO por extensão aditiva do IC-03: `blocking_state: SEM_CUSTO` (editado r01, consumido por M-02 F-04 / M-06). |
| M-05 `abaixo_custo` exclui custo desconhecido | ADR-17 (desconhecido ≠ violação) | consistente por construção nos briefs. |
| M-07 solver margem→preço com degrau em 79 | fórmula IC-04 com taxa fixa <79 / frete ≥79 | satisfazível (busca em domínio com degraus; estado `NO_SOLUTION` contratado p/ os casos sem solução — sem promessa de solução universal). |
| M-09 zero real vs fonte indisponível | ADR-17 | dois estados distintos contratados no brief (`0` com janela vs `null` + reason). |
| M-05/M-06/M-08 leituras cross-módulo | boundary doctrine (ports/APIs públicos, sem SQL alheio) | todos os briefs proíbem query direta; portas nomeadas existem ou são criadas por feature nomeada (§3). |

Nenhum conflito não-resolvido. Resoluções ratificadas NOS artefatos (briefs + IC-03), não empurradas p/ implementação. PASS.

## 3. Prerequisite-Existence

Verificações contra o código (tool runs desta sessão, 2026-07-17, main ≥ f4612be3):

| Assunção | Evidência | Status |
| --- | --- | --- |
| `/mutations*` completo (create/preview/approve/cancel/items/retry) | inventário OpenAPI + `modules/mutations/**` presente | EXISTE |
| Tipo `link_apply` habilitado, executor local | `mutations/domain/protocol.go:37,47`; `application/writer.go:116-128` (linkage writer, não ML) | EXISTE |
| Tipo `price_update` habilitado (executor = ML write) | `protocol.go:35,45`; `adapters/connectors/price_writer.go` | EXISTE (uso gated a `previewed`) |
| `/product-links/*` (candidates/generations/resolutions approve-candidate/manual-resolve/reject-listing/workflows) | inventário OpenAPI | EXISTE |
| `/listings*` (list/by-product/summary/refresh) | inventário OpenAPI | EXISTE |
| `/orders*` (list/import/sankhya-linkage) | inventário OpenAPI | EXISTE |
| `/pricing/simulations` + batch | inventário OpenAPI | EXISTE |
| `/dashboard/summary` | inventário OpenAPI | EXISTE |
| `/market/observations|references` (SAT, intocados) | inventário OpenAPI | EXISTE |
| 15 módulos servidor (incl. mutations, market, listings, orders, pricing, product_links, dashboard, catalog, internal_read, connectors, inventory) | listagem `apps/server_core/internal/modules` | EXISTE |
| `sdk-runtime/src` = `index.ts` + test (flat) → padrão "arquivo novo por milestone + export no barrel" | listagem do diretório | CONFIRMADO (arquivos por milestone = to-be-created pelos próprios chips) |
| Rotas PT-BR + placeholders (/vinculos, /pedidos, /catalogo/produtos/:id, /integracoes) + LegacyRedirects + InstallationContext + primitivas ui + TanStack via web-query | R-03 (`research/w1-merge-addendum-2026-07-17.md`) | EXISTE |
| Migrations: top 0044, fixture=41, blocos 0045+ livres | R-03 | CONFIRMADO |
| Reader ports (`GetCostAsOf`, `GetSellableStock`) em internal_read | R-01/R-03 (uso atual Oracle) | EXISTE (M-01 F-03 conforma adapter) |
| Validador GTIN | não existe — TO-BE-CREATED por M-01 F-01 (nomeado) | ATRIBUÍDO |
| Módulo `erp_import`, tabelas `erp_import_*`, `/erp/imports*` | não existem — TO-BE-CREATED por M-01 F-02 (nomeado) | ATRIBUÍDO |
| Ports IC-06 no connectors | não existem — TO-BE-CREATED por M-02 F-01 (nomeado) | ATRIBUÍDO |
| Tabelas/ports market IC-03, resolver | TO-BE-CREATED por M-02 F-02/F-03/F-04 (nomeados) | ATRIBUÍDO |
| `routes/<area>.tsx` seam | TO-BE-CREATED por M-03 F-02 (nomeado; IC-05) | ATRIBUÍDO |
| `MarginChip`, DataTable/DetailDrawer | TO-BE-CREATED por M-03 F-03 (nomeado) | ATRIBUÍDO |
| Motor `Decompose` + `DifalForUF` + tabelas pricing | TO-BE-CREATED por M-07 F-01 (nomeado; IC-04) | ATRIBUÍDO |
| Lib parser .xlsx no repo | NÃO VERIFICADO nesta sessão — brief M-01 F-02 manda verificar primeiro e `REQUEST` dep ao hub se ausente (protocolo de dependência da doutrina) | TRATADO NO BRIEF |

Nenhuma assunção órfã (★5): tudo verificado ou atribuído a feature criadora nomeada. PASS.

## 4. Riscos novos ratificados nesta rodada

- **R10 (residual, aceito)**: tipo `price_update` tem executor de write ML já habilitado; o
  gate anti-write da demo é comportamental (UI não oferece approve + runbook proíbe). Registrado
  na missão §Risks; MIS-005 M-10 é o dono da habilitação governada.

## 5. Mudanças em artefatos nesta rodada (pós-briefs, pré-manifest)

- IC-03: enum `blocking_state` explicitado + `SEM_CUSTO` (aditivo).
- M-04 F-01: tipo `link_apply` reusado (era tipo novo inventado — defeito corrigido).
- M-07 F-02: teto `previewed` fixado (era "PENDING" vago).
- M-01 F-03: semântica `reservado_desconhecido` desfeita de hedge.
