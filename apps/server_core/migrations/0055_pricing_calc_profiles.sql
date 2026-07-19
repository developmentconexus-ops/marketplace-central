-- CalcProfile (IC-04): one pricing-calculation profile per tenant. The initial
-- state (SIMPLES 4%, origem 'default') is created on first read by the S6 repo,
-- not seeded here. aliquota_pct is bounded 0..35 (transport maps a violation to
-- INVALID_RATE). tarifa_full is NULLABLE — NULL means unset/unknown, never 0
-- (ADR-17); a `full`-modalidade decomposition with a NULL tarifa_full must
-- propagate an unknown component, not compute against zero.
CREATE TABLE IF NOT EXISTS pricing_calc_profiles (
    tenant_id          text NOT NULL,
    regime             text NOT NULL CHECK (regime IN ('SIMPLES', 'PRESUMIDO')),
    aliquota_pct       numeric NOT NULL CHECK (aliquota_pct >= 0 AND aliquota_pct <= 35),
    limiar_verde_pct   numeric NOT NULL,
    limiar_amarelo_pct numeric NOT NULL,
    tarifa_full        numeric,
    difal_enabled      boolean NOT NULL DEFAULT false,
    difal_destino_uf   text CHECK (difal_destino_uf IS NULL OR char_length(difal_destino_uf) = 2),
    origem             text NOT NULL CHECK (origem IN ('default', 'operator')),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id)
);
