package postgres

import (
	"context"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

func (r *Repository) CancelProtocol(ctx context.Context, protocolID string, cancelledAt time.Time) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE mutation_protocols
		SET state=$3, finished_at=$4
		WHERE tenant_id=$1 AND protocol_id=$2 AND state IN ('draft','previewed')`,
		r.tenantID, protocolID, domain.ProtocolStateCancelled, cancelledAt.UTC())
	if err != nil {
		return fmt.Errorf("cancel mutation protocol: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrInvalidStateTransition
	}
	return nil
}
