# Integration Findings — D-120 (base para planning MIS-006)

Findings desta sessão de replanejamento que NÃO estavam nos docs de design quando foram
escritos. Bindam o planning da integração. Base: main @138aac3dff20438d8ddc509daf4171d82e5e45f6.

Relacionado: [SYSTEM-BLUEPRINT.md](../SYSTEM-BLUEPRINT.md),
[STORAGE-SCHEMA.md](../STORAGE-SCHEMA.md), [SCENARIO-WALKTHROUGH.md](../SCENARIO-WALKTHROUGH.md),
[INTEGRATION-DATA-CONTRACT.md](../INTEGRATION-DATA-CONTRACT.md).

---

## F-XLSX-1 — xlsx é UPSERT-MERGE, NUNCA rebuild/wipe (correção do operador)

O texto anterior "cada upload = novo snapshot → rebuild mirror" e "troca de source ativa =
rebuild do mirror" está ERRADO e é anti-usuário. Rebuild = wipe = perde produtos + perde
`product_links` que o operador criou. Modelo correto:

- **Upsert-merge no `products_mirror`, chaveado por `codigo_produto` (ou EAN).**
- Linha existe no arquivo novo → UPDATE campos.
- Linha nova → INSERT.
- **Linha que sumiu do arquivo novo → NÃO deleta.** Permanece no mirror; marca
  `stale_since` (timestamp). Vínculo sobrevive.
- `product_links` intactos: chaveiam em codigo_produto/EAN (estáveis).
- `erp_import_products` continua histórico por protocolo (auditoria "o que veio no arquivo X").
- Decisão do operador: produto descontinuado (some por N uploads) **nunca sofre delete físico**
  — igual anúncio pausado no ML (some da venda, fica no banco). Flag `stale` visível na tela.

Consequência de schema: `products_mirror` ganha `stale_since timestamptz NULL`. STORAGE-SCHEMA.md
§products_mirror deve trocar "rebuild do mirror" por "upsert-merge keep-absent + stale flag".
SYSTEM-BLUEPRINT.md §2 tabela idem (etapa 3 "rebuild mirror" → "upsert-merge").

---

## F-SDK-1 — Estudo do golang-sdk oficial (github.com/mercadolibre/golang-sdk, deprecado)

Repo minúsculo: 1 arquivo real (`sdk/meli.go`). Design deliberadamente CRU/low-level.

**O que é:** transporte burro. 1 `Client` = 1 conta (creds + 1 `Authorization` em memória, não
multiplexa). 4 wrappers `Get(path)`/`Post(path,body string)`/`Put(path,body string)`/`Delete(path)`
retornando `*http.Response` cru. Body de POST/PUT = **string JSON literal** (zero DTO). Token vai
como **query param** `?access_token=` (estilo ML antigo). Refresh via interface `TokenRefresher`
+ `ReceivedAt`/`ExpiresIn`.

**O que NÃO tem (de propósito):** zero DTO tipado, zero paginação, zero retry/429/Retry-After,
zero `context.Context`, zero multi-tenant.

**Lição dupla:**
1. Pró-nosso: o SDK oficial dá NADA de contrato de domínio. Definir DTOs tipados por endpoint no
   adapter é obrigatório — não existe lib oficial que já modele order/item. Mata a fantasia
   "só usa o SDK".
2. Alerta: eles mantêm o transporte burro/stringly-typed de propósito — a API ML muda campo o
   tempo todo. Struct rígido demais quebra no desconhecido. Padrão certo = **struct tipado só
   nos campos consumidos + `json.RawMessage` guardando o payload cru** (não perde o que não
   mapeia, não quebra no novo). Bate direto com o medo do operador: "criamos campo mas está stub
   / não mapeamos um campo".

---

## F-ADAPTER-1 — Veredicto do NOSSO ML adapter vs SDK oficial

Mapeamento file:line em `apps/server_core/internal/modules/connectors/adapters/mercado_livre/`
+ `.../integrations/adapters/mercadolivre/`. **Nosso transporte é MUITO mais maduro que o SDK
oficial** — não é aí que o sistema está bugado. Bugs são de ORQUESTRAÇÃO, não de transporte.

Nós melhor/certo: Bearer header (`capability_adapter.go:716`, não query param) · `context.Context`
em todos wrappers (`:711`) · multi-conta via `AccessTokenResolver` + X-Tenant/X-Installation
headers (`:26,:721`) · 15+ DTOs tipados por endpoint (`:971` item, `:1022` order, shipping,
catalog, pricing, billing) · camada mapper DTO→domínio (`mapListing :753`, `mapOrder :798`) ·
2 estratégias de paginação (offset `:233` + scroll_id `:279`) · redaction de token/secret em erro
(price_writer `:19`) · body LimitReader 1MB (`:744`) · x-format-new + `flexString` tolerante
(`shipping_reader.go:56,155`) · X-Idempotency-Key nos writes (`:724`).

### 3 buracos REAIS do adapter (entram no inventário de refactor)

1. **Sem retry/backoff/Retry-After.** Todo 429 vira `ErrCodeProviderRateLimited` e desiste
   (`capability_adapter.go:653`, price_writer `:79`, listing_writer `:105`). Live: 100 GETs
   concorrentes = zero 429, mas backfill de milhares vai dar. Precisa backoff + honrar
   Retry-After ANTES de sync em lote.
2. **`fee_sync.go:29` semeia 16%/22% ESTÁTICO** — nenhum sync live de tarifa. É o seed que
   `ml_tariffs` mata (ml-tariff-design + R$79 hardcode banido). Confirmado no código.
3. **DTOs rígidos SEM `json.RawMessage` cru** (só `flexString` é tolerante). Campo que o ML
   adiciona/renomeia = perda silenciosa. Recomendação: cada DTO principal (item/order/shipment)
   carrega `Raw json.RawMessage` do corpo original, persistido. Rede de segurança pro medo
   "não mapeamos um campo".

**Direção confirmada pelo estudo:** REUSAR o adapter, não reescrever. O trabalho de integração é
a camada ACIMA (sync engine, sync_state, scheduler, persistir enrich no ingest). Adapter só
ganha: backoff/429, tarifa live, `Raw` nos DTOs. Padrão de DTO vira regra:
tipado-nos-campos-usados + Raw cru guardado.
