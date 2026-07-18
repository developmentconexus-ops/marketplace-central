import type { OrderRead } from "@marketplace-central/sdk-runtime";

export type PedidosTab = "todos" | "atrasados" | "sem_vinculo" | "cancelados";

export interface PedidosTabDefinition {
  value: PedidosTab;
  label: string;
  disabled?: boolean;
}

export const pedidosTabs: PedidosTabDefinition[] = [
  { value: "todos", label: "Todos" },
  { value: "atrasados", label: "Atrasados" },
  { value: "sem_vinculo", label: "Sem vínculo" },
  { value: "cancelados", label: "Cancelados", disabled: true },
];

export function filterOrdersByTab(items: OrderRead[], tab: PedidosTab): OrderRead[] {
  switch (tab) {
    case "atrasados":
      return items.filter((item) => item.sla?.atrasado === true);
    case "sem_vinculo":
      return items.filter((item) => item.vinculo_status === "SEM_VINCULO");
    case "cancelados":
      return [];
    case "todos":
    default:
      return items;
  }
}
