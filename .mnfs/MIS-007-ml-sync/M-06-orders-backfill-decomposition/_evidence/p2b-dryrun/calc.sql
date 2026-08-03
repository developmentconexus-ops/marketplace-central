-- Dry-run da fórmula P2.b contra os 38 pedidos reais.
-- Matriz TGFICM inlined a partir de q04_matriz_recorte.sql (UFORIG=13/MG),
-- resolvida pela lista branca (N/0 + I/<grupo>  ou  I/<grupo> + S/-1).
-- Células ausentes NAO aparecem aqui: viram UNKNOWN por LEFT JOIN nulo.

WITH grupo AS (
    SELECT * FROM (VALUES
        (15956, 311), (20303, 0), (20322, 0), (22467, 122),
        (39563, 122), (39587, 122), (41912, 122), (42194, 122)
    ) AS t(codprod, grupo_icms)
),
matriz AS (
    SELECT * FROM (VALUES
        ('BA', 0,   60, 7.0,  NULL::numeric, NULL::numeric),
        ('BA', 122,  0, 7.0,  20.5,          NULL),
        ('ES', 0,    0, 7.0,  17.0,          NULL),
        ('ES', 122,  0, 7.0,  17.0,          NULL),
        ('ES', 311,  0, 7.0,  17.0,          NULL),
        ('MA', 0,   60, 7.0,  NULL,          NULL),
        ('MA', 122,  0, 7.0,  23.0,          NULL),
        ('MG', 122, 60, 0.0,  NULL,          NULL),
        ('MG', 311, 60, 0.0,  NULL,          NULL),
        ('PR', 0,   60, 12.0, NULL,          NULL),
        ('PR', 122,  0, 12.0, NULL,          NULL),
        ('PE', 0,   60, 7.0,  NULL,          NULL),
        ('PE', 122, 60, 7.0,  NULL,          NULL),
        ('PE', 311,  0, 7.0,  18.0,          NULL),
        ('RJ', 0,    0, 12.0, 18.0,          2.0),
        ('RJ', 122,  0, 12.0, 18.0,          2.0),
        ('RS', 0,    0, 12.0, NULL,          NULL),
        ('RS', 122,  0, 12.0, NULL,          NULL),
        ('SP', 0,    0, 12.0, 18.0,          NULL),
        ('SP', 122,  0, 12.0, 18.0,          NULL),
        ('SP', 311,  0, 12.0, 18.0,          NULL)
    ) AS t(uf, grupo_icms, codtrib, aliquota, aliqintdest, fcp)
),
base AS (
    SELECT o.provider_order_id                       AS pedido,
           o.provider_status                          AS status,
           s.dest_state                               AS uf,
           i.internal_product_id                      AS codprod,
           g.grupo_icms,
           (i.unit_price * i.quantity)::numeric(14,2) AS receita,
           (i.sale_fee_amount * i.quantity)::numeric(14,2) AS comissao,
           s.cost_seller                              AS frete,
           (p.custo * i.quantity)::numeric(14,4)      AS custo,
           m.codtrib, m.aliquota, m.aliqintdest, m.fcp
    FROM orders_marketplace_orders o
    JOIN orders_marketplace_order_items i USING (provider_order_id)
    LEFT JOIN order_shipments s ON s.provider_order_id = o.provider_order_id
    LEFT JOIN products_mirror p ON p.codigo_produto = i.internal_product_id::text
                               AND p.source = 'sankhya'
    LEFT JOIN grupo g ON g.codprod = i.internal_product_id
    LEFT JOIN matriz m ON m.uf = s.dest_state AND m.grupo_icms = g.grupo_icms
),
calc AS (
    SELECT b.*,
           -- ICMS de saida: CODTRIB 60 => 0 explicito (ST). CODTRIB 0 => V x ALIQUOTA.
           CASE WHEN b.codtrib = 60 THEN 0::numeric
                WHEN b.codtrib = 0  THEN round(b.receita * b.aliquota / 100, 2)
           END AS icms,
           -- DIFAL: 60 => 0. 0 com ALIQINTDEST => V x ALIQINTDEST - ICMS. 0 sem => NULL.
           CASE WHEN b.codtrib = 60 THEN 0::numeric
                WHEN b.codtrib = 0 AND b.aliqintdest IS NOT NULL
                     THEN round(b.receita * b.aliqintdest / 100, 2)
                          - round(b.receita * b.aliquota / 100, 2)
           END AS difal,
           CASE WHEN b.codtrib = 60 THEN 0::numeric
                WHEN b.codtrib = 0  THEN round(b.receita * coalesce(b.fcp, 0) / 100, 2)
           END AS fcp_v,
           round(b.receita * 8.62 / 100, 2) AS piscofins,
           round(b.receita * 4.0  / 100, 2) AS imposto_hoje
    FROM base b
),
final AS (
    SELECT pedido, status, uf, codprod, grupo_icms, codtrib,
           receita, comissao, frete, round(custo, 2) AS custo,
           icms, difal, fcp_v AS fcp, piscofins, imposto_hoje,
           round(receita - comissao - coalesce(frete, 0) - custo - imposto_hoje, 2)
               AS margem_hoje,
           CASE WHEN icms IS NULL OR difal IS NULL THEN NULL
                ELSE round(receita - comissao - coalesce(frete, 0) - custo
                           - icms - difal - fcp_v - piscofins, 2) END AS margem_com_pc,
           CASE WHEN icms IS NULL OR difal IS NULL THEN NULL
                ELSE round(receita - comissao - coalesce(frete, 0) - custo
                           - icms - difal - fcp_v, 2) END AS margem_sem_pc
    FROM calc
)
SELECT
    count(*)                                                        AS pedidos,
    count(*) FILTER (WHERE margem_com_pc IS NULL)                   AS nao_calculavel,
    count(*) FILTER (WHERE margem_hoje > 0)                         AS positivos_hoje,
    count(*) FILTER (WHERE margem_com_pc > 0)                       AS positivos_com_pc,
    count(*) FILTER (WHERE margem_sem_pc > 0)                       AS positivos_sem_pc,
    count(*) FILTER (WHERE margem_sem_pc > 0 AND margem_com_pc < 0) AS viram_negativo_so_pelo_pc,
    round(sum(margem_hoje)   FILTER (WHERE margem_com_pc IS NOT NULL), 2) AS soma_hoje,
    round(sum(margem_com_pc), 2)                                    AS soma_com_pc,
    round(sum(margem_sem_pc), 2)                                    AS soma_sem_pc
FROM final;
