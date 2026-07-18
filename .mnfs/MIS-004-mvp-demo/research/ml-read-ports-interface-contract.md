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

## Amendment A1 — 2026-07-17 (hub-ratified, D-12; drafted by CHIP-M02 per D-04 flow)

Freezes the 4 previously underspecified nested shapes. Top-level members above remain FROZEN
unchanged. Package `connectors/domain`. Money = decimal string; nil = null (ADR-17).

**Shared:**
```go
type CatalogAttribute struct {
    ID        string
    Name      string
    ValueName *string // nil = null
}
```
RULING: `value_id` OMITTED (hub-approved) — gate matches on id + normalized value_name
(handoff §5); provider-internal enum key not consumed. Re-add = explicit amendment.

**SearchCatalogByEAN:**
```go
type CatalogSearchResult struct {
    Products  []CatalogCandidate // provider order, no ranking assumed
    FetchedAt time.Time          // ENVELOPE placement (resolves :29-vs-:23 gap)
}
type CatalogCandidate struct {
    CatalogProductID string
    Title            string
    Attrs            []CatalogAttribute
}
```

**GetCatalogProduct:**
```go
type CatalogProduct struct {
    ID           string
    Title        string
    BuyBoxWinner *BuyBoxWinner // nil when null (22/22 in probe)
    Attrs        []CatalogAttribute
    FetchedAt    time.Time
}
type BuyBoxWinner struct {
    Price    *Money // decimal string
    SellerID string
}
```
FLAG (no shape change): live `buy_box_winner` payload carries `item_id` not `seller_id`
(probe:150-152); SellerID may need item-join at wire-up; path unexercised (null 22/22).
If live probe confirms seller_id unobtainable → amendment at F-01 probe (dated deferral).

**GetShipmentInfo:**
```go
type ShipmentInfo struct {
    ID            string
    Status        string
    SLADue        *time.Time
    Delayed       *bool
    Costs         *ShipmentCosts
    DestinationUF *string // 2-letter UF from receiver_address.state → DIFAL (IC-04)
    FetchedAt     time.Time
}
type ShipmentCosts struct {
    GrossAmount  *Money
    ReceiverCost *Money
    SenderCost   *Money // all nil = unknown
}
```

**GetFreeShippingCost:**
```go
type FreeShippingQuery struct{ ItemID string }
type FreeShippingCost struct {
    Cost      *Money
    FetchedAt time.Time
}
```
RULING: minimal `{ItemID}` APPROVED (ML resolves dimensions server-side from listing).
DATED DEFERRAL 2026-07-17: if live-lane probe shows route requires explicit dims/price,
finding folds in at F-01 probe time as amendment A2.

Typed errors unchanged (IC-06:37). DTOs die at adapter; GET-only.
