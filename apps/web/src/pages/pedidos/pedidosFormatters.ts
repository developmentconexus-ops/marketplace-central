import type { OrderBucket } from "@marketplace-central/sdk-runtime";

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

export function formatMoney(value: number | null | undefined): string | null {
  if (value === null || value === undefined) return null;
  return currencyFormatter.format(value);
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
  return `${(value * 100).toFixed(1)}%`;
}
