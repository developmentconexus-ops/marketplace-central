ALTER TABLE market_aggregates
  ADD COLUMN IF NOT EXISTS max_valid_amount NUMERIC,
  ADD COLUMN IF NOT EXISTS max_valid_currency TEXT;

ALTER TABLE market_aggregates
  ADD CONSTRAINT market_aggregates_max_valid_pair_check
    CHECK ((max_valid_amount IS NULL AND max_valid_currency IS NULL) OR (max_valid_amount > 0 AND max_valid_currency = 'BRL'));
