-- 0075: sync_state — cadence-agnostic sync bookkeeping (E8, MIS-006 M-01).
-- Foundation of the sync engine: one row per (tenant, installation, entity)
-- tracks the sync cursor, a generic schedule, last full/incremental timestamps,
-- the last error, and a consecutive-failure counter. MIS-006 M-01 ships only the
-- skeleton (table + scheduler loop + cursor read/write); heavy ML sync jobs land
-- in later missions (M-03 xlsx / M-04 sankhya connect the real job body).
--
-- Cadence-agnostic (D6): `schedule` is a generic JSONB descriptor (e.g.
-- {"kind":"interval","every_seconds":900} or {"kind":"cron","expr":"..."}); no
-- cadence is hardcoded as a fixed business rule here or in application code.
--
-- `entity` is a SEMANTIC enum (products|listings|orders|market|tariffs) validated
-- in the application layer, NOT a DB CHECK — the DB accepts any string so a future
-- entity never requires a migration (F-01).
--
-- ADR-17: cursor/schedule/last_error are nullable JSONB. Absence is a NULL honest
-- unknown, never a fabricated `{}`. consecutive_failures defaults to 0 — an honest
-- count of "zero failures so far", the only lifecycle default on this table.
CREATE TABLE IF NOT EXISTS sync_state (
    tenant_id            text        NOT NULL,
    installation_id      text        NOT NULL,
    entity               text        NOT NULL,
    cursor               jsonb,
    schedule             jsonb,
    last_full_sync_at    timestamptz,
    last_incremental_at  timestamptz,
    last_error           jsonb,
    consecutive_failures integer     NOT NULL DEFAULT 0,
    PRIMARY KEY (tenant_id, installation_id, entity)
);
