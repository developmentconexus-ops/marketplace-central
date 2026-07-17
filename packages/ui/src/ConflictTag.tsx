interface ConflictTagProps {
  detail?: string;
}

export function ConflictTag({ detail }: ConflictTagProps) {
  return (
    <span
      className="inline-flex items-center rounded bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-800"
      title={detail || undefined}
    >
      divergente
    </span>
  );
}
