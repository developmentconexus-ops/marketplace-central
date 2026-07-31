# CHIP-VENDAVEL — regra de sortimento vendável (context pack)

Base de despacho: `main @ cf741f68`. Chip interino pós-MIS-006, ordenado pelo operador antes do
planning da MIS-007. Track único — nenhum chip concorrente em voo.

## Objetivo

O operador abre `/catalogo` e vê 10.538 produtos, a maioria morta para venda. A regra de
"vendável" vira configuração por tenant com efeito em todas as leituras relevantes:

- Card **"Sortimento vendável"** em `/integracoes` (ao lado do card de fonte ativa): 3 toggles
  persistidos no banco + linha "Resultado: N de M produtos" computada ao vivo.
- `/catalogo` abre filtrado com chip contador "Vendáveis N de M" + botão "ver todos" (escape de
  tela, não muda config). Produto do sortimento com estoque zero = badge, nunca corte extra.
- Geração de vínculos considera só vendáveis.

## Regra ratificada (operador, 2026-07-29 — medida na base viva METALPRD)

| toggle | semântica | default | corte medido |
|---|---|---|---|
| `only_revenda` | `TGFPRO.USOPROD = 'R'` | **true** | 10.538 → 10.007 |
| `only_em_estoque` | estoque disponível > 0 em CODEMP IN (1,2); recorte de local = ver EMENDA A7 abaixo | **true** | ver A7 (2.923) |
| `only_ecommerce_eligible` | `NVL(AD_ECOMMERCE,'X') <> 'N'` (live) / `IS NULL OR <> 'N'` (espelho) | **false** | corta só o que o ERP marcou `'N'` |

Fatos que NÃO são graus de liberdade do chip:
- "Tipo de estoque" (`TIPCONTEST`) é pista falsa — tipo de CONTROLE, não vendabilidade. Fora.
- `AD_ECOMMERCE` entra como coluna e toggle desligado, não como filtro duro.
- Critério com dado ausente NÃO exclui (honest-unknown): linha xlsx sem `usoprod` passa o
  filtro de revenda. Ausência ≠ reprovação.
- F-LINK-1 refutado: espelho Sankhya casa 31/31 EANs da conta ML (exato, sem normalizar);
  xlsx casa 0. Nada de normalização de EAN neste chip.

## EMENDA A7 — recorte de localização de estoque (ratificada pelo operador 2026-07-29, pós-despacho)

Verbatim do operador, dito na sessão do chip: *"CODLOCAL é legal: 10101 é estoque revenda, 10108 é
show room não é pra contar, e 10102 é outlet — esses que vendem"*. Perguntado se valia só para o
filtro ou também para o número em tela: **os dois**.

**Disponível vendável = `CODEMP IN (1,2) AND CODLOCAL IN (10101, 10102) AND CODPARC = 0`.**
10108 (show room) fora. Consequências fixadas:
- O pin do Q4 do sync é COMPLETO (empresa **e** local); o must-fail falha se qualquer um dos dois sair.
- `catalog_page.go` hoje usa `CODLOCAL = 10101` → passa a `IN (10101, 10102)` (outlet entra; o
  operador quer essa mudança de número).
- **O recorte vive em DOIS sítios no live** (medido 2026-07-29, ruling #3): a CTE do
  `catalog_page.go` E `domain/internal_stock.go` `DefaultSellableStockPolicy()` (`LocationIDs`),
  da qual `reader.go` e `stock_batch_reader.go` constroem as queries. A policy passa a
  `{10101, 10102}`; `ExcludedLocationIDs` default vira `nil` (blacklist vacuosa — o `NOT IN`
  sobre linhas já restritas à whitelist nunca excluiu nada; whitelist nunca blacklist).
  `domain/contract_test.go` assere os valores novos e é o pin da policy. Divergência de
  quantidade entre `/catalogo` e `EstoqueTab` para produto só-outlet NÃO é intencional — "os
  dois" do operador cobre todo número em tela.
- **Correção do chip ao ruling #3 (aceite, medida nos dois assentos):** `buildNotIntListClause`
  com lista vazia gera ` AND ... NOT IN ()` = ORA-00936, e `stock_batch_reader.go:115` chama SEM
  guarda (reader.go:503 tem). Esvaziar a policy sem conserto rebentaria só contra Oracle vivo.
  Conserto no HELPER (lista vazia → `""`, "sem cláusula"), não no chamador — guarda-no-chamador
  é o padrão que acabou de falhar (dois chamadores, um guardado). `buildIntListClause` NÃO muda
  por simetria: lá, vazio é erro de política e a borda (:490) o exige não-vazio — a assimetria é
  a semântica certa e vai comentada no helper.
- Espelho, catálogo vivo e contador N/M usam a MESMA definição de disponível — divergência = VC-2 falha.
- Lista de CODLOCAL é constante comentada no código; nada de UI de escolha de local (YAGNI).

**O número 3.822 do pack original está MORTO** — foi medido sem o recorte de CODLOCAL. Até a medição
viva voltar do especialista Oracle, o VC-2 se descarrega por concordância (tela == SQL da mesma regra
no mesmo banco), nunca por constante.

### A7 fechada — db-consult respondido e medido ao vivo (METALPRD 2026-07-29, READ-ONLY)

Predicado RATIFICADO, autoridade final do chip:

```sql
p.ATIVO='S' AND p.CODPROD>0 AND p.USOPROD='R'
AND EXISTS (SELECT 1 FROM METALPRD.TGFEST e
  WHERE e.CODPROD=p.CODPROD AND e.CODPARC=0
    AND e.CODEMP IN (1,2) AND e.CODLOCAL IN (10101,10102)
  GROUP BY e.CODPROD HAVING SUM(NVL(e.ESTOQUE,0)-NVL(e.RESERVADO,0))>0)
```

- **Subtrair RESERVADO: SIM.** 871 linhas reservadas / 37.428 un no escopo; **354 produtos zeram só por
  reserva**. "Tem estoque pra vender" = disponível.
- **`CODPARC = 0`: correto e defensivo.** No recorte é o único CODPARC com estoque (4.091 linhas,
  109.494 un); em outras empresas/locais existe estoque de terceiro, então o predicado fica.
- **Número de aceitação: 2.923 produtos vendáveis.** Sem o recorte de CODLOCAL seriam 3.508 — os
  **585 de diferença** entram indevidos hoje (show room, almoxarifado, quebras…), que é exatamente o
  defeito do sync. Decomposição: 10101 = 2.795 · 10102 = 187 (soma > total por sobreposição).
- **Conciliação com o número morto**: 3.822 usava ESTOQUE bruto, sem reserva e sem local. Mesmo escopo
  bruto com local = 3.277; a reserva tira 354 → 2.923.

**CODLOCAL com disponível > 0 nas empresas 1|2 (para o comentário-constante não mentir):**

| CODLOCAL | TGFLOC.DESCRLOCAL | produtos | vende |
|---|---|---|---|
| 10101 | 1_REVENDA | 8.310 | SIM |
| 10102 | 2_OUTLET | 595 | SIM |
| 10108 | 8_SHOW ROMM [sic no ERP] | 1.574 | não — mostruário |
| 10106 | 6_PENDENCIA PRODUTO FORNECEDOR | 673 | não |
| 10107 | 7_QUEBRAS EM PALETES | 553 | não |
| 10105 | 5_ALMOXARIFADO | 273 | não |
| 10103 | 3_AGUARDANDO CONSERTO | 113 | não |
| 10109 | 9_USO E CONSUMO | 70 | não |

**Whitelist, nunca blacklist** — local interno novo entraria vendendo sozinho se a regra fosse por
exclusão. O código carrega a whitelist como constante comentada.

### DR-3 — `AD_ECOMMERCE` é TRI-ESTADO (medido 2026-07-29, db-consult)

`VARCHAR2(10)`, nullable, **sem default e sem CHECK** — o domínio S/N/NULL é convenção, não garantia.
Distribuição em ATIVO='S': NULL 6.939 · `'N'` 2.993 · `'S'` 606. Dentro dos 2.923 vendáveis:
`'N'` 1.595 · NULL 887 · `'S'` 442.

`= 'S'` estrito derrubaria **85% do sortimento vendável** (2.923 → 442). Os três estados carregam
significados diferentes — publicado / recusado / **não-decidido** — e o operador já disse que a flag
"por agora não é confiável mas vai ser": o NULL é o "vai ser".

**Cláusula ratificada, live e espelho iguais (só o negativo explícito corta):**

```sql
AND NVL(AD_ECOMMERCE, 'X') <> 'N'      -- live
AND (m.ad_ecommerce IS NULL OR m.ad_ecommerce <> 'N')   -- espelho
```

Supera a forma `IS NULL OR = 'S'` do ruling inicial: com o domínio não garantido pelo banco, valor
novo que aparecer amanhã não é afirmação de "fora", e honest-unknown manda passar. Espelho guarda os
três estados como texto — **nunca colapsar em booleano**. Resultado com a cláusula ligada: 442 + 887
= ~1.329.

**A opção passa a chamar-se `only_ecommerce_eligible`** ("Somente elegíveis ao e-commerce"): o nome
`only_ecommerce` afirmaria "só os publicados", que é o que a cláusula deliberadamente NÃO faz.

## Defeito incluído no escopo (achado na investigação)

O Q4 de estoque do sync soma TODAS as empresas — `WHERE CODPARC = 0` sem CODEMP em
[sync.go:252-256](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go).
`estoque_total` do espelho está inflado com 901/501/904. **Pin `AND CODEMP IN (1, 2)`.**
(O catálogo live já usa `est.CODEMP IN (1, 2)` em
[catalog_page.go:135](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go) — o pin torna espelho e catálogo consistentes.)

## Decisões de desenho FIXADAS (não rediscutir)

1. **Config = 3 colunas boolean na infraestrutura `tenant_config` existente** (mesmo módulo e
   handler do `active_source` — [active_source.go](../../../apps/server_core/internal/modules/tenant_config/active_source.go)). Sem tabela nova, sem JSONB, sem módulo novo.
2. **Espelho ganha 2 colunas aditivas**: `usoprod`, `ad_ecommerce` (`products_mirror`,
   nullable). EAN já existe. Migração no bloco **0083–0084** (pré-alocado; precisar de mais =
   `REQUEST migration-number`).
3. Sync Sankhya Q1 ([sync.go:101-110](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go)) passa a ler `p.USOPROD, p.AD_ECOMMERCE`. Adapter xlsx: colunas
   OPCIONAIS no parser (planilha sem elas continua importando; ficam NULL).
   **S5B (ruling 2026-07-29):** a cadeia xlsx estava partida no PRIMEIRO salto — o parser lia
   `USOPROD` da folha e o valor morria no `CopyFrom` (`erp_import_products` sem as colunas; VC-6
   "xlsx COM as colunas as popula" já exigia o round-trip, o plano nunca deu dono). **Migração
   0085 concedida** (bloco vira 0083–0085): colunas em `erp_import_products` + CopyFrom + SELECT
   do sync + must-fail do round-trip folha→espelho.
   **Case (ruling): a regra é CASE-INSENSITIVE, canonicalizada na ESCRITA** — `TrimSpace+ToUpper`
   no parser xlsx E no writer Sankhya (S5B, um dono para as duas fontes); predicado live Oracle
   dobra na query (`UPPER(TRIM(...))`, S6 — sem escrita nossa para canonicalizar, e sem a dobra o
   live divergiria do espelho, proibido pelo DR-1); predicado do espelho fica SIMPLES, sem fold
   (S7 — armazenamento já canónico). Canonicalizar caixa NÃO colapsa o tri-estado: valor fora do
   domínio segue passando por honest-unknown.
4. Filtro aplica em TRÊS caminhos de leitura (**DR-1**, emendado 2026-07-29 — o pack original
   listava dois e esquecia o catálogo servido pelo espelho):
   - catálogo live Sankhya: condições na query de página conforme a regra do tenant;
   - catálogo servido pelo espelho (`routing.Reader` → leitor de upload → `MirrorCatalogPage`):
     MESMO predicado, MESMA contagem viva. A regra não pode depender da fonte ativa;
   - geração de vínculos: `MirrorMatcher` ([root.go:490](../../../apps/server_core/internal/composition/root.go)) filtra `products_mirror` pela regra.
5. Contador N/M: endpoint leve computado da fonte que a tela usa (não constante, não cache
   manual). M = população sem regra; N = com regra.
6. FE: card em [IntegracoesPage.tsx](../../../apps/web/src/pages/integracoes/IntegracoesPage.tsx); catálogo em `apps/web/src/pages/` (rota `/catalogo`); SDK client no bloco do domínio.

## Fora de escopo (YAGNI — corte do operador)

- Corte por `AD_DTULTVEND` ("sem venda há N meses") — futuro, colunas não entram agora.
- UI de escolha de empresa/armazém — `CODEMP IN (1,2)` é constante comentada no código.
- Qualquer coisa em análise de mercado / MIS-008.
- `TIPCONTEST` no espelho.

## Números de aceitação (base METALPRD 2026-07-29; tolerância = drift diário ~5 produtos)

- ATIVO=S: 10.538 · ∧R: 10.007.
- ∧estoque>0 em CODEMP (1,2) **sem recorte de CODLOCAL**: 3.822 (com EAN 3.154 = 82,5%) —
  número HISTÓRICO, superado pela EMENDA A7. Não usar como constante de aceitação.
- ∧disponível>0 em CODEMP (1,2) **∧ CODLOCAL IN (10101,10102)**, disponível = ESTOQUE − RESERVADO:
  **2.923** ← número de aceitação vivo do VC-2 (tolerância = drift diário).

## EMENDA 2026-07-29 — toggles do tenant não chegam ao catálogo: lacuna de PLANO; dono = S9 (ESCALATION do chip, ruling do hub)

**Facto (verificado pelo hub no objeto @8303f7c6):** o caminho de catálogo — Oracle
(`catalog_page.go:217`, `catalogAssortmentPredicate(options.IncludeAll)`; e
`buildCatalogAssortmentCountQuery`, que hoje não recebe opção NENHUMA) e espelho (S7,
`defaultSellableAssortment()` hardcode all-true em `erp_import/adapters/internalread/reader.go`)
— nunca lê a linha do tenant. O único sítio que pina a política numa leitura é
`routing/matcher.go:45-48` (caminho de linking). Contra VC-2/VC-3: virar toggle muda a geração
de vínculo e não muda a tela. Nenhum card atribuía "ler `tenant_config` e pinar a política no
caminho de leitura do catálogo". A lacuna atravessa S6 (aceito), S7 (em revisão) e S9 (fiação).

**Ruling 1 — dono é o S9.** A costura é exatamente a que o S9 já possui (root.go, service.go,
routing/reader.go): resolver a política do tenant UMA vez na seam de routing — onde
`matcher.go` já a resolve para linking — e entregá-la aos leitores; um produtor, N
consumidores. Fatia corretiva separada escreveria na mesma assinatura de port que o S9 vai
reescrever na fiação: colisão de costura, um dono só. **Extensão explícita de write-set do
S9:** `internal_read/ports/catalog_page.go`, `internal_read/adapters/oracle/catalog_page.go`
(predicado E a query de contagem), o leitor de catálogo do espelho
(`erp_import/adapters/internalread/reader.go`), além dos arquivos já concedidos.

**Ruling 2 — `IncludeAll` morre do port.** Bool separado + política são dois mecanismos que
têm de concordar — a classe do F-1. O port passa a receber o VALOR da política; "ver todos" e
`CatalogProductFactsByIDs` passam a política all-inclusive (três toggles desligados ⇒ o
predicado não emite nada) via construtor nomeado do domínio — caller nunca monta zero-value à
mão. Página e contagem recebem o MESMO valor. O default (linha do tenant ausente) vive na seam
de load do `tenant_config` — um lugar, com os defaults já ratificados do pack — e o leitor
nunca decide default; `defaultSellableAssortment()` morre. A forma do plumbing (pin de contexto
como o linking, ou campo de options) é escolha do S9, defendida na evidência — mas UM mecanismo
só, resolvido na seam de routing.

**Ruling 3 — lado Oracle entra no MESMO remendo (S9).** Mesmo defeito, mesma assinatura de
port, um dono. S6 não reabre; nada de fatia corretiva. Mirror-first não desculpa shim: VC-2/
VC-3 valem AGORA e o routing pode servir Oracle vivo — o predicado Oracle passa a ser
construído da política (pin por asserção de texto da query na lane unit, forma do S6, já que
Oracle vivo não entra na lane).

**Must-fail do S9 (grau de contrato; soma ao teste de chegada do catalog-503):** pelo leitor
COMPOSTO do root.go, virar um toggle na linha real de `tenant_config` muda página E contagem
juntas; reverter o threading (hardcode all-true) faz o teste falhar NOMEANDO o valor. Os dois
lados cobertos: espelho na lane de integração; Oracle por texto de query.

**S7: ACEITO como está.** Critérios da fatia provados (quatro mutações com morte nomeando
valor, SKIP 0 com env carregado, canonicalização segurando na linha suja, lado produtor da
condição 1). A lacuna é de plano e agora tem dono. Residual do DISTINCT na contagem
('007' vs '7') fica registrado, não-bloqueante.

**Correção do hub em registo:** no P4 do S6 o hub verificou que `IncludeAll` CHEGAVA ao SQL e
parou aí — meia-doutrina. Chegada tem duas metades: o valor chega ao consumidor E vem do
PRODUTOR certo (a linha do tenant). Gate de fatia que plumba opção/config pergunta as duas.

## EMENDA A-19 2026-07-29 — S8 reprovado no P4 do chip: ruling do hub nos três pontos

**Contexto:** worker do S8 cego na lane do sdk-runtime (sandbox nega travessia de diretório ao
esbuild; vitest abortou ANTES da coleta) e declarou o próprio zero — comportamento correto sob
o VC-7. No assento do chip: 2 vermelhos em guards de OUTRAS fatias. Probe do chip (trocar só as
duas refs `ErrorResponse` do bloco novo, medir, reverter) separou as causas: erpImport = defeito
de contrato REAL; activeSource = colisão de janela posicional.

**Ruling A — superfície de erro (recomendação medida do chip APROVADA):**
- Schema plano novo `SellableAssortmentError` (`{error, detail?}`), espelho do `ActiveSourceError`
  — a forma que `tenant_config/transport/http_handler.go` de facto emite.
- GET declara 200+500; PUT declara 200+400+500. O `404` SE DELETA (R-25: falsidade em contrato
  publicado se deleta — sob A-17 o default resolve na seam de load, GET nunca 404).
- Enum mínimo `[invalid_body, internal_error]`. Se a S10 medir necessidade de código novo, o
  contrato emenda por medição na hora — nada especulativo agora (YAGNI).
- O handler da S10 passa a dever exatamente esta superfície; critério entra no card da S10.

**Ruling A2 — janelas posicionais (tratamento SIMÉTRICO):** os DOIS guards que fatiam
`slice(âncora, indexOf("\ncomponents:"))` re-apontam a janela POR VALOR (terminar na PRÓXIMA
âncora de caminho), asserções intactas, movimento declarado em voz alta na evidência. Vale
também para o `fiveHundreds toHaveLength(4)` do erpImport quando o conserto adicionar o 500
novo: a contagem se protege por janela certa, nunca por inflar o número. Guard NOVO nascido
neste conserto já nasce com janela por valor. Fragilidade posicional geral fica como residual
escrito (fatias futuras que apendarem caminhos).

**Ruling B — quem executa: opção (i) GRANTED.** O chip aplica como glue de orquestrador —
grant EXPLÍCITO de estouro do teto de ≤10 linhas (~25 linhas YAML + re-apontamentos de teste).
Condições: commit PRÓPRIO em cima de `20dbc321` (o commit do worker fica como está); evidência
com contagem por linha antes/depois da suíte inteira e o probe registrado; linha de ledger
nomeando o grant. O vermelho que NOMEIA já existe (o guard do erpImport é o must-fail deste
conserto). S9 despacha depois da lane verde — o `include_all` do wire resolve nesta superfície.

**Ruling C — cegueira de sandbox é FINDING de harness** (arquivo próprio na missão,
`FINDING-sandbox-blind-fe-lane.md`): worker despachado não consegue rodar a lane do sdk-runtime
que o VC-7 exige; o orquestrador é a única testemunha de execução. Vinculado CHIP-LOCALMENTE
já: brief de fatia que toca `packages/*` DECLARA a cegueira de antemão e nomeia o assento que
mede (chip roda a lane a cada entrega); zero-observed do worker = BLIND, nunca verde.
Candidatura a profile aguarda ratificação do operador.

## EMENDA A-21 2026-07-29 — S9 reprovada no P4 do chip: F-1 (constante-default falsa) + F-2 (valor mágico no port); ruling do hub

**Contexto (verificado pelo hub em `routing/reader.go` @5adeeb56):** a metade-do-produtor está
CERTA (routing resolve a linha do tenant na mesma seam do matcher; página e contagem da mesma
chamada). Mas: **F-1** — `defaultSellableAssortment()` não morreu; foi exportada como
`ports.DefaultSellableAssortment()` (`{true,true,true}` fixo) e multiplicada para 11 sítios (7
produção). O nome é falso por medição: `tenant_config/repository.go:31` é fail-closed — tenant
sem linha recebe ERRO, nunca default; constante `Default...` que não é política de tenant
nenhum é classe R-24/R-25 (deleta-se). **F-2** — `resolveCatalogAssortment` compara
`requested == AllProductsAssortment()` (o ZERO VALUE): política explícita não-zero passada
pelo chamador é DESCARTADA em silêncio e a do tenant entra no lugar. O booleano voltou vestido
de struct — segunda instância da classe "dois mecanismos que têm de concordar" no chip
(primeira: `IncludeAll` ao lado da política, A-17/A-18). Regra 1 do A-20 invocada pelo chip:
reprova, não aceite-com-nota. Raiz única: o port não tem jeito honesto de dizer "usa a
política armazenada do tenant".

**Ruling — desenho ratificado (recomendação medida do chip):** o port passa a receber
`requested *SellableAssortmentPolicy` — `nil` = política armazenada do tenant; não-nil =
EXATAMENTE esta, honrada. Transporte mapeia na sua seam: `include_all` ausente/false → `nil`;
`true` → `&AllProductsAssortment()`. `AllProductsAssortment()` sobrevive como construtor
nomeado do VALOR honesto "sem corte" (é política real, não sentinela).
`DefaultSellableAssortment()` DELETA — zero referências após o conserto, asserido na evidência
por grep contado.

**Invariante novo (por-construção):** o `nil` resolve UMA vez, na seam de routing — o único
produtor que o A-17 estabeleceu. Abaixo do routing só viaja política CONCRETA; adapter/sítio
que precisaria "decidir default" é erro de programação, não fallback. É isso que torna a
constante desnecessária nos 7 sítios de produção.

**Execução — opção (i) GRANTED:** o chip aplica como conserto de orquestrador sob grant
explícito (escopo maior que glue: assinatura do port + transporte + deleção nos 11 sítios),
commit PRÓPRIO sobre `5adeeb56`, worker commit intacto. Razão contra (ii): re-injetar arrisca
regressão num diff de 13 arquivos que já passa a contraprova, e o assento que mede é o chip.
(iii) rejeitada: S11/S12 são exatamente os chamadores expostos à armadilha — dívida aqui
venceria antes de ser paga.

**Must-fails do conserto (grau de contrato):**
1. M1 re-executada pós-conserto (chegada da política armazenada pelo leitor composto, nas
   DUAS lanes, com env — o `ok` sem env já provou mentir sob mutação);
2. NOVA: política explícita não-nil passada pelo chamador é HONRADA pela cadeia composta
   (mata F-2 para sempre; é a proteção de S11/S12);
3. grep contado `DefaultSellableAssortment` = 0;
4. contagens por linha antes/depois (A-15).

**Disposição A-20 para F-2:** classe "valor mágico/dois mecanismos", segunda ocorrência —
ramo (a) inline: o desenho por ponteiro remove o MECANISMO (não há mais valor válido dobrando
de sentinela), não a instância.

## EMENDA A-22 2026-07-29 — alegação falsa do A-17 corrigida (não há default de sortimento); tenant sem linha = 400 `unknown_erp_source`; grant de seam contrato+SDK ao S10

**Correção de registro (a falsidade era do HUB):** o A-17 escreveu "o default (linha do tenant
ausente) vive na seam de load do `tenant_config`". Refutado por medição do chip contra
`595c15e3`: `tenant_config/repository.go:29-51` (Get → `ErrUnknownActiveSource` sem linha,
"fail-closed — never a silent default") e `:91-104` (SetSellableAssortment = UPDATE, nunca
cria linha). **Este sistema não tem default de sortimento — tem dois estados: configurado e
não-configurado, fail-closed.** A frase do A-17 fica corrigida aqui; a frase do card do S10
que a herdou se corrige no plano. Nota: o 404 deletado no A-19 continua deletado — a razão
correta nunca foi "default existe", é a que segue.

**Ruling — opção (b) do chip APROVADA:** tenant sem linha responde **400 `unknown_erp_source`**
em GET e PUT de `/config/sellable-assortment` — a MESMA resposta que o vizinho
`/config/active-source` do mesmo arquivo já dá para o MESMO estado. Um estado, um código, um
módulo. (a) 500 rejeitada: "erro interno" para tenant novo é falso na tela — o estado é normal
e a S11 renderiza ao lado do seletor de fonte. (c) 200-com-defaults PROIBIDA por nome: é o
`DefaultSellableAssortment` do A-21 voltando pelo transporte.

**Grant de seam (contrato+SDK, dono S8 → estendido ao S10, additive-only, mesmo commit do
handler):** GET ganha `"400"` com `SellableAssortmentError`; enum ganha `unknown_erp_source`
(vale para GET e PUT — `RowsAffected()==0` no PUT é o mesmo estado). Guards do sdk-runtime
re-apontam POR VALOR com movimento declarado (janelas já são posição-independentes pela S8).
A válvula do card ("código além dos dois só por medição") foi usada como desenhada: o chip
MEDIU, pediu, o hub concedeu — a emenda nasce com a medição citada.

**Must-fail:** GET sem linha → 400 corpo plano `unknown_erp_source`; PUT sem linha → idem;
mutação que devolva o interino 500 (ou invente default) morre nomeando o valor. S11 recebe no
brief: 400 `unknown_erp_source` nesta tela = estado "configure a fonte primeiro" (igual ao
card vizinho), nunca tela de erro.

## EMENDA A-23 2026-07-29 — S11: copy dos toggles `only_revenda`/`only_em_estoque` APROVADA COMO PROVISÓRIA; decisões de brief endossadas

**Ruling de copy (R-24: string de tela é superfície do OPERADOR — ratificação final é dele):**
o hub aprova a copy proposta pelo chip como PROVISÓRIA e ela pode ir a commit; o hub leva as
strings ao operador no report do turno e uma discordância vira edit de uma linha, nunca rodada.

- `only_revenda` → rótulo "Somente produtos de revenda"; apoio "Mantém no sortimento apenas o
  que o ERP classifica como revenda."
- `only_em_estoque` → rótulo "Somente com estoque disponível"; apoio "Considera o saldo
  disponível no CD e na loja física. Com a regra desligada, produtos sem saldo continuam no
  sortimento com o aviso Sem estoque."
- `only_ecommerce_eligible` já era ratificado verbatim (BATCH-PLAN.md:1389) — inalterado.

As duas razões do chip ficam ratificadas como FORMA para copy futura deste módulo: (1) rótulo
não cita coluna de ERP (segue a forma do rótulo já ratificado); (2) apoio descreve o efeito do
estado que produz a dúvida real na tela — para o toggle de estoque, o estado DESLIGADO, porque
é onde o badge "Sem estoque" coexiste com o produto no sortimento e a copy tem que impedir a
leitura "badge = segundo filtro".

**Decisões de brief ENDOSSADAS (nenhuma é decisão retida do hub):** três checkboxes
independentes; write sempre com a tripla completa (espelha o 400 `invalid_body` do servidor);
sem otimismo, fieldset desabilitado durante write (paridade com card de fonte ativa);
invalidação ALVO no namespace `catalog` com a razão registrada (sortimento não toca
pedidos/anúncios/instalações — o `invalidateQueries()` cego do active-source NÃO é o padrão a
copiar); chave de cache de counts carrega `erp_source` (classe já paga no CHIP-IMPORT-FIX:
escopo resolvido no servidor TEM que entrar na chave do cache downstream); counts com falha
nunca renderizam `0 de 0` (ADR-17).

**Re-entrega (D-1, 3º caso):** as duas REQUESTs que o chip lista como abertas JÁ FORAM
decididas e commitadas ANTES desta mensagem do chip: linha ausente = 400 `unknown_erp_source`
com grant de seam ao S10 (A-22 @d514e0e1); `removal_owner: "CHIP-ERROR-UNIFY"` na
temporary_exception (ruling S10 @18a9cebb-adjacente, ACK enviado). Drene a fila de entrada
antes de compor o próximo evento.

## EMENDA A-24 2026-07-29 — varredura S9: caminho sem corte morto POR COMPOSIÇÃO, não por asserção; grant de teste de composição; deref pré-existente registrado

**Resultado da varredura aceito como medido:** nenhum handler/serviço de produção alcança
página sem corte fora do routing. (a) rotas legadas só registram sob
`PageReader==nil && !hasRouteClasses`, e o mux de produção (`root.go:280`,
`*httpx.RouteClassMux`) fixa `hasRouteClasses=true` — ramo morto por condição de composição;
(b) pricing segura `catalogSvc` e não chama os métodos base (grep de módulo); (c)
`CatalogProductFactsByIDs` sem corte por desenho (BATCH-PLAN:655-662); (d) demais chamadores
são a própria cadeia de decorators. S11/S12 podem fiar tela.

**Ruling — a garantia não pode repousar em ramo negativo não-assertado. GRANT (i):** o chip
autora, como conserto de orquestrador em commit próprio (mesma moldura do grant A-19), um
TESTE DE COMPOSIÇÃO que monta a forma de produção do mux e afirma que as rotas legadas NÃO
existem. Classe nomeada (stop-the-line C-1): "a composição de produção decide o
comportamento e nenhum teste MONTA a composição de produção" — 2ª ocorrência da família
(1ª: M02, decorator apagou porta opcional; remédio de lá foi compile-assert por porta). O
teste é a forma geral para condição de REGISTRO (não expressável em compile-time).
Must-fail: (m1) registrar o handler legado sob o mux de produção → o teste morre NOMEANDO a
rota; (m2) inverter a condição `hasRouteClasses` → idem. Write-set: arquivo de teste novo no
pacote da composição — disjunto de S10-COND e S11 por construção; item 9 do checklist do chip
cobre.

**Achado pré-existente REGISTRADO (não introduzido pelo chip, não tocar nesta chip):**
`internalReadAvailable==false` deixa `PageReader` nil com caminho de deref no handler de
catálogo — produção degradada derrubaria o handler. Mesma classe (ramo só exercitado em
degradação). Vai ao operador como backlog nomeado; candidato natural: chip de hardening
pós-CHIP-ERROR-UNIFY, ou absorção pelo próprio ERROR-UNIFY se o remédio for o padrão único de
erro no caminho degradado. Decisão do operador, não desta chip.

**ACKs registrados:** A-22 em execução numa peça só com o TERCEIRO assento (re-leitura
pós-PUT) — correto e melhor que o pedido, mesmo estado, mesma mentira do 500; retratação da
falsidade herdada em DOIS sítios do BATCH-PLAN (:657 e :726) com medição ao lado, forma
certa (retratar, não apagar — o leitor futuro vê por que a frase morreu); governança com
ordem de chaves casada aos precedentes e stop-and-report para chave não especificada; D-1
com conserto de processo (drenagem no INÍCIO do passo — raiz nomeada: compor no fim de
revisão longa é o pior momento da caixa); S11 rodada 2 carrega a semântica A-22 com o
checklist P4 anotado com a data da inversão para o gate não ler critério velho.

---

## A-25 (2026-07-29) — R-1 endpoint de teste + R-2 `.js` emitidos em `src/`

### R-1 — RULING: pg-session do PRÓPRIO worktree; endpoint morto não se re-aponta, se RE-CRIA

Medido pelo hub: `scripts/harness.ps1` expõe `pg-session-up`/`pg-session-down`;
`Get-HarnessPostgresSessionContainerName` deriva o nome do contêiner do HASH do repo root —
desenho explícito no código: "hub and chip worktrees each own their session container".
`Start-HarnessPostgresSession` detecta estado stale/contêiner morto e recria do zero
(estado em `scripts/.runs/pg-session.json`). 55864 era porta efêmera de sessão morta.

Autorizado ao chip (não é violação do "chip nunca sobe servidor" — essa regra é do DEV STACK;
a sessão de pg de teste é por-checkout POR DESENHO):
1. `pwsh scripts/harness.ps1 pg-session-up` no worktree do chip (emite `container=` e `port=`).
2. Reconstruir a linha de env do A-15 a partir de `port` + password do `pg-session.json`
   (password NUNCA impresso em evidência — LEN apenas, prática corrente).
3. `go run ./cmd/migrate` de `apps/server_core` contra a base da sessão — `applied N` com N
   conferido contra as migrações da árvore (regra @f1cba2a, metade 2).
4. Re-rodar o pacote `tenant_config` com SKIP=0 contado por linha.
5. `pg-session-down` no fechamento do chip.
NÃO apontar para 5435: é o postgres do dev-stack (dados vivos de dev, seam do hub).
Avisos herdados: `.gomodcache` frio → `HPG_MIGRATION_FAILED` falso (aquecer antes);
primeiro boot → retry do CREATE DATABASE (pg_isready mente).

ACHADO RATIFICADO (classe, entra no HARNESS-DEBTS como B-5): **"variável setada" ≠ "banco
alcançável"** — a linha do A-15 prova o env, não o endpoint. Brief que exige SKIP=0 exige
mecanismo de endpoint VIVO (lane que boota a própria base, ou probe pré-lane com token
distinto), senão só troca pulo silencioso por vermelho ambiental e queima rodada.

Registrado: asserção do sentinela em `repository_test.go:189` NUNCA executou nesta chip —
lacuna declarada, não cobertura. O passo 4 acima a executa; o pacote de evidência do S10-COND
só fecha com ela rodada.

### R-2 — GRANT (a) + retenção (b)

(a) AUTORIZADO: apagar, no worktree do chip, os `.js` (e `.js.map` se houver) UNTRACKED
(`git status --porcelain`, linhas `??`) que tenham irmão `.ts`/`.tsx` de mesmo basename no
mesmo diretório — inclui `apps/web/vite.config.js`. Verificação por arquivo (irmão existe),
contagem reportada (esperado ~186). Tracked não se toca; `git add -A` segue proibido.

(b) RETIDO NO HUB: regra de `.gitignore` e conserto do `tsconfig` (emissão dentro de `src/`)
são seam compartilhado — hub decide PÓS-MERGE (backlog registrado abaixo). Chip não toca
`tsconfig`/`.gitignore`. Root-cause aceito como reportado: `tsc -b` da revisão S11 do próprio
chip emitiu; classe já nomeada no CHIP-IMPORT-CHAIN ("observável que passa nos dois mundos
não é evidência") — agora com árvore fantasma EM DISCO com nome de fonte.

BACKLOG HUB (pós-merge): (i) decidir emissão do `tsc -b` (noEmit na lane de tipos vs outDir
fora de `src/` — atenção: projetos composite exigem emit p/ incremental; medir antes);
(ii) `.gitignore` só DEPOIS do (i) para não mascarar a doença.

### S14 — forma endossada

Teste invisível na varredura `-run TestRootRuntime` (4 linhas `---`, a dele ausente) e visível
nomeado exato: NÃO fechar S14 com a varredura em que ele não apareceu — correto. A varredura
se valida no chamador (regra CHIP-FIM); resultado da re-run entra na evidência do S14.

---

## A-26 (2026-07-31) — decisões da CLOSURE-SPEC (@0fe085c) — pontos 1–3 + ratificação de findings

Spec lida na íntegra pelo hub via `git show 0fe085c`. Endossada como está, incluindo as três
correções de lane (pacotes sem lane rodados explícitos; `-BaseSha` = tip da MAIN re-lido na
execução, nunca base de despacho; forma workspace do vitest).

### Ponto 1 — regeneração de candidatos: QA PRIMEIRO com o banco velho; regeneração DEPOIS, do hub

Regra já ratificada na missão (QA final M-06): **o baseline do QA é o banco velho — não se
limpa estado antes do QA**; o defeito aparecendo na tela é o before/after real. Aplicação:
1. QA de browser roda COM as linhas antigas. Ressalva por escrito no pacote: motivos gravados
   antes do F2 exibem `(unproved)`; motivos produzidos depois exibem a copy nova. VC-7 avalia
   o PRODUTOR (linhas novas limpas = PASS); linha velha na tela não reprova o chip — ela é a
   evidência "before".
2. Regeneração de candidatos = passo do HUB, pós-merge e pós-QA. Escopo: candidatos PENDENTES
   apenas; vínculos aprovados (os 31 do operador + auto-aprovados D-121) intocados. Vira o
   "after" documentado. Operador pode vetar no relatório de QA.

### Ponto 2 — F3 CONFIRMADO, chip executa sob grant

`N≈3.822` medido sem corte de CODLOCAL e sem reserva = prosa falsa em contrato, e prosa falsa
em contrato SE DELETA (R-25, CHIP-ANCHORS). Chip deleta a linha e aponta para
`ADDENDUM-01-codlocal.md:53` + A-3 como fonte, PRESERVANDO a restrição do BATCH-PLAN:1188
(2.923 = Oracle live apenas; espelho tem número próprio). Único toque em texto de critério;
grant registrado no ledger do chip.

### Ponto 3 — fronteira mantida; fila formal no CLOSED

Scrub de PII (`docs/design/evidence/ml-api/`, seam do hub) e unificação da superfície de erro
HTTP (chip próprio, já ratificado como próximo da fila) ficam FORA. "Finalização final" NÃO os
inclui — escopo do chip não muda, nenhum grant novo. Os dois entram FORMALMENTE no payload do
CLOSED como fila pós-merge. Se o operador quiser dentro, é ordem nova ao hub, não alargamento
silencioso.

### Findings ratificados (entram no CLOSED para o profile; 1 e 2 já viram dívida)

1. Descoberta de pacote da lane de integração (5 primeiras linhas, só `internal/modules/**`)
   deixa `tenant_config` e `internal/composition` FORA de toda lane → **B-6 no
   HARNESS-DEBTS.md**. Caso pago: teste de política armazenada movendo página+contagem pulou o
   chip inteiro; rodado vivo: composição RUN=72 PASS=72 SKIP=0.
2. `vitest.config.ts` inclui `feature-products` por NOME EXATO, não glob → segundo arquivo de
   teste roda em lane nenhuma, zero silencioso → **B-7 no HARNESS-DEBTS.md**.
3. Critério vácuo contra o TIPO ("nunca chama X" com X ausente do tipo do colaborador) —
   classe já nomeada no S12; ratificada para o profile: todo critério "nunca chama X" abre o
   tipo e confirma que X existe antes do despacho.

---

## A-27 (2026-07-31) — BLOCKED shell do chip + ratificação do alargamento F4/F4b/F5/F6

### Shell

Wedge é POR SESSÃO: shell do hub probado vivo (`true` + `git --version` OK) no mesmo minuto do
BLOCKED. Conserto = restart da sessão do chip pelo operador (wrapper composto em runtime, sem
arquivo em disco para editar — diagnóstico do chip aceito). Escalado ao operador. Dívida **B-8**
registrada no HARNESS-DEBTS.md (shell morre para a sessão inteira, subagentes incluídos, sem
degradação graciosa; candidato: probe de shell no bootstrap — converge com intervenção 2 da
análise global, atestação no boot).

### Alargamento RATIFICADO, com condições

Cobertura: F4 (dobra da porta `Source = Reader + CatalogPageReader`, `CatalogAssortmentReader`
deletado), F4b (decorator único no cache + asserção `timing.go`), F5 (deleção da rota legada
não-cortada: fork do `Register`, handlers legados, métodos em 2 services + 2 portas + 4 walks +
4 `UnavailableReader`), F6 (guard `q` obrigatório em `handleSearch` → 400 `invalid_q`).

F6 em especial: **conformidade, não mudança de contrato** — hub verificou
`contracts/api/marketplace-central.openapi.yaml:440-442` (`q` `required: true` na rota de
busca). Comportamento atual (q em branco → catálogo inteiro vestido de busca) é a violação.

Condições (nada aceito antes delas):
1. NADA do F4–F6 é aceito até a escada P5 INTEIRA verde do zero — 4 arquivos de teste migrados
   à mão nunca compilados; verde pré-wedge do F4 não vale para a árvore atual.
2. Must-fail do F6 obrigatório (mesma régua do F1): guard removido → teste vermelho NOMEANDO
   `invalid_q`. Escopo do guard = SÓ a rota onde o spec exige (`:442`); as rotas com `q`
   `required: false` (`:573`, `:671`) não ganham guard.
3. F6 evidencia o CHAMADOR FE: prova de que a tela de busca nunca dispara a rota com `q` vazio
   (senão o conserto quebra a tela e o QA reprova o chip inteiro).
4. Forma do erro 400 `invalid_q` = mesma família dos 400 existentes do módulo (`invalid_body`,
   `unknown_erp_source`) — não nasce quinta família; a unificação global segue com o
   CHIP-ERROR-UNIFY.
5. `catalog_routes_test.go`: cabeçalho descreve fork que o F5 deletou = prosa falsa sobre o
   repo (R-25, se DELETA) e a checagem `CATALOG_METHOD_NOT_ALLOWED` virou vácua (string só
   sobrevive em `handleTaxonomy`). Consertar DENTRO do F5 — deletar a prosa e re-apontar ou
   deletar a checagem com observável; contenção do chip (não empilhar mudança não verificável
   com shell morto) foi correta, mas o F5 não fecha sem isso.
6. Governança da escada: `-BaseSha` = tip da main RE-LIDO na execução (agora 774d75aa ou
   posterior), regra A-26 mantida.

### Estado congelado

`STATE-SHELL-BLOCKED.md` (untracked no worktree do chip) aceito como custódia do estado;
primeiro ato pós-restart = commitá-lo no branch do chip antes de retomar.

---

## A-28 (2026-07-31) — A-27-R2: âncoras reconciliadas; F6 alarga para declarar o contrato inteiro do handler

### R1-a — as duas leituras estavam certas, em árvores diferentes

Hub mediu na MAIN (`q` da busca em :440-442); chip mediu na árvore DELE (:467-469) — o OpenAPI
andou no branch. Facto idêntico: `/catalog/products/search` → `q` `required: true`; get-by-id
não tem `q`. Condição 3 RE-ANCORADA POR CONTEÚDO: o guard vale para a rota
`/catalog/products/search` (operationId `searchCatalogProductFacts`), única com `q`
`required: true`; rotas com `q` opcional (`/listings`, `/listings/by-product`, `/orders`) não
ganham guard. REGRA DE FORMA (classe, entra no A-2 do HARNESS-DEBTS): âncora de ruling
cross-árvore é CONTEÚDO (rota/operationId/schema name); linha só como cortesia COM a árvore
nomeada. Mesma lição da autoridade content-addressed (P-C da análise global).

### R1-b — condição 5 CORRIGIDA

Família certa da página de catálogo = `CatalogPageErrorResponse` (enum FECHADO:
`[invalid_cursor, invalid_limit, source_unavailable, deadline_exceeded]`). `invalid_body` e
`unknown_erp_source` são de mutations/config — citação errada do hub, retirada. O `invalid_q`
do chip já sai por `writeCatalogPageError` na forma dos irmãos → condição 5 SATISFEITA no
código; o débito é de DECLARAÇÃO.

### Pergunta 1 — SIM: F6 fecha OpenAPI + SDK no mesmo commit

Regra do repo (AGENTS.md): mudança de API atualiza OpenAPI e `sdk-runtime` juntos. F6 inclui:
`invalid_q` no enum de `CatalogPageErrorResponse` + descrição do 400 da rota de busca
corrigida (hoje "Invalid cursor or limit" — prosa falsa desde o S9; passa a nomear as causas
reais) + regen do SDK. Sem isso o F6 nasceria violando a regra que invoca.

### Pergunta 2 — os 3 pré-existentes ENTRAM no mesmo enum, sob o mesmo F6

`invalid_include_all`, `invalid_erp_source`, `invalid_ids` já são RESPONDIDOS pelo handler e
mudos no contrato = o contrato MENTE sobre o módulo (classe R-25 — só que aqui a falsidade se
conserta DECLARANDO, porque o comportamento é o ratificado; deletar seria quebrar S9/S10).
Não é a doença do CHIP-ERROR-UNIFY (nenhuma família sendo unificada; só declaração do que o
módulo já responde). Condição: cada código declarado prova o sítio EMISSOR (file:line na
árvore do chip) na evidência — nenhum código especulativo entra no enum.

### Pergunta 3 — reconhecido

Regen do SDK exige shell; F6 (e tudo) travado até o restart. Escalação mantida.

Condições 1, 2, 4, 6 aceitas pelo chip sem ressalva — mantidas.

---

## A-29 (2026-07-31) — A-28-R2: censo aceito; `allowed_range` = saída (a), serializar

### Censo de emissores — condição do A-28 CUMPRIDA

`evidence/F6-error-code-census.txt` aceito: 4 declarados + 3 pré-existentes + `invalid_q`,
todos com sítio emissor nomeado, todos pelo mesmo escritor `writeCatalogPageError`. Nenhum
especulativo. Commit do censo entra junto com o F6 quando o shell voltar.

### `allowed_range` — RULING: saída (a), serializar para todo código que o tenha

Achado aceito como classe já conhecida ("descrição verdadeira por acidente do writer, não por
intenção do código" — irmã do "Invalid cursor or limit" do S9, mais escondida). Decisão (a):

1. `writeCatalogPageError` serializa `allowed_range` sempre que o código o carrega — hoje:
   `invalid_limit`, `invalid_erp_source` (`xlsx|catalogo_cliente`), `invalid_ids`.
2. Descrição do schema reescrita para a regra REAL (presente quando o erro tem domínio aceito
   limitado, nomeando-o), não lista hardcoded de um código só.
3. MESMO commit do F6 (enum + descrições + regen SDK) — uma mudança de contrato, não três.
4. Must-fail com DOIS nomes: `invalid_q` e `allowed_range` ausente onde devido.

Razões: aditivo no fio (campo opcional declarado passa a aparecer; nenhum cliente quebra por
campo a mais que o schema já descreve); o dado existe porque alguém quis respondê-lo — (b)
consertaria a prosa apagando a intenção; e o CHIP-ERROR-UNIFY herda a superfície INTEIRA —
herdar campo mutilado é pior que herdar campo servido. Terceira e ÚLTIMA largura do F6: escopo
congela aqui; achado novo nessa superfície vira REPORT para o CHIP-ERROR-UNIFY, não conserto.

### Estado

Tudo segue travado no shell; restart escalado ao operador segue pendente. Contenção do chip
(zero edição não verificável) endossada de novo.

---

## A-30 (2026-07-31) — Takeover do hub executado; escada verde; Sol fora (quota) → waiver §12 re-invocado

### Takeover (ordem do operador)

Sessão do chip abandonada (shell wedge, B-8). Hub = único escritor do worktree
`chip-vendavel` desde então. Trabalho recuperado por custódia @c8f7d90e; F4/F4b/F5
@80f0319c; F6 @86fc6400 (contrato nos 3 lados no MESMO commit + must-fail duplo em
`evidence/F6-must-fail.txt`); cond.6 @83c91556; governança @a2f96fa1; EVIDENCE.md @47dc7dd2
(árvore congelada do gate).

### Escada P5 — verde (EVIDENCE.md no worktree)

migrate 72/0; no-lane (tenant_config+composition) PASS=115 SKIP=0 com banco da sessão;
go vet limpo; go test 107 pacotes ok; tsc 12 pré-existentes (0 em tocado); vitest 67+5+1
files verdes; integração status=passed (sessão `mpc-pg-session-3eee515d`:50265).

### Governança — lane VERMELHA no main (finding novo, B-9)

Baseline em worktree LIMPO no main tip `4ad36272`: status=failed, **51 violações
pré-existentes**. Chip = 50 = subconjunto estrito (remove GOV_MODULE_COVERAGE
tenant_config; adiciona zero). Critério de aceite aplicado: zero violação NOVA por diff de
conjunto (code/id/path). Verde absoluto é inalcançável p/ qualquer chip até o hub saldar a
lane. Registrado como B-9 no HARNESS-DEBTS; caso irmão B-10: lane no checkout do hub trava
>20min varrendo dumps untracked `docs/design/evidence/ml-api/` (filtro da Policy não exclui).

### Gate P6 — Sol MORTO (quota codex até Aug 5) → waiver §12

Assento Sol morreu: `You've hit your usage limit ... try again at Aug 5th, 2026`.
RULING: re-invocado o padrão ratificado do waiver 2026-07-18 (profile §12, contingency
block) — dual gate = cold Opus (`harness:gate-reviewer`) + assento sonnet adversarial
independente — sob a ordem vigente do operador ("use sonnet 5 subagent e orquestre para
finalizarmos rapido"). Opção de review Sol retroativo pós-Aug-5 mantida, como no waiver
original. Transporte: próximo dispatch codex usa stdin `-` (binding 2026-07-29), não argv.

## A-31 (2026-07-31) — Gate P6 fechado (2x APPROVE) e MERGE executado @3ef01f65

Round 0: Opus REJECT (1 BLOCKER + 2 MAJOR + 1 MINOR); sonnet adversarial APPROVE (2 MINOR).
BLOCKER real: a reescrita do writer no F6 introduziu um 9º código (`internal_error`, 500,
fall-through de `writeCatalogPageError` :429) alcançável na rede e não declarado — a mesma
classe que o F6 existia p/ matar. Consertos @a0bdc271: declarado nos 3 lados (enum OpenAPI +
resposta "500" em list/search/counts, counts ganhou "400"/"504" que já alcançava; SDK
`CatalogPageErrorCode`); censo atualizado (9 códigos, todos com emissor); âncoras do
FE-caller re-apontadas p/ os caminhos reais (feature-products, não apps/web/src/pages);
razão do modules.json corrigida ("43" morto → main 51 / chip 50 subconjunto estrito); stubs
mortos deletados dos 2 test files. Marker de supersessão no S10 @14c018c9.

Round 1: Opus re-verify focado — 5/5 DISCHARGED, zero scope creep, APPROVE. Acordo dos dois
assentos → MERGE `--no-ff` @3ef01f65. Conflitos (2, hub-owned): tabela de toggles resolvida
pelo CÓDIGO (nomes/predicados do chip: `only_ecommerce_eligible` NVL, + ponteiro A7 do main);
VC-2 resolvido pelo lado do chip (F3: número fora da linha do critério, nota preservada).
Adições do main confirmadas vivas pós-merge (A-30, B-9/B-10, EMENDA A7).

REPORT novo p/ CHIP-ERROR-UNIFY (achado Opus round 1, pré-existente): rota counts lê
`erp_source` (`http_handler.go:53`) mas não declara o parâmetro no OpenAPI.

### A-31 correção pós-merge — totais da governança

Lane pós-merge em worktree limpo @3ef01f65 (BaseSha main tip 4ad36272): status=failed,
**54 violações**; set-diff vs baseline = ONLY-BASELINE `GOV_MODULE_COVERAGE tenant_config`,
ONLY-POSTMERGE vazio — **zero violação nova**, critério de aceite mantido. Os totais "51/50"
registrados em A-30/A-31 eram soma errada do breakdown do próprio baseline (que sempre somou
55); corrigidos p/ **55/54** no EVIDENCE.md e no reason do modules.json. O set-diff — o
critério — nunca dependeu do total.

### A-32 QA live-drive PASS + encerramento (2026-07-31)

QA browser pós-merge (persona fresca, banco velho per A-26) — **PASS VC-1..VC-7**, medições
em `qa-validation-result.md`. Um defeito achado e consertado DURANTE o QA: counts 503 por
ORA-00937 (COUNT(*) + subquery escalar sem GROUP BY; latente desde 5adeeb56 — substring-pin
de teste não enxerga validade de statement, só o live drive). Fix @f9d6a91, after
`200 {"sellable_count":2945,"total_count":10549}`.

VC-5 LIVE-VERIFIED: primeiro tick do scheduler pós-rebuild (17:10 UTC) populou
usoprod=10549 / ad_ecommerce=3604 no espelho sankhya; predicado no espelho
(`usoprod='R' AND estoque_total>0`) = **2945**, idêntico ao Oracle live — o Q4 do sync
reproduz o corte vendável exatamente.

Encerramento: pg-session mpc-pg-session-3eee515d down; worktree chip-vendavel removido;
branch `worktree-chip-vendavel` deletado com `-d` (mergeado); sessão do chip notificada.
CHIP-VENDAVEL **CLOSED**. Fila herdada: CHIP-ERROR-UNIFY (+ gap `erp_source` no OpenAPI do
counts), regeneração de candidates PENDING, root package.json test script, PII scrub ml-api,
retro Sol pós-Aug-5.
