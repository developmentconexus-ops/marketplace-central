import { useQuery } from "@tanstack/react-query";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";

export const erpImportChainQueryKeys = {
  detail: (importId: string) => ["erp-imports", "chain", importId] as const,
};

function errorStatus(error: unknown): number | undefined {
  if (typeof error !== "object" || error === null || !("status" in error)) return undefined;
  const status = error.status;
  return typeof status === "number" ? status : undefined;
}

export function useErpImportChain(importId: string) {
  const client = useClient();

  return useQuery({
    queryKey: erpImportChainQueryKeys.detail(importId),
    queryFn: () => client.getErpImportChain(importId),
    enabled: Boolean(importId),
    staleTime: QUERY_STALE_TIME.listings,
    retry: (failureCount, error) => {
      // A 4xx is a settled answer; only transient failures get the default single retry.
      const status = errorStatus(error);
      if (status !== undefined && status >= 400 && status < 500) return false;
      return failureCount < 1;
    },
  });
}
