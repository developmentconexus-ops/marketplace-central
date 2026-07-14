import { useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, StatCard, SurfaceCard } from "@marketplace-central/ui";
import type {
  IntegrationInstallation,
  MarketplaceOrder,
  MarketplaceOrderItem,
  OrderLinkQuality,
  ProfitSnapshotFlag,
  ProfitSnapshotQuality,
  ProfitabilityInputKind,
  ProfitabilityInputScope,
  ProfitabilityManualAdjustment,
  ProfitabilityManualAdjustmentCategory,
  ProfitabilityMarginInput,
  ProfitabilityProfitSnapshot,
} from "@marketplace-central/sdk-runtime";
import {
  FreshnessIndicator,
  profitabilityQueryKeys,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface OrdersClient extends RefreshableClient {
  listIntegrationInstallations: () => Promise<{ items: IntegrationInstallation[] }>;
  importMarketplaceOrders: (req: { installation_id: string; limit?: number }) => Promise<unknown>;
  listMarketplaceOrders: (installationId: string, limit?: number) => Promise<{ items: MarketplaceOrder[] }>;
  importProfitabilityMarginInputs: (req: { installation_id: string; limit?: number }) => Promise<unknown>;
  listProfitabilityMarginInputs: (installationId: string, limit?: number) => Promise<{ items: ProfitabilityMarginInput[]; as_of?: string }>;
  createProfitabilityManualAdjustment: (req: {
    installation_id: string;
    idempotency_key: string;
    provider_order_id: string;
    provider_item_id?: string;
    provider_variation_id?: string;
    scope?: ProfitabilityInputScope;
    category?: ProfitabilityManualAdjustmentCategory;
    amount: number;
    currency?: string;
    reason: string;
    actor: {
      actor_type: string;
      actor_id: string;
      actor_name?: string;
    };
  }) => Promise<ProfitabilityManualAdjustment>;
  listProfitabilityManualAdjustments: (installationId: string, limit?: number) => Promise<{ items: ProfitabilityManualAdjustment[] }>;
  calculateProfitabilityProfitSnapshots: (req: { installation_id: string; limit?: number }) => Promise<unknown>;
  listProfitabilityProfitSnapshots: (installationId: string, limit?: number) => Promise<{ items: ProfitabilityProfitSnapshot[] }>;
}

export interface OrdersPageProps {
  client: OrdersClient;
  operator?: Readonly<ProfitabilityActor>;
}

export interface ProfitabilityActor {
  actor_type: string;
  actor_id: string;
  actor_name?: string;
}

type LoadState = "loading" | "ready" | "error";
type QualityFilter = "all" | ProfitSnapshotQuality | "not_calculated";

type AdjustmentFormState = {
  scope: ProfitabilityInputScope;
  category: ProfitabilityManualAdjustmentCategory;
  amount: string;
  reason: string;
  itemIdentity: string;
};

type PendingAdjustmentSubmission = {
  fingerprint: string;
  idempotencyKey: string;
};

const linkSeverity: Record<OrderLinkQuality, number> = {
  resolved: 0,
  rejected: 1,
  missing: 2,
  unresolved: 3,
  conflict: 4,
};

function normalizeError(error: unknown, fallback: string): string {
  if (error && typeof error === "object") {
    const structured = error as { error?: { message?: string }; message?: string };
    return structured.error?.message ?? structured.message ?? fallback;
  }
  if (typeof error === "string" && error.trim()) {
    return error;
  }
  return fallback;
}

function createManualAdjustmentIdempotencyKey(): string {
  return globalThis.crypto?.randomUUID?.() ?? `manual-adjustment-${Date.now()}-${Math.random().toString(36).slice(2)}`;
}

function manualAdjustmentFingerprint(input: {
  installationID: string;
  providerOrderID: string;
  providerItemID?: string;
  providerVariationID?: string;
  scope: ProfitabilityInputScope;
  category: ProfitabilityManualAdjustmentCategory;
  amount: number;
  currency: string;
  reason: string;
  actor: ProfitabilityActor;
}): string {
  return JSON.stringify(input);
}

function money(amount?: number | null, currency = "BRL"): string {
  if (amount == null) {
    return "—";
  }
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

function percent(value?: number | null): string {
  if (value == null) {
    return "—";
  }
  return `${value.toFixed(2)}%`;
}

function formatDate(value?: string | null): string {
  if (!value) {
    return "-";
  }
  return new Date(value).toLocaleString("en-US");
}

function qualityLabel(quality: QualityFilter): string {
  switch (quality) {
    case "complete":
      return "Complete";
    case "incomplete":
      return "Incomplete";
    case "negative_margin":
      return "Negative margin";
    case "not_realized":
      return "Not realized";
    case "not_calculated":
      return "Not calculated";
    default:
      return "All qualities";
  }
}

function qualityTone(quality: QualityFilter): string {
  switch (quality) {
    case "complete":
      return "bg-emerald-100 text-emerald-800";
    case "incomplete":
      return "bg-amber-100 text-amber-800";
    case "negative_margin":
      return "bg-rose-100 text-rose-800";
    case "not_realized":
      return "bg-slate-200 text-slate-800";
    case "not_calculated":
      return "bg-slate-100 text-slate-700";
    default:
      return "bg-slate-100 text-slate-700";
  }
}

function linkLabel(quality: OrderLinkQuality): string {
  switch (quality) {
    case "resolved":
      return "Resolved";
    case "rejected":
      return "Rejected";
    case "missing":
      return "Missing";
    case "unresolved":
      return "Unresolved";
    case "conflict":
      return "Conflict";
    default:
      return quality;
  }
}

function itemIdentity(item: { provider_item_id?: string; provider_variation_id?: string }): string {
  return `${item.provider_item_id ?? ""}::${item.provider_variation_id ?? ""}`;
}

function snapshotIdentity(snapshot: {
  provider_order_id: string;
  provider_item_id?: string;
  provider_variation_id?: string;
  scope: ProfitabilityInputScope;
}): string {
  return `${snapshot.provider_order_id}::${snapshot.scope}::${snapshot.provider_item_id ?? ""}::${snapshot.provider_variation_id ?? ""}`;
}

function adjustmentIdentity(adjustment: {
  provider_order_id: string;
  scope: ProfitabilityInputScope;
  provider_item_id?: string;
  provider_variation_id?: string;
}): string {
  return `${adjustment.provider_order_id}::${adjustment.scope}::${adjustment.provider_item_id ?? ""}::${adjustment.provider_variation_id ?? ""}`;
}

function deriveOrderLinkQuality(items: MarketplaceOrderItem[]): OrderLinkQuality {
  if (items.length === 0) {
    return "missing";
  }
  return [...items]
    .sort((left, right) => linkSeverity[right.link_quality] - linkSeverity[left.link_quality])[0]
    .link_quality;
}

function defaultAdjustmentForm(order?: MarketplaceOrder): AdjustmentFormState {
  return {
    scope: "order",
    category: "freight",
    amount: "",
    reason: "",
    itemIdentity: order?.items[0] ? itemIdentity(order.items[0]) : "",
  };
}

function flagLabel(flag: ProfitSnapshotFlag): string {
  switch (flag) {
    case "order_cancelled":
      return "Order cancelled";
    case "order_state_unknown":
      return "Order state unknown";
    case "negative_margin":
      return "Negative margin";
    default: {
      const label = flag.replace(/_/g, " ");
      return label.charAt(0).toUpperCase() + label.slice(1);
    }
  }
}

function isOperationalFlag(flag: ProfitSnapshotFlag): boolean {
  return flag === "order_cancelled" || flag === "order_state_unknown";
}

function isMissingInputFlag(flag: ProfitSnapshotFlag): boolean {
  return flag.startsWith("missing_");
}

function realizationLabel(snapshot?: ProfitabilityProfitSnapshot | null): string | null {
  switch (snapshot?.realization_state) {
    case "not_realized":
      return "Order not realized";
    case "unknown":
      return "Order state unknown";
    default:
      return null;
  }
}

function inputLabel(kind: ProfitabilityInputKind): string {
  return kind.replace(/_/g, " ");
}

export function OrdersPage({ client, operator }: OrdersPageProps) {
  const queryClient = useQueryClient();
  const pendingAdjustmentSubmission = useRef<PendingAdjustmentSubmission | null>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const [installations, setInstallations] = useState<IntegrationInstallation[]>([]);
  const [orders, setOrders] = useState<MarketplaceOrder[]>([]);
  const [inputs, setInputs] = useState<ProfitabilityMarginInput[]>([]);
  const [adjustments, setAdjustments] = useState<ProfitabilityManualAdjustment[]>([]);
  const [snapshots, setSnapshots] = useState<ProfitabilityProfitSnapshot[]>([]);
  const [state, setState] = useState<LoadState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [pendingAction, setPendingAction] = useState<string | null>(null);
  const [form, setForm] = useState<AdjustmentFormState>(defaultAdjustmentForm());

  const selectedInstallationID = searchParams.get("installation") ?? "";
  const selectedOrderID = searchParams.get("order") ?? "";
  const selectedQuality = (searchParams.get("quality") as QualityFilter | null) ?? "all";

  const profitabilityQuery = useQuery({
    queryKey: profitabilityQueryKeys.marginInputs(selectedInstallationID),
    queryFn: () => client.listProfitabilityMarginInputs(selectedInstallationID, 200),
    staleTime: QUERY_STALE_TIME.pricecost,
    enabled: Boolean(selectedInstallationID),
  });

  useEffect(() => {
    if (profitabilityQuery.data) {
      setInputs(profitabilityQuery.data.items);
    }
    if (profitabilityQuery.error) {
      setError(normalizeError(profitabilityQuery.error, "Failed to load profitability inputs."));
      setState("error");
    }
  }, [profitabilityQuery.data, profitabilityQuery.error]);

  useEffect(() => {
    let cancelled = false;

    async function loadInstallations() {
      try {
        const result = await client.listIntegrationInstallations();
        if (cancelled) {
          return;
        }
        setInstallations(result.items);
        if (!selectedInstallationID && result.items[0]) {
          setSearchParams({ installation: result.items[0].installation_id }, { replace: true });
        }
      } catch (loadError) {
        if (cancelled) {
          return;
        }
        setError(normalizeError(loadError, "Failed to load installations."));
        setState("error");
      }
    }

    void loadInstallations();
    return () => {
      cancelled = true;
    };
  }, [client, selectedInstallationID, setSearchParams]);

  async function loadWorkspace(installationID: string) {
    const [ordersResult, adjustmentsResult, snapshotsResult] = await Promise.all([
      client.listMarketplaceOrders(installationID, 50),
      client.listProfitabilityManualAdjustments(installationID, 200),
      client.listProfitabilityProfitSnapshots(installationID, 200),
    ]);

    setOrders(ordersResult.items);
    setAdjustments(adjustmentsResult.items);
    setSnapshots(snapshotsResult.items);
  }

  useEffect(() => {
    if (!selectedInstallationID) {
      return;
    }

    let cancelled = false;

    async function loadData() {
      setState("loading");
      setError(null);
      try {
        await loadWorkspace(selectedInstallationID);
        if (cancelled) {
          return;
        }
        setState("ready");
      } catch (loadError) {
        if (cancelled) {
          return;
        }
        setError(normalizeError(loadError, "Failed to load orders workspace."));
        setState("error");
      }
    }

    void loadData();
    return () => {
      cancelled = true;
    };
  }, [client, selectedInstallationID]);

  const orderSnapshots = useMemo(() => {
    const map = new Map<string, ProfitabilityProfitSnapshot>();
    for (const snapshot of snapshots) {
      if (snapshot.scope === "order") {
        map.set(snapshot.provider_order_id, snapshot);
      }
    }
    return map;
  }, [snapshots]);

  const itemSnapshots = useMemo(() => {
    const map = new Map<string, ProfitabilityProfitSnapshot>();
    for (const snapshot of snapshots) {
      if (snapshot.scope === "item") {
        map.set(snapshotIdentity(snapshot), snapshot);
      }
    }
    return map;
  }, [snapshots]);

  const filteredOrders = useMemo(() => {
    return orders.filter((order) => {
      if (selectedQuality === "all") {
        return true;
      }
      const snapshot = orderSnapshots.get(order.provider_order_id);
      if (!snapshot) {
        return selectedQuality === "not_calculated";
      }
      return snapshot.quality === selectedQuality;
    });
  }, [orders, orderSnapshots, selectedQuality]);

  const selectedOrder =
    filteredOrders.find((order) => order.provider_order_id === selectedOrderID) ??
    filteredOrders[0] ??
    null;

  useEffect(() => {
    if (!selectedOrder) {
      return;
    }
    setForm((current) => ({
      ...current,
      itemIdentity: current.itemIdentity || (selectedOrder.items[0] ? itemIdentity(selectedOrder.items[0]) : ""),
    }));
  }, [selectedOrder]);

  const selectedOrderSnapshot = selectedOrder ? orderSnapshots.get(selectedOrder.provider_order_id) ?? null : null;
  const selectedOrderAdjustments = selectedOrder
    ? adjustments
        .filter((adjustment) => adjustment.provider_order_id === selectedOrder.provider_order_id)
        .sort((left, right) => right.created_at.localeCompare(left.created_at))
    : [];
  const selectedOrderInputs = selectedOrder
    ? inputs.filter((input) => input.provider_order_id === selectedOrder.provider_order_id)
    : [];
  const selectedOperationalFlags = selectedOrderSnapshot?.flags?.filter(isOperationalFlag) ?? [];
  const selectedMissingInputFlags = selectedOrderSnapshot?.flags?.filter(isMissingInputFlag) ?? [];
  const selectedRealizationLabel = realizationLabel(selectedOrderSnapshot);
  const hasValidOperator = Boolean(operator?.actor_type.trim() && operator.actor_id.trim());

  const summary = useMemo(() => {
    return filteredOrders.reduce(
      (acc, order) => {
        acc.total += 1;
        const snapshot = orderSnapshots.get(order.provider_order_id);
        if (!snapshot) {
          acc.notCalculated += 1;
          return acc;
        }
        if (snapshot.quality === "incomplete") {
          acc.incomplete += 1;
        }
        if (snapshot.quality === "negative_margin") {
          acc.negative += 1;
        }
        if (snapshot.quality === "not_realized") {
          acc.notRealized += 1;
        }
        if (snapshot.quality === "complete") {
          acc.complete += 1;
        }
        return acc;
      },
      { total: 0, incomplete: 0, negative: 0, notRealized: 0, complete: 0, notCalculated: 0 },
    );
  }, [filteredOrders, orderSnapshots]);

  async function runWorkspaceAction(actionKey: string, action: () => Promise<void>, successMessage: string) {
    setPendingAction(actionKey);
    setActionError(null);
    setActionMessage(null);
    try {
      await action();
      setActionMessage(successMessage);
    } catch (actionFailure) {
      setActionError(normalizeError(actionFailure, "Action failed."));
    } finally {
      setPendingAction(null);
    }
  }

  async function refreshCurrentInstallation() {
    if (!selectedInstallationID) {
      return;
    }
    const run = async () => {
      await Promise.all([
        loadWorkspace(selectedInstallationID),
        queryClient.refetchQueries({ queryKey: profitabilityQueryKeys.marginInputs(selectedInstallationID) }),
      ]);
    };
    await (client.withNoCache ? client.withNoCache(run) : run());
  }

  async function handleCreateAdjustment() {
    if (!selectedInstallationID || !selectedOrder || !operator || !hasValidOperator) {
      return;
    }
    const amount = Number(form.amount);
    if (!Number.isFinite(amount) || amount === 0) {
      setActionError("Amount must be a non-zero number.");
      return;
    }

    const selectedItem =
      form.scope === "item"
        ? selectedOrder.items.find((item) => itemIdentity(item) === form.itemIdentity)
        : undefined;

    await runWorkspaceAction(
      "create-adjustment",
      async () => {
        const actor = {
          actor_type: operator.actor_type,
          actor_id: operator.actor_id,
          actor_name: operator.actor_name,
        };
        const reason = form.reason || "Manual operator adjustment";
        const fingerprint = manualAdjustmentFingerprint({
          installationID: selectedInstallationID,
          providerOrderID: selectedOrder.provider_order_id,
          providerItemID: selectedItem?.provider_item_id,
          providerVariationID: selectedItem?.provider_variation_id,
          scope: form.scope,
          category: form.category,
          amount,
          currency: "BRL",
          reason,
          actor,
        });
        const previousSubmission = pendingAdjustmentSubmission.current;
        const idempotencyKey = previousSubmission?.fingerprint === fingerprint
          ? previousSubmission.idempotencyKey
          : createManualAdjustmentIdempotencyKey();
        pendingAdjustmentSubmission.current = { fingerprint, idempotencyKey };
        await client.createProfitabilityManualAdjustment({
          installation_id: selectedInstallationID,
          idempotency_key: idempotencyKey,
          provider_order_id: selectedOrder.provider_order_id,
          provider_item_id: selectedItem?.provider_item_id,
          provider_variation_id: selectedItem?.provider_variation_id,
          scope: form.scope,
          category: form.category,
          amount,
          currency: "BRL",
          reason,
          actor,
        });
        await client.calculateProfitabilityProfitSnapshots({ installation_id: selectedInstallationID, limit: 50 });
        await refreshCurrentInstallation();
        pendingAdjustmentSubmission.current = null;
        setForm(defaultAdjustmentForm(selectedOrder));
      },
      `Manual adjustment saved for order ${selectedOrder.provider_order_id}.`,
    );
  }

  function setParam(updates: Record<string, string | null>) {
    const next = new URLSearchParams(searchParams);
    for (const [key, value] of Object.entries(updates)) {
      if (!value) {
        next.delete(key);
      } else {
        next.set(key, value);
      }
    }
    setSearchParams(next, { replace: true });
  }

  return (
    <div className="space-y-5">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Orders + Margin ML</p>
          <h2 className="text-2xl font-semibold text-slate-900">Orders workspace</h2>
          <p className="text-sm text-slate-500">Inspect persisted order profitability, missing inputs, and operator adjustments.</p>
          <p className="text-sm text-slate-500"><FreshnessIndicator asOf={profitabilityQuery.data?.as_of} /></p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button
            variant="secondary"
            loading={pendingAction === "refresh"}
            disabled={!selectedInstallationID}
            onClick={() => void runWorkspaceAction("refresh", refreshCurrentInstallation, "Orders workspace refreshed.")}
          >
            Refresh
          </Button>
          <Button
            variant="secondary"
            loading={pendingAction === "import-orders"}
            onClick={() =>
              void runWorkspaceAction(
                "import-orders",
                async () => {
                  await client.importMarketplaceOrders({ installation_id: selectedInstallationID, limit: 20 });
                  await refreshCurrentInstallation();
                },
                "Orders imported from Mercado Livre.",
              )
            }
          >
            Import orders
          </Button>
          <Button
            variant="secondary"
            loading={pendingAction === "import-inputs"}
            onClick={() =>
              void runWorkspaceAction(
                "import-inputs",
                async () => {
                  await client.importProfitabilityMarginInputs({ installation_id: selectedInstallationID, limit: 50 });
                  await refreshCurrentInstallation();
                },
                "Margin inputs imported.",
              )
            }
          >
            Import inputs
          </Button>
          <Button
            variant="primary"
            loading={pendingAction === "recalculate"}
            onClick={() =>
              void runWorkspaceAction(
                "recalculate",
                async () => {
                  await client.calculateProfitabilityProfitSnapshots({ installation_id: selectedInstallationID, limit: 50 });
                  await refreshCurrentInstallation();
                },
                "Profit snapshots recalculated.",
              )
            }
          >
            Recalculate snapshots
          </Button>
        </div>
      </div>

      <SurfaceCard className="rounded-[24px]">
        <div className="grid gap-4 lg:grid-cols-[1fr_220px]">
          <label className="space-y-1 text-sm text-slate-700">
            <span className="block font-medium">Installation</span>
            <select
              aria-label="Installation"
              className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
              value={selectedInstallationID}
              onChange={(event) => setParam({ installation: event.target.value, order: null })}
            >
              {installations.map((installation) => (
                <option key={installation.installation_id} value={installation.installation_id}>
                  {installation.display_name}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1 text-sm text-slate-700">
            <span className="block font-medium">Margin quality</span>
            <select
              aria-label="Margin quality"
              className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
              value={selectedQuality}
              onChange={(event) => setParam({ quality: event.target.value === "all" ? null : event.target.value, order: null })}
            >
              <option value="all">All qualities</option>
              <option value="incomplete">Incomplete</option>
              <option value="negative_margin">Negative margin</option>
              <option value="not_realized">Not realized</option>
              <option value="complete">Complete</option>
              <option value="not_calculated">Not calculated</option>
            </select>
          </label>
        </div>
      </SurfaceCard>

      {actionError ? <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{actionError}</div> : null}
      {actionMessage ? <div className="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">{actionMessage}</div> : null}

      <div className="grid gap-4 md:grid-cols-5">
        <StatCard label="Orders" value={summary.total} />
        <StatCard label="Incomplete" value={summary.incomplete} />
        <StatCard label="Negative" value={summary.negative} />
        <StatCard label="Not realized" value={summary.notRealized} />
        <StatCard label="Complete" value={summary.complete} sub={`${summary.notCalculated} not calculated`} />
      </div>

      {state === "loading" ? <SurfaceCard><p className="text-sm text-slate-600">Loading orders workspace...</p></SurfaceCard> : null}
      {state === "error" && error ? <SurfaceCard><p className="text-sm text-red-700">{error}</p></SurfaceCard> : null}
      {state === "ready" && installations.length === 0 ? <SurfaceCard><p className="text-sm text-slate-600">No installations available.</p></SurfaceCard> : null}
      {state === "ready" && filteredOrders.length === 0 ? <SurfaceCard><p className="text-sm text-slate-600">No orders available for the selected installation.</p></SurfaceCard> : null}

      {state === "ready" && filteredOrders.length > 0 && selectedOrder ? (
        <div className="grid gap-6 xl:grid-cols-[1.05fr_0.95fr]">
          <div className="grid gap-4">
            {filteredOrders.map((order) => {
              const snapshot = orderSnapshots.get(order.provider_order_id);
              const orderQuality: QualityFilter = snapshot?.quality ?? "not_calculated";
              const linkQuality = deriveOrderLinkQuality(order.items);
              const selected = order.provider_order_id === selectedOrder.provider_order_id;

              return (
                <button
                  key={order.provider_order_id}
                  type="button"
                  className={`w-full rounded-[24px] border p-5 text-left transition-colors ${selected ? "border-blue-300 bg-blue-50" : "border-slate-200 bg-white hover:border-slate-300"}`}
                  onClick={() => setParam({ order: order.provider_order_id })}
                >
                  <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                    <div className="space-y-2">
                      <div className="flex flex-wrap gap-2">
                        <span className={`rounded-full px-3 py-1 text-xs font-semibold ${qualityTone(orderQuality)}`}>{qualityLabel(orderQuality)}</span>
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">{linkLabel(linkQuality)}</span>
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">{order.provider_status ?? "unknown"}</span>
                      </div>
                      <h3 className="text-lg font-semibold text-slate-900">Order #{order.provider_order_id}</h3>
                      <p className="text-sm text-slate-600">{order.items.length} item(s) · Updated {formatDate(order.provider_updated_at ?? order.fetched_at)}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Contribution</p>
                      <p className="mt-2 text-2xl font-semibold text-slate-900">{money(snapshot?.contribution_amount, snapshot?.currency ?? "BRL")}</p>
                    </div>
                  </div>
                  {snapshot?.flags?.length ? (
                    <div className="mt-4 flex flex-wrap gap-2">
                      {snapshot.flags.map((flag) => (
                        <span key={flag} className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-600">{flagLabel(flag)}</span>
                      ))}
                    </div>
                  ) : null}
                </button>
              );
            })}
          </div>

          <div className="grid gap-4">
            <SurfaceCard className="space-y-4 rounded-[24px]">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Selected order</p>
                  <h3 className="mt-2 text-xl font-semibold text-slate-900">Order #{selectedOrder.provider_order_id}</h3>
                  <p className="text-sm text-slate-500">Status: {selectedOrder.provider_status ?? "unknown"} · Shipping: {selectedOrder.shipping_id ?? "-"}</p>
                </div>
                <span className={`rounded-full px-3 py-1 text-xs font-semibold ${qualityTone(selectedOrderSnapshot?.quality ?? "not_calculated")}`}>
                  {qualityLabel(selectedOrderSnapshot?.quality ?? "not_calculated")}
                </span>
              </div>

              <div className="grid gap-3 sm:grid-cols-3">
                <div className="rounded-2xl bg-slate-50 p-4">
                  <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Revenue</p>
                  <p className="mt-2 text-xl font-semibold text-slate-900">{money(selectedOrderSnapshot?.revenue_amount, selectedOrderSnapshot?.currency ?? "BRL")}</p>
                </div>
                <div className="rounded-2xl bg-slate-50 p-4">
                  <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Contribution</p>
                  <p className="mt-2 text-xl font-semibold text-slate-900">{money(selectedOrderSnapshot?.contribution_amount, selectedOrderSnapshot?.currency ?? "BRL")}</p>
                </div>
                <div className="rounded-2xl bg-slate-50 p-4">
                  <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Margin</p>
                  <p className="mt-2 text-xl font-semibold text-slate-900">{percent(selectedOrderSnapshot?.margin_percent)}</p>
                </div>
              </div>

              {selectedRealizationLabel ? (
                <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-800">
                  <p className="font-medium">{selectedRealizationLabel}</p>
                  {selectedOperationalFlags.length ? (
                    <div className="mt-2 flex flex-wrap gap-2">
                      {selectedOperationalFlags.map((flag) => (
                        <span key={flag} className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-700">{flagLabel(flag)}</span>
                      ))}
                    </div>
                  ) : null}
                </div>
              ) : null}

              {selectedMissingInputFlags.length ? (
                <div className="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
                  <p className="font-medium">Data quality</p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    {selectedMissingInputFlags.map((flag) => (
                      <span key={flag} className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-700">{flagLabel(flag)}</span>
                    ))}
                  </div>
                </div>
              ) : null}
            </SurfaceCard>

            <SurfaceCard className="space-y-4 rounded-[24px]">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Item detail</p>
                  <h3 className="mt-1 text-lg font-semibold text-slate-900">Per-item profitability</h3>
                </div>
              </div>
              <div className="grid gap-3">
                {selectedOrder.items.map((item) => {
                  const snapshot = itemSnapshots.get(
                    snapshotIdentity({
                      provider_order_id: selectedOrder.provider_order_id,
                      provider_item_id: item.provider_item_id,
                      provider_variation_id: item.provider_variation_id,
                      scope: "item",
                    }),
                  );
                  return (
                    <article key={itemIdentity(item)} className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4">
                      <div className="flex flex-col gap-2 lg:flex-row lg:items-start lg:justify-between">
                        <div>
                          <p className="text-sm font-semibold text-slate-900">{item.title || item.provider_item_id}</p>
                          <p className="text-sm text-slate-600">{item.provider_item_id}{item.provider_variation_id ? ` / ${item.provider_variation_id}` : ""}</p>
                        </div>
                        <div className="flex flex-wrap gap-2">
                          <span className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-700">{linkLabel(item.link_quality)}</span>
                          <span className={`rounded-full px-3 py-1 text-xs font-medium ${qualityTone(snapshot?.quality ?? "not_calculated")}`}>{qualityLabel(snapshot?.quality ?? "not_calculated")}</span>
                        </div>
                      </div>
                      <div className="mt-3 grid gap-3 sm:grid-cols-3">
                        <div>
                          <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Revenue</p>
                          <p className="mt-1 text-sm font-semibold text-slate-900">{money(snapshot?.revenue_amount, snapshot?.currency ?? "BRL")}</p>
                        </div>
                        <div>
                          <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Contribution</p>
                          <p className="mt-1 text-sm font-semibold text-slate-900">{money(snapshot?.contribution_amount, snapshot?.currency ?? "BRL")}</p>
                        </div>
                        <div>
                          <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Margin</p>
                          <p className="mt-1 text-sm font-semibold text-slate-900">{percent(snapshot?.margin_percent)}</p>
                        </div>
                      </div>
                    </article>
                  );
                })}
              </div>
            </SurfaceCard>

            <SurfaceCard className="space-y-4 rounded-[24px]">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Input breakdown</p>
                <h3 className="mt-1 text-lg font-semibold text-slate-900">Persisted margin inputs</h3>
              </div>
              {selectedOrderInputs.length === 0 ? (
                <p className="text-sm text-slate-500">No persisted inputs for this order yet.</p>
              ) : (
                <div className="grid gap-3">
                  {selectedOrderInputs.map((input) => (
                    <article key={`${input.scope}:${input.kind}:${input.provider_item_id ?? "order"}`} className="rounded-2xl border border-slate-200 bg-white px-4 py-4">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700">{input.scope}</span>
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700">{inputLabel(input.kind)}</span>
                        <span className={`rounded-full px-3 py-1 text-xs font-medium ${qualityTone(input.quality === "complete" ? "complete" : "incomplete")}`}>{input.quality}</span>
                      </div>
                      <p className="mt-2 text-sm font-semibold text-slate-900">{money(input.amount, input.currency)}</p>
                      <p className="mt-1 text-sm text-slate-600">
                        Source: {input.source_system}
                        {input.source_reference ? ` · ${input.source_reference}` : ""}
                        {input.provider_item_id ? ` · ${input.provider_item_id}` : ""}
                      </p>
                    </article>
                  ))}
                </div>
              )}
            </SurfaceCard>

            <SurfaceCard className="space-y-4 rounded-[24px]">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Manual adjustments</p>
                <h3 className="mt-1 text-lg font-semibold text-slate-900">Manual adjustments</h3>
              </div>

              {selectedOrderAdjustments.length === 0 ? (
                <p className="text-sm text-slate-500">No manual adjustments yet.</p>
              ) : (
                <div className="grid gap-3">
                  {selectedOrderAdjustments.map((adjustment) => (
                    <article key={`${adjustment.adjustment_id}:${adjustmentIdentity(adjustment)}`} className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-700">{adjustment.scope}</span>
                        <span className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-700">{adjustment.category}</span>
                      </div>
                      <p className="mt-2 text-sm font-semibold text-slate-900">{money(adjustment.amount, adjustment.currency)}</p>
                      <p className="mt-1 text-sm text-slate-600">{adjustment.reason}</p>
                      <p className="mt-1 text-xs text-slate-500">{adjustment.actor.actor_name || adjustment.actor.actor_id} · {formatDate(adjustment.created_at)}</p>
                    </article>
                  ))}
                </div>
              )}

              <div className="rounded-[24px] border border-slate-200 bg-slate-50 p-4">
                <div className="grid gap-3 md:grid-cols-2">
                  <label className="space-y-1 text-sm text-slate-700">
                    <span className="block font-medium">Scope</span>
                    <select
                      aria-label="Adjustment scope"
                      className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                      value={form.scope}
                      onChange={(event) => setForm((current) => ({ ...current, scope: event.target.value as ProfitabilityInputScope }))}
                    >
                      <option value="order">order</option>
                      <option value="item">item</option>
                    </select>
                  </label>
                  <label className="space-y-1 text-sm text-slate-700">
                    <span className="block font-medium">Category</span>
                    <select
                      aria-label="Adjustment category"
                      className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                      value={form.category}
                      onChange={(event) => setForm((current) => ({ ...current, category: event.target.value as ProfitabilityManualAdjustmentCategory }))}
                    >
                      <option value="freight">freight</option>
                      <option value="commission">commission</option>
                      <option value="cost">cost</option>
                      <option value="generic_adjustment">generic_adjustment</option>
                    </select>
                  </label>
                  {form.scope === "item" ? (
                    <label className="space-y-1 text-sm text-slate-700 md:col-span-2">
                      <span className="block font-medium">Order item</span>
                      <select
                        aria-label="Adjustment item"
                        className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                        value={form.itemIdentity}
                        onChange={(event) => setForm((current) => ({ ...current, itemIdentity: event.target.value }))}
                      >
                        {selectedOrder.items.map((item) => (
                          <option key={itemIdentity(item)} value={itemIdentity(item)}>
                            {item.title || item.provider_item_id}
                          </option>
                        ))}
                      </select>
                    </label>
                  ) : null}
                  <label className="space-y-1 text-sm text-slate-700">
                    <span className="block font-medium">Amount</span>
                    <input
                      aria-label="Adjustment amount"
                      className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                      value={form.amount}
                      onChange={(event) => setForm((current) => ({ ...current, amount: event.target.value }))}
                    />
                  </label>
                  <label className="space-y-1 text-sm text-slate-700 md:col-span-2">
                    <span className="block font-medium">Reason</span>
                    <input
                      aria-label="Adjustment reason"
                      className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                      value={form.reason}
                      onChange={(event) => setForm((current) => ({ ...current, reason: event.target.value }))}
                    />
                  </label>
                </div>
                <p className="mt-3 text-sm text-slate-700" role="status">
                  {hasValidOperator
                    ? `Operator: ${operator?.actor_name || operator?.actor_id}`
                    : "Manual adjustment creation is disabled because trustworthy operator identity is unavailable."}
                </p>
                <div className="mt-4 flex justify-end">
                  <Button disabled={!hasValidOperator} variant="primary" loading={pendingAction === "create-adjustment"} onClick={() => void handleCreateAdjustment()}>
                    Save adjustment
                  </Button>
                </div>
              </div>
            </SurfaceCard>
          </div>
        </div>
      ) : null}
    </div>
  );
}
