-- 0076: canonical product mirror + active-source config (MIS-006 M-02, bloco B).
--
-- products_mirror is the current-state canonical product (E2.1). It replaces the
-- snapshot-rescan read model: adapters (xlsx M-03, sankhya M-04) upsert-merge
-- rows here keyed by (tenant_id, codigo_produto). erp_import_products stays as
-- immutable per-protocol history — this table is NOT a rename of it.
--
-- Honest-unknown (ADR-17): custo / preco_venda / estoque_total are NUMERIC and
-- nullable with NO default. A missing value stays NULL, never a fabricated 0.
-- (grep guard M02-C3/AC-03: these three columns carry neither DEFAULT 0 nor
-- NOT NULL.)
--
-- Keep-absent (ADR-031): a rebuild never physically deletes. When a product from
-- the previous snapshot is missing from a new one, the merge logic (M-03/M-04)
-- sets absent_in_last_snapshot=true + stale_since=now(); when it reappears the
-- flag clears back to false. M-02 provides the columns + defaults; the merge
-- itself lands with the adapters.
CREATE TABLE IF NOT EXISTS products_mirror (
    tenant_id                TEXT NOT NULL,
    source                   TEXT NOT NULL,
    codigo_produto           TEXT NOT NULL,
    descricao                TEXT,
    referencia               TEXT,
    ean                      TEXT,
    marca                    TEXT,
    grupo_codigo             TEXT,
    grupo_descricao          TEXT,
    ncm                      TEXT,
    custo                    NUMERIC,
    preco_venda              NUMERIC,
    estoque_total            NUMERIC,
    protocol_id              UUID,
    absent_in_last_snapshot  BOOLEAN NOT NULL DEFAULT false,
    stale_since              TIMESTAMPTZ,
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, codigo_produto),
    CONSTRAINT products_mirror_source_check
        CHECK (source IN ('sankhya', 'xlsx', 'catalogo_cliente'))
);

-- Per-location stock breakdown for a mirror product (E2.1). Logical FK to
-- products_mirror via (tenant_id, codigo_produto) — intentionally NOT a physical
-- FK so mirror rebuilds and location loads stay independent.
CREATE TABLE IF NOT EXISTS products_mirror_stock_locations (
    tenant_id        TEXT NOT NULL,
    codigo_produto   TEXT NOT NULL,
    local_codigo     TEXT NOT NULL,
    local_descricao  TEXT,
    quantidade       NUMERIC,
    PRIMARY KEY (tenant_id, codigo_produto, local_codigo)
);

-- active_source is the per-tenant data-source config (E9.0). It replaces the
-- MC_ERP_SOURCE boot env var: the active source is resolved per request from
-- this table, fail-closed (no row → ErrUnknownActiveSource, never a silent
-- default). source_kind classifies the shape (upload_snapshot | live_read_through).
CREATE TABLE IF NOT EXISTS active_source (
    tenant_id      TEXT PRIMARY KEY,
    active_source  TEXT NOT NULL,
    source_kind    TEXT NOT NULL,
    set_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    set_by         TEXT,
    CONSTRAINT active_source_value_check
        CHECK (active_source IN ('sankhya', 'xlsx', 'catalogo_cliente')),
    CONSTRAINT active_source_kind_check
        CHECK (source_kind IN ('upload_snapshot', 'live_read_through'))
);

-- Seed the default tenant to the xlsx (upload_snapshot) source so the existing
-- single-tenant deployment keeps resolving to real ERP data after MC_ERP_SOURCE
-- is removed. This is explicit config, not a code fallback — other tenants
-- without a row still fail closed.
INSERT INTO active_source (tenant_id, active_source, source_kind, set_by)
VALUES ('tenant_default', 'xlsx', 'upload_snapshot', 'migration_0076')
ON CONFLICT (tenant_id) DO NOTHING;
