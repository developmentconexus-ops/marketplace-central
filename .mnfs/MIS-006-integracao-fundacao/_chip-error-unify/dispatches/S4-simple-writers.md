# S4 — tenant_config, profitability, dashboard, marketplaces (+ the erp_import whole-body pin)

Blocks 1, 2, 3, 4 and 5 are in `_common-blocks.md` in this same directory. **Read that
file first, in full.** It binds you exactly as if it were inline here. This file is your
slice card and adds only what is specific to your slice.

## Your helpers

All four are the simple `(w, status, code, message)` form. Delete each one and repoint its
callers at `apierror.Write` directly, EXCEPT dashboard's, which builds `details` and
therefore survives.

1. `internal/modules/tenant_config/transport/http_handler.go:165` — `func writeError(w, status int, code, message string)`.
   Note this module's idiom writes the `detail` key; per the shared rule `detail` becomes
   `message`. Check every caller: some pass a detail string, some do not.
2. `internal/modules/profitability/transport/http_handler.go:176` — `func writeError(w, status int, code, message string)`.
3. `internal/modules/dashboard/transport/http_handler.go:64` — `func writeDashboardError(w, status int, code, message, key string)`.
   This one BUILDS details (`if key != "" { details["key"] = key }`). Keep the function
   and keep that conditional exactly; only its `httpx.WriteJSON(...)` line becomes
   `apierror.Write(w, status, code, message, details)`. The local `dashboardAPIError`
   struct (around `:57-62`) is then orphaned — delete it and prove 0 references.
4. `internal/modules/marketplaces/transport/http_handler.go:42` — `func writeMarketplacesError(w, status int, code, message string)`.

## The erp_import debt this slice also clears

`internal/modules/erp_import/transport/http_handler.go` was migrated in an earlier slice
and is already correct — **do not change it**. What it lacks is the VC-2 whole-body pin:
its tests read `.error.code` and `.error.details["column"]` field by field, which cannot
prove the absence of a stray top-level key.

Add ONE test to `internal/modules/erp_import/transport/http_handler_test.go` named
`TestErpImportErrorEnvelopeWholeBody`, pinning the COMPLETE body of the
`missing_required_column` case (chosen because it carries the migrated `column` field in
`details`). Reuse whatever handler-construction and request helpers that file already has
— read it before writing. Do not modify any existing test in that file.

## write_set

- `apps/server_core/internal/modules/tenant_config/transport/http_handler.go`
- `apps/server_core/internal/modules/tenant_config/transport/http_handler_test.go`
- `apps/server_core/internal/modules/profitability/transport/http_handler.go`
- `apps/server_core/internal/modules/profitability/transport/http_handler_test.go`
- `apps/server_core/internal/modules/dashboard/transport/http_handler.go`
- `apps/server_core/internal/modules/dashboard/transport/http_handler_test.go`
- `apps/server_core/internal/modules/marketplaces/transport/http_handler.go`
- `apps/server_core/internal/modules/marketplaces/transport/http_handler_test.go`
- `apps/server_core/internal/modules/erp_import/transport/http_handler_test.go` (the new test ONLY)

If a test file above does not exist, create it. If an existing test elsewhere breaks
because it asserted the flat/old shape, that is expected — but it is OUTSIDE your
write_set, so STOP and report the file:line rather than editing it. Tests in
`apps/server_core/tests/unit/` and `tests/integration/` are explicitly not yours.

## Your packages for the test command

`./internal/modules/tenant_config/... ./internal/modules/profitability/... ./internal/modules/dashboard/... ./internal/modules/marketplaces/... ./internal/modules/erp_import/...`

**open_questions**: none. If you find one, stop and report rather than deciding.
