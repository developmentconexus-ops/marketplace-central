# CHIP-M08-BUYER — Evidence Pack

**Intent:** Surface buyer fiscal identity on the order-detail drawer for ERP registration, mirroring Mercado Livre's documented two-step billing-info flow. Read-time only, additive, no migration.

**BaseSha:** `0df200dce136ef43257f60252c2db908123c2e74`
**HEAD:** `2f089604`
**Branch:** chip/m08-buyer

## What shipped

### ML two-step flow (faithful to deepmap Q2)
1. `GET /orders/{id}` → decode `buyer.{id, nickname, billing_info.id}`. Empty billing_info.id ⇒ honest absence, step 2 skipped.
2. `GET /orders/billing-info/{SITE_ID}/{billing_info_id}` — **Bearer only, NO `x-format-new` header** (that header lives on `doJSONWithHeaders`, used solely by shipments/costs; billing-info uses plain `doJSON`). Decodes name, last_name, identification.{type,number}, address.{street_name, street_number, city_name, state.{code,name}, zip_code, country_id}.
   - SITE_ID = `firstNonEmpty(order.SiteID, a.siteID)` — no hardcoded "MLB".
   - `identification.type` rendered verbatim as opaque `doc_tipo` — never mapped to a CPF/CNPJ enum (ADR-17).
   - `invoice_type` intentionally NOT decoded.

### Degrade split (exact)
- Adapter `doJSON`: 404 → `ErrCodeProviderInvalidReference`, undecodable → `ErrCodeProviderPayloadInvalid`; `buyer_fiscal_reader.go` maps both to empty+nil **silently** (billing-info 404 = normal "none").
- Other codes (auth/rate-limit/transient/validation) and step-1 failures propagate → `EnrichService.resolveBuyerFiscal` **warns once** (installation_id/provider_order_id/error only) then returns nil. Order render always survives.

### Hexagonal boundaries
- Provider structs (`mlBillingInfo*`, `mlOrderBuyer`) die at the ML adapter → neutral `connectors/domain.BuyerFiscalInfo`.
- Two-layer ports: connectors `BuyerFiscalReader` keyed by `ProviderAccountRef`; orders `BuyerFiscalReader` keyed by `installationID`; composition `ordersBuyerFiscalReaderAdapter` bridges installationID→accountRef (mirrors ShipmentReader).

### Transport / contract / FE
- DTO: additive nullable `comprador_fiscal` block (all fields `omitempty`, pointer).
- OpenAPI `OrderCompradorFiscal(Endereco)` + sdk-runtime interfaces landed in the SAME commit (ccb68d67); field-for-field consistent with the Go DTO.
- FE PedidoDrawer: `Comprador · fiscal (ERP)` section (Nome/Razão, Documento = doc_tipo+doc_numero verbatim, Endereço); plus round-4 fields (destinatário, destino_uf·destino_cep, frete_real{bruto,receiver,sender}, rastreio.transportadora + url_rastreio). Absent fields render honest "—". doc_numero never logged.

## Gate evidence

| Gate | Result |
|------|--------|
| `go build ./apps/server_core/...` | OK |
| `go vet` (connectors/orders/composition) | OK |
| `go test` (ML adapter + orders/…) | all `ok`, no FAIL/panic |
| `gofmt -l` (7 touched Go files) | clean |
| SDK `tsc --noEmit` | clean |
| vitest pedidos | 23/23 pass |
| full web vitest | 322 pass |
| vite build | ✓ |
| Governance lane (`harness.ps1 governance -BaseSha <40hex>`) | `status=passed` (only pre-existing baseline exceptions) |

## P6 dual gate (agreement required)

| Reviewer | Round 1 | Fix | Round 2 |
|----------|---------|-----|---------|
| Cold Opus (independent) | FAIL — SDK carried M-05 listing-signal types with no OpenAPI source (scope contamination + OpenAPI↔SDK invariant break) | commit 2f089604 removes M-05 types, keeps only comprador_fiscal | **PASS** — SDK diffs BaseSha..HEAD with only comprador_fiscal; git grep for M-05 types returns none |
| Adversarial sonnet | FAIL — capability_adapter.go var-block not gofmt-clean | commit 2f089604 gofmt-realigns var block | **PASS** — gofmt -l empty; SDK drift also confirmed resolved |

**Both PASS on the corrected diff → agreement reached.**

## Files (BaseSha..HEAD, code only)

```
connectors/domain/buyer_fiscal.go                         (new)
connectors/ports/buyer_fiscal_read.go                     (new)
connectors/adapters/mercado_livre/buyer_fiscal_reader.go  (new) + _test
connectors/adapters/mercado_livre/capability_adapter.go   (edit: SiteID+Buyer decode, interface assertion)
orders/ports/buyer_fiscal_reader.go                       (new)
orders/application/enrich_service.go                      (edit) + _test
orders/transport/http_handler.go                          (edit: comprador_fiscal DTO) + _test
composition/orders_adapters.go                            (edit: bridge adapter)
composition/root.go                                       (edit: wire reader)
contracts/api/marketplace-central.openapi.yaml            (edit: schemas)
packages/sdk-runtime/src/index.ts                         (edit: interfaces)
apps/web/src/pages/pedidos/PedidoDrawer.tsx               (edit) + PedidosPage.test.tsx
```

Full unified diff: `fix.diff` (this directory).

## Hub owns

- Merge to main.
- Live P7 browser QA drive of the drawer against a real ML order.
