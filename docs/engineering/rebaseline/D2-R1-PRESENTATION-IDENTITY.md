# D2-R1 — Presentation Identity

> **Status:** OPERATOR-APPROVED / CANDIDATE — targeted D6-exposed clarification
> **Parent authority:** `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
> **Consumer:** D6 frontend human-facing Organization/access context
> **Wire realization:** canonical Product OAD under D5

## 1. Why this clarification exists

D6-B1 exposed a bounded presentation gap while mapping the App Shell and access-management interactions: canonical IDs are correct for authority and correlation, but a human client cannot reliably distinguish Organizations, Principals or AccessRoles when the Product read contract exposes only opaque IDs/keys.

This is not a new identity model, directory, profile domain or authorization mechanism. It is presentation metadata needed by already-admitted human Product interactions.

## 2. Governing invariant

> **Human-readable presentation labels may identify an Organization, Principal or AccessRole to a person, but canonical IDs/keys remain the only identity, scope and correlation authority. A presentation label is mutable, non-unique and never authenticates, authorizes, scopes, joins or collapses identities.**

Consequences:

- `organization_id` remains the Organization identity/isolation root;
- `principal_id` remains the accountable MPC Principal identity;
- `role_key` remains the Product-defined AccessRole identity used for assignment/revocation;
- equal `display_name` values never imply equal identities;
- changing `display_name` never changes Membership, RoleAssignment, Permission, audit attribution or historical identity;
- frontend/client-local aliases, OIDC names/emails and provider account labels never substitute for canonical MPC identity;
- exact persistence, refresh and administrative realization of labels remains later technical work and creates no Product 1.0 profile-management surface by implication.

## 3. Product read requirement

For the already-admitted human-facing access flows, the Product read contract must provide non-empty presentation labels proportionately:

- current Principal → `display_name`;
- each Organization in current access context → `display_name`;
- each Organization member returned for access administration → `display_name`;
- each AccessRole returned for access administration → `display_name`.

The client submits and keys interactions by canonical ID/key, never by label.

No new Product operation, ordinary Permission, Principal kind or domain is admitted by this clarification.

## 4. Falsifiers

This clarification fails if any realization permits:

1. Organization scope/cache identity to be keyed by `display_name` rather than `organization_id`;
2. Principal attribution or role mutation to target a person by `display_name` rather than `principal_id`;
3. AccessRole assignment/revocation to use `display_name` rather than `role_key`;
4. two equal labels to collapse distinct canonical identities;
5. a label change to modify access, authority or historical attribution;
6. absence of a required human-readable label to be repaired by a frontend hardcoded/local alias that becomes apparent Product truth.

## 5. Boundary

This amendment changes presentation completeness only. D2 identity/isolation semantics remain otherwise unchanged. **D2-R1 itself admits no new Product operation, ordinary Permission, Principal kind, business domain or authorization authority.** Later bounded amendments may extend the Product surface without changing this presentation-identity rule; D6 remains a client of the resulting authorities rather than becoming identity authority.
