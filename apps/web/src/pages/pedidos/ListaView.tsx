import { useEffect, useState } from "react";
import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { EmptyState } from "@marketplace-central/ui";
import { PedidosTable } from "./PedidosTable";
import { bucketTabCount, filterOrdersByTab, pedidosTabs, type PedidosTab } from "./pedidosTabs";

export interface ListaViewProps {
  items: OrderRead[];
  tab: PedidosTab;
  onTabChange: (tab: PedidosTab) => void;
  onOpenOrder: (orderId: string) => void;
}

const placeholderCopy: Record<string, string> = {
  concluido: "Concluídos — em breve",
  cancelado:
    "Cancelados — em breve: causas (estoque, prazo, arrependimento) e impacto na reputação",
  devolucao: "Devoluções — em breve: fluxo de tratativa e reintegração de estoque",
};

export function ListaView({ items, tab, onTabChange, onOpenOrder }: ListaViewProps) {
  const [selected, setSelected] = useState<Set<string>>(new Set());

  // Selection is scoped to a single tab's rows; switching tabs invalidates the row keys shown.
  useEffect(() => {
    setSelected(new Set());
  }, [tab]);

  const activeDef = pedidosTabs.find((tabItem) => tabItem.value === tab);
  const isPlaceholder = activeDef?.placeholder === true;
  const visibleItems = filterOrdersByTab(items, tab);

  return (
    <div className="flex flex-col gap-3">
      <div aria-label="Filtros de pedidos" role="tablist" className="flex flex-wrap gap-1">
        {pedidosTabs.map((tabItem) => {
          const count = bucketTabCount(items, tabItem.value);
          const label = count !== null ? `${tabItem.label} ${count}` : tabItem.label;
          const active = tab === tabItem.value;
          return (
            <button
              key={tabItem.value}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => onTabChange(tabItem.value)}
              className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                active ? "bg-accent-soft text-accent-ink" : "text-muted hover:text-ink"
              }`}
            >
              {label}
            </button>
          );
        })}
      </div>

      {isPlaceholder ? (
        <div className="rounded-card border border-dashed border-border bg-surface-2 p-8 text-center text-sm text-muted">
          {placeholderCopy[tab] ?? `${activeDef?.label} — em breve`}
        </div>
      ) : visibleItems.length === 0 ? (
        <EmptyState />
      ) : (
        <>
          <PedidosTable
            items={visibleItems}
            selectedKeys={selected}
            onSelectionChange={setSelected}
            onRowClick={(item) => onOpenOrder(item.provider_order_id)}
          />
          {selected.size > 0 ? (
            <div className="flex items-center gap-3 rounded-card border border-border bg-surface px-4 py-2 text-sm text-muted">
              <b className="text-ink">{selected.size} selecionado(s)</b>
              <button
                type="button"
                disabled
                title="disponível em breve"
                className="rounded-md border border-border px-3 py-1 text-xs font-medium text-muted disabled:cursor-not-allowed disabled:opacity-60"
              >
                Etiquetas
              </button>
              <button
                type="button"
                disabled
                title="disponível em breve"
                className="rounded-md border border-border px-3 py-1 text-xs font-medium text-muted disabled:cursor-not-allowed disabled:opacity-60"
              >
                Marcar enviados
              </button>
              <button
                type="button"
                onClick={() => setSelected(new Set())}
                className="ml-auto text-xs text-faint hover:text-ink"
              >
                limpar
              </button>
            </div>
          ) : null}
        </>
      )}
    </div>
  );
}
