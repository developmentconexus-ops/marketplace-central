# M-03 Plan Adjudication — rulings on contract/plan divergences

## #1 — Status codes `execute_required` / `nothing_to_retry` (hub ruling M03-COR-1f, 2026-07-17)

- **Finding (Sol, dual gate round 1):** code returns `execute_required` → 422 (VC C04 text said 400) and `nothing_to_retry` → 409 (VC C09 text said 422).
- **Truth order (profile §8, precedent M-02 PLAN-ADJUDICATION #9):** OpenAPI + SDK > VC text. `contracts/api/marketplace-central.openapi.yaml` (~176-256) documents 422 for `execute_required` and 409 for `nothing_to_retry` — identical to code and SDK.
- **Ruling (hub):** VC stale, NOT a code defect. `validation-contract.md` C04 amended 400→422 and C09 amended 422→409, both citing this ruling. No code change.
- **Recorded:** DISPATCH-LEDGER.md M03-COR-1 row; F-03/validation.md already flagged the discrepancy at feature close.
