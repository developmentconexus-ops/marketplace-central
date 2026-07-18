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
      className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-blue-200 bg-blue-50 px-3 py-2 text-sm"
    >
      <span className="font-medium text-blue-900">{selectedCount} selecionado(s)</span>
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onClear}
          className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-50"
        >
          Limpar seleção
        </button>
        <button
          type="button"
          onClick={onPreview}
          className="rounded-lg border border-blue-600 bg-blue-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-blue-700"
        >
          Pré-visualizar em lote
        </button>
      </div>
    </div>
  );
}
