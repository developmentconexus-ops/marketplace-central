import { useState, type JSX, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import type { ListingReadModel } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { RepricingTable } from "./RepricingTable";
import { OportunidadesTable } from "./OportunidadesTable";
import { MonitoradosTab } from "./MonitoradosTab";
import { buildOppRows, chunk } from "./oportunidades";
import { useMercadoMarketClient, type MercadoMarketClient } from "./marketClient";
import { DASH, formatCollectedAt, MARGIN_MIN_PCT } from "./mercadoFormatters";

type MercadoTab = "reprice" | "opp" | "mon";

// Both /listings and /catalog/products are cursor-paginated. The radar must scan the WHOLE
// catalog, not an arbitrary first page — a partial count/table silently dropping rows past the
// page is a dishonest count (ADR-17). We cursor-walk to exhaustion under a bounded ceiling, and
// if the ceiling is ever hit we surface it honestly in the UI (never a silent truncation).
const PAGE_SIZE = 100;
const MAX_PAGES = 20; // hard ceiling ⇒ ≤2000 items scanned per tab; truncation is surfaced, never silent
const SCAN_CEILING = PAGE_SIZE * MAX_PAGES;
const READ_ID_CHUNK = 200; // /market batch reads cap at MaxReadIDs=200 — chunk the codprod list

interface CursorPage<T> {
  items: T[];
  next_cursor: string | null;
}

/**
 * Cursor-walk a paginated endpoint to exhaustion under a page ceiling. `truncated` is true iff
 * the ceiling was hit while more pages still remained — the caller MUST surface that honestly.
 */
async function walkAll<T>(
  fetchPage: (cursor: string | undefined) => Promise<CursorPage<T>>,
  maxPages: number,
): Promise<{ items: T[]; truncated: boolean }> {
  const items: T[] = [];
  let cursor: string | undefined;
  for (let page = 0; page < maxPages; page++) {
    const res = await fetchPage(cursor);
    items.push(...res.items);
    if (res.next_cursor == null) return { items, truncated: false };
    cursor = res.next_cursor;
  }
  return { items, truncated: true };
}

export interface MercadoPageProps {
  /** Injectable for tests; defaults to the standalone IC-03 market client. */
  marketClient?: MercadoMarketClient;
}

function fmtCount(n: number | null): string {
  return n == null ? DASH : n.toLocaleString("pt-BR");
}

/** Honest banner shown when a tab's cursor-walk hit the scan ceiling (ADR-17 — no silent cap). */
function TruncationNote(): JSX.Element {
  return (
    <p
      role="status"
      className="rounded-card border border-border bg-surface-2 px-3 py-2 text-[12px] text-muted"
    >
      Varredura limitada aos primeiros {SCAN_CEILING.toLocaleString("pt-BR")} itens — refine a busca
      para cobrir todo o catálogo.
    </p>
  );
}

/**
 * MercadoPage is the /mercado "radar": three underline-active tabs
 * (Reprecificação default · Oportunidades · Monitorados) over the ML market
 * evidence. It composes the central useClient() seam (our listings + ERP catalog
 * facts) with the standalone IC-03 market client (aggregates + verdicts). Every
 * unbacked value renders an honest "—" (ADR-17) and every write control is an
 * inert demo affordance (D-57) — the page performs no live Mercado Livre write.
 */
export function MercadoPage({ marketClient: injectedMarket }: MercadoPageProps = {}): JSX.Element {
  const client = useClient();
  const { installationId } = useInstallation();
  const marketClient = useMercadoMarketClient(injectedMarket);
  const [tab, setTab] = useState<MercadoTab>("reprice");

  const listingsQuery = useQuery({
    queryKey: ["mercado", "listings", installationId],
    staleTime: QUERY_STALE_TIME.listings,
    queryFn: async () => {
      // Honest total comes from the backed summary; the rendered rows come from the cursor
      // walk. If the summary read fails, fall back to the walked length (still honest).
      const summary = await client.getListingsSummary(installationId).catch(() => null);
      const walk = await walkAll<ListingReadModel>(
        (cursor) => client.listListings({ installation_id: installationId, limit: PAGE_SIZE, cursor }),
        MAX_PAGES,
      );
      return {
        items: walk.items,
        truncated: walk.truncated,
        total: summary?.total ?? null,
        asOf: summary?.as_of ?? null,
      };
    },
  });

  // Oportunidades: catalog facts joined with market aggregates + verdicts. Only fetched while
  // the tab is active. The catalog is cursor-walked in full; the codprod set is chunked to the
  // /market MaxReadIDs=200 cap, and the per-chunk aggregates/verdicts are concatenated in fact
  // order so buildOppRows' positional fact↔verdict alignment survives the chunking.
  const oppQuery = useQuery({
    queryKey: ["mercado", "oportunidades", installationId],
    enabled: tab === "opp",
    staleTime: QUERY_STALE_TIME.listings,
    queryFn: async () => {
      const walk = await walkAll(
        (cursor) => client.listCatalogProductFacts({ limit: PAGE_SIZE, cursor }),
        MAX_PAGES,
      );
      const facts = walk.items;
      const codprods = facts.map((f) => String(f.internal_product_id));
      if (codprods.length === 0) return { rows: [], truncated: walk.truncated };
      const batches = await Promise.all(
        chunk(codprods, READ_ID_CHUNK).map(async (ids) => {
          const [aggregates, verdicts] = await Promise.all([
            marketClient.listMarketAggregates(ids),
            marketClient.listMarketVerdicts(ids),
          ]);
          return { aggregates, verdicts };
        }),
      );
      const aggregates = batches.flatMap((b) => b.aggregates);
      const verdicts = batches.flatMap((b) => b.verdicts);
      return { rows: buildOppRows(facts, aggregates, verdicts), truncated: walk.truncated };
    },
  });

  // Operator D-120: order anúncios by observed price gap (mediana ML − nosso preço,
  // both real displayed facts — a sort key, never a rendered margin, ADR-17). Listings
  // without a collected market signal sort after, keeping their original order.
  const rawListingRows = listingsQuery.data?.items ?? [];
  const listingGap = (r: (typeof rawListingRows)[number]): number => {
    const meu = r.price == null ? NaN : Number(r.price.amount);
    const mediana =
      r.market_signal?.median == null ? NaN : Number(r.market_signal.median.amount);
    const value = mediana - meu;
    return Number.isFinite(value) ? value : -Infinity;
  };
  const listingRows = [...rawListingRows].sort((a, b) => listingGap(b) - listingGap(a));
  const repriceCount = listingsQuery.data ? listingsQuery.data.total ?? listingRows.length : null;
  const oppRows = oppQuery.data?.rows ?? [];
  const oppCount = oppQuery.data ? oppRows.length : null;
  const collectedAt = formatCollectedAt(listingsQuery.data?.asOf);
  const listingsTruncated = listingsQuery.data
    ? listingsQuery.data.truncated ||
      (listingsQuery.data.total != null && listingsQuery.data.total > listingRows.length)
    : false;
  const oppTruncated = oppQuery.data?.truncated ?? false;

  const tabs: { key: MercadoTab; label: string }[] = [
    { key: "reprice", label: `Reprecificação · nossos anúncios ${fmtCount(repriceCount)}` },
    { key: "opp", label: `Oportunidades · não vendemos ${fmtCount(oppCount)}` },
    { key: "mon", label: "Monitorados" },
  ];

  let body: ReactNode;
  if (tab === "reprice") {
    if (listingsQuery.isPending) body = <LoadingState />;
    else if (listingsQuery.isError)
      body = <ErrorState onRetry={() => void listingsQuery.refetch()} detail="Falha ao carregar os anúncios." />;
    else if (listingRows.length === 0)
      body = (
        <>
          {listingsTruncated && <TruncationNote />}
          <EmptyState hint="Nenhum anúncio com sinal de mercado ainda." />
        </>
      );
    else
      body = (
        <>
          {listingsTruncated && <TruncationNote />}
          <RepricingTable rows={listingRows} />
        </>
      );
  } else if (tab === "opp") {
    if (oppQuery.isPending) body = <LoadingState />;
    else if (oppQuery.isError)
      body = <ErrorState onRetry={() => void oppQuery.refetch()} detail="Falha ao carregar as oportunidades." />;
    else if (oppRows.length === 0)
      body = (
        <>
          {oppTruncated && <TruncationNote />}
          <EmptyState hint="Nenhum produto do catálogo com demanda de mercado observada." />
        </>
      );
    else
      body = (
        <>
          {oppTruncated && <TruncationNote />}
          <OportunidadesTable rows={oppRows} />
        </>
      );
  } else {
    body = <MonitoradosTab />;
  }

  return (
    <section aria-labelledby="mercado-title" className="mx-auto flex max-w-7xl flex-col gap-[14px]">
      {/* Header: title + Categoria filter + margin threshold pill + source/refresh */}
      <div className="flex flex-wrap items-center gap-[10px]">
        <h1 id="mercado-title" className="text-[22px] font-bold tracking-tight text-ink">
          Mercado
        </h1>
        <button
          type="button"
          disabled
          title="filtro de categoria — em breve"
          className="cursor-not-allowed whitespace-nowrap rounded-pill border border-border px-3 py-1 text-[12.5px] text-muted"
        >
          Categoria ▾
        </button>
        <span className="whitespace-nowrap rounded-pill border border-accent bg-accent-soft px-3 py-1 text-[12.5px] font-semibold text-accent-ink">
          Margem mín: {MARGIN_MIN_PCT}%
        </span>
        <span className="ml-auto flex items-center gap-[10px]">
          <span className="whitespace-nowrap text-[11.5px] text-faint">
            fonte: busca pública ML · coletado {collectedAt ?? DASH}
          </span>
          <button
            type="button"
            disabled
            title="nova coleta de mercado — em breve"
            className="cursor-not-allowed whitespace-nowrap rounded-control border border-border px-3 py-[7px] text-[12.5px] text-muted"
          >
            Atualizar agora
          </button>
        </span>
      </div>

      {/* Underline-active tabs */}
      <div role="tablist" aria-label="Seções do radar" className="flex gap-[2px] overflow-x-auto border-b border-border text-[13px]">
        {tabs.map((t) => {
          const active = tab === t.key;
          return (
            <button
              key={t.key}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => setTab(t.key)}
              className={`whitespace-nowrap border-b-2 px-4 py-2 ${
                active
                  ? "border-accent font-semibold text-accent-ink"
                  : "border-transparent font-normal text-muted hover:text-ink"
              }`}
            >
              {t.label}
            </button>
          );
        })}
      </div>

      {body}
    </section>
  );
}
