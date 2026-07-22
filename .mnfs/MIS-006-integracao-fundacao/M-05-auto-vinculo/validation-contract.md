# Milestone Validation Contract — M-05-auto-vinculo

```yaml
id: M-05-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-05
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização.

## Critérios

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M05-C1 | Import xlsx (IO A, produto 90008 EAN único) dispara geração de candidatos internamente, sem passar pelo endpoint HTTP | pós-import (dev stack), `product_links` candidate row existe para o par 90008/listing correspondente SEM ter chamado `POST /product-links/link-candidates/generations` manualmente na sessão de teste | L2 ran | F-01 |
| M05-C2 | Candidato EAN-exato-único (`collisions[ean]==1`) auto-aprova | `SELECT status FROM product_links WHERE ... = 90008` retorna aprovado, sem nenhuma chamada a `ApproveCandidate`/endpoint manual na sessão de teste | L2 ran | F-02 |
| M05-C3 | Auto-aprovação grava linha E10 correta | `SELECT rule_matched, actor, collisions_at_decision FROM <tabela_audit_E10> WHERE link_id = ...` retorna `rule_matched=exact_ean_unique, actor=system, collisions_at_decision=1` | L2 ran | F-03 |
| M05-C4 | EAN duplicado (`collisions[ean]>1`) NÃO auto-aprova — fica REVIEW | fixture com 2 produtos mesmo EAN → candidato gerado permanece `status=REVIEW`; zero linha E10 com `actor=system` para esse par | L2 ran | F-02 |
| M05-C5 | Listing sem EAN fica REVIEW com motivo visível | candidato de listing sem `ean` → `status=REVIEW` + campo de motivo não-vazio (honesto, não silencioso) | L1 ran | F-02 |
| M05-C6 | Sinal de colisão é REUSADO, não reimplementado | diff de `product_links/*` não introduz nova função de contagem de EAN; grep mostra chamada a `validEANCounts`/`identityQuality` (`erp_import/adapters/internalread/reader.go:344-366`) ou ao equivalente Oracle (`internal_read/adapters/oracle/reader.go:70-76`) a partir de `generation_service.go`/`resolution_service.go` | L1 ran | F-02 |
| M05-C7 | Handlers HTTP órfãos permanecem registrados e funcionais (KEEP) | `product_links/transport/http_handler.go:89-90` inalterado em assinatura; `curl POST /product-links/link-candidates/generations` continua retornando 200 pós-milestone | L1 ran | F-01 |
| M05-C8 | Idempotência de geração+approve: 2× trigger (re-import/re-sync) do mesmo produto/listing não duplica vínculo | roda import 2× sobre o mesmo snapshot → `SELECT count(*) FROM product_links WHERE (internal_product_id, provider_listing_id) = (...)` retorna 1 | L1 ran | F-03 |
| M05-C9 | Constraint única A8 existe em `product_links` | `\d product_links` (ou introspecção equivalente) mostra unique constraint em `(internal_product_id, provider_listing_id)` | L1 ran | F-03 |
| M05-C10 | Override manual do operador vence auto-aprovação e não é revertido | candidato auto-aprovado (`actor=system`) → operador aprova manualmente de novo (ou reatribui) → nova linha E10 `actor=operator` + linha anterior recebe `superseded_by` apontando pra ela; re-run do trigger automático NÃO cria terceira linha revertendo pro `system` | L1 ran | F-03 |
| M05-C11 | Falha de geração não derruba o import/sync que a disparou | fixture que força erro no matcher durante a geração pós-import → import permanece `completed` (não regride pra falho), erro de geração é logado isoladamente | L1 ran | F-01 |
| M05-C12 | E10 shape completa conforme contrato | `\d <tabela_audit_E10>` mostra `link_id, rule_matched, actor, collisions_at_decision, created_at, superseded_by NULL` | L1 ran | F-03 |
| M05-C13 | Migração E10/A8 é bloco B+ (aplicada após bloco B de M-02) e só-aditiva | `git diff` da migração mostra só `CREATE TABLE`/`ADD CONSTRAINT`; ordem de bloco confirmada (roda depois de `products_mirror`/`active_source` no histórico de migrações) | L0 ran (diff) | F-03 |

## Critérios de missão cobertos (mapeamento)

- **MC-05** (`validation-contract.md` da missão) — coberto por M05-C1, M05-C2, M05-C3, M05-C4.
  IO A: produto 90008 EAN único → candidato existe E `product_links` aprovado + E10
  `actor=system, rule=exact_ean_unique` (M05-C1/C2/C3). EAN duplicado → REVIEW (M05-C4).
- **MC-06** — coberto por M05-C8, M05-C9, M05-C10. Idempotência: 2× sync mesmo produto → 1 link
  (M05-C8/C9). Override manual vence (M05-C10).

## Anti-critérios (falha se presente)

- AC-01: linha E10 ou `product_links` gravada sem `tenant_id` (auto-approve OU manual).
- AC-03: `collisions_at_decision` ou qualquer campo E10 gravado como `0`/default quando o dado
  real é desconhecido — viola ADR-17 se o sinal de colisão não foi de fato lido.
- AC-04: reimplementação stub do sinal de colisão de EAN (M05-C6) — auto-approve baseado em
  contagem própria não reusada = falha, mesmo que produza resultado aparentemente correto.
- AC-05: auto-approve rodando sobre EAN ambíguo (`collisions[ean]>1`) por qualquer motivo —
  segurança > cobertura é inegociável (ADR-05); qualquer instância = falha automática de M05-C4.
- AC-06: push para remote (profile §9).
- AC-07 (específico M-05): endpoint HTTP órfão removido/desativado em vez de KEPT — viola escopo
  explícito da mission spec ("KEEP handlers órfãos, add invocação interna").

## Gate de milestone

M-05 fecha `validated` só com todos os M05-C1..C13 em `ran` (ou `could-not-run` nomeado com
bloqueio explícito, nunca silenciosamente pulado) E zero anti-critério presente. Evidência de
IO A (90008) e do fixture EAN-duplicado deve citar dados reais de dev stack ou fixture versionada
— nunca verdict assumido sem prova para MC-05/MC-06 (são L2/L1 `ran`, não `assumed`).

## Critérios de user-drive (AMENDMENT D-120 — obrigatório, ratificado pelo operador)

Mesma regra ratificada em M-03 (origem: regressão /catalogo 503 invisível aos gates de código,
hub-fix @2567eb44): fechamento exige dirigir o dev stack como usuário nas telas EXISTENTES.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M05-U1 | Pós-import real: vínculos EAN-único-exato aparecem AUTO-APROVADOS na tela de vínculos como o usuário vê; contadores da tela batem com o DB | browser drive + SELECT de conferência |
| M05-U2 | Produto sem EAN (ou EAN duplicado) aparece como pendente/sem vínculo na tela — nunca falso auto-vínculo | browser drive citando 1 caso real de cada |
| M05-U3 | /anuncios reflete o vínculo novo (coluna produto/vs-mercado) sem intervenção manual além do import | browser drive /anuncios pós-import |
