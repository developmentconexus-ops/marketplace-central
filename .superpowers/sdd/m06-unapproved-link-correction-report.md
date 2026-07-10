# M-06 unapproved-link correction report

## Status

DONE

## Root-cause trace

1. `orders/adapters/productlinks.LinkReader.ResolveLinks` copied `InternalProductID` from every link candidate even though candidate states map only to `unresolved`, `conflict`, or `missing` order-link quality.
2. The same adapter copied `InternalProductID` from every persisted workflow link, including `rejected`, `conflict`, and `unresolved` states.
3. `profitability/application.Service.buildItemInputs` checked only whether `InternalProductID` was nil before calling `internal_read`; it did not require `LinkQuality == resolved`.
4. Therefore an unapproved candidate ID could cross the orders boundary and trigger Oracle cost/tax reads, producing known values for an unresolved link.

The correction removes candidate IDs at the source boundary, exposes a persisted ID only for an explicitly `resolved` workflow link, and independently requires both resolved quality and a nonnil ID before profitability reads internal facts.

## TDD RED evidence

### Orders product-links adapter

Command (from `apps/server_core`, with an absolute local `GOCACHE`):

```text
go test ./internal/modules/orders/adapters/productlinks -run TestResolveLinksExposesInternalProductOnlyForResolvedPersistedLink -count=1 -v
```

Genuine pre-fix failure:

```text
exact-ean-candidate: InternalProductID = 20303, want nil for unresolved link
conflict-candidate: InternalProductID = 20303, want nil for conflict link
unresolved-candidate: InternalProductID = 20303, want nil for unresolved link
rejected-link: InternalProductID = 20303, want nil for rejected link
conflict-link: InternalProductID = 20303, want nil for conflict link
unresolved-link: InternalProductID = 20303, want nil for unresolved link
FAIL marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks
```

The resolved persisted workflow-link case passed during RED, isolating leakage to candidate and non-resolved mappings.

### Profitability application

Command:

```text
go test ./internal/modules/profitability/application -run TestImportMarginInputsDoesNotReadInternalFactsForNonResolvedLinkWithProductID -count=1 -v
```

Genuine pre-fix failure for each `rejected`, `conflict`, `unresolved`, and `missing` case:

```text
internal_read calls = cost:1 tax:1, want 0 for <quality> link
FAIL marketplace-central/apps/server_core/internal/modules/profitability/application
```

## Changed paths

- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader.go`
- `apps/server_core/internal/modules/orders/adapters/productlinks/link_reader_test.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `.superpowers/sdd/m06-unapproved-link-correction-report.md`

No product link, live database row, approval state, staging area, or commit was mutated.

## GREEN evidence

Focused tests passed after their respective minimal production changes:

```text
PASS TestResolveLinksExposesInternalProductOnlyForResolvedPersistedLink
ok marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks

PASS TestImportMarginInputsDoesNotReadInternalFactsForNonResolvedLinkWithProductID
ok marketplace-central/apps/server_core/internal/modules/profitability/application
```

The profitability regression verifies for all four non-resolved qualities that:

- cost reads: 0
- tax reads: 0
- cost, ICMS, IPI, PIS, and COFINS amounts: nil
- input quality: exact mapped link truth (`rejected_link`, `conflict_link`, `unresolved_link`, or `missing`)
- reason: `internal product link is not resolved`

Impacted suites:

```text
go test ./internal/modules/orders/... ./internal/modules/profitability/... -count=1
```

Result: all tested orders and profitability packages passed; domain/ports packages correctly reported no test files.

## Formatting and boundary evidence

```text
gofmt -d <four scoped Go paths>
gofmt -d: clean (no output)
```

Boundary scan:

```text
rg -n 'internal/modules/orders' internal/modules/profitability -g '*.go'
```

Only these matches were returned, all under the permitted adapter:

```text
internal/modules/profitability/adapters/orders/order_reader.go
internal/modules/profitability/adapters/orders/order_reader_test.go
```

## Concerns

- The four scoped Go paths were already part of a heavily dirty shared worktree and currently appear as untracked files. They were not staged or committed, as required.
- Prior live Oracle cost/tax values obtained through the unresolved `MLB4834373620` candidate are invalid as evidence of a resolved-link flow. This correction does not approve or mutate that link, and no new live resolved-link claim is made.
