import { useQuery } from "@tanstack/react-query";
import type { OrderPage, OrderRead } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { ordersQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState, type ReactNode } from "react";
import { useSearchParams } from "react-router-dom";
import type { Client } from "../../app/ClientContext";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { FilaView } from "./FilaView";
import { KanbanView } from "./KanbanView";
import { ListaView } from "./ListaView";
import { PedidoDrawer } from "./PedidoDrawer";
import { formatMoney } from "./pedidosFormatters";
import { bucketTabCount, type PedidosTab } from "./pedidosTabs";

type PedidosView = "fila" | "lista" | "kanban";

const viewOptions: { value: PedidosView; label: string }[] = [
  { value: "fila", label: "Fila" },
  { value: "lista", label: "Lista" },
  { value: "kanban", label: "Kanban" },
];

// The em dash matches UnknownValue's glyph so a KPI reads as an honest unknown, never a
// fabricated 0 (ADR-17).
const UNKNOWN_KPI_VALUE = "—";

// Safety cap on the pagination accumulator below (G2/F02-S5): a backstop against an unbounded
// loop if the API ever returns a next_cursor that never resolves to null. 20 pages is far beyond
// the demo dataset.
const MAX_ORDER_PAGES = 20;

// Fila/Kanban rows, the Lista per-tab counts AND the KPI row's counts all derive from this single
// query's full dataset (each order's OrderRead.bucket, the server-side single authority that now
// consumes the delivered tag + live shipment status). Deriving the KPI cards from the same rows the
// Lista tabs count makes KPI == Lista by construction — the old getOrderSummary by_status path was
// shipment-status-blind (pure SQL over provider_status+shipping_id) and disagreed with the list.
// This accumulator follows OrderPage.next_cursor until it is exhausted (or the safety cap is hit)
// so the single read query returns the complete dataset — no pagination UI, no second query.
async function fetchAllOrders(client: Client, installationId: string): Promise<OrderPage> {
  const items: OrderRead[] = [];
  let cursor: string | undefined;
  for (let page = 0; page < MAX_ORDER_PAGES; page += 1) {
    const result = await client.listOrders({ installation_id: installationId, cursor });
    items.push(...result.items);
    if (!result.next_cursor) {
      return { items, next_cursor: null };
    }
    cursor = result.next_cursor;
  }
  // Cap reached while next_cursor was still non-null: an honest partial-load edge, not a silent
  // "complete" claim. Acceptable for the demo; flagged so it isn't mistaken for a full dataset.
  console.warn("[pedidos] order list pagination hit MAX_ORDER_PAGES, dataset may be incomplete");
  return { items, next_cursor: null };
}

type KpiTone = "default" | "amber";

interface KpiCardProps {
  label: string;
  value: string;
  sub?: string;
  tone?: KpiTone;
  onSelect: () => void;
}

// Design KPI card (Pedidos.dc.html): flex-1, min-width 150, 1.5px border, 10.5px label, 19px
// mono value with an inline nota. The nota is aria-hidden so the button's accessible name stays
// "LABEL VALUE" (the counts the tests assert), not "LABEL nota VALUE".
function KpiCard({ label, value, sub, tone = "default", onSelect }: KpiCardProps) {
  const toneClass =
    tone === "amber" ? "border-amber bg-amber-soft" : "border-border bg-surface";
  return (
    <button
      type="button"
      onClick={onSelect}
      className={`min-w-[150px] flex-1 rounded-[11px] border-[1.5px] ${toneClass} px-[15px] py-[10px] text-left transition-colors hover:border-accent focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent`}
    >
      <span
        className={`block text-[10.5px] font-semibold uppercase tracking-[0.06em] ${
          tone === "amber" ? "text-amber" : "text-faint"
        }`}
      >
        {label}
      </span>
      <span
        className={`mt-1 block whitespace-nowrap font-mono text-[19px] font-bold ${
          tone === "amber" ? "text-amber" : "text-ink"
        }`}
      >
        {value}
        {sub ? (
          <span aria-hidden className="ml-2 font-sans text-[11px] font-semibold text-faint">
            {sub}
          </span>
        ) : null}
      </span>
    </button>
  );
}

export function PedidosPage() {
  const client = useClient();
  const { installationId } = useInstallation();
  const [tab, setTab] = useState<PedidosTab>("novo");
  const [view, setView] = useState<PedidosView>("fila");
  const [searchParams, setSearchParams] = useSearchParams();
  const openOrderId = searchParams.get("order");

  // Deep-link state lives in the `order` query param (not a route path param), so a row/card
  // click and a shared/bookmarked URL open the same drawer. `replace: true` keeps opening and
  // closing the drawer out of browser history (G2).
  const openOrder = (orderId: string) => {
    const next = new URLSearchParams(searchParams);
    next.set("order", orderId);
    setSearchParams(next, { replace: true });
  };
  const closeOrder = () => {
    const next = new URLSearchParams(searchParams);
    next.delete("order");
    setSearchParams(next, { replace: true });
  };

  const ordersQuery = useQuery({
    queryKey: ordersQueryKeys.list(installationId, {}),
    queryFn: () => fetchAllOrders(client, installationId),
    staleTime: QUERY_STALE_TIME.orders,
  });

  const allItems: OrderRead[] = ordersQuery.data?.items ?? [];

  // KPI counts derive from the fully-loaded list's per-order bucket (same source as the Lista
  // tabs), so a card and its tab always agree. Until the list resolves — or if it errors — the
  // count is an honest unknown "—", never a fabricated 0 (ADR-17). A loaded-but-empty dataset
  // legitimately yields 0 in every bucket.
  const ordersReady = !ordersQuery.isPending && !ordersQuery.isError;
  const kpiValue = (bucket: PedidosTab): string =>
    ordersReady ? String(bucketTabCount(allItems, bucket) ?? 0) : UNKNOWN_KPI_VALUE;

  // DIFAL urgency banner is DATA-GATED: today every order's difal.amount is null (decomposer not
  // wired, hub C2), so difalPending is empty and the banner is omitted — never fabricated. Once C2
  // lands with real amounts, the same derivation lights the banner up with no UI change.
  const difalPending = allItems.filter(
    (order) => order.difal?.amount != null && order.difal.paid !== true,
  );
  const difalTotal = difalPending.reduce((sum, order) => sum + (order.difal.amount ?? 0), 0);
  const showDifalBanner = ordersReady && difalPending.length > 0;
  const difalResumo = `${formatMoney(difalTotal) ?? ""} a pagar · ${difalPending.length} agendado(s)`;

  const goToFila = () => setView("fila");
  const goToListaTab = (bucketTab: PedidosTab) => () => {
    setView("lista");
    setTab(bucketTab);
  };

  // Loading/Error gate the whole list area (as before); Empty is per-view since Lista keeps its
  // own filter tabs visible even when the current tab's slice is empty.
  let body: ReactNode;
  if (ordersQuery.isPending) {
    body = <LoadingState />;
  } else if (ordersQuery.isError) {
    body = <ErrorState onRetry={() => void ordersQuery.refetch()} />;
  } else if (view === "fila") {
    body = allItems.length === 0 ? <EmptyState /> : <FilaView items={allItems} onOpenOrder={openOrder} />;
  } else if (view === "kanban") {
    body = <KanbanView items={allItems} onOpenOrder={openOrder} />;
  } else {
    body = <ListaView items={allItems} tab={tab} onTabChange={setTab} onOpenOrder={openOrder} />;
  }

  return (
    <section aria-labelledby="pedidos-title" className="mx-auto flex max-w-7xl flex-col gap-[14px]">
      <div className="flex flex-wrap items-center gap-[10px]">
        <h1 id="pedidos-title" className="text-[22px] font-bold tracking-tight text-ink">
          Pedidos
        </h1>
        <button
          type="button"
          disabled
          title="filtro de período em breve"
          className="whitespace-nowrap rounded-pill border border-border px-3 py-1 text-[12.5px] text-muted disabled:cursor-not-allowed"
        >
          período: 7 dias ▾
        </button>
        <button
          type="button"
          disabled
          title="filtro de logística em breve"
          className="whitespace-nowrap rounded-pill border border-border px-3 py-1 text-[12.5px] text-muted disabled:cursor-not-allowed"
        >
          logística: todas ▾
        </button>
        {showDifalBanner ? (
          <button
            type="button"
            onClick={goToFila}
            className="whitespace-nowrap rounded-pill border border-amber bg-amber-soft px-3 py-1 text-[11.5px] font-semibold text-amber"
          >
            DIFAL: {difalResumo}
          </button>
        ) : null}
        <div
          role="tablist"
          aria-label="Alternar visualização de pedidos"
          className="ml-auto inline-flex overflow-hidden rounded-[9px] border border-border text-[12.5px]"
        >
          {viewOptions.map((option) => (
            <button
              key={option.value}
              type="button"
              role="tab"
              aria-selected={view === option.value}
              onClick={() => setView(option.value)}
              className={`px-4 py-1.5 font-medium transition-colors ${
                view === option.value
                  ? "bg-accent-soft text-accent-ink"
                  : "bg-surface text-muted hover:text-ink"
              }`}
            >
              {option.label}
            </button>
          ))}
        </div>
      </div>

      <div aria-label="Indicadores de pedidos" className="flex flex-wrap gap-[10px]">
        <KpiCard label="NOVOS" value={kpiValue("novo")} sub="aguard. pagto" onSelect={goToListaTab("novo")} />
        <KpiCard label="A FATURAR" value={kpiValue("faturar")} onSelect={goToListaTab("faturar")} />
        <KpiCard label="A ENVIAR" value={kpiValue("enviar")} onSelect={goToListaTab("enviar")} />
        <KpiCard label="ENVIADOS" value={kpiValue("enviado")} sub="últimos 7d" onSelect={goToListaTab("enviado")} />
        <KpiCard
          label="DIFAL A PAGAR"
          value={UNKNOWN_KPI_VALUE}
          sub="decomposição pendente"
          onSelect={goToFila}
        />
      </div>

      <div>{body}</div>

      <p className="text-[11.5px] text-faint">
        retorno = valor − comissão − frete − imposto − DIFAL − custo ERP (chips: verde ≥18% · âmbar
        10–18% · vermelho &lt;10%) · SLA = prazo de despacho do ML · status sincronizado do ML/ERP,
        sem edição inline
      </p>

      <PedidoDrawer orderId={openOrderId} onClose={closeOrder} />
    </section>
  );
}
