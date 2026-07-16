package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type refreshGatewayStub struct {
	account     ports.InstallationAccount
	found       bool
	active      string
	records     []ports.RefreshRun
	recordErrAt int
}

func (s *refreshGatewayStub) InstallationAccount(context.Context, string) (ports.InstallationAccount, bool, error) {
	return s.account, s.found, nil
}
func (s *refreshGatewayStub) BeginRefresh(_ context.Context, run ports.RefreshRun) (ports.RefreshRun, bool, error) {
	if s.active != "" {
		return ports.RefreshRun{OperationRunID: s.active}, true, nil
	}
	s.records = append(s.records, run)
	return run, false, nil
}
func (s *refreshGatewayStub) RecordRefresh(_ context.Context, run ports.RefreshRun) error {
	s.records = append(s.records, run)
	if s.recordErrAt > 0 && len(s.records) == s.recordErrAt {
		return errors.New("persist")
	}
	return nil
}

type pullerStub struct {
	err      error
	contexts []context.Context
}

func (s *pullerStub) Pull(ctx context.Context, _ ports.InstallationAccount) error {
	s.contexts = append(s.contexts, ctx)
	return s.err
}

func TestRefreshServiceValidatesBeforeLookup(t *testing.T) {
	g := &refreshGatewayStub{}
	_, err := NewRefreshService(g, &pullerStub{}, func(func()) {}, time.Now, nil, nil).Start(context.Background(), " ")
	if !errors.Is(err, ErrInstallationIDRequired) {
		t.Fatalf("error = %v", err)
	}
}

func TestRefreshServiceUnknownInstallation(t *testing.T) {
	_, err := NewRefreshService(&refreshGatewayStub{}, &pullerStub{}, func(func()) {}, time.Now, nil, nil).Start(context.Background(), "missing")
	if !errors.Is(err, ErrInstallationNotFound) {
		t.Fatalf("error = %v", err)
	}
}

func TestRefreshServiceLifecycleAndInjectedRunner(t *testing.T) {
	g := &refreshGatewayStub{found: true, account: ports.InstallationAccount{InstallationID: "inst", TenantID: "tenant"}}
	p := &pullerStub{}
	var task func()
	s := NewRefreshService(g, p, func(f func()) { task = f }, func() time.Time { return time.Unix(100, 0) }, func(err error) { t.Error(err) }, nil)
	id, err := s.Start(context.Background(), "inst")
	if err != nil || id == "" || task == nil {
		t.Fatalf("id=%q task=%v err=%v", id, task != nil, err)
	}
	task()
	if got := []string{g.records[0].Status, g.records[1].Status, g.records[2].Status}; got[0] != "queued" || got[1] != "running" || got[2] != "succeeded" {
		t.Fatalf("statuses=%v", got)
	}
}

func TestRefreshServiceFailureTranslationAndTerminalPersistenceReporting(t *testing.T) {
	g := &refreshGatewayStub{found: true, account: ports.InstallationAccount{InstallationID: "inst"}, recordErrAt: 3}
	p := &pullerStub{err: errors.New("provider body must not persist")}
	var task func()
	var reported error
	s := NewRefreshService(g, p, func(f func()) { task = f }, time.Now, func(err error) { reported = err }, func(error) string { return "CONNECTORS_PROVIDER_TRANSIENT" })
	if _, err := s.Start(context.Background(), "inst"); err != nil {
		t.Fatal(err)
	}
	task()
	failed := g.records[2]
	if failed.Status != "failed" || failed.TranslatedErrorCode != "CONNECTORS_PROVIDER_TRANSIENT" {
		t.Fatalf("failed=%+v", failed)
	}
	if reported == nil {
		t.Fatal("terminal persistence failure was not reported")
	}
}

func TestRefreshServiceReturnsActiveRun(t *testing.T) {
	s := NewRefreshService(&refreshGatewayStub{found: true, account: ports.InstallationAccount{InstallationID: "inst"}, active: "op_active"}, &pullerStub{}, func(func()) {}, time.Now, nil, nil)
	_, err := s.Start(context.Background(), "inst")
	var active *RefreshInProgressError
	if !errors.As(err, &active) || active.OperationRunID != "op_active" {
		t.Fatalf("error=%v", err)
	}
}
