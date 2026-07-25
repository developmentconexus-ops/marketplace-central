-- 0078: the product mirror is keyed BY SOURCE (MIS-006 hardening, hub-owned).
--
-- 0076 keyed products_mirror on (tenant_id, codigo_produto) alone, with source
-- as a plain column. That makes the mirror hold at most ONE row per product
-- across ALL sources: importing the same CODPROD from a second source does not
-- add a row, it HIJACKS the existing one and flips its source. Two consequences,
-- both observed live (D-121): switching the active source silently emptied or
-- corrupted screens, and a source's snapshot could be destroyed by an unrelated
-- upload. The active-source switch (E9.0) is only honest if each source owns its
-- own rows, so the key must carry the source.
--
-- The same defect applies to the per-location stock child table, which had no
-- source column at all: locations from a live Sankhya sync and from an xlsx
-- upload collided on (tenant_id, codigo_produto, local_codigo).
--
-- Existing rows keep their current source, so this migration is value-preserving:
-- no row is dropped, no source is rewritten, nothing is fabricated (ADR-17).

ALTER TABLE products_mirror
    DROP CONSTRAINT products_mirror_pkey;

ALTER TABLE products_mirror
    ADD CONSTRAINT products_mirror_pkey
    PRIMARY KEY (tenant_id, source, codigo_produto);

-- Reads that resolve a product by code within the active source, and the
-- keyset-paginated catalog page, both order by codigo_produto inside one
-- source; the primary key above already serves them. This index serves the
-- reverse direction: "which sources carry this product" (source reconciliation).
CREATE INDEX IF NOT EXISTS products_mirror_code_idx
    ON products_mirror (tenant_id, codigo_produto);

-- Location rows inherit the source of the mirror row they belong to. Rows that
-- predate this migration were all written by the sources currently present in
-- products_mirror, so the join resolves them; any orphan location (no mirror row)
-- is deleted rather than assigned an invented source.
ALTER TABLE products_mirror_stock_locations
    ADD COLUMN IF NOT EXISTS source TEXT;

UPDATE products_mirror_stock_locations AS locations
SET source = mirror.source
FROM products_mirror AS mirror
WHERE mirror.tenant_id = locations.tenant_id
  AND mirror.codigo_produto = locations.codigo_produto
  AND locations.source IS NULL;

DELETE FROM products_mirror_stock_locations WHERE source IS NULL;

ALTER TABLE products_mirror_stock_locations
    ALTER COLUMN source SET NOT NULL;

ALTER TABLE products_mirror_stock_locations
    DROP CONSTRAINT products_mirror_stock_locations_pkey;

ALTER TABLE products_mirror_stock_locations
    ADD CONSTRAINT products_mirror_stock_locations_pkey
    PRIMARY KEY (tenant_id, source, codigo_produto, local_codigo);

ALTER TABLE products_mirror_stock_locations
    ADD CONSTRAINT products_mirror_stock_locations_source_check
    CHECK (source IN ('sankhya', 'xlsx', 'catalogo_cliente'));

-- The mirror stored only estoque_total (physical MINUS reserved), computed with
-- a rule that required BOTH components: a product whose spreadsheet row carried
-- physical stock but no reserved value collapsed to NULL and rendered as "no
-- stock" on every screen (observed live: 1845 products with physical stock, 164
-- of them shown as unknown). Storing the two components keeps the fact instead
-- of destroying it, and lets a reader show what is actually known.
ALTER TABLE products_mirror
    ADD COLUMN IF NOT EXISTS estoque_fisico NUMERIC;

ALTER TABLE products_mirror
    ADD COLUMN IF NOT EXISTS estoque_reservado NUMERIC;

-- The per-protocol snapshot history did not persist the ERP sale price, so a
-- mirror rebuild from a stored protocol could only ever restore a NULL price
-- even when the imported spreadsheet carried one. Nullable with no default:
-- an import without a price column keeps NULL (honest-unknown, ADR-17), and
-- every protocol imported before this column existed stays NULL rather than
-- being back-filled with an invented value.
ALTER TABLE erp_import_products
    ADD COLUMN IF NOT EXISTS preco_venda NUMERIC;

ALTER TABLE erp_import_products
    ADD CONSTRAINT erp_import_products_preco_venda_check
    CHECK (preco_venda IS NULL OR preco_venda > 0);

-- erp_import_protocols was keyed on id alone while every other constraint,
-- every foreign key and every query is tenant-scoped. A protocol id therefore
-- had to be unique ACROSS tenants, which is not a property the domain has: two
-- tenants importing concurrently could collide on a primary key that models
-- nothing. Key it the way it is actually used, (tenant_id, id); the redundant
-- unique constraint that already carried that pair is dropped with it.
ALTER TABLE erp_import_products DROP CONSTRAINT erp_import_products_protocol_fk;
ALTER TABLE erp_import_issues DROP CONSTRAINT erp_import_issues_protocol_fk;

ALTER TABLE erp_import_protocols DROP CONSTRAINT erp_import_protocols_pkey;
ALTER TABLE erp_import_protocols DROP CONSTRAINT erp_import_protocols_tenant_id_key;
ALTER TABLE erp_import_protocols ADD CONSTRAINT erp_import_protocols_pkey PRIMARY KEY (tenant_id, id);

ALTER TABLE erp_import_products
    ADD CONSTRAINT erp_import_products_protocol_fk
    FOREIGN KEY (tenant_id, protocol_id) REFERENCES erp_import_protocols (tenant_id, id);
ALTER TABLE erp_import_issues
    ADD CONSTRAINT erp_import_issues_protocol_fk
    FOREIGN KEY (tenant_id, protocol_id) REFERENCES erp_import_protocols (tenant_id, id);
