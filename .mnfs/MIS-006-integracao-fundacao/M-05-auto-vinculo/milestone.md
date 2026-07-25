# M-05 — auto-vínculo

```yaml
id: M-05
type: milestone
mission: MIS-006
status: draft
depends_on: [M-02, M-03]
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Fechar o caminho de vínculo automático produto↔anúncio sem reimplementar nada que já existe e
já funciona: REFACTOR de wiring em `generation_service.go` (trigger interno pós-import/sync, hoje
só alcançável por endpoint órfão) + REFACTOR da âncora de SKU do matcher para `p.CODPROD`
(D-121-1 — o SKU do anúncio no ML É o CODPROD; contra `p.REFFORN` a âncora forte do operador
nunca casaria) + REUSE da contagem de colisão que o próprio gerador já calcula + REFACTOR de
`resolution_service.go` para uma transição auto-approve que reusa a máquina de audit existente,
formalizada como E10 (trail nova). Handlers HTTP órfãos são KEPT (não removidos) — ganham um
caminho de invocação interna adicional.

Política de auto-approve ratificada pelo operador em D-121 (ADR-05 amendado,
`RATIFIED-BY-OPERATOR`), e é ela que os briefs abaixo implementam:

| Sinal no anúncio | Resultado |
|---|---|
| CODPROD válido **e** EAN, mesmo produto | **auto-aprova** (`concordant_codprod_ean`) |
| CODPROD válido, sem EAN | **CONFIRMAÇÃO** — produto proposto + aviso, 1 clique (`exact_codprod_unique`) |
| EAN único, sem CODPROD | **CONFIRMAÇÃO** — produto proposto + aviso, 1 clique (`exact_ean_unique`) |
| CODPROD e EAN apontam produtos diferentes | REVIEW — conflito, sem regra de precedência |
| EAN colidente (>1 produto no ERP) | REVIEW |
| só título | REVIEW, nunca auto-aprova |
| título contradiz (kit/combo/cor/voltagem) | bloqueia tudo acima — hard-negative vence |

**CONFIRMAÇÃO ≠ REVIEW** (D-121-2, ratificado). São dois estados distintos, contados e filtrados
separadamente:

- *confirmação* = "achei o produto, faltou a segunda âncora para corroborar — confirma?". O
  produto já vem escolhido; o operador aprova em 1 clique. Existe produto legítimo cadastrado só
  com CODPROD ou só com EAN, e ele não pode cair na mesma fila de investigação.
- *revisão* = "não sei qual produto é" (colisão, conflito, só título). Exige investigação.

Colapsar os dois em um genérico "pendente" apaga a decisão e é falha de contrato (AC-10).

Por que âncora única deixou de auto-aprovar (D-121-2 revisa D-121): sozinha, ela erra em
silêncio — um CODPROD digitado errado que caia sobre outro código VÁLIDO vincula o produto
errado e ninguém revisa. Com as duas âncoras concordando, o erro teria que acontecer duas vezes
de forma coerente. Trade-off aceito: a cobertura do automático cai e o resto vira 1 clique.

Ver `mission.md` §Milestone Strategy (linha M-05), ADR-05, `interface-contracts-mis006.md` §E4
(E4.1)/§E10, `architecture-map.md` (M-05 depende de M-02+M-03; alimenta M-06),
`research/refactor-inventory-backend.md` §6 (file:line da REFACTOR/REUSE).

## Scope

- REFACTOR `product_links/application/generation_service.go:60-78` (`GenerateLinkCandidates`):
  ADD chamada interna (não-HTTP) disparada ao final de import xlsx (hook de M-03) e de sync
  Sankhya (M-04, quando existir) — lógica de geração em si NÃO muda, só ganha caller automático.
- REFACTOR `internal_read/adapters/oracle/reader.go:451` (D-121-1, `RATIFIED-BY-OPERATOR`):
  `seller_sku` casa `p.CODPROD`, não `p.REFFORN`. No Mercado Livre o SKU do anúncio É o CODPROD
  do ERP — o operador cadastra assim e todos os anúncios já carregam o código. REFFORN (código
  do fabricante) sai como âncora de SKU. O `seller_sku` que não for um CODPROD válido não vira
  cláusula de match (mesma disciplina do EAN não-GTIN em `reader.go:448` — junk nunca alarga a
  query). Rebate no fake reader e no mirror reader do `erp_import` para os dois lados da
  routing manterem a mesma semântica.
- REUSE a contagem de colisão que o gerador JÁ calcula (`len(eanMatches.Products)` /
  `len(skuMatches.Products)` em `buildExactCandidates`, `generation_service.go:194-222`) — é o
  número que a política precisa. NÃO reimplementar, e NÃO reusar `validEANCounts`
  (`erp_import/.../reader.go:344-366`): aquele mede duplicidade de EAN dentro do ARQUIVO xlsx,
  não colisão no ERP — fato diferente (correção de plano D-121).
- REFACTOR `product_links/application/resolution_service.go:129-149` (`ApproveCandidate`): ADD
  transição interna de auto-approve reusando a MESMA máquina de transição/audit hoje só
  acionada por operador manual via transport. Política ratificada (ADR-05 amendado, D-121-2):
  auto-aprova SÓ o concordante (CODPROD e EAN no mesmo produto); âncora única vira CONFIRMAÇÃO
  (proposta + aviso, 1 clique); conflito CODPROD≠EAN, colisão (>1 produto) e hard-negative ficam
  REVIEW.
- ADD o estado de CONFIRMAÇÃO ao domínio de candidatos + um motivo de aviso honesto por caso
  (`sem EAN para corroborar o CODPROD` / `sem CODPROD para corroborar o EAN`). Não invente
  texto genérico: o aviso diz exatamente qual âncora faltou (ADR-17). O transport de
  `product_links` expõe o estado para M-06 filtrar; nenhuma tela é escrita aqui.
- CREATE E10 audit trail (tabela nova, migração bloco B+ — após bloco B de M-02): toda decisão de
  vínculo (manual ou automática) grava linha `rule_matched`, `actor`, `collisions_at_decision`,
  `superseded_by`. `rule_matched` cobre a política ratificada:
  `exact_codprod_unique | exact_ean_unique | concordant_codprod_ean | manual`.
- A8 idempotência: **já satisfeita pela PK existente** de `product_links`
  `(tenant_id, installation_id, provider_item_id, provider_variation_id)` — re-run do trigger
  faz upsert na mesma identidade, não duplica. NÃO criar a chave que o plano original pedia,
  `(internal_product_id, provider_listing_id)`: essa coluna não existe nessa tabela (é
  `provider_item_id`) e a chave proposta perderia a variação, colidindo entre duas variações do
  mesmo anúncio (correção de plano D-121). O que M-05 ADD é a regra de precedência: override
  manual do operador vence auto-aprovação prévia e nunca é sobrescrito de volta pelo automático.
- Um produto ↔ N anúncios sem limite e sem sinalização (D-121, ratificado): mesmo codprod
  anunciado várias vezes é operação normal; auto-approve não checa nada além da identidade do
  anúncio.
- KEEP `product_links/transport/http_handler.go:89-90` (os dois endpoints órfãos) — permanecem
  registrados e curl-acessíveis; ganham companhia do caller interno, não são removidos nem
  substituídos.

## Non-Scope

- Implementação do trigger de sync Sankhya em si (isso é conteúdo de M-04 — M-05 só consome o
  ponto de hook quando M-04 landar; para xlsx o hook já existe via M-03).
- F3.7 discovery EAN→catálogo ML (M-07) — M-05 só trata vínculo produto↔anúncio já existente em
  `listings`, nunca cria/descobre anúncio novo.
- Telas de vínculo / badge auto-aprovado na UI (`VinculosPage`, M-06) — M-05 é 100% backend.
- Qualquer execução de coleta de mercado — MIS-006 só ENFILEIRA (boundary MC-11); auto-vínculo
  não dispara chamada a ML API nenhuma.
- Match fuzzy / por título — o caminho de título continua existindo e continua SEMPRE em REVIEW;
  nenhuma heurística nova de similaridade entra nesta milestone. (A âncora de SKU passou a ser
  ESCOPO em D-121: ver F-04. O que segue fora de escopo é qualquer critério novo além das três
  âncoras já implementadas — sku, ean, título.)
- Correção do cadastro de SKU no Mercado Livre (anúncios cujo SKU não é um CODPROD) — é ação do
  operador no painel do ML, não código. M-05 só garante que um SKU não-CODPROD não vira cláusula
  de match e o candidato cai em REVIEW.

## Feature Briefs

### F-01 — Trigger interno de geração de candidatos (pós-import/sync)

**EARS:**
- WHEN um import xlsx completa (hook de M-03, pós-`import_service.go` merge no mirror) THEN
  `GenerateLinkCandidates` (`generation_service.go:60-78`) é chamado internamente, sem passar
  pelo endpoint HTTP órfão.
- WHEN um sync Sankhya completa (hook de M-04, quando existir) THEN o mesmo trigger interno é
  chamado — o caller é agnóstico à fonte (consome o mirror, não o adapter).
- WHEN o endpoint HTTP órfão (`http_handler.go:89-90`) é chamado diretamente (curl/debug) THEN
  continua funcionando exatamente como hoje — comportamento preservado, só ganha um segundo
  caminho de invocação.
- IF a geração de candidatos falha (erro de matcher, etc.) THEN o import/sync que a disparou NÃO
  falha em cascata — falha de geração é logada e isolada (import/sync já completou com sucesso
  antes do hook disparar).

**Inputs/Outputs (MUST have):**
- Ponto de chamada interno em `product_links/application` (função exportada, ex.
  `TriggerGenerationForTenant(ctx, tenantID, source)`) invocado pelo hook de completion de
  import (M-03) e, futuramente, de sync (M-04).
- `generation_service.go:60-78` inalterado na lógica interna — só ganha um segundo caller além
  do handler HTTP.
- Log/métrica de disparo automático (para diagnosticar `MC-05` — "import dispara geração").

**Negative Scenarios:**
- Import roda sem nenhum listing existente ainda no tenant → geração roda, produz zero
  candidatos, não erra (mirror sem par em `listings` = sem candidato, não é falha).
- Dois imports concorrentes do mesmo tenant → geração não duplica candidatos para o mesmo par
  produto/listing (idempotência de geração já é comportamento hoje via
  `ProductMatcher.FindProductsForLinking`, preservado).
- Handler HTTP chamado manualmente ENQUANTO o trigger interno acabou de rodar → não há corrida
  destrutiva; pior caso é geração redundante, resolvida pela idempotência de F-03 (A8) no
  approve, não aqui.

**Write-set:** `internal/modules/product_links/application/generation_service.go` (novo caller
exportado, lógica interna intocada), ponto de invocação no hook de `erp_import/application`
(M-03 write-set, chamada cross-módulo — decisão de import cycle fica com quem implementa: hook
chama `product_links/application`, nunca o inverso).

---

### F-02 — Auto-approve corroborado + fila de CONFIRMAÇÃO

**EARS:**
- WHEN a geração produz um candidato cujo `seller_sku` casa exatamente 1 produto (CODPROD único)
  E o EAN casa o MESMO produto THEN o candidato transiciona automaticamente para aprovado,
  reusando a máquina de transição de `resolution_service.go:129-149` (`ApproveCandidate`), com
  `actor=system` e `rule_matched=concordant_codprod_ean`. Este é o ÚNICO caminho automático.
- WHEN o `seller_sku` casa exatamente 1 produto E o anúncio não tem EAN válido THEN o candidato
  fica em CONFIRMAÇÃO com o produto já proposto e o aviso `sem EAN para corroborar o CODPROD`
  — NÃO auto-aprova, NÃO cai na fila de revisão.
- WHEN o anúncio não tem `seller_sku` utilizável E o EAN casa exatamente 1 produto THEN o
  candidato fica em CONFIRMAÇÃO com o aviso `sem CODPROD para corroborar o EAN`.
- WHEN o operador confirma um candidato em CONFIRMAÇÃO THEN o vínculo é aprovado com
  `actor=operator` e `rule_matched` = a âncora única que o propôs
  (`exact_codprod_unique`/`exact_ean_unique`) — a linha E10 registra que houve aprovação humana
  sobre âncora não-corroborada.
- WHEN `seller_sku` e EAN casam produtos DIFERENTES THEN o candidato fica em REVIEW como
  conflito — nenhuma das duas âncoras tem precedência sobre a outra (D-121, ratificado).
- WHEN qualquer âncora casa mais de 1 produto (colisão no ERP) THEN o candidato fica em REVIEW —
  auto-approve NUNCA dispara em ambiguidade (segurança > cobertura, ADR-05).
- WHEN o match é só por título THEN vai para REVIEW, nunca para CONFIRMAÇÃO nem auto-approve,
  qualquer que seja o score — título não é âncora.
- IF `detectHardNegative` acusa contradição (kit/combo/cor/voltagem) entre título do anúncio e
  nome do produto THEN o candidato é REJECT/REVIEW mesmo com CODPROD e EAN concordantes, e
  mesmo no caminho de CONFIRMAÇÃO — a contradição vence todas as regras acima (D-121,
  ratificado).
- WHEN o par produto↔anúncio já tem vínculo ativo com `actor=operator` THEN auto-approve não
  roda (ver F-03, override do operador vence).

**Inputs/Outputs (MUST have):**
- A contagem de colisão vem do próprio gerador: `len(skuMatches.Products)` /
  `len(eanMatches.Products)` em `buildExactCandidates`
  (`generation_service.go:194-222`) — já é o número de produtos que a âncora casou no ERP.
  NÃO consumir `validEANCounts`/`identityQuality` (`erp_import/.../reader.go:344-366`): aquele
  conta EANs repetidos DENTRO do arquivo xlsx, não colisão no catálogo — usá-lo aqui seria medir
  o fato errado (correção de plano D-121).
- `collisions_at_decision` gravado em E10 é exatamente esse número (1 no caminho concordante e
  nas confirmações de âncora única).
- Estado de CONFIRMAÇÃO no domínio de candidatos + motivo de aviso por caso, legível pelo
  transport para M-06 filtrar. Contável separado de REVIEW.
- Extensão de `resolution_service.go` com uma transição de auto-approve que reusa o MESMO
  código de mudança de estado + escrita de audit hoje usado pelo caminho manual
  (`ApproveCandidate`), diferindo só no `actor` gravado.

**Negative Scenarios:**
- EAN casando 2 produtos → candidato fica REVIEW, mesmo que um dos dois "pareça" o certo — sem
  heurística de desempate automática. (Dado real do ERP do operador: 91 EANs colidem,
  ex. `7896902180697` casa 4 produtos.)
- `seller_sku` preenchido com algo que não é um CODPROD válido → não vira cláusula de match
  (F-04); se sobrar só o EAN, vira CONFIRMAÇÃO por EAN; se não sobrar âncora, REVIEW
  `unresolved`.
- Anúncio com CODPROD e EAN concordantes MAS título "KIT 2 UN" contra produto unitário →
  hard-negative bloqueia, não auto-aprova.
- Candidato de âncora única entrando na MESMA fila/contador de REVIEW = falha (AC-10): o
  operador não consegue distinguir "confirma?" de "investiga".
- Reimplementação de contagem de colisão dentro de `product_links` (em vez de consumir a que o
  gerador já calcula) = falha de design (ADR-05 "não reimplementar sinal já testado").
- Auto-approve tenta rodar sobre candidato já em estado terminal (aprovado/rejeitado
  manualmente) → não regride nem sobrescreve (ver F-03, override do operador vence).

**Write-set:** `internal/modules/product_links/application/resolution_service.go` (nova função
de transição auto-approve, reusando helpers internos de `ApproveCandidate`),
`internal/modules/product_links/application/generation_service.go` (ponto de chamada da
transição ao final da geração, condicional ao sinal de colisão).

---

### F-03 — E10 audit trail + idempotência A8

**EARS:**
- WHEN qualquer vínculo é aprovado (automático ou manual) THEN uma linha E10 é gravada:
  `link_id, rule_matched, actor (system|operator), collisions_at_decision, created_at,
  superseded_by NULL`.
- WHEN a aprovação é automática (F-02) THEN `actor=system`, `collisions_at_decision=1`,
  `rule_matched=concordant_codprod_ean` — é o único caminho automático (D-121-2).
- WHEN o operador confirma um candidato de âncora única THEN `actor=operator` e `rule_matched` é
  `exact_codprod_unique` ou `exact_ean_unique` — a trilha distingue vínculo corroborado por duas
  âncoras de vínculo aceito por humano sobre uma só.
- WHEN o operador aprova/sobrescreve manualmente um vínculo (inclusive um já auto-aprovado) THEN
  uma NOVA linha E10 é gravada com `actor=operator`, e a linha anterior recebe `superseded_by`
  apontando para a nova — o override do operador VENCE e nunca é revertido de volta pelo caminho
  automático.
- WHEN a mesma identidade de anúncio
  `(tenant_id, installation_id, provider_item_id, provider_variation_id)` é processada de novo
  (re-run de import/sync, F-01 disparando geração outra vez) THEN a PK JÁ EXISTENTE de
  `product_links` impede duplicação — a geração/approve é idempotente, não cria segundo vínculo
  nem segunda linha de audit redundante para o mesmo estado. Nenhum índice novo é criado
  (correção de plano D-121: a chave que o plano pedia,
  `(internal_product_id, provider_listing_id)`, nomeia coluna inexistente e perderia a variação).
- IF `product_links` já tem vínculo ativo para o par E o novo candidato geraria o MESMO resultado
  (mesmo `rule_matched`) THEN nenhuma escrita nova ocorre (no-op idempotente, não erro).

**Inputs/Outputs (MUST have):**
- Tabela E10 nova (migração bloco B+, após bloco B de M-02): shape exata do
  `interface-contracts-mis006.md` §E10 (`link_id, rule_matched, actor, collisions_at_decision,
  created_at, superseded_by NULL`).
- A8 satisfeita pela PK existente
  `(tenant_id, installation_id, provider_item_id, provider_variation_id)` — a milestone PROVA a
  idempotência (teste de re-run), não cria constraint nova. Um produto pode ter N vínculos ativos
  para N anúncios distintos, sem limite (D-121, ratificado) — qualquer unique sobre
  `internal_product_id` sozinho seria um defeito.
- Escrita de audit acoplada à MESMA transação da transição de estado do vínculo (aprovação e
  audit não podem divergir — sem vínculo aprovado sem linha E10 correspondente).

**Negative Scenarios:**
- Override manual sobre vínculo auto-aprovado → 2 linhas E10 no total para o par
  (`system` superseded, `operator` vigente) — nunca 1 linha só sobrescrita in-place (histórico
  imutável, mesmo padrão de `erp_import_protocols`).
- Re-run do trigger duas vezes seguidas sobre o mesmo snapshot → `SELECT count(*) FROM
  product_links WHERE (tenant_id, installation_id, provider_item_id, provider_variation_id) =
  (...)` retorna 1, não 2.
- Duas variações do MESMO anúncio (`provider_variation_id` diferente) → 2 vínculos, e isso é
  correto; a chave do plano original teria colapsado as duas em uma.
- Tentativa de auto-approve escrever por cima de vínculo com `actor=operator` já vigente → é
  bloqueada pela regra de precedência (F-02 nunca chama a transição se já existe vínculo ativo
  com `actor=operator` no par).
- Qualquer escrita de E10 sem `tenant_id` explícito → falha de design (AC-01, todas as tabelas
  novas desta milestone são tenant-scoped via o vínculo pai).

**Write-set:** `apps/server_core/migrations/09xx_product_links_audit.{up,down}.sql` (bloco B+,
tabela E10 SOMENTE — nenhuma alteração em `product_links`),
`internal/modules/product_links/domain` (tipo de audit row),
`internal/modules/product_links/application/resolution_service.go` (escrita de audit acoplada à
transição, manual E automática).

---

### F-04 — Âncora de SKU passa de REFFORN para CODPROD

Brief criado em D-121 pela entrevista com o operador (`RATIFIED-BY-OPERATOR`). Sem ele, F-02 é
inerte no cadastro real: hoje `seller_sku` é comparado com `p.REFFORN` (referência do
fabricante, ex. `L.87.22`), enquanto todo anúncio do operador carrega o CODPROD no campo SKU —
a âncora forte nunca casaria e todo vínculo cairia no caminho fraco do EAN.

**EARS:**
- WHEN o input de match traz `seller_sku` THEN a cláusula gerada é `p.CODPROD = :n`, nunca
  `p.REFFORN = :n` (`internal_read/adapters/oracle/reader.go:451`).
- WHEN `seller_sku` não é um CODPROD sintaticamente válido (não-numérico, vazio, lixo) THEN
  nenhuma cláusula de SKU é adicionada à query — mesma disciplina que `IsValidGTIN` já aplica ao
  EAN em `reader.go:448`; entrada suja nunca alarga a busca.
- WHEN o candidato é gerado a partir do mirror (lado xlsx / routing por `active_source`) THEN a
  mesma semântica vale — os dois readers concordam sobre o que "SKU" significa.
- IF nenhum produto casa o CODPROD informado THEN o resultado é zero matches por SKU (não erro),
  e a decisão recai sobre o EAN ou vira `unresolved`.

**Inputs/Outputs (MUST have):**
- `internal_read/adapters/oracle/reader.go:451` — troca de coluna + guarda de validade.
- Reader espelho do `erp_import` e o fake reader usado nos testes alinhados à mesma coluna, senão
  a suite prova o comportamento errado.
- REFFORN deixa de ser âncora de match; nenhuma outra leitura de REFFORN é removida.

**Negative Scenarios:**
- Anúncio com SKU `L.87.22` (REFFORN legado, não-CODPROD) → sem cláusula de SKU, cai no EAN ou
  em REVIEW; nunca casa por engano.
- Um CODPROD digitado errado que exista como outro produto válido → vincula errado sem revisão.
  Risco explicitamente aceito pelo operador em D-121 (o cadastro do SKU é governado por ele).
- Teste que continue provando `p.REFFORN` = falha de F-04, não baseline.

**Write-set:** `internal/modules/internal_read/adapters/oracle/reader.go`, reader espelho em
`internal/modules/erp_import/adapters/internalread/`, fakes/testes de matcher dos dois pacotes.

## Ownership & Concurrency (six-axis)

| Eixo | M-05 |
|------|------|
| Migração | bloco B+ (E10 audit trail SOMENTE) — após bloco B de M-02 (`products_mirror`/`active_source`); `product_links` NÃO é alterada (A8 já satisfeita pela PK existente) |
| DB shape | tabela de audit E10 — dono; `product_links` lida/escrita sem mudança de shape |
| Módulo Go | `internal/modules/product_links/*` (application: generation_service.go, resolution_service.go; transport: http_handler.go KEPT sem mudança de assinatura) — dono. Além disso, edição PONTUAL da cláusula de match de SKU em `internal_read/adapters/oracle/reader.go` + reader espelho de `erp_import` (F-04) — surface de M-02/M-04 já mergeada, nenhuma wave concorrente escreve nesses arquivos; se alguma abrir, vira additive-lock grant do hub |
| `root.go` | nenhuma mudança de wiring de composição nova (consome geração/resolution já compostos) |
| Contrato/SDK | nenhum endpoint novo nesta milestone (auto-approve é 100% interno; M-06 consome o RESULTADO via chain-read, seção própria) |
| FE surface | nenhuma (M-06 consome via SDK/chain-read) |

Regra de merge: `product_links/*` é escopo exclusivo de M-05 nesta wave — M-02/M-03 não tocam o
pacote (só o hook de completion em `erp_import/application` chama para dentro, na direção
correta, sem import cycle). Migração deve ser alocada pelo hub em bloco POSTERIOR ao bloco B de
M-02 (mirror/active_source já aplicado antes de M-05 rodar).

## Dependencies

- **M-02** — `products_mirror` precisa existir como fonte de candidatos (`FindProductsForLinking`
  lê do mirror via o port formalizado em M-02); sem mirror populado, geração roda mas não produz
  candidatos.
- **M-03** — hook de completion de import xlsx é a ORIGEM concreta do trigger interno (F-01) para
  a primeira fonte disponível; sem M-03, F-01 não tem caller real ainda (só o handler HTTP
  órfão, que já existe hoje independente de M-05).
- Não depende de M-04 (Sankhya) para landar — quando M-04 chegar, adiciona um segundo caller ao
  mesmo trigger interno (`TriggerGenerationForTenant`), sem mudança de shape em M-05.

## Validation

Critérios de missão que M-05 é dono (`validation-contract.md` da missão):

- **MC-05**: cadeia de vínculo automática — import dispara geração; âncora não-ambígua
  (CODPROD-único, EAN-único, ou ambos concordantes) auto-aprova com audit; colisão, conflito
  CODPROD≠EAN e hard-negative → REVIEW.
- **MC-06**: idempotência — re-run não duplica vínculo; override manual vence.

Detalhe binário → `M-05-auto-vinculo/validation-contract.md`.
