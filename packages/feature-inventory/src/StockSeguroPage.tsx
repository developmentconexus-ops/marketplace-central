import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, SurfaceCard } from "@marketplace-central/ui";
import type {
  ApplyInventoryStockActionResponse,
  IntegrationInstallation,
  InventoryStockRiskItem,
} from "@marketplace-central/sdk-runtime";
import {
  FreshnessIndicator,
  inventoryQueryKeys,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface StockSeguroClient extends RefreshableClient {
  listIntegrationInstallations: () => Promise<{ items: IntegrationInstallation[] }>;
  listInventoryStockRisks: (input: {
    installation_id: string;
    state?: InventoryStockRiskItem["state"];
    link_state?: InventoryStockRiskItem["link_state"];
    actionability?: InventoryStockRiskItem["actionability"];
    limit?: number;
  }) => Promise<{ items: InventoryStockRiskItem[]; as_of?: string }>;
  applyInventoryManualStockAction: (req: {
    stock_action_id: string;
    installation_id: string;
    provider_item_id: string;
    provider_variation_id?: string;
    requested_quantity: number;
    reason?: string;
    approval: {
      approved: boolean;
      approved_at: string;
      operator: {
        actor_type: string;
        actor_id: string;
        actor_name?: string;
      };
    };
  }) => Promise<ApplyInventoryStockActionResponse>;
}

export interface StockSeguroPageProps {
  client: StockSeguroClient;
}

type LoadState = "loading" | "ready" | "error";

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

function stateLabel(state: InventoryStockRiskItem["state"]): string {
  switch (state) {
    case "oversell":
      return "Oversell";
    case "undersell":
      return "Undersell";
    case "stale":
      return "Stale";
    case "conflict":
      return "Conflito";
    case "unresolved":
      return "Sem vinculo";
    case "ineligible":
      return "Ineligivel";
    case "unsupported":
      return "Sem suporte";
    default:
      return "Healthy";
  }
}

function stateTone(state: InventoryStockRiskItem["state"]): string {
  switch (state) {
    case "oversell":
      return "bg-red-100 text-red-700";
    case "undersell":
      return "bg-amber-100 text-amber-800";
    case "stale":
      return "bg-slate-200 text-slate-700";
    case "conflict":
      return "bg-rose-100 text-rose-700";
    case "unresolved":
      return "bg-orange-100 text-orange-700";
    case "ineligible":
      return "bg-violet-100 text-violet-700";
    case "unsupported":
      return "bg-slate-100 text-slate-700";
    default:
      return "bg-emerald-100 text-emerald-700";
  }
}

function actionabilityLabel(actionable: boolean): string {
  return actionable ? "Actionable" : "Blocked";
}

function itemKey(item: InventoryStockRiskItem): string {
  return `${item.identity.installation_id}:${item.identity.provider_item_id}:${item.identity.provider_variation_id ?? ""}`;
}

function actorIdFromName(name: string): string {
  return name.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "") || "operator";
}

export function StockSeguroPage({ client }: StockSeguroPageProps) {
  const queryClient = useQueryClient();
  const [searchParams, setSearchParams] = useSearchParams();
  const [installations, setInstallations] = useState<IntegrationInstallation[]>([]);
  const [items, setItems] = useState<InventoryStockRiskItem[]>([]);
  const [state, setState] = useState<LoadState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [openConfirmKey, setOpenConfirmKey] = useState<string | null>(null);
  const [operatorName, setOperatorName] = useState("Leandro");
  const [reason, setReason] = useState("Stock Seguro manual apply");

  const selectedInstallationID = searchParams.get("installation") ?? "";
  const selectedState = (searchParams.get("state") ?? "") as InventoryStockRiskItem["state"] | "";
  const selectedLinkState = (searchParams.get("link_state") ?? "") as InventoryStockRiskItem["link_state"] | "";
  const selectedActionability = (searchParams.get("actionability") ?? "") as InventoryStockRiskItem["actionability"] | "";

  const riskFilters = {
    state: selectedState || undefined,
    link_state: selectedLinkState || undefined,
    actionability: selectedActionability || undefined,
    limit: 50,
  };
  const riskQuery = useQuery({
    queryKey: inventoryQueryKeys.risks(selectedInstallationID, riskFilters),
    queryFn: () => client.listInventoryStockRisks({
      installation_id: selectedInstallationID,
      ...riskFilters,
    }),
    staleTime: QUERY_STALE_TIME.stock,
    enabled: Boolean(selectedInstallationID),
  });
  useEffect(() => {
    if (riskQuery.data) {
      setItems(riskQuery.data.items);
      setError(null);
      setState("ready");
    }
    if (riskQuery.error) {
      setError(normalizeError(riskQuery.error, "Failed to load stock risks."));
      setState("error");
    }
  }, [riskQuery.data, riskQuery.error]);

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

  const displayedItems = riskQuery.data?.items ?? items;
  const displayedError = error ?? (riskQuery.error ? normalizeError(riskQuery.error, "Failed to load stock risks.") : null);
  const displayedState: LoadState = riskQuery.error
    ? "error"
    : riskQuery.data
      ? "ready"
      : state;

  const summary = useMemo(() => {
    return displayedItems.reduce(
      (acc, item) => {
        acc.total += 1;
        if (item.actionable) {
          acc.actionable += 1;
        }
        if (item.state === "oversell") {
          acc.oversell += 1;
        }
        if (!item.actionable) {
          acc.blockers += 1;
        }
        return acc;
      },
      { total: 0, actionable: 0, oversell: 0, blockers: 0 },
    );
  }, [displayedItems]);

  const selected = displayedItems[0];

  async function applyAction(item: InventoryStockRiskItem) {
    if (item.recommended_quantity == null) {
      return;
    }
    const key = itemKey(item);
    setPendingKey(key);
    setActionError(null);
    setActionMessage(null);
    try {
      const result = await client.applyInventoryManualStockAction({
        stock_action_id: `ssa-${Date.now()}`,
        installation_id: item.identity.installation_id,
        provider_item_id: item.identity.provider_item_id,
        provider_variation_id: item.identity.provider_variation_id,
        requested_quantity: item.recommended_quantity,
        reason,
        approval: {
          approved: true,
          approved_at: new Date().toISOString(),
          operator: {
            actor_type: "operator",
            actor_id: actorIdFromName(operatorName),
            actor_name: operatorName,
          },
        },
      });
      setActionMessage(`Action ${result.action.state} for ${item.identity.provider_item_id}.`);
      setOpenConfirmKey(null);
      const refreshed = await client.listInventoryStockRisks({
        installation_id: selectedInstallationID,
        state: selectedState || undefined,
        link_state: selectedLinkState || undefined,
        actionability: selectedActionability || undefined,
        limit: 50,
      });
      setItems(refreshed.items);
      queryClient.setQueryData(inventoryQueryKeys.risks(selectedInstallationID, riskFilters), (current: { items: InventoryStockRiskItem[]; as_of?: string } | undefined) => ({
        ...current,
        items: refreshed.items,
      }));
    } catch (runError) {
      setActionError(normalizeError(runError, "Failed to apply stock action."));
    } finally {
      setPendingKey(null);
    }
  }

  async function refreshRisks() {
    if (!selectedInstallationID) {
      return;
    }
    const run = () => queryClient.refetchQueries({
      queryKey: inventoryQueryKeys.risks(selectedInstallationID, riskFilters),
    });
    await (client.withNoCache ? client.withNoCache(run) : run());
  }

  function updateFilters(patch: Record<string, string>) {
    const next = new URLSearchParams(searchParams);
    Object.entries(patch).forEach(([key, value]) => {
      if (value) {
        next.set(key, value);
      } else {
        next.delete(key);
      }
    });
    setSearchParams(next);
  }

  return (
    <div className="space-y-6">
      <section className="rounded-[28px] border border-slate-200 bg-[linear-gradient(135deg,#10233f_0%,#17325a_50%,#234c7f_100%)] px-6 py-7 text-white shadow-sm">
        <div className="flex flex-col gap-6 xl:flex-row xl:items-end xl:justify-between">
          <div className="max-w-2xl space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">Inventory</p>
            <h1 className="text-3xl font-semibold tracking-tight">Stock Seguro</h1>
            <p className="text-sm leading-6 text-slate-200">
              Scan persisted Mercado Livre stock drift, inspect source evidence, and only apply explicit manual corrections.
            </p>
            <p className="text-sm text-cyan-100"><FreshnessIndicator asOf={riskQuery.data?.as_of} /></p>
          </div>
          <div className="grid gap-3 rounded-2xl border border-white/15 bg-white/8 p-4 sm:grid-cols-2 xl:min-w-[460px]">
            <Button variant="secondary" disabled={!selectedInstallationID || riskQuery.isFetching} loading={riskQuery.isFetching} onClick={() => void refreshRisks()}>Refresh</Button>
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Installation</span>
              <select
                aria-label="Installation"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white"
                value={selectedInstallationID}
                onChange={(event) => updateFilters({ installation: event.target.value })}
              >
                {installations.map((installation) => (
                  <option key={installation.installation_id} value={installation.installation_id} className="text-slate-900">
                    {installation.display_name}
                  </option>
                ))}
              </select>
            </label>
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Risk</span>
              <select
                aria-label="Risk filter"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white"
                value={selectedState}
                onChange={(event) => updateFilters({ state: event.target.value })}
              >
                <option value="" className="text-slate-900">All</option>
                <option value="oversell" className="text-slate-900">Oversell</option>
                <option value="undersell" className="text-slate-900">Undersell</option>
                <option value="healthy" className="text-slate-900">Healthy</option>
                <option value="stale" className="text-slate-900">Stale</option>
                <option value="conflict" className="text-slate-900">Conflito</option>
                <option value="unresolved" className="text-slate-900">Sem vinculo</option>
                <option value="ineligible" className="text-slate-900">Ineligivel</option>
                <option value="unsupported" className="text-slate-900">Sem suporte</option>
              </select>
            </label>
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Link</span>
              <select
                aria-label="Link filter"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white"
                value={selectedLinkState}
                onChange={(event) => updateFilters({ link_state: event.target.value })}
              >
                <option value="" className="text-slate-900">All</option>
                <option value="resolved" className="text-slate-900">Resolved</option>
                <option value="conflict" className="text-slate-900">Conflito</option>
                <option value="unresolved" className="text-slate-900">Sem vinculo</option>
                <option value="rejected" className="text-slate-900">Rejected</option>
              </select>
            </label>
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Actionability</span>
              <select
                aria-label="Actionability filter"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white"
                value={selectedActionability}
                onChange={(event) => updateFilters({ actionability: event.target.value })}
              >
                <option value="" className="text-slate-900">All</option>
                <option value="actionable" className="text-slate-900">Actionable</option>
                <option value="blocked" className="text-slate-900">Blocked</option>
              </select>
            </label>
          </div>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-4">
        <SurfaceCard><p className="text-xs uppercase tracking-[0.18em] text-slate-500">Rows</p><p data-testid="stock-summary-total" className="mt-2 text-3xl font-semibold text-slate-900">{summary.total}</p></SurfaceCard>
        <SurfaceCard><p className="text-xs uppercase tracking-[0.18em] text-slate-500">Actionable</p><p data-testid="stock-summary-actionable" className="mt-2 text-3xl font-semibold text-slate-900">{summary.actionable}</p></SurfaceCard>
        <SurfaceCard><p className="text-xs uppercase tracking-[0.18em] text-slate-500">Oversell</p><p data-testid="stock-summary-oversell" className="mt-2 text-3xl font-semibold text-slate-900">{summary.oversell}</p></SurfaceCard>
        <SurfaceCard><p className="text-xs uppercase tracking-[0.18em] text-slate-500">Blockers</p><p data-testid="stock-summary-blockers" className="mt-2 text-3xl font-semibold text-slate-900">{summary.blockers}</p></SurfaceCard>
      </section>

      {actionError ? <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{actionError}</div> : null}
      {actionMessage ? <div className="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">{actionMessage}</div> : null}

      {displayedState === "loading" ? <SurfaceCard><p className="text-sm text-slate-600">Loading Stock Seguro...</p></SurfaceCard> : null}
      {displayedState === "error" && displayedError ? <SurfaceCard><p className="text-sm text-red-700">{displayedError}</p></SurfaceCard> : null}
      {displayedState === "ready" && displayedItems.length === 0 ? <SurfaceCard><p className="text-sm text-slate-600">No Stock Seguro rows for the selected filters.</p></SurfaceCard> : null}

      {displayedState === "ready" && displayedItems.length > 0 ? (
        <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
          <div className="grid gap-4">
            {displayedItems.map((item) => {
              const key = itemKey(item);
              const confirmOpen = openConfirmKey === key;
              return (
                <SurfaceCard key={key} className="space-y-4 rounded-[24px]">
                  <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                    <div className="space-y-2">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className={`rounded-full px-3 py-1 text-xs font-semibold ${stateTone(item.state)}`}>{stateLabel(item.state)}</span>
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">{actionabilityLabel(item.actionable)}</span>
                        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">{item.link_state}</span>
                      </div>
                      <h2 className="text-lg font-semibold text-slate-900">{item.title || item.identity.provider_item_id}</h2>
                      <p className="text-sm text-slate-600">
                        {item.identity.provider_item_id}
                        {item.identity.provider_variation_id ? ` / ${item.identity.provider_variation_id}` : ""}
                        {item.internal_reference_code ? ` · ${item.internal_reference_code}` : ""}
                      </p>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      <Button
                        variant="primary"
                        disabled={!item.actionable || item.recommended_quantity == null}
                        onClick={() => setOpenConfirmKey(confirmOpen ? null : key)}
                      >
                        Apply recommended quantity
                      </Button>
                    </div>
                  </div>

                  <div className="grid gap-3 sm:grid-cols-3">
                    <div className="rounded-2xl bg-slate-50 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-slate-500">ML quantity</p>
                      <p className="mt-2 text-2xl font-semibold text-slate-900">{item.provider_quantity ?? "-"}</p>
                    </div>
                    <div className="rounded-2xl bg-slate-50 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Internal sellable</p>
                      <p className="mt-2 text-2xl font-semibold text-slate-900">{item.internal_quantity ?? "-"}</p>
                    </div>
                    <div className="rounded-2xl bg-slate-50 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-slate-500">Recommended</p>
                      <p className="mt-2 text-2xl font-semibold text-slate-900">{item.recommended_quantity ?? "-"}</p>
                    </div>
                  </div>

                  {confirmOpen ? (
                    <div className="rounded-[24px] border border-blue-100 bg-blue-50 p-4">
                      <p className="text-sm font-medium text-slate-900">Confirm manual apply</p>
                      <p className="mt-1 text-sm text-slate-600">
                        This will request ML quantity <strong>{item.recommended_quantity ?? "-"}</strong> for this listing using the current policy evidence.
                      </p>
                      <div className="mt-4 grid gap-3 md:grid-cols-2">
                        <label className="space-y-1 text-sm text-slate-700">
                          <span className="block font-medium">Operator</span>
                          <input className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2" value={operatorName} onChange={(event) => setOperatorName(event.target.value)} />
                        </label>
                        <label className="space-y-1 text-sm text-slate-700">
                          <span className="block font-medium">Reason</span>
                          <input className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2" value={reason} onChange={(event) => setReason(event.target.value)} />
                        </label>
                      </div>
                      <div className="mt-4 flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => setOpenConfirmKey(null)}>Cancel</Button>
                        <Button variant="primary" loading={pendingKey === key} onClick={() => void applyAction(item)}>
                          Confirm apply
                        </Button>
                      </div>
                    </div>
                  ) : null}

                  {item.blocking_reason?.message ? (
                    <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-600">
                      <strong className="text-slate-900">Blocked:</strong> {item.blocking_reason.message}
                    </div>
                  ) : null}
                </SurfaceCard>
              );
            })}
          </div>

          {selected ? (
            <SurfaceCard className="space-y-4 rounded-[24px]">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-slate-500">Row detail</p>
                <h2 className="mt-2 text-xl font-semibold text-slate-900">{selected.title || selected.identity.provider_item_id}</h2>
              </div>
              <div className="grid gap-3 text-sm text-slate-700">
                <p><strong>Policy:</strong> {selected.policy_id}</p>
                <p><strong>Provider observed:</strong> {selected.provider_observed_at ?? "-"}</p>
                <p><strong>Internal observed:</strong> {selected.internal_observed_at ?? "-"}</p>
                <p><strong>Internal product:</strong> {selected.internal_product_name || selected.internal_reference_code || "-"}</p>
                <p><strong>Seller SKU:</strong> {selected.seller_sku || "-"}</p>
                <p><strong>EAN:</strong> {selected.ean || "-"}</p>
              </div>
            </SurfaceCard>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}
