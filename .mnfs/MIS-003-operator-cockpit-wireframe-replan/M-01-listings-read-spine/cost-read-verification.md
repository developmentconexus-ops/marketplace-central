# M-01 · Cost-read verification (below_margin / Slice 5 / C09)

## DB specialist consultation address (hub-provided)
- Session `local_ec787804-f8e9-4981-9c12-7d3f45292294` ("Marketplace Central database queries", repo MNOS). READ-ONLY schema METALPRD.
- Consult via `mcp__ccd_session_mgmt__send_message` before finalizing Slice 5 cost-read semantics.
- Request format: (a) intent, (b) CODPROD list or resolution path, (c) as-of semantics, (d) output grain the Go port expects. Replies with bind-var-ready SELECT + joins + reconciliation vs real METALPRD when golden exists.

## Specialist-ratified Oracle facts
- Latest cost is **AS-OF**, not scalar: per CODPROD, row with `MAX(DTATUAL) WHERE DTATUAL <= reference_date` in `TGFCUS`. Batch reads correlated per-id.
- Cost basis = `TGFCUS.CUSSEMICM` (cost **WITHOUT ICMS**). Never `TGFITE.CUSTO` (per-note historical). `CUSVARIAVEL` ≠ Gasto Variável.
- Product identity: `TGFPRO.CODPROD` (int PK). Price list `TGFTAB`/`TGFEXC`; realized `TGFITE`+`TGFCAB`; stock `TGFEST`.

## Verification of `internal_read` BatchReader.GetCostFactsByIDs (VERIFIED — matches specialist)
`apps/server_core/internal/modules/internal_read/adapters/oracle/batch_reader.go:40-96,246-264`:
- **As-of correlated per-id** ✓ — `buildCostFactsQuery` uses `c.DTATUAL = (SELECT MAX(latest.DTATUAL) FROM TGFCUS latest WHERE latest.CODPROD = c.CODPROD AND latest.CODEMP = c.CODEMP AND latest.DTATUAL <= :ref)`; `CODEMP = 1`.
- **Basis CUSSEMICM** ✓ — `SELECT c.CODPROD, c.CUSSEMICM, c.DTATUAL`; domain `CostBasisCUSSEMICM` ("cussemicm").
- Null preserved → `QualityMissingCost` (ADR-17). No per-row query (single batched IN-chunk). Satisfies the D-20 "one bounded Oracle batch cost read".
→ No BLOCKED on as-of/basis. The read port is contract-correct for below_margin cost input.

## OPEN — ICMS comparison-basis mismatch (REQUEST sent to hub)
below_margin literal (D-16): `price < cost × (1 + min_margin%)`. Cost = CUSSEMICM = **net of ICMS**. ML listing price = **gross** (consumer price, tax-inclusive). Comparing gross price to net cost systematically misclassifies. Hub rule: do NOT silently pick. Options tabled to hub:
- (a) net the price side to a comparable ICMS basis before comparison (ICMS available via `GetTaxFactsByIDs`, but that is per-note VW_IMPOSTO_ITEM, not a clean per-product rate — needs a defined rate source);
- (b) document the criterion as gross-price vs net-cost with the known bias (weakens honesty);
- (c) amend the below_margin criterion literal.
Blocks below_margin classification only (Slice 2 enrichment + Slice 5 counter). Slices 1/3/4 and route/OpenAPI work proceed. Awaiting hub ruling.

## Hub ruling (received) — route (a) conditional
Route (a) GRANTED conditionally: net the price side to a CUSSEMICM-comparable basis (only option preserving D-16 honesty; (b) = ADR-17 violation in a footnote, rejected; (c) = no ratified with-ICMS cost source, rejected). AUTHORIZED to consult DB specialist `local_ec787804` with 3 questions: (1) canonical per-CODPROD ICMS rate / gross→CUSSEMICM-comparable adjustment path (product ICMS table? NCM/UF-dependent? fixed company rate?); (2) whether stable enough for a batch read w/ same as-of discipline + exact tables/joins; (3) golden reconciliation for one CODPROD (gross price, ICMS adjustment, net basis, CUSSEMICM, resulting below_margin) vs real METALPRD.
CONDITIONAL: clean per-product path → adopt (a), document formula+source in this file + F-02/validation.md, flag in CLOSED as D-16 basis clarification (criterion literal unchanged, bases made comparable); unknown/missing ICMS rate for a row → below_margin unevaluable, surfaced like unknown cost (never defaulted to a rate). No clean path (ICMS only per-note/UF at sale time) → STOP below_margin, report to hub, hub escalates basis to operator. Do NOT improvise a rate.
STATUS: consultation SENT to specialist `local_ec787804` (this session, 3 questions verbatim per ruling) — below_margin work parked pending reply. F-01 Slice 1 dispatched in parallel (no cost dependency, agent__f01-slice1.log).

## Specialist REPLY — route (a) FAILS the "clean path" condition (received 2026-07-15)
No deterministic per-CODPROD ICMS rate exists to net a gross price down to a CUSSEMICM-comparable basis. Evidence (live METALPRD introspection, specialist):
- Rate engine = `METALPRD.TGFICM`, PK is 7 fields (UFORIG, UFDEST, TIPRESTRICAO, CODRESTRICAO, TIPRESTRICAO2, CODRESTRICAO2, SEQUENCIA) — resolves per grupo × UF-origin × **UF-destination** × CFOP/negociação × sequence. NOT scalar per product/group/company. Group 201, pair 13→8 alone = 150 rows.
- `TGFPRO.GRUPOICMS/NCM/ORIGEMFAT/TEMICMS` only NAMES the fiscal group; yields no rate without UFDEST+CFOP context.
- Same CODPROD 42664, empresa 1, maio/2026 realized: intra-MG (ST) → ALIQICMS **0%**; interstate → **12%**; MG normal-tributed → 7%/18%. Buyer UF unknown at listing time → no single as-of rate. CUSSEMICM(42664)=65,6597 [golden cus1]; adjustment is 0% OR 12% on the SAME product → no deterministic below_margin decision under "net price to CUSSEMICM basis".
- No effective-dated per-product rate table → cost read's MAX(effective)≤ref discipline does NOT transfer. Only DB-derivable batchable value = REALIZED effective rate per CODPROD over a trailing window (approximation, high variance mix of 0%/12%), not authoritative.

CONCLUSION (per hub ruling's "no clean path" branch): STOP below_margin by DB-derived per-product rate. The basis is a BUSINESS decision — escalated to operator. Specialist's three real options (each needs operator to SUPPLY the rate assumption; DB cannot):
- **A** Fixed marketplace ICMS assumption, **net the price** down by it (operator declares rate, e.g. all ML = interstate 12%, or per-UF if ML gives ship-to). Deterministic, batchable.
- **B** Realized effective-rate **proxy** per CODPROD (trailing-window SELECT, provided) — accept approximation + variance + seed rule for never-sold items.
- **C** Gross-up **CUSSEMICM** by the same assumed rate instead of netting price (mirror of A; leaves the verified cost read untouched).
Realized-proxy SELECT (if B chosen) archived from specialist: SUM(VLRICMS)/NULLIF(SUM(BASEICMS),0) over TGFITE⋈TGFCAB, CODEMP=1, STATUSNOTA='L', venda TOPs (101,21,116,303,313,701,203), window bind vars.

→ ESCALATED to operator (this session, precedent D-16 = operator-ratified below_margin). F-01 unaffected.

## Operator steer + HUB escalation (2026-07-15)
Operator did NOT pick A/B/C — directed: do NOT drop below_margin, **SIMULATE** it. ICMS varies by UF-destination (e.g. MG→SP = interstate ICMS + DIFAL). Compute a **worst-case** = max ICMS payable; ideally **per state** (per UF-dest). Conservative upper-bound below_margin flag. Operator available to talk to hub directly.
ESCALATION sent to hub `local_efa46c30` with 3 research asks: (1) is per-(GRUPOICMS, UF-dest) worst-case output-ICMS+DIFAL deterministically derivable from TGFICM with as-of discipline? (2) worst-case single rate vs full per-state matrix (latter expands read contract with a per-UF dimension)? (3) M-01 read-spine scope vs mission-level defer (below_margin is 1st-class IC-02 exception D-17 → dropping risks gate fail; simulating = net-new scope).

## HUB RULING → D-22 (received, operator-ratified 2026-07-15) — below_margin UNPARKED
below_margin SHIPS in M-01 as worst-case ICMS simulation (scalar flag + per-UF matrix). Full spec in DECISIONS.md D-22. Rate source: `MAX(ALIQUFDEST)+NVL(MAX(PERCICMSFCP),0) FROM TGFICM WHERE UFORIG=:uforig GROUP BY UFDEST` (product-independent, one batch, NOT date-versioned → "current-config" rate). Formula ICMS-por-dentro: `price_net_basis = price_gross×(1−worstcase_pct/100)`; `below_margin = price_net_basis < CUSSEMICM×(1+min_margin/100)`. ALIQUFDEST+FCP ceiling ONLY (never ALIQUOTA/REDBASE, never realized TGFITE.ALIQICMS). Unknown cost OR unknown UF ceiling → unevaluable/margin_unknown, never defaulted. Golden 42664/RJ = 22% (12+8+2). Oracle within D-20: cost batch + 1 rate-ceiling batch. Evidence → F-02/validation.md; flag in CLOSED. Deferred to hub queue: per-NCM/ST refinement + configurable-rate UX.
STATUS: below_margin resumes at F-02 Slice 2 (enrichment) + Slice 5 (summary counter) + Slice 3 (per-UF matrix in GET /listings/{id}). F-01 Slice 2 (connectors mapper) unaffected — running.

## HUB RULING → D-24 (received 2026-07-15) — internal_read ceiling read LOCKED
Additive contract-lock GRANTED on internal_read for one method `GetICMSCeilingByOrigin` (port + oracle adapter + domain `ICMSCeiling` + fake stub). Query = specialist bind-var verbatim (`MAX(ALIQUFDEST)+NVL(MAX(PERCICMSFCP),0) GROUP BY UFDEST WHERE UFORIG=:uforig`). Reuses internal_read Oracle infra (semaphore/wrapOracleError/nil-unknown/chunking) — no duplicate connection in listings. TGFICM current-config (no as-of, D-22.4). UF w/o row → nil/unevaluable/never-defaulted. Adapter test MUST negatively assert ALIQUFDEST(+FCP) not ALIQUOTA (REDBASE trap = known wrong-rate failure). Terms = D-3/D-21; lock ends at CLOSED; diff in CLOSED payload; beyond one method = new REQUEST. Full spec DECISIONS D-24. Consumed by F-02 below_margin slices; F-01 unaffected. below_margin read path now fully unblocked.
