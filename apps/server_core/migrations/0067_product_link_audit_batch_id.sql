ALTER TABLE product_link_audit_entries
  ADD COLUMN IF NOT EXISTS batch_id text NOT NULL DEFAULT '';
