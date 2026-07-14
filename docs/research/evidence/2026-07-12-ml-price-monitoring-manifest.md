# Manifesto de evidência — preços concorrentes Mercado Livre

**Data:** 2026-07-12  
**Repositório autoritativo:** marketplace-central  
**Natureza:** pesquisa e probes; nenhuma integração produtiva

## 1. Classes

- **E1:** comando, artefato ou hash reexecutável e persistido.
- **E2:** saída redigida observada; requer repetição independente para gate formal.
- **E3:** documentação primária.
- **B:** tentativa executada que não provou o critério.

Resumo estruturado: [2026-07-12-ml-live-api-v7-summary.json](2026-07-12-ml-live-api-v7-summary.json).

## 2. Integridade

- Probes do Mercado Livre usaram somente GET; nenhum anúncio ou preço foi alterado.
- O ticker produtivo renovou OAuth v6 para v7; nenhuma credencial foi impressa.
- Oracle recebeu somente SELECT.
- Sem retry evasivo, proxy, CAPTCHA solving ou stealth.
- Buyer payloads, tokens e PII não foram persistidos.
- O Postgres terminou parado (`Exited 137` após timeout do stop); nenhum volume foi removido.

## 3. Artefatos e hashes

| Artefato | SHA-256 | Uso |
|---|---|---|
| C:\tmp\ml_price_baseline.py | 4306420A0BE5BB55F3856866F25B97A072A93699C2A91562E80E0E5B7F95BC66 | baseline HTTP |
| C:\tmp\ml_price_poc_frameworks.py | 2DACA1CBF150E24E49023614BAD14F75E6E1B3667334300F9A901E325D35A936 | Scrapling/Playwright |
| C:\tmp\ml_search_cards_poc.py | 3DE10BB4D5F7707CE1BC88EB8B87643ED853CA30FD22BC951C2C96715199FF4E | cards e primeiro 403 |
| probe Go inicial v6 | 1FCBA1563624EAA1CCE6490A47753E4750A84833B9CB7534F5DC3EB6D6A416ED | oito GETs 401 |
| probe Go final v7 | 70E6FAD8ADD0645DC601C2BCB53B396FB7700081FF49AC0303E1DB40C21CC583 | API live |
| overlay Go | 850548379BD5DC21B7F2AA0C71FCB5D954907FE418781B3B70F3AE8344644DAC | sem edição do repo |
| wrapper PowerShell | 701475EE431A5F77499B95576EA0587F06F393FB6438ECDC09032CFFAF82BA4A | ambiente redigido |

O probe foi ampliado durante a pesquisa; os hashes v6 e v7 preservam as duas versões.

## 4. Ambiente E1

~~~text
Python 3.12.13
requests 2.34.2
Scrapling 0.4.10
Playwright Python 1.61.0
Microsoft Edge 150.0.4078.65
Go 1.25.1
Windows
~~~

Versões foram verificadas por pip show, go.mod e metadata do executável.

## 5. PoCs HTML E2

### Baseline

~~~text
official item API: 403 PolicyAgent
product page: 200, redirect account-verification, sem preço
public search: 200, redirect account-verification, sem preço
~~~

### Scrapling

Comando: C:\tmp\ml-scrape-poc-venv\Scripts\python.exe C:\tmp\ml_price_poc_frameworks.py scrapling

~~~text
item: 200 em 284 ms, título presente, preço ausente
busca: 200 em 247 ms, título presente, preço ausente
~~~

Somente Fetcher e curl_cffi; nenhum modo stealth/browser.

### Playwright

Comando: C:\tmp\ml-scrape-poc-venv\Scripts\python.exe C:\tmp\ml_price_poc_frameworks.py playwright

~~~text
duas execuções do item MLB4735328201: 200, aproximadamente 4.752 s
meta=169.99, texto=R$169,99, JSON=169.99
busca: 200, aproximadamente 10,9–11,4 s
tentativa seguinte de cards: 403, zero cards
~~~

O seletor genérico também capturou parcela R$34. Não houve tentativa após o bloqueio.

## 6. API oficial v6/v7

### v6 — B

~~~text
installation=connected/healthy
credential_version=6
token e refresh presentes
oito endpoints: 401 invalid access token
~~~

### v7 — E1/E2

Ao iniciar o Postgres, o ticker renovou a credencial:

~~~text
credential v7 created_at=2026-07-12T15:53:19.199685Z
auth session state=valid
access_token_expires_at=2026-07-12T21:53:18.773285Z
~~~

Dois StartAuthorize posteriores criaram states não consumidos e mudaram a instalação para
pending_connection. A credencial e sessão continuaram válidas.

Regra oficial adicionada ao handoff: o refresh token é de uso único, somente o último permanece
válido e cada refresh devolve um novo par. A integração futura deve serializar refresh por
instalação, persistir access+refresh atomicamente e, em `invalid_grant`, reconciliar uma rotação mais
nova ou exigir reautorização — nunca repetir cegamente um token já consumido.

Comando final: C:\tmp\run_ml_price_probe.ps1

Resultado: PASS, PROBE_EXIT=0.

### Cobertura live

~~~text
own active items=10
catalog linked=9
non catalog=1
winning=3
competing=6
not_listed raw=1
own sale_price: 10/10 HTTP 200
price_to_win: 10/10 HTTP 200
~~~

### PDPs

~~~text
catalogs queried=9
offers declared=256
offers materialized=236
complete catalogs=8/9
own sale_price equals own catalog offer=9/9
price_to_win winner equals matching catalog offer=9/9
~~~

O catálogo com 120 ofertas teve apenas a primeira página de 100. Cobertura: 236/256 = 92,1875%.

### Superfícies negativas

~~~text
cross-seller /items: 1/1 HTTP 403
cross-seller /sale_price: 4/4 HTTP 403
/sites/MLB/search: 2/2 HTTP 403
/products/{id}.buy_box_winner: ausente 9/9 com HTTP 200
seller_sku=15956: total 0
~~~

O 403 é compatível com falta de permissão funcional, mas a causa não foi provada.
A lista de ofertas funciona live apesar do aviso en-US de desligamento.

## 7. Oracle E2

Timestamp Oracle: 2026-07-12 12:35:00.404743; timezone não emitido.

SQL base:

~~~sql
SELECT CODPROD, DESCRPROD, REFERENCIA, REFFORN, MARCA, FABRICANTE, CARACTERISTICAS
FROM METALPRD.TGFPRO
WHERE ATIVO = 'S'
ORDER BY CODPROD;
~~~

Saída resumida:

~~~text
active_total=10519
codprod_unique=10519
valid_gtin_rows=7351
valid_gtin_unique_rows=7163
valid_gtin_duplicate_keys=91
valid_gtin_duplicate_rows=188
refforn_filled_rows=10325
brand_refforn_unique_rows=10120
description_only_no_ref_fields=186
REFERENCIA_EQ_EAN -> 15956
REFFORN_EQ_MANUFACTURER_REF -> 15956
REFFORN_EQ_FUTURE_SKU -> zero
~~~

## 8. Código E1

- [SellerSKU projetado para REFFORN](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go#L384)
- [EAN e ReferenceCode recebem REFERENCIA](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go#L48)
- [M-09 separa identidades](../../../.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md#L154)
- [M-12 prevê escrita/readback do SKU](../../../.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md#L157)
- [StartAuthorize sem gate completo de transição](../../../apps/server_core/internal/modules/integrations/application/auth_flow_service.go#L256)
- [expiração canônica na sessão](../../../apps/server_core/internal/modules/integrations/adapters/postgres/auth_session_repo.go#L25)

## 9. Fontes E3

- [Concorrência pt-BR](https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/concorrencia-em-catalogo)
- [Aviso en-US](https://developers.mercadolivre.com.br/en_us/catalog-competition)
- [API de preços](https://developers.mercadolivre.com.br/pt_br/servicos-gerenciamento-de-contatos/api-de-precos)
- [Permissões funcionais](https://developers.mercadolivre.com.br/pt_br/permissoes-funcionais/)
- [Buscador de produtos](https://developers.mercadolivre.com.br/pt_br/buscador-de-produtos)
- [Autenticação e autorização](https://developers.mercadolivre.com.br/pt_br/gerenciamento-perguntas-respostas/autenticacao-e-autorizacao)
- [Gestão de identidades e acessos](https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/gestao-de-identidades-e-acessos-oauth-e-tokens)
- [Termos do Programa](https://developers.mercadolivre.com.br/pt_br/termos-e-condicoes)
- [Termos gerais](https://www.mercadolivre.com.br/ajuda/991)

## 10. Lacunas

1. instalação pending_connection apesar de sessão v7 válida;
2. permissões funcionais cross-seller não auditadas no DevCenter;
3. substituto durável da lista completa de ofertas desconhecido;
4. paridade visual somente 1/1 exploratória, não 30/30;
5. autorização de retenção/agregação ausente;
6. matching ≥99% sem coorte adjudicada;
7. descoberta fora do catálogo não comprovada;
8. quota e 429 não medidos.
