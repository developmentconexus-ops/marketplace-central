package application

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

var protocolPattern = regexp.MustCompile(`^#([0-9]+)-E$`)

// SyncEnqueuer enqueues accepted product codes into the market sync cursor for
// every installation of the tenant and returns the installations it touched.
type SyncEnqueuer interface {
	EnqueueMarketProducts(ctx context.Context, productCodes []string) (installationIDs []string, err error)
}

// LinkCandidateGenerator generates link candidates for one installation.
type LinkCandidateGenerator interface {
	GenerateLinkCandidates(ctx context.Context, installationID string) error
}

type ImportService struct {
	parser    ports.Parser
	repo      ports.ImportRepository
	tenantID  string
	now       func() time.Time
	newID     func() (domain.ImportID, error)
	enqueuer  SyncEnqueuer
	generator LinkCandidateGenerator
}

type ImportServiceOption func(*ImportService)

func WithClock(now func() time.Time) ImportServiceOption {
	return func(service *ImportService) {
		service.now = now
	}
}

func WithIDGenerator(newID func() (domain.ImportID, error)) ImportServiceOption {
	return func(service *ImportService) {
		service.newID = newID
	}
}

func WithSyncEnqueuer(enqueuer SyncEnqueuer) ImportServiceOption {
	return func(service *ImportService) {
		service.enqueuer = enqueuer
	}
}

func WithLinkCandidateGenerator(generator LinkCandidateGenerator) ImportServiceOption {
	return func(service *ImportService) {
		service.generator = generator
	}
}

func NewImportService(parser ports.Parser, repo ports.ImportRepository, tenantID string, opts ...ImportServiceOption) *ImportService {
	service := &ImportService{
		parser:   parser,
		repo:     repo,
		tenantID: tenantID,
		now:      time.Now,
		newID:    newUUIDv4,
	}
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func (s *ImportService) SetPostImportHooks(enqueuer SyncEnqueuer, generator LinkCandidateGenerator) {
	s.enqueuer = enqueuer
	s.generator = generator
}

func (s *ImportService) RunImport(ctx context.Context, source io.Reader) (domain.ImportReport, error) {
	return s.runImport(ctx, source, domain.SourceXLSX, s.parser.Parse, domain.ValidateRows)
}

// RunImportLenient runs the client-catalog import path: CUSTO and
// ESTOQUE_FISICO may be absent from the workbook and are left as
// honest-unknown (nil/empty) rather than rejected (ADR-17). The strict path
// (RunImport) is untouched.
func (s *ImportService) RunImportLenient(ctx context.Context, source io.Reader) (domain.ImportReport, error) {
	return s.runImport(ctx, source, domain.SourceCatalogoCliente, s.parser.ParseLenient, domain.ValidateRowsLenient)
}

func (s *ImportService) runImport(
	ctx context.Context,
	source io.Reader,
	importSource domain.ImportSource,
	parse func(context.Context, io.Reader) ([]domain.NormalizedRow, []domain.Issue, error),
	validate func([]domain.NormalizedRow) ([]domain.NormalizedRow, []domain.Issue, int),
) (domain.ImportReport, error) {
	buf, err := io.ReadAll(source)
	if err != nil {
		return domain.ImportReport{}, fmt.Errorf("read ERP import source: %w", err)
	}

	digest := sha256.Sum256(buf)
	fileSHA256 := domain.FileSHA256(hex.EncodeToString(digest[:]))
	rows, fileIssues, err := parse(ctx, bytes.NewReader(buf))
	if err != nil {
		var fileErr *ports.FileError
		if errors.As(err, &fileErr) {
			return domain.ImportReport{}, fileErr
		}
		return domain.ImportReport{}, fmt.Errorf("parse ERP import: %w", err)
	}

	accepted, rowIssues, _ := validate(rows)
	// File-level parse issues (e.g. a skipped worksheet) lead the per-row issues.
	issues := append(fileIssues, rowIssues...)
	status := domain.ImportStatusCompleted
	if len(accepted) == 0 {
		status = domain.ImportStatusRejected
	}

	protocol, err := s.nextProtocol(ctx)
	if err != nil {
		return domain.ImportReport{}, err
	}
	id, err := s.newID()
	if err != nil {
		return domain.ImportReport{}, fmt.Errorf("generate ERP import ID: %w", err)
	}

	snapshot := domain.ImportSnapshot{
		ID:           id,
		Protocol:     protocol,
		FileSHA256:   fileSHA256,
		Source:       importSource,
		ImportedAt:   s.now().UTC(),
		Status:       status,
		AcceptedRows: accepted,
		Issues:       issues,
	}
	if err := s.repo.PersistSnapshotAtomically(ctx, s.tenantID, snapshot); err != nil {
		return domain.ImportReport{}, fmt.Errorf("persist ERP import: %w", err)
	}
	if status == domain.ImportStatusCompleted && s.enqueuer != nil && s.generator != nil {
		codes := make([]string, 0, len(accepted))
		for _, row := range accepted {
			codes = append(codes, row.Codprod)
		}

		installationIDs, err := s.enqueuer.EnqueueMarketProducts(ctx, codes)
		if err != nil {
			return domain.ImportReport{}, fmt.Errorf("enqueue market products after ERP import: %w", err)
		}
		for _, installationID := range installationIDs {
			// Candidate generation is a consequence of the import, not part of
			// it: the snapshot is already persisted and correct. A matcher that
			// fails here must be visible in the log and re-runnable (the
			// generation endpoint stays live), never a reason to report an
			// import that landed as failed (M05-C11).
			if err := s.generator.GenerateLinkCandidates(ctx, installationID); err != nil {
				slog.ErrorContext(ctx, "link candidate generation after ERP import failed",
					"installation_id", installationID,
					"protocol", string(snapshot.Protocol),
					"error", err)
			}
		}
	}
	return s.repo.GetImport(ctx, s.tenantID, id)
}

func (s *ImportService) nextProtocol(ctx context.Context) (domain.Protocol, error) {
	imports, err := s.repo.ListImports(ctx, s.tenantID)
	if err != nil {
		return "", fmt.Errorf("list ERP imports for protocol allocation: %w", err)
	}

	maxNumber := 0
	for _, report := range imports {
		matches := protocolPattern.FindStringSubmatch(string(report.Protocol))
		if matches == nil {
			continue
		}
		number, err := strconv.Atoi(matches[1])
		if err == nil && number > maxNumber {
			maxNumber = number
		}
	}

	// Concurrent same-tenant imports fail under the repository lock; the protocol UNIQUE constraint is the residual fail-loud backstop.
	return domain.Protocol(fmt.Sprintf("#%03d-E", maxNumber+1)), nil
}

func newUUIDv4() (domain.ImportID, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	value[6] = value[6]&0x0f | 0x40
	value[8] = value[8]&0x3f | 0x80
	return domain.ImportID(fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])), nil
}
