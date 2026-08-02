import type {
  ListingDetail,
  ListingLinkState,
  ListingMarketSignal,
  ListingMoney,
  ListingReadModel,
  ListingSyncState,
} from "@marketplace-central/sdk-runtime";
import { useQuery } from "@tanstack/react-query";
import {
  ConflictTag,
  ErrorState,
  formatMoney as formatBRL,
  formatSignedPercent,
  FreshnessIndicator,
  LoadingState,
  UnknownValue,
} from "@marketplace-central/ui";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { Link } from "react-router-dom";
import { useClient } from "../app/ClientContext";

export interface ListingDetailPanelProps {
  listingId: string | null;
  onClose: () => void;
}

const linkStateLabels: Record<ListingLinkState, string> = {
  resolved: "vinculado",
  unresolved: "sem vínculo",
  rejected: "rejeitado",
  conflict: "divergente",
};

// SYNC status pill copy/colors mirror the ratified drawer status pill, which
// reuses the table SYNC column's pillFor() mapping — duplicated here (not
// imported) because AnunciosTable.tsx's map is page-internal and out of this
// slice's write_set; see pages/AnunciosTable.tsx:27-53 for the source pattern.
const syncPillLabels: Record<ListingSyncState, string> = {
  synced: "sincronizado",
  error: "com erro",
  stale: "desatualizado",
  queued: "na fila",
  syncing: "sincronizando",
  paused_sync: "pausado",
};

const syncPillClassName: Record<ListingSyncState, string> = {
  synced: "bg-accent-soft text-accent-ink",
  error: "bg-warn-soft text-warn",
  stale: "bg-amber-soft text-amber",
  queued: "bg-info-soft text-info",
  syncing: "bg-surface-2 text-faint",
  paused_sync: "bg-surface-2 text-faint",
};

function stateTag(label: string, className = "bg-surface-2 text-muted") {
  return <span className={`inline-flex whitespace-nowrap rounded-pill px-2 py-0.5 text-xs font-medium ${className}`}>{label}</span>;
}

function renderMargin(item: ListingReadModel) {
  if (item.below_margin_worst_case === null) {
    return item.cost === null ? <UnknownValue hint="sem custo no ERP → não simulado" /> : <UnknownValue />;
  }
  return item.below_margin_worst_case
    ? stateTag("abaixo da margem", "bg-warn-soft text-warn")
    : stateTag("ok", "bg-accent-soft text-accent-ink");
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

// Timeline row color follows the event's own `kind` (real field, never
// inferred/fabricated): sync errors read warn, stale/queue waits read amber,
// everything else (synced/created/...) reads muted — matches the ratified
// drawer's colored timeline text without inventing new event data.
function timelineColor(kind: string) {
  if (kind === "sync_error") return "text-warn";
  if (kind === "stale" || kind === "queued") return "text-amber";
  return "text-muted";
}

function renderEvidenceCount(value: number | null) {
  return value === null ? <UnknownValue /> : value;
}

function formatPriceToWin(signal: ListingMarketSignal) {
  const formatted = signal.price_to_win === null ? null : formatBRL(signal.price_to_win.amount, signal.price_to_win.currency);
  return formatted === null ? <UnknownValue /> : formatted;
}

function formatMoney(money: ListingMoney | null) {
  const formatted = money === null ? null : formatBRL(money.amount, money.currency);
  return formatted === null ? <UnknownValue /> : formatted;
}

// Faixa de mercado: competitor price range min — mediana — max, our own offer
// excluded upstream (own-seller exclusion at the market aggregation seam). Each
// bound is honest per ADR-17 — a null part renders "—", never a fabricated 0,
// even when the other bounds are present.
function FaixaCard({ signal }: { signal: ListingMarketSignal }) {
  return (
    <div className="col-span-2 rounded-lg border border-border-2 px-2.5 py-2">
      <div className="text-[10.5px] text-faint">Faixa de mercado (concorrentes)</div>
      <div className="mt-0.5 flex items-baseline gap-1.5 font-mono text-sm font-semibold text-ink">
        <span>{formatMoney(signal.min_valid)}</span>
        <span className="text-faint">—</span>
        <span>{formatMoney(signal.median)}</span>
        <span className="text-faint">—</span>
        <span>{formatMoney(signal.max_valid)}</span>
      </div>
      <div className="mt-0.5 text-[10.5px] text-faint">mín · mediana · máx</div>
    </div>
  );
}

function formatPosition(signal: ListingMarketSignal) {
  return signal.position ? `${signal.position.rank}/${signal.position.total}` : <UnknownValue />;
}

function formatDelta(signal: ListingMarketSignal) {
  return signal.delta_pct === null
    ? <UnknownValue />
    : formatSignedPercent(signal.delta_pct);
}

// 2x2 (and evidence-grid) bordered card idiom from the ratified drawer:
// faint uppercase-ish label + mono bold value, border-border-2, rounded-lg.
function InfoCard({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-border-2 px-2.5 py-2">
      <div className="text-[10.5px] text-faint">{label}</div>
      <div className="mt-0.5 font-mono text-sm font-semibold text-ink">{children}</div>
    </div>
  );
}

// Evidência vs. mercado: distinct honest rendering per signal_status (ADR-17) —
// SEM_VINCULO shows no market numbers (link to /vinculos to resolve the link),
// NO_PRICE_EVIDENCE (and a missing/undefined signal_status, treated as the
// honest absent state) shows a named "sem evidência" state — never a
// fabricated number — and OK/STALE show the underlying evidence as ADDITIVE
// bordered cards in the SAME grid/token idiom as the facts grid above (small
// n_offers/n_sellers samples included verbatim, never suppressed), with STALE
// additionally marked via FreshnessIndicator so old numbers read as old.
function EvidenceSection({ detail }: { detail: ListingDetail }) {
  const status = detail.signal_status ?? "NO_PRICE_EVIDENCE";
  const heading = (
    <h4 id="listing-detail-evidence-title" className="mb-2 text-[10.5px] font-semibold uppercase tracking-wide text-faint">
      Vs. mercado
    </h4>
  );

  if (status === "SEM_VINCULO") {
    return (
      <section aria-labelledby="listing-detail-evidence-title">
        {heading}
        <p className="text-sm text-muted">
          Anúncio sem vínculo com produto — sem evidência de mercado.{" "}
          <Link to="/vinculos" className="font-medium text-accent-ink underline underline-offset-2 hover:opacity-80">
            Vincular produto
          </Link>
        </p>
      </section>
    );
  }

  const signal = detail.market_signal;
  if (status === "NO_PRICE_EVIDENCE" || !signal) {
    return (
      <section aria-labelledby="listing-detail-evidence-title">
        {heading}
        <p className="text-sm text-muted">
          <UnknownValue hint="sem evidência de preço de mercado" /> Sem evidência de preço de mercado
        </p>
      </section>
    );
  }

  return (
    <section aria-labelledby="listing-detail-evidence-title">
      {heading}
      <div className="grid grid-cols-2 gap-2">
        <InfoCard label="Posição">{formatPosition(signal)}</InfoCard>
        <InfoCard label="Alvo buybox (ML)">{formatPriceToWin(signal)}</InfoCard>
        <InfoCard label="Diferença">{formatDelta(signal)}</InfoCard>
        <InfoCard label="Match">{signal.match_status ?? <UnknownValue />}</InfoCard>
        <FaixaCard signal={signal} />
        <InfoCard label="Ofertas">{renderEvidenceCount(signal.n_offers)}</InfoCard>
        <InfoCard label="Concorrentes">{renderEvidenceCount(signal.n_sellers)}</InfoCard>
        <InfoCard label="Fonte">{signal.evidence?.source ?? <UnknownValue />}</InfoCard>
        <InfoCard label="Coletado em">
          {signal.evidence?.fetched_at ? <FreshnessIndicator asOf={signal.evidence.fetched_at} /> : <UnknownValue />}
        </InfoCard>
      </div>
      {status === "STALE" ? (
        <p className="mt-2 text-xs text-amber">Evidência desatualizada — pode não refletir o mercado atual.</p>
      ) : null}
      {detail.link.product_id ? (
        <Link
          to={`/catalogo/produtos/${detail.link.product_id}`}
          className="mt-2 inline-block text-xs font-medium text-accent-ink underline underline-offset-2 hover:opacity-80"
        >
          Ver produto vinculado
        </Link>
      ) : null}
    </section>
  );
}

// Foto placeholder + título + produto/modalidade + sync-status pill — mirrors
// the ratified drawer's identity row. PRODUTO reuses the same honest
// present/absent rule as AnunciosTable's PRODUTO cell (product_id, else the
// link-state label, else the conflict tag) — see
// pages/AnunciosTable.tsx:64-77 for the source pattern.
function IdentityRow({ detail }: { detail: ListingDetail }) {
  const produtoNode = detail.link.product_id
    ? detail.link.product_id
    : detail.link.state === "conflict"
      ? <ConflictTag />
      : linkStateLabels[detail.link.state];
  const modalidade = detail.listing_type?.label;

  return (
    <div className="flex gap-2.5">
      <div
        aria-hidden="true"
        className="flex h-14 w-14 flex-none items-center justify-center rounded-lg border border-border bg-surface-2 text-[10px] text-faint"
      >
        foto
      </div>
      <div className="min-w-0">
        <p className="text-[12.5px] font-semibold text-ink">{detail.title ?? detail.provider_listing_id}</p>
        <p className="mt-0.5 truncate text-[11.5px] text-faint">
          {produtoNode}
          {modalidade ? ` · ${modalidade}` : ""}
        </p>
        <span
          className={`mt-1 inline-flex items-center gap-1 whitespace-nowrap rounded-pill px-2 py-0.5 text-xs font-medium ${syncPillClassName[detail.sync_state]}`}
        >
          {syncPillLabels[detail.sync_state]}
        </span>
      </div>
    </div>
  );
}

function DetailBody({ detail }: { detail: ListingDetail }) {
  const [technicalDetailsOpen, setTechnicalDetailsOpen] = useState(false);

  return (
    <>
      <IdentityRow detail={detail} />

      {detail.sync_error ? (
        <section aria-labelledby="listing-detail-sync-error-title" className="rounded-lg border border-warn bg-warn-soft p-3">
          <h4 id="listing-detail-sync-error-title" className="text-sm font-semibold text-warn">Erro de sincronização</h4>
          <p className="mt-1 text-sm text-warn">{detail.sync_error.message_pt}</p>
          {detail.sync_error.message_provider !== null ? (
            <details
              className="mt-2 text-sm text-warn"
              onToggle={(event) => setTechnicalDetailsOpen(event.currentTarget.open)}
            >
              <summary className="cursor-pointer" onClick={() => setTechnicalDetailsOpen((open) => !open)}>▸ técnico</summary>
              {technicalDetailsOpen ? (
                <p className="mt-2 break-words font-mono text-xs text-faint">{detail.sync_error.message_provider}</p>
              ) : null}
            </details>
          ) : null}
        </section>
      ) : null}

      <section aria-label="Dados do anúncio">
        <div className="grid grid-cols-2 gap-2">
          <InfoCard label="Preço">{formatMoney(detail.price)}</InfoCard>
          <InfoCard label="Est. publicado">{renderFactValue(detail.published_quantity)}</InfoCard>
          <InfoCard label="Margem est.">{renderMargin(detail)}</InfoCard>
        </div>
      </section>

      <FutureActions detail={detail} />

      <EvidenceSection detail={detail} />

      <section aria-labelledby="listing-detail-timeline-title">
        <h4 id="listing-detail-timeline-title" className="mb-2 text-[10.5px] font-semibold uppercase tracking-wide text-faint">
          Linha do tempo
        </h4>
        <ol className="space-y-1.5">
          {detail.timeline.map((event, index) => (
            <li key={`${event.at}-${index}`} data-testid="timeline-event" className="flex gap-2 text-xs">
              <time dateTime={event.at} className="flex-none font-mono text-faint">{formatEventTime(event.at)}</time>
              <span className={timelineColor(event.kind)}>{event.message_pt}</span>
            </li>
          ))}
        </ol>
      </section>

      {detail.link.product_id ? (
        <Link to={`/catalogo/produtos/${detail.link.product_id}`} className="text-xs font-semibold text-accent-ink hover:underline">
          Abrir edição completa →
        </Link>
      ) : null}
    </>
  );
}

// Action row: kept intentionally inert (D-57 — zero ML writes, no inline edit
// affordance yet). Rendered in the drawer body right after the stat grid, where
// the ratified prototype places Corrigir/Simular/Pausar. The primary action
// label follows the listing's real state (error → corrigir, unlinked → vincular,
// else editar) to mirror the prototype without wiring any mutation.
function FutureActions({ detail }: { detail: ListingDetail }) {
  const primaryLabel = detail.sync_error
    ? "Corrigir anúncio"
    : detail.link.state === "unresolved"
      ? "Vincular produto"
      : "Editar anúncio";
  return (
    <div className="flex flex-wrap gap-1.5">
      <button
        type="button"
        disabled
        title="disponível em breve"
        className="rounded-lg bg-accent px-3 py-1.5 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-50"
      >
        {primaryLabel}
      </button>
      <button
        type="button"
        disabled
        title="disponível em breve"
        className="rounded-lg border border-border px-3 py-1.5 text-xs font-medium text-muted disabled:cursor-not-allowed disabled:opacity-50"
      >
        Simular preço
      </button>
      <button
        type="button"
        disabled
        title="disponível em breve"
        className="rounded-lg border border-border px-3 py-1.5 text-xs font-medium text-muted disabled:cursor-not-allowed disabled:opacity-50"
      >
        Pausar
      </button>
    </div>
  );
}

// Inline detail drawer: a 300px sticky column beside the table (not an overlay),
// matching the ratified prototype (design/handoff/Anuncios.dc.html:102-141).
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

  return (
    <aside
      aria-label="Detalhe do anúncio"
      className="sticky top-4 w-[300px] flex-none overflow-hidden rounded-xl border border-border bg-surface"
    >
      <div className="flex items-center gap-2 border-b border-border bg-surface-2 px-3.5 py-2.5">
        <span className="flex-none font-mono text-[11.5px] text-faint">
          {detail?.provider_listing_id ?? listingId}
        </span>
        {detail?.title ? (
          <b data-testid="listing-detail-header-title" className="min-w-0 flex-1 truncate text-[13px] font-semibold text-ink">
            {detail.title}
          </b>
        ) : (
          <span className="flex-1" />
        )}
        <button type="button" onClick={onClose} aria-label="Fechar painel" className="flex-none text-faint hover:text-ink">
          ✕
        </button>
      </div>
      <div className="flex flex-col gap-3 p-3.5 text-[12.5px]">
        {query.isPending ? (
          <LoadingState />
        ) : query.isError ? (
          <ErrorState onRetry={() => void query.refetch()} />
        ) : detail ? (
          <DetailBody detail={detail} />
        ) : null}
      </div>
    </aside>
  );
}
