import type { ListingSummary } from "@marketplace-central/sdk-runtime";
import { ErrorState, FreshnessIndicator, LoadingState, UnknownValue } from "@marketplace-central/ui";

type SummaryExceptionChipKey = "sem_vinculo" | "abaixo_custo" | "sem_evidencia";

const exceptionChipLabels: Record<SummaryExceptionChipKey, string> = {
  sem_vinculo: "Sem vínculo",
  abaixo_custo: "Abaixo do custo",
  sem_evidencia: "Sem evidência",
};

const exceptionChipKeys: SummaryExceptionChipKey[] = ["sem_vinculo", "abaixo_custo", "sem_evidencia"];

export interface ListingsSummaryProps {
  isPending: boolean;
  isError: boolean;
  data: ListingSummary | undefined;
  onRetry: () => void;
  onExceptionChipClick: (exception: SummaryExceptionChipKey) => void;
  activeException?: string;
}

interface SummaryCounterProps {
  label: string;
  value: number | null;
}

function SummaryCounter({ label, value }: SummaryCounterProps) {
  return (
    <div className="border-l border-slate-200 pl-3 first:border-l-0 first:pl-0">
      <dt className="text-xs text-slate-500">{label}</dt>
      <dd className="mt-1 text-lg font-semibold text-slate-950">
        {value === null ? <UnknownValue /> : value}
      </dd>
    </div>
  );
}

export function ListingsSummary({
  isPending,
  isError,
  data,
  onRetry,
  onExceptionChipClick,
  activeException,
}: ListingsSummaryProps) {
  if (isPending) return <LoadingState />;
  if (isError) return <ErrorState onRetry={onRetry} />;
  if (!data) return null;

  const chips = exceptionChipKeys
    .map((key) => ({ key, count: data.exceptions[key] }))
    .filter((chip): chip is { key: SummaryExceptionChipKey; count: number } => Boolean(chip.count));

  return (
    <div className="mt-3">
      <dl className="grid gap-x-5 gap-y-4 sm:grid-cols-4 lg:grid-cols-8">
        <SummaryCounter label="Total" value={data.total} />
        <SummaryCounter label="Ativos" value={data.active} />
        <SummaryCounter label="Pausados" value={data.paused} />
        <SummaryCounter label="Com erro de sync" value={data.exceptions.sync_error} />
        <SummaryCounter label="Desatualizados" value={data.exceptions.stale} />
        <SummaryCounter label="Sem vínculo" value={data.exceptions.unlinked} />
        <SummaryCounter label="Abaixo da margem" value={data.exceptions.below_margin_worst_case} />
        <SummaryCounter label="Margem desconhecida" value={data.exceptions.margin_unknown} />
      </dl>
      <p className="mt-4 text-xs text-slate-500">
        Atualizado <FreshnessIndicator asOf={data.as_of} />
      </p>
      {chips.length > 0 ? (
        <div className="mt-4 flex flex-wrap gap-2" aria-label="Filtros rápidos de exceção">
          {chips.map(({ key, count }) => {
            const active = activeException === key;
            return (
              <button
                key={key}
                type="button"
                aria-pressed={active}
                onClick={() => onExceptionChipClick(key)}
                className={`rounded-full border px-3 py-1 text-sm font-medium transition-colors ${
                  active
                    ? "border-blue-600 bg-blue-600 text-white"
                    : "border-slate-200 bg-slate-100 text-slate-700 hover:border-slate-300 hover:bg-slate-200"
                }`}
              >
                {`${exceptionChipLabels[key]}: ${count}`}
              </button>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}
