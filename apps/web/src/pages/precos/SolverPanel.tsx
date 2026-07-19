import { useState } from "react";
import type { JSX } from "react";
import { useMutation } from "@tanstack/react-query";
import { ErrorState } from "@marketplace-central/ui";
import type { PricingCalcInput, PricingSolveResponse } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { ptBrRateToDot } from "./ptbrDecimal";

export interface SolverPanelProps {
  /** Selected product to solve for, or null when none is chosen. */
  productId: number | null;
  /** ML commission for the active modalidade (decimal string). */
  comissaoPct: string;
  /** Active modalidade key. */
  modalidade: string;
}

/**
 * SolverPanel is the reverse (bidirectional) leg of the simulator: the operator
 * enters a TARGET margin and the server solves for the price that yields it
 * (`pricingSolveTarget`). When the target is unreachable the response carries
 * `reached=false` plus the achievable `ceiling_pct` — we surface that ceiling and
 * NEVER fabricate a price (ADR-17). SEM_CUSTO (no ERP cost) blocks the solve; the
 * banner names the state. Money/percentages stay decimal strings end-to-end.
 */
export function SolverPanel({ productId, comissaoPct, modalidade }: SolverPanelProps): JSX.Element {
  const client = useClient();
  const [target, setTarget] = useState<string>("");
  const [result, setResult] = useState<PricingSolveResponse | null>(null);

  const solve = useMutation({
    mutationFn: (): Promise<PricingSolveResponse> => {
      const input: PricingCalcInput = {
        margem_alvo_pct: ptBrRateToDot(target),
        comissao_pct: comissaoPct,
        modalidade,
        product_id: productId,
      };
      return client.pricingSolveTarget(input);
    },
    onSuccess: (res) => setResult(res),
  });

  const disabled = productId === null || target.trim() === "" || solve.isPending;
  const semCusto = result?.blocking_state === "SEM_CUSTO";
  const reached = result?.reached === true && result.preco !== null;

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
      ) : null}

      {semCusto ? (
        <p role="alert" data-testid="solver-blocking" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
          Sem custo do ERP para este produto — o preço não pode ser resolvido (SEM_CUSTO).
        </p>
      ) : result && !reached ? (
        <div data-testid="solver-unreachable" className="rounded-md bg-amber-soft px-3 py-2 text-sm text-amber">
          Alvo inatingível — a melhor margem possível para este produto é{" "}
          <span className="font-mono">{result.ceiling_pct}%</span>.
        </div>
      ) : reached && result ? (
        <div className="flex items-baseline gap-2 rounded-md border border-border bg-surface-2 px-3 py-2">
          <span className="text-xs text-muted">Preço sugerido</span>
          <span data-testid="solver-price" className="font-mono text-base text-ink">
            R$ {result.preco}
          </span>
        </div>
      ) : null}
    </div>
  );
}
