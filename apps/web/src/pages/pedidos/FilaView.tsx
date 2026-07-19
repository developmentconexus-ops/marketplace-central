import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";
import { formatDateTime, formatMoney } from "./pedidosFormatters";

export interface FilaViewProps {
  items: OrderRead[];
}

interface FilaAction {
  label: "Faturar" | "Etiqueta";
}

// Client-derived heuristic from real fields only (no fabricated bucket field on OrderRead):
// no nf_state yet -> needs invoicing; nf_state but no rastreio -> needs a shipping label;
// rastreio present -> already shipped, nothing left to action here.
function deriveAction(item: OrderRead): FilaAction | null {
  if (!item.nf_state) return { label: "Faturar" };
  if (!item.rastreio) return { label: "Etiqueta" };
  return null;
}

function sortByUrgency(items: OrderRead[]): OrderRead[] {
  return [...items].sort((a, b) => {
    const aAtrasado = a.sla?.atrasado === true;
    const bAtrasado = b.sla?.atrasado === true;
    if (aAtrasado !== bAtrasado) return aAtrasado ? -1 : 1;

    const aDue = a.sla?.due ? new Date(a.sla.due).getTime() : null;
    const bDue = b.sla?.due ? new Date(b.sla.due).getTime() : null;
    if (aDue !== bDue) {
      if (aDue === null) return 1;
      if (bDue === null) return -1;
      return aDue - bDue;
    }

    const aCreated = a.provider_created_at ? new Date(a.provider_created_at).getTime() : 0;
    const bCreated = b.provider_created_at ? new Date(b.provider_created_at).getTime() : 0;
    return bCreated - aCreated;
  });
}

function orderDescription(item: OrderRead): string {
  const buyer = item.buyer?.display ?? "comprador desconhecido";
  const itemsLabel =
    item.items.length > 1 ? `${item.items.length} itens` : (item.items[0]?.title ?? "item");
  const total = formatMoney(item.total);
  return total ? `${buyer} · ${itemsLabel} · ${total}` : `${buyer} · ${itemsLabel}`;
}

export function FilaView({ items }: FilaViewProps) {
  const ordered = sortByUrgency(items);

  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs text-faint">
        Fila de trabalho · ordenada por urgência (atrasados primeiro, depois SLA de envio) · DIFAL
        indisponível nesta fila
      </p>
      <div className="overflow-x-auto rounded-card border border-border bg-surface">
        <div style={{ minWidth: "820px" }} className="flex flex-col">
          {ordered.map((item, index) => {
            const action = deriveAction(item);
            const atrasado = item.sla?.atrasado === true;
            const due = formatDateTime(item.sla?.due);
            return (
              <div
                key={item.provider_order_id}
                className={`flex items-center gap-3 px-4 py-3 text-xs ${
                  index > 0 ? "border-t border-border-2" : ""
                } ${atrasado ? "bg-red-50" : ""}`}
              >
                <span
                  className={`w-24 flex-none text-[11px] font-bold uppercase tracking-wide ${
                    atrasado ? "text-red-700" : "text-faint"
                  }`}
                >
                  {atrasado ? "atrasado" : (due ?? <UnknownValue hint="sem SLA de envio" />)}
                </span>
                <span className="w-20 flex-none font-mono text-[11.5px] text-faint">
                  {item.provider_code || item.provider_order_id}
                </span>
                <span className="min-w-0 flex-1 truncate">{orderDescription(item)}</span>
                <span className="flex flex-none items-center gap-1 whitespace-nowrap font-mono font-semibold">
                  <UnknownValue hint="retorno depende de decomposição ainda não disponível (DIFAL/custo)" />
                </span>
                {action ? (
                  <button
                    type="button"
                    disabled
                    title="disponível em breve"
                    className="flex-none rounded-md border border-border bg-surface-2 px-3 py-1 text-[11.5px] font-semibold text-muted disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    {action.label}
                  </button>
                ) : (
                  <span className="flex-none text-[11.5px] text-faint">sem ação</span>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
