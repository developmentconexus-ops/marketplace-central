package catalogfeed

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func testTenant(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_metal")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

func TestMapProductCarriesDescriptionAsKnown(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:    4711,
		Descrprod:  sql.NullString{String: "Cafeteira Eletrica", Valid: true},
		Referencia: sql.NullString{String: "7891234567895", Valid: true},
		ReadAt:     time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if obs.Description.State() != fact.Known {
		t.Fatalf("Description state = %v, want Known", obs.Description.State())
	}
	if v, _ := obs.Description.Value(); v != "Cafeteira Eletrica" {
		t.Fatalf("Description = %q", v)
	}
	if obs.Key.ExternalKey() != "4711" {
		t.Fatalf("ExternalKey = %q, want 4711", obs.Key.ExternalKey())
	}
	if len(obs.Identifiers) != 1 || obs.Identifiers[0].Value() != "7891234567895" {
		t.Fatalf("Identifiers = %v", obs.Identifiers)
	}
}

// The whole point of the kernel: a NULL column must not arrive downstream as "".
func TestMapProductTurnsNullDescriptionIntoUnknownNotEmptyString(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:   4712,
		Descrprod: sql.NullString{},
		ReadAt:    time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if obs.Description.State() != fact.Unknown {
		t.Fatalf("Description state = %v, want Unknown", obs.Description.State())
	}
	if _, ok := obs.Description.Value(); ok {
		t.Fatal("Unknown description handed out a value")
	}
	if obs.Description.Reason() == "" {
		t.Fatal("Unknown description carries no reason")
	}
}

// A blank REFERENCIA is not an identifier. Emitting "" would make every
// product with a blank EAN collide with every other one.
func TestMapProductDropsBlankIdentifier(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:    4713,
		Descrprod:  sql.NullString{String: "Chaleira", Valid: true},
		Referencia: sql.NullString{String: "   ", Valid: true},
		ReadAt:     time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if len(obs.Identifiers) != 0 {
		t.Fatalf("Identifiers = %v, want none", obs.Identifiers)
	}
}

// Two reads of the same unchanged row must hash the same, or every sync run
// bumps every version and idempotence means nothing.
func TestMapProductHashIsStableAcrossReadTimes(t *testing.T) {
	base := oracle.ProductRow{
		Codprod:   4714,
		Descrprod: sql.NullString{String: "Torradeira", Valid: true},
		ReadAt:    time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	later := base
	later.ReadAt = base.ReadAt.Add(6 * time.Hour)

	a, err := MapProduct(testTenant(t), "sankhya-it-01", base)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	b, err := MapProduct(testTenant(t), "sankhya-it-01", later)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if a.Evidence.PayloadHash() != b.Evidence.PayloadHash() {
		t.Fatalf("hash moved with read time: %q vs %q",
			a.Evidence.PayloadHash(), b.Evidence.PayloadHash())
	}
	if a.Evidence.ObservedAt().Equal(b.Evidence.ObservedAt()) {
		t.Fatal("ObservedAt did not move, so the test proves nothing")
	}
}

func TestMapProductRejectsZeroCodprod(t *testing.T) {
	_, err := MapProduct(testTenant(t), "sankhya-it-01", oracle.ProductRow{Codprod: 0})
	if err == nil {
		t.Fatal("MapProduct accepted CODPROD 0")
	}
}

// A malformed cursor must fail before the query ever reaches Oracle: the
// zero-value Feed (no db, no instance, no now) proves it, because NextPage
// would nil-pointer-panic on f.db the moment it got past the parse.
func TestNextPageRejectsMalformedCursorWithoutTouchingDB(t *testing.T) {
	f := Feed{}
	page, err := f.NextPage(context.Background(), testTenant(t), port.NewCursor("not-a-number"), 10)
	if err == nil {
		t.Fatal("NextPage accepted a cursor that is not a CODPROD")
	}
	if !strings.Contains(err.Error(), "not-a-number") {
		t.Fatalf("error %q does not name the malformed token", err.Error())
	}
	if len(page.Observations) != 0 || page.Next.Token() != "" || page.Done {
		t.Fatalf("page = %+v, want the zero value", page)
	}
}

var _ = contracts.ProductObservation{}
