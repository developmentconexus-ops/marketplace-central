import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import type { ActiveSourceConfig, ActiveSourceError, ActiveSourceName, SetActiveSourceRequest, SourceKindName } from "./activeSource";

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
    const openapi = readFileSync(resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"), "utf8");
    const region = openapi.slice(openapi.indexOf("    ActiveSourceConfig:"));
    const config = region.slice(region.indexOf("    ActiveSourceConfig:"), region.indexOf("    SetActiveSourceRequest:"));
    const request = region.slice(region.indexOf("    SetActiveSourceRequest:"), region.indexOf("    ActiveSourceError:"));
    const error = region.slice(region.indexOf("    ActiveSourceError:"));

    expect(config).toContain("required: [active_source, source_kind, set_at]");
    expect(config).toContain("enum: [sankhya, xlsx, catalogo_cliente]");
    expect(config).toContain("enum: [upload_snapshot, live_read_through]");
    expect(request).toContain("required: [active_source]");
    expect(request).not.toContain("source_kind:");
    expect(error).toContain("enum: [unknown_erp_source, invalid_active_source, invalid_body, internal_error]");
  });

  it("declares /config/active-source GET and PUT with flat ActiveSourceError responses", () => {
    const openapi = readFileSync(resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"), "utf8");
    const path = openapi.slice(openapi.indexOf("  /config/active-source:"), openapi.indexOf("\ncomponents:"));
    expect(path).toContain("operationId: getActiveSource");
    expect(path).toContain("operationId: setActiveSource");
    expect(path).toContain("$ref: '#/components/schemas/ActiveSourceConfig'");
    const fourHundreds = path.match(/"400":/g) ?? [];
    expect(fourHundreds).toHaveLength(2);
    for (const block of path.split(/"400":/).slice(1)) {
      expect(block.slice(0, 200)).toContain("ActiveSourceError");
    }
  });
});
