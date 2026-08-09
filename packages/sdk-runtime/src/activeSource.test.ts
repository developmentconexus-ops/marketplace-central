import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { createMarketplaceCentralClient } from "./index";
import type {
  ActiveSourceConfig,
  ActiveSourceError,
  ActiveSourceName,
  CatalogAssortmentCounts,
  SellableAssortmentConfig,
  SetActiveSourceRequest,
  SetSellableAssortmentRequest,
  SourceKindName,
} from "./activeSource";

describe("active-source SDK contract", () => {
  it("constructs every exported active-source type", () => {
    const activeSource: ActiveSourceName = "xlsx";
    const sourceKind: SourceKindName = "upload_snapshot";
    const config: ActiveSourceConfig = {
      active_source: activeSource,
      source_kind: sourceKind,
      set_at: "2026-07-20T12:00:00Z",
      set_by: "operator-1",
    };
    const request: SetActiveSourceRequest = {
      active_source: "sankhya",
      set_by: null,
    };
    const error: ActiveSourceError = { error: "unknown_erp_source" };

    expect(config.active_source).toBe("xlsx");
    expect(config.source_kind).toBe("upload_snapshot");
    expect(request.active_source).toBe("sankhya");
    expect(error.error).toBe("unknown_erp_source");
  });

  it("keeps SetActiveSourceRequest server-derive-only: source_kind is not an accepted field", () => {
    // @ts-expect-error source_kind is server-derived (tenant_config.DefaultKind) — never client-supplied.
    const withSourceKind: SetActiveSourceRequest = {
      active_source: "xlsx",
      source_kind: "upload_snapshot",
    };
    expect(withSourceKind.active_source).toBe("xlsx");
  });

  it("keeps the OpenAPI active-source schemas in parity with the SDK", () => {
    const openapi = readFileSync(
      resolve(
        fileURLToPath(new URL(".", import.meta.url)),
        "../../../contracts/api/marketplace-central.openapi.yaml",
      ),
      "utf8",
    );
    const region = openapi.slice(openapi.indexOf("    ActiveSourceConfig:"));
    const config = region.slice(
      region.indexOf("    ActiveSourceConfig:"),
      region.indexOf("    SetActiveSourceRequest:"),
    );
    const request = region.slice(
      region.indexOf("    SetActiveSourceRequest:"),
      region.indexOf("    ActiveSourceError:"),
    );
    const error = region.slice(region.indexOf("    ActiveSourceError:"));

    expect(config).toContain("required: [active_source, source_kind, set_at]");
    expect(config).toContain("enum: [sankhya, xlsx, catalogo_cliente]");
    expect(config).toContain("enum: [upload_snapshot, live_read_through]");
    expect(request).toContain("required: [active_source]");
    expect(request).not.toContain("source_kind:");
    expect(error).toContain(
      "enum: [unknown_erp_source, invalid_active_source, invalid_body, internal_error]",
    );
  });

  it("declares /config/active-source GET and PUT with flat ActiveSourceError responses", () => {
    const openapi = readFileSync(
      resolve(
        fileURLToPath(new URL(".", import.meta.url)),
        "../../../contracts/api/marketplace-central.openapi.yaml",
      ),
      "utf8",
    );
    // Window RE-POINTED BY VALUE (hub ruling A-19), assertions unchanged — still exactly
    // two 400s, still ActiveSourceError in each. It used to run to `\ncomponents:`, so the
    // sibling /config/sellable-assortment path appended after it was counted here.
    const region = openapi.slice(
      openapi.indexOf("  /config/active-source:"),
      openapi.indexOf("\ncomponents:"),
    );
    const nextPath = region.search(/\n {2}\/(?!config\/active-source)\S*:/);
    const path = nextPath === -1 ? region : region.slice(0, nextPath);
    expect(path).toContain("operationId: getActiveSource");
    expect(path).toContain("operationId: setActiveSource");
    expect(path).toContain("$ref: '#/components/schemas/ActiveSourceConfig'");
    const fourHundreds = path.match(/"400":/g) ?? [];
    expect(fourHundreds).toHaveLength(2);
    for (const block of path.split(/"400":/).slice(1)) {
      // Local window bounding to this 400's own response body; 300 covers the
      // longest description (GET's fail-closed note) while still failing if the
      // $ref were absent/wrong — the ref sits at ~char 209 after that description.
      expect(block.slice(0, 300)).toContain("ActiveSourceError");
    }
  });

  it("exposes the sellable-assortment paths, schemas, options, and exact client URLs", async () => {
    const openapi = readFileSync(
      resolve(
        fileURLToPath(new URL(".", import.meta.url)),
        "../../../contracts/api/marketplace-central.openapi.yaml",
      ),
      "utf8",
    );
    const config: SellableAssortmentConfig = {
      only_revenda: true,
      only_em_estoque: true,
      only_ecommerce_eligible: false,
    };
    const request: SetSellableAssortmentRequest = {
      only_revenda: false,
      only_em_estoque: true,
      only_ecommerce_eligible: true,
    };
    const counts: CatalogAssortmentCounts = { sellable_count: 2, total_count: 4 };
    expect(
      config.only_revenda,
      "SDK symbol SellableAssortmentConfig.only_revenda was not exposed",
    ).toBe(true);
    expect(
      request.only_ecommerce_eligible,
      "SDK symbol SetSellableAssortmentRequest.only_ecommerce_eligible was not exposed",
    ).toBe(true);
    expect(
      counts.total_count,
      "SDK symbol CatalogAssortmentCounts.total_count was not exposed",
    ).toBe(4);

    const apiPaths = openapi.slice(0, openapi.indexOf("\ncomponents:"));
    // Bounded by VALUE at birth (hub ruling A-19): this guard is the one written while
    // repairing the two windows that ran to the end of the paths section, so it does not
    // repeat their shape. /config/sellable-assortment is last today; the moment it is not,
    // this slice still describes only its own path.
    const configRegion = apiPaths.slice(apiPaths.indexOf("  /config/sellable-assortment:"));
    const afterConfig = configRegion.search(/\n {2}\/(?!config\/sellable-assortment)\S*:/);
    const configPath = afterConfig === -1 ? configRegion : configRegion.slice(0, afterConfig);
    const countsPath = apiPaths.slice(apiPaths.indexOf("  /catalog/products/counts:"));
    const listPath = apiPaths.slice(
      apiPaths.indexOf("  /catalog/products:"),
      apiPaths.indexOf("  /catalog/products/{id}:"),
    );
    const searchPath = apiPaths.slice(
      apiPaths.indexOf("  /catalog/products/search:"),
      apiPaths.indexOf("  /catalog/products/{id}/enrichment:"),
    );
    const schemas = openapi.slice(openapi.indexOf("    SellableAssortmentConfig:"));
    const configSchema = schemas.slice(
      schemas.indexOf("    SellableAssortmentConfig:"),
      schemas.indexOf("    SetSellableAssortmentRequest:"),
    );
    const requestSchema = schemas.slice(
      schemas.indexOf("    SetSellableAssortmentRequest:"),
      schemas.indexOf("    SellableAssortmentError:"),
    );
    const errorSchema = schemas.slice(
      schemas.indexOf("    SellableAssortmentError:"),
      schemas.indexOf("    CatalogAssortmentCounts:"),
    );
    const countsSchema = schemas.slice(schemas.indexOf("    CatalogAssortmentCounts:"));
    expect(apiPaths, "OpenAPI missing path /config/sellable-assortment").toContain(
      "  /config/sellable-assortment:",
    );
    expect(configPath, "OpenAPI missing GET operationId getSellableAssortment").toContain(
      "operationId: getSellableAssortment",
    );
    expect(configPath, "OpenAPI missing PUT operationId setSellableAssortment").toContain(
      "operationId: setSellableAssortment",
    );
    expect(configPath, "OpenAPI missing SellableAssortmentConfig response schema").toContain(
      "$ref: '#/components/schemas/SellableAssortmentConfig'",
    );
    expect(configPath, "OpenAPI missing SetSellableAssortmentRequest request schema").toContain(
      "$ref: '#/components/schemas/SetSellableAssortmentRequest'",
    );
    // The error surface must be what the server can actually write. The handler that will
    // serve these two operations is tenant_config/transport/http_handler.go, whose
    // writeError emits map[string]string{"error": code} plus an optional "detail" — the
    // FLAT shape, never the nested ErrorResponse. It has no not-found branch at all, and
    // under A-22 an absent tenant row fails closed with unknown_erp_source, so GET can
    // never 404. The contract shipped 404 + nested ErrorResponse on both operations; those
    // were three statements the server cannot honour, and a false statement in a published
    // contract is deleted, not softened.
    expect(
      configPath,
      "OpenAPI /config/sellable-assortment declares a 404 the handler has no branch for",
    ).not.toContain('"404":');
    expect(
      configPath,
      "OpenAPI /config/sellable-assortment refs the nested ErrorResponse; this surface writes a flat body",
    ).not.toContain("ErrorResponse");
    expect(
      configPath.match(/"500":/g) ?? [],
      "OpenAPI /config/sellable-assortment must declare 500 on BOTH operations",
    ).toHaveLength(2);
    // A-22 moves the fail-closed unknown_erp_source response onto GET as well as PUT.
    expect(
      configPath.match(/"400":/g) ?? [],
      "OpenAPI /config/sellable-assortment must declare 400 on GET and PUT",
    ).toHaveLength(2);
    for (const block of configPath.split(/"(?:400|500)":/).slice(1)) {
      expect(
        block.slice(0, 300),
        "every /config/sellable-assortment error resolves to SellableAssortmentError",
      ).toContain("SellableAssortmentError");
    }
    expect(
      errorSchema,
      "OpenAPI SellableAssortmentError must be the flat {error, detail?} body",
    ).toContain("required: [error]");
    // A-22 moves unknown_erp_source into SellableAssortmentError for both verbs.
    expect(errorSchema, "OpenAPI SellableAssortmentError has the wrong code enum").toContain(
      "enum: [unknown_erp_source, invalid_body, internal_error]",
    );
    expect(apiPaths, "OpenAPI missing path /catalog/products/counts").toContain(
      "  /catalog/products/counts:",
    );
    expect(countsPath, "OpenAPI missing operationId getCatalogAssortmentCounts").toContain(
      "operationId: getCatalogAssortmentCounts",
    );
    expect(countsPath, "OpenAPI missing CatalogAssortmentCounts response schema").toContain(
      "$ref: '#/components/schemas/CatalogAssortmentCounts'",
    );
    expect(
      configSchema,
      "OpenAPI SellableAssortmentConfig has the wrong required fields",
    ).toContain("required: [only_revenda, only_em_estoque, only_ecommerce_eligible]");
    expect(
      requestSchema,
      "OpenAPI SetSellableAssortmentRequest has the wrong required fields",
    ).toContain("required: [only_revenda, only_em_estoque, only_ecommerce_eligible]");
    expect(countsSchema, "OpenAPI CatalogAssortmentCounts has the wrong required fields").toContain(
      "required: [sellable_count, total_count]",
    );
    expect(
      countsSchema,
      "OpenAPI CatalogAssortmentCounts is missing non-negative count constraints",
    ).toMatch(
      /sellable_count:[\s\S]*?type: integer[\s\S]*?minimum: 0[\s\S]*?total_count:[\s\S]*?type: integer[\s\S]*?minimum: 0/,
    );
    expect(listPath, "OpenAPI missing include_all on GET /catalog/products").toMatch(
      /name: include_all[\s\S]*?in: query[\s\S]*?type: boolean[\s\S]*?default: false/,
    );
    expect(searchPath, "OpenAPI missing include_all on GET /catalog/products/search").toMatch(
      /name: include_all[\s\S]*?in: query[\s\S]*?type: boolean[\s\S]*?default: false/,
    );

    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        const body = String(input).endsWith("/catalog/products/counts") ? counts : config;
        return new Response(JSON.stringify(body), { status: 200 });
      },
    });

    expect(
      typeof client.getSellableAssortment,
      "SDK symbol getSellableAssortment was not exposed",
    ).toBe("function");
    expect(
      typeof client.setSellableAssortment,
      "SDK symbol setSellableAssortment was not exposed",
    ).toBe("function");
    expect(
      typeof client.getCatalogAssortmentCounts,
      "SDK symbol getCatalogAssortmentCounts was not exposed",
    ).toBe("function");
    await client.getSellableAssortment();
    await client.setSellableAssortment(request);
    const receivedCounts = await client.getCatalogAssortmentCounts();
    await client.listCatalogProductFacts({ cursor: "MTIz", limit: 25, include_all: true });
    await client.searchCatalogProductFacts({
      q: "PARAFUSO",
      cursor: "NDU2",
      limit: 50,
      include_all: true,
    });

    expect(String(requests[0].input), "SDK getSellableAssortment read the wrong URL").toBe(
      "http://localhost:8080/config/sellable-assortment",
    );
    expect(String(requests[1].input), "SDK setSellableAssortment wrote the wrong URL").toBe(
      "http://localhost:8080/config/sellable-assortment",
    );
    expect(String(requests[2].input), "SDK getCatalogAssortmentCounts read the wrong URL").toBe(
      "http://localhost:8080/catalog/products/counts",
    );
    expect(
      String(requests[3].input),
      "SDK listCatalogProductFacts serialized the wrong include_all URL",
    ).toBe("http://localhost:8080/catalog/products?cursor=MTIz&limit=25&include_all=true");
    expect(
      String(requests[4].input),
      "SDK searchCatalogProductFacts serialized the wrong include_all URL",
    ).toBe(
      "http://localhost:8080/catalog/products/search?q=PARAFUSO&cursor=NDU2&limit=50&include_all=true",
    );
    expect(requests[1].init?.method, "SDK setSellableAssortment used the wrong HTTP method").toBe(
      "PUT",
    );
    expect(
      JSON.parse(String(requests[1].init?.body)),
      "SDK setSellableAssortment sent the wrong request body",
    ).toEqual(request);
    expect(
      receivedCounts.sellable_count,
      "SDK CatalogAssortmentCounts read the wrong sellable_count field",
    ).toBe(2);
  });
});
