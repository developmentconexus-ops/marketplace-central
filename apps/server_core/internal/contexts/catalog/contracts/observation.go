package contracts

import (
	"fmt"

	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// ProductObservation is what a source feed hands to Catalog: one product as one
// system saw it at one moment. It is not a product — Catalog decides that.
type ProductObservation struct {
	Key         SourceProductKey
	Description fact.Fact[string]
	Identifiers []Identifier
	Evidence    provenance.Evidence
}

// Validate rejects an observation that cannot be acted on. It deliberately does
// NOT require a known description: a source that returned nothing is a fact
// about the source, and Catalog records it as Unknown rather than refusing the
// product or inventing a name.
func (o ProductObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: source product key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	for i, id := range o.Identifiers {
		if id.Value() == "" {
			return fmt.Errorf("%w: identifier at index %d", ErrBlank, i)
		}
	}
	return nil
}

// Disposition is what ingesting an observation did.
type Disposition string

const (
	// DispositionCreated means a new canonical product was minted.
	DispositionCreated Disposition = "created"
	// DispositionChanged means an existing product gained a new version.
	DispositionChanged Disposition = "changed"
	// DispositionIdempotent means this exact payload was already recorded and
	// nothing moved. Re-polling must be free.
	DispositionIdempotent Disposition = "idempotent"
)

// IngestResult reports what happened, including conflicts a human must see.
//
// DuplicateIdentifiers is not an error. Catalog still creates a distinct
// product, because silently merging two ERP codes that share a bad EAN is a
// data loss that cannot be undone; surfacing the conflict is reversible.
type IngestResult struct {
	ProductID            string
	Disposition          Disposition
	Version              int
	DuplicateIdentifiers []Identifier
}
