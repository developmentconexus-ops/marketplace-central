import { useQuery } from "@tanstack/react-query";
import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { ordersQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { PedidosTable } from "./PedidosTable";
import { filterOrdersByTab, pedidosTabs, type PedidosTab } from "./pedidosTabs";

export function PedidosPage() {
  const client = useClient();
  const { installationId } = useInstallation();
  const [tab, setTab] = useState<PedidosTab>("todos");

  const ordersQuery = useQuery({
    queryKey: ordersQueryKeys.list(installationId, {}),
    queryFn: () => client.listOrders({ installation_id: installationId }),
    staleTime: QUERY_STALE_TIME.orders,
  });

  const allItems: OrderRead[] = ordersQuery.data?.items ?? [];
  const visibleItems = filterOrdersByTab(allItems, tab);

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

      <div className="flex flex-col gap-4 border-b border-border pb-4">
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
      </div>

      <section aria-labelledby="pedidos-list-title" className="rounded-card border border-border bg-surface p-4">
        <h2 id="pedidos-list-title" className="text-sm font-semibold text-ink">
          Lista de pedidos
        </h2>
        <div className="mt-3">
          {ordersQuery.isPending ? (
            <LoadingState />
          ) : ordersQuery.isError ? (
            <ErrorState onRetry={() => void ordersQuery.refetch()} />
          ) : visibleItems.length === 0 ? (
            <EmptyState />
          ) : (
            <PedidosTable items={visibleItems} />
          )}
        </div>
      </section>
    </section>
  );
}
