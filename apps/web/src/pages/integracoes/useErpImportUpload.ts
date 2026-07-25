import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { ErpImportCreated, ErpImportSourceInput } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { erpImportsQueryKeys } from "../vinculos/useErpImports";

export interface ErpImportUploadInput {
  file: File;
  source: ErpImportSourceInput;
}

/**
 * The flat error shape the ERP import endpoint returns on non-2xx. The SDK's
 * multipart helper attaches the parsed body under `body`; we read it to
 * distinguish the operator-facing failure classes (duplicate/in-progress/
 * missing-column) rather than collapsing every failure into one message.
 */
interface ErpImportErrorBody {
  error?: string;
  detail?: string;
  column?: string;
  import_id?: string;
  protocol?: string;
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
  const maybe = err as { status?: number; body?: ErpImportErrorBody } | undefined;
  const body = maybe?.body;
  const code = body?.error;
  switch (code) {
    case "invalid_file":
      return { kind: "invalid_file" };
    case "missing_required_column":
      return { kind: "missing_required_column", column: body?.column };
    case "duplicate_file":
      return { kind: "duplicate_file", protocol: body?.protocol };
    case "import_in_progress":
      return { kind: "import_in_progress", protocol: body?.protocol };
    case "internal_error":
      return { kind: "internal_error" };
    // The route deadline aborts the request at 120s and answers 504. That is not
    // an internal error: the file was too large for the synchronous path, and the
    // operator needs to be told that rather than "erro interno".
    case "deadline_exceeded":
      return { kind: "deadline_exceeded" };
    default:
      // A thrown Response-shaped error with a known status but unmapped body,
      // or a genuine transport failure with no body at all.
      if (maybe?.status === 422) return { kind: "missing_required_column", column: body?.column };
      if (maybe?.status === 409) return { kind: "duplicate_file", protocol: body?.protocol };
      if (maybe?.status === 400) return { kind: "invalid_file" };
      if (maybe?.status === 504) return { kind: "deadline_exceeded" };
      if (typeof maybe?.status === "number") return { kind: "internal_error" };
      return { kind: "network_error" };
  }
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
