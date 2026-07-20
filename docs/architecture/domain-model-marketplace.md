# Modelo de Domínio — Operação de Venda em Marketplace

> **Status:** RASCUNHO v1 (entrevista operador + pesquisa + especialista Sankhya, 2026-07-19) — aguarda ratificação.
> Complementa [produto-anuncio-marketplace-identity.md](produto-anuncio-marketplace-identity.md).
> N:N/kit e catálogo-vs-item resolvidos (§6). Fatos autoritativos do ERP em §8.

## 1. Glossário (linguagem do negócio)

| Termo | Definição |
|---|---|
| **Produto** | Item do ERP (Sankhya). Tem custo, estoque, impostos. Identidade = **CODPROD**. |
| **CODPROD** | Código do produto no ERP. **É o SKU** exibido ao usuário. Nunca MLB. |
| **REFFORN** | Referência do fornecedor/fabricante. Secundária. Não é SKU, não é MLB. |
| **EAN/GTIN** | Código de barras. Usado pra dar match automático produto↔anúncio. |
| **Anúncio (item)** | Oferta do vendedor no marketplace (item MLB). Tem preço, tipo (Clássico/Premium/Full), categoria, status. |
| **Catálogo (ML)** | Produto de catálogo do ML (página do produto). Vários itens/anúncios de vendedores se penduram nele. |
| **Vínculo** | Ligação produto(CODPROD) ↔ anúncio (adapter). Resolve **categoria** → comissão/frete. |
| **Categoria** | Categoria ML do anúncio. Dirige **comissão** e regras. |
| **Evidência de mercado** | Mediana, posição, ofertas concorrentes coletadas do ML. Por categoria/busca. |

## 2. Entidades e relações (v0 — N:N sob pesquisa)

```
 PRODUTO (ERP)                  VÍNCULO (N:N + qtd)         ANÚNCIO / ITEM (ML)         CATÁLOGO (ML)
 internal_product_id=CODPROD ─┐  produto ↔ anúncio (qtd)  ┌─ item MLB                   catalog_product_id
 refforn, ean, descr          ├──────────────────────────┤   tipo (Clássico/Premium)    (vários itens
 custo (ERP)                  │  kit: 1 anúncio ↔ N prod. │   categoria ──► comissão      penduram aqui)
 icms_aliquota (Sankhya)      │  1 produto ↔ N anúncios   │   preço, status
 estoque (ERP)                ┘                           └─ peso/dim ──► frete (CEP)
                                                             │
                                                             └─ EVIDÊNCIA MERCADO: mediana, posição, ofertas
 CONFIG (nosso, não-ML): regime imposto, limiares margem verde/amarela, tabela DIFAL por UF, tarifa Full
```

**Cardinalidades (confirmadas na entrevista):**
- 1 produto → **N** anúncios (sem limite).
- 1 anúncio → **N** produtos (kit/combo). ⇒ **produto ↔ anúncio = N:N** (com quantidade por produto no kit).
- 1 produto → 1 catálogo ML (por marketplace) → N itens/anúncios pendurados.
- Comissão/frete/posição vivem no **anúncio (item)**, resolvidos pela **categoria** — não são atributos do produto.

## 3. Onde mora cada dado (fonte da verdade)

| Dado | Fonte | Nível |
|---|---|---|
| CODPROD, REFFORN, EAN, descrição, custo, estoque | **ERP (Sankhya)** ou planilha import | produto |
| **ICMS alíquota** | **NÃO é atributo do produto.** Motor Sankhya `TGFICM` por par UForig×UFdest (+ST/DIFAL/FCP); realizado no item `TGFITE.ALIQICMS`. Produto só carrega identidade fiscal (`GRUPOICMS`, `NCM`, `TEMICMS`). Ver §8. | operação (UF-par) |
| Comissão (%) | API ML `listing_prices` | anúncio (categoria + tipo) |
| Frete | API ML shipping | anúncio (peso/dim + CEP destino) |
| Tipo (Clássico/Premium/Full), tarifa Full | ML + config | anúncio |
| Mediana / posição / ofertas | Coleta ML | categoria/busca do anúncio |
| **Histórico de preço de mercado** | **FALTA — precisa criar** | série temporal por produto/categoria |
| Regime imposto, limiares margem, DIFAL por UF | **Nosso** (config) | tenant/global |

## 4. Import ERP — heterogeneidade a padronizar

- **Empresa A (Docol/Deca):** Sankhya conectado (Oracle `METALPRD.TGFPRO`). Fonte viva.
- **Empresa B:** Sankhya **não conectado**. Enviou planilha com **3 abas = 3 tipos de produto** distintos → o import precisa de um **schema de mapeamento por aba/tipo** (padronização). Sem ICMS na planilha.
- Regra: campos operacionais desconhecidos ficam honestos (branco/"—"), nunca zero (ADR-17).

## 5. Norte do produto (escopo-alvo do app, da entrevista)

Achar oportunidades de venda · preço de mercado estratégico · controle de estoque vs ERP (não zerar) ·
controlar/publicar anúncios · simular preço+lucro **antes** de entrar · **histórico de preço de mercado** ·
tela de produto · conciliação financeira / resultado · envios / pós-venda. — plataforma completa de
eficiência em marketplace.

## 6. Resolvido pela pesquisa (fontes ML/VTEX/Bling)

1. **N:N/kit** → padrão consolidado = **tabela BOM separada** (VTEX `stockkeepingunitkit`, Bling `estruturas/componentes`). Kit = entidade própria + filhos `(kit, componente, quantidade)`. **NÃO** sobrecarregar `product_links`. `product_links` já é N:N-capaz (várias linhas por produto). Fórmula kit: custo = Σ(custo·qtd), estoque = MIN(estoque/qtd).
2. **Catálogo vs item ML** → `catalog_product_id` é metadado **no item/variação** (N itens → 1 catálogo; buy-box). Persistir em `listings` em vez de rederivar. Reforça evidência de mercado (rivais de buy-box = comparáveis fortes). Fontes: `POST /items/catalog_listings`, `GET /items/{id}/catalog_listing_eligibility`.
3. **Variações** → cada variação ML ↔ 1 SKU interno. `product_links` já tem `provider_variation_id` (granularidade certa). Um SKU nunca deve cobrir várias variações do mesmo item (quebra baixa de estoque).
4. **Economia** → comissão = `GET /sites/MLB/listing_prices?category_id=&listing_type_id=` (categoria+tipo, **sem** param de produto). Frete = dimensões do item + `categories/{id}/shipping_preferences`. Ambos no **anúncio**, nunca no produto.
5. **ICMS + CODPROD real + contrato import** → consulta enviada ao especialista Sankhya (`local_ec787804`).

## 7. Decisões estruturais — FIX-NOW (aditivo, SEM feature de kit)

> Objetivo (operador): backend/lógica/contratos ESTRUTURADOS pra realidade, mesmo sem construir kit.

- **FIX-1 · Identidade (frontend):** coluna SKU = **CODPROD** (`internal_product_id`), nunca `manufacturer_reference`/MLB. REFFORN vira secundário. Aplicar em toda tela de produto.
- **FIX-2 · Dados demo:** trocar CODPROD sintético (90001..) + MLB-no-REFFORN pelos valores REAIS do Sankhya; MLB migra pro vínculo (`product_links.provider_item_id`). ← aguarda especialista.
- **FIX-3 · Contrato aditivo (estrutura pronta, sem feature):**
  - `catalog_product_id TEXT NULL` em `listings` (persistir da sincronização ML).
  - **ICMS ≠ campo de produto** (corrigido pelo especialista §8): produto guarda só identidade fiscal
    (`grupo_icms`, `ncm`, `tem_icms`). A carga tributária é resolvida por **operação** (par UForig→UFdest)
    espelhando `TGFICM` (`ALIQUOTA` interestadual origem, `ALIQUFDEST` interna destino, `PERCICMSFCP`,
    `ALIQSUBTRIB`/`MARGLUCRO` p/ ST). B2C EC 87/2015: carga vendedor = `ALIQUFDEST + FCP`;
    DIFAL = `ALIQUFDEST − ALIQUOTA (+FCP)`. Modelar como tabela config por UF-par (já existe DIFAL por UF
    na tela) + valor realizado; **não** cravar alíquota no produto. `TGFICM` não é versionado → é config vigente.
  - Documentar `listing_kit_components` (BOM) como schema-alvo; **não** construir agora. Código não pode assumir vínculo 1:1.
- **FIX-4 · Verificação — RODADA, achou 1 violação (auditoria read-only 2026-07-19):**
  - ✅ Nenhuma comissão é chaveada por `internal_product_id`/codprod. Caminho decompose
    (`calc_service` → `tarifflive/resolver.QuoteCommission`) resolve por `(category_id, listing_type)` da API ML — **correto**.
    Profitability lê `SaleFeeAmount` do snapshot do pedido — correto. Fallbacks 0.16/0.22 só p/ categoria "default" (bootstrap).
  - ❌ **VIOLAÇÃO — batch orchestrator (alimenta a matrix `/precos`):**
    `pricing/application/batch_orchestrator.go:184,188` busca a comissão por `prod.CategoryID`, que vem de
    `catalog/reader.go:37` = `TaxonomyNodeID` (**taxonomia do CATÁLOGO ERP do produto**), não a `category_id`
    do **anúncio ML**. ERP category ≠ ML listing category → margem do simulador pode sair errada.
    Correção estrutural: no batch, resolver categoria pela cadeia **vínculo → listing → category_id ML**
    (como o decompose já faz), nunca pela taxonomia ERP do produto. Candidato a chip (Go backend).

## 8. Fatos autoritativos do Sankhya (especialista `local_ec787804`, 2026-07-19, read-only)

### 8.1 Contrato de import — tabela.coluna exata (METALPRD)
| Campo import | Fonte Sankhya | Nota |
|---|---|---|
| CODPROD | `TGFPRO.CODPROD` | inteiro real, PK — **não** sintetizar 90001+ |
| EAN/GTIN | `TGFPRO.REFERENCIA` | neste ambiente o EAN mora aqui (`TGFBAR.CODBARRA` vazio) |
| Ref. fornecedor | `TGFPRO.REFFORN` | cód. Docol/Deca — **não** é o MLB |
| Descrição / NCM / Marca | `TGFPRO.DESCRPROD` / `.NCM` / `.CODMARCA→TGFMAR.DESCRICAO` | |
| Identidade fiscal | `TGFPRO.GRUPOICMS(+2)`, `.TEMICMS` | chave p/ regra tributária (não é a alíquota) |
| Custo | `TGFCUS.CUSSEMICM` **as-of** (CODPROD,CODEMP,DTATUAL≤ref) | **não** usar `TGFITE.CUSTO` |
| ICMS / ST / DIFAL | `TGFICM` por UF-par; realizado em `TGFITE.*` | ver §3 e FIX-3 |

### 8.2 Reconciliação demo→real (match por EAN em `TGFPRO.REFERENCIA`)
| EAN | Produto | CODPROD real | REFFORN real | Custo as-of | Obs |
|---|---|---|---|---|---|
| 7898016503522 | Placa SINALIZE Diretoria 14x14 | **20317** | 500AA | 9,89 ⚠️ | custo defasado (2010) — não confiar p/ margem |
| 7894200146179 | Papeleira DECA Flex 2020.C.FLX | **15956** | 2020.C.FLX | 91,57 | |
| 7894200160885 | Misturador DECA Polo cromado | **22467** (canônico) | 1877.C33 | 691,13 | dup 44975 (re-cadastro 2026-07-18) → operador decide |
| 7891461490355 | Torneira DOCOL Lift Ônix | **42519** (canônico) | 008720CE | 478,06 | dup 43534 quase morto → usar 42519 |
| 7891461034580 | Torneira DOCOL Pressmatic Compact | **39563** | 90171606006 | 138,85 | |

- **Bug de identidade confirmado (2 pontos):** demo usa CODPROD sintético (deve ser o real acima) **e** enfia
  o id MLB no `REFFORN` (REFFORN é código do fornecedor). MLB é id externo do canal → tabela de link
  anúncio↔produto (`product_links.provider_item_id`), **nunca** no REFFORN.
- **2 EANs têm CODPROD duplicado** (mesmo EAN, cadastro dobrado) → escolha do operador; recomendação: 22467 e 42519 (canônicos com histórico de venda).
- Especialista oferece gerar o SELECT batch de reconciliação para os 34 (EAN→CODPROD real + custo + carga por UFdest).
