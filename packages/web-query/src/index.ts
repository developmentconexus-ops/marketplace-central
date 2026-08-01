import { QueryClient } from "@tanstack/react-query";
import type { ActiveSourceName } from "@marketplace-central/sdk-runtime";
import { createElement, type ReactNode } from "react";

export const QUERY_STALE_TIME = {
  catalog: 300_000,
  stock: 45_000,
  pricecost: 120_000,
  listings: 45_000,
  mutations: 5_000,
  orders: 120_000,
  sync: 30_000,
  market: 300_000,
} as const;

export const queryKeyNamespaces = {
  catalog: ["catalog"] as const,
  inventory: ["inventory"] as const,
  linkage: ["linkage"] as const,
  profitability: ["profitability"] as const,
  listings: ["listings"] as const,
  mutations: ["mutations"] as const,
  orders: ["orders"] as const,
  sync: ["sync"] as const,
  market: ["market"] as const,
  installations: ["installations"] as const,
} as const;

export const catalogQueryKeys = {
  facts: (params: Record<string, unknown> = {}) => ["catalog", "facts", { params }] as const,
  counts: (erpSource: ActiveSourceName | undefined) =>
    ["catalog", "counts", { erp_source: erpSource ?? null }] as const,
  // The trailing params object is appended only when a caller supplies params,
  // so the pre-toggle two-arg key shape (["catalog","search",q]) stays byte-stable.
  search: (q: string, params?: Record<string, unknown>) =>
    (params && Object.keys(params).length > 0 ? ["catalog", "search", q, { params }] : ["catalog", "search", q]) as
      | readonly ["catalog", "search", string]
      | readonly ["catalog", "search", string, { params: Record<string, unknown> }],
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

export const listingsQueryKeys = {
  page: (installationId: string, filters: Record<string, unknown>) =>
    ["listings", "page", { installation_id: installationId, filters }] as const,
  byProduct: (installationId: string, filters: Record<string, unknown>) =>
    ["listings", "by-product", { installation_id: installationId, filters }] as const,
  detail: (listingId: string) => ["listings", "detail", listingId] as const,
  summary: (installationId: string) =>
    ["listings", "summary", { installation_id: installationId }] as const,
};

export const mutationsQueryKeys = {
  list: (installationId: string, filters: Record<string, unknown>) =>
    ["mutations", "list", { installation_id: installationId, filters }] as const,
  detail: (protocolId: string) => ["mutations", "detail", protocolId] as const,
  items: (protocolId: string) => ["mutations", "items", protocolId] as const,
};

export const ordersQueryKeys = {
  list: (installationId: string, filters: Record<string, unknown>) =>
    ["orders", "list", { installation_id: installationId, filters }] as const,
};

export const syncQueryKeys = {
  runs: (installationId: string, filters: Record<string, unknown>) =>
    ["sync", "runs", { installation_id: installationId, filters }] as const,
};

export const installationsQueryKeys = {
  list: () => ["installations", "list"] as const,
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

/**
 * pt-BR date + time for a timestamp the operator reads as a point in history
 * (an import, a run). Returns null for absent/unparseable input so the caller
 * renders an honest unknown instead of "Invalid Date" or the raw ISO string.
 */
export function formatDateTime(value: string | null | undefined): string | null {
  if (!value) {
    return null;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return date.toLocaleString("pt-BR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
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

export { invalidateAfterMutation, UnknownMutationInvalidationTypeError, type MutationInvalidationType } from "./invalidation";
export { failureCodes, failureCopy, type FailureCode } from "./failureCopy";
export {
  type ActiveSourceClient,
  activeSourceQueryKeys,
  useActiveSourceQuery,
  useSetActiveSourceMutation,
} from "./activeSource";
export {
  type SellableAssortmentClient,
  sellableAssortmentQueryKeys,
  useCatalogAssortmentCountsQuery,
  useSellableAssortmentQuery,
  useSetSellableAssortmentMutation,
} from "./activeSource";
export { type SyncHealthClient, syncHealthQueryKeys, useSyncHealthQuery } from "./syncHealth";
