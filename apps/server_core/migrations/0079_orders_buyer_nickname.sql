-- Mercado Livre returns the buyer nickname on the order detail payload, but the
-- import dropped it, so the orders list had no buyer identity at all and the
-- COMPRADOR column rendered blank for every row. Nullable: an order whose
-- payload carries no buyer stays an honest unknown instead of a fabricated name.
ALTER TABLE orders_marketplace_orders
    ADD COLUMN IF NOT EXISTS buyer_nickname text;
