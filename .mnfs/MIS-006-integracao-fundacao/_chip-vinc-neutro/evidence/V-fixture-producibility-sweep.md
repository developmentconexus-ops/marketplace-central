# EXHAUSTIVE FIXTURE SWEEP — "fixture not producible by the wire"

Hub precondition for round 5. The class has recurred: round 3 claimed it closed, round 4 found two
more — one of them in the very file round 3 edited. That is occurrence 2, and a point-fix per
occurrence is what let it recur. The `Record`-indexed-by-union class was discharged by a
grep-by-STRING sweep that found a fourth site no reviewer named; this class had no equivalent.

Scope: **every** candidate fixture in the chip's write-set. 28 fixtures across the two files that
have any; the other four test files under `pages/vinculos/` construct no candidates.

## The reading error being swept for, named

The previous round's mistake was **reading the generator at the scoring switch and stopping there**.
Scoring (`buildConcordantCandidate`, `applySingleAnchorScore`, `applyAmbiguousCorroborationScore`,
`applyUnresolvedScore`) sets confidence/band/status and SEEDS reasons. But every one of those paths
then returns through **`appendProviderDeclaredUnavailableReasons`** (generation_service.go:648-698),
which walks the provider's DECLARED anchors and appends one reason for each anchor with no
FOR/AGAINST signal. A fixture verified against the switch alone can have a correct
confidence/band/status and an impossible reasons array — which is exactly what shipped.

So this sweep checks every fixture against the FINALIZING step, and the invariant it enforces is:

> **`marca` is in `KnownIdentityAnchors()` (4 members, marketplace_capability.go:40-45), the adapter
> declares all four for every provider (identity_anchor_adapter.go:28-35), `marca` never carries a
> FOR/AGAINST signal, and `classifyProviderIdentityAnchor` emits for it either way — UNAVAILABLE
> when unsupplied (:700-702), INCOMPARABLE when supplied (:706-708, because `marca` has no case in
> `identityAnchorValues`). Therefore EVERY candidate of EVERY provider carries a `marca` reason.**
>
> For `mercado_livre` — the only provider with a capability adapter in this tree, declaring
> `{seller_sku, ean, title}` (mercado_livre/capability_adapter.go:90) — that reason is **UNAVAILABLE**.

## Classification rule (stated before the verdicts, so it cannot be fitted to them)

- **A — CLAIMS PRODUCIBILITY.** The fixture asserts "this is a row the backend emits". Every field
  must be exactly producible. This is the whole golden file: a design gate whose example cannot
  occur gates a fiction.
- **B — DELIBERATELY NON-PRODUCIBLE.** Non-producibility IS the assertion (wire-drift probes). A
  producible fixture here would prove nothing. Requirement: the test must SAY so.
- **C — PRESENTATION ISOLATION.** A minimal fixture exercising one rendering rule. Defective only
  when the test's CONCLUSION depends on a shape the backend cannot produce. Omitting fields the
  test does not read is not a defect; asserting a LAYOUT off an impossible reason set is.

## VERDICTS — VinculosDesign.golden.test.tsx (class A, 4 fixtures + the shared `base()`)

| Fixture | Verdict |
|---|---|
| `base()` L42 — CONFIRM 60/MEDIA `exact_ean`/`ean` | **PRODUCIBLE.** `applySingleAnchorScore` case `ExactEAN` (:539-549): 60/MEDIA/CONFIRM, `ean` FOR "ean corrobora codprod (unproved)", `seller_sku` via `missingMatchedAnchorReason` → INCOMPARABLE/provider "sem CODPROD para corroborar o EAN". Finalizer appends `marca` UNAVAILABLE. All three present. Fixed round 4. |
| L123 `cand_1` EXEMPLO-IO | **PRODUCIBLE** — inherits `base()` whole, overriding only ids and product name. |
| L180 `cand_nc` NO_CANDIDATE | **PRODUCIBLE.** `applyUnresolvedScore` (:620-628) via `LinkCandidateStateUnresolved` (:215/:379): 0/BAIXA/NO_CANDIDATE, `state: "unresolved"`, two INCOMPARABLE/erp absence reasons (product is nil ⇒ `missingMatchedAnchorReason` takes the `product == nil` branch, :636), plus `marca`. Fixed round 4. |
| L238 `cand_1` ACCEPT 95/ALTA | **PRODUCIBLE.** `buildConcordantCandidate` (:493-507): `exact_sku`/`seller_sku`, 95/ALTA/ACCEPT, both FOR reasons verbatim, plus `marca`. `title` correctly absent: product non-nil and both title values non-empty ⇒ `emit=false` (:709-711). Fixed round 4. |
| L256 `cand_nc` (sweep row) | **PRODUCIBLE** — same shape as L180. Fixed round 4. |

**Class A is clean.** Every fixture that claims to be a backend row now carries the finalizer's
output, and the golden asserts the layout that follows from it (two chips + a "+1" toggle, not two
chips alone).

## VERDICTS — QueueTab.test.tsx

### Class B — deliberately non-producible (3). All correctly labelled.

| Fixture | Verdict |
|---|---|
| L659 `cand_drift` `match_status: "PENDING_REVIEW"` | **CORRECT.** The cast is the point and the test says so. |
| L695 `cand_dir` `direction: "PARTIAL"` | **CORRECT.** Same. |
| L737 `cand_band` `confidence_band: "CRITICA"` | **CORRECT.** Same. |

### Class C — presentation isolation. TWO FINDINGS, both mine, neither reported by any reviewer.

| Fixture | Verdict |
|---|---|
| **L192 `cand_inc` — all-INCOMPARABLE, no `marca`** | **FINDING #1 — NOT PRODUCIBLE for `mercado_livre`, and it backs the DECIDING criterion (V2).** ML always emits `marca` UNAVAILABLE, so a 100%-INCOMPARABLE row cannot occur for it. It IS producible for a provider that declares `marca` supplied (then `marca` classifies INCOMPARABLE) — a capability declaration no adapter in this tree has. |
| **L330 `cand_noside` — `marca` INCOMPARABLE with no side** | **FINDING #2 — same class.** For ML, `marca` arrives UNAVAILABLE, not INCOMPARABLE. The side-less INCOMPARABLE path (:706-708) is real and reachable, but only under a declaration that supplies `marca`. |
| L65/L74/L82 `cand_1..3` bands | CLEAN (C). Asserts band pills and one chip each; conclusion does not depend on the reason set being complete. |
| L127 `cand_confirm` 4 reasons incl. `marca` | CLEAN (C). Exercises the +N collapse; carries `marca`. |
| L251 `cand_unresolved` | CLEAN (A-grade, added by this sweep — see below). |
| L295 `cand_mixed` 3 reasons | CLEAN (C). Ranking rule; mixed directions, no impossible claim. |
| L353 `cand_nc` `reasons: []` | CLEAN (C). Asserts the NO_CANDIDATE affordance (sem candidato / Criar produto / Ignorar); reads no reason. Empty is a minimal isolation, not a claimed wire row. |
| L382, L413, L427, L434 provider rows | CLEAN (C). Assert the CANAL cell only. |
| L467, L502, L527, L528, L533 | CLEAN (C). Approve/deep-link/bulk flows; assert no motivo layout. |
| L602 `cand_accept`, L612 `cand_confirm_sku`, L621 `cand_title` | CLEAN (C). Assert the Identificado por cell, which is a pure function of `match_status`+`match_input` — a field pair each fixture sets correctly. The reason arrays are one-element isolations and the conclusion does not read them. |

## WHAT WAS DONE ABOUT FINDINGS #1 AND #2

Not deleted, and not "fixed" by bolting `marca` onto them — that would destroy what each one tests.
An all-INCOMPARABLE row is the exact boundary V2 is about; a side-less INCOMPARABLE is the exact
case V4's honesty rule is about. Both remain valid under a capability declaration, which is a real
future provider, not a fiction.

**What was wrong was the silence.** Both fixtures read as "this is what the wire sends" and neither
said which declaration it assumed. Both now carry that statement.

**And V2 — the deciding criterion — got a producible proof, which it did not have.** New test:
`promotes the actionable absences over the permanent one, on a reason set the backend really emits`
(QueueTab.test.tsx:251). Its fixture is verbatim `applyUnresolvedScore` output for `mercado_livre`:
two INCOMPARABLE/erp absences plus `marca` UNAVAILABLE — the single most common row on this screen.

It is a must-fail, run against the pre-fix string-literal enumeration restored into `QueueRow.tsx`:

    const byDirection = (d: string) => reasons.filter((r) => r.direction === d);
    const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")]
      .slice(0, COMPACT_CHIP_LIMIT);

    ×  QueueTab > promotes the actionable absences over the permanent one, on a reason set the
       backend really emits
       → expected [ <span …(4)></span> ] to have a length of 2 but got 1
    Tests  1 failed | 15 skipped (16)

**One chip.** And the one that survived is `– Marca` — "provider não fornece a âncora marca", the
single motivo nothing can ever be done about — while BOTH actionable absences sat behind "+2".

That is the V2 defect in the form the operator actually meets it. The original V2 fixture proves the
extreme (zero chips, empty cell); this proves the common case (the only useless motivo, promoted).
The extreme needs a capability declaration; this needs nothing but today's backend.

The probe was reverted from a saved copy; `git diff -- QueueRow.tsx` came back empty against HEAD
and `npx --no-install vitest run src/pages/vinculos/` is **34/34 green**.

## RESIDUAL RISK, STATED

This sweep is exhaustive over the chip's **write-set**. It does not cover fixtures in files this
chip does not own, and it verifies against the generator as it stands at `bcab8269`+`main` — a
backend change to the declared-anchor list or to `mercado_livre`'s capability would invalidate class
A without touching any file here. That coupling is real and undetectable from the FE side; it is
filed as FINDING F-07 in EVIDENCE.md rather than left implied.

## CLASS STATUS

Two findings, both in class C, both found by this sweep and not by four reviewer rounds. The class
is **not** declared closed by assertion — it is declared swept, with the sweep's own boundary stated
above. If round 5 finds a third instance INSIDE the write-set, the hub's rule applies: not a
point-fix, but a named mechanism plus exhaustive re-sweep.
