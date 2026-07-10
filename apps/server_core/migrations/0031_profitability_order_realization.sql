ALTER TABLE profitability_profit_snapshots
    ADD COLUMN realization_state text NOT NULL DEFAULT 'unknown';

ALTER TABLE profitability_profit_snapshots
    ALTER COLUMN realization_state DROP DEFAULT;

ALTER TABLE profitability_profit_snapshots
    ADD CONSTRAINT profitability_profit_snapshots_realization_state_valid
        CHECK (realization_state IN ('realized', 'not_realized', 'unknown')) NOT VALID;
