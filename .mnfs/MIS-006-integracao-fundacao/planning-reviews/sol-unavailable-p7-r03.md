# Sol Unavailable — P7 round 03

```yaml
phase: P7
round: 03
reviewer: gpt-5.6-sol (HIGH)
status: unavailable
cause: codex-quota-wall
quota_resets: 2026-07-25
operator_authorization: present (dispatch brief)
```

## Blockage

`/codex:rescue --model gpt-5.6-sol` unreachable — codex usage limit hit 2026-07-18, resets
2026-07-25 (post-demo). Wall is mission-wide, blocks ALL codex roles including the P3/P5/P7 Sol
touchpoints. Recorded in memory `codex-quota-exhausted-jul25`.

## Contingency invoked (operator-authorized)

The mission dispatch brief pre-authorized the codex quota-wall contingency: rebind the Sol-side
review to a second independent COLD Claude `mission-reviewer` crew for the duration of the wall.
This is an explicit operator grant on record for THIS mission, not a silent Claude substitution.

Distinction from the protocol's default failure/skip rule (which forbids Claude substitution and
`status: planned`): that rule governs the UNAUTHORIZED case. Here the owner has authority over the
gate composition and exercised it in advance. The rebind is logged, and the P7 gate remains a
genuine dual **crew** gate (two independent cold crews on the same frozen manifest), just not a
dual **model** gate while the wall stands.

## What the rebound Sol-side crew was

The round-03 focused ★2 re-review (`p7-claude-readiness-r03.md`) plus the held-PASS folds from
rounds 01–02 constitute the primary cold crew. Per operator efficiency directive ("run only the
★2-specific review, don't rerun the full crew when criteria consistently pass"), the Sol-side
rebound crew was NOT re-dispatched as a redundant full pass on round 03 — its round-01/round-02
coverage already folded PASS on the six held criteria, and ★2 (the sole moved criterion) was
re-verified by the focused reviewer.

## Follow-up owed

When the quota wall lifts (≥ 2026-07-25) and if this mission has not yet entered execution, a true
GPT-5.6 Sol HIGH pass MAY be run against the then-current manifest as a confirming (non-blocking)
review. Not a gate reopener unless it surfaces a valid ★ FAIL.
