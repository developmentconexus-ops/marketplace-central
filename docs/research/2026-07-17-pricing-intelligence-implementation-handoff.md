# Handoff de implementação — inteligência de preços do Mercado Livre

> **Escopo:** Mercado Livre Brasil, preço sem frete.  
> **Autoridade detalhada:** [batch pré-anúncio de 50 produtos](2026-07-12-mercado-livre-prelisting-batch50.md), [monitoramento de anúncios ativos](2026-07-12-mercado-livre-competitor-price-monitoring.md) e [protocolo de homologação](2026-07-13-mercado-livre-provider-homologation.md).  
> **Estado em 17/07/2026:** pode-se implementar com segurança o monitoramento de anúncios próprios ativos e o núcleo de identidade. Não existe fonte profissional homologada para precificação automática de produtos ainda não anunciados.

## 1. A decisão central

Há duas features diferentes. Elas não devem ser fundidas.

| Feature | Pergunta respondida | Estado | Fonte atual |
|---|---|---|---|
| **Monitoramento pós-anúncio** | “Meu anúncio ativo está competitivo; qual é o preço vencedor e o alvo?” | **GO para implementação** | API oficial ML: item próprio, `sale_price` e `price_to_win` |
| **Pesquisa pré-anúncio** | “Este produto do ERP vale anunciar; qual faixa o mercado pratica?” | **GO apenas para identificação e `REVIEW`** | API oficial para catálogo/identidade; preço automático depende de fornecedor homologado |

Não apresentar a segunda feature como “preço de mercado automático” enquanto não houver uma fonte de ofertas com licença OEM e matching aprovado. A resposta correta pode ser `NO_PRICE_EVIDENCE` ou `REVIEW`.

## 2. O que foi provado

### 2.1 Monitoramento de anúncios próprios

Probe live, somente leitura, em 10 anúncios ativos:

- `GET /items/{our_item}`: item próprio e vínculo de catálogo;
- `GET /items/{our_item}/sale_price`: 10/10 com sucesso;
- `GET /items/{our_item}/price_to_win?version=v2`: 10/10 com sucesso;
- 9/10 anúncios estavam vinculados ao catálogo; 3 estavam `winning`, 6 `competing` e 1 retornou `not_listed` bruto;
- em dois catálogos, nosso preço era menor que o preço vencedor e ainda assim não vencíamos. Frete, Full, parcelamento e reputação importam.

**Conclusão:** `price_to_win` é o sinal primário para anúncios já existentes. Não equivale a “menor preço” nem promete vitória; é o alvo competitivo calculado para as condições atuais.

### 2.2 Distribuição de ofertas de catálogo

`GET /products/{catalog_product_id}/items` retornou 236 de 256 ofertas declaradas em 9 catálogos, incluindo vendedores e preços. Isso permite derivar distribuição de preços, mas tem duas limitações:

1. a paginação precisa ser completa; um catálogo tinha 120 ofertas e somente 100 foram materializadas;
2. a documentação oficial em inglês marca a rota para desligamento, embora a documentação em português ainda a mostre e o probe tenha funcionado.

**Decisão:** implementar somente atrás de feature flag/adaptador. Nunca fazer dela a única fundação do produto.

### 2.3 Descoberta pré-anúncio

No batch real de 50 produtos do ERP, `GET /products/search?product_identifier={EAN}` retornou 22 produtos de catálogo; 18 chegaram a ter preços pela rota transitória de ofertas. A auditoria manual dos 22 encontrou:

- 12 correspondências aceitáveis;
- 8 que exigiam revisão;
- 2 rejeições claras por colisão de EAN.

Exemplos de colisão: acessório Doka retornando furadeira Menegotti e gancho Doka retornando lanterna traseira Volkswagen. Portanto EAN é sinal forte, mas **não é verdade suficiente**.

O substituto oficial mais durável, `GET /products/{id}` com `buy_box_winner` e `buy_box_winner_price_range`, respondeu 200 em 22/22 produtos, mas ambos os campos vieram nulos em 22/22. Ele deve ser consumido quando existir, sem inventar zero quando vier nulo.

## 3. O que não funcionou — decisões fechadas

| Opção | Teste/evidência | Veredito |
|---|---|---|
| Pricefy | 5 produtos; 14 sugestões; 0 Mercado Livre e 0 preços; domínios irrelevantes | **NO-GO** |
| GeckoAPI | 50 consultas; 49 HTTP 200; 0 EAN verificável e 0 preço aceito | **NO-GO** |
| SearchAPI.io | três canários reais: EAN trouxe biscoito/panos/fita/esparadrapo; versão semântica cruzou variante Crystal → Linea; correção ficou 0/5 com preço seguro | **NO-GO** |
| Busca pública, Playwright, Crawlee, Scrapy, proxies/stealth | técnica instável e/ou proibida pelos termos do ML | **Não implementar** |
| Bright Data, Oxylabs, Zyte, Apify | infraestrutura técnica possível, mas não transfere autorização do ML nem concede OEM padrão | **NO-GO sem aditivo ML + OEM** |
| JoomPulse | bom scanner de catálogo e benchmark manual; termos padrão são uso interno e proíbem redistribuição/API derivada | **Não integrar sem contrato especial** |
| Cargoos, Nubimetrics, Shopping de Preços | úteis para benchmark/manual, mas sem API pré-anúncio OEM comprovada | **Não usar como backend** |

`/sites/MLB/search`, `/items/{competitor}` e `/items/{competitor}/sale_price` retornaram 403 com a aplicação atual. Isso pode ser permissão funcional da app, não prova de proibição global; não colocar a feature em dependência deles.

## 4. Fonte de verdade de identidade

Separar os campos. O mapeamento atual tem defeito estrutural e precisa ser corrigido antes de ligar `seller_sku` ao ERP.

| Significado | Fonte Sankhya atual | Regra |
|---|---|---|
| SKU interno canônico | `CODPROD` | futuro `seller_sku`; único |
| EAN/GTIN | `REFERENCIA` | validar checksum, unicidade e categoria |
| Referência do fabricante | `REFFORN` | não confundir com SKU/EAN |
| Marca, linha, modelo e variante | atributos/descrição normalizados | obrigatórios para confirmar match |

No audit de 10.519 produtos ativos: 7.163 tinham GTIN válido e único (68,10%); havia 91 GTINs válidos repetidos em 188 produtos. `REFFORN` estava preenchido em 98,16% e marca + referência era única em 10.120 casos. A combinação marca + referência tende a ser mais útil que EAN isolado quando há colisão.

## 5. Contrato de matching obrigatório

### Decisões possíveis

- `ACCEPT`: identidade comprovada e sem contradição;
- `REVIEW`: há candidato, mas falta prova suficiente;
- `REJECT`: há contradição determinante;
- `NO_CANDIDATE`: nenhuma evidência de produto/oferta;
- `NO_PRICE_EVIDENCE`: produto identificado, porém sem preço válido;
- `INSUFFICIENT_MARKET`: preço existe, mas menos de cinco sellers válidos.

### Regras

1. Normalizar GTIN e validar checksum.
2. Quando a fonte devolver GTIN, exigir igualdade exata.
3. Validar domínio/família/categoria e marca.
4. Exigir ao menos duas âncoras independentes entre referência/modelo, linha e atributos.
5. Comparar variante: medida, capacidade, cor, acabamento, voltagem, kit, unidade e condição.
6. Contradição determinante vence EAN e produz `REJECT`.
7. Fuzzy, embedding e título servem para recuperar/rankear candidatos; nunca para autoaceitar sem as regras anteriores.

Hard negatives: kit/multipack, quantidade divergente, cor/acabamento, medida/capacidade, modelo, categoria, condição, GTIN duplicado e identidades conflitantes.

Para afirmar precisão automática de 99% com 95% de confiança, a referência usada é 368/368 acertos independentes por classe automática. Até existir gold set desse tamanho, manter casos limítrofes em `REVIEW`.

## 6. Modelo de dados recomendado

Não guardar somente `market_price`. Persistir origem e evidência.

```text
ProductIdentity
  internal_sku (CODPROD)
  gtin
  manufacturer_reference
  brand, title, category, variant_attributes

MarketPriceSnapshot
  product_identity_id
  source = ML_OFFICIAL | PROVIDER
  source_catalog_product_id
  source_item_id
  source_url
  match_status, match_reasons[], contradictions[]
  observed_at, fetched_at, expires_at
  raw_status, provider_request_id

CompetitiveSignal
  our_item_id
  our_sale_price
  winner_item_id, winner_price
  competitive_target_price
  competition_status, boosts
  observed_at

ValidatedOffer
  snapshot_id
  seller_key, item_id, url
  price, original_price, currency, condition
  identity_evidence, match_status, captured_at

MarketPriceAggregate
  snapshot_id
  n_offers, n_sellers
  min, p25, median, trimmed_mean, p75, max
  market_status
```

Nunca sobrescrever o último snapshot válido com zero, preço antigo sem marcação ou falha de rede. Registrar `fetched_at`, erro, expiração e a fonte.

## 7. Semântica dos preços exibidos

| Campo de UI/API | Significado | Não confundir com |
|---|---|---|
| `our_sale_price` | nosso preço efetivo no ML | preço de tabela ERP |
| `winner_price` | preço do anúncio vencedor no catálogo | menor preço absoluto |
| `competitive_target_price` | `price_to_win` calculado pelo ML | “vencedor - R$ 0,01” ou garantia |
| `catalog_min/median/max` | distribuição de ofertas válidas | recomendação de preço isolada |
| `visual_search_card_price` | observação manual no contexto visual | dado oficial de integração |

Para agregados: aceitar somente BRL, preço positivo, condição `new`, timestamp conhecido e identidade `ACCEPT`. Deduplicar por seller, usando a menor oferta válida por seller. Mostrar também `n_offers`, `n_sellers`, idade da coleta e status da amostra.

## 8. Arquitetura que deve ser implementada agora

```text
ERP identity (CODPROD + EAN + REFFORN + atributos)
  -> Identity Resolver / gate determinístico
  -> MercadoLivreOfficialAdapter
       -> anúncios próprios: item + sale_price + price_to_win
       -> catálogo: products/search -> products/{id}
       -> ofertas: products/{id}/items somente feature flag
  -> MarketPriceAggregator
  -> snapshots/evidências/auditoria
  -> UI: sinais separados + confiança + última coleta
```

Separar adapters de provedores externos do domínio. A resposta bruta de DataForSEO, Precifica ou outro fornecedor não deve vazar para a camada de negócio; o provider só entrega candidatos/ofertas normalizados para o mesmo gate interno.

### Fluxo pós-anúncio — implementar primeiro

1. Listar itens próprios ativos.
2. Resolver a identidade do item para `CODPROD` quando a escrita/readback de `seller_sku` estiver corrigida.
3. Buscar item, `sale_price` e `price_to_win`.
4. Salvar `CompetitiveSignal` e exibir preço próprio, estado, vencedor e alvo.
5. Se houver catálogo, coletar ofertas transitórias somente com feature flag e paginação completa.

### Fluxo pré-anúncio — implementar sem prometer preço

1. Buscar catálogo oficial por EAN; fallback por referência/part number/título; depois domínio + atributos.
2. Aplicar gate de identidade.
3. Consultar `products/{id}` e consumir winner/faixa somente se vierem preenchidos.
4. Se não houver preço, retornar `NO_PRICE_EVIDENCE` ou `REVIEW`.
5. Não chamar HTML, busca pública, proxy ou scraper como fallback silencioso.

## 9. Provedor externo: contrato antes do código

O único challenger programático que ainda merece canário é **DataForSEO Merchant API**, pois possui batch Products/Sellers. Ainda não foi executado com dados reais e não está aprovado.

A alternativa comercial brasileira mais alinhada é **Precifica**. Antes de qualquer integração, exigir demonstração dos cinco produtos congelados e resposta escrita sobre:

- origem/autorização dos dados ML;
- API usada na demonstração e em produção;
- EAN, catálogo, item/URL, seller, preço, condição, timestamp e atributos devolvidos;
- DPA/LGPD, subprocessadores, retenção, SLA, rate limits e política de falha;
- direito de armazenar, derivar mínimo/mediana/faixa e exibir no Marketplace Central multi-tenant;
- preço por SKU/consulta/oferta/tenant/atualização e responsabilidade pela licença dos dados.

Qualquer cláusula de “uso interno”, “sem sublicença”, “sem service bureau” ou que transfira para nós a violação de termos da fonte é `NO-GO` sem aditivo.

## 10. Gates antes de ligar pré-anúncio automático

### Canário fixo de 5

Não alterar a coorte após ver respostas: catálogo conhecido, colisão EAN conhecida, produto sem match oficial, variante sensível e cauda longa.

Passa somente se: pelo menos 3/5 matches corretos, zero falso `ACCEPT`, preços e URLs rastreáveis pelo mesmo canal contratado e ganho incremental sobre a API oficial.

### Batch de 50

Executar somente após canário. Medir separadamente cobertura bruta, identidade confirmada, preço disponível, `n >= 5 sellers`, custo, latência e adjudicação humana cega. Sem retries que escondam falhas. Um lote de 50 sem erros não prova 99% de precisão; continua sendo validação inicial.

## 11. Ordem de implementação sugerida

1. Corrigir contrato de identidade: `CODPROD`, EAN e referência de fabricante separados; gravar/readback de `seller_sku=CODPROD`.
2. Criar modelos de snapshot, oferta validada, agregados e auditoria.
3. Implementar adapter oficial de monitoramento pós-anúncio e UI de sinais separados.
4. Implementar `Identity Resolver` determinístico, testes das colisões e testes de variantes/kit.
5. Implementar descoberta pré-anúncio oficial retornando somente `ACCEPT/REVIEW/NO_PRICE_EVIDENCE`; nenhuma promessa de faixa quando não há evidência.
6. Colocar `/products/{id}/items` atrás de feature flag, paginação, telemetria e fallback explícito de indisponibilidade.
7. Somente após contrato/canário aprovado, adicionar `ProviderAdapter` externo e reexecutar batch50.

## 12. Evidências e artefatos reutilizáveis

- [Monitoramento de anúncios ativos](2026-07-12-mercado-livre-competitor-price-monitoring.md)
- [Batch pré-anúncio de 50 produtos](2026-07-12-mercado-livre-prelisting-batch50.md)
- [Protocolo de homologação](2026-07-13-mercado-livre-provider-homologation.md)
- `docs/research/evidence/2026-07-12-ml-live-api-v7-summary.json`: probe de monitoramento ativo.
- `docs/research/evidence/2026-07-12-ml-prelisting-batch50-summary.json`: lote oficial e adjudicação.
- `docs/research/evidence/2026-07-13-searchapi-canary5-adjudication.json`: prova de falso match e falta de repetibilidade.
- `docs/research/probes/*`: probes reproduzíveis; preservar como evidência, não executar contra novos provedores sem autorização por escrito.

## 13. Frase correta para produto e vendas

**Agora:** “Monitoramos o posicionamento competitivo e o preço-alvo dos seus anúncios ativos no Mercado Livre, com sinais oficiais e evidência de coleta.”

**Depois de homologar fornecedor:** “Para produtos do ERP ainda não anunciados, mostramos uma estimativa baseada em ofertas validadas, com tamanho de amostra, fonte e nível de confiança.”

Não prometer “menor preço”, “preço vencedor” ou “previsão de rentabilidade” para produto pré-anúncio sem amostra validada e módulo de custos separado.
