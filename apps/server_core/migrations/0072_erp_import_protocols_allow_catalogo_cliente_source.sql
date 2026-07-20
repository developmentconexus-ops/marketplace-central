-- 0072: allow the lenient client-catalog source on import protocols.
-- The two-source demo needs each import protocol tagged by origin so the dataset
-- selector can distinguish the operator's Sankhya export ('xlsx') from a prospect
-- catalog imported through the lenient path ('catalogo_cliente'). The prior CHECK
-- pinned source to 'xlsx' only, which rejected lenient imports at insert time.
ALTER TABLE erp_import_protocols DROP CONSTRAINT IF EXISTS erp_import_protocols_source_check;
ALTER TABLE erp_import_protocols ADD CONSTRAINT erp_import_protocols_source_check
  CHECK (source IN ('xlsx', 'catalogo_cliente'));
