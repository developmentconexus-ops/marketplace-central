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

### D5-B2 Final Problem / Media Consistency — independent Fable review

The current router is the sole status/next-action authority.

Current authority state entering this review:

- D0→D4 + D4-R1 + D5-B1 are accepted/canonical;
- Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress are accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- `D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md` is NON-AUTHORITATIVE lead evidence;
- canonical W1/W2/W3/W4/Technical Ingress are unchanged by that candidate;
- Product OpenAPI/tooling, D6–D9 and implementation remain blocked.

## GPT — Final Problem / Media independent review handoff (2026-08-19)

Perform **one coherent independent adversarial review** of the final D5-B2 Problem/media package. Do not split PM-C1…PM-C8 into ceremonial micro-reviews and do not treat the lead candidate as authority.

Reconstruct repository authority from scratch in the router's order and follow the canonical Standard Fable review workflow from `developmentconexus-ops/conexus-methodology/README.md`.

Review as one system:

- `D5-B2-FINAL-PROBLEM-MEDIA-CONSISTENCY-REVIEW-CANDIDATE.md` — non-authoritative lead evidence;
- `D4-R1-PUBLICATION-INPUT.md`;
- `D5-API.md`;
- `D5-B2-PRODUCT-OPERATION-SURFACE.md`;
- `D5-B2-OPERATION-ADMISSION-MATRIX.md`;
- canonical W1/W2/W3/W4;
- canonical `D5-B2-TECHNICAL-INGRESS.md`;
- accepted D1/D2/D3/D4 meaning where needed;
- RFC 9110 and RFC 9457 primary text where the HTTP/Problem claim depends on them.

Challenge at minimum:

1. **Standard HTTP failures:** Is RFC 9457 `about:blank` the Global Maximum for status-only 405/413/415 failures, or does an MPC-specific consumer require a custom type? Test compatibility with W2's absolute stable type law and client branching.
2. **Media negative map:** Attack every proposed 400/413/415/422/409/500 classification. Look for syntactic/semantic misclassification, framework dependence, or a missing reachable failure.
3. **Problem identifier:** Determine whether `type` alone remains sufficient or whether any current Product consumer genuinely requires a global `code`. Reject duplication by habit, but do not hide a real client need.
4. **Product/protocol separation:** Ensure technical provider/OAuth routes may reuse a standard representation format without inheriting W2 Product types, Product SDK/OpenAPI membership or business meaning.
5. **Media identity:** Attack the claim that one `listing_intent_media_id` immutably denotes one accepted upload within exactly one ListingIntent. Look for update/replacement semantics already required by authority.
6. **Selection vs deletion:** Test whether unselecting media must preserve historical attempt evidence and whether the proposed D7 garbage-collection fence is sufficient without inventing Product delete/retention authority.
7. **Binary delivery seam:** Challenge the lead's claim that durable authored-media authoring requires an authorized inspectability/access reference while no standalone Product media GET is currently admitted. Look for:
   - an unproven D6 consumer;
   - hidden new Product operation;
   - bearer-locator leakage;
   - storage topology becoming Product meaning;
   - inability to implement safely without prematurely choosing D7;
   - a smaller alternative that remains usable across reads.
8. **Source vs authored media:** Verify that shared presentation primitives cannot collapse source-qualified candidate evidence and ListingIntent-authored media into one generic Media/Asset owner.
9. **Idempotency + revision:** Probe lost-response replay, binary fingerprinting, changed ETag under same key, concurrent upload and whether successful creation advances ListingIntent revision without breaking exact replay ordering.
10. **Untrusted binary boundary:** Challenge size/content-type/content-inspection, filename/path, active content, scanning, logging and PII assumptions. Raise only D5-semantic gaps; do not choose D7 products/topology.
11. **Problem leakage:** Test object-store, CDN, scanner, transformer and provider failures for accidental Product taxonomy or raw detail leakage.
12. **OpenAPI/SDK:** Ensure only the admitted Product media capability appears in Product OpenAPI/SDK and no technical callback/OAuth/provider errors enter it.
13. **Status convergence:** Determine whether removing stale “next work” text from active B2 artifacts is required authority convergence or unnecessary editorial churn. Router remains sole status authority either way.
14. **Operation/Permission count:** Prove the package adds no hidden 96th operation, 30th Permission, media GET, delete/update operation or technical ingress operation.
15. **Structural Inversion:** Assume the current code uses direct object-store URLs, generic asset tables, provider error passthrough, framework default 400s and separate media handlers. Which conclusions still hold?
16. **New findings:** Search beyond PM-C1…PM-C8 for contradictory schema, hidden authority, security hole, unreachable Product flow, unbounded retention or false HTTP semantics.
17. **Parent reopen:** Propose D0→D5-B1/D4-R1 reopen only if materially unavoidable; never from implementation preference, aesthetic purity or hypothetical future clients.

For every material finding provide:

- `APPROVE / REVISE / REJECT`;
- evidence/authority basis;
- root cause;
- corrected invariant/direction;
- credible alternatives;
- Global Maximum versus local maximum;
- essential versus accidental complexity;
- YAGNI/future-cost assessment;
- parent reopen trigger if any.

Do not criticize the absence of a concrete blob store, CDN, scanner, signed-URL scheme, image processor, transaction or deployment topology merely because it is D7 work. Raise it only if the current D5 contract cannot remain correct without a semantic decision now.

### Write scope

You are explicitly authorized to modify **only `AI-DIALOG.md`** for this review cycle.

- append your independent material findings;
- do not edit this GPT turn;
- do not alter router, candidate, D4-R1, D5, Operation Matrix, W1/W2/W3/W4, Technical Ingress, cockpit, code or OpenAPI;
- commit + push only `AI-DIALOG.md` to `docs/global-methodology-alignment`;
- verify the remote branch contains your review commit.

Finish with:

**`HANDOFF → GPT`**

and identify exactly what requires lead adjudication or operator decision.
