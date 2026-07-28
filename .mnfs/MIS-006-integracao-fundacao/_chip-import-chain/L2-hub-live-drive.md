# L2 — drive ao vivo do CHIP-IMPORT-CHAIN, executado pelo HUB

`2026-07-28` · executor: hub · branch dirigido: `chip/import-chain` @ `67e4a3d`

Este é o L2 que cobra o waiver de `2026-07-28`: `GET /erp/imports/{id}/chain` landou no
CHIP-ANCHORS-2 e **nunca tinha sido dirigido ao vivo por ninguém**. O risco aceito no waiver foi
nominal: *um defeito de FIAÇÃO — composition root, decorator perdido — fica escondido até agora*,
que é a classe exata do catalog-503 do M-02.

**Veredito: a fiação está de pé. O waiver está pago.**

## Mecanismo — por que este drive não copiou credencial nenhuma

O waiver existia porque dirigir a partir do worktree do chip exigiria copiar `.env` (Sankhya
Oracle + ML) para um segundo diretório. Não foi o que aconteceu.

O stack subiu do **checkout principal**, com um override que repontua **só o bind mount do
frontend**:

```
docker compose -f docker-compose.yml -f <scratchpad>/l2-import-chain.override.yml up -d
```

O compose base continua sendo o do checkout principal, então `env_file: .env` e o `context:` do
build resolvem lá. Verificado no container:

```
C:/…/.claude/worktrees/sharp-pike-3387c1 -> /workspace
marketplace-central_node_modules       -> /workspace/node_modules
marketplace-central_web_node_modules   -> /workspace/apps/web/node_modules
```

Seguro porque o diff do chip é FE puro — zero Go, zero migrations, zero `contracts/`, zero
`packages/sdk-runtime/` contra `5441fe18`. O backend serve o mesmo código dos dois jeitos, então
apontar só o frontend não falsifica nada.

## Dado usado — real, não fabricado

A base do dev estava com `erp_import_protocols` vazia (o `#003-E` da demo não sobreviveu ao
volume atual). Em vez de inserir linha à mão — que seria fabricar o dado que o teste deveria
observar — subi um import pelo caminho de produção:

```
POST /erp/imports  (multipart, source=xlsx)
  file = .mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/fixtures/example-erp.xlsx
→ 201 em 0.24s
  {"import_id":"eac3ac9e-b87b-43be-959c-7bf1f11a07f1","protocol":"#001-E","status":"COMPLETED"}
```

## Os números na tela, VERBATIM

`/importacoes/eac3ac9e-b87b-43be-959c-7bf1f11a07f1`:

```
Protocolo #001-E
Produtos do import   55
Vinculados            0
Enfileirados         55
Fila lida em: 28/07/2026, 13:32
```

`GET http://localhost:8080/erp/imports/eac3ac9e-…/chain → 200 OK` (observado na aba de rede do
browser, não só no curl — é o que prova o **I4**: a tela consome o endpoint novo, não recalcula a
cadeia no cliente a partir de `listErpImports`).

## Conferência contra o DB — consulta direta, nunca contra a própria tela

O contrato exige que os números batam com o banco **conferidos por consulta direta**. A query que
o endpoint roda está em `query_repository.go:71-109`; reproduzi cada contador por fora:

| Contador | Tela | Consulta direta | Confere |
|---|---|---|---|
| Produtos do import | 55 | `count(*) from erp_import_products where protocol_id=…` → **55** | sim |
| Vinculados | 0 | `count(distinct p.codprod)` juntando `product_links` state=`resolved` → **0** | sim |
| Enfileirados | 55 | `jsonb_array_length(cursor->'pending')` de `sync_state` entity=`market` → **55** | sim |

### O `vinculados = 0` é VERDADEIRO, e não é o defeito de subcontagem

O contrato manda suspeitar de número baixo antes de culpar o FE. Suspeitei e **descartei**, com
duas medições:

- `count(*) filter (where codprod ~ '^0')` no import = **0**. Nenhum CODPROD com zero à esquerda
  neste conjunto, então o defeito que o CHIP-ANCHORS-3 conserta (`links.internal_product_id::text
  = products.codprod`, `query_repository.go:89`) não tem como morder aqui.
- As faixas são **disjuntas**: os CODPROD do import vão de `1001` a `1055`; os
  `internal_product_id` dos 29 vínculos resolvidos vão de `15956` a `42194`. Interseção zero por
  construção — a fixture do M-01 é catálogo sintético, não o catálogo Sankhya.

**Limitação que declaro em vez de esconder:** este drive **não exerceu um `vinculados` diferente
de zero**. Provar o caminho não-zero exige um import cujos CODPROD casem com vínculos resolvidos
de verdade, e isso é dado de cliente, não fixture. Fica registrado como o que o L2 **não** cobriu.

Em compensação ele exerceu o discriminante que importa para o ADR-17 pelo lado oposto: `0` é um
valor **conhecido** aqui, e a tela renderiza `0` — não `—`. Um `—` neste caso seria tão mentiroso
quanto um zero fabricado no caso ausente.

## Id inexistente — erro honesto, sem cadeia de zeros

```
/importacoes/00000000-0000-0000-0000-000000000000
→ "Erro ao carregar. Importação não encontrada."  + botão "Tentar novamente"
GET …/00000000-…/chain → 404 Not Found
```

Nenhum contador na tela. **I6 PASS**, e pelo padrão que o próprio contrato antecipou: 404
renderizado como erro é PASS de FE.

## I2 — os DOIS sítios de render, conferidos na tela

O pack original só citava `VinculosPage.tsx:159`; o segundo sítio (`IntegracoesPage.tsx:449`) foi
achado do chip. Ambos verificados ao vivo:

- **`/vinculos`** — a seção "Importação" **sumiu**. Confirmado lendo a página inteira.
- **`/integracoes`** — a seção "Importação" **está lá**, com o `#001-E`, os contadores
  Aceitos/Rejeitados/Avisos e os botões `Ver detalhes` / `Ver cadeia`. É o recibo do upload que o
  operador acabou de fazer, na tela onde ele fez.

Nenhuma tela sumiu sem registro.

## I1 — a rota e o gate de instalação

`/importacoes` monta e renderiza. Evidência observável da decisão de gate, além da leitura de
código: os links de navegação das rotas **gated** carregam
`?installation=inst-mercado_livre-…`, enquanto `/importacoes`, `/integracoes`, `/catalogo` e
`/estoque` são `href` limpos. Importação de ERP não depende de marketplace conectado, e a
navegação reflete isso.

## I7 / F-02 — o toggle de fonte ativa persiste E a app inteira muda de fonte

Este é o teste que o F-02 sempre prometeu e que nenhuma suíte conseguia provar.

```
antes:  GET /config/active-source → {"active_source":"sankhya","source_kind":"live_read_through"}
clique em "Planilha Sankhya"
depois: PUT /config/active-source → 200
        GET /config/active-source → {"active_source":"xlsx","source_kind":"upload_snapshot"}
        DB: active_source | xlsx | upload_snapshot | 2026-07-28 16:34:00.754941+00
```

E a invalidação global **funciona de verdade**: `/catalogo` deixou de servir o mirror Sankhya
(10529 linhas) e passou a servir o snapshot do xlsx — `REF-1001 Example Product 1001`,
`— (missing_price)`, custo `R$ 13,34`. Não é o card que mudou; é o app inteiro que passou a ler a
outra fonte, que é exatamente o que o `activeSource.ts` promete.

Estado restaurado ao final: `sankhya`, `set_at 2026-07-28T16:37:49Z`.

## Claro e escuro

Sem screenshot — o painel do browser não estava compondo frames — então medi valor computado, que
discrimina melhor que pixel:

| | `data-theme` | `body` background | `body` color | classes de cor literal em `main` |
|---|---|---|---|---|
| claro | `light` | `rgb(251, 250, 247)` | `rgb(37, 41, 31)` | **0** |
| escuro | `dark` | `rgb(22, 24, 20)` | `rgb(233, 232, 226)` | **0** |

A varredura procurou `-(slate|blue|emerald|red|gray|zinc|neutral|amber|orange|green|indigo|purple|pink|teal|cyan|yellow)-\d` em toda classe sob `main`. Zero nos dois temas: a tela nova usa token, não cor literal. Conteúdo e contadores idênticos nos dois.

## Anomalia observada, NÃO reproduzida — registrada mesmo assim

Numa das trocas de fonte ativa, um clique no rádio `sankhya` logo após o carregamento da página
marcou o input no DOM mas **não disparou `PUT` nenhum** (aba de rede vazia para
`/config/active-source` naquele intervalo). Recarreguei e o mesmo clique funcionou; tentei
reproduzir de novo com carga limpa e **não reproduziu**.

Não chamo isso de defeito — uma ocorrência sem repro é anomalia, e reportar anomalia como bug
gasta o chip atrás de fantasma. Registro porque a classe é real (input não-controlado cujo estado
de DOM descola do estado do React perde o `change`), e porque um leitor futuro merece saber que
foi visto uma vez.

**E registro também o meu próprio erro de método**, porque ele quase virou um achado falso: um
clique meu num link de menu não-visível caiu no card `catalogo_cliente` que estava embaixo, e
escreveu a fonte ativa sem eu pedir. Por alguns minutos eu tinha "tela discorda do servidor" na
mão. Discordava porque **eu** tinha escrito. Só o `set_at` no DB desmontou a hipótese. Anotado
como técnica: num drive ao vivo, um clique por coordenada em elemento de menu colapsado é uma
escrita cega.

## Fora do escopo deste drive

- `vinculados` não-zero (justificado acima).
- Campo ausente → `—` (I5): a rota viva sempre devolve inteiro; o discriminante é o teste
  unitário do chip. O live drive prova o lado oposto (`0` conhecido renderiza `0`).
- Correção de backend de qualquer espécie — inclusive a subcontagem, que é do CHIP-ANCHORS-3.
