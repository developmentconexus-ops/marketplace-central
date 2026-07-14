# Monitoramento automatizado de preços concorrentes no Mercado Livre

**Data:** 2026-07-12  
**Workspace:** `marketplace-central`  
**Status:** pesquisa live concluída; decisão arquitetural ainda condicionada aos gates deste relatório  
**Escopo:** Mercado Livre Brasil, preço sem frete, identidade de produto e coleta automatizada  
**Fora de escopo:** lucro, margem, impostos, alteração de preço, publicação, UI e implementação produtiva

> **Correção de escopo:** este relatório validou principalmente concorrência pós-anúncio. O objetivo
> pré-anúncio foi testado depois em 50 produtos reais do ERP. O resultado e o novo veredito estão em
> [Descoberta pré-anúncio — batch real de 50 produtos](2026-07-12-mercado-livre-prelisting-batch50.md).
> Para seleção do que vale anunciar, esse documento posterior é a autoridade.

## 1. Veredito executivo

A solução líder para os produtos que já anunciamos é **API oficial, com
`price_to_win` como sinal primário**. Não devemos construir um scraper HTML.

O fluxo comprovado nos dez anúncios ativos é:

1. localizar nossos anúncios;
2. consultar `/items/{our_item}` e `/items/{our_item}/sale_price`;
3. consultar `/items/{our_item}/price_to_win?version=v2`;
4. guardar status competitivo, preço vencedor, item vencedor, `price_to_win` e boosts;
5. usar identidade de catálogo/GTIN para impedir comparação entre produtos diferentes.

Isso responde melhor à pergunta comercial “qual preço me torna competitivo?” do que simplesmente
pegar o menor valor da busca. O teste real provou que **preço sozinho não determina o vencedor**:
em dois catálogos nosso preço era menor que o preço vencedor e mesmo assim estávamos perdendo.
Frete, Full, parcelamento e reputação participam do algoritmo, como documentado pelo
[Mercado Livre](https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/concorrencia-em-catalogo).

Há três limitações vinculantes:

- `/items/{competitor}` e `/sale_price` de concorrentes retornaram `403` com a aplicação atual;
- `/sites/MLB/search` retornou `403`;
- `/products/{catalog_product_id}/items` funcionou e trouxe ofertas de múltiplos sellers; oito dos
  nove catálogos ficaram completos e um ficou limitado às primeiras 100 de 120 ofertas. Além disso, existe
  [aviso oficial en-US de desligamento](https://developers.mercadolivre.com.br/en_us/catalog-competition),
  enquanto a documentação pt-BR ainda o publica. É uma superfície transitória, não fundação de longo prazo.

Portanto:

- **comprovado hoje:** preço próprio, vencedor concorrente, alvo competitivo e snapshot de ofertas de catálogo;
- **não comprovado:** paridade sistemática com o preço promocional/Pix exibido no card da busca;
- **sem solução automática comprovada:** descoberta de concorrentes para produto fora do catálogo.

Não planejar a integração até os gates da seção 12 passarem.

## 2. Escopo decidido pelo operador

- Somente Mercado Livre.
- Preço exibido sem frete.
- EAN/GTIN é o vínculo atual.
- No futuro, nosso `seller_sku` será o `CODPROD`.
- Lucro e custo serão tratados depois.
- Matching automático exige precisão conservadora mínima de 99%.
- Nenhum scraping evasivo, CAPTCHA solving, proxy rotation ou stealth.

## 3. Evidência live da API oficial

O probe foi executado somente com GET. O token foi resolvido pelo serviço cifrado existente e nunca
foi impresso. O resumo auditável está em
[evidence/2026-07-12-ml-live-api-v7-summary.json](evidence/2026-07-12-ml-live-api-v7-summary.json).

### 3.1 OAuth e estado real

Cronologia confirmada:

1. credencial v6: todos os oito GETs iniciais retornaram `401 invalid access token`;
2. ao iniciar o Postgres, o ticker renovou a credencial para v7;
3. sessão v7: `valid`, expiração canônica `2026-07-12T21:53:18Z`;
4. os probes v7 retornaram `200` nas superfícies próprias, de catálogo e competição;
5. dois `StartAuthorize` posteriores sobrescreveram a instalação para `pending_connection`, embora
   credencial e sessão continuassem válidas.

O `expires_at=null` emitido pelo probe antigo era leitura da fonte errada. A expiração vive em
`integration_auth_sessions.access_token_expires_at`, não no blob cifrado da credencial.

### 3.2 Cobertura dos dez anúncios ativos

| Métrica | Resultado live |
|---|---:|
| anúncios próprios ativos | 10 |
| ligados a catálogo | 9 (90%) |
| fora do catálogo | 1 (10%) |
| `winning` | 3 |
| `competing` | 6 |
| `not_listed` bruto | 1 |
| `/sale_price` próprio | 10/10 com `200` |
| `/price_to_win` | 10/10 com `200` |

`not_listed` foi preservado como valor bruto porque não aparece na lista de status documentados.

### 3.3 Preço próprio, vencedor e alvo

| Nosso item | Catálogo | Nosso preço | Estado | Preço vencedor | `price_to_win` |
|---|---|---:|---|---:|---:|
| `MLB4834219830` | — | R$ 29,90 | `not_listed` | — | — |
| `MLB4735328201` | `MLB22624877` | R$ 169,99 | `winning` | R$ 169,99 | R$ 169,00 |
| `MLB6896039640` | `MLB15531385` | R$ 189,90 | `competing` | R$ 207,81 | R$ 178,00 |
| `MLB6896003832` | `MLB32390517` | R$ 334,90 | `competing` | R$ 308,00 | R$ 319,00 |
| `MLB6896003262` | `MLB19858075` | R$ 729,90 | `competing` | R$ 698,00 | R$ 651,00 |
| `MLB6896001442` | `MLB23427630` | R$ 859,90 | `competing` | R$ 712,00 | R$ 664,00 |
| `MLB4735364085` | `MLB28262741` | R$ 299,90 | `winning` | R$ 299,90 | R$ 299,00 |
| `MLB4735326915` | `MLB60275101` | R$ 546,99 | `competing` | R$ 547,25 | R$ 510,00 |
| `MLB4735324525` | `MLB35928565` | R$ 69,90 | `winning` | R$ 69,90 | R$ 69,90 |
| `MLB4735304125` | `MLB27562650` | R$ 179,99 | `competing` | R$ 147,00 | R$ 157,00 |

`price_to_win` não significa “preço do vencedor menos um centavo” nem garantia imutável de vitória.
É o preço que, nas condições atuais, torna a publicação mais competitiva. Em `MLB32390517`, por
exemplo, o vencedor estava a R$ 308 e o alvo retornado foi R$ 319.

### 3.4 Lista de ofertas por PDP — superfície transitória

| Catálogo | Ofertas declaradas | Materializadas | Sellers observados | Faixa observada | Vencedor |
|---|---:|---:|---:|---:|---:|
| `MLB22624877` | 15 | 15 | 11 | R$ 169,99–307,86 | R$ 169,99 |
| `MLB15531385` | 120 | 100 | 83 | R$ 120,00–463,91 | R$ 207,81 |
| `MLB32390517` | 30 | 30 | 26 | R$ 150,00–723,94 | R$ 308,00 |
| `MLB19858075` | 26 | 26 | 22 | R$ 549,99–1.316,60 | R$ 698,00 |
| `MLB23427630` | 53 | 53 | 38 | R$ 712,00–2.450,86 | R$ 712,00 |
| `MLB28262741` | 6 | 6 | 5 | R$ 299,90–749,88 | R$ 299,90 |
| `MLB60275101` | 3 | 3 | 3 | R$ 546,99–725,00 | R$ 547,25 |
| `MLB35928565` | 1 | 1 | 1 | R$ 69,90 | R$ 69,90 |
| `MLB27562650` | 2 | 2 | 2 | R$ 147,00–179,99 | R$ 147,00 |

Totais: 256 publicações declaradas; 236 materializadas (`92,1875%`). O catálogo com 120 resultados
precisaria da segunda página. Não chamar essas 256 linhas de “concorrentes”: são publicações e podem
existir várias por seller.

Consistência interna observada:

- preço próprio de `/sale_price` = oferta própria da PDP: 9/9;
- `price_to_win.winner.item_id/price` = oferta vencedora da PDP: 9/9.

Essas verificações provam consistência entre superfícies da API, não paridade visual independente.

### 3.5 Testes negativos

| Superfície | Resultado |
|---|---|
| `/items/{competitor}` | `403` em 1/1 |
| `/items/{competitor}/sale_price` | `403` em 4/4 |
| `/sites/MLB/search` | `403` em 2/2 |
| `/products/{id}.buy_box_winner` | ausente em 9/9, apesar de `200` |
| busca própria `seller_sku=15956` | `200`, total 0 |

O `403` pode ser falta de permissão funcional da aplicação atual. A documentação explica que
`PolicyAgent`/`PA_UNAUTHORIZED_RESULT_FROM_POLICIES` exige revisar permissões no DevCenter:
[Permissões funcionais](https://developers.mercadolivre.com.br/pt_br/permissoes-funcionais/).
O teste não prova que o recurso seja globalmente proibido.

## 4. Definição correta dos preços

Guardar sinais diferentes; não colapsá-los em um único campo:

| Sinal | Fonte | Semântica |
|---|---|---|
| `our_sale_price` | `/items/{our_item}/sale_price` | nosso preço efetivo no canal |
| `winner_price` | `price_to_win.winner.price` | preço da publicação que vence a PDP |
| `competitive_target` | `price_to_win.price_to_win` | alvo sugerido para melhorar competitividade |
| `catalog_offer_price` | `/products/{catalog}/items[].price` | snapshot transitório de cada oferta |
| `catalog_min/median/max` | derivado das ofertas paginadas | distribuição, não recomendação isolada |
| `visual_search_card_price` | observação manual autorizada | preço exibido ao comprador no contexto definido |

Não usar:

- parcela como preço total;
- frete;
- zero para indisponível;
- valor antigo silencioso após falha;
- “menor oferta” como sinônimo de vencedor;
- `price_to_win` como promessa de vitória.

## 5. Identidade e matching nos dados reais

### 5.1 Defeito estrutural atual

O futuro `seller_sku` será `CODPROD`, mas o adapter atual interpreta `SellerSKU` como `REFFORN`:
[reader.go](../../apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go#L384).

Além disso, `ProductCandidate.EAN` e `ReferenceCode` recebem ambos `REFERENCIA`:
[reader.go](../../apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go#L48).

Caso real:

| Campo Sankhya | Valor | Semântica |
|---|---|---|
| `CODPROD` | `15956` | nosso SKU interno |
| `REFERENCIA` | `7894200146179` | EAN |
| `REFFORN` | `2020.C.FLX` | referência do fabricante |

No anúncio `MLB4735328201`, o `SELLER_SKU` live ainda é `2020.C.FLX`. Por isso a busca
`seller_sku=15956` retornar zero é o estado atual esperado, não falha da API. M-12 precisa gravar e
confirmar readback de `CODPROD`; M-09 precisa separar identidades canônicas.

### 5.2 Cobertura Oracle

Audit somente `SELECT`, 10.519 produtos ativos:

| Sinal | Resultado observado |
|---|---:|
| `CODPROD` único | 10.519 (100%) |
| `REFERENCIA` preenchida | 7.373 (70,09%) |
| GTIN válido por formato/checksum | 7.351 |
| GTIN válido e único | 7.163 (68,10%) |
| colisões de GTIN válido | 91 códigos / 188 produtos |
| `REFFORN` preenchida | 10.325 (98,16%) |
| `REFFORN` globalmente única | 9.996 |
| marca + `REFFORN` únicos | 10.120 |
| sem `REFERENCIA` e sem `REFFORN` | 186 (1,77%) |

Checksum é necessário, mas não suficiente. Existem GTINs válidos repetidos em produtos diferentes.

### 5.3 Política de matching

Anúncio próprio:

1. `seller_sku == CODPROD` somente após escrita/readback auditados;
2. item pertence à nossa instalação;
3. GTIN, catálogo e variação não contradizem o SKU.

Concorrente:

1. mesma PDP validada (`catalog_product_id`);
2. GTIN válido e único + categoria/variação compatíveis;
3. marca + referência do fabricante + atributos;
4. lexical/fuzzy/embedding apenas para recuperar candidatos.

Hard negatives: kit/multipack, quantidade, cor/acabamento, medida/capacidade, modelo, categoria,
condição, GTIN duplicado ou identidades conflitantes.

Meta estatística: auto-match somente quando o limite inferior bilateral exato de 95% da precisão for
`≥0,99`. Referência conservadora: `368/368` acertos por classe automática produz limite inferior
aproximado `0,990026`. Fuzzy continua manual até obter sua própria coorte.

## 6. Fluxo recomendado

```text
agora:
  listar anúncios próprios ativos
  -> /items/{our_item}
  -> /items/{our_item}/sale_price
  -> /items/{our_item}/price_to_win?version=v2
  -> vencedor + preço vencedor + alvo + boosts
  -> opcional/transitório: /products/{catalog_product_id}/items paginado

depois de M-12:
  CODPROD
  -> seller_sku do anúncio próprio
  -> resolver item
  -> mesmo fluxo acima
```

Fora do catálogo:

1. tentar `/products/search` por GTIN/part number;
2. validar produto encontrado por marca/modelo/categoria/variante;
3. se continuar sem catálogo, marcar `discovery_unproven` e enviar para revisão humana;
4. não cair silenciosamente para HTML automatizado.

## 7. Alternativas avaliadas

| Opção | Evidência | Decisão |
|---|---|---|
| `price_to_win` | `200` em 10/10; vencedor consistente em 9/9 catalogados | **Primário** |
| `/products/{catalog}/items` | live em 9/9; 236 ofertas observadas | Snapshot transitório, atrás de adapter/feature flag |
| `/sale_price` concorrente | `403` em 4/4 com a aplicação atual | Reavaliar permissões; não depender |
| `/sites/MLB/search` | `403` em 2/2 | Reavaliar permissões; não depender |
| `/products/{id}.buy_box_winner` | ausente em 9/9 | Não usar como fonte MLB comprovada |
| Scrapling HTTP | página 200, sem preço | Não escolhido |
| Playwright | extraiu preço duas vezes; depois `403` | Rejeitado sem autorização |
| Crawlee + Playwright | boa orquestração conceitual | Reserva somente com autorização escrita |
| Camoufox/stealth/proxies | foco em evasão | Rejeitado |
| fornecedor comercial | reduz manutenção | Não transfere autorização nem resolve matching |

## 8. Pesquisa de scraping e referências

PoCs executadas:

| PoC | Resultado observado |
|---|---|
| `requests` | item público `403 PolicyAgent`; páginas redirecionaram para verificação |
| Scrapling `0.4.10` Fetcher | HTTP 200 rápido, título presente, nenhum preço |
| Playwright `1.61.0` + Edge 150 | item `R$169,99` em texto, meta e JSON; ~4,75 s |
| segunda execução Playwright | mesmo preço e sinais |
| extração posterior de cards | HTTP `403`, zero cards |

Não houve retry evasivo, stealth, proxy, CAPTCHA bypass ou rotação de identidade.

Projetos consultados:

- [Scrapling 0.4.10](https://pypi.org/project/scrapling/0.4.10/): parser/fetcher testado;
- [Playwright](https://github.com/microsoft/playwright-python): único browser que extraiu o preço;
- [Crawlee](https://github.com/apify/crawlee): orquestrador conceitual se autorizado;
- [Scrapy](https://github.com/scrapy/scrapy) + [scrapy-playwright](https://github.com/scrapy-plugins/scrapy-playwright): alternativa Python não testada;
- [Camoufox](https://github.com/daijro/camoufox/releases/tag/v150.0.2-beta.25): rejeitado por foco anti-detect.

Fornecedores avaliados apenas como referência: Apify, Oxylabs, Bright Data e Zyte. Publicação de um
scraper por terceiro não prova autorização do Mercado Livre. A página de um Actor Apify apresentava
preços conflitantes; custo faturável ficou desconhecido.

## 9. Restrições contratuais

Fontes:

- [Termos do Programa de Desenvolvedores](https://developers.mercadolivre.com.br/pt_br/termos-e-condicoes);
- [Termos gerais do Mercado Livre](https://www.mercadolivre.com.br/ajuda/991).

Os termos consultados proíbem robots, spiders, scraping e tecnologias equivalentes sem autorização
expressa, além de contorno técnico e volume abusivo. Também há restrições de retenção, redistribuição
e usos derivados do conteúdo.

Consequência:

1. não automatizar HTML sem autorização escrita;
2. confirmar por escrito polling, histórico e agregados internos de concorrência;
3. não usar proxies, CAPTCHA solving ou stealth;
4. manter dados internos, tenant-scoped e não redistribuídos;
5. terceiro comercial não elimina a obrigação contratual.

Não é parecer jurídico.

## 10. Achados de confiabilidade OAuth

### 10.1 Estado divergente

O refresh v7 concluiu corretamente, mas dois `StartAuthorize` posteriores mudaram a instalação de
`connected` para `pending_connection`. O fluxo não aplica a matriz de transições do domínio e não
registra ator/idempotência da tentativa.

Um coletor arquiteturalmente correto pode parar porque capabilities projetam `needs_auth`; um probe
que usa apenas `CredentialResolver` pode continuar executando. Isso cria falso positivo de prontidão.

### 10.2 Gate correto de credencial

Antes de coletar preço, exigir em conjunto:

- instalação e snapshot `connected/healthy`;
- capability futura `competitor_price_read` executável;
- credencial ativa, não revogada e igual ao ponteiro da instalação;
- sessão `valid`, conta coerente e expiração posterior ao horizonte da coleta;
- access e refresh token presentes;
- conta alinhada entre instalação, sessão e payload.

A fonte de expiração é `AuthSession.AccessTokenExpiresAt`, não `payload["expires_at"]`.

O Mercado Livre documenta que o `refresh_token` é de uso único, somente o último token emitido
permanece válido e cada refresh devolve um novo par. Portanto, a rotação deve persistir o novo
`access_token` **e** o novo `refresh_token` na mesma transação antes de publicar a nova credencial
como ativa. Também deve existir exclusão mútua por instalação (ou compare-and-swap equivalente),
pois dois refreshes concorrentes podem consumir o mesmo token e deixar a aplicação apontando para
um par já inválido. Um `invalid_grant` não deve entrar em retry cego: após confirmar que não há uma
rotação concorrente mais nova, a instalação deve exigir nova autorização do seller. A atomicidade e
a exclusão mútua são requisitos arquiteturais inferidos da semântica oficial de uso único.

## 11. Confirmado, inferido e desconhecido

### Confirmado

- `price_to_win` funciona nos dez anúncios ativos.
- Nove anúncios possuem catálogo; seis perdem, três ganham.
- Preço próprio e vencedor foram consistentes entre superfícies oficiais em 9/9.
- `sale_price` cross-seller e busca pública estão bloqueados para a aplicação atual.
- A lista de ofertas funciona live, mas só 236/256 linhas foram materializadas por paginação incompleta.
- `buy_box_winner` esteve ausente nos nove detalhes de produto.
- `seller_sku=15956` ainda não existe no anúncio real.

### Inferido

- Falta de permissão funcional é uma causa plausível dos `403`, mas ainda não provada.
- `price_to_win` é o melhor sinal estável para a decisão de preço dos produtos que já vendemos.
- A lista de ofertas deve ser isolada atrás de um adapter substituível enquanto continuar disponível.

### Desconhecido

- Substituto oficial durável para listar todas as ofertas de uma PDP.
- Paridade centavo a centavo com preço visual/Pix em 30 casos.
- Permissão contratual para histórico e agregados internos.
- Quota real e resposta a 429.
- Descoberta automatizada fora do catálogo.
- Precisão cross-marketplace dos métodos sem GTIN.

## 12. Gates antes de planejar integração

### G1 — Prontidão OAuth operacional

- reparar/validar a transição `connected → pending_connection`;
- usar `AuthSession.AccessTokenExpiresAt` como fonte canônica;
- exigir a tríade instalação + sessão + credencial coerente;
- serializar refresh por instalação e persistir atomicamente o novo par access+refresh;
- em `invalid_grant`, reconciliar uma eventual rotação concorrente; se não existir par mais novo,
  interromper a coleta e solicitar reautorização, sem retry do token consumido;
- resultado esperado: `connected/healthy`, sessão válida e capability executável.

**Estado:** falhou operacionalmente, embora os GETs v7 tenham funcionado.

### G2 — Permissões e durabilidade da API

- auditar no DevCenter permissões funcionais de publicação/preços/busca;
- reautorizar somente se as permissões forem alteradas;
- repetir `sale_price` cross-seller e busca;
- obter resposta escrita sobre o substituto de `/products/{id}/items` e sua data efetiva de desligamento.

**Estado:** pendente.

### G3 — Autorização de uso

- confirmação escrita para polling, retenção e agregações internas;
- autorização separada antes de qualquer HTML automatizado.

**Estado:** pendente.

### G4 — Paridade visual

- coorte pré-definida de 30 pares elegíveis;
- contexto neutro: MLB Brasil, BRL, sem frete, login, loyalty ou cupom pessoal;
- API e observação manual em até cinco minutos;
- aprovação somente com `30/30` iguais ao centavo;
- divergência promocional conta como falha, não como explicação para passar.

**Estado:** 1/1 exploratório; gate não aprovado.

### G5 — Matching

- adjudicação por método e estrato;
- hard negatives explícitos;
- limite inferior bilateral CI95% `≥0,99` por classe automática;
- fuzzy/embedding manual até evidência própria.

**Estado:** pendente.

### G6 — Escala e quota

- medir headers/429 e paginação;
- começar com dez anúncios ativos, uma vez ao dia;
- aumentar frequência somente após prova de quota;
- respeitar `Retry-After`, jitter e backoff.

**Estado:** pendente.

### G7 — Segurança

- rotacionar a credencial Oracle exposta no artefato histórico M-04 F-02;
- redigir o artefato e avaliar limpeza de histórico;
- nunca registrar tokens ou buyer PII.

**Estado:** pendente.

## 13. Próximo passo recomendado

Ainda não planejar a feature.

1. corrigir/validar o estado OAuth e a capability operacional;
2. auditar permissões da aplicação e repetir os endpoints bloqueados;
3. obter resposta oficial sobre desligamento/substituto da lista de ofertas;
4. obter autorização de retenção/agregação;
5. executar paridade manual de 30 casos;
6. adjudicar matching por método;
7. somente então planejar a inteligência comercial.

## 14. Fontes primárias

- [Concorrência em catálogo — pt-BR](https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/concorrencia-em-catalogo)
- [Catalog competition — aviso en-US](https://developers.mercadolivre.com.br/en_us/catalog-competition)
- [API de preços](https://developers.mercadolivre.com.br/pt_br/servicos-gerenciamento-de-contatos/api-de-precos)
- [Referências de preços](https://developers.mercadolivre.com.br/pt_br/usuarios-e-aplicativos/referencias-de-precos)
- [Buscador de produtos](https://developers.mercadolivre.com.br/pt_br/buscador-de-produtos)
- [Itens e buscas](https://developers.mercadolivre.com.br/pt_br/itens-e-buscas)
- [Permissões funcionais](https://developers.mercadolivre.com.br/pt_br/permissoes-funcionais/)
- [Autenticação e autorização — refresh token de uso único](https://developers.mercadolivre.com.br/pt_br/gerenciamento-perguntas-respostas/autenticacao-e-autorizacao)
- [Gestão de identidades e acessos](https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/gestao-de-identidades-e-acessos-oauth-e-tokens)
- [Termos do Programa](https://developers.mercadolivre.com.br/pt_br/termos-e-condicoes)
- [Termos gerais](https://www.mercadolivre.com.br/ajuda/991)

## 15. Integridade

- Nenhum write de anúncio, preço ou Oracle foi executado; o ticker realizou um POST OAuth de refresh.
- Oracle foi consultado somente por `SELECT`.
- Probes de dados do Mercado Livre usaram somente GET.
- Tokens nunca foram impressos.
- Nenhum código produtivo foi implementado.
- Mudanças concorrentes fora de `docs/research/` foram preservadas.
- Evidência completa e limitações estão no manifesto e no resumo JSON v7.
- O Postgres local terminou parado; Docker reportou `Exited (137)` após o timeout do stop. Nenhum
  volume foi removido, e backend/frontend/ngrok não foram tocados.
