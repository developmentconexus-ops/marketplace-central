# S5 — classifications, inventory, product_links, pricing

Blocks 1, 2, 3, 4 and 5 are in `_common-blocks.md` in this same directory. **Read that
file first, in full.** It binds you exactly as if it were inline here. This file is your
slice card and adds only what is specific to your slice.

## Your helpers

1. `internal/modules/classifications/transport/http_handler.go:33` — `writeClassificationsError(w, status int, code, message string)`.
2. `internal/modules/inventory/transport/http_handler.go:45` — `writeInventoryError(w, status int, code, message string)`.
3. `internal/modules/product_links/transport/http_handler.go:112` — `writeProductLinksError(w, status int, code, message string)`.
4. `internal/modules/pricing/transport/http_handler.go:36` — `writePricingError(w, status int, code, message string)`.
5. `internal/modules/pricing/transport/calc_handler.go:62` — `func (h Handler) writeCalcError(w http.ResponseWriter, err error)`.
   This one MAPS an error to a status/code, so it survives; only what it writes changes.
   Read its whole body before touching it: if any branch carries an ad hoc top-level
   field (the chip is migrating `allowed_range`, `limit`, `column`, `import_id`,
   `protocol`), that field moves into `details` under the SAME key name.

Numbers 1-4 are the simple form: delete each and repoint its callers at `apierror.Write`
directly, unless the body does something other than write (read it — do not assume). Any
local envelope struct left orphaned gets deleted, with a grep proving 0 references.

`pricing` has two transport files sharing one package. Migrate both in this slice; they
must end up consistent with each other.

## write_set

- `apps/server_core/internal/modules/classifications/transport/http_handler.go`
- `apps/server_core/internal/modules/classifications/transport/http_handler_test.go`
- `apps/server_core/internal/modules/inventory/transport/http_handler.go`
- `apps/server_core/internal/modules/inventory/transport/http_handler_test.go`
- `apps/server_core/internal/modules/product_links/transport/http_handler.go`
- `apps/server_core/internal/modules/product_links/transport/http_handler_test.go`
- `apps/server_core/internal/modules/pricing/transport/http_handler.go`
- `apps/server_core/internal/modules/pricing/transport/calc_handler.go`
- `apps/server_core/internal/modules/pricing/transport/http_handler_test.go`
- `apps/server_core/internal/modules/pricing/transport/calc_handler_test.go`

If a test file above does not exist, create it. If an existing test elsewhere breaks
because it asserted the flat/old shape, that is expected — but it is OUTSIDE your
write_set, so STOP and report the file:line rather than editing it. Tests in
`apps/server_core/tests/unit/` and `tests/integration/` are explicitly not yours.

VC-2 needs FOUR whole-body pins here (classifications, inventory, product_links, pricing),
one per module, per Block 4.

## Your packages for the test command

`./internal/modules/classifications/... ./internal/modules/inventory/... ./internal/modules/product_links/... ./internal/modules/pricing/...`

**open_questions**: none. If you find one, stop and report rather than deciding.
