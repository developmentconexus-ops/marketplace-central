# F-04-preview-confirm-ui

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-03 (lifecycle, poll 2s, failure codes), IC-05 (mutationsQueryKeys, invalidateAfterMutation, failureCopy, `/protocolos/:protocolId` route, state components), R-01 mutation-surfaces inventory.

## Milestone

M-03 mutation-envelope-writes. Depends on F-03 (endpoints) + M-02 F-02 (crosswalk/failureCopy).

## Brief

Wire the write UX: enable Anúncios (2a) bulk action buttons (Atualizar preço, Corrigir estoque, Pausar, Ressincronizar, Vincular, Editar) opening the preview/confirm modal — intent form → `createMutation` → `previewMutation` → per-item before/after preview table with totals → explicit confirm checkbox → `approveMutation {execute:true}` → progress view polling `getMutation` at 2s until terminal → result summary with per-item failures (failureCopy strings) + link to `/protocolos/:id`. Build `/protocolos/:protocolId` detail page: header (type, actor, timestamps, state), items table (before/after, failure code+copy, "▸ técnico" `message_provider`), "Repetir itens com falha" button calling retry → navigates to NEW protocol. On terminal state call `invalidateAfterMutation(qc, type)`.

EARS:
- While the preview modal is open, when preview returns, the operator shall see item count, totals, and per-item before→after rows before any confirm control activates.
- While a protocol is applying, when the operator closes the modal, progress shall continue server-side and `/protocolos/:id` shall show live state (poll while non-terminal).
- While approve returns 409 `preview_stale`, when confirming, the modal shall return to preview step with copy "Prévia expirada. Gere novamente." — never silently re-approve.
- While a terminal state arrives, when the modal or detail page observes it, `invalidateAfterMutation` shall fire exactly once for the protocol type.
- While retry succeeds, when the new protocol is created, the UI shall navigate to the new `/protocolos/:id` and show `retried_from` link back.

## Inputs

- IC-03 lifecycle + endpoints (via sdk-runtime from F-03), IC-05 keys/crosswalk/failureCopy/state components, R-01 wireframe mutation surfaces (modal shape, protocolo screen), M-02 F-03 selection state (ids handed to createMutation).

## Expected Output

- Preview/confirm modal component (shared across types; per-type intent form section).
- `/protocolos/:protocolId` page; protocolos list accessible from Integrações later (M-05) — this feature adds only the detail route per IC-05.
- Poll via TanStack Query `refetchInterval: 2000` while non-terminal; stops at terminal.
- Component tests: preview-before-confirm gating, stale-preview recovery, single invalidation on terminal, retry navigation, failure copy rendering.

## Constraints

- All server state via `mutationsQueryKeys`; no direct fetch; no inline invalidation (crosswalk only).
- Confirm requires explicit checkbox + button (two-step, wireframe pattern) — no one-click apply.
- pt-BR copy; failure strings only from failureCopy module.
- Does not modify router/context/web-query seam files beyond adding the route element at the IC-05-assigned row.

## Negative Scenarios

- `selection_too_large` on preview → modal error state with count and guidance to narrow filter; no protocol left in draft limbo (cancel fired).
- Network drop mid-poll → poll resumes on reconnect (query retry), state from server not client guess.
- Double-click confirm → single approve call (button disabled on first click).

## Interaction Model

Modal state machine: intent → previewing → preview-shown → approving → applying(poll) → terminal(result). Cancel allowed in intent/preview-shown (calls cancelMutation). Modal state is component-local; protocol truth is always server response — UI never computes lifecycle.

## Validation Expectations

- Vitest output: gating, stale-recovery, single-invalidation, double-click tests green.
- Browser proof: full stub-adapter flow screenshots — preview table with before/after, progress, result with failure pt-BR copy, protocolo detail page.
- Reload proof: F5 on `/protocolos/MP-000042` mid-apply keeps polling and reaches terminal.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-03 accepted).
- Next action: compile context pack; read IC-03/IC-05 + modal wireframe section of R-01.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
