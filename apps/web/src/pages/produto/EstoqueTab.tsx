import type { JSX, ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { ErrorState, FreshnessIndicator, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient, type Client } from "../../app/ClientContext";

export interface EstoqueTabProps {
  productId: string;
  /** Injectable for tests; defaults to useClient(). */
  client?: Client;
}

function StockCard({ label, hint, children }: { label: string; hint?: string; children: ReactNode }) {
  return (
    <div className="rounded-lg border border-border-2 px-3 py-3">
      <div className="text-[10.5px] text-faint">{label}</div>
      <div className="mt-0.5 flex items-center gap-2 font-mono text-lg font-semibold text-ink">{children}</div>
      {hint ? <div className="mt-0.5 text-xs text-faint">{hint}</div> : null}
    </div>
  );
}

/**
 * EstoqueTab is the /catalogo/produtos/:id "Estoque" tab panel (S5). FÍSICO
 * mirrors the ProdutoHeader catalog query (same key + Number() conversion so
 * react-query shares one cache entry). RESERVADO/DISPONÍVEL have no MVP data
 * source and always render the honest UnknownValue "—" (ADR-17: never
 * fabricate a derived number from an unknown operand). Concorrência /
 * Pedidos / Histórico / Dados are a deliberate MIS-005 exclusion, not stubs.
 */
export function EstoqueTab({ productId, client: injectedClient }: EstoqueTabProps): JSX.Element {
  const client = injectedClient ?? useClient();
  const numericId = Number(productId);
  const query = useQuery({
    queryKey: ["catalog", "product", productId],
    queryFn: () => client.getCatalogProduct(numericId),
    staleTime: QUERY_STALE_TIME.catalog,
  });

  if (query.isPending) return <LoadingState />;

  if (query.isError) {
    return <ErrorState detail="Não foi possível carregar o estoque." onRetry={() => void query.refetch()} />;
  }

  const product = query.data;
  const stock = product.stock_quantity;
  const hasFisico = stock.value !== null && !product.quality_flags.includes("missing_stock");

  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div className="rounded-lg border border-border-2 px-3 py-3">
        <div className="text-[10.5px] text-faint">FÍSICO</div>
        {hasFisico ? (
          <div className="mt-0.5 flex items-center gap-2 font-mono text-lg font-semibold text-ink">
            <span>{stock.value}</span>
            <FreshnessIndicator asOf={stock.observed_at} />
          </div>
        ) : (
          <div className="mt-2 flex flex-col gap-1 text-sm text-muted">
            <p>Estoque não informado.</p>
            <p className="text-xs text-faint">Importar planilha para habilitar este indicador.</p>
          </div>
        )}
      </div>

      <StockCard label="RESERVADO" hint="DESCONHECIDO">
        <UnknownValue />
      </StockCard>

      <StockCard label="DISPONÍVEL" hint="DESCONHECIDO">
        <UnknownValue />
      </StockCard>
    </div>
  );
}
