# Interface Contracts — MIS-006-integracao-fundacao (v1)

**Extends, não recria.** Base = `docs/design/INTEGRATION-DATA-CONTRACT.md` (E1–E7, DRAFT v0
D-120) + `docs/design/STORAGE-SCHEMA.md`. Este arquivo versiona SÓ o que MIS-006 toca:
E2 (extensão), E4 (extensão), E7 (parcial), e os novos E8/E9/E10. E1/E3/E5/E6 = **intocados**
(pertencem a onboarding/ML-sync/mercado — non-scope).

Regra de verdade (profile §8): `ARCHITECTURE.md`/ADR > OpenAPI+SDK > este contrato de dados.
Planning **PROPÕE**; OpenAPI + `sdk-runtime` landam no MESMO commit da implementação (profile §7).

Regra de compat geral: **aditiva**. Colunas novas nascem `NULL`-able (honesto-desconhecido,
ADR-17 — nunca 0/default). Nenhuma coluna existente muda tipo/semântica. Nenhum enum perde valor.

---

## E2 — Produto interno (products_mirror) · **EXTENSÃO**  ·  `contract-version: E2.1`

**O que muda vs v0:** o modelo canônico atual materializa só 6 dos 10 campos
(`NormalizedRow` em `erp_import/domain/import.go:10-22`: codprod, descrprod, custo, stock,
ean, refforn, marca, ncm). MIS-006 materializa o contrato E2 completo numa TABELA (`products_mirror`),
não mais num rescan de snapshot.

Tabela `products_mirror` (nova — STORAGE-SCHEMA §products_mirror):
```
tenant_id, source (sankhya|xlsx|catalogo_cliente), codigo_produto  [PK: (tenant_id, codigo_produto)],
descricao, referencia, ean, marca, grupo_codigo, grupo_descricao, ncm NULL,
custo NUMERIC NULL, preco_venda NUMERIC NULL, estoque_total NUMERIC NULL,
protocol_id NULL (origem xlsx), absent_in_last_snapshot BOOL DEFAULT false, stale_since timestamptz NULL,
updated_at timestamptz
products_mirror_stock_locations(tenant_id, codigo_produto, local_codigo, local_descricao, quantidade)
```
Campos NOVOS vs modelo atual: `local` (via stock_locations), `preco_venda`, `grupo` formalizado,
`absent_in_last_snapshot`, `stale_since`. `ncm` = passthrough opcional (ver Decisão D3).

**Regra de merge (F-XLSX-1, RATIFICADA):** upsert-merge **keep-absent**. Rebuild de snapshot
NUNCA apaga fisicamente. Produto ausente no snapshot novo → `absent_in_last_snapshot=true` +
`stale_since=now()`. Volta a aparecer → flag limpa. Isto protege `product_links` de cascata de wipe.

**Bloqueio de linha:** só `codigo_produto` ausente rejeita a linha (D3-refinada proposta). Demais
campos ausentes → importa COM badge de completude no protocolo, `NULL` no mirror.

**Compat:** `erp_import_products`/`erp_import_protocols` continuam como HISTÓRICO imutável por
protocolo (auditoria). `products_mirror` = estado corrente. Consumidores migram do rescan para o
mirror; a interface de leitura (`readports.Reader` shape) é preservada (ver E-port abaixo).

**EXEMPLO-IO (marker A2):**
- A) Sankhya 90008 → `products_mirror(tenant,90008)={custo 82.00, estoque_total 120, ean ...017,
  marca Ferragens, grupo Torneiras, preco_venda 169.00}` + `stock_locations(CD-01,120)`.
- B) xlsx 74606 (#005-E) → `products_mirror(tenant,74606)={ean ...745, custo 410.00...}`.
  Prospect real #004-E: custo/estoque **NULL** (honesto-desconhecido, nunca 0).

---

## E-port — ProductSourceAdapter · **FORMALIZAÇÃO**  ·  `contract-version: Eport.1`

Não existe hoje como nome, mas a SHAPE existe: `readports.Reader` +
`FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs` satisfeita por
`erp_import/adapters/internalread.Reader` E `internal_read/adapters/oracle.Reader`.

MIS-006 formaliza o port com DOIS lados:
```
ProductSourceAdapter interface {
  // read-side (preservado — consumidores não mudam)
  readports.Reader
  // write/sync-side (NOVO — alimenta products_mirror)
  Sync(ctx) (SyncResult, error)   // xlsx: do snapshot recém-completado; sankhya: read-through→mirror
  Kind() SourceKind               // upload_snapshot | live_read_through
}
```
`SourceKind` distingue upload-shaped (xlsx, catalogo_cliente — têm protocolo/snapshot) de
live-read-through (sankhya — sem protocolo). Enum `ImportSource` naïve NÃO serve (ver E9).

**Compat:** consumidores (pricing, vínculo, mercado) leem só `readports.Reader` — inalterados.

---

## E4 — Vínculo · **EXTENSÃO**  ·  `contract-version: E4.1`

**O que muda:** política de auto-aprovação + trigger + audit (D4 RATIFICADO).
- Trigger automático de geração de candidatos após import (xlsx) / sync (Sankhya) — hoje o
  gerador (`generation_service.go:60-78`) só é alcançável por endpoint órfão (`http_handler.go:89-90`).
- **Auto-aprovar EAN-exato-único:** ao gerar candidatos, se `collisions[ean]==1` (sinal JÁ
  calculado em `validEANCounts`/`identityQuality`, `reader.go:344-366`) E match é exato por EAN
  → auto-transiciona para aprovado, reusando a máquina de audit de `resolution_service.go:129-149`.
- Anúncio sem EAN → fica REVIEW com motivo visível (honestidade, inalterado).
- **Idempotência (A8):** chave única `(internal_product_id, provider_listing_id)`; re-run não
  duplica; override manual do operador vence auto-aprovação e não é sobrescrito.

**Compat:** modelo de candidatos + aprovação manual preservado; auto-aprovação é caminho ADICIONAL.

---

## E7 — Catálogo ML (discovery) · **EXTENSÃO PARCIAL**  ·  `contract-version: E7.1-partial`

SÓ a resolução de identidade EAN→catalog_product_id (F3.1/F3.2) entra. Tarifas (`ml_tariffs`,
`fee_sync`), buy_box, comissão = **intocados** (missão mercado).

`descobrir_produto_catalogo(ean)` (CONDICIONAL — só se live T13-T16 provar, ver Decisão F3.7):
```
GET /products/search?site_id=MLB&status=active&product_identifier={ean}  (read-only)
→ catalog_product_id | vazio (não-catalogável → honesto-desconhecido)
```
Persiste a identidade + enfileira coleta (a EXECUÇÃO da coleta = mercado mission). Se disprovada:
E7.1-partial é REMOVIDA do escopo, decisão registrada, produto sem anúncio recebe só caminho de
vínculo (não de mercado) nesta missão.

---

## E8 — sync_state · **NOVO**  ·  `contract-version: E8.0`

Coração do motor de sync (STORAGE-SCHEMA §sync_state). MIS-006 entrega SÓ o esqueleto (tabela +
loop + cursor read/write + last_error) — jobs pesados de ML = missões seguintes.
```
tenant_id, installation_id, entity (products|listings|orders|market|tariffs),
cursor JSONB, schedule JSONB (interval/cron GENÉRICO — cadence-agnostic, ver D6),
last_full_sync_at, last_incremental_at, last_error JSONB NULL, consecutive_failures INT
```
**Cadence-agnostic (D6):** `schedule` guarda intervalo/cron por entity genericamente; MIS-006
não hardcoda "diário". Cadência real (D6) ratifica na missão mercado sem mudar esta shape.

---

## E9 — Active-source config (por tenant) · **NOVO**  ·  `contract-version: E9.0`

Mata `MC_ERP_SOURCE` de boot. Fonte ativa = config em banco, por tenant.
```
tenant_id [PK], active_source (sankhya|xlsx|catalogo_cliente), source_kind (upload_snapshot|live_read_through),
set_at, set_by
```
Resolução por-request/por-tenant; ambos adapters construtíveis simultaneamente na composition.
**Regra de cache (lição `chip-import-fix-closed`):** `active_source` do ctx DEVE entrar na chave
de cache downstream (`internal_read/adapters/cache/cache.go` `canonicalKey`); adicionar 3ª source
(sankhya) exige estender a key no MESMO change (senão poluição cross-source). Fail-closed via
`ErrUnknownActiveSource` (já existe `reader.go:53`) — nunca fallback silencioso.

---

## E10 — Auto-link audit trail · **NOVO**  ·  `contract-version: E10.0`

E4 cita "audit" mas nunca deu schema. Formaliza:
```
link_id, rule_matched (exact_ean_unique|manual|sku|...), actor (system|operator),
collisions_at_decision INT, created_at, superseded_by NULL (idempotência/override A8)
```
Toda auto-aprovação grava linha `actor=system, rule_matched=exact_ean_unique`. Operador que
sobrescreve → nova linha `actor=operator` + `superseded_by` na anterior.

---

## Matriz de versão + compat

| Contrato | Ação | Version | Compat rule |
|---|---|---|---|
| E1 Conta ML | intocado | — | — |
| E2 Produto | extensão | E2.1 | aditiva; colunas novas NULL-able; history preservado |
| E-port | formalização | Eport.1 | read-side preservada; write-side aditiva |
| E3 Anúncio | intocado | — | (E4 só LÊ listings existentes) |
| E4 Vínculo | extensão | E4.1 | auto-approve = caminho adicional; manual preservado; idempotente |
| E5 Pedido | intocado | — | — |
| E6 Mercado | intocado | — | (mirror é upstream, não é E6) |
| E7 Catálogo | parcial | E7.1-partial | só discovery EAN; condicional live; tarifas intocadas |
| E8 sync_state | novo | E8.0 | cadence-agnostic |
| E9 active-source | novo | E9.0 | cache-key atomicidade obrigatória |
| E10 auto-link audit | novo | E10.0 | — |

---

## Error Matrix

Todo caso de erro retornado por feature desta missão. Semântica de wire final ratifica no OpenAPI
no commit da implementação (profile §7); esta matriz é a fonte de verdade de PROPOSTA por caso.

| Trigger | Resultado tipado | Contrato/Feature dona | Não pode virar |
|---|---|---|---|
| PUT active-source com valor fora do enum (`sankhya\|xlsx\|catalogo_cliente`) | `400 Bad Request` (enum inválido) | E9 · M-02 F-04 | aceito e gravado cegamente |
| Resolver fonte ativa p/ tenant sem config `active_source` | `ErrUnknownActiveSource` (fail-closed, já existe `reader.go:53`) | E9 · M-02 F-03 | fallback silencioso p/ source default |
| `descobrir_produto_catalogo(ean)` → EAN não-catalogável | vazio honesto (não-erro; honest-unknown ADR-17) | E7.1-partial · M-07 F-01 | erro OU `catalog_product_id` fabricado |
| `descobrir_produto_catalogo(ean)` → ML responde `403 PolicyAgent` | erro tipado, distinguível de "vazio legítimo" | E7.1-partial · M-07 F-01 | mascarado como vazio |
| ML responde `429` (rate limit) no probe/discovery | erro tipado de backoff (retry disciplinado; backoff pleno = F-ADAPTER-1 diferido, ADR-06) | E7.1-partial · M-07 | crash OU loop de retry sem limite |
| Chain-read (`/importacoes`) falha ao buscar contagem | UI mostra `—` honesto (ADR-17), não `0` fabricado | chain-read · M-06 F-01 | `0` fabricado apresentado como dado real |
| Import xlsx com linha sem `codigo_produto` | linha rejeitada no parse (não vira row no mirror) | E2.1 · M-03 | row parcial gravada no mirror |
