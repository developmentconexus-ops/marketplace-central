import { useQuery } from "@tanstack/react-query";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { useSearchParams } from "react-router-dom";
import { useClient } from "../app/ClientContext";
import { useInstallation } from "../app/InstallationContext";
import {
  applyAnunciosQueryState,
  clearFilters,
  parseAnunciosQueryState,
  type AnunciosTab,
} from "./anunciosQueryState";
import { anunciosPageQuery, anunciosSummaryQuery } from "./anunciosQueries";
import { AnunciosTable } from "./AnunciosTable";
import { ListingsSummary } from "./ListingsSummary";

const tabs: Array<{ value: AnunciosTab; label: string }> = [
  { value: "todos", label: "Todos" },
  { value: "ativos", label: "Ativos" },
  { value: "pausados", label: "Pausados" },
  { value: "pendencia", label: "Com pendência" },
];

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
  const pageQuery = useQuery(anunciosPageQuery(client, installationId, state));
  const summaryQuery = useQuery(anunciosSummaryQuery(client, installationId));

  const updateState = (nextState: typeof state, options?: { replace: boolean }) => {
    setSearchParams(applyAnunciosQueryState(searchParams, nextState), options);
  };

  return (
    <section aria-labelledby="anuncios-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-blue-700">Workspace operacional</p>
        <h1 id="anuncios-title" className="text-2xl font-semibold tracking-tight text-slate-950">
          Anúncios
        </h1>
        <p className="max-w-2xl text-sm text-slate-600">
          Acompanhe publicação, vínculo e sincronização dos anúncios da instalação selecionada.
        </p>
      </header>

      <div className="flex flex-col gap-4 border-b border-slate-200 pb-4">
        <div aria-label="Filtros de status" role="tablist" className="flex flex-wrap gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.value}
              type="button"
              role="tab"
              aria-selected={state.tab === tab.value}
              className={`rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
                state.tab === tab.value
                  ? "border-blue-600 bg-blue-600 text-white"
                  : "border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:bg-slate-50"
              }`}
              onClick={() => updateState({ ...state, tab: tab.value })}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <label className="flex max-w-xl flex-col gap-1 text-sm font-medium text-slate-700" htmlFor="anuncios-search">
          Buscar anúncios
          <input
            id="anuncios-search"
            type="search"
            value={state.q}
            onChange={(event) => updateState({ ...state, q: event.target.value }, { replace: true })}
            placeholder="Título, SKU ou MLB"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 font-normal text-slate-900 outline-none transition focus:border-blue-600 focus:ring-2 focus:ring-blue-100"
          />
        </label>
      </div>

      <section aria-labelledby="anuncios-summary-title" className="rounded-xl border border-slate-200 bg-white p-4">
        <h2 id="anuncios-summary-title" className="text-sm font-semibold text-slate-900">
          Resumo
        </h2>
        <ListingsSummary
          isPending={summaryQuery.isPending}
          isError={summaryQuery.isError}
          data={summaryQuery.data}
          onRetry={() => void summaryQuery.refetch()}
        />
      </section>

      <section aria-labelledby="anuncios-list-title" className="rounded-xl border border-slate-200 bg-white p-4">
        <h2 id="anuncios-list-title" className="text-sm font-semibold text-slate-900">
          Lista de anúncios
        </h2>
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
        ) : pageQuery.data && pageQuery.data.items.length === 0 ? (
          <EmptyState
            hint={
              <button
                type="button"
                className="font-medium text-blue-700 underline decoration-blue-300 underline-offset-2 hover:text-blue-800"
                onClick={() => setSearchParams(clearFilters(searchParams))}
              >
                Limpar filtros
              </button>
            }
          />
        ) : pageQuery.data ? (
          <AnunciosTable items={pageQuery.data.items} asOf={pageQuery.data.as_of} />
        ) : null}
      </section>
    </section>
  );
}
