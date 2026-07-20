# Identidade: Produto ↔ Anúncio ↔ Marketplace ↔ Evidência de Mercado

> **Status:** PROPOSTA (aguarda ratificação do operador) — 2026-07-19
> Nasceu da pergunta do operador: "não definimos uma arquitetura do sistema como um
> todo, um padrão entre anúncio, produto e marketplace". Este doc mapeia o que HOJE
> existe (verificado em código), nomeia a fragilidade real, e propõe o alvo durável.

## 1. O que existe hoje (verificado, file:line)

A relação NÃO é disjunta — a cadeia produto→anúncio→mercado fecha ponta a ponta e
funciona viva (bloco 90001–90034 = 7 veredictos OK com mediana real, custo ERP + posição ML).

```
  PRODUTO (ERP)                ANÚNCIO (ML)                 MERCADO
  internal_product_id  ──┐     listings                     market_aggregates
  (codprod)              │     (provider_listing_id, MLB…)  .product_id  ← = internal_product_id
  custo, ean, refforn    │            ▲                            │  (por CONVENÇÃO, sem FK)
                         │            │                            │
                         └── product_links (VÍNCULO) ──────────────┘
                             internal_product_id ↔ (provider_item_id, variation)
                                                                 market_competitive_signals
                                                                 .listing_id  ← anúncio↔mercado
```

| Par | Chave de join | Onde declarado |
|---|---|---|
| produto ↔ anúncio | `product_links.internal_product_id` ↔ `(provider_item_id, provider_variation_id)` | `migrations/0025_product_link_workflows.sql:9` |
| produto ↔ mercado | `market_aggregates.product_id` = `internal_product_id`/codprod (mesmo valor) | `migrations/0053_market_signals_aggregates.sql:25`; `internal/modules/market/application/evidence_read_service.go:78-103` |
| anúncio ↔ mercado | `market_competitive_signals.listing_id` | `migrations/0053_market_signals_aggregates.sql:3` |

Identidade ERP resolvida via ponte CODPROD: `migrations/0035_catalog_codprod_compatibility.sql`,
`0034_catalog_legacy_product_reference_evidence.sql`. Import cru ERP:
`migrations/0046_create_erp_import_products.sql` (`codprod`, `descrprod`, `custo`, `ean`, `refforn`).

## 2. Fragilidade real (o núcleo verdadeiro da preocupação do operador)

1. **Sem tabela canônica única de produto.** Identidade ERP é virtual, espalhada por
   `product_enrichments` / `classification_products` / facts / `erp_import_products`. Não há
   um `catalog_products` que seja a fonte-única-de-verdade de "o que é um produto".
2. **produto↔mercado ligado por convenção, não por FK.** `market_aggregates.product_id` é
   TEXT e "deve" igualar `internal_product_id` — nada no banco força. Um writer que use o
   id errado quebra o join silenciosamente.
3. **Seleção "com preço" depende de ordem de id** (cursor `base64("90000")` no frontend,
   `PricingPage.tsx`), não de filtro semântico. Import futuro que coloque produto cost-only
   acima do bloco reintroduz o sintoma "só anúncios / tudo —".
4. **Identidade de exibição ambígua.** A coluna SKU renderiza `manufacturer_reference`; se
   esse campo carregar o id ML em vez do REFFORN/CODPROD, o produto "parece" anúncio.

## 3. Alvo durável (pós-demo, para ratificar)

- **PROD-1** — Tabela/vista canônica de produto (`catalog_products` ou vista materializada)
  como fonte única de `internal_product_id → {codprod, refforn, ean, descr, custo}`.
- **PROD-2** — FK real (ou constraint de validação) de `market_aggregates.product_id` para a
  identidade canônica; escrita de agregado rejeita id não-resolvível.
- **PROD-3** — Endpoint de catálogo com filtro semântico (`listed=true` /
  `has_market_evidence=true`) — aposenta o cursor-por-id-tail.
- **PROD-4** — Contrato de identidade de exibição: SKU exibido = REFFORN/CODPROD ERP;
  id ML só aparece no contexto de anúncio, nunca como identidade de produto.

## 4. Regra binding provisória (até PROD-1..4)

Até a ratificação: **toda tela que mostra "produto" exibe identidade ERP (CODPROD/REFFORN),
nunca o id MLB.** O id ML pertence ao contexto de anúncio/vínculo. Telas de preço/simulador
consomem o produto pela sua identidade ERP e trazem a evidência de mercado pelo join acima.
