// Package domain holds the channel_fees ledger's provider-agnostic types and
// write-time validation (IC-01, ADR-009). It must never import a Mercado
// Livre / provider-specific type (Q6 architectural boundary, proven by
// boundary_test.go in the parent package) — every value here is generic
// across providers, keyed by (tenant, provider, installation, subject).
package domain

import (
	"encoding/json"
	"time"
)

// Layer values (channel_fees.layer CHECK 1,2,3 — IC-01 Fields).
const (
	LayerCategory = 1
	LayerListing  = 2
	LayerOrder    = 3
)

// FeeKind values (channel_fees.fee_kind CHECK — IC-01 Fields).
const (
	FeeKindCommission = "commission"
	FeeKindFreight    = "freight"
)

// SubjectType values (channel_fees.subject_type CHECK — IC-01 Fields).
const (
	SubjectTypeCategory  = "category"
	SubjectTypeListing   = "listing"
	SubjectTypeOrder     = "order"
	SubjectTypeOrderLine = "order_line"
)

// ValueType values (channel_fees.value_type CHECK — IC-01 Fields).
const (
	ValueTypePercent = "percent"
	ValueTypeAmount  = "amount"
)

// Origem values (channel_fees.origem CHECK — IC-01 Fields). ApiShippingOptions
// is an additive reservation: no producer in this mission writes it (IC-01 —
// listing freight stays honest-unknown until a future producer lands).
const (
	OrigemAPIListingPrices   = "api_listing_prices"
	OrigemAPIShippingOptions = "api_shipping_options"
	OrigemAPIOrder           = "api_order"
	OrigemAPIShipment        = "api_shipment"
	OrigemConfig             = "config"
)

// Fee is one channel_fees row (IC-01 Fields — verbatim binding). Value is a
// plain decimal string (e.g. "15.99"), matching this codebase's money
// convention: exact decimal literals, never float64 (pricing/domain/decimal.go
// G1 — the same rule applies here since channel_fees.value is numeric(14,4)).
type Fee struct {
	TenantID       string
	Provider       string
	InstallationID string
	Layer          int
	FeeKind        string
	SubjectType    string
	SubjectID      string
	ValueType      string
	Value          string
	Currency       *string         // nil when ValueType == percent; NOT NULL when amount (DB CHECK enforces the pairing)
	Detail         json.RawMessage // nil-able; verbatim jsonb decomposition of the source (IC-01 Canonical Examples)
	Origem         string
	ColetadoEm     time.Time // our fetch clock — never equated with SourceTime for convenience (IC-01 Timestamp semantics)
	SourceTime     *time.Time
}

// ResolvedFee is the honest resolution result for one fee_kind under
// ResolveListingFees. Found=false is the typed-absent sentinel (IC-01 Enums
// And Statuses: ledger-only 2→1→absent, never a zero amount, never a
// default) — the consumer (M-07) decides any config fallback, never this
// reader. Detail is the JSONB verbatim from the resolved row (IC-01 Required
// Outputs) — never synthesized, never the config's shape.
type ResolvedFee struct {
	Found      bool
	Layer      int
	ValueType  string
	Value      string
	Currency   *string
	Detail     json.RawMessage
	Origem     string
	ColetadoEm time.Time
}

// ResolvedListingFees is the per-listing output of ResolveListingFees: one
// ladder result per fee_kind. Commission resolves 2→1→absent; Freight
// resolves 2→absent only (IC-01 Enums And Statuses — freight has no layer 1).
// Layer 3 never participates in either ladder (ADR-009 — layer 3 rows may
// exist, written by other producers, but this reader ignores them).
type ResolvedListingFees struct {
	Commission ResolvedFee
	Freight    ResolvedFee
}
