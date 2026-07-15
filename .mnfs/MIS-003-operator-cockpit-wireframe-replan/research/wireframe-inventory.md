# Research Note

```yaml
id: R-01
type: research
status: draft
owner: Codebase Investigator
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Screen inventory of the primary planning input: `C:\Users\leandro.theodoro\Downloads\Wireframes MPC - Standalone.html` (self-extracting single-file wireframe, pt-BR, hand-drawn style; static HTML, navigation via URL fragments only).

## Sources Checked

- Source: the wireframe HTML (content extracted exhaustively from its `__bundler/template` payload, ~1,360 text nodes).
- Why it matters: MIS-003 is wireframe-driven; screens must not be invented.

## Findings

### Structure

- 16 screens in 2 decks. **Deck 2 (2a–2d) is the declared evolution target**: its header states 2a unifies explorations 1a+1b+1d, adds the 9th sidebar item "Mercado", and market intelligence appears in 3 places (tela Mercado 2c, aba Concorrência in 2b, colunas de mercado no Simulador 2d).
- Deck 1 = 8-workspace baseline: `Visão geral · Catálogo · Anúncios · Vínculos & Import. · Estoque · Preços & Simulador · Pedidos · Integrações & Sync` + Editor de anúncio (1l, 9-step wizard).
- Design language: "pt-BR · desktop-first · denso estilo ERP · filtros globais (empresa / marketplace / conta) na barra superior · badges de sistema-mestre · estados de sync explícitos".

### Screens (mission-relevant summary)

| ID | Screen | Key elements |
| --- | --- | --- |
| 2a | Anúncios unificado | Tabela mestre agrupada por produto (colunas: Produto/anúncio, MLB, Modal., Preço, vs mercado, Est., Sync, Qual., Pendência); faixa de exceções clicável (🔴 erro 23 · 🟡 desatualizado 14 · ⚪ sem vínculo 41 · 💰 abaixo da margem 6); tabs Todos/Ativos/Pausados/Com pendência; seleção em massa → [Pausar][Atualizar preço][Re-sync] "com preview e protocolo"; drawer direito com erro traduzido + "▸ técnico", margem, ações Corrigir/Simular/Pausar, timeline. |
| 2b | Detalhe do produto | Breadcrumb Catálogo/; header custo/dimensões/estoque físico-reservado-disponível-seguro/completude; caixa veredicto "Vale a pena vender?"; tabs Dados/Estoque/Anúncios vinculados/Concorrência/Pedidos/Histórico; tabela anúncios vinculados (MLB, Conta, Título, Modal., Preço, vs mercado, Sync, Vendas 30d). |
| 2c | Mercado (novo workspace) | Radar de reprecificação (Meu preço, Menor conc., Mediana, Posição, Vendas 30d, Margem atual, Margem se igualar, Sugestão, Ação) + Oportunidades (SKU sem anúncio × demanda ML) + Monitorados. Freshness "coletado hoje 06:00 · 1×/dia". |
| 2d | Simulador @mercado | Matriz produto × política com pares "@meu / @mercado"; decomposição monospace (preço − comissão − taxa fixa − frete − custo = margem); painel Políticas. |
| 1e | Visão geral | 6 KPIs; "⚠ Precisa de ação agora" por domínio (sync/estoque/margem/vínculo/pedido/pergunta) com contadores e links; saúde da integração; atalhos. |
| 1f | Catálogo | Leitura do ERP (badge "⚑ catálogo: ERP Sankhya"); colunas SKU, Título, GTIN, Var., Custo, Est. fís., Reserv., Disp., Anúncios, Completude; estados explícitos "sem GTIN", "custo?", "sem dimensões". |
| 1g | Vínculos & Importação | 3 grupos: confirmados/revisão/sem correspondência; match SKU → GTIN → título (fuzzy); confiança %, motivo; ações Vincular/Outro/Escolher/Criar produto/Ignorar; batch com "protocolo #12-B". |
| 1h | Estoque | Reconciliação ERP × ML; regra "disponível − seguro → ML"; colunas Físico, Reserv., Disponível, Seguro, Deveria publicar, Publicado ML, Δ dif., Estado; Full = só leitura; [Corrigir no ML] com preview. |
| 1j | Pedidos | Pipeline tabs (Novos/A faturar/A enviar/Enviados/Concluídos/Cancelados/Devoluções); "Status ML → interno" (paid→a faturar); SLA envio; ações Faturar/Etiqueta (FORA do escopo MIS-003). |
| 1k | Integrações & Sync | Cards de conta (OAuth, permissões, mestres por domínio); central de sincronização (Hora, Tipo, Objeto, Direção, Status, Tent., Erro traduzido + ▸ técnico, Reprocessar); "auditoria completa: quem disparou, payload, resposta · retenção 90d". |
| 1l | Editor de anúncio | Wizard 9 etapas; atributos da categoria via API ML, pré-preenchidos do ERP; checklist de publicação; MIS-003 reduz para mini-fluxo "corrigir atributo obrigatório" (decisão P3). |
| 1a–1d | Explorações de Anúncios | Fonte de material para 2a (não construídas separadamente). Vocabulário de sync: sincronizado/com erro/desatualizado/sem vínculo/na fila/sincronizando/pausado. |

### Cross-screen conventions (binding for UI contracts)

- System-of-record badges: `⚑ catálogo/estoque/NF/cadastro: ERP` · `⚑ preço/anúncio: HUB`.
- Context pills top bar: `Empresa: EAN ▾`, `ML: Loja Principal ▾`, `+ conta`; global sync chip (`sync ok · 14:32` / `23 erros`); bell `🔔`.
- Tag colors: ok(green)/warn(amber)/err(red)/mut(gray)/blu(blue: sincronizando, processando, fulfillment, oportunidade).
- Explicit unknown convention: `—`, `custo?`, `sem GTIN`, `sem dimensões`, "sem custo no ERP → não simulado (estado explícito)", "frete indisponível".
- Terminology glossary: produto/produto-mestre (ERP, SKU=CODPROD) ≠ anúncio (ML listing, MLB); vínculo; estoque físico/reservado/disponível/seguro; modalidade (Clássico/Premium/Full); política de preço; mediana ML; completude (cadastro) ≠ qualidade (anúncio); protocolo (recibo de lote); hub = MPC.
- English leaks only as raw provider states ("paid", "Attribute [BRAND] is required") deliberately shown behind "▸ técnico" / "Erro (traduzido)".

### Mutation surfaces (all "via fila de sync com preview e protocolo")

Preço (Aplicar/✎), estoque ([Corrigir no ML]/[Reconciliar tudo]), vínculo (Vincular/Criar produto/batch Aplicar), anúncio (Criar/Corrigir/Pausar/Re-sync/bulk), pedidos (Faturar/Etiqueta — excluded), políticas (editar/+ nova).

## Recommendation

Deck 2 is the build target; deck-1 Anúncios variants are source material only (operator confirmed P1b-3). Mercado UI surfaces (2c, aba Concorrência, colunas @mercado) deferred to enhancement mission (operator decision at P3 review); their data contracts are pinned in MIS-003.

## Impact On Mission

Defines screen scope, state vocabulary, glossary, and mutation surface list consumed by IC-02..IC-05 and milestones M-01..M-06.

## Handoff

- Current status: complete.
- Next owner: Mission Strategist (consumed by mission.md and interface contracts).
- Next action: none.
- Required files/evidence: the wireframe HTML (kept at its Downloads path; not copied into repo).
- Blockers or open decisions: none.
