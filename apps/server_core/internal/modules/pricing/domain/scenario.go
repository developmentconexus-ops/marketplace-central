package domain

import "time"

// Scenario is a saved set of simulation inputs a tenant can re-run. Payload
// holds the opaque input JSON (the engine re-decomposes from it); the schema
// keeps it as-is. CreatedAt orders the list newest-first.
type Scenario struct {
	ID        string
	Name      string
	Payload   string
	CreatedAt time.Time
}
