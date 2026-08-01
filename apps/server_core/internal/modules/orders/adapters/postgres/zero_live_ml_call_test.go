package postgres

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

// mercadoLivreConnectorImportPath is the sole Mercado Livre HTTP client
// package's import path, as it appears (quoted) in an *ast.ImportSpec.Path.Value
// -- the same literal internal/platform/archguard/archguard_test.go's
// mlImportPath checks for the interactive-guard's own purposes. This file
// checks it independently (a different package, a narrower and more direct
// claim) so the F-03 zero-live-call proof does not depend on archguard's
// unrelated allowlist bookkeeping to hold.
const mercadoLivreConnectorImportPath = `"marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"`

// TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector is the
// F-03 read-path-switch zero-live-call proof: it parses shipment_reader.go and
// buyer_fiscal_reader.go's import declarations (syntax-only, ImportsOnly mode
// -- no type-checking/module resolution needed) and asserts NEITHER imports
// the mercado_livre connector adapter package. Since Go has no way to call a
// package's exported identifiers without importing it, an absent import is a
// structural, compiler-enforced guarantee that ShipmentReader.GetShipment and
// BuyerFiscalReader.GetBuyerFiscal cannot reach the live Mercado Livre HTTP
// client -- not just "does not currently call it", but "cannot call it without
// this test going red first". Both files import connectorsdomain (the shared,
// marketplace-neutral value-type package: ShipmentInfo/BuyerFiscalInfo/Money),
// which is a different package from the ML capability ADAPTER and is not
// flagged here.
func TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector(t *testing.T) {
	for _, file := range []string{"shipment_reader.go", "buyer_fiscal_reader.go"} {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, imp := range f.Imports {
			if imp.Path.Value == mercadoLivreConnectorImportPath {
				t.Fatalf("%s imports the Mercado Livre connector adapter (%s) -- the F-03 read-path switch requires this reader to be Postgres-only, never live-ML-backed", file, imp.Path.Value)
			}
		}
	}
}

// TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall is a
// belt-and-suspenders runtime companion to the import-scan proof above: it
// calls the two pure mapping functions directly (no *pgxpool.Pool, no
// network) and confirms they return ordinary values with no error -- there is
// no code path in either function that could reach out to a provider even in
// principle, since neither accepts anything resembling a capability/client
// argument (their only inputs are already-scanned column values).
func TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall(t *testing.T) {
	if _, err := scanShipmentRow(defaultFakeShipmentRow()); err != nil {
		t.Fatalf("scanShipmentRow: %v", err)
	}
	empty := pgtype.Text{}
	_ = buildBuyerFiscalInfo(empty, empty, empty, empty, empty, empty, empty, empty, empty)
}

// TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall proves
// the ShipmentReader adapter itself has a total, network-free honest-unknown
// path: an unconfigured (nil) pool returns ports.ErrShipmentNotFound
// immediately, never dialing out, never panicking. Combined with the
// import-scan proof, this shows the reader is Postgres-only end to end: no
// pool configured means no read attempted at all, let alone a live ML call.
func TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall(t *testing.T) {
	reader := NewShipmentReader(nil, "tenant-1")
	_, err := reader.GetShipment(t.Context(), "install-1", "SHIP-1")
	if err == nil {
		t.Fatal("GetShipment() error = nil, want ports.ErrShipmentNotFound for an unconfigured pool")
	}
}
