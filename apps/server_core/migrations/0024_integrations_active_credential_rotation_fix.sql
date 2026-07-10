CREATE OR REPLACE FUNCTION prevent_deactivating_or_revoking_referenced_credential()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF (NEW.is_active = false OR NEW.revoked_at IS NOT NULL)
     AND EXISTS (
       SELECT 1
       FROM integration_installations i
       WHERE i.tenant_id = NEW.tenant_id
         AND i.installation_id = NEW.installation_id
         AND i.active_credential_id = NEW.credential_id
         AND NOT EXISTS (
           SELECT 1
           FROM integration_credentials replacement
           WHERE replacement.tenant_id = NEW.tenant_id
             AND replacement.installation_id = NEW.installation_id
             AND replacement.credential_id <> NEW.credential_id
             AND replacement.is_active = true
             AND replacement.revoked_at IS NULL
         )
     ) THEN
    RAISE EXCEPTION 'INTEGRATIONS_ACTIVE_CREDENTIAL_STATE_CONFLICT';
  END IF;

  RETURN NEW;
END;
$$;
