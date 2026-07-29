# FINDING — sandbox do worker cega a lane do sdk-runtime; orquestrador vira única testemunha

`status: CANDIDATO a profile — vinculado chip-localmente por A-19; ratificação do operador pendente`
`provenance: 2026-07-29 · CHIP-VENDAVEL S8 · worker 20dbc321 · evidence/S8-orchestrator.txt`

## O facto

`npx --no-install vitest run` de dentro de `packages/sdk-runtime` é impossível para worker
despachado neste sandbox: o esbuild sobe a árvore para resolver `vitest.config.ts` e o sandbox
nega a travessia (`Cannot read directory "../../../../../../..": Access is denied`). O vitest
aborta ANTES da coleta — zero arquivos, zero testes, zero skips. No assento do orquestrador a
mesma lane roda (base 5 arquivos / 77 testes).

O worker do S8 fez o certo: contou o próprio zero, escreveu *"No test was collected or skipped;
no skip name exists"* e PEDIU re-run fora do sandbox em vez de tratar como verde. A regra do
VC-7 (`No test files found` nunca é verde) segurou sob uma falha que o brief não previu. Mas
cego mesmo assim embarcou: a lane que ele não pôde rodar era exatamente a que o teria parado —
o defeito só apareceu no P4 do chip.

## Regra (chip-local por A-19; candidata a profile)

Brief de fatia que toca `packages/*` (ou qualquer lane que o sandbox do worker não executa):
1. DECLARA a cegueira de antemão e nomeia o assento que mede — o orquestrador roda a lane a
   cada entrega, contada por linha;
2. zero-observed do worker = BLIND, nunca verde; a entrega só fecha com a medição do assento
   que enxerga.

Alternativa estrutural (harness upstream, gated): sandbox do worker ganhar leitura da raiz do
repo para resolução de config — enquanto não existir, valem as regras 1–2.
