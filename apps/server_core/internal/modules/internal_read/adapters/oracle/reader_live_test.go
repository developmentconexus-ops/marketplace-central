//go:build cgo

package oracle

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestOracleLiveSmoke(t *testing.T) {
	if os.Getenv("MPC_ORACLE_LIVE_TEST") != "1" {
		t.Skip("set MPC_ORACLE_LIVE_TEST=1 to run live Oracle validation")
	}

	cfg, err := LoadConfigFromEnv(os.Getenv)
	if err != nil {
		t.Fatalf("load live Oracle config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	db, err := OpenDB(ctx, cfg)
	if err != nil {
		t.Fatalf("open live Oracle db: %v", err)
	}
	defer db.Close()

	reader := NewReader(db)

	t.Run("product lookup", func(t *testing.T) {
		input := discoverProductLookupInput(t, ctx, db)
		got, err := reader.FindProductsForLinking(ctx, input)
		if err != nil {
			t.Fatalf("find products for linking: %v", err)
		}
		if len(got) == 0 {
			t.Fatal("expected at least one product candidate")
		}
		if got[0].ProductID <= 0 {
			t.Fatalf("expected positive product id, got %d", got[0].ProductID)
		}
		if got[0].Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
	})

	t.Run("sellable stock", func(t *testing.T) {
		productID := discoverStockProductID(t, ctx, db)
		got, err := reader.GetSellableStock(ctx, ports.SellableStockInput{
			ProductID: productID,
			Policy:    domain.DefaultSellableStockPolicy(),
		})
		if err != nil {
			t.Fatalf("get sellable stock: %v", err)
		}
		if got.ProductID != productID {
			t.Fatalf("expected product %d, got %d", productID, got.ProductID)
		}
		if got.Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
	})

	t.Run("current price", func(t *testing.T) {
		productID, tableID, localID, effectiveAt := discoverCurrentPriceInput(t, ctx, db)
		got, err := reader.GetCurrentPrice(ctx, ports.CurrentPriceInput{
			ProductID: productID,
			Policy: domain.CurrentPricePolicy{
				PriceTableID: tableID,
				LocationID:   localID,
				EffectiveAt:  effectiveAt,
			},
		})
		if err != nil {
			t.Fatalf("get current price: %v", err)
		}
		if got.ProductID != productID {
			t.Fatalf("expected product %d, got %d", productID, got.ProductID)
		}
		if got.Amount == nil {
			t.Fatal("expected non-nil current price amount")
		}
		if got.Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
	})

	t.Run("cost as-of", func(t *testing.T) {
		productID, companyID, effectiveAt := discoverCostInput(t, ctx, db)
		got, err := reader.GetCostAsOf(ctx, ports.CostAsOfInput{
			ProductID: productID,
			Policy: domain.CostAsOfPolicy{
				CompanyID:   companyID,
				EffectiveAt: effectiveAt,
				Basis:       domain.CostBasisCUSSEMICM,
			},
		})
		if err != nil {
			t.Fatalf("get cost as-of: %v", err)
		}
		if got.ProductID != productID {
			t.Fatalf("expected product %d, got %d", productID, got.ProductID)
		}
		if got.Amount == nil {
			t.Fatal("expected non-nil cost amount")
		}
		if got.Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
	})

	t.Run("sales history", func(t *testing.T) {
		input := discoverSalesHistoryInput(t, ctx, db)
		got, err := reader.GetSalesHistory(ctx, input)
		if err != nil {
			t.Fatalf("get sales history: %v", err)
		}
		if len(got.Entries) == 0 {
			t.Fatal("expected at least one sales history entry")
		}
		if got.Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
	})

	t.Run("tax inputs", func(t *testing.T) {
		productID, effectiveAt, incidenceCode := discoverTaxInput(t, ctx, db)
		got, err := reader.GetTaxInputs(ctx, ports.TaxInput{
			ProductID: productID,
			Policy: domain.TaxPolicy{
				EffectiveAt:   effectiveAt,
				IncidenceCode: incidenceCode,
			},
		})
		if err != nil {
			t.Fatalf("get tax inputs: %v", err)
		}
		if got.ProductID != productID {
			t.Fatalf("expected product %d, got %d", productID, got.ProductID)
		}
		if got.Source.FetchedAt.IsZero() {
			t.Fatal("expected source fetched timestamp")
		}
		if got.ICMSAmount == nil && got.IPIAmount == nil && got.PISAmount == nil && got.COFINSAmount == nil {
			t.Fatal("expected at least one tax amount")
		}
	})
}

func discoverProductLookupInput(t *testing.T, ctx context.Context, db *sql.DB) ports.FindProductsInput {
	t.Helper()

	var reference string
	err := db.QueryRowContext(ctx, `
SELECT p.REFERENCIA
FROM METALPRD.TGFPRO p
WHERE p.ATIVO = 'S'
  AND LENGTH(TRIM(p.REFERENCIA)) > 0
FETCH FIRST 1 ROW ONLY`).Scan(&reference)
	if err != nil {
		t.Fatalf("discover product lookup input: %v", err)
	}
	return ports.FindProductsInput{EAN: &reference}
}

func discoverStockProductID(t *testing.T, ctx context.Context, db *sql.DB) int {
	t.Helper()

	var productID int
	err := db.QueryRowContext(ctx, `
SELECT est.CODPROD
FROM METALPRD.TGFEST est
WHERE est.CODPARC = 0
  AND est.CODEMP IN (1, 2)
  AND est.CODLOCAL = 10101
FETCH FIRST 1 ROW ONLY`).Scan(&productID)
	if err != nil {
		t.Fatalf("discover stock product id: %v", err)
	}
	return productID
}

func discoverCurrentPriceInput(t *testing.T, ctx context.Context, db *sql.DB) (int, int, int, time.Time) {
	t.Helper()

	var (
		productID   int
		tableID     int
		localID     int
		effectiveAt time.Time
	)
	err := db.QueryRowContext(ctx, `
SELECT p.CODPROD, p.CODTAB, p.CODLOCAL, p.VIGENTE_DE
FROM LEANDRO.VW_PRECO_TABELA p
WHERE p.VIGENTE_HOJE = 'S'
FETCH FIRST 1 ROW ONLY`).Scan(&productID, &tableID, &localID, &effectiveAt)
	if err != nil {
		t.Fatalf("discover current price input: %v", err)
	}
	return productID, tableID, localID, effectiveAt
}

func discoverCostInput(t *testing.T, ctx context.Context, db *sql.DB) (int, int, time.Time) {
	t.Helper()

	var (
		productID   int
		companyID   int
		effectiveAt time.Time
	)
	err := db.QueryRowContext(ctx, `
SELECT c.CODPROD, c.CODEMP, c.DTATUAL
FROM METALPRD.TGFCUS c
WHERE c.CUSSEMICM IS NOT NULL
ORDER BY c.DTATUAL DESC
FETCH FIRST 1 ROW ONLY`).Scan(&productID, &companyID, &effectiveAt)
	if err != nil {
		t.Fatalf("discover cost input: %v", err)
	}
	return productID, companyID, effectiveAt
}

func discoverSalesHistoryInput(t *testing.T, ctx context.Context, db *sql.DB) ports.SalesHistoryInput {
	t.Helper()

	var (
		productID int
		start     time.Time
		end       time.Time
	)
	err := db.QueryRowContext(ctx, `
SELECT v.CODPROD, MIN(v.DTNEG), MAX(v.DTNEG)
FROM LEANDRO.VW_FAT_VENDA_ITEM v
GROUP BY v.CODPROD
FETCH FIRST 1 ROW ONLY`).Scan(&productID, &start, &end)
	if err != nil {
		t.Fatalf("discover sales history input: %v", err)
	}
	end = end.Add(24 * time.Hour)
	return ports.SalesHistoryInput{
		ProductID: &productID,
		Window: domain.SalesHistoryWindow{
			Start: start,
			End:   end,
		},
	}
}

func discoverTaxInput(t *testing.T, ctx context.Context, db *sql.DB) (int, time.Time, int) {
	t.Helper()

	var (
		productID     int
		effectiveAt   time.Time
		incidenceCode int
	)
	err := db.QueryRowContext(ctx, `
SELECT i.CODPROD, d.DTNEG, d.CODINC
FROM LEANDRO.VW_IMPOSTO_ITEM d
JOIN METALPRD.TGFITE i
  ON i.NUNOTA = d.NUNOTA
 AND i.SEQUENCIA = d.SEQUENCIA
WHERE d.VALOR IS NOT NULL
FETCH FIRST 1 ROW ONLY`).Scan(&productID, &effectiveAt, &incidenceCode)
	if err != nil {
		t.Fatalf("discover tax input: %v", err)
	}
	return productID, effectiveAt, incidenceCode
}
