import { createContext, useContext, useMemo, type ReactNode } from "react";
import { createMarketplaceCentralClient } from "@marketplace-central/sdk-runtime";
import { createRefreshableFetch, type WithNoCache } from "@marketplace-central/web-query";

export type Client = ReturnType<typeof createMarketplaceCentralClient> & { withNoCache: WithNoCache };

const ClientContext = createContext<Client | null>(null);

export function resolveApiBaseUrl(envBaseUrl: string | undefined, isDev: boolean): string {
  if (envBaseUrl && envBaseUrl.trim()) {
    return envBaseUrl;
  }
  if (isDev) {
    return "http://localhost:8080";
  }
  return "";
}

export function ClientProvider({ children }: { children: ReactNode }) {
  const baseUrl = resolveApiBaseUrl(import.meta.env.VITE_API_BASE_URL, import.meta.env.DEV);
  const refreshableFetch = useMemo(() => createRefreshableFetch(), []);
  const client = useMemo(() => ({
    ...createMarketplaceCentralClient({ baseUrl, fetchImpl: refreshableFetch.fetchImpl }),
    withNoCache: refreshableFetch.withNoCache,
  }), [baseUrl, refreshableFetch]);
  return <ClientContext.Provider value={client}>{children}</ClientContext.Provider>;
}

export function useClient(): Client {
  const ctx = useContext(ClientContext);
  if (!ctx) throw new Error("useClient must be used inside <ClientProvider>");
  return ctx;
}
