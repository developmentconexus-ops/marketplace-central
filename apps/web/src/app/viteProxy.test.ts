// @vitest-environment node

import { describe, expect, it } from "vitest";
import config from "../../vite.config";

describe("vite dev server proxy", () => {
  it("routes inventory and product link APIs to the backend without shadowing SPA routes", () => {
    expect(config.server?.proxy).toMatchObject({
      "/inventory/stock-actions": "http://backend:8080",
      "/inventory/stock-risks": "http://backend:8080",
      "/product-links/link-candidates": "http://backend:8080",
      "/product-links/link-resolutions": "http://backend:8080",
      "/product-links/link-workflows": "http://backend:8080",
      "/product-links/listing-snapshots": "http://backend:8080",
    });
    expect(config.server?.proxy).not.toHaveProperty("/inventory");
    expect(config.server?.proxy).not.toHaveProperty("/orders");
    expect(config.server?.proxy).not.toHaveProperty("/profitability");
    expect(config.server?.proxy).not.toHaveProperty("/product-links");
  });
});
