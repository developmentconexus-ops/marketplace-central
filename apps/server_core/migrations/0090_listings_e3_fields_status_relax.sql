-- 0090: listings E3 columns + lifecycle + raw, plus relax the status CHECK —
-- IC-07 (M-04 listings-ddl).
--
-- status: the provider's own vocabulary, verbatim, and it has more values than
-- the 4 this repo guessed at in 0036/0037 (a row observed live today already
-- carries 'paused', and ML's real vocabulary is wider than either CHECK ever
-- enumerated). IC-07 Enums And Statuses: "status de listing: verbatim do
-- provider ... sem CHECK — vocabulário do ML." Dropping the constraint by
-- name, not narrowing it further — the CHECK itself is the defect, any fixed
-- list this repo writes will drift from ML again.
ALTER TABLE listings
    DROP CONSTRAINT IF EXISTS listings_status_check;

-- E3 additive columns (design §7, IC-07 Required Inputs). Every column below
-- is NULLable with no DEFAULT — honest-unknown (ADR-17) until M-04's writer
-- (IngestListing) populates it. IF NOT EXISTS on every ADD COLUMN so a re-run
-- is a no-op.
--
-- Deliberately EXCLUDED (IC-07 decision, planning design §7 item 6):
-- commission_amount, commission_pct, free_shipping_cost — fee data lives only
-- in channel_fees (IC-01, ADR-09 provenance); dual-write here would be a drift
-- seam, not a feature.
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS sold_quantity int;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS category_id text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS condition text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS permalink text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS thumbnail text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS date_created_ml timestamptz;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS tags text[];
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS catalog_product_id text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS shipping_mode text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS free_shipping boolean;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS logistic_type text;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS available_quantity int;

-- Lifecycle support (ADR-06): a row's absence is only ever learned at the end
-- of a COMPLETE run (IC-06) — status is never inferred from absence.
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS last_seen_at timestamptz;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS absent_since timestamptz;

-- raw payload (ADR-03 permits it here, unlike order_shipments 0088 — /items
-- carries no PII) plus a marker for when the stored payload was truncated.
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS raw jsonb;
ALTER TABLE listings
    ADD COLUMN IF NOT EXISTS raw_truncated boolean;
