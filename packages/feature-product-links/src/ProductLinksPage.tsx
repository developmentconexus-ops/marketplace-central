import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { Button, SurfaceCard } from "@marketplace-central/ui";
import type {
  IntegrationInstallation,
  ProductLinkActor,
  ProductLinkCandidateItem,
  ProductLinkState,
  ProductLinkWorkflowItem,
} from "@marketplace-central/sdk-runtime";

export interface ProductLinksClient {
  listIntegrationInstallations: () => Promise<{ items: IntegrationInstallation[] }>;
  listProductLinkWorkflows: (installationId: string, limit?: number) => Promise<{ items: ProductLinkWorkflowItem[] }>;
  approveProductLinkCandidate: (req: {
    candidate_id: string;
    reason?: string;
    actor: ProductLinkActor;
  }) => Promise<unknown>;
  rejectProductLinkListing: (req: {
    installation_id: string;
    provider_code: string;
    provider_item_id: string;
    provider_variation_id?: string;
    reason?: string;
    actor: ProductLinkActor;
  }) => Promise<unknown>;
  manualResolveProductLink: (req: {
    installation_id: string;
    provider_code: string;
    provider_item_id: string;
    provider_variation_id?: string;
    internal_product_id: number;
    internal_product_name?: string;
    internal_reference_code?: string;
    reason?: string;
    actor: ProductLinkActor;
  }) => Promise<unknown>;
}

export interface ProductLinksPageProps {
  client: ProductLinksClient;
}

type LoadState = "loading" | "ready" | "error";

type ManualFormState = {
  internal_product_id: string;
  internal_product_name: string;
  internal_reference_code: string;
  reason: string;
};

const stateTone: Record<ProductLinkState, string> = {
  none: "bg-slate-100 text-slate-700 border-slate-200",
  unresolved: "bg-amber-100 text-amber-800 border-amber-200",
  conflict: "bg-rose-100 text-rose-800 border-rose-200",
  resolved: "bg-emerald-100 text-emerald-800 border-emerald-200",
  rejected: "bg-slate-200 text-slate-700 border-slate-300",
};

function formatStateLabel(state: ProductLinkState): string {
  switch (state) {
    case "unresolved":
      return "Unresolved";
    case "conflict":
      return "Conflict";
    case "resolved":
      return "Resolved";
    case "rejected":
      return "Rejected";
    default:
      return "Pending";
  }
}

function inferWorkflowState(item: ProductLinkWorkflowItem): ProductLinkState {
  if (item.current_link) {
    return item.current_link.state;
  }
  if (item.candidates.some((candidate) => candidate.state === "conflict")) {
    return "conflict";
  }
  return "unresolved";
}

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

function makeDefaultManualForm(item: ProductLinkWorkflowItem): ManualFormState {
  const preferredCandidate = item.candidates.find((candidate) => candidate.internal_product_id != null);
  return {
    internal_product_id: preferredCandidate?.internal_product_id != null ? String(preferredCandidate.internal_product_id) : "",
    internal_product_name: preferredCandidate?.internal_product_name ?? item.current_link?.internal_product_name ?? "",
    internal_reference_code: preferredCandidate?.internal_reference_code ?? item.current_link?.internal_reference_code ?? "",
    reason: "",
  };
}

function buildActor(actorName: string): ProductLinkActor {
  const trimmed = actorName.trim();
  const actorId = trimmed.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "");
  return {
    actor_type: "operator",
    actor_id: actorId || "operator",
    actor_name: trimmed || undefined,
  };
}

function getWorkflowProviderCode(item: ProductLinkWorkflowItem): string {
  return item.current_link?.provider_code ?? item.candidates[0]?.provider_code ?? "";
}

function getWorkflowTitle(item: ProductLinkWorkflowItem): string {
  const identity = item.identity;
  return identity.provider_variation_id
    ? `${identity.provider_item_id} / ${identity.provider_variation_id}`
    : identity.provider_item_id;
}

export function ProductLinksPage({ client }: ProductLinksPageProps) {
  const [searchParams, setSearchParams] = useSearchParams();
  const selectedInstallationId = searchParams.get("installation") ?? "";
  const [installations, setInstallations] = useState<IntegrationInstallation[]>([]);
  const [items, setItems] = useState<ProductLinkWorkflowItem[]>([]);
  const [state, setState] = useState<LoadState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [operatorName, setOperatorName] = useState("Leandro");
  const [actionError, setActionError] = useState<string | null>(null);
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [openManualKey, setOpenManualKey] = useState<string | null>(null);
  const [manualForms, setManualForms] = useState<Record<string, ManualFormState>>({});

  useEffect(() => {
    let cancelled = false;

    async function loadInstallations() {
      setState("loading");
      setError(null);
      try {
        const result = await client.listIntegrationInstallations();
        if (cancelled) {
          return;
        }
        setInstallations(result.items);
        if (!selectedInstallationId && result.items[0]) {
          setSearchParams({ installation: result.items[0].installation_id }, { replace: true });
          return;
        }
        setState("ready");
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
  }, [client, selectedInstallationId, setSearchParams]);

  useEffect(() => {
    if (!selectedInstallationId) {
      setItems([]);
      return;
    }

    let cancelled = false;

    async function loadWorkflows() {
      setState("loading");
      setError(null);
      try {
        const result = await client.listProductLinkWorkflows(selectedInstallationId, 50);
        if (cancelled) {
          return;
        }
        setItems(result.items);
        setManualForms(
          Object.fromEntries(
            result.items.map((item) => [getWorkflowTitle(item), makeDefaultManualForm(item)]),
          ),
        );
        setState("ready");
      } catch (loadError) {
        if (cancelled) {
          return;
        }
        setError(normalizeError(loadError, "Failed to load product link workflows."));
        setState("error");
      }
    }

    void loadWorkflows();
    return () => {
      cancelled = true;
    };
  }, [client, selectedInstallationId]);

  const summary = useMemo(() => {
    return items.reduce(
      (acc, item) => {
        const workflowState = inferWorkflowState(item);
        acc.total += 1;
        acc[workflowState] += 1;
        return acc;
      },
      { total: 0, unresolved: 0, conflict: 0, resolved: 0, rejected: 0 } as Record<"total" | "unresolved" | "conflict" | "resolved" | "rejected", number>,
    );
  }, [items]);

  async function refreshWorkflows() {
    if (!selectedInstallationId) {
      return;
    }
    const result = await client.listProductLinkWorkflows(selectedInstallationId, 50);
    setItems(result.items);
    setManualForms(
      Object.fromEntries(
        result.items.map((item) => [getWorkflowTitle(item), manualForms[getWorkflowTitle(item)] ?? makeDefaultManualForm(item)]),
      ),
    );
  }

  async function runAction(actionKey: string, action: () => Promise<unknown>) {
    setPendingKey(actionKey);
    setActionError(null);
    try {
      await action();
      await refreshWorkflows();
    } catch (runError) {
      setActionError(normalizeError(runError, "Failed to persist link decision."));
    } finally {
      setPendingKey(null);
    }
  }

  function updateManualForm(key: string, patch: Partial<ManualFormState>) {
    setManualForms((current) => ({
      ...current,
      [key]: {
        ...current[key],
        ...patch,
      },
    }));
  }

  const selectedInstallation = installations.find((item) => item.installation_id === selectedInstallationId);
  const actor = buildActor(operatorName);

  return (
    <div className="space-y-6">
      <section className="rounded-[28px] border border-slate-200 bg-[linear-gradient(135deg,#0f172a_0%,#1e293b_55%,#334155_100%)] px-6 py-7 text-white shadow-sm">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div className="max-w-2xl space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-200">Product Links</p>
            <h1 className="text-3xl font-semibold tracking-tight">Operator workflow for real listing resolution</h1>
            <p className="text-sm leading-6 text-slate-200">
              Candidates stay visible as evidence. Operators explicitly approve, reject, or manually resolve the listing truth that downstream modules will trust.
            </p>
          </div>
          <div className="grid gap-3 rounded-2xl border border-white/15 bg-white/8 p-4 sm:grid-cols-2">
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Installation</span>
              <select
                aria-label="Installation"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white"
                value={selectedInstallationId}
                onChange={(event) => setSearchParams({ installation: event.target.value })}
              >
                {installations.map((installation) => (
                  <option key={installation.installation_id} value={installation.installation_id} className="text-slate-900">
                    {installation.display_name}
                  </option>
                ))}
              </select>
            </label>
            <label className="space-y-1 text-sm">
              <span className="block text-slate-300">Operator</span>
              <input
                aria-label="Operator name"
                className="w-full rounded-xl border border-white/15 bg-slate-950/40 px-3 py-2 text-sm text-white placeholder:text-slate-400"
                placeholder="Operator name"
                value={operatorName}
                onChange={(event) => setOperatorName(event.target.value)}
              />
            </label>
          </div>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-4">
        {[
          { label: "Total workflows", value: summary.total },
          { label: "Unresolved", value: summary.unresolved + summary.conflict },
          { label: "Resolved", value: summary.resolved },
          { label: "Rejected", value: summary.rejected },
        ].map((card) => (
          <SurfaceCard key={card.label} className="space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{card.label}</p>
            <p className="text-3xl font-semibold text-slate-900">{card.value}</p>
          </SurfaceCard>
        ))}
      </section>

      {actionError ? (
        <SurfaceCard className="border-rose-200 bg-rose-50 text-rose-900">
          <p className="text-sm font-medium">{actionError}</p>
        </SurfaceCard>
      ) : null}

      {state === "loading" ? (
        <SurfaceCard>
          <p className="text-sm text-slate-600">Loading product link workflows...</p>
        </SurfaceCard>
      ) : null}

      {state === "error" && error ? (
        <SurfaceCard className="border-rose-200 bg-rose-50">
          <p className="text-sm font-medium text-rose-900">{error}</p>
        </SurfaceCard>
      ) : null}

      {state === "ready" && installations.length === 0 ? (
        <SurfaceCard>
          <h2 className="text-lg font-semibold text-slate-900">No installations available</h2>
          <p className="mt-2 text-sm text-slate-600">Connect a marketplace installation first so we can load live listing identities and candidate evidence.</p>
        </SurfaceCard>
      ) : null}

      {state === "ready" && selectedInstallation && items.length === 0 ? (
        <SurfaceCard>
          <h2 className="text-lg font-semibold text-slate-900">No workflows yet</h2>
          <p className="mt-2 text-sm text-slate-600">
            {selectedInstallation.display_name} does not have imported listing identities ready for operator review yet.
          </p>
        </SurfaceCard>
      ) : null}

      {state === "ready" && items.length > 0 ? (
        <div className="grid gap-4">
          {items.map((item) => {
            const workflowState = inferWorkflowState(item);
            const providerCode = getWorkflowProviderCode(item);
            const workflowKey = getWorkflowTitle(item);
            const manualForm = manualForms[workflowKey] ?? makeDefaultManualForm(item);
            const resolvableCandidates = item.candidates.filter((candidate) => candidate.internal_product_id != null);

            return (
              <SurfaceCard key={workflowKey} className="space-y-5">
                <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                  <div className="space-y-2">
                    <div className="flex items-center gap-3">
                      <h2 className="text-lg font-semibold text-slate-900">{workflowKey}</h2>
                      <span className={`inline-flex rounded-full border px-3 py-1 text-xs font-semibold ${stateTone[workflowState]}`}>
                        {formatStateLabel(workflowState)}
                      </span>
                    </div>
                    <p className="text-sm text-slate-600">
                      Provider: <span className="font-medium text-slate-900">{providerCode || "unknown"}</span>
                      {" · "}
                      Candidates: <span className="font-medium text-slate-900">{item.candidates.length}</span>
                      {" · "}
                      Audit events: <span className="font-medium text-slate-900">{item.audit.length}</span>
                    </p>
                    {item.current_link ? (
                      <p className="text-sm text-slate-700">
                        Current truth:{" "}
                        <span className="font-medium">
                          {item.current_link.internal_product_name || item.current_link.internal_reference_code || "No internal product name"}
                        </span>
                        {item.current_link.internal_product_id != null ? ` (#${item.current_link.internal_product_id})` : ""}
                      </p>
                    ) : (
                      <p className="text-sm text-slate-700">Current truth: no persisted operator decision yet.</p>
                    )}
                  </div>

                  <div className="flex flex-wrap gap-2">
                    {resolvableCandidates.slice(0, 2).map((candidate) => (
                      <Button
                        key={candidate.candidate_id}
                        variant="primary"
                        loading={pendingKey === `approve:${candidate.candidate_id}`}
                        onClick={() =>
                          runAction(`approve:${candidate.candidate_id}`, () =>
                            client.approveProductLinkCandidate({
                              candidate_id: candidate.candidate_id,
                              reason: candidate.match_value ? `Approved from ${candidate.match_input}: ${candidate.match_value}` : "Approved from generated candidate",
                              actor,
                            }),
                          )
                        }
                      >
                        Approve {candidate.internal_reference_code || candidate.internal_product_name || candidate.candidate_id}
                      </Button>
                    ))}
                    <Button
                      variant="secondary"
                      onClick={() => setOpenManualKey(openManualKey === workflowKey ? null : workflowKey)}
                    >
                      Manual resolve
                    </Button>
                    <Button
                      variant="danger"
                      loading={pendingKey === `reject:${workflowKey}`}
                      onClick={() =>
                        runAction(`reject:${workflowKey}`, () =>
                          client.rejectProductLinkListing({
                            installation_id: item.identity.installation_id,
                            provider_code: providerCode,
                            provider_item_id: item.identity.provider_item_id,
                            provider_variation_id: item.identity.provider_variation_id,
                            reason: "Operator rejected listing identity for linking",
                            actor,
                          }),
                        )
                      }
                    >
                      Reject listing
                    </Button>
                  </div>
                </div>

                <div className="grid gap-4 xl:grid-cols-[1.2fr_0.8fr]">
                  <div className="space-y-3">
                    <h3 className="text-sm font-semibold uppercase tracking-[0.18em] text-slate-500">Candidate evidence</h3>
                    {item.candidates.length === 0 ? (
                      <p className="rounded-2xl border border-dashed border-slate-200 px-4 py-5 text-sm text-slate-500">
                        No candidate evidence recorded for this listing identity.
                      </p>
                    ) : (
                      <div className="grid gap-3">
                        {item.candidates.map((candidate: ProductLinkCandidateItem) => (
                          <article key={candidate.candidate_id} className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4">
                            <div className="flex flex-wrap items-center gap-2">
                              <span className="text-sm font-semibold text-slate-900">
                                {candidate.internal_product_name || candidate.internal_reference_code || "No internal product suggested"}
                              </span>
                              <span className="rounded-full bg-white px-2 py-1 text-xs font-medium text-slate-600">{candidate.state}</span>
                              <span className="rounded-full bg-white px-2 py-1 text-xs font-medium text-slate-600">{candidate.match_input}</span>
                            </div>
                            <p className="mt-2 text-sm text-slate-600">
                              Internal product ID: {candidate.internal_product_id != null ? candidate.internal_product_id : "not suggested"}
                              {candidate.match_value ? ` · Match value: ${candidate.match_value}` : ""}
                            </p>
                          </article>
                        ))}
                      </div>
                    )}
                  </div>

                  <div className="space-y-3">
                    <h3 className="text-sm font-semibold uppercase tracking-[0.18em] text-slate-500">Audit trail</h3>
                    {item.audit.length === 0 ? (
                      <p className="rounded-2xl border border-dashed border-slate-200 px-4 py-5 text-sm text-slate-500">
                        No audit events yet. The first operator action will establish persisted link truth.
                      </p>
                    ) : (
                      <div className="grid gap-3">
                        {item.audit.slice(0, 4).map((entry) => (
                          <article key={entry.audit_id} className="rounded-2xl border border-slate-200 bg-white px-4 py-4">
                            <p className="text-sm font-semibold text-slate-900">{entry.action}</p>
                            <p className="mt-1 text-sm text-slate-600">
                              {formatStateLabel(entry.previous_state)} to {formatStateLabel(entry.next_state)}
                            </p>
                            <p className="mt-1 text-xs text-slate-500">
                              {entry.actor.actor_name || entry.actor.actor_id} · {new Date(entry.created_at).toLocaleString("pt-BR")}
                            </p>
                          </article>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                {openManualKey === workflowKey ? (
                  <div className="rounded-[24px] border border-slate-200 bg-slate-50 p-4">
                    <div className="grid gap-3 md:grid-cols-2">
                      <label className="space-y-1 text-sm text-slate-700">
                        <span className="block font-medium">Internal product ID</span>
                        <input
                          aria-label={`Internal product ID ${workflowKey}`}
                          className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                          value={manualForm.internal_product_id}
                          onChange={(event) => updateManualForm(workflowKey, { internal_product_id: event.target.value })}
                        />
                      </label>
                      <label className="space-y-1 text-sm text-slate-700">
                        <span className="block font-medium">Reference code</span>
                        <input
                          aria-label={`Reference code ${workflowKey}`}
                          className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                          value={manualForm.internal_reference_code}
                          onChange={(event) => updateManualForm(workflowKey, { internal_reference_code: event.target.value })}
                        />
                      </label>
                      <label className="space-y-1 text-sm text-slate-700 md:col-span-2">
                        <span className="block font-medium">Internal product name</span>
                        <input
                          aria-label={`Internal product name ${workflowKey}`}
                          className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                          value={manualForm.internal_product_name}
                          onChange={(event) => updateManualForm(workflowKey, { internal_product_name: event.target.value })}
                        />
                      </label>
                      <label className="space-y-1 text-sm text-slate-700 md:col-span-2">
                        <span className="block font-medium">Reason</span>
                        <input
                          aria-label={`Reason ${workflowKey}`}
                          className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2"
                          value={manualForm.reason}
                          onChange={(event) => updateManualForm(workflowKey, { reason: event.target.value })}
                        />
                      </label>
                    </div>
                    <div className="mt-4 flex justify-end">
                      <Button
                        variant="primary"
                        loading={pendingKey === `manual:${workflowKey}`}
                        onClick={() =>
                          runAction(`manual:${workflowKey}`, async () => {
                            const internalProductID = Number(manualForm.internal_product_id);
                            if (!Number.isInteger(internalProductID) || internalProductID <= 0) {
                              throw new Error("Internal product ID must be a positive integer.");
                            }
                            await client.manualResolveProductLink({
                              installation_id: item.identity.installation_id,
                              provider_code: providerCode,
                              provider_item_id: item.identity.provider_item_id,
                              provider_variation_id: item.identity.provider_variation_id,
                              internal_product_id: internalProductID,
                              internal_product_name: manualForm.internal_product_name,
                              internal_reference_code: manualForm.internal_reference_code,
                              reason: manualForm.reason || "Manual operator resolution",
                              actor,
                            });
                            setOpenManualKey(null);
                          })
                        }
                      >
                        Save manual resolution
                      </Button>
                    </div>
                  </div>
                ) : null}
              </SurfaceCard>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}
