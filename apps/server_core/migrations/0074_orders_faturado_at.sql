-- 0074: track the operator's "faturado" (invoiced) mark on orders.
-- FATURADO workflow (goal B): an unshipped PAID order shows as "A FATURAR",
-- not "A ENVIAR", until the operator confirms it was invoiced. This column is
-- our-DB only (never a Mercado Livre write). Nullable, no default, no
-- backfill: nil means "not yet marked invoiced" is an honest unknown, never a
-- fabricated timestamp (ADR-17).
ALTER TABLE orders_marketplace_orders ADD COLUMN IF NOT EXISTS faturado_at timestamptz;
