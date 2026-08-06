package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/channelfees/domain"
	"marketplace-central/apps/server_core/internal/modules/channelfees/ports"
)

// Reader is the Postgres-backed ports.ChannelFeeReader.
type Reader struct {
	pool *pgxpool.Pool
}

// NewReader builds a Reader over pool.
func NewReader(pool *pgxpool.Pool) *Reader {
	return &Reader{pool: pool}
}

var _ ports.ChannelFeeReader = (*Reader)(nil)

// resolveLadderSQL picks the highest-layer row among the allowed layer set
// for one (tenant, provider, installation, listing subject, fee_kind).
// ORDER BY layer DESC gives layer 2 priority over layer 1 by construction —
// the ladder order lives in the SQL, not in a second Go-side comparison.
const resolveLadderSQL = `
SELECT layer, value_type, value, currency, detail, origem, coletado_em
FROM channel_fees
WHERE tenant_id = $1 AND provider = $2 AND installation_id = $3
  AND subject_type = $4 AND subject_id = $5
  AND fee_kind = $6 AND layer = ANY($7)
ORDER BY layer DESC
LIMIT 1`

// ResolveListingFees resolves both fee_kinds for one listing, ledger-only
// (IC-01 Enums And Statuses). Commission tries layer 2 then layer 1 (layer 1
// / category is not produced by any producer in this mission — MIS-008 —
// so this branch is a currently-unreachable placeholder for that future
// producer, kept here because the ladder order is decided and pinned by the
// contract, not by this feature). Freight tries layer 2 only: freight has no
// layer 1 (IC-01 — "frete NÃO tem camada 1"). Layer 3 is never queried by
// either ladder (ADR-009).
func (r *Reader) ResolveListingFees(ctx context.Context, tenantID, provider, installationID, providerListingID string) (domain.ResolvedListingFees, error) {
	if tenantID == "" {
		return domain.ResolvedListingFees{}, errors.New("channelfees: empty tenantID")
	}

	commission, err := r.resolveOne(ctx, tenantID, provider, installationID, providerListingID,
		domain.FeeKindCommission, []int16{domain.LayerListing, domain.LayerCategory})
	if err != nil {
		return domain.ResolvedListingFees{}, fmt.Errorf("channelfees: resolve commission: %w", err)
	}

	freight, err := r.resolveOne(ctx, tenantID, provider, installationID, providerListingID,
		domain.FeeKindFreight, []int16{domain.LayerListing})
	if err != nil {
		return domain.ResolvedListingFees{}, fmt.Errorf("channelfees: resolve freight: %w", err)
	}

	return domain.ResolvedListingFees{Commission: commission, Freight: freight}, nil
}

func (r *Reader) resolveOne(ctx context.Context, tenantID, provider, installationID, providerListingID, feeKind string, allowedLayers []int16) (domain.ResolvedFee, error) {
	row := r.pool.QueryRow(ctx, resolveLadderSQL,
		tenantID, provider, installationID, domain.SubjectTypeListing, providerListingID,
		feeKind, allowedLayers)

	var resolved domain.ResolvedFee
	var layer int16
	err := row.Scan(&layer, &resolved.ValueType, &resolved.Value, &resolved.Currency,
		&resolved.Detail, &resolved.Origem, &resolved.ColetadoEm)
	if errors.Is(err, pgx.ErrNoRows) {
		// Honest typed-absent: no ledger row at any allowed layer. Never a
		// zero amount, never a default (IC-01 Enums And Statuses).
		return domain.ResolvedFee{Found: false}, nil
	}
	if err != nil {
		return domain.ResolvedFee{}, err
	}
	resolved.Found = true
	resolved.Layer = int(layer)
	return resolved, nil
}
