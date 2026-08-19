# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 + W3 + W4 CANONICAL; Technical Ingress A+B ACCEPTED IN-STAGE; Whole-Ingress adversarial coherence = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-19

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this router
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
8. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
9. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
10. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
11. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
12. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
13. `docs/engineering/rebaseline/D5-API.md`
14. `docs/engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md`
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` — canonical W1
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — accepted Technical Ingress A+B design home
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

Technical Ingress remains **accepted in-stage**, not canonical, until Whole-Ingress coherence and final operator ratification. `AI-DIALOG.md` is protocol-only unless a review cycle is explicitly opened. Cockpit remains non-authoritative and is synchronized only after canonical status changes.

Legacy/current code, provider routes/topics, OAuth handler shape, middleware, Product OpenAPI and frontend behavior remain evidence only.

---

## 2. Program state

```text
D0 — Product / System Definition                         CLOSED / ACCEPTED
D1 — Domains / Boundaries                                CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership                  CLOSED / ACCEPTED
D3 — Communication / Events                              CLOSED / ACCEPTED
D4 — External Integrations                               CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring            ACCEPTED / CANONICAL
Decision Reconciliation                                  ACCEPTED / CANONICAL
D5 — API                                                  OPEN / ACTIVE
  B1 Semantic API Model                                  ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface                 OPEN / ACTIVE
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Query / Cursor Grammar            ACCEPTED / CANONICAL
      W4 Permission / Client-Class Enforcement           ACCEPTED / CANONICAL
      Technical non-Product ingress                      OPEN / ACTIVE
        A External Acquisition Ingress                   ACCEPTED IN-STAGE / OPERATOR-RATIFIED
        B OAuth / Authorization Ceremony                 ACCEPTED IN-STAGE / OPERATOR-RATIFIED
        Whole-Ingress adversarial coherence              NEXT
      Final Problem/media consistency                    BLOCKED BY INGRESS
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Load-bearing Technical Ingress authority

### 3.1 Lane A — External Acquisition

```text
provider-specific inbound transport adapter
        ↓
MPC-native recoverable acquisition ingress
        ↓
closed typed acquisition request
        ↓
accepted D4 authoritative acquisition/reread
        ↓
consumer-owned semantic translation
        ↓
D1 owner commits meaning
```

Binding laws:

- provider signal is neither provider current truth nor a D3 domain event;
- exact current MarketplaceInstallation/SourceInstance binding determines Organization;
- positive provider acknowledgement means recoverable technical custody/quarantine only;
- push, recovery, polling, reconciliation and cold-start discovery converge on the same typed acquisition path;
- delivery dedup/coalescing never replaces owner semantic idempotency;
- arbitrary provider resource URLs are never fetched directly;
- bounded unbound quarantine cannot dispatch D1 owner state;
- provider PII/raw payload retention is minimized;
- no generic Webhook/ExternalEvent/Integration/ProviderResource authority.

Accepted acquisition families:

```text
AcquireMarketplaceListing
AcquireMarketplacePrice
AcquireMarketplaceSale
AcquireMarketplaceShipment
AcquireMarketplacePayment
AcquireMarketplacePostSaleClaim
AcquireMarketplaceCompetitivePosition
```

Current Mercado Livre admission:

```text
ADMIT:
  orders_v2
  shipments
  items
  items_prices
  payments
  post_purchase:claims
  post_purchase:claims_actions
  catalog_item_competition_status

DEFER:
  stock-location
  fbm_stock_operations
  flex-handshakes
  invoices
  user-products-families
  catalog_suggestions
  price_suggestion
  promotion/best-price families
  provider image-error/picture feeds

REJECT BASELINE:
  legacy orders when orders_v2 selected
  created_orders checkout-only feed
  feedback/reputation
  messages
  questions/Q&A
  leads/vertical signals
  unknown/unclassified topic
```

### 3.2 Lane B — OAuth / Authorization Ceremony

```text
Product-authorized human begin
        ↓
one current server-bound Authorization Attempt
        ↓
provider Authorization Code + state + selected PKCE protection
        ↓
callback-time current MPC authority revalidation
        ↓
server-to-server token exchange
        ↓
provider-authoritative seller/account proof
        ↓
initial bind or same-seller reauthorization only
        ↓
generation-safe complete credential activation
```

Binding laws:

- begin requires current authenticated human + Principal eligibility + Membership + `portfolio.manage` + target Installation;
- OAuth callback does not create Organization or MarketplaceInstallation;
- state is opaque/high-entropy/single-use/finite-lived correlation, never durable authorization;
- one current unfinished attempt per Installation/provider application; a newer begin supersedes the old attempt;
- callback revalidates initiator authority and Installation eligibility regardless of current browser session;
- provider seller identity is authoritative external namespace evidence, never MPC Principal/Organization identity;
- same-seller reauthorization may replace credentials; different-seller result fails closed;
- credential activation is complete-generation only; stale refresh/older attempt cannot overwrite a newer generation;
- Product Problem Details does not absorb provider OAuth protocol errors;
- refresh is outbound D4/D7 credential lifecycle, not ingress or Product operation.

---

## 4. Prohibited now

Until Whole-Ingress coherence converges:

- do not canonicalize Technical Ingress A/B silently;
- do not begin final Problem/media consistency, OpenAPI/tooling, D6–D9 or implementation;
- do not create generic Integration/Ingress/Webhook/ExternalEvent/OAuth workflow authority;
- do not add provider topics by availability/symmetry;
- do not let provider protocol credentials authenticate Product Principals or Product tokens authenticate callbacks;
- do not route OAuth into the acquisition inbox;
- do not choose D7 queue/database/broker/quarantine/secret/lock/refresh realization;
- do not include provider technical routes in the Product SDK/OpenAPI business surface;
- do not reopen D0→W4 for naming/current-code convenience.

---

## 5. Exact next action

**Run one Whole-Ingress adversarial coherence review over Technical Ingress A+B as one external trust-boundary system.**

Challenge at minimum:

1. whether native acquisition ingress is needed as a real MPC seam or is accidental framework machinery;
2. whether seven typed acquisition families are the smallest honest vocabulary and no provider topic leaks into D1 ontology;
3. whether push/poll/recovery convergence preserves D3 occurrence/recovery semantics;
4. whether positive acknowledgement after recoverable custody is both necessary and sufficient without prematurely choosing D7 storage;
5. quarantine correctness, especially no-Organization attribution and eventual correlation/recovery;
6. provider namespace → Installation/SourceInstance → Organization correlation and all ambiguity/mismatch cases;
7. whether acquisition coalescing can accidentally lose material Payment/Post-Sale occurrences;
8. Product API/OpenAPI/W4 separation from provider technical routes;
9. OAuth begin authority, state semantics, one-current-attempt rule and callback-time revocation/current-authority revalidation;
10. provider seller/account proof and initial-bind versus same-seller reauthorization;
11. credential-generation concurrency, especially stale refresh versus newer authorization;
12. replay, partial exchange, provider denial and safe post-callback navigation;
13. whether OAuth and Acquisition lanes share only mechanism or accidentally create a generic ingress authority;
14. Structural Inversion against current webhook/controller/OAuth implementation;
15. whether any finding truly requires D2/D3/D4/W1–W4 reopen rather than a Technical-Ingress-local correction.

Because this package crosses external trust boundaries, independent Fable review is expected after lead adjudication if the lead review leaves no prerequisite split.

After Whole-Ingress convergence + final operator ratification:

1. make `D5-B2-TECHNICAL-INGRESS.md` canonical;
2. archive/remove any review candidate and reset review channel;
3. update cockpit as non-authoritative projection;
4. advance to final Problem/media consistency;
5. then close one Product OpenAPI authority/tooling/minor-version decision.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- `D5-B2-TECHNICAL-INGRESS.md` is the single Technical Ingress design home;
- Acquisition-A and OAuth-B are both accepted in-stage/operator-ratified;
- the acquisition matrix remains closed-world with seven typed families;
- OAuth is a separate lane with current-authority revalidation and generation-safe credentials;
- **Whole-Ingress adversarial coherence is the exact next action**;
- final Problem/media/OpenAPI/D6–D9/implementation remain blocked.

If not, the active authority tree is inconsistent.
