-- 0093: provider_status_detail and cancellation_detail become genuinely
-- nullable on orders_marketplace_orders — ADR-17 reachable in the column, not
-- just the Go type.
--
-- 0027 declared both `text NOT NULL DEFAULT ''`. That constraint is the
-- actual root of the bug: "the provider didn't send this field" and "the
-- provider sent an empty string" were forced to collapse into the same
-- stored value ('') no matter what the application layer did — which is
-- exactly why cancellation_detail showed up populated on 38/38 orders,
-- including the 7 cancelled ones where the real cancellation reason was the
-- thing actually missing.
--
-- domain.MarketplaceOrder's ProviderStatusDetail/CancellationDetail are now
-- *string (nil = absent, see order.go). Without this migration, upserting a
-- nil value would still fail the NOT NULL constraint at the database — the
-- Go-level fix alone cannot make NULL reachable.
--
-- Existing rows keep their stored '' as-is: which historical rows were
-- "provider sent empty" vs "provider omitted the field" is not recoverable
-- from the stored value itself, so this migration does not attempt a
-- backfill (rewriting ambiguous history into NULL would be fabricating a
-- fact, the same failure mode ADR-17 forbids in the other direction).
ALTER TABLE orders_marketplace_orders
    ALTER COLUMN provider_status_detail DROP NOT NULL;
ALTER TABLE orders_marketplace_orders
    ALTER COLUMN provider_status_detail DROP DEFAULT;
ALTER TABLE orders_marketplace_orders
    ALTER COLUMN cancellation_detail DROP NOT NULL;
ALTER TABLE orders_marketplace_orders
    ALTER COLUMN cancellation_detail DROP DEFAULT;
