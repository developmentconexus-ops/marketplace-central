# M-06 cold milestone review — 2026-07-10

## Independent review result

**BLOCKED, not failed.** The final independent quality reviewer found no
actionable code-quality defect in the corrected F-03 profitability path. It
confirmed safe partial-tax propagation, correct cost-adjustment scope,
idempotent adjustment persistence, required OpenAPI/SDK request keys, the UI
retry fingerprint lifecycle, and the profitability/orders boundary.

This is not an approval of M-06 because the live, approved resolved-link
scenario required by the F-03 design and plan has not occurred. A targeted
buyer-PII scope review has now passed: normalized orders intentionally exclude
buyer/contact/address data and retain only the operational shipping identifier
and safe provider endpoint reference.

## Review trail

| Review | Result | Disposition |
| --- | --- | --- |
| Initial cold spec review | BLOCKED | Required result artifact and live resolved-link scenario absent. |
| Initial cold quality review | FAIL | Found partial-tax unknown-to-realized risk, wrong order-scope cost mapping, and non-idempotent adjustments. |
| Correction review | CONDITIONAL PASS | Found a UI pending-key lifecycle defect after a failed request and material form change. |
| Focused UI review | PASS | Fingerprint-scoped retry key regression passed. |
| Final independent quality review | PASS | No actionable code-quality findings. |
| Unapproved-link correction review | PASS | SPEC and QUALITY approved: candidates/non-resolved links expose no product ID; profitability independently requires resolved quality. |
| Buyer-PII scope review | PASS | No buyer/contact/address fields in normalized contract, API, SDK, or Orders UI; raw reference is a safe `/orders/{id}` path. |
| Milestone outcome | BLOCKED | Awaiting explicit link approval/actor and live realized evidence. |

## Evidence integrity

The validation result distinguishes contract/unit evidence from real targets:
Mercado Livre and Oracle are live-provider evidence; PostgreSQL evidence uses
the real Docker database; browser evidence uses the built-in browser. No mock
or compile-only result is being used as proof of an external integration.

See `validation-result.md` for criterion-level commands, observations, and
the exact unblock sequence.
