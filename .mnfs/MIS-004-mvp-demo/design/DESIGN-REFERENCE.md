# DESIGN-REFERENCE — MIS-004 (BINDING para QA visual)

Fonte: pacote de handoff do operador (2026-07-18), `design/handoff/` — 9 protótipos navegáveis
`.dc.html` + `HANDOFF.md` (decisões) + `README.md` (tokens/fidelity) + `API-MAP.md` (UI×API).
**Fidelity declarada: HIGH — cores, tipografia, espaçamentos, estados e copy são FINAIS.**
Protótipos são referência de design, NÃO código de produção (não portar `support.js`).

Status: ratificado pelo hub como referência 1:1 para o gate visual de QA (wave B em diante).
Render verificado no browser (Simulador.dc.html) 2026-07-18.

## Como o QA usa

1. Abrir o `.dc.html` da tela no browser (file://, `support.js` na mesma pasta) lado a lado com o app real (:5174).
2. Comparar por item do checklist da tela (abaixo) + tokens + decisões de não-regressão.
3. Divergência de token/copy/layout = FINDING com screenshot dos dois lados. Divergência de
   dado (protótipo usa dados fake) NÃO é finding — shape/formato é o que conta.
4. Light E dark (`data-theme`), desktop primeiro; tabelas devem ter `overflow-x:auto` + `min-width` interno.

## Tokens (finais — extraídos de `:root` dos .dc.html, conferidos 2026-07-18)

Fontes: **Instrument Sans** (UI) · **IBM Plex Mono** (números/códigos/valores monetários).

| Token | Light | Dark |
|---|---|---|
| --bg | #fbfaf7 | #161814 |
| --surface | #ffffff | #1f221c |
| --surface2 | #f4f2ea | #262a23 |
| --border | #e6e3da | #2e312a |
| --border2 | #f0eee6 | #2a2d26 |
| --ink | #25291f | #e9e8e2 |
| --muted | #6f6d63 | #a5a399 |
| --faint | #8b887c | #8b887c |
| --accent | #4a7c59 | #7fb08c |
| --accent-soft | #edf3ec | #243328 |
| --accent-ink | #33573f | #9ecfab |
| --warn | #a3552e | #d08a63 |
| --warn-soft | #f7ebe2 | #3a2a20 |
| --info | #2f5bb7 | #7fa3e0 |
| --info-soft | #e8eef9 | #222b3c |
| --amber | #8a6d1f | #d4b45a |
| --amber-soft | #fdf3d7 | #37311d |

Raios: 8px botões/inputs · 10–12px cards · 999px chips/pills.
Tabelas: header `--surface2`, 11px, letter-spacing .04em; linhas 12.5px.
Chips de margem: verde ≥18% · âmbar 10–18% · vermelho <10% (limiares configuráveis).
Chips de confiança (vínculos): verde ≥85 · âmbar 50–84 · vermelho <50 (= bandas IC-01).
Shell topo (todas as telas): logo M + nav pills (Visão geral · Anúncios · Mercado · Simulador ·
Pedidos · Repasses) + busca + toggle tema + ⚙ + avatar. Vínculos FORA da nav global.

## Decisões de não-regressão (HANDOFF.md — findings se violadas)

- Sem edição inline em tabelas — detalhe/edição via drawer (300–400px, overlay) ou tela própria.
- Comparação preço mercado: valor R$ + chip % colorido; nunca duas sublinhas por célula.
- Config em 3 camadas: global (Parâmetros de cálculo) → produto → anúncio (comissão real do ML, só exibe).
- Kanban de Pedidos: status vem do ML/ERP; arrastar NÃO muda status.
- Ações mutantes (aplicar preço, faturar, vincular) passam por fila de sync com preview e protocolo.
- Vincular NÃO altera o ML.
- Sem AI-slop: sem gradientes agressivos, sem emojis, sem Inter/Roboto.
- Simulador sem "Impacto/mês" (sem tracking confiável).

## Telas × milestones (wave B em negrito)

| Tela (.dc.html) | Milestone | Checklist QA (essência) |
|---|---|---|
| **Anuncios** | **M-05** | Tabela plana + agrupar por produto; exceções clicáveis; seleção em massa; drawer detalhe SEM edição inline, "Abrir edição completa" → Produto Detalhe; chip % vs mercado |
| **Simulador** | **M-07** | Matriz SKU/DESCRIÇÃO/CUSTO/NOSSO PREÇO/PREÇO MERCADO/MARGEM/VEREDICTO (min-width 900px, scroll-x); painel preço↔margem bidirecional com decomposição real (comissão, taxa fixa <R$79, frete ≥R$79, imposto, custo ERP); Clássico/Premium/Full; drawer "⚙ Parâmetros de cálculo" com recálculo AO VIVO (regime Simples 4%/L.Presumido 9,25% + alíquota editável, limiares de margem, tarifa Full, DIFAL toggle + mini-tabela UF); salvar cenários; aplicar via fila de sync; banner frete ≥R$79 |
| **Pedidos** | **M-08** | KPIs (Novos/A faturar/A enviar/Enviados/DIFAL a pagar); 3 visões: Fila (padrão, urgência SLA→DIFAL→data) / Lista (abas de status c/ contagem, seleção em massa) / Kanban read-only; drawer: itens+vínculo, decomposição c/ retorno real, bloco DIFAL (agendar/lembrete/marcar pago), timeline ML→interno, NF, rastreio, comprador; coluna-chip DIFAL (pago ✓ / vence hoje / ag. dd/mm) |
| Dashboard | M-06 (wave C) | KPIs/agregações; tudo derivado |
| Mercado | M-06 (wave C) | Reprecificação + Oportunidades + Monitorados (variação 7d, alertas) |
| Produto Detalhe | M-06 (wave C) | Header + "Vale a pena vender?"; abas Anúncios vinculados/Concorrência/Estoque/Pedidos/Histórico; "Dados" placeholder |
| Vinculos e Importacao | M-04 (fechado) | Referência retroativa — match SKU→GTIN→título, chips de confiança IC-01, lote+protocolo, desfazer |
| Repasses | M-09 (CORTÁVEL) | KPIs, calendário liberações, vendas−tarifas−fretes−retenções=líquido, chip conciliação, drawer-extrato |
| Configuracoes | transversal | Geral/Notificações/Cálculo & DIFAL (tabela 27 UFs "padrão 2026", exceção via drawer: interna − interestadual = efetivo, chip ajustada/padrão, restaurar)/Integrações |

## Alinhamento com contratos existentes

- Chips de confiança = bandas IC-01 (ALTA ≥85 / MEDIA 50–84 / BAIXA <50) — consistente.
- Decomposição do Simulador = IC-04 (formula única, unknown ≠ zero) — o protótipo mostra "—"
  para componente ausente (ex.: taxa fixa ≥R$79), nunca 0 fabricado — consistente com ADR-17.
- DIFAL: interna − interestadual = efetivo, override por UF — consistente com domain.DifalForUF (M-07 d288cd12).
- `API-MAP.md` é consultivo para implementação (fonte por bloco: ✅ API direta / ⚠️ derivado / ❌ ERP-nosso);
  não substitui OpenAPI/SDK como truth de contrato.

## Gaps conhecidos (não são findings de QA)

- Protótipo usa dados fake; app usa dev stack (import #003-E). Comparar shape/formato, não valores.
- Backlog do próprio design (HANDOFF.md): link "ver tabela completa por UF" → Configurações (hoje #);
  aba "Dados" do Produto Detalhe; drawer "2 candidatos" em Pedidos; variações de layout de Repasses.
- Auto-ACCEPT de vínculos inatingível na conta ML atual (sem seller_sku) — ceiling REVIEW, decisão D-46.
