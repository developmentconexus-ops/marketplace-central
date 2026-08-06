# Protocolo de módulo — arquitetura correta, padronizada, e os instrumentos que a mantêm

**Data:** 2026-08-05 · **Estado:** EM EXECUÇÃO — **Onda 0 FECHADA** (0.0–0.5), Ondas 1–4 abertas

> **Onda 0, o que de fato saiu** (medido, não alegado): 32 ADRs vivos em
> `docs/architecture/decisions/` (o plano previa 17); 001/002/003 recuperados do git;
> `_citations/RENUMBERING-REGISTRY.md` com o crosswalk global; ADR-023 = o protocolo de módulo
> (D-1…D-7 ratificados, §7 com a lista de camadas fechada). Duas afirmações do §1 deste plano
> foram **medidas e refutadas** — estão registradas em "Correções ao plano da Onda 0" no
> registry. O fecho da 0.3 (grafia/ponteiro de citação) está em "Fecho da task 0.3" no mesmo
> arquivo, com 1 pendência escrita.
**Origem:** decisão do operador em 2026-08-05 — parar antes da Fatia B da MIS-008 e consertar a
arquitetura enquanto o código ainda é pequeno.

---

## 0. O que este plano é, e o que ele não é

**É:** fechar o contrato entre módulos. Um módulo passa a ter lógica interna privada e um
consumível público, com um protocolo único de como publica e como se consome. E os instrumentos
que impedem o padrão de se desfazer de novo.

**Não é:** reescrever o sistema. A medição (§1) diz que dois eixos estão ruins, um médio e
**dois estão limpos**. O escopo desta missão são exatamente os números do §1. O que não estiver
neles fica fora, ou entra no §6 como dívida declarada com dono.

**Por que agora, e não depois da Fatia B:** a task B4 da MIS-008 é literalmente
"OpenAPI + SDK + handler, um commit" — ou seja, B4 escreve à mão, de novo, as quatro cópias da
mesma forma que o §1.4 mede como o pior defeito do sistema. Fazer B antes é pagar o defeito mais
uma vez, de propósito, depois de tê-lo medido.

**Custo aceito:** a Fatia A da MIS-008 (commits `476920b6`, `d57ea44b`) fica sem live drive até
esta missão terminar. Ela está commitada e verde em lane; não apodrece. O que apodrece é a
medição do §1 — por isso ela vem primeiro.

---

## 1. A medição

Toda a base deste plano. Medida em 2026-08-05 contra `d57ea44b`, **por mim, na sessão** — não por
subagente. Motivo: no censo anterior eu refutei metade das alegações dos subagentes com
verificação própria (76 dos "122 bypasses" eram `adapters/`, que é o ponto sancionado; as duas
"tabelas com dois escritores" eram particionadas por `source` no PK com CHECK no banco). Plano
construído sobre alegação não verificada é a própria classe de defeito que esta missão existe
para matar.

### 1.1 ADRs — 17 números citados, 8 documentos já existiram, 9 nunca existiram

Números normalizados (`ADR-04` e `ADR-004` são o mesmo). Citações contadas fora de
`scripts/.runs/*/snapshot/`, que duplica árvores inteiras e inflaria tudo.

| Medida | Valor |
|---|---|
| Números de ADR citados | **17** (`ADR-1` … `ADR-17`) |
| Documentos ADR que **já existiram**, em qualquer ponto da história do git | **8** (001–008) |
| Arquivos vivos hoje em `docs/architecture/decisions/` | **5** (004–008) |
| Números que **nunca existiram, em forma nenhuma** | **9** (09, 10, 11, 12, 13, 14, 15, 16, 17) |
| Citações somadas desses 9 órfãos | **1.365** |
| Citações de `ADR-17` sozinho | **1.076** |
| Citações de `ADR-17` **dentro do OpenAPI publicado** | **11** (`contracts/api/marketplace-central.openapi.yaml:5721…5920`) |
| Esquemas de numeração convivendo | **2** — `ADR-04` e `ADR-004`. Só o de 3 dígitos tem arquivo |

`git log --all --diff-filter=A` sobre todos os diretórios de decisão do repo devolve exatamente
oito arquivos: `001`…`008`. **Nenhum `009`–`017` jamais foi criado.** Nenhum `.md` do repo tem
`ADR-17` como título — as 1.076 ocorrências são citações a um documento que nunca houve.

**Os três primeiros têm história, e o caso de cada um é diferente** (deletados pelo commit
`99540107` "refactor(governance): retire brain knowledge store", quando o `.brain/` foi aposentado
e 004–006 foram movidos para `docs/architecture/decisions/`):

| Número | Estado declarado no próprio documento | Citações vivas | O que fazer |
|---|---|---|---|
| ADR-001 | `superseded on 2026-07-07` | 28 | deleção legítima. Restaurar como **superseded**, só para a citação não apontar para o nada |
| ADR-002 | `superseded on 2026-07-07` | 36 | idem |
| ADR-003 | **`accepted`** | 57 | **deletado por engano** junto do `.brain/`. Restaurar |

**Leitura:** a arquitetura deste sistema é governada por folclore. Nove regras — incluindo a mais
citada do repo, repetida 1.076 vezes e citada dentro do contrato publicado — não têm o que ler. Um
agente novo (ou uma pessoa nova) só descobre o que ADR-17 exige inferindo das citações. É a
causa-raiz mais barata de consertar e a que mais rende: **as regras já existem, só não estão
escritas.** E três delas nem precisam ser escritas — só recuperadas do git.

### 1.2 Vocabulário de camada — 12 nomes, 8 diretórios vazios

21 módulos em `apps/server_core/internal/modules/`. Nomes de diretório-camada em uso:

| Nome | Módulos | Situação |
|---|---|---|
| `domain` | 19 | canônico |
| `ports` | 18 | canônico |
| `adapters` | 18 | canônico |
| `transport` | 17 | canônico |
| `application` | 17 | canônico |
| `readmodel` | 4 | **cada um contém 1 arquivo `doc.go` de 1 linha. Zero código.** |
| `events` | 4 | **idem — 1 arquivo `doc.go` de 1 linha. Zero código.** |
| `composition` | 4 | tem código; **colide com o `composition/` de topo** |
| `background` | 2 | tem código (schedulers de `integrations`) |
| `registry` | 1 | tem código (`marketplaces`) |
| `observability` | 1 | tem código (`internal_read`) |
| `integration` | 1 | **só testes** (`mutations/integration/*_test.go`) |

**8 diretórios são andaime vazio** — `package readmodel` e nada mais. Fazem a arquitetura parecer
mais rica do que é, e é exatamente esse tipo de fachada que um revisor lê como "tem CQRS aqui".

Dois módulos fora do molde: `sourcekind` (sem camadas — correto, é o tipo compartilhado sem
dependência, ADR-02) e **`tenant_config`, que só tem `transport/`** — sem domain, sem application,
sem ports.

### 1.3 Fronteiras — 35 violações reais, concentradas em 3 alvos

121 importações cross-module não-teste. **76 partem de `adapters/`**, que é o ponto de tradução
sancionado — não são violação. Sobram **35**, e elas não estão espalhadas:

| Alvo invadido | Vezes | Módulos que invadem |
|---|---|---|
| `connectors/domain` | 12 | integrations, orders, product_links |
| `internal_read/domain` | 8 | **catalog, inventory, listings, orders, product_links, profitability** (6) |
| `sync/application` | 4 | internal_read, listings, market, orders |
| `connectors/application` | 3 | integrations |
| `erp_import/internalread` | 2 | catalog, market |
| `erp_import/domain` | 2 | dashboard |
| `catalog/domain` | 2 | erp_import, internal_read |
| `integrations/domain` | 1 | dashboard |
| `listings/domain` | 1 | mutations |

**Leitura:** não são 35 problemas. São **3 módulos sem consumível publicado** (`connectors`,
`internal_read`, `sync`) e, por isso, todo mundo se serve das entranhas. É a descrição exata do que
o operador nomeou. Conserto = publicar o consumível dos 3, não corrigir 14 chamadores.

Agravante medido: **9 dessas violações partem de `X/ports`** — ou seja, a porta de um módulo está
tipada no `domain` de outro. A porta, que deveria ser a fronteira, é o próprio vazamento
(`orders/ports/*` ×4 para `connectors/domain`, `dashboard/ports/sources.go`,
`inventory/ports/dashboard.go`, `profitability/ports/internal_read.go`).

### 1.4 Contrato e dinheiro — o pior eixo

| Medida | Valor |
|---|---|
| Operações no OpenAPI | 111 |
| Interfaces no SDK | **172, escritas à mão**, em 2.595 linhas. Sem codegen, sem banner `@generated` |
| Arquivos de transport com DTO à mão | 38 (de 62) |
| Checks de **acordo** entre as cópias | **0** |

`GOV_API_SDK_SPLIT` (`scripts/harness/Policy.psm1:458`) exige que OpenAPI e SDK mudem no **mesmo
commit** — atomicidade, não concordância. Nada pega um campo que o SDK declara e o OpenAPI não tem.

**Dinheiro no wire: 4 representações simultâneas.**

| Representação | Onde |
|---|---|
| `string` | decomposição de `pricing` |
| `float64` / `*float64` | **17 campos de valor absoluto**, incluindo `orders/transport/http_handler.go:451-466` — `Comissao`, `Frete`, `Imposto`, `Custo`, `MargemValor`, `TarifaFull` |
| `domain.Money` | 3 campos |
| `int` / `*int` | 4 campos |

O caso que resume tudo: **a mesma decomposição fiscal sai `string` em `pricing` e `float64` em
`orders`.** E ADR-17 (o tal das 1.068 citações) proíbe `float64` no caminho do dinheiro.

Outros dois vazamentos confirmados: `profitability/transport/http_handler.go:62,80,95` serializa
`profitabilityapp.ImportMarginInputsResult` — um tipo de **application** — direto no wire. O
contrato público daquela rota é definido por struct tags de um tipo interno.

### 1.5 Erro — envelope certo, vocabulário quebrado

Envelope: **1 forma só**, `apierror.Write` (`apps/server_core/internal/platform/apierror/apierror.go:25`),
192 chamadas, 1 tipo de erro no SDK. Este eixo está **certo** e não precisa de trabalho estrutural.

O vocabulário, não:

| Medida | Valor |
|---|---|
| Códigos distintos realmente emitidos | **42** |
| Documentados no OpenAPI | **8** |
| **Ausentes do contrato** | **34 (81%)** |

`CATALOG_PRODUCT_NOT_FOUND`, `ORDERS_INTERNAL_ERROR`, `MARKETPLACES_POLICY_INVALID`… o cliente
recebe códigos que o contrato nunca prometeu. O FE não tem como tratá-los sem lê-los do código Go.

### 1.6 Cerimônia — 7 portas fantasma verificadas uma a uma

196 interfaces nos módulos, **126 em `ports/`**. Destas, **7 não têm nenhuma referência não-teste
fora do próprio arquivo de declaração** — nada as implementa nomeadamente, nada as recebe como
parâmetro, nada as declara como campo. Verificadas individualmente com `grep -w` no repo inteiro:

| Módulo | Porta | Arquivo |
|---|---|---|
| pricing | `CalcEngine` | `pricing/ports/calc_ports.go:21` — tem comentário cuidadoso explicando que é "a porta de decomposição" |
| integrations | `AuthProvider` | `integrations/ports/auth_provider.go:9` |
| integrations | `EncryptionService` | `integrations/ports/encryption_service.go:5` |
| integrations | `ProviderRegistry` | `integrations/ports/provider_registry.go:16` |
| inventory | `InternalStockReader` | `inventory/ports/dashboard.go:24` |
| listings | `PageSource` | `listings/ports/ingestion.go:26` |
| mutations | `LinkageReader` | `mutations/ports/linkage.go:9` |

Outras 13 portas têm referência não-teste em **um único arquivo**. Estas **não** são
necessariamente ocas — em Go o adapter satisfaz a interface implicitamente, sem nomeá-la, então
uma porta citada só pelo serviço que a consome é o padrão normal. Ficam para triagem na Onda 3,
não para deleção automática.

`EncryptionService` sem consumidor merece nota à parte: é uma porta de criptografia que ninguém
usa. Verificar na Onda 3 se há segredo de integração sendo persistido sem cifra — se houver, isso
sobe de prioridade imediatamente e sai desta missão como incidente.

### 1.7 Banco — limpo

67 tabelas. **0 cruzamento de escrita** entre módulos. 6 cruzamentos de leitura. 1 vazamento de SQL
(`tenant_config/repository.go:40`). As duas "tabelas com dois escritores" do censo anterior caíram
na verificação: `products_mirror` tem PK `(tenant_id, source, codigo_produto)` com CHECK no banco —
os dois escritores são particionados por `source` **por projeto**, é o design de duas fontes que já
foi ratificado.

**Este eixo não entra no escopo.** Um item avulso (o SQL de `tenant_config`) vai junto da Onda 2.

### 1.8 O gate — existe, está vermelho, é cego, e ninguém o roda

- `scripts/harness.ps1` já é o entrypoint único (`unit|integration|live|browser|provider-write|governance-*`).
- `Test-GovernanceDrift` reporta hoje **55 violações**. Baseline vermelho.
- **Cego por uma palavra:** `Policy.psm1:328` testa `$layer -in @('adapters','transport','registry')`.
  `domain`, `application` e `ports` são estruturalmente invisíveis — as 35 violações do §1.3 não
  podem ser detectadas pelo checador que existe para detectá-las.
- **Legaliza por declaração:** 48 arestas de dependência declaradas em `contracts/governance/modules.json`,
  sem prazo e sem motivo. Declarar uma vez legaliza para sempre.
- **Nada o roda.** `.github/workflows/` tem só `release-images.yml` (build/push, sem testes).
  `.git/hooks/` está vazio.
- **Exige runtime que não declara:** rodado sob Windows PowerShell 5.1, `Test-GovernanceContracts`
  reporta `GOV_SCHEMA_INVALID` nos 6 documentos de governança — porque `Test-Json -SchemaFile` não
  existe no 5.1. Sob `pwsh` 7: passa. O checador mente dizendo "documento inválido" quando o
  problema é o interpretador.

---

## 2. As decisões — o protocolo de módulo

Isto é o que vira ADR na Onda 0. Está escrito aqui em prosa para o operador ratificar antes de
virar documento normativo.

### D-1 · O que é um módulo

Um módulo é uma entidade com **lógica interna privada** e **um consumível público**. Outros módulos
falam com ele só pelo consumível, sempre do mesmo jeito.

### D-2 · A regra de importação (a regra central, e é uma linha)

> **Um módulo é importável por outro módulo APENAS em `X/ports`.
> `X/domain`, `X/application`, `X/adapters` e `X/transport` são privados.**

É a tradução exata do que o operador descreveu, e é verificável com um grep. As 35 violações do
§1.3 são precisamente as violações desta linha.

### D-3 · A porta tem que se bastar

> **Uma interface em `X/ports` só pode ser tipada com: tipos declarados no próprio `X/ports`, tipos
> da biblioteca padrão, ou tipos do núcleo compartilhado (§D-5). Nunca com `Y/domain`.**

Sem isto, D-2 é teatro: o import vira legal mas o acoplamento continua. Mata as 9 portas vazadas.

### D-4 · Quando uma porta existe

> **Uma porta existe quando há dois ou mais consumidores, OU quando o cruzamento é de tecnologia
> (banco, HTTP de provider, Oracle, relógio, fila). Nunca por simetria.**

Contra-regra explícita. Sem ela, "padronizar" produz mais 126 interfaces. As 7 fantasmas do §1.6
existem porque ninguém tinha esta regra escrita.

### D-5 · Núcleo compartilhado, deliberadamente pequeno

Pacotes sem dependência que qualquer módulo pode importar. Hoje: `sourcekind`. Vai receber:
**o tipo de dinheiro** (§D-6). Regra: entra no núcleo só o que **dois ou mais** módulos precisam
declarar em assinatura pública, e o pacote não pode importar módulo nenhum. Lista fechada no ADR;
adicionar exige emenda.

### D-6 · Um dinheiro só

> **No domínio: um tipo `Money` único, no núcleo compartilhado, sem `float64`.
> No wire: decimal como `string`.**

`string` e não número porque JSON number é `double` e ADR-17 proíbe `float64` no caminho do
dinheiro — um campo `float64` no DTO viola a regra mais citada do repo em 17 lugares hoje.

Consequência assumida: mudar `orders` de `float64` para `string` **quebra o contrato**. Sai em um
commit só (OpenAPI + SDK + handler + FE), como B4 já previa fazer, mas uma vez e para valer.

### D-7 · O vocabulário de camada é uma lista fechada

Canônicos: `domain`, `application`, `ports`, `adapters`, `transport`.

Os 7 extras do §1.2 são julgados na Onda 0 e cada um vira uma de três coisas: (a) canônico, com o
ADR dizendo o que é e quando se usa; (b) dobrado numa camada canônica; (c) deletado.
Julgamento provisório, a confirmar lendo o código:
- `readmodel`, `events` → **deletar** (8 diretórios, zero código, andaime).
- `composition` dentro de módulo → renomear; colide com o `composition/` de topo, que é o dono
  da fiação. Provável destino: `application`.
- `background` → provável canônico (schedulers são um tipo real de entrada, como `transport`).
- `registry`, `observability` → provável `adapters`.
- `integration` (só teste) → mover para o layout de teste, não é camada.

Depois disso a lista é fechada e o checador rejeita nome novo.

### D-8 · O contrato é um só, e alguém verifica

Enquanto o SDK não for gerado, **existe um teste que lê o OpenAPI e compara o conjunto de campos
com as json tags do DTO Go**. Drift silencioso vira lane vermelha.

Gerar as 172 interfaces do SDK a partir do OpenAPI é o conserto real, e é **missão própria** —
fica como dívida declarada (§6), não entra aqui.

### D-9 · Todo código de erro emitido está no contrato

Todo literal passado como `code` para `apierror.Write` tem que aparecer no enum do OpenAPI.
Hoje: 8 de 42.

### D-10 · ADR citado tem que existir

Citar `ADR-N` onde não há documento `N` é falsidade em documento normativo. Vira check.

### D-11 · A regra do instrumento

> **Regra sem checador é comentário. Checador sem gate é log. Checador sem teste de must-fail é
> decoração. E checador tem que declarar o runtime que exige e falhar dizendo isso** — nunca
> reportar "documento inválido" quando o problema é o interpretador.

Vale para tudo que esta missão produzir. Nenhum check entra sem um teste que prove que ele
**reprova** o caso que existe para pegar.

---

## 3. As ondas

Ordem obrigatória. Cada onda fecha antes da próxima abrir.

### Onda 0 — escrever as regras (nenhum código de produção muda)

| # | Task | Aceite (medição, não teste verde) |
|---|---|---|
| 0.0 | Restaurar ADR-001/002/003 de `99540107^` — 001 e 002 marcados `superseded`, 003 é `accepted` e caiu por engano | 3 arquivos de volta; 121 citações deixam de apontar para o nada. Custo: `git show`, não redação |
| 0.1 | Colher as 1.365 citações dos **9 números que nunca existiram** (09–17); agrupar por o que cada uma afirma | Um arquivo por número, listando as afirmações distintas com `file:line` |
| 0.2 | Escrever os **9 ADRs faltantes** a partir do que as citações afirmam. Onde as citações se contradizerem, a contradição vai escrita no ADR e o operador decide. ADR-17 primeiro — 1.076 citações e está dentro do OpenAPI | `ls docs/architecture/decisions/` cobre todos os 17 números citados |
| 0.3 | Unificar a numeração (`ADR-04` vs `ADR-004`) e atualizar o README de decisões | 1 esquema, 0 citações no esquema morto |
| 0.4 | Ler as 7 camadas não canônicas e ratificar D-7 caso a caso | Lista fechada escrita no ADR de arquitetura |
| 0.5 | Escrever o **ADR do protocolo de módulo** — D-1 a D-7 | Documento existe e é citável |

**Por que primeiro:** sem D-2/D-3/D-4 escritos, cada fatia das ondas seguintes rediscute a regra do
zero, e o revisor não tem contra o quê reprovar.

### Onda 1 — o consumível dos 3 módulos invadidos (mata o §1.3)

| # | Task | Aceite |
|---|---|---|
| 1.1 | Publicar o consumível de `internal_read` — a fatia do active-source sai de `erp_import/adapters/internalread` para pacote próprio com tipo próprio, desfazendo a colisão de vocabulário `ImportSource` | 6 módulos consumidores param de importar `internal_read/domain`; contagem cai de 8 → 0 |
| 1.2 | Publicar o consumível de `connectors` | 3 consumidores; 12 → 0 |
| 1.3 | Publicar o consumível de `sync` | 4 consumidores; 4 → 0 |
| 1.4 | Fechar as 9 portas tipadas em `domain` alheio (D-3) | 9 → 0 |
| 1.5 | Resto: `erp_import`, `catalog`, `integrations`, `listings` (8 restantes) | 35 → 0, ou o resíduo vira baseline com motivo e dono |
| 1.6 | `tenant_config`: SQL fora do módulo (`repository.go:40`) e decidir se um módulo só-`transport` é legítimo | Decisão escrita; SQL onde a regra manda |

### Onda 2 — o contrato e o dinheiro (o pior eixo, §1.4)

| # | Task | Aceite |
|---|---|---|
| 2.1 | `Money` único no núcleo compartilhado (D-5, D-6) | 1 tipo; `grep float64` no caminho do dinheiro do domínio = 0 |
| 2.2 | Wire: decimal string. Migrar os 17 campos `float64` — **`orders` primeiro**, é o maior | 17 → 0. OpenAPI + SDK + handler + FE no mesmo commit |
| 2.3 | Tipo de application fora do wire (`profitability:62,80,95`) | DTO próprio; nenhum handler serializa tipo de application |
| 2.4 | **Teste de acordo OpenAPI ↔ DTO** (D-8): lê o YAML, compara conjunto de campos com json tags | Must-fail provado: apagar um campo do YAML reprova |
| 2.5 | Enum de código de erro no OpenAPI + check (D-9) | 42/42 documentados; must-fail provado |

### Onda 3 — deleção (§1.6, §1.2)

Onda mais barata, maior ganho de legibilidade. É só apagar.

| # | Task | Aceite |
|---|---|---|
| 3.1 | Deletar as 7 portas fantasma | Compila; 126 → 119 |
| 3.2 | **Antes de 3.1**: verificar se `EncryptionService` sem consumidor significa segredo de integração sem cifra | Se sim, vira incidente e sai desta missão |
| 3.3 | Deletar os 8 diretórios `readmodel/` e `events/` vazios | 12 nomes de camada → 10 |
| 3.4 | Aplicar o ADR de D-7 aos 5 nomes restantes | Lista fechada em vigor |
| 3.5 | Triagem das 13 portas de referência única — **triagem, não deleção automática** | Cada uma: mantida com motivo, ou apagada |

### Onda 4 — o gate (por último, sempre)

| # | Task | Aceite |
|---|---|---|
| 4.1 | Estender `GOV_MODULE_LAYER` para `domain`/`application`/`ports` (`Policy.psm1:328`) | Must-fail: reintroduzir uma das 35 reprova |
| 4.2 | Check de D-3 (porta se basta) | Must-fail provado |
| 4.3 | Check de D-7 (nome de camada fora da lista) | Must-fail provado |
| 4.4 | Check de D-10 (ADR citado sem documento) | Must-fail provado |
| 4.5 | `contracts/governance/modules.json`: toda aresta de dependência exige `reason` + `removal_owner` | 48 arestas justificadas ou removidas |
| 4.6 | Baseline do resíduo, com motivo e dono cada | 0 violação sem justificativa escrita |
| 4.7 | Todo checador declara o runtime que exige e falha dizendo isso (D-11) | Rodar sob 5.1 diz "precisa pwsh 7", não "documento inválido" |
| 4.8 | **Só então:** hook pre-commit via `core.hooksPath` em `scripts/githooks/` — governança + gofmt (segundos). Lanes completas em pre-push / fecho de fatia | Lane verde de verdade antes de ligar |

**4.8 é a última linha da missão.** Gate ligado em lane vermelha não é gate — vira `--no-verify`,
que é pior que gate nenhum.

---

## 4. Como cada task se executa

Vale a regra do operador, sem exceção:

> **Todos os RED primeiro, contra árvore limpa. Depois implementa. Depois verifica o que era pra
> dar green. Se não deu, cria RED para as situações novas — e itera até a implementação condizer
> com o que é pra ser, nunca o teste condizer com a implementação.**

Mais três, aprendidas em custo:
1. **Aceite é medição, não teste verde.** Contagem que muda, `count(*)`, campo que sai do wire.
   Teste verde prova que o teste passou.
2. **Alegação sobre o repo apodrece.** Todo `file:line` deste plano é re-verificado no momento em
   que a task abre. Este plano não é fonte de verdade sobre o código — o código é.
3. **Nenhum check entra sem must-fail provado** (D-11).

---

## 5. Tamanho, honestamente

| Onda | Peso | Risco |
|---|---|---|
| 0 | médio — 3 recuperados do git + 9 escritos, e o material já existe nas citações | baixo; não toca código |
| 1 | **maior da missão** — 3 consumíveis novos, 14 chamadores | médio; mexe em módulo que todo mundo usa |
| 2 | grande — quebra contrato público de `orders` | **maior risco**; FE junto, um commit |
| 3 | pequena — deleção | baixo; compilador é o juiz |
| 4 | média | baixo, mas só funciona se 1–3 fecharam |

Onda 0 e Onda 3 são independentes de tudo e podem correr em paralelo com a 1.
Onda 2 depende da 1 (o consumível tem que existir antes de o tipo de dinheiro atravessá-lo).
Onda 4 depende de todas.

---

## 6. Fora de escopo — dívidas declaradas, com dono

| Dívida | Por que fora | Dono |
|---|---|---|
| **Codegen do SDK a partir do OpenAPI** | 172 interfaces com prosa escrita à mão; é missão própria. A Onda 2.4 põe o teste de acordo no lugar como paliativo verificável | operador |
| `channelfees` e `divergences` sem `application`/`transport` | são `status: planned` da MIS-007 M-05 — trabalho inacabado, não código morto. **Não deletar** | MIS-007 |
| `GOV_MIGRATION_PREFIX` (prefixo duplicado) | `schema_migrations` é chaveado por **filename** (`migrate/runner.go:51`); renomear reaplica a migração | baseline permanente |
| 9 `RCFG_READER_MISSING` em adapters Magalu/Amazon/Shopee | dívida real de conector, não de arquitetura | backlog de conectores |
| 18 `RCFG` em arquivos de ops | **uma** decisão de escopo, não 18. Precisa do operador | operador |
| Fatias B, C, D da MIS-008 | retomam depois, em sessão nova | operador |

---

## 7. O que fecha a missão

1. As 35 violações do §1.3 em zero, ou baseline com motivo e dono por linha.
2. Uma representação de dinheiro no domínio, uma no wire.
3. 42/42 códigos de erro no contrato.
4. Nenhum ADR citado sem documento.
5. Vocabulário de camada fechado e verificado.
6. Todo check com must-fail provado, rodando em hook, com lane verde de verdade.
7. **QA live drive** — precisa da dev stack, que é seam do hub e depende do operador.

E aí, sessão nova, do zero: Fatia B da MIS-008 — agora contra uma arquitetura que tem protocolo
escrito, dinheiro de um tipo só, e um checador que reprova quem sair da linha.
