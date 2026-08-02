import { useQuery } from "@tanstack/react-query";
import {
  hasCode,
  isApiError,
  type ListingException,
  type ListingGroupPage,
  type ListingLinkState,
  type ListingPage,
  type ListingReadModel,
  type ListingSummaryExceptions,
  type ListingSyncState,
  type MutationType,
} from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useEffect, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useClient } from "../app/ClientContext";
import { useInstallation } from "../app/InstallationContext";
import {
  applyAnunciosQueryState,
  clearFilters,
  parseAnunciosQueryState,
  toListingListOptions,
  type AnunciosTab,
} from "./anunciosQueryState";
import { anunciosSummaryQuery } from "./anunciosQueries";
import { AnunciosTable } from "./AnunciosTable";
import { ListingDetailPanel } from "./ListingDetailPanel";
import { ListingsSummary } from "./ListingsSummary";
import { ListingsRefreshControl } from "./ListingsRefreshControl";
import { MutationBulkActions } from "./mutations/MutationBulkActions";
import { MutationPreviewModal } from "./mutations/MutationPreviewModal";

const tabs: Array<{ value: AnunciosTab; label: string }> = [
  { value: "todos", label: "Todos" },
  { value: "ativos", label: "Ativos" },
  { value: "pausados", label: "Pausados" },
  { value: "pendencia", label: "Com pendência" },
];

const exceptionLabels = {
  sync_error: "Erro de sync",
  stale: "Desatualizado",
  unlinked: "Sem vínculo",
  below_margin: "Abaixo da margem",
  sem_vinculo: "Sem vínculo",
  abaixo_custo: "Abaixo do custo",
  sem_evidencia: "Sem evidência",
} satisfies Record<ListingException, string>;

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
  conflict: "divergente",
} satisfies Record<ListingLinkState, string>;

// Header exception chips, promoted from the old Resumo panel to sit inline next
// to the title (design/handoff/Anuncios.dc.html:42-53). Each maps a design
// category to its honest summary count + the URL exception filter it deep-links.
// A chip renders ONLY when its count is a real positive number — a null/unknown
// or zero count is never fabricated into a "0" pill (ADR-17). "sem vínculo"
// prefers the sem_vinculo family count, falling back to the unlinked count.
const headerExceptionChips: Array<{
  filterKey: ListingException;
  label: string;
  glyph: string;
  tone: string;
  count: (exceptions: ListingSummaryExceptions) => number | null | undefined;
}> = [
  { filterKey: "sync_error", label: "erro", glyph: "●", tone: "text-warn", count: (e) => e.sync_error },
  { filterKey: "stale", label: "desatualizado", glyph: "●", tone: "text-amber", count: (e) => e.stale },
  { filterKey: "sem_vinculo", label: "sem vínculo", glyph: "○", tone: "text-muted", count: (e) => e.sem_vinculo ?? e.unlinked },
  { filterKey: "below_margin", label: "abaixo da margem", glyph: "●", tone: "text-amber", count: (e) => e.below_margin_worst_case },
];

const numberFormatter = new Intl.NumberFormat("pt-BR");

function csvEscape(value: string): string {
  return /[",\n]/.test(value) ? `"${value.replace(/"/g, '""')}"` : value;
}

// Client-side CSV of the currently visible listings — a non-mutating export of
// what the table already holds (no new fetch, no ML write). Columns mirror the
// table's own columns; unknown values are left blank rather than fabricated.
function buildListingsCsv(items: ListingReadModel[]): string {
  const header = ["MLB", "Título", "Produto", "Preço", "Estoque", "Sync", "Pendência"];
  const rows = items.map((item) => {
    const produto =
      item.link.product_id ??
      (item.link.state === "conflict" ? linkLabels.conflict : linkLabels[item.link.state] ?? "");
    return [
      item.provider_listing_id ?? "",
      item.title ?? "",
      produto,
      item.price?.amount ?? "",
      item.published_quantity ?? "",
      syncLabels[item.sync_state] ?? item.sync_state,
      item.pending_issue?.message_pt ?? "",
    ]
      .map((cell) => csvEscape(String(cell)))
      .join(",");
  });
  return [header.join(","), ...rows].join("\r\n");
}

function ActiveFilterChip({ kind, value, onDismiss }: { kind: string; value: string; onDismiss: () => void }) {
  return (
    <span className="inline-flex items-center gap-2 rounded-full bg-surface-2 px-3 py-1 text-sm font-medium text-muted">
      <span>{`${kind}: ${value}`}</span>
      <button
        type="button"
        aria-label={`Remover filtro ${value}`}
        onClick={onDismiss}
        className="text-faint hover:text-ink"
      >
        ×
      </button>
    </span>
  );
}

function isInvalidFilterError(error: unknown): boolean {
  return isApiError(error) && hasCode(error, "invalid_filter");
}

export function AnunciosPage() {
  const client = useClient();
  const { installationId } = useInstallation();
  const [searchParams, setSearchParams] = useSearchParams();
  const state = parseAnunciosQueryState(searchParams);
  const [cursor, setCursor] = useState<string | undefined>();
  const [cursorStack, setCursorStack] = useState<Array<string | undefined>>([]);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(() => new Set());
  const [openListingId, setOpenListingId] = useState<string | null>(null);
  const [mutationLaunch, setMutationLaunch] = useState<{
    type: MutationType;
    selectedIds: string[];
  } | null>(null);
  const paginationIdentity = JSON.stringify([installationId, state]);
  const previousPaginationIdentity = useRef(paginationIdentity);
  const selectionInstallationId = useRef(installationId);
  const cursorForQuery = previousPaginationIdentity.current === paginationIdentity ? cursor : undefined;
  const pageOptions = {
    ...toListingListOptions(state, installationId),
    ...(cursorForQuery ? { cursor: cursorForQuery } : {}),
  };
  const pageQuery = useQuery<ListingPage | ListingGroupPage>({
    queryKey: state.grouped
      ? listingsQueryKeys.byProduct(installationId, pageOptions)
      : listingsQueryKeys.page(installationId, pageOptions),
    queryFn: () => (state.grouped ? client.listListingsByProduct(pageOptions) : client.listListings(pageOptions)),
    staleTime: QUERY_STALE_TIME.listings,
  });
  const pageData = pageQuery.data;
  const groups = pageData && "groups" in pageData ? pageData.groups : undefined;
  const flatItems = pageData && "items" in pageData ? pageData.items : undefined;
  const visibleListings: ListingReadModel[] = groups
    ? groups.flatMap((group) => group.listings)
    : flatItems ?? [];
  const isEmptyPage = groups ? groups.length === 0 : (flatItems?.length ?? 0) === 0;
  const summaryQuery = useQuery(anunciosSummaryQuery(client, installationId));

  useEffect(() => {
    if (previousPaginationIdentity.current !== paginationIdentity) {
      previousPaginationIdentity.current = paginationIdentity;
      setCursor(undefined);
      setCursorStack([]);
    }
  }, [paginationIdentity]);

  useEffect(() => {
    if (selectionInstallationId.current !== installationId) {
      selectionInstallationId.current = installationId;
      setSelectedIds(new Set());
      setOpenListingId(null);
      setMutationLaunch(null);
    }
  }, [installationId]);

  const installationChanged = selectionInstallationId.current !== installationId;
  const visibleSelectedIds = installationChanged ? new Set<string>() : selectedIds;
  const visibleMutationLaunch = installationChanged ? null : mutationLaunch;

  const updateState = (nextState: typeof state, options?: { replace: boolean }) => {
    setCursor(undefined);
    setCursorStack([]);
    setSearchParams(applyAnunciosQueryState(searchParams, nextState), options);
  };

  const toggleSelection = (listingId: string) => {
    setSelectedIds((current) => {
      const next = new Set(current);
      if (next.has(listingId)) next.delete(listingId);
      else next.add(listingId);
      return next;
    });
  };

  const togglePageSelection = (listingIds: string[]) => {
    setSelectedIds((current) => {
      const next = new Set(current);
      const allSelected = listingIds.every((listingId) => next.has(listingId));
      for (const listingId of listingIds) {
        if (allSelected) next.delete(listingId);
        else next.add(listingId);
      }
      return next;
    });
  };

  const clearSelection = () => setSelectedIds(new Set());

  const exportVisibleCsv = () => {
    if (visibleListings.length === 0) return;
    const csv = buildListingsCsv(visibleListings);
    const blob = new Blob([`﻿${csv}`], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = "anuncios.csv";
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
    URL.revokeObjectURL(url);
  };

  const goToNextPage = () => {
    const nextCursor = pageQuery.data?.next_cursor;
    if (nextCursor !== null && nextCursor !== undefined) {
      setCursorStack((current) => [...current, cursorForQuery]);
      setCursor(nextCursor);
    }
  };

  const goToPreviousPage = () => {
    if (cursorStack.length === 0) return;
    const previousCursor = cursorStack[cursorStack.length - 1];
    setCursorStack((current) => current.slice(0, -1));
    setCursor(previousCursor);
  };

  const summary = summaryQuery.data;
  const activeException = state.filters.exception;
  const chips = summary
    ? headerExceptionChips
        .map((chip) => ({ ...chip, value: chip.count(summary.exceptions) }))
        .filter((chip): chip is typeof chip & { value: number } => typeof chip.value === "number" && chip.value > 0)
    : [];

  return (
    <section aria-labelledby="anuncios-title" className="mx-auto flex max-w-7xl flex-col gap-3.5">
      {/* Header row: title + total · exceptions + inline exception chips + actions.
          Mirrors the dense workspace header of the ratified prototype — no eyebrow,
          no full-width hero, exceptions promoted up here from the old Resumo panel. */}
      <header className="flex flex-wrap items-center gap-x-3 gap-y-2">
        <h1 id="anuncios-title" className="text-[22px] font-bold tracking-tight text-ink">
          Anúncios
        </h1>
        {summary ? (
          <span className="text-[12.5px] text-faint">
            {`${numberFormatter.format(summary.total)} · ${chips.length > 0 ? "exceções:" : "sem exceções"}`}
          </span>
        ) : null}
        <div className="flex flex-wrap items-center gap-2" aria-label="Exceções">
          {chips.map((chip) => {
            const active = activeException === chip.filterKey;
            return (
              <button
                key={chip.filterKey}
                type="button"
                aria-pressed={active}
                aria-label={`${chip.label} ${chip.value}`}
                onClick={() =>
                  updateState({
                    ...state,
                    filters: { ...state.filters, exception: active ? undefined : chip.filterKey },
                  })
                }
                className={`inline-flex items-center gap-1.5 whitespace-nowrap rounded-full border px-3 py-1 text-xs font-semibold ${chip.tone} ${
                  active ? "border-current bg-surface-2" : "border-border bg-surface hover:bg-surface-2"
                }`}
              >
                <span aria-hidden="true" className="text-[10px]">{chip.glyph}</span>
                <span>{chip.label}</span>
                <span className="font-mono">{chip.value}</span>
              </button>
            );
          })}
        </div>
        <div className="ml-auto flex flex-wrap items-center gap-2">
          <label className="inline-flex cursor-pointer items-center gap-2 whitespace-nowrap rounded-lg border border-border bg-surface px-3 py-[7px] text-[12.5px] font-medium text-muted hover:border-border-2">
            <input
              type="checkbox"
              checked={state.grouped}
              onChange={(event) => updateState({ ...state, grouped: event.target.checked })}
              aria-label="Agrupar por produto"
            />
            Agrupar por produto
          </label>
          <button
            type="button"
            onClick={exportVisibleCsv}
            disabled={visibleListings.length === 0}
            className="whitespace-nowrap rounded-lg border border-border bg-surface px-3 py-[7px] text-[12.5px] font-medium text-muted hover:border-border-2 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Exportar
          </button>
          <button
            type="button"
            disabled
            title="disponível em breve"
            className="whitespace-nowrap rounded-lg bg-accent px-3.5 py-2 text-[12.5px] font-semibold text-white disabled:cursor-not-allowed disabled:opacity-50"
          >
            + Criar anúncio
          </button>
        </div>
      </header>

      {/* Status tabs — underline style of the prototype. */}
      <div aria-label="Filtros de status" role="tablist" className="flex gap-0.5 border-b border-border text-sm">
        {tabs.map((tab) => {
          const active = state.tab === tab.value;
          return (
            <button
              key={tab.value}
              type="button"
              role="tab"
              aria-selected={active}
              className={`whitespace-nowrap border-b-2 px-4 py-2 transition-colors ${
                active
                  ? "border-accent font-semibold text-accent-ink"
                  : "border-transparent font-normal text-muted hover:text-ink"
              }`}
              onClick={() => updateState({ ...state, tab: tab.value })}
            >
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Compact toolbar: search + refresh + active filter chips. The full-width
          Buscar hero is demoted to this inline field so the table leads the page. */}
      <div className="flex flex-wrap items-center gap-2">
        <label htmlFor="anuncios-search" className="sr-only">
          Buscar anúncios
        </label>
        <input
          id="anuncios-search"
          type="search"
          value={state.q}
          onChange={(event) => updateState({ ...state, q: event.target.value }, { replace: true })}
          placeholder="Buscar SKU, MLB, título…"
          className="w-full max-w-xs rounded-lg border border-border bg-surface px-3 py-1.5 text-sm text-ink outline-none transition focus:border-accent focus:ring-2 focus:ring-accent-soft"
        />
        <ListingsRefreshControl installationId={installationId} />
        <div className="flex flex-wrap items-center gap-2" aria-label="Filtros ativos">
          {state.filters.exception ? (
            <ActiveFilterChip
              kind="Exceção"
              value={exceptionLabels[state.filters.exception]}
              onDismiss={() => updateState({ ...state, filters: { ...state.filters, exception: undefined } })}
            />
          ) : null}
          {state.filters.sync_state ? (
            <ActiveFilterChip
              kind="Sync"
              value={syncLabels[state.filters.sync_state]}
              onDismiss={() => updateState({ ...state, filters: { ...state.filters, sync_state: undefined } })}
            />
          ) : null}
          {state.filters.link_state ? (
            <ActiveFilterChip
              kind="Vínculo"
              value={linkLabels[state.filters.link_state]}
              onDismiss={() => updateState({ ...state, filters: { ...state.filters, link_state: undefined } })}
            />
          ) : null}
          {state.filters.listing_type_code ? (
            <ActiveFilterChip
              kind="Modalidade"
              value={state.filters.listing_type_code}
              onDismiss={() => updateState({ ...state, filters: { ...state.filters, listing_type_code: undefined } })}
            />
          ) : null}
        </div>
      </div>

      {/* Bulk selection bar — surfaces only when rows are selected. */}
      {visibleSelectedIds.size > 0 ? (
        <div className="flex flex-wrap items-center gap-3 rounded-lg border border-accent bg-accent-soft px-4 py-2.5 text-sm">
          <b className="text-accent-ink">{visibleSelectedIds.size} selecionados</b>
          <MutationBulkActions
            selectedIds={visibleSelectedIds}
            onOpen={(type) => {
              if (visibleSelectedIds.size === 0) return;
              setMutationLaunch({ type, selectedIds: Array.from(visibleSelectedIds) });
            }}
          />
          <span className="text-xs text-faint">preview antes de aplicar</span>
          <button type="button" onClick={clearSelection} className="ml-auto text-muted hover:text-ink">
            ✕ Limpar
          </button>
        </div>
      ) : null}

      {/* Content row: table card (flex) + inline detail drawer (300px sticky). */}
      <div className="flex min-w-0 items-start gap-3.5">
        <div className="min-w-0 flex-1 overflow-hidden rounded-xl border border-border bg-surface">
          {pageQuery.isPending ? (
            <div className="p-4">
              <LoadingState />
            </div>
          ) : pageQuery.isError ? (
            <div className="p-4">
              <ErrorState
                onRetry={() => {
                  if (isInvalidFilterError(pageQuery.error)) {
                    setSearchParams(clearFilters(searchParams));
                  } else {
                    void pageQuery.refetch();
                  }
                }}
              />
            </div>
          ) : pageQuery.data && isEmptyPage ? (
            <div className="p-4">
              <EmptyState
                hint={
                  <button
                    type="button"
                    className="font-medium text-accent underline decoration-accent-soft underline-offset-2 hover:text-accent-ink"
                    onClick={() => setSearchParams(clearFilters(searchParams))}
                  >
                    Limpar filtros
                  </button>
                }
              />
            </div>
          ) : pageQuery.data ? (
            <>
              <AnunciosTable
                items={flatItems}
                groups={groups}
                asOf={pageQuery.data.as_of}
                selectedIds={visibleSelectedIds}
                onToggle={toggleSelection}
                onTogglePage={togglePageSelection}
                onOpen={setOpenListingId}
              />
              <nav
                className="flex items-center justify-end gap-2 border-t border-border-2 px-3 py-3"
                aria-label="Paginação de anúncios"
              >
                <button
                  type="button"
                  onClick={goToPreviousPage}
                  disabled={cursorStack.length === 0}
                  className="rounded-lg border px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Anterior
                </button>
                <button
                  type="button"
                  onClick={goToNextPage}
                  disabled={pageQuery.data.next_cursor === null}
                  className="rounded-lg border px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Próxima
                </button>
              </nav>
            </>
          ) : null}
        </div>
        <ListingDetailPanel listingId={openListingId} onClose={() => setOpenListingId(null)} />
      </div>

      {/* Demoted Resumo — counters only, secondary, below the table. */}
      <section aria-labelledby="anuncios-summary-title" className="border-t border-border-2 pt-3">
        <h2 id="anuncios-summary-title" className="sr-only">
          Resumo
        </h2>
        <ListingsSummary
          isPending={summaryQuery.isPending}
          isError={summaryQuery.isError}
          data={summaryQuery.data}
          onRetry={() => void summaryQuery.refetch()}
        />
      </section>

      {visibleMutationLaunch ? (
        <MutationPreviewModal
          key={`${visibleMutationLaunch.type}:${visibleMutationLaunch.selectedIds.join(",")}`}
          open
          type={visibleMutationLaunch.type}
          installationId={installationId}
          selectedIds={visibleMutationLaunch.selectedIds}
          onClose={() => setMutationLaunch(null)}
        />
      ) : null}
    </section>
  );
}
