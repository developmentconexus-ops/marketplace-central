import { useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import type {
  ListingGroup,
  ListingLinkState,
  ListingMarketSignal,
  ListingReadModel,
  ListingSyncState,
} from "@marketplace-central/sdk-runtime";
import { ConflictTag, UnknownValue } from "@marketplace-central/ui";
import { FreshnessIndicator } from "@marketplace-central/web-query";

export interface AnunciosTableProps {
  items?: ListingReadModel[];
  groups?: ListingGroup[];
  asOf?: string;
  selectedIds: Set<string>;
  onToggle: (id: string) => void;
  onTogglePage: (ids: string[]) => void;
  onOpen?: (listingId: string) => void;
}

// Ratified 9-column set (design/handoff/Anuncios.dc.html): checkbox · MLB ·
// TÍTULO · PRODUTO · PREÇO · EST. · SYNC · QUAL. · PENDÊNCIA.
const TABLE_COLUMN_COUNT = 9;

const syncLabels = {
  synced: "sincronizado",
  error: "com erro",
  stale: "desatualizado",
  queued: "na fila",
  syncing: "sincronizando",
  paused_sync: "pausado",
} satisfies Record<ListingSyncState, string>;

const linkLabels = {
  resolved: "vinculado",
  unresolved: "sem vínculo",
  rejected: "rejeitado",
} satisfies Record<Exclude<ListingLinkState, "conflict">, string>;

// SYNC status pill color map, per the ratified prototype's pillFor():
// Sincronizado -> accent-soft/accent-ink, Erro -> warn-soft/warn,
// Desatualizado -> amber-soft/amber, Na fila -> info-soft/info,
// Pausado/Sincronizando -> surface-2/faint.
const syncPillClassName = {
  synced: "bg-accent-soft text-accent-ink",
  error: "bg-warn-soft text-warn",
  stale: "bg-amber-soft text-amber",
  queued: "bg-info-soft text-info",
  syncing: "bg-surface-2 text-faint",
  paused_sync: "bg-surface-2 text-faint",
} satisfies Record<ListingSyncState, string>;

function renderSyncPill(state: ListingSyncState) {
  return (
    <span className={`inline-flex items-center gap-1 whitespace-nowrap rounded-pill px-2 py-0.5 text-xs font-medium ${syncPillClassName[state]}`}>
      <span className="h-1 w-1 rounded-pill bg-current" aria-hidden="true" />
      {syncLabels[state]}
    </span>
  );
}

function renderProductCell(item: ListingReadModel) {
  if (item.link.state === "conflict") return <ConflictTag />;
  if (item.link.product_id) {
    return (
      <Link
        to={`/catalogo/produtos/${item.link.product_id}`}
        className="font-mono text-xs text-muted hover:text-ink hover:underline"
      >
        {item.link.product_id}
      </Link>
    );
  }
  return <span className="font-mono text-xs text-faint">{linkLabels[item.link.state as Exclude<ListingLinkState, "conflict">]}</span>;
}

function formatMarketDeltaPct(signal: ListingMarketSignal): string | null {
  if (signal.delta_pct === null) return null;
  return `${signal.delta_pct.startsWith("-") ? "" : "+"}${signal.delta_pct}%`;
}

// PREÇO cell: valor R$ (mono) + a compact %-chip vs mercado folded in from
// the VS MERCADO signal (RULING D-56). Honest per signal_status (ADR-17) —
// never a fabricated chip for an evidence-less/unlinked listing. A missing
// signal_status (optional field) is treated as the honest absent state.
// NUNCA dupla sublinha na célula: exactly one of {chip, freshness link, plain value}.
function renderPriceCell(item: ListingReadModel) {
  const priceValue = item.price === null ? <UnknownValue /> : `R$ ${item.price.amount}`;
  const status = item.signal_status ?? "NO_PRICE_EVIDENCE";

  if (status === "SEM_VINCULO") {
    return (
      <div className="flex flex-col gap-0.5">
        <span className="font-mono tabular-nums">{priceValue}</span>
        <Link to="/vinculos" className="text-xs text-muted hover:text-ink hover:underline">
          sem vínculo
        </Link>
      </div>
    );
  }

  const delta = item.market_signal ? formatMarketDeltaPct(item.market_signal) : null;

  if ((status === "OK" || status === "STALE") && item.market_signal && delta !== null) {
    // Sign-colored per the delta: acima do mercado (positive) -> warn/vermelho,
    // abaixo ou igual ao mercado (<= 0) -> accent/verde; STALE always âmbar.
    const isAboveMarket = Number(item.market_signal.delta_pct) > 0;
    const chipClassName = status === "STALE" ? "bg-amber-soft text-amber" : isAboveMarket ? "bg-warn-soft text-warn" : "bg-accent-soft text-accent-ink";
    return (
      <div className="flex flex-col gap-0.5">
        <span className="flex items-center gap-1.5">
          <span className="font-mono tabular-nums">{priceValue}</span>
          <span className={`inline-flex whitespace-nowrap rounded-pill px-1.5 py-0.5 font-mono text-[11px] font-medium ${chipClassName}`}>
            {delta}
          </span>
        </span>
        {status === "STALE" ? (
          <span className="text-amber">
            <FreshnessIndicator asOf={item.market_signal.evidence?.fetched_at ?? null} />
          </span>
        ) : null}
      </div>
    );
  }

  // NO_PRICE_EVIDENCE (or OK/STALE without a usable delta): honest absent —
  // show the price value only, never 0 and never a fabricated %.
  return <span className="font-mono tabular-nums">{priceValue}</span>;
}

function formatQuality(qualityScore: number | null) {
  // quality_score is an integer 0–100 at the source (0036_listings.sql CHECK
  // BETWEEN 0 AND 100; Go QualityScore *int) — render it directly, never ×100.
  if (qualityScore === null) return <UnknownValue />;
  return `${qualityScore}%`;
}

export function AnunciosTable({ items, groups, asOf, selectedIds, onToggle, onTogglePage, onOpen }: AnunciosTableProps) {
  const headerCheckboxRef = useRef<HTMLInputElement>(null);
  // Agrupar-por-produto renders per-CODPROD groups (group-header row + child
  // listing rows); the flat W1 view (groups absent) is unchanged. Both modes
  // select over the same flattened listing set.
  const flatItems = groups ? groups.flatMap((group) => group.listings) : items ?? [];
  const pageIds = flatItems.map((item) => item.listing_id);
  const selectedOnPage = pageIds.filter((id) => selectedIds.has(id)).length;
  const allSelected = flatItems.length > 0 && selectedOnPage === flatItems.length;
  const partiallySelected = selectedOnPage > 0 && !allSelected;

  useEffect(() => {
    if (headerCheckboxRef.current) headerCheckboxRef.current.indeterminate = partiallySelected;
  }, [partiallySelected]);

  function renderRow(item: ListingReadModel) {
    const isSelected = selectedIds.has(item.listing_id);
    return (
      <tr
        key={item.listing_id}
        className={`align-middle text-ink ${isSelected ? "bg-accent-soft" : "hover:bg-surface-2"}`}
      >
        <td className="px-3 py-3">
          <input
            type="checkbox"
            checked={isSelected}
            onChange={() => onToggle(item.listing_id)}
            aria-label={`Selecionar anúncio ${item.title}`}
          />
        </td>
        <td className="px-3 py-3 font-mono text-xs text-faint">{item.provider_listing_id}</td>
        <td className="px-3 py-3">
          <button
            type="button"
            className="font-medium text-left text-accent-ink hover:underline"
            onClick={() => onOpen?.(item.listing_id)}
          >
            {item.title}
          </button>
        </td>
        <td className="px-3 py-3">{renderProductCell(item)}</td>
        <td className="px-3 py-3" data-testid="preco-cell">{renderPriceCell(item)}</td>
        <td className="px-3 py-3 font-mono tabular-nums">
          {item.published_quantity === null ? <UnknownValue /> : item.published_quantity}
        </td>
        <td className="px-3 py-3">{renderSyncPill(item.sync_state)}</td>
        <td className="px-3 py-3 font-mono tabular-nums text-xs">{formatQuality(item.quality_score)}</td>
        <td className="px-3 py-3 text-xs text-faint">
          {item.pending_issue ? (
            <span title={item.pending_issue.message_pt} className="truncate">
              {item.pending_issue.message_pt}
            </span>
          ) : null}
        </td>
      </tr>
    );
  }

  function renderGroupHeader(group: ListingGroup) {
    const title = group.product_title ?? "Sem produto";
    return (
      <tr key={`group:${group.product_id ?? "sem-produto"}`} className="bg-surface-2 text-ink">
        <td className="px-3 py-2" colSpan={TABLE_COLUMN_COUNT}>
          <span className="text-sm font-semibold">{title}</span>{" "}
          <span className="text-xs font-normal text-faint">({group.listing_count})</span>
        </td>
      </tr>
    );
  }

  return (
    <div className="mt-3 overflow-x-auto">
      <table className="w-full min-w-[820px] border-collapse text-left text-sm">
        <caption className="sr-only">Anúncios{asOf ? `, dados de ${asOf}` : ""}</caption>
        <thead className="border-b border-border text-xs font-medium text-faint">
          <tr>
            <th className="px-3 py-3" scope="col">
              <input
                ref={headerCheckboxRef}
                type="checkbox"
                checked={allSelected}
                disabled={flatItems.length === 0}
                onChange={() => onTogglePage(pageIds)}
                aria-label="Selecionar todos os anúncios desta página"
              />
            </th>
            <th className="px-3 py-3" scope="col">MLB</th>
            <th className="px-3 py-3" scope="col">TÍTULO</th>
            <th className="px-3 py-3" scope="col">PRODUTO</th>
            <th className="px-3 py-3" scope="col">PREÇO</th>
            <th className="px-3 py-3" scope="col">EST.</th>
            <th className="px-3 py-3" scope="col">SYNC</th>
            <th className="px-3 py-3" scope="col">QUAL.</th>
            <th className="px-3 py-3" scope="col">PENDÊNCIA</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border-2">
          {groups
            ? groups.flatMap((group) => [renderGroupHeader(group), ...group.listings.map((item) => renderRow(item))])
            : flatItems.map((item) => renderRow(item))}
        </tbody>
      </table>
    </div>
  );
}
