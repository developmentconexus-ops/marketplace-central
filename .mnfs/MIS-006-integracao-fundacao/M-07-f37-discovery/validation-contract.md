# Milestone Validation Contract — M-07-f37-discovery

```yaml
id: M-07-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-07
validation_level: QA-0
conditional: true
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização.

**Estrutura condicional:** M07-C1 a M07-C3 (Live Gate) são SEMPRE exigíveis — a rodada tem que
rodar (ou nomear bloqueio). M07-C4 em diante (build) só se aplicam SE o gate = PASS; se FAIL,
eles ficam `n/a — gate FAIL` (não é gap, é o próprio veredito honesto do fork).

## Critérios — Live Gate (sempre exigíveis)

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M07-C1 | Rodada T13-T16 executada OU bloqueio nomeado, via REQUEST hub (chip não roda direto) | 4 arquivos `docs/design/evidence/ml-api/T13-*.json`..`T16-*.json` presentes E cada um `ran` OU nomeado `could-not-run` com motivo (ex. "REQUEST sem resposta do hub") | L2 ran/could-not-run | Live Gate |
| M07-C2 | EANs de origem = protocolo #004-E real (não sintéticos) | evidência T13 mostra os EANs usados batem com `SELECT ean FROM erp_import_products WHERE protocol_id='#004-E'` (amostra de 10) | L2 ran | Live Gate |
| M07-C3 | Credencial NUNCA exposta/impressa (AC-05) | grep nos 4 arquivos de evidência + em qualquer log/output do chip: 0 hits de token/access_token/refresh_token em texto plano | L0 ran (grep) | Live Gate |
| M07-C4 | Decisão fork registrada em `mission.md`/`interface-contracts-mis006.md` (PASS→build ou FAIL→REMOVE honest-unknown), independente do resultado | diff de `mission.md` §Decisions Resolution linha F3.7 mostra veredito final (não mais "PENDENTE"); se FAIL, `interface-contracts-mis006.md` §E7 marca `E7.1-partial REMOVIDO` | L1 ran | Live Gate |

## Critérios — Build (só se M07-C1..C4 = PASS; senão `n/a — gate FAIL`)

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M07-C5 | `descobrir_produto_catalogo(ean)` existe, read-only, retorna `catalog_product_id` ou vazio | chamada com EAN catalogável de #004-E → retorna id não-vazio; chamada com EAN não-catalogável → retorna vazio, não erro; grep no código: zero `PUT`/`POST`/`PATCH` para a API ML dentro da função | L2 ran | F-01 |
| M07-C6 | 403 PolicyAgent é propagado tipado, não mascarado como vazio | teste/evidência: resposta 403 simulada ou real (se T15 confirmou o padrão) resulta em erro tipado distinguível de "vazio legítimo" | L1 ran | F-01 |
| M07-C7 | Identidade persistida em tabela dedicada M-07-owned, sem ALTER em `products_mirror` | pós-discovery, `SELECT ... FROM product_catalog_identity WHERE tenant_id=? AND codigo_produto=?` retorna a associação; `git diff` da migração = zero `ALTER TABLE products_mirror` (shape de M-02 intocada) | L2 ran + L0 diff | F-02 |
| M07-C8 | Coleta enfileirada em `sync_state`, NUNCA executada (MC-11 boundary) | `sync_state` tem row `entity=market` chave `(tenant_id, installation_id, entity)` com o produto no `cursor`; `git diff` do milestone = zero linhas tocando `market_aggregates`/`competitor_offers` (write) | L1 ran + L0 diff | F-02 |
| M07-C9 | Idempotência do enqueue (shape E8, uma row por instalação/entity) | 2× discovery+enqueue do mesmo produto → o `codigo_produto` aparece 1× no `cursor` da row `sync_state (tenant_id, installation_id, entity=market)`; NÃO cria 2ª row (E8 não tem coluna `codigo_produto`) | L1 ran | F-02 |
| M07-C10 | `cmd/mlprobe` estendido (não reescrito) com T13-T16 | diff mostra adição incremental em `main.go`/`followup.go` seguindo padrão dos rounds T0-T12 existentes, sem remoção dos testes anteriores | L1 ran | F-03 |

## Anti-critérios (falha se presente)

- AC-03: EAN não-catalogável vira erro ou valor fabricado em vez de vazio honesto (M07-C5).
- AC-04: stub da resposta ML no lugar da prova live T13-T16 sem autorização — o Live Gate É a
  prova de integração real; stub aqui invalida MC-10 inteiro.
- AC-05: credencial ML exposta/impressa em qualquer evidência, log, ou output de chat (M07-C3).
- AC-06: push para remote (profile §9).
- AC-01: qualquer query nova (incl. `product_catalog_identity` e leitura/escrita de `sync_state`)
  sem filtro `tenant_id` — a nova tabela é tenant-scoped, toda query DEVE filtrar por `tenant_id`.
- **Anti-critério específico M-07**: qualquer write a `market_aggregates`/`competitor_offers` no
  diff desta milestone — viola MC-11 mesmo que o gate tenha PASSADO (execução de coleta nunca é
  escopo de MIS-006, PASS só habilita descoberta+enqueue).

## Nota sobre o veredito FAIL (não é gap)

Se T13-T16 disprovar F3.7 (ex. T13 retorna sistematicamente vazio, ou T15 confirma 403
PolicyAgent sem T14 suprir alternativa suficiente), M07-C5 a M07-C10 ficam `n/a — gate FAIL`
registrado explicitamente aqui (não omitidos, não escondidos como "não testado"). O veredito da
milestone nesse caso é **CLOSED-BY-DISPROOF**: objetivo cumprido (decisão tomada com evidência),
entrega parcial (só link-path via M-05) é o resultado ESPERADO da condicionalidade, não uma
falha de execução. `mission.md` §Outcome já antecipa este desfecho ("produto sem anúncio recebe
só caminho de vínculo").

## Critérios de user-drive (AMENDMENT D-120 — obrigatório, ratificado pelo operador)

Mesma regra ratificada em M-03 (origem: regressão /catalogo 503 invisível aos gates de código,
hub-fix @2567eb44).

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M07-U1 | Se gate live T13-T16 PASS: produto SEM anúncio vinculado ganha caminho de mercado visível na UI (oportunidade/coleta aparece como o usuário veria); se CLOSED-BY-DISPROOF: n/a — gate FAIL nomeado, nunca omitido | browser drive de 1 produto sem vínculo OU registro explícito do disproof |
