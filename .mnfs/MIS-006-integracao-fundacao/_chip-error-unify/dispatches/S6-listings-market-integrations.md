# S6 — listings, market, integrations

Blocks 1, 2, 3, 4 and 5 are in `_common-blocks.md` in this same directory. **Read that
file first, in full.** It binds you exactly as if it were inline here. This file is your
slice card and adds only what is specific to your slice.

## Your helpers

**listings** — `internal/modules/listings/transport/http_handler.go`
- `:333 writeListError(w, status, code, message, key string)` — BUILDS details via
  `if key != "" { details["key"] = key }`, then delegates. Keep the function and keep that
  conditional exactly (an empty key must not become `"key": ""`).
- `:340 writeListErrorDetails(w, status, code, message string, details map[string]any)` —
  its whole body is one write, so it becomes one `apierror.Write` call. Its signature is
  already identical to `apierror.Write`'s tail, so prefer deleting it and repointing
  callers directly; keep it only if the call-site count makes deletion churn worse, and
  say which you chose and why.
- The orphaned `listErrorEnvelope` / `listError` structs get deleted, with a grep proving
  0 references.

**market** — three files, one package plus a sibling
- `internal/modules/market/transport/http_handler.go:128 writeMarketQueryError(w, err error)` — a mapper, survives.
- `internal/modules/market/transport/http_handler.go:141 writeMarketServiceError(w, err error)` — a mapper, survives.
- `internal/modules/market/transport/http_handler.go:150 writeMarketError(w, status, code, message string, details map[string]any)` — the terminal writer; its signature already matches `apierror.Write`, so it should die and its callers repoint.
- `internal/modules/market/transport/collection_handler.go:49 writeCollectionError(w, err error)` — a mapper, survives; read every branch, since this is the module most likely to carry an ad hoc top-level field. Any such field moves into `details` under the same key name.

**integrations** — two files, one package
- `internal/modules/integrations/transport/http_handler.go:48 writeIntegrationError(w, status, code, message string)` — simple form.
- `internal/modules/integrations/transport/run_read_handler.go:107 writeRunQueryError(w, err error)` — a mapper, survives.
- `internal/modules/integrations/transport/run_read_handler.go:123 writeRunReadError(w, status, code, message, key string)` — BUILDS details with the same `key` conditional; keep the conditional exactly. The orphaned `apiErrorResponse` / `apiError` structs in that file get deleted, with a grep proving 0 references — **check first whether the sibling `http_handler.go` in the same package also uses them**, because they are package-scoped and deleting a still-referenced type breaks the build.

## write_set

- `apps/server_core/internal/modules/listings/transport/http_handler.go`
- `apps/server_core/internal/modules/listings/transport/http_handler_test.go`
- `apps/server_core/internal/modules/market/transport/http_handler.go`
- `apps/server_core/internal/modules/market/transport/collection_handler.go`
- `apps/server_core/internal/modules/market/transport/http_handler_test.go`
- `apps/server_core/internal/modules/market/transport/collection_handler_test.go`
- `apps/server_core/internal/modules/integrations/transport/http_handler.go`
- `apps/server_core/internal/modules/integrations/transport/run_read_handler.go`
- `apps/server_core/internal/modules/integrations/transport/http_handler_test.go`
- `apps/server_core/internal/modules/integrations/transport/run_read_handler_test.go`

If a test file above does not exist, create it. If an existing test elsewhere breaks
because it asserted the flat/old shape, that is expected — but it is OUTSIDE your
write_set, so STOP and report the file:line rather than editing it. Tests in
`apps/server_core/tests/unit/` and `tests/integration/` are explicitly not yours.

VC-2 needs THREE whole-body pins here (listings, market, integrations), one per module,
per Block 4. For listings and integrations, pin a case that carries the `key` field, so
the test proves `key` landed inside `details` and nowhere else.

## Your packages for the test command

`./internal/modules/listings/... ./internal/modules/market/... ./internal/modules/integrations/...`

**open_questions**: none. If you find one, stop and report rather than deciding.
