-- TariffDefaults (CHIP-T1 Slice A): one pricing-tariff-defaults row per
-- tenant/installation. comissao_classico_pct/comissao_premium_pct are seeded
-- by DB column DEFAULT (13.00/16.00) — the row is materialized on first read
-- via INSERT ... ON CONFLICT DO NOTHING, never a Go literal. frete_estimativa_amount
-- is NULLABLE — NULL means unknown, never 0 (ADR-17); frete_policy tracks
-- whether an estimate is known ('estimativa') or absent ('sem_dados').
-- installation_id '' is the default-installation sentinel for the current
-- single-installation demo; the (tenant_id, installation_id) key already
-- makes this schema multi-installation-ready.
CREATE TABLE IF NOT EXISTS pricing_tariff_defaults (
    tenant_id             text NOT NULL,
    installation_id       text NOT NULL DEFAULT '',
    comissao_classico_pct numeric NOT NULL DEFAULT 13.00 CHECK (comissao_classico_pct >= 0 AND comissao_classico_pct <= 100),
    comissao_premium_pct  numeric NOT NULL DEFAULT 16.00 CHECK (comissao_premium_pct  >= 0 AND comissao_premium_pct  <= 100),
    frete_estimativa_amount numeric CHECK (frete_estimativa_amount IS NULL OR frete_estimativa_amount >= 0),
    frete_policy          text NOT NULL DEFAULT 'sem_dados' CHECK (frete_policy IN ('estimativa','sem_dados')),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, installation_id)
);
