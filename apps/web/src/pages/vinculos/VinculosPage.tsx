import { useQuery } from "@tanstack/react-query";
import type { ProductLinkCandidateItem, ProductLinkWorkflowItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { vinculosQueryKeys } from "./vinculosQueryKeys";

type VinculosTab = "fila" | "resolvidos";

const tabs: Array<{ value: VinculosTab; label: string }> = [
  { value: "fila", label: "Fila" },
  { value: "resolvidos", label: "Resolvidos" },
];

function isResolved(item: ProductLinkWorkflowItem): boolean {
  return item.current_link?.state === "resolved";
}

function isResolvedToday(item: ProductLinkWorkflowItem): boolean {
  const link = item.current_link;
  if (!link || link.state !== "resolved") return false;
  const updatedAt = new Date(link.updated_at);
  if (Number.isNaN(updatedAt.getTime())) return false;
  const now = new Date();
  return (
    updatedAt.getFullYear() === now.getFullYear() &&
    updatedAt.getMonth() === now.getMonth() &&
    updatedAt.getDate() === now.getDate()
  );
}

function countHighConfidence(items: ProductLinkCandidateItem[]): number {
  return items.filter((item) => item.confidence_band === "ALTA").length;
}

export function VinculosPage() {
  const client = useClient();
  const { installationId } = useInstallation();
  const [tab, setTab] = useState<VinculosTab>("fila");

  const queueQuery = useQuery({
    queryKey: vinculosQueryKeys.queue(installationId),
    queryFn: () => client.listProductLinkCandidates(installationId),
    staleTime: QUERY_STALE_TIME.listings,
  });

  const resolvedQuery = useQuery({
    queryKey: vinculosQueryKeys.resolved(installationId),
    queryFn: () => client.listProductLinkWorkflows(installationId),
    staleTime: QUERY_STALE_TIME.listings,
  });

  const queueItems = queueQuery.data?.items ?? [];
  const resolvedItems = (resolvedQuery.data?.items ?? []).filter(isResolved);
  const pendentesCount = queueItems.length;
  const altaConfiancaCount = countHighConfidence(queueItems);
  const resolvidosHojeCount = resolvedItems.filter(isResolvedToday).length;

  return (
    <section aria-labelledby="vinculos-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-blue-700">Workspace operacional</p>
        <h1 id="vinculos-title" className="text-2xl font-semibold tracking-tight text-slate-950">
          Vínculos
        </h1>
        <p className="max-w-2xl text-sm text-slate-600">
          Revise candidatos de vínculo entre anúncios e produtos internos, e acompanhe o que já foi resolvido.
        </p>
      </header>

      <section aria-labelledby="vinculos-kpis-title" className="rounded-xl border border-slate-200 bg-white p-4">
        <h2 id="vinculos-kpis-title" className="text-sm font-semibold text-slate-900">
          Resumo
        </h2>
        <dl className="mt-3 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div className="rounded-lg border border-slate-200 bg-slate-50 p-3">
            <dt className="text-xs font-medium text-slate-500">Pendentes</dt>
            <dd className="text-xl font-semibold text-slate-950">
              {queueQuery.isPending ? "…" : queueQuery.isError ? "—" : pendentesCount}
            </dd>
          </div>
          <div className="rounded-lg border border-slate-200 bg-slate-50 p-3">
            <dt className="text-xs font-medium text-slate-500">Alta confiança</dt>
            <dd className="text-xl font-semibold text-slate-950">
              {queueQuery.isPending ? "…" : queueQuery.isError ? "—" : altaConfiancaCount}
            </dd>
          </div>
          <div className="rounded-lg border border-slate-200 bg-slate-50 p-3">
            <dt className="text-xs font-medium text-slate-500">Resolvidos hoje</dt>
            <dd className="text-xl font-semibold text-slate-950">
              {resolvedQuery.isPending ? "…" : resolvedQuery.isError ? "—" : resolvidosHojeCount}
            </dd>
          </div>
        </dl>
      </section>

      <div className="flex flex-col gap-4 border-b border-slate-200 pb-4">
        <div aria-label="Filtros de vínculos" role="tablist" className="flex flex-wrap gap-2">
          {tabs.map((tabItem) => (
            <button
              key={tabItem.value}
              type="button"
              role="tab"
              aria-selected={tab === tabItem.value}
              className={`rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
                tab === tabItem.value
                  ? "border-blue-600 bg-blue-600 text-white"
                  : "border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:bg-slate-50"
              }`}
              onClick={() => setTab(tabItem.value)}
            >
              {tabItem.label}
            </button>
          ))}
        </div>
      </div>

      {tab === "fila" ? (
        <section aria-labelledby="vinculos-fila-title" className="rounded-xl border border-slate-200 bg-white p-4">
          <h2 id="vinculos-fila-title" className="text-sm font-semibold text-slate-900">
            Fila de candidatos
          </h2>
          {queueQuery.isPending ? (
            <LoadingState />
          ) : queueQuery.isError ? (
            <ErrorState onRetry={() => void queueQuery.refetch()} />
          ) : queueItems.length === 0 ? (
            <EmptyState />
          ) : (
            // TODO(S7): render the queue table rows + drawer detail here.
            <div className="mt-3 text-sm text-slate-500">
              {queueItems.length} candidato(s) na fila — tabela detalhada chega na próxima etapa.
            </div>
          )}
        </section>
      ) : (
        <section aria-labelledby="vinculos-resolvidos-title" className="rounded-xl border border-slate-200 bg-white p-4">
          <h2 id="vinculos-resolvidos-title" className="text-sm font-semibold text-slate-900">
            Resolvidos
          </h2>
          {resolvedQuery.isPending ? (
            <LoadingState />
          ) : resolvedQuery.isError ? (
            <ErrorState onRetry={() => void resolvedQuery.refetch()} />
          ) : resolvedItems.length === 0 ? (
            <EmptyState />
          ) : (
            // TODO(S7): render the resolved links table here.
            <div className="mt-3 text-sm text-slate-500">
              {resolvedItems.length} vínculo(s) — tabela detalhada chega na próxima etapa.
            </div>
          )}
        </section>
      )}
    </section>
  );
}
