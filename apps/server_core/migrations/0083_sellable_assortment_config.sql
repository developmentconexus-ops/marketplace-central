-- 0083: per-tenant sellable assortment defaults (S1-SCHEMA).
--
-- The rule is stored on the existing active_source row so the catalog policy
-- follows the tenant's selected source without introducing a second config
-- table. Defaults are the operator-ratified policy for existing tenants.
ALTER TABLE active_source
    ADD COLUMN only_revenda BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE active_source
    ADD COLUMN only_em_estoque BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE active_source
    ADD COLUMN only_ecommerce BOOLEAN NOT NULL DEFAULT false;
