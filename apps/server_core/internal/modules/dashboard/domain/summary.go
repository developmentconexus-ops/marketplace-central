package domain

import (
	"errors"
	"time"
)

var ErrInstallationNotFound = errors.New("installation_not_found")

type InstallationNotFoundError struct {
	InstallationID string
}

func (e *InstallationNotFoundError) Error() string { return ErrInstallationNotFound.Error() }
func (e *InstallationNotFoundError) Unwrap() error { return ErrInstallationNotFound }

type Summary struct {
	SyncErrors   *int64                `json:"sync_errors"`
	PendingLinks *int64                `json:"pending_links"`
	BelowMargin  *int64                `json:"below_margin"`
	MissingGTIN  *int64                `json:"missing_gtin"`
	OrdersToday  *int64                `json:"orders_today"`
	Orders7d     *int64                `json:"orders_7d"`
	LastSyncAt   map[string]*time.Time `json:"last_sync_at"`
	Degraded     []string              `json:"degraded"`
	AsOf         time.Time             `json:"as_of"`
}
