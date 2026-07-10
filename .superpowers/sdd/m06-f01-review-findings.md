# M-06 F-01 Independent Review Findings

## Review Verdict

- Spec compliance: `REJECTED`
- Code quality: `REJECTED`

## Blocking Finding

- Severity: P1 / Important
- Location: `apps/server_core/internal/modules/orders/ports/order_source.go`, `apps/server_core/internal/modules/orders/application/import_service.go`
- Finding: the orders-owned port and application expose `connectors/domain.OrderSnapshot`, so the consuming module depends on another module's internal domain type. The project architecture requires the orders port to own its contract and the integrations adapter to translate connector snapshots at the boundary.

## Required Correction

- Define an orders-owned ingestion snapshot contract containing only operational order, item, and payment fields.
- Do not include raw provider payload in the orders-owned type.
- Translate `connectors/domain.OrderSnapshot` into the orders-owned contract inside `orders/adapters/integrations`.
- Remove all connectors imports from `orders/application` and `orders/ports`.
- Preserve the already verified timestamp, PII, cancellation, and persistence behavior.

## Review Notes

The reviewer found the atomic upsert, child-write guard, tenant transaction, timestamp semantics, safe provider reference, and cancelled-order mapping correct. The only blocking finding is the module-boundary violation above.
