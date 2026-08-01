# Sol unavailable — P7 r03

```yaml
phase: P7
round: r03
date: 2026-08-01
reason: quota Codex esgotada (hit 2026-07-25, re-hit 2026-07-31); reset 2026-08-05
claude_side_verdict: "Ready" (p7-claude-readiness-r03.md)
status_set: blocked
```

- Regra failure/skip (mission-planning skill): sem resultado Sol válido ⇒ persistir este
  artefato, `status: blocked`, reter `planning_phase: readiness`, escalar. NÃO há
  substituição Claude p/ o veredito CONJUNTO e NÃO há skip p/ `planned`.
- Waiver do operador (ratificado, ver memória/handoff): crews Claude frias substituem Sol
  na frente p/ destravar; Sol retroativo P3 (medium, `p3-input-r01`), P5 (medium,
  `p5-input-r09.sha256`) e P7 (HIGH, árvore pós-r03 — congelar `p7-input-r04.sha256`
  antes do despacho) TODOS obrigatórios antes de `status: planned`.
- Pré-condição extra p/ o P7 Sol: shell funcional (sessão atual com wedge B-8 total —
  sem execução de comando; freeze do manifesto exige sha256).
- Escalação: registrada no Handoff da missão + resposta ao operador nesta sessão.

## SUPERSEDIDO — decisão do operador 2026-08-01

Verbatim: **"Eu aprovo não vai ter Sol"** (resposta direta à escalação acima, mesma data).

- Dispensa TOTAL dos touchpoints Sol desta missão (P3/P5/P7 retroativos INCLUSOS) —
  supersede o waiver de 2026-07-31 que exigia retroativo antes de `planned`.
- Autoridade: dono da missão; a regra failure/skip do skill escala exatamente p/ essa
  decisão — decisão tomada, escalação DESCARREGADA.
- Consequência: veredito conjunto do P7 = lado Claude (Ready, r03) + dispensa registrada
  ⇒ `status: planned`. Ver `readiness-review.md` (fold conjunto).
