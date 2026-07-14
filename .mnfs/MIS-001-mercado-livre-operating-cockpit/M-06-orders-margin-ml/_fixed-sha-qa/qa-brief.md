# M-06 Proportional QA Brief

Validate only frozen SHA `ef4b08c78d30a5e2269e79b051a432c9dc12b58d`
after the independent fixed-SHA review.

Run proportional deterministic QA for the M-06 validation contract and the
resumed assisted Sankhya linkage slice. Capture durable exact command/output
evidence under `_gate-evidence/round-4/`. Validate idempotent order behavior,
stable MPC line identity, assisted exact TOP 313 confirmation, exact TGFVAR
TOP 306 lineage, one-to-many/partial tax aggregation, unknown-value honesty,
API/SDK parity, and append-only audit mechanics.

Do not use live Oracle, PostgreSQL, provider, browser, or production writes.
Their absence must not be upgraded into proof. The F-08 deployment header
field/effective-TOP facts remain unknown and activation must stay fail-closed.
Production authentication and scoped authorization for manual adjustments are
explicitly deferred and must remain a C03 failure; do not consume correction
attempt 2 or implement a fix.

Only QA may update `validation-result.md`. Preserve prior rounds append-only,
record the exact fixed SHA, commands, outputs, target types, verdict, blockers,
and remaining correction budget. No source or Git mutation.
