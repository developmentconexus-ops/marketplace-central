# D5-B2 — Final Problem / Media Consistency Lead Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD REVIEW EVIDENCE — NOT CANONICAL  
> **Parent stage:** D5-B2 Product Operation / Resource Surface  
> **Canonical inputs:** accepted D0→D4 + D4-R1 + D5-B1 + Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Reviewed:** 2026-08-19  
> **Lead outcome:** `RESTRUCTURE NOW — D5-B2-LOCAL PROBLEM/MEDIA CONSISTENCY`  
> **Independent review:** materially warranted before operator ratification

## 1. Purpose and authority fence

This candidate performs the router-mandated final D5-B2 consistency pass over:

1. W2 Product Problem Details versus provider/OAuth protocol-local failures;
2. ListingIntent-scoped authored-media intake, representations and D7-deferred storage/delivery;
3. W1/W2/W3/W4/Technical Ingress terminology and routing;
4. Product OpenAPI/SDK exclusion of technical ingress;
5. duplicate meaning, hidden authority and contradictory wire claims.

It does not amend any canonical artifact. Reviewer output remains evidence until GPT adjudication, explicit operator ratification and canonical filing.

This review does not choose:

- OpenAPI minor, generator, server library or SDK technology;
- blob/object store, CDN, signed-URL, proxy, scanner, image transformer or retention implementation;
- D6 screen/component topology;
- D7 transaction, storage, queue, process, secret or deployment topology;
- new Product operations, Permissions or media-business authority.

Implementation remains blocked until D9.

---

## 2. Evidence basis

### 2.1 Repository authority

- `D4-R1-PUBLICATION-INPUT.md` admits source-qualified media and ListingIntent-scoped MPC-authored media while rejecting ProductAsset/media-master authority.
- `D5-B2-OPERATION-ADMISSION-MATRIX.md` admits `CreateListingIntentMedia` only inside an exact mutable ListingIntent and rejects generic media-library CRUD.
- W1 selects `POST .../listing-intents/{id}:create-media` and typed multipart revision proof.
- W2 selects one binary multipart part + typed `etag`, `200` with a ListingIntent-scoped descriptor, client idempotency, no arbitrary URL and no standalone Product media resource by default.
- W2 makes RFC 9457 `type` the primary Product problem identifier and rejects a duplicate global `code` taxonomy by default.
- Technical Ingress keeps provider callback/notification/OAuth failure vocabulary outside the Product Problem catalog, Product OpenAPI and Product SDK.
- W4 maps 95 Product operations / 29 Permissions; technical OAuth begin is not a 96th operation and `CreateListingIntentMedia` remains the only admitted media mutation.

### 2.2 Normative external evidence

- RFC 9110 defines HTTP method and status semantics, including ordinary transport/representation failures such as 405, 413 and 415.
- RFC 9457 defines Problem Details, makes `type` the primary problem identifier and defines `about:blank` for failures that carry no application-specific semantics beyond the HTTP status.

Normative references:

- https://www.rfc-editor.org/rfc/rfc9110.html
- https://www.rfc-editor.org/rfc/rfc9457.html

No current implementation or legacy OpenAPI shape is used as target authority.

---

## 3. Method result

The canonical package is structurally sound:

- Product problems and valid business outcomes remain separated;
- technical provider/OAuth failures remain protocol-local;
- authored media remains ListingIntent-scoped;
- source media and authored media do not collapse;
- idempotency and revision proof remain independent;
- no generic Media/Asset, Integration or provider error authority exists.

However, final consistency is not yet complete. Three classes remain reachable:

1. **standard HTTP failures lack an explicit disposition** between W2's small custom problem catalog and D5-B1's broader API-failure classes;
2. **media identity, selection lifecycle and binary access are not separated explicitly enough**, allowing either hidden ProductAsset CRUD or an unusable descriptor-only realization;
3. **active canonical artifacts retain superseded sequencing text**, which is non-semantic but can misroute a fresh agent despite the router being sole status authority.

The corrections below are D5-B2-local. No semantic parent reopen is proven.

---

# 4. Converged lead corrections proposed for independent challenge

## PM-C1 — Standard HTTP failure grammar uses `about:blank`; custom Product types remain small

### Finding

W2 correctly defines a small stable MPC Product problem catalog, but it does not explicitly classify ordinary HTTP failures whose entire semantics are already expressed by their status code. D5-B1 includes unsupported HTTP/operation-contract failures, and the admitted multipart operation necessarily exposes content-size and content-type failure classes.

Adding one custom MPC problem type for every standard HTTP status would create duplicate protocol vocabulary. Reusing `malformed-request` or `validation-error` for a wrong request media type or transport-level size rejection would misstate the failure.

### Corrected invariant

```text
When a Product API failure has stable MPC-specific semantics,
use the canonical W2 custom problem type.

When the failure has no additional Product semantics beyond
the standard HTTP status, return RFC 9457 Problem Details with:

  type = about:blank
  status = the applicable standard status
  title = the standard status phrase

and do not create another MPC code/type by symmetry.
```

Current examples include proportionately:

- `405 Method Not Allowed`;
- `413 Content Too Large`;
- `415 Unsupported Media Type`.

A later standard failure receives a custom Product type only if a real programmatic consumer needs stable MPC semantics beyond the status.

### Global Maximum

This uses the standard's own extension boundary, preserves W2's small catalog and avoids a second `code` taxonomy. It is smaller and more interoperable than creating `method-not-allowed`, `request-too-large` and `unsupported-media-type` MPC type families.

---

## PM-C2 — `CreateListingIntentMedia` gets one exact failure map

### Finding

The positive media contract is coherent, but without an explicit negative map, different realizations can classify the same failure as 400, 413, 415, 422, 409 or 500 and leak object-storage/scanner/provider vocabulary into the Product API.

### Corrected invariant

| Failure | Product HTTP disposition |
|---|---|
| Request cannot be parsed as a valid multipart representation | `400 malformed-request` |
| Request representation is not the selected multipart media type | `415` + RFC 9457 `about:blank` |
| Request exceeds the enforced whole-request/file size bound | `413` + RFC 9457 `about:blank` |
| Required binary part or typed `etag` missing/duplicated/invalid; semantic metadata invalid; declared and inspected content materially contradict | `422 validation-error` with bounded part/field diagnostics |
| ListingIntent revision proof is stale | `409 resource-revision-conflict` |
| Required idempotency key absent | existing `400 idempotency-key-required` |
| Same key reused with different binary identity, metadata or revision proof | existing `409 idempotency-key-reused` |
| Equivalent first intake still in progress | existing `409 idempotency-in-progress` |
| Unexpected internal storage/scanning/transformation/runtime failure before a successful Product result can be established | `500 internal-error` with no backend/provider details |

Rules:

- deterministic accepted-media constraints are contract validation, not provider/storage errors;
- exact accepted content families and size bounds must be explicit before implementation, but their numerical/provider realization is not selected here;
- no `blob-upload-failed`, `virus-scanner-error`, `cdn-error`, `provider-image-error` or similar Product problem taxonomy is admitted;
- raw binary, storage key, signed locator, scanner result, provider payload, secret and PII never enter Problem Details.

---

## PM-C3 — `type` remains the sole global Product problem discriminator

### Finding

D5-B1 permits a stable machine-readable MPC extension where a real consumer needs one; W2 later correctly selected the RFC 9457 `type` URI as the primary identifier and rejected a duplicate top-level `code` taxonomy by default. The two statements are compatible but should be closed explicitly so OpenAPI/tooling cannot introduce both by habit.

### Corrected invariant

- Product clients branch on custom problem `type`, or on HTTP `status` when `type = about:blank`;
- there is no global Product `code` duplicate of `type`;
- a bounded problem-specific extension is admitted only for additional structured data that the problem contract actually requires, such as validation issues or idempotency timing;
- human-readable `title`/`detail` are never stable programmatic identifiers.

No D5-B1 reopen is required; this is the already-accepted W2 crystallization.

---

## PM-C4 — Product Problem format and technical protocol vocabulary stay separate

### Finding

Technical Ingress already excludes provider/OAuth failures from W2, but “protocol-local” could be misread as either forbidding standards reuse or authorizing reuse of Product problem types.

### Corrected invariant

```text
Product route
→ Product OpenAPI + W2 Product Problem contract

technical provider/OAuth route
→ provider/protocol-local executable contract
→ never W2 Product type authority
→ never Product SDK operation/error surface
```

A technical route may reuse a standard representation format such as Problem Details when useful, but format reuse does not merge type registries, Product operation identity, SDK exposure or business outcome semantics.

Examples such as `oauth-state-expired`, `provider-code-invalid`, `seller-mismatch` and `provider-origin-invalid` remain forbidden from the Product catalog.

Product UI learns durable current posture only through accepted Product reads such as `GetMarketplaceInstallation`, not by treating a callback redirect/error as Product state.

---

## PM-C5 — ListingIntent authored-media identity is immutable and selection is not deletion

### Finding

W2 establishes `listing_intent_media_id` but does not state the complete lifecycle fence. Without it, an implementation could replace bytes under one ID, delete media merely because it is no longer selected, or invent generic media CRUD/retention authority.

### Corrected invariant

For one ListingIntent:

```text
accepted authored binary + admitted metadata
→ one stable listing_intent_media_id

that ID never rebinds to different bytes/meaning
```

- the ID is scoped to exactly one ListingIntent and cannot be referenced by another;
- the media descriptor is immutable as to the accepted binary identity; a materially different upload creates another ID through the existing create capability;
- current media selection/order/role is ListingIntent desired-state meaning and may change through the existing draft update;
- removing an authored media ID from current selection **does not itself delete** the accepted binary or historical reference;
- submission/publication-attempt snapshots preserve the media references/provenance required for explanation;
- no Product `DeleteMedia`, `UpdateMedia`, media collection or ProductAsset CRUD is admitted;
- D7 may garbage-collect only content that is not current, not historically required and whose retention/privacy obligations permit removal.

This preserves history without turning MPC into a general media archive.

---

## PM-C6 — Descriptor identity and binary delivery are separate; current Product surface gains no hidden GET

### Finding

W2 correctly rejects a standalone media Product resource/GET and defers blob/CDN/storage mechanics. But an overly literal realization can return only an opaque ID that an authorized author cannot inspect across reads, making durable selection/order unusable. The opposite local maximum is to invent a generic media GET/resource or expose stable storage/provider URLs as identity.

### Corrected invariant

`ListingIntentMediaDescriptor` carries Product meaning sufficient to identify the authored media within its ListingIntent and support the accepted authoring consumer. Binary presentation/access is a separate technical capability:

- a server-issued access reference may be exposed in an authorized media descriptor when needed for inspection/preview;
- the access reference is not media identity, not client-authored state and not historical provenance;
- it may change or expire between reads without changing `listing_intent_media_id`;
- it must not expose object-store key/topology or a provider-authored URL as MPC authority;
- bearer-capability locators, if selected, are bounded and must not enter normal logs, Problem Details or durable ListingIntent history;
- exact signed URL, authenticated proxy, CDN and transformation mechanics remain D7;
- exact descriptor field spelling belongs to the single Product OpenAPI closure;
- no standalone Product media GET/download operation is currently admitted.

If D6 later proves that an additional Product retrieval operation is required and an embedded access reference cannot satisfy the consumer safely, reopen only the smallest B2 operation/W4 surface before implementation. D7 may not invent the operation privately.

---

## PM-C7 — Source media and authored media remain two bounded descriptor families

### Finding

The selection union correctly distinguishes `source_media_candidate_key` from `listing_intent_media_id`. Final consistency must prevent shared presentation mechanics from creating a generic `Media`, `Asset`, `URL` or payload-bag authority.

### Corrected invariant

- source media descriptor/candidate remains source-qualified external evidence owned by the accepted Readiness/D4 seam;
- authored media descriptor remains ListingIntent-scoped MPC state;
- selection may use a bounded discriminated union because ListingIntent must choose/order both origins;
- common technical fields such as dimensions/content type/access reference may share schema primitives only where meanings are truly identical;
- source candidate key and authored media ID never substitute for one another;
- arbitrary client URL remains rejected;
- provider image-error feeds remain deferred unless a named media/readiness consumer proves need.

No generic media owner or cross-ListingIntent media library is introduced.

---

## PM-C8 — Active canonical artifacts must stop carrying misleading “next work” text

### Finding

The router is explicitly the sole status/next-action authority, and later appendices currently supersede stale sequencing in W1, the Operation Matrix and W4. The semantics are not contradictory, but active files still contain historical “W3 next”, “Technical Ingress next” and “Wire Contract next” sentences that can misroute an agent reading sections out of order.

### Corrected direction at final filing

- preserve canonical semantic sections;
- remove or replace stale terminal sequencing/status claims in the active B2 artifacts rather than append another superseding status appendix;
- make every artifact defer current status/next action to the router;
- keep Git history as the archive;
- change no Product operation, Permission, owner, path, schema or accepted decision through this cleanup.

This is authority convergence, not documentation aesthetics.

---

# 5. Alternatives rejected

## 5.1 One custom problem type per status

Rejected. It duplicates HTTP semantics, expands client taxonomy and offers no additional Product meaning for ordinary 405/413/415 failures.

## 5.2 One global `code` plus RFC `type`

Rejected. Two global discriminators can drift and create competing client contracts. `type` is already the selected primary identifier.

## 5.3 Provider/storage/scanner error passthrough

Rejected. It leaks implementation, creates unstable external vocabulary, risks secrets/PII and makes a technical dependency a Product error authority.

## 5.4 Generic Product media/asset resource

Rejected. There is no cross-ListingIntent media-library consumer or independent lifecycle. It recreates the ProductAsset/media-master structure already rejected by D4-R1 and the Operation Matrix.

## 5.5 No durable inspectability seam

Rejected as a local maximum. An opaque ID without an authorized presentation path makes durable media selection/order across reads operationally unusable. The seam is required; generic media CRUD is not.

## 5.6 Admit a standalone media GET immediately

Rejected. No separate Product operation/Permission consumer has been admitted, and a bounded descriptor access reference can prepare the seam without inventing the operation.

---

# 6. Cross-artifact disposition

| Artifact | Lead disposition |
|---|---|
| D0→D4 / D4-R1 | **CONFIRM — no reopen** |
| D5-B1 | **CONFIRM — W2 crystallization only** |
| Operation Admission Matrix | **CONFIRM — 95 operations unchanged** |
| W1 | **REVISE locally — status/routing cleanup only** |
| W2 | **REVISE locally — PM-C1…C7 and negative controls** |
| W3 | **CONFIRM — existing cursor problems remain coherent** |
| W4 | **CONFIRM — 29 Permissions unchanged; stale sequencing cleanup only** |
| Technical Ingress | **CONFIRM — add bounded cross-reference only if needed** |
| Product OpenAPI/SDK | **technical ingress remains excluded** |
| D6/D7 | **remain deferred/blocked** |

No Product operation, Permission, Principal class, business owner or technical ingress family is added.

---

# 7. Proof obligations / negative fixtures

Later OpenAPI/D7/D8/implementation proof must make at least these defects red:

1. wrong top-level content type reaches media handler as `422` instead of standard `415`;
2. oversized media is buffered/accepted or represented as `validation-error` instead of bounded `413`;
3. malformed multipart becomes provider/storage failure;
4. missing/duplicate media part or invalid typed ETag is accepted;
5. same idempotency key + changed bytes/metadata/ETag resolves as the prior intake;
6. exact lost-response retry with the original fingerprint creates a second media ID;
7. stale ListingIntent revision creates/associates media;
8. raw object-store/scanner/provider error appears in Product Problem Details;
9. Product clients receive two global problem identifiers (`type` and duplicate `code`);
10. protocol-local OAuth/provider errors appear in Product OpenAPI/SDK;
11. `listing_intent_media_id` from ListingIntent A is accepted in ListingIntent B;
12. accepted media bytes are replaced under an existing ID;
13. unselecting media silently destroys evidence required by a past attempt;
14. access-reference expiry changes media identity or current selection;
15. storage key/provider URL becomes Product media identity;
16. descriptor access capability leaks to normal logs/history/problems;
17. source candidate and authored media ID are interchangeable;
18. arbitrary client URL becomes trusted authored media;
19. D7 privately creates a new media GET/Product operation;
20. stale “next action” text in an active canonical artifact overrides/misroutes the router.

Exact executable mechanisms belong to the responsible later stage.

---

# 8. Independent-review decision

One coherent independent review is warranted because the correction package touches:

- an untrusted binary-upload trust boundary;
- HTTP error interoperability;
- idempotency/revision behavior over binary identity;
- authorization-scoped media access capability;
- retention/history versus deletion;
- Product API versus technical-protocol error authority.

The review must challenge the whole package, not perform eight ceremonial micro-reviews.

Round 2 exists only if a material contradiction survives GPT adjudication.

---

# 9. Lead outcome

```text
current D5-B2 semantic structure                 CONFIRMED
Problem/media final consistency                  RESTRUCTURE NOW — D5-B2 LOCAL
custom Product problem catalog                   KEEP SMALL
standard status-only problems                    RFC 9457 about:blank
Product problem primary identifier               type only
technical provider/OAuth problems                protocol-local / Product-excluded
authored media authority                         ListingIntent-scoped
media identity                                   immutable per accepted upload
selection removal                                not deletion
binary delivery                                  technical access seam, not identity
standalone Product media GET                     NOT ADMITTED
Product operations                               95 unchanged
ordinary Permissions                             29 unchanged
parent reopen                                    NONE
independent review                               NEXT
implementation                                   BLOCKED UNTIL D9
```

This candidate is review evidence only. Canonical artifacts remain unchanged pending independent challenge, GPT adjudication and explicit operator ratification.
