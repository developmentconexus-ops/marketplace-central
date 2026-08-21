# D6 — Final Independent Fable Challenge

> Review branch only: `review/d6-final-fable`
> Candidate branch: `stage/d6-frontend`
> Candidate HEAD expected at handoff: `8273da58155828c1bcd34305c0398ce778247c29`
> Candidate PR expected: #54 — `docs(d6): open frontend authority stage`
> Base `main` expected: `9d2c81e175bc39ac388c9d8924ddad21f2a86480`
> Expected Product candidate: 99 Product operations / 30 ordinary Permissions / Principal kinds H-A-S only
> D6-R1: approved/reviewed
> D6-B1: operator-ratified
> D6-B2: operator-ratified
> D7–D9: blocked
> Product implementation: blocked until D9

## Purpose

Run the isolated **final D6 adversarial challenge** against the exact current candidate before any D6 closeout or D7 opening.

This is not a preference review. Attempt to falsify the combined D5-R1 + D6-R1 + D6-B1 + D6-B2 candidate as one coherent frontend architecture and Product consumer. Search for material contradictions, hidden authority, missing semantics, proof blind spots, YAGNI violations and premature D7 mechanism choices.

Do not optimize for agreement with GPT, the operator, prior Fable reviews, CodeRabbit, prior chat history or the candidate. Reconstruct the reasoning independently from current repository authority and proportional primary evidence.

Reviewer output is **Evidence, not authority**. Do not edit the candidate branch or PR #54. Write only below `## Fable response` in this file on `review/d6-final-fable`.

## Mandatory revalidation before analysis

Revalidate independently and record the exact values:

1. remote `main` HEAD;
2. `stage/d6-frontend` HEAD;
3. PR #54 base/head/state/draft/mergeability;
4. changed files and count;
5. GitHub Actions / required CI on the exact candidate HEAD;
6. this review branch ancestry/tree relation to the exact candidate;
7. that `candidate...review` differs by **only** `docs/work/current/ai-dialog.md`.

If the candidate HEAD differs from `8273da58155828c1bcd34305c0398ce778247c29`, **stop** and report the stale-review condition. Do not review a moved candidate.

## Strict reading discipline

Start exactly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D6-FRONTEND.md`
5. `ARCHITECTURE.md`

Then switch bounded packs only when a concrete falsifier needs them. Do not recursively ingest the rebaseline.

For the final combined challenge, the most likely bounded owners are:

- `docs/engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md` for human-session / A-S bearer semantics;
- `docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md` for the bounded 99/30 Performance amendment;
- `docs/engineering/rebaseline/D6-B1-INTERACTION-MAP.md` and `qualification/d6-wireframes/index.html` for 99-operation interaction proof;
- `contracts/api/product/openapi.yaml` when exact Product wire meaning is material;
- `docs/development/engineering-rules.md` for the review-isolation/proof protocol.

Use accepted D0–D5 narrative only when these owners plus the OAD cannot answer a concrete semantic question. Use current official upstream evidence only where a technology/provider-dependent claim materially depends on it.

## Candidate invariants to attack

### 1. Authority and topology

Challenge whether the ratified hybrid actually preserves:

```text
app/routes
  -> features/<human-lens-or-flow>
  -> api/<owner-or-operation-family>
  -> api/transport
  -> api/generated

features -> ui
```

Attempt to prove that this shape creates any of the following:

- feature-to-feature semantic coupling;
- owner-adapter-to-owner-adapter private coupling;
- screen-shaped/business API authority in features;
- a generic workflow/action/client-domain layer by another name;
- business decisions in transport;
- generated OpenAPI shapes becoming a client-side domain model;
- duplicated business semantics between lenses and owner adapters;
- unavoidable circular dependencies or ownership ambiguity.

If a problem exists, identify the **smallest** D6 correction. Do not replace the architecture because another folder taxonomy is prettier.

### 2. Human authentication vs machine authentication

Falsify the claim that D6 correctly consumes D5-R1:

```text
H browser -> external OIDC login -> server-side exchange/session -> Secure HttpOnly application session + CSRF on unsafe Product requests
A / S     -> Client Credentials -> audience-bound bearer
```

Attack at least:

- any path that makes browser JavaScript own/store OIDC access or refresh tokens;
- accidental use of `MpcMachineBearerAuth` as the H-browser carrier;
- CSRF confused with identity/authorization;
- hidden D7 commitment to session-store, CSRF bootstrap, Keycloak realm or backend framework mechanics;
- transport abstractions that make the H and A/S client classes indistinguishable where they must differ.

### 3. Product wire discipline

Challenge `openapi-typescript` + `openapi-fetch` as the smallest safe profile.

Attempt to find an admitted Product wire behavior that cannot be preserved without unsafe shadow code, including:

- path/query/body serialization;
- `Idempotency-Key`;
- `If-Match` / ETag;
- RFC 9457 Problem Details;
- multipart/custom serialization such as ListingIntent media;
- blob/stream/arrayBuffer/raw-response needs;
- session/CSRF headers;
- source-qualified identities.

A technology-dependent finding must cite current primary upstream evidence. Distinguish an architectural blocker from an implementation-version detail.

### 4. TanStack Router vs TanStack Query ownership

Attempt to prove a second server-state authority is implicit.

Challenge:

- router loaders as a competing data cache;
- organization / installation / period / filters split incorrectly between URL and server state;
- malformed/stale URL state becoming Product truth;
- exact-required Installation silently defaulting;
- Organization switch retaining incompatible cached/navigation state;
- cursor state becoming durable business identity;
- query keys omitting canonical semantic scope inputs.

### 5. Retry, idempotency, ambiguity and concurrency

Falsify the promise that frontend convenience cannot corrupt consequential semantics.

Construct counterexamples for:

- automatic mutation retry after ambiguous potentially accepted write;
- regenerating an idempotency key for the same semantic retry;
- reusing one key for a materially changed request;
- stale `If-Match` / precondition failure collapsed into business/provider rejection;
- generic invalidation/event-bus behavior hiding operation-specific convergence;
- optimistic UI that presents accepted/pending/ambiguous as completed/applied.

Decide whether D6 is precise enough without prematurely selecting D7 mechanisms.

### 6. Permission, identity, Organization and source qualification

Attack:

- route/button visibility treated as authorization;
- `performance.read` leaking Market/Economics/Offering/Sales/Availability/write authority;
- `display_name` becoming identity;
- remembered Organization/Installation becoming ambient authority;
- provider/native keys becoming global identifiers;
- SourceInstance/Marketplace Installation qualification being lost in URL, query identity or feature composition;
- cross-owner read composition turning into write authority.

### 7. Knowledge, freshness and outcome honesty

Construct cases where the frontend could collapse:

- known zero vs unknown;
- empty vs partial;
- unavailable vs unsupported;
- stale vs current;
- provider-reported evidence vs MPC-authored fact;
- accepted vs pending vs rejected vs ambiguous.

Check whether the chosen state ownership/topology creates pressure to normalize away these distinctions.

### 8. D6-R1 Performance consumption

Challenge the final frontend use of the bounded Performance authority:

- exact Installation requirement;
- current/comparison period handling;
- comparability gates;
- Listing Performance remaining Offering-owned context rather than ownership transfer;
- Retail Media campaign/listing/catalog/family scope preservation;
- provider CVR/ROAS not reconstructed or relabelled;
- source-qualified historical evidence remaining honest;
- no cross-provider KPI equivalence by naming;
- no `Strategy`, generic Analytics/Metric domain, Ads mutation, signals/recommendation or AI authority created in the frontend.

### 9. D6-B1 99-operation interaction coverage

Do not trust the `99/99` statement by itself. Independently challenge the interaction map and wireframes for:

- orphan Product operations;
- duplicated homes that imply multiple client authorities;
- material operations technically mapped but not realistically reachable;
- hidden #100 operation implied by UX symmetry;
- fake generic Refresh/Sync/Close/Connect/SetPrice/SetAvailableQuantity actions;
- Technical Ingress represented as Product operation;
- exact-required contextual qualifiers missing from user flows;
- hidden write capability behind read-only strategic lenses.

The goal is semantic reachability/coherence, not visual polish.

### 10. D6-B2 dependency necessity / YAGNI

For each selected dependency, require a current proven consumer/property:

- React + TypeScript strict;
- TanStack Query;
- TanStack Router;
- `openapi-typescript`;
- `openapi-fetch`.

Attempt to prove one is unnecessary, overlapping or materially insufficient. Conversely, attempt to prove one deferred dependency is already required now.

Do **not** recommend Redux/Zustand server mirrors, generic generated query SDKs, Axios, universal form/schema/UI frameworks, microfrontends, SSR, offline-first, realtime, design-system platforms or meta-frameworks without a concrete falsifier proving the accepted smaller profile insufficient.

### 11. D7 mechanism leakage

Search the entire D6 candidate diff for hidden commitment to:

- server HTTP router/mux/framework;
- database/schema/RLS realization;
- transaction/outbox implementation;
- scheduler/worker/queue;
- Keycloak realm/session-store topology;
- deployment/process topology;
- exact dependency versions without a version-specific architectural reason;
- lint/plugin implementation details prematurely frozen.

Classify true future obligations as `D7_OBLIGATION`, not D6 blockers, unless D6 already depends on the mechanism.

### 12. Proof quality

Challenge whether the green candidate proves what it claims.

At minimum inspect whether current executable controls genuinely establish:

- accepted D5 baseline 95/29 non-regression;
- D6-R1 99/30 surface;
- D5-R1 auth profile;
- generated TypeScript/Go projection determinism/compilability;
- Performance knowledge controls;
- review isolation;
- no active legacy runtime.

Look for tautological/self-referential fixtures, subject-population blind spots, bundle pruning that hides source drift, or negative controls that never exercise the real enforcement path.

Do not turn generated-contract proof into a claim of D7 runtime behavior.

## Targeted falsifiers

Attempt at least these statements and record PASS/FAIL/INSUFFICIENT-EVIDENCE with reasoning:

1. H-browser code never needs OIDC bearer/refresh-token custody.
2. A/S bearer support does not leak into H-browser auth architecture.
3. CSRF is preserved without selecting D7 realization.
4. Every one of 99 Product operations has one coherent human interaction home where applicable.
5. No 100th Product operation is implied by the UX.
6. Exact Organization/Installation/source context cannot silently default.
7. `display_name` cannot become identity/authorization authority.
8. `performance.read` remains isolated.
9. Performance does not steal Market/Economics/Sales/Offering/Availability meaning.
10. Retail Media scope cannot collapse to Listing by convenience.
11. Provider CVR/ROAS cannot be reconstructed client-side.
12. Unknown/partial/unavailable/stale evidence cannot collapse into zero/complete/current.
13. Accepted/pending/ambiguous outcomes cannot collapse into completed.
14. Consequential mutations cannot be blindly retried.
15. Idempotency identity survives same-semantic-request retry and changes for materially different requests.
16. Concurrency/precondition failure stays distinct.
17. Router cannot become a second server-state cache.
18. TanStack Query key identity includes every semantic scope input.
19. `openapi-typescript` remains the only generated wire shape authority.
20. `openapi-fetch` remains a thin Product transport and not a business abstraction.
21. No handwritten DTO/schema shadow authority is required.
22. No cross-owner adapter dependency is required for coherent screens.
23. No generic client workflow/action bus is required.
24. No D7 mechanism has been selected by D6.
25. No deferred framework/dependency has a proven current consumer.
26. Current selected dependencies each have a proven current consumer.
27. Product implementation remains blocked until D9.
28. The final D6 design is proportionate/YAGNI under current Product 1.0 needs.

## Output contract

Append below `## Fable response` using this structure.

### 1. Verdict

Choose exactly one:

- `ACCEPT`
- `ACCEPT WITH BOUNDED FIXES`
- `REOPEN SMALLEST AUTHORITY`
- `REJECT / RECONSTRUCT`

### 2. Revalidation record

Record exact main SHA, candidate SHA, PR state/base/head, changed-file count, CI status and review-isolation result.

### 3. Executive coherence assessment

Concisely assess Product/frontend coherence, DevelopmentConexus Method alignment, YAGNI/proportionality and whether the candidate is ready for D6 closeout adjudication.

### 4. Material findings

Number findings highest severity first. For each include:

- **classification:** `D5_FIX`, `D6_FIX`, `D7_OBLIGATION`, `REPOSITORY_FIX`, `LATER_NON_BLOCKING`, `AUTHORITY_CONTRADICTION`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- current primary external evidence when technology/provider-dependent;
- concrete counterexample/failure;
- smallest correction;
- why it belongs in that exact stage/owner.

If there are no material findings, state that explicitly. Do not manufacture style findings.

### 5. Authentication assessment

Explicitly adjudicate H session+CSRF vs A/S bearer consumption and D7 boundary.

### 6. Interaction / Product-surface assessment

Explicitly adjudicate 99-operation coverage, contextual qualifiers, permission/identity/knowledge/outcome honesty and hidden-operation risk.

### 7. Frontend topology / dependency assessment

Explicitly adjudicate lens/flow + stateless owner-adapter topology, TanStack Router/Query split, generated-wire discipline and selected/deferred dependencies.

### 8. Retry / concurrency assessment

Explicitly adjudicate idempotency, ambiguous outcomes, preconditions/concurrency and client retry/invalidation behavior.

### 9. Performance assessment

Explicitly adjudicate D6-R1 frontend consumption, owner boundaries, Retail Media scope, comparability/provenance and absence of Analytics/Ads/AI authority creep.

### 10. Proof assessment

State what the current green proof genuinely establishes, what it does not establish and any material blind spots.

### 11. D7 leakage / future obligations

Separate actual D6 defects from correct D7 obligations.

### 12. Reconstruction decision

Answer explicitly: is there any material reason to reconstruct/reopen broader accepted architecture before D6 closeout?

### 13. Continuation recommendation

State the smallest exact action after GPT adjudication. No review output authorizes merge, D7 opening or Product implementation.

---

## Interaction rule

Fable writes **only** to this file on `review/d6-final-fable`. Do not edit PR #54, `stage/d6-frontend`, `main` or any other review branch.

GPT will independently adjudicate every material finding against current repository authority and executable evidence. Round 2 is justified only for a surviving material contradiction after adjudication/bounded fixes.

---

## Fable response

<!-- Claude Fable 5: append independent final D6 review here. -->

### 1. Verdict

**ACCEPT WITH BOUNDED FIXES**

No material reason exists to reopen D5-R1, D6-R1, D6-B1 or D6-B2 semantics, the 13th boundary, the 99/30 surface, the auth split, the lens/owner-adapter topology or the selected dependency profile. The bounded fixes below are repository/verification hygiene, not architecture.

### 2. Revalidation record

Independently revalidated from a fresh session before analysis:

```text
remote main HEAD                 9d2c81e175bc39ac388c9d8924ddad21f2a86480
stage/d6-frontend HEAD           8273da58155828c1bcd34305c0398ce778247c29  (matches expected; not stale)
PR #54                           open · Draft · base main@9d2c81e1 · head stage/d6-frontend@8273da58 · mergeable_state=clean
changed files (main...candidate) 19 files · +4461 / −781 (matches PR metadata)
GitHub Actions @ 8273da58        ci / required = SUCCESS · pr-title / conventional-title = SUCCESS
review branch                    review/d6-final-fable @ 7d03e327; parent commit is exactly 8273da58
review isolation                 git diff candidate..review = only docs/work/current/ai-dialog.md (+371)
candidate tree hygiene           0 tracked docs/work/** paths on stage/d6-frontend; bootstrap pack 19994/20480 bytes
```

Independent proof re-execution in this session (node 22.x, go1.24.7 host with GOTOOLCHAIN go1.25.1; pwsh unavailable here, so `scripts/gate.ps1` itself was exercised only by CI on the exact HEAD):

```text
node scripts/verify-product-oad.mjs                     exit 0
  product_oad_baseline_non_regression=PASS              95/95 ops · 29/29 Permissions · 12/12 negative controls
  product_oad_performance_repair=PASS                   99/99 ops · 30/30 Permissions · 28/28 List/Search · 7/7 negative controls
  product_oad_auth_profile=PASS                         5/5 auth negative controls · deterministic TS + Go, both compile
  go_minimum=go1.25.1                                   colon-suffix custom-mux wrapper test executed
node scripts/verify-performance-evidence-knowledge.mjs  exit 0 · 2/2 negative controls PASS
```

Direct census performed independently of the verifiers (raw grep over path files, fixture block excluded): exactly 99 operationIds, 59 Q + 40 C, 30 ordinary Permissions + `authenticated`, principal-kind vectors only over {H, A, S} (57×[H,A,S], 14×[H,A], 2×[H,S], 26×[H]), exactly 4 Performance Qs under `performance.read`, all Installation-nested.

### 3. Executive coherence assessment

The combined D5-R1 + D6-R1 + D6-B1 + D6-B2 candidate is coherent as one frontend architecture and Product consumer. The interaction authority (B1), the realization topology (B2), the auth carrier split (D5-R1) and the Performance surface (D6-R1) reference each other without circular authority: the OAD stays the single wire authority, the roadmap stays the single status authority, and the frontend never acquires business meaning. Method alignment is good — expectations in the proof chain are pinned to accepted stage documents (W4 §8 matrix, W2 Problem catalog, admission-matrix safety sweep, D6-R1 marker-delimited matrix), not to the artifact under test, so the green result is not self-referential. The dependency profile is proportionate: five dependencies, each with a ratified consumer, and a explicit reject list with reopen triggers. The candidate is ready for D6 closeout adjudication after the bounded fixes below are dispositioned.

### 4. Material findings

Numbered highest severity first. No Critical findings.

**FF-1 — `npm run gate` is broken on the candidate (documented verification entry point fails at parameter binding)**

- **classification:** REPOSITORY_FIX
- **severity:** Important
- **location:** `package.json` scripts (`"gate": "pwsh ... -Lane all && ..."`), `scripts/gate.ps1:2-3` (`[ValidateSet('quick', 'full')] [string]$Lane = 'quick'`)
- **governing authority:** `AGENTS.md` §Verification and `docs/development/engineering-rules.md` §Local verification both document `npm run gate` and `npm run gate:full` as the shared local/CI gate.
- **counterexample:** `main`'s gate.ps1 declared `[ValidateSet('all', 'full')]` with default `'all'`; the candidate rewrote gate.ps1 to `ValidateSet('quick', 'full')` but left `package.json`'s `gate` script passing `-Lane all`. PowerShell parameter binding rejects `all` before the script body runs, so `npm run gate` exits non-zero on every invocation. CI is green only because it invokes `gate:full`. Caveat recorded honestly: pwsh is unavailable in this review environment, so this is established by static analysis of ValidateSet semantics plus the main-vs-candidate diff, not by execution.
- **secondary observation folded in:** `$Lane` is referenced nowhere except the final `Write-Host "gate lane: $Lane"` — `quick` and `full` run an identical body, so the lane distinction is currently decorative.
- **smallest correction:** one line in `package.json`: `-Lane all` → `-Lane quick` (or re-admit `all` in the ValidateSet). Optionally either give the lanes a real behavioral difference or collapse to one lane.
- **why this stage/owner:** repository verification interface, candidate-introduced in this PR's own gate rewrite; no Product/architecture meaning involved. It fails loudly (red, not silently green), so it does not corrupt any proof claim — but a candidate whose own documented default gate command cannot run should not close a stage unfixed.

**FF-2 — D6-B1 artifact status header contradicts the candidate's own ratification claim**

- **classification:** REPOSITORY_FIX
- **severity:** Minor
- **location:** `docs/engineering/rebaseline/D6-B1-INTERACTION-MAP.md:3` — "Status: CANDIDATO DERIVADO — ... **ainda não é ratificação de D6-B1**"
- **governing authority:** `docs/roadmap.md` (sole mutable status authority) and `D6-FRONTEND.md` §8 both state D6-B1 was operator-ratified on 2026-08-21; AGENTS.md requires surfacing material authority conflicts rather than silently choosing.
- **counterexample:** a fresh actor routed into D6-B1 by `docs/index.md` reads, in the owning artifact itself, that the map is not yet ratified — the opposite of the current program truth. The `docs/index.md` frozen-snapshot rule covers *pre-standard* status prose; this header was authored inside the current candidate, so the exemption reads awkwardly here.
- **smallest correction:** one header line in D6-B1 recording the ratified status (or an explicit pointer that status lives only in the roadmap).
- **why this stage/owner:** documentation-consistency inside the D6 candidate; roadmap authority already resolves the conflict, so Minor, not blocking.

**FF-3 — D6-B1 derived screen-state total is arithmetically wrong (40 claimed, 39 present)**

- **classification:** REPOSITORY_FIX
- **severity:** Minor
- **location:** `docs/engineering/rebaseline/D6-B1-INTERACTION-MAP.md` §5 — "Total derivado: 40 estados de tela/rota"
- **counterexample:** the §5 table contains exactly 39 `S..` rows (mechanically counted: `grep -c '^| S[0-9]'` = 39; manual enumeration agrees).
- **smallest correction:** 40 → 39.
- **why it matters at all:** the load-bearing numbers (99/99 operation coverage, per-owner group counts) were independently reconciled and are exact — see §6 below — so this is a non-load-bearing derived count, and B1 itself declares the state count is not an implementation target. Still, an operator-ratified proof artifact should not contain a wrong derived total.

**FF-4 — comparison-available schemas allow silently absent change objects**

- **classification:** LATER_NON_BLOCKING
- **severity:** Minor
- **location:** `contracts/api/product/paths-performance.yaml` — `TrafficPerformanceComparisonAvailable.visits_change`, `SalesActivityPerformanceComparisonAvailable.sales_count_change/units_count_change`, `RetailMediaPerformanceComparisonAvailable.changes`, `RetailMediaSummaryComparisonAvailable.spend_change/return_on_ad_spend_change` are all optional.
- **governing authority:** D6-R1 §9 (W2 comparisons): "A change is emitted only when current/comparison evidence is sufficiently complete and the measurement basis is comparable. **Otherwise the change is explicitly `insufficient_evidence` or `not_comparable`**."
- **counterexample:** the change unions already carry explicit `comparable | insufficient_evidence | not_comparable` states, but because the field itself is optional, a server may return both periods available and omit the change entirely — a fourth, undeclared "no statement" state the W2 text says should be explicit. No client-side dishonesty can result (the D6-B1 law "delta only when the Product returns comparable" fails safe to showing no delta), which is why this is non-blocking.
- **smallest correction if adjudication wants precision now:** make the change discriminator required inside each `*ComparisonAvailable` branch (schema-only tightening inside the already-open D6-R1 surface; no operation, Permission or ownership change; the gate re-proves in minutes). Otherwise document omission-semantics explicitly and defer.

**FF-5 — two cosmetic proof-reporting observations**

- **classification:** LATER_NON_BLOCKING
- **severity:** Minor
- (a) `scripts/gate.ps1:335` prints `legacy_runtime_population: 0` as a hardcoded literal; the actual guard is the forbidden-roots tracked-file scan above it (which does Fail when populated, and which `verify-product-oad-baseline.mjs` duplicates executably), so the claim is enforced — but the printed number is not computed from the check it reports.
- (b) The bootstrap authority pack sits at 19994/20480 bytes — 486 bytes of headroom before the gate's own cap turns any small AGENTS/index/roadmap edit red. Not a defect; recorded so the next docs change is not surprised.

No other material findings. I did not manufacture style findings; everything else attacked below held.

### 5. Authentication assessment

**Adjudication: the candidate correctly consumes D5-R1; no reopen.**

- The OAD root admits exactly `MpcHumanSessionAuth` (apiKey/cookie `__Host-mpc_session`) OR `MpcMachineBearerAuth` (http bearer, no bearerFormat), with zero operation-local `security` overrides — verified both by executing `validateAuth` and by reading the source. The `x-mpc-authentication-profile` extension pins H→session/CSRF (method set DELETE/PATCH/POST/PUT), A/S→client-credentials bearer, `browser_oidc_token_exposure: none`, and full `__Host-` cookie profile. All five auth negative controls genuinely fire (I mutated nothing myself; the verifier's mutations were re-executed and each failed as required).
- **No path makes browser JavaScript an OIDC token holder:** D6-B2 §9.6 binds the human React client to the session carrier and explicitly forbids ordinary browser code injecting/storing a bearer or refresh token; D6-FRONTEND imported invariant 13 repeats it; the wireframes expose no token surface.
- **`MpcMachineBearerAuth` as accidental H carrier:** the wire contract cannot itself refuse a bearer-authenticated human — that refusal is W4 step-6 principal-kind/carrier enforcement, which is runtime and therefore a **D7 proof obligation**, not a D6 defect. The projection-level guards that exist (machine profile limited to [A,S], negative control 4 fails if H is added) are the strongest available at this stage.
- **CSRF is correctly scoped as request trust, not identity/Permission** (D5-R1 §3.1, D6-B2 §9.6, B1 §7.3), and no CSRF bootstrap/storage/rotation mechanism is selected — correctly left to D7. The generated types intentionally do not model the CSRF header; injecting it belongs to the one transport's middleware, which upstream `openapi-fetch` documentation confirms is supported (`onRequest` header mutation).
- **No hidden D7 commitment found:** no session-store, Keycloak realm/client, serving topology or OIDC library appears anywhere in the diff. Keycloak remains "preferred first candidate" only, which pre-dates this candidate.
- The transport keeps H and A/S client classes distinguishable where they must differ: the SPA is only ever the session-class client; machine clients are separate consumers of the same OAD, not the same transport instance.

### 6. Interaction / Product-surface assessment

**Adjudication: 99/99 coverage is real, not narrative.**

- I reconciled the B1 §6 coverage table against an independent owner census of the OAD: Identity 5, Portfolio 6, Readiness 5, Offering 12, Availability 9, Market 3, Performance 4, Economics 11, Governance 7, Sales 3, Materialization 7, Fulfillment 17, Post-Sale 3, Work 7 — all fourteen groups match exactly and sum to 99.
- Every one of the 99 operationIds is named in B1 either literally or by its explicit compound/pluralized shorthand ("List/Get/CreatePriceIntent", "5 InventorySource + 2 policy ops", "4 AuthorizationDelegation ops", "2 target ops", "Get/Update/Deactivate" etc.); the 11 IDs not matched literally are all accounted for by such compounds, and each compound's arity matches the census. No orphan operation, no duplicated home implying a second client authority.
- **No hidden 100th operation:** the flows resolve every UX need to an existing Q/C or to Technical Non-Product Ingress (OIDC login/logout, marketplace OAuth ceremony, advertiser binding); the shell needs only `GetCurrentAccessContext`; Overview composes existing reads. The wireframes demonstrate absent capabilities as explicitly disabled "não existe" controls (SetAvailableQuantity, close Work, close resolution) rather than hiding them — a good visible negative control.
- The 30-Permission vocabulary used across B1 screens is exactly the OAD census (30 ordinary + `authenticated` for `GetCurrentAccessContext`); `performance.read` maps only to the four Performance Qs; permission-conditioned visibility is declared usability-only with server authority on deep links.
- Exact-required qualifiers survive the flows: Performance, Listings, Sales, Shipments, Preparação all pin exact Installation; Readiness search without `source_instance_id` is a bounded multi-source search, never a hidden default; results stay SourceInstance-qualified.
- Knowledge/outcome honesty is carried into the UX artifact itself: the wireframes render partial/unknown/unavailable/not-comparable states, blocked comparisons (both coverage-blocked and basis-blocked variants), preserved-history provenance and accepted≠applied language. Reachability, not polish, was assessed.
- Two truthful caveats: (i) the low-fi HTML merges Aprovações and Acesso into one demo screen and does not render every one of the 39 B1 states — B1 §9 only requires the 13 representative sets, which are all present; (ii) FF-2/FF-3 above are the only defects found in the B1 artifact.

### 7. Frontend topology / dependency assessment

**Adjudication: the ratified hybrid holds; no smaller safe profile exists among the considered options, and nothing deferred is already required.**

- **Topology:** the closed-world direction (app/routes → features → api/owner → transport → generated; features/routes → ui) admits no cycles by construction; cross-owner composition is pinned to features/routes; owner adapters are stateless consumption modules barred from importing each other's internals and from becoming client domain models; `transport` is restricted to HTTP/session/CSRF/Problem/header mechanics. I attempted the failure modes in the contract (feature-to-feature coupling, screen-shaped API authority, generic workflow layer by another name, generated shapes becoming a domain model, duplicated semantics between lenses and adapters) and each is either structurally impossible under the stated direction or explicitly prohibited with a reopen trigger. The honest gap: none of this is mechanically enforced yet — correctly so, since no source tree exists; the default-deny import rule is a stated obligation for first implementation and should be treated as a **D7/implementation-entry obligation**.
- **Dependencies, each against a proven current consumer:** React+TS strict (pre-accepted client technology, ARCHITECTURE §3), TanStack Query (ADR-021, server-state class in the ratified state model), TanStack Router (the B1 URL grammar requires typed/validated Organization/Installation/period/filter deep links; upstream docs confirm `validateSearch` performs runtime validation whose failure routes to an error boundary rather than becoming state — exactly the "malformed URL fails to a safe navigation state" law), `openapi-typescript` (generated wire shapes, already exercised by the proof chain at 7.13.0), `openapi-fetch` (thin OAD-bound client). 
- **Wire discipline, technology-verified:** I checked the current upstream `openapi-fetch` documentation (openapi-ts.dev API + middleware pages, via current doc index) for every admitted Product wire behavior: path/query/body serialization (`querySerializer`/`pathSerializer`/`bodySerializer`), multipart for `CreateListingIntentMedia` (custom `bodySerializer` returning `FormData`; Content-Type is then browser-set, and the `etag` string part defaults to text/plain per multipart semantics — matching the OAD's encoding declaration), `Idempotency-Key`/`If-Match`/CSRF headers (plain fetch `headers` option and `onRequest` middleware), Problem Details (typed non-2xx response schemas in the generated `paths`), and blob/stream/arrayBuffer (`parseAs: "json" | "text" | "arrayBuffer" | "blob" | "stream"`). No admitted behavior requires unsafe shadow code. Note the OAD currently has **no** binary Product response at all (fulfillment artifact reads return JSON descriptors), so the stream/blob clause is forward cover, not a current dependency.
- **Router vs Query ownership:** loaders are restricted to prewarming the same TanStack Query contract (`ensureQueryData`), which upstream documents as the canonical integration; cursor state is barred from durable identity; Organization switch invalidates incompatible state; query identity must include every canonical semantic input. No second server-state authority is implicit.
- **Deferred set:** I attempted to prove one deferred dependency already necessary. The closest candidates — a form framework for the ListingIntent editor and a schema validator for search params — both fail: form draft is an accepted state class needing no framework at design stage, and TanStack Router's `validateSearch` accepts plain functions. Nothing else has a current consumer. Conversely no selected dependency is removable: dropping the router loses validated URL scope state; dropping openapi-fetch recreates manual path/param drift the OAD proof exists to prevent.

### 8. Retry / concurrency assessment

**Adjudication: precise enough for D6 without selecting D7 mechanisms.**

- The contract grammar carries everything the frontend laws need: 14 mandatory `Idempotency-Key` carriers (census matches the admission matrix exactly), required `If-Match` on the 8 update operations, typed `etag` in the 17 command bodies plus correspondence/authorization/superseded-ref variants, and a Problem catalog that keeps `idempotency-key-required` (400), `idempotency-request-in-progress` (409), `idempotency-key-reused` (422), `precondition-required` (428) and precondition-failed distinct from business rejection.
- The frontend laws consume that grammar correctly: same key only for the same semantic retry, new key for a materially different request, stale precondition ≠ business/provider rejection, no automatic replay of consequential mutations, no global retry policy, no generic invalidation event bus, mutation success invalidating only semantically affected reads, and no optimistic UI presenting accepted/pending/ambiguous as completed (B1 §7.3, B2 §9.7, wireframe sale/fulfillment screens demonstrate the "aceita ≠ aplicada" language and the absence of generic retry buttons).
- What D6 deliberately does not prove — that a client actually behaves this way — cannot be proven before a source tree exists; the laws are stated in enforceable terms (each maps to a code-reviewable/lintable property). No D7 mechanism (retry library, mutation queue, cache persistence) was selected.

### 9. Performance assessment

**Adjudication: D6-R1 consumption is honest; no authority creep.**

- All four Performance operations are exact-Installation-nested Qs under `performance.read`; required half-open period params, optional paired comparison params, canonical limit/cursor only, and no metrics/group_by/dimensions/granularity/sort vocabulary — all re-verified executably (7/7 negative controls) and by direct read.
- Listing Performance stays Offering-context, not ownership transfer: `ListMarketplaceListingPerformance` enumerates all known Listings (anti-survivorship law is in the operation description and schema: population coverage is distinct from per-item evidence), and the wireframe keeps rows with unknown/unavailable evidence in the table.
- Retail Media scope is a closed four-variant union (campaign/listing/catalog-group/family-group) with per-variant native keys; the family-group wireframe row explicitly refuses single-listing attribution; no `{entity_type, entity_id}` genericism; `advertiser_id` remains a D4 technical namespace bound via Technical Ingress under `portfolio.manage`, surfaced as "configuração necessária" knowledge rather than 403.
- Provider CVR/ROAS cannot be reconstructed: rates/multiples arrive as unit-carrying decimal strings with `ProviderMeasurementBasis` (provider_reported origin, optional attribution window, basis revision), comparisons are server-owned with explicit `comparable | insufficient_evidence | not_comparable`, and the knowledge verifier proves available-evidence-requires-measure / unavailable-evidence-forbids-measure structurally (2/2). Historical values remain `preserved_source_evidence` custody with provenance — storage does not convert them to MPC-authored fact.
- No cross-provider KPI equivalence: `marketplace_kind` enum is `mercado_livre` only; the roadmap/B1 forbid "Todos os marketplaces" aggregates; the Overview wireframe deliberately renders Performance as a per-Installation entry point instead of a number.
- No `Strategy` domain, Analytics/Metric vocabulary, Ads mutation, `signals[]`, score, time series or AI authority exists anywhere in the surface (checked by path/schema sweep and by the executable forbidden-name controls). The single residual precision item is FF-4 above.

### 10. Proof assessment

**What the green proof genuinely establishes:** deterministic Redocly lint/bundle of the exact source tree; a 95/29 baseline projection (performance paths stripped by marker, access paths re-pointed to the retained frozen fixtures) validated against expectations parsed from the accepted W4/W2/admission-matrix documents — i.e., pinned to independent stage authority, not to the artifact under test; the 99/30 pre-auth surface validated the same way plus D6-R1's own marker-delimited matrix; the current split-auth profile with five falsifiable negative controls; effective Permission vocabulary in AccessContext/AccessRole reads equal to the operation census; deterministic, strict-compilable TypeScript and buildable+tested Go projections including multipart and colon-suffix route registration actually executed through a mux; source-tree hygiene (no servers, no foreign extensions, no remote refs, single OpenAPI entrypoint, zero legacy runtime population) — 12 + 7 + 5 + 2 + 1 negative controls all of which I re-executed or re-derived.

**What it does not establish (and does not claim):** any runtime behavior — W4 carrier/kind enforcement, CSRF enforcement, session issuance, idempotency conflict behavior, RLS/isolation, Performance data acquisition, measure-by-scope availability. The scripts print `runtime_schema_enforcement=NOT_CLAIMED_D7` and `router_selection=NONE_D7`, which is the honest scope.

**Blind spots found, all bounded:** (i) the baseline projection is synthesized from the *current* shared files, so the 95/29 proof establishes "the current tree still satisfies the accepted D5 authority", not byte-identity with the historically merged tree — acceptable because the expectations are externally pinned, but worth knowing; (ii) the retained access-fixture pathItems in `paths-identity-portfolio-readiness.yaml` duplicate the five live access operations and could drift from them in detail richer than the W4-checked properties — both copies are authority-checked, so drift cannot change proven semantics, but the duplication is a standing maintenance hazard worth a comment-level guard someday; (iii) the repository-level gate has exactly one executable negative control (the review-isolation predicate) — its other invariants are direct guards, which the gate's own comment declares deliberately; (iv) FF-1: the `gate` lane of the verification interface cannot currently run at all, and quick/full are behaviorally identical; (v) `gate.ps1`'s `legacy_runtime_population: 0` line is a printed literal (FF-5a). None of these makes a green claim false.

### 11. D7 leakage / future obligations

**No D6 defect from leakage was found.** Correct **D7_OBLIGATIONs** (all already recorded in candidate authority, listed here for the adjudicator's ledger):

1. W4 step-1/6 runtime enforcement that the machine-bearer carrier cannot authenticate an H Principal (D5-R1 §4).
2. Session persistence/expiry/rotation, CSRF issuance/bootstrap/rotation, OIDC exchange mechanics, Keycloak realm/deployment, same-origin serving topology (D5-R1 §8).
3. Mercado Livre measure-by-scope availability proof for every Retail Media scope; refresh-window/attribution semantics into coverage/basis (D6-R1 §6).
4. Default-deny mechanical import-boundary enforcement when a real frontend source tree exists (roadmap / B2 §9.4) — an implementation-entry obligation.
5. D2-R2 physical persistence/retention/indexing of preserved Performance evidence.

**Not leakage:** the pinned tool versions in the proof scripts (redocly 2.45.0, openapi-typescript 7.13.0, tsc 5.9.3, oapi-codegen 2.8.0, go 1.25.1) are proof-determinism pins, not D7 runtime selections; the Go std-http-server/strict-server generation is expressibility evidence explicitly disclaimed as `router_selection=NONE_D7`; `.node-version`/engines pin the proof toolchain, not a server runtime.

### 12. Reconstruction decision

**No.** There is no material reason to reconstruct or reopen broader accepted architecture before D6 closeout. Every finding above is a bounded repository/precision fix; none touches accepted Product semantics, ownership/trust boundaries, a safety invariant, the admitted API surface or an operator-ratified decision. I specifically attempted and failed to falsify: the auth split (holds against current IETF browser-app BCP direction and admits A/S unchanged), the 13th boundary (no owner overlap found), the four-Q surface (no missing strategic read forced frontend authority), the hybrid topology (no unavoidable cycle or duplication), and the dependency profile (no removable or missing member).

### 13. Continuation recommendation

Smallest exact action after GPT adjudication:

1. Adjudicate FF-1..FF-5; apply on `stage/d6-frontend` at minimum FF-1 (one line in `package.json`) and preferably FF-2/FF-3 (two documentation lines) as bounded fixes; decide FF-4 explicitly (tighten now inside the D6-R1 surface, or record omission-semantics and defer).
2. Re-run `ci / required` (and a local `npm run gate` to prove FF-1's fix) on the new candidate HEAD.
3. Then, and only with explicit operator authorization: close D6 in `docs/roadmap.md` and adjudicate D7 opening. No review output — including this one — authorizes merge of PR #54, D7 opening or Product implementation. Round 2 of this review is warranted only if adjudication surfaces a surviving material contradiction, which I do not currently expect.

— Fable, final independent D6 challenge, 2026-08-21, from candidate `8273da58` via `review/d6-final-fable`.
