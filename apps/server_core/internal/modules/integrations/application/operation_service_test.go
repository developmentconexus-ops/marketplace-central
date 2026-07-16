package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

func TestOperationServiceListByInstallationReturnsRuns(t *testing.T) {
	t.Parallel()

	store := &stubOperationRunListStore{
		runs: []domain.OperationRun{
			{
				OperationRunID: "run_001",
				InstallationID: "inst_001",
			},
		},
	}
	svc := NewOperationService(store, "tenant-default")

	runs, err := svc.ListByInstallation(context.Background(), "inst_001")
	if err != nil {
		t.Fatalf("ListByInstallation() error = %v", err)
	}

	if len(runs) != 1 {
		t.Fatalf("runs = %d, want 1", len(runs))
	}
	if got, want := runs[0].OperationRunID, "run_001"; got != want {
		t.Fatalf("operation_run_id = %q, want %q", got, want)
	}
	if got, want := runs[0].InstallationID, "inst_001"; got != want {
		t.Fatalf("installation_id = %q, want %q", got, want)
	}
	if got, want := store.lastInstallationID, "inst_001"; got != want {
		t.Fatalf("forwarded installation_id = %q, want %q", got, want)
	}
}

func TestOperationServiceListByInstallationRejectsEmptyInstallationID(t *testing.T) {
	t.Parallel()

	svc := NewOperationService(&stubOperationRunListStore{}, "tenant-default")

	_, err := svc.ListByInstallation(context.Background(), " ")
	if err == nil {
		t.Fatal("ListByInstallation() error = nil, want invalid input error")
	}
	if got, want := err.Error(), "INTEGRATIONS_OPERATION_INVALID"; got != want {
		t.Fatalf("ListByInstallation() error = %q, want %q", got, want)
	}
}

func TestOperationServiceListRunsDelegatesToReadStore(t *testing.T) {
	t.Parallel()

	startedAt := time.Unix(100, 0).UTC()
	store := &stubOperationRunListStore{readPage: ports.RunPage{Items: []ports.RunReadModel{{
		OperationRunID: "run_read_001",
		InstallationID: "inst_001",
		Module:         "listing_read",
		StartedAt:      &startedAt,
	}}}}
	svc := NewOperationService(store, "tenant-default")
	query := ports.RunListQuery{InstallationID: "inst_001", Limit: 10, Module: "listing_read"}

	page, err := svc.ListRuns(context.Background(), query)
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if !reflect.DeepEqual(page, store.readPage) {
		t.Fatalf("page = %#v, want %#v", page, store.readPage)
	}
	if !reflect.DeepEqual(store.lastRunQuery, query) {
		t.Fatalf("query = %#v, want %#v", store.lastRunQuery, query)
	}
}

func TestOperationServiceLatestRunsByModuleDelegatesToStore(t *testing.T) {
	t.Parallel()

	startedAt := time.Unix(300, 0).UTC()
	successfulAt := time.Unix(400, 0).UTC()
	want := []ports.LatestRunByModule{{
		OperationType:      "listing_read",
		LatestAttemptedAt:  &startedAt,
		LatestSuccessfulAt: &successfulAt,
	}}
	store := &stubOperationRunListStore{latestRuns: want}
	svc := NewOperationService(store, "tenant-default")

	got, err := svc.LatestRunsByModule(context.Background(), "inst_001")
	if err != nil {
		t.Fatalf("LatestRunsByModule() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("runs = %#v, want %#v", got, want)
	}
	if got, want := store.lastLatestRunsInstallationID, "inst_001"; got != want {
		t.Fatalf("forwarded installation_id = %q, want %q", got, want)
	}
}

func TestOperationServiceLatestRunsByModulePropagatesStoreError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("store unavailable")
	store := &stubOperationRunListStore{latestRunsErr: wantErr}
	svc := NewOperationService(store, "tenant-default")

	_, err := svc.LatestRunsByModule(context.Background(), "inst_001")
	if !errors.Is(err, wantErr) {
		t.Fatalf("LatestRunsByModule() error = %v, want %v", err, wantErr)
	}
}

func TestOperationServiceBeginExclusiveReturnsActiveRun(t *testing.T) {
	store := &stubOperationRunListStore{runs: []domain.OperationRun{{OperationRunID: "op_active", InstallationID: "inst", OperationType: "listings_refresh", Status: domain.OperationRunStatusRunning}}}
	svc := NewOperationService(store, "tenant")
	run, active, err := svc.BeginExclusive(context.Background(), RecordOperationInput{OperationRunID: "op_new", InstallationID: "inst", OperationType: "listings_refresh", Status: domain.OperationRunStatusQueued})
	if err != nil || !active || run.OperationRunID != "op_active" || len(store.runs) != 1 {
		t.Fatalf("run=%+v active=%v count=%d err=%v", run, active, len(store.runs), err)
	}
}

type stubOperationRunListStore struct {
	runs                         []domain.OperationRun
	lastInstallationID           string
	readPage                     ports.RunPage
	lastRunQuery                 ports.RunListQuery
	latestRuns                   []ports.LatestRunByModule
	latestRunsErr                error
	lastLatestRunsInstallationID string
}

func (s *stubOperationRunListStore) BeginExclusive(_ context.Context, run domain.OperationRun) (domain.OperationRun, bool, error) {
	for _, existing := range s.runs {
		if existing.InstallationID == run.InstallationID && existing.OperationType == run.OperationType && (existing.Status == domain.OperationRunStatusQueued || existing.Status == domain.OperationRunStatusRunning) {
			return existing, true, nil
		}
	}
	s.runs = append(s.runs, run)
	return run, false, nil
}

func (s *stubOperationRunListStore) SaveOperationRun(_ context.Context, run domain.OperationRun) error {
	s.runs = append(s.runs, run)
	return nil
}

func (s *stubOperationRunListStore) ListByInstallation(_ context.Context, installationID string) ([]domain.OperationRun, error) {
	s.lastInstallationID = installationID
	return append([]domain.OperationRun(nil), s.runs...), nil
}

func (s *stubOperationRunListStore) ListRuns(_ context.Context, query ports.RunListQuery) (ports.RunPage, error) {
	s.lastRunQuery = query
	return s.readPage, nil
}

func (s *stubOperationRunListStore) LatestRunsByModule(_ context.Context, installationID string) ([]ports.LatestRunByModule, error) {
	s.lastLatestRunsInstallationID = installationID
	return s.latestRuns, s.latestRunsErr
}

func TestOperationServiceRecordMapsRichFields(t *testing.T) {
	t.Parallel()

	startedAt := time.Unix(100, 0).UTC()
	completedAt := time.Unix(200, 0).UTC()
	store := &stubOperationRunListStore{}
	svc := NewOperationService(store, "tenant-default")

	run, err := svc.Record(context.Background(), RecordOperationInput{
		OperationRunID:      "run_002",
		InstallationID:      "inst_002",
		OperationType:       "pricing_fee_sync",
		Status:              domain.OperationRunStatusRunning,
		ResultCode:          "INTEGRATIONS_OPERATION_RUNNING",
		FailureCode:         "INTEGRATIONS_OPERATION_TIMEOUT",
		TranslatedErrorCode: "CONNECTORS_PROVIDER_TIMEOUT",
		AttemptCount:        3,
		ActorType:           "system",
		ActorID:             "scheduler",
		ProviderEvidence: map[string]any{
			"provider_status": "timeout",
		},
		DurationMs:  3210,
		StartedAt:   &startedAt,
		CompletedAt: &completedAt,
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	if got, want := run.FailureCode, "INTEGRATIONS_OPERATION_TIMEOUT"; got != want {
		t.Fatalf("failure_code = %q, want %q", got, want)
	}
	if got, want := run.ActorType, "system"; got != want {
		t.Fatalf("actor_type = %q, want %q", got, want)
	}
	if got, want := run.ActorID, "scheduler"; got != want {
		t.Fatalf("actor_id = %q, want %q", got, want)
	}
	if got, want := run.TranslatedErrorCode, "CONNECTORS_PROVIDER_TIMEOUT"; got != want {
		t.Fatalf("translated_error_code = %q, want %q", got, want)
	}
	if got, want := run.DurationMs, int64(3210); got != want {
		t.Fatalf("duration_ms = %d, want %d", got, want)
	}
	if got := run.ProviderEvidence["provider_status"]; got != "timeout" {
		t.Fatalf("provider_evidence = %#v", run.ProviderEvidence)
	}
	if run.StartedAt == nil || !run.StartedAt.Equal(startedAt) {
		t.Fatalf("started_at = %v, want %v", run.StartedAt, startedAt)
	}
	if run.CompletedAt == nil || !run.CompletedAt.Equal(completedAt) {
		t.Fatalf("completed_at = %v, want %v", run.CompletedAt, completedAt)
	}
	if len(store.runs) != 1 {
		t.Fatalf("saved runs = %d, want 1", len(store.runs))
	}
	if got, want := store.runs[0], run; !reflect.DeepEqual(got, want) {
		t.Fatalf("saved run = %#v, want %#v", got, want)
	}
}
