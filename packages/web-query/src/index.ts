import { QueryClient } from "@tanstack/react-query";
import { createElement, type ReactNode } from "react";

export const QUERY_STALE_TIME = {
  catalog: 300_000,
  stock: 45_000,
  pricecost: 120_000,
} as const;

export const queryKeyNamespaces = {
  catalog: ["catalog"] as const,
  inventory: ["inventory"] as const,
  linkage: ["linkage"] as const,
  profitability: ["profitability"] as const,
} as const;

export const catalogQueryKeys = {
  facts: (params: Record<string, unknown> = {}) => ["catalog", "facts", { params }] as const,
  search: (q: string) => ["catalog", "search", q] as const,
};

export const inventoryQueryKeys = {
  risks: (installation_id: string, filters: Record<string, unknown>) =>
    ["inventory", { installation_id, filters }] as const,
};

export const linkageQueryKeys = {
  workflows: (installation_id: string) => ["linkage", { installation_id }] as const,
};

export const profitabilityQueryKeys = {
  marginInputs: (installation_id: string) => ["profitability", { installation_id }] as const,
};

export function createWebQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: 1,
        refetchOnWindowFocus: false,
      },
    },
  });
}

export function formatAsOf(asOf: string | null | undefined): string {
  if (!asOf) {
    return "dados de desconhecido";
  }
  const date = new Date(asOf);
  if (Number.isNaN(date.getTime())) {
    return "dados de desconhecido";
  }
  return `dados de ${date.toLocaleTimeString("pt-BR", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  })}`;
}

export type WithNoCache = <T>(operation: () => Promise<T>) => Promise<T>;

export interface RefreshableClient {
  withNoCache?: WithNoCache;
}

export function createRefreshableFetch(baseFetch: typeof fetch = fetch): {
  fetchImpl: typeof fetch;
  withNoCache: WithNoCache;
} {
  let noCacheDepth = 0;

  const fetchImpl: typeof fetch = (input, init) => {
    if (noCacheDepth === 0 || (init?.method && init.method.toUpperCase() !== "GET")) {
      return baseFetch(input, init);
    }
    const headers = new Headers(init?.headers);
    headers.set("Cache-Control", "no-cache");
    return baseFetch(input, { ...init, headers });
  };

  const withNoCache: WithNoCache = async <T>(operation: () => Promise<T>) => {
    noCacheDepth += 1;
    try {
      return await operation();
    } finally {
      noCacheDepth -= 1;
    }
  };

  return { fetchImpl, withNoCache };
}

export function FreshnessIndicator({ asOf }: { asOf: string | null | undefined }): ReactNode {
  return createElement("span", { "aria-label": "Data freshness" }, formatAsOf(asOf));
}
