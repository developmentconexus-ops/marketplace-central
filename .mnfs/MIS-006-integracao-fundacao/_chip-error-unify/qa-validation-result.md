# QA live-drive VC-9 — CHIP-ERROR-UNIFY pós-merge (hub, 2026-07-31)

Merge @6bc30c4c (--no-ff, gate P6 waiver §12: round 0 Opus REJECT + sonnet APPROVE →
remediação 42808712/bce555a8 → round 1 Opus APPROVE). Escada pós-merge VERDE nas duas
lanes (D-10): Go 110 pkgs ok · integração status=passed run_id=f8bac93e85c246acbc528a4bb2f9f976 ·
web tsc 12 exatos · sdk tsc prod 0 + test-lane 7 exatos · vitest web 566 + sdk 85 ·
governança 54 (breakdown 19/13/10/6/5/1), GOV_PRODUCTION_PANIC=0, GOV_SCHEMA_INVALID=0,
exceção production-panic-apierror-recover-abort-handler reconhecida. Dev stack rebuildado
do checkout principal (backend/frontend recriados; postgres preservado, banco velho A-26).

## Telas dirigidas (browser real :5174, backend :8080)

| # | Tela / erro forçado | Fio | UI | Veredito |
|---|---|---|---|---|
| 1 | `/catalogo/produtos/abc` (id inválido) | 400 `{"error":{"code":"invalid_identity","message":"productId must be a positive integer","details":{}}}` | "Produto não encontrado" + ErrorState "Identificador de produto inválido." + link de volta | PASS |
| 2 | `/catalogo/produtos/999999999` (recurso inexistente) | 404 `{"error":{"code":"CATALOG_PRODUCT_NOT_FOUND","message":"product not found","details":{}}}` | envelope correto no fio; UI renderiza workspace vazio (ver REPORT abaixo — pré-existente, não é do chip) | PASS (fio) |
| 3 | `/integracoes` upload de .xlsx lixo via UI (file input + Importar) | 400 `{"error":{"code":"invalid_file","message":"invalid xlsx file","details":{}}}` | "Arquivo inválido. Envie um .xlsx exportado do ERP." — mensagem tipada, sem string crua, sem tela branca | PASS |

Console: **zero erro** na sessão inteira (todas as telas). Nenhuma tela branca/undefined em
rota casada. Recover central: provado por pin de unidade (recover_test.go, 5 testes; panic
route de teste removida antes do merge, conforme plano).

Controles negativos observados de graça: FE **saneia** query params inválidos em vez de
repassar (`?erp_source=invalido` e `?filter.bogus=1` nunca chegam ao backend — só filtros
conhecidos vão ao fio); botão Atualizar de /anuncios desabilita durante refresh (409
refresh_in_progress inalcançável por duplo-clique).

## REPORTs (pré-existentes, fora do chip — fila do hub)

- **ProdutoPage não trata 404 do produto**: guard cobre só id não-numérico
  (ProdutoPage.tsx:45); fetch 404 não vira estado not-found — página mostra workspace vazio
  "Produto 999999999" com tudo "—". Desenho partialFailure antigo; melhoria = missão futura.
- **Router sem fallback 404**: rota não-casada (ex. `/produto/abc`) = tela branca
  (AppRouter.tsx sem catch-all).

## Veredito

**QA PASS — VC-9 cumprido. CHIP-ERROR-UNIFY encerrado.**
