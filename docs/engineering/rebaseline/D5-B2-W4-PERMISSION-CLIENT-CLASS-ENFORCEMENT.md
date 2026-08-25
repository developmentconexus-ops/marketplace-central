# D5-B2 — W4 Permission / Principal-Class Enforcement

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Identity/access authority:** D2  
> **Exact current operation mapping:** canonical Product OAD, mechanically proved

## 1. Governing invariant

> **Authentication, current Principal eligibility, Organization Membership, exact ordinary Permission, operation-specific Principal-class admission, owner business disposition, Governance authorization and execution-time validity are separate gates. Possessing one never implies the others.**

W4 owns the ordinary Permission vocabulary and access/admission laws. The exact current per-operation mapping is emitted in the canonical OAD from the accepted operation surface and verified against this current vocabulary/special rules; no second hand-maintained 106-row wire catalog is kept here.

## 2. Principal kinds

Exactly three D2 Principal kinds exist:

```text
H = human
A = automation
S = bounded system
```

- H is accountable interactive human Principal from the accepted OIDC/session path.
- A is non-human business automation/policy execution, including future AI/agent/repricer automation when authorized.
- S is bounded machine/system behavior and is not a generic substitute for A.

OAuth grant type is not Principal kind.

General admission defaults:

- ordinary Q reads may admit H/A/S when their owner semantics allow it;
- consequential business authoring usually admits H/A, unless a more specific accepted operation says otherwise;
- human-only administration/governance/physical-input operations remain H-only;
- stateless side-effect-free `EvaluatePriceScenario` may admit H/A/S despite POST/C encoding.

Exact operation admission remains explicit; there is no inference from HTTP method or Permission name.

## 3. Flat Permission model

Permissions are exact Product capabilities. There is no hierarchy or wildcard implication:

```text
*.manage != *.read
*.execute != *.read
sales.manage != sales.*
```

AccessRole may bundle exact Permissions; role membership never changes Permission meaning. A Permission may authorize invoking a Product capability but does not itself prove business action validity/approval/execution authority.

`authenticated` is a bounded operation condition, not a stored Permission.

## 4. Current ordinary Permission vocabulary — 31

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
performance.read
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
notifications.manage
```

No `notifications.read`, `authorization_requests.read`, generic `admin`, wildcard or provider-specific Permission exists.

A count/name change is a material W4/operation-surface decision and must be explicit.

## 5. Current access evaluation

For an Organization-scoped operation, enforcement evaluates proportionately:

```text
trusted authentication
→ MPC Principal + exact kind H/A/S
→ current Principal Product-access eligibility
→ path Organization
→ current Organization Membership
→ exact ordinary Permission or admitted authenticated/self condition
→ operation-specific Principal-class admission
→ owner business disposition / readiness / evidence sufficiency
→ Governance authorization when required
→ execution-time validity/safety when consequential
```

Access failures do not transfer business authority to middleware. Organization/resource privacy may fail as 404 where contract requires it.

## 6. Platform self-discovery

`GetCurrentAccessContext` is the bounded platform-scoped self discovery operation:

- Principal comes from trusted auth context;
- no caller Principal or Organization selector;
- current Principal eligibility is still required;
- returns only currently visible Organization access context.

It is not a generic user/profile/IAM endpoint.

## 7. Personal Notifications access

### 7.1 Self Inbox

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
```

are:

```text
H only
required_permission = authenticated   // no ordinary Permission
exact self recipient
current path-Organization Membership required
```

Notification possession/read state never grants source capability. Continuation to source is authorized normally against current source owner access.

No `notifications.read` Permission exists.

### 7.2 Organization routing administration

```text
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

are H-only and require `notifications.manage`.

Recipient-candidate discovery remains a bounded routing-support read even though IdentityAccess supplies candidate identity/presentation meaning; it must not require `access.read` or disclose role/Permission internals merely for routing.

`notifications.manage` does not imply access to source business records or other Organization-member administration.

## 8. AuthorizationRequest / Governance access

Current actionable-request operations:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
CreateAuthorizationDecision
```

are all:

```text
H only
governance.decide
ControlledActionGovernance owner
```

The actionable collection/detail population is already filtered to Requests the exact current Principal may decide. `governance.read` is not substituted for this self-actionable decision surface; automation is not admitted to the human approval path.

`CreateAuthorizationDecision` additionally requires current Request decision eligibility/material validity and its typed Request ETag/idempotency semantics. Permission alone never decides the case.

General AuthorizationDecision history remains under the separately admitted `governance.read` operations; actionable decision access does not grant general Governance history browse.

## 9. Performance

The four MarketplacePerformanceIntelligence Q operations require `performance.read` and admit H/A/S under exact Installation/period semantics. `performance.read` does not imply Economics/Market/Offering access.

## 10. Permission-conditioned UI

Frontend may hide/disable controls based on current access context for usability, but this is not authorization. Server operation enforcement remains authoritative. No client-side role-name branching becomes business authority.

## 11. Machine access

A/S authentication resolves one MPC Principal. Non-human access still requires current Principal eligibility, Organization Membership, exact Permission and operation principal-class admission.

Machine bearer does not authorize human-only Governance, access administration, physical checkpoints or Personal Notification self flows by convention.

## 12. Forbidden access shortcuts

Reject:

- role name hard-coded as operation authorization instead of exact Permission;
- `manage` implying `read` or wildcard inheritance;
- IdP role/claim becoming MPC business Permission by convenience;
- caller-supplied Principal identity;
- ambient/default Organization;
- Notification routing revealing role/Permission internals;
- Permission replacing Governance/business disposition;
- business automation impersonating a human;
- generic super-admin bypass not explicitly modeled/authorized.

## 13. Reopen trigger

Reopen W4 only when a real admitted Product operation requires a materially different ordinary capability or Principal-class access law. UI convenience, provider roles or implementation middleware shape are not evidence for new Permission/client classes.
