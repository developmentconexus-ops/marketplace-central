-- Reconciliacao: formula P2.b (matriz TGFICM vigente) x notas realmente emitidas.
-- Notas de q06_notas.sql, TOP 306 apenas (a 313 e a mesma operacao duplicada
-- pelo TGFVAR e inflaria a amostra), CLASSIFICMS='C' (PF nao-contribuinte = venda ML).

WITH nota(nunota, uf, codprod, grupo, base, icms_nota, difal_nota, fcp_nota) AS (VALUES
    (894164, 'RJ', 15956, 311, 169.99, 20.40,  11.90,  3.40),
    (893184, 'SP', 15956, 311, 169.90, 20.39,  10.19,  0.00),
    (893459, 'SP', 15956, 311, 169.90, 20.39,  10.19,  0.00),
    (894077, 'SP', 15956, 311, 169.99, 20.40,  10.20,  0.00),
    (894549, 'SP', 15956, 311, 339.98, 40.80,  20.40,  0.00),
    (894733, 'SP', 15956, 311, 169.99, 20.40,  10.20,  0.00),
    (892337, 'MG', 22467, 122,   0.00,  0.00,   0.00,  0.00),
    (892946, 'SP', 22467, 122,1489.90,178.79,  89.39,  0.00),
    (895436, 'BA', 39563, 122, 189.90, 13.29,  25.64,  0.00),
    (895435, 'PR', 39563, 122, 189.90, 22.79,  11.39,  0.00),
    (895190, 'PR', 39587, 122,1459.80,175.18,  87.59,  0.00),
    (893516, 'RS', 39587, 122, 729.90, 87.59,  36.50,  0.00),
    (895507, 'BA', 41912, 122, 299.90, 20.99,  29.99,  0.00),
    (893649, 'RJ', 41912, 122, 299.90, 35.99,  20.99,  6.00),
    (894296, 'SP', 41912, 122, 599.80, 71.98,  35.99,  0.00),
    (894287, 'SP', 41912, 122, 299.90, 35.99,  17.99,  0.00),
    (891855, 'ES', 42194, 122, 309.90, 21.69,  30.99,  0.00),
    (891752, 'RJ', 42194, 122, 309.90, 37.19,  18.59,  6.20),
    (891753, 'RJ', 42194, 122, 309.90, 37.19,  18.59,  6.20),
    (892339, 'RJ', 42194, 122, 309.90, 37.19,  18.59,  6.20),
    (892338, 'RJ', 42194, 122,1859.40,223.13, 111.56, 37.19)
),
matriz(uf, grupo, codtrib, aliquota, aliqintdest, aliqufdest, fcp) AS (VALUES
    ('BA', 0,   60,  7.0, NULL::numeric, 20.5, NULL::numeric),
    ('BA', 122,  0,  7.0, 20.5,          20.5, NULL),
    ('ES', 0,    0,  7.0, 17.0,          19.5, NULL),
    ('ES', 122,  0,  7.0, 17.0,          19.5, NULL),
    ('ES', 311,  0,  7.0, 17.0,          19.5, NULL),
    ('MA', 0,   60,  7.0, NULL,          22.0, NULL),
    ('MA', 122,  0,  7.0, 23.0,          22.0, NULL),
    ('MG', 122, 60,  0.0, NULL,          NULL, NULL),
    ('MG', 311, 60,  0.0, NULL,          NULL, NULL),
    ('PR', 0,   60, 12.0, NULL,          19.0, NULL),
    ('PR', 122,  0, 12.0, NULL,          19.0, NULL),
    ('PE', 0,   60,  7.0, NULL,          20.5, NULL),
    ('PE', 122, 60,  7.0, NULL,          20.5, NULL),
    ('PE', 311,  0,  7.0, 18.0,          20.5, NULL),
    ('RJ', 0,    0, 12.0, 18.0,          20.0, 2.0),
    ('RJ', 122,  0, 12.0, 18.0,          20.0, 2.0),
    ('RJ', 311, NULL, NULL, NULL,        NULL, NULL),
    ('RS', 0,    0, 12.0, NULL,          17.0, NULL),
    ('RS', 122,  0, 12.0, NULL,          17.0, NULL),
    ('SP', 0,    0, 12.0, 18.0,          18.0, NULL),
    ('SP', 122,  0, 12.0, 18.0,          18.0, NULL),
    ('SP', 311,  0, 12.0, 18.0,          18.0, NULL)
),
prev AS (
    SELECT n.*, m.codtrib, m.aliquota, m.aliqintdest, m.fcp AS fcp_pct,
           CASE WHEN m.codtrib = 60 THEN 0::numeric
                WHEN m.codtrib = 0  THEN round(n.base * m.aliquota / 100, 2)
           END AS icms_prev,
           -- variante A: diferenca das aliquotas, arredondada uma vez
           CASE WHEN m.codtrib = 60 THEN 0::numeric
                WHEN m.codtrib = 0 AND m.aliqintdest IS NOT NULL
                     THEN round(n.base * (m.aliqintdest - m.aliquota) / 100, 2)
           END AS difal_A,
           -- variante B: total arredondado menos ICMS arredondado
           CASE WHEN m.codtrib = 60 THEN 0::numeric
                WHEN m.codtrib = 0 AND m.aliqintdest IS NOT NULL
                     THEN round(n.base * m.aliqintdest / 100, 2)
                          - round(n.base * m.aliquota / 100, 2)
           END AS difal_B,
           CASE WHEN m.codtrib = 60 THEN 0::numeric
                WHEN m.codtrib = 0  THEN round(n.base * coalesce(m.fcp, 0) / 100, 2)
           END AS fcp_prev
    FROM nota n LEFT JOIN matriz m ON m.uf = n.uf AND m.grupo = n.grupo
)
SELECT nunota, uf, codprod, grupo, base,
       icms_nota, icms_prev,
       CASE WHEN icms_prev IS NULL THEN 'DESCONHECIDO'
            WHEN icms_prev = icms_nota THEN 'ok'
            ELSE 'DIFF ' || (icms_prev - icms_nota)::text END AS v_icms,
       difal_nota, difal_A, difal_B,
       CASE WHEN difal_A IS NULL THEN 'DESCONHECIDO'
            WHEN difal_A = difal_nota THEN 'ok'
            ELSE 'DIFF ' || (difal_A - difal_nota)::text END AS v_difal_A,
       CASE WHEN difal_B IS NULL THEN 'DESCONHECIDO'
            WHEN difal_B = difal_nota THEN 'ok'
            ELSE 'DIFF ' || (difal_B - difal_nota)::text END AS v_difal_B,
       CASE WHEN fcp_prev IS NULL THEN 'DESCONHECIDO'
            WHEN fcp_prev = fcp_nota THEN 'ok'
            ELSE 'DIFF ' || (fcp_prev - fcp_nota)::text END AS v_fcp
FROM prev
ORDER BY uf, codprod, nunota;
