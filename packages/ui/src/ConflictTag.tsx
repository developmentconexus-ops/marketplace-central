interface ConflictTagProps {
  detail?: string;
}

export function ConflictTag({ detail }: ConflictTagProps) {
  return (
    <span
      className="inline-flex items-center rounded-pill bg-amber-soft px-2 py-0.5 text-xs font-medium text-amber"
      title={detail || undefined}
    >
      divergente
    </span>
  );
}
