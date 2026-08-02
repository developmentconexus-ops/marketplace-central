-- 0096: orders.currency existia no contrato de leitura sem produtor nenhum —
-- null em 38/38. A fonte é o currency_id que o ML já manda em cada pedido e
-- que o DTO passou a declarar (Task 3 desta fatia). Nullable: pedido antigo
-- não tem a informação e desconhecido nunca vira 'BRL' por conveniência
-- (ADR-17).
ALTER TABLE orders_marketplace_orders ADD COLUMN IF NOT EXISTS currency text;
