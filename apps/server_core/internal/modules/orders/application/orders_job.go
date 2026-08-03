package application

import (
	"context"
	"encoding/json"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// OrdersCursor é o estado persistido em sync_state.cursor para a entidade
// orders. O campo Phase usa o vocabulário do ADR-07 (backfill|incremental|sweep)
// porque o scheduler lê exatamente esse campo para decidir se o ciclo avança
// last_incremental_at.
type OrdersCursor struct {
	Phase             string     `json:"phase"`
	LastUpdatedAt     *time.Time `json:"last_updated_at,omitempty"`
	Offset            int        `json:"offset"`
	RunStartedAt      *time.Time `json:"run_started_at,omitempty"`
	LastRunEnumerated int        `json:"last_run_enumerated"`
	LastRunImported   int        `json:"last_run_imported"`
	LastRunSkipped    int        `json:"last_run_skipped"`
}

const (
	phaseBackfill    = "backfill"
	phaseIncremental = "incremental"
)

// OrdersImporter é o pedaço de ImportService que o job consome. Declarado aqui
// (e não no pacote do consumidor) porque o job é quem define o contrato de que
// precisa.
type OrdersImporter interface {
	Import(ctx context.Context, input ImportOrdersInput) (domain.ImportResult, error)
}

// NewOrdersJob devolve o corpo de sync para a entidade orders.
//
// Duas invariantes valem mais que qualquer otimização aqui:
//
//  1. Página cheia não avança a marca d'água. Uma página cheia diz "a janela
//     ainda tem mais"; avançar a marca aí pula pedidos em silêncio, e o
//     silêncio é o defeito caro — a tela continuaria verde.
//  2. Marca d'água só recebe instante MEDIDO. Provider sem date_last_updated é
//     desconhecido; desconhecido não vira now() nem época zero (ADR-17). O
//     custo de não avançar é reprocessar uma janela, e reprocessar é barato
//     porque o ingest é upsert idempotente (provado em
//     tests/integration/orders_reingest_test.go).
//
// A sobreposição (overlap) recua a janela para sobreviver a skew de relógio
// entre nós e o provider, pelo mesmo motivo: reprocessar é barato, perder não é.
func NewOrdersJob(importer OrdersImporter, installationID string, pageSize int, overlap time.Duration, now func() time.Time) syncapp.JobFunc {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	return func(ctx context.Context, raw json.RawMessage) (json.RawMessage, error) {
		cursor := parseOrdersCursor(raw, now)

		input := ImportOrdersInput{
			InstallationID: installationID,
			Limit:          pageSize,
			Offset:         cursor.Offset,
		}
		if cursor.Phase == phaseIncremental && cursor.LastUpdatedAt != nil {
			from := cursor.LastUpdatedAt.Add(-overlap)
			input.UpdatedAfter = &from
		}

		result, err := importer.Import(ctx, input)
		if err != nil {
			// Cursor de volta INALTERADO: nil apagaria o estado, e uma marca
			// d'água apagada faz o próximo ciclo varrer tudo de novo achando
			// que é a primeira vez.
			return raw, err
		}

		next := cursor
		next.LastRunEnumerated = result.EnumeratedCount
		next.LastRunImported = result.ImportedCount
		next.LastRunSkipped = result.SkippedCount

		if result.EnumeratedCount >= pageSize {
			next.Offset = cursor.Offset + pageSize
			return marshalCursor(next, raw)
		}

		next.Offset = 0
		switch {
		case result.MaxProviderUpdatedAt != nil:
			next.LastUpdatedAt = result.MaxProviderUpdatedAt
		case cursor.Phase == phaseBackfill:
			// Backfill drenado sem nenhuma data do provider: a marca é o
			// instante em que o backfill COMEÇOU — "enumeramos tudo até aqui" é
			// um fato medido, diferente de chutar now().
			next.LastUpdatedAt = cursor.RunStartedAt
		}
		next.Phase = phaseIncremental
		return marshalCursor(next, raw)
	}
}

// parseOrdersCursor é tolerante de propósito: cursor ausente, vazio ou ilegível
// significa "nunca rodou", e a resposta certa para isso é um backfill, não um
// erro que pinta a tela de vermelho por causa de um JSON.
func parseOrdersCursor(raw json.RawMessage, now func() time.Time) OrdersCursor {
	var c OrdersCursor
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &c); err == nil && (c.Phase == phaseBackfill || c.Phase == phaseIncremental) {
			if c.Phase == phaseBackfill && c.RunStartedAt == nil {
				started := now().UTC()
				c.RunStartedAt = &started
			}
			return c
		}
	}
	started := now().UTC()
	return OrdersCursor{Phase: phaseBackfill, RunStartedAt: &started}
}

// marshalCursor devolve o cursor anterior se a serialização falhar: o ciclo
// realmente teve sucesso, e reportar erro aqui mentiria sobre o que aconteceu
// com os pedidos.
func marshalCursor(c OrdersCursor, previous json.RawMessage) (json.RawMessage, error) {
	encoded, err := json.Marshal(c)
	if err != nil {
		return previous, nil
	}
	return encoded, nil
}
