import { defineConfig } from "vitest/config";

export default defineConfig({
  esbuild: {
    jsx: "automatic",
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["../../node_modules/@testing-library/jest-dom/dist/vitest.mjs"],
    include: [
      "src/**/*.test.ts",
      "src/**/*.test.tsx",
      "../../packages/feature-products/src/CatalogPage.test.tsx",
      "../../packages/web-query/src/**/*.test.{ts,tsx}",
      "../../packages/ui/src/**/*.test.{ts,tsx}",
    ],
  },
});
