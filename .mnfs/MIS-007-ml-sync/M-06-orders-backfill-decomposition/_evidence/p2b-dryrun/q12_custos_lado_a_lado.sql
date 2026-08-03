WITH ult AS (
  SELECT c.CODPROD, c.CODEMP, c.DTATUAL,
         c.CUSMED, c.CUSMEDICM, c.CUSSEMICM, c.CUSVARIAVEL, c.CUSREP, c.CUSGER,
         ROW_NUMBER() OVER (PARTITION BY c.CODPROD, c.CODEMP ORDER BY c.DTATUAL DESC) AS RN
  FROM METALPRD.TGFCUS c
  WHERE c.CODPROD IN (15956, 20303, 20322, 22467, 39563, 39587, 41912, 42194)
    AND c.CODEMP = 1
    AND c.DTATUAL <= SYSDATE
)
SELECT u.CODPROD,
       SUBSTR(p.DESCRPROD, 1, 26)              AS PRODUTO,
       p.TIPSUBST                              AS TSUB,
       p.ORIGPROD                              AS ORIG,
       p.GRUPOICMS                             AS GRUPO,
       TO_CHAR(u.DTATUAL, 'YYYY-MM-DD')        AS DTATUAL,
       ROUND(u.CUSMED, 2)                      AS CUSMED,
       ROUND(u.CUSMEDICM, 2)                   AS CUSMEDICM,
       ROUND(u.CUSSEMICM, 2)                   AS CUSSEMICM,
       ROUND(u.CUSVARIAVEL, 2)                 AS CUSVARIAVEL,
       ROUND(u.CUSREP, 2)                      AS CUSREP,
       ROUND(u.CUSGER, 2)                      AS CUSGER
FROM ult u
JOIN METALPRD.TGFPRO p ON p.CODPROD = u.CODPROD
WHERE u.RN = 1
ORDER BY u.CODPROD
