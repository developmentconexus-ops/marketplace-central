import type {
  ListingDetail,
  ListingLinkState,
  ListingReadModel,
  ListingStatus,
} from "@marketplace-central/sdk-runtime";
import { useQuery } from "@tanstack/react-query";
import {
  ConflictTag,
  DetailPanel,
  ErrorState,
  FreshnessIndicator,
  LoadingState,
  UnknownValue,
} from "@marketplace-central/ui";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { useClient } from "../app/ClientContext";

export interface ListingDetailPanelProps {
  listingId: string | null;
  onClose: () => void;
}

const statusLabels: Record<ListingStatus, string> = {
  active: "ativo",
  paused: "pausado",
  closed: "encerrado",
  unknown: "desconhecido",
  under_review: "em análise",
  inactive: "inativo",
  payment_required: "pagamento necessário",
  not_yet_active: "ainda não ativo",
};

const linkLabels: Record<ListingLinkState, string> = {
  resolved: "vinculado",
  unresolved: "sem vínculo",
  rejected: "rejeitado",
  conflict: "divergente",
};

function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>{label}</span>;
}

function renderLinkState(state: ListingLinkState) {
  if (state === "conflict") return <ConflictTag />;
  return stateTag(linkLabels[state]);
}

function renderMargin(item: ListingReadModel) {
  if (item.below_margin_worst_case === null) {
    return item.cost === null ? <UnknownValue hint="sem custo no ERP → não simulado" /> : <UnknownValue />;
  }
  return item.below_margin_worst_case
    ? stateTag("abaixo da margem", "bg-red-100 text-red-800")
    : stateTag("ok", "bg-emerald-100 text-emerald-800");
}

function renderFactValue(value: string | number | null, hint?: string) {
  return value === null ? <UnknownValue hint={hint} /> : value;
}

function formatEventTime(at: string) {
  const date = new Date(at);
  if (Number.isNaN(date.getTime())) return at;
  return date.toLocaleString("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  });
}

function Fact({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-slate-100 pb-3 last:border-b-0 last:pb-0">
      <dt className="text-sm text-slate-500">{label}</dt>
      <dd className="text-right text-sm font-medium text-slate-900">{children}</dd>
    </div>
  );
}

function DetailBody({ detail }: { detail: ListingDetail }) {
  const [technicalDetailsOpen, setTechnicalDetailsOpen] = useState(false);

  return (
    <>
      <section aria-labelledby="listing-detail-facts-title">
        <h4 id="listing-detail-facts-title" className="mb-3 text-xs font-semibold uppercase tracking-wide text-slate-500">
          Dados do anúncio
        </h4>
        <dl className="space-y-3">
          <Fact label="Status">{stateTag(statusLabels[detail.status])}</Fact>
          <Fact label="Vínculo">{renderLinkState(detail.link.state)}</Fact>
          <Fact label="Preço">
            {detail.price === null ? <UnknownValue /> : `R$ ${detail.price.amount}`}
          </Fact>
          <Fact label="Estoque">{renderFactValue(detail.published_quantity)}</Fact>
          <Fact label="Vendas 30d">{renderFactValue(detail.sales_30d)}</Fact>
          <Fact label="Qualidade">{renderFactValue(detail.quality_score)}</Fact>
          <Fact label="Margem">{renderMargin(detail)}</Fact>
        </dl>
      </section>

      <section aria-labelledby="listing-detail-freshness-title">
        <h4 id="listing-detail-freshness-title" className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
          Atualização
        </h4>
        <p className="text-sm text-slate-600">
          {detail.fetched_at === null ? <UnknownValue /> : <FreshnessIndicator asOf={detail.fetched_at} />}
        </p>
      </section>

      {detail.sync_error ? (
        <section aria-labelledby="listing-detail-sync-error-title" className="rounded-lg border border-red-200 bg-red-50 p-3">
          <h4 id="listing-detail-sync-error-title" className="text-sm font-semibold text-red-900">Erro de sincronização</h4>
          <p className="mt-1 text-sm text-red-800">{detail.sync_error.message_pt}</p>
          {detail.sync_error.message_provider !== null ? (
            <details
              className="mt-2 text-sm text-red-900"
              onToggle={(event) => setTechnicalDetailsOpen(event.currentTarget.open)}
            >
              <summary className="cursor-pointer" onClick={() => setTechnicalDetailsOpen((open) => !open)}>▸ técnico</summary>
              {technicalDetailsOpen ? (
                <p className="mt-2 break-words font-mono text-xs">{detail.sync_error.message_provider}</p>
              ) : null}
            </details>
          ) : null}
        </section>
      ) : null}

      <section aria-labelledby="listing-detail-timeline-title">
        <h4 id="listing-detail-timeline-title" className="mb-3 text-xs font-semibold uppercase tracking-wide text-slate-500">
          Histórico
        </h4>
        <ol className="space-y-3">
          {detail.timeline.map((event, index) => (
            <li key={`${event.at}-${index}`} data-testid="timeline-event" className="border-l-2 border-slate-200 pl-3">
              <time dateTime={event.at} className="block text-xs text-slate-500">{formatEventTime(event.at)}</time>
              <p className="mt-1 text-sm text-slate-800">{event.message_pt}</p>
            </li>
          ))}
        </ol>
      </section>
    </>
  );
}

function FutureActions() {
  return (
    <div className="flex gap-2">
      {(["Corrigir", "Simular", "Pausar"] as const).map((label) => (
        <button
          key={label}
          type="button"
          disabled
          title="disponível em breve"
          className="flex-1 rounded-lg border border-slate-200 px-2 py-2 text-sm font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {label}
        </button>
      ))}
    </div>
  );
}

export function ListingDetailPanel({ listingId, onClose }: ListingDetailPanelProps) {
  const client = useClient();
  const query = useQuery({
    queryKey: listingsQueryKeys.detail(listingId as string),
    queryFn: () => client.getListing(listingId as string),
    staleTime: QUERY_STALE_TIME.listings,
    enabled: listingId !== null,
  });

  if (listingId === null) return null;

  const detail = query.data;
  const title = detail?.title ?? detail?.provider_listing_id ?? listingId;

  return (
    <DetailPanel
      open={true}
      onClose={onClose}
      closeLabel="Fechar painel"
      title={title}
      subtitle={detail?.provider_listing_id}
      footer={<FutureActions />}
    >
      {query.isPending ? <LoadingState /> : query.isError ? <ErrorState onRetry={() => void query.refetch()} /> : detail ? <DetailBody detail={detail} /> : null}
    </DetailPanel>
  );
}
