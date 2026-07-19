import type { JSX } from "react";
import { UnknownValue, FreshnessIndicator } from "@marketplace-central/ui";

/**
 * Tarifa carimbo carried per component in the /pricing/solve and /pricing/decompose
 * responses (design DESIGN-TARIFAS-ML §4.4). `valor` nil / `sem_dados` = NO-DATA
 * honest state (ADR-17) — never 0, never a green figure. Shape is tolerant/optional so
 * the SDK's stricter PricingTarifaComponent/PricingTarifaFrete assign to it directly.
 */
export interface TariffComponent {
  valor?: string | null;
  fonte?: string | null; // VENDA | COTACAO | CATEGORIA | PADRAO
  degrau?: number | null; // 1-4
  data?: string | null; // event date (venda) or fetched_at (cotação)
  estimativa?: boolean;
  sem_dados?: boolean;
}

export interface TariffBlock {
  comissao?: TariffComponent | null;
  frete?: TariffComponent | null;
}

const FONTE_LABELS: Record<string, string> = {
  VENDA: "Venda",
  COTACAO: "Cotação",
  CATEGORIA: "Categoria",
  PADRAO: "Padrão",
};

export function fonteLabel(fonte: string | null | undefined): string | null {
  if (!fonte) return null;
  return FONTE_LABELS[fonte] ?? fonte;
}

/** NO-DATA = explicit sem_dados, nil valor, OR an empty/blank valor string. */
export function isNoData(comp: TariffComponent): boolean {
  return (
    comp.sem_dados === true ||
    comp.valor == null ||
    (typeof comp.valor === "string" && comp.valor.trim() === "")
  );
}

interface TariffCarimboProps {
  comp: TariffComponent | null | undefined;
  testId?: string;
}

/**
 * The provenance carimbo cluster beside an already-rendered value: fonte + degrau +
 * ESTIMATIVA pill + freshness. Returns null when the component is absent or NO-DATA —
 * a source/degrau/ESTIMATIVA next to an unknown value reads as if it were resolved.
 */
export function TariffCarimbo({ comp, testId }: TariffCarimboProps): JSX.Element | null {
  if (!comp || isNoData(comp)) return null;
  const fonte = fonteLabel(comp.fonte);
  const isEstimativa = comp.estimativa === true || comp.degrau === 4;
  return (
    <span data-testid={testId} className="flex items-center gap-1.5">
      {fonte ? <span className="rounded bg-surface px-1.5 py-0.5 text-muted">{fonte}</span> : null}
      {comp.degrau != null ? <span className="text-faint">degrau {comp.degrau}</span> : null}
      {isEstimativa ? (
        <span className="rounded bg-amber-soft px-1.5 py-0.5 font-medium text-amber">ESTIMATIVA</span>
      ) : null}
      {comp.data ? <FreshnessIndicator asOf={comp.data} /> : null}
    </span>
  );
}

interface TariffBadgeProps {
  testId: string;
  label: string;
  comp: TariffComponent | null | undefined;
  /** Comissão is a percentage; frete is money. */
  percent?: boolean;
}

/**
 * Full labeled badge (SolverPanel result panel): label + value + carimbo cluster.
 * NO-DATA renders UnknownValue ("—") with the carimbos suppressed.
 */
export function TariffBadge({ testId, label, comp, percent }: TariffBadgeProps): JSX.Element | null {
  if (!comp) return null;
  const noData = isNoData(comp);
  return (
    <span data-testid={testId} className="flex items-center gap-1.5">
      <span className="text-muted">{label}</span>
      {noData ? (
        <UnknownValue hint="Sem dados — cadastre dimensões ou vincule um anúncio ML" />
      ) : (
        <>
          <span className="font-mono text-ink">{percent ? `${comp.valor}%` : `R$ ${comp.valor}`}</span>
          <TariffCarimbo comp={comp} />
        </>
      )}
    </span>
  );
}
