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
