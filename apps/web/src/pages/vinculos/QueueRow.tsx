import type {
  ProductLinkCandidateItem,
  ProductLinkConfidenceBand,
  ProductLinkReason,
  ProductLinkReasonDirection,
} from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";
import { useState } from "react";

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

// The wire anchors are machine names (`seller_sku`, `refforn`); the table shows
// the operator-facing short form. Unknown anchors fall through verbatim — never
// hidden, never renamed into something the wire didn't say.
const anchorShortLabels: Record<string, string> = {
  seller_sku: "SKU",
  ean: "EAN",
  title: "Título",
  marca: "Marca",
  refforn: "Ref. forn.",
};

// Direction is carried by colour AND by a glyph, so the FOR/AGAINST/sem-dado
// distinction survives without relying on colour alone.
const directionGlyphs: Record<ProductLinkReasonDirection, string> = {
  FOR: "✓",
  AGAINST: "✕",
  UNAVAILABLE: "–",
};

function anchorShortLabel(anchor: string): string {
  return anchorShortLabels[anchor] ?? anchor;
}

/**
 * Compact chip label for the table. FOR/UNAVAILABLE motivos restate the anchor
 * ("seller_sku resolve exato para codprod") so the detail adds nothing at a
 * glance and lives in the tooltip. An AGAINST detail is the *reason the vínculo
 * is blocked* — it stays inline (CSS-truncated to one line, full text in the
 * tooltip and in the expanded view).
 */
function compactChipLabel(reason: ProductLinkReason): string {
  const head = `${directionGlyphs[reason.direction]} ${anchorShortLabel(reason.anchor)}`;
  if (reason.direction === "AGAINST" && reason.detail) {
    return `${head}: ${reason.detail}`;
  }
  return head;
}

function compactChip(reason: ProductLinkReason) {
  return (
    <span
      title={reasonChipLabel(reason)}
      className={`inline-flex max-w-[12rem] items-center truncate whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ${directionClasses[reason.direction]}`}
    >
      {compactChipLabel(reason)}
    </span>
  );
}

const COMPACT_CHIP_LIMIT = 2;

/**
 * Motivo cell. Collapsed by default to keep the row one line tall (a 4-chip
 * stack of full sentences made every row ~5 lines and pushed the Ação column
 * off-screen). AGAINST motivos are ranked first — a blocking contradiction is
 * what the operator must read — and nothing is silently dropped: the remainder
 * is behind a "+N" toggle that renders every motivo in its full wire form.
 */
function AnchorChips({ reasons }: { reasons: ProductLinkReason[] }) {
  const [expanded, setExpanded] = useState(false);

  if (reasons.length === 0) {
    return <UnknownValue hint="sem sinais de correspondência" />;
  }

  if (expanded) {
    return (
      <div className="flex flex-col items-start gap-1">
        <ul className="flex flex-col items-start gap-1">
          {reasons.map((reason, index) => (
            <li key={`${reason.anchor}-${index}`}>
              {/* Expandido: o texto integral do motivo pode quebrar linha — o
                  chip nowrap do modo compacto seria cortado pelo limite do td. */}
              <span
                className={`inline-block rounded-2xl px-2 py-0.5 text-xs font-medium ${directionClasses[reason.direction]}`}
              >
                {reasonChipLabel(reason)}
              </span>
            </li>
          ))}
        </ul>
        <button
          type="button"
          className="rounded-control px-1 text-xs font-medium text-muted underline hover:text-ink"
          onClick={() => setExpanded(false)}
        >
          menos
        </button>
      </div>
    );
  }

  // Rank AGAINST → FOR → UNAVAILABLE and slice. Ranking (never filtering) is
  // what keeps at least one motivo on screen even for a row whose only signals
  // are UNAVAILABLE ones (ADR-17 — motivo sempre visível).
  const byDirection = (direction: ProductLinkReasonDirection) =>
    reasons.filter((reason) => reason.direction === direction);
  const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")].slice(
    0,
    COMPACT_CHIP_LIMIT,
  );
  const hidden = reasons.length - shown.length;

  return (
    <div className="flex min-w-0 items-center gap-1">
      <ul className="flex min-w-0 items-center gap-1">
        {shown.map((reason, index) => (
          <li key={`${reason.anchor}-${index}`} className="min-w-0">
            {compactChip(reason)}
          </li>
        ))}
      </ul>
      {hidden > 0 ? (
        <button
          type="button"
          aria-label={`Mostrar todos os ${reasons.length} motivos`}
          title={reasons.map(reasonChipLabel).join(" · ")}
          className="shrink-0 rounded-full border border-border px-2 py-0.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
          onClick={() => setExpanded(true)}
        >
          +{hidden}
        </button>
      ) : null}
    </div>
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

      {/* AÇÃO — sticky à direita: as ações nunca saem da tela quando a tabela
          rola horizontalmente. */}
      <td className="sticky right-0 border-l border-border bg-surface px-3 py-3">
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
