CREATE TABLE IF NOT EXISTS erp_import_products (
    tenant_id TEXT NOT NULL,
    protocol_id UUID NOT NULL,
    codprod TEXT NOT NULL,
    descrprod TEXT NOT NULL,
    custo NUMERIC NOT NULL,
    stock_physical INTEGER NOT NULL,
    stock_reserved INTEGER,
    ean TEXT,
    refforn TEXT,
    marca TEXT,
    ncm TEXT,
    PRIMARY KEY (tenant_id, protocol_id, codprod),
    CONSTRAINT erp_import_products_protocol_fk FOREIGN KEY (tenant_id, protocol_id) REFERENCES erp_import_protocols(tenant_id, id),
    CONSTRAINT erp_import_products_custo_check CHECK (custo > 0),
    CONSTRAINT erp_import_products_stock_physical_check CHECK (stock_physical >= 0),
    CONSTRAINT erp_import_products_stock_reserved_check CHECK (stock_reserved IS NULL OR stock_reserved >= 0)
);

