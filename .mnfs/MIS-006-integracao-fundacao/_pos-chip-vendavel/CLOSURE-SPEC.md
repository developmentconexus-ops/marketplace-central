# CHIP-VENDAVEL — spec de fechamento

Data: 2026-07-31. Aprovada pelo operador nesta sessão.
Localização: este arquivo vive no pacote `.mnfs`, não em `docs/superpowers/specs/`, porque
`AGENTS.md` vincula evidência ao `.mnfs` e instrução do repositório vence default de skill.

Estado na abertura desta spec: `f2dfb7a9` (worktree-chip-vendavel), árvore limpa, 35 commits
desde a base de despacho `554788d5`. S0–S12, S10-COND e S14 fechados. Tip da main: `664ac48e`.

## 0. Divisão de trabalho — decidida

O chip fecha após P5 (escada) e P6 (dual gate), e emite `CLOSED`. O hub faz merge, sync live
contra Oracle, e QA de browser com persona fresca. Chip não sobe servidor, não toca `:8080`,
não lê `.env*`, não dá push, não commita na main, não dirige o browser.

Razão: só QA passa, e quem construiu não valida a própria obra. VC-1..VC-7 exigem persona
fresca; um live-drive meu seria autocertificação.

## 1. Consertos antes da escada

Os dois consertos de código entram ANTES do P5. Se entrassem depois, a escada rodaria duas
vezes: ambos tocam Go e invalidam a lane de integração e a unitária.

### F1 — S13-CLOCK (grant G-1)

Sítio: `apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go:74-85`.

Defeito: `before := time.Now()` / `after := time.Now()` lêem o relógio do **host Windows**,
enquanto `queue_read_at` vem de `statement_timestamp()` (`query_repository.go:162`), que é o
relógio do **container Postgres**. Dois domínios de relógio comparados com tolerância zero.
Sob Docker Desktop os dois derivam; a flake disparou 3× e passou em re-run (ledger:535).

Conserto: ler os dois lados do sanduíche do PRÓPRIO banco, na mesma pool, em vez do host.
Mesmo domínio de relógio → tolerância zero continua correta e a asserção continua estrita.

Rejeitado: tolerância de ±2s. O número seria arbitrário e passaria a esconder regressão real
— um campo congelado ou hardcoded cabe dentro de qualquer janela generosa.

Must-fail obrigatório: congelar `queue_read_at` numa constante na query e mostrar o teste
VERMELHO nomeando o intervalo. Asserção de presença não pega valor errado; o vermelho é o que
prova que o intervalo ainda discrimina.

### F2 — `(unproved)` vazando para a tela (grant novo, registrado no ledger)

Sítios produtores: `apps/server_core/internal/modules/product_links/application/generation_service.go:524`
e `:577`, que gravam o marcador dentro do `Detail` do motivo do candidato — texto que a tela
`/vinculos` renderiza no chip de motivo.

Por que é escopo deste chip e não resíduo da MIS-006: o próprio
`validation-contract.md` fecha com "nenhum marcador de dev tipo `(unproved)` em copy", e VC-7
varre `/vinculos`. É violação do nosso contrato, não de outro.

Write set do conserto (alargamento além do write set original — daí o grant):
- `internal/modules/product_links/application/generation_service.go` (2 strings)
- `internal/modules/product_links/application/generation_service_test.go` (fixture)
- `apps/web/src/pages/vinculos/QueueTab.test.tsx` (4 sítios de fixture)
- `apps/web/src/pages/vinculos/VinculosDesign.golden.test.tsx` (3 sítios, incluindo a asserção
  de `title` em `:168`)

Texto novo: `ean corrobora o mesmo codprod, unicidade não comprovada` (e a variante curta
`ean corrobora codprod, unicidade não comprovada` no sítio `:577`). Copy provisória, herda A-23.

**Não tocar** `ReadErrorUniquenessUnproved = "uniqueness_unproved"`
(`internal/modules/internal_read/domain/sankhya_linkage.go:12`). É código de erro no fio,
consumido por comparação, não texto de tela. Trocá-lo quebraria contrato para consertar copy.

Prova: os goldens FE assertam o `title` do chip, então as fixtures são a asserção — mudá-las
sem mudar o produtor deixa o teste vermelho. Somar um grep provando zero ocorrência de
`unproved` em posição de copy, com o código de erro explicitamente excluído e nomeado.

**Ressalva de dado, para o hub:** motivos já gravados no banco carregam o texto velho. Conserto
de código não reescreve linha velha. O QA verá as duas formas até candidatos serem regenerados.
Regeneração mexe em dado vivo de vínculo — decisão do hub, não deste chip.

### F3 — número morto no contrato

`validation-contract.md:9` ainda carrega `N≈3.822`. Esse número foi medido sem o corte de
`CODLOCAL` e sem reserva. `ADDENDUM-01-codlocal.md:53` e A-3 (`BATCH-PLAN.md:1133`) já
ratificaram **2.923**, e `BATCH-PLAN.md:1188` restringe: 2.923 é o número do **Oracle live
apenas**; o caminho do espelho tem número próprio.

Ação: deletar a prosa morta e apontar para o adendo como fonte do número. Não é reescrever
critério de aceite — a emenda já existe e foi ratificada; o arquivo é que ficou mentindo, e
prosa falsa em contrato se deleta.

## 2. P5 — escada de verificação

### 2.1 Pré-condição: container de sessão

O `mpc-pg-session-*` desta sessão não existe mais. Recriar sob grant A-25 R-1: endpoint NOVO,
nunca re-apontar o dev stack, porta 5435 proibida. `CREATE DATABASE` com retry (`pg_isready`
mente no primeiro boot), depois `go run ./cmd/testdb migrate` esperando `applied 72` e, na
segunda passada, `applied 0`. Nome do banco tem que casar `^mpc_test_[0-9a-f]{32}$`.

Nunca imprimir `$env:MPC_TEST_DATABASE_URL` — carrega a senha do container. Provar que está
setado por booleano ou comprimento.

### 2.2 Lanes

Comandos do card S13-VERIFY estão obsoletos em dois pontos e são substituídos aqui.

| Lane | Comando | Verde é |
|---|---|---|
| Go build | `cd apps/server_core`, caches locais, `go build ./...` | exit 0 |
| Go test | `go test ./... -count=1` | exit 0 |
| **Pacotes sem lane** | `go test ./internal/modules/tenant_config/... ./internal/composition/... -v -count=1` | RUN=PASS, **SKIP=0** |
| tsc | `apps/web`, `tsc --noEmit -p tsconfig.json` | ≤12 erros, todos em arquivos alheios; **zero** nos arquivos tocados |
| vitest web | `npm run test --workspace @marketplace-central/web` (run cheio) | contagem por linha, SKIP=0 |
| vitest sdk | `packages/sdk-runtime` (run cheio) | idem |
| integração | `npm run harness:integration` | relatório prova que os testes nomeados rodaram |
| governança | `npm run harness:governance -- -BaseSha 664ac48e…` (tip da main, 40-hex, RE-LIDO no momento da execução) | limpo |

Correção 1: os pacotes `tenant_config` e `internal/composition` não são descobertos por lane
nenhuma. A descoberta varre as **primeiras cinco linhas** de `internal/modules/**/*_test.go`
procurando `//go:build integration` (`scripts/harness/Postgres.psm1:42-61`), e esses dois
pacotes caem fora — `internal/composition` nem está sob `internal/modules/`. Rodar explícito,
com o banco vivo, contando SKIP. Sem isso o verde da escada é fantasma sobre a fiação de
produção.

Correção 2: `-BaseSha` do card aponta para a base de DESPACHO (`554788d5`). A main já andou
(`664ac48e`). Diff contra base de despacho esconde revert — a base é o tip do ALVO. Verificar
`merge-base` antes de rodar.

Correção 3: `npx --no-install vitest run` do card foi substituído pela forma de workspace. A
forma do card já produziu verde vacuoso neste repositório.

Nenhuma lane pode imprimir `No test files found`. Isso é filtro batendo em nada, não verde.

Antes das lanes: conferir que não há `.js`/`.js.map` emitidos não-rastreados sombreando fonte
(A-25 R-2a). Verificado limpo em `f2dfb7a9`; re-verificar depois dos consertos.

## 3. P6 — dual gate

Sobre árvore CONGELADA, depois da escada verde. Dois assentos, em paralelo entre si, nunca em
paralelo com a escada: veredito sobre SHA que a escada ainda pode invalidar não vale nada.

- Assento 1 — Opus frio, fisicamente read-only (`harness:gate-reviewer`).
- Assento 2 — GPT-5.6 Sol medium, OS-process: prompt em arquivo, stdin fechado
  (`@() | codex exec`), `-c model_reasoning_effort=medium`, log teed e `-o <...>.last.md`.

Os dois recebem: o contrato de validação mais o adendo, o range do diff contra o tip da main
(não o pacote — gate lê o diff), e a ordem de STOP-THE-LINE por classe: defeito reincidente
vira conserto geral ou dívida registrada, não remendo pontual.

Concordância obrigatória. Discordância = outra rodada com o desacordo nomeado, não desempate
por moeda nem por antiguidade.

## 4. Entrega

`EVIDENCE.md` no diretório do pacote, escrito durante o `go test` (só os números entram no
fim), antes de qualquer evento `CLOSED`.

Evento `CLOSED` ao hub `local_99feb041-a5b3-4161-b6dc-bd38e65b6156`, carregando:

1. Range dos commits e caminho do `EVIDENCE.md`.
2. `REQUEST`: merge; sync live contra Oracle; QA de browser com persona fresca cobrindo
   VC-1..VC-7 nas três telas (`/integracoes`, `/catalogo`, `/vinculos`).
3. Ressalva do F2: linhas antigas de motivo ainda dizem `(unproved)` até regeneração.
4. Findings para o profile:
   - descoberta de pacote da lane de integração por `//go:build integration` nas 5 primeiras
     linhas deixa `tenant_config` e `internal/composition` fora de qualquer lane;
   - `apps/web/vitest.config.ts` inclui `feature-products` por nome de arquivo exato, não por
     glob, então um segundo arquivo de teste ali roda em lane nenhuma;
   - critério de aceite pode ser vácuo contra o TIPO (achado do S12): "nunca chama X" não
     reprova se o tipo do colaborador não declara X.
5. Fila pós-merge: chip de unificação da superfície de erro HTTP, já ratificado pelo operador
   como o próximo.

`pg-session-down` no fecho.

## 5. Fora de escopo, explicitamente

- Scrub de PII dos dumps `docs/design/evidence/ml-api/` (untracked na main, seam do hub).
- Regeneração de candidatos para limpar motivos antigos (decisão de dado).
- Unificação da superfície de erro HTTP (chip próprio, já enfileirado).
- Qualquer edição em `.gitignore` ou `tsconfig` (retidos como seam do hub por A-25 R-2b).

## 6. Ordem de execução

1. Recriar container de sessão, migrar, provar `applied 72` → `applied 0`.
2. F1 com must-fail; commit.
3. F2 com prova por golden + grep; commit com o grant registrado.
4. F3; commit.
5. Escrever `EVIDENCE.md` (números pendentes).
6. Escada P5 completa; colar saídas no `EVIDENCE.md`.
7. Congelar árvore. P6 nos dois assentos.
8. Tratar vereditos. Rodada extra se discordarem.
9. `CLOSED` + `REQUEST`. `pg-session-down`.
