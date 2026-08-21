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
