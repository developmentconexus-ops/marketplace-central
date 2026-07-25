import type { JSX } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { CanonicalNumericSourceFact, MarketPriceIntelVerdict } from "@marketplace-central/sdk-runtime";
import { Button, ErrorState, LoadingState, UnknownValue, formatMoney } from "@marketplace-central/ui";
import { useProdutoMarketClient, type ProdutoMarketClient } from "./marketClient";

const BLOCKING_STATE_COPY: Record<string, string> = {
  NO_CANDIDATE: "sem candidato de mercado — veredicto não avaliável",
  NO_PRICE_EVIDENCE: "sem evidência de preço — veredicto não avaliável",
  INSUFFICIENT_MARKET: "mercado insuficiente — veredicto não avaliável",
  SEM_CUSTO: "sem custo — margem não avaliável",
};

export interface VeredictoBoxProps {
  productId: string;
  costFact?: CanonicalNumericSourceFact;
  /** Injectable for tests; defaults to the standalone IC-03 market client. */
  client?: ProdutoMarketClient;
}

function Money({ value }: { value: { amount: string; currency: string } | null }): JSX.Element {
  const formatted = value === null ? null : formatMoney(value.amount, value.currency);
  if (formatted === null) return <UnknownValue />;
  return <span className="font-mono text-ink">{formatted}</span>;
}

function Headline({ verdict }: { verdict: MarketPriceIntelVerdict }): JSX.Element {
  if (verdict.price_evidence_status === "INSUFFICIENT_MARKET") {
    if (verdict.market_range === null) {
      return (
        <p className="text-sm font-semibold text-warn">
          vendedores insuficientes — mínimo 5
        </p>
      );
    }
    return (
      <p className="text-sm font-semibold text-warn">
        {verdict.market_range.n_sellers} vendedores — mínimo 5
      </p>
    );
  }
  if (verdict.price_evidence_status === "NO_PRICE_EVIDENCE") {
    return <p className="text-sm font-semibold text-muted">Este produto ainda não tem evidência de preço coletada.</p>;
  }
  return <p className="text-sm font-semibold text-accent-ink">Mercado com evidência de preço.</p>;
}

/**
 * VeredictoBox is the /catalogo/produtos/:id "Veredicto" tab panel. It reads
 * the M-02 verdict for the product and offers a single synchronous "coletar
 * agora" action (no polling — B.1's collection resolves in the same request).
 * verdict_label is always null pre-M-07: the margin row always renders the
 * honest UnknownValue, never a fabricated saudável/apertado label (ADR-17).
 */
export function VeredictoBox({ productId, costFact, client: injected }: VeredictoBoxProps): JSX.Element {
  const marketClient = useProdutoMarketClient(injected);
  const queryClient = useQueryClient();

  const verdictQuery = useQuery({
    queryKey: ["market", "verdict", productId],
    queryFn: async () => {
      const [verdict] = await marketClient.listMarketVerdicts([productId]);
      return verdict ?? null;
    },
  });

  const collectMutation = useMutation({
    mutationFn: () => marketClient.collectMarketPriceIntel(productId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["market", "verdict", productId] });
      void queryClient.invalidateQueries({ queryKey: ["listings", "by-product"] });
    },
  });

  if (verdictQuery.isLoading) return <LoadingState />;
  if (verdictQuery.isError) {
    return <ErrorState detail="Não foi possível carregar o veredicto de mercado." onRetry={() => void verdictQuery.refetch()} />;
  }

  const verdict = verdictQuery.data;
  if (!verdict) {
    return <p className="text-sm text-muted">Sem veredicto disponível para este produto.</p>;
  }

  const marketRange = verdict.market_range;

  return (
    <div className="flex flex-col gap-4 rounded-xl border border-border bg-surface p-4">
      <div className="flex items-start justify-between gap-3">
        <Headline verdict={verdict} />
        <Button
          type="button"
          variant="primary"
          disabled={collectMutation.isPending}
          onClick={() => collectMutation.mutate()}
        >
          {collectMutation.isPending ? "coletando…" : "coletar agora"}
        </Button>
      </div>

      {verdict.blocking_state ? (
        <p role="note" className="text-xs text-warn">
          {BLOCKING_STATE_COPY[verdict.blocking_state] ?? verdict.blocking_state}
        </p>
      ) : null}

      {collectMutation.isError ? (
        <p role="alert" className="rounded-md bg-warn-soft px-3 py-2 text-xs text-warn">
          Falha ao coletar preços de mercado.
        </p>
      ) : null}

      <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm sm:grid-cols-3">
        <dt className="text-muted">Mínimo válido</dt>
        <dd className="text-right"><Money value={marketRange?.min_valid ?? null} /></dd>
        <dt className="text-muted">Mediana</dt>
        <dd className="text-right"><Money value={marketRange?.median ?? null} /></dd>
        <dt className="text-muted">Vendedores</dt>
        <dd className="text-right font-mono text-ink">
          {marketRange ? marketRange.n_sellers : <UnknownValue />}
        </dd>
        <dt className="text-muted">Custo ERP</dt>
        <dd className="text-right font-mono text-ink">
          {costFact === undefined || costFact.value === null ? (
            <UnknownValue hint={costFact?.quality_reason ?? undefined} />
          ) : (
            costFact.value
          )}
        </dd>
        <dt className="text-muted">Veredicto de margem</dt>
        <dd className="text-right">
          <UnknownValue hint="veredicto de margem — M-07" />
        </dd>
      </dl>

      <div className="flex flex-col gap-1">
        <h4 className="text-xs font-semibold text-muted">Evidência</h4>
        {Object.keys(verdict.inputs_used).length === 0 ? (
          <p className="text-xs text-faint">Sem evidência coletada.</p>
        ) : (
          <ul className="flex flex-col gap-0.5">
            {Object.entries(verdict.inputs_used).map(([key, ref]) => (
              <li key={key} className="flex justify-between text-xs">
                <span className="font-mono text-ink">{ref.source}</span>
                <span className="font-mono text-muted">{ref.as_of}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
