import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { DataTable, UnknownValue, type DataTableColumn } from "@marketplace-central/ui";
import { formatDateTime } from "./pedidosFormatters";

export interface PedidosTableProps {
  items: OrderRead[];
}

function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  );
}

function renderVinculo(status: string | undefined) {
  if (!status) return <UnknownValue />;
  if (status === "SEM_VINCULO") return stateTag("sem vínculo", "bg-amber-100 text-amber-800");
  return stateTag(status.toLowerCase(), "bg-emerald-100 text-emerald-800");
}

function renderSla(item: OrderRead) {
  const due = formatDateTime(item.sla?.due);
  return (
    <div className="flex flex-col gap-1">
      <span>{due ?? <UnknownValue />}</span>
      {item.sla?.atrasado ? stateTag("atrasado", "bg-red-100 text-red-800") : null}
    </div>
  );
}

function renderCusto(item: OrderRead) {
  if (item.componentes_desconhecidos && item.componentes_desconhecidos.length > 0) {
    return stateTag("custo incompleto", "bg-amber-100 text-amber-800");
  }
  return null;
}

const columns: DataTableColumn<OrderRead>[] = [
  {
    key: "pedido",
    header: "Pedido",
    render: (item) => (
      <div className="flex flex-col">
        <span className="font-medium">{item.provider_code || item.provider_order_id}</span>
        <span className="text-xs text-faint">{item.provider_order_id}</span>
      </div>
    ),
  },
  {
    key: "status",
    header: "Status",
    render: (item) => stateTag(item.status),
  },
  {
    key: "comprador",
    header: "Comprador",
    render: (item) => item.buyer?.display ?? <UnknownValue />,
  },
  {
    key: "destino",
    header: "Destino",
    render: (item) => item.destino_uf ?? <UnknownValue />,
  },
  {
    key: "vinculo",
    header: "Vínculo",
    render: (item) => renderVinculo(item.vinculo_status),
  },
  {
    key: "sla",
    header: "SLA",
    render: renderSla,
  },
  {
    key: "custo",
    header: "Custo",
    render: renderCusto,
  },
];

export function PedidosTable({ items }: PedidosTableProps) {
  return <DataTable<OrderRead> columns={columns} rows={items} rowKey={(item) => item.provider_order_id} />;
}
