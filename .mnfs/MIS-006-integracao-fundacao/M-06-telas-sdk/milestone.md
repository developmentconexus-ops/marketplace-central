# M-06 — telas + SDK + chain-viz

```yaml
id: M-06
type: milestone
mission: MIS-006
status: draft
depends_on: [M-01, M-02, M-05]
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Tornar visível — pela primeira vez — a cadeia real que a fundação (M-01/M-02/M-03/M-04/M-05)
constrói: N produtos importados → N vinculados → N coletas enfileiradas. Hoje essa cadeia não
existe como tela (`/importacoes` não é uma rota — `AppRouter.tsx:42-61`); o pedaço que existe
(`ImportacaoSection.tsx:115-146`) só lê estágio de parse, nunca junta vínculo nem enfileiramento,
e está duplicado entre `/vinculos` e `/integracoes`. M-06 fecha o ciclo: cria a rota, promove o
componente compartilhado, troca a config de fonte ativa de localStorage para banco-por-tenant, e
estende o SDK para carregar os dados reais que a tela precisa.

Ver `mission.md` §Milestone Strategy (linha M-06), `interface-contracts-mis006.md` §E8/§E9,
`architecture-map.md` (M-06 = onda 4, depende de M-01+M-02+M-05), `research/refactor-inventory-frontend.md`
§1-5 (evidência file:line de todo REFACTOR/CREATE abaixo).

## Scope

- CREATE rota `/importacoes` (`AppRouter.tsx`) + link de nav (`Header.tsx`).
- PROMOVE `ImportacaoSection.tsx` (hoje em `pages/vinculos/`, importado cross-page por
  `pages/integracoes/IntegracoesPage.tsx:7-8`) para módulo compartilhado — resolve a colisão de
  import relativo entre as duas páginas (`refactor-inventory-frontend.md` §2).
- REFACTOR `ImportacaoSection` p/ renderizar a cadeia `N-imported → N-linked → N-enqueued` (join
  `sync_state` + `product_links`), não mais só contagens de parse (`accepted/rejected/warning`).
- REFACTOR `ActiveSourceCard` (`IntegracoesPage.tsx:280-319`) de `localStorage` para config em
  banco por tenant, via `useActiveErpSource` (pacote `@marketplace-central/web-query`) apontando
  pro endpoint GET/PUT de M-02 (E9.0) — mata a leitura/escrita local (`refactor-inventory-frontend.md`
  §3 linha 28).
- REFACTOR `listErpImports` (SDK, `sdk-runtime/src/index.ts:1855`) para shape rica (contagens de
  vínculo + coleta por protocolo, não só parse); CREATE métodos novos de active-source (GET/PUT) e
  chain-read.
- KEEP `/integracoes` com os 4 cards existentes (`ActiveSourceCard`, `UploadCard`,
  `ProviderConnectCard`, `ImportacaoSection` promovido) — nenhum card removido, só o de active-source
  troca a fonte de dado.
- REFACTOR `VinculosPage`/`QueueRow` para badge "auto-aprovado" nos vínculos com
  `rule_matched=exact_ean_unique` (E10, produzidos por M-05).

## Non-Scope

- Backend-ização de `Oportunidades`/`buildOppRows` (`oportunidades.ts:49-101`) — REFACTOR
  identificado em `refactor-inventory-frontend.md` §1, mas é `/mercado`, não `/importacoes`;
  entra numa missão de mercado futura. Se aparecer no diff desta milestone, é escopo vazando.
- Scope "monitorado" (`MonitoradosTab.tsx`) — depende de backend inexistente, non-scope da missão
  inteira (`mission.md` §Non-Scope).
- Onboarding saga passos 2-6 (backfill anúncios/pedidos, progresso pós-conexão ML) —
  `IntegracoesPage.tsx:328-383` `ProviderConnectCard` fica como está (KEEP), sem novos passos.
- Implementação do endpoint active-source em si (M-02 já entrega GET/PUT) — M-06 só CONSOME.
- Implementação da geração/auto-aprovação de vínculo (M-05) — M-06 só EXIBE o resultado (badge).
- F3.7 discovery / tela de "produto sem anúncio" (M-07, condicional).

## Feature Briefs

### F-01 — Rota `/importacoes` + `ImportacaoSection` promovido + cadeia real

**EARS:**
- WHEN o operador navega para `/importacoes` THEN a tela renderiza, por protocolo de import, três
  contadores: `N importados` (linhas aceitas no protocolo), `N vinculados` (linhas com
  `product_links` resolvido, qualquer `rule_matched`), `N enfileirados` (linhas com row em
  `sync_state`/fila de coleta correspondente) — vindos de um join real, nunca hardcode.
- WHEN `ImportacaoSection` é usado em `/vinculos` E em `/integracoes` THEN ambos consomem o MESMO
  módulo compartilhado (import único, não duas cópias) — resolve a duplicação
  (`refactor-inventory-frontend.md` §2 linha 21).
- WHEN a lista de protocolos de `/importacoes` (`listErpImports`) é retornada THEN é ordenada
  `imported_at DESC` (protocolo mais recente primeiro) — ordem declarada e estável, nunca ordem de
  retorno de banco indefinida; a lista Fila/Resolvidos de `/vinculos` (F-04) ordena
  `created_at DESC`.
- WHEN um protocolo não tem nenhum vínculo ainda (import muito recente, M-05 não rodou) THEN
  `N vinculados`/`N enfileirados` mostram `0` real (contagem vazia, não erro) — distinto de dado
  ausente/desconhecido, que mostra "—" (ADR-17 aplica a CAMPO ausente, não a CONTAGEM zerada
  legítima).
- IF a chamada ao endpoint de chain-read falha (erro de rede/servidor) THEN a tela mostra "—" nas
  três contagens com indicação de erro, nunca `0` fabricado (ADR-17).

**Inputs/Outputs (MUST have):**
- SDK: `listErpImports()` (refactor, F-03) retorna, por protocolo, `{protocol_id, imported_count,
  linked_count, enqueued_count, imported_at, source}` — shape estendida vs hoje (só
  `accepted/rejected/warning_count/imported_at`, `ImportacaoSection.tsx:92-109`).
- Rota: `apps/web/src/app/AppRouter.tsx` ganha entry `path: "/importacoes"` render
  `ImportacaoPage` (nova, componente-casca fino que monta o `ImportacaoSection` promovido).
- Componente promovido: novo caminho compartilhado (ex.
  `apps/web/src/pages/shared/ImportacaoSection.tsx` ou `apps/web/src/components/importacao/` —
  decisão de pasta exata não bloqueante, DEVE deixar de viver sob `pages/vinculos/`).

**Negative Scenarios:**
- Import de protocolo malformado que zera todas as linhas (`accepted_count=0`) → cadeia mostra
  `0 importados → 0 vinculados → 0 enfileirados`, sem crash e sem "—" fabricado (contagem real
  zero é honesta, distinta de indisponibilidade).
- Dois componentes `ImportacaoSection` divergentes remanescendo pós-merge (um em `pages/vinculos`,
  outro promovido) = falha desta feature — grep deve mostrar UM arquivo fonte, dois import-sites.

**State Model:**
- Fonte de verdade das 3 contagens = servidor (join `sync_state`+`product_links`), nunca
  recomputado no cliente a partir de listas já carregadas para outro propósito (repete o erro
  golden-rule já flagueado em Oportunidades, `refactor-inventory-frontend.md` §1).
- Refetch: ao entrar em `/importacoes` (mount) e após qualquer ação de upload bem-sucedida
  disparada a partir da própria tela (reusa o padrão de invalidação de query já usado por
  `useErpImportUpload.ts`); não há polling automático nesta milestone (fora de escopo — sync_state
  cadence-agnostic é backend, M-01).
- Stale-state: se o operador troca a fonte ativa (F-02) enquanto `/importacoes` está aberta, a
  lista de protocolos NÃO se filtra automaticamente por fonte nesta milestone (protocolos são
  histórico imutável, `mission.md` ADR-01) — a troca de fonte afeta o que É importado dali em
  diante, não reescreve o histórico exibido.

**Write-set:** `apps/web/src/app/AppRouter.tsx` (nova rota), `apps/web/src/app/Header.tsx` (nav
link), novo módulo compartilhado do `ImportacaoSection` (path exato a definir na implementação,
fora de `pages/vinculos/`), `apps/web/src/pages/integracoes/IntegracoesPage.tsx` (import
atualizado pro novo caminho), `apps/web/src/pages/vinculos/VinculosPage.tsx` (idem).

---

### F-02 — `ActiveSourceCard` localStorage → DB-por-tenant

**EARS:**
- WHEN `ActiveSourceCard` monta THEN lê a fonte ativa do endpoint GET de M-02 (E9.0) via
  `useActiveErpSource`, NUNCA de `localStorage` (`IntegracoesPage.tsx:270-274` comentário atual
  removido).
- WHEN o operador troca a fonte no card THEN a troca chama o PUT do endpoint de M-02, e só reflete
  na UI após a resposta confirmar (sem update otimista silencioso que mascare falha de servidor).
- WHEN a troca é bem-sucedida THEN qualquer view dependente de fonte ativa já montada (ex.
  `IntegracoesPage`) refaz a leitura da fonte ativa a partir do hook (fonte única, não duas cópias
  de estado — local vs banco).
- IF o PUT retorna erro (rede, 400, etc.) THEN a UI mantém a fonte anterior visível e mostra erro,
  nunca assume a troca como aplicada.

**Inputs/Outputs (MUST have):**
- `useActiveErpSource` (`@marketplace-central/web-query`): hoje só lê/escreve `localStorage`
  (`refactor-inventory-frontend.md` §3 linha 28); REFACTOR para GET `/tenants/{tenant_id}/active-source`
  no mount e PUT no mesmo path na troca (path exato = o publicado por M-02 F-04).
- `ActiveSourceCard` (`IntegracoesPage.tsx:280-319`): mesma UI, troca só a fonte de dado
  (hook atualizado); nenhuma mudança visual fora de estado de loading/erro.

**Negative Scenarios:**
- Tenant sem `active_source` configurado ainda (M-02 `ErrUnknownActiveSource`) → card mostra
  estado "não configurado" explícito, nunca assume um default silencioso (ecoa AC-03/M02-C12).
- Troca de fonte em aba/sessão diferente do mesmo tenant → outra aba, ao refetch, vê o valor novo
  (prova de que o estado é server-side, não por-browser como hoje).

**State Model:**
- Ownership da fonte ativa passa de `localStorage` (client, por-browser) para o servidor
  (por-tenant); `useActiveErpSource` é o único ponto de leitura/escrita no FE — nenhum outro
  componente lê `localStorage` para isso pós-refactor.
- Refetch: no mount de qualquer consumidor do hook; após PUT bem-sucedido, o hook invalida seu
  próprio cache local (query) para refletir o valor confirmado pelo servidor, não o valor otimista.

**Write-set:** `packages/web-query/src/*` (arquivo de `useActiveErpSource`, path exato a
localizar na implementação — grep `useActiveErpSource` no pacote), `apps/web/src/pages/integracoes/IntegracoesPage.tsx:280-319`
(`ActiveSourceCard`, remoção do comentário/lógica de `localStorage`).

**Forbidden:** endpoint GET/PUT em si (`internal/composition`, migração `active_source`) —
pertence a M-02; F-02 só consome.

---

### F-03 — SDK: `listErpImports` rica + métodos active-source/chain

**EARS:**
- WHEN `listErpImports()` é chamado THEN o payload retornado inclui, por protocolo,
  `linked_count` e `enqueued_count` além dos campos já existentes
  (`accepted_count/rejected_count/warning_count/imported_at`) — shape estendida, não substituída
  (compat aditiva, ecoa MC-12).
- WHEN um método novo de active-source é adicionado ao SDK THEN ele espelha exatamente o schema
  publicado por M-02 F-04 (`{active_source, source_kind, set_at, set_by}`), sem campo inventado.
- WHEN o SDK e o OpenAPI da seção chain-read (M-06) são publicados THEN landam no MESMO commit
  (profile §7) — a seção é DISJUNTA da seção active-source de M-02 (path/schema sem overlap,
  `architecture-map.md` §Superfícies compartilhadas).

**Inputs/Outputs (MUST have):**
- `listErpImports()` (`sdk-runtime/src/index.ts:1855`): tipo de retorno `ErpImportSummary`
  estendido com `linked_count: number`, `enqueued_count: number` (NULL-safe: ausência de dado de
  vínculo/coleta para protocolos antigos pré-M-05 deve ser um valor explícito, ex. `0` real de
  contagem — não confundir com campo desconhecido; se o backend não tiver o dado ainda, retorna
  `null` e a UI honra ADR-17, não assume `0`).
- Novo método `getActiveSource(tenantId)` / `setActiveSource(tenantId, source)` — consumido por
  F-02.
- Novo método de chain-read (ex. `getImportChain(protocolId)` ou equivalente ao path real de
  M-06 no OpenAPI) — consumido por F-01.
- OpenAPI: seção nova para os endpoints acima que M-06 publica (chain-read); path exato definido
  na implementação, mas DEVE ser disjunto de `/tenants/{tenant_id}/active-source` (path de M-02).

**Negative Scenarios:**
- Commit que só toca `sdk-runtime` sem o OpenAPI correspondente (ou vice-versa) = falha de MC-12,
  igual à regra já aplicada em M-02 F-04.
- Seção OpenAPI de M-06 colide em path/schema com a seção active-source de M-02 = falha de
  ownership (`architecture-map.md` contract-lock).

**Write-set:** `packages/sdk-runtime/src/index.ts` (extensão de `listErpImports`, novos métodos),
spec OpenAPI do repo (seção chain-read, path exato a confirmar contra `contracts/` real).

**Parallel-safe with:** F-01/F-02/F-04 consomem os tipos/métodos que F-03 produz — serializar
F-03 antes ou em paralelo estreito (mesmo PR) com F-01/F-02, nunca depois isolado (as outras
features não compilam sem os tipos novos).

---

### F-04 — `VinculosPage` badge "auto-aprovado"

**EARS:**
- WHEN um vínculo em `listProductLinkWorkflows`/Resolvidos tem `rule_matched=exact_ean_unique` E
  `actor=system` (E10, produzido por M-05) THEN a linha exibe um badge visual distinto (ex.
  "Auto-aprovado") ao lado do status.
- WHEN um vínculo foi aprovado manualmente (`actor=operator`) THEN NENHUM badge de auto-aprovação
  aparece — distinção honesta de proveniência.
- IF o campo `rule_matched`/`actor` vier ausente (registro antigo pré-M-05, sem audit trail) THEN
  a linha não exibe nem "auto" nem "manual" — omite o badge (não assume manual por default).

**Inputs/Outputs (MUST have):**
- Payload de `listProductLinkWorkflows` (via SDK) precisa carregar `rule_matched`/`actor` (E10)
  por vínculo — se o método já retorna esses campos, F-04 só consome; se não, é dependência de SDK
  (F-03 ou M-05 SDK — a decidir na implementação, registrar como `REQUEST` se não existir).
- UI: badge em `QueueRow.tsx` (ou componente equivalente de Resolvidos), reusando o padrão visual
  de pill já existente em `QueueRow.tsx:155` (`bandLabels`/confidence-band).

**Negative Scenarios:**
- Vínculo com `collisions_at_decision > 1` (não deveria ter sido auto-aprovado, ADR-05) aparecendo
  com badge "auto-aprovado" = falha de dado, não de UI — reportar como REQUEST ao dono de M-05, não
  mascarar no FE.

**Write-set:** `apps/web/src/pages/vinculos/QueueRow.tsx`, `apps/web/src/pages/vinculos/VinculosPage.tsx`
(leitura do campo novo, sem mudar a lógica de filtro Fila/Resolvidos existente).

## Ownership & Concurrency (six-axis)

| Eixo | M-06 |
|------|------|
| Migração | nenhuma (M-06 é FE + SDK puro) |
| DB shape | nenhuma (consome `sync_state`/`product_links`/`active_source` de M-01/M-02/M-05, não cria tabela) |
| Módulo Go | nenhum de domínio novo; toca só a camada OpenAPI (spec) para publicar a seção chain-read |
| `root.go` | nenhum |
| Contrato/SDK | seção chain-read — **contract-lock**; disjunta da seção active-source de M-02 (`architecture-map.md` §Superfícies compartilhadas) |
| FE surface | **dono**: `apps/web` `AppRouter`/nav (`Header.tsx`), `ImportacaoSection` promovido, `pages/integracoes/*`, `pages/vinculos/QueueRow.tsx`/`VinculosPage.tsx` (leitura de badge), `packages/web-query` `useActiveErpSource` |

- Exclusive surfaces (only this milestone writes): `apps/web/src/app/AppRouter.tsx` (rota
  `/importacoes`), `apps/web/src/app/Header.tsx` (nav), módulo compartilhado novo do
  `ImportacaoSection`, `packages/web-query` (arquivo de `useActiveErpSource`), OpenAPI seção
  chain-read, `packages/sdk-runtime/src/index.ts` (métodos novos + `listErpImports` estendido).
- Migration block: none.
- Predicted seam locks: OpenAPI (additive, seção chain-read nova — não edita a seção active-source
  de M-02); `packages/sdk-runtime` (additive, novos métodos + campo estendido em tipo existente,
  sem remover campo).
- Runs in parallel with: nenhum milestone da mesma onda (M-06 é onda 4, sozinho); pode rodar
  paralelo a M-07 (onda condicional, superfícies disjuntas — M-07 não toca FE de importação).
- Internal feature DAG: F-03 (SDK) precede ou anda em lockstep com F-01/F-02 (consomem seus
  tipos/métodos); F-04 é independente das demais (`F-04 ∥ F-01 ∥ F-02`, todas dependem só de F-03
  já existir/andar junto).

## Risks

| Risco | Prob | Impacto | Mitigação | Trigger | Owner |
|-------|------|---------|-----------|---------|-------|
| `ImportacaoSection` promovido diverge visualmente entre `/vinculos` e `/integracoes` (props/contexto diferentes hoje) | Baixa | Médio | reusar mesmo componente sem props condicionais por página; se a diferença for real, documentar como duas variantes explícitas, não um componente com `if (page === ...)` | review do diff | F-01 |
| `useActiveErpSource` tem outros consumidores além de `IntegracoesPage` ainda não mapeados | Média | Médio | grep `useActiveErpSource` em todo `apps/web`/`packages` antes de refatorar; qualquer consumidor extra entra no write-set | grep na implementação | F-02 |
| Contagens de vínculo/coleta exigem endpoint backend que M-06 não implementa (só consome) | Média | Alto (F-01 bloqueia sem endpoint) | dependência explícita em M-02 (config) + M-01 (`sync_state`) + M-05 (`product_links`) já modeladas em `depends_on`; se o endpoint de chain-read específico não existir em nenhuma milestone anterior, é gap de decomposição — REQUEST ao hub antes de codar | ausência do endpoint na implementação | hub |

## Done Means

- `/importacoes` navegável, mostra cadeia real (3 contagens de um join, não hardcode/client-only).
- `ImportacaoSection` existe em UM lugar só, usado por `/vinculos` e `/integracoes` sem duplicação.
- `ActiveSourceCard` lê/escreve banco por tenant; zero leitura de `localStorage` para fonte ativa.
- `listErpImports` + métodos novos de SDK publicados junto do OpenAPI no mesmo commit.
- Badge de auto-aprovado visível em `/vinculos` para vínculos `rule_matched=exact_ean_unique`.
- `/integracoes` mantém os 4 cards, nenhuma regressão visual (design parity light+dark).

## Handoff

- Current status: planned (aguardando P7 readiness da missão).
- Next owner: hub (execução por milestone) após P7 Ready.
- Next action: dispatch de chip(s) para F-01..F-04 (onda 4, após M-01+M-02+M-05 mergeados).
- Required files/evidence: este `milestone.md`; `validation-contract.md` (M-06); P7 browser QA
  screenshots light+dark (`/importacoes`, `/integracoes`, `/vinculos`).
- Blockers or open decisions: endpoint exato de chain-read (path/schema) fica a definir na
  implementação (F-03), desde que disjunto de M-02; nenhum blocker de escopo.

## Correction Handoff

- QA failure summary: n/a (planning).
- Correction scope: n/a.
- Attempts used/remaining: n/a.
- Next artifact: n/a.
- Revalidation evidence required: n/a.

## Decisão do operador D-120 (2026-07-22) — toggle de fonte ativa

- Regressão registrada: desde M-02 (@49ab3bdd) a fonte ativa é resolvida pelo BANCO
  (routing.Reader sobrescreve o ctx vindo de query-param) — o rádio "Fonte ativa" de
  /integracoes (localStorage) está MORTO. Operador decidiu: NENHUM bridge FE antes de M-06;
  o wiring rádio→PUT /config/active-source fecha AQUI (ver M06-U1 no validation-contract).
  Até lá, troca de fonte só via API — limitação conhecida e aceita.
- Drift de contrato a corrigir no dispatch: F-02 cita GET/PUT `/tenants/{tenant_id}/active-source`,
  mas o endpoint LANDADO por M-02 é `GET/PUT /config/active-source` (single-tenant, tenant fixo
  server-side; OpenAPI :3207/:3230 + sdk-runtime activeSource.ts). O chip consome o landado.

## Correção do hub D-122 (2026-07-28) — briefs desatualizados, verificados no repo

Este milestone é de D-120. O repo andou. Três afirmações acima **não são mais verdade em
`5441fe18`**, e ficam corrigidas aqui em vez de anotadas (R-25) — quem despacha consome esta seção,
não o brief original.

- **F-04, predicado impossível.** O EARS pede `rule_matched=exact_ean_unique` **E** `actor=system`.
  A CHECK de `migrations/0082_product_link_decisions.sql:54` é
  `CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean')` — ator `system` só admite
  `concordant_codprod_ean`. O brief é da política D-120; **D-121 estreitou** (só CODPROD+EAN
  concordantes auto-aprovam) e o brief não acompanhou.
- **F-04, campo fora do wire.** `rule_matched` não existe em `contracts/marketplace-central.openapi.yaml`
  nem em `packages/sdk-runtime/src/`. Vive só no DB e num read per-link
  (`link_candidate_repo.go:411 ListDecisionsForLink`) que rota nenhuma expõe. **O caminho que existe**
  custa zero backend: a auto-aprovação grava a auditoria com
  `ActorType: "system", ActorID: "auto_linker"` (`resolution_service.go:280`) e
  `item.audit[].actor.actor_type` já está no wire. Predicado do badge = a entrada de auditoria que
  resolveu tem `actor_type === "system"`.
- **F-02 já está descarregado.** `packages/web-query/src/activeSource.ts` e
  `pages/integracoes/IntegracoesPage.tsx:297-346` já implementam active-source do DB sobre
  `GET/PUT /config/active-source`. A decisão de operador D-120 acima descreve uma regressão **já
  corrigida**; o chip da onda 2 **verifica e declara**, não reconstrói.
- **Blocker de F-03 resolvido.** "endpoint exato de chain-read fica a definir na implementação" está
  fechado: `GET /erp/imports/{id}/chain` landou no CHIP-ANCHORS-2 com OpenAPI e SDK
  (`getErpImportChain`, `sdk-runtime/src/index.ts:1901`). Falta o **consumo no FE**, que é o F-01.
- **Ponteiro de linha podre:** F-04 cita `QueueRow.tsx:155` como o padrão de pill; hoje `:155` é o
  comentário do ranking ADR-17 e os pills estão em `:20-38`. Verifique por string, nunca por linha.

**D-E fatia este milestone por TELA, não por camada** (ver `DECISOES-D122-anchors-telas.md`):
`/vinculos` (F-04 + F-05) = CHIP-VINC-NEUTRO; `/importacoes` + `/integracoes` (F-01 + F-03) =
CHIP-IMPORT-CHAIN. Os packs são a autoridade de dispatch; este arquivo é o contexto deles.
