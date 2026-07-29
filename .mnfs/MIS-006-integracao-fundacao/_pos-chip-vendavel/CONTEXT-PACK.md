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
| `only_em_estoque` | estoque disponível > 0 (ver EMENDA A7 abaixo para o recorte) | **true** | ver A7 |
| `only_ecommerce` | `TGFPRO.AD_ECOMMERCE = 'S'` | **false** | flag rala hoje (606); cliente vai passar a manter |

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
