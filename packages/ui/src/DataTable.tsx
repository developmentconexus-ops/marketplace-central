import { useRef, type ReactNode } from "react";

export interface DataTableColumn<T> {
  key: string;
  header: ReactNode;
  render: (row: T) => ReactNode;
  align?: "left" | "right";
  sortable?: boolean;
}

export interface DataTableProps<T> {
  columns: DataTableColumn<T>[];
  rows: T[];
  rowKey: (row: T) => string;
  onRowClick?: (row: T) => void;
  selectedKeys?: ReadonlySet<string>;
  onSelectionChange?: (next: Set<string>) => void;
  sortKey?: string | null;
  sortDir?: "asc" | "desc";
  onSortChange?: (key: string) => void;
  loading?: boolean;
  emptyState?: ReactNode;
  stickyHeader?: boolean;
}

export function DataTable<T>({
  columns,
  rows,
  rowKey,
  onRowClick,
  selectedKeys,
  onSelectionChange,
  sortKey,
  sortDir,
  onSortChange,
  loading = false,
  emptyState,
  stickyHeader = false,
}: DataTableProps<T>): JSX.Element {
  const headerCheckboxRef = useRef<HTMLInputElement | null>(null);
  const selectionEnabled = onSelectionChange !== undefined;
  const keys = rows.map(rowKey);
  const selectedCount = keys.filter((key) => selectedKeys?.has(key)).length;
  const allSelected = keys.length > 0 && selectedCount === keys.length;
  const someSelected = selectedCount > 0 && !allSelected;

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12 text-sm text-muted">
        <div className="animate-spin rounded-full h-5 w-5 border-2 border-border border-t-accent" role="status" aria-label="Carregando" />
      </div>
    );
  }

  if (rows.length === 0) {
    return <div className="flex items-center justify-center py-16 text-sm text-faint">{emptyState ?? "Nenhum item."}</div>;
  }

  function handleSelectAll() {
    if (!onSelectionChange) return;
    onSelectionChange(allSelected ? new Set() : new Set(keys));
  }

  function handleRowSelection(key: string) {
    if (!onSelectionChange) return;
    const next = new Set(selectedKeys ?? []);
    if (next.has(key)) {
      next.delete(key);
    } else {
      next.add(key);
    }
    onSelectionChange(next);
  }

  return (
    <div className="overflow-x-auto border border-border rounded-card">
      <table className="w-full text-sm text-left">
        <thead className={`bg-surface-2 border-b border-border ${stickyHeader ? "sticky top-0" : ""}`}>
          <tr>
            {selectionEnabled && (
              <th className="px-3 py-2 text-xs font-medium text-muted uppercase tracking-wide">
                <input
                  ref={(element) => {
                    headerCheckboxRef.current = element;
                    if (element) element.indeterminate = someSelected;
                  }}
                  type="checkbox"
                  checked={allSelected}
                  onChange={handleSelectAll}
                  onClick={(event) => event.stopPropagation()}
                  className="rounded border-border"
                  aria-label="Select all rows"
                />
              </th>
            )}
            {columns.map((column) => {
              const alignClass = column.align === "right" ? " text-right" : "";
              const isActive = sortKey === column.key;

              return (
                <th key={column.key} className={`px-3 py-2 text-xs font-medium text-muted uppercase tracking-wide${alignClass}`}>
                  {column.sortable ? (
                    <button
                      type="button"
                      className="inline-flex items-center gap-1 cursor-pointer"
                      onClick={() => onSortChange?.(column.key)}
                    >
                      {column.header}
                      <span className={isActive ? "text-ink" : "text-faint"} aria-hidden="true">
                        {isActive ? (sortDir === "asc" ? "▲" : "▼") : "↕"}
                      </span>
                    </button>
                  ) : (
                    column.header
                  )}
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => {
            const key = rowKey(row);
            const selected = selectedKeys?.has(key) ?? false;

            return (
              <tr
                key={key}
                className={`border-b border-border last:border-0 ${onRowClick ? "cursor-pointer hover:bg-surface-2" : ""} ${selected ? "bg-accent-soft" : ""}`}
                onClick={onRowClick ? () => onRowClick(row) : undefined}
              >
                {selectionEnabled && (
                  <td className="px-3 py-2 text-sm text-ink">
                    <input
                      type="checkbox"
                      checked={selected}
                      onChange={() => handleRowSelection(key)}
                      onClick={(event) => event.stopPropagation()}
                      className="rounded border-border"
                      aria-label={`Select row ${key}`}
                    />
                  </td>
                )}
                {columns.map((column) => (
                  <td
                    key={column.key}
                    className={`px-3 py-2 text-sm text-ink${column.align === "right" ? " text-right" : ""}`}
                  >
                    {column.render(row)}
                  </td>
                ))}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
