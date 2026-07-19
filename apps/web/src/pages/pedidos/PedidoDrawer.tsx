import type { OrderBucket, OrderLinkQuality, OrderRead, OrderReadItem } from "@marketplace-central/sdk-runtime";
import { useQuery } from "@tanstack/react-query";
import { DetailDrawer, ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { actionLabelForBucket, formatDateTime, formatMoney } from "./pedidosFormatters";

export interface PedidoDrawerProps {
  orderId: string | null;
  onClose: () => void;
}

// ordersQueryKeys (web-query barrel) has only `.list`, no detail key — a local, stable key
// mirrors the pattern PedidosPage already uses for the summary query (no barrel key exists).
const orderDetailKey = (installationId: string, orderId: string) =>
  ["orders", "detail", installationId, orderId] as const;

// Mirrors PedidosTable's local stateTag idiom (not exported from packages/ui) — kept as a small
// local duplicate here, same as ListingDetailPanel/VinculoDrawer already do for their own drawers.
function stateTag(label: string, className = "bg-slate-100 text-slate-700") {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${className}`}>
      {label}
    </span>
  );
}

const bucketStatusLabels: Record<OrderBucket, string> = {
  novo: "aguard. pagamento",
  faturar: "pago · falta NF",
  enviar: "NF emitida",
  enviado: "enviado",
  cancelado: "cancelado",
};

const linkQualityLabels: Record<OrderLinkQuality, string> = {
  resolved: "vinculado",
  rejected: "rejeitado",
  conflict: "divergente",
  unresolved: "sem vínculo",
  missing: "sem SKU",
};

const linkQualityClasses: Record<OrderLinkQuality, string> = {
  resolved: "bg-emerald-100 text-emerald-800",
  rejected: "bg-red-100 text-red-800",
  conflict: "bg-amber-100 text-amber-800",
  unresolved: "bg-slate-100 text-slate-700",
  missing: "bg-slate-100 text-slate-500",
};

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="flex flex-col gap-2">
      <h4 className="text-[10px] font-semibold uppercase tracking-wide text-faint">{title}</h4>
      {children}
    </section>
  );
}

function ItemRow({ item }: { item: OrderReadItem }) {
  const title = item.title || item.seller_sku || item.provider_item_id;
  const lineTotal = item.unit_price === undefined ? null : formatMoney(item.quantity * item.unit_price);
  const custo = item.custo_unitario === undefined ? null : formatMoney(item.custo_unitario);
  return (
    <div className="flex flex-col gap-1 border-b border-border-2 pb-2 last:border-b-0 last:pb-0">
      <div className="flex items-baseline gap-2">
        <span className="min-w-0 flex-1 truncate text-sm">{title}</span>
        <span className="flex-none font-mono text-xs">
          {item.quantity}× {lineTotal ?? <UnknownValue />}
        </span>
      </div>
      <div className="flex flex-wrap items-center gap-2 text-[11px] text-faint">
        {stateTag(linkQualityLabels[item.link_quality], linkQualityClasses[item.link_quality])}
        <span>custo unit.: {custo ?? <UnknownValue hint="sem custo ERP no vínculo" />}</span>
        {item.internal_product_id !== undefined ? <span>CODPROD {item.internal_product_id}</span> : null}
      </div>
    </div>
  );
}

function ItemsSection({ order }: { order: OrderRead }) {
  return (
    <Section title="Itens">
      <div className="flex flex-col gap-2">
        {order.items.length === 0 ? (
          <p className="text-xs text-faint">Sem itens.</p>
        ) : (
          order.items.map((item) => (
            <ItemRow key={`${item.provider_item_id}-${item.provider_variation_id ?? ""}`} item={item} />
          ))
        )}
        {order.componentes_desconhecidos && order.componentes_desconhecidos.length > 0
          ? stateTag("custo incompleto", "bg-amber-100 text-amber-800")
          : null}
      </div>
    </Section>
  );
}

function DecomposicaoSection() {
  // Gated (slice C, IC-04 open): OrderRead has no retorno/margem/decomposição/DIFAL fields —
  // an honest "não disponível" block, never fabricated numbers (ADR-17).
  return (
    <Section title="Decomposição + DIFAL">
      <div className="rounded-lg border border-dashed border-border bg-surface-2 p-3 text-xs text-muted">
        <p>
          Retorno líquido, margem e decomposição de custos (comissão, frete, imposto, DIFAL) ainda
          não estão disponíveis nesta instalação.
        </p>
        <dl className="mt-2 space-y-1">
          <div className="flex items-center justify-between gap-3">
            <dt>Retorno líquido</dt>
            <dd>
              <UnknownValue hint="decomposição de custos ainda não disponível (IC-04)" />
            </dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt>Margem</dt>
            <dd>
              <UnknownValue hint="decomposição de custos ainda não disponível (IC-04)" />
            </dd>
          </div>
          <div className="flex items-center justify-between gap-3">
            <dt>DIFAL</dt>
            <dd>
              <UnknownValue hint="DIFAL ainda não decomposto" />
            </dd>
          </div>
        </dl>
      </div>
    </Section>
  );
}

interface TimelineEvent {
  label: string;
  when: string | null;
}

// Client-derived from present timestamps only (ruling-6) — no fabricated dates; events whose
// timestamp is absent are omitted entirely rather than shown with a placeholder date. "Enviado"
// has no timestamp field on OrderRastreio, so it renders with status only, no invented "quando".
function buildTimeline(order: OrderRead): TimelineEvent[] {
  const events: TimelineEvent[] = [];
  if (order.provider_created_at) {
    events.push({ label: "Recebido", when: formatDateTime(order.provider_created_at) });
  }
  if (order.provider_updated_at) {
    events.push({ label: "Atualizado", when: formatDateTime(order.provider_updated_at) });
  }
  if (order.rastreio) {
    events.push({ label: `Enviado · ${order.rastreio.status}`, when: null });
  }
  if (order.provider_closed_at) {
    events.push({ label: "Fechado", when: formatDateTime(order.provider_closed_at) });
  }
  return events;
}

function TimelineSection({ order }: { order: OrderRead }) {
  const events = buildTimeline(order);
  return (
    <Section title="Linha do tempo · ML → interno">
      {events.length === 0 ? (
        <p className="text-xs text-faint">Sem eventos com data disponível.</p>
      ) : (
        <ol className="flex flex-col gap-1.5 text-xs">
          {events.map((event) => (
            <li key={event.label} className="flex items-baseline gap-2">
              <span className="flex-1">{event.label}</span>
              <span className="flex-none font-mono text-[10.5px] text-faint">{event.when ?? ""}</span>
            </li>
          ))}
        </ol>
      )}
    </Section>
  );
}

function FactsSection({ order }: { order: OrderRead }) {
  const local = [order.buyer?.city, order.buyer?.uf].filter(Boolean).join("/");
  return (
    <section className="flex flex-col gap-2 border-t border-border-2 pt-3 text-xs text-muted">
      <div className="flex items-start justify-between gap-3">
        <span>Nota fiscal</span>
        {/* nf_state is the vínculo 'linked' marker on OrderRead, not a real NF number — rendered
            honestly, never as a fabricated NF number/deep-link (ruling-6). */}
        <span className="text-right text-ink">
          {order.nf_state ? `vínculo: ${order.nf_state}` : <UnknownValue hint="ainda não emitida" />}
        </span>
      </div>
      <div className="flex items-start justify-between gap-3">
        <span>Rastreio</span>
        <span className="text-right font-mono text-[11px] text-ink">
          {order.rastreio ? `${order.rastreio.shipment_id} · ${order.rastreio.status}` : <UnknownValue />}
        </span>
      </div>
      <div className="flex items-start justify-between gap-3">
        <span>Destino</span>
        <span className="text-right text-ink">{order.destino_uf ?? <UnknownValue />}</span>
      </div>
      <div className="flex items-start justify-between gap-3">
        <span>Código de rastreio</span>
        {/* No carrier tracking code on OrderRead — honest unknown, never fabricated (ADR-17). */}
        <span className="text-right text-ink">
          <UnknownValue hint="código da transportadora não disponível no backend" />
        </span>
      </div>
      <div className="flex items-start justify-between gap-3">
        <span>Comprador</span>
        <span className="text-right text-ink">
          {order.buyer?.display ?? <UnknownValue />}
          {order.buyer ? (
            <>
              <br />
              <span className="text-[11px] text-faint">{local || <UnknownValue />}</span>
            </>
          ) : null}
        </span>
      </div>
      <div className="flex items-start justify-between gap-3">
        <span>Documento</span>
        {/* No buyer document (CPF/CNPJ) on OrderRead — honest unknown, never fabricated (ruling-6). */}
        <span className="text-right text-ink">
          <UnknownValue hint="documento do comprador não disponível no backend" />
        </span>
      </div>
    </section>
  );
}

function DrawerBody({ order }: { order: OrderRead }) {
  return (
    <div className="flex flex-col gap-4">
      <ItemsSection order={order} />
      <DecomposicaoSection />
      <TimelineSection order={order} />
      <FactsSection order={order} />
    </div>
  );
}

// Design's footer buttons per bucket (Faturar via ERP / Etiqueta / Marcar enviado / DIFAL
// agendar / Devolução…) — ALL disabled here (READ-ONLY RICH bar, HUB ruling D-57): no working
// mutations, no write SDK calls this slice.
function DrawerActions({ bucket }: { bucket: OrderBucket }) {
  const primaryLabel = actionLabelForBucket(bucket);
  const buttons: string[] = [];
  if (primaryLabel === "Faturar") buttons.push("Faturar via ERP");
  if (primaryLabel === "Etiqueta") buttons.push("Etiqueta");
  if (bucket === "enviar") buttons.push("Marcar enviado");
  buttons.push("DIFAL agendar", "Devolução…");
  return (
    <div className="flex flex-wrap gap-2">
      {buttons.map((label) => (
        <button
          key={label}
          type="button"
          disabled
          title="disponível em breve"
          className="flex-1 rounded-md border border-border bg-surface-2 px-3 py-2 text-xs font-semibold text-muted disabled:cursor-not-allowed disabled:opacity-60"
        >
          {label}
        </button>
      ))}
    </div>
  );
}

export function PedidoDrawer({ orderId, onClose }: PedidoDrawerProps) {
  const client = useClient();
  const { installationId } = useInstallation();

  const query = useQuery({
    queryKey: orderDetailKey(installationId, orderId ?? ""),
    queryFn: () => client.getOrder(installationId, orderId as string),
    enabled: orderId !== null,
    staleTime: QUERY_STALE_TIME.orders,
  });

  if (orderId === null) return null;

  const order = query.data;
  const title = order?.provider_code || orderId;
  const subtitle = order ? bucketStatusLabels[order.bucket] : undefined;

  return (
    <DetailDrawer
      open
      onClose={onClose}
      closeLabel="Fechar detalhe do pedido"
      title={title}
      subtitle={subtitle}
      width={372}
      actions={order ? <DrawerActions bucket={order.bucket} /> : undefined}
    >
      {query.isPending ? (
        <LoadingState />
      ) : query.isError ? (
        <ErrorState onRetry={() => void query.refetch()} />
      ) : order ? (
        <DrawerBody order={order} />
      ) : null}
    </DetailDrawer>
  );
}
