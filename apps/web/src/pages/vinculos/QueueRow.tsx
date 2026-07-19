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
  onReject: (candidate: ProductLinkCandidateItem) => void;
  pending?: boolean;
  selected?: boolean;
  onToggleSelect?: (candidateId: string) => void;
}

const bandLabels: Record<ProductLinkConfidenceBand, string> = {
  ALTA: "ALTA",
  MEDIA: "MEDIA",
  BAIXA: "BAIXA",
};

// Confidence bands keep their SEMANTIC mapping (≥85 verde / 50-84 âmbar / <50
// vermelho) but through the paper+green design tokens, never literal Tailwind.
const bandClasses: Record<ProductLinkConfidenceBand, string> = {
  ALTA: "bg-accent-soft text-accent-ink",
  MEDIA: "bg-amber-soft text-amber",
  BAIXA: "bg-warn-soft text-warn",
};

const directionClasses: Record<ProductLinkReasonDirection, string> = {
  FOR: "bg-accent-soft text-accent-ink",
  AGAINST: "bg-warn-soft text-warn",
  UNAVAILABLE: "bg-surface-2 text-faint",
};

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

export function QueueRow({ candidate, onOpen, onApprove, onReject, pending, selected, onToggleSelect }: QueueRowProps) {
  const noCandidate = candidate.match_status === "NO_CANDIDATE";
  // GTIN "✓ igual" is only honest when the listing matched the product on EAN
  // (match_input === "ean"); otherwise the GTIN relationship is UNKNOWN → "—",
  // never fabricated (ADR-17).
  const gtinEqual = candidate.match_input === "ean" && Boolean(candidate.match_value);

  return (
    <tr className="align-top text-ink" data-testid="queue-row" data-match-status={candidate.match_status}>
      {/* Seleção em lote */}
      <td className="px-3 py-3">
        <input
          type="checkbox"
          className="accent-[var(--accent)]"
          aria-label={`Selecionar ${candidate.provider_item_id}`}
          checked={selected ?? false}
          disabled={noCandidate || !onToggleSelect}
          onChange={() => onToggleSelect?.(candidate.candidate_id)}
        />
      </td>

      {/* ANÚNCIO ML (listing) */}
      <td className="px-3 py-3">
        <div className="font-mono text-sm font-medium text-ink">{candidate.provider_item_id}</div>
      </td>

      {/* SKU ML (código do anúncio no provider) */}
      <td className="px-3 py-3">
        <span className="font-mono text-xs text-muted">
          {candidate.provider_code ? candidate.provider_code : <UnknownValue />}
        </span>
      </td>

      {/* PRODUTO SUGERIDO (produto interno) */}
      <td className="px-3 py-3">
        {noCandidate ? (
          <div className="flex flex-col gap-1">
            {pill("Sem candidato", "bg-surface-2 text-faint")}
            <span className="text-xs text-faint">
              Nenhuma correspondência para {candidate.provider_item_id}.
            </span>
          </div>
        ) : (
          <div className="font-medium text-ink">
            {candidate.internal_product_name ? (
              candidate.internal_product_name
            ) : (
              <UnknownValue hint="sem descrição no ERP" />
            )}
          </div>
        )}
      </td>

      {/* SKU HUB (CODPROD interno) */}
      <td className="px-3 py-3">
        <span className="font-mono text-sm text-ink">
          {candidate.internal_product_id === undefined ? (
            <UnknownValue hint="sem CODPROD" />
          ) : (
            candidate.internal_product_id
          )}
        </span>
      </td>

      {/* GTIN ("✓ igual" | "—") */}
      <td className="px-3 py-3">
        {gtinEqual ? (
          <span className="whitespace-nowrap text-xs font-medium text-accent-ink">✓ igual</span>
        ) : (
          <UnknownValue />
        )}
      </td>

      {/* CONFIANÇA (banda + %) */}
      <td className="px-3 py-3">
        {noCandidate ? (
          <UnknownValue hint="sem confiança sem candidato" />
        ) : (
          <div className="flex items-center gap-2">
            {pill(bandLabels[candidate.confidence_band], bandClasses[candidate.confidence_band])}
            <span className="font-mono text-xs font-medium tabular-nums text-muted">
              {confidencePercent(candidate.confidence)}
            </span>
          </div>
        )}
      </td>

      {/* MOTIVO (sinais / anchor chips) */}
      <td className="px-3 py-3">
        <AnchorChips reasons={candidate.reasons} />
      </td>

      {/* AÇÃO */}
      <td className="px-3 py-3">
        <div className="flex justify-end gap-2">
          {noCandidate ? (
            <>
              {/* Criação de produto a partir do anúncio ainda não tem seam de
                  escrita — afordância honesta (desabilitada), nunca sucesso falso. */}
              <button
                type="button"
                disabled
                title="Criação de produto a partir do anúncio ainda não disponível"
                className="cursor-not-allowed rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-faint opacity-60"
              >
                Criar produto
              </button>
              <button
                type="button"
                disabled={pending}
                className="rounded-control border border-warn/40 px-2.5 py-1.5 text-xs font-medium text-warn hover:bg-warn-soft disabled:cursor-not-allowed disabled:opacity-50"
                onClick={() => onReject(candidate)}
              >
                Ignorar
              </button>
            </>
          ) : (
            <>
              <button
                type="button"
                className="rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
                onClick={() => onOpen(candidate.candidate_id)}
              >
                Outro…
              </button>
              <button
                type="button"
                disabled={pending}
                className="rounded-control bg-accent px-2.5 py-1.5 text-xs font-medium text-accent-ink hover:bg-accent/90 disabled:cursor-not-allowed disabled:opacity-50"
                onClick={() => onApprove(candidate)}
              >
                Vincular
              </button>
            </>
          )}
        </div>
      </td>
    </tr>
  );
}
