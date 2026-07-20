import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import type {
  ErpImportConflict,
  ErpImportCreated,
  ErpImportDetail,
  ErpImportError,
  ErpImportIssue,
  ErpImportIssueCode,
  ErpImportList,
  ErpImportStatus,
  ErpImportSummary,
} from "./erpImport";

describe("ERP import SDK contract", () => {
  it("constructs every exported ERP import type", () => {
    const status: ErpImportStatus = "COMPLETED";
    const issueCode: ErpImportIssueCode = "INVALID_EAN";
    const issue: ErpImportIssue = {
      row: 2,
      code: issueCode,
      detail: "ean is not a valid GTIN",
      column: "EAN",
      offending_value: "bad-ean",
    };
    const summary: ErpImportSummary = {
      import_id: "11111111-1111-1111-1111-111111111111",
      protocol: "#101-E",
      file_sha256: "hash",
      source: "xlsx",
      imported_at: "2026-07-17T12:00:00Z",
      status,
      accepted_count: 1,
      rejected_count: 0,
      warning_count: 1,
    };
    const detail: ErpImportDetail = {
      ...summary,
      rejected_rows: [],
      warnings: [issue],
    };
    const list: ErpImportList = { items: [summary] };
    const created: ErpImportCreated = {
      import_id: summary.import_id,
      protocol: summary.protocol,
      status,
    };
    const error: ErpImportError = { error: "invalid_file", detail: "not an xlsx file" };
    const conflict: ErpImportConflict = {
      error: "duplicate_file",
      import_id: summary.import_id,
      protocol: summary.protocol,
    };

    expect(detail.warnings[0]).toEqual(issue);
    expect(list.items[0]).toEqual(summary);
    expect(created.status).toBe("COMPLETED");
    expect(error.error).toBe("invalid_file");
    expect(conflict.protocol).toBe(summary.protocol);
  });

  it("keeps the OpenAPI ERP import schemas in parity with the SDK", () => {
    const openapi = readFileSync(resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"), "utf8");
    const region = openapi.slice(openapi.indexOf("    ErpImportStatus:"));
    const issueCodeSchema = region.slice(region.indexOf("    ErpImportIssueCode:"), region.indexOf("    ErpImportIssue:"));
    const summary = region.slice(region.indexOf("    ErpImportSummary:"), region.indexOf("    ErpImportIssueCode:"));
    const issue = region.slice(region.indexOf("    ErpImportIssue:"), region.indexOf("    ErpImportDetail:"));
    const error = region.slice(region.indexOf("    ErpImportError:"), region.indexOf("    ErpImportConflict:"));
    const conflict = region.slice(region.indexOf("    ErpImportConflict:"));

    expect(region).toContain("enum: [COMPLETED, REJECTED]");
    expect(summary).toContain("required: [import_id, protocol, file_sha256, source, imported_at, status, accepted_count, rejected_count, warning_count]");
    expect(summary).toContain("file_sha256:");
    expect(summary).toContain("source:");
    expect(summary).not.toContain("filename:");
    expect(issueCodeSchema).toContain("enum: [EMPTY_CODPROD, DUPLICATE_CODPROD, EMPTY_DESCRPROD, INVALID_CUSTO, INVALID_ESTOQUE, INVALID_EAN, INVALID_NCM, MISSING_CUSTO, MISSING_ESTOQUE, MISSING_REQUIRED_COLUMN]");
    expect(error).toContain("enum: [invalid_file, missing_required_column, import_not_found, internal_error]");
    expect(error).not.toContain("$ref: '#/components/schemas/ErrorResponse'");
    expect(issue).toContain("required: [row, code, detail]");
    expect(issue).toMatch(/column:\s*\n\s*type: string\s*\n\s*nullable: true/);
    expect(issue).toMatch(/offending_value:\s*\n\s*type: string\s*\n\s*nullable: true/);
    expect(conflict).toContain("import_id:");
    expect(conflict).toContain("protocol:");
    expect(conflict).toContain("enum: [duplicate_file, import_in_progress]");
  });

  it("declares a flat 500 ErpImportError on every erp import operation", () => {
    const openapi = readFileSync(resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"), "utf8");
    const paths = openapi.slice(openapi.indexOf("  /erp/imports:"), openapi.indexOf("\ncomponents:"));
    // handler emits {"error":"internal_error"} (500) on POST, list, and detail — spec must cover all three.
    const fiveHundreds = paths.match(/"500":/g) ?? [];
    expect(fiveHundreds).toHaveLength(3);
    // no 500 may resolve to the legacy nested ErrorResponse — erp errors are flat.
    expect(paths).not.toContain("ErrorResponse");
    for (const block of paths.split(/"500":/).slice(1)) {
      expect(block.slice(0, 200)).toContain("ErpImportError");
    }
  });
});
