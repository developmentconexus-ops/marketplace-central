import type { OrderBucket, OrderRead } from "@marketplace-central/sdk-runtime";

export type PedidosTab = "novo" | "faturar" | "enviar" | "enviado" | "concluido" | "cancelado" | "devolucao";

export interface PedidosTabDefinition {
  value: PedidosTab;
  label: string;
  placeholder?: boolean;
}

export const pedidosTabs: PedidosTabDefinition[] = [
  { value: "novo", label: "Novos" },
  { value: "faturar", label: "A faturar" },
  { value: "enviar", label: "A enviar" },
  { value: "enviado", label: "Enviados" },
  { value: "concluido", label: "Concluídos", placeholder: true },
  // cancelado exists as an OrderRead.bucket value, but the demo dataset is REVIEW-capped (no
  // seller SKUs on the ML account) and there is no honest cause/reputation breakdown to show yet
  // — treated as a placeholder tab alongside concluido/devolucao rather than a live filter.
  { value: "cancelado", label: "Cancelados", placeholder: true },
  { value: "devolucao", label: "Devoluções", placeholder: true },
];

const liveBucketTabs = new Set<PedidosTab>(["novo", "faturar", "enviar", "enviado"]);

export function isLiveBucketTab(tab: PedidosTab): tab is OrderBucket & PedidosTab {
  return liveBucketTabs.has(tab);
}

export function filterOrdersByTab(items: OrderRead[], tab: PedidosTab): OrderRead[] {
  if (!isLiveBucketTab(tab)) return [];
  return items.filter((item) => item.bucket === tab);
}

// Counts derive from the already-loaded list's `bucket` field, so a tab's label always matches
// the rows it renders. The KPI cards (PedidosPage) derive from this same helper over the same
// list, so KPI == Lista by construction.
export function bucketTabCount(items: OrderRead[], tab: PedidosTab): number | null {
  if (!isLiveBucketTab(tab)) return null;
  return items.filter((item) => item.bucket === tab).length;
}
