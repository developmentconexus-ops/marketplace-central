import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isRefreshInProgressError } from "@marketplace-central/sdk-runtime";
import { formatRelativeAge, queryKeyNamespaces, syncQueryKeys } from "@marketplace-central/web-query";
import { useEffect, useRef, useState } from "react";
import { useClient } from "../app/ClientContext";

type RunStatus = "queued" | "running" | "succeeded" | "failed" | "cancelled";

const statusCopy: Record<RunStatus, string> = {
  queued: "na fila",
  running: "em andamento",
  succeeded: "concluído",
  failed: "falhou",
  cancelled: "cancelado",
};

const REFRESH_OPERATION_TYPE = "listings_refresh";

function isRefreshRun(run: { operation_type: string }): boolean {
  return run.operation_type === REFRESH_OPERATION_TYPE;
}

function isActive(status: string | undefined): boolean {
  return status === "queued" || status === "running";
}


function attachedRunId(error: unknown): string | null {
  return isRefreshInProgressError(error) ? error.details.operation_run_id : null;
}

export function ListingsRefreshControl({ installationId }: { installationId: string }) {
  const client = useClient();
  const queryClient = useQueryClient();
  const [observedRunId, setObservedRunId] = useState<string | null>(null);
  const [startError, setStartError] = useState<string | null>(null);
  const currentInstallationId = useRef(installationId);
  const invalidatedRunId = useRef<string | null>(null);

  useEffect(() => {
    if (currentInstallationId.current !== installationId) {
      currentInstallationId.current = installationId;
      invalidatedRunId.current = null;
      setObservedRunId(null);
      setStartError(null);
    }
  }, [installationId]);

  const refreshMutation = useMutation({
    mutationFn: (requestedInstallationId: string) =>
      client.refreshListings({ installation_id: requestedInstallationId }),
    onSuccess: (accepted, requestedInstallationId) => {
      if (requestedInstallationId !== currentInstallationId.current) return;
      setStartError(null);
      setObservedRunId(accepted.operation_run_id);
    },
    onError: (error, requestedInstallationId) => {
      if (requestedInstallationId !== currentInstallationId.current) return;
      const operationRunId = attachedRunId(error);
      if (operationRunId !== null) {
        setStartError(null);
        setObservedRunId(operationRunId);
        return;
      }
      setObservedRunId(null);
      setStartError("Falha ao iniciar atualização.");
    },
  });

  // A full pull of a large account takes minutes. The run lives server-side, so
  // the control watches the runs list even before this tab started one: a reload
  // (or landing here from another screen) during a refresh must show the run in
  // progress instead of an idle button that 409s on the next click.
  const runsQuery = useQuery({
    queryKey: syncQueryKeys.runs(installationId, { operation_run_id: observedRunId }),
    queryFn: () => client.listIntegrationOperationRuns(installationId),
    enabled: installationId !== "",
    // React Query pauses polling while the document is hidden, and the app
    // disables refetch-on-focus — so a user who switched tabs mid-refresh came
    // back to a status pill frozen at "na fila" and listings that never
    // invalidated. The run keeps going server-side; keep watching it.
    refetchIntervalInBackground: true,
    refetchInterval: (query) => {
      const items = query.state.data?.items ?? [];
      if (observedRunId !== null) {
        const observedRun = items.find((item) => item.operation_run_id === observedRunId);
        return isActive(observedRun?.status) || observedRun === undefined ? 2_000 : false;
      }
      // Nothing started here: keep polling only while somebody else's run is live.
      return items.some((item) => isRefreshRun(item) && isActive(item.status)) ? 2_000 : false;
    },
  });

  const runs = runsQuery.data?.items ?? [];
  // A run this tab did not start still belongs to this account — adopt it so the
  // button is disabled and the pill honest while it finishes.
  const adoptedRunId =
    runs.find((item) => isRefreshRun(item) && isActive(item.status))?.operation_run_id ?? null;
  const trackedRunId = observedRunId ?? adoptedRunId;
  const observedRun = trackedRunId === null
    ? undefined
    : runs.find((item) => item.operation_run_id === trackedRunId);
  const status: RunStatus | null = trackedRunId === null ? null : observedRun?.status ?? "queued";

  useEffect(() => {
    if (
      trackedRunId !== null &&
      status === "succeeded" &&
      invalidatedRunId.current !== trackedRunId
    ) {
      invalidatedRunId.current = trackedRunId;
      void queryClient.invalidateQueries({ queryKey: queryKeyNamespaces.listings });
    }
  }, [queryClient, status, trackedRunId]);

  const runIsActive = isActive(status ?? undefined);

  // Ticks only while a run is live, so the elapsed label stays truthful without
  // re-rendering an idle screen.
  const [now, setNow] = useState<number>(() => Date.now());
  useEffect(() => {
    if (!runIsActive) return;
    setNow(Date.now());
    const timer = setInterval(() => setNow(Date.now()), 1_000);
    return () => clearInterval(timer);
  }, [runIsActive]);
  const elapsed = runIsActive && observedRun?.started_at ? formatRelativeAge(observedRun.started_at, now) : null;

  const terminalError = status === "failed"
    ? "Atualização falhou."
    : status === "cancelled"
      ? "Atualização cancelada."
      : null;

  return (
    <div className="flex items-center gap-2">
      <button
        type="button"
        className="rounded-lg border px-3 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-50"
        disabled={refreshMutation.isPending || runIsActive}
        onClick={() => {
          setStartError(null);
          refreshMutation.mutate(installationId);
        }}
      >
        Atualizar
      </button>
      {status !== null ? (
        <small className="text-muted">
          {statusCopy[status]}
          {elapsed !== null ? ` ${elapsed}` : ""}
        </small>
      ) : null}
      {startError !== null ? <small role="alert" className="text-warn">{startError}</small> : null}
      {terminalError !== null ? <small role="alert" className="text-warn">{terminalError}</small> : null}
    </div>
  );
}
