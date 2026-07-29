-- 0085: preserve ERP sellability fields on the imported product snapshot.
--
-- Both source values are nullable because an ERP row may not provide them;
-- missing facts remain SQL NULL under ADR-17.
ALTER TABLE erp_import_products
    ADD COLUMN usoprod TEXT;

ALTER TABLE erp_import_products
    ADD COLUMN ad_ecommerce TEXT;
