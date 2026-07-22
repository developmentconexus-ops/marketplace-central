# Milestone Validation Contract — M-06-telas-sdk

```yaml
id: M-06-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-06
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização.

## Milestone ID

M-06

## QA Level

QA-0

## Required Outcome

`/importacoes` existe como rota própria e renderiza a cadeia real (N-imported→N-linked→N-enqueued)
a partir de um join server-side, light+dark; `ActiveSourceCard` lê/escreve config por tenant no
banco (não mais `localStorage`); SDK (`listErpImports` + métodos novos) e OpenAPI da seção
chain-read landam no mesmo commit, disjunta da seção active-source de M-02; `/integracoes` mantém
os 4 cards sem regressão; `VinculosPage` mostra badge de auto-aprovado para vínculos
`exact_ean_unique`.

## Critérios

## Criterion: Rota /importacoes existe e renderiza
ID: M06-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: navegar para `/importacoes` (browser QA)
- Expected: rota resolve (sem 404/redirect), tela renderiza `ImportacaoSection` promovido
- Actual:
- Artifact:
Blocking failure: rota ausente, ou `/importacoes` cai em fallback/redirect para `/integracoes`
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant seed com ≥1 protocolo de import já concluído (reusar seed de `#004-E`/`#005-E`
  citado em `mission.md`/contrato, ou equivalente do dev stack hub-owned)
- Steps:
  - open /importacoes
  - assert url ~ "/importacoes"
  - assert text "N importados"
- Expected: tela carrega sem erro, mostra ao menos 1 protocolo
Owner: QA Validator

## Criterion: Cadeia real (join, não hardcode)
ID: M06-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: inspecionar payload de rede da chamada que popula `/importacoes` (DevTools/network capture)
- Expected: response contém `linked_count`/`enqueued_count` numéricos por protocolo, valores
  distintos de `accepted_count` (prova de join real, não eco do mesmo número três vezes)
- Actual:
- Artifact:
Blocking failure: `linked_count`/`enqueued_count` ausentes do payload, ou sempre iguais a
`accepted_count` (sinal de hardcode/eco client-side)
Blocking failure observed: No
Owner: QA Validator

## Criterion: ImportacaoSection compartilhado (sem duplicação)
ID: M06-C3
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `grep -rn "function ImportacaoSection" apps/web/src`
- Expected: exatamente 1 definição do componente; `pages/vinculos/*` e `pages/integracoes/*`
  importam do mesmo caminho promovido
- Actual:
- Artifact:
Blocking failure: 2+ definições de `ImportacaoSection` (duplicação não resolvida)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Cadeia visível em light+dark (design parity)
ID: M06-C4
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: browser QA `/importacoes` com tema claro e escuro (DESIGN-REFERENCE)
- Expected: 3 contadores legíveis (contraste ok) nos dois temas, sem elemento cortado/sobreposto
- Actual:
- Artifact:
Blocking failure: contraste insuficiente ou layout quebrado em qualquer um dos dois temas
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: mesma de M06-C1
- Steps:
  - open /importacoes
  - assert text "N vinculados"
  - assert text "N enfileirados"
- Expected: contadores visíveis e legíveis nos dois temas (screenshot par light/dark)
Owner: QA Validator

## Criterion: Chain honesta em erro de rede
ID: M06-C5
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: simular falha da chamada de chain-read (network throttling/mock 500 no dev stack)
- Expected: as 3 contagens mostram "—" (honesto), nunca `0` fabricado (ADR-17)
- Actual:
- Artifact:
Blocking failure: qualquer contagem cai para `0` em vez de "—" quando o fetch falha (viola AC-03)
Blocking failure observed: No
Owner: QA Validator

## Criterion: ActiveSourceCard lê/escreve banco, não localStorage
ID: M06-C6
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `grep -rn "localStorage" packages/web-query/src apps/web/src/pages/integracoes/IntegracoesPage.tsx`
- Expected: 0 hits relacionados a `active_source`/fonte ativa
- Actual:
- Artifact:
Blocking failure: qualquer leitura/escrita de `localStorage` remanescente para fonte ativa
Blocking failure observed: No
Owner: QA Validator

## Criterion: Toggle de fonte ativa refetch cross-tab (server-side de fato)
ID: M06-C7
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA — trocar fonte em `/integracoes`, abrir segunda aba do mesmo tenant, recarregar
- Expected: segunda aba mostra o valor novo após reload (prova que o estado é server-side por
  tenant, não local-per-browser)
- Actual:
- Artifact:
Blocking failure: segunda aba mantém valor antigo após reload (indício de estado ainda local)
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant seed com `active_source` já configurado (M-02 seed); duas abas mesma sessão
- Steps:
  - open /integracoes
  - click "Alternar fonte ativa"
  - assert text "<nova fonte>"
- Expected: troca confirmada na UI após resposta do PUT (sem update otimista sem confirmação)
Owner: QA Validator

## Criterion: ActiveSourceCard trata tenant sem config
ID: M06-C8
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA `/integracoes` com tenant sem row em `active_source` (M-02 `ErrUnknownActiveSource`)
- Expected: card mostra estado "não configurado" explícito, nunca assume default silencioso
- Actual:
- Artifact:
Blocking failure: card renderiza uma fonte qualquer como se estivesse configurada (fallback silencioso)
Blocking failure observed: No
Owner: QA Validator

## Criterion: listErpImports shape estendida
ID: M06-C9
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: leitura de tipo `ErpImportSummary` em `packages/sdk-runtime/src/index.ts`
- Expected: campos `linked_count`/`enqueued_count` presentes, campos existentes
  (`accepted_count/rejected_count/warning_count/imported_at`) preservados sem mudança de tipo
- Actual:
- Artifact:
Blocking failure: campo existente removido/tipo alterado (quebra de compat aditiva)
Blocking failure observed: No
Owner: QA Validator

## Criterion: SDK+OpenAPI mesmo commit (chain-read)
ID: M06-C10
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git show --stat <commit>` do PR/chip que publica a seção chain-read
- Expected: diff do mesmo commit inclui spec OpenAPI E arquivo(s) `sdk-runtime` gerado
- Actual:
- Artifact:
Blocking failure: SDK e OpenAPI landam em commits separados
Blocking failure observed: No
Owner: QA Validator

## Criterion: Seção chain-read disjunta de active-source (M-02)
ID: M06-C11
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: diff do OpenAPI spec — comparar path/schema da seção chain-read (M-06) contra a seção
  active-source (M-02, `/tenants/{tenant_id}/active-source`)
- Expected: nenhum overlap de path ou schema entre as duas seções
- Actual:
- Artifact:
Blocking failure: seção chain-read reusa/edita path ou schema já publicado por M-02
Blocking failure observed: No
Owner: QA Validator

## Criterion: /integracoes preserva os 4 cards (sem regressão)
ID: M06-C12
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: browser QA `/integracoes` light+dark
- Expected: `ActiveSourceCard`, `UploadCard`, `ProviderConnectCard`, `ImportacaoSection`
  (promovido) todos presentes e funcionais, layout idêntico ao pré-M-06 exceto fonte de dado do
  ActiveSourceCard
- Actual:
- Artifact:
Blocking failure: qualquer um dos 4 cards ausente ou quebrado
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: mesma de M06-C1
- Steps:
  - open /integracoes
  - assert text "Conectar"
  - assert text "Fonte ativa"
- Expected: 4 cards visíveis, tema claro e escuro sem quebra visual
Owner: QA Validator

## Criterion: Badge auto-aprovado em /vinculos
ID: M06-C13
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA `/vinculos` tab Resolvidos, com ≥1 vínculo `rule_matched=exact_ean_unique`
  (seed produzido por M-05, ex. EAN único como no IO A da missão)
- Expected: linha do vínculo mostra badge "Auto-aprovado"; vínculo aprovado manualmente na mesma
  lista NÃO mostra o badge
- Actual:
- Artifact:
Blocking failure: badge ausente onde deveria aparecer, ou aparece em vínculo `actor=operator`
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant seed com ≥1 vínculo auto-aprovado (M-05) e ≥1 vínculo manual, mesma tela
- Steps:
  - open /vinculos
  - click "Resolvidos"
  - assert text "Auto-aprovado"
- Expected: badge só na linha correta, ausente nas demais
Owner: QA Validator

## Criterion: Badge omitido quando audit ausente
ID: M06-C14
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: inspecionar vínculo antigo pré-M-05 (sem `rule_matched`/`actor`) em `/vinculos`
- Expected: linha não mostra badge "auto" nem "manual" — omissão honesta, não default assumido
- Actual:
- Artifact:
Blocking failure: vínculo sem audit trail exibe qualquer badge de proveniência (assume manual por default)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Nenhuma escrita fora do write-set desta milestone
ID: M06-C15
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git diff --stat` do chip/PR de M-06 contra `Ownership & Concurrency` do milestone.md
- Expected: todos os arquivos tocados pertencem às exclusive surfaces listadas (FE `apps/web`,
  `packages/web-query`, `packages/sdk-runtime`, seção OpenAPI chain-read); zero migração, zero
  módulo Go de domínio novo
- Actual:
- Artifact:
Blocking failure: diff toca migração, `composition/root.go`, ou seção active-source do OpenAPI (M-02)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Toda evidência de UI (M06-C1, M06-C4, M06-C7, M06-C8, M06-C12, M06-C13) exige screenshot ou
  transcript de browser QA light+dark salvo em `docs/design/evidence/`.
- Toda evidência de payload/rede (M06-C2, M06-C5, M06-C9) exige captura do response JSON, não
  descrição textual.
- Toda evidência de commit (M06-C10) exige `git show --stat` salvo, não afirmação.

## Blocking Failures

- Qualquer contagem fabricada (`0` no lugar de "—" em erro, ou eco do mesmo número nas 3 posições)
  = blocking (viola ADR-17/M06-C2/M06-C5).
- `localStorage` remanescente para fonte ativa = blocking (M06-C6).
- SDK/OpenAPI em commits separados, ou overlap de seção com M-02 = blocking (M06-C10/M06-C11).
- Regressão em `/integracoes` (card ausente/quebrado) = blocking (M06-C12).
- Badge de auto-aprovado incorreto (presente onde não deveria, ou default assumido em audit
  ausente) = blocking (M06-C13/M06-C14).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (dispatch de chip pós P7 Ready, onda 4).
- Next action: aguardar merge de M-01+M-02+M-05 antes de dispatch.
- Required files/evidence: este arquivo; `M-06-telas-sdk/milestone.md`; evidência viva em
  `docs/design/evidence/` (screenshots light+dark, payload captures, `git show --stat`).
- Blockers or open decisions: nenhum — dependências (M-01, M-02, M-05) e boundary (chain-viz aqui,
  Oportunidades backend-ization em missão futura) já explicitados.

## Critérios de user-drive (AMENDMENT D-120 — obrigatório, ratificado pelo operador)

Mesma regra ratificada em M-03 (origem: regressão /catalogo 503 invisível aos gates de código,
hub-fix @2567eb44). M-06 é a milestone MAIS user-facing — estes critérios são o coração do gate.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M06-U1 | Toggle "Fonte ativa" em /integracoes grava via PUT /config/active-source NO BANCO e /catalogo reflete o flip imediatamente — fecha a regressão do toggle morto (D-120; decisão do operador 2026-07-22: bridge deferido para cá, nenhum wiring FE antes desta milestone) | browser drive flip completo ida-e-volta + row active_source no DB |
| M06-U2 | /importacoes chain-viz mostra N-importado, N-vinculado, N-enfileirado com números REAIS batendo com o DB | browser drive + 3 SELECTs de conferência |
| M06-U3 | Perfil limpo (localStorage vazio): todas as telas funcionam só com estado do banco — nenhum comportamento depende de localStorage para fonte ativa | browser drive em perfil/aba limpa |
| M06-U4 | Zero regressão nas telas existentes: /catalogo, /integracoes, /anuncios, /precos, /pedidos, /mercado carregam e operam sem erro novo | browser drive (6 telas) light+dark |
