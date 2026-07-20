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

### 8.2 Reconciliação demo→real — batch dos 34 COMPLETO
Mapa autoritativo (codprod_demo → codprod_real + refforn_real + custo + carga por UF):
**[`demo-reconciliation-34.tsv`](demo-reconciliation-34.tsv)** (especialista, 31/34 casados por EAN).

- **31 casaram** por `TGFPRO.REFERENCIA`=EAN (placas SINALIZE, Deca, Docol, Lorenzetti).
- **3 sem match** (EAN vazio no import): 90004 Puxador Feng, 90005 Fechadura Imab, 90007 Toalheiro Soul Zen → título/manual.
- **2 EANs com CODPROD duplicado** → canônico (mais vendas): Polo=**22467** (dup 44975), Lift Ônix=**42519** (dup 43534).
- **Bug de identidade confirmado (2 pontos):** CODPROD sintético (90001+) **e** id MLB enfiado no `REFFORN`
  (visto cru: `REFFORN='MLB3758134295'`). CODPROD real e REFFORN real no TSV; MLB → link externo
  (`product_links.provider_item_id`), nunca `REFFORN`.
- ⚠️ **90001 (20317)** custo vigente é de 2010 (9,89) — defasado, não confiar p/ margem.

### 8.3 Carga ICMS B2C — por UF de DESTINO, não por produto (EC 87/2015)
Teto da carga = `ALIQUFDEST + FCP` da UF destino (não-contribuinte, partilha 100%→destino). **Produto-independente**;
ST/redução só ABAIXA o teto. Saindo de MG(13): **SP = 18%** (FCP 0), **RJ = 22%** (dest 20 + FCP 2).
Aplicar: `preco_liquido = preco_bruto × (1 − carga/100)`, comparar com `CUSSEMICM`.
⇒ confirma FIX-3(c): ICMS é config por UF-par, nunca coluna do produto.
