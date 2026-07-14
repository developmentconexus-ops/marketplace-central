# Descoberta pré-anúncio no Mercado Livre — batch real de 50 produtos

**Data:** 2026-07-12  
**Escopo:** produtos ativos do ERP sem anúncio ativo na conta  
**Método:** API oficial, somente `GET`, preço sem frete  
**Veredito atualizado em 2026-07-13:** viável parcialmente pela API oficial; GTIN sozinho não
autoriza match automático e nenhuma fonte externa está homologada para produção/OEM.

## 1. Pergunta testada

Para um produto que existe no ERP, mas ainda não está anunciado, conseguimos localizar o mesmo
produto no Mercado Livre e obter uma distribuição de preços útil para uma futura simulação de
rentabilidade?

Este é o funil correto do produto:

```text
ERP: CODPROD + EAN
  → /products/search por product_identifier
  → candidato de catálogo com o mesmo GTIN
  → validação semântica de domínio/marca/modelo/variante
  → /products/{catalog_product_id}: winner/faixa somente quando não nulos
  → /products/{catalog_product_id}/items apenas como evidência legacy, nunca contrato durável
  → fonte externa homologada para ofertas/preços quando a rota oficial não cobrir
  → futura simulação de custos e margem
```

`price_to_win` não participa da descoberta inicial; ele continua sendo uma capacidade posterior,
quando já existe um item do seller.

## 2. Seleção e integridade

- Oracle: `SELECT` somente.
- Reserva inicial: 60 produtos ativos, 60 `CODPROD` e 60 EANs únicos.
- EANs: formato e checksum GTIN válidos; unicidade verificada entre produtos ativos.
- Seleção: determinística, com `CODPROD 15956` excluído por ser o vínculo já comprovado no estudo anterior.
- Antes do batch, a API listou os anúncios ativos e excluiu qualquer coincidência por EAN.
- Resultado final: 50 produtos; `already_active=0` em 50/50.
- API Mercado Livre: somente `GET`; apenas o EAN foi transmitido.
- Nenhum token, segredo, URL de banco ou buyer PII foi persistido.

## 3. Resultado bruto

| Métrica | Resultado |
|---|---:|
| buscas por EAN | 50/50 HTTP 200 |
| candidatos com GTIN idêntico | 22/50 (44%) |
| sem candidato GTIN exato | 28/50 (56%) |
| candidato com ofertas e preços | 18/50 (36%) |
| candidato sem ofertas (`404`) | 4/50 (8%) |
| ofertas declaradas/materializadas | 132/132 |
| sellers somados por catálogo | 103 |
| candidatos de preço com ≥5 sellers | 11 |
| candidatos de preço com 3–4 sellers | 1 |
| candidatos de preço com ≤2 sellers | 6 |

Todos os 22 candidatos responderam `200` em `/products/{id}` e em
`/pricing-automation/products/{id}/rules`. A existência da regra `INT`/`INT_EXT` não fornece o
benchmark da tela; apenas informa quais regras estão disponíveis para o produto.

## 4. Adjudicação adversarial dos 22 candidatos

O algoritmo do probe exigiu GTIN exato, mas o resultado provou que isso é insuficiente.

### Contradições claras — rejeitar

| CODPROD | ERP | Catálogo retornado |
|---:|---|---|
| `45477` | toalheiro inox Doka | parafusadeira/furadeira Menegotti |
| `44812` | cabide Rainbow Gold Doka | lanterna traseira automotiva |

Esses dois casos têm EAN idêntico na resposta, mas domínio e natureza do produto são incompatíveis.
Isso pode indicar GTIN incorreto no ERP, associação incorreta no catálogo ou reutilização indevida do
código. A origem exata ainda é desconhecida; o efeito prático é confirmado: GTIN não pode vencer uma
contradição semântica.

### Ambíguos — revisão obrigatória

| CODPROD | Motivo |
|---:|---|
| `34602` | ERP indica acabamento Deca `GD` para 1¼–1½; catálogo indica `PQ`/até 1. |
| `30585` | ERP indica KDT 2113/440 L; título do catálogo informa 8800 W, mas não prova modelo/capacidade completos. |

### Triagem inicial por título — posteriormente superada

- A primeira triagem considerou 18/22 candidatos compatíveis pelo título e domínio.
- Depois, `/products/{id}` forneceu atributos estruturados para 22/22 e revelou seis casos que o
  título não permitia validar com segurança.
- Os números 18/50, 15/50 e 10/50 são preservados apenas como evidência da triagem inicial; não são
  mais o veredito operacional.

Nenhum resultado da triagem inicial está autorizado para auto-match produtivo.

## 5. O que conseguimos calcular

Para cada catálogo com ofertas, o probe produziu em dois níveis:

- por anúncio: mínimo, P25, mediana, média, média aparada 10%, P75 e máximo;
- por seller: menor oferta válida de cada seller e a mesma distribuição.

A visão por seller evita que um vendedor com anúncios duplicados domine a estatística.

Exemplos semanticamente prováveis com evidência mais forte:

| CODPROD | Produto | Sellers | Mínimo | Mediana por seller | Máximo |
|---:|---|---:|---:|---:|---:|
| `39352` | Coifa Franke Linea Touch Ilha 90 cm | 8 | R$ 2.219,00 | R$ 2.998,00 | R$ 5.219,90 |
| `41937` | Acabamento Fani 509 preto | 6 | R$ 209,94 | R$ 223,54 | R$ 241,71 |
| `33590` | Torneira Deca Soul 1198.C38 | 18 | R$ 459,90 | R$ 889,06 | R$ 1.749,88 |
| `39452` | Chuveiro Deca Flex 1965.C | 11 | R$ 739,90 | R$ 1.198,49 | R$ 1.518,73 |
| `25070` | Chuveiro Deca Cubo | 9 | R$ 960,00 | R$ 1.637,47 | R$ 3.200,90 |
| `39342` | Torneira Deca Soul 1178.GL38.RD | 10 | R$ 949,00 | R$ 1.041,19 | R$ 1.616,68 |

Faixas muito abertas não devem virar recomendação automática. Podem refletir outliers, promoções,
duplicidade, condições comerciais ou erros de variante ainda não capturados.

## 6. Decisão

### Confirmado

- O fluxo pré-anúncio funciona sem item próprio ativo.
- A busca oficial por EAN respondeu em 50/50.
- Houve cobertura bruta de catálogo de 44% e cobertura bruta de preço de 36%.
- A lista de ofertas permitiu calcular distribuições completas em 18 produtos.
- O mínimo isolado é uma referência ruim; mediana por seller é mais defensável quando há `n≥5`.

### Refutado

- “EAN exato basta para aceitar o match.” Dois falsos matches claros refutam essa regra.
- “Encontrou preço, então temos mercado suficiente.” Seis dos 18 produtos precificados tinham no
  máximo dois sellers.
- “O batch prova precisão de 99%.” Mesmo 50/50 acertos dariam limite inferior bilateral de 95% de
  aproximadamente 92,9%; esta amostra teve contradições.

### Desconhecido

- Cobertura representativa de todos os 10.519 produtos; esta coorte contém somente GTIN válido e único.
- Origem dos GTINs contraditórios.
- Solução durável após eventual desligamento de `/products/{id}/items`.
- Benchmark externo `INT_EXT` antes de existir um item do seller.
- Previsão de unidades vendidas/demanda; este batch mede preço, não volume.

## 7. Fallbacks oficiais para os 28 EANs sem catálogo

Uma segunda execução testou três rotas sem enviar descrição ou referência do ERP:

| Rota | Resultado |
|---|---:|
| `/product-identifier/validator` | 28/28 HTTP 200 |
| `/products/search?...&q={EAN}` | 28/28 HTTP 200; zero match com EAN exato |
| `/sites/MLB/search?q={EAN}` | 0/28 utilizável; 28/28 HTTP 403 |

**Conclusão confirmada:** repetir o EAN como texto livre não recupera os 28 produtos ausentes do
índice de `product_identifier`. O validador confirma a forma do código, não a existência de catálogo.
A busca geral de anúncios continua indisponível para a aplicação atual. Esse fallback acrescentou
zero cobertura e deve ser descartado.

O próximo fallback oficial que merece teste usa identificadores públicos diferentes do EAN:
referência do fabricante/part number e, depois, descrição + atributos. Esses métodos precisam passar
pelo mesmo gate semântico e não podem transformar candidato textual em match automático.

### Canário oficial por referência, domínio e atributos

Três produtos sem match foram testados em domínios distintos:

- coifa Franke `FCR925`, referência `14568`;
- argamassa Ceramfix AC-III, referência `603005`;
- acabamento Docol New Edge, referência `90010032072`.

Resultados:

- `q={referência}` e `q={marca+referência}` responderam `200`, mas geraram ruído severo: a referência
  da coifa trouxe para-choque, alternador e farol; a referência da argamassa trouxe kit automotivo;
- `domain_discovery` identificou corretamente coifas, não encontrou domínio para a argamassa e
  produziu domínios genéricos/plausíveis para o acabamento Docol;
- os metadados das categorias confirmaram os IDs públicos de `BRAND`, `MODEL`, `WIDTH` e/ou `MPN`;
- três buscas documentadas `POST /products/search` por domínio + três atributos responderam `200` e
  retornaram zero candidatos;
- nenhum anúncio ou recurso comercial foi criado ou alterado.

**Conclusão:** referência textual isolada não é matching e o POST por atributos é seguro como busca,
mas não recuperou cobertura adicional neste canário. O catálogo oficial continua sendo um ótimo
matcher quando contém o produto, porém não resolve sozinho a cauda fora do catálogo.

O estágio GET avançado foi então ampliado para os 28 produtos sem EAN no catálogo:

- `domain_discovery` encontrou ao menos um domínio para 18/28;
- busca por referência exata retornou algum resultado em 6/28, todos ruído ou produto incompatível;
- busca marca+referência restrita ao domínio retornou candidatos em 5/28;
- auditoria conservadora: 0 `ACCEPT`, 0 `REVIEW`, 10 `REJECT` e 18 `NO_CANDIDATE`;
- cobertura incremental confirmada: **0/28**.

Mesmo os candidatos aparentemente próximos falharam em variantes decisivas: coifa de parede versus
ilha `FCR925`, linhas/funções diferentes para a saboneteira Zen e kit de dois chuveiros cromados em
vez do acabamento Docol New Edge ouro escovado. A busca oficial avançada foi esgotada nesta coorte
sem ampliar a cobertura.

## 8. Gate semântico profissional proposto e reauditoria por atributos

O gate deve ser triclasse e determinístico:

- **ACCEPT:** GTIN exato, domínio compatível, nenhuma contradição e pelo menos duas âncoras
  independentes entre marca, referência/modelo, família e atributos de variante;
- **REVIEW:** família compatível, mas um atributo determinante está ausente ou divergente;
- **REJECT:** conflito de domínio, marca, modelo, medida, capacidade, cor, acabamento, voltagem,
  unidade ou composição do kit.

O detalhe público de catálogo respondeu `200` em 22/22 e forneceu:

| Atributo | Cobertura |
|---|---:|
| marca | 22/22 |
| modelo | 20/22 |
| cor | 15/22 |
| acabamento | 9/22 |
| voltagem | 3/22 |
| unidades por kit | 3/22 |
| part number | 1/22 |

Com esses atributos, a reauditoria conservadora separou:

- **12 `ACCEPT`**: domínio, marca e atributos discriminantes coerentes;
- **8 `REVIEW`**: evidência incompleta ou conflito de variante/título/atributo;
- **2 `REJECT`**: produto e domínio incompatíveis.

Dos 12 `ACCEPT`, 11 possuem algum preço e apenas 7 possuem ao menos cinco sellers. Portanto, o
resultado operacional conservador desta coorte é:

- match aceitável: 12/50 (24%);
- match aceitável com preço: 11/50 (22%);
- match aceitável com preço e evidência de ≥5 sellers: 7/50 (14%).

Os oito casos enviados para revisão foram `39352`, `17530`, `41937`, `34602`, `30585`, `39957`,
`39197` e `17284`. Um exemplo importante: `34602` possui modelo estruturado coerente com o ERP,
mas o título do próprio catálogo afirma outra variante. Inconsistência interna do catálogo nunca
pode ser resolvida silenciosamente.

O resultado demonstra que a regra é promissora, não que está validada: ela foi derivada e avaliada
na mesma amostra. Mesmo 12/12 não estabelece precisão de 99%.

O vetor redigido de atributos públicos foi persistido no artefato temporário e deve fazer parte de
qualquer futura prova. O próximo gate estatístico ainda exige regras congeladas e coorte independente.

### PoC executável conservadora

Uma regra local genérica foi executada sobre os 22 candidatos, sem hardcode de decisão por SKU. Ela
normaliza marca/referência, mapeia família do ERP para domínio Mercado Livre, exige múltiplas âncoras
e faz contradições vencerem o GTIN.

| Decisão da PoC | Resultado |
|---|---:|
| `ACCEPT` | 8 |
| `REVIEW` | 12 |
| `REJECT` | 2 |

Comparada à adjudicação conservadora 12/8/2, a PoC:

- aceitou 8/12 casos aceitos;
- enviou 4/12 aceitos para revisão;
- enviou 8/8 ambíguos para revisão;
- rejeitou 2/2 falsos;
- não promoveu falso ou ambíguo para `ACCEPT` nesta amostra.

No batch total, a cobertura automática exploratória seria 8/50 (16%); 7/50 (14%) possuem algum
preço e 5/50 (10%) possuem preço com pelo menos cinco sellers. Isso é comportamento conservador
adequado, mas não prova precisão: regras e avaliação usam a mesma coorte e o mapa de domínios ainda
é allowlist incompleta.

## 9. Alternativas profissionais para ampliar cobertura

Shortlist para benchmark, sem aprovação produtiva:

1. **API oficial + matcher próprio:** núcleo autorizado e já testado, mas com cobertura e continuidade
   insuficientes como fonte única de preços.
2. **[DataForSEO Merchant API](https://dataforseo.com/apis/merchant-api):** melhor challenger
   programável; canário real ainda pendente.
3. **[JoomPulse](https://joompulse.com/scanner-de-catalogos):** melhor ajuste funcional encontrado
   para catálogo ERP pré-anúncio; scanner em lote e MCP existem, mas os termos padrão bloqueiam OEM.
4. **[Precifica](https://precifica.com.br/monitoramento-de-ofertas-mercado-livre/):** melhor candidata
   brasileira para contrato enterprise/API sob medida.
5. **[Nubimetrics](https://centrodepartners.mercadolivre.com.br/apps/nubimetrics):** referência manual
   de inteligência de mercado; EAN batch/API/OEM não comprovados.

O benchmark GeckoAPI foi executado ao vivo com os mesmos 50 GTINs válidos e únicos: 49 respostas
HTTP 200 e uma 502, sem retry. As 49 respostas bem-sucedidas declararam entre 8 e 53 resultados
(média 43,69), mas **nenhum item retornou EAN estruturado idêntico ao consultado**. Portanto, sob o
contrato testado (`PLP` por keyword EAN), a cobertura verificável foi 0/50 e a cobertura de preço
aceitável também 0/50. Isso não prova que não existam anúncios: prova que os resultados textuais do
fornecedor não carregaram o identificador necessário para confirmar identidade. Aceitar seus preços
sem matching semântico independente repetiria justamente o erro já observado na API oficial.

**Decisão:** GeckoAPI não deve ser comprado nem integrado com esse fluxo como solução de matching
por EAN. Uma nova avaliação só faria sentido se o fornecedor demonstrar endpoint que preserve GTIN
verificável ou retorno de PDP com atributos suficientes para o gate semântico, sob novo orçamento de
créditos e protocolo separado.

### Busca pública e scraping próprio

Um canário no navegador com o EAN `7892162167027` confirmou que a busca pública do Mercado Livre
renderiza resultados utilizáveis: 8 resultados declarados, títulos, links e preço atual separado do
preço anterior. Os primeiros resultados eram semanticamente compatíveis com a coifa Franke do ERP,
enquanto resultados patrocinados posteriores já introduziam produtos incompatíveis. Logo, a coleta
é **tecnicamente possível**, mas exige matching e não pode aceitar a posição da busca como verdade.

Ela não é uma opção profissional autorizada hoje. As cláusulas 7.5 e 7.6 dos
[Termos do Programa de Desenvolvedores](https://developers.mercadolivre.com.br/pt_br/termos-e-condicoes)
proíbem contornar limitações e usar robôs, spiders ou scraping para acessar o site/conteúdo ou obter
informações fora do programa. Crawlee, Playwright e Scrapy são frameworks maduros, mas não resolvem
essa restrição; stealth/proxy agravariam o risco. Repositórios específicos auditados são exemplos
educacionais ou coletores de PDP/ofertas e nenhum prova matching por EAN. Portanto, scraper próprio
só deve ser reconsiderado com autorização escrita do Mercado Livre.

### Segundo challenger externo

[ShopAPIS](https://shopapis.com/api/getting-started) documenta consulta direta por GTIN/EAN e retorno
normalizado de `identifiers.gtin`, título, marca, preço e URL para Mercado Livre. É mais alinhado que
a rota textual da Gecko, mas continua **não validado**: não há prova independente encontrada, preços
e créditos exatos ficam no painel e os termos/política não estavam disponíveis na auditoria. O gate
correto é trial com os mesmos 50 EANs, sem retry, contando apenas GTIN idêntico e depois aplicando o
gate semântico. Não comprar antes de cobertura incremental, zero falso aceito e origem/licença
contratualmente comprovadas.

Uma diligência posterior rebaixou a ShopAPIS antes de qualquer envio de dados: o site não apresenta
entidade legal/endereço verificável, Termos, Política de Privacidade ou DPA; preços/créditos ficam
ocultos no painel, há inconsistência de marca no contato e as alegações de cobertura/uptime não têm
evidência independente. Ela não deve receber dados nem ser contratada até sanar esses pontos.

### Via brasileira certificada: Shopping de Preços

O [Shopping de Preços](https://centrodepartners.mercadolivre.com.br/apps/shopping-de-precos) é o
challenger contratual mais forte encontrado: aplicativo Platinum na Central oficial, trial de 15
dias, monitor de concorrentes, entidade/canais brasileiros e selo SOC 2 Type 2 declarado. Sua
[documentação](https://documentacao.shoppingdeprecos.com.br/pt/article/introducao-a-api-do-shopping-de-precos)
confirma API aberta mediante usuário/senha, token `X-Client`, cadastro empresarial e NDA.

Há uma limitação decisiva ainda não resolvida: a documentação técnica pública da API cobre HUB,
produtos, integrações, anúncios e pedidos, mas **não publica endpoints do Monitor de Concorrentes**.
A interface do monitor busca concorrentes por palavra-chave e requer confirmação manual; a análise
ampla também é limitada ao nicho/categorias de sellers com vendas mínimas. Antes de contrato, o
fornecedor deve demonstrar em nossos 5 EANs: pesquisa pré-anúncio, retorno de GTIN/MLB/título/preço,
licença para exibir os dados dentro de nossa plataforma, batch/rate limits, retenção, SLA e preço.
Se passar, executar os mesmos 50 produtos sob NDA; se não houver API redistribuível, serve apenas
como benchmark manual e não como componente da arquitetura.

### Rota oficial de autorização

O Mercado Livre mantém o [Developer Partner Program](https://developers.mercadolivre.com.br/pt_br/dpp),
mas exige escala relevante (no Brasil, GMVe mensal de USD 2,5 milhões), assessment de segurança e
iniciativas atribuídas; não é solução imediata. O próximo passo oficial é abrir ticket no
[Suporte Brasil para desenvolvedores](https://developers.mercadolivre.com.br/pt_br/suporte-brasil)
perguntando por endpoint pré-listagem, limites de exibição/cache e parceiro com API redistribuível.
Certificação de parceiro reduz risco operacional, mas os próprios termos esclarecem que não é
endosso, garantia de precisão ou licença automática de redistribuição.

#### Perguntas obrigatórias ao Shopping de Preços

1. A API oferece pesquisa de produtos **antes de existir anúncio do seller**, por EAN/GTIN, marca,
   referência e atributos? Enviar exemplos reais de request/response.
2. O retorno inclui GTIN verificado, MLB/MLBU, título, condição, preço atual, preço anterior, seller
   anonimizado e timestamp? Há múltiplas ofertas ou somente um resultado?
3. O contrato permite incorporar e exibir min/mediana/faixa desses dados a clientes da nossa
   plataforma? Quais limites de retenção, cache e redistribuição?
4. Qual é a origem autorizada dos dados do monitor, SLA, frequência, rate limit, batch, custo por
   consulta/conta e tratamento de 404/429?
5. Aceitam executar um canário sob NDA com 5 EANs e, se aprovado, os mesmos 50 produtos do benchmark?

#### Perguntas obrigatórias ao suporte do Mercado Livre

1. Existe endpoint autorizado para obter referência/faixa de preço de mercado **pré-listagem** por
   `product_id`, GTIN ou atributos, sem publicar um item?
2. O benchmark `INT_EXT` mostrado na criação do anúncio pode ser consumido por aplicativos? Qual
   recurso, elegibilidade e contrato?
3. É permitido mostrar aos sellers OAuth estatísticas agregadas (mínimo, mediana e faixa), com qual
   retenção/cache?
4. Qual parceiro certificado possui API B2B licenciada/redistribuível para esse caso?
5. Se não houver endpoint, há autorização formal possível para coleta limitada da busca pública?

### APIs brasileiras por EAN que não cobrem o canal-alvo

Duas APIs legítimas foram avaliadas e rejeitadas **por incompatibilidade de cobertura**, não por falta
de recursos técnicos:

- **Cnova Tech Data Market:** entidade brasileira identificada, EAN validado, min/máximo/média,
  plano grátis de 50 consultas e documentação REST. A oferta pública cobre redes de supermercado,
  e-commerce alimentar e delivery. Mercado Livre aparece apenas no produto separado Omnichannel,
  sem prova de que seus dados sejam expostos na API. Não executar batch.
- **Falcon Data Hub:** endpoint EAN, faixa/média/histórico e trial. Um GET público executado na lista
  de fontes retornou exatamente 7 redes (Atacadão, Carrefour, Higa Atacado, Hortifruti, Sam's Club,
  São Vicente e Super Pague Menos), nenhuma Mercado Livre. Os próprios materiais apontam foco em
  supermercado/FMCG e baixa cobertura de cauda longa. Não executar batch.

Ambos também deixam a licença de redistribuição dependente das restrições da fonte original; isso
não satisfaz o contrato de exibição em nosso SaaS. São possíveis fontes futuras de preço geral para
FMCG, mas não respondem ao requisito atual de preço **dentro do Mercado Livre**.

### Candidata SaaS testada e reprovada: Pricefy

[Pricefy](https://www.pricefy.io/features/marketplace-competitors-analysis) é o primeiro fornecedor
encontrado que documenta literalmente a combinação necessária: MercadoLibre em todos os países,
descoberta por EAN/GTIN, preços e disponibilidade por seller, atualização diária, mínimo/mediana,
histórico e API REST. O plano gratuito anunciado comporta 50 SKUs, 5 concorrentes e uma atualização
diária, exatamente suficiente para repetir o benchmark sem compra.

Isso ainda não equivale a aprovação produtiva. A ingestão/matching aparenta ser assíncrona, não uma
chamada síncrona EAN→ofertas; é preciso confirmar MLB Brasil, tempo de descoberta, GTIN/MLB no
retorno e disponibilidade da API no plano gratuito. Mais importante: os Termos padrão concedem uso
somente interno, pessoal e intransferível, sem sublicença/revenda do serviço ou componentes. Logo,
um canário pode validar a tecnologia, mas **incorporar resultados no nosso SaaS exige contrato OEM
ou licença de dados por escrito**.

Gate proposto: importar 5 produtos reais (incluindo dois falsos GTIN já conhecidos), aguardar o ciclo
de matching, exportar resposta da API e exigir GTIN idêntico + gate semântico. Só então importar os
50. Go técnico requer zero falso aceito e cobertura incremental; go comercial requer entidade/DPA,
origem autorizada dos dados, SLA, retenção e direito explícito de exibição multi-tenant.

O canário foi congelado antes de qualquer transferência em
[`evidence/2026-07-12-pricefy-canary5-input.json`](evidence/2026-07-12-pricefy-canary5-input.json).
A OpenAPI confirma que `monitored_urls` inclui `code` observado do concorrente, além de título,
marca, URL, preço, quantidade e `analyzed_at`; portanto `code == EAN` pode ser um gate real. O arquivo
permanece marcado `external_transfer_authorized: false` até autorização específica do usuário.

O coletor read-only foi preparado em `C:\tmp\pricefy_canary5_read.py`: exatamente 5 GETs em
`/v1/merchant/products/{sku}`, sem writes/retries, chave somente em memória, remoção de query strings
e persistência apenas de URL, título, marca, código observado, preço, status, quantidade e timestamp.
Dry-run e fixture de igualdade de código/domínio/sanitização passaram. SHA-256:
`871AADD3F15049F457A1A0DAA545FFDB14A4C0D21CC9C026C690C25283809C05`.

A conta real foi validada por `GET /v1/merchant/account`: HTTP 200, plano Free, limite de 50
produtos, 5 competidores, 100 URLs monitoradas, 3 automatches e 20 single-discoveries/dia; uso atual
zero. O status de assinatura retornou `inactive`, apesar da API key válida, e deve ser observado no
upload. Para reduzir exposição, os scripts foram revisados para transmitir somente IDs sintéticos
`PF-CANARY-1..5`, nunca o SKU ERP. Hashes atuais: upload
`941F0873EFA09678528A6802D52D0ED00F8A4E65C0D09FA6B0921AA0C606B6ED`; leitura
`CA2BF69025DEFA52AD8B69EB15ECEAF754F8A96F4F6D1077CF799060D732DB1A`.

#### Resultado live do canário Pricefy

O canário foi executado no dashboard real com autorização informada. A API Free aceitou leitura,
mas rejeitou PUT (`400`, `catalog readonly access via API`); os cinco produtos foram então criados
pela UI com IDs sintéticos e confirmados por GET, com 5/5 códigos válidos. O Differential Discovery
gratuito foi iniciado com `Title + Code` e custo zero.

O resultado foi materialmente inválido: em vez de Mercado Livre, encontrou `melbourneit.au`,
`registry.gov.in`, `stcu.org`, `videomagazine.io` e `expolingua.com`. Todos foram deixados sem seleção
para não consumir Automatch. Após mais de 110 segundos (estimativa 35 s), o job ainda indicava busca.
Na revalidação final autenticada, o processamento assíncrono já constava como concluído, o catálogo
mostrava os cinco produtos e a UI anunciava 14 novas sugestões; porém a revisão apresentava somente
os mesmos cinco domínios irrelevantes, todos `Pending Approval`, com zero URL monitorada, zero preço e
zero domínio Mercado Livre. A conta já estava configurada para Brazil; somente a moeda foi corrigida
de USD para BRL, o que não explica a seleção de domínios.

Uma revalidação autenticada em 13/07 eliminou a hipótese de conta errada: o e-mail retornado pela
API corresponde ao e-mail visível no dashboard. Inicialmente o dashboard exibiu `0 Products` e
`Last import: Never`, enquanto `GET /v1/merchant/account` informou uso de 5 produtos e os cinco
`GET /v1/merchant/products/PF-CANARY-*` responderam `200`; depois a UI convergiu para 5 produtos.
Logo, a divergência era transitória/eventual, não uma segunda conta. O retorno final permaneceu com
zero URL monitorada, zero domínio Mercado Livre e zero código observado exato, sem alterar o
`NO-GO` técnico.

**Decisão:** Pricefy é `NO-GO` para o fluxo testado. Não executar batch50, não pagar e não negociar
OEM enquanto o fornecedor não demonstrar, em nossos mesmos EANs, discovery restrito a
`mercadolivre.com.br` e resultado verificável. Evidência estruturada:
[`evidence/2026-07-12-pricefy-canary5-result.json`](evidence/2026-07-12-pricefy-canary5-result.json).

### JoomPulse: melhor aderência funcional, bloqueada para OEM no contrato padrão

A diligência de 13/07 encontrou uma solução que implementa quase literalmente a dor relatada. O
[Scanner de Catálogos](https://joompulse.com/scanner-de-catalogos) recebe um ou vários PDFs de
fornecedor, extrai nome, referência, custo, EAN, peso e dimensões e cruza todos os itens com produtos
similares do Mercado Livre. A saída anunciada inclui top 10 similares, preço individual, preço médio
de mercado, receita média, média de vendas, seller e exportação XLS. O
[guia oficial](https://blog.joompulse.com/2026/05/09/scanner-catalogo-fornecedor-joompulse/)
descreve explicitamente catálogo com centenas de itens ainda não validados pelo seller.

Há também um [MCP OAuth oficial](https://joompulse.com/mcp-connector) para ChatGPT/Claude. O
[repositório público de skills](https://github.com/joomcode/pulse-skills) comprova pesquisa semântica
de produtos e processamento de linhas CSV/XLSX/campos digitados. A skill de análise por produto
documenta preço atual, mínimo histórico, vendas/receita estimadas, status de catálogo, preço de
buy-box, quantidade de sellers e links Mercado Livre/JoomPulse. A skill de produto idêntico exige
validar marca, modelo, variante, cor, tamanho, capacidade e kit, e recomenda precisão acima de
recall. Isso valida que o fornecedor possui uma superfície automatizada real, não apenas dashboard.

Ainda não está comprovado que EAN seja preservado como chave de matching no MCP ou na resposta do
scanner; os próprios materiais dizem “similares”. Tampouco há prova pública de mediana nativa,
ofertas deduplicadas por seller, API REST server-to-server, SLA ou DPA para nossa operação.

Um probe público do protocolo confirmou `POST /mcp -> 401` com metadata OAuth válida, escopo `mcp`,
PKCE S256 e grants `authorization_code` + `refresh_token`; `client_credentials` não é anunciado.
No commit público `75f9070254a642a1e1e0d96f456ab72048f442e0`, as 17 skills não contêm os termos
EAN, GTIN, SKU, UPC ou barcode. O fluxo tabular usa título como busca quando não recebe MLB. Logo o
MCP é tecnicamente real, mas não prova a rota automática `EAN -> produto exato`, e a autenticação
inicial é interativa por conta.

O bloqueio decisivo é contratual. Os [Termos de 06/05/2026](https://joompulse.com/termos-e-condicoes)
limitam a licença a uso comercial interno, declaram os outputs propriedade da JoomPulse e proíbem
copiar/distribuir, criar API derivada, automação não autorizada e revender/explorar os dados. O MCP é
um canal autorizado para o assinante, não uma licença OEM para exibir resultados aos clientes do
Marketplace Central.

**Decisão:** `GO` como canário técnico/manual dos mesmos cinco produtos e, se passar, batch50;
`NO-GO` para integração no contrato padrão. A adoção só pode ocorrer com aditivo escrito concedendo
API/MCP server-to-server, EAN/SKU batch, exibição multi-tenant, cache/derivados, SLA/DPA e direitos de
redistribuição. Evidência:
[`evidence/2026-07-13-ml-partners-prelisting-audit.json`](evidence/2026-07-13-ml-partners-prelisting-audit.json).

### Nubimetrics e parceiros certificados

A [Nubimetrics](https://www.nubimetrics.com/br/product/mercado) permite descobrir produtos antes de
anunciar por palavra-chave e mostra preço médio, último preço, vendas e concorrência. Porém não há
prova pública de ingestão EAN/SKU em lote ou API; seus
[termos](https://www.nubimetrics.com/pdf/terminos_pt.pdf) proíbem acesso automatizado e extração sem
autorização escrita. Serve como referência manual/enterprise, não como dependência do produto.

O Shopping de Preços permite busca por keyword e exportação Excel, mas somente dentro de categorias
elegíveis pelo histórico da conta do seller; sua API pública cobre operação/hub, não foi localizado
endpoint de inteligência pré-listagem. Seconds não documenta intake EAN/SKU/API para produto
arbitrário, e Mercado Cuentas vincula concorrentes às publicações existentes. Nenhum supera
JoomPulse como benchmark funcional nem DataForSEO como próximo teste programático.

### Alternativa enterprise brasileira: Precifica

[Precifica](https://precifica.com.br/monitoramento-de-ofertas-mercado-livre/) documenta serviço
específico para Mercado Livre, descoberta a partir de SKUs/categorias/termos e entrega por painel,
Excel, API ou FTP. É uma empresa brasileira mais adequada a negociação customizada, mas não publica
EAN, preço, trial nem licença OEM. Fica como segunda conversa comercial se Pricefy falhar técnica ou
contratualmente; não há base para integração antes de demonstração nos mesmos 5 EANs.

Uma diligência posterior elevou a Precifica a **challenger enterprise número 1**. O FAQ oficial
declara matching por EAN, nome, especificações técnicas e fotos; a entrega pode ocorrer por API para
ERP/sistemas internos e o serviço localiza links correspondentes no Mercado Livre a partir da lista
do cliente. Isso satisfaz teoricamente o fluxo pré-anúncio. Continuam desconhecidos: schema real,
precisão, faixa/mediana, SLA, preço, DPA, origem autorizada dos dados e direito OEM/multi-tenant.
Portanto, a Precifica está aprovada apenas para demonstração/canário contratual, não para produção.

Os probes públicos passivos confirmaram `monitor.precifica.com.br` ativo (`200`) e hosts de app/API,
mas nenhum Swagger, OpenAPI, sandbox ou endpoint público utilizável. Casos publicados corroboram
escala e pesquisa pré-cadastro: a Pill recebeu preços antes de seu e-commerce entrar no ar e chegou a
540 mil preços/mês; a Cobasi declara 4 mil itens e 770 mil preços/mês. A implantação divulgada é de
aproximadamente seis semanas, sinal de que o matching pode incluir configuração e homologação, não
uma simples chamada síncrona `EAN -> ofertas`.

Gate contratual obrigatório: a API precisa aceitar ID sintético + EAN + marca + referência +
descrição, restringir a fonte ao Mercado Livre Brasil e devolver oferta rastreável (MLB/URL, seller,
preço atual/original e timestamp). O contrato precisa autorizar SaaS multi-tenant, cache, derivados e
exibição aos clientes do Marketplace Central. Nada disso está comprovado publicamente.

### Nova rota oficial durável e teste executado

Em 2025 o Mercado Livre documentou `buy_box_winner` e `buy_box_winner_price_range` no
`GET /products/{product_id}`. Em tese, esses campos oferecem preço vencedor e faixa min/max antes de
o seller publicar. Essa rota é mais durável que `/products/{product_id}/items`, cujo desligamento foi
anunciado, mas não fornece mediana, média nem todos os concorrentes.

O probe live foi repetido com OAuth renovado, somente `GET`, nos 22 catálogos do lote:

| Métrica | Resultado live |
|---|---:|
| `/products/{id}` HTTP 200 | 22/22 |
| `buy_box_winner` não nulo | 0/22 |
| `buy_box_winner_price_range` não nulo | 0/22 |

**Confirmado:** todos os catálogos continuam ativos, porém nenhum expôs winner/faixa para esta
coorte. Logo, a rota é oficial, mas não preservou cobertura de preço e não substitui a lista de
ofertas no nosso caso real.

A auditoria da documentação vigente confirmou que essa continua sendo a melhor cadeia oficial:
`products/search` (EAN, part number/título ou POST por domínio+atributos) → gate forte de identidade
→ `products/{id}`. `price_to_win` e referências de preço exigem item existente; o novo MCP oficial
pesquisa documentação, não produtos. A documentação EN marca `/products/{id}/items` para shutdown,
enquanto a PT ainda o mostra e o runtime MLB respondeu no probe: conflito suficiente para classificá-lo
como legado atrás de health-check, nunca dependência principal. Evidência:
[`evidence/2026-07-13-ml-official-current-route-audit.json`](evidence/2026-07-13-ml-official-current-route-audit.json).

Também foi publicada a API oficial de
[referências de preços](https://developers.mercadolivre.com.br/pt_br/referencias-de-precos):
`GET /suggestions/user/{seller_id}/items` e `GET /suggestions/items/{item_id}/details`. Ela retorna
preço sugerido, menor preço, referência interna e metadados de similares, mas exige item existente e
token do proprietário. Isso é compatível com o modal mostrado pelo usuário, porém a equivalência
continua sendo **inferida**, não documentada.

O item `MLB3474723449` da captura foi consultado com a conta OAuth atualmente conectada:

- `GET /items/MLB3474723449`: `403 access_denied`;
- `GET /suggestions/items/MLB3474723449/details`: `401`, `Caller is not the item's owner`.

**Confirmado:** esse item não pode ser testado com a conta conectada. **Desconhecido:** se pertence a
outra conta do usuário e se a resposta reproduziria R$ 971,01. O teste correto exige conectar a conta
proprietária; criar anúncios pausados/estoque zero apenas como sonda não está autorizado nem
documentado e foi rejeitado.

### Cargoos: canário funcional aprovado, OEM ainda reprovado

[Cargoos](https://cargoos.com.br/) corresponde literalmente à dor: pesquisa em lote de EANs antes
de comprar estoque, faixa de preço, concorrentes, demanda e posicionamento. A
[Chrome Web Store](https://chromewebstore.google.com/detail/cargoos-%E2%80%94-intelig%C3%AAncia-pa/hhpmjdebejnbpidgokkeikoainkcmlgi)
confirma versão 1.0.15 em 12/07/2026, 222 usuários, 23 avaliações e retorno de preço,
vendas e concorrência por múltiplos EANs. A empresa informa Understack Tecnologia LTDA, CNPJ e
política LGPD; plano Pro anunciado a R$ 49,99/mês com sete dias de teste.

Limites decisivos: nenhuma API/export/webhook pública, termos de uso apontando para página inválida,
sem SLA, DPA ou licença OEM. Portanto:

- `GO` apenas como benchmark manual dos mesmos cinco EANs;
- `NO-GO` como dependência do Marketplace Central sem contrato/API.

Antes de criar conta, a extensão pública foi baixada diretamente do serviço oficial da Chrome Web
Store e auditada estaticamente, sem instalação ou execução. O artefato atual é versão `1.0.15`, CRX3,
1.266.510 bytes, SHA-256
`2F0AC889F52F323F44FE87B8D13F91ED2C5DF027449BB3F94BA83D9F18D84AB0`.

A extensão contém um fluxo antigo que contradiz a promessa comercial:

- a UI aceita EAN-8/EAN-13, remove duplicados, mostra créditos, histórico e import/export JSON;
- há backend autenticado em `https://prod.cargooscorp.com.br/api/v1` e endpoints para saldo e
  persistência das análises;
- o resolvedor usado pelo batch retorna um objeto vazio incondicionalmente;
- o orquestrador chama esse resolvedor sem enviar o lote atual e, portanto, marca todas as linhas
  como `not_found`;
- a função seguinte devolve sucesso sintético com visitas/vendas zero e um seller, sem request de
  análise real;
- nenhum endpoint de rede para resolver/analisar EAN foi encontrado nesse fluxo da extensão.

Uma segunda auditoria encontrou a explicação: a pesquisa antiga foi explicitamente migrada para o
dashboard. O bundle atual de `app.cargoos.com.br` possui um fluxo server-side real:

- `POST /ean-analyses` cria análise; há listagem, detalhe, execução, retry, histórico e exportação;
- aceita CSV/XLSX e limita cada coleção a 1.000 EANs;
- a UI faz polling a cada 5 segundos e calcula Lite = 1 crédito/EAN e Deep = 3 créditos/EAN;
- o resultado observado contém `catalogProductId`, `catalogName`, um `salePrice`, total de sellers,
  visitas e vendas estimadas, logística, reviews e custos;
- o contrato batch não expõe ofertas individuais, preço mínimo, média ou mediana;
- o modo Lite/Deep aparece na UI, mas o request de criação observado não transmite `depth`, possível
  inconsistência que precisa ser validada ao vivo;
- parte das métricas requer conexão OAuth com Mercado Livre e a arquitetura combina API oficial com
  extração de páginas; há dependência declarada de `/products/{id}/items`, rota de continuidade
  arriscada.

**Decisão corrigida:** `GO` para canário manual dos cinco EANs no dashboard, sem extensão. `NO-GO`
para integração produtiva até existir API server-to-server, DPA, SLA, licença OEM/redistribuição e
uma fonte de distribuição de preços. Evidências:
[`extension`](evidence/2026-07-13-cargoos-extension-audit.json) e
[`dashboard`](evidence/2026-07-13-cargoos-dashboard-audit.json).

### DataForSEO Merchant API: melhor challenger programável

A [Merchant API](https://dataforseo.com/apis/merchant-api) oferece uma rota profissional indireta:
buscar o produto no Google Shopping por EAN + marca + referência, obter o identificador global do
produto e então recuperar os sellers/ofertas, filtrando somente domínios do Mercado Livre Brasil.
Isso não exige anúncio próprio e não depende da busca interna por EAN do Mercado Livre.

**Confirmado nas páginas e documentação oficiais:** 

- `POST /v3/merchant/google/products/task_post` aceita keyword, localização Brasil, idioma e até
  100 tasks por chamada; o limite geral divulgado é 2.000 chamadas/minuto;
- o retorno de produtos inclui título, descrição, preço, ranking, rating e domínio relacionado;
- o endpoint Sellers retorna os dez principais vendedores de um produto, com preço base/total,
  frete, condições de compra, disponibilidade, domínio e URL direta da oferta;
- normal custa cerca de US$ 0,001 por produto/seller/SERP e alta prioridade US$ 0,002; há depósito
  mínimo de US$ 50, crédito inicial de US$ 1 e sandbox gratuito com dados sintéticos;
- o fornecedor declara que seu modelo é voltado a software builders e que clientes exibem dados a
  seus próprios end-users; os Termos incluem DPA, mas transferem ao cliente o risco de uso dos dados
  de mecanismos de busca.

**Ainda não provado:** Google Shopping Brasil pode não indexar todos os produtos/lojas do Mercado
Livre; `keyword=EAN` não é um parâmetro GTIN dedicado; e a permissão geral para software não substitui confirmação específica de
nosso caso de uso multi-tenant. O sandbox exige conta e devolve dummy data, portanto não mede
cobertura.

**Decisão:** melhor próximo canário automatizado, acima de Bright Data e Cargoos para arquitetura.
Enviar os mesmos cinco EANs em uma única task batch, buscar produtos e sellers em `pt-BR/Brasil`,
aceitar somente URL/domínio `mercadolivre.com.br` com gate semântico e, se passar, repetir os 50.
Não criar conta nem transmitir dados sem autorização específica no momento da execução.

Após o falso cluster Crystal→Linea observado na SearchAPI, o contrato atual da DataForSEO foi
revalidado via Context7. A documentação recomenda enviar juntos `product_id`, `data_docid` e `gid`;
Sellers também aceita `pvf`, filtro de variante retornado pelo Products. O probe final está em
[`probes/dataforseo_product_details_canary5.py`](probes/dataforseo_product_details_canary5.py) e
implementa `Products → Product Info → gate de identidade → Sellers com todos os IDs + pvf → gate
independente por oferta`. Confiança do Product Info nunca é herdada pelo preço final.

`--self-test`, dry-run e compilação passaram. O máximo é 10 tasks Products, 10 Product Info e 10
Sellers, mas Sellers só roda para candidatos aprovados. Nenhum request externo foi feito. A execução
real exige `--execute`, `DATAFORSEO_LOGIN`/`DATAFORSEO_PASSWORD` e a trava
`MPC_DATAFORSEO_TRANSFER_AUTHORIZED=YES`. URL `mercadolivre.com.br`, marca e referência ou ≥75% da
linha/variante continuam obrigatórios para alimentar estatísticas.

### SearchAPI.io: canário live executado e reprovado

A [Google Shopping API](https://www.searchapi.io/docs/google-shopping) da SearchAPI.io fornece
`product_token`; a rota [Google Product Offers](https://www.searchapi.io/docs/google-product-offers)
usa esse token recém-gerado e devolve ofertas com merchant, URL direta e preço numérico. Ambas
aceitam `gl=br`, `hl=pt-br` e autenticação Bearer. O plano inicial oferece 100 requests sem cartão;
os pagos começam em US$ 40/10 mil buscas e anunciam SLA de 99,9%.

O canário live autorizado mostrou que a API é tecnicamente real, mas não serve para este fluxo:

1. EAN + Product Offers: 20 requests, 0/5 preços aceitos. Os EANs levaram a biscoito, panos de
   limpeza, fita e esparadrapo no Mercado Livre; nenhum preço contaminou a saída.
2. Busca semântica + Google Product Details + Offers: 17 requests e duas aceitações aparentes. A
   auditoria manual encontrou um falso automático: o detalhe era Franke Crystal/FCR925, mas as
   ofertas pertenciam a Linea/16702. O cluster Product→Offers não preservou a variante.
3. Reexecução com gate independente em cada oferta: 19 requests, 0 falsos aceites, porém 0/5 preços
   aceitos. A oferta Linea correta observada na execução anterior desapareceu, provando também falta
   de repeatabilidade.

O probe profundo
[`probes/searchapi_product_details_canary5.py`](probes/searchapi_product_details_canary5.py)
revalida marca, referência, título/linha/variante, abreviações, tensão e kit na **oferta final**, não
herda confiança do cluster Google. Esse ajuste corrigiu o falso positivo, mas reduziu a cobertura
segura a zero. Portanto SearchAPI.io é `NO-GO`; não executar batch50 nem pagar. Seus termos também
continuam exigindo permissão expressa para exploração/OEM. Evidência de adjudicação:
[`evidence/2026-07-13-searchapi-canary5-adjudication.json`](evidence/2026-07-13-searchapi-canary5-adjudication.json).

### Batch50 de descoberta por mecanismo de busca externo

Para testar a arquitetura sem depender de fornecedor, foram executadas 100 consultas primárias nos
50 produtos reais: EAN exato e marca + referência restritos a `mercadolivre.com.br`, mais 16 buscas
adversariais para refutar candidatos ambíguos.

| Resultado conservador | Produtos | Cobertura |
|---|---:|---:|
| `ACCEPT` com MLB/MLBU, identidade e preço utilizável | 4 | 8% |
| `REVIEW` | 5 | 10% |
| `REJECT`/sem evidência | 41 | 82% |
| falso aceite observado entre os quatro aceitos | 0 | — |

EAN exato não produziu anúncio individual utilizável em nenhum dos 50. Um EAN apareceu dentro de
bundles, mas os preços dos kits foram corretamente descartados. Marca + referência recuperou alguns
produtos, porém snippets e páginas de lista frequentemente mostraram preço sem vínculo estável a
um MLB ativo. Exemplos aceitos: Franke 16702 (`MLB21014314`, R$ 2.289), Zen ZP5402.A00
(`MLB26670537`, R$ 63,24), Deca 1198.C38 (`MLB31007956`, R$ 573,61) e Deca 1178.GL38.RD
(`MLB5287842998`, R$ 899).

**Decisão:** busca externa serve somente para gerar candidatos em fallback/`REVIEW`. Com 8% de
cobertura automática segura, índice desatualizável e ausência de distribuição de sellers, não é a
fonte profissional principal. O canário DataForSEO continua válido porque testa Google Shopping
estruturado + Sellers, não snippets orgânicos.

### Challengers adicionais e rejeições

| Candidato | Evidência favorável | Bloqueio atual | Decisão |
|---|---|---|---|
| DataForSEO Merchant API | API batch, produtos+sellers, preço baixo, sandbox, uso por software builders | cobertura Mercado Livre por EAN e direito específico multi-tenant ainda não testados | challenger programável nº 1; canário de 5 |
| SearchAPI.io | Product→Offers estruturado, URL/preço/merchant e detalhes ricos | EAN trouxe itens errados; cluster cruzou Crystal→Linea; 0/5 após correção e sem repeatabilidade | `NO-GO`; batch50 proibido |
| JoomPulse Scanner + MCP | batch de catálogo pré-anúncio, EAN extraído, preço médio, XLS, MCP OAuth e matching semântico | matching EAN exato não provado; termos limitam uso interno e proíbem API derivada/redistribuição | challenger funcional nº 1; canário técnico, sem OEM |
| Bright Data | scraper explícito `mercadolivre.com.br products`, schema de preço/URL/ID, API/webhook/batch e SLA enterprise | EAN nativo não provado; licença padrão é uso interno; scraping segue proibido pela cláusula 7.6 do ML | melhor infraestrutura técnica, mas `NO-GO` sem autorização ML + SoW OEM |
| Oxylabs | API server-to-server, JS, batch/scheduler, trial, DPA | sem contrato GTIN; AUP exige cumprir os termos do alvo; revenda só por SoW | `NO-GO` padrão para ML |
| Zyte | extração product/productList, URL/preço/GTIN/MPN, DPA e alta vazão | sem extrator ML provado; restringe processamento para terceiros e pode suspender por violação do alvo | `NO-GO` OEM/fonte |
| Apify | Actors ML por keyword, API/OpenAPI e datasets JSON/CSV | Community Actors, matching EAN não provado e responsabilidade transferida ao usuário | POC apenas; não base profissional |
| Minderest | matching EAN/UPC + IA/revisão humana, ERP/BI, ISO 27001 declarada | termos padrão proíbem service bureau/redistribuição; MLB não nomeado | contrato especial obrigatório |
| InfoPrice | upload de EAN em escala e API OAuth | cobertura Mercado Livre e cauda de casa/construção não provadas | prioridade secundária |
| Rally ML Catalog API | API barata, SLA publicado, uso SaaS alegado | apenas keyword em pool pequeno; origem inclui scraping; sem GTIN; termos ainda template | rejeitado |
| Pricenly | Mercado Livre, preço médio/mínimo e pesquisa por keyword | sem EAN/API/OEM e foco em URL/termo | rejeitado |
| Retail Shake | afirma GTIN/EAN e batch de 50 | Mercado Livre Brasil não comprovado; exploração comercial restrita | rejeitado até prova |

Apify, scrapers comunitários, Playwright/Crawlee/Scrapy e fornecedores genéricos de scraping não
ganham conformidade por terceirizar a coleta. Os termos do Mercado Livre continuam proibindo robots,
spiders e scraping; contrato do coletor não substitui autorização da fonte. A auditoria comparativa
está em [`evidence/2026-07-13-professional-scraping-providers-audit.json`](evidence/2026-07-13-professional-scraping-providers-audit.json).

### Arquitetura global provisória baseada nas provas

```text
ERP (SKU + EAN + marca + referência + atributos)
  -> API oficial products/search por GTIN
  -> gate determinístico: domínio + marca + modelo/ref + variante; contradição vence GTIN
  -> se ACCEPT:
       /products/{id} winner/faixa (quando não nulo)
       fonte de ofertas autorizada (a homologar) para min/mediana/n/sellers
  -> se sem catálogo ou sem preço:
       DataForSEO ou fornecedor contratual homologado
       JoomPulse somente sob licença OEM escrita
  -> normalizar ofertas e deduplicar por seller
  -> estatísticas com n, timestamp, fonte e confiança
  -> REVIEW quando identidade/variante ou amostra forem insuficientes
  -> nunca inventar preço; “sem evidência” é resultado válido
```

Escolha global provisória: híbrido oficial-first + fornecedor contratual. Scraper próprio teria maior
cobertura aparente e menor custo inicial, mas perde em autorização, estabilidade e responsabilidade.
Um único SaaS caixa-preta seria simples, porém Pricefy e Gecko demonstraram como cobertura sem
identidade verificável produz números inúteis. O híbrido preserva provenance e permite trocar o
fornecedor sem alterar o matching.

Scrapers comunitários e Actors pequenos do Apify não atendem ao nível profissional exigido. Um
fornecedor de scraping também não elimina o risco contratual: os termos do Programa de Desenvolvedores
proíbem [scraping automatizado](https://developers.mercadolivre.com.br/pt_br/termos-e-condicoes). Qualquer fornecedor deve oferecer garantia de origem/licença, DPA/LGPD,
retenção, SLA, segurança e indenização.

Para entity resolution, a pesquisa favorece **[Splink 4](https://github.com/moj-analytical-services/splink)** como ferramenta de diagnóstico probabilístico
por atributo e **[dedupe OSS](https://github.com/dedupeio/dedupe)** para active learning após existir gold set. [Sentence Transformers](https://github.com/huggingface/sentence-transformers)
multilíngue pode recuperar/rankear candidatos, mas nunca neutralizar contradições. Ditto,
DeepMatcher e Magellan não foram escolhidos como núcleo devido a manutenção, dependências antigas,
necessidade de treino e menor explicabilidade. O benchmark metodológico recomendado é WDC Products,
mantendo um gold PT-BR próprio e split por entidade não vista.

## 10. Próximo gate correto

1. executar demonstração contratual da Precifica com os cinco EANs congelados, exigindo o mesmo
   retorno/API que seria contratado e direito OEM;
2. executar o canário manual Cargoos com os mesmos cinco EANs, medindo match, `salePrice`, sellers,
   latência e diferença Lite/Deep; não confundir esse resultado com aprovação OEM;
3. conectar a conta proprietária de `MLB3474723449` e repetir `/suggestions/.../details`, sem criar
   anúncio novo;
4. abrir consulta formal ao Mercado Livre sobre endpoint pré-item, TTL/derivados e substituto da lista
   de ofertas;
5. somente se um challenger passar zero falso aceite + pelo menos 3/5 matches corretos, executar os
   mesmos 50 produtos;
6. congelar o matcher e validar em coorte independente antes de planejar integração produtiva.

## 11. Evidência

- Resumo estruturado: [evidence/2026-07-12-ml-prelisting-batch50-summary.json](evidence/2026-07-12-ml-prelisting-batch50-summary.json)
- Input local temporário SHA-256: `3BFE157AF20FCBBE8552B63D193ED1B61423BF5312B96056007BCEA5D611FB3B`
- Resultado local temporário SHA-256: `6314818BAECADB2CA7C505C17D1FB9C5557F76325BB8B44ABA55DB09F651F27F`
- Probe SHA-256: `901060BF8B7900478DBA491998147A1C43BE5A03AF6F13E5889F8A80AA2AD77D`
- Comando wrapper SHA-256: `F7CF309362B126BDBE65CF6D1683DE4E2AA6314A8595288429AFB369EAB164A7`
- Execução: `PASS`, 50/50, teste `172,32 s`, pacote `192,925 s`.
- Fallback raw SHA-256: `B4D8751170A47C4BB7D85DF88BC37EE271EE6FD7650E1C98CCECF04AB86267FA`
- Fallback probe SHA-256: `4800532A3809C8BE6398AE84EBA04190981DB5C4550170CEAC08BA24DBC29A9C`
- Fallback wrapper SHA-256: `1A4A4933F21BD43D2D1E8BFD0EFFD4FB13C8BDF7239A7EA57B21C78F2BF871C3`
- Fallback: `PASS`, 28/28, teste `80,72 s`, pacote `84,927 s`.
- Canário GET avançado raw SHA-256: `652F407DF33B5BF0DF50D09B4E20B716A886C5EFC8DFCF42922B66A9BC67D2B5`
- Canário GET executado probe SHA-256: `FDB20F311650A1619AA7CD4CC377B9E716A528C97D1C7CF86DF8E5B1F2393E75`
- Canário POST por atributos raw SHA-256: `73F359F2248BBEF7BD1238659CB42E2143DA70FEF41D28D954C8E659B66ADC90`
- Canário POST probe SHA-256: `B3FF6700D9C8D12C05B7E24D5D0260992D09EE87CAD17B8266B697B9979CE364`
- Canário POST: 3/3 HTTP 200, zero candidatos, pacote `5,300 s`.
- Batch oficial avançado raw SHA-256: `5736049D3C4B915A5EED018C115392A6358CF51021B3E6BF8FF565F46BC8B6EF`
- Batch oficial avançado probe SHA-256: `A302566023603A60EBA810CB525DE443901F08B1ED6A7470D73F7B8C700C5D9E`
- Batch avançado: 28/28, 0 `ACCEPT`, 10 `REJECT`, 18 `NO_CANDIDATE`, cobertura adicional zero.
- Atributos públicos raw SHA-256: `46CBD2F34F5BDEB550BC45A85FE124DDF2C9F8622AE0DD2F6A962570F26E2401`
- Probe de atributos SHA-256: `13763FBA0D258E6A16A932A4930DF370708364FF82FA296764E57BEDA88B4B92`
- Wrapper de atributos SHA-256: `797E2BEA7150726F2BBEC836F30C6FE033B8F32BEC10547B6DDE262CE4E61AE3`
- Atributos: `PASS`, 22/22, teste `4,72 s`, pacote `8,702 s`.
- PoC local do gate SHA-256: `C3BD9B5716B7E5428F9218D4ECEC32ECD2E5E5759B8B510425BBF4B4E2A0CA16`
- PoC local: 8 `ACCEPT`, 12 `REVIEW`, 2 `REJECT`; zero falso/ambíguo aceito na coorte usada para derivá-la.
- GeckoAPI benchmark live: 50 requests; 49 HTTP 200, 1 HTTP 502, 0 EANs exatos e 0 preços aceitos.
- GeckoAPI output sanitizado SHA-256: `9FC64B9857CBD3ABE4D6416E949F029B842D16BAB47D8D53FED1FE83FC27CC05`.
- GeckoAPI probe executado SHA-256: `34EC2A82AB3824825EFE2CD24B4FF371743170AA2CDFF87A9CCBE5A300851D55`.
- GeckoAPI wrapper SHA-256: `821AA757EC63F8F005DF2ECFCC421E97B345C44535CBA24F6890F0324D680495`.
- Pricefy canário live final: 5 produtos, discovery concluído, 14 sugestões anunciadas, 5 domínios
  visíveis irrelevantes, 0 Mercado Livre, 0 URL monitorada e 0 preço.
- Resultado Pricefy sanitizado SHA-256: `435C27054B6371E1D2B0FDDB10EC5BB84B5E40B370349CFA4D896C7734167FA4`.
- Probe oficial durável: 22/22 HTTP 200; 0 winner; 0 faixa; item mostrado 403 e referência 401.
- Resultado oficial durável SHA-256: `D63509F3261D370227F6CF8890B1B2BBB040BBA22B985D6395892B873BE70D21`.
- Resultado estruturado: [evidence/2026-07-13-ml-official-durable-price-result.json](evidence/2026-07-13-ml-official-durable-price-result.json).
- Probe reproduzível: [probes/ml_official_durable_price_probe_test.go](probes/ml_official_durable_price_probe_test.go).
- Auditoria Cargoos: extensão 1.0.15 com fluxo legado stubado; dashboard atual com batch EAN
  server-side real; nenhum EAN enviado no probe estático.
- Evidência extensão Cargoos: [evidence/2026-07-13-cargoos-extension-audit.json](evidence/2026-07-13-cargoos-extension-audit.json),
  SHA-256 `13A401AD29B1AB4CA2C3D76FCEC268AB8048A60A4823CF950B8ED99BB61601FE`.
- Evidência dashboard Cargoos: [evidence/2026-07-13-cargoos-dashboard-audit.json](evidence/2026-07-13-cargoos-dashboard-audit.json),
  SHA-256 `F9195D14893F2B8B17428282E9835A4761DA98F06BC1E7B54F404C9C26013CD2`.
- Auditoria DataForSEO Merchant API: [evidence/2026-07-13-dataforseo-merchant-api-audit.json](evidence/2026-07-13-dataforseo-merchant-api-audit.json),
  SHA-256 `99BD8CC2095663987EC25AB8F96E680300A4D44514A438B3721BA264DD45F8B8`.
- Probe DataForSEO: [probes/dataforseo_canary5.py](probes/dataforseo_canary5.py),
  SHA-256 `DCA120B8AE5C35CBDC2BA1385F935A7E260194109BB4D54A9E0CA22C207A6F17`; self-test/dry-run/compilação `PASS`.
- Probe DataForSEO com Product Info + `pvf` + gate por oferta:
  [probes/dataforseo_product_details_canary5.py](probes/dataforseo_product_details_canary5.py),
  SHA-256 `B02DFF755C12B3FA5A3FFAF0A9B4A69C4D8D222C061A852AB8DA0FB0EB0DEA3A`;
  self-test/dry-run/compilação `PASS`.
- Auditoria Google Shopping APIs: [evidence/2026-07-13-google-shopping-api-challengers-audit.json](evidence/2026-07-13-google-shopping-api-challengers-audit.json),
  SHA-256 `7AE2D475AD40DC0C1C494F6DAE16507984CE11AEFA669F0FE30EE91D4A28707A`.
- Probe SearchAPI.io: [probes/searchapi_canary5.py](probes/searchapi_canary5.py),
  SHA-256 `CC7536B0631B892490C758387087622EDCAD7D5802B5992B1C164910F4D370E2`; self-test/dry-run/compilação/live `PASS`.
- Resultado SearchAPI EAN/Offers adicional: [evidence/2026-07-13-searchapi-canary5-result.json](evidence/2026-07-13-searchapi-canary5-result.json),
  SHA-256 `28C331784FE18DC4031332658C31A74D07D03D5F9E676CEAAFBAFDEECA6664D7`; 10 Shopping 200,
  9 Offers 200, 0 `ACCEPT`, 4 ofertas ML sem identidade verificável.
- Probe SearchAPI.io com detalhe + gate por oferta:
  [probes/searchapi_product_details_canary5.py](probes/searchapi_product_details_canary5.py),
  SHA-256 `5FBF13698108964BA1B7DE2C9320980550166681242D49C30ADB0B37E1B94E88`; self-test/dry-run/compilação `PASS`.
- Adjudicação SearchAPI live:
  [evidence/2026-07-13-searchapi-canary5-adjudication.json](evidence/2026-07-13-searchapi-canary5-adjudication.json),
  SHA-256 `F1AAAD33D23930AFA6CCB84E275A99AD42600011A65D3EFCCC5DD53D8AA97130`; `FAIL`.
- SearchAPI EAN/Offers SHA-256: `B90979EA478F2275C1846AE342BF17E7C665F094586FC0A786E4A463FA9EA372`.
- SearchAPI Details v1 SHA-256: `5807C1AA790C2AD0648060401ED109EB00461696269120944F538C2E594D5A1B`;
  1 falso automático em 2 aceitações aparentes.
- SearchAPI Details v2 SHA-256: `64ACA78902C7F7D20BAFD7E4D97AB2F0969FC1C635397CD4E996F37889CD5DEE`;
  0 falso automático, 0/5 preços aceitos.
- Auditoria de infraestrutura profissional de scraping:
  [evidence/2026-07-13-professional-scraping-providers-audit.json](evidence/2026-07-13-professional-scraping-providers-audit.json),
  SHA-256 `5A5AA9460819379349D9D27D1908AA53CE286363E33CEAB19DA1507C62DF8F82`.
- Auditoria das rotas oficiais vigentes:
  [evidence/2026-07-13-ml-official-current-route-audit.json](evidence/2026-07-13-ml-official-current-route-audit.json),
  SHA-256 `7FCB71817B272918DD9F6577A785621790A0C71E23EF66F247ACB0FA2E4DE46E`.
- Batch50 de busca externa: [evidence/2026-07-13-external-search-batch50.json](evidence/2026-07-13-external-search-batch50.json),
  SHA-256 `E6F522DEC07E3347579360BACBA2B6DC5328C9681A0B319CAEC72BFB9B9C608C`.
- Auditoria JoomPulse/Nubimetrics/parceiros ML:
  [evidence/2026-07-13-ml-partners-prelisting-audit.json](evidence/2026-07-13-ml-partners-prelisting-audit.json),
  SHA-256 `E713080A6E454E796C7BF44643CB42A1BFB3F080A49A664A587FD3D90D99989B`.
- Canário JoomPulse preparado, não enviado:
  [evidence/2026-07-13-joompulse-canary5-prepared.json](evidence/2026-07-13-joompulse-canary5-prepared.json),
  PDF SHA-256 `4BB05DD120185F8CF1135FD33BC42754AD3A61B5B784022FBF61E491B039688F`.
