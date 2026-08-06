// Package postgres is Catalog's only writer. No other package in the platform
// issues DML against the catalog schema, and no other schema holds a foreign
// key into it.
package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Repository implements application.Store. The context's published reader is
// SummaryReader at the bottom of this file, which wraps it.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository wires the repository to a pool.
func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// withTenant runs fn in a transaction whose session is scoped to t.
//
// SET LOCAL and not SET: the scope dies with the transaction, so a pooled
// connection handed to the next request never carries the previous tenant.
func (r *Repository) withTenant(ctx context.Context, t tenant.ID, fn func(pgx.Tx) error) error {
	if t.IsZero() {
		return errors.New("catalog/postgres: tenant is required")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", t.String()); err != nil {
		return fmt.Errorf("catalog/postgres: scope tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// factColumns renders a Fact[string] as its three columns.
func factColumns(f fact.Fact[string]) (state string, value *string, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

// Insert writes a new product and everything that hangs off it, atomically.
func (r *Repository) Insert(ctx context.Context, p domain.Product) error {
	return r.withTenant(ctx, p.Tenant(), func(tx pgx.Tx) error {
		return writeProduct(ctx, tx, p, true)
	})
}

// Update writes a new version of an existing product, atomically.
func (r *Repository) Update(ctx context.Context, p domain.Product) error {
	return r.withTenant(ctx, p.Tenant(), func(tx pgx.Tx) error {
		return writeProduct(ctx, tx, p, false)
	})
}

func writeProduct(ctx context.Context, tx pgx.Tx, p domain.Product, insert bool) error {
	state, value, reason := factColumns(p.Description())
	if insert {
		_, err := tx.Exec(ctx, `
			INSERT INTO catalog.products
				(tenant_id, product_id, version, description_state, description_value,
				 description_reason, last_payload_hash)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			p.Tenant().String(), p.ID().String(), p.Version(),
			state, value, reason, p.LastPayloadHash())
		if err != nil {
			return fmt.Errorf("catalog/postgres: insert product: %w", err)
		}
	} else {
		tag, err := tx.Exec(ctx, `
			UPDATE catalog.products
			   SET version = $3, description_state = $4, description_value = $5,
			       description_reason = $6, last_payload_hash = $7, updated_at = now()
			 WHERE tenant_id = $1 AND product_id = $2`,
			p.Tenant().String(), p.ID().String(), p.Version(),
			state, value, reason, p.LastPayloadHash())
		if err != nil {
			return fmt.Errorf("catalog/postgres: update product: %w", err)
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("catalog/postgres: update touched %d rows for %s",
				tag.RowsAffected(), p.ID())
		}
	}

	for _, id := range p.Identifiers() {
		if _, err := tx.Exec(ctx, `
			INSERT INTO catalog.product_identifiers (tenant_id, product_id, kind, value)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (tenant_id, product_id, kind, value) DO NOTHING`,
			p.Tenant().String(), p.ID().String(), string(id.Kind()), id.Value()); err != nil {
			return fmt.Errorf("catalog/postgres: insert identifier %s: %w", id, err)
		}
	}

	for _, k := range p.SourceKeys() {
		if _, err := tx.Exec(ctx, `
			INSERT INTO catalog.source_product_keys
				(tenant_id, source_system, source_instance, object_kind, external_key, product_id)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (tenant_id, source_system, source_instance, object_kind, external_key)
			DO UPDATE SET product_id = EXCLUDED.product_id`,
			k.Tenant().String(), k.System(), k.Instance(), k.ObjectKind(), k.ExternalKey(),
			p.ID().String()); err != nil {
			return fmt.Errorf("catalog/postgres: insert source key %s: %w", k, err)
		}
	}

	e := p.LastEvidence()
	if e.IsZero() {
		return fmt.Errorf("catalog/postgres: product %s has no evidence to record", p.ID())
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO catalog.source_observations
			(tenant_id, product_id, payload_hash, source_system, object_kind,
			 external_key, observed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (tenant_id, product_id, payload_hash) DO NOTHING`,
		p.Tenant().String(), p.ID().String(), e.PayloadHash(), e.System(),
		e.ObjectKind(), e.ExternalKey(), e.ObservedAt().UTC()); err != nil {
		return fmt.Errorf("catalog/postgres: record observation: %w", err)
	}
	return nil
}

// BySourceKey resolves a source address to the product it maps to.
func (r *Repository) BySourceKey(ctx context.Context, k contracts.SourceProductKey) (domain.Product, bool, error) {
	var out domain.Product
	var found bool
	err := r.withTenant(ctx, k.Tenant(), func(tx pgx.Tx) error {
		var productID string
		row := tx.QueryRow(ctx, `
			SELECT product_id FROM catalog.source_product_keys
			 WHERE tenant_id = $1 AND source_system = $2 AND source_instance = $3
			   AND object_kind = $4 AND external_key = $5`,
			k.Tenant().String(), k.System(), k.Instance(), k.ObjectKind(), k.ExternalKey())
		switch err := row.Scan(&productID); {
		case errors.Is(err, pgx.ErrNoRows):
			return nil
		case err != nil:
			return fmt.Errorf("catalog/postgres: resolve source key: %w", err)
		}
		p, ok, err := loadProduct(ctx, tx, k.Tenant(), productID)
		if err != nil {
			return err
		}
		out, found = p, ok
		return nil
	})
	return out, found, err
}

// loadProduct rebuilds an aggregate from its rows.
func loadProduct(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) (domain.Product, bool, error) {
	var version int
	var state, hash string
	var value, reason *string
	var recordedAt time.Time
	row := tx.QueryRow(ctx, `
		SELECT version, description_state, description_value, description_reason,
		       last_payload_hash, updated_at
		  FROM catalog.products WHERE tenant_id = $1 AND product_id = $2`,
		t.String(), productID)
	switch err := row.Scan(&version, &state, &value, &reason, &hash, &recordedAt); {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.Product{}, false, nil
	case err != nil:
		return domain.Product{}, false, fmt.Errorf("catalog/postgres: load product: %w", err)
	}

	// Rehydration goes back through the same constructors that guard the write
	// path. A row that cannot be turned back into a valid aggregate is a defect
	// we want to hear about here, not three layers later.
	id, err := domain.NewProductID(productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	ident, err := loadIdentifiers(ctx, tx, t, productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	keys, err := loadSourceKeys(ctx, tx, t, productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	if len(keys) == 0 {
		return domain.Product{}, false, fmt.Errorf("catalog/postgres: product %s has no source key", productID)
	}

	e, err := loadLatestEvidence(ctx, tx, t, productID, hash)
	if err != nil {
		return domain.Product{}, false, err
	}
	desc, err := rehydrateFact(state, value, reason, e)
	if err != nil {
		return domain.Product{}, false, err
	}

	obs := contracts.ProductObservation{
		Key: keys[0], Description: desc, Identifiers: ident, Evidence: e,
	}
	p, err := domain.NewProduct(id, obs)
	if err != nil {
		return domain.Product{}, false, err
	}
	// Every key, restored directly. Replaying Apply here would report Idempotent
	// on the identical payload hash and drop every key after the first.
	p = domain.ReconstituteSourceKeys(p, keys)
	return domain.ReconstituteVersion(p, version, hash), true, nil
}

// loadLatestEvidence returns how we actually learned this product's current
// version. It reads the observation that produced the stored payload hash, so
// the rehydrated Fact carries the source system that observed it — not this
// package's own name, which is what it used to fabricate.
func loadLatestEvidence(ctx context.Context, tx pgx.Tx, t tenant.ID, productID, hash string) (provenance.Evidence, error) {
	var system, objectKind, externalKey string
	var observedAt time.Time
	row := tx.QueryRow(ctx, `
		SELECT source_system, object_kind, external_key, observed_at
		  FROM catalog.source_observations
		 WHERE tenant_id = $1 AND product_id = $2 AND payload_hash = $3`,
		t.String(), productID, hash)
	switch err := row.Scan(&system, &objectKind, &externalKey, &observedAt); {
	case errors.Is(err, pgx.ErrNoRows):
		// Not a missing optional. A product whose current version has no recorded
		// observation cannot say where it came from, and inventing one here is the
		// exact fabrication this change removes.
		return provenance.Evidence{}, fmt.Errorf(
			"catalog/postgres: product %s version hash %s has no source observation", productID, hash)
	case err != nil:
		return provenance.Evidence{}, fmt.Errorf("catalog/postgres: load observation: %w", err)
	}
	return provenance.NewEvidence(system, objectKind, externalKey, observedAt.UTC(), hash)
}

func loadIdentifiers(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) ([]contracts.Identifier, error) {
	rows, err := tx.Query(ctx, `
		SELECT kind, value FROM catalog.product_identifiers
		 WHERE tenant_id = $1 AND product_id = $2 ORDER BY kind, value`,
		t.String(), productID)
	if err != nil {
		return nil, fmt.Errorf("catalog/postgres: load identifiers: %w", err)
	}
	defer rows.Close()
	var out []contracts.Identifier
	for rows.Next() {
		var kind, value string
		if err := rows.Scan(&kind, &value); err != nil {
			return nil, err
		}
		id, err := contracts.NewIdentifier(contracts.IdentifierKind(kind), value)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func loadSourceKeys(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) ([]contracts.SourceProductKey, error) {
	rows, err := tx.Query(ctx, `
		SELECT source_system, source_instance, object_kind, external_key
		  FROM catalog.source_product_keys
		 WHERE tenant_id = $1 AND product_id = $2
		 ORDER BY source_system, source_instance, object_kind, external_key`,
		t.String(), productID)
	if err != nil {
		return nil, fmt.Errorf("catalog/postgres: load source keys: %w", err)
	}
	defer rows.Close()
	var out []contracts.SourceProductKey
	for rows.Next() {
		var system, instance, kind, key string
		if err := rows.Scan(&system, &instance, &kind, &key); err != nil {
			return nil, err
		}
		k, err := contracts.NewSourceProductKey(t, system, instance, kind, key)
		if err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// rehydrateFact turns three columns back into a Fact, through the same
// constructors the write path uses.
func rehydrateFact(state string, value, reason *string, e provenance.Evidence) (fact.Fact[string], error) {
	text := ""
	if value != nil {
		text = *value
	}
	why := ""
	if reason != nil {
		why = *reason
	}
	switch state {
	case fact.Known.String():
		return fact.NewKnown(text, e)
	case fact.Estimated.String():
		return fact.NewEstimated(text, why, e)
	case fact.NotApplicable.String():
		return fact.NewNotApplicable[string](why, e)
	default:
		return fact.NewUnknown[string](why, e)
	}
}

// ByIdentifier returns every product of this tenant carrying the identifier.
func (r *Repository) ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error) {
	var out []domain.Product
	err := r.withTenant(ctx, t, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT product_id FROM catalog.product_identifiers
			 WHERE tenant_id = $1 AND kind = $2 AND value = $3 ORDER BY product_id`,
			t.String(), string(id.Kind()), id.Value())
		if err != nil {
			return fmt.Errorf("catalog/postgres: lookup identifier: %w", err)
		}
		var ids []string
		for rows.Next() {
			var productID string
			if err := rows.Scan(&productID); err != nil {
				rows.Close()
				return err
			}
			ids = append(ids, productID)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		for _, productID := range ids {
			p, ok, err := loadProduct(ctx, tx, t, productID)
			if err != nil {
				return err
			}
			if ok {
				out = append(out, p)
			}
		}
		return nil
	})
	return out, err
}

// ByProductID reads one product as a summary. SummaryReader forwards to it.
func (r *Repository) ByProductID(ctx context.Context, t tenant.ID, productID string) (port.Summary, bool, error) {
	var out port.Summary
	var found bool
	err := r.withTenant(ctx, t, func(tx pgx.Tx) error {
		p, ok, err := loadProduct(ctx, tx, t, productID)
		if err != nil || !ok {
			return err
		}
		out, found = summarise(p), true
		return nil
	})
	return out, found, err
}

// summarise flattens an aggregate for a consumer. The knowledge state travels
// with the value: this is the ONLY place the empty string is allowed to stand
// for an unknown description, and it is allowed only because DescriptionState
// says so beside it.
func summarise(p domain.Product) port.Summary {
	desc, known := p.Description().Value()
	if !known {
		desc = ""
	}
	return port.Summary{
		ProductID:                 p.ID().String(),
		Description:               desc,
		DescriptionState:          p.Description().State().String(),
		DescriptionEvidenceSystem: p.Description().Evidence().System(),
		Identifiers:               p.Identifiers(),
		SourceKeys:                p.SourceKeys(),
		Version:                   p.Version(),
	}
}

// ULIDFactory mints canonical identifiers. The value is random and carries no
// source semantics whatsoever, which is the whole property.
type ULIDFactory struct{}

// NewULIDFactory builds the factory.
func NewULIDFactory() *ULIDFactory { return &ULIDFactory{} }

// NewProductID mints prd_ + 32 hex characters from crypto/rand.
func (f *ULIDFactory) NewProductID() (domain.ProductID, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return domain.ProductID{}, fmt.Errorf("catalog: mint id: %w", err)
	}
	return domain.NewProductID(domain.ProductIDPrefix + hex.EncodeToString(b[:]))
}

// SummaryReader adapts the repository to port.Reader without a second method of
// the same name on the same type.
type SummaryReader struct{ repo *Repository }

// NewSummaryReader wraps the repository as the context's published reader.
func NewSummaryReader(r *Repository) *SummaryReader { return &SummaryReader{repo: r} }

// ByProductID implements port.Reader.
func (s *SummaryReader) ByProductID(ctx context.Context, t tenant.ID, productID string) (port.Summary, bool, error) {
	return s.repo.ByProductID(ctx, t, productID)
}

// ByIdentifier implements port.Reader.
func (s *SummaryReader) ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]port.Summary, error) {
	products, err := s.repo.ByIdentifier(ctx, t, id)
	if err != nil {
		return nil, err
	}
	out := make([]port.Summary, 0, len(products))
	for _, p := range products {
		out = append(out, summarise(p))
	}
	return out, nil
}

var (
	_ application.Store     = (*Repository)(nil)
	_ application.IDFactory = (*ULIDFactory)(nil)
	_ port.Reader           = (*SummaryReader)(nil)
)
