# QA live-drive — CHIP-VENDAVEL pós-merge (hub, 2026-07-31)

Persona fresca (localStorage/sessionStorage limpos), banco velho (A-26: baseline é o estado
que existe), dev stack rebuildado da main mergeada (@3ef01f65 + fix @f9d6a91). Frontend
:5174, backend :8080, postgres :5435 (`marketplace-central-postgres-1`).

## Defeito achado e consertado DURANTE o QA (before/after real, não FAIL)

`GET /catalog/products/counts` → **503** nas telas /integracoes e /catalogo.
Root cause: ORA-00937 — `COUNT(*)` ao lado de subquery escalar no mesmo select list sem
GROUP BY. Latente desde 5adeeb56 (pré-takeover); os testes de unidade pinam SUBSTRING do SQL
e não enxergam validade de statement — só o live drive pega. Fix @f9d6a91: dois scalar
subqueries FROM DUAL. After: `code=200 {"sellable_count":2945,"total_count":10549}`.
Lição de classe registrada em memória (substring-pin não prova statement válido).

## Veredito por critério

| VC | Veredito | Evidência medida |
|----|----------|------------------|
| VC-1 | PASS | Toggle `only_ecommerce_eligible` via UI → `active_source` DB t/t/t; localStorage+sessionStorage limpos + reload → toggles [t,t,t] sobreviveram; `localStorage.length===0`. Restaurado t/t/f. |
| VC-2 | PASS | Card /integracoes "Resultado: 2945 de 10549" = chip /catalogo "Vendáveis 2945 de 10549" = endpoint = SQL live (concordância tela↔SQL; ≈2.923 ratificado 2026-07-29 + drift diário). 3º toggle on → 1341 (≈1.329 DR-3). |
| VC-3 | PASS | /catalogo abre filtrado por default; "Ver todos" → `include_all=true`, primeira página com IDs não-vendáveis (1016/2592/4259/5306) ausentes do modo filtrado; botão vira "Ver sortimento filtrado". `only_em_estoque` off → chip 10017/10549 e 13 linhas com badge "Sem estoque" permanecem na lista. Defaults restaurados (t/t/f, 2945). |
| VC-4 | PASS (lane) | Geração de vínculo ignora não-vendável — teste de integração verde na escada P5 pós-merge (status=passed, run_id 2ddf6e7668ee49fba891ccb770d1fc08). |
| VC-5 | PASS | Sync live pós-merge contra Oracle populou as colunas novas do espelho: ver seção LIVE-VERIFIED abaixo. Pin do predicado Q4 = testes de unidade verdes na escada. |
| VC-6 | PASS (lane) | Colunas xlsx opcionais — vitest/go lanes verdes na escada P5. |
| VC-7 | PASS | Zero erro de console nas 3 telas (/integracoes, /catalogo, /vinculos); network sweep só 200/304 pós-fix (único não-2xx da sessão foi o counts 503 pré-fix, acima); tsc teto 12 pré-existentes; vitest web 566 + sdk-runtime 78 + feature-products 13. |

## VC-5 LIVE-VERIFIED (tick do scheduler 2026-07-31 17:10 UTC)

Scheduler de products (15min, composition root.go:671) tickou ~15min após o rebuild do
backend e o sync live contra o Oracle populou as colunas novas do chip:

- `sync_state` cursor: `{"source":"sankhya","processed":10549,"completed_at":"2026-07-31T17:10:38Z"}`
- `products_mirror` source=sankhya: 10549 rows, `COUNT(usoprod)=10549`, `COUNT(ad_ecommerce)=3604`, `MAX(updated_at)=17:10:34Z`
- Distribuição usoprod: R=10017 (bate com o toggle only_em_estoque off na tela), C=350, I=172, …
- **Concordância espelho↔Oracle**: `usoprod='R' AND COALESCE(estoque_total,0)>0` no espelho
  = **2945** — idêntico ao count live do Oracle (2945/10549). O predicado Q4 do sync
  (whitelist CODEMP 1,2 / CODLOCAL 10101,10102 embutida no estoque_total) reproduz o corte
  vendável byte-a-byte.

Observação (não é regressão do chip): `estoque_fisico`/`estoque_reservado` nunca foram
populadas por sync nenhum (0 antes e depois); o disponível vive em `estoque_total`.

## Estado final

Toggles restaurados aos defaults ratificados (t/t/f), counts 2945/10549. Nenhuma escrita
viva em ML. Worktree chip-vendavel: pg-session container derrubado (mpc-pg-session-3eee515d),
worktree removido, branch `worktree-chip-vendavel` deletado com `-d` (mergeado).
