import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import type { DashboardLastImport, DashboardOverview } from "./dashboard";

describe("Dashboard SDK contract", () => {
  it("models DashboardOverview with last_import and erp_import degraded", () => {
    const lastImport: DashboardLastImport = {
      at: "2026-07-19T00:00:00Z",
      age_seconds: 3600,
    };
    const overview = {
      sync_errors: 2,
      pending_links: null,
      below_margin: 1,
      missing_gtin: 0,
      orders_today: 0,
      orders_7d: 12,
      anuncios_ativos: 42,
      last_sync_at: { listings: "2026-07-19T00:00:00Z" },
      degraded: ["erp_import"],
      last_import: lastImport,
      as_of: "2026-07-19T01:00:00Z",
    } satisfies DashboardOverview;

    expect(overview.last_import).toEqual(lastImport);
    expect(overview.anuncios_ativos).toBe(42);
    expect(overview.degraded).toContain("erp_import");
  });

  it("keeps the OpenAPI DashboardSummary schema in parity with the SDK (additive, old fields preserved)", () => {
    const openapi = readFileSync(
      resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"),
      "utf8",
    );
    const schema = openapi.slice(
      openapi.indexOf("    DashboardSummary:"),
      openapi.indexOf("    SyncRun:"),
    );

    expect(schema).toContain(
      "required: [sync_errors, pending_links, below_margin, missing_gtin, orders_today, orders_7d, anuncios_ativos, last_sync_at, degraded, last_import, as_of]",
    );
    expect(schema).toContain("anuncios_ativos:");
    expect(schema).toMatch(/last_import:\s*\n\s*type: object\s*\n\s*nullable: true/);
    expect(schema).toContain("enum: [listings, linkage, orders, sync, erp_import]");
  });
});
