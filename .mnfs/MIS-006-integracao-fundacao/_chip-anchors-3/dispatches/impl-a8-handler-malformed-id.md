# A8 — malformed {id} answers 4xx on BOTH routes — implementer artifact

File owned and edited: `apps/server_core/internal/modules/erp_import/transport/http_handler_test.go`

No other file was touched. `http_handler.go`, `domain/import.go`, the OpenAPI yaml, and
everything under `packages/` were read-only context.

## 1. The test added (verbatim)

Also required a one-field extension to the existing `fakeImportQuerier` stub in the same
file (no new mock type introduced — `chainID` already existed for the chain route; `getID`
mirrors it for `GetImport` so the id passed into `GetImport` is observable the same way):

```go
type fakeImportQuerier struct {
	items    []domain.ImportReport
	item     domain.ImportReport
	chain    domain.ImportChain
	listErr  error
	getErr   error
	chainErr error
	getID    domain.ImportID
	chainID  domain.ImportID
}

func (f *fakeImportQuerier) ListImports(context.Context) ([]domain.ImportReport, error) {
	return f.items, f.listErr
}

func (f *fakeImportQuerier) GetImport(_ context.Context, id domain.ImportID) (domain.ImportReport, error) {
	f.getID = id
	return f.item, f.getErr
}
```

(`GetImportChain` was already recording into `chainID` — unchanged.)

The new test itself:

```go
func TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery(t *testing.T) {
	const malformedID = "not-a-uuid"
	tests := []struct {
		name string
		path string
	}{
		{name: "get import", path: "/erp/imports/" + malformedID},
		{name: "get import chain", path: "/erp/imports/" + malformedID + "/chain"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			querier := &fakeImportQuerier{}
			response := performRequest(t, &fakeImportRunner{}, querier, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", response.Code)
			}
			if got, want := response.Body.String(), "{\"error\":\"invalid_import_id\"}\n"; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if querier.getID != "" {
				t.Fatalf("GetImport ran with id %q; the malformed path value must be rejected before any query", querier.getID)
			}
			if querier.chainID != "" {
				t.Fatalf("GetImportChain ran with id %q; the malformed path value must be rejected before any query", querier.chainID)
			}
		})
	}
}
```

## 2. Exact commands run and their real output

Command 1 — full package suite, from `apps/server_core`:

```
cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/erp_import/transport/... -run '.*' -v
```

Output (tail, full run, all green):

```
=== RUN   TestHandlerPostImport
=== RUN   TestHandlerPostImport/completed
=== RUN   TestHandlerPostImport/all_rejected
--- PASS: TestHandlerPostImport (0.00s)
    --- PASS: TestHandlerPostImport/completed (0.00s)
    --- PASS: TestHandlerPostImport/all_rejected (0.00s)
=== RUN   TestHandlerPostImportLenientTrigger
=== RUN   TestHandlerPostImportLenientTrigger/default_strict
=== RUN   TestHandlerPostImportLenientTrigger/source_catalogo_cliente
=== RUN   TestHandlerPostImportLenientTrigger/mode_lenient
=== RUN   TestHandlerPostImportLenientTrigger/unrelated_field
--- PASS: TestHandlerPostImportLenientTrigger (0.00s)
    --- PASS: TestHandlerPostImportLenientTrigger/default_strict (0.00s)
    --- PASS: TestHandlerPostImportLenientTrigger/source_catalogo_cliente (0.00s)
    --- PASS: TestHandlerPostImportLenientTrigger/mode_lenient (0.00s)
    --- PASS: TestHandlerPostImportLenientTrigger/unrelated_field (0.00s)
=== RUN   TestHandlerPostImportErrors
=== RUN   TestHandlerPostImportErrors/invalid_file
=== RUN   TestHandlerPostImportErrors/missing_required_column
=== RUN   TestHandlerPostImportErrors/duplicate
=== RUN   TestHandlerPostImportErrors/import_in_progress
=== RUN   TestHandlerPostImportErrors/internal
--- PASS: TestHandlerPostImportErrors (0.00s)
    --- PASS: TestHandlerPostImportErrors/invalid_file (0.00s)
    --- PASS: TestHandlerPostImportErrors/missing_required_column (0.00s)
    --- PASS: TestHandlerPostImportErrors/duplicate (0.00s)
    --- PASS: TestHandlerPostImportErrors/import_in_progress (0.00s)
    --- PASS: TestHandlerPostImportErrors/internal (0.00s)
=== RUN   TestHandlerPostImportMalformedRequests
=== RUN   TestHandlerPostImportMalformedRequests/missing_file_field
=== RUN   TestHandlerPostImportMalformedRequests/non_multipart
--- PASS: TestHandlerPostImportMalformedRequests (0.00s)
    --- PASS: TestHandlerPostImportMalformedRequests/missing_file_field (0.00s)
    --- PASS: TestHandlerPostImportMalformedRequests/non_multipart (0.00s)
=== RUN   TestHandlerGetImportListPreservesServiceOrder
--- PASS: TestHandlerGetImportListPreservesServiceOrder (0.00s)
=== RUN   TestHandlerGetImportDetailSplitsIssuesAndEmitsEmptyArrays
--- PASS: TestHandlerGetImportDetailSplitsIssuesAndEmitsEmptyArrays (0.00s)
=== RUN   TestHandlerGetImportDetailUnknownID
--- PASS: TestHandlerGetImportDetailUnknownID (0.00s)
=== RUN   TestHandlerPostImportOversizeRejected
--- PASS: TestHandlerPostImportOversizeRejected (0.03s)
=== RUN   TestHandlerImportedAtIsRFC3339UTC
--- PASS: TestHandlerImportedAtIsRFC3339UTC (0.00s)
=== RUN   TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse
=== RUN   TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse/not_found
=== RUN   TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse/internal
--- PASS: TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse (0.00s)
    --- PASS: TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse/not_found (0.00s)
    --- PASS: TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse/internal (0.00s)
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import_chain
--- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery (0.00s)
    --- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import (0.00s)
    --- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import_chain (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/transport	2.099s
```

Command 2 — new test run in isolation (`-run` scoped to only this test name), from `apps/server_core`:

```
cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/erp_import/transport/... -run 'TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery' -v
```

Output (full, real):

```
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import
=== RUN   TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import_chain
--- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery (0.00s)
    --- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import (0.00s)
    --- PASS: TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery/get_import_chain (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/transport	1.824s
```

## 3. Which assertion breaks if the validation is removed

If `readImportID` (and its call at the top of `handleGetImport`/`handleGetImportChain` in
`http_handler.go`) were removed or bypassed so that `not-a-uuid` reached the query layer
unchecked, the fake querier's `GetImport`/`GetImportChain` methods would run and set
`f.getID`/`f.chainID` to `"not-a-uuid"` (a non-empty value) before returning their zero-value
`domain.ImportReport{}`/`domain.ImportChain{}` with a `nil` error — which the handler would
then encode as a `200 OK` JSON body instead of a `400`. That breaks two of this test's
assertions at once: `response.Code != http.StatusBadRequest` (200 != 400) fires first, and
even if status handling were patched to still emit 400 by some other path, the
`querier.getID != ""` / `querier.chainID != ""` check is the one that specifically catches
"query ran before validation" — it is the assertion this test exists for, not the status
code check (a handler that queries first and only fails after would still pass a status-only
test but fails this one because `getID`/`chainID` would already be populated).

## 4. What could NOT be done / limitations

Nothing was left incomplete. Both routes (`GET /erp/imports/{id}` and
`GET /erp/imports/{id}/chain`) are covered in one table-driven test, each subtest asserts
status 400, exact body `{"error":"invalid_import_id"}`, and that neither `GetImport` nor
`GetImportChain` ever received a non-empty id (i.e. the query never ran). The full package
test suite passes, and the new test passes in isolation via `-run`. Only the one file listed
above (`http_handler_test.go`) was edited; no production code, OpenAPI, or `packages/` files
were touched.
