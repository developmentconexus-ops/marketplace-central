import { useQuery } from "@tanstack/react-query";
import type { DashboardDegradedSource, DashboardOverview } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState, StatCard } from "@marketplace-central/ui";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { Atalhos } from "./Atalhos";
import { FilaDeAtencao } from "./FilaDeAtencao";
import { PedidosRecentes } from "./PedidosRecentes";

// Honest-absent glyph: matches the UNKNOWN_KPI_VALUE convention in pedidos/PedidosPage.tsx:33
// (cited/reused, not forked) so a null source and a real 0 are never visually confused (ADR-17).
const UNKNOWN_KPI_VALUE = "—";

function formatKpi(value: number | null): string {
  return value === null ? UNKNOWN_KPI_VALUE : String(value);
}

function isDegraded(degraded: DashboardDegradedSource[], source: DashboardDegradedSource): boolean {
  return degraded.includes(source);
}

interface KpiCardConfig {
  key: string;
  label: string;
  sub?: string;
  value: number | null;
  source: DashboardDegradedSource;
}

function KpiCard({ config, degraded, onRetry }: {
  config: KpiCardConfig;
  degraded: DashboardDegradedSource[];
  onRetry: () => void;
}) {
  if (config.value === null && isDegraded(degraded, config.source)) {
    return (
      <div className="bg-surface border border-border rounded-card p-5">
        <p className="text-xs font-medium text-muted uppercase tracking-wide">{config.label}</p>
        <div className="mt-2">
          <ErrorState onRetry={onRetry} />
        </div>
      </div>
    );
  }
  return <StatCard label={config.label} value={formatKpi(config.value)} sub={config.sub} />;
}

function formatLastImportAge(ageSeconds: number): string {
  if (ageSeconds < 60) return "há menos de 1 min";
  const minutes = Math.floor(ageSeconds / 60);
  if (minutes < 60) return `há ${minutes} min`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `há ${hours} h`;
  const days = Math.floor(hours / 24);
  return `há ${days} d`;
}

function UltimoImportCard({ summary, onRetry }: { summary: DashboardOverview; onRetry: () => void }) {
  if (summary.last_import) {
    const at = new Date(summary.last_import.at);
    return (
      <StatCard
        label="Último import ERP"
        value={at.toLocaleString("pt-BR")}
        sub={formatLastImportAge(summary.last_import.age_seconds)}
      />
    );
  }
  if (isDegraded(summary.degraded, "erp_import")) {
    return (
      <div className="bg-surface border border-border rounded-card p-5">
        <p className="text-xs font-medium text-muted uppercase tracking-wide">Último import ERP</p>
        <div className="mt-2">
          <ErrorState onRetry={onRetry} />
        </div>
      </div>
    );
  }
  return (
    <div className="bg-surface border border-border rounded-card p-5">
      <p className="text-xs font-medium text-muted uppercase tracking-wide">Último import ERP</p>
      <div className="mt-2">
        <EmptyState hint="Nenhuma importação ainda" />
      </div>
    </div>
  );
}

export function DashboardPage() {
  const client = useClient();
  const { installationId } = useInstallation();

  const query = useQuery({
    queryKey: ["dashboard", "summary", { installation_id: installationId }],
    queryFn: async () =>
      (await client.getDashboardSummary(installationId)) as unknown as DashboardOverview,
    refetchOnWindowFocus: false,
  });

  return (
    <section aria-labelledby="dashboard-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-accent-ink">Visão geral</p>
        <h1 id="dashboard-title" className="text-2xl font-semibold tracking-tight text-ink">
          Visão geral
        </h1>
        <p className="max-w-2xl text-sm text-muted">
          Resumo do dia: vendas, anúncios com exceção, vínculos e o último import ERP.
        </p>
      </header>

      {query.isPending ? (
        <LoadingState />
      ) : query.isError ? (
        <ErrorState onRetry={() => void query.refetch()} />
      ) : (
        <>
          <div
            aria-label="Indicadores do painel"
            className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4"
          >
            <KpiCard
              config={{
                key: "orders_today",
                label: "Vendas hoje",
                sub: "hoje",
                value: query.data.orders_today,
                source: "orders",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "orders_7d",
                label: "Vendas 7 dias",
                sub: "últimos 7 dias",
                value: query.data.orders_7d,
                source: "orders",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "anuncios_ativos",
                label: "Anúncios ativos",
                value: query.data.anuncios_ativos,
                source: "listings",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "sync_errors",
                label: "Anúncios com erro de sync",
                value: query.data.sync_errors,
                source: "listings",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "below_margin",
                label: "Anúncios abaixo da margem",
                value: query.data.below_margin,
                source: "listings",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "pending_links",
                label: "Produtos sem vínculo",
                value: query.data.pending_links,
                source: "linkage",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <KpiCard
              config={{
                key: "missing_gtin",
                label: "Sem GTIN",
                value: query.data.missing_gtin,
                source: "linkage",
              }}
              degraded={query.data.degraded}
              onRetry={() => void query.refetch()}
            />
            <UltimoImportCard summary={query.data} onRetry={() => void query.refetch()} />
          </div>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
            <div className="flex flex-col gap-6 lg:col-span-2">
              <FilaDeAtencao summary={query.data} />
              <PedidosRecentes />
            </div>
            <Atalhos />
          </div>
        </>
      )}
    </section>
  );
}
