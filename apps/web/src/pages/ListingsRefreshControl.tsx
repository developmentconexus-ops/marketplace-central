import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { queryKeyNamespaces, syncQueryKeys } from "@marketplace-central/web-query";
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

function attachedRunId(error: unknown): string | null {
  if (typeof error !== "object" || error === null) return null;
  const candidate = error as {
    status?: number;
    error?: { code?: string; details?: { operation_run_id?: unknown } };
  };
  const operationRunId = candidate.error?.details?.operation_run_id;
  return candidate.status === 409 &&
    candidate.error?.code === "refresh_in_progress" &&
    typeof operationRunId === "string"
    ? operationRunId
    : null;
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

  const runsQuery = useQuery({
    queryKey: syncQueryKeys.runs(installationId, { operation_run_id: observedRunId }),
    queryFn: () => client.listIntegrationOperationRuns(installationId),
    enabled: observedRunId !== null,
    refetchInterval: (query) => {
      const observedRun = query.state.data?.items.find(
        (item) => item.operation_run_id === observedRunId,
      );
      return observedRun?.status === "succeeded" ||
        observedRun?.status === "failed" ||
        observedRun?.status === "cancelled"
        ? false
        : 2_000;
    },
  });

  const observedRun = observedRunId === null
    ? undefined
    : runsQuery.data?.items.find((item) => item.operation_run_id === observedRunId);
  const status: RunStatus | null = observedRunId === null ? null : observedRun?.status ?? "queued";

  useEffect(() => {
    if (
      observedRunId !== null &&
      status === "succeeded" &&
      invalidatedRunId.current !== observedRunId
    ) {
      invalidatedRunId.current = observedRunId;
      void queryClient.invalidateQueries({ queryKey: queryKeyNamespaces.listings });
    }
  }, [observedRunId, queryClient, status]);

  const runIsActive = status === "queued" || status === "running";
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
      {status !== null ? <small className="text-muted">{statusCopy[status]}</small> : null}
      {startError !== null ? <small role="alert" className="text-warn">{startError}</small> : null}
      {terminalError !== null ? <small role="alert" className="text-warn">{terminalError}</small> : null}
    </div>
  );
}
