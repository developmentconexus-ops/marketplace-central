# F-01-product-links-api-gaps

```yaml
id: F-01
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

Fechar gaps da API product_links p/ a tela de vínculos: candidatos com confiança + motivo por âncoras (tabela-verdade IC-01), aprovação em lote LOCAL dentro do próprio product_links (batch-preview dry-run itemizado → batch apply com trilha de auditoria própria; zero write ML), undo de resolução. API existente (candidates/generations/resolutions/workflows) preservada aditivamente. Envelope de mutações NÃO é usado (P5-F-12: selection_resolver do envelope não aceita payload por item p/ candidatos, e o stub `MPC_PROVIDER_WRITES_ENABLED=false` desliga o writer local junto — mecanismo local preserva o resultado aprovado).

## Inputs

- IC-01 — âncoras/confiança (≥85 / 50–84 / <50) / enums; MESMA tabela-verdade do resolver M-02 F-03 (implementações separadas, verdade única via fixtures compartilhadas).
- OpenAPI atual `/product-links/*` (candidates, generations, resolutions approve-candidate/manual-resolve/reject-listing, workflows).
- Fluxos de resolução existentes (approve-candidate/manual-resolve/reject-listing) — o batch reusa a MESMA lógica de aplicação por item.
- Migrations bloco 0065–0069 (`apps/server_core/migrations/`) — inclui tabela de auditoria do batch. Toda tabela nova e todo ALTER deste bloco carrega/preserva `tenant_id`; toda query nova escopa `tenant_id` (invariante da missão).

## Expected Output

- Candidatos enriquecidos: `confidence: number (0-100)`, `confidence_band: ALTA|MEDIA|BAIXA`, `reasons: [{anchor, direction: FOR|AGAINST|UNAVAILABLE, detail}]` — ALTER aditivo em `product_link_candidates` + backfill de score p/ candidatos existentes na geração seguinte (não retroativo).
- `POST /product-links/resolutions/batch-preview` `{approvals: [{candidate_id}]}` ⇒ 200 dry-run itemizado `{items: [{candidate_id, status OK|FAILED, cause?}]}` — NADA aplicado.
- `POST /product-links/resolutions/batch` `{approvals: [{candidate_id}]}` ⇒ 200 `{batch_id, applied: [...], failed: [{candidate_id, cause}]}`: aplica cada item LOCALMENTE pela mesma lógica de `approve_candidate` individual; item falho não bloqueia os demais; registro de auditoria por batch (quem, quando, itens, resultados) em tabela própria do módulo — "quem" no MVP sem auth = ator fixo `operator` (valor literal; auth real = MIS-005). Envelope `/mutations*` NÃO participa — nenhuma chamada ML por construção.
- `POST /product-links/resolutions/{id}/undo` ⇒ 200: reverte resolução (candidato volta à fila), auditável.
- Seção OpenAPI `/product-links/*` aditiva + `sdk-runtime/src/productLinks.ts`.
- EARS: While candidato tem EAN+marca a favor, when geração computa confiança, the sistema shall atribuir banda ALTA com ambas as âncoras em `reasons`. While lote de N aprovações é submetido, when POST batch-preview chega, the sistema shall responder dry-run itemizado sem aplicar nada. While lote confirmado é submetido, when POST batch chega, the sistema shall aplicar item a item localmente, responder applied/failed itemizados e gravar auditoria do batch. While resolução é desfeita, when undo chega, the sistema shall restaurar o candidato na fila com histórico da resolução preservado.

## Inputs/Outputs

Shapes detalhados: derivar do OpenAPI atual + IC-01. Batch é local ao product_links: preview 200 dry-run itemizado; apply 200 com `applied`/`failed` itemizados (falha parcial NÃO é erro HTTP). Status codes fixos: candidato inexistente no lote ⇒ item FAILED itemizado, lote continua; undo de resolução inexistente ⇒ 404; undo de resolução já desfeita ⇒ 409 `ALREADY_UNDONE`.

## Negative Scenarios

- Aprovar candidato já resolvido ⇒ 409 `ALREADY_RESOLVED` (individual); no lote ⇒ item FAILED itemizado.
- Lote vazio ⇒ 422.
- Undo após nova resolução do mesmo produto ⇒ 409 `SUPERSEDED` (não reverte silenciosamente o vínculo novo).
- Candidato com EAN em colisão (flag M-01) ⇒ banda máxima MEDIA + reason AGAINST citando colisão.

## Constraints

- ZERO write ML — resoluções são estado local; sincronização com ML é MIS-005.
- `modules/mutations/**` é forbidden — consumo via API/port público apenas.
- Enums/strings EXATOS IC-01.

## Ownership

- Owned paths: `modules/product_links/**`, `apps/server_core/migrations/0065*–0067*` (+ fixture `apps/server_core/internal/platform/migrate/runner_test.go` bump), seção `/product-links/*` OpenAPI (aditiva), `sdk-runtime/src/productLinks.ts`.
- Forbidden paths: `modules/mutations/**`, `modules/market/**`, barrel SDK (hub), `apps/web/**` (F-02).
- Parallel-safe with: none — F-02 depende da seção OpenAPI deste F (pode iniciar quando ela estiver comitada).

## Validation Expectations

- Fixtures IC-01 no ranking: ≥8 casos com `confidence_band` + `reasons` exatos esperados.
- Transcript batch: batch-preview de 3 aprovações (1 inválida) ⇒ 200 com 2 OK + 1 FAILED com causa e NADA aplicado; batch apply ⇒ 200 com 2 applied + 1 failed; GET fila não contém os 2; linha de auditoria do batch inspecionável.
- Transcript undo ⇒ 200, candidato de volta na fila; segundo undo ⇒ 409 `ALREADY_UNDONE`.
- Zero chamadas a connectors nos testes (mock provider não tocado — prova local-only).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-04).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` com transcripts acima.
- Blockers or open decisions: none.
