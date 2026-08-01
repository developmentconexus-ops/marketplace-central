package application

import (
	"context"
	"errors"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type fakeOrderDetailReader struct {
	detail connectorsdomain.OrderDetail
	err    error
	calls  []string
}

func (f *fakeOrderDetailReader) GetOrderDetail(_ context.Context, _ string, providerOrderID string) (connectorsdomain.OrderDetail, error) {
	f.calls = append(f.calls, providerOrderID)
	return f.detail, f.err
}

type fakeShipmentDetailReader struct {
	detail connectorsdomain.ShipmentDetail
	err    error
	calls  []string
}

func (f *fakeShipmentDetailReader) GetShipmentDetail(_ context.Context, _ string, shipmentID string) (connectorsdomain.ShipmentDetail, error) {
	f.calls = append(f.calls, shipmentID)
	return f.detail, f.err
}

type fakeIngestBuyerFiscalReader struct {
	info  connectorsdomain.BuyerFiscalInfo
	err   error
	calls int
}

func (f *fakeIngestBuyerFiscalReader) GetBuyerFiscal(context.Context, string, string) (connectorsdomain.BuyerFiscalInfo, error) {
	f.calls++
	return f.info, f.err
}

type fakeOrderIngestStore struct {
	faturadoAt   *time.Time
	faturadoErr  error
	ingestErr    error
	ingestCalled bool
	gotOrder     ordersdomain.MarketplaceOrder
	gotShipment  *ordersdomain.OrderShipment
}

func (f *fakeOrderIngestStore) GetFaturadoAt(context.Context, string, string) (*time.Time, error) {
	return f.faturadoAt, f.faturadoErr
}

func (f *fakeOrderIngestStore) IngestOrder(_ context.Context, order ordersdomain.MarketplaceOrder, shipment *ordersdomain.OrderShipment) error {
	f.ingestCalled = true
	f.gotOrder = order
	f.gotShipment = shipment
	return f.ingestErr
}

func baseOrderDetail(providerOrderID string) connectorsdomain.OrderDetail {
	now := time.Date(2026, 7, 30, 10, 0, 0, 0, time.UTC)
	return connectorsdomain.OrderDetail{
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   providerOrderID,
		ProviderStatus:    "paid",
		ProviderUpdatedAt: &now,
		FetchedAt:         now,
		Items: []connectorsdomain.OrderDetailItem{{
			ProviderItemID: "MLB1",
			SellerSKU:      "SKU-1",
			Quantity:       1,
		}},
	}
}

func newTestIngestService(t *testing.T, detailReader *fakeOrderDetailReader, shipmentReader *fakeShipmentDetailReader, fiscalReader *fakeIngestBuyerFiscalReader, store *fakeOrderIngestStore) *IngestService {
	t.Helper()
	return NewIngestService(IngestServiceConfig{
		OrderDetail:    detailReader,
		ShipmentDetail: shipmentReader,
		BuyerFiscal:    fiscalReader,
		Links:          stubLinkReader{links: map[string]ordersdomain.ItemLink{}},
		Store:          store,
		Now:            func() time.Time { return time.Date(2026, 7, 30, 11, 0, 0, 0, time.UTC) },
	})
}

// TestIngestOrderUnavailableOnOrderAuthErrorWritesNothing is scenario (e) from the brief: a
// third-party order (403) must classify as domain.ErrOrderUnavailable and never reach the store.
func TestIngestOrderUnavailableOnOrderAuthErrorWritesNothing(t *testing.T) {
	detailReader := &fakeOrderDetailReader{err: connectorsdomain.ErrUnauthorized}
	shipmentReader := &fakeShipmentDetailReader{}
	fiscalReader := &fakeIngestBuyerFiscalReader{}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, detailReader, shipmentReader, fiscalReader, store)

	err := service.IngestOrder(context.Background(), "inst-1", "order-403")
	if !errors.Is(err, ordersdomain.ErrOrderUnavailable) {
		t.Fatalf("IngestOrder() error = %v, want ErrOrderUnavailable", err)
	}
	if !errors.Is(err, connectorsdomain.ErrUnauthorized) {
		t.Fatalf("IngestOrder() error = %v, want to also unwrap the underlying connectors error", err)
	}
	if store.ingestCalled {
		t.Fatal("store.IngestOrder was called on a 403 order — must be zero writes")
	}
	if shipmentReader.calls != nil || fiscalReader.calls != 0 {
		t.Fatalf("shipment/fiscal fetched after order fetch failed: shipment calls=%v fiscal calls=%d", shipmentReader.calls, fiscalReader.calls)
	}
}

// TestIngestOrderUnavailableOnOrderNotFoundWritesNothing covers the 404 half of the same
// negative scenario (403 covered above uses the same code path via errors.Is on either sentinel).
func TestIngestOrderUnavailableOnOrderNotFoundWritesNothing(t *testing.T) {
	detailReader := &fakeOrderDetailReader{err: connectorsdomain.ErrNotFound}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, detailReader, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, store)

	err := service.IngestOrder(context.Background(), "inst-1", "order-404")
	if !errors.Is(err, ordersdomain.ErrOrderUnavailable) {
		t.Fatalf("IngestOrder() error = %v, want ErrOrderUnavailable", err)
	}
	if store.ingestCalled {
		t.Fatal("store.IngestOrder was called on a 404 order — must be zero writes")
	}
}

// TestIngestOrderPersistsWithoutShipmentOnHonestAbsence is scenario (d): a 404/410 shipment
// (Found: false, nil error, per the F-01 reader's documented contract) must still let the order
// persist, with provider_shipment_id set for a future retry but no order_shipments row.
func TestIngestOrderPersistsWithoutShipmentOnHonestAbsence(t *testing.T) {
	detail := baseOrderDetail("order-1")
	shipID := "ship-404"
	detail.ShippingID = &shipID
	detailReader := &fakeOrderDetailReader{detail: detail}
	shipmentReader := &fakeShipmentDetailReader{detail: connectorsdomain.ShipmentDetail{Found: false}}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, detailReader, shipmentReader, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-1"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if !store.ingestCalled {
		t.Fatal("store.IngestOrder was not called — order must persist despite shipment honest-absence")
	}
	if store.gotShipment != nil {
		t.Fatalf("gotShipment = %#v, want nil on honest-absence", store.gotShipment)
	}
	if store.gotOrder.ProviderShipmentID == nil || *store.gotOrder.ProviderShipmentID != shipID {
		t.Fatalf("ProviderShipmentID = %#v, want %q set for future retry", store.gotOrder.ProviderShipmentID, shipID)
	}
}

// TestIngestOrderPersistsShipmentWhenFound covers the positive counterpart of the above:
// order_shipments row IS built and passed to the store when the shipment is found.
func TestIngestOrderPersistsShipmentWhenFound(t *testing.T) {
	detail := baseOrderDetail("order-2")
	shipID := "ship-ok"
	detail.ShippingID = &shipID
	detailReader := &fakeOrderDetailReader{detail: detail}
	shipmentReader := &fakeShipmentDetailReader{detail: connectorsdomain.ShipmentDetail{
		ProviderShipmentID: shipID,
		Status:             "shipped",
		Found:              true,
		FetchedAt:          time.Date(2026, 7, 30, 10, 5, 0, 0, time.UTC),
	}}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, detailReader, shipmentReader, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-2"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if store.gotShipment == nil {
		t.Fatal("gotShipment = nil, want a built OrderShipment when Found=true")
	}
	if store.gotShipment.ProviderShipmentID != shipID || store.gotShipment.Status != "shipped" {
		t.Fatalf("gotShipment = %#v", store.gotShipment)
	}
	// A shipped shipment status must win the bucket over a bare "paid" provider status
	// (domain.DeriveOrderBucket's own documented priority) — proves shipmentStatus was
	// actually threaded into the bucket derivation, not silently dropped.
	if store.gotOrder.Bucket != ordersdomain.BucketEnviado {
		t.Fatalf("Bucket = %q, want %q", store.gotOrder.Bucket, ordersdomain.BucketEnviado)
	}
}

// TestIngestOrderSkipsShipmentFetchWhenNoShippingID: no shipping.id on the order means
// GetShipmentDetail must never be called at all (feature.md Negative Scenarios).
func TestIngestOrderSkipsShipmentFetchWhenNoShippingID(t *testing.T) {
	detail := baseOrderDetail("order-3")
	detailReader := &fakeOrderDetailReader{detail: detail}
	shipmentReader := &fakeShipmentDetailReader{}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, detailReader, shipmentReader, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-3"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if shipmentReader.calls != nil {
		t.Fatalf("shipment reader was called = %v, want no call when ShippingID is nil", shipmentReader.calls)
	}
	if store.gotOrder.ProviderShipmentID != nil {
		t.Fatalf("ProviderShipmentID = %#v, want nil", store.gotOrder.ProviderShipmentID)
	}
}

// TestIngestOrderAlwaysFetchesBuyerFiscal documents the judgment call: unlike
// EnrichService.Enrich's list path (FINDING-M08-LIST-TIMEOUT), IngestOrder is off the
// interactive request path and always fetches fiscal (IC-03 Operations table).
func TestIngestOrderAlwaysFetchesBuyerFiscal(t *testing.T) {
	detail := baseOrderDetail("order-4")
	name := "Comprador Teste"
	fiscalReader := &fakeIngestBuyerFiscalReader{info: connectorsdomain.BuyerFiscalInfo{Name: &name}}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, &fakeShipmentDetailReader{}, fiscalReader, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-4"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if fiscalReader.calls != 1 {
		t.Fatalf("fiscal reader calls = %d, want exactly 1", fiscalReader.calls)
	}
	if store.gotOrder.BuyerName == nil || *store.gotOrder.BuyerName != name {
		t.Fatalf("BuyerName = %#v, want %q", store.gotOrder.BuyerName, name)
	}
}

// TestIngestOrderPropagatesShipmentRealErrorWithoutWriting: a genuine shipment-fetch failure
// (not the 404/410 honest-absence the F-01 reader degrades internally) must stop the ingest
// before any write, same as the order-level 403/404 case — nothing is written until every fetch
// has already succeeded.
func TestIngestOrderPropagatesShipmentRealErrorWithoutWriting(t *testing.T) {
	detail := baseOrderDetail("order-5")
	shipID := "ship-5"
	detail.ShippingID = &shipID
	shipmentReader := &fakeShipmentDetailReader{err: connectorsdomain.ErrProviderUnavailable}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, shipmentReader, &fakeIngestBuyerFiscalReader{}, store)

	err := service.IngestOrder(context.Background(), "inst-1", "order-5")
	if !errors.Is(err, connectorsdomain.ErrProviderUnavailable) {
		t.Fatalf("IngestOrder() error = %v, want the real shipment error to propagate", err)
	}
	if store.ingestCalled {
		t.Fatal("store.IngestOrder was called after a real shipment fetch error — must be zero writes")
	}
}

// TestIngestOrderPropagatesFaturadoLookupError: the store's own faturado read failing must not
// be silently treated as "not faturado" — an honest read failure aborts before any write.
func TestIngestOrderPropagatesFaturadoLookupError(t *testing.T) {
	detail := baseOrderDetail("order-6")
	wantErr := errors.New("db unavailable")
	store := &fakeOrderIngestStore{faturadoErr: wantErr}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, store)

	err := service.IngestOrder(context.Background(), "inst-1", "order-6")
	if !errors.Is(err, wantErr) {
		t.Fatalf("IngestOrder() error = %v, want %v", err, wantErr)
	}
	if store.ingestCalled {
		t.Fatal("store.IngestOrder was called after GetFaturadoAt failed")
	}
}

// TestIngestOrderDerivesFaturarBucketWhenPaidAndNotFaturado / TestIngestOrderDerivesEnviarBucketWhenPaidAndFaturado
// are the two truth-table arms from the brief's required scenario (a) (paid+faturado case here;
// the delivered-tag arm is TestIngestOrderPersistsShipmentWhenFound's shipped-status sibling
// covered above, plus the tag-driven case below).
func TestIngestOrderDerivesFaturarBucketWhenPaidAndNotFaturado(t *testing.T) {
	detail := baseOrderDetail("order-7")
	store := &fakeOrderIngestStore{faturadoAt: nil}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-7"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if store.gotOrder.Bucket != ordersdomain.BucketFaturar {
		t.Fatalf("Bucket = %q, want %q", store.gotOrder.Bucket, ordersdomain.BucketFaturar)
	}
}

func TestIngestOrderDerivesEnviarBucketWhenPaidAndFaturado(t *testing.T) {
	detail := baseOrderDetail("order-8")
	faturadoAt := time.Date(2026, 7, 29, 9, 0, 0, 0, time.UTC)
	store := &fakeOrderIngestStore{faturadoAt: &faturadoAt}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-8"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if store.gotOrder.Bucket != ordersdomain.BucketEnviar {
		t.Fatalf("Bucket = %q, want %q", store.gotOrder.Bucket, ordersdomain.BucketEnviar)
	}
}

func TestIngestOrderDerivesEnviadoBucketFromDeliveredTag(t *testing.T) {
	detail := baseOrderDetail("order-9")
	detail.Tags = []string{"delivered"}
	store := &fakeOrderIngestStore{}
	service := newTestIngestService(t, &fakeOrderDetailReader{detail: detail}, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, store)

	if err := service.IngestOrder(context.Background(), "inst-1", "order-9"); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}
	if store.gotOrder.Bucket != ordersdomain.BucketEnviado {
		t.Fatalf("Bucket = %q, want %q", store.gotOrder.Bucket, ordersdomain.BucketEnviado)
	}
}

func TestIngestOrderRejectsBlankIdentifiers(t *testing.T) {
	service := newTestIngestService(t, &fakeOrderDetailReader{}, &fakeShipmentDetailReader{}, &fakeIngestBuyerFiscalReader{}, &fakeOrderIngestStore{})
	if err := service.IngestOrder(context.Background(), "", "order-1"); err == nil {
		t.Fatal("IngestOrder() error = nil, want blank installation id rejected")
	}
	if err := service.IngestOrder(context.Background(), "inst-1", ""); err == nil {
		t.Fatal("IngestOrder() error = nil, want blank provider order id rejected")
	}
}
