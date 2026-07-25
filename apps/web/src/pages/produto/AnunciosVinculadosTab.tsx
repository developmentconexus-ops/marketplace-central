import type { JSX } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import type { ListingGroupPage, ListingMarketSignal, ListingMoney, ListingReadModel } from "@marketplace-central/sdk-runtime";
import { DataTable, EmptyState, ErrorState, LoadingState, UnknownValue, formatMoney, formatSignedPercent, type DataTableColumn } from "@marketplace-central/ui";
import { FreshnessIndicator, listingsQueryKeys } from "@marketplace-central/web-query";
import { useClient, type Client } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";

export interface AnunciosVinculadosTabProps {
  productId: string;
  /** Injectable for tests; defaults to useClient(). */
  client?: Client;
  /** Injectable for tests; defaults to useInstallation(). */
  installationId?: string;
}

function usableSignal(item: ListingReadModel): ListingMarketSignal | null {
  const signal = item.market_signal;
  if (!signal || signal.status === "SEM_VINCULO") return null;
  return signal;
}

function Money({ value }: { value: ListingMoney | null }): JSX.Element {
  const formatted = value === null ? null : formatMoney(value.amount, value.currency);
  if (formatted === null) return <UnknownValue />;
  return <span className="font-mono tabular-nums">{formatted}</span>;
}

function PositionCell({ item }: { item: ListingReadModel }): JSX.Element {
  const signal = usableSignal(item);
  if (!signal || signal.position === null) return <UnknownValue />;
  return (
    <span className="font-mono tabular-nums">
      {signal.position.rank}/{signal.position.total}
    </span>
  );
}

function DeltaCell({ item }: { item: ListingReadModel }): JSX.Element {
  const signal = usableSignal(item);
  if (!signal || signal.delta_pct === null) return <UnknownValue />;
  return <span className="font-mono tabular-nums">{formatSignedPercent(signal.delta_pct)}</span>;
}

function FreshnessCell({ item }: { item: ListingReadModel }): JSX.Element {
  const signal = usableSignal(item);
  if (!signal || signal.evidence === null) return <UnknownValue />;
  return <FreshnessIndicator asOf={signal.evidence.fetched_at} />;
}

const columns: DataTableColumn<ListingReadModel>[] = [
  { key: "title", header: "Anúncio", render: (item) => item.title },
  { key: "status", header: "Status", render: (item) => item.status },
  { key: "position", header: "Posição", render: (item) => <PositionCell item={item} />, align: "right" },
  {
    key: "price_to_win",
    header: "Preço p/ vencer",
    render: (item) => <Money value={usableSignal(item)?.price_to_win ?? null} />,
    align: "right",
  },
  { key: "delta_pct", header: "Δ vs. mercado", render: (item) => <DeltaCell item={item} />, align: "right" },
  { key: "freshness", header: "Atualizado", render: (item) => <FreshnessCell item={item} />, align: "right" },
];

/**
 * AnunciosVinculadosTab is the /catalogo/produtos/:id "Anúncios vinculados"
 * tab panel (S4). It scopes listListingsByProduct to this productId and
 * renders the resolved group's listings; any missing market_signal field
 * renders the honest UnknownValue "—", never a fabricated 0 (ADR-17).
 */
export function AnunciosVinculadosTab({
  productId,
  client: injectedClient,
  installationId: injectedInstallationId,
}: AnunciosVinculadosTabProps): JSX.Element {
  const client = injectedClient ?? useClient();
  const installationId = injectedInstallationId ?? useInstallation().installationId;
  const options = { installation_id: installationId, product_id: productId };

  const query = useQuery<ListingGroupPage>({
    queryKey: listingsQueryKeys.byProduct(installationId, options),
    queryFn: () => client.listListingsByProduct(options),
  });

  if (query.isLoading) return <LoadingState />;
  if (query.isError) {
    return <ErrorState detail="Não foi possível carregar os anúncios vinculados." onRetry={() => void query.refetch()} />;
  }

  const groups = query.data?.groups ?? [];
  const group =
    groups.find((candidate) => candidate.product_id === productId) ??
    (groups.length === 1 && (groups[0].product_id === productId || groups[0].product_id === null) ? groups[0] : undefined);
  const listings = group?.listings ?? [];

  return (
    <DataTable
      columns={columns}
      rows={listings}
      rowKey={(item) => item.listing_id}
      emptyState={
        <EmptyState
          hint={
            <Link to="/vinculos" className="font-medium text-accent underline decoration-accent-soft underline-offset-2 hover:text-accent-ink">
              vincular anúncios
            </Link>
          }
        />
      }
    />
  );
}
