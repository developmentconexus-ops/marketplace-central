# F-01-ml-adapter-read-ports

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-02-price-intel-core.

## Brief

Estender o capability_adapter Mercado Livre em `connectors` com os 7 ports read-only do IC-06: `GetOwnItemPricing`, `GetPriceToWin`, `SearchCatalogByEAN`, `GetCatalogProduct`, `ListCatalogOffers` (flag-gated), `GetShipmentInfo`, `GetFreeShippingCost`. DTO cru do ML nunca sai do adapter; nullable fica null; erros tipados.

## Inputs

- IC-06 (`research/ml-read-ports-interface-contract.md`) — shapes, rotas, flag, erros (fonte única; não redecidir).
- Adapter ML existente em `modules/connectors/**` (OAuth via CredentialResolver — reutilizar fluxo, zero token novo).
- Research de pricing (probe evidence: buy_box_winner null 22/22; /sites/MLB/search 403 PROIBIDO).

## Expected Output

- 7 ports em `modules/connectors/ports` com shapes EXATOS do IC-06, todos com `FetchedAt`.
- Flag `MC_ML_CATALOG_OFFERS_ENABLED` default OFF; OFF ou rota falha ⇒ `ErrCatalogOffersUnavailable`.
- `ListCatalogOffers`: paginação COMPLETA obrigatória (loop até esgotar; parcial = erro, nunca resultado truncado silencioso) + telemetria por chamada (counter com status/página).
- Erros tipados: `ErrUnauthorized`, `ErrNotFound`, `ErrRateLimited`, `ErrCatalogOffersUnavailable`, `ErrProviderUnavailable`.
- EARS: While flag OFF, when `ListCatalogOffers` é chamado, the adapter shall retornar `ErrCatalogOffersUnavailable` sem chamar o ML. While ML responde 429, when qualquer port é chamado, the adapter shall retornar `ErrRateLimited` (sem retry interno silencioso). While `buy_box_winner` vem null, when `GetCatalogProduct` responde, the adapter shall manter `BuyBoxWinner: nil` (nunca preço 0).

## Inputs/Outputs

Shapes por port: IC-06 §Operations (referência, não restate). Mapeamento rota ML → shape normalizado é decisão DESTE feature apenas onde IC-06 deixar campo interno livre; shapes públicos são fixos.

## Negative Scenarios

- Item inexistente ⇒ `ErrNotFound` (não struct vazia).
- Token inválido/expirado ⇒ `ErrUnauthorized`; refresh segue fluxo CredentialResolver existente — falhou refresh, propaga tipado.
- Rota price_to_win indisponível p/ item sem catálogo ⇒ `ErrNotFound` mapeado, consumidor decide estado honesto.
- Timeout/5xx ML ⇒ `ErrProviderUnavailable` com causa embrulhada (log estruturado, sem payload cru em log).

## Constraints

- Read-only absoluto: nenhum PUT/POST ao ML; `PUT /items/{id}` existente intocado.
- PROIBIDO `/sites/MLB/search`.
- Provider payload não vaza do adapter (doutrina AGENTS.md).
- Mocks provam shape do port; lane live-provider-read prova integração (installation real).

## Ownership

- Owned paths: `apps/server_core/internal/modules/connectors/**` (ports novos + client ML + telemetria + testes).
- Forbidden paths: `modules/market/**` (F-02..F-04), OpenAPI, SDK, migrations, demais módulos.
- Parallel-safe with: F-02 (disjoint: módulos distintos, sem arquivo compartilhado).

## Validation Expectations

- Mock-contract tests por port: request esperado + response fixture ⇒ shape normalizado exato (incl. caso null → nil).
- Teste paginação: fixture 3 páginas ⇒ resultado com soma exata das 3; página 2 falha ⇒ erro, zero resultado parcial.
- Transcript flag OFF ⇒ `ErrCatalogOffersUnavailable` sem hit HTTP (mock não tocado).
- Lane live-provider-read: `GetOwnItemPricing` + `GetPriceToWin` contra anúncio próprio ativo real — valores não-nulos + `FetchedAt` recente no log.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-02). Preflight R4 (permissões conta) na primeira hora.
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + log lane live-provider-read.
- Blockers or open decisions: none.
