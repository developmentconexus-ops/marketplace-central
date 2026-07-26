# CHIP-ANCHORS — evidence pack

```yaml
chip: CHIP-ANCHORS
branch: chip/anchors
base_sha: 917f7bb58e385847fba5612201823f9db48791c6
level: QA-0
status: in-progress
```

BASE-SHA: 917f7bb58e385847fba5612201823f9db48791c6
CONTRATO: .mnfs/MIS-006-integracao-fundacao/_chip-anchors/validation-contract.md
EXEMPLO-IO: MLB4735326915 · "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO" vs ERP product written in cm

## Note on the pack location

The chip pack (`chip.md` + `validation-contract.md`) was authored by the hub on `main` at
`f81b8975`, i.e. ONE COMMIT AFTER the accepted BASE-SHA `917f7bb5` this branch is based on. The
chip did NOT rebase (the prompt binds the 40-hex BASE-SHA as the governance BaseSha). Both pack
files were read from `f81b8975` via `git show`, byte-verified against that commit. They arrive on
this branch through the hub's merge, not through the chip's diff.

## Dispatch ledger

Rows are written AT DISPATCH TIME (core §1 Claude-side/OS-process ledger rule), completed with
the verdict/output artifact afterwards.

| # | Phase | Role | Model / effort | Path | Prompt artifact | Output artifact | Status |
|---|-------|------|----------------|------|-----------------|-----------------|--------|
| D1 | P0 | Investigator (module map) | Explore subagent (sonnet) | Agent tool, sync | inline (recorded §Dispatch prompts) | file:line map, consumed into this pack | completed |
| D2 | P2 | Feature planner | gpt-5.6-sol / medium | OS-process codex | `scratchpad/prompt-p2-planner.md` | `scratchpad/agent__p2-planner.last.md` + `.log` | completed — plan accepted with 4 chip adjudications (below) |
| D3 | P3 | Implementer S1 (F-01) | gpt-5.6-luna / high | OS-process codex | `scratchpad/prompt-s1.md` (sha256 `cbb6bcff…bb80`) | `scratchpad/agent__s1-impl.last.md` + `.log` | GREEN, committed `5bc55219`, chip-verified independently |
| D4 | P3 | Implementer S2 (F-01) | gpt-5.6-luna / high | OS-process codex | `scratchpad/prompt-s2.md` | `scratchpad/agent__s2-impl.last.md` + `.log` | GREEN, committed `b9da6d2`, chip-verified independently |
| D5 | P3 | Implementer S3 (F-02) | gpt-5.6-sol / low (complex slice) | OS-process codex | `scratchpad/prompt-s3.md` | `scratchpad/agent__s3-impl.last.md` + `.log` | code GREEN; worker reported BLOCKED on a defective card instruction (below); chip verified and committed `11e68f6f` |
| D6 | P3 | Implementer S4 (F-03) | gpt-5.6-luna / high | OS-process codex | `scratchpad/prompt-s4.md` | `scratchpad/agent__s4-impl.last.md` + `.log` | GREEN, committed `030fa58c`, chip-verified independently |
| D7 | P4 | Adversarial reviewer, feature F-01 | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/ac5c97165db9fcd08.output` | completed — **PASS-WITH-FINDINGS** (2 SHOULD-FIX, 1 NIT) |
| D8 | P4 | Adversarial reviewer, feature F-02 | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/a051c9ecbaf2f1f10.output` | completed — **FAIL** (1 BLOCKING, 2 SHOULD-FIX) |
| D9 | P4 | Adversarial reviewer, feature F-03 | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/a4ef477310cfa1b64.output` | completed — **PASS-WITH-FINDINGS** (1 SHOULD-FIX) |
| D10 | P3 | Implementer S5 (corrective: F-02 blocking + F-01 layering) | gpt-5.6-sol / low (complex slice) | OS-process codex | `scratchpad/prompt-s5.md` | `scratchpad/agent__s5-impl.last.md` + `.log` | GREEN, committed `f92ca9c7`, chip-verified independently |
| D11 | P3 | Implementer S6 (F-03 test hardening) | gpt-5.6-luna / high | OS-process codex | `scratchpad/prompt-s6.md` | `scratchpad/agent__s6-impl.last.md` + `.log` | GREEN, committed `1df627f8`, chip-verified independently |
| D12 | P4 | Re-reviewer, feature F-02 over `f92ca9c7` | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/af63a82825465373b.output` | completed — **PASS-WITH-FINDINGS** (2 SHOULD-FIX both pre-existing, 1 NIT); blocking defect confirmed closed |
| D14 | P6 | Gate reviewer, GPT side of the dual gate | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol.md` | `scratchpad/agent__p6-sol.last.md` + `.log` | dispatched |
| D15 | P6 | Gate reviewer, COLD Opus side | Opus, `harness:gate-reviewer` (no Edit/Write/Bash by construction) | Agent tool, async | inline brief | `tasks/a93b3ffaca9b176b7.output` | dispatched |
| D13 | P4 | Re-reviewer, feature F-01 over `f92ca9c7` | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/ad9c0421c9c2a538c.output` | completed — **PASS**, zero findings |

## Slice S3 — chip verification, and a defect in the chip's own slice card

The worker returned **BLOCKED, uncommitted**, with the implementation and all required tests
complete and green. The block was caused by the CHIP's card, not by the code: the card listed

```
gofmt -l internal/modules/product_links/application
```

with "expected EMPTY". That directory contains NINE pre-existing files whose only gofmt diff is
`^M` — CRLF from the Windows autocrlf checkout, zero token changes. Verified by the chip with
`gofmt -d ... | cat -A`, which shows `^M$` on every line and no reordering. This is precisely the
false alarm the profile already records as field finding F-ENV-M01 (§3, 2026-07-18: "scope gofmt
gates to milestone-authored dirs"). The check could never come back empty, so the worker was
correct to stop rather than reformat nine files outside its write_set — it obeyed the write_set
over the command, which is the right precedence.

Chip actions: (a) verified the work independently, (b) committed it as `11e68f6f` with the cause
recorded in the commit message, (c) corrected the S4 card to scope gofmt to the two authored files
and to carry an explicit CRLF warning — `resolution_service.go` and `resolution_service_test.go`
are BOTH in the CRLF set, so a `gofmt -w` there would have turned a ~40-line change into a
~900-line diff.

| Check | Command | Result |
|---|---|---|
| Write-set exactness | `git status --porcelain` before commit | exactly the 3 declared files |
| L0 build | `go build ./...` | exit 0 |
| L0 vet | `go vet ./...` (full tree) | exit 0 |
| L1 tests | `go test ./internal/modules/product_links/...` | all `ok`, exit 0 |
| L0 format, authored files only | `gofmt -l <the 3 files>` | empty |

Chip ruling A1 verified applied: `hardNegativeDimension` now returns a comparison key AND a
display string (`generation_service.go:726-770`), and `detectHardNegative:698-700` formats the
operator message from the DISPLAY values, so the queue keeps reading `50cm≠500mm` rather than a
canonical `mm:...` form.

## C12 — proven under the hub's revised proof form (ruling of 2026-07-26, contract @ `c9fecbc`)

`gofmt -w internal/composition/root.go` applied and committed as `633bf9fa`, per the hub ruling.

| Proof | Command | Result |
|---|---|---|
| additive-only in TOKENS | `git diff -w 917f7bb5..HEAD -- .../root.go \| grep '^-'` | **no removed lines** |
| file is gofmt-clean | `gofmt -l internal/composition/root.go` | empty |
| region diff quoted raw | `git diff --unified=2 917f7bb5..HEAD -- .../root.go` | quoted below, realignment visible |

```
+	productlinksconnectors "…/product_links/adapters/connectors"                                   (import block)
+	productLinkIdentityAnchorReader := productlinksconnectors.NewIdentityAnchorAdapter(marketplaceCapabilities)
 	productLinkGenerationSvc := productlinksapp.NewGenerationService(productlinksapp.GenerationServiceConfig{
-		Snapshots:    productLinkSnapshotRepo,
-		Matcher:      productMatcher,
-		Store:        productLinkCandidateRepo,
-		AutoApprover: productLinkResolutionSvc,
+		Snapshots:       productLinkSnapshotRepo,
+		Matcher:         productMatcher,
+		Store:           productLinkCandidateRepo,
+		IdentityAnchors: productLinkIdentityAnchorReader,
+		AutoApprover:    productLinkResolutionSvc,
 	})
```

The four `-`/`+` pairs are the gofmt column widening forced by the added key. They are shown, not
hidden: `git diff -w` collapses them to nothing, which is the token-level test the criterion now
uses.

## Slice S2 — chip verification (not the worker's word)

| Check | Command | Result |
|---|---|---|
| Write-set exactness | `git diff --name-only 917f7bb5..HEAD` | 13 files, all declared across S1+S2, zero undeclared |
| L0 build | `go build ./...` | exit 0 |
| L0 vet | `go vet ./...` (full tree) | exit 0 |
| L1 tests | `go test ./internal/modules/product_links/... ./internal/modules/connectors/... ./internal/composition` | all `ok`, exit 0 |
| C2 symbol | `grep -rn mandatoryUnavailableReasons .../application` | **0 hits** |
| C2 anchor literals | `grep -rn '"marca"\|"refforn"' .../application` (production only) | **0 hits** |
| C3 provider branch | `grep -rn 'mercado_livre\|ProviderCode ==' product_links` (production only) | **0 hits** |
| C12 additive-only | `git diff --unified=0 917f7bb5..HEAD -- root.go \| grep '^-'` | **no removed lines** |

The worker again reported plain `go build ./...` blocked by VCS stamping; again it does not
reproduce chip-side (exit 0). Profile §3 (2026-07-17, CHIP-M03 F-ENV) already ratifies that the
chip's re-run outside the codex sandbox is the verification of record for exactly this class of
sandbox false alarm.

Substance spot-checks on the diff:

- All six scoring branches route through one finalizer,
  `appendProviderDeclaredUnavailableReasons` (`generation_service.go:624-635`), called at
  `:505, :561, :575, :587, :602, :621`. `mandatoryUnavailableReasons()` is gone.
- The provider-level detail is `fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor)`
  (`:633`) — R1 honoured: the distinction lives in `detail`, the direction enum is untouched.
- The failure names the provider: `PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE for
  provider %q: %w` (`:181`) — chip ruling A2 applied.
- Resolution happens once per distinct trimmed provider code (`:149-178`) before any scoring or
  persistence, and every failure mode — empty code, adapter error, nil declaration, unknown
  anchor — returns the error rather than an empty declaration (`:161, :165, :168, :172`). That is
  the ADR-17 point: an unresolved declaration must not decay into "the provider supplies none".
- `:76` keeps the pre-existing `PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED` behaviour when the
  reader is absent.

### RESOLVED — hub ruling on root.go gofmt vs additive-only

The hub ruled **(b)**: the chip runs `gofmt -w internal/composition/root.go` and commits the
realignment. Not an exception to the grant — the hub's reading of it. The grant forbids editing
another owner's LOGIC in the seam; whitespace realignment forced by a line the chip was authorized
to add is nobody's logic.

The hub rejected the chip's recommendation (a) with a better argument than the chip had:
option (a) does not avoid the churn, it **displaces the authorship** — the next writer with
format-on-save delivers those four lines inside their diff, at a gate that cannot tell where they
came from. And a hub-side `gofmt -w` at merge is worse still: it puts an unreviewed edit inside
the merge commit, outside the dual gate that reviews the chip's diff.

**C12's proof form changed** (the criterion, not the chip's code) — committed to `main` @
`c9fecbc`:

> C12 | Grant de `root.go` foi **additive-only em tokens**, e o arquivo fica `gofmt`-limpo |
> `git diff -w <base>..HEAD -- internal/composition/root.go` **sem nenhuma linha removida**
> (o `-w` ignora o realinhamento de whitespace) + `gofmt -l internal/composition/root.go` sem
> saída + diff da região citado no payload do CLOSED

So C12 is now proven with `git diff -w` (token-level, which is what "additive-only" always meant
to protect) plus a clean `gofmt -l`, and the CLOSED payload quotes the RAW diff with the
realignment visible — not hidden.

Application deferred until slice S3 lands: S3's codex worker is mid-flight in this same worktree
and a concurrent chip commit races its `.git/index.lock`. Recorded here at ruling time so the
ordering is auditable rather than looking like the chip sat on the ruling.

### (superseded by the ruling above) REQUEST as originally sent

The grant's three added lines are committed with zero removals, so C12 passes as written. But the
added key `IdentityAnchors` is the longest in the `GenerationServiceConfig` literal, so gofmt
wants to widen the alignment column of the four existing keys — `gofmt -l root.go` reports the
file. Whitespace only, zero token change; no alignment exists that is both gofmt-stable and
leaves neighbours untouched.

Profile §2 L0 is `go build` + `go vet`; gofmt is not a ladder gate, so the committed state fails
nothing this chip is measured on. The chip committed the literal-additive reading and did NOT
reformat, because widening the column would put four removed lines in the C12 diff and that is
the hub's call, not the chip's. REQUEST sent to `local_99feb041`. Whichever way it rules, the C12
evidence quotes the real diff — a realigned file will not be presented as "additive-only" unless
the hub records it as one.

## Slice S1 — chip verification (not the worker's word)

Every claim below was re-run by the chip in the worktree, after the worker exited.

| Check | Command | Result |
|---|---|---|
| Write-set exactness | `git diff --name-only 917f7bb5..HEAD` | exactly the 8 declared files, zero undeclared |
| L0 build | `go build ./...` (absolute GOCACHE/GOMODCACHE) | exit 0 |
| L0 vet | `go vet ./internal/modules/connectors/... ./internal/modules/product_links/...` | exit 0 |
| L1 tests | `go test ./internal/modules/connectors/... ./internal/modules/product_links/...` | all `ok`, exit 0 |
| L0 format | `gofmt -l <the 8 files>` | empty |

Worker-reported caveat NOT reproduced: it claimed plain `go build ./...` was "blocked by
environment VCS stamping" and passed only with `-buildvcs=false`. Run from the worktree outside
the codex sandbox, plain `go build ./...` exits 0. This was a sandbox artifact of the worker's
own run, not a property of the repo, and is recorded as such rather than carried forward as a
caveat.

Substance spot-checks the chip made on the diff itself:

- `connectors/ports/marketplace_capability.go` — vocabulary is the FULL set of five anchors the
  engine reasons about, not the subset some provider supplies; `KnownIdentityAnchors()` returns a
  defensive copy.
- `connectors/application/marketplace_capability_service.go:129-155` — nil declaration ⇒
  `unsupported(providerCode, "identity anchors")`, NOT an empty slice. This is the ADR-17 hinge of
  the whole feature: "nobody told us" must not decay into "the provider supplies none". Unknown
  and duplicate anchors also error; the return is a clone.
- `connectors/adapters/mercado_livre/capability_adapter.go:79-91` — declares exactly
  `seller_sku`, `ean`, `title`. The `-` lines in this file's diff are gofmt field realignment
  from the added field, not logic edits. (C12's additive-only constraint binds `root.go`, which
  this slice does not touch.)
- `product_links/adapters/connectors/identity_anchor_adapter.go` — projects the complete
  vocabulary, holds `*MarketplaceCapabilityService` exactly as the cited precedent
  `inventory/adapters/connectors/stock_writer.go:12` does, carries the compile-time assertion
  `var _ productlinksports.ProviderIdentityAnchorReader = IdentityAnchorAdapter{}`, and contains
  no provider-name literal and no `providerCode ==` comparison (R2 / C3).
- Production `ProviderCapabilitySet()` declaration sweep: one site repo-wide, the ML adapter.

## Chip adjudications on the P2 plan

The planner's five decisions (Q1–Q5) are accepted as the design, with four corrections the chip
made as slice-card authority before dispatch. Recorded here because they change what the
implementers build, not just how it is phrased.

**A1 — canonical dimension form is for COMPARISON ONLY; the operator-facing detail keeps the
original tokens.** The planner proposed a reduced-rational signature `mm:<num>/<den>` as the
value compared inside `hardNegativeDimension`. But at `generation_service.go:654` the signature
IS the message: `fmt.Sprintf("hard-negative: medida/dimensão divergente %s≠%s", stDim, inDim)`.
Shipping the plan as written would turn the live queue message `50cm≠40cm` (M-05 EVIDENCE.md:379
records this exact text on the real account) into `mm:500/1≠mm:400/1`. That is a readability
regression on an operator-facing string, introduced by a fix meant to help the operator. S3 must
therefore return BOTH a canonical comparison key and the original display tokens, and its C6 test
asserts the detail for `50cm` vs `40cm` still contains `50cm` and `40cm` and does NOT contain
`mm:`.

**A2 — the identity-anchor lookup failure names the provider code.** `PRODUCT_LINKS_PROVIDER_
IDENTITY_ANCHORS_UNAVAILABLE` with an unnamed provider is an unactionable error on a run that
touches N providers. ADR-17 honest-failure means the failure says which provider it could not
resolve.

**A3 — the day-1 golden test uses the REAL M-05 pair, not the chip prompt's wording.** The
dispatch prompt's EXEMPLO-IO has the listing and the ERP product swapped. `M-05-auto-vinculo/
EVIDENCE.md:377-380` records the live pair: listing `Toalheiro Simples Soul Zen 50cm …` (cm)
against ERP product `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` (mm), CODPROD 33698 exact SKU
match, listing carries NO EAN. So the post-fix expectation is **CONFIRM / confidence 70 / band
MEDIA** — the third CONFIRM M-05 predicted — and NOT `ALTA/ACCEPT`, which the planner's synthetic
concordant-CODPROD+EAN fixture produces. S3 carries both cases: the real golden (CONFIRM) and the
synthetic concordant (ACCEPT).

**A4 — Q5 accepted, with the residual stated rather than hidden.** Independent defaults
(candidates 20 / links 2000 / audits 10000) make all 29 links reach `/vinculos` with no wire
change, and leave `mutations/adapters/productlinks/writer.go:163,191` (`Limit: 2000`) compiling
AND behaviourally unchanged. Residual: above 2000 resolved links the list truncates silently
again. This chip cannot close that — R1 forbids the response-shape change a truncation flag
needs — and the contract's own U3 note assigns the never-lies-about-truncation screen to M-06.
Stated here so the closure does not read as if the class were eliminated.

**A5 — `math/big.Rat` accepted over integer micrometres.** The inch factor 25.4 plus arbitrary
decimal input reintroduces rounding in any fixed-scale integer representation; `big.Rat` is
stdlib and exact. Recorded because it is the kind of decision a reviewer should see argued.

## Slice S4 — chip verification (F-03)

Worker returned GREEN, committed `030fa58c`. Verified independently of the worker's claims.

| Check | Command | Result |
|---|---|---|
| Write-set exactness | `git show 030fa58c --stat` | exactly the 2 declared files |
| L0 vet | `go vet ./internal/modules/product_links/...` | exit 0 |
| L1 tests | `go test -count=1 ./internal/modules/product_links/application` | `ok` 1.612s |
| C8 guard is not theatre | read `limitingCandidateStore` / `limitingWorkflowStore` | both RECORD the observed limit AND truncate (`items = items[:limit]`); the base stubs take `_ int` and discard it, so the wrappers are what make the assertion load-bearing |
| C8 pins the specified defaults | read the assertions | asserts `20 / 2000 / 10000` as exact literals, not "≥29"; the old `limit*5` bug would yield `100` audits and fail here |

The worker again reported plain `go build ./...` failing with VCS stamping and passing only under
`-buildvcs=false`. It does NOT reproduce chip-side (exit 0). Profile §3 already ratifies the
chip's re-run outside the codex sandbox as the verification of record for this class. Recorded as
a sandbox artifact; NOT carried forward as a repo caveat.

### A gofmt violation the CRLF noise was hiding — commit `f6736d0`

The S4 card scoped `gofmt -l` to the two authored files, which was the right fix for the S3
card defect but is still blind: both files are in the nine-file CRLF set, so `gofmt -l` names them
whatever their content. The chip re-ran gofmt against a **CR-stripped copy** in the scratchpad,
which isolates token changes from the line-ending artifact, and that surfaced a real defect —
`ProviderCode:` was written with one space too many in both 29-link fixtures.

Fixed by hand in `f6736d0`, NOT with `gofmt -w`, which would have rewritten all ~900 lines from
CRLF to LF. Verified after: CR-stripped `gofmt -l` empty, `go test -count=1` still `ok`,
`git diff --stat` exactly 2 insertions / 2 deletions.

**Technique worth carrying upstream:** scoping a gofmt gate to authored files is necessary but not
sufficient on a CRLF checkout. The check that actually discriminates is
`tr -d '\r' < file > tmp && gofmt -l tmp`. The F-03 reviewer independently re-ran exactly this and
found zero diff, which is external confirmation the fix is correct.

## P4 — per-feature adversarial review (§4 obligation 3)

One independent reviewer per FEATURE, over all of that feature's slices together, before the
feature closes. Reviewers were given the defect, the contract and the attack surface, but were
**deliberately NOT given the chip's own three suspicions about F-02**, so that a pass would be
genuine independent evidence rather than an echo.

| Feature | Commits reviewed | Verdict | Findings |
|---|---|---|---|
| F-01 | `5bc55219` + `b9da6d2` + `633bf9fa` | PASS-WITH-FINDINGS | 2 SHOULD-FIX, 1 NIT |
| F-02 | `11e68f6f` | **FAIL** | 1 BLOCKING, 2 SHOULD-FIX |
| F-03 | `030fa58c` | PASS-WITH-FINDINGS | 1 SHOULD-FIX |

The F-01 reviewer re-diffed `base..b9da6d2` explicitly rather than `base..633bf9fa`, because the
naive range interleaves F-02's dimension work into the same file and would have conflated the two
features. Recorded because it is the correct method for a chip whose features share a file.

### Convergence check on the withheld suspicions

All three concerns the chip withheld were found independently by the F-02 reviewer:

| Chip's withheld concern | Reviewer's finding | Status |
|---|---|---|
| widened unit pattern creates false positives | F1, BLOCKING, with a concrete `"N in 1"` scenario the chip had not constructed | **converged, and the reviewer's case is stronger** |
| `display` sliced from `original` using offsets computed on `lowered` | F2, SHOULD-FIX, same analysis, plus the panic path when `lowered` is longer | converged |
| parse failure abandons the whole side (fail-open) | named as a maintenance-coupling risk, agreed direction, correctly assessed as unreachable today | converged |

Three-for-three convergence from a reviewer working blind is what makes the F-02 FAIL credible
rather than a matter of taste, and it is the reason the corrective is being treated as blocking.

### Chip adjudications on the P4 findings

**A6 — F-02 F1 (BLOCKING) accepted; fix by REVERTING the token set, not by adding a boundary.**
`git show 917f7bb5` proves the base alternation was
`...(?:mm|cm|pol|")|\d+(?:[.,]\d+)?\s*m\b` — it never contained `inches`, `inch` or `in`. S3's
mandate was unit EQUIVALENCE (`50cm` == `500MM`), which needs mm/cm/m/pol only; the inch aliases
were unrequested scope that introduced a live regression. A `\b` after the unit would ALSO break
`50cmx30cm`, which the base pattern handles today, so reverting is both the smaller change and the
one that preserves existing behaviour byte-for-byte. Verified no test depends on the aliases: the
only inch-family case (`generation_service_test.go:585`) is KEYED `"inches"` but its DATA is
`"Produto 2pol"` / `"Produto 50.8mm"`, i.e. it exercises `pol`. Dispatched as S5 part A, together
with fail-CLOSED parse fallback (F2 direction), matching on `original` with `(?i)` to kill the
offset aliasing, and key/display pair compaction.

**A7 — F-01 layering leak (SHOULD-FIX) accepted; the redundant validation is removed.**
`generation_service.go:14` imports `connectorsports` into the product_links APPLICATION layer, and
`resolveIdentityAnchors` re-validates anchor names against `connectorsports.KnownIdentityAnchors()`
at `:150-153` and `:169-172`. `product_links/ports/provider_identity_anchor_reader.go` types its
anchor as a plain `string` precisely so product_links does not depend on the connectors vocabulary
type — this validation reaches around that decoupling. It is also dead three times over: the
connectors service already rejects unknown and duplicate anchors, and the adapter BUILDS its result
by iterating the known set and structurally cannot emit anything outside it. Ruling:
module-boundary integrity outranks a redundant unreachable check, and a product_links-local copy of
the vocabulary would be worse (two lists to keep in sync). Dispatched as S5 part B. The `nil` and
empty-`providerCode` checks are KEPT — different failure, and they are the fail-closed behaviour
C4 depends on.

**A8 — F-01 duplicate contradictory reasons (SHOULD-FIX) DEFERRED to the hub, not fixed here.**
A provider declaring only `title` would produce, for the same anchor, both
`{seller_sku, UNAVAILABLE, "seller_sku sem correspondência"}` (implies a search happened) and
`{seller_sku, UNAVAILABLE, "provider não fornece a âncora seller_sku"}` (says no search was
possible). Real, and semantically contradictory to an operator. Deferred because: it is DORMANT —
Mercado Livre is the only registered provider and supplies all three of seller_sku/ean/title;
fixing it means inventing suppression semantics for a provider that does not exist yet, against
"no abstraction without a second named consumer"; and C4 is satisfied as written (the two cases are
distinguishable via `Detail`, per hub ruling R1 which forbids a 4th enum value). Carried to the hub
in the CLOSED payload as a named follow-up, NOT left as a TODO in the code.

**A9 — F-03 test-gap (SHOULD-FIX) accepted, queued behind S5.**
The default-path test pins `2000`/`10000` only at `Limit: 20`, so a future regression deriving the
link/audit defaults from `Limit` by a formula that happens to agree at 20 would pass. The
implementation is correct today (literal independent constants, verified by reading
`resolution_service.go:369-381`); this is a hole in the regression net, not a live defect. Queued
rather than dispatched immediately because S5 is writing in the same worktree and two concurrent
codex workers race `.git/index.lock` — a MIS-004 field finding this chip has been serialising
against throughout.

**NITs accepted as-is, no change:** the `declaration == nil` defensive branch in
`resolveIdentityAnchors` is unreachable via the real adapter but cheap and fail-closed; voltage
remains lexically extracted and its decimal blindness (`110.0V`) is an out-of-scope sibling
finding, recorded for the hub.

## Slice S5 — chip verification (corrective, `f92ca9c7`)

Worker returned GREEN. Verified independently of every claim it made.

| Check | Command / method | Result |
|---|---|---|
| Write-set exactness | `git show f92ca9c7 --stat` | exactly `generation_service.go` + `generation_service_test.go`, 92 insertions / 31 deletions |
| A1 — inch aliases gone | `grep -c 'inches\|"inch"\|"in"' generation_service.go` | **0** |
| A1 — token set back to base | read `:655` | `(?i)\d+\s*x\s*\d+(?:\s*x\s*\d+)?\|\d+(?:[.,]\d+)?\s*(?:mm\|cm\|pol\|")\|\d+(?:[.,]\d+)?\s*m\b` — identical to `917f7bb5` apart from the `(?i)` prefix |
| A1 — candidate list | read `:735` | `{"pol", "mm", "cm", `"`, "m"}`, longest-suffix-first preserved so `m` cannot shadow `mm`/`cm` |
| A2 — fail CLOSED | read the `SetString` branch | on parse failure appends `{key: token, display: originalToken}` and `continue`s — the side's other tokens survive; the old `return "", "", false` is gone |
| A3 — no offset aliasing | read the match loop | matches `original` directly; `originalToken` sliced with offsets native to that same string; lowering is per-token, not whole-string |
| A4 — key/display correspondence | read the tail | one `dimensionPair` per token, `slices.SortFunc` + `slices.CompactFunc` both keyed on `key`, then the two slices are projected from the SAME compacted pairs — they cannot diverge |
| B1 — layering leak closed | `grep -rn 'connectors/ports' product_links/application/*.go` | only `generation_integration_test.go:11`, a PRE-EXISTING test-only import; zero hits in production code |
| L0 build | `go build ./...` (plain, no `-buildvcs=false`) | **exit 0** |
| L0 vet | `go vet ./...` (full tree) | exit 0 |
| L1 tests | `go test -count=1 ./internal/modules/{product_links,connectors,mutations}/... ./internal/composition/...` | all `ok`, exit 0, 0 failures |
| Class re-sweep | CR-stripped `gofmt -l` over all 15 write-set files | swept=15, token-dirty=**0** |

New tests inspected for substance, not counted: `TestNInOneTitleIsNotReadAsADimension` drives the
full generator on the reviewer's exact scenario; `TestDimensionDisplaySurvivesLengthChangingCaseFold`
calls `detectHardNegative("İnox 50cm", "İnox 40cm")` and asserts the detail contains `50cm`/`40cm`
AND contains neither `mm:` nor the truncated fragments `ox 5`/`ox 4` — that last assertion is what
makes it a real guard against the offset bug rather than a smoke test.

**Honest limitations the worker reported, verified true and accepted:**
- the fail-closed fallback is unreachable by construction with the five-token alternation, so
  there is NO must-fail proof for it — stated, not faked;
- the `connectors/ports` grep is not zero directory-wide because a pre-existing test file imports
  it; changing that file was outside the frozen write set. Correct call.

**Worker sandbox artifact, third occurrence:** plain `go build ./...` again reported VCS-stamping
exit 128 inside codex and passed only with `-buildvcs=false`. Chip-side it is exit 0, as it was for
S1, S2 and S4. Profile §3 ratifies the chip's re-run outside the sandbox as the verification of
record for this class. Recorded; NOT carried forward as a repo caveat.

## Slice S6 — chip verification (F-03 test hardening, `1df627f8`)

Worker returned GREEN. Verified independently of every claim it made.

| Check | Command / method | Result |
|---|---|---|
| Write-set exactness | `git show 1df627f8 --stat` | exactly `resolution_service_test.go`, **41 insertions / 0 deletions** — test-only, as the card froze it |
| Production code untouched | `git diff --stat 030fa58c..1df627f8 -- …/resolution_service.go` | **empty** — the must-fail mutation was fully reverted |
| Test is substantive, not a smoke test | read the diff | table over `5`/`20`/`500` with `LinkLimit`/`AuditLimit` **left zero**, asserting the candidate limit tracks the input while `2000`/`10000` stay constant across every row — exactly the invariant `Limit`-derived formulas would break |
| No new stub types | read the diff | reuses `limitingCandidateStore` / `limitingWorkflowStore`; existing tests unmodified |
| L1 tests | `go test -count=1 ./internal/modules/product_links/...` | all 5 packages `ok`, exit 0 |

**Must-fail proof (R5), as reported and accepted:** replacing the literals with `candidateLimit*100`
and `candidateLimit*500` turns the new test RED on the `5` and `500` rows while the `20` row still
passes **and both pre-existing default tests stay GREEN** — which is the actual point of the proof:
it demonstrates the old net was blind to this regression class and the new one is not. A proof that
had reddened the old tests too would have shown nothing about the gap.

**Tree hygiene note.** After S6 the production file showed as dirty despite having no content diff.
Byte comparison isolated it: worktree 38905 bytes / 898 CR against a HEAD blob of 38007 bytes / 0 CR
— a difference of exactly 898 bytes, i.e. line endings and nothing else. Cause: S6's must-fail cycle
restored the file through git, so checkout re-materialised CRLF where S4's worker had left LF.
Resolved with `git add`, which normalised CRLF→LF, produced a blob identical to HEAD, staged nothing
and left the tree clean. `git checkout --` / discard was deliberately NOT used — doctrine forbids
deleting unknown state, and at diagnosis time the 898 bytes were not yet known to be inert.

## P4 round 2 — re-review of F-02 over the corrective

`f92ca9c7` reopened F-02, so the feature was re-reviewed end-to-end rather than spot-checked.
Verdict **PASS-WITH-FINDINGS**; the BLOCKING defect is closed. F-01 was likewise re-reviewed over
the same commit (D13) and came back PASS with zero findings.

Every claim the corrective made was independently confirmed: token set back to the exact base
alternation (only `(?i)` differs), unit candidates reduced to `{pol, mm, cm, ", m}`, no `(?i)`
collision with the clothing-size grade, C6 golden `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm`
PASS at CONFIRM/MEDIA/70, C6 must-fail (`50cm` vs `40CM`) still contradicting, C7 kit/voltage/NxM
all still capped BAIXA/REJECT, `TestNInOneTitleIsNotReadAsADimension` PASS.

Worth recording as method: the reviewer did not take the display-guard test on trust. It rebuilt the
PRE-corrective function standalone and ran the same input through it, confirming it returned
`display="50c"` — truncated — for `"İnox 50cm"`. That is what distinguishes a guard proven
load-bearing from one merely asserted to be, and it is the same must-fail discipline R5 requires,
applied by a reviewer to someone else's test.

### A10 — the `pol` prefix finding: REPORTED to the hub, deliberately not fixed here

The re-review's SHOULD-FIX F1: `pol` is unbounded, so a digit followed by any word starting `pol`
is read as inches. Repro
`detectHardNegative("Blusa 30 poliester 70 algodao Azul", "Blusa 40 poliester 60 algodao Azul")`
→ `"hard-negative: medida/dimensão divergente 30 pol≠40 pol"`. A fabricated measurement driving a
REJECT is an ADR-17 violation in effect, so it is real.

It is the same **shape** as the `in` defect A6 fixed — a unit token matching as the prefix of an
ordinary Portuguese word. Two instances of one shape, so §11 clause 1 has not fired; but the
mechanism is now known, so the class was swept exhaustively by tool rather than argued
(`unit-prefix-sweep.go.txt`, alongside this pack, run against the alternation copied verbatim from
`generation_service.go:655`; kept as `.go.txt` so no Go tool ever picks it up as package source):

| token | prefix-vulnerable | evidence |
|---|---|---|
| `mm` | theoretical only | `30 mmhg` matches, but no real PT-BR retail word starts `mm` after a digit |
| `cm` | theoretical only | `30 cmyk` matches; same reasoning |
| `pol` | **LIVE** | `30 poliester` → `30 pol`; `2 polias` → `2 pol` — both ordinary retail phrasing |
| `m` | no | already carries `\b` in the base pattern |
| `"` | no | not a word prefix |

The base author guarded `m` with a boundary and missed `pol`. **One live member; class closed.**

Candidate fix, verified against the full probe set:
`(?:mm|cm|polegadas|polegada|pol\b|")` — kills `30 poliester` and `2 polias`, PRESERVES the true
measurements `24 polegadas` and `2 pol`, and leaves the golden `500MM CR/POLIDO` byte-identical.
(A bare `\b` after `pol` was rejected: it would also stop matching `24 polegadas`, trading a false
positive for a false negative.)

Not applied in this chip, for four reasons stated together because no one of them carries it alone:
it is **pre-existing at `917f7bb5`** and provably not worsened by the corrective (the reviewer
verified the golden is unaffected because `POLIDO` follows `CR/`, not a digit); it is **outside
C6/C7**, which scope F-02 to unit equivalence; the independent reviewer that found it explicitly
recommended a follow-up ticket over blocking; and a third change to F-02's dimension logic reopens
the feature review for a third time, which is precisely the ceremony §11 exists to make expensive.
Carried to the hub in the CLOSED payload **with the fix and the sweep attached**, so the ruling is
a short one, not a re-investigation.

The re-review's other two findings need no adjudication: the fail-closed branch being unreachable
and therefore untested was already declared honestly by both the worker and the chip (S5 table
above), and the display-token ordering NIT follows the key sort by construction.

## ROUND-1 FULL ANALYSIS — the gofmt/CRLF signal-masking class (profile §11, hub ruling R4)

Triggered by clause 1 of the third-round rule — a third defect of the **same shape**, in three
different artifacts. Not clause 2: no single criterion has reached a third correction round.

### The shape, in one sentence

*On this repo `gofmt`'s output carries a constant, content-independent CRLF signal, so any gate
phrased as "`gofmt -l` must be empty" is either unsatisfiable-by-construction or blind to real
token violations — and both failure directions have now occurred in this chip.*

### The three instances

| # | Artifact | Direction of failure | Outcome |
|---|---|---|---|
| 1 | `internal/composition/root.go` | gofmt-clean vs additive-only-in-tokens pulled against each other; no alignment satisfies both | escalated as REQUEST; hub ruled `gofmt -w` + commit, and REVISED C12's proof form to `git diff -w` on main @ `c9fecbc` |
| 2 | S3 slice card | **unsatisfiable**: card demanded a directory-wide `gofmt -l` come back empty; nine files there are CRLF, so it never can | worker correctly returned BLOCKED rather than reformat outside its write_set; chip verified and committed the work as `11e68f6f` |
| 3 | S4 output, `resolution_service_test.go` | **blind**: a real misalignment (`ProviderCode:` with one space too many, twice) sat inside the `^M` noise and passed the file-scoped check | shipped green in `030fa58c`; caught only by a CR-stripped re-run; fixed in `f6736d0` |

Each was point-fixed correctly on its own terms. Instance 2 even produced a *narrower* gate
(file-scoped instead of directory-scoped) — which is exactly what let instance 3 through, because
narrowing the scope does nothing about the masking. That is the local-maximum failure the rule
exists to catch.

### The discriminating check

```bash
tr -d '\r' < <file> > <tmp> && gofmt -l <tmp>
```

Empty output here means token-clean regardless of line endings. This is the only form of the check
that carries information on this repo.

### The sweep — exhaustive, by tool, including the clean sites

Every `.go` file in the chip's write-set (`git diff --name-only 917f7bb5..HEAD -- '*.go'`), each
run through both the raw and the discriminating check. `raw` = `gofmt -l <file>`;
`tok` = `gofmt -l` on the CR-stripped copy.

| Endings | raw | tok | File | Classification |
|---|---|---|---|---|
| LF | 0 | 0 | `internal/composition/root.go` | correct |
| LF | 0 | 0 | `connectors/adapters/mercado_livre/capability_adapter.go` | correct |
| LF | 0 | 0 | `connectors/adapters/mercado_livre/capability_adapter_test.go` | correct |
| LF | 0 | 0 | `connectors/application/marketplace_capability_service.go` | correct |
| LF | 0 | 0 | `connectors/application/marketplace_capability_service_test.go` | correct |
| LF | 0 | 0 | `connectors/ports/marketplace_capability.go` | correct |
| LF | 0 | 0 | `product_links/adapters/connectors/identity_anchor_adapter.go` | correct |
| LF | 0 | 0 | `product_links/adapters/connectors/identity_anchor_adapter_test.go` | correct |
| LF | 0 | 0 | `product_links/application/auto_link_policy_test.go` | correct |
| LF | 0 | 0 | `product_links/application/generation_integration_test.go` | correct |
| LF | 0 | 0 | `product_links/application/generation_service.go` | correct |
| LF | 0 | 0 | `product_links/application/generation_service_test.go` | correct |
| **CRLF** | **1** | 0 | `product_links/application/resolution_service.go` | correct — raw hit is `^M` only |
| **CRLF** | **1** | 0 | `product_links/application/resolution_service_test.go` | **was defective, fixed in `f6736d0`**; raw hit is now `^M` only |
| LF | 0 | 0 | `product_links/ports/provider_identity_anchor_reader.go` | correct |

**Class closed: 15 sites enumerated, 15 token-clean, 1 defect found and fixed.** Two sites remain
raw-`gofmt -l`-dirty and that is CORRECT — they are CRLF and token-clean.

### Two facts the sweep established that were being assumed wrong

1. **The CRLF set is a property of individual files, not of the directory.**
   `product_links/application` holds nine CRLF files (`batch_service{,_test}.go`,
   `decision_trail_test.go`, `import_service{,_test}.go`, `resolution_service{,_test}.go`,
   `summary_service{,_test}.go`) and four LF files
   (`auto_link_policy_test.go`, `generation_integration_test.go`, `generation_service{,_test}.go`).
   The S5 card was written asserting `generation_service*.go` were CRLF. They are LF. The error is
   harmless — the card told the worker not to run `gofmt -w` and to verify via CR-strip, and
   `tr -d '\r'` is a no-op on an LF file — but it is recorded rather than quietly corrected,
   because an unrecorded wrong premise in a card is how instance 2 happened.

2. **No worker converted line endings inside any commit**, verified by
   `git diff --numstat 917f7bb5..HEAD`: every file changes proportionately to its real edit
   (largest: `generation_service.go` 149/76 on an ~800-line file). A CRLF↔LF conversion inside a
   commit would show every line as changed. The chip's diff is reviewable line-by-line.

### CORRECTION to this section, found after it was first written

The original fact 2 above also asserted that `.gitattributes` covers only `*.sh`, that Go files
therefore get no normalisation, that the endings live in the blob, and that the mixed working-tree
endings are "PRE-EXISTING on main". **That reasoning was wrong**, and the remedy it produced would
have sent the hub a no-op backlog item. Corrected here rather than silently rewritten, because the
pack is the audit trail.

What actually holds, measured:

```
git config --get core.autocrlf   ->  true
git show HEAD:<each of the nine "CRLF" files> | tr -cd '\r' | wc -c   ->  0
```

- **Every blob in the repository is LF.** There is no CRLF anywhere in git's object store; the
  nine "CRLF files" are LF in the repo.
- `core.autocrlf=true` materialises CRLF **at checkout**. A file a tool later rewrites (the Edit
  tool, a codex worker) lands as LF and stays LF. So the working-tree mix is generated LOCALLY by
  which files happen to have been rewritten since checkout — it is not a property of `main` and
  says nothing about history.
- The mechanism statement at the top of this section is UNAFFECTED and still correct: `gofmt` reads
  the working tree, sees CRLF there, and flags unconditionally. The CR-stripped check remains the
  only informative form.
- The diff-reviewability conclusion is unaffected and in fact STRONGER: because blobs are always
  LF, line-ending noise can never enter a commit at all. That is also why `f6736d0` was 2/2 lines.

**How it surfaced.** After S6 committed, `git status` reported
`M product_links/application/resolution_service.go` while `git diff --numstat` reported nothing.
Byte comparison against the HEAD blob: worktree 38905 bytes with 898 CR, blob 38007 bytes with 0 CR
— a difference of exactly 898 bytes, i.e. line endings and nothing else. Cause: S6's must-fail
cycle edited that production file and restored it through git, so checkout re-materialised CRLF
where S4's worker had left LF; the stat cache then went stale. Resolved with `git add`, which
normalised CRLF→LF, produced a blob identical to HEAD, staged nothing, and left the tree clean.
The production file is byte-identical to `f92ca9c7` — S6's claim of "no content diff" was true, and
its commit is correctly test-only (41 insertions, one file).

**Revised remedy for the hub** (replaces the version below): a single `.gitattributes` line,
`*.go text eol=lf`, forces LF at checkout and removes the masking at its source. It is NOT a
repo-wide renormalise commit and touches no file's content, because the blobs are already LF —
far cheaper than this section originally claimed. Still hub-owned (`.gitattributes` is a shared
seam), still outside a chip that owns `product_links`.

### Remedy: informational, not structural

The rule requires asking whether the fix is a guard, a funnel, a type, a lint rule, or a corrected
sentence — and requires an independent adversarial judgement briefed AGAINST any abstraction before
one is added. **No abstraction is being added, so no such dispatch is owed.** Reasoning, recorded
so a reviewer can disagree with it:

- A committed helper script or a lint rule would be a second source of truth about formatting,
  living beside `gofmt` itself, owned by nobody, and would have to be maintained past the day the
  line endings get normalised — the same drift trap that killed CHIP-M05's proposed domain
  constructor.
- The actual root cause is upstream of this chip and outside its ownership: nine files carry CRLF
  in the blob. Normalising them is a one-line `.gitattributes` change plus a repo-wide renormalise
  commit — a shared-seam edit, hub-owned, and far outside a chip that owns `product_links`.
- What this chip can honestly produce is the discriminating command and the evidence that it
  works. That goes to the hub as a proposed extension of profile §3 field finding F-ENV-M01, which
  today says only "scope gofmt gates to milestone-authored dirs" — advice that instance 3 proves
  insufficient, because scoping does not unmask.

Carried in the CLOSED payload as a named follow-up: (a) amend F-ENV-M01 with the CR-stripped form
and the reason scoping alone is not enough; (b) hub-owned backlog item to normalise the nine CRLF
files, with `.gitattributes`, in a commit that touches nothing else.

## Criteria table

All evidence below is against the final tip `1df627f8` (code) with the pack filed at `8e37958`.
Base is `917f7bb58e385847fba5612201823f9db48791c6` throughout.

| ID | Verdict | Evidence |
|----|---------|----------|
| C1 | **PASS** | declaration `connectors/ports/marketplace_capability.go:22-41` (`IdentityAnchor` type, the five constants, `KnownIdentityAnchors()`); product_links' own port `product_links/ports/provider_identity_anchor_reader.go`; adapter `product_links/adapters/connectors/identity_anchor_adapter.go:9-23`; wiring `internal/composition/root.go:387` + `:543` |
| C2 | **PASS** | `git grep -n mandatoryUnavailableReasons -- 'apps/server_core/**/*.go'` = **0** in tracked source (base had **8**, all in `generation_service.go`). `"marca"`/`"refforn"` in `product_links/application` production code = **0**; the 8 remaining hits are all in `generation_service_test.go` (`:88,:89,:143,:144,:170,:173,:498,:502`), which the criterion permits and which are the tests that PROVE the names now arrive from the declaration |
| C3 | **PASS** | `git grep 'mercado_livre\|ProviderCode *==\|ProviderCode *!=\|switch .*[Pp]rovider' -- 'product_links/**/*.go' ':!*_test.go'` = **ZERO HITS**, over the **34 tracked `.go` files** listed in the C3 sweep block above (the criterion requires the swept list, not just the count) |
| C4 | **PASS** | the two details are distinct and both present: `generation_service.go:536-537` / `:609-610` emit `"<anchor> sem correspondência"` (provider supplies it, this listing has no value) and `:624` emits `"provider não fornece a âncora %s"` (provider does not supply it). Tests `generation_service_test.go:135` and `:1005`. Per R1 the distinction is carried in `Detail`; no 4th enum value exists |
| C5 | **PASS** | `git diff --stat <base>..HEAD -- openapi/ packages/sdk-runtime/` → **empty**. Also visible in the C10 name list: no path under either tree appears |
| C6 | **PASS** | `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm` (`:527`, the EXEMPLO-IO golden: `500MM` listing vs `50cm` ERP → CONFIRM/MEDIA/70), `TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate` (`:558`), `TestDimensionCanonicalizationUsesExactMillimetres` (`:582`). Contradiction still fires: `TestDifferentCanonicalDimensionsStillRejectConcordantCandidate` (`:653`). Must-fail below |
| C7 | **PASS** | corroborated path still blocked: `TestCase6DokaKitHardNegativeCapsBaixaReject` (`:902`), `TestCase7VoltageHardNegativeCapsBaixaReject` (`:933`), `TestCase10DimensionHardNegativeCapsBaixaReject` (`:1026`), `TestDimensionPresenceAndGradeRulesRemainNonBlocking` (`:682`) — all with CODPROD+EAN concordant |
| C8 | **PASS** | `resolution_service.go:369-381` three independent literals; `TestListLinkWorkflowsUsesIndependentDefaultLimitsAndReturnsAll29Links` (`resolution_service_test.go:361`) asserts `len(links)==29` on a 29-link fixture; `TestListLinkWorkflowsDefaultsDoNotVaryWithLimit` (`:406`) pins the independence at 5/20/500; `TestListLinkWorkflowsHonorsIndependentExplicitLimits` (`:447`). Must-fail below |
| C9 | **PASS** | `git diff --name-only <base>..HEAD \| grep -i migration` → **zero**; the C10 list contains no `migrations/` path and no `product_links/adapters/postgres/` file at all, so the chip writes no `UPDATE` and no backfill. Characterised honestly: a pre-existing `ON CONFLICT … DO UPDATE SET reasons = EXCLUDED.reasons` does exist at `link_candidate_repo.go:79`, but it is the candidate **regeneration** upsert, untouched by this chip and depended on by the hub's own U1 ("depois de regerar candidatos"). R3 forbids retro-editing persisted motivos; it does not forbid regeneration from producing fresh ones |
| C10 | **PASS** | full `git diff --name-only <base>..HEAD` pasted in the C10 block above — 17 paths, of which 2 are this evidence pack. `grep -c '^apps/web/'` = **0** |
| C11 | **PASS** — L0+L1 green chip-side; governance rung hub-run per R-b, **differential PASS** (53 violations at chip tip == 53 at BASE, outputs identical line for line; lane is red on `main` and not green here — see the honesty note) | see the Ladder section below + hub evidence `main@d36d89a` |
| C12 | **PASS** | `git diff -w <base>..HEAD -- internal/composition/root.go` → **0 removed lines**, exactly 3 added, quoted verbatim in the C12 block above and in the CLOSED payload; `gofmt -l` clean (proven CR-stripped per the ROUND-1 analysis — the naive file-scoped check is uninformative on this checkout) |

## Ladder

Profile §2 bindings. `GOCACHE`/`GOMODCACHE` bound ABSOLUTE before every Go command, as §2 requires
on Windows (echoed into the transcript:
`GOCACHE=/c/…/.claude/worktrees/chip-anchors/apps/server_core/.gocache`).

| Rung | Command | Result |
|---|---|---|
| L0a | `go build ./...` (full tree, plain — no `-buildvcs=false`) | **exit 0** |
| L0b | `go vet ./...` (full tree) | **exit 0** |
| L1 | `go test -count=1 ./internal/modules/product_links/... ./internal/modules/connectors/... ./internal/modules/mutations/... ./internal/composition/...` | **exit 0**, 27 packages, all `ok`, zero failures |

L1 scope is touched-packages-plus-guard-suites, not a full sweep, per §2: "full sweep only when
migrations/platform touched" — and this chip has **zero** migrations (C9). The four trees are the
complete set the diff reaches: `product_links` (owner), `connectors` (the declaration side),
`mutations` (the only external consumer of `ListLinkWorkflowsInput`, via
`adapters/productlinks/writer.go:163,191`), and `composition` (the root.go grant).

**No allowlist entry was needed.** Zero failures, so profile §2's known-failure list is not being
leaned on. Recorded explicitly because a GREEN-with-allowlist run and a genuinely green run are
different claims and the contract asks for the citation only when the former applies. `apps/web`
lanes are not owed: zero web files in the diff (C10), so the `TS2688` allowlist entry is likewise
untouched here.

**Governance rung: REQUEST ao hub por R-b; L0+L1 PASS chip-side com as saídas acima.**

The conflict, for the record: §2 and C11 both require the lane from a **clean detached worktree**
with the 40-hex BaseSha, because a main checkout sweeps `.claude/worktrees/*` and false-fails. The
dispatch prompt independently forbids this chip from creating another worktree or running
`git checkout` in the primary repo, naming it the mission's most-recurring launch hazard (4/4
chips). Both cannot be satisfied by the chip. Per "disagreement = BLOCKED with evidence, never a
unilateral decision", the chip sent `REQUEST governance-lane` with both readings and the exact
command rather than picking one.

**HUB RULING (received 2026-07-26): R-b — the hub runs the lane.** Decided on ownership, not on the
literal reading of the ban. Three reasons, in the hub's order of weight: (1) the governance lane is
already a hub seam by charter — `hub-ops` carries "governance-lane runs in a clean worktree" in its
own description, same family as the dev stack; (2) worktree topology is hub-owned (profile §6), and
the ban is the barrier that stops a chip discovering the safe path wrongly — "barreira não vira
permissão só porque existe um caminho seguro dentro dela"; (3) provisioning cost — a fresh worktree
inherits neither `node_modules` nor `.gomodcache` (profile §3 forbids the junction), so the chip
would pay a full `npm ci` to run one scan, while the hub's setup already exists.

Consequences, as directed: C11 is **OPEN** on this rung from the chip's side — **not** PASS, **not**
BLOCKED. The hub runs the lane in parallel on tip `8e37958` and attaches the result at acceptance.
CLOSED is explicitly **not** held for it; drift, if the lane reports any, returns as a post-CLOSED
correction rather than a closure blocker. The P6 gate judges what is the chip's.

**Upstream amendment candidate, registered by the hub:** the dispatch prompt and C11 contradict each
other in *any* chip that carries a governance rung. C11's text should read "hub-run" from dispatch
rather than leaving the chip to discover the conflict in flight. Hub's characterisation: a doctrine
defect, not a code defect.

### Governance rung — RUN BY THE HUB, differential PASS

Hub evidence, committed to `main` @ `d36d89a`: `_chip-anchors/hub-governance-lane.md`. Cited here,
**not re-proven** — the chip did not run the lane and does not own it.

Method: clean detached worktree created OUTSIDE `.claude/worktrees/` (the scanner sweeps that
directory and falsifies its own result), and the SAME lane run twice — once at chip tip `8e37958a`,
once at BASE-SHA `917f7bb5…`.

| Run | Result |
|---|---|
| chip tip `8e37958a` | `status=failed`, exit 1, **53 violations** |
| BASE-SHA `917f7bb5…` | `status=failed`, exit 1, **53 violations** |
| `Compare-Object` of the two outputs (175 lines each) | **empty — identical line for line**, including all 14 `baseline_exception` entries |

No violation names a file this chip wrote. `GOV_MODULE_COVERAGE` fires on `sourcekind` +
`tenant_config` (M-02 debt, no `modules.json` entry); `RCFG_*` fires on the legacy
magalu/amazon/shopee adapters and `MC_ERP_SOURCE` in `root.go`.

**Hub verdict: C11's governance rung PASSES on the differential reading — zero new violations
introduced.**

**Honesty note, at the hub's explicit instruction and not buried:** this is **NOT a green lane**.
The lane is red on `main` and was red before this chip existed. Accepting a chip over a red lane is
legitimate here *only* because the difference is null and the nullity is proven by a line-for-line
comparison of two runs — not because the red was waved through. The 53 violations are the hub's
debt, not this chip's, and they remain outstanding.

**Synchrony check (hub asked, chip verified by tool rather than asserting).** The hub measured at
`8e37958a`; the branch is now at `b954783d`. `git diff --name-only 8e37958..b954783` returns exactly
one path — this `EVIDENCE.md`. Zero `.go` files, and `git diff --stat 8e37958..b954783 -- apps/` is
empty, so the entire code tree is byte-identical between the measured tip and the current one. The
two intervening commits (`078ae7be`, `b954783d`) are both documentation. The differential therefore
still holds and **no re-run is needed before merge**. One correction to the hub's framing: S4
(`030fa58c`) sits far below `8e37958a` and was already inside the measured tip — nothing after that
tip is code at all.

## Dual gate

Two reviewers on the fixed-SHA diff `917f7bb5..1df627f8`, dispatched concurrently and blind to each
other: the GPT side (`gpt-5.6-sol` / medium, OS-process, `--sandbox read-only`) and a COLD Opus side
on `harness:gate-reviewer`, which has no Edit/Write/Bash **by construction** rather than by
instruction.

**The two briefs are deliberately not identical, and the asymmetry is recorded so the reconciliation
can be read correctly.** Both got the same contract, the same attack list and the same hub rulings.
The GPT side was additionally handed the four already-declared open items (the `pol` token, the
dormant duplicate reasons, the unprovable fail-closed branch, the un-run governance lane) with an
instruction to RULE on the chip's handling rather than re-report them as discoveries. The Opus side
was given no such list: it has to find them in the pack itself and judge them cold. So the GPT side
is the tighter check on the contract, and the Opus side is the better check on whether the pack's
own declarations survive a reader who was not told where to look. A disagreement between them about
a declared item is therefore informative, not noise.

Neither reviewer can run the governance lane either — same worktree constraint, and hub ruling R-b
has since made that rung hub-owned outright — so C11's governance rung stays OPEN across the gate
regardless of the verdicts, by ruling rather than by omission.

(verdicts + reconciliation — pending)

## Must-fail proofs (R5)

Every guard this chip touched, with the mutation that must redden it. Each was run and observed, not
reasoned about.

| Guard | Mutation applied | Observed |
|---|---|---|
| C6 unit canonicalisation | remove the mm/cm/m/pol scaling so tokens compare lexically | golden `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm` goes RED (`500MM` vs `50cm` reads as a contradiction again) |
| C6 contradiction still fires | — (inverse check, no mutation) | `50cm` vs `40CM` still REJECTs, so the canonicalisation did not simply disable the rule |
| C8 independent limits | restore the shared `limit*5` derivation | the 29-link fixture returns fewer than 29 → RED |
| C8 defaults independent of `Limit` | `linkLimit = candidateLimit*100`, `auditLimit = candidateLimit*500` | new test RED on the `5` and `500` rows, `20` row still green, **and both pre-existing default tests stay GREEN** — proving the old net was blind to the class |
| F-02 display offsets | (proved by the reviewer, not the author) rebuilt the pre-corrective function standalone | returned `display="50c"`, truncated, for `"İnox 50cm"` — the guard is load-bearing |

**One guard has NO must-fail proof, and this is stated rather than faked:** the fail-closed parse
fallback in `hardNegativeDimension` is unreachable by construction with the five-token alternation,
so no input reddens it. Both the implementing worker and the F-02 re-reviewer reached the same
conclusion independently. R5 asks for a proof per guard TOUCHED; the honest report is that this
branch cannot be exercised, not that it was exercised.
