# M-07 — chicken-egg break + F3.7 discovery (CONDICIONAL)

```yaml
id: M-07
type: milestone
mission: MIS-006
status: out_of_scope   # operador 2026-07-28: MIS-006 é ERP + planilha + vínculo; descoberta EAN→catálogo é caminho de MERCADO e migra para a missão mercado junto com MC-11. Nada aqui é construído; o gate live T13-T16 NÃO roda nesta missão. Ver ../CORTE-YAGNI.md
depends_on: [M-02]
gated_by: [live T13-T16 REQUEST hub]
conditional: true
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Fechar o chicken-egg do namesake da missão (`Collect(codprod)` só resolve mercado via
`LinkedListings(codprod)`, `collection_pipeline_service.go:269` — produto sem vínculo não tem
caminho de mercado) para o caso de produto **sem anúncio nosso**: descobrir a identidade de
catálogo ML por EAN (`descobrir_produto_catalogo(ean)`, F3.7/E7.1-partial) e enfileirar coleta.

Esta milestone é **CONDICIONAL**: nada do bloco "Feature Briefs" é construído sem prova live
prévia. Ver `mission.md` §Decisions Resolution (F3.7 — PENDENTE, owner-authority), §Real-integration
bindings item 2, §Risks ("F3.7 disprovada em T13-T16"), `research/contracts-decisions-scenario.md`
§F3.7 viability, `docs/design/SCENARIO-WALKTHROUGH.md` Adendo A1 (chicken-egg) + §Rodada LIVE 2
(T13-T16), `docs/design/ML-API-QUERY-CATALOG.md` F3.1 (`identidade_catalogo(ean)`, marcado `DOC`
não `LIVE`).

## Scope / Non-Scope

**Scope (só se o Live Gate PASSAR):**
- CREATE `descobrir_produto_catalogo(ean)` — GET read-only, EAN→`catalog_product_id` ou vazio.
- Persistir identidade descoberta (liga `products_mirror` row ao `catalog_product_id`, quando
  encontrado).
- Enfileirar coleta em `sync_state` (E8) para o produto-sem-anúncio recém-identificado — a
  EXECUÇÃO da coleta (escrever `market_aggregates`/`competitor_offers`) fica fora (MC-11).
- `cmd/mlprobe` (round 3, evidência T13-T16) — ferramenta de prova, não produto.

**Scope (sempre, independente do gate):**
- Rodar o Live Gate (T13-T16) e registrar o veredito, PASS ou FAIL, com evidência.
- Se o gate FALHAR: REMOVE F3.7 do escopo (nenhum código de F-01/F-02 é escrito), registra
  decisão honest-unknown em `mission.md`/`interface-contracts-mis006.md` (ADR-17), e documenta
  que produto-sem-anúncio, nesta missão, só recebe o caminho de VÍNCULO (M-05, contra listings
  já existentes) — não o caminho de MERCADO. Mercado-path para produto nunca-listado vira escopo
  da missão mercado.

**Non-Scope (sempre):**
- Execução de coleta (`market_aggregates`/`competitor_offers` write) — missão mercado (MC-11).
- Tarifas (`ml_tariffs`), buy_box, comissão — E7 full, missão mercado.
- Qualquer sync de anúncio/pedido (E3/E5) — missão ML-sync.
- Onboarding saga, writes ML, webhooks — MIS-005/onboarding.
- Rodar a rodada live diretamente no chip — é REQUEST ao hub (chip não executa; ver
  `mission.md` §Real-integration bindings item 2).

## Live Gate (T13-T16) — bloqueia todo Feature Brief abaixo

Rodada live, read-only, via `apps/server_core/cmd/mlprobe` (round 3; hoje untracked no worktree
que `hub-erp-main` monta — permanece untracked até promovido por decisão do hub). Chip **não
executa diretamente**: emite `REQUEST live-round-T13-T16` ao hub; hub decide ambiente/credencial
e roda (ou delega a quem tem acesso à conta ML ativa).

**Fonte dos EANs:** `erp_import_products`, protocolo **#004-E** (prospect, 2012 produtos).
Amostra: 10 EANs reais, mix com/sem anúncio nosso (evita viés — precisa provar o caso
"sem anúncio", que é o caso de interesse do chicken-egg).

**Credencial:** conta ML ativa já conectada no banco (não uma nova). **NUNCA exposta/impressa**
em log, evidência, ou output de chat (AC-05) — `mlprobe` lê do storage de token existente, nunca
de `.env*`, nunca em stdout.

| Teste | Chamada | Decide |
|---|---|---|
| T13 | `GET /products/search?site_id=MLB&status=active&product_identifier={EAN}` × 10 EANs #004-E | F3.7 existe? shape da resposta? taxa de acerto EAN→catalog_product_id? |
| T14 | p/ cada `catalog_product_id` achado em T13: `GET /products/{id}` + `GET /products/{id}/items` | dados de demanda (ofertas/preços) suficientes p/ compor uma linha de Oportunidade? |
| T15 | fallback `GET /sites/MLB/search?q=` p/ EAN sem catálogo em T13 | aberta ou bloqueada por PolicyAgent (precedente T8/F2: 403 em leitura de terceiro já provado neste projeto)? |
| T16 | simulação de margem fim-a-fim p/ produto B (`74606`, xlsx #005-E): agregados de T14 + `ml_tariffs` sweep + frete (F4.2) | margem potencial é computável fim-a-fim mesmo sem venda nossa? |

**Evidência:** cada teste grava output em `docs/design/evidence/ml-api/T13-*.json` ...
`T16-*.json` (mesmo padrão de `T0-T12` já existentes ali). Tipo de evidência por teste:
`ran` (executado, output salvo) ou `could-not-run` (bloqueado — nomear bloqueio, ex.
"REQUEST hub sem resposta", "credencial indisponível").

**Decisão fork (binária, registrada em `mission.md` pós-rodada):**
- **PASS** (T13 taxa de acerto > 0 E T15 não-403 OU T14 supre dados suficientes mesmo com T15
  bloqueado) → Feature Briefs F-01/F-02/F-03 abaixo são construídos.
- **FAIL** (T13 endpoint não existe/vazio sistemático, OU T15 confirma 403 PolicyAgent E T14 não
  supre alternativa) → F3.7 é REMOVIDA do escopo da missão; `interface-contracts-mis006.md` §E7
  marca `E7.1-partial` como `REMOVIDO — honest-unknown (T13-T16 disprovaram)`; produto sem
  anúncio permanece com rota SÓ de vínculo (M-05) nesta missão; mercado-path vira item de
  backlog explícito para a missão mercado.

Este fork É o `MC-10` da missão (`validation-contract.md` §MC-10) e a prova mínima do M07-C1
desta milestone.

## Feature Briefs (CONDICIONAL — só se o Live Gate = PASS)

### F-01 — `descobrir_produto_catalogo(ean)` (read-only discovery)

**EARS:**
- WHEN `descobrir_produto_catalogo(ean)` é chamada para um EAN presente em `products_mirror`
  THEN executa `GET /products/search?site_id=MLB&status=active&product_identifier={ean}`
  (read-only, sem side-effect na conta ML) e retorna `catalog_product_id` OU vazio.
- IF o EAN não é catalogável (resposta vazia/404) THEN o resultado é honesto-vazio — NÃO um
  erro, NÃO um valor default fabricado (ADR-17) — o produto segue existindo em `products_mirror`
  sem `catalog_product_id`.
- WHEN a chamada recebe `403 PolicyAgent` (padrão já observado em T8/F2 para leitura de
  terceiro) THEN o erro é propagado tipado e documentado — nunca mascarado como "vazio" (evita
  confundir "não catalogável" com "bloqueado pela ML").

**Inputs/Outputs (MUST have):**
- Input: `ean string` (de `products_mirror.ean`).
- Output: `catalog_product_id *string` (nil = não-catalogável) + erro tipado distinto para
  bloqueio de policy vs resposta vazia legítima.
- Read-only: nenhuma escrita na conta ML; único side-effect é local (persist identidade, F-02).

**Negative Scenarios:**
- EAN inválido/vazio na row de origem → função rejeita antes de chamar a API (não gasta rate
  limit em input já sabido inválido).
- Resposta 403 PolicyAgent → erro tipado propagado, log/evidência registra o bloqueio nomeado,
  NUNCA vira `catalog_product_id=""` silencioso.
- Rate limit (429) → respeita backoff do cliente ML existente (adapter maduro reusado,
  `mission.md` ADR-06 — não reimplementa client HTTP).

**Write-set:** `internal/modules/market/application/*` (novo arquivo de discovery, half
"discovery" do owns de `architecture-map.md` §Superfícies compartilhadas).

---

### F-02 — Persist identidade + enqueue coleta (produto sem anúncio)

**EARS:**
- WHEN `descobrir_produto_catalogo(ean)` retorna `catalog_product_id` não-vazio para um produto
  de `products_mirror` sem vínculo (`product_links` ausente) THEN a identidade é persistida
  (associação `codigo_produto ↔ catalog_product_id`) E uma entrada é enfileirada em `sync_state`
  (E8, entity=`market`, cadence-agnostic conforme D6) — isto é o "caminho de mercado" que quebra
  o chicken-egg.
- WHEN a coleta é enfileirada THEN NENHUM write ocorre em `market_aggregates`/
  `competitor_offers` nesta milestone — só o enqueue (MC-11, boundary explícito com a missão
  mercado, que é quem drena a fila e executa).
- IF o produto já tem vínculo (`product_links` existe, `LinkedListings` não-vazio) THEN o
  caminho normal de M-05/coleta-via-vínculo já cobre — F-02 só atua no caso "sem anúncio",
  não duplica enqueue.

**Inputs/Outputs (MUST have):**
- Persistência: `catalog_product_id` associado ao `codigo_produto` numa **tabela dedicada NOVA,
  propriedade de M-07** (ex. `product_catalog_identity(tenant_id, codigo_produto,
  catalog_product_id, discovered_at)`), tenant-scoped. **NUNCA ALTER em `products_mirror`** (shape
  de M-02) — evita colidir com superfície de outro milestone (ver ★3 seam; sem grant aditivo em
  tabela alheia). Não colide com `product_links` (vínculo a LISTING, não a catalog_product).
- `sync_state` row nova/atualizada seguindo o shape E8 SEM alteração: chave
  `(tenant_id, installation_id, entity=market)` (E8 é uma row por instalação/entity, NÃO por
  produto); o `codigo_produto`/`catalog_product_id` do produto descoberto é acumulado no
  `cursor JSONB`, `schedule` genérico (D6 não hardcoda cadência).

**Negative Scenarios:**
- Produto já enfileirado (re-run de discovery) → idempotente: o `codigo_produto` não é
  re-adicionado ao `cursor` da row `sync_state (tenant_id, installation_id, entity=market)` (mesma
  disciplina A8/M-05); NÃO cria uma 2ª row de `sync_state` (E8 não tem coluna `codigo_produto`).
- Produto com `catalog_product_id` vazio (não-catalogável) → nada é enfileirado; row de
  `products_mirror` permanece sem mercado-path nesta missão (honest-unknown, não erro).
- Diff da milestone toca `market_aggregates`/`competitor_offers` (write) → viola MC-11, é
  anti-critério (falha imediata de revisão).

**Write-set:** `internal/modules/market/application/*` (enqueue); migração **bloco C**
(condicional — só se o Live Gate PASS e F-02 materializar) criando a tabela dedicada NOVA
`product_catalog_identity` — aditiva, propriedade de M-07, nunca ALTER em `products_mirror`.

---

### F-03 — `cmd/mlprobe` evidência (round 3, T13-T16)

**EARS:**
- WHEN o Live Gate roda THEN `cmd/mlprobe` é estendido (não reescrito) com os 4 testes T13-T16,
  seguindo o padrão dos rounds anteriores (`main.go`/`followup.go`, T0-T12 já presentes em
  `docs/design/evidence/ml-api/`).
- WHEN qualquer teste roda THEN o output vai para `docs/design/evidence/ml-api/T1{3-6}-*.json`
  — nunca para stdout do chat/log persistido em texto livre com credencial embutida (AC-05).

**Inputs/Outputs (MUST have):**
- Extensão de `apps/server_core/cmd/mlprobe/{main.go,followup.go}` com os 4 casos T13-T16.
- 4 arquivos de evidência (ou menos, se algum `could-not-run` — nomear o bloqueio no lugar do
  arquivo ausente).

**Negative Scenarios:**
- `mlprobe` lido/mantido untracked — não vira parte do build de produção; se hub decidir
  promovê-lo, é decisão fora desta milestone (ADR-06 nota "diferido").
- Evidência com token/credencial visível em qualquer campo do JSON → falha de AC-05, deve ser
  redigido antes de commitar.

**Write-set:** `apps/server_core/cmd/mlprobe/*` (untracked, evidência não é código de produto),
`docs/design/evidence/ml-api/T13-*.json` .. `T16-*.json`.

## Ownership & Concurrency (six-axis)

| Eixo | M-07 |
|------|------|
| Migração | **bloco C** (condicional — só se Live Gate PASS): tabela NOVA `product_catalog_identity`, propriedade de M-07; nenhuma se gate FAIL/build não ocorrer; **nunca ALTER `products_mirror`** (M-02) |
| DB shape | tabela dedicada `product_catalog_identity` (`codigo_produto↔catalog_product_id`, M-07-owned, se PASS); `sync_state` entry (entity=market) reusa E8 shape `(tenant,installation_id,entity)` de M-01, não redefine e não adiciona coluna |
| Módulo Go | `market/application` — metade "discovery/enqueue"; **execução** de coleta (metade que falta) é fora, missão ML-sync — dono desta metade só |
| `root.go` | nenhum toque esperado (discovery é chamada sob demanda, não boot-time wiring) |
| Contrato/SDK | nenhum endpoint novo previsto nesta milestone (discovery é interno, chamado por enqueue, não por FE) — se M-06 quiser expor status na tela, é FK/leitura de `sync_state` já exposta por M-06, não uma seção nova de M-07 |
| FE surface | nenhuma — M-06 já cobre a visualização de `sync_state`/cadeia |

M-07 roda **∥ ondas 2-4** (paralelo a M-03/M-04/M-05/M-06), gate real é temporal (live round),
não estrutural — pode ser dispatchada cedo mas só materializa código pós-veredito do gate.

## Dependencies

- **M-02** (`products_mirror` deve existir para `descobrir_produto_catalogo` ter o que
  enriquecer; `sync_state` de M-01 deve existir para o enqueue ter alvo).
- **Live Gate T13-T16** (REQUEST hub) — dependência não-estrutural, temporal/decisória; bloqueia
  só os Feature Briefs, não a rodada do gate em si.

## Validation

Critérios de missão que M-07 é dono:
- **MC-10**: prova live T13-T16 decide construir-ou-remover; sem execução de coleta.
- **MC-11** (co-dono, com todas as milestones): boundary respeitado — nenhum write a
  `market_aggregates`/`competitor_offers` no diff.

Detalhe binário → `M-07-f37-discovery/validation-contract.md`.
