# Interface Contract — channel_fees (3 camadas com proveniência)

```yaml
id: IC-01
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

Dados de comissão/frete ML entre: produtores (M-05 ingest de listings → camada 2; M-06 ingest
de pedidos → camada 3), consumidor de pricing (M-07 resolver), auditoria 3→2 (M-06 →
divergences IC-02) e telas (/anuncios, /precos). Schema e ports nascem no M-02.

## Why This Contract Exists

Dois produtores em milestones distintos e um consumidor em um terceiro; sem shape único,
cada um inventa representação (percent × amount, unidade × total, listing × variação).
ADR-09: todo fee carrega (camada, origem, coletado_em); `sale_fee` do ML é POR UNIDADE
(fato live T2) — o lugar onde essa multiplicação se pina é aqui.

## Resources Or Entities

- Tabela `channel_fees` (nova, M-02, range 0086-0089).
- Ports Go (núcleo, provider-agnóstico, M-02): `ChannelFeeWriter` (upsert), `ChannelFeeReader`
  (resolução p/ pricing). Nomes finais dos ports = M-02; semântica abaixo é binding.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| UpsertFee | ingest listings (camada 2) / ingest orders (camada 3) | row completa (ver Fields) | — | upsert current-value na chave natural; NUNCA insert duplicado |
| ResolveListingFees | pricing read (M-07) | tenant, provider, installation, provider_listing_id | comissão + frete com proveniência | ordem de resolução: ver Enums; camada 3 NUNCA entra |
| AuditRealizedVsEstimated | ingest de pedido com sale_fee (M-06) | linha camada 3 + row camada 2 correspondente | divergence IC-02 quando delta > tolerância | nunca sobrescreve camada 2 |

## Fields

### Required Inputs (colunas)

- `id bigserial PK`; `tenant_id` NOT NULL (tenant-scoped, sempre no WHERE);
  `provider text` NOT NULL; `installation_id` NOT NULL.
- `layer smallint` NOT NULL CHECK (1,2,3).
- `fee_kind text` NOT NULL CHECK (`commission`,`freight`).
- `subject_type text` NOT NULL CHECK (`category`,`listing`,`order`,`order_line`).
- `subject_id text` NOT NULL. Formatos pinados:
  - listing → `<provider_listing_id>` (ex.: `MLB123`);
  - order → `<provider_order_id>`;
  - order_line → `<provider_order_id>:<provider_item_id>` (separador `:` — MESMO formato
    em IC-02 divergences);
  - category → `<category_id>` (MIS-008).
- `value_type text` NOT NULL CHECK (`percent`,`amount`).
- `value numeric(14,4)` NOT NULL; `currency char(3)` NULL (NOT NULL quando `amount`; NULL
  quando `percent`).
- `detail jsonb` NULL — decomposição da fonte (ver Canonical Examples).
- Proveniência (ADR-09, todos NOT NULL exceto source_time): `origem text` CHECK
  (`api_listing_prices`,`api_shipping_options`,`api_order`,`api_shipment`,`config`);
  `coletado_em timestamptz` NOT NULL (hora do NOSSO fetch); `source_time timestamptz` NULL
  (timestamp do lado provider quando existir — nunca fabricado).
- `api_shipping_options` no CHECK é RESERVA aditiva: NENHUM produtor nesta missão (frete de
  listing = honesto-desconhecido); vocabulário exposto em DTOs desta missão =
  `api_listing_prices | api_order | api_shipment | config` (auditoria P5 F-8).

### Required Outputs (resolução)

`ResolveListingFees` retorna por fee_kind: `{value, value_type, currency, layer, detail,
origem, coletado_em}` — proveniência SEMPRE junto do número; consumidor que exibe número
sem proveniência reprova milestone (ADR-09). `detail` = o jsonb VERBATIM da row resolvida
(camada 2 = tupla canônica de 5 chaves; NULL quando a row não tem detail) — declarado UMA
vez aqui, escopado por camada, nunca re-inventado por consumidor; sem ele o braço camada-2
do M-07 (`detail.percentage_fee/fixed_fee`) seria insatisfazível contra o port (auditoria
P5 r06 F-r06-2).

## Enums And Statuses

Ordem de resolução p/ PRICING (comissão): camada 2 (listing) → camada 1 (category — ausente
nesta missão, MIS-008) → fallback `pricing_tariff_defaults` exposto com proveniência
`config`. **Divisão de dono (auditoria P5 F-10)**: o port
`ChannelFeeReader.ResolveListingFees` (M-02) resolve SÓ o ledger — camada 2 → 1 → **ausente
TIPADO** (nunca zero/default); o degrau `config` é COMPOSIÇÃO do consumidor pricing (M-07
F-01, reusa `pricingtariffdefaults.NewResolver`) — nunca dentro do reader, senão o braço
honest-absent do M-07 fica inalcançável. Frete: camada 2 → **honesto-desconhecido** (frete
NÃO tem camada 1 — design §4; NUNCA número inventado). Camada 3 NÃO participa de resolução
de pricing; alimenta margem realizada + auditoria 3→2.

## Error Cases

Sem superfície HTTP própria (dados servidos via DTOs de listings/pricing existentes).

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| — | — | — | sem endpoint próprio; erros de pricing/listings mantêm envelope apierror vigente |

## Persistence Expectations

- UNIQUE natural key: `(tenant_id, provider, installation_id, subject_type, subject_id,
  layer, fee_kind)` — no máximo 1 row corrente; upsert atualiza `value/detail/coletado_em/
  source_time`. Sem histórico nesta missão (aditivo depois se precisar).
- Camada 3 comissão: `value` = **TOTAL da linha** = `sale_fee_unit × quantity`;
  `detail` OBRIGATÓRIO com `{"sale_fee_unit": x, "quantity": n}` (R-4: por-unidade é
  medição live, doc não declara — o teste de regressão lê o detail).
- Camada 3 frete (produtor M-06 F-02): `value` = custo seller do shipment;
  `subject_type=order`, `origem=api_shipment`; `detail` NULL permitido — NÃO existe
  decomposição sale_fee_unit/quantity p/ frete (F-r06-1).
- `pricing_tariff_defaults` (0068) NÃO migra para channel_fees: permanece tabela própria,
  re-rotulada fallback `config` na resolução. `baseline_commission_percent: 0.16`
  (`auth_adapter.go:42-48`) é METADATA de catálogo do provider (contrato publicado
  wiki/OpenAPI/SDK) SEM call site em pricing — fica intocada e NUNCA entra na resolução
  (auditoria P5 r02 N-2).
- `FeeSyncScheduler`/`RegisterFeeSyncerFactory`: NÃO adotados, NÃO estendidos (ADR-09).

## Canonical Examples

Camada 2 comissão (produtor M-05):

```json
{"layer":2,"fee_kind":"commission","subject_type":"listing","subject_id":"MLB123",
 "value_type":"amount","value":15.99,"currency":"BRL","origem":"api_listing_prices",
 "detail":{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,
 "price_used":79.90,"listing_type_id":"gold_special"}}
```

(Exemplo aritmeticamente consistente com a fórmula ratificada da auditoria 3→2 —
`value = percentage_fee × price_used/100 + fixed_fee + financing_add_on_fee` =
12.5% × 79.90 + 6.00 + 0 = 15.99; auditoria P5 r05 F-r05-2.)

Camada 3 comissão (produtor M-06; sale_fee unitário 8.11, qty 3):

```json
{"layer":3,"fee_kind":"commission","subject_type":"order_line",
 "subject_id":"200012345:MLB123","value_type":"amount","value":24.33,"currency":"BRL",
 "origem":"api_order","detail":{"sale_fee_unit":8.11,"quantity":3}}
```

Rejeição canônica: upsert de camada 3 **fee_kind=`commission`** com `detail` sem
`sale_fee_unit`/`quantity` → writer recusa (constraint de aplicação; teste nomeia). O
mandato de detail é da COMISSÃO camada 3 (§Persistence Expectations) — row camada 3
fee_kind=`freight` do shipment (produtor M-06 F-02: subject_type=`order`,
origem=`api_shipment`) NÃO tem decomposição sale_fee_unit/quantity, `detail` NULL é
aceito (auditoria P5 r06 F-r06-1).

## Database Shape

- Tabela: `channel_fees`; migração no range do M-02 (0086-0089).
- CHECKs: layer, fee_kind, subject_type, value_type, origem (acima); currency NOT NULL
  quando value_type='amount' (CHECK).
- Timestamps: `timestamptz` UTC.

## Seed Data

Nenhum seed. Fixture de validação: rows camada 2+3 do mesmo listing/pedido p/ auditoria 3→2
nas 2 direções.

## Timestamp And ID Semantics

- `coletado_em` = clock nosso no fetch; `source_time` = do provider ou NULL; nunca igualar
  os dois por conveniência.
- `subject_id` sempre id do PROVIDER (nunca id interno).

## Compatibility Rules

- Camada 1 (category) entra na MIS-008 SEM mudança de shape — colunas já comportam.
- Histórico (tabela de eventos irmã) é extensão aditiva futura.

## Route Namespace

Nenhuma rota própria. Valores chegam ao FE via DTOs de /listings (comissão/frete por anúncio
+ data de coleta — design §8) e /pricing (M-07). ADR-14 governa o par OpenAPI/SDK.

## Transport And Integration

N/A (sem origem cruzada nova).

## Must Preserve

- `MissingSaleFee` (profitability `service.go:1014-1015`) continua disparando — NUNCA
  satisfeito por row fabricada.
- Proveniência viaja com o valor até a tela.
- Tolerância da auditoria 3→2: delta > R$ 0.01 (absoluto por linha) gera divergência;
  ≤ 0.01 = arredondamento, não diverge.

## Must Not Decide In Feature Execution

- Representação percent×amount por camada; formato de subject_id; ordem de resolução;
  multiplicação por unidade; intocabilidade do `baseline_commission_percent` (metadata de
  catálogo, N-2); não-adoção do FeeSyncScheduler.

## Validation Impact

- Teste de resolução: camada 2 presente vence config; ausente cai p/ config COM proveniência
  visível; frete ausente = desconhecido (nunca 0).
- Must-fail: fixture com sale_fee unitário × qty>1 — asserção no TOTAL nomeia a falha se
  alguém gravar o unitário.
- Divergência 3→2 provada nas 2 direções (IC-02).
