# Milestone Validation Contract — M-02-mirror-port-active-source

```yaml
id: M-02-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-02
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização.

## Critérios

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M02-C1 | `products_mirror` existe com shape E2.1 completo (10 campos de negócio + PK + flags) | migração aplicada; `\d products_mirror` (ou introspecção equivalente) mostra `tenant_id, source, codigo_produto, descricao, referencia, ean, marca, grupo_codigo, grupo_descricao, ncm, custo, preco_venda, estoque_total, protocol_id, absent_in_last_snapshot, stale_since, updated_at`; PK = `(tenant_id, codigo_produto)` | L1 ran | F-01 |
| M02-C2 | `products_mirror_stock_locations` existe, escopada por `(tenant_id, codigo_produto)` | `\d products_mirror_stock_locations` mostra `tenant_id, codigo_produto, local_codigo, local_descricao, quantidade` | L1 ran | F-01 |
| M02-C3 | `custo`/`preco_venda`/`estoque_total` são honesto-NULL, nunca 0 por schema (ADR-17) | grep na migração: 0 hits de `DEFAULT 0` ou `NOT NULL` nessas 3 colunas; coluna aceita `NULL` explicitamente | L1 ran | F-01 |
| M02-C4 | `absent_in_last_snapshot`/`stale_since` existem com defaults corretos (F-XLSX-1 base) | DDL mostra `absent_in_last_snapshot BOOL DEFAULT false`, `stale_since timestamptz NULL` (sem default, honesto até o primeiro ausente) | L1 ran | F-01 |
| M02-C5 | `NormalizedRow` estendido para os 10 campos E2 (não mais 6) | diff de `erp_import/domain/import.go` mostra `local`/`preco_venda`/`grupo_codigo`/`grupo_descricao` adicionados aos 6 campos existentes | L1 ran | F-01 |
| M02-C6 | Port `ProductSourceAdapter` compila e embute `readports.Reader` sem alterar assinaturas existentes | `go build ./...` verde; diff mostra a interface embutindo o tipo existente, zero mudança de assinatura em `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs` | L1 ran | F-02 |
| M02-C7 | `Sync()`/`Kind()` declarados no port (sem implementação concreta obrigatória nesta milestone) | assinatura `Sync(ctx) (SyncResult, error)` e `Kind() SourceKind` presentes na interface | L1 ran | F-02 |
| M02-C8 | `SourceKind` é tipo dedicado com exatamente `{upload_snapshot, live_read_through}`, não reaproveita `ImportSource` | grep do tipo novo; `ImportSource` (`erp_import/domain/import.go:64-67`) permanece com seus 2 valores originais, sem merge de conceito | L1 ran | F-02 |
| M02-C9 | `active_source` config existe, PK `tenant_id`, shape E9.0 | `\d active_source` mostra `tenant_id [PK], active_source, source_kind, set_at, set_by` | L1 ran | F-03 |
| M02-C10 | `MC_ERP_SOURCE` removido — zero hits no código | `grep -rn MC_ERP_SOURCE apps/server_core` retorna 0 | L0 ran (grep) | F-03 |
| M02-C11 | `composition/root.go` constrói AMBOS adapters simultaneamente, sem boot-time branch de fonte única | diff de `root.go` mostra `erpSource(getenv)` (linha ~772) removido; wiring instancia reader xlsx E reader oracle na composição, sem `if source == X` decidindo qual existe | L1 ran | F-03 |
| M02-C12 | Fail-closed: tenant sem `active_source` retorna erro tipado, nunca fallback silencioso | teste unitário: resolução para `tenant_id` sem row em `active_source` retorna `ErrUnknownActiveSource` (reusado de `reader.go:53`); teste passa | L1 ran | F-03 |
| M02-C13 | `active_source` do ctx acessível para downstream (pré-requisito de cache-key, extensão real fica p/ M-04) | grep/leitura mostra função de leitura de `active_source` a partir do `ctx` exportada/consumível fora do pacote de resolução | L1 ran | F-03 |
| M02-C14 | Migração é só-aditiva | `git diff` da migração mostra somente `CREATE TABLE`/`ADD COLUMN`; zero `ALTER COLUMN ... TYPE`, zero `DROP` | L0 ran (diff) | F-01,F-03 |
| M02-C15 | OpenAPI + `sdk-runtime` do endpoint active-source landam no mesmo commit | `git show --stat <commit>` do PR/chip mostra spec OpenAPI E arquivo(s) `sdk-runtime` gerado no mesmo commit | L1 ran | F-04 |
| M02-C16 | Endpoint active-source rejeita valor fora do enum ratificado | teste: PUT com `active_source="foo"` → 400, nada gravado | L1 ran | F-04 |

## Anti-critérios (falha se presente)

- AC-01: query tenant-scoped sem `tenant_id` (mirror, stock_locations, active_source — todas
  as 3 tabelas novas desta milestone são tenant-scoped).
- AC-03: coluna nova com `DEFAULT 0` fingindo honesto-desconhecido (M02-C3).
- AC-04: adapter concreto stubado só p/ satisfazer o port sem autorização — NÃO É ESCOPO de
  M-02 (adapters ficam p/ M-03/M-04); se algum stub aparecer aqui, é escopo vazando, não prova.
- AC-06: push para remote (profile §9).

## Nota de escopo (não confundir com MC-01 pleno)

M02-C1/C2 provam SHAPE, não FEEDING. `MC-01` da missão ("alimentado por AMBOS adapters") só
fecha quando M-03 (xlsx Sync) e M-04 (Sankhya Sync) escreverem rows reais. M-02 sozinho não tenta
provar `SELECT ... FROM products_mirror WHERE source IN (...)` com dados reais — isso é `could-not-run`
nomeado aqui até M-03/M-04 landarem, não um gap a esconder.
