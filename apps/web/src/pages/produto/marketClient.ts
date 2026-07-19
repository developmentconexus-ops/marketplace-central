import { useMemo } from "react";
import { createMarketPriceIntelClient } from "@marketplace-central/sdk-runtime";

/** Return type of the standalone IC-03 client — the app's useClient() carries no market seam. */
export type ProdutoMarketClient = ReturnType<typeof createMarketPriceIntelClient>;

/**
 * Resolve the IC-03 base URL from the build env, mirroring MarketComparison.tsx
 * (apps/web/src/pages/precos/MarketComparison.tsx:21-25) and useClient's rule.
 */
export function apiBaseUrl(): string {
  const env = import.meta.env.VITE_API_BASE_URL;
  if (env && env.trim()) return env;
  return import.meta.env.DEV ? "http://localhost:8080" : "";
}

/** Injectable for tests; defaults to a standalone IC-03 client. */
export function useProdutoMarketClient(injected?: ProdutoMarketClient): ProdutoMarketClient {
  return useMemo<ProdutoMarketClient>(
    () => injected ?? createMarketPriceIntelClient({ baseUrl: apiBaseUrl() }),
    [injected],
  );
}
