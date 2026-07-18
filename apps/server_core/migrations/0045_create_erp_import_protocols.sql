CREATE TABLE IF NOT EXISTS erp_import_protocols (
    id UUID PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    file_sha256 TEXT NOT NULL,
    protocol TEXT NOT NULL,
    source TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    accepted_count INTEGER NOT NULL,
    rejected_count INTEGER NOT NULL,
    warning_count INTEGER NOT NULL,
    CONSTRAINT erp_import_protocols_file_sha256_key UNIQUE (tenant_id, file_sha256),
    CONSTRAINT erp_import_protocols_protocol_key UNIQUE (tenant_id, protocol),
    CONSTRAINT erp_import_protocols_tenant_id_key UNIQUE (tenant_id, id),
    CONSTRAINT erp_import_protocols_protocol_check CHECK (protocol ~ '^#[0-9]{3,}-E$'),
    CONSTRAINT erp_import_protocols_source_check CHECK (source = 'xlsx'),
    CONSTRAINT erp_import_protocols_status_check CHECK (status IN ('COMPLETED', 'REJECTED'))
);

