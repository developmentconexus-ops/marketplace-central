# S7 — orders, mutations, connectors, and the platform deadline writer

Blocks 1, 2, 3, 4 and 5 are in `_common-blocks.md` in this same directory. **Read that
file first, in full.** It binds you exactly as if it were inline here. This file is your
slice card and adds only what is specific to your slice.

## Your helpers

**orders** — two files, one package
- `internal/modules/orders/transport/http_handler.go:359 writeSummaryError(w, err error)` — mapper, survives.
- `:688 writeOrdersError(w, status, code, message string)` — simple form.
- `:692 writeOrderReadError(w, err error)` — mapper, survives.
- `:709 writeOrdersErrorDetails(w, status, code, message string, details map[string]any)` — terminal writer, signature already matches `apierror.Write`; it should die and its callers repoint.
- `internal/modules/orders/transport/sankhya_linkage_handler.go:174 writeSankhyaLinkageError(w, err error)` — mapper, survives. Read every branch: this handler is linkage-related and is the most likely place in the tree to carry an ad hoc top-level field. Any such field moves into `details` under the same key name.

**mutations** — this module is the family-A REFERENCE the rest of the chip was modelled on
(`transport/errors.go:20-35`), so its wire shape should not change at all. Your job is to
make it produce that shape through `apierror.Write` instead of its own local structs.
- `internal/modules/mutations/transport/errors.go:37 writeMutationError(w, err error)` — mapper, survives.
- `:44 writeMutationMethodError(w)` — simple.
- `internal/modules/mutations/transport/query_handler.go:326 writeMutationQueryMethodError(w)` — simple.
- The local `mutationErrorEnvelope` / `mutationError` structs get deleted once orphaned,
  with a grep proving 0 references.
- **Because this module's shape does not change, its existing tests must keep passing
  UNEDITED.** If one goes red, that is a real regression in your change, not a stale
  test — fix your change, do not touch the test. Report it either way.

**connectors** — two files in different layers
- `internal/modules/connectors/transport/http_handler.go:42 writeConnectorsError(w, status, code, message string)` — simple form. Note (from this chip's context pack) this site currently emits NO `details` key at all, so it is the one that made `.error.details` `undefined` on some routes and an object on others; after migration it must emit `{}`. Pin that in the whole-body test.
- `internal/modules/connectors/adapters/melhorenvio/oauth.go:249 writeOAuthError(w, status, code, message string)` — this one is in an ADAPTER, not a transport. It still answers HTTP (an OAuth callback), so it migrates too. Keep the adapter-boundary rule in mind: you are changing how it writes, not moving provider payloads around.

**platform — the finding this slice must close**
`internal/platform/httpx/route_deadline.go:129`:
```go
func writeDeadlineExceeded(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusGatewayTimeout)
	_, _ = io.WriteString(w, `{"error":"deadline_exceeded"}`)
}
```
This is a FLAT body written by the platform router, so it fires on every deadline-bounded
route regardless of module — the highest-reach flat producer in the tree, and no
`write*Error` grep in a module directory could ever see it.

It must emit the envelope: `504` with
`{"error":{"code":"deadline_exceeded","message":"<pt-BR>","details":{}}}`.

**Hard constraint**: `httpx` must NOT import `apierror`. The dependency direction is
`apierror -> httpx`, and reversing it is an import cycle. Write the body as a **constant
string literal**, exactly the way `internal/platform/httpx/json.go` already does it for
its encode-failure fallback — read that file and match its idiom rather than inventing a
second one. Say in your report that you checked the cycle direction and how.

Add a unit test in the httpx package pinning the complete 504 body, and assert the
`Content-Type` header is still `application/json`.

## write_set

- `apps/server_core/internal/modules/orders/transport/http_handler.go`
- `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler.go`
- `apps/server_core/internal/modules/orders/transport/http_handler_test.go`
- `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler_test.go`
- `apps/server_core/internal/modules/mutations/transport/errors.go`
- `apps/server_core/internal/modules/mutations/transport/query_handler.go`
- `apps/server_core/internal/modules/mutations/transport/errors_test.go`
- `apps/server_core/internal/modules/connectors/transport/http_handler.go`
- `apps/server_core/internal/modules/connectors/transport/http_handler_test.go`
- `apps/server_core/internal/modules/connectors/adapters/melhorenvio/oauth.go`
- `apps/server_core/internal/modules/connectors/adapters/melhorenvio/oauth_test.go`
- `apps/server_core/internal/platform/httpx/route_deadline.go`
- `apps/server_core/internal/platform/httpx/route_deadline_test.go`

If a test file above does not exist, create it. If an existing test elsewhere breaks
because it asserted the flat/old shape, that is expected — but it is OUTSIDE your
write_set, so STOP and report the file:line rather than editing it. Tests in
`apps/server_core/tests/unit/` and `tests/integration/` are explicitly not yours.
`internal/platform/httpx/json.go` is owned by another slice — read it, never edit it.

VC-2 needs whole-body pins for orders, mutations, connectors and the httpx deadline
writer, per Block 4.

## Your packages for the test command

`./internal/modules/orders/... ./internal/modules/mutations/... ./internal/modules/connectors/... ./internal/platform/httpx/...`

**open_questions**: none. If you find one, stop and report rather than deciding.
