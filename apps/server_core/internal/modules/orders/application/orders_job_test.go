package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

func fixedNow() time.Time {
	return time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
}

func ptr[T any](v T) *T {
	return &v
}

func mustUnmarshal(t *testing.T, raw json.RawMessage, v interface{}) {
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

func cursorJSON(t *testing.T, phase string, lastUpdatedAt *time.Time, offset int) json.RawMessage {
	c := application.OrdersCursor{
		Phase:         phase,
		LastUpdatedAt: lastUpdatedAt,
		Offset:        offset,
	}
	raw, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal cursor: %v", err)
	}
	return raw
}

type fakeImporter struct {
	calls   []application.ImportOrdersInput
	results []domain.ImportResult
	errs    []error
}

func (f *fakeImporter) Import(_ context.Context, in application.ImportOrdersInput) (domain.ImportResult, error) {
	f.calls = append(f.calls, in)
	i := len(f.calls) - 1
	var err error
	if i < len(f.errs) {
		err = f.errs[i]
	}
	var res domain.ImportResult
	if i < len(f.results) {
		res = f.results[i]
	}
	return res, err
}

// 1. Cursor ausente entra em backfill SEM janela.
func TestOrdersJobStartsInBackfillWithoutWindow(t *testing.T) {
	imp := &fakeImporter{results: []domain.ImportResult{{EnumeratedCount: 0}}}
	job := application.NewOrdersJob(imp, "inst-1", 50, 5*time.Minute, fixedNow)

	raw, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	if imp.calls[0].UpdatedAfter != nil {
		t.Fatalf("backfill nao pode ter janela; recebi %v", imp.calls[0].UpdatedAfter)
	}
	var c application.OrdersCursor
	mustUnmarshal(t, raw, &c)
	if c.Phase != "incremental" {
		t.Fatalf("pagina parcial (0 de 50) drena a janela e vira incremental; phase=%q", c.Phase)
	}
	if c.LastUpdatedAt == nil || !c.LastUpdatedAt.Equal(fixedNow()) {
		t.Fatalf("sem date_last_updated nenhum, a marca d'agua e o run_started_at MEDIDO; recebi %v", c.LastUpdatedAt)
	}
}

// 2. Página cheia NÃO avança a marca d'água — só o offset.
func TestOrdersJobFullPageAdvancesOffsetNotWatermark(t *testing.T) {
	previous := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC)
	imp := &fakeImporter{results: []domain.ImportResult{{
		EnumeratedCount:      50,
		ImportedCount:        50,
		MaxProviderUpdatedAt: &newer,
	}}}
	job := application.NewOrdersJob(imp, "inst-1", 50, 5*time.Minute, fixedNow)

	raw, err := job(context.Background(), cursorJSON(t, "incremental", &previous, 0))
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	var c application.OrdersCursor
	mustUnmarshal(t, raw, &c)
	if !c.LastUpdatedAt.Equal(previous) {
		t.Fatalf("pagina cheia significa janela NAO drenada; avancar a marca pula pedidos. quero %s, recebi %s", previous, c.LastUpdatedAt)
	}
	if c.Offset != 50 {
		t.Fatalf("offset: quero 50, recebi %d", c.Offset)
	}
}

// 3. Página parcial avança a marca d'água e zera o offset.
func TestOrdersJobPartialPageAdvancesWatermarkAndResetsOffset(t *testing.T) {
	previous := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC)
	imp := &fakeImporter{results: []domain.ImportResult{{
		EnumeratedCount:      7,
		ImportedCount:        7,
		MaxProviderUpdatedAt: &newer,
	}}}
	job := application.NewOrdersJob(imp, "inst-1", 50, 5*time.Minute, fixedNow)

	raw, err := job(context.Background(), cursorJSON(t, "incremental", &previous, 0))
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	var c application.OrdersCursor
	mustUnmarshal(t, raw, &c)
	if !c.LastUpdatedAt.Equal(newer) {
		t.Fatalf("pagina parcial significa janela drenada; avancar a marca. quero %s, recebi %s", newer, c.LastUpdatedAt)
	}
	if c.Offset != 0 {
		t.Fatalf("offset: quero 0, recebi %d", c.Offset)
	}
}

// 4. Nenhum date_last_updated em fase incremental NÃO move a marca d'água.
func TestOrdersJobKeepsWatermarkWhenProviderOmitsUpdatedAt(t *testing.T) {
	previous := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	imp := &fakeImporter{results: []domain.ImportResult{{
		EnumeratedCount:      3,
		ImportedCount:        3,
		MaxProviderUpdatedAt: nil,
	}}}
	job := application.NewOrdersJob(imp, "inst-1", 50, 5*time.Minute, fixedNow)

	raw, err := job(context.Background(), cursorJSON(t, "incremental", &previous, 0))
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	var c application.OrdersCursor
	mustUnmarshal(t, raw, &c)
	if !c.LastUpdatedAt.Equal(previous) {
		t.Fatalf("desconhecido nao vira now(): quero a marca anterior %s, recebi %s", previous, c.LastUpdatedAt)
	}
}

// 5. Erro do importador devolve o cursor RECEBIDO, byte a byte.
func TestOrdersJobErrorReturnsCursorUnchanged(t *testing.T) {
	in := cursorJSON(t, "incremental", ptr(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)), 100)
	imp := &fakeImporter{errs: []error{errors.New("token expirado")}}
	job := application.NewOrdersJob(imp, "inst-1", 50, 5*time.Minute, fixedNow)

	out, err := job(context.Background(), in)
	if err == nil {
		t.Fatalf("falha do provider tem que virar erro do ciclo — e' o que pinta a tela de vermelho")
	}
	if !bytes.Equal(out, in) {
		t.Fatalf("cursor no erro: devolver nil APAGA o estado. quero %s, recebi %s", in, out)
	}
}
