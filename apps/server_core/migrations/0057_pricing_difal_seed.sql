-- R-04 27-UF DIFAL seed (origem_versao 'padrao-2026') for the default tenant.
-- interna_pct is the only source datum per UF; interestadual_pct (12% for
-- {MG,PR,RJ,RS,SC,SP}, else 7%) and efetivo_pct (max(interna-interestadual, 0))
-- are DERIVED in SQL — never hand-duplicated — so the seed cannot drift from the
-- rule (mirrors domain.buildDifalSeed). Overrides start NULL. Idempotent.
--
-- Single-tenant demo scope: seeded for 'tenant_default' only. Per-tenant
-- provisioning for future tenants is handled app-side (S6 seeds on first read).
INSERT INTO pricing_difal_rates
    (tenant_id, uf, interna_pct, interestadual_pct, efetivo_pct, origem_versao)
SELECT
    'tenant_default',
    s.uf,
    s.interna_pct,
    inter.interestadual_pct,
    GREATEST(s.interna_pct - inter.interestadual_pct, 0),
    'padrao-2026'
FROM (VALUES
    ('AC', 19.0), ('AL', 20.5), ('AM', 20.0), ('AP', 18.0), ('BA', 20.5),
    ('CE', 20.0), ('DF', 20.0), ('ES', 17.0), ('GO', 19.0), ('MA', 23.0),
    ('MG', 18.0), ('MS', 17.0), ('MT', 17.0), ('PA', 19.0), ('PB', 20.0),
    ('PE', 20.5), ('PI', 22.5), ('PR', 19.5), ('RJ', 20.0), ('RN', 20.0),
    ('RO', 19.5), ('RR', 20.0), ('RS', 17.0), ('SC', 17.0), ('SE', 19.0),
    ('SP', 18.0), ('TO', 20.0)
) AS s(uf, interna_pct)
CROSS JOIN LATERAL (
    SELECT CASE WHEN s.uf IN ('MG', 'PR', 'RJ', 'RS', 'SC', 'SP') THEN 12 ELSE 7 END
        AS interestadual_pct
) AS inter
ON CONFLICT (tenant_id, uf) DO NOTHING;
