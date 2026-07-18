export type ErpImportStatus = "COMPLETED" | "REJECTED";

export type ErpImportIssueCode =
  | "EMPTY_CODPROD"
  | "DUPLICATE_CODPROD"
  | "EMPTY_DESCRPROD"
  | "INVALID_CUSTO"
  | "INVALID_ESTOQUE"
  | "INVALID_EAN"
  | "INVALID_NCM";

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
  source: "xlsx";
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
