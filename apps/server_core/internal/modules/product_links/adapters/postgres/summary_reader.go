package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

var _ ports.SummaryReader = (*SummaryReader)(nil)

type SummaryReader struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewSummaryReader(pool *pgxpool.Pool, tenantID string) *SummaryReader {
	return &SummaryReader{pool: pool, tenantID: strings.TrimSpace(tenantID)}
}

func (r *SummaryReader) GetLinkageSummary(ctx context.Context, installationID string) (ports.LinkageSummary, error) {
	if r.pool == nil {
		return ports.LinkageSummary{}, errors.New("product links summary reader: database is not configured")
	}

	var summary ports.LinkageSummary
	err := r.pool.QueryRow(ctx, `
		WITH settled AS (
			-- A listing the operator already answered is not pending, whatever
			-- its candidate rows still say. Approving or rejecting settles
			-- product_links but never retires the candidate, and regeneration
			-- re-materializes the same candidate from the same anchors on every
			-- run — so a candidate-only count never goes down (round-3 B1).
			SELECT provider_item_id, provider_variation_id
			FROM product_links
			WHERE tenant_id = $1
			  AND installation_id = $2
			  AND state IN ($8, $9)
		), pending_links AS (
			SELECT provider_item_id, provider_variation_id
			FROM product_links
			WHERE tenant_id = $1
			  AND installation_id = $2
			  AND state IN ($3, $4)
			UNION
			SELECT c.provider_item_id, c.provider_variation_id
			FROM product_link_candidates c
			WHERE c.tenant_id = $1
			  AND c.installation_id = $2
			  AND c.state IN ($5, $6)
			  AND NOT EXISTS (
				  SELECT 1 FROM settled s
				  WHERE s.provider_item_id = c.provider_item_id
					AND s.provider_variation_id = c.provider_variation_id
			  )
			UNION
			-- A candidate awaiting CONFIRMAÇÃO is pending too. It resolved an
			-- anchor, so its STATE is exact_sku/exact_ean and the branches
			-- above cannot see it — but nothing is linked until the operator
			-- answers, and a counter that reads zero would report that work as
			-- done (D-121-2).
			SELECT c.provider_item_id, c.provider_variation_id
			FROM product_link_candidates c
			WHERE c.tenant_id = $1
			  AND c.installation_id = $2
			  AND c.match_status = $7
			  AND NOT EXISTS (
				  SELECT 1 FROM settled s
				  WHERE s.provider_item_id = c.provider_item_id
					AND s.provider_variation_id = c.provider_variation_id
			  )
		)
		SELECT
			(SELECT COUNT(*) FROM pending_links),
			(
				SELECT COUNT(*)
				FROM product_link_listing_snapshots
				WHERE tenant_id = $1
				  AND installation_id = $2
				  AND NULLIF(BTRIM(COALESCE(ean, '')), '') IS NULL
			)
	`, r.tenantID,
		strings.TrimSpace(installationID),
		domain.ProductLinkStateUnresolved,
		domain.ProductLinkStateConflict,
		domain.LinkCandidateStateUnresolved,
		domain.LinkCandidateStateConflict,
		domain.LinkCandidateMatchStatusConfirm,
		domain.ProductLinkStateResolved,
		domain.ProductLinkStateRejected,
	).Scan(&summary.PendingLinks, &summary.MissingGTIN)
	if err != nil {
		return ports.LinkageSummary{}, err
	}
	return summary, nil
}
