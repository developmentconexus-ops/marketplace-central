# Método de Revisão — Marketplace Central

`status: VIGENTE desde 2026-08-03 · aplica-se ao fechamento de toda onda e antes do despacho da onda seguinte`
`doutrina superior: AGENTS.md · docs/HARNESS-PROFILE.md · ARCHITECTURE.md + ADRs`

## §0. Por que este documento existe

Uma onda fechada com "todos os testes passaram" não responde à única pergunta que
importa para o operador: **o que mudou na tela, de onde vem cada número, e o que
acontece quando esse número não existe.**

Verde de lane prova que o código faz o que o teste diz. Não prova que o teste diz a
coisa certa, que o campo tem produtor real, que o valor exibido não é um default
plausível, nem que existe alguém consumindo o que foi construído. Esta revisão
existe para fechar essa distância, e produz **dois artefatos por onda**:

- **A. Ficha de Entrega** (`docs/entregas/ONDA-<n>-ENTREGA.md`) — tela por tela,
  campo por campo, o que o operador vai ver e de onde cada coisa vem. Escrita em
  linguagem de operação, não de código.
- **B. Laudo de Revisão** (`docs/entregas/ONDA-<n>-LAUDO.md`) — o veredito técnico
  nos seis eixos abaixo, com `file:line` e saída de comando.

Regra que governa os dois: **não escrito = não aconteceu; não medido = não é fato.**
Toda linha de qualquer um dos dois artefatos aponta para uma medição.

---

## §1. Artefato A — Ficha de Entrega (a "cadeia do campo")

Este é o eixo principal e o que o operador lê primeiro. Para **cada tela** tocada
pela onda, uma seção; para **cada campo visível**, uma linha.

### 1.1 Formato da linha

| Campo (rótulo na tela) | Onde aparece | Operação SDK | `operationId` | Handler Go | Coluna Postgres | Origem externa | Se desconhecido | Veredito |
|---|---|---|---|---|---|---|---|---|

- **Rótulo na tela**: o texto literal que o operador lê, verbatim.
- **Onde aparece**: `apps/web/src/pages/<tela>/<Componente>.tsx:<linha>`.
- **Operação SDK**: método de `packages/sdk-runtime`.
- **`operationId`**: âncora em `contracts/api/marketplace-central.openapi.yaml`
  (âncora por conteúdo, nunca por linha — linha apodrece entre árvores).
- **Handler Go** → **serviço** → **porta** → **adaptador**: a travessia hexagonal.
  Se algum salto não existe, a cadeia está quebrada e o veredito não pode ser VIVO.
- **Coluna Postgres**: `tabela.coluna`. Postgres é o único estado canônico.
- **Origem externa**: endpoint do Mercado Livre (`GET /orders/search?...`), tabela
  Oracle (`TGFCUS`), planilha, ou **cálculo** (e então qual função pura o produz).
- **Se desconhecido**: o que a tela mostra quando o valor não existe. Esta coluna é
  a que pega violação de ADR-17. Respostas aceitáveis: `—`, "sem dado", badge de
  pendência. Resposta que reprova: `0`, `R$ 0,00`, `""`, `false`, data de hoje.
- **Veredito**: um dos cinco de §1.2.

### 1.2 Vocabulário de veredito (por campo)

| Veredito | Significa | Como se prova |
|---|---|---|
| **VIVO** | cadeia completa e valor real medido na tela com dado de produção | print/leitura da tela + a linha correspondente no banco |
| **ÓRFÃO** | a tela mostra, mas não há produtor a montante — o valor nasce no FE | busca do campo no handler/serviço não encontra produtor |
| **FANTASMA** | existe no contrato/banco, nenhuma tela consome | `operationId` sem chamador no `apps/web` |
| **MENTIRA** | desconhecido renderizado como valor plausível | ADR-17: o produtor devolve default em vez de ausência |
| **MUDO** | cadeia existe, mas na base real a coluna é sempre nula/vazia | `SELECT count(*) ... WHERE col IS NOT NULL` = 0 |

ÓRFÃO e MENTIRA são defeito. FANTASMA é desperdício (ou fatia futura, e então
tem que estar nomeada como tal). MUDO é o mais perigoso: passa em toda lane e a
tela fica verde — só aparece contando linhas na base real.

### 1.3 O que a Ficha responde ao operador

Cada seção de tela termina com três parágrafos em linguagem de operação:

1. **O que você vai ver** — o passo a passo do que abrir e o que aparece.
2. **O que isso te permite fazer** — a decisão que o operador consegue tomar agora
   e não conseguia antes.
3. **O que ainda não dá** — o limite honesto da fatia, com a dívida ou fatia futura
   nomeada (`D-xx` / `F-xx`).

### 1.4 Versão prospectiva (antes de despachar a onda seguinte)

A mesma ficha, escrita **antes** da implementação, com a coluna Veredito trocada
por **Prometido**. Serve de contrato de aceitação: no fechamento, cada linha
Prometido tem que virar VIVO ou virar dívida nomeada. Linha que some sem virar
nenhum dos dois é escopo perdido em silêncio.

---

## §2. Eixo 1 — Arquitetura e fronteiras

O que se mede, com `file:line`:

- **Camadas**: `domain / ports / application / adapters/<tech> / transport /
  composition`. `application` não importa `adapters` nem `transport`. Camada-alvo
  proibida: `adapters`, `transport`, `registry`. `composition` é legal.
- **Payload de provedor fica no adaptador**: nenhum DTO do Mercado Livre atravessa
  para `application` ou `domain`. O domínio fala a língua do nosso sistema, não a
  do fornecedor — é isso que faz a plataforma ser nativa e não um proxy.
- **`tenant_id` em toda consulta multi-tenant**, sem exceção.
- **Escrita em provedor** exige linkage resolvido, política/tempo de origem
  explícitos, proteção contra duplicata e auditoria.
- **Atomicidade de contrato**: OpenAPI e `packages/sdk-runtime` no mesmo commit
  (`GOV_API_SDK_SPLIT`); o SDK é mantido à mão.
- **Oracle é somente leitura**; Postgres é o único estado canônico (mirror-first).

Governança roda como lane, e o critério é **diff de conjunto**, nunca "zero
violação" — o main tem violações herdadas e verde absoluto é inalcançável:

```bash
cd "$(git rev-parse --show-toplevel)" && pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha "$(git rev-parse main)"
```

Dois cuidados que já custaram rodada:

- **Violação de governança é relação, não propriedade de arquivo.**
  `GOV_MODULE_COVERAGE` e `RCFG_READER_MISSING` são reportadas no nível de
  módulo/registro. Filtrar por caminho de arquivo é estruturalmente cego a elas. O
  único critério sadio é o diff de `(error_code, id)` entre base e merge, **medidos
  em ambientes equivalentes** (worktree destacado fora da árvore montada, gocache
  espelhado, base = tip do alvo do merge).
- **A lane varre `.claude/worktrees/`** e emite veredito sobre a árvore errada
  (B-10b: 88% dos blocos vinham de worktree alheio). Medir fora da árvore montada.

---

## §3. Eixo 2 — Verdade dos dados (ADR-17)

Regra: **fato operacional desconhecido nunca vira zero, `""`, `false` ou default
plausível.** Onde procurar, em ordem de retorno:

1. `coalesce(...)` em SQL sobre coluna que carrega fato de provedor.
2. `?? 0`, `|| ""`, `|| 0` no FE sobre campo vindo do SDK.
3. Struct Go com campo de valor (não ponteiro) onde a ausência é possível —
   ausência precisa ser alcançável no tipo, senão o zero é indistinguível do real.
4. Formatador que devolve `R$ 0,00` para `null`.
5. Cursor/timestamp que ganha "agora" quando a origem não informou.

Contrapartida do outro lado do ADR-17: **zero conhecido é `0`**. Suprimir um zero
verdadeiro é o mesmo defeito na direção oposta.

Anti-padrão já pago: guard parcial sob frase total. Um guard que cobre 3 dos 5
sítios produtores, sob um texto que promete cobertura total, é pior que guard
nenhum — cria confiança falsa. Tupla completa por sítio produtor, ou nada.

---

## §4. Eixo 3 — Máximo local, redundância e código morto

Esta é a pergunta "a gente fez um bom trabalho ou fez trabalho a mais?". Cinco
buscas, todas com resultado numérico:

1. **Motor duplicado.** Existe uma segunda implementação da mesma regra ao lado de
   uma que já funcionava? (Caso pago: o plano P2.b construía um segundo motor de
   preço ao lado do `pricing/domain` que já acertava.) Sinal: duas funções que
   calculam a mesma grandeza com assinaturas parecidas em módulos diferentes.
2. **Fórmula sem consumidor.** Para cada função exportada de `domain`, quem chama?
   Zero chamadores = fórmula morta. (Caso pago: margem com 4 fórmulas, 2 mortas.)
3. **Operação de contrato sem consumidor.** Cada `operationId` do OpenAPI tem
   chamador no `apps/web` ou consumidor externo declarado? (Caso pago: 11 operações
   sem consumidor.)
4. **Campo sem produtor.** O inverso: veredito ÓRFÃO da §1.2.
5. **Abstração que não abstrai.** Porta com uma única implementação e um único
   chamador, onde a interface só repete o concreto, é cerimônia. Porta com
   implementação real + implementação de teste **não** é isso — é o padrão certo.
   O critério é: a porta protege o domínio de uma dependência que muda? Se sim,
   fica mesmo com um implementador.

Regra de corte: YAGNI vale contra o que **ainda não** tem consumidor; não vale
contra o que já tem consumidor e está feio. Feio com consumidor = dívida nomeada.

---

## §5. Eixo 4 — Execução real (live drive)

Nenhuma fatia fecha sem alguém dirigir a tela no navegador contra dado real. As
regras que existem porque cada uma já falhou aqui:

- **Vermelho antes do verde, com controle negativo nomeado.** Um verde que também
  passaria no mundo sem a mudança não é evidência. O vermelho tem que **nomear** o
  teste/campo que quebrou.
- **Binário velho faz o live drive mentir.** Antes de medir, comparar a data de
  criação do container com o primeiro commit da fatia. Stack de pé ≠ stack com o
  teu código. Registrar o container id no veredito, para a janela ser
  reconstituível depois.

  ```bash
  docker inspect --format '{{.Created}}' <container>
  ```

  (Sempre com `--format` alvo; `docker inspect` cru despeja segredos no transcript.)
- **Tela verde não prova visibilidade de falha.** Provar que a falha aparece exige
  provocar a falha, não presumi-la.
- **Pulado e verde são byte-idênticos.** Lane precisa de contagem por linha
  (RUN/PASS/SKIP/FAIL) e de `failure_token=test=` para ser atribuível.
- **Não limpe o estado antes do QA.** O banco velho é a linha de base; limpar antes
  apaga exatamente a evidência que o QA existe para achar.
- **Escrita ao vivo no Mercado Livre exige autorização explícita do operador**, por
  escrito, uma por operação.

---

## §6. Eixo 5 — Lanes e evidência

Ordem de execução no fechamento de onda (todas com saída completa colada no laudo):

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go build ./...
```
```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test ./...
```
```bash
npm run test --workspace @marketplace-central/web
```
```bash
npm run harness:integration
```
```bash
cd "$(git rev-parse --show-toplevel)" && pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha "$(git rev-parse main)"
```

Armadilhas conhecidas das lanes:

- Lane de integração **não é superconjunto** da de unidade: só `./tests/integration`
  e só com `//go:build integration` nas 5 primeiras linhas; `transport` nunca compila lá.
- `npx --no-install tsc` pode passar por vacuidade; `tsc` **não é lane declarada**
  neste repo — erro de `tsc` vira dívida, não bloqueio, a menos que esteja em
  arquivo tocado pela onda.
- Worktree sem `node_modules` resolve `@mc/*` para o branch do main — medição
  cross-branch contaminada.
- Sempre `cd apps/server_core` antes de usar `GOCACHE`, por causa do gitignore da raiz.

**Evidência válida** = comando + saída completa + `file:line` + contagem. Alegação
sem proveniência de medição não entra em laudo, nem em card, nem em brief.

---

## §7. Eixo 6 — Ordem de verdade e conflitos

Quando duas fontes discordam, a ordem é:

`ARCHITECTURE.md`/ADRs → OpenAPI + `sdk-runtime` → `contracts/governance/` → wiki →
`.mnfs/` → testes/builds/commits.

Conflito de arquitetura, contrato, runtime, propriedade ou verificação **para a
linha**: classificar antes de continuar. Defeito reincidente vira parada com
root-cause e conserto geral ou dívida registrada — nunca remendo local repetido.

---

## §8. Checklist de fechamento de onda

- [ ] Ficha de Entrega escrita, uma seção por tela, uma linha por campo
- [ ] Todo campo com veredito; ÓRFÃO/MENTIRA viraram defeito aberto; FANTASMA/MUDO viraram dívida nomeada
- [ ] Live drive executado em navegador real, com container id e `Created` no veredito
- [ ] Vermelho-antes-do-verde registrado com controle negativo nomeado
- [ ] Cinco buscas do §4 executadas, com números
- [ ] Governança por diff de conjunto contra o tip do alvo, medida fora da árvore montada
- [ ] Lanes do §6 com saída colada e contagem por linha
- [ ] Dívidas novas em `.mnfs/HARNESS-DEBTS.md` com numeração conferida contra o arquivo (três `D-19` colidiram em 2026-08-03)
- [ ] Ficha prospectiva da onda seguinte escrita antes do despacho

## §9. Checklist de abertura de onda

- [ ] Ficha prospectiva (§1.4) com coluna Prometido preenchida por tela e por campo
- [ ] Para cada campo prometido: origem externa já medida (endpoint existe, devolve o campo, com payload real anexado) — promessa sobre API não medida é a fonte nº 1 de fatia estourada
- [ ] Colisão de seam declarada contra as outras fatias da onda
- [ ] Critério de aceitação de cada fatia legível contra a **assinatura** do tipo
      (critério que o tipo não consegue reprovar é vácuo)
