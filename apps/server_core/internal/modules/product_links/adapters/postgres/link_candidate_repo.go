package postgres

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

var _ ports.LinkCandidateStore = (*LinkCandidateRepository)(nil)
var _ ports.ProductLinkWorkflowStore = (*LinkCandidateRepository)(nil)

type LinkCandidateRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewLinkCandidateRepository(pool *pgxpool.Pool, tenantID string) *LinkCandidateRepository {
	return &LinkCandidateRepository{pool: pool, tenantID: tenantID}
}

func (r *LinkCandidateRepository) ReplaceLinkCandidates(ctx context.Context, installationID string, identities []domain.ListingIdentity, candidates []domain.LinkCandidate) error {
	if r.pool == nil {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, identity := range identities {
		if _, err := tx.Exec(ctx, `
			DELETE FROM product_link_candidates
			WHERE tenant_id = $1
			  AND installation_id = $2
			  AND provider_item_id = $3
			  AND provider_variation_id = $4
		`, r.tenantID, installationID, identity.ProviderItemID, identity.ProviderVariationID); err != nil {
			return err
		}
	}

	for _, candidate := range candidates {
		reasons, err := marshalReasons(candidate.Reasons)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO product_link_candidates (
				tenant_id, candidate_id, installation_id, provider_code, provider_item_id, provider_variation_id,
				internal_product_id, internal_product_name, internal_reference_code, state, match_input, match_value,
				source_snapshot_fetched_at, confidence, confidence_band, match_status, reasons, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11, $12,
				$13, $14, $15, $16, $17, $18, $19
			)
			ON CONFLICT (tenant_id, candidate_id) DO UPDATE SET
				provider_code = EXCLUDED.provider_code,
				provider_item_id = EXCLUDED.provider_item_id,
				provider_variation_id = EXCLUDED.provider_variation_id,
				internal_product_id = EXCLUDED.internal_product_id,
				internal_product_name = EXCLUDED.internal_product_name,
				internal_reference_code = EXCLUDED.internal_reference_code,
				state = EXCLUDED.state,
				match_input = EXCLUDED.match_input,
				match_value = EXCLUDED.match_value,
				source_snapshot_fetched_at = EXCLUDED.source_snapshot_fetched_at,
				confidence = EXCLUDED.confidence,
				confidence_band = EXCLUDED.confidence_band,
				match_status = EXCLUDED.match_status,
				reasons = EXCLUDED.reasons,
				created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at
		`, r.tenantID,
			candidate.CandidateID,
			candidate.InstallationID,
			candidate.ProviderCode,
			candidate.ProviderItemID,
			candidate.ProviderVariationID,
			candidate.InternalProductID,
			candidate.InternalProductName,
			candidate.InternalReferenceCode,
			candidate.State,
			candidate.MatchInput,
			candidate.MatchValue,
			timestamptzArg(candidate.SourceSnapshotFetchedAt),
			candidate.Confidence,
			candidate.ConfidenceBand,
			candidate.MatchStatus,
			reasons,
			candidate.CreatedAt,
			candidate.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *LinkCandidateRepository) ListLinkCandidates(ctx context.Context, installationID string, limit int) ([]domain.LinkCandidate, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			candidate_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			internal_product_id, internal_product_name, internal_reference_code, state, match_input, match_value,
			source_snapshot_fetched_at, confidence, confidence_band, match_status, reasons, created_at, updated_at
		FROM product_link_candidates c
		WHERE c.tenant_id = $1
		  AND c.installation_id = $2
		  -- An anúncio the operator already decided leaves the queue. Candidates
		  -- outlive the decision (they are only rewritten by the next generation
		  -- run), so without this the row the operator just resolved sat in the
		  -- fila as if nothing had happened.
		  AND NOT EXISTS (
			SELECT 1 FROM product_links l
			WHERE l.tenant_id = c.tenant_id
			  AND l.installation_id = c.installation_id
			  AND l.provider_item_id = c.provider_item_id
			  AND l.provider_variation_id = c.provider_variation_id
		  )
		-- Actionable first. The queue is a work list, not a log: a run over a whole
		-- account produces mostly NO_CANDIDATE rows, and ordering by recency buried
		-- every real suggestion behind thousands of variations nobody can act on.
		-- Confidence already encodes the priority (ACCEPT > REVIEW > REJECT/conflict
		-- > no candidate), so ranking by it puts the decisions the operator can make
		-- on the first page.
		ORDER BY confidence DESC, updated_at DESC, candidate_id DESC
		LIMIT $3
	`, r.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := make([]domain.LinkCandidate, 0)
	for rows.Next() {
		var candidate domain.LinkCandidate
		var internalProductID pgtype.Int4
		var sourceSnapshotFetchedAt pgtype.Timestamptz
		var reasons []byte
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz

		if err := rows.Scan(
			&candidate.CandidateID,
			&candidate.InstallationID,
			&candidate.ProviderCode,
			&candidate.ProviderItemID,
			&candidate.ProviderVariationID,
			&internalProductID,
			&candidate.InternalProductName,
			&candidate.InternalReferenceCode,
			&candidate.State,
			&candidate.MatchInput,
			&candidate.MatchValue,
			&sourceSnapshotFetchedAt,
			&candidate.Confidence,
			&candidate.ConfidenceBand,
			&candidate.MatchStatus,
			&reasons,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		candidate.InternalProductID = scanInt4Ptr(internalProductID)
		candidate.SourceSnapshotFetchedAt = scanTimestamptz(sourceSnapshotFetchedAt)
		candidate.Reasons, err = unmarshalReasons(reasons)
		if err != nil {
			return nil, err
		}
		candidate.CreatedAt = createdAt.Time.UTC()
		candidate.UpdatedAt = updatedAt.Time.UTC()
		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

func (r *LinkCandidateRepository) GetLinkCandidate(ctx context.Context, candidateID string) (domain.LinkCandidate, bool, error) {
	if r.pool == nil {
		return domain.LinkCandidate{}, false, nil
	}
	row := r.pool.QueryRow(ctx, `
		SELECT
			candidate_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			internal_product_id, internal_product_name, internal_reference_code, state, match_input, match_value,
			source_snapshot_fetched_at, confidence, confidence_band, match_status, reasons, created_at, updated_at
		FROM product_link_candidates
		WHERE tenant_id = $1
		  AND candidate_id = $2
	`, r.tenantID, strings.TrimSpace(candidateID))

	var candidate domain.LinkCandidate
	var internalProductID pgtype.Int4
	var sourceSnapshotFetchedAt pgtype.Timestamptz
	var reasons []byte
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(
		&candidate.CandidateID,
		&candidate.InstallationID,
		&candidate.ProviderCode,
		&candidate.ProviderItemID,
		&candidate.ProviderVariationID,
		&internalProductID,
		&candidate.InternalProductName,
		&candidate.InternalReferenceCode,
		&candidate.State,
		&candidate.MatchInput,
		&candidate.MatchValue,
		&sourceSnapshotFetchedAt,
		&candidate.Confidence,
		&candidate.ConfidenceBand,
		&candidate.MatchStatus,
		&reasons,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.LinkCandidate{}, false, nil
		}
		return domain.LinkCandidate{}, false, err
	}
	candidate.InternalProductID = scanInt4Ptr(internalProductID)
	candidate.SourceSnapshotFetchedAt = scanTimestamptz(sourceSnapshotFetchedAt)
	parsedReasons, err := unmarshalReasons(reasons)
	if err != nil {
		return domain.LinkCandidate{}, false, err
	}
	candidate.Reasons = parsedReasons
	candidate.CreatedAt = createdAt.Time.UTC()
	candidate.UpdatedAt = updatedAt.Time.UTC()
	return candidate, true, nil
}

func (r *LinkCandidateRepository) GetProductLink(ctx context.Context, identity domain.ListingIdentity) (domain.ProductLink, bool, error) {
	if r.pool == nil {
		return domain.ProductLink{}, false, nil
	}
	row := r.pool.QueryRow(ctx, `
		SELECT
			installation_id, provider_code, provider_item_id, provider_variation_id, state,
			source_candidate_id, internal_product_id, internal_product_name, internal_reference_code,
			created_at, updated_at
		FROM product_links
		WHERE tenant_id = $1
		  AND installation_id = $2
		  AND provider_item_id = $3
		  AND provider_variation_id = $4
	`, r.tenantID, identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID)
	var link domain.ProductLink
	var internalProductID pgtype.Int4
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	if err := row.Scan(
		&link.InstallationID,
		&link.ProviderCode,
		&link.ProviderItemID,
		&link.ProviderVariationID,
		&link.State,
		&link.SourceCandidateID,
		&internalProductID,
		&link.InternalProductName,
		&link.InternalReferenceCode,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.ProductLink{}, false, nil
		}
		return domain.ProductLink{}, false, err
	}
	link.InternalProductID = scanInt4Ptr(internalProductID)
	link.CreatedAt = createdAt.Time.UTC()
	link.UpdatedAt = updatedAt.Time.UTC()
	return link, true, nil
}

func (r *LinkCandidateRepository) ApplyProductLinkTransition(ctx context.Context, transition domain.ProductLinkTransition) error {
	if r.pool == nil {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO product_links (
			tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			state, source_candidate_id, internal_product_id, internal_product_name, internal_reference_code,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12
		)
		ON CONFLICT (tenant_id, installation_id, provider_item_id, provider_variation_id) DO UPDATE SET
			provider_code = EXCLUDED.provider_code,
			state = EXCLUDED.state,
			source_candidate_id = EXCLUDED.source_candidate_id,
			internal_product_id = EXCLUDED.internal_product_id,
			internal_product_name = EXCLUDED.internal_product_name,
			internal_reference_code = EXCLUDED.internal_reference_code,
			updated_at = EXCLUDED.updated_at
	`, r.tenantID,
		transition.Link.InstallationID,
		transition.Link.ProviderCode,
		transition.Link.ProviderItemID,
		transition.Link.ProviderVariationID,
		transition.Link.State,
		transition.Link.SourceCandidateID,
		transition.Link.InternalProductID,
		transition.Link.InternalProductName,
		transition.Link.InternalReferenceCode,
		transition.Link.CreatedAt,
		transition.Link.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO product_link_audit_entries (
			tenant_id, audit_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			action, reason, source_candidate_id, actor_type, actor_id, actor_name,
			previous_state, next_state, previous_internal_product_id, next_internal_product_id, batch_id, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18
		)
	`, r.tenantID,
		transition.Audit.AuditID,
		transition.Audit.InstallationID,
		transition.Audit.ProviderCode,
		transition.Audit.ProviderItemID,
		transition.Audit.ProviderVariationID,
		transition.Audit.Action,
		transition.Audit.Reason,
		transition.Audit.SourceCandidateID,
		transition.Audit.Actor.ActorType,
		transition.Audit.Actor.ActorID,
		transition.Audit.Actor.ActorName,
		transition.Audit.PreviousState,
		transition.Audit.NextState,
		transition.Audit.PreviousInternalProductID,
		transition.Audit.NextInternalProductID,
		transition.Audit.BatchID,
		transition.Audit.CreatedAt,
	)
	if err != nil {
		return err
	}

	if transition.Decision != nil {
		if err := insertDecision(ctx, tx, r.tenantID, *transition.Decision); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// insertDecision appends one E10 row and stamps whatever decision it replaces.
// Both statements run inside the transition's transaction: a decision that
// landed without its link, or a link approved without its decision, would make
// the trail a claim rather than a record.
func insertDecision(ctx context.Context, tx pgx.Tx, tenantID string, decision domain.ProductLinkDecision) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO product_link_decisions (
			tenant_id, decision_id, installation_id, provider_item_id, provider_variation_id,
			rule_matched, actor, collisions_at_decision, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, tenantID,
		decision.DecisionID,
		decision.InstallationID,
		decision.ProviderItemID,
		decision.ProviderVariationID,
		decision.RuleMatched,
		decision.Actor,
		decision.CollisionsAtDecision,
		decision.CreatedAt,
	); err != nil {
		return err
	}
	// Stamp the rows this one replaces only after it exists — superseded_by
	// references a real decision, so the other order fails the foreign key.
	// Excluding the new row keeps it from stamping itself.
	_, err := tx.Exec(ctx, `
		UPDATE product_link_decisions
		SET superseded_by = $3
		WHERE tenant_id = $1 AND link_id = $2 AND superseded_by IS NULL AND decision_id <> $3
	`, tenantID, decision.LinkID, decision.DecisionID)
	return err
}

func (r *LinkCandidateRepository) ListDecisionsForLink(ctx context.Context, identity domain.ListingIdentity) ([]domain.ProductLinkDecision, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT decision_id, installation_id, provider_item_id, provider_variation_id, link_id,
		       rule_matched, actor, collisions_at_decision, created_at, COALESCE(superseded_by, '')
		FROM product_link_decisions
		WHERE tenant_id = $1 AND link_id = $2
		ORDER BY created_at ASC, decision_id ASC
	`, r.tenantID, domain.LinkID(identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	decisions := make([]domain.ProductLinkDecision, 0)
	for rows.Next() {
		var decision domain.ProductLinkDecision
		var collisions pgtype.Int4
		var createdAt pgtype.Timestamptz
		if err := rows.Scan(
			&decision.DecisionID,
			&decision.InstallationID,
			&decision.ProviderItemID,
			&decision.ProviderVariationID,
			&decision.LinkID,
			&decision.RuleMatched,
			&decision.Actor,
			&collisions,
			&createdAt,
			&decision.SupersededBy,
		); err != nil {
			return nil, err
		}
		decision.CollisionsAtDecision = scanInt4Ptr(collisions)
		decision.CreatedAt = createdAt.Time.UTC()
		decisions = append(decisions, decision)
	}
	return decisions, rows.Err()
}

// InsertBatch persists the module-owned batch audit row for a completed
// ApplyBatch run (S3, C02 "tabela própria").
func (r *LinkCandidateRepository) InsertBatch(ctx context.Context, batch domain.ProductLinkBatch) error {
	if r.pool == nil {
		return nil
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO product_link_batches (
			tenant_id, batch_id, installation_id, actor_type, actor_id, actor_name,
			requested_count, applied_count, failed_count, status, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11
		)
		ON CONFLICT (tenant_id, batch_id) DO UPDATE SET
			installation_id = EXCLUDED.installation_id,
			actor_type = EXCLUDED.actor_type,
			actor_id = EXCLUDED.actor_id,
			actor_name = EXCLUDED.actor_name,
			requested_count = EXCLUDED.requested_count,
			applied_count = EXCLUDED.applied_count,
			failed_count = EXCLUDED.failed_count,
			status = EXCLUDED.status
	`, r.tenantID,
		batch.BatchID,
		batch.InstallationID,
		batch.Actor.ActorType,
		batch.Actor.ActorID,
		batch.Actor.ActorName,
		batch.RequestedCount,
		batch.AppliedCount,
		batch.FailedCount,
		batch.Status,
		batch.CreatedAt,
	)
	return err
}

func (r *LinkCandidateRepository) ListProductLinks(ctx context.Context, installationID string, limit int) ([]domain.ProductLink, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			installation_id, provider_code, provider_item_id, provider_variation_id, state,
			source_candidate_id, internal_product_id, internal_product_name, internal_reference_code,
			created_at, updated_at
		FROM product_links
		WHERE tenant_id = $1
		  AND installation_id = $2
		ORDER BY updated_at DESC, provider_item_id DESC, provider_variation_id DESC
		LIMIT $3
	`, r.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ProductLink, 0)
	for rows.Next() {
		var link domain.ProductLink
		var internalProductID pgtype.Int4
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		if err := rows.Scan(
			&link.InstallationID,
			&link.ProviderCode,
			&link.ProviderItemID,
			&link.ProviderVariationID,
			&link.State,
			&link.SourceCandidateID,
			&internalProductID,
			&link.InternalProductName,
			&link.InternalReferenceCode,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		link.InternalProductID = scanInt4Ptr(internalProductID)
		link.CreatedAt = createdAt.Time.UTC()
		link.UpdatedAt = updatedAt.Time.UTC()
		items = append(items, link)
	}
	return items, rows.Err()
}

func (r *LinkCandidateRepository) ListProductLinkAuditEntries(ctx context.Context, installationID string, limit int) ([]domain.ProductLinkAuditEntry, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			audit_id, installation_id, provider_code, provider_item_id, provider_variation_id,
			action, reason, source_candidate_id, actor_type, actor_id, actor_name,
			previous_state, next_state, previous_internal_product_id, next_internal_product_id, batch_id, created_at
		FROM product_link_audit_entries
		WHERE tenant_id = $1
		  AND installation_id = $2
		ORDER BY created_at DESC, audit_id DESC
		LIMIT $3
	`, r.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ProductLinkAuditEntry, 0)
	for rows.Next() {
		var audit domain.ProductLinkAuditEntry
		var previousInternalProductID pgtype.Int4
		var nextInternalProductID pgtype.Int4
		var createdAt pgtype.Timestamptz
		if err := rows.Scan(
			&audit.AuditID,
			&audit.InstallationID,
			&audit.ProviderCode,
			&audit.ProviderItemID,
			&audit.ProviderVariationID,
			&audit.Action,
			&audit.Reason,
			&audit.SourceCandidateID,
			&audit.Actor.ActorType,
			&audit.Actor.ActorID,
			&audit.Actor.ActorName,
			&audit.PreviousState,
			&audit.NextState,
			&previousInternalProductID,
			&nextInternalProductID,
			&audit.BatchID,
			&createdAt,
		); err != nil {
			return nil, err
		}
		audit.PreviousInternalProductID = scanInt4Ptr(previousInternalProductID)
		audit.NextInternalProductID = scanInt4Ptr(nextInternalProductID)
		audit.CreatedAt = createdAt.Time.UTC()
		items = append(items, audit)
	}
	return items, rows.Err()
}

// scanAuditRow scans one product_link_audit_entries row in the column order
// shared by GetAuditEntry / LatestAuditForIdentity / ListAuditByBatch.
func scanAuditRow(row interface {
	Scan(dest ...any) error
}) (domain.ProductLinkAuditEntry, error) {
	var audit domain.ProductLinkAuditEntry
	var previousInternalProductID pgtype.Int4
	var nextInternalProductID pgtype.Int4
	var createdAt pgtype.Timestamptz
	err := row.Scan(
		&audit.AuditID,
		&audit.InstallationID,
		&audit.ProviderCode,
		&audit.ProviderItemID,
		&audit.ProviderVariationID,
		&audit.Action,
		&audit.Reason,
		&audit.SourceCandidateID,
		&audit.Actor.ActorType,
		&audit.Actor.ActorID,
		&audit.Actor.ActorName,
		&audit.PreviousState,
		&audit.NextState,
		&previousInternalProductID,
		&nextInternalProductID,
		&audit.BatchID,
		&createdAt,
	)
	if err != nil {
		return domain.ProductLinkAuditEntry{}, err
	}
	audit.PreviousInternalProductID = scanInt4Ptr(previousInternalProductID)
	audit.NextInternalProductID = scanInt4Ptr(nextInternalProductID)
	audit.CreatedAt = createdAt.Time.UTC()
	return audit, nil
}

const auditEntryColumns = `
	audit_id, installation_id, provider_code, provider_item_id, provider_variation_id,
	action, reason, source_candidate_id, actor_type, actor_id, actor_name,
	previous_state, next_state, previous_internal_product_id, next_internal_product_id, batch_id, created_at
`

// GetAuditEntry loads a single audit row by AuditID (S4 undo target lookup),
// tenant-scoped.
func (r *LinkCandidateRepository) GetAuditEntry(ctx context.Context, auditID string) (domain.ProductLinkAuditEntry, bool, error) {
	if r.pool == nil {
		return domain.ProductLinkAuditEntry{}, false, nil
	}
	row := r.pool.QueryRow(ctx, `
		SELECT `+auditEntryColumns+`
		FROM product_link_audit_entries
		WHERE tenant_id = $1
		  AND audit_id = $2
	`, r.tenantID, strings.TrimSpace(auditID))
	audit, err := scanAuditRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ProductLinkAuditEntry{}, false, nil
		}
		return domain.ProductLinkAuditEntry{}, false, err
	}
	return audit, true, nil
}

// LatestAuditForIdentity returns the most recent audit row for a listing
// identity (S4 undo ordering decision), tenant-scoped.
func (r *LinkCandidateRepository) LatestAuditForIdentity(ctx context.Context, identity domain.ListingIdentity) (domain.ProductLinkAuditEntry, bool, error) {
	if r.pool == nil {
		return domain.ProductLinkAuditEntry{}, false, nil
	}
	row := r.pool.QueryRow(ctx, `
		SELECT `+auditEntryColumns+`
		FROM product_link_audit_entries
		WHERE tenant_id = $1
		  AND installation_id = $2
		  AND provider_item_id = $3
		  AND provider_variation_id = $4
		ORDER BY created_at DESC, audit_id DESC
		LIMIT 1
	`, r.tenantID, identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID)
	audit, err := scanAuditRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ProductLinkAuditEntry{}, false, nil
		}
		return domain.ProductLinkAuditEntry{}, false, err
	}
	return audit, true, nil
}

// ListAuditByBatch returns every audit row tagged with batchID, newest
// first (S4 batch-undo fan-out), tenant-scoped.
func (r *LinkCandidateRepository) ListAuditByBatch(ctx context.Context, batchID string) ([]domain.ProductLinkAuditEntry, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT `+auditEntryColumns+`
		FROM product_link_audit_entries
		WHERE tenant_id = $1
		  AND batch_id = $2
		ORDER BY created_at DESC, audit_id DESC
	`, r.tenantID, strings.TrimSpace(batchID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ProductLinkAuditEntry, 0)
	for rows.Next() {
		audit, err := scanAuditRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, audit)
	}
	return items, rows.Err()
}

func marshalReasons(reasons []domain.LinkCandidateReason) ([]byte, error) {
	if reasons == nil {
		reasons = []domain.LinkCandidateReason{}
	}
	return json.Marshal(reasons)
}

func unmarshalReasons(raw []byte) ([]domain.LinkCandidateReason, error) {
	reasons := make([]domain.LinkCandidateReason, 0)
	if len(raw) == 0 {
		return reasons, nil
	}
	if err := json.Unmarshal(raw, &reasons); err != nil {
		return nil, err
	}
	return reasons, nil
}
