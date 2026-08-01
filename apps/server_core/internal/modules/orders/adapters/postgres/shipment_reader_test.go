package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// fakeShipmentRow feeds a fixed 12-column sequence into scanShipmentRow
// without a database. Column order mirrors the SELECT in
// ShipmentReader.GetShipment exactly: provider_shipment_id, status,
// substatus, sla_limit_at, dest_state, dest_city, dest_zip, receiver_name,
// cost_gross::text, cost_seller::text, currency, fetched_at.
type fakeShipmentRow struct {
	providerShipmentID string
	status             string
	substatus          pgtype.Text
	slaLimitAt         pgtype.Timestamptz
	destState          pgtype.Text
	destCity           pgtype.Text
	destZip            pgtype.Text
	receiverName       pgtype.Text
	costGross          pgtype.Text
	costSeller         pgtype.Text
	currency           pgtype.Text
	fetchedAt          pgtype.Timestamptz
}

func (f fakeShipmentRow) Scan(dest ...any) error {
	*dest[0].(*string) = f.providerShipmentID
	*dest[1].(*string) = f.status
	*dest[2].(*pgtype.Text) = f.substatus
	*dest[3].(*pgtype.Timestamptz) = f.slaLimitAt
	*dest[4].(*pgtype.Text) = f.destState
	*dest[5].(*pgtype.Text) = f.destCity
	*dest[6].(*pgtype.Text) = f.destZip
	*dest[7].(*pgtype.Text) = f.receiverName
	*dest[8].(*pgtype.Text) = f.costGross
	*dest[9].(*pgtype.Text) = f.costSeller
	*dest[10].(*pgtype.Text) = f.currency
	*dest[11].(*pgtype.Timestamptz) = f.fetchedAt
	return nil
}

func defaultFakeShipmentRow() fakeShipmentRow {
	return fakeShipmentRow{
		providerShipmentID: "SHIP-1",
		status:             "shipped",
		fetchedAt:          pgtype.Timestamptz{Time: time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC), Valid: true},
	}
}

// TestScanShipmentRow_MoneyAmountSatisfiesValidateMoneyRegex is the F-03
// brief's load-bearing evidence requirement: cost_gross/cost_seller are cast
// to ::text in SQL specifically to avoid a float64 round-trip that could
// produce scientific notation (e.g. "4.25e+07") and fail
// connectorsdomain.ValidateMoney's decimal regex. This test proves the
// resulting Money.Amount actually satisfies that regex for a realistic
// numeric(14,2) text value.
func TestScanShipmentRow_MoneyAmountSatisfiesValidateMoneyRegex(t *testing.T) {
	row := defaultFakeShipmentRow()
	row.costGross = pgtype.Text{String: "42.50", Valid: true}
	row.costSeller = pgtype.Text{String: "13.45", Valid: true}
	row.currency = pgtype.Text{String: "BRL", Valid: true}

	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	if info.Costs == nil {
		t.Fatal("Costs = nil, want populated (cost_gross and cost_seller and currency all present)")
	}
	if err := connectorsdomain.ValidateMoney("costs.gross_amount", info.Costs.GrossAmount); err != nil {
		t.Fatalf("ValidateMoney(GrossAmount) = %v, want nil (Amount=%q)", err, info.Costs.GrossAmount.Amount)
	}
	if err := connectorsdomain.ValidateMoney("costs.sender_cost", info.Costs.SenderCost); err != nil {
		t.Fatalf("ValidateMoney(SenderCost) = %v, want nil (Amount=%q)", err, info.Costs.SenderCost.Amount)
	}
	if info.Costs.GrossAmount.Amount != "42.50" {
		t.Fatalf("GrossAmount.Amount = %q, want %q (no reformatting, verbatim from ::text cast)", info.Costs.GrossAmount.Amount, "42.50")
	}
	// ReceiverCost has no backing column in 0088 — must stay nil (accepted gap).
	if info.Costs.ReceiverCost != nil {
		t.Fatalf("ReceiverCost = %+v, want nil (0088 carries no receiver-cost column)", info.Costs.ReceiverCost)
	}
}

func TestScanShipmentRow_CostAbsentWithoutCurrencyStaysHonestUnknown(t *testing.T) {
	row := defaultFakeShipmentRow()
	row.costGross = pgtype.Text{String: "42.50", Valid: true}
	row.currency = pgtype.Text{} // NULL currency

	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	if info.Costs != nil {
		t.Fatalf("Costs = %+v, want nil (a cost value without a currency must not fabricate a Money)", info.Costs)
	}
}

func TestScanShipmentRow_NoCostColumnsStaysNilCosts(t *testing.T) {
	row := defaultFakeShipmentRow()
	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	if info.Costs != nil {
		t.Fatalf("Costs = %+v, want nil", info.Costs)
	}
}

func TestScanShipmentRow_AlwaysNilFieldsPer0088Gap(t *testing.T) {
	row := defaultFakeShipmentRow()
	row.costGross = pgtype.Text{String: "1.00", Valid: true}
	row.costSeller = pgtype.Text{String: "1.00", Valid: true}
	row.currency = pgtype.Text{String: "BRL", Valid: true}

	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	// Delayed/CarrierName/TrackingURL/Costs.ReceiverCost are ALWAYS nil: 0088
	// has no such columns (accepted, documented UX gap in validation.md).
	if info.Delayed != nil {
		t.Fatalf("Delayed = %v, want nil", info.Delayed)
	}
	if info.CarrierName != nil {
		t.Fatalf("CarrierName = %v, want nil", info.CarrierName)
	}
	if info.TrackingURL != nil {
		t.Fatalf("TrackingURL = %v, want nil", info.TrackingURL)
	}
	if info.Costs.ReceiverCost != nil {
		t.Fatalf("Costs.ReceiverCost = %v, want nil", info.Costs.ReceiverCost)
	}
}

func TestScanShipmentRow_DestinationAndReceiverFieldMapping(t *testing.T) {
	row := defaultFakeShipmentRow()
	row.substatus = pgtype.Text{String: "in_hub", Valid: true}
	row.destState = pgtype.Text{String: "SP", Valid: true}
	row.destCity = pgtype.Text{String: "Sao Paulo", Valid: true}
	row.destZip = pgtype.Text{String: "01000-000", Valid: true}
	row.receiverName = pgtype.Text{String: "Fulano de Tal", Valid: true}
	slaDue := time.Date(2026, 7, 25, 18, 0, 0, 0, time.UTC)
	row.slaLimitAt = pgtype.Timestamptz{Time: slaDue, Valid: true}

	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	if info.ID != "SHIP-1" || info.Status != "shipped" || info.Substatus != "in_hub" {
		t.Fatalf("ID/Status/Substatus = %q/%q/%q, want SHIP-1/shipped/in_hub", info.ID, info.Status, info.Substatus)
	}
	if info.SLADue == nil || !info.SLADue.Equal(slaDue) {
		t.Fatalf("SLADue = %v, want %v", info.SLADue, slaDue)
	}
	if info.DestinationUF == nil || *info.DestinationUF != "SP" {
		t.Fatalf("DestinationUF = %v, want SP", info.DestinationUF)
	}
	if info.DestinationCity == nil || *info.DestinationCity != "Sao Paulo" {
		t.Fatalf("DestinationCity = %v, want Sao Paulo", info.DestinationCity)
	}
	if info.DestinationZip == nil || *info.DestinationZip != "01000-000" {
		t.Fatalf("DestinationZip = %v, want 01000-000", info.DestinationZip)
	}
	if info.ReceiverName == nil || *info.ReceiverName != "Fulano de Tal" {
		t.Fatalf("ReceiverName = %v, want Fulano de Tal", info.ReceiverName)
	}
}

func TestScanShipmentRow_AbsentDestinationFieldsStayNil(t *testing.T) {
	row := defaultFakeShipmentRow()
	info, err := scanShipmentRow(row)
	if err != nil {
		t.Fatalf("scanShipmentRow returned error: %v", err)
	}
	if info.SLADue != nil {
		t.Fatalf("SLADue = %v, want nil", info.SLADue)
	}
	if info.DestinationUF != nil || info.DestinationCity != nil || info.DestinationZip != nil || info.ReceiverName != nil {
		t.Fatalf("destination/receiver fields = %+v, want all nil", info)
	}
	if info.Substatus != "" {
		t.Fatalf("Substatus = %q, want empty", info.Substatus)
	}
}
