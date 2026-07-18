ALTER TABLE product_link_candidates
  ADD COLUMN IF NOT EXISTS confidence integer NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS confidence_band text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS match_status text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS reasons jsonb NOT NULL DEFAULT '[]'::jsonb;
