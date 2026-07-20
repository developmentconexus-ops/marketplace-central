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
só alcançável por endpoint órfão) + REUSE do sinal de unicidade de EAN já calculado
(`validEANCounts`/`identityQuality`) + REFACTOR de `resolution_service.go` para uma transição
auto-approve que reusa a máquina de audit existente, formalizada como E10 (trail nova). Handlers
HTTP órfãos são KEPT (não removidos) — ganham um caminho de invocação interna adicional.

Ver `mission.md` §Milestone Strategy (linha M-05), ADR-05, `interface-contracts-mis006.md` §E4
(E4.1)/§E10, `architecture-map.md` (M-05 depende de M-02+M-03; alimenta M-06),
`research/refactor-inventory-backend.md` §6 (file:line da REFACTOR/REUSE).

## Scope

- REFACTOR `product_links/application/generation_service.go:60-78` (`GenerateLinkCandidates`):
  ADD chamada interna (não-HTTP) disparada ao final de import xlsx (hook de M-03) e de sync
  Sankhya (M-04, quando existir) — lógica de geração em si NÃO muda, só ganha caller automático.
- REUSE sinal de colisão de EAN já implementado e testado: `validEANCounts`/`identityQuality`
  (`erp_import/adapters/internalread/reader.go:344-366`) do lado xlsx; equivalente já existe
  simétrico em `internal_read/adapters/oracle/reader.go:70-76` do lado Sankhya. NÃO reimplementar
  a contagem de colisão — só consumir o `QualityFlags` resultante no momento da geração.
- REFACTOR `product_links/application/resolution_service.go:129-149` (`ApproveCandidate`): ADD
  transição interna de auto-approve (`collisions[ean]==1` + match exato por EAN) que reusa a
  MESMA máquina de transição/audit hoje só acionada por operador manual via transport.
- CREATE E10 audit trail (tabela nova, migração bloco B+ — após bloco B de M-02): toda decisão de
  vínculo (manual ou automática) grava linha `rule_matched`, `actor`, `collisions_at_decision`,
  `superseded_by`.
- CREATE idempotência A8: chave única `(internal_product_id, provider_listing_id)` em
  `product_links` — re-run do trigger não duplica vínculo; override manual do operador vence
  auto-aprovação prévia e nunca é sobrescrito de volta pelo automático.
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
- Mudança de política de match (fuzzy, SKU, título) — só o caminho EAN-exato-único ganha
  auto-approve; qualquer outro critério permanece manual (fora de escopo desta milestone).

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

### F-02 — Auto-approve EAN-exato-único (REUSE sinal de colisão)

**EARS:**
- WHEN a geração de candidatos produz um candidato com match exato por EAN E
  `collisions[ean] == 1` (sinal de `validEANCounts`/`identityQuality`,
  `reader.go:344-366`, já calculado — NÃO recalculado aqui) THEN o candidato transiciona
  automaticamente para aprovado, reusando a máquina de transição de
  `resolution_service.go:129-149` (`ApproveCandidate`), com `actor=system`.
- WHEN `collisions[ean] > 1` (EAN duplicado/ambíguo entre produtos) THEN o candidato permanece
  em REVIEW — auto-approve NUNCA dispara em ambiguidade (segurança > cobertura, ADR-05).
- WHEN o listing não tem EAN (`ean` ausente/vazio) THEN o candidato fica em REVIEW com motivo
  visível (`"sem EAN"` ou equivalente honesto) — comportamento hoje inalterado, só formalizado
  como caso negativo explícito desta milestone.
- IF o match não é exato por EAN (ex. candidato por título/SKU futuro) THEN nunca auto-aprova,
  independentemente de colisão — auto-approve é EXCLUSIVO do caminho EAN-exato-único.

**Inputs/Outputs (MUST have):**
- Consumo do `QualityFlags`/contagem de colisão já produzido por
  `erp_import/adapters/internalread/reader.go:344-366` (lado xlsx) e
  `internal_read/adapters/oracle/reader.go:70-76` (lado Sankhya) — nenhuma nova lógica de
  contagem de EAN é escrita nesta milestone.
- Extensão de `resolution_service.go` com uma transição de auto-approve que reusa o MESMO
  código de mudança de estado + escrita de audit hoje usado pelo caminho manual
  (`ApproveCandidate`), diferindo só no `actor` gravado.

**Negative Scenarios:**
- EAN com `collisions[ean] == 2` → candidato fica REVIEW, mesmo que um dos dois produtos "pareça"
  o certo — sem heurística de desempate automática.
- Reimplementação de contagem de colisão dentro de `product_links` (em vez de consumir o sinal
  existente) = falha de design (viola ADR-05 "não reimplementar sinal já testado").
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
  `link_id, rule_matched (exact_ean_unique|manual|...), actor (system|operator),
  collisions_at_decision, created_at, superseded_by NULL`.
- WHEN a aprovação é automática (F-02) THEN `rule_matched=exact_ean_unique`, `actor=system`,
  `collisions_at_decision=1`.
- WHEN o operador aprova/sobrescreve manualmente um vínculo (inclusive um já auto-aprovado) THEN
  uma NOVA linha E10 é gravada com `actor=operator`, e a linha anterior recebe `superseded_by`
  apontando para a nova — o override do operador VENCE e nunca é revertido de volta pelo caminho
  automático.
- WHEN o mesmo par `(internal_product_id, provider_listing_id)` é processado de novo (re-run de
  import/sync, F-01 disparando geração outra vez) THEN a constraint única de `product_links`
  impede duplicação — a geração/approve é idempotente, não cria segundo vínculo nem segunda
  linha de audit redundante para o mesmo estado.
- IF `product_links` já tem vínculo ativo para o par E o novo candidato geraria o MESMO resultado
  (mesmo `rule_matched`) THEN nenhuma escrita nova ocorre (no-op idempotente, não erro).

**Inputs/Outputs (MUST have):**
- Tabela E10 nova (migração bloco B+, após bloco B de M-02): shape exata do
  `interface-contracts-mis006.md` §E10 (`link_id, rule_matched, actor, collisions_at_decision,
  created_at, superseded_by NULL`).
- Constraint única em `product_links`: `(internal_product_id, provider_listing_id)` — chave A8.
- Escrita de audit acoplada à MESMA transação da transição de estado do vínculo (aprovação e
  audit não podem divergir — sem vínculo aprovado sem linha E10 correspondente).

**Negative Scenarios:**
- Override manual sobre vínculo auto-aprovado → 2 linhas E10 no total para o par
  (`system` superseded, `operator` vigente) — nunca 1 linha só sobrescrita in-place (histórico
  imutável, mesmo padrão de `erp_import_protocols`).
- Re-run do trigger duas vezes seguidas sobre o mesmo snapshot → `SELECT count(*) FROM
  product_links WHERE (internal_product_id, provider_listing_id) = (...)` retorna 1, não 2.
- Tentativa de auto-approve escrever por cima de vínculo com `actor=operator` já vigente → é
  bloqueada pela regra de precedência (F-02 nunca chama a transição se já existe vínculo ativo
  com `actor=operator` no par).
- Qualquer escrita de E10 sem `tenant_id` explícito → falha de design (AC-01, todas as tabelas
  novas desta milestone são tenant-scoped via o vínculo pai).

**Write-set:** `apps/server_core/migrations/09xx_product_links_audit.{up,down}.sql` (bloco B+,
tabela E10 + constraint única A8 em `product_links`), `internal/modules/product_links/domain`
(tipo de audit row), `internal/modules/product_links/application/resolution_service.go` (escrita
de audit acoplada à transição, manual E automática).

## Ownership & Concurrency (six-axis)

| Eixo | M-05 |
|------|------|
| Migração | bloco B+ (E10 audit trail + constraint A8 em `product_links`) — após bloco B de M-02 (`products_mirror`/`active_source`) |
| DB shape | `product_links` (constraint nova A8), tabela de audit E10 — dono |
| Módulo Go | `internal/modules/product_links/*` (application: generation_service.go, resolution_service.go; transport: http_handler.go KEPT sem mudança de assinatura) — dono |
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

- **MC-05**: cadeia de vínculo automática — import dispara geração; EAN-exato-único auto-aprova
  com audit; EAN duplicado → REVIEW.
- **MC-06**: idempotência — re-run não duplica vínculo; override manual vence.

Detalhe binário → `M-05-auto-vinculo/validation-contract.md`.
