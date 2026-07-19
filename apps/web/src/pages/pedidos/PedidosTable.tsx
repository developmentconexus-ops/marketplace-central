import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { DataTable, UnknownValue, type DataTableColumn } from "@marketplace-central/ui";
import { actionLabelForBucket, formatDateTime, formatMoney } from "./pedidosFormatters";

export interface PedidosTableProps {
  items: OrderRead[];
  selectedKeys?: ReadonlySet<string>;
  onSelectionChange?: (next: Set<string>) => void;
  onRowClick?: (item: OrderRead) => void;
}

function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  );
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

function renderItens(item: OrderRead) {
  if (item.items.length === 0) return <UnknownValue />;
  if (item.items.length > 1) return `${item.items.length} itens`;
  return item.items[0]?.title ?? <UnknownValue />;
}

function renderAcao(item: OrderRead) {
  const label = actionLabelForBucket(item.bucket);
  if (!label) return <span className="text-xs text-faint">—</span>;
  return (
    <button
      type="button"
      disabled
      title="disponível em breve"
      className="flex-none rounded-md border border-border bg-surface-2 px-3 py-1 text-[11.5px] font-semibold text-muted disabled:cursor-not-allowed disabled:opacity-60"
    >
      {label}
    </button>
  );
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
    key: "data",
    header: "Data",
    render: (item) => formatDateTime(item.provider_created_at) ?? <UnknownValue />,
  },
  {
    key: "comprador",
    header: "Comprador",
    render: (item) => item.buyer?.display ?? <UnknownValue />,
  },
  {
    key: "itens",
    header: "Itens",
    render: renderItens,
  },
  {
    key: "valor",
    header: "Valor",
    render: (item) => formatMoney(item.total) ?? <UnknownValue />,
  },
  {
    key: "retorno",
    header: "Retorno",
    // Gated (slice C, IC-04 open): no retorno/margem field on OrderRead — honest unknown, never
    // fabricated (ADR-17).
    render: () => <UnknownValue hint="retorno depende de decomposição ainda não disponível" />,
  },
  {
    key: "sla",
    header: "SLA",
    render: renderSla,
  },
  {
    key: "difal",
    header: "DIFAL",
    // Gated (slice C): no difal field on OrderRead — honest unknown, never fabricated.
    render: () => <UnknownValue hint="DIFAL ainda não decomposto" />,
  },
  {
    key: "acao",
    header: "Ação",
    render: renderAcao,
  },
];

export function PedidosTable({ items, selectedKeys, onSelectionChange, onRowClick }: PedidosTableProps) {
  return (
    <div className="overflow-x-auto">
      <div style={{ minWidth: "980px" }}>
        <DataTable<OrderRead>
          columns={columns}
          rows={items}
          rowKey={(item) => item.provider_order_id}
          selectedKeys={selectedKeys}
          onSelectionChange={onSelectionChange}
          onRowClick={onRowClick}
        />
      </div>
    </div>
  );
}
