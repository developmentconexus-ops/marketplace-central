package postgres

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func TestBuildBuyerFiscalInfo_AllNullIsHonestAbsence(t *testing.T) {
	info := buildBuyerFiscalInfo(
		pgtype.Text{}, pgtype.Text{}, pgtype.Text{},
		pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{},
	)
	if info.HasData() {
		t.Fatalf("HasData() = true, want false for a row with all 9 buyer_* columns NULL")
	}
	if info.Address != nil {
		t.Fatalf("Address = %+v, want nil (no address field present)", info.Address)
	}
	if info.Name != nil || info.DocType != nil || info.DocNumber != nil {
		t.Fatalf("Name/DocType/DocNumber = %v/%v/%v, want all nil", info.Name, info.DocType, info.DocNumber)
	}
}

func TestBuildBuyerFiscalInfo_NameOnlyHasDataTrueNoAddress(t *testing.T) {
	info := buildBuyerFiscalInfo(
		pgtype.Text{String: "Fulano de Tal", Valid: true}, pgtype.Text{}, pgtype.Text{},
		pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{}, pgtype.Text{},
	)
	if !info.HasData() {
		t.Fatal("HasData() = false, want true (Name is present)")
	}
	if info.Address != nil {
		t.Fatalf("Address = %+v, want nil (no address column present, only name)", info.Address)
	}
	if info.Name == nil || *info.Name != "Fulano de Tal" {
		t.Fatalf("Name = %v, want Fulano de Tal", info.Name)
	}
}

func TestBuildBuyerFiscalInfo_PartialAddressOnlyIncludesPresentFields(t *testing.T) {
	info := buildBuyerFiscalInfo(
		pgtype.Text{}, pgtype.Text{}, pgtype.Text{},
		pgtype.Text{String: "Rua das Flores", Valid: true}, pgtype.Text{}, pgtype.Text{},
		pgtype.Text{String: "SP", Valid: true}, pgtype.Text{}, pgtype.Text{},
	)
	if !info.HasData() {
		t.Fatal("HasData() = false, want true (Address is non-nil)")
	}
	if info.Address == nil {
		t.Fatal("Address = nil, want non-nil (street and state_code are present)")
	}
	if info.Address.StreetName == nil || *info.Address.StreetName != "Rua das Flores" {
		t.Fatalf("Address.StreetName = %v, want Rua das Flores", info.Address.StreetName)
	}
	if info.Address.StreetNumber != nil {
		t.Fatalf("Address.StreetNumber = %v, want nil (column absent)", info.Address.StreetNumber)
	}
	if info.Address.City != nil {
		t.Fatalf("Address.City = %v, want nil (column absent)", info.Address.City)
	}
	if info.Address.StateCode == nil || *info.Address.StateCode != "SP" {
		t.Fatalf("Address.StateCode = %v, want SP", info.Address.StateCode)
	}
	if info.Address.ZipCode != nil || info.Address.CountryID != nil {
		t.Fatalf("Address.ZipCode/CountryID = %v/%v, want both nil", info.Address.ZipCode, info.Address.CountryID)
	}
}

// TestBuildBuyerFiscalInfo_StateNameAlwaysNil is the F-03 brief's explicit
// verification: migration 0089 has no uf_nome/state_name column, so
// Address.StateName must always come back nil regardless of what StateCode
// carries. This is VERIFIED safe (not just documented) because
// apps/web/src/pages/pedidos/PedidoDrawer.tsx's formatEndereco only reads
// uf_codigo, never uf_nome — see buyer_fiscal_reader.go's doc comment.
func TestBuildBuyerFiscalInfo_StateNameAlwaysNil(t *testing.T) {
	info := buildBuyerFiscalInfo(
		pgtype.Text{String: "Fulano de Tal", Valid: true}, pgtype.Text{String: "CPF", Valid: true}, pgtype.Text{String: "12345678900", Valid: true},
		pgtype.Text{String: "Rua das Flores", Valid: true}, pgtype.Text{String: "100", Valid: true}, pgtype.Text{String: "Sao Paulo", Valid: true},
		pgtype.Text{String: "SP", Valid: true}, pgtype.Text{String: "01000-000", Valid: true}, pgtype.Text{String: "BR", Valid: true},
	)
	if info.Address == nil {
		t.Fatal("Address = nil, want non-nil (all address fields present)")
	}
	if info.Address.StateName != nil {
		t.Fatalf("Address.StateName = %v, want nil (0089 has no state_name/uf_nome column)", info.Address.StateName)
	}
	if info.Address.StateCode == nil || *info.Address.StateCode != "SP" {
		t.Fatalf("Address.StateCode = %v, want SP", info.Address.StateCode)
	}
}

// TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData is the
// read-side half of the F-03 golden-test requirement (the transport-side half
// -- mapCompradorFiscal producing byte-identical JSON -- lives in
// orders/transport/http_handler_test.go, since mapCompradorFiscal is
// unexported there and buildBuyerFiscalInfo is unexported here; together they
// prove the whole DB-row -> comprador_fiscal chain is unchanged).
//
// oldAdapterShape hand-mirrors mercado_livre/buyer_fiscal_reader.go's
// mapBuyerFiscalInfo + mapBuyerFiscalAddress field-by-field (Name = trimmed
// "first last", DocType/DocNumber verbatim, Address built when the ML payload
// carried it, StateCode from state.code) for a billing-info payload equivalent
// to the seeded row below, EXCLUDING state.name -- 0089 deliberately excludes
// uf_nome, a VERIFIED-safe ratified shape difference (PedidoDrawer.tsx's
// formatEndereco only reads uf_codigo), not an accepted-gap regression. Both
// sides' FetchedAt are zeroed before comparing: the old adapter always stamps
// a.now() while buildBuyerFiscalInfo never sets it, and this field is
// verified dead downstream (compradorFiscalDTO has no fetched_at key).
func TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData(t *testing.T) {
	name := "Fulano de Tal"
	docType := "CPF"
	docNumber := "12345678900"
	street := "Rua das Flores"
	number := "100"
	city := "Sao Paulo"
	stateCode := "SP"
	zip := "01000-000"
	country := "BR"

	oldAdapterShape := connectorsdomain.BuyerFiscalInfo{
		Name:      &name,
		DocType:   &docType,
		DocNumber: &docNumber,
		Address: &connectorsdomain.BuyerFiscalAddress{
			StreetName:   &street,
			StreetNumber: &number,
			City:         &city,
			StateCode:    &stateCode,
			StateName:    nil,
			ZipCode:      &zip,
			CountryID:    &country,
		},
	}

	newDBReaderShape := buildBuyerFiscalInfo(
		pgtype.Text{String: name, Valid: true}, pgtype.Text{String: docType, Valid: true}, pgtype.Text{String: docNumber, Valid: true},
		pgtype.Text{String: street, Valid: true}, pgtype.Text{String: number, Valid: true}, pgtype.Text{String: city, Valid: true},
		pgtype.Text{String: stateCode, Valid: true}, pgtype.Text{String: zip, Valid: true}, pgtype.Text{String: country, Valid: true},
	)

	oldJSON, err := json.Marshal(oldAdapterShape)
	if err != nil {
		t.Fatalf("marshal old-adapter-shape: %v", err)
	}
	newJSON, err := json.Marshal(newDBReaderShape)
	if err != nil {
		t.Fatalf("marshal new-db-reader-shape: %v", err)
	}
	if string(oldJSON) != string(newJSON) {
		t.Fatalf("BuyerFiscalInfo diverged across the read-path switch for equivalent data (FetchedAt excluded by construction -- both zero-value here):\n  old (live ML adapter shape): %s\n  new (Postgres reader shape): %s", oldJSON, newJSON)
	}
}

func TestBuildBuyerFiscalInfo_FullRowMapsAllNineColumns(t *testing.T) {
	info := buildBuyerFiscalInfo(
		pgtype.Text{String: "Fulano de Tal", Valid: true}, pgtype.Text{String: "CPF", Valid: true}, pgtype.Text{String: "12345678900", Valid: true},
		pgtype.Text{String: "Rua das Flores", Valid: true}, pgtype.Text{String: "100", Valid: true}, pgtype.Text{String: "Sao Paulo", Valid: true},
		pgtype.Text{String: "SP", Valid: true}, pgtype.Text{String: "01000-000", Valid: true}, pgtype.Text{String: "BR", Valid: true},
	)
	want := struct {
		name, docType, docNumber, street, number, city, state, zip, country string
	}{"Fulano de Tal", "CPF", "12345678900", "Rua das Flores", "100", "Sao Paulo", "SP", "01000-000", "BR"}

	if info.Name == nil || *info.Name != want.name {
		t.Fatalf("Name = %v, want %v", info.Name, want.name)
	}
	if info.DocType == nil || *info.DocType != want.docType {
		t.Fatalf("DocType = %v, want %v", info.DocType, want.docType)
	}
	if info.DocNumber == nil || *info.DocNumber != want.docNumber {
		t.Fatalf("DocNumber = %v, want %v", info.DocNumber, want.docNumber)
	}
	if info.Address.StreetNumber == nil || *info.Address.StreetNumber != want.number {
		t.Fatalf("Address.StreetNumber = %v, want %v", info.Address.StreetNumber, want.number)
	}
	if info.Address.City == nil || *info.Address.City != want.city {
		t.Fatalf("Address.City = %v, want %v", info.Address.City, want.city)
	}
	if info.Address.ZipCode == nil || *info.Address.ZipCode != want.zip {
		t.Fatalf("Address.ZipCode = %v, want %v", info.Address.ZipCode, want.zip)
	}
	if info.Address.CountryID == nil || *info.Address.CountryID != want.country {
		t.Fatalf("Address.CountryID = %v, want %v", info.Address.CountryID, want.country)
	}
}
