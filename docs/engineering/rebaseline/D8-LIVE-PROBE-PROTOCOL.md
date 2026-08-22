# D8 — Controlled Live Probe Protocol

> **Status:** EXECUTION AUTHORIZED / PROBES NOT YET EXECUTED  
> **Parent:** [D8-R1 Proof Closure & Implementation-Readiness Coherence](D8-R1-PROOF-CLOSURE-COHERENCE.md)  
> **Operator authorization:** 2026-08-22 — execute the real probes needed to avoid theory-only architecture closure  
> **Scope:** D4-deferred D8 external-contract probes only; no Product implementation  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose

This protocol converts the D8 P1–P6 probe ledger into a controlled execution contract. The goal is not to make the accepted architecture pass. The goal is to **falsify it with real Mercado Livre / Sankhya behavior before D6-R2, implementation-readiness closure and D9**.

A real provider/system result that contradicts accepted authority is a finding. Do not work around it, weaken the claim, invent a fallback, or mutate unrelated state merely to obtain a green result.

## 2. Execution boundary

### 2.1 What is authorized

The operator authorizes bounded real external effects required by P1/P2/P3/P5 and any P4/P6 branch that the selected fixture actually triggers.

Authorization does **not** waive:

- exact Organization / Installation / SourceInstance binding;
- current provider/business-system capability checks;
- no-blind-retry after possible acceptance;
- authoritative reread / correlation;
- customer/master-data safety;
- provider PII minimization;
- Sankhya API-Gateway-only target transport;
- explicit abort conditions below.

### 2.2 Credential rule

Credentials/tokens remain deployment/operator secrets. They must be supplied only through the already-authorized credentialed execution environment (for example protected environment/secret injection or the existing encrypted Installation/SourceInstance credential path).

Never:

- paste credentials into Git, evidence, PR comments, chat or command history;
- add a credential to this repository merely to run D8;
- log raw Authorization headers, cookies, OAuth material or Sankhya secret headers;
- revive Direct Oracle as a probe path.

The current ChatGPT execution environment used to prepare this protocol has **no Mercado Livre/Sankhya credentials and no connected provider runner**. That is an execution-environment fact, not evidence that the provider capability is absent. Real execution must occur from the credentialed operator environment.

## 3. Universal probe laws

Every live probe follows this sequence:

```text
PRE-FLIGHT READ-ONLY
→ freeze exact fixture + intended effect + abort predicates
→ snapshot only material before-state
→ ONE consequential dispatch attempt
→ classify the observed transport result no stronger than evidence allows
→ authoritative reread / correlation
→ record semantic result
→ optional restoration only through a NEW explicitly safe effect
→ authoritative reread of restoration
```

Binding rules:

1. **No blind retry.** Timeout, connection loss, process failure or unknown provider response after dispatch is `POSSIBLE_ACCEPTANCE`; do not resend. Reread/reconcile only.
2. **No fake success.** HTTP/provider success that ignores a field, returns a warning, or creates only an intermediate result is not convergence.
3. **No unrelated mutation to force success.** Never add a dummy mutable field merely to obtain provider `200` if the intended field is rejected/ignored.
4. **No destructive cleanup by assumption.** Restoration is another write and executes only if current authoritative reread proves what is being restored and its blast radius remains safe.
5. **One fixture, one declared scope.** If a provider shared resource widens the physical effect, stop unless that widened scope was measured and explicitly admitted before dispatch.
6. **No customer-impact surprise.** A production Sale/Shipment fixture must be a legitimate current business transaction. Do not fabricate a marketplace buyer transaction or invoice an already-materialized sale just to satisfy the probe.
7. **Evidence is bounded.** Record IDs/keys only where needed for correlation, redact/minimize PII, and never retain raw payloads merely because they are available.
8. **A failure may be the correct result.** `UNSUPPORTED`, `EXTERNAL_REQUIRED`, `AMBIGUOUS`, provider rejection, or a material contract contradiction is a valid D8 finding.

## 4. Common preflight

Before any write, record:

```text
D8 candidate SHA
probe ID
UTC execution time
Organization ID (opaque)
Marketplace Installation / SourceInstance ID when applicable
fixture external IDs required for correlation
current credential-binding identity evidence (redacted)
current relevant Product/D4 assumption being falsified
exact intended effect
expected authoritative reread surface
abort predicates
```

Fail closed if:

- authenticated provider/source namespace does not match the bound Installation/SourceInstance;
- selected resource belongs to another seller/source namespace;
- material current state cannot be established;
- a prior MPC/native/provider effect for the same intended fixture cannot be ruled out sufficiently;
- required owner meaning/authorization is unknown;
- write target or blast radius is ambiguous;
- the only path would require Direct Oracle or customer/master-data corruption.

---

# 5. P1 — Mercado Livre Price / Availability effect

## 5.1 Claim under test

For the selected first Mercado Livre lane, MPC can issue the accepted Offering Price meaning and Availability meaning through a currently writable provider path, preserve provider controls/blast radius, and establish the result from authoritative reread rather than transport status.

## 5.2 Fixture selection

Prefer one **existing, active, seller-owned, low-blast-radius Item** from the already-selected first lane.

Preflight must establish at least:

- exact seller/account binding;
- exact Item and current User Product/shared-resource relations relevant to the fields being changed;
- current active/paused state and current sellable quantity;
- current authoritative price using the current price-read surface, not a deprecated field by convenience;
- current pricing-automation status and current promotion/offer constraints;
- no variation/multi-origin/provider-managed-stock behavior outside the selected lane;
- original values needed for a later safe restoration.

If no fixture satisfies the lane, P1 returns a **capability/fixture finding**; do not broaden to another provider mode merely to make the write possible.

## 5.3 Current provider-contract preflight

Current Mercado Livre documentation must be rechecked immediately before execution. At protocol preparation time it establishes:

- `available_quantity` remains writable for ordinary admitted Item stock through `PUT /items/{item_id}`;
- price updates are subject to current pricing-automation/promotion rules;
- current official price documentation warns that price-only `PUT /items` requests can be rejected and that requests carrying price alongside other attributes may succeed while **ignoring the price** under the described current conditions;
- pricing automation can be inspected through the current pricing-automation Item surface.

Therefore:

> **P1 never sends an unrelated second field to bypass a price-only rejection. A rejected or ignored Price write is provider evidence and may falsify/narrow the accepted Price-control lane.**

## 5.4 Execution

Use two semantically separate effects unless the current provider contract proves they must be physically combined:

### A. Availability

1. authoritative before-read;
2. choose the smallest non-destructive quantity change that remains legitimate against real sellable inventory and cannot accidentally pause/out-of-stock the Item;
3. one write attempt;
4. authoritative reread of exact Item/provider stock surface;
5. verify shared-User-Product/member blast radius where current provider relations make it material;
6. classify `CONVERGED`, `REJECTED`, `PENDING/PROPAGATING`, `AMBIGUOUS`, `UNSUPPORTED`, or `CONTRACT_CHANGED`.

### B. Price

1. authoritative current-price read;
2. pricing-automation and promotion/control preflight;
3. if current provider authority says direct write is not admissible for the fixture, record `UNSUPPORTED/EXTERNAL_CONTROLLED` — do **not** disable provider automation to make MPC succeed;
4. otherwise choose a small legitimate price change within current provider/business constraints;
5. one write attempt;
6. authoritative price reread plus warning/effect inspection;
7. classify independently from Availability.

P1 passes the architecture claim only if the selected normal lane has a legitimate controllable effect for each meaning D4 currently claims MPC controls. A provider-level denial may require the smallest D4/Product capability reopen rather than a probe workaround.

## 5.5 Restoration

Restore an original price/quantity only when:

- the first attempt's state is known from authoritative reread;
- restoration is still legal/current;
- no intervening sale/provider automation/other actor has changed the subject;
- the restoration blast radius is still bounded.

Restoration uses a fresh explicit effect and is itself reread/reconciled. If those conditions are false, stop and surface operator Work rather than overwrite newer authority.

---

# 6. P2 / P3 / P5 — one controlled Sale → Sankhya → fiscal → label lane

## 6.1 Composition rule

P2, P3 and P5 should share **one qualifying real Sale fixture** when doing so reduces irreversible state without weakening any proof. Each probe receives its own recorded verdict even when physically executed in one progression.

Do not manufacture a Mercado Livre purchase merely to obtain a fixture. Use a legitimate current Sale that has not already been materialized/invoiced through another path.

## 6.2 Preflight before any Sankhya write

Establish from authoritative Product/provider/source evidence:

- exact Organization + Marketplace Installation + Sale native key;
- exact Selling Entity attribution;
- exact Sankhya SourceInstance binding;
- Product/quantity identities and current availability/materialization prerequisites;
- current Party Resolution state;
- current delivery destination evidence;
- absence of an already-created native Business Order for the same semantic intent;
- current e-commerce SourceInstance binding (including the accepted current TOP/series/negotiation-type facts);
- Fulfillment state needed before invoicing;
- exact Shipment/provider lane and label/artifact prerequisites;
- whether any P6 unexercised fiscal branch/component becomes material for this fixture.

Any duplicate/ambiguous native correlation aborts external creation and becomes an explicit finding/Work condition.

## 6.3 P5 — alternate destination/contact realization

P5 is falsifiable only with a Sale whose delivery destination materially differs from the reusable Partner master state.

Rules:

1. never overwrite the Partner master address merely to represent this Sale;
2. never create another Partner merely to hold another delivery address;
3. prefer the least-destructive sanctioned contact/destination representation whose current SourceInstance behavior is actually proven;
4. current Sankhya official API includes a sanctioned client-contact creation surface, but **contact existence alone does not prove order/fiscal destination propagation**;
5. snapshot the Partner/master/contact state before the probe;
6. create/reuse the bounded destination representation only when the selected Party resolution is exact and the Sale evidence is sufficient;
7. verify the order references/reconstructs the intended destination through authoritative reread;
8. after P3, verify the destination survives fiscal progression/XML/document meaning where applicable;
9. verify unrelated Partner/master state and unrelated sales were not changed.

If no safe representation is established, the correct verdict is `EXTERNAL_REQUIRED / CAPABILITY_NOT_PROVEN`, not destructive fallback.

## 6.4 Business Order materialization

Use the sanctioned Sankhya API Gateway only.

For the accepted current lane, D4 evidence names the MGECOM movement-creation/progression path. Current official Sankhya documentation continues to expose `CACSP.incluirNota` through the MGECOM Gateway for movement creation.

Execution laws:

- one semantic Business Order attempt;
- exact current binding values only; no generic provider-field bag;
- one network dispatch attempt;
- timeout/connection loss after dispatch => no second create; authoritative reread/correlation only;
- native `NUNOTA` is an external/source-qualified result, never MPC BusinessOrder identity;
- provider success does not close materialization until the authoritative point/reread state is sufficient.

## 6.5 Physical conference gate

Before P3 invoicing:

- the exact items/quantities for the controlled Sale must be physically established under accepted Fulfillment authority;
- ordinary automation cannot self-assert physical readiness;
- if the fixture is not physically ready, stop before fiscal progression.

D8 does not fake physical evidence to reach the invoice step.

## 6.6 P3 — irreversible `313 → 306` fiscal progression

Current official Sankhya documentation continues to identify `SelecaoDocumentoSP.faturar` as the MGECOM service for order invoicing and requires the source order to be confirmed plus the relevant TOP configuration.

Execution:

1. authoritative reread of the source TOP-313 order immediately before dispatch;
2. confirm all physical/fiscal/SourceInstance prerequisites are still valid;
3. one `SelecaoDocumentoSP.faturar` attempt for the exact source order under the accepted target configuration;
4. if acceptance becomes ambiguous, do not call `faturar` again; use source reread/correlation;
5. obtain the generated fiscal `NUNOTA`/result correlation from sanctioned source evidence;
6. authoritative reread of source order + result document;
7. verify distinct origin/result identities, expected TOP/result kind, item/quantity correlation and no unrelated order mutation;
8. verify fiscal/XML/document consequence required by the selected lane;
9. record `CONVERGED`, `REJECTED`, `AMBIGUOUS`, `PARTIAL`, or `CONTRACT_CHANGED`.

Because this effect is legal/irreversible, there is no artificial "restore" operation. Correction follows ordinary lawful fiscal/business procedures if the real system requires it; D8 must not invent a reversal just for test cleanup.

## 6.7 P2 — selected-lane fiscal / invoice / label progression

After the fiscal result is established, prove the provider-required selected fulfillment lane rather than assuming invoice success means dispatch readiness.

At minimum:

- reread the exact Mercado Livre Shipment with the current required format/header semantics;
- confirm payment/Shipment/provider readiness required for label acquisition;
- acquire the current label/artifact through the sanctioned provider surface for that Shipment/lane;
- bind the artifact to the exact Sale/Shipment and preserve source qualification;
- verify the artifact/readiness does not become Fulfillment business authority by itself;
- verify deadline/SLA evidence remains provider-authoritative and separate from MPC internal targets.

If the real selected lane requires a provider step not represented by accepted D4 semantics, stop and raise the smallest D4 finding.

---

# 7. P4 — native Party mutation (conditional)

P4 triggers only when the selected real Sale has **zero** sufficiently established compatible native Party matches and every identity-bearing fact required for safe creation/update is known from legitimate transaction evidence.

Do not trigger P4 when:

- exactly one compatible Party already exists — reuse under accepted compatibility rules;
- multiple plausible Parties exist — verdict is `AMBIGUOUS`, never first-result-wins;
- fiscal/customer identity evidence is incomplete/contradictory;
- mutation would be performed only to make P5 easier.

If triggered:

1. snapshot exact lookup evidence;
2. enforce duplicate/repeat protection around the same fiscal identity;
3. use only sanctioned Sankhya API/Gateway capability;
4. one create/update dispatch attempt;
5. authoritative reread by native identity/correlation;
6. verify no duplicate Party was created concurrently;
7. preserve bounded resolution lineage needed to prevent later re-guessing.

---

# 8. P6 — unexercised fiscal branch/component (conditional)

Before P3, determine whether the selected Sale materially depends on any still-unproven fiscal branch/component recorded by D4, including the current bounded Unknowns such as a distinct interstate PJ-contributor branch or fiscal components whose visibility was not established by earlier probes.

If none is material, record:

```text
NOT_TRIGGERED
+ exact fixture facts showing why
```

If one is material, **P6 becomes a blocking precondition to claiming the branch**. Prove it through the smallest sanctioned Expected-Tax/fiscal/source surface before relying on equivalence. Do not infer it from an HTTP `200`, another customer class, or realized tax from an unrelated transaction.

---

# 9. Evidence record

Each executed probe produces a concise, scrubbed record containing:

```text
probe_id
candidate_sha
executed_at
fixture identifiers needed for correlation
preflight facts / authoritative sources
intended semantic effect
exact write class (not secret-bearing raw request)
transport observation
possible_acceptance? true|false
post-write authoritative reread facts
blast-radius result
restoration result when applicable
verdict
architecture implication
reopen_owner when verdict falsifies accepted authority
```

Never record access tokens, secret headers, buyer PII, arbitrary raw provider payloads or unnecessary fiscal/customer personal data.

## 10. Verdict grammar

Use one of:

```text
PASS_CONVERGED
PASS_FAIL_CLOSED
REJECTED_DEFINITIVELY
AMBIGUOUS_RECONCILIATION_REQUIRED
UNSUPPORTED_OR_EXTERNAL_CONTROLLED
CONTRACT_CHANGED
CAPABILITY_NOT_PROVEN
NOT_TRIGGERED
```

A `PASS_FAIL_CLOSED` result is valid when the architecture explicitly requires the attempted unsafe/unsupported path to be denied.

## 11. D8 closeout rule after execution

After the authorized run:

1. update the P1–P6 ledger in D8-R1 with exact dispositions and evidence references;
2. classify every contradiction against the smallest owning authority;
3. if no accepted assumption changes, re-run the repository gate and close D8 through operator ratification;
4. if a material D0–D7 assumption changes, reopen only the smallest owner, repair it, then boundedly revalidate the affected D8 flow before D8 close;
5. only after D8 closes may D6-R2 begin.

No probe result by itself authorizes D9 or Product implementation.
