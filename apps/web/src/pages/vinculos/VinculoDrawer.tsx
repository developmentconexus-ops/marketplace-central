import type {
  ProductLinkCandidateItem,
  ProductLinkConfidenceBand,
} from "@marketplace-central/sdk-runtime";
import { DetailPanel, UnknownValue } from "@marketplace-central/ui";

export interface VinculoDrawerProps {
  candidateId: string | null;
  candidates: ProductLinkCandidateItem[];
  onClose: () => void;
  onApprove: (candidate: ProductLinkCandidateItem) => void;
  onReject: (candidate: ProductLinkCandidateItem) => void;
  pending?: boolean;
}

const bandClasses: Record<ProductLinkConfidenceBand, string> = {
  ALTA: "bg-emerald-100 text-emerald-800",
  MEDIA: "bg-amber-100 text-amber-800",
  BAIXA: "bg-slate-100 text-slate-700",
};

const directionClasses = {
  FOR: "bg-emerald-100 text-emerald-800",
  AGAINST: "bg-red-100 text-red-800",
  UNAVAILABLE: "bg-slate-100 text-slate-600",
} as const;

function pill(label: string, className: string) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  );
}

function confidencePercent(confidence: number): string {
  if (!Number.isFinite(confidence)) return "";
  return `${Math.round(confidence)}%`;
}

function Fact({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-slate-100 pb-2 last:border-b-0 last:pb-0">
      <dt className="text-xs text-slate-500">{label}</dt>
      <dd className="text-right text-xs font-medium text-slate-900">{children}</dd>
    </div>
  );
}

function CandidateCompareCard({
  candidate,
  rank,
  isSelected,
  onApprove,
  pending,
}: {
  candidate: ProductLinkCandidateItem;
  rank: number;
  isSelected: boolean;
  onApprove: (candidate: ProductLinkCandidateItem) => void;
  pending?: boolean;
}) {
  const noCandidate = candidate.match_status === "NO_CANDIDATE";
  return (
    <section
      data-testid="drawer-candidate"
      className={`rounded-lg border p-3 ${isSelected ? "border-blue-400 bg-blue-50/40" : "border-slate-200"}`}
    >
      <div className="mb-2 flex items-center justify-between">
        <span className="text-xs font-semibold uppercase tracking-wide text-slate-500">
          Candidato #{rank}
          {isSelected ? " (selecionado)" : ""}
        </span>
        {noCandidate
          ? pill("Sem candidato", "bg-slate-100 text-slate-600")
          : (
            <span className="flex items-center gap-1">
              {pill(candidate.confidence_band, bandClasses[candidate.confidence_band])}
              <span className="text-xs font-medium tabular-nums text-slate-700">
                {confidencePercent(candidate.confidence)}
              </span>
            </span>
          )}
      </div>

      <dl className="space-y-2">
        <Fact label="CODPROD">
          {candidate.internal_product_id === undefined ? <UnknownValue hint="sem CODPROD" /> : candidate.internal_product_id}
        </Fact>
        <Fact label="Descrição">
          {candidate.internal_product_name ? candidate.internal_product_name : <UnknownValue hint="sem descrição no ERP" />}
        </Fact>
        <Fact label="EAN/refforn">
          {candidate.internal_reference_code
            ? candidate.internal_reference_code
            : candidate.match_input === "ean" && candidate.match_value
              ? candidate.match_value
              : <UnknownValue />}
        </Fact>
        <Fact label="Entrada de match">{candidate.match_input}</Fact>
        <Fact label="Preço">
          {/* Sem campo de preço no candidato → estado honesto, nunca 0 falso (ADR-17). */}
          <UnknownValue hint="preço não disponível no candidato" />
        </Fact>
      </dl>

      {candidate.reasons.length > 0 ? (
        <ul className="mt-2 flex flex-wrap gap-1">
          {candidate.reasons.map((reason, index) => (
            <li key={`${reason.anchor}-${index}`}>
              {/* IC-01: motivo (anchor) sempre visível; detail (com %) anexado — % nunca sozinho. */}
              {pill(
                reason.detail ? `${reason.anchor}: ${reason.detail}` : reason.anchor,
                directionClasses[reason.direction],
              )}
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-2 text-xs text-slate-400">Sem sinais de correspondência.</p>
      )}

      <button
        type="button"
        disabled={noCandidate || pending}
        className="mt-3 w-full rounded-lg border border-blue-600 bg-blue-600 px-2.5 py-1.5 text-xs font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
        onClick={() => onApprove(candidate)}
      >
        Aprovar este candidato
      </button>
    </section>
  );
}

export function VinculoDrawer({
  candidateId,
  candidates,
  onClose,
  onApprove,
  onReject,
  pending,
}: VinculoDrawerProps) {
  if (candidateId === null) return null;

  const selected = candidates.find((item) => item.candidate_id === candidateId);

  // Deep-link to a candidate that is no longer in the queue (stale ?candidate=…):
  // show an honest "not found" state instead of a blank drawer.
  if (!selected) {
    return (
      <DetailPanel open onClose={onClose} closeLabel="Fechar painel" title="Candidato não encontrado">
        <p className="text-sm text-slate-600">
          Este candidato não está mais na fila. Ele pode já ter sido resolvido ou removido.
        </p>
      </DetailPanel>
    );
  }

  // All ranked candidates that target the same listing (provider_item_id), best first.
  const ranked = candidates
    .filter(
      (item) =>
        item.provider_item_id === selected.provider_item_id &&
        item.provider_variation_id === selected.provider_variation_id,
    )
    .sort((a, b) => b.confidence - a.confidence);

  return (
    <DetailPanel
      open
      onClose={onClose}
      closeLabel="Fechar painel"
      title={selected.provider_item_id}
      subtitle={selected.provider_code}
      width={360}
      footer={
        <button
          type="button"
          disabled={pending}
          className="w-full rounded-lg border border-red-200 px-2.5 py-2 text-sm font-medium text-red-700 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50"
          onClick={() => onReject(selected)}
        >
          Rejeitar anúncio
        </button>
      }
    >
      <p className="text-xs text-slate-500">
        Comparação campo a campo do anúncio com {ranked.length} candidato(s) ranqueado(s).
      </p>
      <div className="space-y-3">
        {ranked.map((candidate, index) => (
          <CandidateCompareCard
            key={candidate.candidate_id}
            candidate={candidate}
            rank={index + 1}
            isSelected={candidate.candidate_id === selected.candidate_id}
            onApprove={onApprove}
            pending={pending}
          />
        ))}
      </div>
    </DetailPanel>
  );
}
