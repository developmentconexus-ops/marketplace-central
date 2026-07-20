// Package domain holds the sync_state model (E8): cadence-agnostic bookkeeping
// for the ERP/ML sync engine. MIS-006 M-01 ships the skeleton — table shape,
// entity enum, and error value — while heavy sync jobs land in later missions.
package domain

import (
	"encoding/json"
	"errors"
	"time"
)

// Entity is the semantic kind of stream a sync job tracks. It is validated in
// the application layer (not enforced by a DB constraint) so new entities never
// need a migration — see migration 0075 and F-01.
type Entity string

const (
	EntityProducts Entity = "products"
	EntityListings Entity = "listings"
	EntityOrders   Entity = "orders"
	EntityMarket   Entity = "market"
	EntityTariffs  Entity = "tariffs"
)

// ErrUnknownEntity is returned when a caller registers or persists an entity
// outside the documented enum. Mirrors the fail-closed posture of integrations'
// ErrUnknownActiveSource — no silent acceptance of an unknown value.
var ErrUnknownEntity = errors.New("sync: unknown entity")

// Valid reports whether e is one of the documented sync entities.
func (e Entity) Valid() bool {
	switch e {
	case EntityProducts, EntityListings, EntityOrders, EntityMarket, EntityTariffs:
		return true
	default:
		return false
	}
}

// SyncState is one row of the sync_state table (E8): the bookkeeping for a
// single (tenant, installation, entity) sync stream. Cursor, Schedule, LastError,
// and the timestamps are nil/absent when the underlying JSONB/column is SQL NULL
// — an honest unknown, never a fabricated value (ADR-17).
type SyncState struct {
	TenantID       string
	InstallationID string
	Entity         Entity

	// Cursor is opaque, entity-specific progress state (nil = never synced).
	Cursor json.RawMessage
	// Schedule is a generic cadence descriptor (nil = use the caller default).
	// Cadence-agnostic (D6): no "daily"/cron literal is baked into the type.
	Schedule json.RawMessage

	LastFullSyncAt    *time.Time
	LastIncrementalAt *time.Time

	// LastError carries a generic message + timestamp only. It MUST NOT contain
	// provider payloads, credentials, or tokens.
	LastError *SyncError

	ConsecutiveFailures int
}

// SyncError is the JSONB persisted in sync_state.last_error. Message is generic
// operator-facing text; it never carries provider credentials or payloads.
type SyncError struct {
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}
