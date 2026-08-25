# D6-R2 P8 — B10 Correspondence Region Revalidation

> **Status:** CANDIDATE / OPERATOR WALKTHROUGH REQUIRED — no re-LOCK claimed
> **Block:** B10 — Preparação / R10, correspondence region only
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Candidate evidence:** `qualification/d6-r2-wireframes/b10-preparation.html`
> **Wire prerequisite:** PR #70 ACCEPTED / INTEGRATED / PROOF PASS
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Bounded decision

**Outcome: `CURRENT STRUCTURE CONFIRMED` — repair only the correspondence region.**

The accepted B10 hierarchy, marketplace-requirement table and ListingIntent handoff remain valid. PR #70 repaired the missing read projection that previously forced a human to act on an opaque `candidate_key`.

Target invariant:

```text
human recognition / selection
  = current owner-supplied display_label adjacent to candidate_key

correspondence write
  = candidate_key + correspondence_etag only

presentation
  != identity
  != write authority
```

No new Product operation, path, Permission, Principal kind, owner, state machine or frontend business authority is required.

## 2. Authority consumed

`GetProductChannelReadiness` now carries:

```text
subject_presentation
correspondence
correspondence_candidate_population
  known(candidates[])
  unknown
  unavailable
correspondence_etag
```

Each known candidate is exactly:

```text
candidate_key
display_label
```

`ResolveProductChannelCorrespondence` continues to accept only the canonical `candidate_key` with the exact subject and current `correspondence_etag`. The candidate label is read presentation and is never submitted.

## 3. Candidate interaction

The bounded candidate renders:

- resolved correspondence with current human label and explicit change/remove actions;
- unresolved known-populated candidates as a native radio group;
- conflicting correspondence through the same explicit candidate choice, with conflict explanation;
- known-empty population as a completed read with zero candidates;
- unknown population as not yet determinable;
- unavailable population as a failed consultation, never as known-empty;
- no preselected candidate;
- resolve disabled until explicit selection;
- consequential resolve/clear followed by a mandatory authoritative reread before continuation;
- `candidate_key` as the only simulated write carrier.

Continuation to ListingIntent remains enabled only after current correspondence truth is resolved. B10 still performs navigation only and does not call `CreateListingIntentDraft`.

## 4. Protected structure

This candidate does not reopen:

- `OFERTA > Preparação` placement;
- exact Organization + Marketplace Installation context;
- bounded multi-source product search;
- selected source-qualified subject;
- the four-column marketplace-requirement projection;
- downstream ListingIntent boundary;
- technical-detail secondary disclosure;
- B00 shell/navigation meaning;
- B20 or any unopened P8 block.

## 5. Proof

Browser-operable checks on the local deterministic artifact demonstrated:

- selection enables resolve only after one radio choice;
- the simulated effect carries `candidate_key` and no display label;
- continuation remains blocked until reread;
- reread returns to resolved current presentation;
- known-empty, unknown, unavailable and conflicting states remain distinct;
- the candidate list fits a 390px viewport without horizontal document overflow;
- the native mobile navigation law remains active.

`scripts/verify-d6-r-b10-correspondence-wireframe.mjs` protects the bounded structure with eight deterministic negative controls and is wired into `npm run gate` whenever this surface or its verifier changes.

## 6. P8 operator gate

Current disposition:

```text
CANDIDATE
operator walkthrough required
LOCK / REVISE / UPSTREAM FINDING
```

The assistant, verifier and reviewer cannot set `LOCKED`. P9 must not rerun until the operator explicitly re-LOCKs this affected region.

## 7. Reopen triggers

Reopen only if operation of this candidate proves that:

- `display_label` is insufficient to distinguish real correspondence candidates;
- a resolved candidate cannot be presented from current Readiness truth;
- pagination/search is materially required for a real candidate population;
- continuation safety cannot be decided from current correspondence/readiness truth;
- a write would require presentation rather than canonical identity;
- responsive or accessible operation cannot preserve the candidate-selection job.
