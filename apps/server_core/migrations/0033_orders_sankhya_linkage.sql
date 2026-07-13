ALTER TABLE orders_marketplace_order_items
    ADD COLUMN IF NOT EXISTS mpc_line_id text,
    ADD COLUMN IF NOT EXISTS reconciliation_state text;

UPDATE orders_marketplace_order_items
SET mpc_line_id = 'mpl_' || replace(gen_random_uuid()::text, '-', ''),
    reconciliation_state = 'legacy_unresolved'
WHERE mpc_line_id IS NULL OR reconciliation_state IS NULL;

ALTER TABLE orders_marketplace_order_items
    ALTER COLUMN mpc_line_id SET NOT NULL,
    ALTER COLUMN reconciliation_state SET NOT NULL,
    ADD CONSTRAINT orders_marketplace_order_items_mpc_line_format
        CHECK (mpc_line_id ~ '^mpl_[0-9a-f]{32}$'),
    ADD CONSTRAINT orders_marketplace_order_items_reconciliation_state
        CHECK (reconciliation_state IN ('stable', 'legacy_unresolved', 'ambiguous')),
    ADD CONSTRAINT orders_marketplace_order_items_scoped_line_identity
        UNIQUE (tenant_id, installation_id, provider_order_id, mpc_line_id);

CREATE OR REPLACE FUNCTION orders_reject_mpc_line_id_update()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.mpc_line_id IS DISTINCT FROM OLD.mpc_line_id THEN
        RAISE EXCEPTION 'orders mpc line identity is immutable' USING ERRCODE = '23514';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER orders_marketplace_order_items_mpc_line_immutable
    BEFORE UPDATE OF mpc_line_id ON orders_marketplace_order_items
    FOR EACH ROW EXECUTE FUNCTION orders_reject_mpc_line_id_update();

CREATE TABLE orders_sankhya_linkage_events (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    event_id text NOT NULL,
    event_type text NOT NULL,
    provider_order_id text NOT NULL,
    external_order_key text NOT NULL,
    internal_document_id bigint NOT NULL,
    actor_type text NOT NULL,
    actor_id text NOT NULL,
    reason text NOT NULL,
    idempotency_key text NOT NULL,
    source_at timestamptz NOT NULL,
    recorded_at timestamptz NOT NULL,
    previous_event_id text NULL,
    configuration_revision text NOT NULL,
    evidence_state text NOT NULL,
    evidence_reference text NOT NULL DEFAULT '',
    PRIMARY KEY (tenant_id, installation_id, event_id),
    CONSTRAINT orders_sankhya_linkage_event_type
        CHECK (event_type IN ('confirmed')),
    CONSTRAINT orders_sankhya_linkage_document_positive
        CHECK (internal_document_id > 0),
    CONSTRAINT orders_sankhya_linkage_external_order_key_shape
        CHECK (external_order_key = 'ml:v1:' || installation_id || ':' || provider_order_id),
    CONSTRAINT orders_sankhya_linkage_audit_nonempty
        CHECK (
            btrim(event_id) <> '' AND btrim(external_order_key) <> '' AND
            btrim(actor_type) <> '' AND
            btrim(actor_id) <> '' AND btrim(reason) <> '' AND
            btrim(idempotency_key) <> '' AND btrim(configuration_revision) <> ''
        ),
    CONSTRAINT orders_sankhya_linkage_evidence_state
        CHECK (evidence_state IN ('unknown', 'exact')),
    CONSTRAINT orders_sankhya_linkage_order_fk
        FOREIGN KEY (tenant_id, installation_id, provider_order_id)
        REFERENCES orders_marketplace_orders (tenant_id, installation_id, provider_order_id),
    CONSTRAINT orders_sankhya_linkage_previous_event_fk
        FOREIGN KEY (tenant_id, installation_id, previous_event_id)
        REFERENCES orders_sankhya_linkage_events (tenant_id, installation_id, event_id),
    CONSTRAINT orders_sankhya_linkage_one_current_order
        UNIQUE (tenant_id, installation_id, provider_order_id),
    CONSTRAINT orders_sankhya_linkage_external_order_key
        UNIQUE (tenant_id, installation_id, external_order_key),
    CONSTRAINT orders_sankhya_linkage_scoped_idempotency
        UNIQUE (tenant_id, installation_id, idempotency_key),
    CONSTRAINT orders_sankhya_linkage_line_parent_identity
        UNIQUE (tenant_id, installation_id, event_id, provider_order_id, internal_document_id)
);

CREATE TABLE orders_sankhya_linkage_lines (
    tenant_id text NOT NULL,
    installation_id text NOT NULL,
    event_id text NOT NULL,
    provider_order_id text NOT NULL,
    mpc_line_id text NOT NULL,
    internal_document_id bigint NOT NULL,
    internal_line_number integer NOT NULL,
    PRIMARY KEY (tenant_id, installation_id, event_id, mpc_line_id),
    CONSTRAINT orders_sankhya_linkage_line_document_positive
        CHECK (internal_document_id > 0 AND internal_line_number > 0),
    CONSTRAINT orders_sankhya_linkage_line_event_fk
        FOREIGN KEY (tenant_id, installation_id, event_id, provider_order_id, internal_document_id)
        REFERENCES orders_sankhya_linkage_events (
            tenant_id, installation_id, event_id, provider_order_id, internal_document_id
        ),
    CONSTRAINT orders_sankhya_linkage_line_order_identity_fk
        FOREIGN KEY (tenant_id, installation_id, provider_order_id, mpc_line_id)
        REFERENCES orders_marketplace_order_items (
            tenant_id, installation_id, provider_order_id, mpc_line_id
        ) DEFERRABLE INITIALLY DEFERRED,
    CONSTRAINT orders_sankhya_linkage_one_origin_per_mpc_line
        UNIQUE (tenant_id, installation_id, provider_order_id, mpc_line_id),
    CONSTRAINT orders_sankhya_linkage_no_active_origin_reuse
        UNIQUE (tenant_id, installation_id, internal_document_id, internal_line_number)
);

CREATE OR REPLACE FUNCTION orders_reject_sankhya_linkage_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'orders sankhya linkage ledger is append-only' USING ERRCODE = '23514';
END;
$$;

CREATE TRIGGER orders_sankhya_linkage_events_append_only
    BEFORE UPDATE OR DELETE ON orders_sankhya_linkage_events
    FOR EACH ROW EXECUTE FUNCTION orders_reject_sankhya_linkage_mutation();

CREATE TRIGGER orders_sankhya_linkage_lines_append_only
    BEFORE UPDATE OR DELETE ON orders_sankhya_linkage_lines
    FOR EACH ROW EXECUTE FUNCTION orders_reject_sankhya_linkage_mutation();
