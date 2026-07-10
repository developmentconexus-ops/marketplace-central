import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../..");
const redirectUri = process.env.MPC_OAUTH_REDIRECT_URI;
const proxyTarget = process.env.MPC_WEB_PROXY_TARGET ?? "http://backend:8080";

function additionalAllowedHosts(): string[] {
  if (!redirectUri) {
    return [];
  }

  try {
    return [new URL(redirectUri).host];
  } catch {
    return [];
  }
}

export default defineConfig({
  plugins: [tailwindcss()],
  esbuild: {
    jsx: "automatic",
  },
  server: {
    host: "0.0.0.0",
    port: 5174,
    allowedHosts: additionalAllowedHosts(),
    proxy: {
      "/catalog": proxyTarget,
      "/classifications": proxyTarget,
      "/connectors": proxyTarget,
      "/healthz": proxyTarget,
      "/integrations/auth": proxyTarget,
      "/integrations/installations": proxyTarget,
      "/integrations/providers": proxyTarget,
      "/inventory/stock-actions": proxyTarget,
      "/inventory/stock-risks": proxyTarget,
      "/marketplaces/accounts": proxyTarget,
      "/marketplaces/definitions": proxyTarget,
      "/marketplaces/fee-schedules": proxyTarget,
      "/marketplaces/policies": proxyTarget,
      "/pricing": proxyTarget,
      "/product-links/link-candidates": proxyTarget,
      "/product-links/link-resolutions": proxyTarget,
      "/product-links/link-workflows": proxyTarget,
      "/product-links/listing-snapshots": proxyTarget,
    },
  },
  preview: {
    port: 3002,
  },
  test: {
    environment: "jsdom",
    globals: true,
    dir: rootDir,
    setupFiles: [path.resolve(rootDir, "node_modules/@testing-library/jest-dom/dist/vitest.mjs")],
    include: ["apps/web/src/**/*.test.ts", "apps/web/src/**/*.test.tsx", "packages/**/*.test.ts", "packages/**/*.test.tsx"],
    exclude: ["**/node_modules/**", "**/.worktrees/**"],
  },
});
