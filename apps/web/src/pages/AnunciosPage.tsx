import { useQuery } from "@tanstack/react-query";
import type {
  ListingException,
  ListingGroupPage,
  ListingLinkState,
  ListingPage,
  ListingSyncState,
  MutationType,
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
  return (
    typeof error === "object" &&
    error !== null &&
    (error as { error?: { code?: string } }).error?.code === "invalid_filter"
  );
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

  return (
    <section aria-labelledby="anuncios-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-accent">Workspace operacional</p>
        <h1 id="anuncios-title" className="text-2xl font-semibold tracking-tight text-ink">
          Anúncios
        </h1>
        <p className="max-w-2xl text-sm text-muted">
          Acompanhe publicação, vínculo e sincronização dos anúncios da instalação selecionada.
        </p>
      </header>

      <div className="flex flex-col gap-4 border-b border-border pb-4">
        <div aria-label="Filtros de status" role="tablist" className="flex flex-wrap gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.value}
              type="button"
              role="tab"
              aria-selected={state.tab === tab.value}
              className={`rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
                state.tab === tab.value
                  ? "border-accent bg-accent text-white"
                  : "border-border bg-white text-muted hover:border-border-2 hover:bg-surface-2"
              }`}
              onClick={() => updateState({ ...state, tab: tab.value })}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <div className="flex flex-wrap gap-2" aria-label="Filtros ativos">
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
        <label className="flex items-center gap-2 text-sm font-medium text-muted">
          <input
            type="checkbox"
            checked={state.grouped}
            onChange={(event) => updateState({ ...state, grouped: event.target.checked })}
            aria-label="Agrupar por produto"
          />
          Agrupar por produto
        </label>
        <label className="flex max-w-xl flex-col gap-1 text-sm font-medium text-muted" htmlFor="anuncios-search">
          Buscar anúncios
          <input
            id="anuncios-search"
            type="search"
            value={state.q}
            onChange={(event) => updateState({ ...state, q: event.target.value }, { replace: true })}
            placeholder="Título, SKU ou MLB"
            className="rounded-lg border border-border bg-white px-3 py-2 font-normal text-ink outline-none transition focus:border-accent focus:ring-2 focus:ring-accent-soft"
          />
        </label>
      </div>

      <section aria-labelledby="anuncios-summary-title" className="rounded-xl border border-border bg-white p-4">
        <h2 id="anuncios-summary-title" className="text-sm font-semibold text-ink">
          Resumo
        </h2>
        <ListingsSummary
          isPending={summaryQuery.isPending}
          isError={summaryQuery.isError}
          data={summaryQuery.data}
          onRetry={() => void summaryQuery.refetch()}
          activeException={state.filters.exception}
          onExceptionChipClick={(exception) =>
            updateState({ ...state, filters: { ...state.filters, exception } })
          }
        />
      </section>

      <section aria-labelledby="anuncios-list-title" className="rounded-xl border border-border bg-white p-4">
        <h2 id="anuncios-list-title" className="text-sm font-semibold text-ink">
          Lista de anúncios
        </h2>
        <div className="mt-3 flex flex-wrap items-center gap-2" aria-label="Ações de seleção">
          <ListingsRefreshControl installationId={installationId} />
          {visibleSelectedIds.size > 0 ? (
            <span className="rounded-full bg-accent-soft px-3 py-1 text-sm font-medium text-accent-ink">
              {visibleSelectedIds.size} selecionado(s)
            </span>
          ) : null}
          <MutationBulkActions
            selectedIds={visibleSelectedIds}
            onOpen={(type) => {
              if (visibleSelectedIds.size === 0) return;
              setMutationLaunch({ type, selectedIds: Array.from(visibleSelectedIds) });
            }}
          />
        </div>
        {pageQuery.isPending ? <LoadingState /> : pageQuery.isError ? (
          <ErrorState
            onRetry={() => {
              if (isInvalidFilterError(pageQuery.error)) {
                setSearchParams(clearFilters(searchParams));
              } else {
                void pageQuery.refetch();
              }
            }}
          />
        ) : pageQuery.data && isEmptyPage ? (
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
            <nav className="mt-4 flex items-center justify-end gap-2" aria-label="Paginação de anúncios">
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
      </section>
      <ListingDetailPanel listingId={openListingId} onClose={() => setOpenListingId(null)} />
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
