# Token-Economy Analysis — MIS-004 wave B (2026-07-18)

Método: usage somado dos transcripts jsonl (dedupe por message.id — eventos de streaming
duplicam usage ~2,7×; números abaixo são reais por chamada API). Perfil de tool-calls por
parse do content. Preços de lista Anthropic como pesos relativos (Opus $15/$75 in/out,
cacheR 10%, cacheW 125%; Sonnet $3/$15; Haiku $1/$5).

## Medição (sessões principais, 1ª sessão de cada chip)

| Sessão | API calls | Output | Cache-read | Ctx médio/call |
|---|---|---|---|---|
| CHIP-M08 main (Opus) | 582 | 815k | 74,5M | **132k** |
| CHIP-M07 main (Opus) | 460 | 543k | 57,9M | **129k** |
| CHIP-M05 main (Opus) | 407 | 521k | 51,6M | **131k** |
| HUB v2 (Fable, 2 sessões) | 176 | 131k | 20,8M | 122k |

Sessões -b (continuações) somam ~+35%; subagents dos chips: sonnet/haiku, output menor que o main.

Custo relativo (pesos de lista, chip M-08 main): output $61 · cacheR $112 · cacheW $44.
**Cache-read ≈ 50%+ do custo; contexto médio 130k/call = sessão vive no teto.**

## Diagnóstico — onde está caro (evidência)

1. **Orquestrador Opus implementando** (contingência D-23 codex + F-quota-1 sonnet wall):
   output do main = maior share do chip; desenho original punha P3 em codex (custo Anthropic
   zero). Anti-padrão de contingência, não de arquitetura.
2. **Contexto no teto × muitos calls**: 130k × 400–580 calls. Doutrina 76KB + mission
   artifacts + board + histórico de 1.000+ turns numa sessão só. §8 (context discipline)
   existe mas não força quebra por fase.
3. **Turn count**: 431–636 tool calls/chip; Bash 148–166 (comandos 1-a-1, não em batch);
   tool_results 0,9–1,4M chars injetados (outputs sem head_limit/quiet).
4. **Re-reads**: DISPATCH-LEDGER.md relido 9× por chip; arquivos de código 3–6×.
5. **Crew scribe por resume**: filing de 1 linha de ledger custou 64k→67k tokens de subagent
   (transcript cresce a cada resume) vs ~2–3k se o hub arquiva direto. Scribe morreu de
   transcript perdido — modelo resume-forever não paga pra tarefas de 1 linha.
6. **Redundância de gates em modo Claude-only**: P4 por slice (sonnet) + P6 dual (Opus frio +
   sonnet adversarial) + P7 QA = cada linha revisada ≥4×. Evidência de campo já na doutrina:
   crew fria de 5 no M-01 achou ZERO defeitos; o live-drive achou os 5 reais. Review frio
   redundante rende pouco; QA vivo rende muito.

## O que já está certo (não mexer)

- Hub economy validada: hub 131k output vs 6,9M dos chips (bruto); decisões caras, mecânico barato.
- P7 QA live-drive: maior detector de defeito real por token gasto (5 defeitos M-01).
- Evento-grammar batched; collision matrix; grants aditivos (evitam round-trips de REQUEST).
- Haiku investigators (M-05 usou: 87 calls baratos).

## Proposta — máximo global (ordem de impacto)

| # | Mudança | Mecanismo | Economia estimada | Classe |
|---|---|---|---|---|
| 1 | **Orquestrador de chip NÃO implementa** (>~20 linhas ⇒ dispatch worker sonnet c/ context pack; orquestrador só decide/aceita) | corta output Opus (5× sonnet) e reduz calls do main | −40–50% output do chip | regra de prompt (hub authority) |
| 2 | **Sessão fresh por fase** (P2/P3 → P5/P6 → P7) c/ handoff §8 | ctx médio 130k→~60–70k; cacheW do reboot ≪ cacheR acumulado | −40% do cache-read (componente dominante) | regra de prompt |
| 3 | **P4 batched por feature quando implementer = próprio Opus** (self-review slice-a-slice é redundante c/ P6 frio; manter P4 por slice quando implementer é worker barato) | menos passes de review s/ perder gate frio | −10–15% calls | amendment §4 (operator ratifica) |
| 4 | **Higiene de turns**: Bash em batch, head_limit/quiet em outputs, investigators retornam tabela comprimida, ledger do chip em append-only (sem reler 9×) | menos calls × menos bytes injetados | −15–20% calls/results | regra de prompt |
| 5 | **Crew: hub arquiva rows de 1 linha direto; scribe só p/ lotes**; ops mantém lanes | 64k→3k por filing | marginal mas grátis | prática de hub |
| 6 | **Codex volta 2026-07-25**: P2/P3/P6-Sol de volta pra GPT | implementação sai da conta Anthropic | estrutural (−60%+ por chip) | já doctrine; aguardar reset |

Combinado (1+2+4, aplicável já na wave C): **−50–60% custo por chip** vs baseline wave B.

## Eval A/B (wave C)

Baseline = números wave B acima. Wave C roda com regras 1/2/4 nos prompts. Medir igual
(script de dedupe): API calls, output por papel (main vs subs), cacheR/call médio.
Targets: main output share <35% · ctx médio/call <80k · calls/chip <300 por fase.
Bater targets ⇒ propor amendment upstream mnfs-harness (com estes dados como evidência).

## Limitações

- Preços de lista como proxy; operador está em subscription (proporções valem, $ absoluto não).
- Sessões -b e subagents não deduplicados individualmente (fator ~2,7 verificado nos mains).
- Sessão fable 3a34ff8c (967 calls, 125M cacheR bruto) não identificada — fora do baseline.
