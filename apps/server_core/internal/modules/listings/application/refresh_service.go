package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type refreshPuller interface {
	Pull(context.Context, ports.InstallationAccount) error
}
type BackgroundRunner func(func())
type RefreshErrorReporter func(error)
type ErrorCodeTranslator func(error) string

type RefreshInProgressError struct{ OperationRunID string }

func (e *RefreshInProgressError) Error() string { return "listings refresh already in progress" }

type RefreshService struct {
	gateway   ports.RefreshIntegrationGateway
	ingestion refreshPuller
	run       BackgroundRunner
	now       func() time.Time
	report    RefreshErrorReporter
	translate ErrorCodeTranslator
}

func NewRefreshService(gateway ports.RefreshIntegrationGateway, ingestion refreshPuller, run BackgroundRunner, now func() time.Time, report RefreshErrorReporter, translate ErrorCodeTranslator) *RefreshService {
	return &RefreshService{gateway: gateway, ingestion: ingestion, run: run, now: now, report: report, translate: translate}
}

func (s *RefreshService) Start(ctx context.Context, installationID string) (string, error) {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return "", ErrInstallationIDRequired
	}
	account, found, err := s.gateway.InstallationAccount(ctx, installationID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrInstallationNotFound
	}
	id, err := newOperationRunID()
	if err != nil {
		return "", err
	}
	queued, active, err := s.gateway.BeginRefresh(ctx, ports.RefreshRun{OperationRunID: id, InstallationID: installationID, Status: "queued"})
	if err != nil {
		return "", err
	}
	if active {
		return "", &RefreshInProgressError{OperationRunID: queued.OperationRunID}
	}
	if s.run == nil || s.ingestion == nil || s.now == nil {
		return "", errors.New("listing refresh is not configured")
	}
	s.run(func() { s.execute(context.Background(), account, queued.OperationRunID) })
	return queued.OperationRunID, nil
}

func (s *RefreshService) execute(ctx context.Context, account ports.InstallationAccount, id string) {
	started := s.now().UTC()
	if err := s.gateway.RecordRefresh(ctx, ports.RefreshRun{OperationRunID: id, InstallationID: account.InstallationID, Status: "running", StartedAt: &started}); err != nil {
		s.reportError(err)
		return
	}
	err := s.ingestion.Pull(ctx, account)
	completed := s.now().UTC()
	run := ports.RefreshRun{OperationRunID: id, InstallationID: account.InstallationID, Status: "succeeded", StartedAt: &started, CompletedAt: &completed}
	if err != nil {
		run.Status = "failed"
		if s.translate != nil {
			run.TranslatedErrorCode = s.translate(err)
		}
	}
	if persistErr := s.gateway.RecordRefresh(ctx, run); persistErr != nil {
		s.reportError(persistErr)
	}
}

func (s *RefreshService) reportError(err error) {
	if s.report != nil {
		s.report(err)
	}
}

func newOperationRunID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return "op_" + hex.EncodeToString(raw), nil
}
