package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type captureWriter struct{ calls int }

func (w *captureWriter) Apply(context.Context, ports.WriteItem) (ports.WriteOutcome, error) {
	w.calls++
	return ports.WriteOutcome{}, nil
}

type captureLink struct{ calls int }

func (w *captureLink) ApplyLink(_ context.Context, in ports.LinkageInput) (ports.LinkageOutcome, error) {
	w.calls++
	return ports.LinkageOutcome{State: domain.ItemStateSkipped}, nil
}

func TestWriterRouterDispatchesEveryEnabledIntentAndPreservesSkipped(t *testing.T) {
	price, stock, listing, resync, link := &captureWriter{}, &captureWriter{}, &captureWriter{}, &captureWriter{}, &captureLink{}
	r := NewWriterRouter(price, stock, listing, resync, link,
		func(context.Context, string) (bool, error) { return true, nil }, func(context.Context, string) (bool, error) { return true, nil })
	tests := []struct {
		name  string
		typ   domain.ProtocolType
		after string
	}{
		{"price_update", domain.ProtocolTypePriceUpdate, `{}`}, {"stock_correct", domain.ProtocolTypeStockCorrect, `{}`},
		{"listing_pause", domain.ProtocolTypeListingPause, `{}`}, {"listing_edit", domain.ProtocolTypeListingEdit, `{"product_id":7,"attributes":[{"ID":"SELLER_SKU","ValueName":"7"}]}`},
		{"listing_resync", domain.ProtocolTypeListingResync, `{}`},
		{"approve_candidate", domain.ProtocolTypeLinkApply, `{"Action":"approve_candidate"}`},
		{"manual_resolve", domain.ProtocolTypeLinkApply, `{"Action":"manual_resolve"}`},
		{"reject_listing", domain.ProtocolTypeLinkApply, `{"Action":"reject_listing"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := r.ApplyItem(context.Background(), ports.WriteItem{ProtocolType: tt.typ, ListingID: "l", IdempotencyKey: tt.name, After: json.RawMessage(tt.after)})
			if err != nil || out.Failure != nil {
				t.Fatalf("out=%+v err=%v", out, err)
			}
			if tt.typ == domain.ProtocolTypeLinkApply && out.State != domain.ItemStateSkipped {
				t.Fatalf("state=%q", out.State)
			}
		})
	}
	if price.calls != 1 || stock.calls != 1 || listing.calls != 2 || resync.calls != 1 || link.calls != 3 {
		t.Fatalf("calls=%d/%d/%d/%d/%d", price.calls, stock.calls, listing.calls, resync.calls, link.calls)
	}
}

func TestWriterRouterRejectsSKUInvariantBeforeAdapter(t *testing.T) {
	listing := &captureWriter{}
	r := NewWriterRouter(nil, nil, listing, nil, nil, func(context.Context, string) (bool, error) { return true, nil }, nil)
	out, err := r.ApplyItem(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingEdit, After: json.RawMessage(`{"product_id":7,"attributes":[{"ID":"SELLER_SKU","ValueName":"8"}]}`)})
	if err != nil || out.Failure == nil || out.Failure.Code != domain.FailureCodeSKUInvariantViolation || listing.calls != 0 {
		t.Fatalf("out=%+v calls=%d err=%v", out, listing.calls, err)
	}
}

func TestWriterRouterListingCreateValidatesThenReturnsTypedNotEnabled(t *testing.T) {
	r := &WriterRouter{}
	_, err := r.ApplyItem(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingCreate, After: json.RawMessage(`{}`)})
	var gate *GateError
	if !errors.As(err, &gate) || gate.Code != domain.FailureCodeProviderValidation {
		t.Fatalf("err=%v", err)
	}
	out, err := r.ApplyItem(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingCreate, After: json.RawMessage(`{"product_id":7,"listing_type_code":"gold","price":10,"category_id":"c","attributes":[]}`)})
	if err != nil || out.Failure == nil || out.Failure.Code != domain.FailureCodeTypeNotEnabled {
		t.Fatalf("out=%+v err=%v", out, err)
	}
}

func TestWriterRouterProtocolRejectsRemainTypedErrors(t *testing.T) {
	now := time.Now()
	r := &WriterRouter{}
	err := r.ValidateProtocol(ProtocolWriteValidation{Actor: "", Execute: boolPtr(true), SourceAsOf: &now, Now: now, MaxAge: time.Minute})
	var gate *GateError
	if !errors.As(err, &gate) || gate.Code != FailureCodeActorRequired {
		t.Fatalf("err=%v", err)
	}
}
