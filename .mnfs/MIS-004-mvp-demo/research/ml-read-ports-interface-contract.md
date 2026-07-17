# Interface Contract

```yaml
id: IC-06
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

connectors/mercado_livre (M-02, dono único do adapter) ↔ market (interno M-02) ↔ pricing (M-07) ↔ orders (M-08).

## Why This Contract Exists

M-07/M-08 consomem leituras ML sem tocar o adapter (ADR-05). Shapes normalizados aqui; DTO cru do ML NUNCA sai de connectors.

## Operations (ports Go em `modules/connectors/ports`, todos read-only, todos retornam `FetchedAt`)

| Port | Rota ML | Output normalizado | Notes |
| --- | --- | --- | --- |
| `GetOwnItemPricing(itemID)` | `/items/{id}/sale_price` | `{ItemID, Amount, Currency, RegularAmount *decimal, FetchedAt}` | preço vigente do NOSSO anúncio |
| `GetPriceToWin(itemID)` | `/items/{id}/price_to_win?version=v2` | `{ItemID, Status, CurrentPrice, TargetPrice *decimal, Position *{Rank,Total}, FetchedAt}` | alvo competitivo ML — não é "menor preço" |
| `SearchCatalogByEAN(ean)` | `/products/search` | `[]{CatalogProductID, Title, Attrs}` | pré-anúncio por EAN; ordem = relevância do provider, NÃO garantida (consumidor não assume ranking) |
| `GetCatalogProduct(id)` | `/products/{id}` | `{ID, Title, BuyBoxWinner *{Price, SellerID}, Attrs, FetchedAt}` | BuyBoxWinner null fica null (22/22 no probe) |
| `ListCatalogOffers(id)` | `/products/{id}/items` | paginado `[]{SellerID, Price, Condition, ShippingMode}` | **FLAG `MC_ML_CATALOG_OFFERS_ENABLED` default OFF**; paginação COMPLETA obrigatória; telemetria por chamada (counter + status); flag OFF ou rota falha ⇒ erro tipado `ErrCatalogOffersUnavailable` (consumidor converte em NO_PRICE_EVIDENCE); ordem = provider (paginação), NÃO garantida — dedupe/menor-oferta é responsabilidade do consumidor (IC-03) |
| `GetShipmentInfo(shipmentID)` | `/shipments/{id}` (+costs/delays) | `{ID, Status, SLADue *time, Delayed *bool, Costs *{...}, DestinationUF *string, FetchedAt}` | UF destino alimenta chip DIFAL (IC-04) |
| `GetFreeShippingCost(userID, item)` | `/users/{id}/shipping_options/free` | `{Cost *decimal, FetchedAt}` | frete do vendedor p/ preço ≥79 |

## Enums And Statuses

Campos nullable permanecem null (ADR-17 — nunca zero). Erros tipados: `ErrUnauthorized`, `ErrNotFound`, `ErrRateLimited`, `ErrCatalogOffersUnavailable`, `ErrProviderUnavailable` — consumidor decide estado honesto, adapter não inventa fallback.

## Must Preserve

- OAuth via CredentialResolver existente (integrations); nenhum token fora do fluxo atual.
- PROIBIDO: `/sites/MLB/search` (403 provado), scraping, provedor não homologado.
- `PUT /items/{id}` existente (stock, gated) NÃO é tocado neste MIS (zero writes).

## Must Not Decide In Feature Execution

Shapes dos ports, nome/default da flag, semântica de erro tipado, quais rotas ML existem no MVP.

## Validation Impact

Lane live-provider-read com installation real (M-02); telemetria da rota flag inspecionável; mocks provam shape do port, nunca integração viva.
