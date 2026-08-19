# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 + W3 + W4 CANONICAL; Technical Ingress A External Acquisition ACCEPTED IN-STAGE; OAuth / Authorization Ceremony lane = NEXT**  
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
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — accepted technical-ingress design home; Ingress-A accepted in-stage, OAuth lane next
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

Former Whole-W2/W3/W4 review candidates/staging are absent from the active authority tree; Git history is the archive. `AI-DIALOG.md` is protocol-only and is not architecture authority.

`docs/engineering/rebaseline/cockpit.html` is a **non-authoritative visual projection** and is synchronized only after canonical status changes.

Legacy/current code, OpenAPI, IdP roles/scopes, middleware, provider routes/topics and frontend guards remain evidence only.

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
        B OAuth / Authorization Ceremony                 NEXT
        Whole-Ingress coherence                          BLOCKED BY B
      Final Problem/media consistency                    BLOCKED BY SEQUENCE
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2, Whole-W3 and Whole-W4 adversarial review cycles are complete, operator-ratified and incorporated into canonical authority.

---

## 3. Load-bearing current authority

### 3.1 D2 / W4 access

- D2 owns current Principal access eligibility as Principal-scoped revocable identity/access state.
- Principal kinds remain exactly `human | automation | system`.
- W4 maps **95/95** admitted Product operations with **29 flat exact Permissions**.
- Product AuthN, Principal eligibility, Membership, Permission, Principal kind, physical qualification, owner business disposition and Governance are separate gates.
- IdP/OAuth/provider roles/scopes never independently grant Product access.

### 3.2 Canonical Product wire

- W1 owns Product resource/path/HTTP grammar; Product business paths use `/organizations/{organization_id}/...`, with `/access-context` as the bounded self-only platform Q.
- W2 owns request/response schema grammar and the single Product Problem Details catalog.
- W3 owns all 26 admitted List/Search collection/query/cursor semantics.
- W4 owns Permission→operation/client-class enforcement.
- Provider/business-system protocol ingress remains outside Product API roots and Product SDK/OpenAPI business operations.

### 3.3 Technical Ingress-A — accepted in-stage

Ingress-A establishes the MPC-native external acquisition seam:

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

- provider signal/notification is not automatically provider truth or a D3 domain event;
- Organization derives only from an exact current MPC MarketplaceInstallation/SourceInstance binding;
- positive provider acknowledgement means recoverable technical custody/quarantine only, never business acceptance/convergence;
- push, missed-feed/recovery, polling/reconciliation and cold-start discovery converge on the same typed acquisition family;
- provider delivery dedup/coalescing never replaces owner semantic idempotency;
- arbitrary provider `resource` URLs are never fetched directly; adapter uses closed provider grammar and sanctioned D4 reads;
- bounded technical quarantine may hold unattributable signals but cannot dispatch D1 owner work before exact correlation;
- raw provider DTO/PII retention is minimized;
- no generic `Webhook`, `ExternalEvent`, `Integration`, `ProviderResource` or subscription business authority is introduced.

Accepted native acquisition families:

```text
AcquireMarketplaceListing
AcquireMarketplacePrice
AcquireMarketplaceSale
AcquireMarketplaceShipment
AcquireMarketplacePayment
AcquireMarketplacePostSaleClaim
AcquireMarketplaceCompetitivePosition
```

Current Mercado Livre admission baseline:

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

OAuth/authorization callbacks are deliberately excluded from the acquisition inbox and remain the next technical-ingress lane.

---

## 4. Prohibited now

While OAuth / Authorization Ceremony is next:

- do not begin Whole-Ingress review, final Problem/media consistency, OpenAPI/tooling closure, D6–D9 or implementation;
- do not reopen D0→W4 or Ingress-A for naming/style/current-code convenience;
- do not create generic Integration/Ingress/Webhook/Event business authority or public generic provider resource/event model;
- do not add provider topics by availability/symmetry; Ingress-A admission matrix remains closed-world;
- do not let provider credentials/signatures authenticate Product Principals or Product bearer tokens authenticate provider callbacks;
- do not let OAuth callback enter the External Acquisition Signal inbox merely for uniformity;
- do not choose D7 queue/database/broker/worker/quarantine/secret/retry/deployment realization;
- do not include provider technical routes in the Product SDK/OpenAPI business operation surface;
- do not choose exact OpenAPI minor/generator or numeric pagination defaults yet.

---

## 5. Exact next action

**Derive Technical Ingress-B — OAuth / Authorization Ceremony lane from D2/D4/D5/W4 without turning provider authorization protocol into a Product business operation or acquisition event.**

The OAuth lane must decide proportionately:

1. which Product-authorized human action may initiate authorization/re-authorization for an existing MarketplaceInstallation and what current W4 access is revalidated;
2. the boundary between the Product-triggered **begin ceremony** and provider-facing **callback**;
3. server-issued opaque `state` / correlation semantics, replay/expiry/single-use requirements and what authority it may carry;
4. whether Organization/MarketplaceInstallation/initiating Principal/post-callback destination are bound server-side rather than trusted from callback query data;
5. callback verification and code exchange semantics under provider protocol authority, never Product bearer authority;
6. authoritative provider seller/account identity proof after token exchange;
7. same-seller reauthorization versus different-seller mismatch/rebinding rules under canonical D4 Installation identity;
8. revalidation at callback time of current Principal access eligibility, current Organization Membership and `portfolio.manage` so a stale initiation cannot become eternal authority;
9. credential/token secrecy and what may be persisted/audited without exposing secrets;
10. callback replay, stale/expired state, provider denial/error, token exchange failure and ambiguous partial-completion semantics;
11. callback acknowledgement/redirect result semantics without leaking provider protocol vocabulary into Product API authority;
12. route/host boundary and CSRF/open-redirect/tenant-confusion negative controls without choosing D7 framework/middleware implementation.

OAuth begin/callback must remain distinct from External Acquisition Ingress-A while sharing the same non-Product protocol boundary.

After OAuth lane convergence:

1. run one Whole-Ingress adversarial coherence review across Acquisition + OAuth lanes;
2. obtain final operator ratification and consolidate Technical Ingress as canonical;
3. continue final Problem/media consistency;
4. then decide one machine-readable Product OpenAPI authority/tooling/minor version.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4 are canonical;
- `D5-B2-TECHNICAL-INGRESS.md` is the single Technical Ingress design home;
- Ingress-A External Acquisition is accepted in-stage/operator-ratified;
- its seven native acquisition families and Mercado Livre ADMIT/DEFER/REJECT matrix are current;
- OAuth callbacks are not acquisition signals and have not yet been accepted;
- **OAuth / Authorization Ceremony lane is the exact next action**;
- Whole-Ingress/final Problem-media/OpenAPI/D6–D9/implementation remain blocked by sequence.

If not, the active authority tree is inconsistent.
