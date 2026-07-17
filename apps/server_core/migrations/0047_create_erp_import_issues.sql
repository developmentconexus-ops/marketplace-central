CREATE TABLE IF NOT EXISTS erp_import_issues (
    tenant_id TEXT NOT NULL,
    protocol_id UUID NOT NULL,
    row_number INTEGER NOT NULL,
    column_name TEXT,
    kind TEXT NOT NULL,
    code TEXT NOT NULL,
    detail TEXT,
    offending_value TEXT,
    CONSTRAINT erp_import_issues_protocol_fk FOREIGN KEY (tenant_id, protocol_id) REFERENCES erp_import_protocols(tenant_id, id),
    CONSTRAINT erp_import_issues_kind_check CHECK (kind IN ('REJECTION', 'WARNING')),
    CONSTRAINT erp_import_issues_code_check CHECK (code IN (
        'EMPTY_CODPROD',
        'DUPLICATE_CODPROD',
        'EMPTY_DESCRPROD',
        'INVALID_CUSTO',
        'INVALID_ESTOQUE',
        'INVALID_EAN',
        'INVALID_NCM',
        'MISSING_REQUIRED_COLUMN',
        'INVALID_FILE',
        'IMPORT_IN_PROGRESS',
        'DUPLICATE_FILE'
    ))
);

CREATE INDEX erp_import_protocols_latest_idx
    ON erp_import_protocols (tenant_id, status, imported_at DESC);

CREATE INDEX erp_import_issues_protocol_idx
    ON erp_import_issues (tenant_id, protocol_id);
