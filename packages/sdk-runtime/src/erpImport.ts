export type ErpImportStatus = "COMPLETED" | "REJECTED";

export type ErpImportIssueCode =
  | "EMPTY_CODPROD"
  | "DUPLICATE_CODPROD"
  | "EMPTY_DESCRPROD"
  | "INVALID_CUSTO"
  | "INVALID_ESTOQUE"
  | "INVALID_EAN"
  | "INVALID_NCM"
  | "MISSING_CUSTO"
  | "MISSING_ESTOQUE"
  | "MISSING_REQUIRED_COLUMN";

/**
 * Import path selector sent as the multipart `source` field.
 * - "catalogo_cliente" → lenient client-catalog path (CUSTO/ESTOQUE_FISICO may
 *   be absent and stay honest-unknown, ADR-17).
 * - "xlsx" → strict Sankhya cost+stock path (rejects missing required columns).
 */
export type ErpImportSourceInput = "xlsx" | "catalogo_cliente";

export interface ErpImportIssue {
  row: number;
  code: ErpImportIssueCode;
  detail: string;
  column?: string | null;
  offending_value?: string | null;
}

export interface ErpImportSummary {
  import_id: string;
  protocol: string;
  file_sha256: string;
  source: ErpImportSourceInput;
  imported_at: string;
  status: ErpImportStatus;
  accepted_count: number;
  rejected_count: number;
  warning_count: number;
}

export interface ErpImportDetail extends ErpImportSummary {
  rejected_rows: ErpImportIssue[];
  warnings: ErpImportIssue[];
}

export interface ErpImportList {
  items: ErpImportSummary[];
}

export interface ErpImportCreated {
  import_id: string;
  protocol: string;
  status: ErpImportStatus;
}

export interface ErpImportError {
  error: "invalid_file" | "missing_required_column" | "import_not_found" | "internal_error";
  detail?: string;
  column?: string;
}

export interface ErpImportConflict {
  error: "duplicate_file" | "import_in_progress";
  import_id?: string;
  protocol?: string;
}
