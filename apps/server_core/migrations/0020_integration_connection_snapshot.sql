ALTER TABLE integration_installations
  ADD COLUMN IF NOT EXISTS connection_snapshot_json jsonb NOT NULL DEFAULT '{}'::jsonb;

UPDATE integration_installations
SET connection_snapshot_json = jsonb_build_object(
  'state',
    CASE status
      WHEN 'draft' THEN 'draft'
      WHEN 'pending_connection' THEN 'pending_connection'
      WHEN 'connected' THEN 'connected'
      WHEN 'degraded' THEN 'degraded'
      WHEN 'requires_reauth' THEN 'needs_reauth'
      ELSE 'disconnected'
    END,
  'health', health_status,
  'provider_code', provider_code,
  'external_account_id', external_account_id,
  'external_account_name', external_account_name,
  'auth_strategy', 'unknown',
  'last_verified_at', to_jsonb(last_verified_at),
  'expires_at', NULL,
  'next_action',
    CASE status
      WHEN 'draft' THEN 'configure'
      WHEN 'pending_connection' THEN 'none'
      WHEN 'connected' THEN 'none'
      WHEN 'degraded' THEN 'retry'
      WHEN 'requires_reauth' THEN 'reauth'
      ELSE 'authorize'
    END,
  'reauth_reason', ''
)
WHERE connection_snapshot_json = '{}'::jsonb;
