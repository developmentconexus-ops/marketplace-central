import type {
  ProductLinkCandidateItem,
  ProductLinkCandidateMatchInput,
  ProductLinkConfidenceBand,
  ProductLinkMatchStatus,
  ProductLinkReason,
  ProductLinkReasonDirection,
  ProductLinkReasonSide,
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

// INCOMPARABLE gets its OWN token pair (info), never accent (FOR — corroborates)
// and never warn (AGAINST — blocks) and never the muted surface of UNAVAILABLE.
// D-122/D-B created the state precisely to separate "the provider does not
// supply this anchor, ever" (UNAVAILABLE, nothing to do) from "the provider does
// supply it and the value is missing HERE" (INCOMPARABLE, go register it):
// opposite operator actions, so painting them alike would erase the distinction
// the state exists to show.
export const directionClasses: Record<ProductLinkReasonDirection, string> = {
  FOR: "bg-accent-soft text-accent-ink",
  AGAINST: "bg-warn-soft text-warn",
  UNAVAILABLE: "bg-surface-2 text-faint",
  INCOMPARABLE: "bg-info-soft text-info",
};

/**
 * Provider display name for a wire `provider_code`. The wire value is a slug
 * ("mercado_livre"); rendering it raw is the bug that hit CHIP-PED-FILA across
 * four surfaces. An unknown provider falls through verbatim — same honesty rule
 * as `anchorShortLabel`: never rename a value into something the wire did not say.
 */
const providerDisplayNames: Record<string, string> = {
  mercado_livre: "Mercado Livre",
};

export function providerDisplayName(providerCode: string): string {
  return providerDisplayNames[providerCode] ?? providerCode;
}

/**
 * "Identificado por" (D-122): the anchors that DECIDED, joined by ` + `.
 *
 * This is not a new rule invented on the screen. `decisionRuleForCandidate`
 * (resolution_service.go:812-835) is a PURE function of `match_status` and
 * `match_input` — both already on the wire — and it is the same function that
 * writes `rule_matched` into the decision trail. Reading the same two inputs
 * lands on the same rule, so the column says what the trail will say.
 *
 * Anything that is not an ACCEPT or a single-anchor CONFIRM files as `manual`:
 * no anchor carried it, so the column names none. That is what separates this
 * column from Motivo — Motivo lists everything that OPINED (including every
 * UNAVAILABLE and INCOMPARABLE), this lists only what DECIDED. A `title FOR`
 * candidate (state TitleMatch → REVIEW, ranking-only and never ACCEPT under
 * D-121) therefore keeps its motivo and shows NOTHING here.
 *
 * Both maps are keyed by the full union rather than switched on string
 * literals: a sixth match_status or a new match_input has to fail the compiler,
 * not fall silently through a `default` — the QueueRow:159 lesson.
 */
const confirmDecidingAnchors: Record<ProductLinkCandidateMatchInput, string[]> = {
  seller_sku: ["CODPROD"],
  ean: ["EAN"],
  title: [],
  manual: [],
  none: [],
};

const statusDecidingAnchors: Record<
  ProductLinkMatchStatus,
  (matchInput: ProductLinkCandidateMatchInput) => string[]
> = {
  // Corroborated: CODPROD and EAN both named this product. Hiding the second
  // anchor would erase from the screen the very thing that separates
  // auto-approved from sent-to-confirmation (D-121).
  ACCEPT: () => ["CODPROD", "EAN"],
  CONFIRM: (matchInput) => confirmDecidingAnchors[matchInput],
  // The anchors disagreed, collided, or never ran. Nothing decided.
  REVIEW: () => [],
  REJECT: () => [],
  NO_CANDIDATE: () => [],
};

export function decidingAnchors(candidate: ProductLinkCandidateItem): string[] {
  return statusDecidingAnchors[candidate.match_status](candidate.match_input);
}

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

// `side` is the second half of D-122/D-B: WHICH side is missing the value, and
// therefore where the operator has to go fix it. Read from the field, never
// parsed out of `detail` — parsing a Portuguese sentence breaks silently the day
// someone rewrites the sentence, which is exactly why D-B put it in the data.
const incomparableSideLabels: Record<ProductLinkReasonSide, string> = {
  provider: "falta no anúncio",
  erp: "falta no ERP",
  both: "falta nos dois lados",
};

/**
 * The `side` text for a reason, or undefined when there is none to tell.
 *
 * The direction is checked FIRST: the frozen D-122 contract emits `side` only
 * for INCOMPARABLE, so a stray value on any other direction is wire noise, not
 * a fact to render. And an INCOMPARABLE whose `side` is absent stays silent —
 * `classifyProviderIdentityAnchor` has a real path that emits INCOMPARABLE with
 * no side (the `!readable` branch, generation_service.go:711, e.g. `marca`),
 * and inventing a side there would fabricate the one thing D-B added (ADR-17).
 */
export function reasonSideLabel(reason: ProductLinkReason): string | undefined {
  if (reason.direction !== "INCOMPARABLE") return undefined;
  return reason.side ? incomparableSideLabels[reason.side] : undefined;
}

/**
 * IC-01 presentation rule: a reason chip ALWAYS shows its motivo (anchor) text.
 * When a `detail` (which may carry a %) is present it is rendered joined to the
 * motivo — so a bare % is never shown on its own; the motivo is always visible.
 */
export function reasonChipLabel(reason: ProductLinkReason): string {
  const side = reasonSideLabel(reason);
  const head = side ? `${reason.anchor} (${side})` : reason.anchor;
  return reason.detail ? `${head}: ${reason.detail}` : head;
}

// The wire anchors are machine names (`seller_sku`, `refforn`); the table shows
// the operator-facing short form. Unknown anchors fall through verbatim — never
// hidden, never renamed into something the wire didn't say.
//
// `refforn` is KEPT even though D-A (D-122) removed it from the cross-side
// anchor vocabulary, so no candidate generated today carries it. D-A also
// decided that already-persisted reasons are NOT migrated, so a decided
// candidate can still hold a `refforn` motivo — audit data that is meant to
// survive verbatim. Dropping the entry would degrade that historical row from
// "Ref. forn." to the raw machine name for no gain.
const anchorShortLabels: Record<string, string> = {
  seller_sku: "SKU",
  ean: "EAN",
  title: "Título",
  marca: "Marca",
  refforn: "Ref. forn.",
};

// Direction is carried by colour AND by a glyph, so the FOR/AGAINST/sem-dado
// distinction survives without relying on colour alone. INCOMPARABLE is "?":
// the comparison could not be made — distinct from "–" (UNAVAILABLE: there is
// nothing to compare, and never will be).
const directionGlyphs: Record<ProductLinkReasonDirection, string> = {
  FOR: "✓",
  AGAINST: "✕",
  UNAVAILABLE: "–",
  INCOMPARABLE: "?",
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
  const side = reasonSideLabel(reason);
  // The side rides inline, not only in the tooltip: "where do I go fix it" is
  // the whole point of an INCOMPARABLE, and a tooltip is not readable on a scan.
  const anchor = side ? `${anchorShortLabel(reason.anchor)} (${side})` : anchorShortLabel(reason.anchor);
  const head = `${directionGlyphs[reason.direction]} ${anchor}`;
  if (reason.direction === "AGAINST" && reason.detail) {
    return `${head}: ${reason.detail}`;
  }
  return head;
}

function compactChip(reason: ProductLinkReason) {
  return (
    <span
      data-testid="motivo-chip"
      data-direction={reason.direction}
      title={reasonChipLabel(reason)}
      className={`inline-flex max-w-[12rem] items-center truncate whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ${directionClasses[reason.direction]}`}
    >
      {compactChipLabel(reason)}
    </span>
  );
}

const COMPACT_CHIP_LIMIT = 2;

// Display priority for the collapsed cell. AGAINST first (a blocking
// contradiction is what the operator must read), then FOR (what corroborates),
// then INCOMPARABLE — actionable absence, the operator CAN go register the
// value — and last UNAVAILABLE, the only permanent one and therefore the only
// one nothing can be done about.
const directionRank: Record<ProductLinkReasonDirection, number> = {
  AGAINST: 0,
  FOR: 1,
  INCOMPARABLE: 2,
  UNAVAILABLE: 3,
};

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

  // Rank and slice. Ranking (never filtering) is what keeps at least one motivo
  // on screen even for a row whose only signals are absence ones (ADR-17 —
  // motivo sempre visível).
  //
  // This USED to be `[...byDirection("AGAINST"), ...byDirection("FOR"),
  // ...byDirection("UNAVAILABLE")]` — an enumeration by string literal. It is
  // type-correct, so the compiler said nothing when D-B added a fourth
  // direction, and the "ranking" quietly became a FILTER: a row whose motivos
  // are all INCOMPARABLE produced an empty `shown` and rendered a lone "+N"
  // button with zero chips — the exact opposite of the invariant above.
  //
  // A sort over a Record<Direction, number> is total by construction: every
  // reason keeps a place in the ordering, so `shown` is empty only when
  // `reasons` is, and a future fifth direction stops COMPILING here instead of
  // disappearing from the cell.
  const shown = reasons
    .map((reason, index) => ({ reason, index }))
    // The index tiebreak keeps the wire order inside a direction, so the cell
    // does not reshuffle between renders for reasons that rank equally.
    .sort((a, b) => directionRank[a.reason.direction] - directionRank[b.reason.direction] || a.index - b.index)
    .slice(0, COMPACT_CHIP_LIMIT)
    .map((entry) => entry.reason);
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
  const decided = decidingAnchors(candidate);

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

      {/* ANÚNCIO (id do anúncio no provider) — rótulo estrutural neutro; de QUAL
          provider ele é fica na coluna Canal, que é onde esse dado mora. */}
      <td className="px-3 py-3">
        <div className="font-mono text-sm font-medium text-ink">{candidate.provider_item_id}</div>
      </td>

      {/* CANAL (provider do anúncio). A célula mostrava `provider_code` cru sob
          um cabeçalho "SKU ML": o wire carrega aqui o SLUG do marketplace
          ("mercado_livre"), nunca um SKU — o SKU do vendedor não está no
          contrato do candidato. Rótulo agora diz o que o dado é, e o valor sai
          pelo nome de exibição em vez do slug. */}
      <td className="px-3 py-3">
        <span className="text-xs text-muted">
          {candidate.provider_code ? providerDisplayName(candidate.provider_code) : <UnknownValue />}
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

      {/* IDENTIFICADO POR — as âncoras que DECIDIRAM, unidas por " + " (D-122).
          Substitui a coluna GTIN: "✓ igual" era a leitura de UMA âncora, e o
          conjunto que decidiu é a informação que a supera. Vazio (REVIEW,
          REJECT, NO_CANDIDATE, CONFIRM por título) é "—": nada decidiu ainda,
          e inventar uma âncora aqui seria afirmar uma decisão que não houve. */}
      <td className="px-3 py-3">
        {decided.length > 0 ? (
          <span
            className="whitespace-nowrap text-xs font-medium text-accent-ink"
            data-testid="identificado-por"
          >
            {decided.join(" + ")}
          </span>
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
