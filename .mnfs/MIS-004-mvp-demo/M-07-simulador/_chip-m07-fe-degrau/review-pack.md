# Dual-gate review pack — CHIP-M07-FE-DEGRAU

**Fixed SHA:** `72190e6f81fb141f5b4ea21755be84d6e38fd4c7` (base `9e84fe0f`, branch `chip/m07-fe-degrau`)
**Scope owned:** `apps/web/src/pages/precos/**` only.

## Defect being fixed (hub-diagnosed, P7 live evidence)
`PricingPage.tsx` hardcoded ML commission per modalidade (classico=12, premium/full=17)
and ALWAYS sent `comissao_pct` on `/pricing/decompose` and to `SolverPanel`
(`/pricing/solve`). Backend treats any present `comissao_pct` as a MANUAL override
(degrau 0), gating OFF the real tariff resolver chain (degrau 3 = live ML COTACAO fee
quote; degrau 4 = tenant PADRAO default). Live proof: decompose for product 90001
WITHOUT `comissao_pct` → comissao 12 fonte COTACAO degrau 3; WITH it → MANUAL degrau 0.

## Contract fact
`PricingCalcInput.comissao_pct?: string` is OPTIONAL (`packages/sdk-runtime/src/index.ts:1485`).
Empty/absent ⇒ backend resolver chain runs (`calc_service.go resolveTariff`,
`reqComissao == ""` enables degrau 3). Confirmed released contract at base SHA.

## The fix (should be)
- Remove `MODALIDADES[].comissaoPct` + `comissaoFor()`; stop sending `comissao_pct` on
  decompose and on solve; drop the `comissaoPct` prop from `SolverPanel`.
- `modalidade` (contract-required) is still sent on both.
- Tests adjusted to assert `comissao_pct` is OMITTED from request payloads.

## Out of scope (NOT defects — do not request)
- Any backend/OpenAPI/SDK change.
- A manual-override input UI (post-demo feature).
- Other pages / other precos components.
- Known honest post-fix behavior: solve for ERP-xlsx products returns degrau 4 PADRAO
  (13.00) — solve has no price basis and ERP import carries no sale price. Honest by design.

## Your task
Adversarial review per REVIEW-STANDARD (design → correctness → complexity → tests →
naming → docs; style is machine-owned). The DIFF is the unit under review (below).
Verdict = PASS or FAIL. Every finding: two-axis severity
(`blocking|important|suggestion|nit|question`) + anchored `path:line`, anchor-or-abstain.
FAIL only on a blocking/important defect you can quote. Report:
1. VERDICT (PASS/FAIL)
2. Findings (severity + path:line + one-line each), or "none"
3. One-line confirmation the fix actually removes the MANUAL-override cause.

## Diff under review
See `diff-full.patch` in this same directory (also inline below).
