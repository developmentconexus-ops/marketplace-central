import type { ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { useSearchParams } from "react-router-dom";
import { QueueRow } from "./QueueRow";
import { VinculoDrawer } from "./VinculoDrawer";
import { rejectInputForCandidate, useVinculosQueue } from "./useVinculosQueue";

const CANDIDATE_PARAM = "candidate";

export interface QueueTabProps {
  installationId: string;
}

export function QueueTab({ installationId }: QueueTabProps) {
  const { queueQuery, approve, reject, conflict } = useVinculosQueue(installationId);
  const [searchParams, setSearchParams] = useSearchParams();

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
          className="rounded-lg border border-amber-300 bg-amber-50 px-3 py-2 text-sm text-amber-900"
        >
          Este anúncio já foi resolvido por outra pessoa. A fila foi atualizada.
        </div>
      ) : null}

      <div className="overflow-x-auto">
        <table className="w-full min-w-[900px] border-collapse text-left text-sm">
          <caption className="sr-only">Fila de candidatos de vínculo</caption>
          <thead className="border-b border-slate-200 text-xs font-medium text-slate-500">
            <tr>
              <th className="px-3 py-3" scope="col">Produto</th>
              <th className="px-3 py-3" scope="col">Melhor candidato</th>
              <th className="px-3 py-3" scope="col">Sinais</th>
              <th className="px-3 py-3" scope="col">Confiança</th>
              <th className="px-3 py-3 text-right" scope="col">Ações</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {items.map((candidate) => (
              <QueueRow
                key={candidate.candidate_id}
                candidate={candidate}
                onOpen={openDrawer}
                onApprove={handleApprove}
                pending={mutating}
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
    </div>
  );
}
