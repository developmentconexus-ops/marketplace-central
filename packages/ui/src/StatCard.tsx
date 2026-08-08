interface StatCardProps {
  label: string;
  value: string | number;
  sub?: string;
  className?: string;
}

export function StatCard({ label, value, sub, className = "" }: StatCardProps) {
  return (
    <div className={`bg-surface border border-border rounded-card p-5 ${className}`}>
      <p className="text-xs font-medium text-muted uppercase tracking-wide">{label}</p>
      <p
        className="mt-1 text-2xl font-semibold text-ink"
        style={{ fontFamily: "var(--font-mono)" }}
      >
        {value}
      </p>
      {sub && <p className="mt-0.5 text-xs text-faint">{sub}</p>}
    </div>
  );
}
