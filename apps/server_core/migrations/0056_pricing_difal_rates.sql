-- DIFAL rates (IC-04 / R-04 seed): one row per (tenant, UF). efetivo_pct is the
-- derived max(interna - interestadual, 0); it is stored for read convenience but
-- the authoritative value is recomputed by domain.DifalForUF on read (S6), which
-- also applies an active override.
--
-- The override is NOT a 4th table (grant names only 3): it is the nullable
-- override_* columns on this table. An operator override at 2dp (e.g. 19.05)
-- must survive as-entered — override_interna_pct is numeric (no scale cap).
-- override_actor + override_updated_at record the audited change (IC-04: an
-- override persists only when |Δ| > 0.049pp, enforced app-side in S7).
CREATE TABLE IF NOT EXISTS pricing_difal_rates (
    tenant_id            text NOT NULL,
    uf                   text NOT NULL CHECK (char_length(uf) = 2),
    interna_pct          numeric NOT NULL,
    interestadual_pct    numeric NOT NULL,
    efetivo_pct          numeric NOT NULL,
    origem_versao        text NOT NULL,
    override_interna_pct numeric,
    override_updated_at  timestamptz,
    override_actor       text,
    PRIMARY KEY (tenant_id, uf),
    CONSTRAINT pricing_difal_rates_override_paired CHECK (
        (override_interna_pct IS NULL AND override_updated_at IS NULL AND override_actor IS NULL)
        OR (override_interna_pct IS NOT NULL AND override_updated_at IS NOT NULL AND override_actor IS NOT NULL)
    )
);
