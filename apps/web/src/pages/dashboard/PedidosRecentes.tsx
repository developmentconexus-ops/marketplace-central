import { useQuery } from "@tanstack/react-query";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { Link } from "react-router-dom";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";

const RECENT_ORDERS_LIMIT = 5;

export function PedidosRecentes() {
  const client = useClient();
  const { installationId } = useInstallation();

  // Second, dedicated query: recent orders are not part of the summary endpoint (only aggregate
  // counts live there), so the smallest honest way to show them is a direct read of the public
  // orders list — no backend change, no reshaping of the summary contract.
  const query = useQuery({
    queryKey: ["dashboard", "pedidos-recentes", { installation_id: installationId }],
    queryFn: () => client.listOrders({ installation_id: installationId, limit: RECENT_ORDERS_LIMIT }),
    refetchOnWindowFocus: false,
  });

  return (
    <section aria-labelledby="pedidos-recentes-title" className="rounded-card border border-border bg-surface p-4">
      <div className="flex items-center justify-between">
        <h2 id="pedidos-recentes-title" className="text-sm font-semibold text-ink">
          Pedidos recentes
        </h2>
        <Link to="/pedidos" className="text-sm font-medium text-accent-ink hover:underline">
          Ver todos
        </Link>
      </div>

      <div className="mt-2">
        {query.isPending ? (
          <LoadingState />
        ) : query.isError ? (
          <ErrorState onRetry={() => void query.refetch()} />
        ) : query.data.items.length === 0 ? (
          <EmptyState />
        ) : (
          <ul className="flex flex-col gap-2">
            {query.data.items.slice(0, RECENT_ORDERS_LIMIT).map((order) => (
              <li
                key={order.provider_order_id}
                className="flex items-center justify-between text-sm text-ink"
              >
                <span className="truncate">{order.buyer_nickname ?? order.provider_order_id}</span>
                <span className="font-semibold" style={{ fontFamily: "var(--font-mono)" }}>
                  {order.total !== null ? order.total.toLocaleString("pt-BR", { style: "currency", currency: order.currency ?? "BRL" }) : "—"}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}
