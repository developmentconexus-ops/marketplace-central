# Interface Contract — divergences (estoque ERP×ML; tarifa 3→2)

```yaml
id: IC-02
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Tabela única entre 2 produtores em cadências diferentes (M-05 estoque no ingest diário de
listings; M-06 tarifa no ingest 5min de pedidos) e 2 leitores FE (/anuncios badge+filtro;
/precos aviso ⚠). Schema nasce no M-02, ANTES de qualquer produtor.

## Why This Contract Exists

ADR-10: sem shape decidido, produtor diário escreve append-events e produtor 5min escreve
one-open-row — badge e filtro "divergentes" perdem significado e o produtor 5min cresce sem
bound. Gate P1: tabela dedicada.

## Resources Or Entities

- Tabela `divergences` (nova, M-02, range 0086-0089).
- Ports Go (núcleo, M-02): `DivergenceRecorder` (avaliar+upsert+auto-resolve),
  `DivergenceReader` (flags/contagens p/ DTOs).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| Evaluate | TODO ingest da entidade (ADR-10: detecção no INGEST, nunca no read) | esperado + observado + fontes + timestamps de observação | upsert de row aberta OU resolve | comparar → divergente: abre/atualiza; igual (dentro da tolerância): row aberta ganha `resolved_at` |
| ListOpenByEntity | read de /anuncios, /pedidos, /precos | tenant, entity refs | kinds abertos por entidade | ordenação: `detected_at` desc |

## Fields

### Required Inputs

- `id bigserial PK`; `tenant_id` NOT NULL; `provider text` NOT NULL.
- `entity_type text` NOT NULL CHECK (`listing`,`order_line`).
- `entity_id text` NOT NULL — formatos pinados (MESMOS de IC-01):
  - estoque: `<provider_listing_id>` sem variação; `<provider_listing_id>:<variation_id>`
    com variação (grão = onde o estoque ML vive);
  - tarifa: `<provider_order_id>:<provider_item_id>`.
- `kind text` NOT NULL CHECK (`estoque`,`tarifa`) — extensível aditivamente.
- `expected_value numeric(14,4)` NOT NULL; `observed_value numeric(14,4)` NOT NULL.
- `expected_source text` NOT NULL; `observed_source text` NOT NULL (ex.:
  `mirror:sankhya:vendavel`, `api_listing:available_quantity`, `channel_fees:layer2`,
  `api_order:sale_fee`).
- `expected_observed_at timestamptz` NOT NULL; `observed_observed_at timestamptz` NOT NULL
  — R-5: os DOIS lados sempre datados; staleness distinguível de desacordo.
- `detected_at timestamptz` NOT NULL (primeira detecção — imutável no upsert);
  `last_evaluated_at timestamptz` NOT NULL (toda reavaliação); `resolved_at timestamptz`
  NULL.

### Required Outputs

DTOs de listings/orders ganham
`divergences: [{kind, expected_value, observed_value, expected_observed_at,
observed_observed_at, detected_at}]` — projeção direta das rows abertas
(`resolved_at IS NULL`; nomes = colunas de §Required Inputs), array vazio = sem divergência
aberta; flag persistida, zero cálculo no read. (P7 r01 B-2: esta é a projeção COMPLETA que o
painel de divergência do M-05/F-03 renderiza — os dois lados + timestamps; `{kind,
detected_at}` sozinho era insatisfazível contra o brief.)

## Enums And Statuses

Row "aberta" = `resolved_at IS NULL`. Sem coluna status — o NULL é o estado.

## Error Cases

Sem superfície HTTP própria de escrita.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| — | — | — | leitura via DTOs existentes; envelope apierror vigente |

## Persistence Expectations

- **Partial unique**: `UNIQUE (tenant_id, provider, entity_type, entity_id, kind) WHERE
  resolved_at IS NULL` — no máximo 1 row aberta por (entidade, kind).
- Upsert de reavaliação divergente: atualiza values/sources/observed_at/last_evaluated_at;
  `detected_at` NÃO muda.
- Convergência: `resolved_at = now()`; row fica como histórico (nada de delete físico).
  Flapping reabre ROW NOVA (detected_at novo) — o partial unique permite.
- Semântica por kind:
  - `estoque`: expected = disponível vendável ERP (estoque−reservado, corte vendável) via
    vínculo product_links; observed = `available_quantity` ML. Tolerância: 0 (qualquer
    diferença diverge).
  - `tarifa`: expected = camada 2 (IC-01) valorada no preço da venda; observed = camada 3
    total da linha. Tolerância: R$ 0.01 (IC-01 Must Preserve).
- Sem vínculo aprovado → NÃO avaliar estoque (não existe "esperado"); nunca divergência
  fabricada contra NULL.

## Canonical Examples

Abertura (estoque):

```json
{"entity_type":"listing","entity_id":"MLB123","kind":"estoque",
 "expected_value":12,"observed_value":9,
 "expected_source":"mirror:sankhya:vendavel","observed_source":"api_listing:available_quantity",
 "expected_observed_at":"2026-07-31T09:00:12Z","observed_observed_at":"2026-07-31T09:03:44Z"}
```

Rejeição canônica: Evaluate sem `expected_observed_at` → recusa do recorder (teste nomeia);
NUNCA default p/ now().

## Database Shape

- Tabela `divergences`; migração no range do M-02.
- Índices: partial unique acima + `(tenant_id, resolved_at) WHERE resolved_at IS NULL`
  (contagem de badge barata).
- Timestamps `timestamptz` UTC.

## Seed Data

Nenhum. Fixture: par criar→resolver por kind (prova nas 2 direções, design §9).

## Timestamp And ID Semantics

- `*_observed_at` = quando CADA lado foi observado na sua fonte, não quando comparados.
- ids sempre do provider (formatos acima).

## Compatibility Rules

- `kind` novo = ADR + CHECK aditivo; shape não muda.
- Resolução manual futura = coluna aditiva (`resolved_by`), sem migração destrutiva.

## Route Namespace

Nenhuma rota nova. `/listings` ganha o filtro `filter.divergentes=true` — idioma do
transport de listings: SÓ chaves prefixadas `filter.` são interpretadas
(`listings/transport/query.go:31`); a chave entra em `domain.FilterKeys`
(`listings/domain/filter.go:9`); um `?divergentes=true` sem prefixo seria silenciosamente
ignorado. Filtro server-side — lista só entidades com row aberta; DTOs ganham o array
`divergences`. Par OpenAPI/SDK no milestone FE respectivo (ADR-14).

## Transport And Integration

N/A.

## Must Preserve

- Detecção SEMPRE no ingest (cadência do produtor); read só lê flag.
- Divergência é informação: NUNCA sobrescrever camada 2 com camada 3 (IC-01).
- Timestamps dos 2 lados NOT NULL — schema, não convenção.

## Must Not Decide In Feature Execution

- Shape da tabela; grão do entity_id; tolerâncias; auto-resolve; comportamento sem vínculo.

## Validation Impact

- 2 direções na MESMA drive de browser: criar → badge aparece em /anuncios; convergir →
  `resolved_at` gravado E badge some.
- Must-fail R-5: fixture com lados observados a horas de distância e valores iguais → NÃO
  diverge; valores diferentes → diverge com os 2 timestamps visíveis na row.
