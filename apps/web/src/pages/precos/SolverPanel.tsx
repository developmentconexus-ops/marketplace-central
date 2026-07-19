import { useEffect, useState } from "react";
import type { JSX } from "react";
import { useMutation } from "@tanstack/react-query";
import { ErrorState, UnknownValue, FreshnessIndicator } from "@marketplace-central/ui";
import type { PricingCalcInput, PricingSolveResponse } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { ptBrRateToDot } from "./ptbrDecimal";

export interface SolverPanelProps {
  /** Selected product to solve for, or null when none is chosen. */
  productId: number | null;
  /** Active modalidade key. */
  modalidade: string;
}

/**
 * Tarifa carimbo carried per component in the `/pricing/solve` response
 * (design DESIGN-TARIFAS-ML §4.4 Layer 2/3). Owned by the backend chip CHIP-T1;
 * declared locally + read tolerantly here until the SDK type lands the block.
 * `valor` nil / `sem_dados` = NO-DATA honest state (ADR-17) — never 0, never green.
 */
interface TariffComponent {
  valor?: string | null;
  fonte?: string | null; // VENDA | COTACAO | CATEGORIA | PADRAO
  degrau?: number | null; // 1-4
  data?: string | null; // event date (venda) or fetched_at (cotação)
  estimativa?: boolean;
  sem_dados?: boolean;
}

interface TariffBlock {
  comissao?: TariffComponent | null;
  frete?: TariffComponent | null;
}

/** Response shape widened with the design §4.4 tarifa block (SDK type lags). */
type SolveResult = PricingSolveResponse & { tarifa?: TariffBlock | null };

const FONTE_LABELS: Record<string, string> = {
  VENDA: "Venda",
  COTACAO: "Cotação",
  CATEGORIA: "Categoria",
  PADRAO: "Padrão",
};

function fonteLabel(fonte: string | null | undefined): string | null {
  if (!fonte) return null;
  return FONTE_LABELS[fonte] ?? fonte;
}

/**
 * SolverPanel is the reverse (bidirectional) leg of the simulator: the operator
 * enters a TARGET margin and the server solves for the price that yields it
 * (`pricingSolveTarget`). Result branches key on the backend `code`
 * (design §4.4 Layer 3): SEM_CUSTO blocks; SEM_FRETE gives actionable guidance
 * instead of a false "inatingível"; UNREACHABLE_TARGET surfaces the ceiling ONLY
 * when it is non-empty — a blank ceiling never renders as a lone "%". We NEVER
 * fabricate a price (ADR-17). Money/percentages stay decimal strings end-to-end.
 */
export function SolverPanel({ productId, modalidade }: SolverPanelProps): JSX.Element {
  const client = useClient();
  const [target, setTarget] = useState<string>("");
  const [result, setResult] = useState<SolveResult | null>(null);

  const solve = useMutation({
    mutationFn: (): Promise<PricingSolveResponse> => {
      // comissao_pct OMITTED — lets the backend resolver chain run (COTACAO/PADRAO);
      // sending it would force a MANUAL override and hide the real tariff.
      const input: PricingCalcInput = {
        margem_alvo_pct: ptBrRateToDot(target),
        modalidade,
        product_id: productId,
      };
      return client.pricingSolveTarget(input);
    },
    onSuccess: (res) => setResult(res as SolveResult),
  });

  // A stale result must never render under a different product/modalidade — clear it
  // when the selection changes so the operator never reads product A's price on B.
  useEffect(() => {
    setResult(null);
  }, [productId, modalidade]);

  const disabled = productId === null || target.trim() === "" || solve.isPending;

  const code = result?.code;
  const desconhecidos = result?.desconhecidos ?? [];
  const ceiling = result?.ceiling_pct?.trim() ?? "";
  const reached = result?.reached === true && result.preco !== null;

  // A genuine reached+priced result always wins over any lingering banner code.
  // Custo unknown = BLOCKING structural (margin unknowable at any price) — wins over frete.
  const semCusto =
    !reached &&
    (code === "SEM_CUSTO" || result?.blocking_state === "SEM_CUSTO" || desconhecidos.includes("custo"));
  // Frete unknown = segment-conditional — actionable guidance, not "inatingível".
  const semFrete = !reached && !semCusto && (code === "SEM_FRETE" || desconhecidos.includes("frete"));
  // Legitimate unreachable: explicit UNREACHABLE_TARGET (or legacy no-code) WITH a real
  // ceiling. An incomplete-data code carrying a ceiling must never masquerade as one.
  const unreachable =
    !reached &&
    !semCusto &&
    !semFrete &&
    result != null &&
    ceiling !== "" &&
    (code === "UNREACHABLE_TARGET" || code == null);
  // Not reached, no known cause / no legitimate ceiling: honest generic banner, never blank.
  const incomplete = !reached && !semCusto && !semFrete && !unreachable && result != null;

  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-wrap items-end gap-3">
        <label className="flex flex-col text-xs text-muted">
          Margem alvo
          <span className="mt-1 flex items-center gap-1 text-ink">
            <input
              inputMode="decimal"
              value={target}
              onChange={(e) => setTarget(e.target.value)}
              placeholder="0,0"
              aria-label="Margem alvo"
              className="w-24 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink"
            />
            %
          </span>
        </label>
        <button
          type="button"
          data-testid="solver-submit"
          disabled={disabled}
          onClick={() => solve.mutate()}
          className="rounded-md border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2 disabled:opacity-40"
        >
          {solve.isPending ? "Calculando…" : "Calcular preço"}
        </button>
      </div>

      {solve.isError ? (
        <ErrorState onRetry={() => solve.mutate()} detail="Não foi possível calcular o preço para a margem alvo." />
      ) : semCusto ? (
        <p role="alert" data-testid="solver-blocking" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
          Sem custo do ERP para este produto — o preço não pode ser resolvido (SEM_CUSTO).
        </p>
      ) : semFrete ? (
        <p role="alert" data-testid="solver-sem-frete" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
          Sem dados de frete para este produto. Cadastre dimensões (peso, altura, largura, comprimento) OU
          vincule um anúncio ML para cotar o frete.
        </p>
      ) : unreachable ? (
        <div data-testid="solver-unreachable" className="rounded-md bg-amber-soft px-3 py-2 text-sm text-amber">
          Alvo inatingível — a melhor margem possível para este produto é{" "}
          <span className="font-mono">{ceiling}%</span>.
        </div>
      ) : incomplete ? (
        <p role="alert" data-testid="solver-incomplete" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
          Dados incompletos para resolver o preço deste produto.
        </p>
      ) : reached && result ? (
        <div className="flex flex-col gap-2 rounded-md border border-border bg-surface-2 px-3 py-2">
          <div className="flex items-baseline gap-2">
            <span className="text-xs text-muted">Preço sugerido</span>
            <span data-testid="solver-price" className="font-mono text-base text-ink">
              R$ {result.preco}
            </span>
          </div>
          {result.tarifa ? (
            <div data-testid="solver-tarifa" className="flex flex-wrap gap-4 text-xs">
              <TariffBadge testId="tarifa-comissao" label="Comissão" comp={result.tarifa.comissao} percent />
              <TariffBadge testId="tarifa-frete" label="Frete" comp={result.tarifa.frete} />
            </div>
          ) : null}
        </div>
      ) : null}
    </div>
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
 * Renders one resolved tariff component with its fonte/degrau/data carimbos and an
 * ESTIMATIVA pill when the value is a labeled degrau-4 estimate. NO-DATA (nil value
 * or `sem_dados`) renders `UnknownValue` ("—") — never R$0, never a misleading pill.
 */
function TariffBadge({ testId, label, comp, percent }: TariffBadgeProps): JSX.Element | null {
  if (!comp) return null;

  // NO-DATA = explicit sem_dados, nil, OR an empty/blank value string (ceiling_pct uses ""
  // as its blank sentinel; a tarifa valor may too). Never render "R$ " / a false pill.
  const noData =
    comp.sem_dados === true ||
    comp.valor == null ||
    (typeof comp.valor === "string" && comp.valor.trim() === "");
  const isEstimativa = !noData && (comp.estimativa === true || comp.degrau === 4);
  const fonte = fonteLabel(comp.fonte);

  return (
    <span data-testid={testId} className="flex items-center gap-1.5">
      <span className="text-muted">{label}</span>
      {noData ? (
        // Carimbos are suppressed beside "—": a source/degrau/ESTIMATIVA next to an
        // unknown value reads as if the unknown were a real resolved figure.
        <UnknownValue hint="Sem dados — cadastre dimensões ou vincule um anúncio ML" />
      ) : (
        <>
          <span className="font-mono text-ink">
            {percent ? `${comp.valor}%` : `R$ ${comp.valor}`}
          </span>
          {fonte ? <span className="rounded bg-surface px-1.5 py-0.5 text-muted">{fonte}</span> : null}
          {comp.degrau != null ? <span className="text-faint">degrau {comp.degrau}</span> : null}
          {isEstimativa ? (
            <span className="rounded bg-amber-soft px-1.5 py-0.5 font-medium text-amber">ESTIMATIVA</span>
          ) : null}
          {comp.data ? <FreshnessIndicator asOf={comp.data} /> : null}
        </>
      )}
    </span>
  );
}
