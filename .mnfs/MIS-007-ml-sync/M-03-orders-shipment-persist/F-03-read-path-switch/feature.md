# F-03-read-path-switch

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-03 orders-shipment-persist.

## Brief

Replacement-before-deletion (ADR-05): enrich de `GET /orders/{id}` lê shipment + fiscal do
POSTGRES (rows do F-02) em vez dos readers vivos. `NewEnrichServiceWithReaders`
(root.go:597-599) troca `ordersShipmentReaderAdapter`/`ordersBuyerFiscalReaderAdapter`
(root.go:591-592) por readers de banco; sites vivos A (shipment) e B (fiscal) DELETADOS no
mesmo diff; allowlist do guard (M-02 F-04) encolhe 2 entradas — guard must-fail prova que
reintroduzir import vivo no módulo orders quebra build. DTO `comprador_fiscal` byte-idêntico
(shape `compradorFiscalDTO` `http_handler.go:462-467` + endereco `:446-454` — mapeamento
agora de colunas, mesmo JSON). Summary path (`order_repo.go:378`) passa a ler bucket
PERSISTIDO em vez de re-derivar com shipmentStatus="" — before/after das contagens
capturado (mudanças são correções, não regressões).

EARS:
- While pedido ingerido, when GET /orders/{id} responde, the payload shall vir 100% do
  Postgres (zero GET vivo ML — contador de adapter no teste = 0).
- While pedido ANTIGO sem row de shipment, when tela lê, the campos de shipment shall ser
  honest-unknown (null), nunca erro nem valor inventado.
- While summary agrega, when bucket é lido, the valor shall ser a coluna persistida.

## Inputs

F-02 (rows + colunas); `enrich_service.go:194-198` (`EnrichOne`; degrade documentado no
doc comment `:188-193` — sweep residual F-r08-3);
`PedidoDrawer.tsx:355-368` (renderer — campos exibidos: nome, doc_tipo+doc_numero,
endereco composto; uf_nome NÃO renderiza); M-02 F-04 (allowlist).

## Expected Output

Readers de banco (adapters postgres do módulo orders), enrich re-fiado, sites vivos
deletados, allowlist -2, before/after de summary.

## Constraints

- Contrato /orders INALTERADO exceto campos aditivos já no DTO — FE não muda (margem = M-06).
- Campos aditivos novos no DTO (pack_id etc.): par OpenAPI+SDK no MESMO commit (AGENTS).
- Deleção só DEPOIS do replacement no mesmo milestone — nunca janela sem caminho.

## Inputs/Outputs

In: rows de orders/order_shipments (F-02). Out: GET /orders/{id} com `comprador_fiscal`
BYTE-IDÊNTICO ao baseline (`compradorFiscalDTO` `http_handler.go:462-467` + endereco
`:446-454` — golden obrigatório); campos aditivos novos (pack_id etc.) shape IC-03
§Required Outputs;
list/summary: mesmo envelope, bucket da coluna persistida. Par OpenAPI+SDK mesmo commit.

## Negative Scenarios

- Backfill 12m ainda não rodou (M-06) → pedidos velhos ficam honest-unknown até lá; tela não
  pode quebrar com null (teste com row sem shipment).
- Import 202 concorrente com GET → leitura consistente (tx de F-02 garante).

## Ownership

- Owned paths: `orders/application/enrich_service.go`, `orders/adapters/postgres/` (readers),
  `orders/transport/http_handler.go` (aditivo), região orders root.go, allowlist (remoção A/B),
  par OpenAPI+SDK /orders.
- Forbidden paths: connectors; listings; pricing.
- Parallel-safe with: none — depends on F-02.

## Validation Expectations

- Teste com adapter vivo instrumentado: GET /orders/{id} → 0 chamadas.
- Must-fail do guard: import vivo reintroduzido → allowlist test VERMELHO nomeando o site.
- DTO golden: comprador_fiscal antes/depois byte-igual p/ pedido com fiscal completo.
- Before/after summary por bucket (contagens + explicação de cada delta).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-02.
- Required files/evidence: `validation.md`; before/after.
- Blockers or open decisions: none.
