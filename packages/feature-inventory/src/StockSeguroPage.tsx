import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button, SurfaceCard } from "@marketplace-central/ui";
import type {
  ApplyInventoryStockActionResponse,
  IntegrationInstallation,
  InventoryStockRiskItem,
} from "@marketplace-central/sdk-runtime";
import {
  formatDateTime,
  FreshnessIndicator,
  inventoryQueryKeys,
  queryKeyNamespaces,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface StockSeguroClient extends RefreshableClient {
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
  installations: IntegrationInstallation[];
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

// "Oversell"/"Undersell" stay as-is: they are the operator's working vocabulary
// for the drift direction. Everything else reads in pt-BR.
function stateLabel(state: InventoryStockRiskItem["state"]): string {
  switch (state) {
    case "oversell":
      return "Oversell";
    case "undersell":
      return "Undersell";
    case "stale":
      return "Desatualizado";
    case "conflict":
      return "Conflito";
    case "unresolved":
      return "Sem vínculo";
    case "ineligible":
      return "Inelegível";
    case "unsupported":
      return "Sem suporte";
    default:
      return "Saudável";
  }
}

function stateTone(state: InventoryStockRiskItem["state"]): string {
  switch (state) {
    case "oversell":
    case "conflict":
      return "bg-warn-soft text-warn";
    case "undersell":
    case "unresolved":
      return "bg-warn-soft text-warn";
    case "stale":
    case "ineligible":
    case "unsupported":
      return "bg-surface-2 text-muted";
    default:
      return "bg-accent-soft text-accent-ink";
  }
}

function linkStateLabel(linkState: InventoryStockRiskItem["link_state"]): string {
  switch (linkState) {
    case "resolved":
      return "Vinculado";
    case "conflict":
      return "Conflito";
    case "rejected":
      return "Rejeitado";
    default:
      return "Sem vínculo";
  }
}

function actionabilityLabel(actionable: boolean): string {
  return actionable ? "Acionável" : "Bloqueado";
}

// The action state comes from the backend enum; an unmapped value is shown raw
// rather than guessed at.
function actionStateLabel(state: string): string {
  switch (state) {
    case "applied":
      return "aplicada";
    case "failed":
      return "falhou";
    case "pending":
      return "pendente";
    default:
      return state;
  }
}

// The API's blocking reason carries a stable code plus an English engineering
// message. The code is the contract, so the screen translates it; an unmapped
// code falls back to the API message instead of inventing copy.
const blockingReasonCopy: Record<string, string> = {
  unresolved_link: "Anúncio sem vínculo com um produto do ERP.",
  rejected_link: "O vínculo com o produto do ERP foi rejeitado.",
  conflict_link: "O vínculo com o produto do ERP está em conflito.",
  missing_internal_product_link: "O vínculo está resolvido mas não aponta para um produto do ERP.",
  missing_internal_stock: "O ERP não informou o estoque disponível deste produto.",
  missing_provider_quantity: "O Mercado Livre não informou a quantidade deste anúncio.",
  stale_internal_source: "O estoque do ERP está mais antigo que o permitido pela política.",
  stale_provider_source: "O estoque do Mercado Livre está mais antigo que o permitido pela política.",
  ineligible_product: "Produto fora da política de Stock Seguro.",
};

function blockingReasonText(reason: { code?: string; message?: string } | undefined): string | null {
  if (!reason) {
    return null;
  }
  const mapped = reason.code ? blockingReasonCopy[reason.code] : undefined;
  return mapped ?? reason.message ?? null;
}

function itemKey(item: InventoryStockRiskItem): string {
  return `${item.identity.installation_id}:${item.identity.provider_item_id}:${item.identity.provider_variation_id ?? ""}`;
}

function actorIdFromName(name: string): string {
  return name.trim().toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "") || "operator";
}

export function StockSeguroPage({ client, installations }: StockSeguroPageProps) {
  const queryClient = useQueryClient();
  const [searchParams, setSearchParams] = useSearchParams();
  const [items, setItems] = useState<InventoryStockRiskItem[]>([]);
  const [state, setState] = useState<LoadState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [openConfirmKey, setOpenConfirmKey] = useState<string | null>(null);
  // The apply is audited under whoever approves it, so the operator types their
  // own name — a prefilled one would attribute the action to someone else.
  const [operatorName, setOperatorName] = useState("");
  const [reason, setReason] = useState("Aplicação manual Stock Seguro");

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
  const stockActionMutation = useMutation({
    mutationFn: (req: Parameters<StockSeguroClient["applyInventoryManualStockAction"]>[0]) =>
      client.applyInventoryManualStockAction(req),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: queryKeyNamespaces.inventory }),
  });
  useEffect(() => {
    if (riskQuery.data) {
      setItems(riskQuery.data.items);
      setError(null);
      setState("ready");
    }
    if (riskQuery.error) {
      setError(normalizeError(riskQuery.error, "Não foi possível carregar os riscos de estoque."));
      setState("error");
    }
  }, [riskQuery.data, riskQuery.error]);

  const displayedItems = riskQuery.data?.items ?? items;
  const displayedError =
    error ??
    (riskQuery.error ? normalizeError(riskQuery.error, "Não foi possível carregar os riscos de estoque.") : null);
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
    if (item.recommended_quantity == null || !operatorName.trim()) {
      return;
    }
    const key = itemKey(item);
    setPendingKey(key);
    setActionError(null);
    setActionMessage(null);
    try {
      const result = await stockActionMutation.mutateAsync({
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
      setActionMessage(`Ação ${actionStateLabel(result.action.state)} para ${item.identity.provider_item_id}.`);
      setOpenConfirmKey(null);
    } catch (runError) {
      setActionError(normalizeError(runError, "Não foi possível aplicar a ação de estoque."));
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

  const selectClass = "w-full rounded-control border border-border bg-surface px-3 py-2 text-sm text-ink";

  return (
    <div className="space-y-5">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div className="max-w-2xl">
          <p className="text-xs uppercase tracking-[0.2em] text-faint">Estoque</p>
          <h1 className="text-2xl font-semibold text-ink">Stock Seguro</h1>
          <p className="text-sm text-muted">
            Compara o estoque publicado no Mercado Livre com o estoque interno e aplica correções manuais explícitas.
          </p>
          <p className="text-sm text-muted">
            <FreshnessIndicator asOf={riskQuery.data?.as_of} />
          </p>
        </div>
        <Button
          variant="secondary"
          disabled={!selectedInstallationID || riskQuery.isFetching}
          loading={riskQuery.isFetching}
          onClick={() => void refreshRisks()}
        >
          Atualizar
        </Button>
      </div>

      <SurfaceCard className="rounded-card">
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <label className="space-y-1 text-sm text-ink">
            <span className="block font-medium">Conta</span>
            <select
              aria-label="Conta"
              className={selectClass}
              value={selectedInstallationID}
              onChange={(event) => updateFilters({ installation: event.target.value })}
            >
              {installations.map((installation) => (
                <option key={installation.installation_id} value={installation.installation_id}>
                  {installation.external_account_name?.trim() || installation.display_name}
                </option>
              ))}
            </select>
          </label>
          <label className="space-y-1 text-sm text-ink">
            <span className="block font-medium">Risco</span>
            <select
              aria-label="Filtro de risco"
              className={selectClass}
              value={selectedState}
              onChange={(event) => updateFilters({ state: event.target.value })}
            >
              <option value="">Todos</option>
              <option value="oversell">Oversell</option>
              <option value="undersell">Undersell</option>
              <option value="healthy">Saudável</option>
              <option value="stale">Desatualizado</option>
              <option value="conflict">Conflito</option>
              <option value="unresolved">Sem vínculo</option>
              <option value="ineligible">Inelegível</option>
              <option value="unsupported">Sem suporte</option>
            </select>
          </label>
          <label className="space-y-1 text-sm text-ink">
            <span className="block font-medium">Vínculo</span>
            <select
              aria-label="Filtro de vínculo"
              className={selectClass}
              value={selectedLinkState}
              onChange={(event) => updateFilters({ link_state: event.target.value })}
            >
              <option value="">Todos</option>
              <option value="resolved">Vinculado</option>
              <option value="conflict">Conflito</option>
              <option value="unresolved">Sem vínculo</option>
              <option value="rejected">Rejeitado</option>
            </select>
          </label>
          <label className="space-y-1 text-sm text-ink">
            <span className="block font-medium">Ação</span>
            <select
              aria-label="Filtro de ação"
              className={selectClass}
              value={selectedActionability}
              onChange={(event) => updateFilters({ actionability: event.target.value })}
            >
              <option value="">Todos</option>
              <option value="actionable">Acionável</option>
              <option value="blocked">Bloqueado</option>
            </select>
          </label>
        </div>
      </SurfaceCard>

      <section className="grid gap-4 md:grid-cols-4">
        <SurfaceCard>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Linhas</p>
          <p data-testid="stock-summary-total" className="mt-2 text-3xl font-semibold text-ink">{summary.total}</p>
        </SurfaceCard>
        <SurfaceCard>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Acionáveis</p>
          <p data-testid="stock-summary-actionable" className="mt-2 text-3xl font-semibold text-ink">{summary.actionable}</p>
        </SurfaceCard>
        <SurfaceCard>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Oversell</p>
          <p data-testid="stock-summary-oversell" className="mt-2 text-3xl font-semibold text-ink">{summary.oversell}</p>
        </SurfaceCard>
        <SurfaceCard>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Bloqueadas</p>
          <p data-testid="stock-summary-blockers" className="mt-2 text-3xl font-semibold text-ink">{summary.blockers}</p>
        </SurfaceCard>
      </section>

      {actionError ? (
        <div className="rounded-card border border-border bg-warn-soft px-4 py-3 text-sm text-warn">{actionError}</div>
      ) : null}
      {actionMessage ? (
        <div className="rounded-card border border-border bg-accent-soft px-4 py-3 text-sm text-accent-ink">{actionMessage}</div>
      ) : null}

      {displayedState === "loading" ? (
        <SurfaceCard><p className="text-sm text-muted">Carregando Stock Seguro...</p></SurfaceCard>
      ) : null}
      {displayedState === "error" && displayedError ? (
        <SurfaceCard><p className="text-sm text-warn">{displayedError}</p></SurfaceCard>
      ) : null}
      {displayedState === "ready" && displayedItems.length === 0 ? (
        <SurfaceCard><p className="text-sm text-muted">Nenhuma linha de Stock Seguro para os filtros selecionados.</p></SurfaceCard>
      ) : null}

      {displayedState === "ready" && displayedItems.length > 0 ? (
        <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
          <div className="grid gap-4">
            {displayedItems.map((item) => {
              const key = itemKey(item);
              const confirmOpen = openConfirmKey === key;
              return (
                <SurfaceCard key={key} className="space-y-4 rounded-card">
                  <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
                    <div className="space-y-2">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className={`rounded-pill px-3 py-1 text-xs font-semibold ${stateTone(item.state)}`}>{stateLabel(item.state)}</span>
                        <span className="rounded-pill bg-surface-2 px-3 py-1 text-xs font-semibold text-muted">{actionabilityLabel(item.actionable)}</span>
                        <span className="rounded-pill bg-surface-2 px-3 py-1 text-xs font-semibold text-muted">{linkStateLabel(item.link_state)}</span>
                      </div>
                      <h2 className="text-lg font-semibold text-ink">{item.title || item.identity.provider_item_id}</h2>
                      <p className="text-sm text-muted">
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
                        Aplicar quantidade recomendada
                      </Button>
                    </div>
                  </div>

                  <div className="grid gap-3 sm:grid-cols-3">
                    <div className="rounded-card bg-surface-2 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-faint">Qtd. no ML</p>
                      <p className="mt-2 text-2xl font-semibold text-ink">{item.provider_quantity ?? "—"}</p>
                    </div>
                    <div className="rounded-card bg-surface-2 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-faint">Estoque interno</p>
                      <p className="mt-2 text-2xl font-semibold text-ink">{item.internal_quantity ?? "—"}</p>
                    </div>
                    <div className="rounded-card bg-surface-2 p-4">
                      <p className="text-xs uppercase tracking-[0.16em] text-faint">Recomendado</p>
                      <p className="mt-2 text-2xl font-semibold text-ink">{item.recommended_quantity ?? "—"}</p>
                    </div>
                  </div>

                  {confirmOpen ? (
                    <div className="rounded-card border border-border bg-surface-2 p-4">
                      <p className="text-sm font-medium text-ink">Confirmar aplicação manual</p>
                      <p className="mt-1 text-sm text-muted">
                        Isto vai pedir ao Mercado Livre a quantidade <strong>{item.recommended_quantity ?? "—"}</strong> para
                        este anúncio, usando a evidência da política atual.
                      </p>
                      <div className="mt-4 grid gap-3 md:grid-cols-2">
                        <label className="space-y-1 text-sm text-ink">
                          <span className="block font-medium">Operador</span>
                          <input
                            className="w-full rounded-control border border-border bg-surface px-3 py-2"
                            placeholder="Seu nome"
                            value={operatorName}
                            onChange={(event) => setOperatorName(event.target.value)}
                          />
                        </label>
                        <label className="space-y-1 text-sm text-ink">
                          <span className="block font-medium">Motivo</span>
                          <input
                            className="w-full rounded-control border border-border bg-surface px-3 py-2"
                            value={reason}
                            onChange={(event) => setReason(event.target.value)}
                          />
                        </label>
                      </div>
                      {!operatorName.trim() ? (
                        <p className="mt-2 text-xs text-muted">
                          A ação fica registrada no nome de quem aprova — informe o operador para confirmar.
                        </p>
                      ) : null}
                      <div className="mt-4 flex justify-end gap-2">
                        <Button variant="secondary" onClick={() => setOpenConfirmKey(null)}>Cancelar</Button>
                        <Button
                          variant="primary"
                          disabled={!operatorName.trim()}
                          loading={pendingKey === key}
                          onClick={() => void applyAction(item)}
                        >
                          Confirmar aplicação
                        </Button>
                      </div>
                    </div>
                  ) : null}

                  {blockingReasonText(item.blocking_reason) ? (
                    <div className="rounded-card border border-border bg-surface-2 px-4 py-3 text-sm text-muted">
                      <strong className="text-ink">Bloqueado:</strong> {blockingReasonText(item.blocking_reason)}
                    </div>
                  ) : null}
                </SurfaceCard>
              );
            })}
          </div>

          {selected ? (
            <SurfaceCard className="space-y-4 rounded-card">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-faint">Detalhe da linha</p>
                <h2 className="mt-2 text-xl font-semibold text-ink">{selected.title || selected.identity.provider_item_id}</h2>
              </div>
              <div className="grid gap-3 text-sm text-muted">
                <p><strong className="text-ink">Política:</strong> {selected.policy_id}</p>
                <p><strong className="text-ink">Observado no ML:</strong> {formatDateTime(selected.provider_observed_at) ?? "—"}</p>
                <p><strong className="text-ink">Observado no ERP:</strong> {formatDateTime(selected.internal_observed_at) ?? "—"}</p>
                <p><strong className="text-ink">Produto interno:</strong> {selected.internal_product_name || selected.internal_reference_code || "—"}</p>
                <p><strong className="text-ink">SKU:</strong> {selected.seller_sku || "—"}</p>
                <p><strong className="text-ink">EAN:</strong> {selected.ean || "—"}</p>
              </div>
            </SurfaceCard>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}
