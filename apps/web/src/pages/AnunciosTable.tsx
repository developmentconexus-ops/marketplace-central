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

const TABLE_COLUMN_COUNT = 11;

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

function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>{label}</span>;
}

function renderLinkState(state: ListingLinkState) {
  if (state === "conflict") return <ConflictTag />;
  return stateTag(linkLabels[state]);
}

function renderMargin(item: ListingReadModel) {
  if (item.below_margin_worst_case === null) {
    return item.cost === null ? <UnknownValue hint="sem custo no ERP → não simulado" /> : <UnknownValue />;
  }
  return item.below_margin_worst_case
    ? stateTag("abaixo da margem", "bg-red-100 text-red-800")
    : stateTag("ok", "bg-emerald-100 text-emerald-800");
}

function formatMarketPosition(signal: ListingMarketSignal): string {
  return signal.position ? `${signal.position.rank}/${signal.position.total}` : "—";
}

function formatMarketDelta(signal: ListingMarketSignal): string {
  return signal.delta_pct === null
    ? "—"
    : `${signal.delta_pct.startsWith("-") ? "" : "+"}${signal.delta_pct}%`;
}

function renderMarketSignalValue(signal: ListingMarketSignal) {
  return (
    <span className="font-mono tabular-nums">
      <span>{formatMarketPosition(signal)}</span> · <span>{formatMarketDelta(signal)}</span>
    </span>
  );
}

// VS MERCADO column: distinct honest rendering per signal_status (ADR-17) —
// never a fabricated number for an evidence-less/unlinked listing. A missing
// signal_status (optional field) is treated as the honest absent state.
function renderMarketSignal(item: ListingReadModel) {
  const status = item.signal_status ?? "NO_PRICE_EVIDENCE";

  if (status === "SEM_VINCULO") {
    return (
      <Link
        to="/vinculos"
        className="inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium bg-surface-2 text-muted hover:text-ink"
      >
        sem vínculo
      </Link>
    );
  }

  if (status === "STALE" && item.market_signal) {
    return (
      <div className="flex flex-col gap-0.5">
        {renderMarketSignalValue(item.market_signal)}
        <span className="text-amber">
          <FreshnessIndicator asOf={item.market_signal.evidence?.fetched_at ?? null} />
        </span>
      </div>
    );
  }

  if (status === "OK" && item.market_signal) {
    return renderMarketSignalValue(item.market_signal);
  }

  return <UnknownValue hint="sem evidência de preço de mercado" />;
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
    return (
      <tr key={item.listing_id} className="align-middle text-slate-700">
        <td className="px-3 py-3">
          <input
            type="checkbox"
            checked={selectedIds.has(item.listing_id)}
            onChange={() => onToggle(item.listing_id)}
            aria-label={`Selecionar anúncio ${item.title}`}
          />
        </td>
        <td className="px-3 py-3">
          <button
            type="button"
            className="font-medium text-left text-blue-700 hover:text-blue-800 hover:underline"
            onClick={() => onOpen?.(item.listing_id)}
          >
            {item.title}
          </button>
          <div className="mt-0.5 text-xs text-slate-500">{item.provider_listing_id}</div>
        </td>
        <td className="px-3 py-3">{item.listing_type ? stateTag(item.listing_type.label) : <UnknownValue />}</td>
        <td className="px-3 py-3">{renderLinkState(item.link.state)}</td>
        <td className="px-3 py-3 tabular-nums">{item.price === null ? <UnknownValue /> : `R$ ${item.price.amount}`}</td>
        <td className="px-3 py-3 tabular-nums">{item.published_quantity === null ? <UnknownValue /> : item.published_quantity}</td>
        <td className="px-3 py-3 tabular-nums">{item.sales_30d === null ? <UnknownValue /> : item.sales_30d}</td>
        <td className="px-3 py-3">{renderMargin(item)}</td>
        <td className="px-3 py-3">{renderMarketSignal(item)}</td>
        <td className="px-3 py-3">
          {stateTag(syncLabels[item.sync_state], item.sync_state === "error" ? "bg-red-100 text-red-800" : undefined)}
        </td>
        <td className="px-3 py-3">
          {item.pending_issue ? (
            <span aria-label={item.pending_issue.message_pt} className="inline-flex h-6 min-w-6 items-center justify-center rounded-full bg-amber-100 px-1.5 text-amber-800" title={item.pending_issue.message_pt}>⚑</span>
          ) : null}
        </td>
      </tr>
    );
  }

  function renderGroupHeader(group: ListingGroup) {
    const title = group.product_title ?? "Sem produto";
    return (
      <tr key={`group:${group.product_id ?? "sem-produto"}`} className="bg-slate-50 text-slate-900">
        <td className="px-3 py-2" colSpan={TABLE_COLUMN_COUNT}>
          <span className="text-sm font-semibold">{title}</span>{" "}
          <span className="text-xs font-normal text-slate-500">({group.listing_count})</span>
        </td>
      </tr>
    );
  }

  return (
    <div className="mt-3 overflow-x-auto">
      <table className="w-full min-w-[980px] border-collapse text-left text-sm">
        <caption className="sr-only">Anúncios{asOf ? `, dados de ${asOf}` : ""}</caption>
        <thead className="border-b border-slate-200 text-xs font-medium text-slate-500">
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
            <th className="px-3 py-3" scope="col">Anúncio</th>
            <th className="px-3 py-3" scope="col">Modalidade</th>
            <th className="px-3 py-3" scope="col">Vínculo</th>
            <th className="px-3 py-3" scope="col">Preço</th>
            <th className="px-3 py-3" scope="col">Estoque</th>
            <th className="px-3 py-3" scope="col">Vendas 30d</th>
            <th className="px-3 py-3" scope="col">Margem</th>
            <th className="px-3 py-3" scope="col">Vs. mercado</th>
            <th className="px-3 py-3" scope="col">Sync</th>
            <th className="px-3 py-3" scope="col"><span className="sr-only">Pendência</span>⚑</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {groups
            ? groups.flatMap((group) => [renderGroupHeader(group), ...group.listings.map((item) => renderRow(item))])
            : flatItems.map((item) => renderRow(item))}
        </tbody>
      </table>
    </div>
  );
}
