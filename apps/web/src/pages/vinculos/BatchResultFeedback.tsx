import type { ApplyProductLinkBatchResponse } from "@marketplace-central/sdk-runtime";

export interface BatchResultFeedbackProps {
  result: ApplyProductLinkBatchResponse;
  onDismiss: () => void;
  onViewResolved: () => void;
}

/**
 * Inline banner (page-local — no toast primitive exists in packages/ui).
 * Partial failure is normal: applied[] and failed[]{candidate_id,cause} are
 * both rendered itemized; a single failure never hides the applied results.
 */
export function BatchResultFeedback({
  result,
  onDismiss,
  onViewResolved,
}: BatchResultFeedbackProps) {
  return (
    <div role="status" className="rounded-card border border-border bg-surface-2 px-4 py-3 text-sm">
      <div className="flex items-center justify-between gap-2">
        <p className="font-medium text-ink">
          {result.applied.length} aplicado(s), {result.failed.length} falha(s)
        </p>
        <button
          type="button"
          onClick={onDismiss}
          className="text-xs font-medium text-faint hover:text-ink"
        >
          Fechar
        </button>
      </div>

      {result.failed.length > 0 ? (
        <ul className="mt-2 space-y-1 text-xs text-warn" data-testid="batch-result-failed">
          {result.failed.map((item) => (
            <li key={item.candidate_id}>
              {item.candidate_id}: {item.cause}
            </li>
          ))}
        </ul>
      ) : null}

      {result.applied.length > 0 ? (
        <ul className="mt-2 space-y-1 text-xs text-accent-ink" data-testid="batch-result-applied">
          {result.applied.map((item) => (
            <li key={item.candidate_id}>{item.candidate_id}: aplicado</li>
          ))}
        </ul>
      ) : null}

      {result.applied.length > 0 ? (
        <button
          type="button"
          onClick={onViewResolved}
          className="mt-2 inline-flex text-xs font-medium text-accent hover:underline"
        >
          Ver em Resolvidos
        </button>
      ) : null}
    </div>
  );
}
