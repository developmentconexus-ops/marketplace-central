# D5 Product OAD — Independent Fable Review

> Review branch only: `review/d5-fable`
> Candidate branch: `d5-product-oad`
> Candidate HEAD at handoff: `f9f581e901df384f4dff6acdbabe9ed57b867036`
> Base `main` at handoff: `4f614e1e289e817ed6d2ca9bcdaf3b97311e6c3d`
> Product implementation: BLOCKED UNTIL D9

## Review protocol

This is a material independent architecture/contract review. Do not optimize for agreement with the candidate, prior GPT reasoning, CodeRabbit, Codex, chat history, or this handoff. Repository authority and executable evidence outrank reviewer preference.

Before analysis, revalidate remote `main`, candidate HEAD, PR #52, changed files, CI and current review threads. If the candidate moved, record the new exact SHA before reviewing.

Read strictly through the repository router:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. the smallest D5 Product OAD authority subpack routed by `docs/index.md`

Default to the bounded pack. Do not recursively read D0–D4, ADR history, old runtime, Git history or broad review dialogue unless a concrete contradiction requires it. Accepted authority is not reopened by preference.

For engineering method / evidence discipline, also use the repository-routed local guide:

- `docs/development/evidence-grounded-production-engineering-for-llm-agents.md`

and the canonical organizational authorities referenced by `AGENTS.md`:

- `developmentconexus-ops/conexus-methodology/METHOD.md`
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md`

## Candidate under review

The candidate is the canonical Product OpenAPI Description authoring/proof for D5 only.

Expected fixed contract:

- exactly 95 Product operations;
- exactly 29 ordinary Permissions;
- Principal kinds H / A / S only;
- stable Product Problem origin `https://conexus.fun`;
- Technical Ingress excluded from Product OAD/SDK;
- active legacy runtime population 0;
- no Product runtime/router/database/deployment selection;
- implementation blocked until D9.

PR #52 currently changes only:

- `.github/workflows/ci.yml`
- `contracts/api/product/components.yaml`
- `contracts/api/product/openapi.yaml`
- `contracts/api/product/paths-economics-governance-sales-materialization.yaml`
- `contracts/api/product/paths-fulfillment-postsale-work.yaml`
- `contracts/api/product/paths-identity-portfolio-readiness.yaml`
- `contracts/api/product/paths-offering-availability-market.yaml`
- `contracts/api/product/redocly.yaml`
- `scripts/verify-product-oad.mjs`

No accepted D0–D5 authority document is changed by the candidate.

The latest executable proof at handoff is green and reports:

- OAS 3.1.2;
- 95/95 Product operations;
- 29/29 ordinary Permissions;
- H/A/S only;
- stable origin `https://conexus.fun`;
- 14/14 mandatory Idempotency-Key carrier operations;
- 26/26 collection operations;
- 10/10 Product OAD negative controls;
- generated TypeScript and Go projection semantics PASS;
- exact Go minimum 1.25.1 plus current Go 1.27.0 proof;
- active legacy runtime population 0;
- repository full gate PASS;
- runtime schema enforcement explicitly NOT CLAIMED — D7;
- router selection explicitly NONE — D7.

Do not accept those claims merely because CI is green. Determine whether the proof is capable of falsifying the accepted contract rather than only proving its own fixture.

## Required independent review

### 1. Product/architecture coherence

Check whether the candidate faithfully projects the accepted D5 semantics and all imported D0–D4 invariants without inventing new authority. In particular challenge:

- Organization scope and cross-Organization fail-closed behavior;
- Principal / Permission / business-authority separation;
- source-qualified external identity;
- Q/C/P knowledge, freshness, ambiguity and convergence semantics;
- idempotency and revision-precondition carrier choices;
- Problem Details identity/status/media grammar;
- collection/cursor grammar;
- Technical Ingress and authored-media delivery exclusion;
- exact W4 operation / class / Permission / Principal mapping;
- absence of a second wire authority or hand-authored Product SDK.

If you find a contradiction, identify the smallest owning authority that must reopen. Do not reopen the whole D5 package by default.

### 2. Method / YAGNI / overengineering

Review the work against the DevelopmentConexus Method and the evidence-grounded production-engineering guide.

Explicitly decide whether the candidate is:

- the smallest sufficient contract/proof;
- over-engineered, under-proved, or duplicating authority;
- carrying speculative abstraction for future providers/clients;
- accidentally selecting D6/D7 mechanics early;
- retaining unnecessary compatibility/legacy cost;
- using generated/proof machinery that protects real accepted properties rather than implementation aesthetics.

Do not recommend refactors, generic platforms, additional layers, or new technology merely because they look more professional. Every expansion needs a concrete accepted consumer, failure class or falsifier.

### 3. Technology/reference freshness review

Independently verify the bounded tooling choices using current primary upstream evidence, not memory or secondary summaries.

Current accepted tooling baseline to challenge:

- OpenAPI 3.1.2;
- Redocly CLI 2.45.0;
- openapi-typescript 7.13.0;
- oapi-codegen v2.8.0;
- `github.com/oapi-codegen/runtime` v1.7.0;
- Go minimum 1.25.1 and current/forward proof on Go 1.27.0.

Use official specifications, official documentation, official release notes and upstream source repositories as primary references. Check current support/compatibility/security status and whether the pinned baseline still preserves the accepted properties.

Important: newer does not mean better. Do **not** reopen a pinned tool/version only because a newer release exists. Recommend a change only if current primary evidence shows a material correctness, compatibility, security, maintenance or proof problem for this contract.

Do not select a Go router/server framework, database, deployment model, runtime validator, frontend transport, BFF topology or other D6/D7 mechanism. If evidence creates a constraint for D7, route it explicitly as a later obligation rather than implementing/selecting it now.

### 4. Executable-proof quality

Challenge `scripts/verify-product-oad.mjs` and CI as falsifiers.

Check for at least:

- real source-closure traversal;
- exact 95-operation equality, not count-only proof;
- exact W4 class/Permission/Principal mapping;
- local refs only and extension allowlist enforcement;
- deterministic bundle and generated projections;
- representative preservation of `oneOf + const`, exact decimal strings, knowledge-state unions, multipart file+etag and parent ListingIntent validator;
- mandatory idempotency carrier semantics, not header-name presence only;
- ETag / quantity grammar strength and negative controls;
- no checked-in generated second authority;
- no false claim that generated Go types prove D7 runtime rejection;
- real minimum Go toolchain execution;
- negative controls that would actually fail if the protected property regressed.

Known review dialogue may contain correct findings and false positives. Perform an independent pass first, then inspect prior review comments only as secondary challenge material. One current reviewer claim specifically asks whether ETag/non-negative-quantity regex regression controls are strong enough; adjudicate it rather than inheriting it.

### 5. Rebuild / technology-reselection decision

Give an explicit answer to this question:

> Based on accepted Product authority, current executable proof and current primary technology evidence, is there any material reason to reconstruct this Product OAD approach or reselect the bounded D5 tooling before merge?

Possible answers must be evidence-backed:

- `NO — candidate architecture/tooling remains sufficient`;
- `BOUNDED D5 FIXES ONLY`;
- `REOPEN <smallest authority>`;
- `RECONSTRUCT` only if a material falsifier proves the current approach cannot preserve accepted meaning.

A preference for different files, prettier code, a newer release, a popular framework or another reviewer style is not a reconstruction reason.

## Output contract

Append your review below `## Fable response` in this file. Do not modify any other file on `review/d5-fable`.

Use this structure:

1. **Verdict** — `ACCEPT`, `ACCEPT WITH BOUNDED FIXES`, `REOPEN SMALLEST AUTHORITY`, or `REJECT`.
2. **Material findings** — numbered, highest severity first. For each finding provide:
   - classification: `D5_FIX`, `D7_OBLIGATION`, `LATER_NON_BLOCKING`, `AUTHORITY_CONTRADICTION`, or `REVIEW_FALSE_POSITIVE`;
   - exact candidate location;
   - governing repository authority;
   - external primary evidence when technology-dependent;
   - concrete failure/counterexample;
   - smallest correction and why it belongs in that phase.
3. **Method/YAGNI assessment** — explicit conclusion on overengineering/underengineering.
4. **Technology freshness assessment** — explicit conclusion for each bounded tool/version above and whether any reselection is justified now.
5. **Executable-proof assessment** — what the proof genuinely establishes and what remains D7/D8.
6. **Reconstruction decision** — explicit yes/no and evidence.
7. **Final merge recommendation** — what, if anything, must change before operator ratification/merge.

If there are no material findings, say so explicitly. Do not manufacture stylistic findings to fill the review.

## Interaction rule

Fable writes only to this dialogue file. GPT will independently adjudicate every finding against repository authority before changing the candidate. No Fable suggestion directly authorizes candidate changes, phase progression, merge, D6/D7 work or Product implementation.

---

## Fable response

<!-- Fable: append independent review here. -->
