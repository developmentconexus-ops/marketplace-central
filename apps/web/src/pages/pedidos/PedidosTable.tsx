import type { OrderRead } from "@marketplace-central/sdk-runtime";
import { DataTable, UnknownValue, type DataTableColumn } from "@marketplace-central/ui";
import {
  actionLabelForBucket,
  formatDateTime,
  formatMoney,
  formatPercent,
  marginBandClass,
} from "./pedidosFormatters";

export interface PedidosTableProps {
  items: OrderRead[];
  selectedKeys?: ReadonlySet<string>;
  onSelectionChange?: (next: Set<string>) => void;
  onRowClick?: (item: OrderRead) => void;
}

function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return (
    <span
      className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}
    >
      {label}
    </span>
  );
}

function renderSla(item: OrderRead) {
  const due = formatDateTime(item.sla?.due);
  const atrasado = item.sla?.atrasado === true;
  return (
    <div className="flex flex-col gap-1">
      <span className={atrasado ? "font-semibold text-warn" : undefined}>
        {due ?? <UnknownValue />}
      </span>
      {atrasado ? stateTag("atrasado", "bg-warn-soft text-warn") : null}
    </div>
  );
}

function renderItens(item: OrderRead) {
  if (item.items.length === 0) return <UnknownValue />;
  if (item.items.length > 1) return `${item.items.length} itens`;
  return item.items[0]?.title ?? <UnknownValue />;
}

// Real-ready wiring (F02-S6): reads order.retorno_liquido/margem_pct — null (honest-empty, hub
// C2 not wired yet) renders UnknownValue, a real value renders formatted with no further UI
// change once the decomposer lands.
function renderRetorno(item: OrderRead) {
  const retorno = formatMoney(item.retorno_liquido);
  const ratio = item.margem_pct;
  const margem = formatPercent(ratio);
  return (
    <div className="flex items-center gap-1.5 whitespace-nowrap">
      <span className="font-mono">
        {retorno ?? <UnknownValue hint="decomposição de custos ainda não disponível (hub C2)" />}
      </span>
      {margem != null && ratio != null ? (
        <span
          className={`rounded-pill px-1.5 py-0.5 text-[10.5px] font-semibold ${marginBandClass(ratio)}`}
        >
          {margem}
        </span>
      ) : null}
    </div>
  );
}

// Real-ready wiring (F02-S6): reads order.difal — null amount renders UnknownValue honestly; a
// real value renders a token-colored chip (pago ✓ / agendado / valor) mirroring the design.
function renderDifal(item: OrderRead) {
  const difal = item.difal;
  if (!difal || difal.amount == null) {
    return <UnknownValue hint="DIFAL ainda não decomposto (hub C2)" />;
  }
  if (difal.paid === true) {
    return stateTag("pago ✓", "bg-accent-soft text-accent-ink");
  }
  const due = formatDateTime(difal.due_date);
  const amount = formatMoney(difal.amount);
  const label = due ? `ag. ${due}` : (amount ?? "agendado");
  return stateTag(label, "bg-info-soft text-info");
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
    // The order number is provider_order_id; provider_code is the channel slug ("mercado_livre")
    // for this connector, so it is not shown as the order label (design parity — the Pedido cell is
    // the order number).
    render: (item) => (
      <span className="font-mono text-[11.5px] text-faint">{item.provider_order_id}</span>
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
    // null until the decomposer is wired (hub C2); renders '—' honestly, real value when present.
    render: renderRetorno,
  },
  {
    key: "sla",
    header: "SLA",
    render: renderSla,
  },
  {
    key: "difal",
    header: "DIFAL",
    // null until the decomposer is wired (hub C2); renders '—' honestly, real value when present.
    render: renderDifal,
  },
  {
    key: "acao",
    header: "Ação",
    render: renderAcao,
  },
];

export function PedidosTable({
  items,
  selectedKeys,
  onSelectionChange,
  onRowClick,
}: PedidosTableProps) {
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
