import type { CanonicalNumericSourceFact, QualityFlag } from "@marketplace-central/sdk-runtime";
import { useQuery } from "@tanstack/react-query";
import { ErrorState, FreshnessIndicator, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";

export interface ProdutoHeaderProps {
  productId: string;
}

function renderText(value: string | null) {
  return value === null ? <UnknownValue /> : value;
}

function renderNumericFact(fact: CanonicalNumericSourceFact) {
  if (fact.value === null) {
    return <UnknownValue hint={fact.quality_reason ?? undefined} />;
  }
  return (
    <span className="inline-flex items-center gap-2">
      <span>{fact.value}</span>
      <FreshnessIndicator asOf={fact.observed_at} />
    </span>
  );
}

function QualityChip({ flag }: { flag: QualityFlag }) {
  const className = flag === "complete" ? "bg-accent-soft text-accent-ink" : "bg-warn-soft text-warn";
  return (
    <span className={`inline-flex whitespace-nowrap rounded-pill px-2 py-0.5 text-xs font-medium ${className}`}>
      {flag}
    </span>
  );
}

function isErrorWithStatus(error: unknown, status: number): boolean {
  return typeof error === "object" && error !== null && (error as { status?: unknown }).status === status;
}

export function ProdutoHeader({ productId }: ProdutoHeaderProps) {
  const client = useClient();
  const numericId = Number(productId);
  const query = useQuery({
    queryKey: ["catalog", "product", productId],
    queryFn: () => client.getCatalogProduct(numericId),
    staleTime: QUERY_STALE_TIME.catalog,
  });

  if (query.isPending) return <LoadingState />;

  if (query.isError) {
    if (isErrorWithStatus(query.error, 404)) return null;
    return <ErrorState detail="Não foi possível carregar o produto." onRetry={() => void query.refetch()} />;
  }

  const product = query.data;

  return (
    <header aria-labelledby="produto-header-title" className="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
      <div className="flex flex-wrap items-baseline gap-2">
        <span className="font-mono text-xs text-faint">CODPROD {product.internal_product_id}</span>
        <h2 id="produto-header-title" className="text-lg font-semibold tracking-tight text-ink">
          {renderText(product.name)}
        </h2>
      </div>

      <div className="flex flex-wrap items-center gap-2 text-sm text-muted">
        {product.ean === null ? (
          <span className="inline-flex whitespace-nowrap rounded-pill bg-warn-soft px-2 py-0.5 text-xs font-medium text-warn">
            sem EAN — matching limitado a REVIEW
          </span>
        ) : (
          <span className="font-mono">{product.ean}</span>
        )}
        <span className="inline-flex items-center gap-1">
          <span className="text-[10.5px] font-semibold uppercase text-faint">REF</span>
          <span>{renderText(product.manufacturer_reference)}</span>
        </span>
      </div>

      <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
        <div className="rounded-lg border border-border-2 px-2.5 py-2">
          <div className="text-[10.5px] text-faint">Marca</div>
          <div className="mt-0.5 text-sm font-semibold text-ink">{renderText(product.brand_name)}</div>
        </div>
        <div className="rounded-lg border border-border-2 px-2.5 py-2">
          <div className="text-[10.5px] text-faint">NCM</div>
          <div className="mt-0.5 font-mono text-sm font-semibold text-ink">{renderText(product.ncm)}</div>
        </div>
        <div className="col-span-2 rounded-lg border border-border-2 px-2.5 py-2">
          <div className="text-[10.5px] text-faint">Custo</div>
          <div className="mt-0.5 font-mono text-sm font-semibold text-ink">{renderNumericFact(product.cost_amount)}</div>
        </div>
      </div>

      <div className="flex flex-wrap gap-1.5" aria-label="Completude do cadastro">
        {product.quality_flags.map((flag) => (
          <QualityChip key={flag} flag={flag} />
        ))}
      </div>
    </header>
  );
}
