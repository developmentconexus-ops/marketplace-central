// Package sankhyaoracle is the only importable surface of the Sankhya adapter.
// Wire rows and SQL live under internal/oracle and cannot be reached from any
// other package in the module, including the composition root.
//
// The driver is NOT imported here. Registering godror is the composition
// root's job; this package speaks database/sql and therefore builds without
// cgo, which is what keeps its tests from passing vacuously on a host with no
// Oracle client.
package sankhyaoracle

import (
	"database/sql"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/catalogfeed"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
)

// Bundle is everything this adapter offers, typed by the ports its consumers
// own. A field typed by catalogfeed.Feed would publish this adapter's concrete
// type to every caller; a field typed by port.ProductFeed publishes only the
// question Catalog asks.
type Bundle struct {
	Instance    string
	CatalogFeed port.ProductFeed
}

// New builds the bundle for one Sankhya installation.
func New(db *sql.DB, instance string, now func() time.Time) (Bundle, error) {
	feed, err := catalogfeed.NewFeed(db, instance, now)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{Instance: instance, CatalogFeed: feed}, nil
}
