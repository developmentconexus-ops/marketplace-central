# P2 — verificação de premissas (2026-08-02)

Medido contra os 38 pedidos reais do dev stack, `installation_id=inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0`.

## Step 1 — margem NÃO é o problema

```
GET /orders?limit=50&installation_id=inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0
count: 38
margem_pct nonnull = 33
margem_pct null = 5
```

**33/38 — bate com o plano.** Confirma que `BuildProfitability` já funciona; nada a reimplementar.

## Step 2 — os seis buracos

```
currency: nonnull=0 null=38
fulfillment: nonnull=0 null=38
provider_status_detail: "" em 38/38 (nunca null, nunca preenchido de fato)
cancellation_detail: campo NEM EXISTE no JSON HTTP hoje — 0/7 cancelados têm motivo
```

```sql
select count(*)||'|'||count(logistic_type)||'|'||count(tracking_method) from order_shipments;
-- 38|0|0
```

**Bate com o plano nos três campos e nos dois campos de shipment.**

## Step 3 — perfil fiscal ausente

```sql
select count(*) from pricing_calc_profiles;
-- 0
```

**Bate com o plano.** P2.b (perfil fiscal, regime, DIFAL) ainda não começou. `imposto`/`difal` seguem fora de escopo desta fatia.

## Veredito

**PREMISSA CONFIRMADA.** Nenhuma divergência com o plano. Prosseguindo com as tasks 2-6 como escrito, sem tocar em `pricingtax/`, `calc_repository.go` ou perfil fiscal.
