# Handoff: Marketplace Central — Hub de vendas ML × ERP

## Overview
Protótipo completo de um hub de vendas em marketplaces (foco Mercado Livre) integrado a ERP: dashboard, gestão de anúncios, radar de concorrência, simulador de preço/margem, pedidos (com DIFAL), vínculos anúncio↔produto, repasses/conciliação e configurações. Idioma do produto: **português (BR)**.

## About the Design Files
Os arquivos `.dc.html` deste pacote são **referências de design criadas em HTML** — protótipos navegáveis que mostram aparência e comportamento pretendidos, **não código de produção para copiar**. A tarefa é **recriar estas telas no ambiente do codebase alvo** (React etc., pasta `marketplace-central` do usuário) usando os padrões e bibliotecas já estabelecidos lá — ou, se não houver ambiente definido, escolher o framework mais adequado e implementar nele.

Cada `.dc.html` abre direto no navegador (com `support.js` na mesma pasta). O markup relevante está dentro de `<x-dc>...</x-dc>`; a lógica/dados de exemplo estão no `<script data-dc-script>` no fim de cada arquivo. Estilos são todos inline com CSS vars — fáceis de extrair.

## ⚠️ Leia primeiro: API-MAP.md
`API-MAP.md` mapeia **cada bloco de cada tela** para o endpoint real do Mercado Livre / Mercado Pago, classificado como ✅ API direta, ⚠️ derivado (agregação/snapshot nosso) ou ❌ ERP/feature própria. **Nada na UI foi inventado sem esse mapeamento** — implemente cada bloco contra a fonte indicada. Os pontos de maior atenção arquitetural:
1. Snapshots diários próprios para preço de mercado/monitorados (não há histórico retroativo na API).
2. Release report (repasses) é assíncrono/CSV — ingestão agendada, não fetch on-demand.
3. DIFAL é 100% domínio nosso (tabela por UF + agendamento + lembretes).
4. Comissões/frete: sempre via `listing_prices`/`shipping_options`, nunca hardcodadas.
5. Webhooks (`orders_v2`, `shipments`, `claims`, `items`, `payments`) como fonte de eventos; GET como reconciliação.

## Fidelity
**High-fidelity.** Cores, tipografia, espaçamentos, estados e copy são finais. Recriar fielmente usando os componentes do codebase.

## Design Tokens
- Fontes: **Instrument Sans** (UI) e **IBM Plex Mono** (números, códigos, valores monetários)
- Light: `--bg #fbfaf7 · --surface #fff · --surface2 #f4f2ea · --border #e6e3da · --border2 #f0eee6 · --ink #25291f · --muted #6f6d63 · --faint #8b887c · --accent #4a7c59 · --accent-soft #edf3ec · --accent-ink #33573f · --warn #a3552e · --warn-soft #f7ebe2 · --info #2f5bb7 · --info-soft #e8eef9 · --amber #8a6d1f · --amber-soft #fdf3d7`
- Dark (via `data-theme="dark"`): ver `:root`/`[data-theme]` em qualquer arquivo
- Raios: 8px (botões/inputs), 10-12px (cards), 999px (chips/pills)
- Chips de margem: verde ≥18% · âmbar 10–18% · vermelho <10% (limiares configuráveis em Configurações)
- Tabelas: header `--surface2` 11px letter-spacing .04em; linhas 12.5px; `overflow-x:auto` com `min-width` interno
- Shell topo em todas as telas: logo M + nav pills (Visão geral · Anúncios · Mercado · Simulador · Pedidos · Repasses) + toggle tema + ⚙ (Configurações) + avatar

## Screens (arquivos neste pacote)
1. **Dashboard.dc.html** — visão geral com KPIs e agregações.
2. **Anuncios.dc.html** — tabela de anúncios, agrupar por produto, exceções, seleção em massa, drawer de detalhe (sem edição inline; edição completa → Produto Detalhe).
3. **Produto Detalhe.dc.html** — header do produto + box "Vale a pena vender?"; abas: Anúncios vinculados, Concorrência, Estoque (físico/reservado/disponível/Full + alerta de cobertura + movimentações), Pedidos (mini-métricas + lista), Histórico (preço vs. mercado + auditoria). Aba "Dados" ainda placeholder.
4. **Mercado.dc.html** — radar: Reprecificação (nossos anúncios), Oportunidades (não vendemos), Monitorados (vendedores/termos/anúncios com filtros, variação 7d e alertas).
5. **Simulador.dc.html** — matriz de produtos + painel preço↔margem bidirecional com decomposição real (comissão, taxa fixa, frete, imposto, custo ERP), modalidades Clássico/Premium/Full, cenários; **drawer "⚙ Parâmetros de cálculo"** (regime/alíquota, limiares de margem, tarifa Full, DIFAL on/off) com recálculo ao vivo.
6. **Pedidos.dc.html** — KPIs + 3 visões (Fila por urgência SLA→DIFAL, Lista com abas de status inclusive Cancelados/Devoluções reais, Kanban read-only); drawer de detalhe com decomposição, bloco DIFAL (agendar/lembrar/marcar pago), timeline, NF, rastreio, comprador.
7. **Vinculos e Importacao.dc.html** — matching anúncio↔produto ERP (SKU→GTIN→título) com chips de confiança, vinculação manual via drawer, lote com protocolo. Vincular não altera o ML.
8. **Repasses.dc.html** — KPIs, calendário de liberações, tabela vendas−tarifas−fretes−retenções=líquido com chip de conciliação, drawer-extrato com retenções e conciliação ERP.
9. **Configuracoes.dc.html** — Geral (tema, visão padrão), Notificações, Cálculo & DIFAL (tabela 27 UFs "padrão 2026" com ajuste de exceção via drawer: interna − interestadual = DIFAL efetivo), Integrações.

`HANDOFF.md` (raiz do pacote) tem o histórico de decisões de produto/design — **não regredir** nas decisões listadas lá (ex.: sem edição inline em tabelas, status do Kanban vem do ML, config em 3 camadas global→produto→anúncio).

## Interactions & Behavior
- Drawers laterais (300–400px, overlay rgba escuro) para todo detalhe/edição; fechar por ✕ ou clique no overlay.
- Navegação entre telas por links relativos (mesmos nomes de arquivo).
- Estados hover em linhas de tabela (`--surface2`) e botões (opacity .9 / border accent).
- Tema light/dark por `data-theme` na raiz; toggle no header de cada tela.
- Ações mutantes (aplicar preço, faturar, vincular) passam por "fila de sync" com preview e protocolo — padrão do produto.
- Recálculo ao vivo: parâmetros do Simulador afetam matriz, chips e veredictos imediatamente.

## State Management (essência por tela)
Ver o `class Component` de cada arquivo: o estado usado no protótipo (seleção, drawer aberto, overrides, filtros) é o estado mínimo que a tela real precisa. Dados de exemplo (catálogos, pedidos, repasses) mostram o shape esperado das entidades.

## Assets
Nenhuma imagem — placeholders hachurados indicam onde entram fotos de produto e gráficos (histograma de preços, linha preço vs. mercado). Fontes via Google Fonts.

## Files
- `API-MAP.md` — mapeamento UI × API (leia primeiro)
- `HANDOFF.md` — decisões de produto/design
- `*.dc.html` — as 9 telas
- `support.js` — runtime dos protótipos (só para abrir os .dc.html no navegador; não portar)
