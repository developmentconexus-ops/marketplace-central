# D5-B2 — Technical Non-Product Ingress

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-API.md`  
> **External semantics:** D4 + D4-R1  
> **Product wire census:** Technical Ingress remains outside the current **106 Product operations**

## 1. Purpose

This owner defines inbound technical surfaces that are **not ordinary Product Principals invoking Product API operations**.

Two lanes remain distinct:

```text
A. External acquisition ingress
provider discovery/notification/recovery signal
→ bounded MPC custody/correlation
→ authoritative D4 acquisition/reread
→ consumer-owned translation
→ D1 owner commits meaning
→ D3 communication only after owner commit

B. OAuth/authorization ceremony
Product-authorized human initiation
→ server-bound Authorization Attempt
→ provider authorization protocol
→ callback
→ current MPC authority revalidation
→ provider seller/account proof
→ safe credential generation activation
```

No generic Integration/OAuth/Event/Workflow business domain is introduced.

## 2. Governing invariants

1. D4 owns provider semantics; Technical Ingress owns only wire/trust-boundary mechanism.
2. Provider transport/DTO/topic/auth remains adapter-local.
3. Provider signal != D3 domain event and != current provider truth.
4. Organization is never trusted from provider payload. Marketplace ingress resolves exact MarketplaceInstallation/source binding first.
5. Unknown provider topic/resource never becomes generic `ExternalEvent` or provider property bag.
6. Positive callback/notification acknowledgement means at most recoverable technical custody/disposition, never business acceptance/completion/convergence.
7. Push/poll/recovery mechanisms exist only where D4 proves provider capability; no invented completeness.
8. Delivery dedupe is not business idempotency.
9. Product authentication and provider callback trust are separate.
10. OAuth `state` is correlation, not standing Product authorization; callback completion revalidates current authority.
11. Credential activation is generation-safe; stale authorization/refresh cannot overwrite newer active generation.
12. Provider-specific protocol may reuse infrastructure seams but never becomes one speculative universal ingress framework.

## 3. Lane A — external acquisition

### 3.1 Provider adapter

The provider adapter may own:

- actual supported HTTP/queue/event delivery mechanics;
- provider-origin verification/auth evidence;
- DTO parsing;
- closed provider topic/resource grammar;
- provider delivery discriminator when supplied;
- conversion of provider resource syntax into bounded native identity/reference.

It does not:

- infer Product owner semantics from topic names;
- infer Organization from payload;
- commit D1 owner state;
- expose raw provider DTOs outside the boundary;
- fetch arbitrary callback-supplied URLs;
- fabricate provider capability/coverage.

### 3.2 MarketplaceInstallation-scoped custody

Marketplace acquisition is correlated to an exact current/historical MarketplaceInstallation namespace binding before Organization-owned processing.

A valid provider-native signal may be:

```text
accepted into recoverable attributed custody
quarantined before safe attribution under bounded technical policy
terminally rejected as invalid/unprocessable technical input
```

Acknowledgement semantics never claim Product/business success.

### 3.3 Typed acquisition requests

After admissible attribution, internal acquisition uses closed typed families tied to accepted D4 authoritative read/coverage contracts. Provider topic vocabulary does not become the internal Product domain model.

Examples may include Listing/Sale/Shipment/Claim/payment/marketplace account evidence only when a named D4 consumer contract exists. No generic `ProviderResource {type,payload}` exists.

### 3.4 Authoritative reread / translation

Provider push evidence normally triggers an accepted D4 reread/acquisition when current material truth is needed. Consumer owner then translates/commits its own meaning. Only after owner commit may D3 emit a domain fact/attention reaction.

### 3.5 Duplicate/order/recovery

Provider duplicate/out-of-order delivery remains acquisition behavior. MPC technical custody may dedupe/coalesce work, but owner/source identity, current reread and D3 semantic idempotency protect business correctness.

Missed notification never proves no source change; polling/backfill/enumeration/recovery exists only where the provider contract proves it. Provider/webhook arrival is not source completeness.

## 4. Pre-attribution quarantine

Quarantine is a narrow technical state for input that cannot yet be safely attributed/processed. It is not Organization Work, Product backlog or provider-payload archive.

Requirements:

- bounded retention and minimized sensitive data;
- no business meaning/owner assignment before safe attribution;
- no automatic replay into an arbitrary Organization;
- observability sufficient for operational/security diagnosis;
- explicit terminal disposition when safe processing cannot be established.

Exact storage/TTL/operator mechanism remains D7.

## 5. Deactivation / retained attribution

Marketplace Installation deactivation removes current business participation/authority. Historical non-secret namespace attribution/correlation required for evidence recovery may remain.

Restoring same-seller technical credentials may be possible only under the accepted current Installation/provider account binding and does not silently reactivate listings, automation, routing or provider writes.

## 6. Lane B — OAuth / provider authorization ceremony

### 6.1 Product-authorized initiation

Only a current Product-authorized human may initiate provider authorization for the exact Organization/MarketplaceInstallation. Browser never supplies/stores provider client secret/refresh token.

The server creates one bounded Authorization Attempt containing only the correlation/current-context facts required for the ceremony. Attempt state is not a Product business resource/API by symmetry.

### 6.2 Transaction binding

Provider authorization uses the strongest accepted provider-supported transaction binding needed for callback injection/replay safety. Current Mercado Livre flow uses state plus PKCE where supported/selected.

Callback input is protocol evidence, not authorization to bind an arbitrary Organization/Installation.

### 6.3 Callback revalidation

Before credential activation, the server revalidates proportionately:

```text
valid unconsumed Authorization Attempt
+ current initiating Principal/Organization authority where required
+ exact target MarketplaceInstallation still eligible
+ provider callback state/PKCE/protocol checks
+ provider-authoritative seller/account identity proof
+ namespace compatibility with existing binding when reauthorizing
```

A provider login to another seller/account cannot silently rebind an existing Installation namespace.

### 6.4 Credential generation safety

Credential storage/activation uses generation-safe semantics:

- complete credential generation activates atomically enough for consumers;
- older Authorization Attempt cannot overwrite a newer generation;
- stale refresh cannot overwrite a newer generation;
- current Mercado Livre single-use refresh requires serialized correctness at generation level;
- secret values never enter Product API, logs, commits, documentation or non-secret lineage.

D7 chooses lock/CAS/transaction/secret-store realization.

## 7. Historical lineage

Keep only bounded non-secret authorization/binding lineage needed for explanation/recovery. Historical lineage is never the current namespace/credential authority; current Installation/D4 binding is.

## 8. Authored-media delivery fence

ListingIntent authored-media byte delivery is a separately justified authenticated technical presentation surface under current Product authorization. It does not become a Product business operation, generic Asset API or provider ingress lane.

Source-media locators, authored-media access references and provider-observed media remain distinct trust concepts under W2/D4.

## 9. Product/OAD separation

Technical Ingress routes/mechanics:

- are outside `contracts/api/product/openapi.yaml` unless a later explicit Product client operation is admitted;
- do not add ordinary Permissions by HTTP symmetry;
- do not generate Product SDK operations;
- may reuse technical server/platform infrastructure after D9 without merging business authority.

## 10. Explicit non-goals

Reject by default:

- generic webhook/ExternalEvent business domain;
- generic Integration/Product credential CRUD API;
- provider topic = MPC domain event;
- one provider endpoint per D1 owner;
- generic OAuth workflow engine;
- public sync/refresh/recovery Product operations for technical machinery;
- generic source callback seam without a real provider/consumer;
- payload-derived Organization;
- callback-delivered arbitrary fetch URLs;
- universal subscription/provider-resource model.

## 11. Reopen trigger

Reopen Technical Ingress only when a real provider/source requires a materially different trust/acquisition/authorization contract that cannot be represented by these two lanes without weakening authority/safety. Provider implementation variety alone does not justify a generic framework.
