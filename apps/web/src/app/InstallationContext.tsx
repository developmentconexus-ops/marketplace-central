import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { installationsQueryKeys } from "@marketplace-central/web-query";
import { createContext, useCallback, useContext, useEffect, useMemo, type ReactNode } from "react";
import { useSearchParams } from "react-router-dom";
import { useClient } from "./ClientContext";

type InstallationStatus = "loading" | "ready" | "empty" | "error";

interface InstallationContextValue {
  installationId: string;
  setInstallationId: (installationId: string) => void;
  installations: IntegrationInstallation[];
  status: InstallationStatus;
}

const InstallationContext = createContext<InstallationContextValue | null>(null);

export function InstallationProvider({ children }: { children: ReactNode }) {
  const client = useClient();
  const [searchParams, setSearchParams] = useSearchParams();
  const query = useQuery({
    queryKey: installationsQueryKeys.list(),
    queryFn: () => client.listIntegrationInstallations(),
  });
  const installations = query.data?.items ?? [];
  const requestedInstallationId = searchParams.get("installation");
  const requestedInstallationExists = installations.some(
    (installation) =>
      installation.installation_id === requestedInstallationId &&
      installation.status !== "draft" &&
      installation.status !== "pending_connection",
  );
  // Row existence is not an account. "Conectar" creates the installation BEFORE
  // the operator authorizes anything, so an abandoned authorization leaves a
  // pending_connection row behind — and with only that row the workspace opened
  // and reported 0 sales, 0 listings, 0 sync errors and "tudo em dia", an
  // all-clear nothing had measured. Only a state that was authorized at least
  // once counts; degraded and requires_reauth accounts still hold real data and
  // render their own warnings, so they keep the workspace open. Among those,
  // default to a connected one — pending rows sort ahead of the real accounts,
  // so taking the first item opened the workspace on an empty account.
  const authorizedInstallations = installations.filter(
    (installation) =>
      installation.status !== "draft" && installation.status !== "pending_connection",
  );
  const defaultInstallation =
    authorizedInstallations.find((installation) => installation.status === "connected") ??
    authorizedInstallations[0];
  const installationId = requestedInstallationExists
    ? (requestedInstallationId ?? "")
    : (defaultInstallation?.installation_id ?? "");

  useEffect(() => {
    if (!query.isSuccess || !defaultInstallation || requestedInstallationExists) return;
    const nextParams = new URLSearchParams(searchParams);
    nextParams.set("installation", defaultInstallation.installation_id);
    setSearchParams(nextParams, { replace: true });
  }, [
    defaultInstallation,
    installations,
    query.isSuccess,
    requestedInstallationExists,
    searchParams,
    setSearchParams,
  ]);

  const setInstallationId = useCallback(
    (nextInstallationId: string) => {
      const nextParams = new URLSearchParams(searchParams);
      nextParams.set("installation", nextInstallationId);
      setSearchParams(nextParams);
    },
    [searchParams, setSearchParams],
  );

  let status: InstallationStatus = "loading";
  if (query.isError) status = "error";
  else if (query.isSuccess && authorizedInstallations.length === 0) status = "empty";
  else if (query.isSuccess) status = "ready";

  const value = useMemo(
    () => ({ installationId, setInstallationId, installations, status }),
    [installationId, installations, setInstallationId, status],
  );

  return <InstallationContext.Provider value={value}>{children}</InstallationContext.Provider>;
}

export function InstallationGate({ children }: { children: ReactNode }) {
  const { status } = useInstallation();
  const queryClient = useQueryClient();

  if (status === "loading") return <LoadingState />;
  if (status === "error") {
    return (
      <ErrorState
        onRetry={() => void queryClient.refetchQueries({ queryKey: installationsQueryKeys.list() })}
      />
    );
  }
  if (status === "empty") return <EmptyState hint="Conecte uma conta em Integrações" />;

  return <>{children}</>;
}

export function useInstallation(): InstallationContextValue {
  const context = useContext(InstallationContext);
  if (!context) throw new Error("useInstallation must be used inside <InstallationProvider>");
  return context;
}
