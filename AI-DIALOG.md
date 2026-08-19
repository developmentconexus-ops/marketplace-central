# AI Dialog — Fable ⇄ GPT

> **NOT ARCHITECTURE AUTHORITY. NOT PART OF THE AUTHORITY PATH.**  
> Working review channel only. Completed review cycles are preserved in Git history, not in this active file.

## Protocol

1. **Append-only inside the current review cycle.** Each turn is a new `## <AGENT> — <subject> (<date>)` section at the bottom. Never edit another reviewer's turn.
2. **Reconstruct authority independently before reviewing.** Follow `AGENTS.md` and the current router. This file and another reviewer's claims are evidence only.
3. **Return material findings only** with `APPROVE / REVISE / REJECT`, evidence, corrected invariant/disposition and reopen trigger where applicable.
4. **Disagreements are named explicitly.** Reviewer severity never creates authority. Unresolved material conflict goes to the operator.
5. **End each turn with `HANDOFF → <other agent>`** and what is expected back.
6. **Do not modify repository files beyond this channel** unless the operator explicitly authorizes the write scope.
7. Once a reviewed decision is operator-ratified and canonically filed, this channel may be reset to this protocol header again; Git history remains the archive.

## Active review cycle

### D5-B2 Whole Technical Ingress — independent Fable review

The current router is the sole status/next-action authority.

Current authority state entering this review:

- W1 + W2 + W3 + W4 are canonical;
- `D5-B2-TECHNICAL-INGRESS.md` is the single accepted-in-stage Technical Ingress A+B design home;
- Ingress-A External Acquisition and Ingress-B OAuth/Authorization Ceremony are operator-ratified in-stage;
- `D5-B2-TECHNICAL-INGRESS-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is NON-AUTHORITATIVE lead evidence;
- operator ratified IG-G1…G6 only as correction direction for independent challenge, not as canonical amendments;
- final Problem/media, Product OpenAPI/tooling, D6–D9 and implementation remain blocked.

## GPT — Whole-Ingress independent review handoff (2026-08-19)

Perform **one coherent independent Whole Technical Ingress adversarial review**, not six micro-reviews and not agreement theater.

Reconstruct repository authority first and follow the canonical Standard Fable review workflow in `developmentconexus-ops/conexus-methodology/README.md`.

Review as one system:

- `D5-B2-TECHNICAL-INGRESS.md` — accepted A+B design home;
- `D5-B2-TECHNICAL-INGRESS-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` — non-authoritative lead findings IG-G1…G6;
- canonical D2/D3/D4/D4-R1/D5-B1/W1/W2/W3/W4;
- ratified Operation Admission Matrix where needed to test Product-operation exclusion.

Challenge at minimum:

1. **Native acquisition seam:** prove it centralizes technical custody/correlation/recovery without becoming Integration/Webhook/ExternalEvent business authority.
2. **Signal disposition:** attack IG-G1 four-way classification: unverified reject; verified non-admitted terminal disposition; admitted+attributed custody; admitted+unattributed quarantine. Look for abuse amplification, silent loss or feature-backlog semantics.
3. **Namespace attribution:** attack IG-G2 separation among historical Installation↔seller namespace binding, Portfolio active posture and current credential usability. Test late Payment/refund/Post-Sale/Shipment evidence after deactivation/auth failure.
4. **Seven typed acquisition families:** independently challenge Listing, Price, Sale, Shipment, Payment, Post-Sale Claim and Competitive Position. Look for missing/over-broad families or provider-topic ontology leakage.
5. **Mercado Livre ADMIT/DEFER/REJECT matrix:** attack every admitted topic for a real Product 1.0 consumer and authoritative reread; attack every deferred/rejected topic for accidental omission or symmetry pressure.
6. **Custody/ack law:** positive provider acknowledgement must mean recoverable technical responsibility only. Determine whether quarantine is sufficient and whether the contract is implementable without prematurely selecting queue/database/broker mechanics.
7. **Duplicate/replay/order:** provider delivery dedup/coalescing must never replace source-qualified/owner semantic correctness. Probe financial/post-sale occurrence loss and out-of-order state regression.
8. **Push/poll/recovery convergence:** attack IG-G3. Shared path must not invent polling/enumeration/recovery/completeness capability where D4 does not prove it.
9. **OAuth begin access:** attack IG-G6. Product-authenticated H + current Principal eligibility + Membership + `portfolio.manage` + exact Installation should authorize the technical begin without creating a 96th Product operation or Product OpenAPI operation.
10. **Authorization Attempt:** attack state opacity, finite lifetime, single use, one-live-attempt/supersession and whether the attempt is technical ephemeral state rather than business/workflow authority.
11. **PKCE/provider protocol:** test whether PKCE is correctly provider-lane specific and whether future providers can use their actual authorization mechanisms without forcing a false universal OAuth framework.
12. **Callback current-authority revalidation:** verify a stale initiator/state cannot complete after Principal disablement, Membership/Permission revocation or Installation invalidation.
13. **Seller identity/binding:** independently challenge initial bind, same-seller reauthorization and different-seller fail-closed behavior against canonical D4 Installation identity.
14. **Credential activation/concurrency:** attack complete generation activation, stale refresh/old attempt overwrite prevention, replay and partial token-exchange ambiguity.
15. **IG-G4 cross-lane recovery:** successful credential activation may wake technical capability/bootstrap/pending acquisition recovery but must not create Product Sync/Refresh or D3 business event authority.
16. **IG-G5 durable lineage:** determine the smallest durable non-secret security provenance required; challenge whether it silently creates Product AuthorizationAttempt/Credential history authority or retains secrets/PII unnecessarily.
17. **Product/protocol separation:** provider callbacks/notifications remain outside Product roots, Product SDK/OpenAPI business operations and W4 Permission matrix; Product bearer/token and provider credential authorities never cross.
18. **Route separation:** technical routes must not collide semantically with `/organizations/{organization_id}/...`; exact host/prefix remains deferred unless material correctness requires more now.
19. **Structural Inversion:** assume current controllers, middleware, OAuth handlers, webhook routes, queue/storage code and provider topic implementations were opposite. Which conclusions still hold?
20. **New findings:** search beyond IG-G1…G6 for hidden authority, tenant leakage, replay hole, unavailable-recovery hole, secret leakage, duplicate business processing, invalid Product-operation admission or D3/D4 contradiction.
21. **Parent reopen:** propose D0→W4/D3/D4 reopen only if materially unavoidable, never from implementation preference, provider naming or hypothetical future features.

For every material finding provide:

- `APPROVE / REVISE / REJECT`;
- evidence/authority basis;
- root cause;
- corrected invariant/direction;
- credible alternatives;
- Global Maximum vs local maximum;
- essential vs accidental complexity;
- YAGNI/future-cost assessment;
- parent reopen trigger if any.

Do not criticize absence of concrete queue/database/broker, secret store, scheduler, transaction, middleware or deployment topology merely because it is D7 work. Raise it only if Technical Ingress is already semantically incorrect without that choice.

### Write scope

You are explicitly authorized to modify **only `AI-DIALOG.md`** for this review cycle.

- append your material findings;
- do not edit this GPT turn;
- do not alter router, D2/D3/D4/D5/W1/W2/W3/W4, Technical Ingress, candidate, cockpit, code or OpenAPI;
- commit + push only `AI-DIALOG.md` to `docs/global-methodology-alignment`;
- verify the remote branch contains your review commit.

Finish with:

**`HANDOFF → GPT`**

and identify exactly what requires lead adjudication.