-- 0073: allow the lenient path's honest-unknown issue codes.
-- The client-catalog (lenient) import emits MISSING_CUSTO / MISSING_ESTOQUE as
-- WARNING issues when a workbook omits those columns (ADR-17: recorded as
-- unknown, never fabricated). The prior code CHECK did not list them, so the
-- issue insert failed and rolled the whole import back.
ALTER TABLE erp_import_issues DROP CONSTRAINT IF EXISTS erp_import_issues_code_check;
ALTER TABLE erp_import_issues ADD CONSTRAINT erp_import_issues_code_check
  CHECK (code = ANY (ARRAY[
    'EMPTY_CODPROD', 'DUPLICATE_CODPROD', 'EMPTY_DESCRPROD',
    'INVALID_CUSTO', 'INVALID_ESTOQUE', 'INVALID_EAN', 'INVALID_NCM',
    'MISSING_REQUIRED_COLUMN', 'INVALID_FILE', 'IMPORT_IN_PROGRESS',
    'DUPLICATE_FILE', 'MISSING_CUSTO', 'MISSING_ESTOQUE'
  ]::text[]));
