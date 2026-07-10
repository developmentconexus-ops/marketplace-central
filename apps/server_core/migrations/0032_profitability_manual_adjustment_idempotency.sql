ALTER TABLE profitability_manual_adjustments
    ADD COLUMN IF NOT EXISTS idempotency_key text NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS uq_profitability_manual_adjustments_idempotency
    ON profitability_manual_adjustments (tenant_id, installation_id, idempotency_key)
    WHERE idempotency_key <> '';
