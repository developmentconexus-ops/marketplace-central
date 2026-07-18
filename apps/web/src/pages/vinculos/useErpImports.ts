import { useQuery } from "@tanstack/react-query";
import { QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";

/**
 * PAGE-LOCAL query keys for the read-only "Importação" history section
 * (S9). This is protocol read-only: no upload mutation lives here.
 */
export const erpImportsQueryKeys = {
  list: () => ["erp-imports", "list"] as const,
  detail: (importId: string) => ["erp-imports", "detail", importId] as const,
};

export function useErpImportsList() {
  const client = useClient();

  return useQuery({
    queryKey: erpImportsQueryKeys.list(),
    queryFn: () => client.listErpImports(),
    staleTime: QUERY_STALE_TIME.listings,
  });
}

export function useErpImportDetail(importId: string | undefined) {
  const client = useClient();

  return useQuery({
    queryKey: erpImportsQueryKeys.detail(importId ?? ""),
    queryFn: () => client.getErpImport(importId as string),
    enabled: Boolean(importId),
    staleTime: QUERY_STALE_TIME.listings,
  });
}
