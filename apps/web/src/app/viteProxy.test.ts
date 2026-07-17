// @vitest-environment node

import { describe, expect, it } from "vitest";
import config from "../../vite.config";

describe("vite dev server proxy", () => {
  it("routes inventory and product link APIs to the backend without shadowing SPA routes", () => {
    const proxy = config.server?.proxy;

    expect(proxy).toMatchObject({
      "/inventory/stock-actions": "http://backend:8080",
      "/inventory/stock-risks": "http://backend:8080",
      "/product-links/link-candidates": "http://backend:8080",
      "/product-links/link-resolutions": "http://backend:8080",
      "/product-links/link-workflows": "http://backend:8080",
      "/product-links/listing-snapshots": "http://backend:8080",
      "/listings": "http://backend:8080",
      "/mutations": "http://backend:8080",
      "/market": "http://backend:8080",
      "/profitability": "http://backend:8080",
      "/dashboard": "http://backend:8080",
      "/sync": "http://backend:8080",
    });
    expect(proxy).not.toHaveProperty("/inventory");
    expect(proxy).not.toHaveProperty("/product-links");

    const ordersProxy = proxy?.["/orders"];
    expect(ordersProxy).toMatchObject({ target: "http://backend:8080" });
    expect(ordersProxy).toHaveProperty("bypass", expect.any(Function));

    const bypass = (ordersProxy as {
      bypass: (req: { url?: string; headers: { accept?: string } }) => string | undefined;
    }).bypass;
    expect(bypass({
      url: "/orders?installation=x",
      headers: { accept: "text/html,application/xhtml+xml" },
    })).toBe("/orders?installation=x");
    expect(bypass({
      url: "/orders",
      headers: { accept: "application/json" },
    })).toBeUndefined();

    expect(Object.keys(proxy ?? {})).not.toEqual(expect.arrayContaining([
      "/anuncios",
      "/catalogo",
      "/estoque",
      "/precos",
      "/pedidos",
      "/integracoes",
      "/vinculos",
    ]));
  });
});
