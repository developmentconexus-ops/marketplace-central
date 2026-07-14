# Protocolo de homologação — preços pré-anúncio no Mercado Livre

**Estado:** congelado para próximos fornecedores  
**Escopo:** descoberta de preços antes de existir anúncio próprio; preço sem frete  
**Objetivo:** impedir que cobertura aparente, EAN colidido ou termos comerciais inadequados virem
uma integração errada.

## 1. Canário fixo

Usar os cinco registros autorizados em
[`evidence/2026-07-12-pricefy-canary5-input.json`](evidence/2026-07-12-pricefy-canary5-input.json),
sempre com IDs externos sintéticos `PF-CANARY-1..5`:

1. produto de catálogo conhecido e semanticamente compatível;
2. colisão conhecida: EAN aponta para produto de domínio incompatível;
3. produto sem match oficial, com marca e referência fortes;
4. produto sensível a variante/acabamento;
5. produto de cauda longa.

O conjunto não pode ser trocado depois de ver a resposta do fornecedor. Ajustar o teste à saída seria
data leakage e invalidaria a comparação.

## 2. Gate jurídico/comercial antes de batch

O fornecedor precisa responder por escrito:

- origem e autorização de coleta dos dados do Mercado Livre;
- entidade contratante, DPA/LGPD, subprocessadores, localização e retenção;
- SLA, frequência, rate limits, suporte e tratamento de falha;
- direito de uso em SaaS multi-tenant/OEM;
- direito de exibir, armazenar temporariamente e derivar mínimo/mediana/faixa;
- TTL/cache e obrigação de atribuição/remoção;
- indenização ou alocação de responsabilidade por direitos da fonte;
- preço por SKU, consulta, oferta, tenant e atualização.

Termos “uso interno”, “não sublicenciável”, “sem service bureau” ou responsabilidade integral do
cliente pela violação da fonte são `NO-GO`, salvo aditivo específico.

## 3. Contrato técnico mínimo de resposta

Cada candidato/oferta deve preservar:

```json
{
  "synthetic_id": "PF-CANARY-1",
  "requested_gtin": "...",
  "source": "mercadolivre.com.br",
  "catalog_product_id": "MLB...",
  "item_id": "MLB...",
  "url": "https://...",
  "observed_gtin": "...",
  "title": "...",
  "brand": "...",
  "manufacturer_reference": "...",
  "variant_attributes": {},
  "condition": "new",
  "price": 0.0,
  "original_price": 0.0,
  "currency": "BRL",
  "seller_key": "stable-or-anonymized",
  "captured_at": "ISO-8601",
  "match_method": "gtin|catalog|semantic",
  "provider_confidence": 0.0
}
```

GTIN ausente não pode virar auto-aceite. Pode seguir para `REVIEW` se item/catálogo fornecer
atributos suficientes para validação independente.

## 4. Gate de identidade

Ordem obrigatória:

1. normalizar e validar checksum do GTIN;
2. exigir igualdade exata quando o fornecedor devolver GTIN;
3. validar domínio/família do produto;
4. validar marca;
5. exigir pelo menos duas âncoras independentes entre referência/modelo, linha e atributos;
6. comparar variantes determinantes: medida, capacidade, cor, acabamento, voltagem, kit e unidade;
7. qualquer contradição determinante vence o EAN e produz `REJECT`;
8. ausência de prova produz `REVIEW`, nunca `ACCEPT`.

Decisões: `ACCEPT`, `REVIEW`, `REJECT`, `NO_CANDIDATE`. O fornecedor não decide sozinho; o gate é
nosso e deve registrar razões positivas e negativas.

## 5. Gate de preço

Somente ofertas `condition=new`, BRL, preço positivo e timestamp conhecido. Deduplicar anúncios pelo
seller e usar a menor oferta válida de cada seller. Produzir:

- `n_offers`, `n_sellers` e idade da coleta;
- mínimo, P25, mediana, média aparada, P75 e máximo;
- fonte e nível de confiança;
- `INSUFFICIENT_MARKET` quando houver menos de cinco sellers;
- `NO_PRICE_EVIDENCE` quando não houver oferta validada.

Mínimo isolado não é “preço recomendado”. A feature futura deve mostrar amostra e dispersão.

## 6. Critérios objetivos

### Entrada no batch de 50

- zero falso auto-aceite nos cinco;
- pelo menos três matches corretos;
- pelo menos um ganho incremental sobre a API oficial;
- preços acessíveis pelo mesmo canal/API que seria contratado;
- latência e passos manuais registrados;
- gate jurídico/comercial sem bloqueio.

### Batch de 50

- exatamente os 50 GTINs congelados; sem retry que masque falha;
- cobertura bruta, identidade confirmada, preço e `n>=5 sellers` separados;
- zero falso auto-aceite; colisões conhecidas obrigatoriamente rejeitadas;
- comparação com API oficial e adjudicação humana cega;
- custo e latência medidos por produto.

Zero erro em 50 casos não prova 99% de precisão: o limite inferior bilateral de 95% é cerca de 92,9%.
Para alegar erro abaixo de 1% com zero falhas e 95% unilateral, a regra de três exige cerca de 300
aceites independentes. Até lá, matches fora do gate conservador permanecem em revisão.

## 7. Ranking atual

1. **API oficial + matcher próprio:** núcleo autorizado e provado; cobertura conservadora limitada e
   lista de ofertas sem substituto documentado.
2. **DataForSEO Merchant API:** melhor challenger programável; batch de produtos+sellers e uso por
   software builders confirmados, mas cobertura ML Brasil por EAN e direito específico ainda não
   passaram canário.
3. **JoomPulse Scanner + MCP:** melhor aderência funcional ao catálogo ERP; scanner batch, XLS e MCP
   OAuth confirmados, mas matching EAN exato e licença OEM não; termos padrão são `NO-GO`.
4. **Precifica:** melhor challenger enterprise brasileiro; canário técnico e licença OEM pendentes.
5. **Cargoos dashboard:** batch EAN server-side real e canário manual válido; sem API pública,
   distribuição min/média/mediana, SLA/DPA ou OEM.
6. **Bright Data Partner / Oxylabs Custom Dataset:** somente contrato especial com autorização da
   fonte, GTIN verificável e derived-data/OEM explícitos.
7. **Nubimetrics:** bom benchmark manual por keyword, preço médio e vendas; sem EAN batch/API/OEM e
   com automação proibida nos termos padrão.
8. **Shopping de Preços:** parceiro Platinum, mas a busca fica restrita às categorias elegíveis pelo
   histórico da conta e a API pública não expõe inteligência pré-listagem.

SearchAPI.io, Pricefy, GeckoAPI, Rally, ShopAPIS, Pricenly, scrapers próprios/comunitários e APIs de supermercado
estão reprovados pelas evidências do relatório principal.

### Canário DataForSEO (preparado, não executado)

- Products em batch com EAN e marca+referência+título, `location=Brazil`, `language=pt`;
- selecionar até dois candidatos por produto e chamar Product Info antes de Sellers;
- preservar juntos `product_id`, `data_docid`, `gid` e o filtro de variante `pvf` em Sellers;
- reter apenas oferta cujo domínio/URL final seja `mercadolivre.com.br` e revalidar a própria oferta;
- exigir marca + referência ou marca + ≥75% de título/linha/variante, sem contradição de tensão/kit;
- nunca herdar confiança do Product Info para o preço final;
- medir cobertura, falso aceite, latência e custo; repetir os dois EANs de colisão conhecida;
- `GO`: pelo menos 3/5 produtos corretos, zero falso automático, preço e URL/MLB rastreáveis;
- antes da produção, confirmar por escrito armazenamento, derivados e exibição multi-tenant.

### Canário SearchAPI.io (executado — `FAIL`)

- EAN→Offers: 0/5 preços aceitos; resultados Mercado Livre semanticamente errados;
- Product Details→Offers v1: duas aceitações automáticas, mas uma falsa (`Crystal/FCR925` recebeu
  ofertas `Linea/16702`), violando o gate de zero falso;
- gate corrigido e repetição v2: 0 falso, mas 0/5 preço aceito e oferta correta anterior ausente;
- `FAIL` em cobertura, falso aceite e repeatabilidade; não executar batch50 nem contratar;
- adjudicação estruturada em
  `evidence/2026-07-13-searchapi-canary5-adjudication.json`.

### Canário JoomPulse (não executado)

- artefato offline pronto em `output/pdf/joompulse-canary5-catalog.pdf`, SHA-256
  `4BB05DD120185F8CF1135FD33BC42754AD3A61B5B784022FBF61E491B039688F`; 1 página, 5 IDs sintéticos,
  5 EANs e zero SKU ERP; não enviado e sem autorização de transferência para a JoomPulse;
- usar o PDF sintético preparado com os mesmos cinco IDs externos, EAN, marca, referência e
  descrição autorizados; nunca incluir SKU ERP, custo ou informação do fornecedor;
- enviar pelo Scanner de Catálogos e, separadamente, consultar as mesmas cinco linhas pelo MCP OAuth;
- registrar se o EAN é preservado, candidato MLB/MLBU, URL, título/atributos, preço atual, mínimo
  histórico, preço médio, sellers, vendas estimadas e export XLS;
- tratar “produto similar” como candidato, nunca match exato; aplicar integralmente o gate de
  identidade e exigir que a colisão conhecida permaneça `REJECT`;
- comparar cobertura e preços contra os sete casos oficiais já validados;
- `GO` técnico: pelo menos 3/5 corretos, zero falso automático e preços rastreáveis;
- mesmo com `GO`, não integrar nem reproduzir outputs no Marketplace Central sem aditivo OEM escrito.

## 8. Pedido de demonstração — Precifica (rascunho, não enviado)

> Temos cinco produtos reais ainda não anunciados no Mercado Livre, identificados por EAN/GTIN,
> marca, referência e atributos. Precisamos que a demonstração use o mesmo fluxo e a mesma API que
> seriam contratados. Para cada produto, devolvam candidatos do Mercado Livre Brasil com GTIN
> observado ou evidência equivalente, MLB/MLBU, URL, título, marca, referência/atributos, preço atual,
> preço anterior, condição, seller anonimizado, timestamp e todas as ofertas disponíveis. Um dos
> cinco EANs possui colisão conhecida e deve ser rejeitado semanticamente. Também precisamos do
> schema/request-response, limites, frequência, SLA, origem autorizada e cláusula que permita exibir e
> derivar agregados em nosso SaaS multi-tenant. Após aprovação do canário, repetiremos lote fixo de 50.

Não enviar CNPJ, contato pessoal, arquivos ou os cinco produtos a um novo fornecedor sem confirmação
no momento da ação.

## 9. Perguntas ao suporte Mercado Livre (rascunho, não enviado)

1. Existe endpoint suportado que aceite GTIN/catalog product/atributos e devolva antes do item os
   mesmos `suggested_price`, `lowest_price`, `internal_price` e `metadata.graph` do fluxo web?
2. `/suggestions/items/{item_id}/details` alimenta o modal de criação? Em qual estado o item já existe?
3. É suportado para item pausado ou sem estoque? É permitido criar item exclusivamente para simulação?
4. Qual substituto de `/products/{product_id}/items` fornece ofertas/contagem/faixa pré-anúncio?
5. Como interpretar `buy_box_winner_price_range` e por que pode ser nulo com ofertas ativas?
6. Podemos exibir e persistir agregados para sellers OAuth? Quais TTL, atribuição e limites multi-tenant?
7. Há data partner certificado e redistribuível para esse caso?
8. Quais quotas de `products/search`, `/products/{id}`, `suggestions` e batch/multiget?
9. Como reportar colisões falsas de GTIN no catálogo?
