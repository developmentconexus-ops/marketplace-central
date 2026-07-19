import { useQuery } from "@tanstack/react-query";
import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState, StatCard } from "@marketplace-central/ui";
import { ordersQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState, type ReactNode } from "react";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { FilaView } from "./FilaView";
import { PedidosTable } from "./PedidosTable";
import { filterOrdersByTab, pedidosTabs, type PedidosTab } from "./pedidosTabs";

type PedidosView = "fila" | "lista" | "kanban";

const viewOptions: { value: PedidosView; label: string }[] = [
  { value: "fila", label: "Fila" },
  { value: "lista", label: "Lista" },
  { value: "kanban", label: "Kanban" },
];

const viewHeadings: Record<PedidosView, string> = {
  fila: "Fila de trabalho",
  lista: "Lista de pedidos",
  kanban: "Kanban",
};

// UnknownValue's rendered node ("—") can't be passed through StatCard.value (string | number
// only) without forking packages/ui/StatCard.tsx, which this slice may only consume. The literal
// em dash below matches UnknownValue's glyph so the KPI still reads as an honest unknown, not a 0.
const UNKNOWN_KPI_VALUE = "—";

interface KpiCardProps {
  label: string;
  value: string;
  sub?: string;
  onSelect: () => void;
}

function KpiCard({ label, value, sub, onSelect }: KpiCardProps) {
  return (
    <button
      type="button"
      onClick={onSelect}
      className="w-full rounded-card text-left transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
    >
      <StatCard label={label} value={value} sub={sub} />
    </button>
  );
}

function KanbanPlaceholder() {
  return (
    <div className="rounded-card border border-dashed border-border bg-surface-2 p-8 text-center text-sm text-muted">
      Kanban por etapa — em breve (F02-S3)
    </div>
  );
}

export function PedidosPage() {
  const client = useClient();
  const { installationId } = useInstallation();
  const [tab, setTab] = useState<PedidosTab>("todos");
  const [view, setView] = useState<PedidosView>("fila");

  const ordersQuery = useQuery({
    queryKey: ordersQueryKeys.list(installationId, {}),
    queryFn: () => client.listOrders({ installation_id: installationId }),
    staleTime: QUERY_STALE_TIME.orders,
  });

  // No summary key exists yet in the web-query barrel (checked packages/web-query/src/index.ts);
  // a local, stable key is used here per the dispatch pack's fallback instruction.
  const summaryQuery = useQuery({
    queryKey: ["orders", "summary", installationId, "by_status"] as const,
    queryFn: () => client.getOrderSummary(installationId, { by: "status" }),
    staleTime: QUERY_STALE_TIME.orders,
  });

  const allItems: OrderRead[] = ordersQuery.data?.items ?? [];
  const visibleItems = filterOrdersByTab(allItems, tab);
  const byStatus = summaryQuery.data?.by_status;

  const goToLista = () => setView("lista");
  const goToFila = () => setView("fila");

  // Loading/Error gate the whole list area (as before); Empty is per-view since Lista keeps its
  // own filter tabs visible even when the current tab's slice is empty.
  let body: ReactNode;
  if (ordersQuery.isPending) {
    body = <LoadingState />;
  } else if (ordersQuery.isError) {
    body = <ErrorState onRetry={() => void ordersQuery.refetch()} />;
  } else if (view === "fila") {
    body = allItems.length === 0 ? <EmptyState /> : <FilaView items={allItems} />;
  } else if (view === "kanban") {
    body = <KanbanPlaceholder />;
  } else {
    body = (
      <div className="flex flex-col gap-4">
        <div aria-label="Filtros de pedidos" role="tablist" className="flex flex-wrap gap-2">
          {pedidosTabs.map((tabItem) => (
            <button
              key={tabItem.value}
              type="button"
              role="tab"
              aria-selected={tab === tabItem.value}
              disabled={tabItem.disabled}
              title={tabItem.disabled ? "disponível em breve" : undefined}
              className={`rounded-lg border px-3 py-2 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
                tab === tabItem.value
                  ? "border-accent bg-accent text-white"
                  : "border-border bg-surface text-ink hover:border-accent"
              }`}
              onClick={() => {
                if (!tabItem.disabled) setTab(tabItem.value);
              }}
            >
              {tabItem.label}
            </button>
          ))}
        </div>
        {visibleItems.length === 0 ? <EmptyState /> : <PedidosTable items={visibleItems} />}
      </div>
    );
  }

  return (
    <section aria-labelledby="pedidos-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-accent-ink">Workspace operacional</p>
        <h1 id="pedidos-title" className="text-2xl font-semibold tracking-tight text-ink">
          Pedidos
        </h1>
        <p className="max-w-2xl text-sm text-muted">
          Acompanhe os pedidos recebidos, o vínculo com produtos internos e os prazos de SLA.
        </p>
      </header>

      <div
        aria-label="Indicadores de pedidos"
        className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5"
      >
        <KpiCard
          label="NOVOS"
          value={byStatus ? String(byStatus.novo) : UNKNOWN_KPI_VALUE}
          onSelect={goToLista}
        />
        <KpiCard
          label="A FATURAR"
          value={byStatus ? String(byStatus.faturar) : UNKNOWN_KPI_VALUE}
          onSelect={goToLista}
        />
        <KpiCard
          label="A ENVIAR"
          value={byStatus ? String(byStatus.enviar) : UNKNOWN_KPI_VALUE}
          onSelect={goToLista}
        />
        {/* seven_days on OrderSummary is a TOTAL-orders count, not an enviado-in-7d bucket — a
            "últimos 7d" sub-note here would misrepresent it (ADR-17), so it is omitted. */}
        <KpiCard
          label="ENVIADOS"
          value={byStatus ? String(byStatus.enviado) : UNKNOWN_KPI_VALUE}
          onSelect={goToLista}
        />
        <KpiCard
          label="DIFAL A PAGAR"
          value={UNKNOWN_KPI_VALUE}
          sub="decomposição pendente"
          onSelect={goToFila}
        />
      </div>

      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border pb-4">
        <h2 id="pedidos-list-title" className="text-sm font-semibold text-ink">
          {viewHeadings[view]}
        </h2>
        <div
          role="tablist"
          aria-label="Alternar visualização de pedidos"
          className="inline-flex overflow-hidden rounded-lg border border-border text-sm"
        >
          {viewOptions.map((option) => (
            <button
              key={option.value}
              type="button"
              role="tab"
              aria-selected={view === option.value}
              onClick={() => setView(option.value)}
              className={`px-4 py-2 font-medium transition-colors ${
                view === option.value
                  ? "bg-accent-soft text-accent-ink"
                  : "bg-surface text-muted hover:text-ink"
              }`}
            >
              {option.label}
            </button>
          ))}
        </div>
      </div>

      <section aria-labelledby="pedidos-list-title" className="rounded-card border border-border bg-surface p-4">
        <div className="mt-1">{body}</div>
      </section>
    </section>
  );
}
