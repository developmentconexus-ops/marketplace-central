# Research Note

```yaml
id: R-02
type: research
status: draft
owner: Codebase Investigator
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Extração semântica das telas MVP do pacote de design high-fidelity `docs/design/handoff-2026-07/` (Anúncios, Produto Detalhe, Simulador, Pedidos, Vínculos, Dashboard, Configurações§DIFAL) para autoria de feature briefs.

## Sources Checked

- Source: as 7 telas .dc.html lidas integralmente (investigador dedicado, 2026-07-17). Markup `<x-dc>`, dados em `<script data-dc-script>`.
- Why it matters: briefs FE precisam de colunas, chips, drawers e estados exatos; mocks não são código de produção.

## Findings

### Anúncios

- Blocos: header (contagem + 4 chips de exceção erro/desat/vinc/marg + Agrupar/Exportar/+Criar); tabs Todos/Ativos/Pausados/Com pendência; barra bulk condicional; tabela; drawer 300px.
- Tabela: (checkbox) MLB · TÍTULO · PRODUTO · PREÇO · EST. · SYNC · QUAL. · PENDÊNCIA; linhas grupo (chevron, meta "ERP est. N · M anúncios", pill ok/N erro) ou item.
- Entidade: mlb, name, produto|"sem vínculo", modal(Clássico/Premium), price, stock(+"0 ⚠"), sync(Sincronizado=verde/Erro=warn/Desatualizado=amber/Na fila=info/Pausado), qual%, pend, margin%+marginOk, error{title,body,tech}, timeline[], noLink.
- Drawer: id+título, produto+modalidade, status, banner de erro (com linha técnica), grid Preço/Est./Margem/Qualidade, ações dinâmicas (Corrigir/Vincular/Editar + Simular preço + Pausar), timeline, "Abrir edição completa →" → Produto Detalhe.
- Interações: seleção em massa → Pausar/Atualizar preço/Re-sincronizar ("preview antes de aplicar"); toggle agrupar por produto; clique linha = drawer.
- Negativos: "sem vínculo", estoque 0 ⚠, erro de atributo obrigatório ML.

### Produto Detalhe

- Header produto: foto, nome, GTIN, NCM, custo, dimensões, estoque físico/reservado/disponível/seguro, barra Completude % + campos faltantes. Card "Vale a pena vender?" (amber).
- Tabs: Dados(placeholder) / Estoque / Anúncios vinculados(default) / Concorrência / Pedidos / Histórico. MVP = header + box veredicto + Anúncios vinculados + Estoque.
- Anúncios vinculados: MLB·CONTA·TÍTULO·MODAL.·PREÇO·VS MERCADO(chip "+18% ↑med")·SYNC·VENDAS 30D·(abrir).
- Estoque: KPIs FÍSICO ERP · RESERVADO(link Pedidos) · DISPONÍVEL=físico−reservado−seguro · FULL ML; banner cobertura ~9d + CTA ERP; tabela movs DATA·TIPO(chip venda/reserva/entrada/envio Full/devolução/ajuste)·ORIGEM·QTDE·SALDO.
- (Concorrência/Pedidos/Histórico = MIS-005.)

### Simulador

- Header: picker "Produtos: N ▾" (busca/filtros, listas NA ANÁLISE + SUGESTÕES), chip CEP, resumo params + botão "⚙ Parâmetros de cálculo".
- Tabela: SKU·DESCRIÇÃO·CUSTO·NOSSO PREÇO·PREÇO MERCADO(+posição "19º/23")·MARGEM(retorno+chip%)·VEREDICTO. "novo" quando sem anúncio.
- Painel Simular (380px): Preço↔Margem bidirecional (busca binária); atalhos mediana/menor/"cobrir verde%"; modalidades Clássico/Premium/Full; breakdown mono: Preço −Comissão% −Taxa fixa −Frete −Imposto% −DIFAL(uf) −TarifaFull −Custo ERP = Retorno/un; nota frete (≥79 grátis obrigatório/<79 taxa fixa 6,50); ação primária "Aplicar no anúncio (fila de sync)"|"Criar anúncio com este preço" + Salvar cenário; lista de cenários.
- Drawer Parâmetros: IMPOSTO (Simples 4%/Presumido 9,25% + alíquota editável); LIMIARES margem verde/amarela; MODALIDADES (comissões read-only "vem do ML", Tarifa Full editável); box regras ML não editáveis; DIFAL toggle + destino + preview UFs + link tabela completa; Restaurar padrão.
- Fórmula mock: comissão=preço×pct; fixo=preço<79?6,50; frete=preço≥79?frete_produto; imposto=preço×aliq; difal=preço×pct_uf; margem=preço−tudo−custo. Faixas ≥verde/≥amarela/vermelho.
- Veredictos mock: "saudável — manter/aplicar|anunciar" / "viável no preço de mercado" / "apertado — revisar custo/frete" / "não vale — custo alto p/ ML".
- DEFEITO DO MOCK (não copiar): difalTable 3 UFs hardcoded, sempre aplica índice [0]=SP ignorando destino real. Implementação real usa fonte única (tabela UF config) + destino real.

### Pedidos

- Header: filtros período/logística, banner DIFAL condicional → Fila, switch Fila/Lista/Kanban. 5 KPIs clicáveis: NOVOS("aguard. pagto") · A FATURAR(nota urgente) · A ENVIAR · ENVIADOS(7d) · DIFAL A PAGAR(soma, "N hoje").
- Fila: ordenada rank 0 SLA hoje(err) → 1 DIFAL hoje(âmbar) → 2 SLA amanhã → 3 faturar/enviar → 4 aguardando pagto; linha: tag, num, comprador·item·valor, retorno+chip%, ação|nota.
- Lista: 7 sub-tabs (Novos/A faturar/A enviar/Enviados/Concluídos/Cancelados/Devoluções). Cancelados/Devoluções = STUB explícito no mock (placeholderMsg "em breve") → MIS-005. Tabs padrão: (checkbox)·PEDIDO·DATA·COMPRADOR·ITENS·VALOR·RETORNO(+chip%)·SLA ENVIO·DIFAL(chip)·AÇÃO; bulk Imprimir etiquetas/Marcar enviados.
- Kanban: 4 colunas (Novos/A faturar/A enviar/Enviados), drag desabilitado ("status vem do ML/ERP").
- Drawer: num+status+logística; ITENS (sku, desc, valor, "↳ vínculo"); decomposição Valor −Comissão% −Taxa fixa −Frete −Imposto −DIFAL −Custo ERP = Retorno real; card DIFAL (rota, valor, status agendado/vence hoje/pago, ações) — MVP: chip informativo da seed, SEM agendar/lembrar/pagar; timeline ML→interno (Pagamento→NF→Etiqueta→Despachado→Entregue); NF+link ERP, rastreio, comprador (doc mascarado); ação primária por status (Faturar via ERP→Etiqueta→Enviado ✓).
- Status chips: novo="aguard. pagamento"(cinza) · faturar="pago·falta NF"(info) · enviar="NF emitida"(âmbar) · enviado(verde) · concluído="entregue"(verde).

### Vínculos e Importação

- Header: chip importação #N + conta, "N anúncios lidos · data", nota método "SKU→GTIN→título(fuzzy)", botão Nova importação. 3 KPIs=tabs: CONFIRMADOS("automático") · REVISÃO NECESSÁRIA(âmbar) · SEM CORRESPONDÊNCIA(vermelho).
- Tabela: (checkbox)·ANÚNCIO ML·SKU ML·PRODUTO SUGERIDO·SKU HUB·GTIN("✓ igual"|"—")·CONFIANÇA(% chip: ≥85 verde/50-84 âmbar/<50 vermelho)·MOTIVO·AÇÃO.
- Estados de linha: pendente por tipo (vincular→Vincular+Outro…; escolher→Escolher… [ambíguo]; criar→Criar produto+Ignorar) vs resolvido ("vinculado ✓"/"ignorado"/"produto criado ✓" + desfazer).
- Bulk: "Vincular sugeridos"/"Ignorar" → "preview: N registro(s) na fila · protocolo #N-X" → Aplicar → banner "✓ Protocolo aplicado — N vínculos. Nada foi alterado no ML".
- Drawer escolha manual: origem (título+mlb+SKU ML), candidatos (sku, nome, chip confiança, "GTIN idêntico · custo · estoque", Vincular este), fallback "Criar produto a partir do anúncio".
- Negativos: sem correspondência (conf=0, "nenhum — é um kit/combo", motivo), só Criar/Ignorar.

### Dashboard

- MOCK DESATUALIZADO (multi-canal Amazon/Shopee/Magalu, domínio genérico) — NÃO copiar conteúdo; recriar no domínio MPC/ML-only com agregações reais dos dados já trazidos. Estrutura aproveitável: 4 KPI cards + tabela pedidos recentes + donut saúde dos anúncios + lista mais vendidos + rail resumo.

### Configurações — Cálculo & DIFAL (só o necessário no MVP)

- Tabela DIFAL: 27 UFs [uf, nome, interna%, interestadual%], origem SC; interestadual 12% Sul/Sudeste(MG,PR,RJ,RS,SC,SP) / 7% resto; DIFAL = interna−interestadual (computado); chip "ajustada"(info)/"padrão 2026"; drawer exceção: input alíquota interna 0-35%, breakdown ao vivo, Restaurar/Salvar (persiste se Δ>0,049pp); overrides = mapa esparso {uf: interna}.
- Card params globais read-only (edição via drawer do Simulador). Tela Configurações completa = MIS-005; MVP consome a TABELA (seed + leitura) p/ Simulador+Pedidos.

### Shell / nav

- Header comum: logo M + wordmark, nav pills, busca, toggle tema (data-theme, default via prop), ⚙ → Configurações, avatar. Fontes Instrument Sans + IBM Plex Mono.
- Inconsistência nos mocks (Anúncios/Dashboard com nav divergente, "Pedidos" 2×, "Estoque"→#): nav canônica = HANDOFF.md: Visão geral · Anúncios · Mercado · Simulador · Pedidos · Repasses · ⚙. Vínculos FORA da nav global (entrada via Anúncios/importação + Configurações→Integrações).
- Contagens agregadas dos mocks (23/14/41/6, "1180+", Cancelados=3) são placeholder — sempre derivar de dados reais.

## Recommendation

Briefs FE citam este note por tela; mocks são referência de aparência/comportamento, defeitos listados (DIFAL SP-hardcoded, Dashboard stale, nav divergente, contagens fake) NÃO se reproduzem.

## Impact On Mission

Define densidade dos briefs FE (State/Interaction Model por tela), fonte única DIFAL, decisão Dashboard-MPC, nav canônica.

## Handoff

- Current status: completo.
- Next owner: Mission Strategist (P3/P5).
- Next action: usar por tela nos feature briefs.
- Required files/evidence: docs/design/handoff-2026-07/*.dc.html; HANDOFF.md.
- Blockers or open decisions: Mercado/Repasses pills no MVP — decidir estado "em breve" vs ocultas (assunção registrada em P1).
