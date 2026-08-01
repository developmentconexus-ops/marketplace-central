# F-02-items-multiget-raw-dto

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-01 ml-client-hardening.

## Brief

Leitor multiget `GET /items?ids=` (batches de 20 — fato #13,
`research/external-ml-api-facts.md`) em ARQUIVO NOVO no package
`mercadolivre` (precedente `shipping_reader.go`/`buyer_fiscal_reader.go` — ADR-02/A-02:
`capability_adapter.go` congela). Estabelece a regra DTO da missão: campos consumidos
tipados + `Raw json.RawMessage` capturado (cap 1MB via LimitReader mantido; oversize →
marcador de truncamento explícito no DTO, nunca prefixo silencioso).

EARS:
- While ids > 20, when hidratação pedida, the reader shall particionar em batches de 20 e
  retornar DTOs na mesma ordem dos ids.
- While um id do batch é inválido/404 no multiget, when o ML retorna erro por item, the
  reader shall retornar o DTO de erro POR ITEM sem falhar o batch inteiro.
- While payload > 1MB, when lido, the reader shall truncar Raw e marcar
  `RawTruncated=true` — nunca Raw parcial sem marca.

## Inputs

- ADR-03 (raw seletivo) — persistência é decisão dos milestones de ingest; este feature só
  garante captura no DTO.
- Fato #13 (multiget 20/call) em `research/external-ml-api-facts.md`.
- Todas as chamadas passam por `doRawWithHeaders` (herdam F-01 automaticamente — sem retry
  próprio).

## Expected Output

- Arquivo novo (ex.: `items_multiget_reader.go`) + DTOs tipados com `Raw json.RawMessage`
  + `RawTruncated bool`.
- DTO cobre os campos E3 de IC-07 (sold_quantity, category_id, condition, permalink,
  thumbnail, date_created, tags, catalog_product_id, shipping, available_quantity,
  variations com price/available_quantity/sold_quantity/seller_sku/attributes).
- **Fonte de fee camada 2 (decisão de posse, auditoria P5 F-9)**: a execução VERIFICA se o
  multiget expõe `sale_price` com `?context=channel_marketplace` (memória
  `ml-catalog-offers-pricing-api`) e REGISTRA o desfecho na validation. Se NÃO expõe, este
  feature TAMBÉM entrega reader dedicado de prices (`GET /items/{id}/prices`, arquivo novo
  no mesmo package) com DTO tipado — p/ que M-05 F-01 consuma UMA fonte pronta e nunca
  escreva no package do adapter (posse M-01).

## Constraints

- ZERO edição em `capability_adapter.go` (isso é do F-01).
- Nenhum consumo/persistência aqui — M-04 consome.
- Multiget de shipments NÃO existe (fato T5) — não inventar.

## Inputs/Outputs

- In: `[]providerItemID`, installation ctx. Out: `[]ItemDTO` alinhado por id, erro por item
  embutido, `Raw` por item.
- Resposta multiget do ML = array de `{code, body}` por id — mapear code≠200 p/ erro por
  item.

## Negative Scenarios

- Batch inteiro 429 → propaga erro do F-01 (sem tratamento próprio).
- Body de item não-JSON → erro por item nomeado, batch segue.

## Ownership

- Owned paths: arquivo(s) novo(s) em
  `apps/server_core/internal/modules/connectors/adapters/mercado_livre/`.
- Forbidden paths: `capability_adapter.go`, `refresh_policy.go`.
- Parallel-safe with: F-01 (eixo files — só arquivos novos).

## Validation Expectations

- 45 ids → exatamente 3 chamadas HTTP (20+20+5), ordem preservada — contagem no transport
  fake.
- Item 404 dentro do batch → 44 DTOs válidos + 1 erro por item; batch NÃO falha.
- Fixture >256KiB (por item) → `RawTruncated=true` assertado — ver amendment abaixo.

**Amendment (2026-08-01, orchestrator M-01, self-classified spec conflict, achado por
review adversarial de F-02):** este brief e ADR-03 citam "1MB" para o cap de Raw, mas a
implementação usa `itemMultigetRawCap = 256KiB` POR ITEM, não 1MB. Razão: o cap de 1MB é do
`doRawWithHeaders`/`io.LimitReader` (`capability_adapter.go:747`) sobre a RESPOSTA INTEIRA do
multiget, compartilhada por até 20 itens do batch — um cap por-item de 1MB nunca seria
alcançável nem testável nessa leitura (um item sozinho >1MB já estouraria o cap externo antes
do array chegar parseável a este reader, virando falha de BATCH, não por-item). 256KiB por
item é folga generosa sobre qualquer payload real de item ML (poucas KB) mantendo espaço para
um batch de 20 itens caber no teto externo de 1MB. Comentário do código
(`items_multiget_reader.go` doc de `capRawMessage`) e teste (`items_multiget_reader_test.go`)
corrigidos para citar 256KiB explicitamente em vez de "(1MB)"/referência tautológica à
constante. `capability_adapter.go`'s outer 1MB cap in itself is unaffected — compliant with
ADR-03 as written for the whole-response case.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` = execução.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` na execução.
- Blockers or open decisions: none.
