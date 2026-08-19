# D5-B2 — Final Problem / Media Consistency Consolidated Review Candidate

> **Status:** NON-AUTHORITATIVE CONSOLIDATED REVIEW EVIDENCE — OPERATOR RATIFICATION NEXT  
> **Parent stage:** D5-B2 Product Operation / Resource Surface  
> **Canonical inputs:** accepted D0→D4 + D4-R1 + D5-B1 + Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Lead review:** complete  
> **Independent Fable review:** complete  
> **GPT adjudication:** converged  
> **Round 2:** not required  
> **Consolidated:** 2026-08-19

## 1. Purpose and authority fence

This candidate consolidates the final D5-B2 consistency review over:

1. W2 Product Problem Details versus provider/OAuth protocol-local failures;
2. ListingIntent-scoped authored-media intake, representation, identity, revision and presentation access;
3. W1/W2/W3/W4/Technical Ingress terminology and routing;
4. Product OpenAPI/SDK exclusion of technical ingress and media-delivery mechanics;
5. duplicate meaning, hidden authority and contradictory wire claims.

It is review evidence only. It amends no canonical artifact until the operator explicitly ratifies this converged package and the repository performs substitutive canonical filing.

This review does not choose:

- OpenAPI minor, generator, server library or SDK technology;
- blob/object store, CDN, image transformer, scanner or persistence topology;
- D6 screen/component topology;
- D7 transaction, queue, process, secret or deployment topology;
- a new Product operation, Permission, Principal class or media-business authority.

Implementation remains blocked until D9.

---

## 2. Evidence and method outcome

### 2.1 Repository authority

- D4-R1 admits source-qualified media and ListingIntent-scoped MPC-authored media while rejecting ProductAsset/media-master authority.
- The Operation Admission Matrix admits `CreateListingIntentMedia` only inside an exact mutable ListingIntent and rejects generic media-library CRUD.
- W1 selects `POST .../listing-intents/{id}:create-media`, typed multipart revision proof and one opaque owner validator authority.
- W2 selects one binary multipart part + typed `etag`, client idempotency, no arbitrary URL, no standalone Product media resource and authored-media descriptors inside the ListingIntent representation.
- W2 makes RFC 9457 `type` the primary Product problem identifier and rejects a duplicate global `code` taxonomy.
- W4 maps 95 Product operations / 29 Permissions; the media capability is `listing.manage`, and `GetListingIntent` is `offering.read`.
- Technical Ingress keeps provider callback/notification/OAuth vocabulary outside Product OpenAPI, SDK and Product Problem authority.
- D5-B1 requires every externally reachable route to be classified as Product API, provider/business-system protocol ingress or a separately justified technical surface.

### 2.2 Normative evidence

- RFC 9110 defines standard HTTP status and header semantics, including `405 + Allow`, `409`, `413` and `415`.
- RFC 9457 makes `type` the primary problem identifier and defines `about:blank` for a problem with no additional semantics beyond the HTTP status.

### 2.3 Consolidated disposition

```text
current D5-B2 semantic structure                 CONFIRMED
Problem/media final consistency                  RESTRUCTURE NOW — D5-B2 LOCAL
Fable material findings                          ADJUDICATED
Round 2                                           NOT REQUIRED
Product operations                               95 unchanged
ordinary Permissions                             29 unchanged
parent semantic reopen                           NONE
operator final ratification                      NEXT
implementation                                   BLOCKED UNTIL D9
```

The independent review found material defects in the lead candidate, especially an idempotency transcription conflict and an unclassified authority boundary for media-byte delivery. Both are resolved below without reopening semantic parents or adding a Product operation.

---

# 3. Consolidated corrections

## PM-C1 — Standard status-only failures use `about:blank` without waiving HTTP obligations

When a Product API failure has stable MPC-specific semantics, use the canonical W2 custom problem type.

When the failure has no additional Product semantics beyond the standard HTTP status, use RFC 9457 Problem Details with:

```text
type   = about:blank
status = applicable standard status
title  = recommended status phrase
```

Do not mint an MPC problem type merely to repeat standard HTTP meaning.

`about:blank` governs only the Problem Details body. It never waives status-specific HTTP obligations:

- `405 Method Not Allowed` includes the current `Allow` header value required by RFC 9110;
- `415 Unsupported Media Type` may include an honest `Accept` or `Accept-Encoding` value when the supported request representation/coding can be expressed without inventing part-level semantics;
- `413 Content Too Large` includes `Retry-After` only when the refusal is genuinely temporary; a standing Product size bound normally does not.

Clients branch on the HTTP status when `type = about:blank`.

---

## PM-C2 — Exact `CreateListingIntentMedia` failure grammar

### 3.2.1 Transport and representation failures

| Failure | Product HTTP disposition |
|---|---|
| Request cannot be parsed as a multipart representation | `400 malformed-request` |
| Top-level request representation is not the selected multipart media type | `415` + RFC 9457 `about:blank` |
| Binary file format is unsupported, undecodable or materially contradicts its declared content type after inspection | `415` + RFC 9457 `about:blank` |
| Request/file exceeds the enforced bound | `413` + RFC 9457 `about:blank` |

The selected operation contains exactly one binary file part. Therefore an unsupported file representation remains a `415` representation-format failure rather than being relabeled as semantic validation merely to obtain a part pointer. Exact accepted media families and size bounds must be explicit before implementation; their values and enforcement mechanism remain later realization.

### 3.2.2 Contract, revision, idempotency and lifecycle failures

| Failure | Product HTTP disposition |
|---|---|
| Required binary part or typed `etag` missing, duplicated or contract-invalid; semantic metadata invalid | `422 validation-error` with bounded field/part diagnostics |
| Supplied ListingIntent validator is stale | `409 resource-revision-conflict` |
| Current ListingIntent state does not admit media creation, after operation-specific exact-repeat handling | `409 resource-state-conflict` |
| Required `Idempotency-Key` absent/invalid | `400 idempotency-key-required` |
| Same key reused with materially different bytes, metadata, target or revision proof | canonical `422 idempotency-key-reused` |
| Equivalent prior intake still processing | canonical `409 idempotency-request-in-progress` |
| Unexpected internal storage/scanning/transformation/runtime failure before a successful Product result is established | `500 internal-error` |

`400 idempotency-key-required` is a new explicit status assignment in this consistency package, not a pre-existing W2 mapping. Its basis is that the required intake header is absent before semantic fingerprint/revision evaluation. The already-ratified `422 idempotency-key-reused` and exact type name `idempotency-request-in-progress` remain unchanged.

### 3.2.3 `resource-state-conflict`

Add one bounded Product Problem family:

```text
resource-state-conflict → HTTP 409
```

It means the server evaluated current authoritative resource state and that state does not admit the requested capability. It is distinct from:

```text
resource-revision-conflict → the supplied typed validator is stale
resource-state-conflict    → the supplied/current view may be current, but the capability is no longer admissible
```

Apply it only after any operation-specific exact-repeat or structural-idempotency rule has resolved a harmless repeat. This closes the already-admitted custom-capability population without introducing one state-specific type per owner.

### 3.2.4 Transport-guard ordering

A transport-level refusal for wrong representation or excess size:

- may occur before Product authentication;
- must be enforceable without full-body buffering merely to discover the violation;
- uses `about:blank` and discloses no Product resource existence, Organization, Membership, Permission or business state.

W2's idempotency processing order governs semantic evaluation after transport admission; it does not require an unbounded request body to pass authentication first.

### 3.2.5 Leakage fence

Do not add Product problem types such as:

```text
blob-upload-failed
virus-scanner-error
cdn-error
provider-image-error
```

Raw binary, storage key, access locator, scanner result, provider payload, secret, stack detail and arbitrary PII never enter Product Problem Details.

---

## PM-C3 — `type` is the sole global Product problem discriminator

- Product clients branch on custom problem `type`, or on HTTP `status` when `type = about:blank`.
- There is no global Product `code` duplicating `type`.
- A bounded problem-specific extension exists only for additional structured data genuinely required by that problem, such as validation locations or a stale-validator pointer.
- `title` and `detail` are human-facing and never stable programmatic identifiers.

No D5-B1 reopen is required; this is the accepted W2 crystallization.

---

## PM-C4 — Product and technical-protocol problem registries stay separate

```text
Product route
→ Product OpenAPI + W2 Product Problem contract

provider/OAuth technical route
→ protocol-local executable contract
→ never W2 Product type authority
→ never Product SDK operation/error surface
```

A technical route may reuse `application/problem+json` as a representation format. Format reuse does not merge type registries, operation identity, SDK exposure or business meaning.

Technical-protocol problem types must use a namespace disjoint from the Product problem-type namespace. Examples such as `oauth-state-expired`, `provider-code-invalid`, `seller-mismatch` and `provider-origin-invalid` remain forbidden from the Product catalog.

Product UI learns durable current posture through accepted Product reads, not callback navigation/error vocabulary.

---

## PM-C5 — Authored-media identity, revision and retention fence

For one ListingIntent:

```text
accepted authored binary + admitted metadata
→ one stable listing_intent_media_id
```

- the ID belongs to exactly one ListingIntent and cannot be referenced from another;
- it never rebinds to different bytes or material meaning;
- a materially different upload creates another ID through the existing capability;
- current selection/order/role is ListingIntent desired-state meaning and changes only through the accepted draft update;
- unselecting an authored ID does not by itself delete its bytes or historical reference;
- publication-attempt history retains only the identity/provenance needed to explain what was attempted;
- no Product `UpdateMedia`, `DeleteMedia`, media collection or ProductAsset CRUD is admitted.

### 3.5.1 ListingIntent validator advancement

Successful media creation changes the ListingIntent representation because authored-media descriptors are one of its read axes. Therefore it advances the ListingIntent's strong opaque owner validator.

The successful custom-method result is operation-specific and returns:

```text
media                         → stable identity/provenance descriptor
listing_intent_etag           → current parent ListingIntent validator after acceptance
```

Do not place the parent ListingIntent validator in an HTTP `ETag` header on the distinct `:create-media` request URI. It remains typed result data carrying the one parent-resource revision authority.

Consequences:

- two concurrent creates against the same ListingIntent validator serialize;
- the first accepted intake advances the validator;
- a materially different concurrent loser receives `409 resource-revision-conflict` and re-reads/rebases;
- a caller can chain sequential successful uploads using the returned `listing_intent_etag`, without a mandatory GET after every success;
- an exact lost-response retry resolves the already accepted intake before stale-revision re-evaluation and returns the same intake result rather than creating another media ID.

No multi-file/bulk intake is admitted. Reopen only if D6 proves real multi-image authoring is unusable under this serialized model.

### 3.5.2 Retention/erasure residual

Product 1.0 admits no client-facing authored-media erasure operation. Unselection is not erasure.

No accepted authority currently establishes a universal legal/contractual retention duration or an Organization-lifetime immortality rule. Retention/erasure therefore remains **Unknown / deferred to D2 data ownership and D7 realization**, subject to these D5 constraints:

- content required by current selection or material historical explanation cannot be silently removed;
- a future retention/erasure rule must preserve enough historical identity/provenance to explain prior consequential attempts;
- a legal, privacy, contractual, operator or material cost obligation reopens the smallest D2/D7 scope rather than automatically creating Product delete CRUD.

---

## PM-C6 — Descriptor identity and authored-media byte delivery authority

Identity/provenance and presentation access are separate schemas:

```text
ListingIntentMediaDescriptor
  → stable authored-media identity + bounded content/provenance facts
  → eligible for current selection and historical attempt basis

ListingIntentMediaPresentationDescriptor
  → ListingIntentMediaDescriptor + volatile access reference
  → response-only presentation aid
  → never persisted into history, idempotency fingerprint, logs or Problem Details
```

### 3.6.1 D5 authority property

Authored-media byte delivery is a **separately justified technical presentation surface** under D5-B1 boundary classification. It is:

- not a 96th Product operation;
- not a Technical Ingress A/B lane;
- not a Product SDK/OpenAPI business operation;
- not a generic Media/Asset service;
- not authorized by stable media ID alone.

The baseline authority law is:

> **A caller that cannot currently obtain the corresponding authored-media presentation descriptor through `GetListingIntent` under the exact path Organization and `offering.read` cannot obtain the bytes.**

The technical delivery realization therefore reuses current Product authentication, unique Principal binding, Principal eligibility, Organization Membership and exact `offering.read` authorization for the referenced ListingIntent/media relationship. Exact route, proxy/storage/CDN topology, streaming and transformation mechanics remain D7.

The `CreateListingIntentMedia` response under `listing.manage` returns the stable descriptor and new parent validator, not a durable presentation locator. The client already possesses the submitted bytes; later server-side retrieval is a read concern through `GetListingIntent` and `offering.read`.

A durable, anonymous or freely forwardable object-store/CDN locator is not baseline. A delegated bearer capability may be reconsidered only by an explicit smallest-scope D5/D7 reopen proving that authenticated delivery is materially unsuitable and preserving tenant/scope/expiry/non-enumerability and credential-handling constraints.

Technical media-delivery failures remain technical-surface failures; they do not expand the W2 Product Problem catalog.

If D6 proves that an embedded presentation reference plus the bounded technical surface cannot satisfy the real consumer, reopen only the smallest B2 operation/W4 surface before implementation. D7 may not invent a Product media GET privately.

---

## PM-C7 — Source and authored media remain distinct descriptor families

- source media remains source-qualified external evidence owned by the Readiness/D4 seam;
- authored media remains ListingIntent-scoped MPC state;
- ListingIntent selection uses a bounded discriminated union because it must choose/order both origins;
- dimensions/content-type primitives may be reused only where their meaning is genuinely identical;
- source-media locators and authored-media access references are distinct types even if their JSON shape happens to match, because issuer, trust, lifetime and governing authority differ;
- source candidate key and authored media ID never substitute for one another;
- arbitrary client URL remains rejected;
- provider image-error feeds remain deferred until a named consumer proves need.

No generic Media/Asset owner or cross-ListingIntent library is introduced.

---

## PM-C8 — Active canonical artifacts converge on the router

At final canonical filing:

- preserve accepted semantic sections;
- remove/replace stale terminal sequencing/status claims in active B2 artifacts instead of appending another superseding status layer;
- make every artifact defer current status and exact next action to the router;
- keep Git history as the archive;
- change no Product operation, Permission, owner or accepted semantic decision through this cleanup.

This is authority convergence, not editorial redesign.

---

# 4. Proof obligations / negative fixtures

Later OpenAPI/D7/D8/implementation proof must make at least these defects red:

1. a `405` response omits or lies in `Allow`;
2. wrong top-level multipart representation is classified as `422` rather than `415`;
3. unsupported/undecodable or inspected-mismatching binary format is accepted or classified as semantic metadata validation;
4. excess-size body is fully buffered before the transport guard can refuse it;
5. malformed multipart becomes a provider/storage failure;
6. missing/duplicate file part or invalid typed ETag is accepted;
7. absent `Idempotency-Key` creates intake or receives an undeclared status;
8. same key + changed bytes/metadata/ETag is not `422 idempotency-key-reused`;
9. in-progress exact intake uses a non-canonical problem type;
10. a current but lifecycle-inadmissible capability is confused with stale revision or silently succeeds;
11. exact lost-response replay creates a second media ID or fails only because the first intake advanced revision;
12. successful media creation fails to advance/return the current parent ListingIntent validator;
13. two concurrent creates with one validator both succeed as independent current-state writes;
14. `listing_intent_media_id` from ListingIntent A is accepted in B;
15. accepted bytes are replaced under an existing ID;
16. unselection destroys material historical evidence;
17. a presentation access reference reaches durable history, idempotency fingerprint, logs or Problem Details;
18. stable media ID alone retrieves bytes;
19. caller without current Organization Membership or `offering.read` retrieves authored bytes;
20. anonymous/durable/freely forwardable object-store locator becomes baseline access;
21. source locator and authored access reference become one authority/type;
22. raw storage/scanner/provider error appears in Product Problem Details;
23. technical OAuth/provider types use the Product problem namespace or appear in Product OpenAPI/SDK;
24. Product clients receive duplicate global `type` and `code` discriminators;
25. D7 privately creates a media Product GET, Delete or Update operation;
26. stale active “next action” text misroutes a fresh session away from the router.

A green artifact that did not execute the protected subject is no proof.

---

# 5. Cross-artifact disposition after ratification

| Artifact | Substitutive filing disposition |
|---|---|
| D0→D4 / D4-R1 | confirm; no reopen |
| D5-B1 | confirm; bounded technical-presentation classification uses its existing route taxonomy |
| Operation Admission Matrix | confirm; 95 operations unchanged |
| W1 | remove stale sequencing and add bounded cross-reference only where needed |
| W2 | incorporate PM-C1…PM-C7, exact statuses, `resource-state-conflict`, media result/descriptor and negative controls |
| W3 | confirm cursor/problem interaction unchanged |
| W4 | confirm 29 Permissions; add bounded delivery-authority cross-reference without a new Permission/operation |
| Technical Ingress | confirm A/B unchanged; media delivery is explicitly not ingress |
| Product OpenAPI/SDK | include only admitted Product capability; exclude technical delivery route and ingress |
| Router/cockpit | advance only after canonical filing; cockpit remains non-authoritative |

After filing, remove this candidate from the active tree and reset `AI-DIALOG.md` to protocol-only. Git history remains the review archive.

---

# 6. Final converged outcome

```text
current semantic structure                     CONFIRMED
Problem/media consistency                      RESTRUCTURE NOW — D5-B2 LOCAL
custom Product problems                        small; type URI is primary
standard status-only problems                  about:blank + required HTTP headers
idempotency mappings                           canonical W2 preserved; missing-key 400 declared
current lifecycle conflict                     resource-state-conflict / 409
authored media identity                        immutable, ListingIntent-scoped
authored media creation                        advances and returns parent validator
selection removal                              not deletion
authored-media byte access                     current Organization + offering.read technical surface
standalone Product media GET                   NOT ADMITTED
Product operations                             95 unchanged
ordinary Permissions                           29 unchanged
parent reopen                                  NONE
Round 2                                        NOT REQUIRED
operator ratification                          NEXT
implementation                                 BLOCKED UNTIL D9
```

The package is converged but remains non-authoritative pending explicit operator ratification.
