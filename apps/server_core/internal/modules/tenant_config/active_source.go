// Package tenant_config holds the per-tenant active-source configuration: which
// ERP/catalog source (sankhya, xlsx, catalogo_cliente) is authoritative for a
// tenant's product data, and the source-kind (upload_snapshot vs
// live_read_through) that follows from it.
package tenant_config

import (
	"context"
	"time"

	internalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

// ActiveSource names the source a tenant has selected as authoritative for
// product data.
type ActiveSource string

const (
	SourceSankhya         ActiveSource = "sankhya"
	SourceXLSX            ActiveSource = "xlsx"
	SourceCatalogoCliente ActiveSource = "catalogo_cliente"
)

// SellableAssortment is the tenant's product-assortment policy.
type SellableAssortment struct {
	OnlyRevenda           bool
	OnlyEmEstoque         bool
	OnlyEcommerceEligible bool
}

// Valid reports whether s is one of the three ratified sources. Anything else
// (including the zero value) is invalid.
func (s ActiveSource) Valid() bool {
	switch s {
	case SourceSankhya, SourceXLSX, SourceCatalogoCliente:
		return true
	default:
		return false
	}
}

// DefaultKind returns the SourceKind a source implies when no kind is set
// explicitly: sankhya is read live (no protocol/snapshot), xlsx and
// catalogo_cliente are upload-shaped. An invalid source yields the zero
// SourceKind ("") — callers must not coerce an unknown source to a default
// kind.
func (s ActiveSource) DefaultKind() sourcekind.SourceKind {
	switch s {
	case SourceSankhya:
		return sourcekind.LiveReadThrough
	case SourceXLSX, SourceCatalogoCliente:
		return sourcekind.UploadSnapshot
	default:
		return sourcekind.SourceKind("")
	}
}

// Config is a tenant's active-source selection as stored/read.
type Config struct {
	TenantID           string
	Source             ActiveSource
	Kind               sourcekind.SourceKind
	SetAt              time.Time
	SetBy              string
	SellableAssortment SellableAssortment
}

// ErrUnknownActiveSource is returned when a tenant has no active-source row.
// It aliases the existing erp_import/adapters/internalread sentinel (reused,
// not reimplemented) so both the routing reader and this package's Repository
// share one identity for transport to map to 400. Fail-closed: unknown never
// silently defaults (ADR-17).
var ErrUnknownActiveSource = internalread.ErrUnknownActiveSource

// ActiveSourceLookup resolves a tenant's active-source Config. Implementations
// return ErrUnknownActiveSource when the tenant has no row. Consumers: the
// routing reader (S4, next slice) and this package's Repository.
type ActiveSourceLookup interface {
	Get(ctx context.Context, tenantID string) (Config, error)
}
