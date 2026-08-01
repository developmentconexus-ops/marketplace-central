# F-02-stock-divergence

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-05 listings-fees-divergence.

## Brief

Avaliador de divergência de estoque (arquivos novos em listings): compara disponível do
espelho ERP (products_mirror, predicado vendável FECHADO: disponível = ESTOQUE−RESERVADO,
memória `vendavel-cut-flink1-resolved`) vs `available_quantity` do anúncio ML — grão
VARIAÇÃO quando listing_variations tem rows p/ o anúncio, senão grão listing (IC-07). Só
avalia anúncio COM vínculo aprovado (LinkedListings — sem vínculo não há o que comparar,
IC-02). Divergente (tolerância 0) → upsert row aberta via porta IC-02
(`entity_type=listing`, kind=`estoque`, expected/observed + AMBOS observed_at NOT NULL:
mirror synced_at e listing last_seen_at); convergiu → auto-resolve (`resolved_at`). Roda
como passo pós-sweep do ingest (mesma composição do M-04 F-04) — nunca no request.

EARS:
- While anúncio vinculado com estoque ML ≠ disponível ERP, when avaliador roda, the sistema
  shall upsert UMA row aberta por (listing, estoque) com os dois lados + timestamps.
- While divergência aberta e valores convergem, when avaliador re-roda, the row shall ganhar
  resolved_at (nunca DELETE).
- While anúncio sem vínculo, when avaliador roda, the sistema shall pular sem row.

## Inputs

IC-02 (integral); IC-07 (grão variação); M-04 (available_quantity + listing_variations
ingeridos); products_mirror (M-02-MIS-006, leitura); LinkedListings existente (vínculo).

## Expected Output

Avaliador + fiação pós-sweep + leitura de mirror por porta (NUNCA SQL cross-module direto —
porta de leitura no módulo listings com adapter). A fiação vive DENTRO do package de
composition de listings, sob o lock aditivo registrado do M-04 — `root.go` NÃO é tocado
pelo M-05 (matriz: célula root.go = `—`; auditoria P5 r03 P-5).

## Constraints

- Um lado ausente (mirror sem produto OU listing sem last_seen) → NÃO avaliar (IC-02 §both
  timestamps NOT NULL é CHECK — ausência de observação não é divergência).
- Tolerância estoque = 0 (IC-02 binding).
- Partial unique `WHERE resolved_at IS NULL` (0087) é a proteção contra duplicata — upsert
  usa ON CONFLICT nela.

## Inputs/Outputs

Row canonical: IC-02 §Canonical Examples kind=estoque (binding).

## Negative Scenarios

- Anúncio com variações: divergência avaliada POR VARIAÇÃO (subject carrega variation_id no
  formato IC-02); comparar agregado vs variação = bug nomeado no teste negativo.
- Mirror stale (synced_at velho): avalia mesmo assim — timestamps na row deixam o leitor
  julgar (honest observation, não censura).

## Ownership

- Owned paths: avaliador novo + porta de leitura mirror em listings; fiação pós-sweep.
- Forbidden paths: schema divergences (M-02); products_mirror (read-only); transport (F-03).
- Parallel-safe with: F-01 (write-sets disjuntos).

## Validation Expectations

- Fixture 2 direções: plantar divergência → row aberta; corrigir → resolved_at (IC-02 exige
  as duas).
- Anúncio com 2 variações, 1 divergente → exatamente 1 row, subject da variação certa.
- Sem vínculo → 0 rows.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
