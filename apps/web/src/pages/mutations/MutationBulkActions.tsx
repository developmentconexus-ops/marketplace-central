import type { MutationType } from "@marketplace-central/sdk-runtime";

interface MutationBulkActionsProps {
  selectedIds: Set<string>;
  onOpen: (type: MutationType) => void;
}

const actions: Array<{ label: string; type: MutationType }> = [
  { label: "Atualizar preço", type: "price_update" },
  { label: "Corrigir estoque", type: "stock_correct" },
  { label: "Pausar", type: "listing_pause" },
  { label: "Ressincronizar", type: "listing_resync" },
  { label: "Vincular", type: "link_apply" },
  { label: "Editar", type: "listing_edit" },
];

export function MutationBulkActions({ selectedIds, onOpen }: MutationBulkActionsProps) {
  const disabled = selectedIds.size === 0;

  return (
    <>
      {actions.map((action) => (
        <button
          key={action.type}
          type="button"
          disabled={disabled}
          title={disabled ? "Selecione ao menos um anúncio" : undefined}
          onClick={() => onOpen(action.type)}
          className="rounded-lg border px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-50"
        >
          {action.label}
        </button>
      ))}
    </>
  );
}
