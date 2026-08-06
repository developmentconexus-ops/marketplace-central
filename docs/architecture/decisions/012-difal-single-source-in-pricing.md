# ADR-012: DIFAL has a single source of truth inside `pricing`

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-004's MVP demo scope but no document was
ever written under a stable number. It is reconstructed here from Assertion A2 of the
citations of `ADR-10`/`ADR-010` gathered in
`docs/architecture/decisions/_citations/adr-010-citations.md` — a distinct rule from the
MIS-007 divergences decision (this repo's `ADR-011`) that happened to share the same
token before the renumbering registry existed.

## Context

DIFAL (the interstate tax differential) was, before this decision, computed
inconsistently in three places: a mock in the Simulador, a second guess in Pedidos, and a
hardcoded São Paulo rate standing in for every destination. None of these agreed with each
other and none had a named source. MIS-004 needed one number both screens could trust
for the demo, seeded with a defensible default and adjustable for the exceptions that
already existed in 2026 (namely a rate transition effective 2026-04-01).

## Decision

**DIFAL rates are seeded and persisted once, inside the `pricing` module, as a 27-UF
table ("padrão 2026") with sparse per-UF overrides. The Simulador applies DIFAL using the
order's real destination UF against this table and is the only place the table is edited
(via its drawer); Pedidos consumes a read-only chip derived from the same table. No
second implementation of the DIFAL calculation is permitted.**

**§1 — One table, one owner module, one edit surface.** The 27-UF seed plus overrides
lives in `pricing`; the Simulador is the only editor, via its drawer, and Pedidos only
reads.
> `.mnfs/MIS-004-mvp-demo/mission.md:94` — "ADR-10 DIFAL fonte única no `pricing`: seed 27
> UFs \"padrão 2026\" + overrides esparsos; Simulador aplica destino REAL; Pedidos consome
> chip read-only; rotular \"seed — não é orientação fiscal\" | ratificada | 3
> implementações divergentes; SP hardcoded | interna−interestadual; origem SC; 12%
> MG/PR/RJ/RS/SC/SP, 7% demais (lista exata: IC-04); taxa desconhecida ⇒ unknown, não 0% |
> UF de exceção ajustada reflete em Simulador E Pedidos."
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:74-79` — "### ADR-10:
> DIFAL fonte única mínima — Decisão: tabela 27 UFs seed \"padrão 2026\" + overrides (mapa
> esparso) persiste no módulo `pricing` (dono dos parâmetros de cálculo no MVP); Simulador
> aplica DIFAL pelo destino REAL selecionado; Pedidos mostra chip informativo por pedido
> (UF do shipment × tabela). SEM agendamento/lembrete/marcar-pago."

**§2 — An unknown rate is `unknown`, never `0%`.** A UF without a resolvable rate must
render as absent, applying `ADR-017`'s general rule to this specific table.
> `.mnfs/MIS-004-mvp-demo/mission.md:94` — "taxa desconhecida ⇒ unknown, não 0%."

**§3 — The seed is explicitly not tax guidance.** The value is labeled as a seed, so it is
not mistaken for a certified fiscal source.
> `.mnfs/MIS-004-mvp-demo/mission.md:94` — "rotular \"seed — não é orientação fiscal\"."

**§4 — An edited exception UF is reflected in both consumers from the single table.** A
manual override made in the Simulador's drawer must be visible in the Pedidos chip too —
there is no second copy of the table to fall out of sync.
> `.mnfs/MIS-004-mvp-demo/mission.md:94` — "UF de exceção ajustada reflete em Simulador E
> Pedidos" (validation column).
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:79` — "Validation
> impact: teste com UF de exceção ajustada refletindo em Simulador e chip Pedidos."

**§5 — Scheduling, reminders, and payment-marking are explicitly out of scope for this
seed.** The decision covers the rate table and its two read/write surfaces only.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:75` — "SEM
> agendamento/lembrete/marcar-pago."

## Amendments

A post-decomposition audit found the mission summary's stated rate list ("12%
Sul/Sudeste, 7% resto") imprecise against the exact per-UF list already normative in
IC-04 (MG/PR/RJ/RS/SC/SP at 12%; ES, part of the Sudeste region, at 7%). The audit treated
this as a documentation-alignment fix rather than a new decision — IC-04 already governed
ES correctly and remains the normative source; only the mission-summary prose was
corrected to match it.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p7-claude-readiness-r02.md:44` — "ADR-10
> (mission.md) \"12% Sul/Sudeste, 7% resto\" impreciso vs lista exata IC-04
> (MG/PR/RJ/RS/SC/SP — ES fica em 7%) — sumário alinhado à lista exata, IC-04 permanece
> normativo (nenhuma decisão nova: IC-04 já decidia ES)."

A 2026-04-01 rate transition and an operator override to the seed were folded into the
same single-source table rather than becoming a second table or a runtime branch.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p7-sol-readiness-r05.md:28` — "The AL ledger
> entry at `mission.md:55` records the 2026-04-01 transition, post-increase seed, and
> operator override. It matches R-04's caveat at `research/difal-interna-rates-2026.md:98`
> –102 and ADR-10's single pricing source/override model at `mission.md:94`."

## Contradictions

**Number collision with an unrelated decision.** The same token (`ADR-10`/`ADR-010`) also
names the MIS-007 divergences-ledger rule reconstructed as this repo's `ADR-011`
(`011-divergences-one-open-row-per-entity-kind.md`), and, more thinly, a MIS-001/MIS-003
"mocks never claim live integration" rule out of scope for this reconstruction. Nothing
in the citing text disambiguates the number by itself; only the mission the citing file
belongs to does.

**The candidate number itself churned during MIS-004 planning.** The reconciliation step
between two independent planning candidates carried this exact rule forward as
`ADR-10` on one side and `ADR-11` on the other before it was ratified into the mission
table as `ADR-10` — the number was not settled until reconciliation.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:18` — "DIFAL fonte
> única no `pricing` (seed 27 UFs + overrides), destino real, toggle no Simulador, chip
> read-only em Pedidos; sem agendar/pagar | ADR-10 | ADR-11 | idêntico."

## Rationale

Three divergent implementations of the same tax calculation are three chances to disagree
with each other and with reality, and a hardcoded single-state rate is silently wrong for
every other destination. Collapsing to one table with one editable surface removes the
possibility of the two screens drifting apart, and treating an unresolved rate as
`unknown` rather than `0%` keeps a missing fact from masquerading as "no tax applies" —
the same reasoning `ADR-017` states generally, applied here to a rate table instead of a
money field.

## Consequences

- Both the Simulador and Pedidos read from the same `pricing`-owned table; there is no
  code path that computes DIFAL independently of it.
- The Simulador's drawer is the only write surface; Pedidos is read-only by construction,
  not by convention.
- A UF with no seeded or overridden rate must render as unknown on both screens
  simultaneously, never as a default percentage.
- Full DIFAL configuration (scheduling, reminders, payment tracking) remains explicitly
  out of scope for this decision and is deferred past MIS-004.

## Alternatives Considered

**Keep per-screen calculation, reconciled by convention.** Rejected: this is the exact
state that produced three divergent implementations and a hardcoded SP rate before this
decision.

**Default an unresolved UF to the most common rate (7% or 12%) instead of unknown.**
Rejected: a defaulted rate is a fabricated fact indistinguishable from a measured one once
it reaches the screen — the same defect `ADR-017` names generally.

## Unverified claims

None. All anchors cited for Assertion A2 were read and confirmed to match verbatim at the
cited line.
