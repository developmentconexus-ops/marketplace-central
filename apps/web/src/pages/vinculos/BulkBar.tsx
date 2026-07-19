export interface BulkBarProps {
  selectedCount: number;
  onPreview: () => void;
  onClear: () => void;
}

/** Local multi-select summary bar over the queue rows (MutationBulkActions precedent). */
export function BulkBar({ selectedCount, onPreview, onClear }: BulkBarProps) {
  if (selectedCount === 0) return null;

  return (
    <div
      role="toolbar"
      aria-label="Ações em lote"
      className="flex flex-wrap items-center justify-between gap-2 rounded-control border border-accent/30 bg-accent-soft px-3 py-2 text-sm"
    >
      <span className="font-medium text-accent-ink">{selectedCount} selecionado(s)</span>
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onClear}
          className="rounded-control border border-border bg-surface px-3 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
        >
          Limpar seleção
        </button>
        <button
          type="button"
          onClick={onPreview}
          className="rounded-control bg-accent px-3 py-1.5 text-xs font-medium text-accent-ink hover:bg-accent/90"
        >
          Pré-visualizar em lote
        </button>
      </div>
    </div>
  );
}
