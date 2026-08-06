// Package catalogfeed translates Sankhya rows into Catalog observations. It is
// the only place in the platform that knows both vocabularies at once.
package catalogfeed

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// SourceSystem names Sankhya in every SourceProductKey this package mints.
const SourceSystem = "sankhya"

// ObjectKind names the TGFPRO row kind.
const ObjectKind = "product"

// Feed pages products out of Sankhya as Catalog observations.
type Feed struct {
	db       *sql.DB
	instance string
	now      func() time.Time
}

// NewFeed builds the feed. instance identifies WHICH Sankhya installation this
// is, because two installations both call their first product CODPROD 1.
func NewFeed(db *sql.DB, instance string, now func() time.Time) (Feed, error) {
	if db == nil {
		return Feed{}, fmt.Errorf("catalogfeed: db is required")
	}
	if strings.TrimSpace(instance) == "" {
		return Feed{}, fmt.Errorf("catalogfeed: instance is required")
	}
	if now == nil {
		now = time.Now
	}
	return Feed{db: db, instance: strings.TrimSpace(instance), now: now}, nil
}

// NextPage implements port.ProductFeed. The Sankhya cursor is CODPROD, and that
// fact stops at this function: the token crossing the port is a string, so the
// day a second ERP pages by something else, nothing in Catalog changes.
func (f Feed) NextPage(ctx context.Context, t tenant.ID, after port.Cursor, limit int) (port.Page, error) {
	var afterCodprod int64
	if !after.IsStart() {
		parsed, err := strconv.ParseInt(after.Token(), 10, 64)
		if err != nil {
			return port.Page{}, fmt.Errorf("catalogfeed: cursor %q is not a CODPROD: %w", after.Token(), err)
		}
		afterCodprod = parsed
	}

	rows, err := oracle.FetchActiveProducts(ctx, f.db, afterCodprod, limit, f.now())
	if err != nil {
		return port.Page{}, err
	}
	out := make([]contracts.ProductObservation, 0, len(rows))
	var last int64
	for _, r := range rows {
		obs, err := MapProduct(t, f.instance, r)
		if err != nil {
			return port.Page{}, err
		}
		out = append(out, obs)
		last = r.Codprod
	}
	page := port.Page{Observations: out, Done: len(rows) < limit}
	if !page.Done {
		page.Next = port.NewCursor(strconv.FormatInt(last, 10))
	}
	return page, nil
}

var _ port.ProductFeed = Feed{}

// MapProduct turns one TGFPRO row into one observation.
func MapProduct(t tenant.ID, instance string, row oracle.ProductRow) (contracts.ProductObservation, error) {
	if row.Codprod <= 0 {
		return contracts.ProductObservation{},
			fmt.Errorf("catalogfeed: CODPROD must be positive, got %d", row.Codprod)
	}
	externalKey := strconv.FormatInt(row.Codprod, 10)

	key, err := contracts.NewSourceProductKey(t, SourceSystem, instance, ObjectKind, externalKey)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	hash := payloadHash(row)
	e, err := provenance.NewEvidence(SourceSystem, ObjectKind, externalKey, row.ReadAt.UTC(), hash)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	desc, err := mapDescription(row.Descrprod, e)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	var identifiers []contracts.Identifier
	if ean := strings.TrimSpace(row.Referencia.String); row.Referencia.Valid && ean != "" {
		id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
		if err != nil {
			return contracts.ProductObservation{}, err
		}
		identifiers = append(identifiers, id)
	}

	obs := contracts.ProductObservation{
		Key: key, Description: desc, Identifiers: identifiers, Evidence: e,
	}
	if err := obs.Validate(); err != nil {
		return contracts.ProductObservation{}, err
	}
	return obs, nil
}

// mapDescription is where ADR-017 stops being a document. A NULL column becomes
// Unknown with a reason, never the empty string.
func mapDescription(col sql.NullString, e provenance.Evidence) (fact.Fact[string], error) {
	if !col.Valid {
		return fact.NewUnknown[string]("TGFPRO.DESCRPROD is NULL", e)
	}
	if v := strings.TrimSpace(col.String); v != "" {
		return fact.NewKnown(v, e)
	}
	return fact.NewUnknown[string]("TGFPRO.DESCRPROD is blank", e)
}

// payloadHash digests the business columns and NOTHING else.
//
// ReadAt is deliberately excluded: including it would make every sync run
// produce a new hash, which would defeat idempotence and bump every version on
// every pass.
func payloadHash(row oracle.ProductRow) string {
	h := sha256.New()
	fmt.Fprintf(h, "codprod=%d\n", row.Codprod)
	for _, f := range []struct {
		name string
		col  sql.NullString
	}{
		{"descrprod", row.Descrprod},
		{"referencia", row.Referencia},
		{"marca", row.Marca},
		{"ncm", row.NCM},
		{"ativo", row.Ativo},
	} {
		if f.col.Valid {
			fmt.Fprintf(h, "%s=%s\n", f.name, f.col.String)
		} else {
			// A distinct marker: NULL must not hash the same as "".
			fmt.Fprintf(h, "%s=<null>\n", f.name)
		}
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
