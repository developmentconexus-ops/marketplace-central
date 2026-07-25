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
| M05-C2 | Candidato EAN-único-sem-SKU vai para CONFIRMAÇÃO, NÃO auto-aprova (D-121-2) | candidato com `ean=7909251260214` e sem SKU → estado CONFIRMAÇÃO + aviso `sem CODPROD para corroborar o EAN`; `product_links` SEM vínculo aprovado e ZERO linha E10 `actor=system` até o operador confirmar | L2 ran | F-02 |
| M05-C3 | Auto-aprovação grava linha E10 correta | `SELECT rule_matched, actor, collisions_at_decision FROM <tabela_audit_E10> WHERE link_id = ...` retorna `actor=system, rule_matched=concordant_codprod_ean, collisions_at_decision=1` no caso de C15 — e NENHUMA linha `actor=system` existe para os casos de C2/C14 (âncora única) | L2 ran | F-03 |
| M05-C4 | EAN colidente (>1 produto no ERP) NÃO auto-aprova — fica REVIEW | fixture com 2 produtos mesmo EAN → candidato gerado permanece `status=REVIEW`; zero linha E10 com `actor=system` para esse par. Caso real no ERP do operador: `7896902180697` casa 4 produtos | L2 ran | F-02 |
| M05-C5 | Listing sem EAN e sem SKU utilizável fica REVIEW com motivo visível | candidato sem âncora → `status=REVIEW` + campo de motivo não-vazio (honesto, não silencioso) | L1 ran | F-02 |
| M05-C6 | Contagem de colisão é a do gerador, não uma nova nem a do xlsx | diff de `product_links/*` não introduz nova função de contagem; a decisão lê `len(skuMatches.Products)`/`len(eanMatches.Products)` de `buildExactCandidates` (`generation_service.go:194-222`). `validEANCounts`/`identityQuality` NÃO é consumido aqui — mede duplicidade DENTRO do arquivo xlsx, fato diferente (correção D-121); grep provando ausência dessa dependência em `product_links/*` | L1 ran | F-02 |
| M05-C14 | CODPROD-único SEM EAN vai para CONFIRMAÇÃO, NÃO auto-aprova (D-121-2) | `seller_sku` = CODPROD válido e único, sem EAN → estado CONFIRMAÇÃO + aviso `sem EAN para corroborar o CODPROD`; zero vínculo aprovado, zero E10 `actor=system` | L2 ran | F-02 |
| M05-C15 | CODPROD + EAN concordantes é o ÚNICO caminho automático | as duas âncoras apontando o mesmo produto → aprovado + E10 `rule_matched=concordant_codprod_ean, actor=system` | L2 ran | F-02 |
| M05-C21 | CONFIRMAÇÃO e REVIEW são estados distintos e contáveis separados | o candidato de C2/C14 e o candidato de C4/C16 estão em estados diferentes na leitura do transport; um filtro/contagem por estado devolve os dois grupos separados, com o aviso de C2/C14 preservado. Must-fail: mapear CONFIRMAÇÃO para o mesmo valor de REVIEW quebra este critério | L1 ran | F-02 |
| M05-C22 | Confirmação do operador aprova com trilha honesta | operador confirma o candidato de C14 → vínculo aprovado + E10 `actor=operator, rule_matched=exact_codprod_unique` (nunca `actor=system`, nunca `concordant_codprod_ean`) | L1 ran | F-03 |
| M05-C16 | CODPROD e EAN apontando produtos DIFERENTES → REVIEW, sem precedência | fixture de conflito → `status=REVIEW`, zero linha E10 `actor=system`; o teste prova que nenhuma âncora venceu a outra | L1 ran | F-02 |
| M05-C17 | Hard-negative bloqueia mesmo com âncoras concordantes | fixture título "KIT 2 UN" vs produto unitário, CODPROD+EAN concordantes → NÃO auto-aprova; must-fail: remover a checagem de `detectHardNegative` faz o teste passar a auto-aprovar | L1 ran | F-02 |
| M05-C18 | `seller_sku` casa `p.CODPROD`, não `p.REFFORN` | SQL gerado pelo matcher contém `p.CODPROD = :n` para o input de SKU e nenhuma cláusula `p.REFFORN`; teste de reader casando um CODPROD real do ERP | L1 ran | F-04 |
| M05-C19 | `seller_sku` inválido não vira cláusula de match | input `seller_sku="L.87.22"` (REFFORN legado) → nenhuma cláusula de SKU na query (mesma disciplina do `IsValidGTIN` para EAN); candidato decide por EAN ou fica `unresolved` | L1 ran | F-04 |
| M05-C20 | Um produto ↔ N anúncios: sem limite, sem flag | 2 anúncios distintos com o mesmo CODPROD → 2 vínculos aprovados, nenhum marcado como suspeito/pendente por duplicidade de produto | L1 ran | F-02 |
| M05-C7 | Handlers HTTP órfãos permanecem registrados e funcionais (KEEP) | `product_links/transport/http_handler.go:89-90` inalterado em assinatura; `curl POST /product-links/link-candidates/generations` continua retornando 200 pós-milestone | L1 ran | F-01 |
| M05-C8 | Idempotência de geração+approve: 2× trigger (re-import/re-sync) da mesma identidade de anúncio não duplica vínculo | roda import 2× sobre o mesmo snapshot → `SELECT count(*) FROM product_links WHERE (tenant_id, installation_id, provider_item_id, provider_variation_id) = (...)` retorna 1 | L1 ran | F-03 |
| M05-C9 | A8 satisfeita pela PK EXISTENTE — nenhuma constraint nova em `product_links` | `\d product_links` mostra a PK `(tenant_id, installation_id, provider_item_id, provider_variation_id)` e o diff da migração NÃO altera `product_links`. Duas variações do mesmo anúncio geram 2 vínculos (a chave que o plano original pedia teria colapsado as duas — correção D-121) | L1 ran | F-03 |
| M05-C10 | Override manual do operador vence auto-aprovação e não é revertido | candidato auto-aprovado (`actor=system`) → operador aprova manualmente de novo (ou reatribui) → nova linha E10 `actor=operator` + linha anterior recebe `superseded_by` apontando pra ela; re-run do trigger automático NÃO cria terceira linha revertendo pro `system` | L1 ran | F-03 |
| M05-C11 | Falha de geração não derruba o import/sync que a disparou | fixture que força erro no matcher durante a geração pós-import → import permanece `completed` (não regride pra falho), erro de geração é logado isoladamente | L1 ran | F-01 |
| M05-C12 | E10 shape completa conforme contrato | `\d <tabela_audit_E10>` mostra `link_id, rule_matched, actor, collisions_at_decision, created_at, superseded_by NULL`; `rule_matched` aceita `exact_codprod_unique | exact_ean_unique | concordant_codprod_ean | manual` | L1 ran | F-03 |
| M05-C13 | Migração E10 é bloco B+ (aplicada após bloco B de M-02) e só-aditiva | `git diff` da migração mostra só `CREATE TABLE` (zero `ALTER TABLE product_links`); ordem de bloco confirmada (roda depois de `products_mirror`/`active_source` no histórico de migrações) | L0 ran (diff) | F-03 |

## Critérios de missão cobertos (mapeamento)

- **MC-05** (`validation-contract.md` da missão) — coberto por M05-C1..C6 e C14..C22.
  Único caminho automático (concordante) → M05-C15 + C3. Fila de CONFIRMAÇÃO → M05-C2
  (EAN sem CODPROD), M05-C14 (CODPROD sem EAN), M05-C21 (é estado distinto de REVIEW), M05-C22
  (confirmação humana com trilha honesta). Os três caminhos de REVIEW → M05-C4 (colisão),
  M05-C16 (conflito CODPROD≠EAN), M05-C17 (hard-negative). A âncora de SKU corrigida →
  M05-C18/C19.
- **MC-06** — coberto por M05-C8, M05-C9, M05-C10, M05-C20. Idempotência: 2× sync mesma
  identidade de anúncio → 1 link (M05-C8/C9); mesmo produto em 2 anúncios → 2 links, e isso NÃO
  é violação de idempotência (M05-C20). Override manual vence (M05-C10).

## Anti-critérios (falha se presente)

- AC-01: linha E10 ou `product_links` gravada sem `tenant_id` (auto-approve OU manual).
- AC-03: `collisions_at_decision` ou qualquer campo E10 gravado como `0`/default quando o dado
  real é desconhecido — viola ADR-17 se o sinal de colisão não foi de fato lido.
- AC-04: contagem de colisão própria/stub dentro de `product_links` (M05-C6), OU consumo de
  `validEANCounts`/`identityQuality` como se fosse colisão de ERP — ambos = falha, mesmo que
  produzam resultado aparentemente correto.
- AC-05: auto-approve rodando sobre âncora ambígua (qualquer âncora casando >1 produto), sobre
  conflito CODPROD≠EAN, ou sobre candidato com hard-negative — segurança > cobertura é
  inegociável (ADR-05 amendado); qualquer instância = falha automática de M05-C4/C16/C17.
- AC-08 (D-121): regra de precedência entre CODPROD e EAN em conflito — o operador ratificou
  "ninguém vence, vai para revisão". Código que eleja um vencedor = falha, ainda que documentado.
- AC-09 (D-121): qualquer unique/limite que impeça um produto de ter N vínculos ativos para N
  anúncios distintos.
- AC-10 (D-121-2): candidato de âncora única colapsado no mesmo estado/contador de REVIEW, ou
  auto-aprovado sem corroboração. As duas filas são distintas por decisão do operador; qualquer
  uma das duas violações = falha.
- AC-11 (D-121-2): aviso de confirmação genérico ("verifique") em vez de nomear a âncora que
  faltou — viola ADR-17 (motivo honesto e específico).
- AC-06: push para remote (profile §9).
- AC-07 (específico M-05): endpoint HTTP órfão removido/desativado em vez de KEPT — viola escopo
  explícito da mission spec ("KEEP handlers órfãos, add invocação interna").

## Gate de milestone

M-05 fecha `validated` só com todos os M05-C1..C22 em `ran` (ou `could-not-run` nomeado com
bloqueio explícito, nunca silenciosamente pulado) E zero anti-critério presente. Evidência de
IO A (90008) e do fixture EAN-colidente deve citar dados reais de dev stack ou fixture versionada
— nunca verdict assumido sem prova para MC-05/MC-06 (são L2/L1 `ran`, não `assumed`).

Estado do dev stack em D-121, medido pelo hub (define o que é fixture vs live): `products_mirror`
10529 linhas source=sankhya (7361 com EAN), `listings` 0, `product_link_listing_snapshots` 0,
`product_link_candidates` 0, `product_links` 0, `integration_installations` 1 linha em
`pending_connection`. Consequência: os critérios de código rodam contra fixture versionada
(permitido pelo gate), e tudo que exige anúncio real é `could-not-run` até o operador autorizar a
conta ML.

## Critérios de user-drive (AMENDMENT D-120 — obrigatório, ratificado pelo operador)

Mesma regra ratificada em M-03 (origem: regressão /catalogo 503 invisível aos gates de código,
hub-fix @2567eb44): fechamento exige dirigir o dev stack como usuário nas telas EXISTENTES.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M05-U1 | Pós-import real: vínculos corroborados (CODPROD+EAN) aparecem AUTO-APROVADOS na tela de vínculos como o usuário vê; contadores da tela batem com o DB | browser drive + SELECT de conferência |
| M05-U2 | Produto sem âncora (ou EAN colidente, ou conflito CODPROD≠EAN) aparece como pendente/sem vínculo na tela — nunca falso auto-vínculo | browser drive citando 1 caso real de cada |
| M05-U4 | Candidato de âncora única aparece como CONFIRMAÇÃO, com o produto já proposto e o aviso do que faltou — visualmente separado da fila de revisão | browser drive citando 1 caso de cada fila |
| M05-U3 | /anuncios reflete o vínculo novo (coluna produto/vs-mercado) sem intervenção manual além do import | browser drive /anuncios pós-import |

**Bloqueio conhecido (D-121):** M05-U1 e M05-U3 exigem anúncio real, e a conta ML ainda não está
autorizada (`integration_installations` = `pending_connection`, `listings` = 0). O chip fecha os
dois como `could-not-run` COM esse bloqueio nomeado — nunca como PASS por inspeção de código. O
hub re-dirige U1/U3 no dev stack depois que o operador completar o OAuth. M05-U2 e M05-U4 são
executáveis contra fixture na tela de vínculos — U4 é a prova de que CONFIRMAÇÃO e REVIEW se
distinguem para o usuário, não só no banco.
