# Modelo de Domínio — Operação de Venda em Marketplace

> **Status:** RASCUNHO v0 (entrevista com operador 2026-07-19) — aguarda ratificação.
> Complementa [produto-anuncio-marketplace-identity.md](produto-anuncio-marketplace-identity.md).
> Ponto N:N/kit e níveis catálogo-vs-item estão sob PESQUISA (external-researcher em curso).

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
| **ICMS alíquota** | **Sankhya** (tem por produto, calcula) — planilha atual não tem → manual/branco honesto | produto |
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
  - `icms_aliquota` (nullable, honesto) no contrato de produto — aguarda Sankhya; branco/manual até lá.
  - Documentar `listing_kit_components` (BOM) como schema-alvo; **não** construir agora. Código não pode assumir vínculo 1:1.
- **FIX-4 · Verificação:** garantir que nenhuma resolução de comissão usa `internal_product_id` — sempre transita pela `category_id` do anúncio.
