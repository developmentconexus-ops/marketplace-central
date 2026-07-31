import { useMutation, useQueryClient } from "@tanstack/react-query";
import { hasCode, isApiError, type ErpImportCreated, type ErpImportSourceInput } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { erpImportsQueryKeys } from "../vinculos/useErpImports";

export interface ErpImportUploadInput {
  file: File;
  source: ErpImportSourceInput;
}

export interface ErpImportUploadError {
  /** i18n-free reason token used to pick the PT message at the call site. */
  kind:
    | "invalid_file"
    | "missing_required_column"
    | "duplicate_file"
    | "import_in_progress"
    | "deadline_exceeded"
    | "internal_error"
    | "network_error";
  /** The missing column name, when the backend reports one (422). */
  column?: string;
  /** The pre-existing protocol, when the file duplicates a prior import (409). */
  protocol?: string;
}

function classifyUploadError(err: unknown): ErpImportUploadError {
  // A fetch-level failure (no response reached the client at all) never
  // becomes a MarketplaceCentralClientError — that is the honest signal for
  // "no network", distinct from every status-bearing server answer below.
  if (!isApiError(err)) return { kind: "network_error" };
  if (hasCode(err, "missing_required_column")) {
    return { kind: "missing_required_column", column: err.details.column };
  }
  if (hasCode(err, "duplicate_file")) {
    return { kind: "duplicate_file", protocol: err.details.protocol };
  }
  if (hasCode(err, "invalid_file")) return { kind: "invalid_file" };
  if (hasCode(err, "import_in_progress")) {
    // The server's import_in_progress branch passes nil details (no
    // guaranteed protocol) — leave it unknown rather than fabricate one.
    return { kind: "import_in_progress" };
  }
  // The route deadline aborts the request at 120s and answers 504
  // deadline_exceeded. That is not an internal error: the file was too large
  // for the synchronous path, and the operator needs to be told that rather
  // than "erro interno".
  if (hasCode(err, "deadline_exceeded")) return { kind: "deadline_exceeded" };
  return { kind: "internal_error" };
}

/**
 * useErpImportUpload posts the selected workbook to POST /erp/imports and, on
 * success, invalidates the shared ERP-import history list so the new protocol
 * appears immediately. Failures are normalized to ErpImportUploadError so the
 * page can render the right PT message per class (409 duplicate, 422 missing
 * column, etc.) instead of a single opaque error.
 */
export function useErpImportUpload() {
  const client = useClient();
  const queryClient = useQueryClient();

  return useMutation<ErpImportCreated, ErpImportUploadError, ErpImportUploadInput>({
    mutationFn: async ({ file, source }) => {
      try {
        return await client.createErpImport(file, source, file.name);
      } catch (err) {
        throw classifyUploadError(err);
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: erpImportsQueryKeys.list() });
    },
  });
}
