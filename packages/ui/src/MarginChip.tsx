interface MarginChipProps {
  marginPct: number | null;
  thresholds?: { healthy: number; tight: number };
}

const DEFAULT_THRESHOLDS = { healthy: 18, tight: 10 };

type MarginBand = "healthy" | "tight" | "warn" | "unknown";

const bandClasses: Record<MarginBand, string> = {
  healthy: "bg-accent-soft text-accent-ink",
  tight: "bg-amber-soft text-amber",
  warn: "bg-warn-soft text-warn",
  unknown: "bg-surface-2 text-muted",
};

function hasValidThresholds(thresholds: { healthy: number; tight: number }): boolean {
  return Number.isFinite(thresholds.healthy) && Number.isFinite(thresholds.tight) && thresholds.healthy > thresholds.tight;
}

export function MarginChip({ marginPct, thresholds }: MarginChipProps): JSX.Element {
  const validThresholds = thresholds === undefined || hasValidThresholds(thresholds);

  if (thresholds !== undefined && !validThresholds && import.meta.env.DEV) {
    console.warn("MarginChip received invalid thresholds; using defaults.");
  }

  const resolvedThresholds = validThresholds && thresholds !== undefined ? thresholds : DEFAULT_THRESHOLDS;
  let band: MarginBand = "unknown";

  if (marginPct !== null && Number.isFinite(marginPct)) {
    band = marginPct >= resolvedThresholds.healthy
      ? "healthy"
      : marginPct >= resolvedThresholds.tight
        ? "tight"
        : "warn";
  }

  return (
    <span className={`inline-flex items-center rounded-pill px-2 py-0.5 text-xs font-mono font-medium ${bandClasses[band]}`}>
      {band === "unknown" ? "—" : `${marginPct}%`}
    </span>
  );
}
