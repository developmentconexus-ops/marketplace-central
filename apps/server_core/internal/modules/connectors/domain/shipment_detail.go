package domain

import "time"

// ShipmentDetail is the FULL persistence-shaped shipment fact for `order_shipments` (migration
// 0088, IC-03) — its fields mirror those DB columns 1:1. This is a DIFFERENT, richer type than
// ShipmentInfo (shipping_read.go), which serves a narrower, still-LIVE call site (own-cost
// simulation) — do not merge or repurpose either type.
//
// Every optional/unknown field is a nil pointer (ADR-17): provider absence, a masked value, or
// an unconfirmed JSON key all degrade to nil, never a fabricated zero/blank.
//
// order_shipments carries NO raw column by design (0088, P7 r01 B-7): a shipment payload
// carries delivery PII. There is no RawShipment field here on purpose — the typed fields below
// are the COMPLETE persistence of shipment for this schema.
type ShipmentDetail struct {
	ProviderShipmentID string
	Status             string
	Substatus          *string

	// LogisticType is `logistic.type` on GET /shipments/{id} (self_service, drop_off,
	// xd_drop_off, fulfillment…) — confirmed against a captured fixture (mercado_livre
	// package's testdata/shipment_body.json) and now decoded by
	// mlIngestShipmentResponse/mapShipmentDetail (shipment_ingest_reader.go). Previously
	// undeclared on the wire DTO, so encoding/json silently dropped it (0/38 rows) — same
	// class of defect fixed twice before in this slice for the order DTO. A legacy `mode`
	// field (values like "me2") is a DIFFERENT, older ML concept (shipping method version,
	// not this v2 logistic_type vocabulary) and stays unaliased. Nil only when the provider
	// payload omits `logistic` or sends a blank `type` (ADR-17).
	//
	// TrackingMethod is the top-level `tracking_method` field on the same payload — also
	// confirmed against the fixture above and now decoded. Nil only when absent/blank.
	LogisticType   *string
	TrackingMethod *string

	// TrackingNumber IS confirmed: docs/design/handoff/API-MAP.md (doc review July 2026)
	// lists `tracking_number` as a direct top-level field of GET /shipments/$ID (✅ API
	// direta). Decoded from that key.
	TrackingNumber *string

	// SLAStatus maps to IC-03's stated source, the separate GET /shipments/{id}/sla
	// sub-resource. That endpoint's existence is live-verified (research/external-ml-api-
	// facts.md fact #15), but its response shape could not be confirmed in this worktree
	// (the referenced evidence file F5-shipment-sla-*.json is not present). Rather than call
	// an endpoint and guess its keys, this reader does not call it — stays honest-unknown
	// nil. Follow-up: confirm the shape against a live/fixture capture and wire the call
	// (validation.md gap, does not block this feature per feature.md's own escape hatch for
	// unconfirmable keys).
	SLAStatus *string

	// SLALimitAt does NOT come from the /sla sub-resource (IC-03 is wrong on this specific
	// point — docs/ outranks .mnfs/research/ per this repo's truth order). It comes from the
	// SAME primary GET /shipments/{id} call (x-format-new: true) already made for the rest of
	// this struct: `lead_time.estimated_delivery_limit.date`. Confirmed live: this exact field
	// is already decoded today by the existing, separate mapShipmentInfo/ShipmentInfo.SLADue
	// path (shipping_reader.go:236-237, mlShipmentLeadTime/mlEstimatedDeliveryLimit types
	// defined at shipping_reader.go:77-83) and documented at
	// docs/design/handoff-2026-07/API-MAP.md:18. Reused here via the same private types
	// (same package), no new endpoint call.
	SLALimitAt *time.Time

	// CostGross/CostSeller are decimal-string amounts (order_shipments.cost_gross/
	// cost_seller are numeric(14,2)); Currency is the single shared ISO code for both,
	// resolved from the shipment's site (the /costs payload itself carries no currency
	// field — same fact already established by shipping_read.go's mlSiteCurrency). nil
	// when no cost value was decodable, never a fabricated "0.00"/"".
	CostGross  *string
	CostSeller *string
	Currency   *string

	ReceiverName *string
	DestCity     *string
	DestState    *string
	DestZip      *string

	// SourceTime is the provider's own timestamp for this shipment fact (IC-03 Timestamp
	// And ID Semantics — `date_created` on the primary shipment payload, verbatim). Nil
	// when the payload carries none.
	SourceTime *time.Time
	FetchedAt  time.Time

	// Found is false on an honest-absence 404/410 (deleted/inaccessible shipment,
	// feature.md Negative Scenarios) — the zero-value ShipmentDetail carries no meaning on
	// its own; callers MUST check Found before trusting any other field.
	Found bool
}
