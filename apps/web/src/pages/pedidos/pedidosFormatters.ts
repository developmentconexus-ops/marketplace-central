import type { OrderBucket, OrderRead } from "@marketplace-central/sdk-runtime";

const currencyFormatter = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL",
});

// Shared by PedidosTable (AÇÃO column) and KanbanView (card action) so both render the same
// disabled label for a given bucket instead of forking the mapping.
export function actionLabelForBucket(bucket: OrderBucket): "Faturar" | "Etiqueta" | null {
  if (bucket === "novo" || bucket === "faturar") return "Faturar";
  if (bucket === "enviar") return "Etiqueta";
  return null;
}

// Margin band → pct-chip token classes. margem_pct is a ratio (0.18 = 18%), matching the design
// legend: verde ≥18% · âmbar 10–18% · vermelho <10%. No --err token in app @theme, so the <10%
// band uses warn (reddish-brown) rather than a hardcoded red.
export function marginBandClass(ratio: number): string {
  if (ratio >= 0.18) return "bg-accent-soft text-accent-ink";
  if (ratio >= 0.1) return "bg-amber-soft text-amber";
  return "bg-warn-soft text-warn";
}

// Fila urgency tier: the leftmost tag of a work-queue row. Real-payload only (sla.due + atrasado,
// then bucket) — never a fabricated relative label. Returns the display text + token color class.
export function frontTier(item: OrderRead): { text: string; className: string } {
  if (item.sla?.atrasado === true) {
    return { text: "ATRASADO", className: "text-warn" };
  }
  const due = formatDate(item.sla?.due);
  if (due) {
    return { text: due, className: "text-muted" };
  }
  const byBucket: Record<OrderBucket, string> = {
    novo: "NOVO",
    faturar: "A FATURAR",
    enviar: "A ENVIAR",
    enviado: "ENVIADO",
    cancelado: "CANCELADO",
  };
  return { text: byBucket[item.bucket], className: "text-faint" };
}

// Buyer/recipient name for a row's description. Prefers destinatario (the shipment receiver, the
// real name we hold) then the masked buyer.display; null when neither is present (ADR-17 — never a
// fabricated name). City/UF is intentionally NOT here: the design (Pedidos.dc.html) keeps the
// Fila/Kanban description to comprador · itens · valor and shows destino only in the drawer.
export function orderComprador(item: OrderRead): string | null {
  return item.destinatario ?? item.buyer?.display ?? null;
}

// Item summary for a row's description, mirroring the design's descDe: a single item shows its
// title, many collapse to "N itens". Empty items → null (honest unknown, not "0 itens").
export function orderItensLabel(item: OrderRead): string | null {
  if (item.items.length === 0) return null;
  if (item.items.length > 1) return `${item.items.length} itens`;
  return item.items[0]?.title ?? null;
}

// Fila row description = comprador · itens · valor (design descDe). Absent segments drop so an
// order with real total but no item titles reads "Marcia Rocha · R$ 339,98" rather than injecting
// a fabricated placeholder mid-string; all-absent → null for an honest UnknownValue (ADR-17).
export function orderFilaDesc(item: OrderRead): string | null {
  const parts = [orderComprador(item), orderItensLabel(item), formatMoney(item.total)].filter(
    Boolean,
  );
  return parts.length > 0 ? parts.join(" · ") : null;
}

export function formatMoney(value: number | null | undefined): string | null {
  if (value === null || value === undefined) return null;
  return currencyFormatter.format(value);
}

// The ML dispatch SLA is a deadline DAY — its wire value is midnight local, so
// rendering it with a time printed "23/03/2026, 00:00" on every row and read as
// a real hour the seller had to hit. Date only says exactly what ML gave us.
export function formatDate(value: string | null | undefined): string | null {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;
  return date.toLocaleDateString("pt-BR", { dateStyle: "short" });
}

export function formatDateTime(value: string | null | undefined): string | null {
  if (!value) return null;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;
  return date.toLocaleString("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  });
}

// margem_pct is a ratio (retorno/preço), not a 0-100 percent-point number on the wire — IC-04
// (.mnfs/MIS-004-mvp-demo/research/pricing-difal-interface-contract.md:27), the single shared
// decomposition formula for Simulador AND Pedidos, states `margem_pct = retorno/preço`
// (dimensionless, e.g. 0.18). That is consistent with packages/ui/src/MarginChip.tsx, which
// expects its `marginPct` prop already as a percent-point number (18 → "18%") — so callers
// multiply the wire ratio by 100 before display. Multiplying here matches both: 0.18 → "18.0%".
// Never fabricate a value for null (ADR-17).
export function formatPercent(value: number | null | undefined): string | null {
  if (value === null || value === undefined) return null;
  // pt-BR decimal separator: the operator reads "18,0%", never "18.0%".
  return `${(value * 100).toLocaleString("pt-BR", { minimumFractionDigits: 1, maximumFractionDigits: 1 })}%`;
}
