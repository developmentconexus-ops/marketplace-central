import type { OrderBucket, OrderLinkQuality, OrderRead, OrderReadItem } from "@marketplace-central/sdk-runtime";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { DetailDrawer, ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { actionLabelForBucket, formatDateTime, formatMoney, formatPercent } from "./pedidosFormatters";

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

// order.fulfillment is a stable provider code (domain layer, untouched). This map is a
// presentation-only PT-BR label for the codes we've actually seen; a code not yet mapped
// falls back to the raw code itself — it's a known fact, just not yet given a pretty label
// (only a genuinely absent value renders UnknownValue, ADR-17).
const fulfillmentLabels: Record<string, string> = {
  fulfillment: "Fulfillment (FBA-like)",
  drop_off: "Coleta no ponto (drop-off)",
  xd_drop_off: "Cross-docking (drop-off)",
  self_service: "Envio próprio (flex)",
};

function formatFulfillment(code: string | null): string | null {
  if (!code) return null;
  return fulfillmentLabels[code] ?? code;
}

// order.cancellation_detail is stored as "<requested_by>:<code>" (domain layer, untouched).
// These maps are presentation-only PT-BR labels for the values we've actually seen; anything
// unmapped falls back to the raw token rather than hiding a real fact (ADR-17).
const cancellationRequesterLabels: Record<string, string> = {
  buyer: "comprador",
  seller: "vendedor",
  meli: "Mercado Livre",
};

const cancellationCodeLabels: Record<string, string> = {
  feedback_unavailable_product: "produto indisponível",
  mediations: "mediação",
  buyer_cancel_express: "cancelamento expresso",
};

function formatCancellationDetail(detail: string | null): string | null {
  if (!detail) return null;
  const sep = detail.indexOf(":");
  if (sep === -1) return detail;
  const requestedBy = detail.slice(0, sep);
  const code = detail.slice(sep + 1);
  const requesterLabel = cancellationRequesterLabels[requestedBy] ?? requestedBy;
  const codeLabel = cancellationCodeLabels[code] ?? code;
  return `${requesterLabel} · ${codeLabel}`;
}

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
        <span>
          {/* "≈" is not decoration: the ERP source keeps one snapshot per
              product, so a cost read for a past sale is the closest observation,
              not the cost on that date. The instant is spelled out rather than
              letting the number pass as exact. */}
          custo unit.: {custo === null ? <UnknownValue hint="sem custo ERP no vínculo" /> : item.custo_aproximado ? `≈ ${custo}` : custo}
        </span>
        {custo !== null && item.custo_aproximado ? (
          <span>
            snapshot ERP de {formatDateTime(item.custo_observado_em) ?? <UnknownValue />}
          </span>
        ) : null}
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

// Definition-list row: value → formatted text, null/undefined → UnknownValue (ADR-17). Never
// hardcodes "—"; the hint explains why a component is unknown today (F01-C1 honest-empty) so
// the same row lights up with a real number once the hub wires the decomposer (C2), no UI change.
function DecompRow({
  label,
  value,
  hint,
}: {
  label: string;
  value: string | null;
  hint?: string;
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <dt>{label}</dt>
      <dd className="font-mono text-[11px]">{value ?? <UnknownValue hint={hint} />}</dd>
    </div>
  );
}

// Restituição ST is the one row in Decomposição that ADDS to margem instead of subtracting
// (BuildProfitability, D-41) — every other row here is a deduction. A bare formatMoney would
// render a positive restituição identically to a cost line, so a "+ " prefix marks a REALIZED
// credit visibly. An explicit 0 (e.g. intra-MG, where restituição never applies) is a real fact,
// not a credit that fired, so it stays a plain "R$ 0,00" — no "+".
function formatCredit(value: number | null | undefined): string | null {
  const formatted = formatMoney(value);
  if (formatted === null || value === null || value === undefined || value <= 0) return formatted;
  return `+ ${formatted}`;
}

// Real-ready wiring (F02-S6): every value here is read FROM order.decomposicao/difal/
// retorno_liquido/margem_pct (F01-C1, additive on OrderRead). Comissão/Taxa fixa/Frete/Tarifa
// Full/Custo/Margem are still honestly null (the cost decomposer isn't wired, hub C2) and render
// UnknownValue via DecompRow/formatMoney — never a hardcoded "—" string. ICMS de saída, DIFAL,
// PIS/COFINS and Restituição ST ARE fed now, by the D-41 per-item tax engine (P2.b T5) —
// Imposto (the old D-38 field) is permanently retired by that same change: it stays in the
// payload for old consumers but this UI never reads it as a live number again.
function DecomposicaoSection({ order }: { order: OrderRead }) {
  const { decomposicao, difal } = order;
  const pending = decomposicao.componentes_desconhecidos.length > 0;
  const difalHint = "DIFAL ainda não decomposto (hub C2)";
  const custoHint = "decomposição de custos ainda não disponível (hub C2)";
  const impostoHint =
    "substituído por ICMS de saída + DIFAL (D-41) — este campo não é mais calculado por esta via";
  // ICMS saída/DIFAL share one resolution path (pricing/domain.TaxesForItem): both go unknown
  // together when a product has no vínculo interno / fiscal ERP snapshot, when the ICMS matrix
  // cell for (destino, grupo_icms) is ausente/ambígua, or when the destino UF has no linha na
  // tabela de alíquota interna (D-37).
  const icmsCellHint =
    "produto sem vínculo interno, ou célula ICMS ausente/ambígua para este destino/grupo, ou UF sem alíquota interna cadastrada (matriz D-37)";
  const pisCofinsHint =
    "mesma lacuna do ICMS de saída — produto sem vínculo interno, célula ICMS não resolvida, ou UF sem alíquota interna (matriz D-37)";
  // Restituição ST only depends on UF destino + o fato do produto (não passa pela célula
  // codtrib/ambígua) — a lacuna aqui é sempre o item sem vínculo interno ou sem snapshot
  // fiscal ERP (products_mirror), nunca a matriz.
  const restituicaoHint =
    "produto sem vínculo interno, ou sem snapshot fiscal ERP (products_mirror), para algum item do pedido";

  return (
    <Section title="Decomposição + DIFAL">
      <div className="rounded-lg border border-border bg-surface-2 p-3 text-xs text-muted">
        {pending ? (
          <p className="mb-2 text-[11px] text-faint">
            Pendente: {decomposicao.componentes_desconhecidos.join(", ")} — estes componentes mostram
            "—" até a decomposição ser calculada.
          </p>
        ) : null}
        <dl className="space-y-1">
          <DecompRow label="Comissão" value={formatMoney(decomposicao.comissao)} hint={custoHint} />
          <DecompRow label="Taxa fixa" value={formatMoney(decomposicao.taxa_fixa)} hint={custoHint} />
          <DecompRow label="Frete" value={formatMoney(decomposicao.frete)} hint={custoHint} />
          <DecompRow label="Imposto" value={formatMoney(decomposicao.imposto)} hint={impostoHint} />
          <DecompRow label="ICMS saída" value={formatMoney(decomposicao.icms_saida)} hint={icmsCellHint} />
          <DecompRow label="DIFAL" value={formatMoney(decomposicao.difal)} hint={icmsCellHint} />
          <DecompRow label="PIS/COFINS" value={formatMoney(decomposicao.pis_cofins)} hint={pisCofinsHint} />
          <DecompRow
            label="Restituição ST"
            value={formatCredit(decomposicao.restituicao_st)}
            hint={restituicaoHint}
          />
          <DecompRow label="Tarifa Full" value={formatMoney(decomposicao.tarifa_full)} hint={custoHint} />
          <DecompRow label="Custo" value={formatMoney(decomposicao.custo)} hint={custoHint} />
          <div className="my-1 border-t border-border-2" />
          <DecompRow label="Margem valor" value={formatMoney(decomposicao.margem_valor)} hint={custoHint} />
          <DecompRow label="Margem %" value={formatPercent(decomposicao.margem_pct)} hint={custoHint} />
          <DecompRow label="Retorno líquido" value={formatMoney(order.retorno_liquido)} hint={custoHint} />
        </dl>
      </div>
      <div className="rounded-lg border border-border bg-surface-2 p-3 text-xs text-muted">
        <h5 className="mb-2 text-[10px] font-semibold uppercase tracking-wide text-faint">DIFAL</h5>
        <dl className="space-y-1">
          <DecompRow label="Valor" value={formatMoney(difal.amount)} hint={difalHint} />
          <DecompRow label="Rota" value={difal.uf_route} hint="rota UF ainda não disponível (hub C2)" />
          <DecompRow label="Vencimento" value={formatDateTime(difal.due_date)} hint={difalHint} />
          <div className="flex items-center justify-between gap-3">
            <dt>Pago</dt>
            <dd className="font-mono text-[11px]">
              {difal.paid == null ? <UnknownValue hint={difalHint} /> : difal.paid ? "sim" : "não"}
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
    const sub = order.rastreio.substatus ? ` (${order.rastreio.substatus})` : "";
    events.push({ label: `Enviado · ${order.rastreio.status}${sub}`, when: null });
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

// FactRow: label → value; null/undefined → honest UnknownValue (ADR-17). Never a
// hardcoded "—" string.
function FactRow({
  label,
  value,
  hint,
  mono,
}: {
  label: string;
  value: React.ReactNode;
  hint?: string;
  mono?: boolean;
}) {
  const empty = value === null || value === undefined || value === "";
  return (
    <div className="flex items-start justify-between gap-3">
      <span>{label}</span>
      <span className={`text-right text-ink${mono ? " font-mono text-[11px]" : ""}`}>
        {empty ? <UnknownValue hint={hint} /> : value}
      </span>
    </div>
  );
}

function FactsSection({ order }: { order: OrderRead }) {
  const local = [order.buyer?.city, order.buyer?.uf].filter(Boolean).join("/");
  const destino = [order.destino_uf, order.destino_cep].filter(Boolean).join(" · ");
  const carrier = order.rastreio?.transportadora;
  const trackUrl = order.rastreio?.url_rastreio;
  const frete = order.frete_real;
  return (
    <section className="flex flex-col gap-2 border-t border-border-2 pt-3 text-xs text-muted">
      <FactRow
        label="Nota fiscal"
        /* nf_state is the vínculo 'linked' marker on OrderRead, not a real NF number — rendered
           honestly, never as a fabricated NF number/deep-link (ruling-6). */
        value={order.nf_state ? `vínculo: ${order.nf_state}` : null}
        hint="ainda não emitida"
      />
      <FactRow
        label="Rastreio"
        mono
        value={
          order.rastreio
            ? `${order.rastreio.shipment_id} · ${order.rastreio.status}${order.rastreio.substatus ? ` · ${order.rastreio.substatus}` : ""}`
            : null
        }
      />
      <FactRow
        label="Modalidade de envio"
        value={formatFulfillment(order.fulfillment)}
        hint="modalidade de envio ainda não disponível"
      />
      <FactRow
        label="Transportadora"
        value={
          carrier ? (
            <>
              <span>{carrier}</span>
              {trackUrl ? (
                <>
                  {" · "}
                  <a
                    href={trackUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="underline underline-offset-2 hover:text-ink"
                  >
                    rastrear
                  </a>
                </>
              ) : null}
            </>
          ) : null
        }
        hint="transportadora ainda não atribuída"
      />
      <FactRow
        label="Motivo do cancelamento"
        value={formatCancellationDetail(order.cancellation_detail)}
        hint="sem cancelamento"
      />
      <FactRow label="Destino" value={destino || null} />
      <FactRow label="Destinatário" value={order.destinatario ?? null} hint="destinatário ainda não disponível" />
      <FactRow
        label="Frete real"
        mono
        value={
          frete && (frete.bruto != null || frete.receiver != null || frete.sender != null) ? (
            <>
              <span>{formatMoney(frete.bruto ?? null) ?? <UnknownValue />}</span>
              <br />
              <span className="text-[11px] text-faint">
                receiver {formatMoney(frete.receiver ?? null) ?? "—"} · sender{" "}
                {formatMoney(frete.sender ?? null) ?? "—"}
              </span>
            </>
          ) : null
        }
        hint="custos de frete não disponíveis"
      />
      <FactRow
        label="Comprador"
        value={
          order.buyer ? (
            <>
              {order.buyer.display || <UnknownValue />}
              <br />
              <span className="text-[11px] text-faint">{local || <UnknownValue />}</span>
            </>
          ) : null
        }
      />
    </section>
  );
}

// Composes the buyer billing address into a single honest line. Absent parts are
// dropped; an all-absent address yields null (caller renders UnknownValue). Never
// fabricates a value (ADR-17).
function formatEndereco(end: NonNullable<OrderRead["comprador_fiscal"]>["endereco"]): string | null {
  if (!end) return null;
  const linha1 = [end.logradouro, end.numero].filter(Boolean).join(", ");
  const linha2 = [end.cidade, end.uf_codigo].filter(Boolean).join("/");
  const parts = [linha1, linha2, end.cep, end.pais].filter(Boolean);
  return parts.length > 0 ? parts.join(" · ") : null;
}

// Buyer fiscal identity for ERP registration (name, opaque document, billing
// address). Additive comprador_fiscal block on OrderRead — absent (buyer without
// billing data / masked until payment) renders honest "—", never fabricated
// (ADR-17). doc_tipo is rendered VERBATIM (opaque — never mapped to a CPF/CNPJ
// enum). The document number is rendered for the operator only; it is never
// logged (LGPD).
function CompradorFiscalSection({ order }: { order: OrderRead }) {
  const cf = order.comprador_fiscal;
  const doc = cf?.doc_tipo || cf?.doc_numero ? [cf?.doc_tipo, cf?.doc_numero].filter(Boolean).join(" ") : null;
  const endereco = formatEndereco(cf?.endereco);
  return (
    <Section title="Comprador · fiscal (ERP)">
      <dl className="flex flex-col gap-2 text-xs text-muted">
        <FactRow label="Nome/Razão" value={cf?.nome ?? null} hint="nome fiscal não disponível" />
        <FactRow label="Documento" value={doc} hint="documento não disponível" mono />
        <FactRow label="Endereço" value={endereco} hint="endereço fiscal não disponível" />
      </dl>
    </Section>
  );
}

function DrawerBody({ order }: { order: OrderRead }) {
  return (
    <div className="flex flex-col gap-4">
      <ItemsSection order={order} />
      <DecomposicaoSection order={order} />
      <CompradorFiscalSection order={order} />
      <TimelineSection order={order} />
      <FactsSection order={order} />
    </div>
  );
}

// Footer buttons per bucket. Only "Foi faturado" is a working mutation — it records our own
// faturado_at fact (goal B), an OUR-DB write that moves the order from "A FATURAR" to "A ENVIAR";
// it is NOT a Mercado Livre write. The provider/ERP actions (Faturar via ERP / Etiqueta / Marcar
// enviado / DIFAL / Devolução…) stay inert per HUB ruling D-57 — no provider-mutating calls.
function DrawerActions({
  bucket,
  onFaturado,
  isFaturando,
  faturarFailed,
}: {
  bucket: OrderBucket;
  onFaturado: () => void;
  isFaturando: boolean;
  faturarFailed: boolean;
}) {
  const primaryLabel = actionLabelForBucket(bucket);
  const inert: string[] = [];
  if (primaryLabel === "Faturar") inert.push("Faturar via ERP");
  if (primaryLabel === "Etiqueta") inert.push("Etiqueta");
  if (bucket === "enviar") inert.push("Marcar enviado");
  inert.push("DIFAL agendar", "Devolução…");
  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-wrap gap-2">
        {bucket === "faturar" ? (
          <button
            type="button"
            onClick={onFaturado}
            disabled={isFaturando}
            className="flex-1 rounded-md border border-accent bg-accent-soft px-3 py-2 text-xs font-semibold text-accent-ink transition-colors hover:bg-accent-soft disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isFaturando ? "Marcando…" : "Foi faturado"}
          </button>
        ) : null}
        {inert.map((label) => (
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
      {faturarFailed ? (
        <small role="alert" className="text-warn">
          Falha ao marcar como faturado.
        </small>
      ) : null}
    </div>
  );
}

export function PedidoDrawer({ orderId, onClose }: PedidoDrawerProps) {
  const client = useClient();
  const { installationId } = useInstallation();
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: orderDetailKey(installationId, orderId ?? ""),
    queryFn: () => client.getOrder(installationId, orderId as string),
    enabled: orderId !== null,
    staleTime: QUERY_STALE_TIME.orders,
  });

  // "Foi faturado" records our own faturado_at fact (OUR-DB only, never an ML write); on success
  // we invalidate the whole "orders" namespace so this drawer's bucket AND the list/KPI counts
  // (PedidosPage) re-derive — the order leaves "A FATURAR" and lands in "A ENVIAR".
  const faturarMutation = useMutation({
    mutationFn: () => client.markOrderFaturado(installationId, orderId as string),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["orders"] }),
  });

  if (orderId === null) return null;

  const order = query.data;
  // The drawer title is the order number. provider_code is the channel slug ("mercado_livre")
  // on real ML data — never the order label (parity with the list views, which all render
  // provider_order_id). orderId is already the provider_order_id the row passed in.
  const title = order?.provider_order_id || orderId;
  const subtitle = order ? bucketStatusLabels[order.bucket] : undefined;

  return (
    <DetailDrawer
      open
      onClose={onClose}
      closeLabel="Fechar detalhe do pedido"
      title={title}
      subtitle={subtitle}
      width={372}
      actions={
        order ? (
          <DrawerActions
            bucket={order.bucket}
            onFaturado={() => faturarMutation.mutate()}
            isFaturando={faturarMutation.isPending}
            faturarFailed={faturarMutation.isError}
          />
        ) : undefined
      }
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
