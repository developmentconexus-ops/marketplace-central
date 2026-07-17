package application

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type readFake struct {
	protocols ProtocolPage
	protocol  ports.Protocol
	found     bool
	items     ItemPage
	itemCalls int
}

func (f *readFake) ListProtocols(context.Context, ProtocolQuery) (ProtocolPage, error) {
	return f.protocols, nil
}
func (f *readFake) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	return f.protocol, f.found, nil
}
func (f *readFake) ListItems(context.Context, ItemQuery) (ItemPage, error) {
	f.itemCalls++
	return f.items, nil
}

func TestReadServiceReturnsRepositoryPages(t *testing.T) {
	wantProtocols := ProtocolPage{Items: []ports.Protocol{{ProtocolID: "MP-2"}}}
	wantItems := ItemPage{Items: []ports.MutationItem{{Seq: 1}}}
	fake := &readFake{protocols: wantProtocols, protocol: ports.Protocol{ProtocolID: "MP-2"}, found: true, items: wantItems}
	service := NewReadService(fake)
	gotProtocols, err := service.ListProtocols(context.Background(), ProtocolQuery{InstallationID: "inst", Limit: 20})
	if err != nil || !reflect.DeepEqual(gotProtocols, wantProtocols) {
		t.Fatalf("page=%+v err=%v", gotProtocols, err)
	}
	got, err := service.GetProtocol(context.Background(), "MP-2")
	if err != nil || got.ProtocolID != "MP-2" {
		t.Fatalf("protocol=%+v err=%v", got, err)
	}
	gotItems, err := service.ListItems(context.Background(), ItemQuery{ProtocolID: "MP-2", Limit: 20})
	if err != nil || !reflect.DeepEqual(gotItems, wantItems) {
		t.Fatalf("page=%+v err=%v", gotItems, err)
	}
}

func TestReadServiceUnknownProtocolStopsItemRead(t *testing.T) {
	fake := &readFake{}
	_, err := NewReadService(fake).ListItems(context.Background(), ItemQuery{ProtocolID: "MP-missing", Limit: 20})
	if !errors.Is(err, ErrProtocolNotFound) || fake.itemCalls != 0 {
		t.Fatalf("error=%v calls=%d", err, fake.itemCalls)
	}
}
