import { useEffect, useState } from "react";
import type { JSX } from "react";
import { useMutation } from "@tanstack/react-query";
import { ErrorState } from "@marketplace-central/ui";
import type { PricingCalcInput, PricingSolveResponse } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { ptBrRateToDot } from "./ptbrDecimal";
import { TariffBadge } from "./tariffBadge";
import type { TariffBlock } from "./tariffBadge";

export interface SolverPanelProps {
  /** Selected product to solve for, or null when none is chosen. */
  productId: number | null;
  /** Active modalidade key. */
  modalidade: string;
  /**
   * Controlled target margin (raw string). When provided together with
   * `onTargetChange`, the margin field is lifted to the parent (the sim panel's
   * Preço/Margem grid) and this panel renders only the action + result. Omit both
   * to keep the panel's own internal target input (backward compatible).
   */
  target?: string;
  onTargetChange?: (value: string) => void;
  /** Hide the built-in margin input (it lives in the parent grid when controlled). */
  hideInput?: boolean;
}

/** Response shape widened with the design §4.4 tarifa block (SDK type lags). */
type SolveResult = PricingSolveResponse & { tarifa?: TariffBlock | null };

/**
 * SolverPanel is the reverse (bidirectional) leg of the simulator: the operator
 * enters a TARGET margin and the server solves for the price that yields it
 * (`pricingSolveTarget`). Result branches key on the backend `code`
 * (design §4.4 Layer 3): SEM_CUSTO blocks; SEM_FRETE gives actionable guidance
 * instead of a false "inatingível"; UNREACHABLE_TARGET surfaces the ceiling ONLY
 * when it is non-empty — a blank ceiling never renders as a lone "%". We NEVER
 * fabricate a price (ADR-17). Money/percentages stay decimal strings end-to-end.
 */
export function SolverPanel({
  productId,
  modalidade,
  target: targetProp,
  onTargetChange,
  hideInput = false,
}: SolverPanelProps): JSX.Element {
  const client = useClient();
  const [targetState, setTargetState] = useState<string>("");
  const controlled = onTargetChange !== undefined;
  const target = controlled ? targetProp ?? "" : targetState;
  const setTarget = controlled ? onTargetChange : setTargetState;
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
        {hideInput ? null : (
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
        )}
        <button
          type="button"
          data-testid="solver-submit"
          disabled={disabled}
          onClick={() => solve.mutate()}
          className="rounded-control border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2 disabled:opacity-40"
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
