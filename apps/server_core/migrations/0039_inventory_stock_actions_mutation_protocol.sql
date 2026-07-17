ALTER TABLE inventory_stock_actions
    ADD COLUMN IF NOT EXISTS mutation_protocol_id TEXT;
