ALTER TABLE erp_import_products
  ALTER COLUMN custo DROP NOT NULL,
  ALTER COLUMN stock_physical DROP NOT NULL;

ALTER TABLE erp_import_products
  DROP CONSTRAINT IF EXISTS erp_import_products_custo_check;
ALTER TABLE erp_import_products
  ADD CONSTRAINT erp_import_products_custo_check CHECK (custo IS NULL OR custo > 0);

ALTER TABLE erp_import_products
  DROP CONSTRAINT IF EXISTS erp_import_products_stock_physical_check;
ALTER TABLE erp_import_products
  ADD CONSTRAINT erp_import_products_stock_physical_check CHECK (stock_physical IS NULL OR stock_physical >= 0);
