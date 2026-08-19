# D5-B2 — Whole Technical Ingress Coherence Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD REVIEW CANDIDATE  
> **Review subject:** accepted-in-stage Technical Ingress A External Acquisition + B OAuth / Authorization Ceremony  
> **Authority:** none — review evidence only until operator ratification and canonical filing  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Reviewed:** 2026-08-19

## 1. Review outcome

**Method outcome:** `RESTRUCTURE NOW — TECHNICAL-INGRESS-LOCAL`.

The A+B architecture is globally sound and removes the prior direct provider-route→business-processing defect. No current finding requires D0→W4, D3 or D4 semantic parent reopen. Six bounded Technical-Ingress corrections are required before independent review/canonical acceptance.

## 2. Strong challenges that passed

The lead review attacked the package as one external trust-boundary system and confirmed:

- a native MPC acquisition seam is justified mechanism, not a new D1 business authority;
- provider-specific transports remain adapter-local and a second provider need not emulate HTTP webhook semantics;
- the seven accepted acquisition families are materially tied to current Product 1.0 consumers and do not create a generic provider topic/event ontology;
- provider notification remains an acquisition hint by default; authoritative D4 reread/owner translation retains current truth authority;
- provider delivery dedup/coalescing remains distinct from owner semantic idempotency;
- Payment/Post-Sale occurrence preservation remains guarded by D3/D4 evidence-edge rules rather than webhook event sourcing;
- OAuth is correctly a separate authorization ceremony lane, not an acquisition signal;
- one current Authorization Attempt per Installation/application is a bounded concurrency rule rather than a workflow framework;
- callback-time revalidation of the initiating Principal's current access prevents `state` from becoming eternal Product authority;
- same-seller reauthorization versus different-seller fail-closed behavior preserves canonical D4 Installation identity;
- complete credential-generation activation plus stale-refresh protection is necessary and remains D7-mechanism-neutral;
- provider technical routes remain outside Product OpenAPI/SDK and provider credentials never authenticate Product Principals;
- Structural Inversion against current controller/webhook/OAuth implementation: **PASS**.

---

## 3. IG-G1 — Signal disposition must distinguish unverified, non-admitted, admitted-attributed and admitted-unattributed — REVISE

### Gap

The accepted custody law requires recoverable responsibility before positive acknowledgement, while quarantine currently names correlation failures and a temporarily unavailable selected mapping. Read literally, this can pressure the system to durably quarantine **unknown/DEFER/REJECT topics** or even protocol-unverified traffic merely to avoid provider retries.

That would turn quarantine into a latent feature backlog and durable attacker-controlled payload sink.

### Corrected invariant

Inbound signal handling has four semantic classes:

```text
A. malformed / protocol-unverified / origin-authentication failed
   → protocol rejection/fail-closed; no MPC acquisition custody

B. protocol-admissible but topic/resource is not admitted by the closed ingress matrix
   → deliberate terminal technical non-processing disposition;
     provider acknowledgement may be positive when the provider protocol requires it;
     no durable acquisition/quarantine obligation is created

C. admitted signal + exact namespace binding
   → recoverable attributed acquisition custody before positive acknowledgement

D. admitted, protocol-admissible signal + correlation ambiguous/missing/contradictory
   → bounded technical quarantine before positive acknowledgement
```

A non-admitted topic never becomes future work merely because it arrived. A protocol-unverified signal never gains quarantine durability by convenience.

### Alternatives

1. quarantine everything unknown — rejected: generic integration backlog + abuse amplification;
2. non-2xx every unsupported topic forever — rejected: creates pointless retry storms for deliberately unconsumed provider surface;
3. **closed disposition matrix above — selected Global Maximum.**

No parent reopen.

---

## 4. IG-G2 — Namespace attribution must not collapse into Installation business activation or credential availability — REVISE

### Gap

Ingress-A repeatedly says a signal requires a "current compatible binding". That is correct for namespace correlation but may be misread as requiring an **active/authorized business Installation**.

A deactivated Installation, temporarily auth-invalid credential set, or seller connection awaiting reauthorization can still receive late provider Payment/refund/Post-Sale/Shipment signals belonging to historical transactions. If deactivation/auth failure erased attribution, MPC could lose material evidence precisely when recovery is needed.

### Corrected invariant

Separate three meanings:

```text
namespace attribution binding
    != Portfolio business participation/active posture
    != current provider credential usability
```

- an unambiguous retained Installation↔provider seller/account namespace binding may still attribute late signals after business deactivation;
- deactivation never erases historical namespace identity/correlation;
- an attributed signal whose authoritative reread is temporarily impossible because credentials are auth-invalid/unavailable remains **Organization-attributed recoverable acquisition**, not unbound quarantine;
- owner/D4 semantics decide whether the late evidence is relevant/actionable; ingress does not silently discard it because the Installation is not currently sell-enabled;
- OAuth begin/reauthorization may remain disallowed for a deactivated Installation under its own lane; this does not weaken acquisition attribution.

This preserves D2/D4 identity/history without granting deactivated Installations new business authority.

No parent reopen.

---

## 5. IG-G3 — Push/poll/recovery convergence is path convergence, not discovery-source equivalence — REVISE

### Gap

The current diagram lists webhook, missed-feed, scheduled reconciliation, cold-start scan and manual recovery converging on the same acquisition family. That is a sound architecture shape, but it can be read as promising that **every acquisition family has every discovery/recovery source**.

D4 explicitly leaves coverage operation-scoped; some Payment/Claim/competition surfaces may lack a complete enumeration or recovery feed.

### Corrected invariant

> **When multiple discovery mechanisms are actually admitted by D4 for the same semantic acquisition, they converge on one typed acquisition path. The existence of that shared path never invents polling/enumeration/recovery capability or completeness for a family whose provider contract does not establish it.**

Each acquisition family must retain an explicit source/recovery availability statement during D7/D8 proof. Missing recovery coverage remains partial/Unknown rather than a generic scan capability.

No parent reopen.

---

## 6. IG-G4 — Successful OAuth credential activation must re-enable technical acquisition recovery without creating Sync/domain-event authority — REVISE

### Gap

A and B are individually coherent but currently under-specify their interaction.

An attributed acquisition may be blocked because provider credentials are expired/revoked/auth-invalid. A later successful same-seller reauthorization activates a usable credential generation, but nothing in the combined contract says what awakens the pending acquisition/reconciliation work. Depending only on a future unrelated notification can create a silent stall.

### Corrected invariant

> **Successful credential-generation activation may awaken Installation-scoped technical capability revalidation, bootstrap and recovery of already-admitted acquisition work. This is technical recovery machinery, not a Product `Sync`/`Refresh` operation, not a D3 business event and not proof that any owner meaning changed.**

Rules:

- initial authorization may make the Installation eligible for bounded initial acquisition/capability discovery only through sources already admitted by D4/Ingress-A;
- same-seller reauthorization may wake pending attributed acquisition/reconciliation blocked on auth;
- no owner state is committed until the normal authoritative D4 acquisition/translation path runs;
- exact scheduler/worker/queue realization remains D7.

This closes a cross-lane recovery gap without adding Product surface.

No parent reopen.

---

## 7. IG-G5 — OAuth binding changes need durable non-secret security provenance — REVISE

### Gap

Authorization Attempt is correctly technical and secrets are correctly excluded from Product/history/log surfaces. But a successful provider seller/account binding or credential-generation replacement materially changes an external trust boundary. If the ephemeral attempt disappears entirely, later operators cannot explain **who initiated the binding, which seller was proven, when generation changed, or why a reauthorization replaced prior credentials**.

### Corrected invariant

Preserve the smallest durable **non-secret authorization/binding lineage** required for security/audit explanation, proportionately including:

```text
MarketplaceInstallation
provider application/config identity
initiating MPC Principal
provider seller/account identity when proven
initial-bind | same-seller-reauthorization | superseded/failed technical outcome
attempt/generation correlation
material timestamps
```

Never preserve authorization code, client secret, access token, refresh token or PKCE verifier in that lineage.

This lineage is D4/D2 technical trust/configuration provenance, not a Product AuthorizationAttempt resource, business event store or generic credential history platform. Exact persistence/retention mechanics remain D7.

No parent reopen.

---

## 8. IG-G6 — Product-authenticated OAuth begin needs an explicit technical-route access law without becoming a 96th Product operation — REVISE

### Gap

The accepted OAuth begin consumes current Product AuthN, human Principal semantics, Membership and `portfolio.manage`, but the route is intentionally **not** a Product API operation. Without an explicit fence, implementations can reach one of two wrong conclusions:

1. add `BeginOAuth` as a 96th Product operation to make W4 authorization convenient; or
2. treat the technical begin as outside Product access rules and authorize it only through browser/session/provider state.

### Corrected invariant

> **OAuth begin is a Product-authenticated technical action, not a Product business operation. It consumes the existing D2/W4 current-access facts (`H`, Principal eligibility, path/declared Organization Membership, `portfolio.manage`, exact Installation scope) as an initiation gate without entering the 95-operation Product matrix or Product OpenAPI.**

Additional route law:

- the technical begin must carry explicit Organization + MarketplaceInstallation scope and fail closed on mismatch;
- provider-facing callback/acquisition routes use provider protocol trust/correlation, not Product bearer authority;
- technical routes must be collision-separated from Product roots so `/organizations/...` Product paths cannot be confused with callbacks;
- exact technical prefix/host spelling is **DEFER SAFELY** to final ingress wire closure; separate host deployment is D7, not required now.

No new Permission is introduced and W4's 95/95 Product mapping remains intact.

No parent reopen.

---

## 9. Global Maximum after IG-G1…G6

```text
provider protocol
        ↓
protocol verification + closed admission
        ↓
[non-admitted terminal disposition]
        OR
[admitted exact namespace attribution]
        OR
[admitted bounded quarantine]
        ↓
recoverable typed acquisition
        ↓
authoritative D4 read
        ↓
owner meaning

Product-authenticated H + portfolio.manage
        ↓
technical OAuth begin
        ↓
server-bound attempt
        ↓
provider callback + seller proof
        ↓
current-authority revalidation
        ↓
generation-safe credential activation
        ↓
non-secret binding lineage
        ↓
technical catch-up/recovery wake-up
```

while preserving:

```text
no generic Integration/Webhook/ExternalEvent authority
no provider topic ontology in D1
no Product Sync/Refresh operation
no OAuth callback in acquisition inbox
no 96th Product operation
no secret-bearing audit history
no deactivation-induced loss of historical source attribution
no invented provider completeness/recovery capability
```

---

## 10. Reopen classification

- D0/D1/D2/D3/D4/D4-R1/D5-B1/W1/W2/W3/W4/Operation Matrix: **NO REOPEN**;
- Technical Ingress A+B: **targeted local corrections IG-G1…G6**;
- D7 owns concrete custody/quarantine/credential generation/locking/scheduler/storage realization;
- D8 owns real provider proof of acknowledgement/recovery/OAuth/concurrency paths.

---

## 11. Recommended next action

Operator reviews/ratifies/revises **IG-G1…IG-G6** as the Whole-Ingress lead correction direction.

If ratified, because the package crosses provider authentication, tenant attribution, durable custody and credential trust boundaries, run **one independent Fable Whole-Ingress challenge** over corrected direction before canonical filing.

Do not micro-review individual findings. Reviewer output remains evidence only.
