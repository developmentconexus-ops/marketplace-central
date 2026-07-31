-- 0084: preserve ERP sellability fields on the product mirror (S1-SCHEMA).
--
-- Both source values are nullable because an ERP row may not provide them;
-- missing facts remain SQL NULL under ADR-17.
ALTER TABLE products_mirror
    ADD COLUMN usoprod TEXT;

ALTER TABLE products_mirror
    ADD COLUMN ad_ecommerce TEXT;
