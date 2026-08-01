-- 0092: listing_variations lookup index — IC-07 (M-04 listings-ddl).
--
-- Supports the read path that fans a listing out to its variations without
-- the installation_id/provider prefix of the PK.
CREATE INDEX IF NOT EXISTS idx_listing_variations_tenant_provider_listing
    ON listing_variations (tenant_id, provider_listing_id);
