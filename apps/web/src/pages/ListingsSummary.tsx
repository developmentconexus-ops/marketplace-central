import type { ListingSummary } from "@marketplace-central/sdk-runtime";
import {
  ErrorState,
  FreshnessIndicator,
  LoadingState,
  UnknownValue,
} from "@marketplace-central/ui";

export interface ListingsSummaryProps {
  isPending: boolean;
  isError: boolean;
  data: ListingSummary | undefined;
  onRetry: () => void;
}

interface SummaryCounterProps {
  label: string;
  value: number | null;
}

function SummaryCounter({ label, value }: SummaryCounterProps) {
  return (
    <div className="flex items-baseline gap-1.5">
      <dt className="text-xs text-faint">{label}</dt>
      <dd className="font-mono text-sm font-semibold text-ink">
        {value === null ? <UnknownValue /> : value}
      </dd>
    </div>
  );
}

// Demoted, counters-only summary strip. The clickable exception chips that used
// to live here were promoted into the page header (design/handoff/Anuncios.dc.html
// renders exceptions inline beside the title, not in a stat panel); this strip is
// now secondary context rendered below the table, never leading the page.
export function ListingsSummary({ isPending, isError, data, onRetry }: ListingsSummaryProps) {
  if (isPending) return <LoadingState />;
  if (isError) return <ErrorState onRetry={onRetry} />;
  if (!data) return null;

  return (
    <div className="flex flex-wrap items-center gap-x-5 gap-y-2">
      <dl className="flex flex-wrap items-center gap-x-5 gap-y-2">
        <SummaryCounter label="Total" value={data.total} />
        <SummaryCounter label="Ativos" value={data.active} />
        <SummaryCounter label="Pausados" value={data.paused} />
        <SummaryCounter label="Com erro de sync" value={data.exceptions.sync_error} />
        <SummaryCounter label="Desatualizados" value={data.exceptions.stale} />
        <SummaryCounter label="Sem vínculo" value={data.exceptions.unlinked} />
        <SummaryCounter label="Abaixo da margem" value={data.exceptions.below_margin_worst_case} />
        <SummaryCounter label="Margem desconhecida" value={data.exceptions.margin_unknown} />
      </dl>
      <p className="text-xs text-faint">
        Atualizado <FreshnessIndicator asOf={data.as_of} />
      </p>
    </div>
  );
}
