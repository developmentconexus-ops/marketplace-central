# MIS-007 — P7 Readiness Review (fold conjunto)

```yaml
mission: MIS-007
rounds_consumed: 3/3
final_round: r03
manifest_final: NENHUM no r03 (desvio registrado — shell wedge B-8; último
  content-addressed = planning-reviews/p7-input-r02.sha256, digest
  e238f447823ab01eaa1780deb84d0db3280488ab4eda864033d51134f8a4ebaa; deltas r02→r03 =
  17 reparos verificados um a um em planning-reviews/p7-verify-repairs-r03.md)
claude_artifact: planning-reviews/p7-claude-readiness-r03.md
claude_verdict: Ready
sol_artifact: planning-reviews/sol-unavailable-p7-r03.md (supersedido — dispensa)
sol_verdict: DISPENSADO — decisão do operador 2026-08-01, verbatim "Eu aprovo não vai
  ter Sol" (supersede waiver 2026-07-31 que exigia retroativo)
joint_verdict: Ready (lado Claude Ready + lado Sol dispensado por autoridade do dono)
status_set: planned
date: 2026-08-01
```

## Trilha das rodadas

| Round | Manifesto | Crew | Verdict | Blocking |
|---|---|---|---|---|
| r01 | `p7-input-r01.sha256` | fria ×5 | Needs revision | B-1..B-9 (9 loci) — reparados |
| r02 | `p7-input-r02.sha256` (e238f447…, re-verificado 65/65) | fria ×5 | Needs revision | ★2-A (PK `listing_variations` 5-col vs 4-col), ★2-B (slug `mercadolivre` vs `mercado_livre`) — seat 5 double-pass; reparados + 12 advisories |
| r03 | nenhum (B-8; desvio registrado) | enxuta (ordem do operador: conferência 17/17 + 1 assento ★2 sonnet) | Ready | nenhum |

## União de blocking findings

Todas descarregadas: B-1..B-9 (r01), ★2-A/★2-B (r02). Verificação final:
`p7-verify-repairs-r03.md` (17/17 PRESENT) + `p7-seat-star2-r03.md` (★2 PASS, sweep
adversarial tree-wide sem divergência nova).

## Advisories que permanecem (nunca mudam veredito)

- **A-9 (r02, seats 4+5):** `notifications_inbox` sem retention/prune + amplificação de
  refetch (forja plausível custa 1 GET autenticado do token-bucket M-01). SEM valor
  decidido — dívida declarada; decisão do operador entra quando quiser (rows inertes,
  M08-C3).

## Disposição de reparo

r01: 9 loci + advisories baratas. r02: 2 loci ★2 + 12 advisories
(A-1/2/4/5/6/7/8/10/11/12/13/14). Todos locais, dentro da repair-authority (nenhum mudou
outcome/escopo/arquitetura/contrato ratificado; ★2-A restaurou a leitura RATIFICADA da
A-13; ★2-B alinhou ao provider_code já decidido).

## Desvios de procedimento (ambos com autoridade registrada)

1. Rodada r03 enxuta — ordem do operador (custo), verbatim em
   `p7-claude-readiness-r03.md`.
2. Sem manifesto r03 — shell da sessão morto (classe B-8); imutabilidade por janela
   única + lista prescrita substituíram o content-addressing.
3. Sol dispensado — decisão do operador (acima).

## Pós-planned (housekeeping, não bloqueia)

- Remover scratch `freeze_tmp.py`/`freeze_srv.py` (raiz da missão) quando houver shell.
- Commit dos artefatos de planejamento (shell morto nesta sessão; nunca push sem
  permissão do operador).
- PII scrub `docs/design/evidence/ml-api/` antes de qualquer commit desses dumps.
