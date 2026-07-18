# Handoff — Marketplace Central (protótipo UI)

## Produto
Hub de vendas em marketplaces (foco Mercado Livre) integrado a ERP: anúncios, sync de preço/estoque, radar de concorrência, simulador de margem. Idioma: **português (BR)**. Wireframes base: `Wireframes MPC.dc.html`.

## Design system (aprovado)
Direção "refinado minimal" em tons de papel + verde. Tudo inline-styled com CSS vars no `<helmet>` de cada tela (copiar de qualquer .dc.html existente):
- Fontes: **Instrument Sans** (UI) + **IBM Plex Mono** (números/códigos)
- Light/dark via `data-theme` + toggle no header; prop `theme` (enum) em cada DC
- Tokens: `--bg #fbfaf7, --surface #fff, --surface2 #f4f2ea, --border #e6e3da, --ink #25291f, --muted, --faint, --accent #4a7c59, --accent-soft/-ink, --warn #a3552e, --amber, --info` (+ variantes dark)
- Padrões: shell topo (logo M + nav pills + busca + toggle tema + avatar AL), chips/pills de status com dot, tabelas em grid com header `--surface2`, drawer lateral 300–380px, chips de margem: verde ≥18% · âmbar 10–18% · vermelho <10%

## Telas prontas (aprovadas)
- `Dashboard.dc.html` — visão geral
- `Anuncios.dc.html` — tabela plana + agrupar por produto opcional, exceções clicáveis, seleção em massa, drawer de detalhe (sem edição inline; "Abrir edição completa" → Produto Detalhe)
- `Produto Detalhe.dc.html` — header produto + box "Vale a pena vender?", abas: Anúncios vinculados, Concorrência (Dados/Estoque/Pedidos/Histórico = placeholder)
- `Mercado.dc.html` — radar: Reprecificação + Oportunidades (+ Monitorados placeholder)
- `Pedidos.dc.html` — hub multi-view: KPIs (Novos/A faturar/A enviar/Enviados/DIFAL a pagar) + 3 visões via toggle: **Fila** (padrão — ordenada por urgência SLA → DIFAL → data), **Lista** (abas de status c/ contagem; Cancelados/Devoluções placeholder; seleção em massa: etiquetas/enviados), **Kanban** (status vem do ML/ERP, arrastar não muda). Drawer de detalhe: itens+vínculo, decomposição c/ retorno real, bloco DIFAL (agendamento+lembrete+marcar pago), timeline ML→interno, NF/ERP, rastreio, comprador. Props: theme, visaoPadrao, lembretesDifal. **DIFAL**: cadastro por pedido c/ data agendada, lembretes no topo + coluna-chip (pago ✓ / vence hoje / ag. dd/mm) — pedido do usuário p/ adiar e lembrar pagamento.
- `Vinculos e Importacao.dc.html` — wireframe 1g: KPIs confirmados/revisão/sem correspondência, match SKU→GTIN→título, tabela c/ chip de confiança (verde ≥85 · âmbar 50–84 · vermelho <50), ações Vincular/Outro…/Escolher…/Criar produto/Ignorar (+desfazer), drawer de escolha manual de candidato, lote + Aplicar c/ protocolo. Vincular não altera o ML.
- `Simulador.dc.html` — **tela crítica, muito iterada.** Drawer "⚙ Parâmetros de cálculo" (formato 1e, escolhido): recalcula ao vivo — regime (Simples 4% / L.Presumido 9,25%) + alíquota editável, limiares de margem verde/âmbar (mudam chips, veredictos e atalho "cobrir X%"), tarifa Full editável, comissões/regras ML read-only, DIFAL toggle (desconta UF destino na decomposição) c/ mini-tabela UF + link "ver tabela completa por UF" → futura tela 1f. Restaurar padrão. Matriz: SKU | DESCRIÇÃO | CUSTO | NOSSO PREÇO | PREÇO MERCADO | MARGEM | VEREDICTO (min-width 900px, scroll-x). Painel lateral: preço↔margem bidirecional com cálculo real, atalhos, Clássico/Premium/Full, decomposição (comissão, taxa fixa <R$79, frete ≥R$79 vendedor paga, imposto 4%, custo ERP), salvar cenários, aplicar via fila de sync. Sem "Impacto/mês" (sem tracking confiável). Botão "⚙ Parâmetros de cálculo" ainda sem tela.

## Decisões de design (não regredir)
- Sem edição inline em tabelas — detalhe via drawer/tela
- Comparações preço mercado: valor R$ + chip % colorido, nunca duas sublinhas por célula
- Config em 3 camadas: global (Parâmetros de cálculo) → produto (cadastro: peso/custo/overrides) → anúncio (comissão real vem do ML, só exibe)
- Tabelas: `overflow-x:auto` + `min-width` interno; nav pills com `white-space:nowrap`
- Sem AI-slop: nada de gradientes agressivos, emojis, Inter/Roboto

## Telas novas deste ciclo
- `Configuracoes.dc.html` — engrenagem ⚙ no header (todas as telas linkam): Geral (tema, visão padrão Pedidos), Notificações (DIFAL/ruptura/radar), Cálculo & DIFAL (resumo dos parâmetros globais + link p/ drawer do Simulador; tabela DIFAL 27 UFs pré-carregada "padrão 2026", ajuste de exceção via drawer c/ interna − interestadual = DIFAL efetivo, chip ajustada/padrão, restaurar), Integrações (ERP + ML, link p/ Vínculos). Formato 1f entregue aqui.
- `Repasses.dc.html` — versão simples aprovada: KPIs (próximo repasse, a receber 30d, retido/disputa, conciliação ERP), calendário de liberações, tabela vendas−tarifas−fretes−retenções=líquido c/ chip de conciliação (bateu/divergente/a conciliar/previsto), drawer-extrato c/ retenções detalhadas, conciliação ERP (ação "Conciliar" p/ divergentes) e maiores pedidos.
- Produto Detalhe: abas reais — Estoque (físico/reservado/disponível/Full + alerta cobertura 9d + movimentações), Pedidos (mini-métricas + lista c/ link), Histórico (preço nosso vs. mercado 90d + auditoria quem-mudou-o-quê). Só "Dados" segue placeholder.
- Pedidos: Cancelados/Devoluções reais — motivo ML c/ chip "rep ↓", status reembolso, reversa/estoque, banner de reputação, ação Contestar no ML.
- Mercado: Monitorados real — mistura vendedor/termo/anúncio c/ filtros, leitura atual, variação 7d, alerta (alertas viram itens do radar), + Monitorar….
- Decisão: Vínculos fora da nav global (entrada via Anúncios/importação + link em Configurações → Integrações). Nav agora: Visão geral · Anúncios · Mercado · Simulador · Pedidos · Repasses · ⚙.

## Backlog (próximas telas)
1. Apontar o link "ver tabela completa por UF" do drawer do Simulador para Configuracoes.dc.html (hoje é #)
2. Aba "Dados" do Produto Detalhe (cadastro/ficha técnica)
3. Pedidos: drawer "2 candidatos" ▾ na coluna itens
4. Repasses: usuário pediu p/ explorar variações de layout depois ("Explore a few options" marcado)

## Nav
Shell topo agora: Visão geral · Anúncios · Mercado · Simulador · Pedidos · Repasses(#). Vínculos tem pill própria só na sua tela (entrada via Anúncios/importação — decidir se entra na nav global).

## Contexto de trabalho
- Time comenta nas telas (Leandro etc.) — tratar comentários como feedback de UX
- Usuário gosta de: entrevistas via perguntas antes de telas críticas, opções lado a lado em canvas p/ decisões visuais (`Simulador Opcoes.dc.html` é exemplo)
- Pasta local `marketplace-central` anexada ao projeto (não explorada — usuário disse que código será adaptado ao frontend depois)
- Pesquisa de referências em `research/`
