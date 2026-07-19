import type { ApplyProductLinkBatchResponse, ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { BatchPreviewModal } from "./BatchPreviewModal";
import { BatchResultFeedback } from "./BatchResultFeedback";
import { BulkBar } from "./BulkBar";
import { QueueRow } from "./QueueRow";
import { VinculoDrawer } from "./VinculoDrawer";
import { rejectInputForCandidate, useVinculosQueue } from "./useVinculosQueue";

const CANDIDATE_PARAM = "candidate";

export interface QueueTabProps {
  installationId: string;
  onViewResolved?: () => void;
}

export function QueueTab({ installationId, onViewResolved }: QueueTabProps) {
  const { queueQuery, approve, reject, conflict } = useVinculosQueue(installationId);
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [batchModalOpen, setBatchModalOpen] = useState(false);
  const [batchCandidateIds, setBatchCandidateIds] = useState<string[]>([]);
  const [batchResult, setBatchResult] = useState<ApplyProductLinkBatchResponse | null>(null);

  const openCandidateId = searchParams.get(CANDIDATE_PARAM);

  const openDrawer = (candidateId: string) => {
    setSearchParams(
      (prev) => {
        const next = new URLSearchParams(prev);
        next.set(CANDIDATE_PARAM, candidateId);
        return next;
      },
      { replace: false },
    );
  };

  const closeDrawer = () => {
    setSearchParams(
      (prev) => {
        const next = new URLSearchParams(prev);
        next.delete(CANDIDATE_PARAM);
        return next;
      },
      { replace: true },
    );
  };

  const handleApprove = (candidate: ProductLinkCandidateItem) => {
    approve.mutate({ candidate_id: candidate.candidate_id });
  };

  const handleReject = (candidate: ProductLinkCandidateItem) => {
    reject.mutate(rejectInputForCandidate(candidate));
    closeDrawer();
  };

  // Row-level "Ignorar" (NO_CANDIDATE) — reject the listing without opening the drawer.
  const handleRejectRow = (candidate: ProductLinkCandidateItem) => {
    reject.mutate(rejectInputForCandidate(candidate));
  };

  const toggleSelect = (candidateId: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(candidateId)) next.delete(candidateId);
      else next.add(candidateId);
      return next;
    });
  };

  const clearSelection = () => setSelectedIds(new Set());

  // Selection is cleared right after enqueue (the moment the batch is sent to
  // preview) — the modal keeps its own snapshot (batchCandidateIds) so it
  // still has the right candidates to dry-run even though selectedIds is
  // already empty. Also cleared (defensively) after apply.
  const openBatchPreview = () => {
    setBatchCandidateIds(Array.from(selectedIds));
    setBatchModalOpen(true);
    clearSelection();
  };

  const closeBatchPreview = () => setBatchModalOpen(false);

  const handleBatchApplied = (result: ApplyProductLinkBatchResponse) => {
    setBatchResult(result);
    setBatchModalOpen(false);
    clearSelection();
  };

  if (queueQuery.isPending) return <LoadingState />;
  if (queueQuery.isError) return <ErrorState onRetry={() => void queueQuery.refetch()} />;

  const items = queueQuery.data?.items ?? [];
  if (items.length === 0) return <EmptyState />;

  const mutating = approve.isPending || reject.isPending;

  return (
    <div className="mt-3 flex flex-col gap-3">
      {conflict ? (
        <div
          role="alert"
          className="rounded-control border border-amber/30 bg-amber-soft px-3 py-2 text-sm text-amber"
        >
          Este anúncio já foi resolvido por outra pessoa. A fila foi atualizada.
        </div>
      ) : null}

      {batchResult ? (
        <BatchResultFeedback
          result={batchResult}
          onDismiss={() => setBatchResult(null)}
          onViewResolved={() => {
            setBatchResult(null);
            onViewResolved?.();
          }}
        />
      ) : null}

      <BulkBar selectedCount={selectedIds.size} onPreview={openBatchPreview} onClear={clearSelection} />

      <div className="overflow-x-auto">
        <table className="w-full min-w-[900px] border-collapse text-left text-sm">
          <caption className="sr-only">Fila de candidatos de vínculo (anúncio → produto sugerido)</caption>
          <thead className="border-b border-border bg-surface-2 text-xs font-medium tracking-[0.04em] text-faint">
            <tr>
              <th className="px-3 py-3" scope="col">
                <span className="sr-only">Selecionar</span>
              </th>
              <th className="px-3 py-3" scope="col">Anúncio ML</th>
              <th className="px-3 py-3" scope="col">SKU ML</th>
              <th className="px-3 py-3" scope="col">Produto sugerido</th>
              <th className="px-3 py-3" scope="col">SKU HUB</th>
              <th className="px-3 py-3" scope="col">GTIN</th>
              <th className="px-3 py-3" scope="col">Confiança</th>
              <th className="px-3 py-3" scope="col">Motivo</th>
              <th className="px-3 py-3 text-right" scope="col">Ação</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-2">
            {items.map((candidate) => (
              <QueueRow
                key={candidate.candidate_id}
                candidate={candidate}
                onOpen={openDrawer}
                onApprove={handleApprove}
                onReject={handleRejectRow}
                pending={mutating}
                selected={selectedIds.has(candidate.candidate_id)}
                onToggleSelect={toggleSelect}
              />
            ))}
          </tbody>
        </table>
      </div>

      <VinculoDrawer
        candidateId={openCandidateId}
        candidates={items}
        onClose={closeDrawer}
        onApprove={(candidate) => {
          handleApprove(candidate);
          closeDrawer();
        }}
        onReject={handleReject}
        pending={mutating}
      />

      <BatchPreviewModal
        open={batchModalOpen}
        candidateIds={batchCandidateIds}
        onClose={closeBatchPreview}
        onApplied={handleBatchApplied}
      />
    </div>
  );
}
