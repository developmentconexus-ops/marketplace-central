import type {
  ProductLinkCandidateItem,
  ProductLinkConfidenceBand,
  ProductLinkReason,
  ProductLinkReasonDirection,
} from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";

export interface QueueRowProps {
  candidate: ProductLinkCandidateItem;
  onOpen: (candidateId: string) => void;
  onApprove: (candidate: ProductLinkCandidateItem) => void;
  pending?: boolean;
  selected?: boolean;
  onToggleSelect?: (candidateId: string) => void;
}

const bandLabels: Record<ProductLinkConfidenceBand, string> = {
  ALTA: "ALTA",
  MEDIA: "MEDIA",
  BAIXA: "BAIXA",
};

const bandClasses: Record<ProductLinkConfidenceBand, string> = {
  ALTA: "bg-emerald-100 text-emerald-800",
  MEDIA: "bg-amber-100 text-amber-800",
  BAIXA: "bg-slate-100 text-slate-700",
};

const directionClasses: Record<ProductLinkReasonDirection, string> = {
  FOR: "bg-emerald-100 text-emerald-800",
  AGAINST: "bg-red-100 text-red-800",
  UNAVAILABLE: "bg-slate-100 text-slate-600",
};

function pill(label: string, className: string) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  );
}

function confidencePercent(confidence: number): string {
  if (!Number.isFinite(confidence)) return "";
  return `${Math.round(confidence * 100)}%`;
}

/**
 * IC-01 presentation rule: a reason chip ALWAYS shows its motivo (anchor) text.
 * When a `detail` (which may carry a %) is present it is rendered joined to the
 * motivo — so a bare % is never shown on its own; the motivo is always visible.
 */
function reasonChipLabel(reason: ProductLinkReason): string {
  return reason.detail ? `${reason.anchor}: ${reason.detail}` : reason.anchor;
}

function AnchorChips({ reasons }: { reasons: ProductLinkReason[] }) {
  if (reasons.length === 0) {
    return <UnknownValue hint="sem sinais de correspondência" />;
  }
  return (
    <ul className="flex flex-wrap gap-1">
      {reasons.map((reason, index) => (
        <li key={`${reason.anchor}-${index}`}>
          {pill(reasonChipLabel(reason), directionClasses[reason.direction])}
        </li>
      ))}
    </ul>
  );
}

export function QueueRow({ candidate, onOpen, onApprove, pending, selected, onToggleSelect }: QueueRowProps) {
  const noCandidate = candidate.match_status === "NO_CANDIDATE";
  const codprod = candidate.internal_product_id;
  const descricao = candidate.internal_product_name;
  const eanRefforn =
    candidate.internal_reference_code ??
    (candidate.match_input === "ean" ? candidate.match_value : undefined);

  return (
    <tr className="align-top text-slate-700" data-testid="queue-row" data-match-status={candidate.match_status}>
      {/* Seleção em lote */}
      <td className="px-3 py-3">
        <input
          type="checkbox"
          aria-label={`Selecionar ${candidate.provider_item_id}`}
          checked={selected ?? false}
          disabled={noCandidate || !onToggleSelect}
          onChange={() => onToggleSelect?.(candidate.candidate_id)}
        />
      </td>

      {/* Produto (interno) */}
      <td className="px-3 py-3">
        <div className="font-medium text-slate-900">
          {codprod === undefined ? <UnknownValue hint="sem CODPROD" /> : codprod}
        </div>
        <div className="mt-0.5 text-xs text-slate-600">
          {descricao ? descricao : <UnknownValue hint="sem descrição no ERP" />}
        </div>
        <div className="mt-0.5 text-xs text-slate-500">
          EAN/refforn: {eanRefforn ? eanRefforn : <UnknownValue />}
        </div>
      </td>

      {/* Melhor candidato (anúncio) */}
      <td className="px-3 py-3">
        {noCandidate ? (
          <div className="flex flex-col gap-1">
            {pill("Sem candidato", "bg-slate-100 text-slate-600")}
            <span className="text-xs text-slate-500">
              Nenhuma correspondência encontrada para {candidate.provider_item_id}.
            </span>
          </div>
        ) : (
          <>
            <div className="font-medium text-slate-900">{candidate.provider_item_id}</div>
            <div className="mt-0.5 text-xs text-slate-500">{candidate.provider_code}</div>
            <div className="mt-1 text-xs text-slate-500">
              {/* Preço genuinamente ausente no candidato → estado honesto, nunca 0 falso (ADR-17). */}
              Preço: <UnknownValue hint="preço não disponível no candidato" />
            </div>
          </>
        )}
      </td>

      {/* Sinais (anchor chips) */}
      <td className="px-3 py-3">
        <AnchorChips reasons={candidate.reasons} />
      </td>

      {/* Confiança / banda % */}
      <td className="px-3 py-3">
        {noCandidate ? (
          <UnknownValue hint="sem confiança sem candidato" />
        ) : (
          <div className="flex items-center gap-2">
            {pill(bandLabels[candidate.confidence_band], bandClasses[candidate.confidence_band])}
            <span className="text-xs font-medium tabular-nums text-slate-700">
              {confidencePercent(candidate.confidence)}
            </span>
          </div>
        )}
      </td>

      {/* Ações */}
      <td className="px-3 py-3">
        <div className="flex justify-end gap-2">
          <button
            type="button"
            className="rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-medium text-slate-700 hover:border-slate-300 hover:bg-slate-50"
            onClick={() => onOpen(candidate.candidate_id)}
          >
            Abrir
          </button>
          <button
            type="button"
            disabled={noCandidate || pending}
            className="rounded-lg border border-blue-600 bg-blue-600 px-2.5 py-1.5 text-xs font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
            onClick={() => onApprove(candidate)}
          >
            Aprovar
          </button>
        </div>
      </td>
    </tr>
  );
}
