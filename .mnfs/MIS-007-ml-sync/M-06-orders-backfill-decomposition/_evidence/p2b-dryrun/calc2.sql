WITH prod AS (
    SELECT * FROM (VALUES
        (15956, 311, 5,  91.57::numeric,  15.1860::numeric,  24.8740::numeric),
        (20303,   0, 0,   5.84,            0.0000,            0.0000),
        (20322,   0, 0,   9.89,            0.0000,            0.0000),
        (22467, 122, 5, 691.13,           63.5800,           67.0382),
        (39563, 122, 5, 138.85,           19.5386,           35.1243),
        (39587, 122, 5, 409.68,           52.8736,           95.0500),
        (41912, 122, 0, 154.53,           22.5754,           37.6944),
        (42194, 122, 5, 179.02,           15.7450,           32.9538)
    ) AS t(codprod, grupo_icms, origprod, custo_unit, s_unit, r_unit)
),
codtrib AS (
    SELECT * FROM (VALUES
        ('BA',   0, 60), ('BA', 122,  0),
        ('ES',   0,  0), ('ES', 122,  0), ('ES', 311, 0),
        ('MA',   0, 60), ('MA', 122,  0),
        ('MG', 122, 60), ('MG', 311, 60),
        ('PR',   0, 60), ('PR', 122,  0),
        ('PE',   0, 60), ('PE', 122, 60), ('PE', 311, 0),
        ('RJ',   0,  0), ('RJ', 122,  0),
        ('RS',   0,  0), ('RS', 122,  0),
        ('SP',   0,  0), ('SP', 122,  0), ('SP', 311, 0)
    ) AS t(uf, grupo_icms, codtrib)
),
aliq_legal AS (
    SELECT * FROM (VALUES
        ('AC',19.0::numeric),('AL',21.5),('AM',20.0),('AP',18.0),('BA',20.5),
        ('CE',20.0),('DF',20.0),('ES',17.0),('GO',19.0),('MA',23.0),
        ('MT',17.0),('MS',17.0),('MG',18.0),('PA',19.0),('PB',20.0),
        ('PR',19.5),('PE',20.5),('PI',22.5),('RN',20.0),('RS',17.0),
        ('RJ',22.0),('RO',19.5),('RR',20.0),('SC',17.0),('SP',18.0),
        ('SE',20.0),('TO',20.0)
    ) AS t(uf, a_int)
),
base AS (
    SELECT o.provider_order_id                          AS pedido,
           s.dest_state                                 AS uf,
           i.internal_product_id                        AS codprod,
           i.quantity                                   AS qtd,
           (i.unit_price * i.quantity)::numeric(14,2)   AS receita,
           (i.sale_fee_amount * i.quantity)::numeric(14,2) AS comissao,
           coalesce(s.cost_seller, 0)::numeric          AS frete,
           p.grupo_icms, p.origprod,
           (p.custo_unit * i.quantity)::numeric         AS custo,
           (p.s_unit     * i.quantity)::numeric         AS s_val,
           (p.r_unit     * i.quantity)::numeric         AS r_val,
           ct.codtrib,
           al.a_int
    FROM orders_marketplace_orders o
    JOIN orders_marketplace_order_items i USING (provider_order_id)
    LEFT JOIN order_shipments s ON s.provider_order_id = o.provider_order_id
    LEFT JOIN prod p       ON p.codprod = i.internal_product_id
    LEFT JOIN codtrib ct   ON ct.uf = s.dest_state AND ct.grupo_icms = p.grupo_icms
    LEFT JOIN aliq_legal al ON al.uf = s.dest_state
),
calc AS (
    SELECT b.*,
           CASE WHEN b.origprod IN (1,2,3,8) THEN 4.0
                WHEN b.uf IN ('SP','RJ','PR','SC','RS') THEN 12.0
                ELSE 7.0 END AS a_inter,
           CASE
             WHEN b.codtrib = 60 THEN 0::numeric
             WHEN b.codtrib = 0 AND b.uf = 'MG' THEN 0.18
             WHEN b.codtrib = 0 AND b.a_int IS NOT NULL THEN
                  (b.a_int/100) * (1 - (CASE WHEN b.origprod IN (1,2,3,8) THEN 0.04
                                             WHEN b.uf IN ('SP','RJ','PR','SC','RS') THEN 0.12
                                             ELSE 0.07 END))
                  / (1 - b.a_int/100)
           END AS a
    FROM base b
),
fin AS (
    SELECT c.*,
           round(c.receita * c.a, 2)                                    AS icms_total,
           greatest(0, c.receita * (1 - c.a) - c.s_val)                 AS base_pc,
           round(0.0925 * greatest(0, c.receita * (1 - c.a) - c.s_val), 2) AS piscofins,
           CASE WHEN c.uf = 'MG' THEN 0::numeric ELSE c.r_val END       AS r_apl,
           round(c.receita * 0.04, 2)                                   AS imposto_hoje
    FROM calc c
),
res AS (
    SELECT f.*,
           round(f.receita - f.comissao - f.frete - f.custo - f.imposto_hoje, 2) AS margem_hoje,
           CASE WHEN f.a IS NULL OR f.custo IS NULL THEN NULL
                ELSE round(f.receita - f.comissao - f.frete - f.custo
                           - f.icms_total - f.piscofins + f.r_apl, 2) END        AS margem_nova
    FROM fin f
)
SELECT count(*)                                            AS pedidos,
       count(*) FILTER (WHERE margem_nova IS NULL)         AS nao_calculavel,
       count(*) FILTER (WHERE margem_hoje > 0)             AS positivos_hoje,
       count(*) FILTER (WHERE margem_nova > 0)             AS positivos_novo,
       count(*) FILTER (WHERE margem_hoje > 0 AND margem_nova < 0) AS viram_negativo,
       round(sum(margem_hoje) FILTER (WHERE margem_nova IS NOT NULL), 2) AS soma_hoje,
       round(sum(margem_nova), 2)                          AS soma_nova
FROM res;
