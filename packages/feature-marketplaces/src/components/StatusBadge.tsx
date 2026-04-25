interface StatusBadgeProps {
  status: string;
}

const STYLES: Record<string, { badge: string; dot: string }> = {
  active:           { badge: "bg-emerald-100 text-emerald-700", dot: "bg-emerald-500" },
  connected:        { badge: "bg-emerald-100 text-emerald-700", dot: "bg-emerald-500" },
  available:        { badge: "bg-emerald-100 text-emerald-700", dot: "bg-emerald-500" },
  inactive:         { badge: "bg-slate-100 text-slate-500", dot: "bg-slate-400" },
  disconnected:     { badge: "bg-slate-100 text-slate-500", dot: "bg-slate-400" },
  pending_connection:{ badge: "bg-amber-100 text-amber-700", dot: "bg-amber-500" },
  requires_reauth:  { badge: "bg-amber-100 text-amber-700", dot: "bg-amber-500" },
  degraded:         { badge: "bg-orange-100 text-orange-700", dot: "bg-orange-500" },
  warning:          { badge: "bg-orange-100 text-orange-700", dot: "bg-orange-500" },
  failed:           { badge: "bg-red-100 text-red-700", dot: "bg-red-500" },
  blocked:          { badge: "bg-red-100 text-red-700", dot: "bg-red-500" },
  unavailable:      { badge: "bg-red-100 text-red-700", dot: "bg-red-500" },
  planned:          { badge: "bg-indigo-100 text-indigo-700", dot: "bg-indigo-500" },
};

const DEFAULT_STYLE = { badge: "bg-slate-100 text-slate-500", dot: "bg-slate-400" };

export function StatusBadge({ status }: StatusBadgeProps) {
  const { badge, dot } = STYLES[status] ?? DEFAULT_STYLE;
  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium ${badge}`}
      aria-label={`Status: ${status}`}
    >
      <span className={`w-1.5 h-1.5 rounded-full ${dot}`} />
      {status}
    </span>
  );
}
