import type { OrderBucket, OrderRead } from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";
import {
  actionLabelForBucket,
  formatMoney,
  orderComprador,
  orderItensLabel,
} from "./pedidosFormatters";

export interface KanbanViewProps {
  items: OrderRead[];
  onOpenOrder: (orderId: string) => void;
}

const columnDefs: { bucket: OrderBucket; title: string }[] = [
  { bucket: "novo", title: "NOVOS" },
  { bucket: "faturar", title: "A FATURAR" },
  { bucket: "enviar", title: "A ENVIAR" },
  { bucket: "enviado", title: "ENVIADOS" },
];

function KanbanCard({
  item,
  onOpenOrder,
}: {
  item: OrderRead;
  onOpenOrder: (orderId: string) => void;
}) {
  const label = actionLabelForBucket(item.bucket);
  const money = formatMoney(item.total);
  // Card description mirrors the design (comprador · itens); valor sits in the row above, destino
  // lives in the drawer. Both parts null → honest UnknownValue (ADR-17).
  const desc = [orderComprador(item), orderItensLabel(item)].filter(Boolean).join(" · ") || null;
  const atrasado = item.sla?.atrasado === true;
  return (
    <div
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
      className={`flex cursor-pointer flex-col gap-1 rounded-lg border-[1.5px] bg-surface p-3 text-xs hover:border-accent ${
        atrasado ? "border-warn" : "border-border"
      }`}
    >
      <div className="flex items-center gap-2">
        <span className="font-mono text-[11px] text-faint" title={item.provider_order_id}>
          {item.provider_order_id}
        </span>
        <span className="ml-auto font-mono font-semibold">{money ?? <UnknownValue />}</span>
      </div>
      <span className="truncate text-muted">
        {desc ?? <UnknownValue hint="comprador/itens ainda não disponíveis" />}
      </span>
      <span className={`text-[11px] ${atrasado ? "font-semibold text-warn" : "text-faint"}`}>
        {atrasado ? "SLA atrasado" : formatSlaLabel(item)}
      </span>
      {label ? (
        <button
          type="button"
          disabled
          title="disponível em breve"
          className="rounded-md border border-border bg-surface-2 px-2 py-1 text-center text-[11px] font-semibold text-muted disabled:cursor-not-allowed disabled:opacity-60"
        >
          {label}
        </button>
      ) : null}
    </div>
  );
}

function formatSlaLabel(item: OrderRead): string {
  return item.sla?.due ? "com prazo de envio" : "sem SLA de envio";
}

export function KanbanView({ items, onOpenOrder }: KanbanViewProps) {
  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs text-faint">
        visão por etapa · o status vem do ML/ERP · sem arrastar · ações nos cards em breve
      </p>
      <div className="flex items-start gap-3 overflow-x-auto pb-1">
        {columnDefs.map((column) => {
          const cards = items.filter((item) => item.bucket === column.bucket);
          return (
            <div
              key={column.bucket}
              style={{ minWidth: "258px" }}
              className="flex flex-none flex-col gap-2"
            >
              <div className="text-[11px] font-bold tracking-wide text-faint">
                {column.title} · {cards.length}
              </div>
              <div className="flex flex-col gap-2">
                {cards.map((item) => (
                  <KanbanCard key={item.provider_order_id} item={item} onOpenOrder={onOpenOrder} />
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
