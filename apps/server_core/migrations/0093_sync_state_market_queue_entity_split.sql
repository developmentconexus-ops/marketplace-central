-- 0093: split the "market" sync_state entity into two.
--
-- entity='market' was shared by two unrelated producers: erp_import's
-- MarketEnqueuer (a pending-codigos queue, keyed by real installation_id,
-- never records a success) and F-A3's periodic market-evidence collection job
-- (tenant-wide, sentinel installation_id). Same (tenant_id, entity) key,
-- different installation_id, so no PK collision — but /sync/health returns
-- both rows with entity="market" and no installation_id in the DTO, so the
-- frontend cannot tell them apart. React then keys two rows on the same
-- "market" string, which makes their rendering order/identity undefined.
--
-- The queue keeps no cursor semantics a sync scheduler understands (no
-- success, no schedule) — it only borrows sync_state's cursor JSONB as
-- storage. Renaming its rows to entity='market_queue' gives it back a
-- distinct identity so the health screen shows one row per real concept.
UPDATE sync_state
   SET entity = 'market_queue'
 WHERE entity = 'market'
   AND installation_id <> 'market';
