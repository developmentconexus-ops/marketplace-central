import type { ReactNode } from "react";

interface EmptyStateProps {
  hint?: ReactNode;
}

export function EmptyState({ hint }: EmptyStateProps) {
  return (
    <div className="flex flex-col gap-1 text-sm text-muted">
      <p>Nenhum registro encontrado.</p>
      {hint ? <p className="text-xs text-faint">{hint}</p> : null}
    </div>
  );
}
