import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { installationsQueryKeys } from "@marketplace-central/web-query";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  type ReactNode,
} from "react";
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
    (installation) => installation.installation_id === requestedInstallationId,
  );
  const installationId = requestedInstallationExists
    ? requestedInstallationId ?? ""
    : installations[0]?.installation_id ?? "";

  useEffect(() => {
    if (!query.isSuccess || installations.length === 0 || requestedInstallationExists) return;
    const nextParams = new URLSearchParams(searchParams);
    nextParams.set("installation", installations[0].installation_id);
    setSearchParams(nextParams, { replace: true });
  }, [
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
  else if (query.isSuccess && installations.length === 0) status = "empty";
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
        onRetry={() =>
          void queryClient.refetchQueries({ queryKey: installationsQueryKeys.list() })
        }
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
