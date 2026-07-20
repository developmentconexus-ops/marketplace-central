import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";
import {
  actionLabelForBucket,
  formatMoney,
  formatPercent,
  frontTier,
  marginBandClass,
  orderFilaDesc,
} from "./pedidosFormatters";

export interface FilaViewProps {
  items: OrderRead[];
  onOpenOrder: (orderId: string) => void;
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

// Retorno + margem-pct chip (design's retDe group): mono value + a token-banded pill. Both derive
// from the order's own retorno_liquido/margem_pct — null (decomposer not wired, hub C2) renders an
// honest UnknownValue with no pill, never a fabricated 0 or band (ADR-17).
function FilaRetorno({ item }: { item: OrderRead }) {
  const retorno = formatMoney(item.retorno_liquido);
  const ratio = item.margem_pct;
  const margem = formatPercent(ratio);
  return (
    <span className="flex flex-none items-center gap-[5px] whitespace-nowrap">
      <span className="font-mono font-semibold">
        {retorno ?? (
          <UnknownValue hint="retorno depende de decomposição ainda não disponível (DIFAL/custo)" />
        )}
      </span>
      {margem != null && ratio != null ? (
        <span className={`rounded-pill px-[7px] text-[10.5px] font-bold ${marginBandClass(ratio)}`}>
          {margem}
        </span>
      ) : null}
    </span>
  );
}

export function FilaView({ items, onOpenOrder }: FilaViewProps) {
  const ordered = sortByUrgency(items);

  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs text-faint">
        fila de trabalho · ordenada por urgência (atrasados → SLA de envio → data) · clique abre o
        detalhe
      </p>
      <div className="overflow-x-auto rounded-card border border-border bg-surface">
        <div style={{ minWidth: "820px" }} className="flex flex-col">
          {ordered.map((item) => {
            const actionLabel = actionLabelForBucket(item.bucket);
            const atrasado = item.sla?.atrasado === true;
            const tier = frontTier(item);
            const desc = orderFilaDesc(item);
            return (
              <div
                key={item.provider_order_id}
                role="button"
                tabIndex={0}
                aria-label={`Abrir detalhe do pedido ${item.provider_order_id}`}
                onClick={() => onOpenOrder(item.provider_order_id)}
                onKeyDown={(event) => {
                  if (event.key === "Enter" || event.key === " ") {
                    event.preventDefault();
                    onOpenOrder(item.provider_order_id);
                  }
                }}
                className={`flex cursor-pointer items-center gap-3 border-t border-border-2 px-4 py-[11px] text-[12.5px] hover:bg-surface-2 ${
                  atrasado ? "bg-warn-soft" : ""
                }`}
              >
                <span className={`w-[110px] flex-none truncate text-[11px] font-bold ${tier.className}`}>
                  {tier.text}
                </span>
                <span
                  className="w-[78px] flex-none truncate font-mono text-[11.5px] text-faint"
                  title={item.provider_order_id}
                >
                  {item.provider_order_id}
                </span>
                <span className="min-w-0 flex-1 truncate">
                  {desc ?? <UnknownValue hint="comprador/itens/valor ainda não disponíveis" />}
                </span>
                <FilaRetorno item={item} />
                {actionLabel ? (
                  <button
                    type="button"
                    disabled
                    title="disponível em breve"
                    className="flex-none rounded-[7px] border border-border bg-surface-2 px-3 py-[5px] text-[11.5px] font-semibold text-muted disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    {actionLabel}
                  </button>
                ) : (
                  <span className="flex-none whitespace-nowrap text-[11.5px] text-faint">sem ação</span>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
