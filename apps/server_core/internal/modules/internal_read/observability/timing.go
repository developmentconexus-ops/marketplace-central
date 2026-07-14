package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var safeOracleCodePattern = regexp.MustCompile(`^oracle error code=(-?[0-9]+)$`)

type TimingReader struct {
	next          ports.Reader
	logger        *slog.Logger
	slowThreshold time.Duration
}

func NewTimingReader(next ports.Reader, logger *slog.Logger, slowThreshold time.Duration) *TimingReader {
	if logger == nil {
		logger = slog.Default()
	}
	if slowThreshold <= 0 {
		slowThreshold = defaultSlowQueryThreshold
	}
	return &TimingReader{next: next, logger: logger, slowThreshold: slowThreshold}
}

func (r *TimingReader) FindProductsForLinking(ctx context.Context, input ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	var result []domain.ProductCandidate
	err := r.observe("FindProductsForLinking", func() error {
		var err error
		result, err = r.next.FindProductsForLinking(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) GetSellableStock(ctx context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	var result domain.SellableStock
	err := r.observe("GetSellableStock", func() error {
		var err error
		result, err = r.next.GetSellableStock(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) GetCurrentPrice(ctx context.Context, input ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	var result domain.CurrentPrice
	err := r.observe("GetCurrentPrice", func() error {
		var err error
		result, err = r.next.GetCurrentPrice(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) GetCostAsOf(ctx context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	var result domain.CostAsOf
	err := r.observe("GetCostAsOf", func() error {
		var err error
		result, err = r.next.GetCostAsOf(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) GetSalesHistory(ctx context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	var result domain.SalesHistory
	err := r.observe("GetSalesHistory", func() error {
		var err error
		result, err = r.next.GetSalesHistory(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) GetTaxInputs(ctx context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	var result domain.TaxInputs
	err := r.observe("GetTaxInputs", func() error {
		var err error
		result, err = r.next.GetTaxInputs(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingReader) observe(method string, call func() error) error {
	started := time.Now()
	err := call()
	duration := time.Since(started)
	attrs := []any{"method", method, "duration_ms", duration.Milliseconds()}
	if err != nil {
		attrs = append(attrs, "oracle_code", safeOracleCode(err))
		r.logger.Error("oracle_read", attrs...)
		return err
	}
	attrs = append(attrs, "slow_query", duration > r.slowThreshold)
	r.logger.Info("oracle_read", attrs...)
	return nil
}

type TimingSankhyaLinkageReader struct {
	next          ports.SankhyaLinkageReader
	logger        *slog.Logger
	slowThreshold time.Duration
}

func NewTimingSankhyaLinkageReader(next ports.SankhyaLinkageReader, logger *slog.Logger, slowThreshold time.Duration) *TimingSankhyaLinkageReader {
	if logger == nil {
		logger = slog.Default()
	}
	if slowThreshold <= 0 {
		slowThreshold = defaultSlowQueryThreshold
	}
	return &TimingSankhyaLinkageReader{next: next, logger: logger, slowThreshold: slowThreshold}
}

func (r *TimingSankhyaLinkageReader) ValidateConfiguration(ctx context.Context) error {
	return r.observe("ValidateConfiguration", func() error { return r.next.ValidateConfiguration(ctx) })
}

func (r *TimingSankhyaLinkageReader) FindCandidates(ctx context.Context, input ports.SankhyaCandidateInput) ([]domain.SankhyaDocumentCandidate, error) {
	var result []domain.SankhyaDocumentCandidate
	err := r.observe("FindCandidates", func() error {
		var err error
		result, err = r.next.FindCandidates(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingSankhyaLinkageReader) ListDescendants(ctx context.Context, input ports.SankhyaDescendantInput) (domain.SankhyaLineage, error) {
	var result domain.SankhyaLineage
	err := r.observe("ListDescendants", func() error {
		var err error
		result, err = r.next.ListDescendants(ctx, input)
		return err
	})
	return result, err
}

func (r *TimingSankhyaLinkageReader) observe(method string, call func() error) error {
	started := time.Now()
	err := call()
	duration := time.Since(started)
	attrs := []any{"method", method, "duration_ms", duration.Milliseconds()}
	if err != nil {
		attrs = append(attrs, "oracle_code", safeOracleCode(err))
		r.logger.Error("oracle_read", attrs...)
		return err
	}
	attrs = append(attrs, "slow_query", duration > r.slowThreshold)
	r.logger.Info("oracle_read", attrs...)
	return nil
}

type oracleCoded interface {
	Code() int
}

func safeOracleCode(err error) string {
	var coded oracleCoded
	if errors.As(err, &coded) {
		return strconv.Itoa(coded.Code())
	}
	for current := err; current != nil; current = errors.Unwrap(current) {
		match := safeOracleCodePattern.FindStringSubmatch(current.Error())
		if len(match) == 2 {
			return match[1]
		}
	}
	return "unknown"
}

func formatSafeOracleCode(code int) error {
	return fmt.Errorf("oracle error code=%d", code)
}
