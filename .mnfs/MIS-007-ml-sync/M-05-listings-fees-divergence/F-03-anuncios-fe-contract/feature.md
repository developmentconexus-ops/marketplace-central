# F-03-anuncios-fe-contract

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-05 listings-fees-divergence.

## Brief

Contrato + FE: `ListingReadModel` (`read_model.go:122-149`) ganha campos ADITIVOS
`tarifa` (objeto: `amount` string decimal, `detail` canônico IC-01 COMPLETO {percentage_fee,
fixed_fee, financing_add_on_fee, price_used, listing_type_id — F-r06-4}
quando presente, `origem`, `coletado_em` — COMPOSTO de channel_fees camada 2 no read
service, IC-07: fee NUNCA vira coluna de listings; shape = IC-01 canonical camada 2) e
`divergences[]` (IC-02 §Required Outputs — projeção exata `{kind, expected_value,
observed_value, expected_observed_at, observed_observed_at, detected_at}` das rows abertas;
P7 r01 B-2); param
novo `filter.divergentes=true` (entra em `domain.FilterKeys`, `filter.go:9`). Query do repo
faz LEFT JOIN lateral em channel_fees/divergences (cuidado ORDER BY: memória
`pg-orderby-output-column-trap`). OpenAPI + SDK (tipos de `listListings` `index.ts:2130`)
MESMO commit. FE: AnunciosTable (`AnunciosTable.tsx:277-294`) ganha coluna `TARIFA` (amount
formatado R$; percentage_fee do detail como texto secundário quando presente; sem
observação → `—`) e badge de divergência na coluna PENDÊNCIA existente (ou badge própria —
spec decide contra DESIGN-REFERENCE); filtro divergentes no anunciosQueryState;
ListingDetailPanel mostra os dois lados + timestamps da divergência
(expected_value/observed_value + expected_observed_at/observed_observed_at da projeção
IC-02 §Required Outputs).

EARS:
- While anúncio com fee camada 2, when /listings responde, the item shall carregar tarifa
  {amount, detail, origem, coletado_em}.
- While filter.divergentes=true, when lista pagina, the página shall conter SÓ itens com
  divergência aberta (fixture >1 página).
- While anúncio sem fee observada, when tela renderiza, the célula shall mostrar `—`
  (nunca 0, nunca R$79).

## Inputs

F-01 + F-02 (dados); IC-01/IC-02/IC-07 (shapes binding); fato #3 de
`research/p5-prerequisites.md` (envelope, params, sort,
colunas atuais verbatim); DESIGN-REFERENCE @8144238 (gate visual).

## Expected Output

Read service + repo + transport aditivos; par OpenAPI+SDK; AnunciosTable + panel + query
state.

## Constraints

- ADITIVO estrito — nenhum campo/param existente muda (baseline `read_model.go:122-149`
  congelado por golden test).
- Composição de tarifa no READ service (camada application) — transport não consulta
  channel_fees direto.
- Paginação: filtro novo respeita cursor existente (`DecodeListingCursor`) — página 1 com
  fixture >1 página (lição CHIP-MERCADO).

## Inputs/Outputs

DTO canonical: IC-02 §Required Outputs + IC-01 §Canonical Examples camada 2 (resolução em
cascata NÃO se aplica aqui — /anuncios mostra SÓ camada 2 observada; cascata é do
pricing/M-07).

## Interaction Model

- Filtro `divergentes` vive no anunciosQueryState (URL-sync como os filtros existentes);
  mudar filtro reseta cursor de paginação p/ página 1.
- Badge de divergência clica → abre ListingDetailPanel na seção de divergência (mesmo
  padrão de abertura do panel atual; sem rota nova).
- Refetch: dados de tarifa/divergência vêm na MESMA query de lista (sem query paralela);
  panel usa o item da lista + fetch de detalhe existente — stale segue a política do
  react-query atual do módulo (sem invalidação nova).
- Sem observação de fee: célula `—` estática (não loading, não 0).

## Negative Scenarios

- Divergência resolvida → some do array (só rows abertas).
- listing com variação divergente → badge no anúncio, panel discrimina a variação.

## Ownership

- Owned paths: `listings/transport/`, `listings/application/read_*`,
  `listings/adapters/postgres/repository.go` (query read aditiva), par OpenAPI+SDK
  /listings, `apps/web/src/pages/AnunciosPage.tsx` (fiação do filtro/badge na página),
  `AnunciosTable.tsx`, `anunciosQueries.ts`,
  `anunciosQueryState.ts`, `ListingDetailPanel.tsx`.
- Forbidden paths: ingest (F-01/F-02); schema; outras rotas FE.
- Parallel-safe with: none — depends on F-01 + F-02.

## Validation Expectations

- Golden: item SEM dados novos → JSON byte-igual ao baseline.
- Fixture >1 página com divergentes espalhados → filtro correto nas 2 páginas.
- tsc verde; live-drive do hub: tarifa real de anúncio real na tela (amount/detail
  batem com o observado no fixture de ingest).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01+F-02.
- Required files/evidence: `validation.md`; screenshot-métrica do live-drive (hub).
- Blockers or open decisions: none.
