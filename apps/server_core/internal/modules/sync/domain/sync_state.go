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
	// EntityMarketQueue is the ERP-import pending-codigos queue (erp_import's
	// MarketEnqueuer). It is a queue, not a sync stream — it never records a
	// success — and shares no state with EntityMarket's periodic collection job.
	// Kept as its own entity so two unrelated producers never collide on the
	// same (tenant, installation, entity) sync_state row (migration 0093).
	EntityMarketQueue Entity = "market_queue"
	EntityTariffs     Entity = "tariffs"
	// EntityICMSMatrix é o espelho de (uf_origem, uf_destino, grupo_icms) lido
	// do TGFICM (icms_matrix_mirror). Entidade própria, e não um estágio do
	// stream de produtos, porque a matriz é fiscal e muda por vigência de
	// legislação — a cadência dela não é a do catálogo, e a falha de uma não
	// pode se esconder no sucesso da outra.
	EntityICMSMatrix Entity = "icms_matrix"
)

// ErrUnknownEntity is returned when a caller registers or persists an entity
// outside the documented enum. Mirrors the fail-closed posture of integrations'
// ErrUnknownActiveSource — no silent acceptance of an unknown value.
var ErrUnknownEntity = errors.New("sync: unknown entity")

// Valid reports whether e is one of the documented sync entities.
func (e Entity) Valid() bool {
	switch e {
	case EntityProducts, EntityListings, EntityOrders, EntityMarket, EntityMarketQueue, EntityTariffs, EntityICMSMatrix:
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
