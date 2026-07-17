# F-02-vinculos-screen

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-04-vinculos-import-ui.

## Brief

Tela `/vinculos` real (substitui placeholder W1 via `routes/vinculos.tsx`): fila de resolução com candidatos rankeados, drawer de comparação lado-a-lado (nosso produto vs anúncio/catálogo ML) com chips de âncora, aprovação individual, bulk select→preview (dry-run)→aplicar, tab Resolvidos com undo.

## Inputs

- Design handoff (tela Vínculos) + leitura R-02 `research/design-screens-2026-07-17.md`.
- `sdk-runtime/src/productLinks.ts` (F-01) — client tipado.
- IC-05 — seam `routes/vinculos.tsx`, tokens, primitivas (ConflictTag, EmptyState, banda de confiança), TanStack Query via web-query.
- IC-01 — bandas/motivos (apresentação fiel: motivo SEMPRE visível, % nunca sozinho).

## Expected Output

- Tabs `Fila` / `Resolvidos`; KPIs no topo (pendentes, alta confiança, resolvidos hoje).
- Linha da fila: produto (CODPROD + descrição + EAN/refforn), melhor candidato (título + preço + banda %), chips de âncora (verde FOR / vermelho AGAINST / cinza UNAVAILABLE), ações aprovar/rejeitar/abrir.
- Drawer: comparação campo a campo, todos os candidatos rankeados, decisão manual (manual-resolve / reject-listing das APIs existentes).
- Bulk: seleção múltipla ⇒ modal de preview via `batch-preview` (itens + falhas previstas, nada aplicado) ⇒ confirma ⇒ `batch` apply ⇒ toast com applied/failed + link p/ tab Resolvidos; fila refetch após apply.
- Resolvidos: lista com quem/quando/como (auto/manual) + undo.
- EARS: While item aprovado individualmente, when confirmação retorna 2xx, the sistema shall remover o item da fila e incrementar KPI resolvidos (refetch, não mutação otimista de lista). While produto sem candidato, when fila renderiza, the sistema shall mostrar estado NO_CANDIDATE com ação "coletar de novo" (dispara geração). While bulk preview contém falhas, when modal renderiza, the sistema shall listar cada falha com causa e permitir prosseguir só com os válidos.

## Negative Scenarios

- API 409 ALREADY_RESOLVED na aprovação ⇒ toast de conflito + refetch da fila (item some, sem crash).
- Fila vazia ⇒ EmptyState com orientação (importar planilha / gerar candidatos).
- Undo 409 SUPERSEDED ⇒ mensagem explicando vínculo mais novo.

## State / Interaction Model

- Server-state 100% TanStack Query (keys: `['product-links','queue',filters]`, `['product-links','resolved']`); invalidação: aprovar individual/undo/batch apply ⇒ invalidate queue+resolved+KPIs.
- Seleção bulk é estado local da página (não URL, não storage); limpa após enfileirar.
- Drawer aberto = candidato na URL (`?candidate=`) p/ deep-link/F5.
- Sem estado otimista de escrita no MVP — sempre refetch pós-2xx (simplicidade > latência aqui).

## Constraints

- Só `routes/vinculos.tsx` + `pages/vinculos/**` — AppRouter/Layout intocados (IC-05).
- Motivo de confiança sempre renderizado com o % (doutrina anti-caixa-preta do research).
- Sem chamadas diretas fetch — só client SDK.

## Ownership

- Owned paths: `apps/web/src/pages/vinculos/**` (novo), `apps/web/src/routes/vinculos.tsx` (conteúdo).
- Forbidden paths: `apps/web/src/app/**`, `packages/ui/**` (mudança de primitiva = REQUEST), outros `routes/*`, `sdk-runtime/**`.
- Parallel-safe with: none — depends on F-01 (`productLinks.ts`) + M-03 seams.

## Validation Expectations

- Screenshot fila com ≥3 itens: bandas distintas + chips de âncora visíveis; drawer aberto com comparação.
- Transcript navegador: aprovar individual ⇒ request/response 2xx + item fora da fila após refetch.
- Bulk 3 itens ⇒ modal preview (dry-run visível) ⇒ apply ⇒ toast applied/failed ⇒ fila reduzida + itens na tab Resolvidos.
- Deep-link `/vinculos?candidate=X` + F5 ⇒ drawer aberto no candidato X.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-04, após F-01 OpenAPI).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots + transcripts.
- Blockers or open decisions: none.
