import { useMemo } from "react";
import { createMarketPriceIntelClient } from "@marketplace-central/sdk-runtime";

/** Return type of the standalone IC-03 client — the app's useClient() carries no market seam. */
export type MercadoMarketClient = ReturnType<typeof createMarketPriceIntelClient>;

/**
 * Resolve the IC-03 base URL from the build env, mirroring MarketComparison.tsx
 * and produto/marketClient.ts (the app's central useClient() exposes no market
 * seam, so the radar talks to IC-03 over its own standalone HTTP client).
 */
export function apiBaseUrl(): string {
  const env = import.meta.env.VITE_API_BASE_URL;
  if (env && env.trim()) return env;
  return import.meta.env.DEV ? "http://localhost:8080" : "";
}

/** Injectable for tests; defaults to a standalone IC-03 client. */
export function useMercadoMarketClient(injected?: MercadoMarketClient): MercadoMarketClient {
  return useMemo<MercadoMarketClient>(
    () => injected ?? createMarketPriceIntelClient({ baseUrl: apiBaseUrl() }),
    [injected],
  );
}
