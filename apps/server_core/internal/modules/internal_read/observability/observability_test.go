package observability

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type fakeReader struct {
	delay time.Duration
	err   error
}

func (r fakeReader) FindProductsForLinking(context.Context, ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	return nil, r.result()
}
func (r fakeReader) GetSellableStock(context.Context, ports.SellableStockInput) (domain.SellableStock, error) {
	return domain.SellableStock{}, r.result()
}
func (r fakeReader) GetCurrentPrice(context.Context, ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	return domain.CurrentPrice{}, r.result()
}
func (r fakeReader) GetCostAsOf(context.Context, ports.CostAsOfInput) (domain.CostAsOf, error) {
	return domain.CostAsOf{}, r.result()
}
func (r fakeReader) GetSalesHistory(context.Context, ports.SalesHistoryInput) (domain.SalesHistory, error) {
	return domain.SalesHistory{}, r.result()
}
func (r fakeReader) GetTaxInputs(context.Context, ports.TaxInput) (domain.TaxInputs, error) {
	return domain.TaxInputs{}, r.result()
}
func (r fakeReader) result() error {
	if r.delay > 0 {
		time.Sleep(r.delay)
	}
	return r.err
}

type codedError struct{}

func (codedError) Error() string {
	return "ORA-00942: SELECT secret FROM credentials WHERE password='secret'"
}
func (codedError) Code() int { return 942 }

type fakePoolStatsSource struct{}

func (fakePoolStatsSource) Stats() sql.DBStats {
	return sql.DBStats{OpenConnections: 7, InUse: 3, WaitCount: 11}
}

func TestTimingReaderLogsLatencySlowFlagAndSanitizedError(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	ctx := context.Background()

	fast := NewTimingReader(fakeReader{}, logger, time.Second)
	if _, err := fast.GetSellableStock(ctx, ports.SellableStockInput{}); err != nil {
		t.Fatalf("fast fake reader: %v", err)
	}

	slow := NewTimingReader(fakeReader{delay: 1100 * time.Millisecond}, logger, time.Second)
	if _, err := slow.GetSellableStock(ctx, ports.SellableStockInput{}); err != nil {
		t.Fatalf("slow fake reader: %v", err)
	}

	_, err := NewTimingReader(fakeReader{err: codedError{}}, logger, time.Second).GetSellableStock(ctx, ports.SellableStockInput{})
	if err == nil {
		t.Fatal("expected fake reader error")
	}

	transcript := output.String()
	if got := strings.Count(transcript, "oracle_read"); got != 3 {
		t.Fatalf("oracle_read lines = %d, want 3; transcript=%q", got, transcript)
	}
	if !strings.Contains(transcript, "method=GetSellableStock") || !strings.Contains(transcript, "duration_ms=") {
		t.Fatalf("latency fields missing: %q", transcript)
	}
	if !strings.Contains(transcript, "slow_query=true") {
		t.Fatalf("slow-query flag missing: %q", transcript)
	}
	if !strings.Contains(transcript, "oracle_code=942") {
		t.Fatalf("numeric Oracle code missing: %q", transcript)
	}
	for _, forbidden := range []string{"SELECT", "credentials", "password", "secret", "ORA-00942"} {
		if strings.Contains(transcript, forbidden) {
			t.Fatalf("unsafe %q in log transcript: %q", forbidden, transcript)
		}
	}
	t.Logf("captured sanitized slog transcript: %s", strings.TrimSpace(transcript))
}

func TestPoolStatsLoopEmitsAndStops(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	loop := NewPoolStatsLoop(fakePoolStatsSource{}, logger, time.Hour)
	loop.Start(context.Background())
	loop.Stop()
	loop.Stop()

	transcript := output.String()
	if !strings.Contains(transcript, "pool_stats") {
		t.Fatalf("pool stats line missing: %q", transcript)
	}
	for _, field := range []string{"open=7", "in_use=3", "wait_count=11"} {
		if !strings.Contains(transcript, field) {
			t.Fatalf("pool stats field %q missing: %q", field, transcript)
		}
	}
	t.Logf("captured pool stats transcript: %s", strings.TrimSpace(transcript))
}

func TestSafeOracleCodeDoesNotExposeUnknownError(t *testing.T) {
	err := fmt.Errorf("driver text with SELECT and password: %w", errors.New("network failure"))
	if got := safeOracleCode(err); got != "unknown" {
		t.Fatalf("safeOracleCode = %q, want unknown", got)
	}
}
