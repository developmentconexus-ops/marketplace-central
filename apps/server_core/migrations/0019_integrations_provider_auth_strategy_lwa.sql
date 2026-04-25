ALTER TABLE integration_provider_definitions
  DROP CONSTRAINT IF EXISTS integration_provider_definitions_auth_strategy_check;

ALTER TABLE integration_provider_definitions
  ADD CONSTRAINT integration_provider_definitions_auth_strategy_check
  CHECK (auth_strategy IN ('oauth2', 'api_key', 'token', 'lwa', 'none', 'unknown'));
