import { useEffect, useRef } from "react";
import type { ApplyProductLinkBatchResponse } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { useVinculosBatch } from "./useVinculosBatch";

export interface BatchPreviewModalProps {
  open: boolean;
  candidateIds: string[];
  onClose: () => void;
  onApplied: (result: ApplyProductLinkBatchResponse) => void;
}

/**
 * Page-local `role="dialog"` modal (MutationPreviewModal precedent).
 *
 * Contract fidelity: opening the modal fires ONLY `previewProductLinkBatch`
 * (a pure dry-run — nothing is applied). `applyProductLinkBatch` is fired
 * exclusively from the "prosseguir só com válidos" button, and only with the
 * OK-status subset of the preview — predicted failures are never sent to
 * apply.
 */
export function BatchPreviewModal({
  open,
  candidateIds,
  onClose,
  onApplied,
}: BatchPreviewModalProps) {
  const { preview, apply } = useVinculosBatch();
  const requestedFor = useRef<string | null>(null);

  useEffect(() => {
    if (!open) {
      requestedFor.current = null;
      return;
    }
    const key = candidateIds.join(",");
    if (requestedFor.current === key) return;
    requestedFor.current = key;
    preview.mutate(candidateIds.map((candidate_id) => ({ candidate_id })));
    // preview/apply mutation identities are stable across renders; only re-run
    // the dry-run when the modal opens for a (possibly new) selection.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, candidateIds]);

  if (!open) return null;

  const items = preview.data?.items ?? [];
  const validApprovals = items
    .filter((item) => item.status === "OK")
    .map((item) => ({ candidate_id: item.candidate_id }));
  const failedItems = items.filter((item) => item.status === "FAILED");

  const handleClose = () => {
    preview.reset();
    apply.reset();
    onClose();
  };

  const handleApply = () => {
    apply.mutate(validApprovals, {
      onSuccess: (result) => onApplied(result),
    });
  };

  return (
    <div
      className="fixed inset-0 z-40 flex items-center justify-center bg-ink/50 p-4"
      role="presentation"
    >
      <section
        role="dialog"
        aria-modal="true"
        aria-labelledby="batch-preview-modal-title"
        className="flex max-h-[90vh] w-full max-w-3xl flex-col overflow-hidden rounded-card border border-border bg-surface shadow-lg"
      >
        <header className="border-b border-border px-5 py-4">
          <h2 id="batch-preview-modal-title" className="text-lg font-semibold text-ink">
            Pré-visualizar aprovação em lote
          </h2>
          <p className="mt-1 text-sm text-muted">
            {candidateIds.length} candidato(s) selecionado(s)
          </p>
        </header>

        <div className="min-h-0 flex-1 overflow-y-auto p-5">
          {preview.isPending ? <LoadingState /> : null}
          {preview.isError ? (
            <ErrorState
              onRetry={() => preview.mutate(candidateIds.map((candidate_id) => ({ candidate_id })))}
            />
          ) : null}
          {preview.isSuccess ? (
            <div className="space-y-4">
              <div className="flex flex-wrap gap-2" aria-label="Totais da prévia">
                <span className="rounded-full bg-accent-soft px-3 py-1 text-sm font-medium text-accent-ink">
                  {validApprovals.length} válido(s)
                </span>
                <span className="rounded-full bg-warn-soft px-3 py-1 text-sm font-medium text-warn">
                  {failedItems.length} previsto(s) para falhar
                </span>
              </div>
              <ul
                className="divide-y divide-border-2 rounded-card border border-border"
                data-testid="batch-preview-items"
              >
                {items.map((item) => (
                  <li
                    key={item.candidate_id}
                    className="flex items-center justify-between px-4 py-2 text-sm"
                  >
                    <span className="font-mono font-medium text-ink">{item.candidate_id}</span>
                    {item.status === "OK" ? (
                      <span className="text-xs font-medium text-accent-ink">OK</span>
                    ) : (
                      <span className="text-xs text-warn">{item.cause ?? "falha prevista"}</span>
                    )}
                  </li>
                ))}
              </ul>
              {apply.isError ? (
                <ErrorState detail="Não foi possível aplicar o lote." onRetry={handleApply} />
              ) : null}
            </div>
          ) : null}
        </div>

        <footer className="flex flex-wrap justify-end gap-2 border-t border-border px-5 py-4">
          <button
            type="button"
            onClick={handleClose}
            disabled={apply.isPending}
            className="rounded-control border border-border px-4 py-2 text-sm font-medium text-muted hover:bg-surface-2 hover:text-ink disabled:opacity-50"
          >
            Cancelar
          </button>
          {preview.isSuccess ? (
            <button
              type="button"
              onClick={handleApply}
              disabled={validApprovals.length === 0 || apply.isPending}
              className="rounded-control bg-accent px-4 py-2 text-sm font-medium text-accent-ink disabled:cursor-not-allowed disabled:opacity-50"
            >
              Prosseguir só com válidos
            </button>
          ) : null}
        </footer>
      </section>
    </div>
  );
}
