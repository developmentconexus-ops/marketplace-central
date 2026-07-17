# F-04 — preview-confirm-ui

Feature: `F-04` · Milestone: `M-03` · Mission: `MIS-003`

## Base assumptions

- F-03 and the M-02 frontend seam are accepted in the current base. The current OpenAPI and SDK already expose `createMutation`, `previewMutation`, `approveMutation`, `cancelMutation`, `retryMutationFailures`, `getMutation`, and `listMutationItems`; `web-query` already exports `mutationsQueryKeys`, `QUERY_STALE_TIME.mutations`, `invalidateAfterMutation`, and `failureCopy`. This feature consumes those surfaces without editing them.
- The Anúncios launch surface sends an explicit selection: `{mode: "explicit", listing_ids: [...]}`. The actor label is the contract value `operator_supplied_unverified`; installations come only from `useInstallation()`.
- The six actions map to `price_update`, `stock_correct`, `listing_pause`, `listing_resync`, `link_apply`, and `listing_edit`, using the IC-03 intent shapes. The server response remains the only lifecycle authority; terminal states are `applied`, `partially_failed`, `failed_preserved`, and `cancelled`.
- The modal owns only its interaction step (`intent → previewing → preview-shown → approving → applying → terminal`). Every request, protocol, preview, and item collection is managed by TanStack Query; every protocol/items query uses `mutationsQueryKeys`. A feature-local protocol hook shares the 2-second polling rule and a per-`QueryClient`/protocol terminal-invalidation guard between modal and detail consumers.
- No dependency, contract, SDK, context, layout, state-component, or `web-query` change is needed. The only seam edit is replacing the assigned `/protocolos/:protocolId` element in `AppRouter.tsx`; no additive contract lock is required.
- Write DAG: S1 → S2 → S3 and S2 → S4 → S5. S3 and S4 are disjoint after S2 and may be reviewed independently.

## Verification map

| Brief criterion | Automated proof | QA proof |
| --- | --- | --- |
| Six intent forms, create/preview, before→after rows, totals, confirm gated behind preview | `MutationPreviewModal.test.tsx` | Open each Anúncios action and capture the preview table before confirmation |
| Explicit confirm, one approve on double-click, stale-preview recovery without re-approve | `MutationPreviewModal.test.tsx` | Confirm once; exercise stale preview and regenerate manually |
| Poll every 2s while non-terminal, reconnect/reload recovery, exactly-one terminal invalidation | `useMutationProtocol.test.tsx`, `ProtocoloPage.test.tsx` | Close modal during apply; F5 `/protocolos/MP-000042`; observe terminal result |
| `selection_too_large` shows count/guidance and cancels the draft | `MutationPreviewModal.test.tsx` | Preview an over-limit explicit selection and verify no draft remains |
| Result/detail failure copy and raw provider message only behind “▸ técnico” | `MutationPreviewModal.test.tsx`, `ProtocoloPage.test.tsx` | Capture result and expanded technical detail |
| Retry creates a new protocol, navigates to it, and shows `retried_from` back-link | `ProtocoloPage.test.tsx`, `AppRouter.test.tsx` | Retry failures and follow both new- and old-protocol links |
| Cancel is available only in intent/preview-shown and calls cancel for an existing draft | `MutationPreviewModal.test.tsx` | Cancel before preview and after preview; verify the modal closes |

### S1 — Intent and preview state machine

Goal: Build the shared modal’s six pt-BR intent forms and the create→preview half of the local interaction state machine, including totals and before/after rows before confirmation is enabled.

Files touched:

- `apps/web/src/pages/mutations/mutationPresentation.ts` (NEW)
- `apps/web/src/pages/mutations/MutationIntentForm.tsx` (NEW)
- `apps/web/src/pages/mutations/MutationItemsTable.tsx` (NEW)
- `apps/web/src/pages/mutations/MutationPreviewModal.tsx` (NEW)
- `apps/web/src/pages/mutations/MutationPreviewModal.test.tsx` (NEW)

Test-first spec: Start with failing component cases proving that no confirmation control is enabled before `previewMutation` resolves; each action submits its exact IC-03 intent plus explicit selected IDs, installation, and actor; preview totals and before→after rows render afterward. Add negative cases for canceling an existing draft and for `selection_too_large`: show the selected count and pt-BR narrowing guidance, call `cancelMutation` once, and leave no confirmable draft.

Done-when: All six forms use one modal; create/preview/cancel are `useMutation` operations; the modal never derives a protocol lifecycle; loading/error rendering uses the shared state vocabulary; cancel before draft creation closes locally, while cancel after creation calls the SDK; no direct fetch or installation lookup exists.

Complexity: `complex`

Dependencies: F-03 and M-02 seam only.

Validation kind: component.

Commands: `npm run test --workspace @marketplace-central/web -- --run src/pages/mutations/MutationPreviewModal.test.tsx`

Expected artifacts: the five source/test files above and captured red→green Vitest output.

Write set: exactly the five paths above.

Open questions: none.

### S2 — Approval, polling, terminal result, and invalidation guard

Goal: Complete approval and application: explicit checkbox, double-click-safe `{execute:true}`, 2-second server polling, stale-preview recovery, terminal summary, and once-only crosswalk invalidation.

Files touched:

- `apps/web/src/pages/mutations/useMutationProtocol.ts` (NEW)
- `apps/web/src/pages/mutations/useMutationProtocol.test.tsx` (NEW)
- `apps/web/src/pages/mutations/MutationResultSummary.tsx` (NEW)
- `apps/web/src/pages/mutations/MutationPreviewModal.tsx`
- `apps/web/src/pages/mutations/MutationPreviewModal.test.tsx`

Test-first spec: First fail tests that (1) keep approve disabled until the explicit checkbox is checked, (2) turn two immediate confirm clicks into one `approveMutation(protocolId, {execute:true})`, (3) on 409/`preview_stale` reset confirmation and return to preview with exactly “Prévia expirada. Gere novamente.” and no automatic second approve, (4) poll `getMutation` every 2000 ms only while the returned server state is non-terminal and resume through Query retry/reconnect behavior, and (5) call `invalidateAfterMutation(queryClient, protocol.type)` once when either consumer first observes that protocol terminal. Add a result failure proving visible copy comes from `failureCopy(code)` and not `message_pt`.

Done-when: The shared hook uses `mutationsQueryKeys.detail(protocolId)` and `QUERY_STALE_TIME.mutations`; terminal invalidation is guarded per QueryClient/protocol across modal and detail consumers; item results use `mutationsQueryKeys.items(protocolId)`; polling stops on terminal truth; closing the modal does not cancel an approved/applying protocol; failures and a `/protocolos/:id` link render in the terminal result.

Complexity: `complex`

Dependencies: S1.

Validation kind: component with fake timers and QueryClient spies.

Commands: `npm run test --workspace @marketplace-central/web -- --run src/pages/mutations/useMutationProtocol.test.tsx src/pages/mutations/MutationPreviewModal.test.tsx`

Expected artifacts: the three new files, two updated modal files, and captured red→green Vitest output.

Write set: exactly the five paths above.

Open questions: none.

### S3 — Enable all Anúncios bulk launchers

Goal: Replace the disabled bulk-action affordances with all six action launchers and hand the current installation plus the accumulated opaque listing IDs to the shared modal.

Files touched:

- `apps/web/src/pages/mutations/MutationBulkActions.tsx` (NEW)
- `apps/web/src/pages/AnunciosPage.tsx`
- `apps/web/src/pages/AnunciosSelection.test.tsx`

Test-first spec: Replace the existing “keeps bulk actions disabled” test with failing cases that require Atualizar preço, Corrigir estoque, Pausar, Ressincronizar, Vincular, and Editar to remain disabled at zero selection, open the matching modal after selection, and pass IDs accumulated across pages without fetching installations. Prove changing installation clears selection and closes any unsubmitted modal.

Done-when: The six pt-BR actions are available only with a non-empty selection; they map to the six contract mutation types; the modal receives `installationId` solely from `useInstallation()` and a snapshot of `selectedIds`; existing pagination, selection, refresh, and detail-panel behavior stays green.

Complexity: `standard`

Dependencies: S2.

Validation kind: component/regression.

Commands: `npm run test --workspace @marketplace-central/web -- --run src/pages/AnunciosSelection.test.tsx src/pages/AnunciosPage.test.tsx src/pages/AnunciosTable.test.tsx`

Expected artifacts: the new launcher component, the Anúncios integration diff, and captured red→green Vitest output.

Write set: exactly the three paths above.

Open questions: none.

### S4 — Reload-safe protocolo detail

Goal: Build `/protocolos/:protocolId` content with protocol header, timestamps/state, retried-from link, paged items, failure presentation, and server-truth polling that works when mounted directly mid-apply.

Files touched:

- `apps/web/src/pages/mutations/ProtocoloPage.tsx` (NEW)
- `apps/web/src/pages/mutations/ProtocoloPage.test.tsx` (NEW)

Test-first spec: Start at a direct protocol route with `getMutation` returning `applying`, advance fake timers by 2000 ms to a terminal response, and assert polling stops and the shared invalidation guard fires once. Add header/timestamp/type/state assertions, loading/error states, paged `listMutationItems` via `mutationsQueryKeys.items(protocolId)`, before/after values with `UnknownValue` for null facts, `failureCopy(code)` output, hidden raw `message_provider`, and disclosure only after clicking “▸ técnico”.

Done-when: Direct mount/F5 needs no modal memory; protocol and items are exclusively TanStack queries using the existing builders; the header exposes type, actor, timestamps, server state, and `retried_from`; item pages are ordered/displayed without inventing unknown values; provider text is never visible outside the technical disclosure.

Complexity: `complex`

Dependencies: S2.

Validation kind: component with fake timers and direct-route mount.

Commands: `npm run test --workspace @marketplace-central/web -- --run src/pages/mutations/ProtocoloPage.test.tsx src/pages/mutations/useMutationProtocol.test.tsx`

Expected artifacts: the detail page, its test, and captured red→green Vitest output.

Write set: exactly the two paths above.

Open questions: none.

### S5 — Retry navigation and assigned route swap

Goal: Finish the detail workflow by retrying failed items into a new protocol and mount the real detail page at the single AppRouter row assigned to M-03.

Files touched:

- `apps/web/src/pages/mutations/ProtocoloPage.tsx`
- `apps/web/src/pages/mutations/ProtocoloPage.test.tsx`
- `apps/web/src/app/AppRouter.tsx`
- `apps/web/src/app/AppRouter.test.tsx`

Test-first spec: Fail a detail-page test where “Repetir itens com falha” calls `retryMutationFailures(oldId)` once, navigates to `/protocolos/<newId>` while preserving the installation query string, loads the new protocol, and renders its `retried_from` link back. Replace the router placeholder assertion with a failing deep-link test for the real protocol page.

Done-when: Retry is mutation-pending guarded, navigation uses the SDK-returned new `protocol_id`, the new detail links back through `retried_from`, and `AppRouter.tsx` changes only its protocol route element plus the import needed to render it. No other route, router/context/layout, or `web-query` seam changes.

Complexity: `standard`

Dependencies: S4.

Validation kind: component/router regression.

Commands: `npm run test --workspace @marketplace-central/web -- --run src/pages/mutations/ProtocoloPage.test.tsx src/app/AppRouter.test.tsx`; final feature sweep: `npm run test --workspace @marketplace-central/web`; `npm run build --workspace @marketplace-central/web`.

Expected artifacts: retry/navigation and route-row diffs, full web Vitest output, web build output, and later `validation.md` browser screenshots for preview, applying, terminal failure, detail, and F5 recovery.

Write set: exactly the four paths above.

Open questions: none.

## Open questions

None.
