# AI Dialog

Candidate: `developmentconexus-ops/marketplace-central` / PR #62 / `2aadfcdeb74d707a3f08fa3d1ee48951876c30d4`  
Candidate base: `main@b54d17bfe6d794645d198a9160f4a2a1c63647e8`  
Round: `R1`  
Methodology: `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`

## Findings

| ID | Severity | Finding |
| --- | --- | --- |
| R1-F1 | IMPORTANT | Retained legacy guide carries a live floating-`main` methodology link and a dead parent-method link, while `scripts/gate.ps1` simultaneously lists the guide as a required file — the gate manufactures a permanent consumer for a surface declared "keep only while it has a real reference consumer". |
| R1-F2 | IMPORTANT | Exact-pin consumption law is enforced on `AGENTS.md` only; 9 further pinned occurrences (5 deep links in `docs/index.md`, plus roadmap/engineering-rules/frontend pointer) are unguarded against pin drift on upgrade. Also a METHOD FINDING candidate. |
| R1-F3 | MINOR | The `methodology exact pin` negative control proves only the missing-pin branch of `Assert-MethodologyBootstrap` fires; the repository-name and ROUTER-route branches are presence-checked but never demonstrated to fire. |
| R1-F4 | MINOR | The "v2.3 generation" cited in the supersession pointer and in pinned `FRONTEND-METHOD.md` lineage has no committed byte provenance in either repository; Marketplace history holds only v2.1. All MPC historical artifacts cite v2.1 (reachable), so the local reachability law holds; the residue is canonical-repo lineage prose. |
| R1-F5 | MINOR | Bootstrap trio is 18,984 / 20,480 bytes (92.7%) after adoption; headroom for future roadmap growth is thin. Observation with a reopen trigger, not a defect. |

## Dialogue

### R1-H0 — LEAD ROUTING ONLY

Revalidate remote repository, PR #62, exact candidate HEAD, and current authority before reviewing. Use the repository bootstrap and the pinned `ROUTER.md` → `METHOD.md` + `ADVERSARIAL-REVIEW-METHOD.md` profile. Load another specialist only when a concrete attacked claim requires it.

Attack the exact methodology-adoption candidate as repository-governance work: exact pin/consumption law, fresh-session recovery, authority duplication/supersession, acceptance-increment semantics, local specialization, context budget, aggregate verification, review hygiene, and whether Product/frontend authority was accidentally reinterpreted.

Do not modify `chore/methodology-v1-1-adoption`. Publish the Challenger turn only in this temporary file on the review branch. Reviewer findings are Evidence, not acceptance authority.

### R1-C1 — CHALLENGER

Independent challenge of the exact candidate `2aadfcde` (PR #62, base `main@b54d17bf`). Revalidated against remote state 2026-08-24: main HEAD, PR #62 base/head/state (OPEN, MERGEABLE, required CI SUCCESS), PR #63 base/head, and review isolation all match the handoff. The candidate was not modified. Context: repository bootstrap route + pinned `ROUTER.md` → `METHOD.md` v1.1.0 + `ADVERSARIAL-REVIEW-METHOD.md` v1.0.0; `FRONTEND-METHOD.md` and `REPOSITORY-STANDARD.md` v1.1.0 loaded because attacked claims 4, 7, 8 and 13 require them.

**Executed proof.** `scripts/gate.ps1` run in review mode against the exact candidate (`GATE_BASE_REF=chore/methodology-v1-1-adoption`, `GATE_BASE_SHA=2aadfcde`): `gate: PASS`, `review_mode: True`, `diff_range: 2aadfcde..HEAD changed_files: 1`, `negative_controls: 2/2`, `bootstrap_bytes: 18984 / 20480`, `methodology_pin: developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`, Product OAD proof `106/106` PASS. CI wiring (`.github/workflows/ci.yml`) supplies `GATE_BASE_REF`/`GATE_HEAD_REF` from the PR event, so PR #63 exercises review mode against the exact candidate and PR #62 exercises candidate hygiene — both required checks are green on the current heads.

**Candidate boundary held.** `git diff --stat b54d17bf..2aadfcde` touches exactly the six declared files; no `contracts/`, `docs/engineering/`, `ARCHITECTURE.md`, ADR, or frontend P8 HTML change. Roadmap retains 106/31/H-A-S, active runtime NONE, implementation BLOCKED UNTIL D9, and all stage-marker guards still pass. Product non-regression (attack 11) is proved by diff scope plus the executed OAD proof, not narrative.

Attack results, keyed to the handoff:

1. **Exact consumption law — holds, with one enforcement gap (R1-F2).** The pin is immutable, discoverable from bootstrap, and stated identically across `AGENTS.md`, `docs/index.md` (5 deep links), `docs/roadmap.md`, `docs/development/engineering-rules.md`, the supersession pointer, and `scripts/gate.ps1` — 10 occurrences, zero drift today. Consumption is canonical reference, not copied bytes. `ROUTER.md` is reachable from a fresh bootstrap via the pinned link in `docs/index.md` step 4 (verified by fetching it at the exact pin). Gap: see R1-F2 below.

2. **Fresh-session recovery — PASS.** From remote state alone: `AGENTS.md` → `docs/index.md` → `docs/roadmap.md` names the current increment (PR #62), the target pin, the exact next action, and every block. This review itself was reconstructed from that route without chat archaeology, and the route stays selective (this reviewer read 3 repository files + 5 pinned Method files + the candidate diff).

3. **Authority uniqueness — PASS.** `docs/roadmap.md` remains the sole mutable status owner; the gate still enforces its marker set and rejects status duplication in README/AGENTS/index. The old duplicate "Method v1.0.0 / Repository Standard v1.0.0" declarations were removed from all three active surfaces and the gate now fails if they reappear there (verified in gate source). The historical `> Method: DevelopmentConexus Engineering Method v1.0.0` banners inside accepted D-stage artifacts are truthful frozen snapshots covered by the existing frozen-routing-snapshot doctrine in `docs/index.md`; they are correctly excluded from the scan and were correctly not rewritten.

4. **Frontend-method supersession — CORRECT.** The local v2.1 method (1,513 lines at `b54d17bf`) was compared section-by-section against pinned `FRONTEND-METHOD.md` v1.0.0: operator-only LOCK, coverage-before-layout, no screen-shaped API / no backend-shaped UX, P8 functional-HTML medium, data-feasibility blocking law before P8 LOCK, walkthrough evidence, block operating protocol, pattern graduation, smallest-reopen law, and accessibility/responsive-as-structure are all absorbed; the global additionally carries the LOCK-impact sweep and P11 assembled-fidelity law the local v2.1 lacked. No CURRENT MPC-specific normative requirement was destroyed: MPC status stays in the roadmap, MPC specialization in engineering-rules. `SUPERSEDED_BY_GLOBAL` with a lineage pointer is the right disposition; 5+ accepted D6-R2 artifacts link the old relative path, so the pointer's retention rationale is factually grounded (verified by grep across the candidate tree).

5. **Historical Evidence preservation — PASS, one provenance residue (R1-F4).** Accepted artifacts still truthfully cite "Method v2.1" and were not rewritten to claim later-Method parentage; the pointer states this explicitly. v2.1 bytes are reachable in Git history (`d409e643`, `ee3d1ec5`). Residue: neither repository's history contains v2.3 bytes, yet both the pointer ("v2.1/v2.3-era") and canonical `FRONTEND-METHOD.md` lineage ("through the operator-ratified MetalDocs v2.3 generation") cite that generation.

6. **Marketplace specialization — PASS.** Every retained engineering-rules block is genuinely MPC-specific (Sankhya gateway mandate, ML operator authorization, Organization isolation, PII, OAD/dependency scope, path specializations with consumer + removal trigger in the Repository Standard §14 deviation shape). No required Product/security/tooling rule was found removed: the AGENTS hard-rail set is preserved or strengthened, and generic law now correctly lives only in the pinned suite. The removed local "Material stop conditions" and "surface authority conflicts" prose is owned by pinned `METHOD.md` (STOP/SPLIT outcome; conflict-surfacing law) plus the per-specialization "path aesthetics" removal triggers — absorbed, not lost.

7. **Acceptance-increment semantics — PASS.** `AGENTS.md` now states one branch/Draft PR per acceptance increment with stages remaining open across integrated increments; roadmap models exactly that (PR #61 integrated, D6-R2 still OPEN, PR #62 a separate increment). No active "one PR per whole stage/gate" authority survives (grep clean); the only stage-shaped naming left is historical.

8. **Context budget — PASS with thin headroom (R1-F5).** Bootstrap trio 18,984 / 20,480 bytes, gate-enforced. Repository-local pack ≤ 5 files with the Method profile correctly excluded from the count in both `docs/index.md` and pinned `ROUTER.md` — the two statements agree, no duplication. Routers are compact and non-overlapping (index routes tasks, roadmap owns status).

9. **Aggregate verification quality — substantively good; two bounded gaps (R1-F2, R1-F3).** The pin guard fires (negative control executed), candidate/main ai-dialog hygiene is enforced in non-review mode, review mode proves the exact-candidate isolation including candidate-tree `docs/work` contamination, and PR-diff proof runs over the intended base…candidate range. The gate proves behavior, not just string presence, for review isolation and the OAD surface. Gaps are recorded as findings below.

10. **Review hygiene — PASS.** Candidate tree contains no `docs/work/**` or ai-dialog (verified via `ls-tree` in the executed gate run). This review branch differs from the exact candidate by `docs/work/current/ai-dialog.md` only (verified by diff and by the review-mode gate). Nothing here can become Product/status authority: the file is unrouted from durable docs and mechanically blocked from candidate/main.

11. **Product non-regression — PROVED.** See "Candidate boundary held" above.

12. **Retained derived guide — the one real hidden-second-Method vector found (R1-F1).** The guide itself is declared non-authoritative twice (its own header, engineering-rules) and the index row demotes it to "optional reference only". But three facts compose badly: (a) line 662 links `conexus-methodology/blob/main/METHOD.md` — a live floating-`main` consumption path that AGENTS.md expressly prohibits; today `main == 9c7210d1`, so it is silently correct, and after the next canonical commit it silently serves un-adopted Method text; (b) its header's parent-method link `./root-cause-global-maximum-method.md` is dead (file retired), and the gate's link checker only fails links that resolve to tracked-durable paths, so dead links pass silently; (c) `scripts/gate.ps1` lists the guide under `$requiredFiles`, so the gate fails if it is ever deleted — the aggregate verification itself now manufactures the "real reference consumer" that engineering-rules says is the only reason to keep it. This is a concrete current route (index row "Production coding…" → guide → floating-main METHOD.md), not tree-cleanliness preference. Deletion is not demanded; neutralizing the consumption path is.

13. **Consumption machinery / YAGNI — PASS.** No submodule, sync bot, projection, central verifier, or duplicated Method files were added. Canonical-reference-only matches Repository Standard §13. Nothing speculative found.

14. **METHOD FINDING — one candidate.** Repository Standard §13 requires the pin to be discoverable from bootstrap but states no law that *every* pinned deep link across a consumer's surfaces must equal the accepted pin, nor that active surfaces must be floating-`main`-free, nor that a consumer's aggregate verification should census this. Marketplace Central now de facto needs exactly that safeguard (R1-F1/R1-F2 are both instances of the missing law: "single-pin census — any `conexus-methodology/blob/<sha>` reference on an active surface equals the accepted pin; `blob/main` is forbidden outside explicitly disclaimed historical bytes"). This is cross-repository by construction and should be surfaced to the canonical methodology rather than remaining a silent local-only guard once implemented here.

15. **Global Maximum — the six-file shape is the smallest sustainable solution.** It removes duplicate active authority without touching accepted semantics, preserves genuine specialization, and adds no ceremony. It misses one necessary property (pin-census enforcement, R1-F2) and leaves one legacy consumption path live (R1-F1); both are bounded gate/doc corrections on the candidate, not restructuring. No over-removal was found (see 4 and 6). No duplicate authority survives in active surfaces.

**Detailed findings.**

**R1-F1 — IMPORTANT — retained guide: live floating-`main` link, dead parent link, gate-required presence.**
Claim attacked: "the retained guide is historical/portable non-authoritative reference, not a current Method selector."
Evidence: `docs/development/evidence-grounded-production-engineering-for-llm-agents.md:662` (`blob/main/METHOD.md` URL); its header parent link to nonexistent `./root-cause-global-maximum-method.md`; `scripts/gate.ps1` `$requiredFiles` including the guide; index route row naming it; gate link checker passing dead links by construction.
Why it matters: the only prohibited consumption mode (floating `main`) survives as a working link on a routed surface, and the gate hard-codes the file as load-bearing while local authority says it is kept only while a real consumer exists.
Smallest scope: that one file's two link lines + `$requiredFiles` membership decision + optionally one gate scan line.
Disposition/falsifier: bounded candidate fix — neutralize or pin the line-662 link and annotate/remove the dead parent link (bytes may change since the guide's *authority* is what was frozen, or a one-line erratum block may be appended if byte-freezing is preferred); decide explicitly whether the guide stays in `$requiredFiles` or is demoted to optional-tracked; alternatively record an explicit §14-shape deviation accepting the residue. Falsifier for the fix: gate scan failing on `conexus-methodology/blob/main` in any routed durable doc.

**R1-F2 — IMPORTANT — exact-pin law enforced on one of six surfaces.**
Claim attacked: "exact pin actually enforced."
Evidence: `Assert-MethodologyBootstrap` receives only `$agent`; grep census shows 10 pin occurrences across 6 files, including 5 deep links in `docs/index.md` that route fresh agents directly to Method bytes.
Why it matters: this is the upgrade-drift class — a future pin move that edits `AGENTS.md` but misses one `docs/index.md` URL passes the gate while fresh sessions silently load the old Method profile from the stale link. No current drift exists; the gap is the absent guard, not present bytes.
Smallest scope: `scripts/gate.ps1` only.
Disposition/falsifier: extend the gate with a pin census (every `conexus-methodology/blob/<40-hex>` occurrence in tracked active surfaces equals `$methodologyPin`; forbid `blob/main` outside the explicitly disclaimed guide until R1-F1 is resolved) plus one deterministic negative control with a mismatched-pin fixture. Also the METHOD FINDING carrier (attack 14).

**R1-F3 — MINOR — pin negative control proves one of three branches.**
Evidence: the `Expect-Failure 'methodology exact pin'` fixture (`'developmentconexus-ops/conexus-methodology ROUTER.md'`) contains the repository string and `ROUTER.md`, so only the missing-pin branch ever throws; the repository-name and ROUTER-route checks are never demonstrated to fire.
Why it matters: "a control counts only when its firing can be demonstrated" (METHOD.md Enforcement); two of three branches are presence-only today.
Smallest scope: the one negative-control block.
Disposition: either two additional one-line fixtures or an explicit statement that the pin branch is the only branch counted as a falsifier. Compatible with the gate's stated reusable-predicate counting policy.

**R1-F4 — MINOR — v2.3 byte provenance is nowhere reachable.**
Evidence: `git log --all` over the local method path shows only v2.1 commits; the conexus-methodology tree at and before the pin never contains a v2.3 file; both the new pointer and canonical `FRONTEND-METHOD.md` lineage cite the "v2.3 generation".
Why it matters: Repository Standard §9 makes Git history a sufficient archive only while required history is reachable. For Marketplace Central the law is satisfied — every accepted artifact cites v2.1, which is reachable. The unreachable generation is cited only as lineage prose.
Smallest scope: canonical-repo lineage wording; no Marketplace change required.
Disposition: DEFER to the canonical methodology owner — either record where v2.3 provenance durably lives or soften the lineage claim to what is reachable. Not a Marketplace candidate defect.

**R1-F5 — MINOR — bootstrap headroom 7.3%.**
Evidence: executed gate output `bootstrap_bytes: 18984 / 20480`; growth came mostly from the five long pinned URLs in `docs/index.md`.
Why it matters: the next few roadmap-table rows may force reactive compression during unrelated work.
Smallest scope: none now.
Disposition: NO CHANGE REQUIRED; reopen trigger = any gate failure on the 20 KiB budget, resolved by URL/prose compression rather than budget raise.

**Verdict.** No MATERIAL finding. The candidate is a genuine, correctly bounded repository-governance increment: exact-pin adoption by canonical reference, correct supersession with truthful history, no Product/frontend/LOCK mutation (proved), review isolation proved by the executed gate. R1-F1 and R1-F2 are bounded IMPORTANT corrections on the candidate (plus one METHOD FINDING to surface upstream); neither falsifies the central decision.

`LEAD RESPONSE REQUIRED`
