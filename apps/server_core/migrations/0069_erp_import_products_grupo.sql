ALTER TABLE erp_import_products
  ADD COLUMN IF NOT EXISTS grupo text,
  ADD COLUMN IF NOT EXISTS descrgrupo text;
