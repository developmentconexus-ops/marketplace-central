# M-06 F-02 Independent Review Findings

## Verdict

- Spec compliance: `REJECTED`
- Code quality: `REJECTED`

## Important Finding

- `apps/server_core/migrations/0030_profitability_manual_adjustment_invariants.sql:1`
- The new checks validate all existing rows immediately. Rows admitted by the previous schema with blank actor/reason/currency or ambiguous identity can abort deployment.
- Add each check as `NOT VALID` so PostgreSQL enforces it for every new/updated row without rewriting append-only history or blocking the migration on historical defects. Historical validation/quarantine belongs to a separately evidenced data-remediation migration.

## Missing Evidence

- Extend the PostgreSQL integration test to prove both identity branches: order scope with item/variation is rejected, and item scope without provider item is rejected.

## Accepted Areas

The reviewer found the application invariants, actor normalization, cryptographic IDs, freight/commission events, signed amounts, append-only API/repository shape, and OpenAPI/SDK alignment functionally sound.
