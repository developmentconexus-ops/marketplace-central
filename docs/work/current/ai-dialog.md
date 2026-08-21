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

> **Reviewing actor:** Claude Opus 5 (`claude-opus-5`), operator-dispatched into the Fable independent-challenger slot. Recorded explicitly because "independent review" is only meaningful if the independent actor is named. All findings below were produced by a first-principles pass against repository authority and re-measured locally; prior review dialogue was opened only afterwards, as secondary challenge material, and is adjudicated in §2.9.
>
> **Revalidated at review time:** remote `main` `4f614e1e289e817ed6d2ca9bcdaf3b97311e6c3d`; PR #52 head `f9f581e901df384f4dff6acdbabe9ed57b867036` (unchanged from handoff); PR #52 OPEN / MERGEABLE, base `main`; checks `required` PASS, `conventional-title` PASS, CodeRabbit complete. Review branch differs from the candidate by `docs/work/current/ai-dialog.md` only. `origin/docs/d5-product-oad` (`90ba1f2a`) exists but is not the PR head and was not reviewed.
>
> **Local measurement environment:** Windows 11, Go 1.26.4, Node 26.3.0, `@redocly/cli@2.45.0`, `oapi-codegen@v2.8.0`. Every measurement below is reproducible from the candidate tree.

## 1. Verdict

**ACCEPT WITH BOUNDED FIXES.**

The Product OAD faithfully projects the accepted D5 semantics and imports the D0–D4 invariants without inventing authority. The 95-operation / 29-Permission / H-A-S / stable-origin core is not merely counted — it is proved by exact set equality against the owning authority documents, which is a materially stronger design than a fixture. No authority reopen is required, and no reconstruction is justified.

Five bounded defects must land before merge (F-1…F-5). All five live inside files the candidate already changes, except one line of `scripts/gate.ps1`. Four further items are non-blocking, and one is a D7 obligation created by upstream evidence dated one day before this review.

The single most important correction: **three of the properties the proof reports as green are not actually protected, and one accepted grammar element is absent from the canonical artifact entirely.** The proof is strong where it derives expectations from authority and weak exactly where it falls back to text proximity or a self-referential fixture.

## 2. Material findings

### F-1 — the accepted 405 + `Allow` grammar has no wire realization, and its only control is a source-text regex — `D5_FIX`

**Location.** `contracts/api/product/components.yaml:44` (`headers.Allow`), `:58-62` (`responses.MethodNotAllowed`), `:1006` (`schemas.AboutBlank405Problem`); control at `scripts/verify-product-oad.mjs:274`.

**Governing authority.** `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` Problem catalog and the accepted final Problem/media consistency cycle, which ratified RFC 9110 §15.5.6 (`Allow` is required on 405).

**Measured evidence.** All three components are unreferenced by any operation and are therefore **pruned by the Redocly bundle** — they never enter the canonical resolved document:

```text
headers:   source=3   bundled=2   ORPHAN=1 -> Allow
responses: source=34  bundled=33  ORPHAN=1 -> MethodNotAllowed
schemas:   source=284 bundled=282 ORPHAN=2 -> OrganizationMembersCollection, AboutBlank405Problem
```

`grep -c MethodNotAllowed contracts/api/product/paths-*.yaml` returns `0,0,0,0`. `grep -c AboutBlank405Problem product.gen.go` returns `0`. By contrast `AboutBlank413Problem` and `AboutBlank415Problem` **are** present in the bundle and are asserted against the bundle at `verify-product-oad.mjs:275-278`.

**Concrete failure.** The 405 control at line 274 reads `componentsText` — the raw source file — not `document`:

```js
assert(/MethodNotAllowed:[\s\S]{0,1200}?Allow:/.test(componentsText), '405 MethodNotAllowed source response lacks Allow');
```

It had to be special-cased against source text *precisely because the component does not exist in the bundle*. So the check is green today while the property it names is absent from the Product wire contract, from the TypeScript projection and from the Go projection. It is also weak on its own terms: a 1200-character proximity window would accept an `Allow:` belonging to a different response. `OrganizationMembersCollection` is a benign alias collapse (`components.yaml:147` refers to `OrganizationMemberCollectionBase`, confirmed present in the bundle) and is not part of this finding.

**Smallest correction.** Decide and encode, in D5, which of the two the contract means: either reference `MethodNotAllowed` from the operations so the grammar reaches the bundle and both projections, or delete the three orphan components and route 405 + `Allow` explicitly as a D7 runtime-level obligation. Either way, delete the source-text regex and replace it with a bundle-level assertion. This belongs in D5 because it is the wire contract's own content, not a runtime mechanism.

### F-2 — the colon-suffix dispatch fixture is decoupled from the artifact it claims to protect — `D5_FIX`

**Location.** `scripts/verify-product-oad.mjs:428` — the generated `colon_suffix_test.go` / `TestCanonicalColonSuffixDispatch`.

**Governing authority.** `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md:415`: "the runtime must provide partial-segment pattern support **or a compatible implementation of the generated mux interface**. The D5 fixture uses a bounded custom mux **to prove dispatch** without selecting Chi, Echo, Gin, Gorilla or another D7 framework." The authority authorises a custom mux — this finding is not about that choice, which is correct and D7-neutral.

**Concrete failure.** The fixture's `exactSuffixMux.ServeHTTP` reduces to `r.Method == m.method && r.URL.Path == m.path`, where both sides were initialised from the same string literal. It never implements the generated `ServeMux` interface, is never passed as `StdHTTPServerOptions.BaseRouter`, and never calls `HandlerWithOptions`. It touches no generated code at all. It proves only that `net/http` preserves a colon through request-URL parsing — real, but not the property named.

Because the fixture never reads the generated registration, a change in oapi-codegen's pattern emission — percent-encoding the colon, or splitting `{id}:submit` into two segments — would leave **every** current assertion green: models unchanged, `StrictServerInterface` present, boundary operation IDs present, compile clean, fixture passing. The canonical `:verb` URI would have silently lost its generated realization.

**Measured evidence that the surface is live, not theoretical.** Registering the candidate's own generated routes into a real `http.ServeMux` aborts:

```text
REAL http.ServeMux PANICKED on generated registration:
parsing "POST /organizations/{organization_id}/authorization-delegations/{authorization_delegation_id}:revoke":
at offset 64: bad wildcard segment (must end with '}')
```

24 of 80 canonical paths carry a `:verb` suffix; 20 of those sit on a `{wildcard}:verb` segment. Registration panics on the first, so all 95 operations become unreachable. This confirms the constraint already recorded at `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md:415` and `D5-B2-WIRE-CONTRACT.md:315` — **no reopen is needed**, and the panic itself is correctly a D7 concern.

**Smallest correction.** Make the fixture the thing the authority actually names: a bounded type implementing the generated `ServeMux` interface, passed as `StdHTTPServerOptions.BaseRouter` to `HandlerWithOptions`, asserting that the captured pattern set equals the 95 contract paths and that one `{id}:verb` request dispatches to the correct generated wrapper. This stays inside the authorised D7-neutral envelope, costs a few lines, and — unlike the present fixture — would bite.

### F-3 — the accepted cross-Organization privacy invariant has no falsifier — `D5_FIX`

**Location.** `scripts/verify-product-oad.mjs:151-297` (`validate()` — the omission).

**Governing authority.** `D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md:522`: an access-eligible Principal with no Membership in the path Organization receives a "privacy-preserving `404 resource-not-found` so the API does not confirm another Organization's existence". Reinforced by the standing `AGENTS.md` rail "Organization isolation fails closed across Organizations."

**Measured evidence.** The property **holds** in the artifact: all 94 org-scoped operations declare `404` (95 total, `/access-context` excepted). The defect is that nothing protects it. `validate()` never reads `'404'`, and Redocly's `operation-4xx-response: error` is fully satisfied by the universal `401`/`403`.

**Concrete failure.** Delete the `404` response from any single org-scoped operation. The whole proof stays green — lint passes, all 10 negative controls pass, both projections regenerate and compile. The wire contract would then admit a 403-only outcome for a non-member, which discloses the existence of another Organization's resource: exactly the disclosure W4:522 exists to prevent.

**Smallest correction.** Assert `404` on every org-scoped operation and add an eleventh negative control that removes one. This is the highest-value missing control in the package, because Organization isolation is a standing safety rail rather than an ordinary contract detail — and because the neighbouring invariants *are* well protected (measured: all 8 If-Match operations carry `412` + `428` + an `ETag` response header; all 14 Idempotency-Key operations carry `409` + `422`; `401`/`403`/`500` are universal).

### F-4 — the canonical Product OAD proof sits outside the declared shared gate and cannot run on the declared local platform — `D5_FIX`

**Location.** `.github/workflows/ci.yml:33-34`; `scripts/verify-product-oad.mjs:26-27`; `package.json` scripts.

**Governing authority.** `AGENTS.md` §Verification names `npm run gate` / `npm run gate:full` and states "`scripts/gate.ps1` is the shared local/CI gate."

**Measured evidence — not wired.** `grep -rn verify-product-oad` across the tree returns exactly one hit: `.github/workflows/ci.yml:34`. `package.json` defines only `gate` and `gate:full`, both `pwsh … scripts/gate.ps1`. `gate.ps1` invokes no `node` and no `.mjs`. The declared shared gate therefore does not cover the Product OAD at all; CI runs the proof as a *sibling step* to the gate, so local and CI are not the same gate.

**Measured evidence — not runnable locally.** Executing the proof on the operator's platform fails immediately:

```text
Error: npx.cmd failed to start: spawnSync npx.cmd EINVAL
    at redoclyProof (scripts/verify-product-oad.mjs:346:3)
```

Cause reproduced directly: `spawnSync('npx.cmd', ['--version'], {shell:false})` gives `EINVAL`; with `{shell:true}` it succeeds. Node 18.20.2/20.12.0 and later refuse to spawn `.cmd`/`.bat` without a shell (CVE-2024-27980 mitigation); `.node-version` pins `26.3.0`. The `process.platform === 'win32'` branch at lines 26-27 is therefore dead code that advertises a local capability the script does not have.

**Concrete failure.** An operator runs `npm run gate:full` on Windows, sees green, and has proved nothing about the canonical Product wire authority — which is the single artifact this stage exists to produce. The only substrate on which it is ever proved is ubuntu CI.

**Smallest correction.** Invoke the proof from a `gate.ps1` lane so `npm run gate:full` covers it, and spawn npx with `shell: true` on win32 (or resolve `npx-cli.js` and spawn `node` directly, which also avoids Node's `shell:true` argument-escaping deprecation). This belongs in D5 because the candidate is what introduced the second, unshared proof path.

### F-5 — optional `sensitivity` contradicts the contract's own absence grammar — `D5_FIX`

**Location.** `contracts/api/product/components.yaml:880-884`.

```yaml
required: [artifact_key, kind, content_type, recorded_at]
properties: {..., sensitivity: {type: string, enum: [ordinary, restricted]}, ...}
```

**Governing authority.** `AGENTS.md` rail: "Unknown, absent, partial or unavailable facts never become plausible known defaults." `D5-B2-W2-SCHEMA-GRAMMAR.md:1108`: fulfillment descriptors "preserve only needed kind/content type/source/provenance/**sensitivity**" — sensitivity is a named descriptor, not optional decoration.

**Concrete failure.** An artifact arrives without `sensitivity`. A consumer cannot distinguish "classified ordinary" from "never classified", so a **restricted** artifact that omits the label is handled as unclassified. Everywhere else this contract refuses exactly that inference and models unknowability explicitly — `DesiredQuantity`, `ProviderQuantity`, `EconomicMoney`, `Correspondence` and `EconomicsConclusion` all use closed `known`/`unknown`/`unavailable`/`not_applicable` unions with `const` discriminators. This one safety label is the sole place where absence is left to mean whatever the reader assumes.

**Smallest correction.** Add `sensitivity` to `required`, or model it as an explicit knowledge state consistent with the rest of the grammar. Optionality in general is fine in this contract (`detail`, `sku`, `gtin`, `observed_price` are legitimately optional); a safety classification is not an ordinary optional field.

### F-6 — the declared minimum Go toolchain went out of support one day before this review — `D7_OBLIGATION`

**Location.** `scripts/verify-product-oad.mjs:437-441` (`go 1.25.1` module directive; `GOTOOLCHAIN: 'go1.25.1'`); baseline at `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md:417-422`.

**External primary evidence** (`https://go.dev/doc/devel/release`, `https://go.dev/dl/?mode=json`, measured today): a Go major is supported until two newer majors exist. **Go 1.27.0 shipped 2026-08-19**, so supported majors are now **1.26 and 1.27 only** — Go 1.25 is end-of-life. `go.dev/dl?mode=json` returns exactly `go1.27.0` and `go1.26.7` as current stable. The final 1.25 patch is `go1.25.14` (2026-08-19), and the 2026-08-19 point releases carried `net/http` fixes — the very package this contract's routing question lives in.

**Assessment.** `GOTOOLCHAIN=go1.25.1` executes a toolchain that is now both end-of-life and 13 patch releases behind the last 1.25. There is **no live exposure**: no runtime is deployed, no Go tree is active, and implementation is blocked until D9. This is therefore **not a merge blocker**, and it is not a "newer exists" objection — Go 1.27.0 landed after the handoff was authored, so this is newly created evidence.

**Route as a D7 obligation.** When D7 selects the runtime, set the floor to a supported major. Note separately that the minimum-toolchain leg gains nothing from `.1` over `.14`; pinning the earliest patch buys no additional proof and forgoes every 1.25.x fix. CI's `go-version: stable` for the forward leg is **correct** and should not change — authority asks for "current/newer toolchain: forward-compatibility evidence", and `stable` delivers exactly that.

### F-7 — `product_oad_negative_controls=10/10` is a literal, not a measurement — `LATER_NON_BLOCKING`

`scripts/verify-product-oad.mjs:486` prints a hardcoded string. The count is accurate today — 8 `expectFailure` calls at lines 335-342, plus the Redocly unresolved-ref fixture (line 358) and the source-tree quoted-remote fixture (line 461) totals 10. But deleting an `expectFailure` still prints `10/10`. The same applies to `product_oad_principal_kinds`, `product_oad_stable_origin` and `product_oad_generated_projection_semantics=PASS`, while `product_oad_operations=${result.all.length}/95` is properly derived. The PR body quotes these banner lines as evidence, so the literals read as measurements. Smallest fix: increment a counter inside `expectFailure` and print it. This is the repository's own "presence is not execution" class applied to the report itself.

### F-8 — `typescript@5.9.3` is a proof-gating pin outside the ratified tooling baseline — `LATER_NON_BLOCKING`

`scripts/verify-product-oad.mjs:388` runs `npx --yes -p typescript@5.9.3 tsc --noEmit --strict`. The accepted baseline block at `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md:390-395` lists only `@redocly/cli`, `openapi-typescript`, `oapi-codegen` and `oapi-codegen/runtime`. Current upstream `typescript` is **7.0.2** (measured: `npm view typescript version`). The strict-compile leg gates the proof, so its compiler version belongs in the declared baseline. Not a blocker — the input is a `.d.ts` and 5.9.3 is a defensible conservative floor. Route "does the projection compile under the compiler D6 actually selects" as a D6 obligation.

### F-9 — "deterministic bundle/projection" is proved only within a single process — `LATER_NON_BLOCKING`

`redoclyProof` (lines 347-352), `typescriptProof` (368-372) and `goProof` (415-420) each generate twice back-to-back from the same npx-resolved install and compare SHA-256. That detects per-invocation nondeterminism only. `npx --yes @redocly/cli@2.45.0` pins the CLI but lets its transitive tree float, and no lockfile covers these tools — Redocly 2.46.1 was itself nothing but an `@redocly/ajv` range bump, which is direct evidence that transitive resolution moves under a fixed CLI pin. The honest claim is "stable within one run", not "reproducible across time and machines". State it that way, or pin the toolchain in a lockfile if cross-run reproducibility is actually wanted.

### F-10 — `component()` fuzzy name resolution can bind an assertion to the wrong schema — `LATER_NON_BLOCKING`

`scripts/verify-product-oad.mjs:107` falls back to `candidate.endsWith('_' + name) || candidate.endsWith(name)`. Measured: harmless today — the bundle contains **zero** mangled component names, and exactly one schema each ends with `Permission`, `StrongETag` and `Money`. But a future schema such as `DelegatedPermission`, if it sorted first, would silently satisfy `component(document, 'schemas', 'Permission')`, and the 29-Permission equality would then be asserted against the wrong enum while reporting green. The alias-collapse case this fallback was presumably written for (`OrganizationMembersCollection` resolving to `OrganizationMemberCollectionBase`) is not solved by `endsWith` anyway. Smallest fix: require exact names and drop the fallback.

### 2.9 Prior review dialogue — adjudicated — `REVIEW_FALSE_POSITIVE`

I performed the pass above before opening PR #52's threads. Six prior claims were then checked; five are stale against this HEAD, having been remediated in `464de3db` and `6fcc9462`, and one stands (reported above as F-5).

**The claim the handoff specifically asked me to adjudicate — that the ETag / non-negative-quantity regression controls are not strong enough — is REFUTED, by mutation test rather than by argument.**

Grammars, read straight from `components.yaml:89,91` and executed:

```text
StrongETag  ^"[!#-~-ÿ]*"$
  "abc" => true    abc => false    W/"x" => false
NonNegativeExactDecimalString  ^[0-9]+(?:\.[0-9]+)?$
  0 => true   1.25 => true   -1 => false   1e3 => false   .5 => false   +1 => false
```

`StrongETag` is character-for-character RFC 9110 §8.8.3 `opaque-tag` (`etagc = %x21 / %x23-7E / obs-text`), correctly anchored so `W/` cannot precede the opening quote. It admits the empty tag `""`, which RFC 9110 permits (`*etagc`) — conformant, and I do not raise it as a finding.

Mutation test of the controls at `verify-product-oad.mjs:266-272`, applied to the real bundle:

```text
baseline etag              : PASSED
baseline qty               : PASSED
MUT etag -> minLength:3    : FAILED as intended -> StrongETag must be a patterned string
MUT etag -> allow weak     : FAILED as intended -> StrongETag must admit only quoted non-weak entity tags
MUT qty  -> signed decimal : FAILED as intended -> desired quantity must remain exact non-negative decimal string
```

The controls bite on every regression tried. Remaining adjudications:

- **"Keep the bearer token format unspecified"** — already done in `6fcc9462`; `openapi.yaml` carries no `bearerFormat`, and line 165 now *asserts its absence*. Stale.
- **"Constrain StrongETag to strong entity-tag syntax"** — `minLength: 3` no longer exists. Stale.
- **"Reject negative desired availability quantities"** — `DesiredQuantityKnown.quantity` references `NonNegativeExactDecimalString` (`components.yaml:457`). Stale.
- **"Model WorkAssignment states as discriminated variants"** — implemented as a closed `oneOf` over `WorkAssignmentAssigned` / `WorkAssignmentUnassigned` (`components.yaml:943-954`), with the split asserted at lines 262-265. Stale.
- **"Enforce the extension allowlist across the entire OAD"** — now `walk(document, …)` over the whole bundle (line 168), plus a source-text sweep (line 449) and an `x-enum-varnames`-on-`Money` negative control (line 342). Stale.
- **"Decide whether `sensitivity` is optional"** — **stands**; reported as F-5 with independent grounding in the absence-grammar rail.

## 3. Method / YAGNI assessment

**The candidate is close to the smallest sufficient contract and proof, and errs toward under-proving rather than over-building.** No overengineering found.

Positively: 3,519 added lines for 95 operations and 282 schemas is lean. There is no second wire authority, no hand-authored SDK, no checked-in generated artifact (`git ls-files` totals 60 files; exactly one tracked `openapi*.yaml`), no `servers` block, no `/v1` prefix, no provider plugin seam, no speculative multi-marketplace abstraction, and no runtime, router, database or deployment selection. `redocly.yaml` is 14 lines and rules-only. Technical Ingress is excluded by a path regex plus a live negative control (lines 202, 337).

The proof's central design decision is genuinely good and worth protecting: it **parses the owning authority documents** — W4's 95-row matrix, W2's 15-slug Problem catalog, the admission matrix's 14 mandatory-idempotency dispositions — and asserts exact set equality against the OAD, rather than against a hand-copied fixture. That is the difference between a contract test and a tautology, and it is why the count-based claims in the PR body are trustworthy where the banner literals in F-7 are not. Row-count anchors (`=== 95`, `=== 29`, `=== 15`, `=== 14`) correctly prevent silent parser drift.

Under-engineered in exactly three places, all reported: F-1 (a ratified grammar element with no wire realization), F-3 (a safety-rail invariant with no falsifier), F-4 (the proof outside the shared gate).

One decorative control worth deleting: the forbidden-import assertion at line 411. Generated model and interface code cannot import `os/exec`, `unsafe` or `syscall`; the assertion cannot fail, and it inflates the apparent strength of the Go leg.

## 4. Technology freshness assessment

Measured today against upstream registries and official release documentation, not memory.

| Tool | Pinned | Current upstream (measured) | Verdict |
|---|---|---|---|
| OpenAPI | 3.1.2 | 3.1.2 is the latest 3.1.x; 3.2.0 exists (2025-09-19) | **KEEP.** 3.2.0 is not consumable by oapi-codegen v2.8.0 / openapi-typescript 7.13.0. Newer, not better. |
| `@redocly/cli` | 2.45.0 | 2.46.2 (2026-08-19) | **KEEP.** 2.46.x is configurable-rule error messaging, a `no-duplicated-enum-values` fix (that rule is not enabled here), and an `@redocly/ajv` range bump. None touches bundling, `$ref` resolution, 3.1 validation or determinism. No material reason to move. |
| `openapi-typescript` | 7.13.0 | 7.13.0 | **KEEP — already current.** |
| `oapi-codegen` | v2.8.0 | v2.8.0 (2026-07-17) | **KEEP — already current.** |
| `github.com/oapi-codegen/runtime` | v1.7.0 | v1.7.0 (2026-08-16) | **KEEP — already current.** Version asserted at line 435. |
| Go minimum | 1.25.1 | **1.25 EOL since Go 1.27.0, 2026-08-19**; supported = 1.26, 1.27; last 1.25 patch = 1.25.14 | **D7 obligation — see F-6.** Not a D5 blocker. |
| Go current/forward | `stable`, resolved to 1.27.0 | 1.27.0 | **KEEP.** `stable` is the right choice for a forward-compatibility leg; the reported `1.27.0` is an observation of that run, not a pin, and should not be read as a fixed property. |
| `typescript` | 5.9.3 (**undeclared**) | 7.0.2 | **See F-8** — declare it in the baseline. |

**No tooling reselection is justified now.** Four of the five ratified pins are exactly at current upstream; the fifth (Redocly) is one patch line behind with no material delta. The only genuine freshness event is Go 1.25's end of support, which post-dates the handoff and is correctly a D7 concern.

## 5. Executable-proof assessment

**What the proof genuinely establishes.** Exact operation-set equality against W4 (not a count) with per-operation class, Permission, Principal-kind and physical-qualification equality; the 29-Permission enum as a set; H/A/S closure; the 15-slug Problem catalog under the stable `https://conexus.fun` origin; 14 mandatory Idempotency-Key carriers with exactly one required header each; the 8-member If-Match carrier set; typed-`etag` requirements across 17 custom methods, both correspondence verbs, every `AuthorizationTarget` branch and the superseded-PriceIntent ref; 26 collection operations with `limit`/`cursor`, no offset/sort/total grammar, and `next_cursor` never required at exhaustion; closed objects and `const`-discriminated `oneOf` branches throughout (`allOf` count is 0, so the classic `allOf` + `additionalProperties:false` footgun is absent); RFC-exact decimal, non-negative-quantity and strong-ETag grammars, all mutation-tested above; local-refs-only with an extension allowlist over the whole bundle; deterministic-within-run bundle and both projections; multipart `file` + typed `etag` and the parent `listing_intent_etag` surviving into both TypeScript and Go; real execution under `GOTOOLCHAIN=go1.25.1` (verified by parsing `go version` output, not by assumption); zero legacy runtime population by `git ls-files` prefix sweep; and a repository-cleanliness check that the proof did not dirty the tree.

That is a substantially stronger falsifier than most stage gates, and the eight semantic negative controls are real — each mutates the bundle and requires `validate()` to reject.

**What it does not establish, despite appearances.** The 405 + `Allow` grammar (F-1 — checked against source text for a component absent from the canonical artifact). Colon-suffix dispatch of the generated handler (F-2 — the fixture never touches generated code). Cross-Organization 404 privacy (F-3 — no assertion at all). Any property at all on the operator's local platform (F-4). The negative-control count (F-7 — a literal).

**Correctly not claimed.** The proof is honest where it matters most: `product_oad_runtime_schema_enforcement=NOT_CLAIMED_D7` and `product_oad_router_selection=NONE_D7`. Generated Go types do **not** prove runtime rejection, and the script never asserts that they do. `std-http-server: true` is a necessary input to obtain `StrictServerInterface` from oapi-codegen, not a router selection — though F-2's fix should make that distinction provable rather than merely stated.

**Remains D7/D8.** Runtime schema validation and rejection; router selection satisfying the measured `{id}:verb` constraint; the supported Go floor; persistence, transactions and idempotent-intake behavior; real Mercado Livre and Sankhya integration evidence; and every golden-flow property, none of which a wire contract can carry.

## 6. Reconstruction decision

> Is there any material reason to reconstruct this Product OAD approach or reselect the bounded D5 tooling before merge?

**NO — the candidate architecture and tooling remain sufficient.**

Evidence: the accepted 95/29/H-A-S/stable-origin contract is projected exactly and is proved by set equality against the owning authority rather than against a fixture; no second wire authority, hand-authored SDK or checked-in generated artifact exists; no D6/D7 mechanism is selected; the source closure is local-ref-only with a working remote-ref negative control; four of five ratified tool pins sit exactly at current upstream and the fifth has no material delta. Every defect found is a bounded correction inside files the candidate already touches — none requires a different structure, a different generator, an additional layer or a new technology.

**No authority reopen is required.** Both facts that might have justified one are already ratified: the `net/http.ServeMux` `{id}:verb` constraint at `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md:415` and `D5-B2-WIRE-CONTRACT.md:315`, and the 404 cross-Organization privacy rule at `D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md:522`. The candidate's gaps are failures to *enforce* accepted authority, not contradictions of it. The smallest owning authority for every finding is the candidate itself.

## 7. Final merge recommendation

**Merge after F-1 through F-5. Do not merge on the current proof alone** — not because the artifact is wrong, but because three of its reported-green properties are unprotected and one accepted grammar element is missing from the canonical document.

Blocking, in order:

1. **F-1** — decide 405 + `Allow`: wire `MethodNotAllowed` into the operations so it reaches the bundle and both projections, **or** delete the three orphan components and route the obligation to D7. Replace the source-text regex with a bundle-level assertion either way.
2. **F-2** — rebuild the colon-suffix fixture as a `ServeMux`-interface implementation passed as `BaseRouter` to the generated `HandlerWithOptions`; assert the registered pattern set equals the 95 contract paths and dispatch one `{id}:verb` request through it.
3. **F-3** — assert `404` on all 94 org-scoped operations and add the corresponding negative control.
4. **F-4** — run the proof from a `gate.ps1` lane, and fix the win32 spawn so the declared local gate actually executes it.
5. **F-5** — make `sensitivity` required, or model it as an explicit knowledge state.

Non-blocking, land with or after the above: **F-7** (derive the negative-control count), **F-8** (declare the TypeScript pin), **F-9** (state determinism precisely), **F-10** (exact component names), and deletion of the decorative forbidden-import assert at line 411.

Record as a routed later obligation: **F-6** — supported Go floor at D7, given Go 1.25's end of support on 2026-08-19.

No change to any accepted D0–D5 artifact is required, and D6–D9 remain closed. Product implementation stays blocked until `docs/roadmap.md` records D9 as accepted.

Nothing in this review authorizes candidate changes, phase progression, merge or implementation; GPT adjudicates each finding against repository authority.
