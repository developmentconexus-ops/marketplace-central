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
  ALTA: "bg-accent-soft text-accent-ink",
  MEDIA: "bg-amber-soft text-amber",
  BAIXA: "bg-warn-soft text-warn",
};

const directionClasses = {
  FOR: "bg-accent-soft text-accent-ink",
  AGAINST: "bg-warn-soft text-warn",
  UNAVAILABLE: "bg-surface-2 text-faint",
} as const;

function pill(label: string, className: string) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ${className}`}>
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
    <div className="flex items-start justify-between gap-4 border-b border-border-2 pb-2 last:border-b-0 last:pb-0">
      <dt className="text-xs text-faint">{label}</dt>
      <dd className="text-right text-xs font-medium text-ink">{children}</dd>
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
      className={`rounded-card border p-3 ${isSelected ? "border-accent bg-accent-soft/50" : "border-border"}`}
    >
      <div className="mb-2 flex items-center justify-between">
        <span className="text-xs font-semibold uppercase tracking-wide text-faint">
          Candidato #{rank}
          {isSelected ? " (selecionado)" : ""}
        </span>
        {noCandidate
          ? pill("Sem candidato", "bg-surface-2 text-faint")
          : (
            <span className="flex items-center gap-1">
              {pill(candidate.confidence_band, bandClasses[candidate.confidence_band])}
              <span className="font-mono text-xs font-medium tabular-nums text-muted">
                {confidencePercent(candidate.confidence)}
              </span>
            </span>
          )}
      </div>

      <dl className="space-y-2">
        <Fact label="SKU HUB (CODPROD)">
          {candidate.internal_product_id === undefined ? (
            <UnknownValue hint="sem CODPROD" />
          ) : (
            <span className="font-mono">{candidate.internal_product_id}</span>
          )}
        </Fact>
        <Fact label="Produto sugerido">
          {candidate.internal_product_name ? candidate.internal_product_name : <UnknownValue hint="sem descrição no ERP" />}
        </Fact>
        <Fact label="GTIN / refforn">
          {candidate.internal_reference_code
            ? <span className="font-mono">{candidate.internal_reference_code}</span>
            : candidate.match_input === "ean" && candidate.match_value
              ? <span className="font-mono">{candidate.match_value}</span>
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
        <p className="mt-2 text-xs text-faint">Sem sinais de correspondência.</p>
      )}

      <button
        type="button"
        disabled={noCandidate || pending}
        className="mt-3 w-full rounded-control bg-accent px-2.5 py-1.5 text-xs font-medium text-accent-ink hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-50"
        onClick={() => onApprove(candidate)}
      >
        Vincular este candidato
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
        <p className="text-sm text-muted">
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
          className="w-full rounded-control border border-warn/40 px-2.5 py-2 text-sm font-medium text-warn hover:bg-warn-soft disabled:cursor-not-allowed disabled:opacity-50"
          onClick={() => onReject(selected)}
        >
          Ignorar anúncio
        </button>
      }
    >
      <p className="text-xs text-faint">
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
