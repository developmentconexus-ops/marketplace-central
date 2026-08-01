package domain

import "errors"

// ErrOrderUnavailable classifies an order-level ingest fetch failure (third-party 403, order
// 404) that IngestOrder must treat as a NO-WRITE outcome (feature.md Negative Scenarios: "order
// 403/404 -> typed error, no writes at all"). The application layer wraps the underlying
// connectors sentinel (connectorsdomain.ErrUnauthorized/ErrNotFound) with %w so callers can
// still errors.Is against the specific cause; ImportService only needs THIS classification to
// count the attempt as a skip instead of failing the whole run (IC-06: "erro por item no
// ingest: registra, segue o batch").
var ErrOrderUnavailable = errors.New("orders: order unavailable")
