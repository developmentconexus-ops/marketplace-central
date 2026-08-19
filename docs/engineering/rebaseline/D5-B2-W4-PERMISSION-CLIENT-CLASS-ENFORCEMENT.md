# D5-B2 — W4 Permission → Operation / Client-Class Enforcement

> **Status:** ACCEPTED / CANONICAL — Whole-W4 operator-ratified  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Client/Auth authority:** `D5-B2-PRODUCT-OPERATION-SURFACE.md` B2-A  
> **Identity/access authority:** canonical D2  
> **Schema/collection authorities:** canonical W2 + W3  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical W1/W2/W3  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Whole-W4 final ratification incorporated:** 2026-08-19

## 1. Purpose and authority

W4 is the **single canonical Product API ordinary-access Permission + allowed Principal/client-class enforcement authority** for every admitted Product 1.0 operation.

It consolidates the accepted W4 design plus the operator-ratified Whole-W4 lead/Fable/GPT adversarial review. Git history preserves review snapshots; they are not parallel active authorities.

W4 does not choose Keycloak realm/client/deployment topology, token-claim mapping, database ACL/RLS, middleware package shape, cache/storage or other D7 realization.

> **Authentication, Principal binding, current Principal access eligibility, Organization Membership, exact ordinary Permission, operation-specific Principal-class admission, physical epistemic qualification, owner business disposition and Governance authorization are separate gates. Possessing one never implies the others.**

---

## 2. Principal/client classes

D2 remains authoritative for exactly three Principal kinds:

```text
H = human
A = automation
S = system
```

Do not introduce Product Principal kinds such as `physical_system`, `agent`, `service_account`, `robot` or `integration` by convenience.

### 2.1 Meaning of the classes

- **H — human:** accountable interactive human Principal resolved from the accepted external OIDC path.
- **A — automation:** non-human Principal executing a business automation/policy-owned path. AI/agent/repricer/business automation belongs here rather than pretending to be `system`.
- **S — system:** bounded machine/system actor for system-established actions and reads. A system Principal is not business automation by default.

OAuth grant type is not Principal kind. Client Credentials authenticates a confidential machine client; MPC still resolves one MPC-owned non-human Principal and its kind.

### 2.2 Exact refinement of B2 `both`

The Operation Matrix used `both` as a deliberately coarse pre-wire admission label. W4 makes it exact:

- ordinary **Q/read** operations admitted to `both` → `H | A | S`;
- side-effect-free stateless `EvaluatePriceScenario` → `H | A | S` despite its C classification;
- consequential/business-authoring/coordination **C** operations admitted to `both` → `H | A` unless the matrix explicitly names a narrower/different class;
- `S` is never a generic substitute for `A` on business authoring.

This is a W4 crystallization of already-distinct D2 Principal meanings, not a new actor taxonomy.

---

## 3. Ordinary Permission model

### 3.1 Permission is flat and exact

Permissions are stable Product capabilities checked **exactly at the operation boundary**.

There is no implicit hierarchy:

```text
*.manage  != *.read
*.execute != *.read
```

and no wildcard/prefix implication:

```text
sales.manage   != sales.*
economics.read != economics.*
```

An `AccessRole` may bundle several exact Permissions. Role bundling never changes the meaning of any Permission.

Permission names do not reserve authorization for future operations. Every new Product operation requires explicit admission + mapping.

### 3.2 Permission vocabulary

`authenticated` for `/access-context` is a special operation condition, **not** a stored Permission.

The canonical ordinary Product Permission vocabulary contains **29 Permissions**:

```text
access.read
access.manage
portfolio.read
portfolio.manage
readiness.read
readiness.manage
offering.read
listing.manage
price.manage
availability.read
availability.manage
market.read
economics.read
economics.policy.manage
economics.reconcile
governance.read
governance.decide
governance.manage
sales.read
sales.manage
materialization.read
materialization.resolve
fulfillment.read
fulfillment.execute
fulfillment.manage
post_sale.read
post_sale.manage
work.read
work.manage
```

No generic ACL/ReBAC/policy expression language, Permission DSL or IdP/provider scope ontology is admitted.

### 3.3 Mutation/capability response disclosure is operation-scoped

A mutation/capability Permission authorizes the W1/W2 response representation for that **same operation and exact operation subject**.

That response disclosure:

- does not grant the corresponding `Get`/`List`/`Search` operation;
- does not authorize another subject/resource;
- creates no Permission inheritance;
- does not require an additional read Permission merely to receive the operation's normal W2 representation/validator;
- does not change 200/204 response shape according to the caller's separate read Permissions.

Therefore `portfolio.manage` does not grant `portfolio.read` **operations**, even though an admitted portfolio mutation may return its own updated MarketplaceInstallation representation.

---

## 4. Canonical enforcement sequence

For an Organization-scoped Product operation, enforcement is proportionately:

```text
1. accept/authenticate the external credential under B2-A
2. resolve MPC Principal binding
   - zero accepted Principal bindings → 401 authentication-required
   - more than one binding → 500 internal-error; never choose one
3. require current D2 Principal access eligibility
4. identify the Product operation
5. resolve current path-Organization Membership
6. enforce allowed Principal kind
7. enforce the exact ordinary Permission
8. when required, resolve current server-owned operation-specific physical qualification
9. resolve path/body/query resources and secondary refs fail-closed in Organization scope
10. apply W1/W2 contract, validation, idempotency and revision/precondition grammar
11. evaluate current owner business validity/disposition
12. evaluate Governance authorization when required
13. establish durable owner intake/effect
```

The exact middleware/code decomposition remains D7/implementation. The semantic gates and their non-equivalence are binding.

### 4.1 Current-authority property

Every mutable/revocable access fact is evaluated against its **current authoritative state** at the Product boundary.

- external credential validity and accepted binding evidence follow the accepted authentication authority;
- MPC-owned Principal access eligibility, Membership, RoleAssignment/Permission, Principal kind and operation-specific physical qualification follow current MPC authority;
- stale token claims, cached snapshots or retired provisioning records never remain access authority after revocation;
- D7 may cache only while preserving this revocation property.

This is a contract property. Storage/cache/notification mechanics remain D7.

### 4.2 `/access-context` exception

`GET /access-context` is the one platform-scoped self-only Product Q.

It requires:

```text
valid accepted authentication
+ exactly one MPC Principal binding
+ current Principal access eligibility
```

It waives only Organization Membership and Product Permission because its purpose is to discover the current Principal's visible Organizations and effective ordinary Permissions.

An eligible Principal with zero memberships receives a successful representation containing zero visible Organizations. An ineligible Principal receives `403 access-denied` and no Membership/Role/Permission disclosure.

`authenticated` never becomes an AccessRole/Permission definition.

---

## 5. Authentication/Permission authority fences

### 5.1 IdP role/scope is not Product Permission

OIDC/OAuth/Keycloak owns authentication protocol identity and token issuance. MPC owns Principal, current Principal access eligibility, Membership, RoleAssignment, AccessRole and Permission.

A realm/client role or OAuth scope named `price.manage`, `admin`, `fulfillment` or similar never independently grants a Product operation.

```text
valid token
!= current Principal access eligibility
!= Membership
!= Permission
!= allowed Principal kind
!= physical qualification
!= business disposition
!= Governance authorization
!= executable now
```

### 5.2 Current revocation remains authoritative

Current Principal eligibility, Membership, RoleAssignment/Permission and operation-specific physical qualification remain authoritative even when a previously issued token or cached snapshot still exists.

Revoking one of those current gates changes future access without rewriting historical actor attribution or past accepted occurrences.

---

## 6. Client-class admission is independent from business/Governance authority

Client-class admission is a Product API boundary.

A caller that has the ordinary Permission but whose Principal kind is not admitted for that operation receives ordinary access denial; the request does not enter owner business semantics merely to be rejected there.

Governance cannot widen the Product client-class surface:

```text
Governance delegation/decision
!= permission to invoke a human-only operation as automation/system
```

Conversely, ordinary Permission never implies Governance authorization or business disposition.

Valid post-access business meanings such as `approval-required`, owner rejection/prohibition, unavailable evidence or external-required remain Product semantics under W2 and are not converted into access denial merely because the outcome is negative.

---

## 7. Physical / epistemic authority

Holding `fulfillment.execute` does not prove epistemic ability to establish a physical fact.

A physically capable machine is represented as:

```text
Principal.kind = system
+ current Fulfillment-recognized operation-specific physical-establisher qualification
```

This qualification:

- is **not** a fourth Principal kind;
- is **not** an AccessRole/Permission;
- is **not** Governance authorization;
- is **not** inferred from Client Credentials, IdP/OAuth/provider roles/scopes or arbitrary token claims;
- is bounded to the exact physical checkpoint capability whose evidence source has been proven;
- is server-resolved only after allowed kind + exact `fulfillment.execute` admission;
- must still be current when the call is admitted.

Request body, client-declared station/device/evidence source, frontend state or other caller-controlled data cannot self-assert trusted/qualified physical authority. W2 server attribution of effective Principal/source/time remains binding.

W4 freezes only the required predicate. D7 owns provisioning/binding/runtime mechanics. Do not create a generic machine-capability graph or `system_capabilities[]` Product model.

Baseline physical operation admission:

```text
RecordSeparation           → H only
RecordPacking              → H only
RecordPhysicalConference   → H OR currently qualified S
RecordDispatchHandoff      → H OR currently qualified S
```

An ordinary `A` with `fulfillment.execute` cannot establish any of these physical checkpoint facts. An unqualified/retired `S` cannot establish PhysicalConference or DispatchHandoff.

If a real machine later must establish Separation/Packing, reopen only those operation client-class admissions with evidence.

---

# 8. Exact 95-operation enforcement matrix

The ratified Product 1.0 operation inventory contains **95 admitted Product operations**. W4 maps every one; deferred/rejected candidates remain outside the runtime contract.

## 8.1 Identity / ordinary access — 5

| Operation | Class | Required ordinary access | Allowed Principal |
|---|---|---|---|
| `GetCurrentAccessContext` | Q | authenticated special condition | H / A / S |
| `ListOrganizationMembers` | Q | `access.read` | H |
| `ListAccessRoles` | Q | `access.read` | H |
| `AssignAccessRole` | C | `access.manage` | H |
| `RevokeAccessRole` | C | `access.manage` | H |

Role revocation remains fail-safe/monotonic for the targeted standing grant after ordinary caller admission. That safety rule does not bypass current Principal eligibility/Membership/Permission/client-class checks.

## 8.2 Marketplace Portfolio — 6

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListMarketplaceInstallations` | Q | `portfolio.read` | H / A / S |
| `GetMarketplaceInstallation` | Q | `portfolio.read` | H / A / S |
| `CreateMarketplaceInstallation` | C | `portfolio.manage` | H |
| `UpdateMarketplaceInstallationConfiguration` | C | `portfolio.manage` | H |
| `DeactivateMarketplaceInstallation` | C | `portfolio.manage` | H |
| `ListSellingEntities` | Q | `portfolio.read` | H / A / S |

## 8.3 Product & Channel Readiness — 5

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `SearchSourceProductsForMarketplace` | Q | `readiness.read` | H / A / S |
| `GetProductChannelReadiness` | Q | `readiness.read` | H / A / S |
| `GetPublicationRequirements` | Q | `readiness.read` | H / A / S |
| `ResolveProductChannelCorrespondence` | C | `readiness.manage` | H / A |
| `ClearProductChannelCorrespondence` | C | `readiness.manage` | H / A |

Automation remains unable to silently supersede a standing human correspondence decision; that is owner business validity after ordinary access, not a second Permission.

## 8.4 Marketplace Listing observation — 2

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListMarketplaceListings` | Q | `offering.read` | H / A / S |
| `GetMarketplaceListing` | Q | `offering.read` | H / A / S |

## 8.5 ListingIntent authoring/media — 7

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListListingIntents` | Q | `offering.read` | H / A / S |
| `GetListingIntent` | Q | `offering.read` | H / A / S |
| `CreateListingIntentDraft` | C | `listing.manage` | H / A |
| `UpdateListingIntentDraft` | C | `listing.manage` | H / A |
| `DiscardListingIntentDraft` | C | `listing.manage` | H / A |
| `SubmitListingIntent` | C | `listing.manage` | H / A |
| `CreateListingIntentMedia` | C | `listing.manage` | H / A |

`listing.manage` never implies `price.manage` or any Availability Permission.

## 8.6 PriceIntent — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListPriceIntents` | Q | `offering.read` | H / A / S |
| `GetPriceIntent` | Q | `offering.read` | H / A / S |
| `CreatePriceIntent` | C | `price.manage` | H / A |

## 8.7 Availability — 9

The two policy operations retain their accepted Operation Matrix labels until final OpenAPI operation-name spelling; W4 does not invent a name merely to complete this table.

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListSellableAvailability` | Q | `availability.read` | H / A / S |
| `GetSellableAvailability` | Q | `availability.read` | H / A / S |
| `ListInventorySources` | Q | `availability.read` | H / A / S |
| `GetInventorySource` | Q | `availability.read` | H / A / S |
| `CreateInventorySource` | C | `availability.manage` | H |
| `UpdateInventorySource` | C | `availability.manage` | H |
| `DeactivateInventorySource` | C | `availability.manage` | H |
| get effective allocation/scope policy | Q | `availability.read` | H / A / S |
| update allocation/scope policy | C | `availability.manage` | H |

## 8.8 Market Intelligence — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListCompetitivePositions` | Q | `market.read` | H / A / S |
| `GetCompetitivePosition` | Q | `market.read` | H / A / S |
| `ListComparableOffers` | Q | `market.read` | H / A / S |

## 8.9 Expected Economics / scenario — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListExpectedEconomics` | Q | `economics.read` | H / A / S |
| `GetExpectedEconomics` | Q | `economics.read` | H / A / S |
| `EvaluatePriceScenario` | side-effect-free C | `economics.read` | H / A / S |

No `economics.scenario.execute` Permission is introduced: the operation is stateless/non-consequential and the accepted `economics.read` capability is sufficient.

## 8.10 Sale Economics — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListSaleEconomics` | Q | `economics.read` | H / A / S |
| `GetSaleEconomics` | Q | `economics.read` | H / A / S |
| `GetEconomicPerformanceSummary` | Q | `economics.read` | H / A / S |

## 8.11 Commercial Policy / Economic Attribution — 5

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `GetCommercialPolicy` | Q | `economics.read` | H / A / S |
| `UpdateCommercialPolicy` | C | `economics.policy.manage` | H |
| `ListEconomicAttributions` | Q | `economics.read` | H / A / S |
| `GetEconomicAttribution` | Q | `economics.read` | H / A / S |
| `ResolveEconomicAttribution` | C | `economics.reconcile` | H |

## 8.12 Controlled Action Governance — 7

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListAuthorizationDecisions` | Q | `governance.read` | H / A / S |
| `GetAuthorizationDecision` | Q | `governance.read` | H / A / S |
| `CreateAuthorizationDecision` | C | `governance.decide` | H |
| `ListAuthorizationDelegations` | Q | `governance.manage` | H |
| `EstablishAuthorizationDelegation` | C | `governance.manage` | H |
| `UpdateAuthorizationDelegation` | C | `governance.manage` | H |
| `RevokeAuthorizationDelegation` | C | `governance.manage` | H |

`governance.decide`/`governance.manage` never grant target-domain write/execute Permissions. Delegation revocation retains its already-ratified fail-safe/monotonic concurrency disposition.

## 8.13 Marketplace Sales — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListMarketplaceSales` | Q | `sales.read` | H / A / S |
| `GetMarketplaceSale` | Q | `sales.read` | H / A / S |
| `ResolveSaleSellingEntityAttribution` | C | `sales.manage` | H |

## 8.14 Business-System Materialization — 5

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListBusinessOrderIntents` | Q | `materialization.read` | H / A / S |
| `GetBusinessOrderIntent` | Q | `materialization.read` | H / A / S |
| `GetBusinessSystemPartyResolution` | Q | `materialization.read` | H / A / S |
| `ResolveBusinessSystemPartyResolution` | C | `materialization.resolve` | H |
| `GetDestinationRealization` | Q | `materialization.read` | H / A / S |

`ResolveDestinationRealization` remains deferred/conditioned on D8 proof and therefore has no W4 runtime mapping.

## 8.15 InvoicingIntent — 2

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListInvoicingIntents` | Q | `materialization.read` | H / A / S |
| `GetInvoicingIntent` | Q | `materialization.read` | H / A / S |

## 8.16 Fulfillment lifecycle / nodes / targets / checkpoints — 13

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListFulfillmentStates` | Q | `fulfillment.read` | H / A / S |
| `GetFulfillmentState` | Q | `fulfillment.read` | H / A / S |
| `ListFulfillmentNodes` | Q | `fulfillment.read` | H / A / S |
| `GetFulfillmentNode` | Q | `fulfillment.read` | H / A / S |
| `CreateFulfillmentNode` | C | `fulfillment.manage` | H |
| `UpdateFulfillmentNode` | C | `fulfillment.manage` | H |
| `DeactivateFulfillmentNode` | C | `fulfillment.manage` | H |
| `GetFulfillmentOperatingTargets` | Q | `fulfillment.read` | H / A / S |
| `UpdateFulfillmentOperatingTargets` | C | `fulfillment.manage` | H |
| `RecordSeparation` | C | `fulfillment.execute` | H |
| `RecordPhysicalConference` | C | `fulfillment.execute` | H OR currently qualified S |
| `RecordPacking` | C | `fulfillment.execute` | H |
| `RecordDispatchHandoff` | C | `fulfillment.execute` | H OR currently qualified S |

## 8.17 Fulfillment artifacts — 2

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListFulfillmentArtifacts` | Q | `fulfillment.execute` | H / A / S |
| `GetFulfillmentArtifact` | Q | `fulfillment.execute` | H / A / S |

This is an explicit W4 resolution of the ratified matrix's ambiguous `human/both read` label toward the read-class baseline H/A/S. It grants artifact visibility only; it confers no physical checkpoint authority.

## 8.18 Shipment observation — 2

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListShipments` | Q | `fulfillment.read` | H / A / S |
| `GetShipment` | Q | `fulfillment.read` | H / A / S |

## 8.19 Post-Sale — 3

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListPostSaleResolutions` | Q | `post_sale.read` | H / A / S |
| `GetPostSaleResolution` | Q | `post_sale.read` | H / A / S |
| `CreatePostSaleResolution` | C | `post_sale.manage` | H / A |

## 8.20 Operational Work — 7

| Operation | Class | Permission | Allowed Principal |
|---|---|---|---|
| `ListWork` | Q | `work.read` | H / A / S |
| `GetWork` | Q | `work.read` | H / A / S |
| `AssignWork` | C | `work.manage` | H / A |
| `ClearWorkAssignment` | C | `work.manage` | H / A |
| `HoldWork` | C | `work.manage` | H / A |
| `ResumeWork` | C | `work.manage` | H / A |
| `EscalateWork` | C | `work.manage` | H / A |

A machine performing business Work coordination is modeled/provisioned as `automation`, not generic `system`.

### 8.21 Matrix coverage

```text
admitted Product operations                           95
W4 operation enforcement rows                        95
unmapped admitted operations                           0
new Product operations added by W4                     0
stored ordinary Permissions                            29
caller Permission hierarchies/wildcards                 0
```

---

## 9. Access/privacy failure grammar

W2 Problem Details remains the one Product problem catalog.

- missing/invalid/untrusted/wrong-audience credential → `401 authentication-required`;
- accepted credential resolving to **zero** MPC Principals → `401 authentication-required`;
- accepted credential resolving to **more than one** MPC Principal → fail closed as existing W2 `500 internal-error`; never select one and never expose conflicting identity details;
- exactly one Principal resolved but current Principal access eligibility denies Product access → `403 access-denied`;
- current access-eligible Principal has no Membership in the path Organization → privacy-preserving `404 resource-not-found` so the API does not confirm another Organization's existence;
- current Membership exists but exact Permission is absent → `403 access-denied`;
- current Membership/Permission exists but Principal kind or required current physical qualification is not admitted → `403 access-denied`;
- path/body/query secondary reference outside the path Organization fails closed without disclosing its real owner; existing W1/W2 cross-Organization privacy grammar remains authoritative.

Do not create `wrong-client-class`, `physical-authority-missing`, duplicate-binding-specific public problems, IdP/provider-specific access errors or a second access taxonomy.

Business `approval-required`, domain rejection/prohibition, unknown/unavailable evidence or external-required is **not** converted into 403 merely because the business says “no”.

---

## 10. Least privilege and sensitive-read review

### 10.1 Permission splits preserved

W4 deliberately preserves these distinct Permissions:

```text
listing.manage != price.manage
fulfillment.read != fulfillment.execute != fulfillment.manage
economics.read != economics.policy.manage != economics.reconcile
governance.read != governance.decide != governance.manage
```

Current names such as `sales.manage` and `post_sale.manage` may be broader linguistically than their one current mutation, but have no wildcard semantics and are not renamed merely for aesthetics.

`offering.read` intentionally covers Listing/ListingIntent/PriceIntent reads under the accepted Offering surface; a separate `price.read` is not invented by symmetry.

### 10.2 Elevated/special read cases

All current reads whose Permission differs from the domain's ordinary read Permission or whose payload is materially sensitive are explicitly reviewed:

1. **FulfillmentArtifacts** — `List/GetFulfillmentArtifacts` require `fulfillment.execute`, not `fulfillment.read`, because labels/handoff artifacts may contain operational/PII-sensitive material. H/A/S is the read-class admission; this grants no checkpoint-establishment authority. A real consumer needing artifact visibility without execute access reopens only this split.
2. **AuthorizationDelegations** — `ListAuthorizationDelegations` remains H + `governance.manage`; standing delegation topology is authorization-management-sensitive and no proven read-only auditor consumer currently justifies a separate Permission. A real auditor consumer needing visibility without mutation authority reopens only this read boundary.
3. **Materialization PII-bearing reads** — `GetBusinessSystemPartyResolution` and `GetDestinationRealization` remain under `materialization.read`; that Permission is explicitly PII-bearing and must be assigned accordingly. A real consumer needing Materialization tracking while excluding party/destination detail reopens the smallest read/Permission boundary.

No speculative `fulfillment.artifact.read`, `governance.delegation.read`, `materialization.pii.read` or equivalent Permission is introduced.

---

## 11. Monotonic/fail-safe revocation is not access bypass

`RevokeAccessRole` and `RevokeAuthorizationDelegation` preserve their ratified monotonic/fail-safe **target-state** semantics: a stale target snapshot must not keep authority alive merely because it changed while revocation was in flight.

They still require at call time:

- valid authenticated human Principal;
- current Principal access eligibility;
- current path-Organization Membership;
- exact `access.manage` or `governance.manage` Permission respectively;
- ordinary same-Organization/reference validity.

Their special concurrency disposition applies only after caller admission.

---

## 12. Negative controls

Later OpenAPI/runtime proof must make at least these defects invalid/unreachable:

1. Keycloak realm/client role `admin` independently authorizes Product operation;
2. OAuth scope named like `price.manage` becomes MPC Permission authority;
3. `portfolio.manage` grants `portfolio.read` **operations** by implication;
4. `fulfillment.manage` grants `fulfillment.execute` or `fulfillment.read` **operations** by implication;
5. mutation response disclosure is treated as a general Get/List/Search grant or Permission inheritance;
6. `listing.manage` permits PriceIntent creation;
7. `governance.decide` permits execution/mutation of the target Intent;
8. Governance delegation lets automation/system invoke a human-only Product operation;
9. ordinary automation with `fulfillment.execute` establishes a physical checkpoint;
10. unqualified or qualification-revoked system with `fulfillment.execute` establishes PhysicalConference/DispatchHandoff;
11. caller body/token/role/scope/client-declared station or evidence source self-asserts physical authority;
12. system Principal is used for Listing/Price/Work business automation instead of an automation Principal;
13. no-Membership caller learns that another Organization exists;
14. foreign Organization body/query reference leaks its real owner;
15. frontend route/button visibility becomes authorization authority;
16. provider seller/user role/scope becomes Product Permission;
17. permission prefix/wildcard grants a new operation automatically;
18. stale token/cached role preserves authority after Principal eligibility/Membership/RoleAssignment/Permission/physical-qualification revocation;
19. disabled Principal can use `/access-context` to enumerate Organizations/Roles/Permissions;
20. duplicate external identity binding is returned as routine 401 or one Principal is selected arbitrarily;
21. `/access-context` requires an Organization before it can discover Organizations;
22. business `approval-required` becomes HTTP access denial;
23. owner business rejection/prohibition is confused with missing ordinary Permission;
24. invalid Principal class is routed into owner business logic instead of failing ordinary access;
25. `physical_system` is introduced as a fourth Principal kind merely for checkpoint authorization;
26. a generic `system_capabilities[]`/machine-capability graph is introduced without a real cross-operation consumer;
27. operation permission is inferred from D1 package name, current handler/middleware or provider endpoint rather than this mapping;
28. deferred/rejected Product operation receives a Permission mapping by symmetry;
29. sensitive read Permissions are split speculatively without a real least-privilege consumer.

---

## 13. W4 outcome / reopen triggers

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED` with bounded D2 identity/access confirmation and Whole-W4 local corrections incorporated.

> **Use external AuthN + unique MPC Principal binding + current Principal access eligibility + current Organization Membership + flat exact Permission + operation-specific Principal-class admission. Keep owner business validity and Governance separate. Require current server-resolved Fulfillment epistemic qualification only for the two system-admitted physical checkpoints; never turn it into a generic IAM/capability framework.**

Reopen only the smallest implicated scope when material evidence shows:

- a real system Principal must perform a currently human-only operation such as Separation/Packing;
- a real automation consumer needs an operation currently human-only;
- a Permission split/name actually prevents least-privilege assignment for a real client rather than looking aesthetically broad;
- a real consumer requires AuthorizationDelegation visibility without governance mutation authority;
- a real consumer requires FulfillmentArtifact visibility without `fulfillment.execute`;
- a real consumer requires Materialization tracking while excluding party/destination PII;
- a physical evidence source cannot be represented/proven without a materially new reusable qualification meaning;
- a real third-party delegated-user client introduces OAuth/consent/client-class requirements outside B2-A.

Do not reopen for Keycloak implementation preference, frontend convenience, middleware shape, role-name aesthetics or hypothetical consumers.

---

## 14. Whole-W4 closure

```text
admitted Product operations                       95 / 95 mapped
stored ordinary Permissions                       29 / 29 used
Principal kinds                                   H / A / S only
Permission hierarchy / wildcard relations         0
additional Product operations                     0
Whole-W4 lead review                              COMPLETE
Fable independent review                          COMPLETE
GPT final adjudication                            CONVERGED
Round 2                                           NOT REQUIRED
D0/D1/D3/D4/D4-R1/D5-B1/W1/W2/W3 reopen          NONE
D2                                                bounded authority confirmation incorporated
W4                                                ACCEPTED / CANONICAL
```

Whole-W4 is closed. The next Wire obligation is **technical non-Product ingress classification**, as routed by the rebaseline router.

Do not begin final OpenAPI/tooling, D6–D9 or implementation out of sequence. Implementation remains blocked until D9.
