ALTER TABLE integration_operation_runs
  ADD COLUMN IF NOT EXISTS translated_error_code text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS provider_evidence_json jsonb NOT NULL DEFAULT '{}'::jsonb,
  ADD COLUMN IF NOT EXISTS duration_ms bigint NOT NULL DEFAULT 0;
