import { useQuery } from "@tanstack/react-query";
import type { ProductLinkCandidateItem, ProductLinkWorkflowItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { ImportacaoSection } from "./ImportacaoSection";
import { QueueTab } from "./QueueTab";
import { ResolvidosTab } from "./ResolvidosTab";
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

function Kpi({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="rounded-card border border-border bg-surface-2 p-3">
      <dt className="text-xs font-medium text-faint">{label}</dt>
      <dd className="font-mono text-xl font-semibold text-ink">{value}</dd>
    </div>
  );
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

  const kpiValue = (query: { isPending: boolean; isError: boolean }, value: number) =>
    query.isPending ? "…" : query.isError ? "—" : value;

  return (
    <section aria-labelledby="vinculos-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-accent">Workspace operacional</p>
        <h1 id="vinculos-title" className="text-2xl font-semibold tracking-tight text-ink">
          Vínculos
        </h1>
        <p className="max-w-2xl text-sm text-muted">
          Revise candidatos de vínculo entre anúncios e produtos internos, e acompanhe o que já foi resolvido.
        </p>
      </header>

      <section aria-labelledby="vinculos-kpis-title" className="rounded-card border border-border bg-surface p-4">
        <h2 id="vinculos-kpis-title" className="text-sm font-semibold text-ink">
          Resumo
        </h2>
        <dl className="mt-3 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <Kpi label="Pendentes" value={kpiValue(queueQuery, pendentesCount)} />
          <Kpi label="Alta confiança" value={kpiValue(queueQuery, altaConfiancaCount)} />
          <Kpi label="Resolvidos hoje" value={kpiValue(resolvedQuery, resolvidosHojeCount)} />
        </dl>
      </section>

      <div
        aria-label="Filtros de vínculos"
        role="tablist"
        className="inline-flex w-fit overflow-hidden rounded-lg border border-border text-sm"
      >
        {tabs.map((tabItem) => (
          <button
            key={tabItem.value}
            type="button"
            role="tab"
            aria-selected={tab === tabItem.value}
            className={`px-4 py-2 font-medium transition-colors ${
              tab === tabItem.value
                ? "bg-accent-soft text-accent-ink"
                : "bg-surface text-muted hover:text-ink"
            }`}
            onClick={() => setTab(tabItem.value)}
          >
            {tabItem.label}
          </button>
        ))}
      </div>

      {tab === "fila" ? (
        <section aria-labelledby="vinculos-fila-title" className="rounded-card border border-border bg-surface p-4">
          <h2 id="vinculos-fila-title" className="text-sm font-semibold text-ink">
            Fila de candidatos
          </h2>
          {queueQuery.isPending ? (
            <LoadingState />
          ) : queueQuery.isError ? (
            <ErrorState onRetry={() => void queueQuery.refetch()} />
          ) : queueItems.length === 0 ? (
            <EmptyState />
          ) : (
            <QueueTab installationId={installationId} onViewResolved={() => setTab("resolvidos")} />
          )}
        </section>
      ) : (
        <section aria-labelledby="vinculos-resolvidos-title" className="rounded-card border border-border bg-surface p-4">
          <h2 id="vinculos-resolvidos-title" className="text-sm font-semibold text-ink">
            Resolvidos
          </h2>
          <ResolvidosTab installationId={installationId} />
        </section>
      )}

      <ImportacaoSection />
    </section>
  );
}
