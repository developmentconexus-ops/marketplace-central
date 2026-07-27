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
          RESOLVES AT: main@f81b8975 (`git show f81b8975:<path>`), NOT on this branch — see
          "Note on the pack location" below. The round-3 GPT gate reported the path as absent from
          `git ls-tree ac72eb82`, which is correct and is what the note already says; the pointer is
          added here so the header line is self-sufficient and a reader need not find the note first.
EXEMPLO-IO: MLB4735326915 · listing "Toalheiro Simples Soul Zen 50cm Cromado" (cm)
            vs ERP product "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO" (mm) → CONFIRM/MEDIA/70
            asserted by TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm
            (run it: go test -run TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm ./...)
            NOTE: the dispatch prompt states this pair the other way round — the 500MM name as the
            LISTING title. Adjudication A3 records the prompt's orientation as the error; the ERP
            product is the all-caps 500MM one. This header line still carried the inverted form
            after A3 corrected the C6 cell; fixed in ROUND 3, flagged by the cold P6 gate.

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
| D14 | P6 | Gate reviewer, GPT side of the dual gate | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol.md` | `scratchpad/agent__p6-sol.last.md` + `.log` | completed — **FAIL** (2 BLOCKING; 1 of them a stale-sha artifact of the chip's own dispatch) |
| D15 | P6 | Gate reviewer, COLD Opus side | Opus, `harness:gate-reviewer` (no Edit/Write/Bash by construction) | Agent tool, async | inline brief | `tasks/a93b3ffaca9b176b7.output` | completed — **FAIL** (2 BLOCKING, incl. the false provenance claim; 1 SHOULD-FIX, 3 NIT) |
| D16 | P3 | Implementer S7 (R5 coverage for 3 reachable guards + clamped-formula row) | gpt-5.6-luna / high | OS-process codex | `scratchpad/prompt-s7.md` | `scratchpad/agent__s7-impl.last.md` + `.log` | GREEN; worker **could not commit** (sandbox `index.lock: Permission denied`, profile §3 class) — chip verified and committed `e1195806` |
| D13 | P4 | Re-reviewer, feature F-01 over `f92ca9c7` | Claude sonnet subagent | Agent tool, async | inline brief | `tasks/ad9c0421c9c2a538c.output` | completed — **PASS**, zero findings |
| D17 | P3 | Implementer S8 (hub rulings R-6 dedup/precedence + R-7 parser seam) | gpt-5.6-sol / low (complex slice) | OS-process codex | `scratchpad/prompt-s8.md` | `scratchpad/agent__s8-impl.last.md` + `.log` | committed `d9952509`; **R-7 accepted, R-6 REGRESSES the TitleMatch path** — chip-verified by differential probe, `ESCALATION` sent, corrective pending the hub ruling |
| D18 | P3 | Implementer S9 (hub ruling R-6a: UNAVAILABLE exclusive, FOR+AGAINST coexist) | gpt-5.6-sol / low (complex slice) | OS-process codex | `scratchpad/prompt-s9.md` | `scratchpad/agent__s9-impl.last.md` + `.log` | committed `2921d563`; matrix reported total, no uncovered case; **chip-verified** — regression closed by own probe, rule-2 must-fail re-run by hand, L0+L1 green at tip |
| D19 | P6 | Dual gate round 2 — COLD Opus side (`harness:gate-reviewer`, Read/Grep/Glob only, no Bash) | Opus / cold | Agent tool, async | `scratchpad/prompt-p6-opus-r2.md` | verdict transcribed below | completed — **FAIL** on the pack, not on the behaviour: "the code is in materially better shape than the pack" |
| D20 | P6 | Dual gate round 2 — GPT side | gpt-5.6-sol / medium | OS-process codex | `scratchpad/prompt-p6-sol-r2.md` | `scratchpad/agent__p6-sol-r2.last.md` + `.log` | completed — **FAIL**; C6/C8/C11 FAIL, 3 BLOCKING, all three about evidence decay, none about code |
| D22 | P6 | Dual gate round 3 — COLD Opus side (`harness:gate-reviewer`, Read/Grep/Glob only, no Bash) | Opus / cold | Agent tool, async | `scratchpad/prompt-p6-opus-r3.md` | verdict transcribed in ROUND-4 below | completed — **FAIL** on the pack; found the ROUND-3 closure sentence overclaiming its own scope |
| D23 | P6 | Dual gate round 3 — GPT side | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol-r3.md` | `scratchpad/agent__p6-sol-r3.last.md` + `.log` | completed — **FAIL**; C5 FAIL, 3 BLOCKING; C1–C4, C6–C12 PASS, all re-anchoring samples resolved |
| D24 | P6 | Mechanical citation audit — the D21 remedy executed, not restated | chip, `cite-audit.py` (deterministic, no model) | local script, read-only | `_chip-anchors/cite-audit.py` (committed pack artifact per R-11) | `_chip-anchors/cite-audit.txt` — the artifact carries its own COVERAGE block (resolved / unresolvable); this row does not restate the number, per R-14 | completed, then **SUPERSEDED** — round 5 showed the coverage total counted unresolved rows as resolved and whole citation forms fell outside the regex. Replaced by `cite-table.py` + `coordinates.txt`, which round 6 then put on the ladder as rung PACK and R-24 took back off — it is a report, not a gate. Both round-4 files stay in the pack, each carrying a SUPERSEDED banner in its own first lines |
| D25 | P6 | Dual gate round 4 — both sides re-pinned to rulings R-6a, R-9, R-10, R-11, R-12 @ `main 371c91d` | COLD Opus + gpt-5.6-sol / medium | Agent tool + OS-process codex | `scratchpad/prompt-p6-opus-r4.md`, `scratchpad/prompt-p6-sol-r4.md` | transcribed / `scratchpad/agent__p6-sol-r4.last.md` | see ROUND-4 below |
| D21 | §11 | Independent adversarial judgement — structural vs informational (briefed AGAINST the abstraction; "change no rule" set as the default verdict) | Claude sonnet subagent | Agent tool, async | inline brief, quoted in ROUND-3 | `scratchpad/agent__p11-judge-r3.md` | completed — **REJECT-IN-FAVOUR-OF-mechanical-table-regeneration-at-final-tip**; the chip's own proposal was refused |
| D26 | P6 | Dual gate round 5 — COLD Opus side. **Area: SEMANTICS** per R-18 (the two sides are given DIFFERENT areas, because agreement between gates looking at the same thing is worth less than disagreement between gates looking at different things). Does the PROSE match the line the audit printed, and is the behaviour in the production code | Opus, `harness:gate-reviewer` (Read/Grep/Glob only, no Bash by construction) | Agent tool, async | `scratchpad/prompt-p6-opus-r5.md` | verdict transcribed below by the chip — the Agent tool's own output file is EMPTY, so no artifact exists to point at | completed — **FAIL** on the pack, code cleared |
| D27 | P6 | Dual gate round 5 — GPT side. **Area: DERIVATION** per R-18. Re-derive every count and listing from `git` at the frozen SHAs, re-run the audit, attack the R-14 pointer discipline | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol-r5.md` | `scratchpad/agent__p6-sol-r5.last.md` + `.log` | completed — **FAIL**, four BLOCKING, all structural: false coverage total, decayed unpinned counts, a wrong count at a frozen SHA, and a listing whose command binds no tree-ish |
| D28 | P6 | Dual gate round 6 — COLD Opus side. **Area: SEMANTICS.** The prose was rewritten wholesale under R-22 and now claims BEHAVIOUR anchored to test and symbol names; the question is whether each named test asserts what the sentence says, and whether the behaviour is in the production code | Opus, `harness:gate-reviewer` (Read/Grep/Glob only, no Bash by construction) | Agent tool, async | `scratchpad/prompt-p6-opus-r6.md` | `_chip-anchors/verdict-r6-cold.md` — transcribed by the chip on arrival, because the Agent tool leaves an empty output file (D30) | completed — **PASS**. No BLOCKING semantic finding: no test cited as proving what it does not assert, no guard described as covering what it does not cover. Sampled ~25 table rows plus all 20 backticked `Test*` names; every one resolved to the stated file AND line. One SHOULD-FIX (the tool's silent anchor drops), one NIT (the U2 inversion) |
| D29 | P6 | Dual gate round 6 — GPT side. **Area: DERIVATION AND THE APPARATUS.** Re-derive every frozen-SHA fact; then turn on the tool that now carries the coverage claim — run the PACK rung, re-run every recorded mutation (four that must FAIL, one that must NOT fire, and one proving tolerance did not cost detection), and try to find text the checks cannot see | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol-r6.md` | `scratchpad/agent__p6-sol-r6.last.md` + `.log`, transcribed into `_chip-anchors/verdict-r6-gpt.md` | completed — **FAIL**, five BLOCKING. Two were pack errors (a frozen-SHA grep sentence omitting its test exclusion; a write-set range pinned two commits past the slices it names, whose count matched anyway by coincidence); three were the tool's own coverage claims being total in the wording and partial in the code. Every frozen-SHA count in the pack re-derived correctly, including the round-5 corrections |
| D30 | P6 | Cold verdict transcription — the remedy for the artifact asymmetry found in round 5. Not a dispatch: a chip step the ledger names so it can be checked | chip | direct write | — | `_chip-anchors/verdict-r6-cold.md` | completed on arrival. R-23 then added a second holder: the same text travelled to the hub VERBATIM in the event, and the hub commits its own copy. The GPT verdict was given the same treatment for the same reason — its codex `-o` artifact is real but lives in a session-scoped scratchpad that dies exactly when a transcript does |
| D31 | P6 | Round-6 corrective under R-24. Two prose corrections; the tool reports what it could not resolve and stops claiming totality, WITHOUT widening any recognizer; the PACK rung leaves the ladder; R-23 and R-24 recorded; the U2 inversion fixed chip-side. ZERO production code — that is the property the hub verifies before any AGREEMENT, not one it takes on trust | chip | direct edit | hub ruling R-24, `main@874d00e5` | this pack + `cite-table.py` + `coordinates.txt` | completed — corrective tip `f1397cf7`; `git diff --name-only 85b6c367 f1397cf7 -- apps contracts packages` returns **nothing** |
| D32 | P6 | **Narrow re-check of a delivered FAIL — not a seventh round.** The gate that filed the five blockers re-checks ITS OWN five at the corrected tip, scoped to exactly what it filed. There is no new surface, because the corrective REMOVES claims rather than adding them. The cold gate does NOT re-run: it passed, its area is untouched, and re-running a passing gate over an unchanged area is the redundancy R-18 devalued | gpt-5.6-sol / medium | OS-process codex, `--sandbox read-only` | `scratchpad/prompt-p6-sol-recheck.md` | `scratchpad/agent__p6-sol-recheck.last.md` + `.log`, transcribed verbatim into `_chip-anchors/verdict-r6-recheck.md` because the codex `-o` artifact lives in a scratchpad that dies with the session | completed — **FAIL**. Three of five closed: two RESOLVED on re-derived facts, one DISSOLVED-BY-R-24 and accepted as closure by the gate that filed it. Two STAND, for one reason: the corrective removed the totality claims from what the tool PRINTS and left them in what the pack and the implementation SAY. No new defect introduced by the corrective; the gate states the survivors predate it. Sent to the hub raw, no AGREEMENT string authored, the three lines deliberately NOT patched while §11 is a live question |
| D33 | P6 | **Discharge of the two survivors under R-25.** Not a round and not a dispatch: a remedy fully determined by its finding — two files, exact strings, replacement target already certified accurate by the gate. *When a remedy is fully determined by its finding, verification is a check of ABSENCE — one bit — and that wants a shell, not a second model.* The comment DELETED rather than reworded, so no rewrite could introduce a fresh falsehood; the two sentences re-aimed at what the report prints; the disclosure paragraph removed with the falsehoods it disclosed | chip, under hub ruling R-25 | direct edit | hub ruling R-25, `main@4e9c5ca4` | this pack + `cite-table.py` + `verdict-r6-recheck.md` | completed — the hub verifies absence itself at the final tip, **by string and not by line**, having found the gate's own coordinate two lines off from the text it cited |

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
| Write-set exactness | the slice commit's own `--stat`, read after committing | exactly the 3 declared files |
| L0 build | `go build ./...` | exit 0 |
| L0 vet | `go vet ./...` (full tree) | exit 0 |
| L1 tests | `go test ./internal/modules/product_links/...` | all `ok`, exit 0 |
| L0 format, authored files only | `gofmt -l <the 3 files>` | empty |

Chip ruling A1 verified applied: `hardNegativeDimension` returns THREE values — a canonical
comparison key, a display string, and the found flag — and `detectHardNegative` formats the operator
message from the DISPLAY values, not the canonical ones. So the comparison happens in millimetres
while the queue still reads `50cm≠500mm`, which is what an operator can act on; a canonical `mm:...`
form in the motivo would have been technically correct and operationally useless.

## C12 — proven under the hub's revised proof form (ruling of 2026-07-26, contract @ `c9fecbc`)

`gofmt -w internal/composition/root.go` applied and committed as `633bf9fa`, per the hub ruling.

| Proof | Command | Result |
|---|---|---|
| additive-only in TOKENS | `git diff -w 917f7bb5 633bf9fa -- .../root.go \| grep '^-'` | **no removed lines** |
| file is gofmt-clean | `gofmt -l internal/composition/root.go` | empty |
| region diff quoted raw | `git diff --unified=2 917f7bb5 633bf9fa -- .../root.go` | quoted below, realignment visible |

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
hidden: a whitespace-insensitive diff collapses them to nothing, which is the token-level test the criterion now
uses.

## Slice S2 — chip verification (not the worker's word)

| Check | Command | Result |
|---|---|---|
| Write-set exactness | `git diff --name-only 917f7bb5 b9da6d2e` | 13 files, all declared across S1+S2, zero undeclared |
| L0 build | `go build ./...` | exit 0 |
| L0 vet | `go vet ./...` (full tree) | exit 0 |
| L1 tests | `go test ./internal/modules/product_links/... ./internal/modules/connectors/... ./internal/composition` | all `ok`, exit 0 |
| C2 symbol | `grep -rn mandatoryUnavailableReasons .../application` | **0 hits** |
| C2 anchor literals | `grep -rn '"marca"\|"refforn"' .../application` (production only) | **0 hits** |
| C3 provider branch | `grep -rn 'mercado_livre\|ProviderCode ==' product_links` (production only) | **0 hits** |
| C12 additive-only | `git diff --unified=0 917f7bb5 633bf9fa -- root.go \| grep '^-'` | **no removed lines** |

**The write-set row above was re-pinned in the round-6 corrective, and the reason is worth more than
the fix.** It previously ended the range at `633bf9fa`, which is not where S2 ends: that range walks
S1 `5bc55219`, S2 `b9da6d2e`, S3 `11e68f6f`, and then the composition-formatting commit. It
nevertheless returned 13 paths — the same 13 — because S3 only touched paths S1 and S2 had already
written. So the number was right while the range was wrong, and nothing in the output could say so.
A count that is correct for the wrong reason does not announce itself; it reads exactly like a
correct one, which is why the round-6 derivation gate filed it as blocking and why the hub verified
it before ruling. Both endpoints are stated here rather than the wrong one being quietly swapped
out, because the coincidence is the finding.

The worker again reported plain `go build ./...` blocked by VCS stamping; again it does not
reproduce chip-side (exit 0). Profile §3 (2026-07-17, CHIP-M03 F-ENV) already ratifies that the
chip's re-run outside the codex sandbox is the verification of record for exactly this class of
sandbox false alarm.

Substance spot-checks on the diff:

- **One finalizer, six branches.** Every scoring branch routes its reasons through
  `appendProviderDeclaredUnavailableReasons` — six call sites in the generation service, verifiable as
  a count with `git grep -c appendProviderDeclaredUnavailableReasons 85b6c367 -- <that file>` (seven,
  the seventh being the declaration itself). `mandatoryUnavailableReasons` is GONE, not bypassed.
- **R1 honoured.** The provider-level detail is `fmt.Sprintf("provider não fornece a âncora %s", …)`.
  The distinction lives in `Detail`; the direction enum is untouched, so no 4th value exists.
- **The failure names the provider** — `unavailableIdentityAnchorsError` wraps the cause as
  `PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE for provider %q: %w`, chip ruling A2 applied.
  `%w` matters: the cause survives `errors.Is`, so a caller can still distinguish an adapter outage
  from a missing declaration.
- **Resolution happens once per distinct trimmed provider code**, in `resolveIdentityAnchors`, BEFORE
  any scoring or persistence. All three failure modes inside that function — empty provider code, an
  adapter error, and a nil declaration — return the error through that one wrapper rather than an
  empty declaration, which is ADR-17: an unresolved declaration must never decay into "the provider
  supplies none". The FOURTH mode, an unknown anchor, is not rejected here at all; it is rejected
  upstream in `MarketplaceCapabilityService`, where an anchor absent from `KnownIdentityAnchors` and a
  duplicate anchor both fail closed as `unsupported`. Named by test:
  `TestGenerateLinkCandidatesRefusesWithoutIdentityAnchorReader` for the reader guard,
  `TestMarketplaceCapabilityServiceRejectsUnknownIdentityAnchor` and
  `TestMarketplaceCapabilityServiceRejectsDuplicateIdentityAnchor` for the upstream pair.
  **ROUND-4 correction, cold gate:** this list previously cited four bare line numbers with no file,
  one of which was the function's SUCCESS return — so the sentence offered `return resolved, nil` as
  proof of failure handling, while the real empty-code return went uncited. Neither the code nor the
  conclusion was wrong; the citations were, and a bare number inherits whichever file the reader last
  saw named. That is why R-17 banned the form and why coordinates are out of the prose entirely now.
- **`GenerateLinkCandidates` fails closed when the reader is absent** — the pre-existing
  `PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED` guard gained `s.identityAnchors == nil` as a fourth
  disjunct, so a composition root that forgets the wiring gets an error instead of candidates scored
  as though no provider supplied anything. This is the guard whose must-fail proof was missing through
  round 4 and is proven below.

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

```quoted-contract
C12 | Grant de `root.go` foi **additive-only em tokens**, e o arquivo fica `gofmt`-limpo |
`git diff -w <base>..HEAD -- internal/composition/root.go` **sem nenhuma linha removida**
(o `-w` ignora o realinhamento de whitespace) + `gofmt -l internal/composition/root.go` sem
saída + diff da região citado no payload do CLOSED
```

Quoted verbatim as the hub wrote it, inside a fence that declares it as quoted — the criterion's
wording is not the chip's to edit. The chip DISCHARGES it against a named tree-ish rather than a
moving ref: every C12 run in this pack names its commit on both ends.

So C12 is now proven by a whitespace-insensitive diff (token-level, which is what "additive-only" always meant
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
| Write-set exactness | `git diff --name-only 917f7bb5 5bc55219` | exactly the 8 declared files, zero undeclared |
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
- `MarketplaceCapabilityService.Capability` — nil declaration ⇒
  `unsupported(providerCode, "identity anchors")`, NOT an empty slice. This is the ADR-17 hinge of
  the whole feature: "nobody told us" must not decay into "the provider supplies none". Unknown
  and duplicate anchors also error; the return is a clone.
- `connectors/adapters/mercado_livre/capability_adapter.go` — declares exactly
  `seller_sku`, `ean`, `title`. The `-` lines in this file's diff are gofmt field realignment
  from the added field, not logic edits. (C12's additive-only constraint binds `root.go`, which
  this slice does not touch.)
- `product_links/adapters/connectors/identity_anchor_adapter.go` — projects the complete
  vocabulary, holds `*MarketplaceCapabilityService` exactly as the cited precedent
  `inventory/adapters/connectors/stock_writer.go` does, carries the compile-time assertion
  `var _ productlinksports.ProviderIdentityAnchorReader = IdentityAnchorAdapter{}`, and contains
  no provider-name literal and no `providerCode ==` comparison (R2 / C3).
- Production `ProviderCapabilitySet()` declaration sweep: one site repo-wide, the ML adapter.

## Chip adjudications on the P2 plan

The planner's five decisions (Q1–Q5) are accepted as the design, with four corrections the chip
made as slice-card authority before dispatch. Recorded here because they change what the
implementers build, not just how it is phrased.

**A1 — canonical dimension form is for COMPARISON ONLY; the operator-facing detail keeps the
original tokens.** The planner proposed a reduced-rational signature `mm:<num>/<den>` as the
value compared inside `hardNegativeDimension`. But the signature IS the message —
`fmt.Sprintf("hard-negative: medida/dimensão divergente %s≠%s", …)`, and at BASE `917f7bb5` BOTH
arguments were the compared keys. *(A1 is the reason `detectHardNegative` now passes the display
tokens instead. The claim is about the state that MOTIVATED the ruling, and is labelled as such —
readable at base with `git show 917f7bb5:<that file>`, not at the reviewed tip.)*
Shipping the plan as written would turn the live queue message `50cm≠40cm` — which M-05's own
evidence pack records verbatim off the real account — into `mm:500/1≠mm:400/1`. That is a readability
regression on an operator-facing string, introduced by a fix meant to help the operator. S3 must
therefore return BOTH a canonical comparison key and the original display tokens, and its C6 test
asserts the detail for `50cm` vs `40cm` still contains `50cm` and `40cm` and does NOT contain
`mm:`.

**A2 — the identity-anchor lookup failure names the provider code.** `PRODUCT_LINKS_PROVIDER_
IDENTITY_ANCHORS_UNAVAILABLE` with an unnamed provider is an unactionable error on a run that
touches N providers. ADR-17 honest-failure means the failure says which provider it could not
resolve.

**A3 — the day-1 golden test uses the REAL M-05 pair, not the chip prompt's wording.** The
dispatch prompt's EXEMPLO-IO has the listing and the ERP product swapped. M-05's evidence pack records the live pair: listing `Toalheiro Simples Soul Zen 50cm …` (cm)
against ERP product `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` (mm), CODPROD 33698 exact SKU
match, listing carries NO EAN. So the post-fix expectation is **CONFIRM / confidence 70 / band
MEDIA** — the third CONFIRM M-05 predicted — and NOT `ALTA/ACCEPT`, which the planner's synthetic
concordant-CODPROD+EAN fixture produces. S3 carries both cases: the real golden (CONFIRM) and the
synthetic concordant (ACCEPT).

**A4 — Q5 accepted, with the residual stated rather than hidden.** Independent defaults
(candidates 20 / links 2000 / audits 10000) make all 29 links reach `/vinculos` with no wire
change, and leave the two `Limit: 2000` call sites in `mutations/adapters/productlinks/writer.go` compiling
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
the commit's own `--stat` reports exactly 2 insertions / 2 deletions.

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
only inch-family subtest is KEYED `"inches"` but its DATA is
`"Produto 2pol"` / `"Produto 50.8mm"`, i.e. it exercises `pol`. Dispatched as S5 part A, together
with fail-CLOSED parse fallback (F2 direction), matching on `original` with `(?i)` to kill the
offset aliasing, and key/display pair compaction.

**A7 — F-01 layering leak (SHOULD-FIX) accepted; the redundant validation is removed.**
**Cited at the state that MOTIVATED the ruling, not at the reviewed tip — the ruling's whole point is
that this no longer exists.** At `f92ca9c7^`, `generation_service.go` imported `connectorsports` into the
product_links APPLICATION layer and `resolveIdentityAnchors` re-validated anchor names against
`connectorsports.KnownIdentityAnchors()`. At the reviewed tip both are GONE from PRODUCTION, and the
proof is an absence with a control:
`git grep -c connectorsports 85b6c367 -- <the application package> ':!*_test.go'` returns nothing,
while the same command at `f92ca9c7^` returns hits — so the emptiness is the fix, not a mistyped
pattern.

**That sentence used to omit the test-file exclusion, and the omission made it false.** Without
`':!*_test.go'` the command reports two hits in `generation_integration_test.go` and exits 0: the integration
fixture legitimately names `connectorsports.IdentityAnchorSellerSKU` to build its input, which is a
test constructing a value, not the application layer reaching for the connectors vocabulary. The
pack disclosed those two occurrences in its own import table further down while this paragraph
claimed the grep was empty, so the document contradicted itself — a reader believing either half
was misled by the other. The round-6 derivation gate filed it, the hub re-ran it before ruling, and
the fix is to state the command that was actually run rather than the tidier one.

`product_links/ports/provider_identity_anchor_reader.go` types its
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
implementation is correct today — three literal constants guarded independently inside
`ListLinkWorkflows`, verified by reading it; this is a hole in the regression net, not a live defect. Queued
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
| A1 — token set back to base | read `hardNegativeDimensionPattern` | `(?i)\d+\s*x\s*\d+(?:\s*x\s*\d+)?\|\d+(?:[.,]\d+)?\s*(?:mm\|cm\|pol\|")\|\d+(?:[.,]\d+)?\s*m\b` — identical to `917f7bb5` apart from the `(?i)` prefix |
| A1 — candidate list | read the suffix loop in `normalizeDimensionToken` | `{"pol", "mm", "cm", `"`, "m"}`, longest-suffix-first preserved so `m` cannot shadow `mm`/`cm` |
| A2 — fail CLOSED | read the `SetString` branch | on parse failure appends `{key: token, display: originalToken}` and `continue`s — the side's other tokens survive; the old `return "", "", false` is gone |
| A3 — no offset aliasing | read the match loop | matches `original` directly; `originalToken` sliced with offsets native to that same string; lowering is per-token, not whole-string |
| A4 — key/display correspondence | read the tail | one `dimensionPair` per token, `slices.SortFunc` + `slices.CompactFunc` both keyed on `key`, then the two slices are projected from the SAME compacted pairs — they cannot diverge |
| B1 — layering leak closed | `grep -rn 'connectors/ports' product_links/application/*.go` | ~~only `generation_integration_test.go`, a PRE-EXISTING test-only import; zero hits in production code~~ — **THIS ROW WAS FALSE. See the correction below.** The production half (zero hits) holds; the provenance half does not |
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
- ~~the `connectors/ports` grep is not zero directory-wide because a pre-existing test file imports
  it; changing that file was outside the frozen write set. Correct call.~~ **FALSE — corrected
  immediately below.**

**A second defect in the two A1 rows above, found by the R-22 tool during this rewrite and by nothing
else.** Both rows used to read "read `generation_service_test.go:<line>`" while quoting content that
lives in `generation_service.go` — the production file. The quoted content was correct in both rows;
the FILE named was wrong in both. No gate reached this in five rounds, and no citation audit could
have: every one of those audits asked whether a coordinate RESOLVES, and both of these resolved
cleanly — to a real line, in a real file, holding something else entirely. That is the strongest
argument in this pack for R-22's actual mechanism. A coordinate is checkable for existence but not
for MEANING; a symbol name carries its own meaning, so `hardNegativeDimensionPattern` cannot be
silently attributed to the test file, because there is no such symbol there to resolve to.

### CORRECTION — the "pre-existing import" claim above was false, and I authored it

Found by the cold Opus P6 gate, re-verified by me against base. Recorded here as a correction rather
than fixed by editing the original rows, because a pack that quietly repairs its own false
statements is worth less than one that shows them.

What the two rows above claimed: the `connectors/ports` import in `generation_integration_test.go` was
pre-existing, and the file was outside the frozen write set.

What is actually true, by tool:

| Check | Result |
|---|---|
| `git show 917f7bb5:…/generation_integration_test.go` imports | `internal_read` domain+ports, `product_links` adapters/postgres + domain, `testsupport/postgres`. **Neither connectors package.** |
| tip imports | adds BOTH `connectorsapp` and `connectorsports`, and uses them — `NewMarketplaceCapabilityService` with an `IdentityAnchor` slice |
| `git diff --stat 917f7bb5 85b6c367 -- <that file>` | non-empty on both axes — **the file IS in the chip write set** |

(The insertion count this row used to state was 10, measured when the row was written. At the final
code tip the same command reports 17, because a later slice touched the file again. Nothing was wrong
with the measurement and the conclusion never moved; the NUMBER decayed, which is the mutable-axis
defect in miniature, found in the pack's own self-correction. The predicate — non-empty, therefore in
the write set — is what the row needed, and it is stable.)

So **this chip introduced the connectors vocabulary into `product_links/application`**, in a
build-tagged integration test, and the pack then explained away the resulting grep hit by
attributing it to someone else. The row sat under a heading that reads "Verified independently of
every claim it made." It was not: the worker's claim was promoted to chip-verified without being
checked against base. That is the failure mode this chip's whole method exists to prevent, and it
is the CHIP-IMPORT-FIX class already in the mission ledger (self-exculpating prose in an evidence
pack).

**The code is fine; only the justification was false.** An integration test that wires the real
`IdentityAnchorAdapter` to the real `MarketplaceCapabilityService` genuinely needs both imports —
the P6 gate that caught the false claim says so explicitly, calling the shape defensible. The
production-code half of the B1 claim also still holds: zero `connectors/ports` hits in
`product_links/application` production code.

**The full truth, which neither the worker nor I stated correctly:** there IS a genuinely
pre-existing connectors import in that package's PRODUCTION code — `import_service.go` imports
`connectors/domain`, present at base and untouched. It is simply not the import the B1 row cited.
Verified at both ends — the frozen tip and the base commit — output quoted:

```quoted-output
TIP:  generation_integration_test.go:10,11  (connectors/application, connectors/ports)  <- CHIP-ADDED
      import_service.go:9                   (connectors/domain)                          <- pre-existing
      import_service_test.go:8              (connectors/domain)                          <- pre-existing
BASE: import_service.go:9, import_service_test.go:8 only
```

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

## Slice S7 — chip verification (P6 corrective, `e1195806`)

Test-only slice closing the R5 coverage gap both gates converged on. Verified by the chip, not
accepted on the worker's report.

| Claim | Chip's check | Result |
|---|---|---|
| tests green | `go test -count=1 ./internal/modules/product_links/... ./internal/modules/connectors/...` re-run by the chip | all `ok`, zero failures |
| production untouched | a `--stat` at the frozen tip over the three production files the mutations touched | **empty** |
| formatting | CR-stripped `gofmt -l` per file (the only informative form on this checkout — ROUND-1) | `[]` for all three |
| the 1b must-fail actually reddens | chip re-applied the mutation and ran the test itself | RED, output pasted in the must-fail table |

**Two things the worker got right that are worth recording, because both were live traps.**

First, the stub. The card told it to verify that `stubIdentityAnchorReader` really returns
`(nil, nil)` for an unregistered provider code before relying on that, and to extend the stub if it
did not. It checked and reported "stub already returned `(nil, nil)`; no stub change was needed" —
which is the difference between a test that drives the nil-declaration guard and a test that quietly
drives some other path.

Second, it caught its own EOL damage. Its must-fail mutations left `\r` differences in two
production files even after the semantic revert, and it normalised them back to the exact `HEAD`
blobs before staging rather than declaring victory on an empty diff. Worth being precise about
why that was diligence rather than necessity: with `core.autocrlf=true` those worktree-only EOL
differences would have been normalised on commit anyway, so the blob was never at risk. But the
worker could not know that from inside its sandbox, and cleaning up is the right instinct — the
ROUND-1 lesson on this repo is exactly that EOL state hides signal.

**The worker could not commit.** `fatal: Unable to create …\.git\worktrees\chip-anchors\index.lock:
Permission denied` — the codex workspace-write sandbox cannot write the worktree's git index. This
is the profile §3 sandbox class, now recorded a fifth time on this chip. The chip verified the
working tree and committed on the worker's behalf as `e1195806`; the ledger row (D16) says so
rather than implying the worker committed.

## Slice S8 — chip verification (hub rulings R-6 + R-7, `d9952509`)

Implemented `gpt-5.6-sol`/low. Verified by the chip, and the verification found a regression the
worker's own report did not — recorded here in full because the failure mode is instructive.

### R-7 — accepted

`normalizeDimensionToken` extracted from `hardNegativeDimension`, which now calls it per token. The extraction is behaviour-preserving by inspection: the moved body
is character-identical apart from returning `(key, parsed)` instead of appending to `pairs`, and the
`x`-token early return keeps its position. `TestNormalizeDimensionTokenFailsClosed` calls the helper
directly with `"abcmm"` — a token the alternation cannot emit — and asserts the raw normalised token
is KEPT as the signature key. Must-fail observed: changing the fallback to drop the key returns
`key=""` against expected `"abcmm"` → RED. That is the seam R-7 asked for, and it now exists where
none was possible before.

The public-path half is still unreachable, and the worker said so instead of contriving it: two
titles carrying different UNPARSEABLE tokens cannot be built through the current regexp. Correct
call — the card explicitly permitted reporting that rather than manufacturing a case.

### R-6 — implemented, and it REGRESSES a path the ruling did not consider

The dedup itself is real: `appendProviderDeclaredUnavailableReasons` now indexes reasons by anchor
and collapses duplicates, seed positions are preserved, new declaration anchors append in
declaration order, and the ordering test survives `-count=5`. The `seller_sku`/`ean` contradiction
the hub ruled on is gone.

But "at most one motivo per anchor", applied literally, also fires on a pair that is not
contradictory at all. In `applySingleAnchorScore` under state `TitleMatch`, the seed
carries `{title, FOR}` and the hard-negative branch in the same function then appends
`{title, AGAINST}`. Same anchor, so the new dedup drops one — and since it keeps the first unless
that first is `UNAVAILABLE`, the one it drops is the **AGAINST**. Losing the FOR would be a scoring
bug; losing the AGAINST is worse, because it removes the reason an operator needs to see.

**Differential proof — chip-run, same probe and fixture, only `generation_service.go` swapped
between the two commits.** Listing title `"Kit 3 Toalheiro Simples 50cm Cromado"` against ERP
`"SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO"`, state `TitleMatch`:

```
pre-S8  (08308afb)  status=REJECT confidence=25
  title       FOR         match por título (ranking-only, nunca ACCEPT)
  seller_sku  UNAVAILABLE seller_sku sem correspondência
  ean         UNAVAILABLE ean sem correspondência
  title       AGAINST     hard-negative: kit/combo divergente entre título do anúncio e produto interno
  marca       UNAVAILABLE provider não fornece a âncora marca
  refforn     UNAVAILABLE provider não fornece a âncora refforn

post-S8 (d9952509)  status=REJECT confidence=25
  ...identical, MINUS the `title AGAINST` line.
```quoted-output

The candidate is still rejected and still scored 25 — the verdict is right. What is lost is the
REASON: an operator reading `/vinculos` sees a rejection whose only `title` reason is FOR, with no
contradiction stated anywhere. That is the ADR-17 shape this chip exists to remove, reintroduced by
the fix for a different instance of it.

Reachable in ordinary use: a title-only match that also trips kit/combo/cor/voltagem/dimensão. **No
existing test caught it** because every hard-negative test drives a concordant SKU+EAN state, where
`title` is not in the seed and so never collides — the suite was green through the regression.

**How the chip found it, since the method is the transferable part.** The gate briefs for round 2
had already been written, and one of the attacks they instruct is "does the fix silently DROP a
reason that should have survived? Losing that would quietly undo F-01 while looking like a dedup
fix." Writing that attack down is what made it obvious to run it against the chip's own tip before
dispatching anyone. The probe was then run twice against the two commits rather than reasoned about,
which is what turned a suspicion into a differential.

**One vacuous-probe near-miss, recorded because it nearly produced a false clean.** The first probe
run PASSED — with no output, because `buildCandidatesFromProducts` returned zero candidates
(`canonicalProductID` requires a non-nil `InternalProductID`, which the fixture lacked). A test that
asserts over an empty slice asserts nothing and reports success. It was caught only by adding an
explicit `if len(candidates) == 0 { t.Fatalf }` guard. Same family as the confounded-assertion class
already recorded on this chip: the assertion has to be able to fail before its passing means
anything.

### THE FINDING — the hard-negative suite could not see the TitleMatch path

Recorded as the finding in its own right, at the hub's direction, because the code defect was the
smaller half. The larger half is a COVERAGE hole that predates this chip:

**Written in the PAST TENSE deliberately — this describes the suite AS S8 FOUND IT, and the hole is
closed at tip.** The round-4 cold gate raised the tense as a NIT: at `2921d563`
`TestTitleMatchHardNegativeKeepsTitleForAndAgainstInSeedOrder`
drives exactly this state. It was already disclosed two paragraphs
below; stating it here removes the need to read on.

**Every hard-negative test in this package DROVE a state where `title` was not in the reason seed.**
The suite exercised `buildConcordantCandidate` (SKU and EAN concordant) and the ExactSKU / ExactEAN
branches of `applySingleAnchorScore`. In all of those the seed names `seller_sku` and `ean`, and the
hard-negative branch appends `{title, AGAINST}` — a fresh anchor, no collision. The one state whose
seed already contains `title` is `TitleMatch`, and before this chip **no** hard-negative test drove
it.

So the whole class of "what happens when the hard-negative anchor is ALREADY in the seed" was
untested before this chip touched anything. That is why a change to reason finalisation could delete
an operator-facing reason on a REJECT and leave the suite green. The defect S8 introduced was
detectable in principle by reading the code; it was undetectable by running the tests, and the
second fact is the one worth carrying forward.

What closes it is the regression test mandated by R-6a — `title` FOR and `title` AGAINST surviving
together on a REJECT — not merely the precedence fix. A fix without that test would leave the hole
exactly as wide as it was.

### METHOD — the differential probe, second time it decided a question argument could not

Recorded at the hub's direction as a reusable technique and an upstream amendment candidate.

**Shape:** hold the fixture, the probe and the command fixed; swap ONE file between two commits; run
the same probe twice. The delta is then attributable to that file and nothing else.

Both uses on this chip settled a question that had resisted argument:

| Use | Question that would not settle by reasoning | What the differential showed |
|---|---|---|
| governance lane | is the lane red because of this chip, or red already? | identical 53 violations at chip tip and at BASE, line for line → red already |
| R-6 regression | is the dropped `title AGAINST` a real behaviour change, or was it never emitted? | pre-S8 emits `hard-negative: kit/combo divergente…`; post-S8 does not → clean regression |

Why it works where argument does not: both questions are of the form "did MY change cause X", and
the honest answer requires observing the counterfactual. Reading the diff tells you what changed;
it does not tell you what the change DID. Reviewers on this chip — including both P6 gates and the
chip itself — repeatedly produced confident, plausible, opposite readings of the same diff. The
probe produced one output.

Cost is low enough that it should be reached for early rather than as a last resort: two blob reads at named commits
redirections, one throwaway `_test.go`, two runs, then restore. The whole R-6 differential took
minutes and replaced an argument nobody could win.

**Preconditions, learned the hard way here.** The probe must be able to FAIL — the first R-6 probe
PASSED while ranging over an empty candidate slice and reported success, and was only caught by
adding `if len(candidates) == 0 { t.Fatalf }`. And when the swapped file has a test file that
references new symbols, that test file must be swapped to the same commit or the package will not
build; a build failure is not a result.

### ESCALATED, not decided in-chip — RESOLVED by hub ruling R-6a

Two readings of R-6 are defensible and they differ exactly on this case:

- **(a) literal** — "no máximo um motivo por âncora por candidato", which is what S8 built. Then only
  precedence is open, and `AGAINST` must outrank `FOR` so a REJECT keeps its explanation.
- **(b) narrow** — the hub's rationale said "dois motivos **contraditórios** pra mesma âncora".
  `title` FOR + `title` AGAINST are not contradictory; they are complementary facts about one anchor.
  Under this reading the rule is "at most one UNAVAILABLE reason per anchor", the ruled-on
  `seller_sku`/`ean` duplication is still fixed, and this path is untouched.

The chip recommends (b) and sent `ESCALATION` with the differential above rather than choosing.
Picking either one IS the decision, and it changes persisted operator-facing data — "disagreement =
BLOCKED with evidence, never a unilateral decision".

**A defect in the chip's own slice card, stated plainly.** The S8 card listed "`FOR` or `AGAINST`"
as a SINGLE precedence rung and never said what happens when both land on the same anchor. The
worker had no rule to follow at exactly the point where the regression occurred. That gap is the
chip's, not the worker's, and it is the second time on this chip that a defective card instruction
produced a defective slice (see S3). The hub traced it one step further at ruling time: the card
inherited the gap from R-6's own wording, which was broader than the defect it was written for. So
the chain is ruling → card → slice, and the chip's card is the middle link rather than the origin.

### HUB RULING R-6a — reading (b), amending R-6 (`main@1be1aa9`)

The hub upheld the chip's recommendation and characterised the defect as its own wording, not the
implementation: R-6 said "no máximo um motivo por âncora" while its rationale said "dois motivos
**contraditórios**", and S8 implemented the sentence. **The invariant below REPLACES the R-6
wording.**

1. **`UNAVAILABLE` is exclusive per anchor** — emitted only if that anchor has no `FOR` and no
   `AGAINST`, and then at most once. Asserting "there is no signal" beside an actual signal is the
   real contradiction. This is A8.
2. **`FOR` and `AGAINST` coexist** on one anchor when both are fact. Zero dedup.
3. **Two `UNAVAILABLE` on one anchor** collapse to one, with deterministic tested precedence —
   declaration-derived beats seed-derived.

Note this is STRICTLY STRONGER than "at most one UNAVAILABLE per anchor", which is how the chip
phrased reading (b) when it escalated. Rule 1 additionally SUPPRESSES a declaration-derived
`UNAVAILABLE` on an anchor that already carries a `FOR` or `AGAINST` — a case the chip's own
formulation left open. S8's code happens to satisfy that half already; the ruling makes it a stated
invariant with a test rather than an accident of ordering.

**Both tests are mandatory by ruling**, and the hub was explicit about why the second one is not
optional: "o buraco é de cobertura, não só de código". Implemented in slice S9, dispatched with the
full collision matrix written out so no case is left to the worker's judgement — the specific
failure that produced this round.

### Not yet verified

The integration-test duplicate-anchor assertion was ADDED but **NOT RUN** — no database is
provisioned and the dev stack is a hub seam this chip may not boot. Recorded as NOT RUN. It is not
being counted as passing, and the worker reported it the same way rather than quietly omitting it.

## Slice S9 — chip verification (hub ruling R-6a, `2921d563`)

Dispatched with the 11-row collision matrix written out in full, because the two preceding slices on
this same function each broke on exactly one case their card left undefined. S9 reported the matrix
total and found no uncovered case.

Write set respected — `git show 2921d563 --stat` lists exactly the two permitted files, 20
insertions / 13 deletions in production. Zero migration, zero `apps/web`, zero contract/SDK.

### The implementation, read against the matrix

`appendProviderDeclaredUnavailableReasons`. Three passes, none of which ranges a map for output:

1. A PRE-PASS builds `hasSignal[anchor]` over all reasons. Because it is a pre-pass, rule
   1 is **order-independent**: a seed `UNAVAILABLE` that PRECEDES the `FOR` on the same anchor is
   still suppressed. The chip verified this case directly (below); S9's own tests do not pin it.
2. Non-`UNAVAILABLE` reasons are appended verbatim with **zero dedup**, so same-anchor
   `FOR` and `AGAINST` both survive in seed order (rule 2). `UNAVAILABLE` reasons are skipped when
   `hasSignal`, else the first is kept and its index recorded (rule 3, second half).
3. The declaration loop skips on `anchor.Supplied || hasSignal[anchor.Anchor]`, else replaces at the
   recorded seed index or appends in declaration order (rule 3, declaration wins).

All eleven matrix rows are satisfied by reading. Output order derives from slice iteration only.

### Chip-run probe — the regression is gone

Independent of S9's tests. A throwaway probe (`zz_chip_probe_test.go`, run then deleted, never
committed) drove state `TitleMatch` with a hard negative through the public
`buildCandidatesFromProducts` path and printed every reason. Verbatim:

```quoted-output
=== RUN   TestChipProbeTitleForAndAgainstSurvive
    status="REJECT"
    reason[0] anchor="title"      direction="FOR"         detail="match por título (ranking-only, nunca ACCEPT)"
    reason[1] anchor="seller_sku" direction="UNAVAILABLE" detail="seller_sku sem correspondência"
    reason[2] anchor="ean"        direction="UNAVAILABLE" detail="ean sem correspondência"
    reason[3] anchor="title"      direction="AGAINST"     detail="hard-negative: kit/combo divergente entre título do anúncio e produto interno"
    reason[4] anchor="marca"      direction="UNAVAILABLE" detail="provider não fornece a âncora marca"
    reason[5] anchor="refforn"    direction="UNAVAILABLE" detail="provider não fornece a âncora refforn"
--- PASS
```

`title FOR` at index 0 and `title AGAINST` at index 3 both survive on a REJECT. The S8 regression —
a REJECT with no stated contradiction and a `FOR` on the contradicting anchor — is closed.

The same output independently re-proves **R1** at the wire level: `seller_sku`/`ean` carry the SEED
detail ("sem correspondência" = the provider supplies this anchor, this listing has no value) while
`marca`/`refforn` carry the DECLARATION detail ("provider não fornece"). The distinction lives in
`detail`, with no 4th enum value — and it survives the finalizer rather than being flattened by it.

The probe asserted `len(candidates) == 0` as a hard failure before reading anything, per the vacuous
pass that already happened once on this chip.

### Chip-run probe — rule 1 is order-independent

Second probe, calling the choke point directly with the `UNAVAILABLE` seeded BEFORE the `FOR`:

```quoted-output
=== RUN   TestChipProbeUnavailableBeforeForIsSuppressed
    reason[0] anchor="ean" direction="FOR" detail="EAN confere"
--- PASS
```quoted-output

One reason out, the `FOR`. This is matrix row "1 × `UNAVAILABLE` + 1 × `FOR` | either", which S9's
suite covers only in the opposite seed order.

### Chip-run must-fail — the rule-2 mutation, re-run by the chip

S8's report was clean while its code regressed, so the highest-value proof was re-run by hand rather
than accepted. Mutation: restore blanket per-anchor dedup on the signal side of pass 2.

```quoted-output
--- FAIL: TestTitleMatchHardNegativeKeepsTitleForAndAgainstInSeedOrder (0.00s)
    generation_service_test.go:335: title reasons=[]domain.LinkCandidateReason{domain.LinkCandidateReason{Anchor:"title", Direction:"FOR", Detail:"match por título (ranking-only, nunca ACCEPT)"}}, want FOR then AGAINST
```

RED, and **not confounded** — the failure message names the missing `AGAINST` specifically, so the
test is pinning the regression itself and not some incidental difference. Mutation reverted from a
byte copy taken before it was applied; the file matches the frozen tip afterwards, byte for byte.

### ~~Reported by S9, not independently re-run by the chip~~ — SUPERSEDED IN ROUND 3, corrected in ROUND 4

**This heading was true when written and false by the time it was read, and it contradicted the
must-fail table.** When S9 landed, three proofs stood on the worker's report alone: the rule-1
suppression mutation, the rule-3 precedence inversion, and the map-iteration ordering mutation.
ROUND 3 re-ran all three by hand and pasted their outputs — see *Pasted must-fail outputs — the
ROUND-3 re-runs* below, where each is labelled `chip-re-run`. This section was not updated to match,
so the pack asserted both provenances for the same three proofs. Current status, single-sourced:
**all three are chip-re-run, outputs pasted.** Neither statement was a fabrication; the defect is
that a corrective updated one site and not its counterpart, which is the same shape as the stale
citations — a fact recorded twice decays in one place first.

The `-count` finding stands unchanged and is now chip-observed rather than S9-reported: the ordering
mutation goes **RED on 3 of 10 runs**, so a `-count=5` lane would have passed a genuinely
non-deterministic ordering. The `-count=10` in the card was load-bearing by luck, and a future
ordering guard on persisted rows should not be run at a lower count.

The integration-test assertion remains **NOT RUN** — no database, dev stack is a hub seam.

### Ladder at the frozen tip

Run by the chip at `2921d563`, GOCACHE/GOMODCACHE absolute:

- `go build ./...` — clean, no output, no VCS-stamping artifact this run.
- `go vet ./internal/modules/product_links/... ./internal/modules/connectors/...` — clean.
- `go test -count=10 ./internal/modules/product_links/... ./internal/modules/connectors/...` —
  `EXIT=0`, 9 packages `ok`, including `product_links/application` (6.878s) and
  `connectors/application` (5.949s).
- Nothing uncommitted when the tip was frozen — checkable now only as the property that every file
  named in this section appears in a commit reachable from the frozen SHA, not as a working-tree
  observation a later reader could repeat.

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
(`unit-prefix-sweep.go.txt`, alongside this pack, run against `hardNegativeDimensionPattern` copied verbatim; kept as `.go.txt` so no Go tool ever picks it up as package source):

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

## The class, five rounds of it, and what finally removed it

Five P6 rounds failed on this pack. **None failed on the code.** The cold reviewer of round 5, having
read the production source, reported F-01, F-02 and F-03 real, tested, and their guards load-bearing;
that verdict has not been in dispute since round 1. What failed, five times, was the pack's account of
itself, and the durable value of this chip is the diagnosis rather than any one fix.

### One class, five disguises

| Round | What the gates found | Root cause |
|---|---|---|
| 1 | Three formatting defects where a CRLF diff hid a real gofmt violation | A signal-masking change and a substantive one landed in the same commit, so the noisy one hid the quiet one |
| 2 | The pack asserted checks that had never run | Prose authored in the shape of a verdict, before the verdict existed |
| 3 | Citations correct when typed, wrong when read | A corrective inserted lines above them; the pack was pinned to a tip that had moved |
| 4 | The remedy for round 3 was itself silently narrow | The recognizer required a filename, so bare citations were invisible while the closure sentence claimed totality |
| 5 | The remedy for round 4 was silently narrow in the same way | The recognizer required backticks, so an unbackticked citation, a comma-list and a citation inside a fence were invisible — one of them stale on the pack's own marker line |

Rounds 3, 4 and 5 are the same defect three times, each one **inside the fix for the previous one**.
Every fix widened a recognizer, which is why every fix failed the same way: while coverage is defined
by what the recognizer matches, whatever falls outside it is both uncovered and silent, and the tool's
coverage line reads clean either way. A wider recognizer is a smaller blind spot, never no blind spot,
and a pack cannot tell the difference from the inside.

### The finding that ended it, which was the hub's and not a gate's

Each gate sees one round; the shape only appears across rounds. The hub's reading: the pack had grown
an audit of citations inside the document containing the audit, exemptions declared as line ranges of
that same document, and a count of citations in the text that stated the count. **Verifying the pack
had become more expensive than verifying the code** — which inverts what a pack is for, since it
exists to make verification cheaper than redoing the work.

The rule outlives this chip and is carried upstream: *an evidence apparatus that costs more to verify
than the thing it evidences has inverted its purpose, and the answer to that is not more apparatus.*

### What R-22 changed here

Coordinates left the prose. Claims are behavioural and anchored to test and symbol names, which are
stable under line shifts and checkable by running the named test. The coordinate table stayed, and is
generated from the code by `cite-table.py` into `coordinates.txt`, regenerated and diffed by the lane
so drift is a failing diff rather than something a reviewer has to notice.

This does not overturn the §11 judge, who rejected replacing the coordinate with the symbol. The
anchor's type is still a coordinate and still mechanical. What is gone is the **duplication** of the
coordinate into prose — the last place a decaying copy still lived.

Three further rules landed with it, all ratified by the hub:

- **Mutable-axis commands are not evidence.** Freezing a SHA in prose does not freeze the command.
  This pack carried write-set counts derived from a range ending at the current branch tip; they were
  true when written and had decayed to different numbers by round 5. The rule: every evidence command
  names its own tree-ish or it is not evidence. The report reads part of that rule and prints, beside
  its count, the forms it does not see — a bare `HEAD`, `git rev-parse`, a mutable named ref, a plain
  worktree `grep`, and a line where an unrelated SHA satisfies its hex-ish test.
- **Declarations travel with the content they describe.** The first draft of the round-5 tool declared
  its exemptions as pack line ranges — which is the same defect one level down, since such a
  declaration decays the moment a corrective inserts a paragraph above it and then silently excuses
  the wrong lines. A historical citation now carries its SHA inline; quoted program output declares
  itself with a fence tag. Never by position.
- **Some errors are irreducible.** Round 5 found a grep count at a frozen SHA that was simply wrong
  when typed — neither decay nor silence. Nothing structural prevents that, which is the argument for
  the generated artifact being regenerated and diffed by the lane rather than trusted.

### The four pasted listings are complete, and one command was wrong about why

Round 4 found a listing pasting 4 of 14 paths with no declaration that it was sampling — silent
truncation, the CHIP-MERCADO page-one class. The remedy was to verify every listing in the pack of
three or more repository paths. All four are complete; none samples:

| Listing | Paths | Command |
|---|---:|---|
| Swept `product_links` Go files | 34 | `git ls-tree -r --name-only 85b6c367 -- apps/server_core/internal/modules/product_links` |
| Contract and SDK existence witness | 14 | `git ls-tree -r --name-only 85b6c367 -- contracts/api packages/sdk-runtime` |
| Code paths in the chip diff | 15 | `git diff --name-only 917f7bb5 85b6c367 -- apps/ contracts/ packages/` |
| Governance synchrony differential | 5 | `git diff --name-only 8e37958a 93c90330 -- apps/server_core` |

The first row is a correction the round-5 derivation gate earned. The pack previously named
the index-reading listing form for it and called the result a fact at a frozen SHA. It is not: that
form reads the index, takes no tree-ish at all, and so can bind no SHA. The pasted content was correct and the
stated mechanism was wrong — which is precisely why the mutable-axis rule is read by a report now and
not left to care. I reported this defect against my own table; no gate had reached it.

### Method that outlived the rounds

The differential probe decided two questions argument could not: whether the governance lane was
already red at base, and whether a corrective had regressed a matching path while its own report read
clean. Hold the fixture, the probe and the command fixed; swap exactly one file between two commits;
run twice; the delta is attributable to that file alone. Its two preconditions, both hit here, are
recorded with the technique in `UPSTREAM-CANDIDATES.md` — swap the sibling test file to the same
commit or the probe reports a compile error instead of a behavioural delta, and assert the probe's own
non-vacuity first, because a fixture producing zero candidates makes a silent loop pass.


## Criteria table

**Every `file:line` below is against CODE TIP `2921d563`** — the last commit that touches code. Pack
commits after it modify only `.mnfs/…/_chip-anchors/`, so code line numbers do not move again.
Base is `917f7bb58e385847fba5612201823f9db48791c6` throughout.

This line previously read "against the final tip `1df627f8`" and was left standing through three
correctives while S7, S8 and S9 shifted the code underneath it. Both round-2 gates blocked on it. See
the ROUND-3 FULL ANALYSIS for the classification of each decayed citation, and the ROUND-4 section
*Class closed — SCOPE AND TOOL NAMED* for the sweep itself, which names its scope, its tool and its
tip. ROUND 3's sweep covered the criteria table only, so "every citation" is what ROUND 4 did, not
what round 3 did; a citation being right is recorded as explicitly as a citation being wrong.

### Evidence blocks the criteria table cites

Written 2026-07-26 as a ROUND-2 correction. The first draft of the criteria table cited these four
blocks as if they were in this pack; they were only ever in the chip's session transcript. Both P6
gates caught it. The outputs below are the real thing, pasted, so each criterion's *prova mínima* is
actually satisfiable by reading this file.

**C3 — the swept-file list.** The criterion requires "com a lista dos arquivos varridos", not just a
count. `git ls-tree -r --name-only 85b6c367 -- apps/server_core/internal/modules/product_links` →
**34 files**, module prefix stripped for width, `diff`ed empty against the command's live output:

```
adapters/connectors/identity_anchor_adapter.go          adapters/connectors/identity_anchor_adapter_test.go
adapters/postgres/link_candidate_repo.go                adapters/postgres/listing_snapshot_repo.go
adapters/postgres/summary_reader.go                     adapters/postgres/summary_reader_integration_test.go
application/auto_link_policy_test.go                    application/batch_service.go
application/batch_service_test.go                       application/decision_trail_test.go
application/generation_integration_test.go              application/generation_service.go
application/generation_service_test.go                  application/import_service.go
application/import_service_test.go                      application/resolution_service.go
application/resolution_service_test.go                  application/summary_service.go
application/summary_service_test.go                     composition/refresher.go
composition/refresher_test.go                           domain/internal_product_id.go
domain/internal_product_id_test.go                      domain/link_candidate.go
domain/listing_snapshot.go                              domain/product_link.go
domain/product_link_decision.go                         ports/link_candidate_store.go
ports/listing_snapshot_store.go                         ports/provider_identity_anchor_reader.go
ports/summary_reader.go                                 ports/workflow_store.go
transport/http_handler.go                               transport/http_handler_test.go
```

```
$ git grep -n 'mercado_livre\|ProviderCode *==\|ProviderCode *!=\|switch .*[Pp]rovider' 85b6c367 \
    -- 'apps/server_core/internal/modules/product_links/**/*.go' ':!*_test.go'
  ZERO HITS
```

**The two chip-added connectors imports, listed and classified — hub instruction of 2026-07-26.**
The false provenance row corrected above was covering *precisely* this spot, so the C3 block must
carry it explicitly rather than let a production-only grep imply it away. Every connectors reference
inside `product_links/application`, at tip:

| Site | Import | Origin | Classification |
|---|---|---|---|
| `generation_integration_test.go`, `connectorsapp` | `connectors/application` | **CHIP-ADDED** | test fixture — builds a real `MarketplaceCapabilityService` |
| `generation_integration_test.go`, `connectorsports` | `connectors/ports` | **CHIP-ADDED** | test fixture — names one anchor constant, `IdentityAnchorSellerSKU` |
| `import_service.go` | `connectors/domain` | pre-existing at base | production, unrelated to anchors, untouched |
| `import_service_test.go` | `connectors/domain` | pre-existing at base | test, untouched |

Why the two chip-added ones are **not** an R2 violation, proven rather than asserted. R2 forbids
*branching by provider inside `product_links`*. The usage, quoted from that test's fixture setup:

```go
capabilities := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
    ProviderCode:    "mercado_livre",
    IdentityAnchors: []connectorsports.IdentityAnchor{connectorsports.IdentityAnchorSellerSKU},
}})
svc := NewGenerationService(GenerationServiceConfig{
    …
    IdentityAnchors: productlinksconnectors.NewIdentityAnchorAdapter(capabilities),
})
```

That is a test ARRANGING data and wiring the real adapter — the provider code is an input value in a
fixture, not a branch in a code path. The three properties that make it safe are checkable, and each
was checked rather than asserted:

```quoted-output
$ git show 85b6c367:…/generation_integration_test.go | sed -n 1p
//go:build integration                       <-- not in the production build at all

$ git grep -c 'mercado_livre' 85b6c367 -- 'apps/server_core/internal/modules/product_links/**/*.go'
  84 hits, ALL in *_test.go — ZERO in production source.
  Every one is a DATA position: a struct field value (`ProviderCode: "mercado_livre"`),
  a map key in the generation-service tests, or a JSON request body in the transport
  handler tests. None is a comparison.

$ git grep -c 'mercado_livre' 85b6c367 \
    -- 'apps/server_core/internal/modules/product_links/**/*.go' ':!*_test.go'
  ZERO — no production file in the module contains the string at all.

$ git grep -nE 'ProviderCode *(==|!=)|switch .*[Pp]roviderCode' 85b6c367 \
    -- 'apps/server_core/internal/modules/product_links/**/*.go'
  ZERO HITS — production and test alike.
```

**Round-5 correction, and it is not a decay case.** This block previously reported **86** hits at
`08308afb`. That SHA is frozen and the number was simply WRONG when typed — re-running the same
command against the same commit returns 82. Neither line drift nor a silent recognizer produced it;
someone miscounted, and no amount of pinning prevents that. It is the argument for the generated
artifact: the block is now re-pinned to the final code tip, where the count is 84, and the count is
produced by `-c` and summed rather than eyeballed off a listing. The conclusion — all hits in test
files, none a comparison — was correct in every round and is what the criterion turns on.

The last one is the load-bearing check for R2, and it is stronger than the production-only grep
higher up: not merely "production does not branch on provider", but "nothing in this module branches
on provider anywhere, so there is no test pinning a branch that production could grow back". The layering rule the chip enforced separately — that `product_links/application`
**production** code must not import the connectors vocabulary type — also holds: the only production
connectors import in the package is `connectors/domain` in `import_service.go`, pre-existing and
unrelated.

Honest note on scope: this is the seam where R2's letter (no provider branching) and the layering
concern (no connectors vocabulary in the application layer) come apart. An integration test is
allowed to reach for both because its job is to prove the real wiring works; production code is not.
Nothing here is implied — the build tag, the usage sites and the grep are each independently
checkable above.

**C10 — the code path list**, pasted COMPLETE as the criterion demands, at the FINAL code tip:
`git diff --name-only 917f7bb5..85b6c367 -- apps/ contracts/ packages/` → **15 paths**, all listed
below with nothing elided.

**Why this is scoped to code paths, and why that is not a dodge.** Through ROUND 3 this block counted
the whole write set: *17 paths = 15 code + 2 pack*. That was TRUE when written and FALSE by the time a
gate read it, because the commit installing the round-3 remedy added `cite-audit.py` and
`cite-audit.txt` — making it 19. The number was falsified by the act of remedying the thing it was
part of. Under R-14 there are two ways out, and this block uses the second: scope the claim to an axis
the remaining work cannot move. A `.mnfs/`-only commit cannot add a path under `apps/`, `contracts/`
or `packages/`, so the code-path count is stable against every edit still to come to this pack. The
hub ruled the move legitimate on the same ground as the earlier C12 reframe: C10 is about code
collision axes, pack paths were never its object, so this names the scope the criterion always had. Had
the criterion been about the whole write set, re-scoping would have been a dodge and was to be denied.

`grep -c '^apps/web/'` over the **whole** write set (not the scoped one) = **0**, so the zero-web-file
guarantee is not an artefact of the narrower scope. The 15 code paths are byte-identical to the lists
produced at `1df627f8`, `93c90330` and `ac72eb82` plus `generation_service_test.go`, which was already
among them — the R5 corrective modified an existing file and added none, so the count did not move:

**Diff SIZE, per SHA, because the path count is stable and the line count is not.** The hub measured
2817 insertions at the code tip while the chip's CLOSED reported 3219; both are right about different
SHAs, and the pack states which is which rather than picking one:

| Range | Insertions | What it measures |
|---|---|---|
| `917f7bb5..2921d563` | 2817 (+167 del) | the code tip rounds 2–4 read — a HISTORICAL statement about a named SHA |
| `917f7bb5..85b6c367 -- apps/server_core` | **1356 (+167 del)** | production + test code alone at the FINAL code tip |

**Two rows were removed here rather than updated, and the reason is the ruling.** Both measured a
range whose far end was the branch's moving tip — figures that grow with every edit to this pack,
including the edit that states them. Under R-14 a fact the pack copies from an artifact generated
out of the pack is self-invalidating, so restating them at all is the defect; freezing them under a
promise not to edit further is what R-15 denied, because a freeze is discipline and a pointer is
structure. Whoever wants the whole-write-set size runs
`git diff --shortstat 917f7bb5..<the SHA being merged>` and reads it — this pack does not carry a
copy of it.

The row that survives is scoped to `apps/server_core` at a fixed code SHA, which is R-14's second
mechanism: no pack edit can move it. **The code-path count is 15 and is stable for the same reason;
only line counts move.** Same lesson as the lane rung: *a measurement names a SHA, or it names
nothing* — and for a size figure, the range as well.

```
apps/server_core/internal/composition/root.go
apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go
apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_test.go
apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go
apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go
apps/server_core/internal/modules/connectors/ports/marketplace_capability.go
apps/server_core/internal/modules/product_links/adapters/connectors/identity_anchor_adapter.go
apps/server_core/internal/modules/product_links/adapters/connectors/identity_anchor_adapter_test.go
apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go
apps/server_core/internal/modules/product_links/application/generation_integration_test.go
apps/server_core/internal/modules/product_links/application/generation_service.go
apps/server_core/internal/modules/product_links/application/generation_service_test.go
apps/server_core/internal/modules/product_links/application/resolution_service.go
apps/server_core/internal/modules/product_links/application/resolution_service_test.go
apps/server_core/internal/modules/product_links/ports/provider_identity_anchor_reader.go
```

15 code paths. `grep -c '^apps/web/'` → **0**. No `migrations/` path (also serves C9). No
`contracts/` or `packages/sdk-runtime/` path (also serves C5).

C10 is not itself an emptiness proof — the list is non-empty and pasted whole, so R-10 has nothing to
bite on. But the three criteria that LEAN on it (C5, C9, and the `apps/web` half of C10) do read
absence out of it, and each carries its own witness: C5 above, C9 immediately below, and for
`apps/web` the control is the list itself — 15 non-web paths printed by the same command that
prints nothing for `^apps/web/`.

**C9 — no migration, no persisted-reason rewrite, WITH THE R-10 WITNESS.**

```
$ git ls-tree -r --name-only ac72eb82 -- apps/server_core/migrations | wc -l    # (a) target EXISTS
75

$ git diff --numstat 917f7bb5..ac72eb82 -- apps/server_core/internal/modules/product_links
39      0       …/adapters/connectors/identity_anchor_adapter.go      # (b) CONTROL: same form,
44      0       …/adapters/connectors/identity_anchor_adapter_test.go  #     other path, DOES report
42      9       …/application/auto_link_policy_test.go
17      2       …/application/generation_integration_test.go
190     79      …/application/generation_service.go
619     49      …/application/generation_service_test.go
16      6       …/application/resolution_service.go
190     0       …/application/resolution_service_test.go
10      0       …/ports/provider_identity_anchor_reader.go

$ git diff --numstat 917f7bb5..ac72eb82 -- apps/server_core/migrations
  (empty)                                  # the proof: not one of the 75 files moved
```

75 migration files exist at the frozen tip, so the diff's emptiness is about the chip, not about a
mistyped path. Corroborated a second way by the C10 list, which contains no `migrations/` path and no
`product_links/adapters/postgres/` file at all — so the chip writes no `UPDATE` and no backfill.
Honest note, unchanged: a pre-existing `ON CONFLICT … DO UPDATE SET reasons = EXCLUDED.reasons` does
live in `link_candidate_repo.go`, untouched by this chip and depended on by the hub's own U1
("depois de regerar candidatos"). R3 forbids retro-editing persisted motivos; it does not forbid
regeneration from producing fresh ones.

**C5 — OpenAPI / SDK byte-identity, WITH THE EXISTENCE WITNESS RULING R-10 REQUIRES.**

The previous version of this block was pinned to `1df627f8`, and the criterion row cited
`-- openapi/`, **a directory that does not exist in this repo**. That command returns empty because
the path is absent, not because nothing changed — it would have passed identically had the OpenAPI
document been rewritten. Re-run at the frozen tip against the real paths, with the witness:

```
$ git ls-tree -r --name-only 85b6c367 -- contracts/api packages/sdk-runtime   # (a) targets EXIST
contracts/api/marketplace-central.openapi.yaml
packages/sdk-runtime/package.json
packages/sdk-runtime/src/activeSource.test.ts
packages/sdk-runtime/src/activeSource.ts
packages/sdk-runtime/src/dashboard.test.ts
packages/sdk-runtime/src/dashboard.ts
packages/sdk-runtime/src/erpImport.test.ts
packages/sdk-runtime/src/erpImport.ts
packages/sdk-runtime/src/index.test.ts
packages/sdk-runtime/src/index.ts
packages/sdk-runtime/src/listings-signals.test.ts
packages/sdk-runtime/src/market.ts
packages/sdk-runtime/tsconfig.json
packages/sdk-runtime/vitest.config.ts
                                                        # 14 paths, COMPLETE — nothing elided

$ git diff --stat 917f7bb5..85b6c367 -- apps/server_core      # (b) CONTROL: same command, other scope
 15 files changed, 1356 insertions(+), 167 deletions(-)

$ git diff --stat 917f7bb5..85b6c367 -- contracts/ packages/sdk-runtime/     # the proof
  (empty)
```

(a) proves the target is there to be diffed; (b) proves the command form produces output when
something did change. Only with both does the third command's emptiness mean *unchanged*. R1 holds.

**ROUND-4 correction — this paste showed 4 of 14 paths, with no marker.** The cold gate blocked on
it and the hub split the question (R-16), because the two gates were right about different things
and the difference decides the scope of the fix:

- *R-10 was discharged.* A witness's job is to prove the target EXISTS, because the defect R-10 was
  written to kill is `openapi/` — empty-because-absent. Existence is proven by one file; four prove
  it. The GPT gate passed C5 on exactly this reading and was correct.
- *The artifact was still defective.* Pasting 4 of 14 and presenting it as the listing is **silent
  truncation**, which violates a separate and older rule — no silent omission — and blocks on its
  own merit. Precedent in this mission: CHIP-MERCADO's page-1 truncation, which both self-gates
  missed for the same reason, that a truncated result looks exactly like a complete one.

So the remedy is NOT "re-audit every other R-10 witness for completeness". It is wider on one axis
and narrower on another: **every pasted listing in the pack is complete, or declares what it
samples** — all pastes, not just R-10 witnesses; and the other witnesses remain valid *as witnesses*.

**R-16 sweep.** SCOPE: every fenced block in `EVIDENCE.md` carrying three or more repository paths,
found by walking the pack's code fences programmatically rather than by recalling which blocks paste
listings. TOOL: for each block, the pasted lines sorted and `diff`ed against the live command's output
at the frozen code tip. All four are complete; none samples. The result table lives in
*The four pasted listings are complete, and one command was wrong about why* — stated once, there,
because one of its rows carries a correction the round-5 derivation gate earned and a second copy of
a corrected table is exactly the decaying duplicate R-22 removes.

**C12 — the root.go grant region**, quoted verbatim, re-run at the frozen tip **with the R-10
witness pair.** The hub applied R-10 retroactively to its own C12 amendment: "0 removed lines under
`-w`" is an emptiness proof on a PATH, and a mistyped path passes it identically.

```
$ git ls-tree -r --name-only ac72eb82 -- apps/server_core/internal/composition/root.go   # (a) exists
apps/server_core/internal/composition/root.go

$ git diff -w --numstat 917f7bb5..ac72eb82 -- …/product_links/application/generation_service.go
174     63                                      # (b) CONTROL: the same -w form DOES report removals

$ git diff -w --numstat 917f7bb5..ac72eb82 -- apps/server_core/internal/composition/root.go
3       0                                       # the proof: 3 added, 0 removed
```

The three added lines:

```diff
+	productlinksconnectors "marketplace-central/apps/server_core/internal/modules/product_links/adapters/connectors"
+	productLinkIdentityAnchorReader := productlinksconnectors.NewIdentityAnchorAdapter(marketplaceCapabilities)
+		IdentityAnchors: productLinkIdentityAnchorReader,
```

(The `-w` matters: the third line joins the `GenerationServiceConfig` literal whose four pre-existing
keys were realigned, which is whitespace-only and therefore not a token removal. The cold gate
verified this independently against base — that literal carried exactly four keys at base, all four
still present at the tip, none removed.)

Every cell states BEHAVIOUR and names the test or symbol that carries it. Positions are not repeated
here: `coordinates.txt` resolves THE ANCHORS IT LISTS to file and line, generated from the code, and
prints under UNRESOLVED the anchors it could not resolve. That split is R-22 — the anchor's type is
still a coordinate, but the coordinate is no longer duplicated into prose, which is the last place a
decaying copy lived.

**This paragraph used to say the table resolves EVERY anchor below, and R-24 struck that word.** The
generator dropped a name it could not find and its footer still printed a resolved count, so the
claim was total in the wording and partial in the code — the same shape as the two recognizer claims
struck alongside it. The remedy is not a wider matcher: the table is NAVIGATION, so a name it misses
costs a reader one manual lookup by name and cannot cost anyone a false proof. Where a cell's anchor is
absent from the table, resolve it by name; the UNRESOLVED and NOT-TREATED-AS-ANCHORS sections say
which those are instead of leaving the reader to discover the silence.

| ID | Verdict | Evidence |
|----|---------|----------|
| C1 | **PASS** | The capability is DATA end to end. `IdentityAnchor` is a string type in the connectors ports package with five declared constants and a backing slice; `KnownIdentityAnchors` returns a COPY of that slice, so no caller can mutate the vocabulary. `product_links` consumes it through its own port, `ProviderIdentityAnchorReader`, implemented by `IdentityAnchorAdapter` in the product_links connectors adapter — which carries a **compile-time assertion** that it satisfies the port, so the wiring cannot silently rot into a nil capability. The composition root constructs that adapter and hands it to the generation service as the `IdentityAnchors` field of `GenerationServiceConfig`; those three added lines are the whole grant, quoted verbatim in the C12 block |
| C2 | **PASS on the *prova mínima*, with a declared divergence — see below** | The hardcoded list is GONE from production, not merely bypassed: `git grep -n mandatoryUnavailableReasons 85b6c367 -- 'apps/server_core/**/*.go'` = **0** in tracked source, where base carried it in the generation service alone. `"marca"`/`"refforn"` as quoted literals in `product_links` production code = **0**, swept across the whole module and not just `application/`. Every surviving hit is in a TEST file, which the criterion permits and which are the tests that PROVE the names now arrive from the declaration — `TestGenerateLinkCandidatesUsesProviderDeclarationForUnavailableReasons` and its precedence siblings in the generation-service tests, plus the adapter's own `identity_anchor_adapter_test.go`. Counting method is stated because it changes the number: quoted literals only. A bare-word grep returns more lines — identifiers and prose — which is why the cold gate reported a larger figure and neither of us was miscounting |
| C3 | **PASS** | No provider branch exists anywhere in the module's production code. `git grep 'mercado_livre\|ProviderCode *==\|ProviderCode *!=\|switch .*[Pp]rovider' 85b6c367 -- 'product_links/**/*.go' ':!*_test.go'` = **ZERO HITS**, over the complete swept list in the C3 evidence block above — the criterion requires the list, not just the count, and R2 is satisfied by construction rather than by discipline: the capability arrives as data through the port, so there is no provider identity in scope to branch on |
| C4 | **PASS** | The two details are distinct and both present, per R1 carried in `Detail` with no 4th enum value. The scoring functions (`applySingleAnchorScore`, `applyUnresolvedScore`) emit `"<anchor> sem correspondência"` — the provider SUPPLIES the anchor, this listing has no value for it. `appendProviderDeclaredUnavailableReasons` emits `"provider não fornece a âncora %s"` — the provider does not supply it at all — and it OVERWRITES the seed form at the same index when both would apply, so the declaration wins rather than duplicating. Named by test: `TestGenerateLinkCandidatesUsesProviderDeclarationForUnavailableReasons` asserts the declaration form, and the precedence tests assert that a real mismatch outranks both. **Both forms are visible side by side in the S9 probe output** (S9 section) on a single candidate: `seller_sku`/`ean` carry the seed form, `marca`/`refforn` the declaration form — which is the distinction being legible in one payload rather than in two runs |
| C5 | **PASS** | The contract surface is untouched: `git diff --stat 917f7bb5 85b6c367 -- contracts/ packages/sdk-runtime/` → **empty**, with the R-10 witness pair pasted in the C5 evidence block — both trees EXIST at the tip (complete `git ls-tree` listing, not a sample), and the same command form over `apps/server_core` reports a large non-empty diff, so the emptiness is a fact about the chip and not a form that can only print zero. R1's contract-lock for M-06 therefore holds. Corroborated independently by the C10 list, in which no path under either tree appears |
| C6 | **PASS** | The golden case from the dispatch prompt is a real test: `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm` drives listing `Toalheiro Simples Soul Zen 50cm Cromado` (**cm**) against ERP `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` (**mm**) and asserts CONFIRM / MEDIA / score 70 — orientation per adjudication A3, NOT the dispatch prompt's EXEMPLO-IO, which has listing and product swapped. `TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate` proves the rule no longer rejects on unit spelling alone, and `TestDimensionCanonicalizationUsesExactMillimetres` covers decimal comma, decimal point, inches and metres as four subtests. Conversion is EXACT, not floating: `normalizeDimensionToken` scales inches by `big.NewRat(127, 5)`, and its suffix list puts the bare metre suffix LAST so it cannot shadow `mm` or `cm` by prefix — the ordering is load-bearing and asserted, not incidental. The rule still contradicts when it should: `TestDifferentCanonicalDimensionsStillRejectConcordantCandidate`. Must-fail run PASTED below |
| C7 | **PASS** | Hard negatives still cap the CORROBORATED path — the path where SKU and EAN both agree, which is the only one that could auto-approve. `TestCase7VoltageHardNegativeCapsBaixaReject` and `TestDifferentCanonicalDimensionsStillRejectConcordantCandidate` both drive SKU+EAN concordant and assert BAIXA / REJECT. `TestCase6DokaKitHardNegativeCapsBaixaReject` and `TestCase10DimensionHardNegativeCapsBaixaReject` are EAN-only paths — real coverage, but NOT the corroborated path, and an earlier version of this cell wrongly presented all four as corroborated. `TestDimensionPresenceAndGradeRulesRemainNonBlocking` holds the non-blocking side. **ROUND-4 correction, kept visible:** the ROUND-3 sweep classified the `cor` branch (`hardNegativeColor`) as having "NO test at all". FALSE — `TestHardNegativeKindsBlockConcordantSKUAndEAN` in `auto_link_policy_test.go` covers `"cor"` with `{"PUXADOR DHARMA AZUL", "PUXADOR DHARMA PRETO"}` on the corroborated path, alongside kit/combo and voltage. The GPT gate caught it; the chip verified it against the file rather than taking the gate's word. The wrong claim is left standing rather than deleted, because a sweep that INVENTS an absence is the same class of defect as one that misses a presence, and it arrived in the round whose whole purpose was to eliminate the class. C7's only remaining declared gap is the grade-token finding below |
| C8 | **PASS** | The three limits are independent by construction: `ListLinkWorkflows` defaults three separate locals from three separate `ListLinkWorkflowsInput` fields (`Limit`→20, `LinkLimit`→2000, `AuditLimit`→10000), each guarded by its own `<= 0` test, with no shared derivation — so no one input can move another's default. `TestListLinkWorkflowsUsesIndependentDefaultLimitsAndReturnsAll29Links` asserts all 29 links come back from a 29-link fixture; `TestListLinkWorkflowsDefaultsDoNotVaryWithLimit` pins that independence across four caller limits; `TestListLinkWorkflowsHonorsIndependentExplicitLimits` proves explicit values are honoured rather than clamped. The defaults are asserted as EXACT literals, never as "≥29" — the old `limit*5` derivation would satisfy a floor and fail these. Must-fail run PASTED below |
| C9 | **PASS** | ZERO migration, as the contract requires: `git diff --numstat 917f7bb5 85b6c367 -- apps/server_core/migrations` → **empty**, with the R-10 witness pair in the C9 evidence block — the migrations directory EXISTS and is populated at the frozen tip, and the same `--numstat` form over `product_links` reports several changed files, so the emptiness is about the chip and not about a mistyped path. Corroborated independently by the C10 list, which contains no `migrations/` path and no `product_links/adapters/postgres/` file at all: the chip writes no `UPDATE` and no backfill, which is R3. Characterised honestly rather than claimed absolutely — a pre-existing `ON CONFLICT … DO UPDATE SET reasons = EXCLUDED.reasons` DOES exist in `link_candidate_repo.go`, untouched by this chip and depended on by the hub's own U1 ("depois de regerar candidatos"). It is the candidate **regeneration** upsert. R3 forbids retro-editing persisted motivos; it does not forbid regeneration from producing fresh ones |
| C10 | **PASS** | ZERO `apps/web` file, and the write set proves it as its own control: `git diff --name-only 917f7bb5 85b6c367 -- apps/ contracts/ packages/` prints a non-empty list of code paths — pasted COMPLETE in the C10 evidence block above, at the final code tip — while the same command's output contains no path under `apps/web/` at all. A form that can print paths, printing none there, is the R-10 pair collapsed into one listing. **Re-scoped in ROUND 4 (R-14 mechanism 2):** this cell used to count the whole write set including pack artifacts, and that number was falsified by the very commit that installed the round-3 remedy. The criterion is about code collision axes and pack paths were never its object, so scoping to code paths names the scope the criterion always had rather than narrowing it to escape the proof — the same reframe the hub applied to C12. The axis is now one the remaining work cannot move: a `.mnfs/`-only commit cannot add a path under `apps/`, `contracts/` or `packages/` by construction |
| C11 | **PASS** | L0+L1 re-run by the chip at the frozen code tip: `go build ./...` clean, `go vet ./...` clean (both exit codes read unpiped, after the vacuous-pass incident recorded in the Ladder), `go test -count=10` → `EXIT=0`, 9 packages ok. The repeat count is deliberate: it is what makes a map-iteration-order dependency fail rather than pass four times out of five. Governance rung: hub-owned per R-b, was OPEN and is now **CLOSED** — the hub re-ran the lane at code tip `2921d563` on the `REQUEST` this pack carried, 53 violations, output IDENTICAL to BASE and to the earlier measurement, so the chip introduced none. Evidence `main` @ `7c54bef`, `_chip-anchors/hub-governance-lane.md`. Cited, not re-proven |
| C12 | **PASS** | The additive grant is exactly three added lines and nothing removed: `git diff -w --numstat 917f7bb5 85b6c367 -- apps/server_core/internal/composition/root.go` → **3 added, 0 removed**, with the R-10 witness pair in the C12 evidence block (the path exists at the tip; the same `-w` form over the generation service reports substantial removals, so it is not a form that can only print zero). The hub applied R-10 to its own C12 amendment for exactly this reason. The three lines are quoted verbatim in the C12 evidence block above and in the CLOSED payload; `gofmt -l` clean, proven CR-stripped, because the naive file-scoped check is uninformative on this checkout |

### C2 — the criterion's statement and its *prova mínima* do not cover the same set

Found by the chip on 2026-07-26 while scoping R-6, after both gates had passed C2. Disclosed rather
than resolved in-chip, because resolving it is a scope decision that is the hub's.

C2 reads, in two parts:

> `mandatoryUnavailableReasons()` não existe mais e **nenhum nome de âncora** está hardcoded no
> gerador | grep do símbolo = 0 hits; grep de `"marca"`/`"refforn"` em `product_links/application` =
> 0 hits em código de produção

The *prova mínima* names `marca` and `refforn`. Both are zero, so C2 passes as written and as gated.
But the STATEMENT says "nenhum nome de âncora", and the generator does still contain hardcoded anchor
names. `"seller_sku"`, `"ean"` and `"title"` appear as reason seeds in exactly five functions —
`applySingleAnchorScore`, `buildConflictCandidates`, `buildConcordantCandidate`,
`buildCollisionCandidates` and `applyUnresolvedScore` — enumerated by walking the file's function
boundaries rather than recalled, because the previous version of this sentence listed line numbers and
a reader had no way to tell which function each belonged to. Read literally, the statement is not
satisfied.

The chip's reading — offered as reasoning, not as a ruling — is that the *prova mínima* is the
operative one, because the two name groups are different kinds of thing:

- `marca` / `refforn` were **fabrications about the provider**. The old
  `mandatoryUnavailableReasons()` asserted "the provider does not supply `marca`" for every provider,
  unconditionally, having never asked. That is the ADR-17 violation F-01 exists to delete, and it is
  gone.
- `seller_sku` / `ean` / `title` are the **core's own search vocabulary**, fixed by IC-01 Amendment
  A2: these are the only cross-side anchors the matcher searches
  on. Naming them in a reason that reports *what the matcher itself did* is a statement about our own
  behaviour, not an unasked claim about the provider's.

The distinction is exactly what R-6's precedence encodes: where the two collide, the provider's
declaration outranks the core's "I searched and missed", because only one of them is a fact about the
provider. After S8, the core never claims to have searched an anchor the provider does not supply.

What the chip did NOT do: widen F-01 to make the core's search vocabulary provider-declared too.
That is a real design question — it would mean the matcher's anchor loop is driven by the
declaration rather than by A2 — and it is well outside this chip's contract. Flagged to the hub as a
finding, not actioned. **C2 stands as PASS on its stated proof; the hub decides whether the
statement wants amending or the scope wants widening.**

#### Hub ruling at acceptance — C2 PASS, the criterion's PROSE was the defect

The hub ruled on this at acceptance, having read the code at the tip rather than accepting the
description. **C2 PASS as it stands; no corrective.** The hub named the defect as its own and as the
same shape as R-6 — a ruling/criterion whose WORDING is wider than the invariant it meant to state.

The invariant, restated by the hub and narrower than C2's prose, is **unidirectional**: *no assertion
ABOUT THE PROVIDER that does not derive from the provider's own declaration.* Hardcoded
`seller_sku`/`ean`/`title` are the core's comparison vocabulary — what it knows how to compare
against the ERP — not claims about the provider, so they do not touch the invariant.

Proof the hub cited and the chip re-read at the code tip, in the declaration loop of
`appendProviderDeclaredUnavailableReasons`: when an anchor is not `Supplied`, the
core-seeded reason is **overwritten** — `finalized[index] = declarationReason` — by
`provider não fornece a âncora %s`. The core's "I searched and missed" cannot outlive the provider's
declaration. That was A8, and it is closed.

#### Hub finding, NOT this chip's, referred to the mission backlog

Raised by the hub at acceptance, verified in the code, explicitly **not** a corrective on this diff:

`appendProviderDeclaredUnavailableReasons` opens its declaration loop with
`if anchor.Supplied || hasSignal[anchor.Anchor] { continue }`. So an anchor the provider **declares
it supplies** but which the core has no comparison for — an `mpn`, a second marketplace's `gtin14` —
yields **no reason at all**. Chip-verified at `2921d563`: the seeding paths only ever emit
`seller_sku`, `ean` and `title`, so nothing else can reach `hasSignal`, and the `Supplied` branch
then skips it.

Not a fabricated value, so not an ADR-17 violation — it is silent OMISSION: the screen says nothing
about an anchor that exists on the provider's side. The hub's characterisation: a cousin of ADR-17
and the next hole on this trail. Filed to the mission backlog, out of scope here.

## Ladder

Profile §2 bindings. `GOCACHE`/`GOMODCACHE` bound ABSOLUTE before every Go command, as §2 requires
on Windows (echoed into the transcript:
`GOCACHE=/c/…/.claude/worktrees/chip-anchors/apps/server_core/.gocache`).

| Rung | Command | Result |
|---|---|---|
| L0a | `go build ./...` (full tree, plain — no `-buildvcs=false`) | **exit 0** |
| L0b | `go vet ./...` (full tree) | **exit 0** |
| L1 | `go test -count=1 ./internal/modules/product_links/... ./internal/modules/connectors/... ./internal/modules/mutations/... ./internal/composition/...` | **exit 0**, 27 packages, all `ok`, zero failures |
| GOV | `npm run harness:governance -- -BaseSha 917f7bb5…` from a clean detached worktree | **hub-run per R-b**, differential PASS — see below |
| LIVE | U1–U3 browser live-drive on the connected ML account | **PASS — hub-run, ruling R-8** — see below |

**The PACK rung was on this ladder and R-24 took it off.** It is not a rung, it is a report, and the
row is removed rather than annotated because a row in a ladder table is a merge gate by position.
See the section below for what the tool does now and why nothing depends on its exit code.

### `cite-table.py` — a report, not a rung (R-24)

R-22 made the coordinate table a generated artifact and put the check on the ladder. Round 6 then
failed on that check — not on the code, and not on the prose it was built to police. All three
blocking findings were claims THE TOOL MADE ABOUT ITS OWN COVERAGE: a prose scan that called itself a
SUPERSET while a negative lookbehind hid `123:456`; a mutable-axis ban that called itself a ban and
did not see bare `HEAD`, `git rev-parse`, or a mutable named ref; a table that claimed to resolve
every anchor and silently dropped the ones it could not find. Same shape three times — total in the
wording, partial in the code.

The hub's ruling widened nothing. It removed the totality:

> A VERIFICATION ARTIFACT KEEPS ONLY THE CLAIMS IT CAN MAKE TOTALLY.
> A CLAIM IT CANNOT MAKE TOTALLY BECOMES A REPORT, NOT A GATE.

A report that says *this I resolved, this I did not* has no totality claim left to falsify, so the
three findings cease to exist by construction rather than by pardon — the false claim is deleted, not
excused. It is ADR-17's honest-unknown rule turned on our own evidence apparatus: the same discipline
this product applies to data it does not have.

What the tool does now:

- **Coordinates in prose, mutable-axis commands.** Both passes still run and both still print. What
  changed is that each states what it does NOT see, in the same output as what it does — the
  lookbehind, the exact git verbs, the fact that ANY `quoted-*` tag is honoured and not only
  `quoted-output`. The RULE is untouched and whole: every evidence command names its own tree-ish or
  it is not evidence. The pass is a partial reading of it and now says so.
- **Unresolved anchors are reported.** `if not locs: continue` was the defect; the name now lands in
  an UNRESOLVED section of the generated table with the number of prose lines that name it. Backticked
  spans carrying a call suffix — never anchors to the recognizer, one of them real and found by the
  cold gate — are listed too. **R-17 is ratified a fourth time**: the requirement to say what could
  not be resolved stands. What R-24 changed is the CONSEQUENCE — a report line, never an exit code.
- **No widening.** `TOPLEVEL` is not taught grouped `var (…)`. The git-verb list is not extended. The
  lookbehind is untouched. Widening a recognizer whose output no longer gates anything buys nothing,
  and it is the exact move that produced rounds 3, 4, 5 and 6.
- **Drift still fails, for the author.** The two modes stay split — a checker that regenerated the
  artifact could never detect drift in it. `--strict` still exits 1. Nothing merges on that exit code.

R-11's mechanism retires with its object. It made citation validity mechanical because the cold gate
has no shell, and that was about claims made IN PROSE, which R-22 removed. The table is navigation
now: a wrong row costs a reader a wrong jump, not a false proof. The insight survives, the enforcement
does not.

This section restates no count from the generated artifact. That is R-14, and reproducing it inside
R-22's own remedy is the recurrence rounds 3, 4 and 5 kept demonstrating.

`cite-audit.py` and `cite-audit.txt` — the round-4 apparatus — are **superseded** by this rung. They
stay in the pack, and are not deleted: the D24 ledger row points at them, and they are the record of
what round 4 found. Each carries a SUPERSEDED banner in its own first lines, so a reader who opens one
without reading this section still learns it is history.

### Rung LIVE — U1–U3, HUB-RUN per R-8, PASS

`LIVE-VERIFIED: U1-U3 hub-run na conta ML conectada (R-8), evidência
.mnfs/MIS-006-integracao-fundacao/_chip-anchors/hub-live-drive-u1-u3.md @ main 40623b5`

Hub evidence, committed to `main` @ `40623b57`. **Cited here, not re-proven** — the chip never boots a
server, never binds port 8080, and never loads `.env*`; the dev stack is a hub-owned seam. The chip
neither ran nor could have run this rung.

**How the hub reached the chip's code without merging it, recorded because it is the reusable part:**
the primary checkout's compose stack (where `.env` resolves) was brought up with an override changing
ONLY the source of the `/workspace` bind mount, re-pointing it at
`.claude/worktrees/chip-anchors`. Compose merges volumes by CONTAINER path, so the `.env` continues to
resolve against the primary while the served tree is the chip's. Confirmed by the hub before driving:
`identity_anchor_adapter.go` present inside `/workspace`. No credential was copied into the worktree
and the chip's tree was not written to. Connected ML account, 34 listings, candidates **regenerated**
(`generated_count: 38`) — not backfilled, so R3 holds through the live drive too.

| Rung | Result | What was actually driven |
|---|---|---|
| **U1** — F-01 reaches the screen | **PASS** | Before: `marca inexistente no lado provider` — the CORE asserting a fact about the provider. Now: `marca: provider não fornece a âncora marca`, derived from the ML adapter's declaration of `[seller_sku, ean, title]`. A second provider that declares `marca` makes the row disappear on its own, with no code change — which is what "capability is data" means in the screen rather than in the type system |
| **U2** — F-02 reaches the verdict | **PASS** | Listing `MLB4735326915` "Toalheiro Simples Soul Zen 50cm Cromado" (**cm**), matched against ERP product "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO" (**mm**), went from **BAIXA 25%, rejected by `Título hard-negative: medidas`** to **`exact_sku` CONFIRM, MEDIA 70%**. The EXEMPLO-IO case from the dispatch, live, in the orientation adjudication A3 fixed. **And the guard did not loosen**: ambiguous-EAN rows in the same view stay BAIXA 20% with `âncora ambígua` |
| **U3** — F-03 reaches the list | **PASS** | Three independent sources agree on 29: `select count(*)` = 29 resolved; `GET /link-workflows` with NO `limit` = 29 items carrying `current_link`; `document.querySelectorAll('table tbody tr').length` on the Resolvidos tab = 29. KPI "Resolvidos hoje" = 29. The 9 that the old shared default of 20 was eating are back |

U2 is the strongest single result in this chip: the golden case named in the dispatch on day 1,
driven end to end on the real account, moving in the intended direction **while the hard negative it
could have broken stays blocking**. A pass that only showed the CONFIRM would not have proved that.

**The U2 cell above was corrected in the round-6 corrective, and the same inversion was in the hub's
own live artifact.** It attached the ERP description to the listing ID — the inversion adjudication
A3 had already fixed once, in the header, reappearing later in the same document. The golden test is
the authority: it sets the listing title to the `50cm` form and puts the all-caps `500MM` name on the
ERP side under `sku:33698`. The chip flagged the hub-owned copy rather than editing it, and the hub
corrected it on `main` in the same changeset as R-24 — by a NEW COMMIT, never by an edit under the
existing pin, because a pin names its SHA and editing beneath one makes the pin lie. The pack's
pointer to the earlier SHA therefore stays valid and stays where it is.

**Hub findings from the drive — three, none blocking, none this chip's.** Recorded because the hub
raised them, and explicitly NOT actioned here: acting on them would be the scope leak the dispatch
forbids.

1. **`refforn` carries no information as an anchor.** It is a supplier reference belonging to OUR
   ERP; no provider will ever declare that it supplies one, so it is now a permanent `UNAVAILABLE`
   for every provider forever. `marca` is genuinely different — a concept both sides hold, so a
   provider CAN declare it. F-01 made the declaration honest; it did not ask whether every anchor in
   the vocabulary deserves to be there. **That is the next question on this trail, and it belongs to
   whoever owns the vocabulary.**
2. **A `Supplied` anchor the core cannot compare vanishes silently** — the declaration loop skips it
   on `anchor.Supplied` and nothing downstream reports the gap.
   Already filed in this pack from the hub's earlier message, under the C2 section — same finding,
   independently re-surfaced by the live drive.
3. **The "SKU ML" column renders `provider_code`.** M-06's F-05, not this chip's; `apps/web` is zero
   paths in this diff.

**Why this section was written after the gates, not before.** The rung was held in
`scratchpad/live-rung-pending.md` while the round-3 gates were reading, because the COLD gate reads
the WORKING TREE rather than a pinned SHA — editing the pack mid-review changes the document under
the reviewer, which is exactly what corrupted one gate's input in round 1. The hub's drive had
finished by then, so the hold was about the gate, not about the stack, and it was stated that way to
the hub at the time.

**How the rung got here — the R-8 history, kept because the reasoning is the point.** The cold Opus
gate raised a SHOULD-FIX:
the diff touches provider-adapter scope (the ML `capability_adapter.go`) with no `LIVE-VERIFIED:` or
`LIVE-WAIVED-BY-OPERATOR:` marker. The chip requested a waiver — the change is a static declaration
with no provider I/O, and U1–U3 are hub-run by dispatch — and **the hub DENIED it**, ruling R-8: the
waiver is not needed because the rung is not the chip's to close. The chip recorded it OPEN; the hub
filled it at acceptance when U1–U3 passed — which is what the section above is. **No merge before
that**, and none happened. The chip did not attempt U1–U3 and did not self-grant the waiver. Same
shape as the governance rung under R-b: a rung the chip cannot close honestly is recorded OPEN and
named, never inferred PASS and never silently omitted — and in both cases the hub's own run later
CONFIRMED what the chip had declined to assert. Two for two is not proof the instinct is always
right, but it is the reason the rule costs nothing: the confirmations arrived anyway, and the pack
never had to be walked back.

L1 scope is touched-packages-plus-guard-suites, not a full sweep, per §2: "full sweep only when
migrations/platform touched" — and this chip has **zero** migrations (C9). The four trees are the
complete set the diff reaches: `product_links` (owner), `connectors` (the declaration side),
`mutations` (the only external consumer of `ListLinkWorkflowsInput`, via its two
`Limit: 2000` call sites in `adapters/productlinks/writer.go`), and `composition` (the root.go
grant).

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

*(State as of the ownership handover. The rung went OPEN → hub PASS at `8e37958a` → reopened in
ROUND 3 → hub PASS at the final code tip `2921d563`. Final state is two sections below.)*

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

**~~Synchrony check (hub asked, chip verified by tool rather than asserting).~~ SUPERSEDED — the
claim below was true when written and is FALSE at the reviewed tip. Struck through in place rather
than rewritten, because how it went stale is the point of ROUND 3.**

> ~~The hub measured at `8e37958a`; the branch is now at `b954783d`. `git diff --name-only
> 8e37958..b954783` returns exactly one path — this `EVIDENCE.md`. Zero `.go` files … The
> differential therefore still holds and **no re-run is needed before merge**.~~

**Correction (ROUND 3).** That was verified against `b954783d`. S7, S8 and S9 landed afterwards.
At the reviewed tip the same command tells the opposite story — run by the chip, pasted:

```
$ git diff --name-only 8e37958a..93c9033052f896712a209cbfb35016ea0ea2c5f1 -- apps/server_core
apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go
apps/server_core/internal/modules/product_links/application/generation_integration_test.go
apps/server_core/internal/modules/product_links/application/generation_service.go
apps/server_core/internal/modules/product_links/application/generation_service_test.go
apps/server_core/internal/modules/product_links/application/resolution_service_test.go
```

Five code files, `generation_service.go` among them — the very file the governance lane would scan.
**The governance differential does NOT cover the reviewed tip, and "no re-run is needed before merge"
is withdrawn.** Found by the GPT gate in round 2, verified by the chip with the command above before
being accepted.

Consequence, stated rather than smoothed over: **the governance rung was OPEN at the final tip.** It
was not counted as PASS. The rung is hub-owned by ruling R-b, so the chip could not close it itself;
a `REQUEST` for a re-run at the final code tip `2921d563` went to the hub. The earlier differential
result (53 == 53 at `8e37958a`) stands as evidence about `8e37958a` and nothing later.

### Rung CLOSED — hub re-ran the lane at the final code tip

Hub evidence, committed to `main` @ `7c54bef`: `_chip-anchors/hub-governance-lane.md`, updated.
Cited here, **not re-proven** — same ownership as before. Identical method: clean detached worktree
outside `.claude/worktrees/`, same BaseSha `917f7bb5…`, measured at `2921d563` — the SHA that gets
merged.

| Run | Result |
|---|---|
| code tip `2921d563` | `status=failed`, exit 1, **53 violations**, 175 lines |
| vs BASE-SHA output | **IDENTICAL** |
| vs the first measurement at `8e37958a` | **IDENTICAL** — nothing in S7/S8/S9 touched a surface the lane measures |

**Hub verdict: C11's governance rung PASSES on the differential reading, now anchored to the SHA
that will be merged.** The honesty note above still applies in full — the lane is red on `main`, was
red before this chip existed, and 53 violations remain the hub's outstanding debt.

**The premise was dead and the conclusion survived anyway — which is why this is recorded and not
quietly dropped.** The re-run confirmed what the withdrawn synchrony claim had asserted. That does
not retroactively make the claim evidence: it was proven between two PACK commits, so it never
covered S7/S8/S9 regardless of how the re-run came out. The hub named the defect in its own
instruction (`<sha medido>..<sha a mergear>` is the only valid window) and made the fix a standing
rule rather than a request: **a hub-run lane rung re-measures at the final tip by default — a
measurement is valid for a SHA, never for a branch.** Carried upstream with the C11 amendment.

## Dual gate

Two reviewers on a fixed-SHA diff, dispatched concurrently and blind to each other. **Round 1 ran on
`917f7bb5..1df627f8`; round 2 on `917f7bb5..93c90330` over code tip `2921d563`** — the SHA in this
paragraph is round 1's, and every round-1 finding below is scoped to it. The GPT side (`gpt-5.6-sol` / medium, OS-process, `--sandbox read-only`) and a COLD Opus side
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

### ROUND 1 — both gates FAIL

| Gate | Verdict | Headline |
|---|---|---|
| GPT side (`gpt-5.6-sol`/medium) | **FAIL** | C4 + C11 marked FAIL; 2 BLOCKING (stale pack tail; R5 gap on reachable guards); ruled 2 of 4 declared open items blocking |
| Cold Opus (`harness:gate-reviewer`) | **FAIL** | C3 FAIL on form, C10 FAIL as unevidenced, C11 FAIL as unrun; 2 BLOCKING (missing evidence blocks; **false provenance claim**); 1 SHOULD-FIX, 3 NIT |

The asymmetry designed into the briefs paid off exactly where predicted — on the declared open
items, where the two gates split:

| Declared item | GPT (told to rule) | Opus (found it cold) | Resolution |
|---|---|---|---|
| `pol` unbounded prefix | acceptable — pre-existing, classified, correctly referred out | acceptable — **independently re-derived** that base already has unbounded `pol`, and that canonicalisation does not worsen the verdict (762 vs 1016 mm diverge exactly as `30pol` vs `40pol` did) | **agreed, deferred** |
| A8 duplicate contradictory reasons | **BLOCKING** — subset declarations are valid data, so the generic path must not emit contradictions | acceptable — dormant, R1 forbids the 4th-enum fix, cases distinguishable via `Detail` | **DISAGREEMENT → hub ruling requested** |
| unreachable parse fallback | **BLOCKING** — honest disclosure does not satisfy R5; remove or restructure | acceptable — and **proved** unreachability: every token the alternation emits reduces to `\d+([.,]\d+)?`, which `big.Rat.SetString` always accepts | **DISAGREEMENT → hub ruling requested**; chip leans Opus, which demonstrated where GPT asserted |
| governance rung | acceptable handling, correct ownership escalation | acceptable as a deferral, blocking as a criterion | agreed on handling; resolved separately by hub ruling R-b |

**Chip's own dispatch defect, recorded because it corrupts one gate's input.** The GPT gate's first
BLOCKING — that the pack tail is placeholders — is TRUE at `8e37958`, the sha the chip pinned it to.
Both gates were dispatched while the criteria and ladder sections were still uncommitted in the
working tree. The Opus gate read the working tree and saw them; the GPT gate read the pinned commit
and did not. That finding is an artifact of the chip's dispatch, not a defect in the pack as it
stands, and the re-gate runs against a frozen tip. Its OTHER blocking finding (the R5 gap) was
derived from code and is real.

**Findings accepted and fixed this round:** the four missing evidence blocks; the false provenance
claim; the C6 golden-pair inversion; the R5 coverage gap on three reachable guards (S7); the
clamped-formula hole in the limits table (S7).

**Findings accepted as accurate but not acted on:** the Opus NIT that a non-nil EMPTY declaration is
accepted where `nil` errors. Correct, and correctly classed non-blocking — it is silence rather than
a fabricated fact, the connectors layer deliberately makes empty legal
(the empty-declaration sub-case of
`TestMarketplaceCapabilityServiceIdentityAnchorsRequiresExplicitDeclaration`), and the real adapter
cannot produce it
(`IdentityAnchorAdapter` always returns the full vocabulary). Recorded for the hub.

**Open request to the hub:** the Opus gate's SHOULD-FIX that the diff touches provider-adapter scope
(the ML `capability_adapter.go`) with no `LIVE-VERIFIED:` / `LIVE-WAIVED-BY-OPERATOR:` marker. The
change is a static declaration with no provider I/O and U1–U3 are hub-run, so a waiver appears
right — but the gate's point stands that it must be recorded rather than inferred. Waiver requested;
not self-granted.

### Hub rulings on the two disagreements — R-6, R-7, R-8

Both gate disagreements went to the hub with evidence rather than being decided in-chip
("disagreement = BLOCKED with evidence, never a unilateral decision"). The hub ruled at
`main@56de5b8`, filed in `_chip-anchors/hub-rulings.md`. The hub went AGAINST the chip's stated lean
on both, and both rulings stand.

| Ruling | Disagreement | Hub's decision | Chip's lean before the ruling |
|---|---|---|---|
| **R-6** | A8 duplicate contradictory reasons | GPT side upheld. It is a **dedup/precedence bug in the generation loop**, not a wire-shape problem. Rule: at most one reason per anchor per candidate, precedence deterministic and tested. Impossible without a wire change ⇒ BLOCKED to the hub, not a chip decision. | leaned Opus (dormant, distinguishable via `Detail`) — **overruled** |
| **R-7** | unreachable parse fallback | GPT side upheld in outcome, Opus in method: **keep** the fail-closed fallback, and add a **direct seam test on the helper** — call the parser with a value the current alternation cannot produce and assert fail-closed. | leaned Opus (unreachability proved, disclosure sufficient) — **overruled** |
| **R-8** | live-verification waiver on the ML `capability_adapter.go` | **Waiver DENIED and not needed.** The chip records `rung LIVE: hub-run, U1-U3, ruling R-8` as OPEN on its side; the hub fills it at acceptance once U1–U3 pass. No merge before that. | waiver requested, not self-granted — **request refused, handling upheld** |

Why the chip's lean was wrong on R-6, stated plainly because the reasoning error is the reusable
part: the chip weighed the defect as *dormant* (Mercado Livre supplies all three anchors, so no
production candidate hits it today) and let that dormancy carry the decision. Dormancy is a fact
about today's only provider, not about the code. F-01 exists precisely to make the anchor set
provider-declared DATA — so a subset declaration is not an exotic input, it is the feature's whole
point, and the first provider that declares a subset would have shipped contradictions. Worse, the
chip's own `generation_integration_test.go` fixture already declares a subset (`seller_sku` only), so
the defect was live in the chip's own fixture and simply unasserted. "Dormant" was measured against
the wrong thing.

R-6 and R-7 are implemented in slice S8. R-8 is recorded on the ladder, not implemented.

### ROUND 2 — dispatched on the frozen tip

Round 2 was deliberately HELD while S8's regression was open, rather than dispatched against a tip
the chip already knew to be defective. The hub endorsed holding it. Both gates were dispatched only
after S9 landed, the chip verified it independently, and the tree went clean.

Frozen tip: the pack commit recorded below. Both gates run blind to each other, both over code AND
pack at that one SHA, with R-6a, R-7 and R-8 stated as authority rather than open questions — and
with the instruction that a ruling in the chip's favour is not evidence the code is right.

Frozen tip: `93c9033052f896712a209cbfb35016ea0ea2c5f1` (pack) over code tip `2921d563`. Both briefs
were amended before dispatch: the round-1 briefs still stated **R-6**, which S9 no longer implements.
Dispatched unpatched, both reviewers would have failed the chip for keeping same-anchor
`FOR`+`AGAINST` — the exact behaviour R-6a now mandates. Both were re-pinned to R-6a, told to hunt
the DROP direction (a collapsed `FOR`+`AGAINST` pair; a declaration-unsupplied anchor losing its
reason), and told to run the ordering test at `-count=10` because the map-iteration mutation only
reddened on run 7 of 10. Recorded as an upstream candidate: **re-dispatching a gate after a
corrective re-pins BOTH the tip and every ruling the corrective changed**, not just the SHA.

| Gate | Verdict | Where it failed |
|---|---|---|
| COLD Opus (`harness:gate-reviewer`) — D19 | **FAIL** | the PACK, explicitly not the behaviour: "the code is in materially better shape than the pack… The chip fails on evidence, not on behaviour." Six stale `file:line` citations; the criteria table still pinned to `1df627f8`; C11's synchrony claim now false; must-fail rows asserted without artifacts; the obsolete R-7 "no proof possible" paragraph; the inverted `EXEMPLO-IO` header |
| GPT side (`gpt-5.6-sol`/medium) — D20 | **FAIL** | C6, C8 and C11 FAIL; 3 BLOCKING, verbatim in `scratchpad/agent__p6-sol-r2.last.md`. C1–C5, C7, C9, C10, C12 PASS |

**Both gates AGREE, and that is the reconciliation.** Neither found a defect in the code. Neither
disputed R-6a's implementation, the finalizer, the unit canonicalisation or the independence of the
three limits — the Opus side having been given no open-items list and having had to find its own
targets. Every blocking finding on both sides reduces to one shape: *the pack's evidence is pinned to
a tip that has moved.* There is nothing to reconcile between them; there was a defect to fix, and it
was mine.

Concrete claims the GPT side made that the chip re-derived rather than accepted:

- `git diff --stat 8e37958a..93c90330 -- apps/server_core` → **five** changed code files. Verified;
  this is what falsified C11's synchrony claim, and the governance rung was reopened because of it.
- the full final diff has **17** paths, not the 15 the C10 cell called "the full final diff".
  Verified: 15 code + 2 pack.
- `go build ./...` returning `error obtaining VCS status: exit status 128` — a worktree/VCS-stamp
  condition, profile §3 class, not a compile failure; `go vet ./...` and the 27-package L1 lane are
  green at the same tip and are what the ladder rung actually rests on.

Both FAILs triggered profile §11 for the **third** time on this chip, on the shape recorded in
ROUND-3 above. Per hub ruling R4 the patching stops there: the class was swept, and the
structural-vs-informational question was sent to an independent judgement (D21) rather than
self-adjudicated. **No third gate round is claimed.** The chip goes to the hub with both round-2
verdicts standing as FAIL-on-evidence, the sweep that answers them, and the hub's own decision on
whether the re-anchored pack now clears — which is the hub's call, not a verdict the chip may issue
about itself.

> **ROUND-4 correction to the sentence above.** It originally read "the class was swept
> exhaustively, every citation re-derived at the code tip". The sweep covered the CRITERIA TABLE
> only; citations in the narrative, the ledger and the review sections were untouched, and ten of
> them were stale. Round 3's cold gate found that, and the scope word "every" is what made it a
> false claim rather than a partial one. Corrected in place, and the class it belongs to is now
> named in the sweep table as *CLAIMED RE-RUN, NOT RE-RUN*. This is the sentence that produced
> ruling R-12 — a closure sentence names the scope swept and the tool that produced it, or it does
> not exist. The R-12-conforming replacement is under *Class closed — SCOPE AND TOOL NAMED*.

### ROUND 3 — dispatched on `ac72eb82`, code tip `2921d563` unchanged

The code did not move between rounds 2 and 3. Only the pack did — which is itself the finding, and
the reason the round-3 briefs carried a two-layer instruction: layer 1, confirm the re-anchoring is
real by sampling citations from at least FIVE different criteria (a sweep that fixed C6/C7 and left
C1 stale is a failed sweep); layer 2, stated as mattering more, **verify the code anyway**, treating
"round 2 said the code was fine" as worth nothing. Hub-closed items (the governance rung at
`main@7c54bef`, C2 under the hub's unidirectional invariant) were fenced off so the gates would not
spend their attack budget on settled ground. Both briefs were re-pinned to the RULINGS in force, not
only to the SHA — upstream candidate #5 applied to myself, at the hub's instruction.

| Gate | Verdict | Where it failed |
|---|---|---|
| COLD Opus (`harness:gate-reviewer`) — D22 | **FAIL** | on the pack again, and on the sharpest possible target: the ROUND-3 closure sentence claiming "every citation re-derived at the code tip" when only the criteria table had been swept. Ten stale citations outside that table, all subsequently verified by the chip and corrected |
| GPT side (`gpt-5.6-sol`/medium) — D23 | **FAIL** | **C5 FAIL** — the pasted artifact still pinned to `1df627f8` while the row claimed a re-run, and the row naming `openapi/`, which does not exist. Plus two more BLOCKING: the false colour-branch classification, and the contract path being unreachable from the frozen tree. C1–C4, C6–C12 PASS; every re-anchoring sample resolved |

**Both gates agree for the third time that the code is not the problem.** Neither round-3 reviewer
disputed the finalizer, R-6a's implementation, the unit canonicalisation, or the independence of the
three limits — and the GPT side re-derived C6's and C8's must-fail mutations itself rather than
reading the pasted outputs. What failed, all three rounds, is the pack's account of its own
verification.

Two round-3 findings were checked and only PARTLY conceded, because a gate is not authority either:

- The GPT side reported the contract file missing from `git ls-tree ac72eb82`. Correct — and the pack
  already said so at "Note on the pack location": *"They arrive on this branch through the hub's
  merge, not through the chip's diff."* The gate had read past a disclosure that was already there.
  Its derived fix was still worth taking: the header's `CONTRATO:` line now carries the
  `main@f81b8975` pointer so it is self-sufficient. Finding refuted as stated, remedy adopted.
- The colour classification, by contrast, was a real chip defect and is corrected in the C7 cell.

### ROUND 4 — dispatched on pack tip `01198a0f`, code tip `2921d563` unchanged (D25)

Both briefs re-pinned to the rulings in force — R-6a, R-9, R-10, R-11, R-12 @ `main 371c91d` — not
only to the SHA. The cold brief was re-aimed at SEMANTICS per R-11: the mechanical half is the
script's job now, so the reviewer was told to test whether claims are TRUE, not whether citations
RESOLVE.

| Gate | Verdict | Where it failed |
|---|---|---|
| COLD Opus (`harness:gate-reviewer`, Read/Grep/Glob, no Bash) | **FAIL** | **C5 NOT-EVIDENCED** — the R-10 witness pastes four paths under a tree-listing label; at least fourteen tracked files live under those two trees. **R-12 violated again** — the closure sentence claims every `file:line` citation while `cite-audit.py`'s regex requires a filename, so bare `:NNN` citations are structurally invisible to it; two survivors demonstrated. Plus SHOULD-FIX on the `55`-vs-`65` count and a stale-but-disclosed NIT |
| GPT side (`gpt-5.6-sol`/medium) | **FAIL** | **C10 FAIL** — the "full list" pastes 17 paths, `git diff --name-only 917f7bb5..01198a0f` returns 19; the two omitted are `cite-audit.py`/`cite-audit.txt`, the round's own artifacts. Plus five more BLOCKING: the `55`/`65` count, the false multiplicity claim, the surviving R-12 overclaim in the base-SHA paragraph (the gate cited it by pack line; pack lines are unstable, so it is named by content here), and **two on the code** — three guards touched by this chip with no must-fail row |

**This is the round the dual gate earned its cost, and R-18 is the hub's reading of why.** The two
gates blocked on almost DISJOINT sets. Cold — which cannot run a command — found C5's witness by
globbing the tree and comparing, and found the regex hole by READING the script; GPT — which can —
found C10 by running the diff and the guards by grepping for mutation rows. Neither found the
other's. *Agreement between gates looking at the same thing is worth less than disagreement between
gates looking at different things.*

**Both gates cleared the code, for the fourth time — and cold said so explicitly**, having attacked
every target the brief named: R-6a rules 1/2/3 all hold, ordering is slice-derived with maps used
only for lookup, the `127/5` factor is exact, the suffix order is longest-first so `m` cannot shadow
`mm`, fail-closed persists nothing and auto-approves nothing, and the three limits are independent
rather than coincidentally equal. What FAILED in round 4 is again the pack — except for the R5 gap,
which is the first CODE-adjacent finding since round 1 and is why the round produced a commit.

**The chip verified both blockers before conceding, and both were real.** The tree listing at the
frozen tip returns 14
paths, not 4 — the paste was truncated with no marker. A line cited as a failure return was in fact
`resolveIdentityAnchors`'s SUCCESS return; and the declaration skip sat one line OUTSIDE the range
that claimed to cover it, so the citation excluded the very line it was offered to prove. Corrected
in place, and the R-16/R-17
remedies that followed are structural rather than instance-level.

**One item the cold gate raised and explicitly declined to score**, recorded here rather than
dropped: `hardNegativeDimension` sorts its pairs with `slices.SortFunc`, which is documented as NOT
stable, then `CompactFunc` keeps the first of each equal-key group. A title carrying one measurement
in two units yields two pairs with equal `key` and different `display`, so which `display` survives
is unspecified by the API. The comparison key — hence every verdict — is unaffected; only the
operator-facing detail string is at stake. The gate could not run it, so it filed it as a SUSPICION
rather than a finding, which is the honest shape. **Referred to the mission backlog beside
`refforn`**, not fixed here: it is pre-existing, it changes no verdict, and a chip on its fifth gate
round does not widen its own diff.

### ROUND 5 — dispatched on code tip `85b6c367`, gates split by AREA per R-18 (D26, D27)

R-18 was applied deliberately for the first time: the two gates were given DIFFERENT briefs rather
than the same one twice. Cold took SEMANTICS — does the prose match the line, and is the behaviour
in the production code. GPT took DERIVATION — re-derive every count and listing from the frozen
SHAs and attack the pointer discipline.

| Gate | Area | Verdict | Where it failed |
|---|---|---|---|
| COLD Opus (`harness:gate-reviewer`, Read/Grep/Glob, no Bash) | SEMANTICS | **FAIL** | On the pack, not the code: citations whose sentences did not match the line they resolved to. The code cleared again — F-01, F-02 and F-03 read as real, tested, and their guards load-bearing |
| GPT side (`gpt-5.6-sol`/medium) | DERIVATION | **FAIL** | Four BLOCKING, all structural: the audit's coverage total counted unresolved rows as resolved; unpinned counts had decayed; a count derived at a FROZEN SHA was simply wrong when typed; and the 34-path listing named a command that reads the index and can bind no tree-ish at all |

The derivation side's own summary, quoted from its artifact:

```quoted-verdict
# VERDICT: FAIL

The four required listings and SHA-integrity checks reproduce exactly. However, the derivation
structure still fails:

- The citation audit's `222 resolved, 0 unresolvable` claim is false.
- Several `HEAD`-relative counts have decayed.
- A frozen-SHA grep claims 86 hits but returns 82.
- The 34-path listing uses `git ls-files`, so its stated frozen-SHA binding is not encoded in the
  command.
```

Its four completeness comparisons all reproduced empty, which is the R-16 remedy holding. What it
killed was the layer above: the pack's account of its own coverage.

**Both gates FAILED and neither failed on the code**, which is the fifth consecutive round of that
result and the reason the chip escalated structure rather than shipping a sixth patch.

#### The cold verdict has no artifact, and that is a defect in the ledger, not a detail

The GPT side dispatches as an OS process and writes `agent__p6-sol-r5.last.md`, so its verdict above
is quoted from a file the hub can open. The cold side dispatches through the Agent tool, whose
per-task output file is EMPTY — all twelve agent runs on this chip wrote zero bytes — so a cold
verdict exists only in a session transcript, which compaction is entitled to drop.

So the cold row above is the chip's account of a verdict it received, not a quotation from an
artifact, and it is marked as such because the reader cannot otherwise tell the two apart. That
distinction is the whole subject of this chip: a claim and a checkable claim are different things
even when both are true.

**Remedy, applied from round 6 on:** the chip transcribes a cold verdict into a pack file the moment
it arrives, and the ledger row points at that file. The reviewer cannot do it — `harness:gate-reviewer`
has no Write tool by construction, which is what makes it cold — so the write is the chip's step, and
naming it as the chip's is what keeps it honest.

### ROUND 6 — dispatched on pack tip `a23aee3a`, code tip `85b6c367` unchanged (D28, D29)

Gates split by area for the second time, and further apart than in round 5: cold took SEMANTICS
against a prose rewritten wholesale under R-22, GPT took DERIVATION **and the apparatus itself** —
because the pack's coverage claim now rested entirely on `cite-table.py`, and five rounds had died on
checks that were narrow and silent about being narrow.

**The gates did not agree, and no AGREEMENT string was authored.** Both verdicts went to the hub raw
and verbatim, and both are committed here in full, one file each, unedited inside a fence:
[verdict-r6-cold.md](verdict-r6-cold.md) and [verdict-r6-gpt.md](verdict-r6-gpt.md).

| Gate | Area | Verdict | Where it landed |
|---|---|---|---|
| COLD Opus (`harness:gate-reviewer`, Read/Grep/Glob, no Bash) | SEMANTICS | **PASS** | No BLOCKING semantic mismatch. Sampled ~25 coordinate rows plus all 20 backticked `Test*` names; every one resolved to the stated file AND line, production symbols to production files — the round-5 defect class is gone from what it checked. One SHOULD-FIX, one NIT |
| GPT side (`gpt-5.6-sol`/medium) | DERIVATION + THE APPARATUS | **FAIL** | Five BLOCKING. Every frozen-SHA count re-derived correctly, including the round-5 corrections. Two blockers were pack errors; three were the tool's claims about its own coverage |

**The two pack errors, re-derived by the chip before relaying and by the hub before ruling.** A
frozen-SHA grep sentence claimed an empty result while omitting the test-file exclusion that makes it
empty — and the pack disclosed those same two test occurrences elsewhere, so the document contradicted
itself. And a write-set range labelled S1+S2 ran two commits past S2, returning the right count for
the wrong reason because the intervening slice touched only paths already listed. Both are corrected
above, at the rows that carried them, with the coincidence stated rather than quietly swapped out.

**The three apparatus blockers are the subject of R-24** and are not fixed by widening anything; see
the `cite-table.py` section under the ladder for what replaced them.

#### The convergence, which is the round's most useful result

Both gates found the silent anchor drop INDEPENDENTLY, from opposite directions and by methods that
do not touch: cold reading symbols and noticing two load-bearing sentences whose anchors were absent
from the table, GPT attacking the tool and proving an injected name changes nothing. The hub's
addition to R-18 came from that: **convergence across differently-scoped gates is the strongest signal
a dual gate can emit**, stronger than agreement between gates looking at the same thing, because the
two lenses share no method and the finding cannot be an artifact of either.

Note also the SEVERITY split — cold filed SHOULD-FIX where GPT filed BLOCKING. The hub's reading is
that this is not miscalibration: each gate scored the finding against its own scope, both were right
about their own half, and only the hub sees both.

#### What round 6 proved about round 5's remedy

R-22 moved coordinates out of the prose. Rounds 2 through 5 all failed on those coordinates decaying;
**round 6 has none of that class**. It failed instead on the tool R-22 created to enforce the rule —
real progress and a real new defect, and the hub declined to relabel either half as the other. What it
refused was a seventh round on the same reasoning: *"each round fails somewhere new" IS the infinite
regress.* The stop is structural, not a resolution — the rung is off the ladder, so there is nothing
left there to fail.

### ROUND-6 RE-CHECK — the gate that filed the FAIL, on its own five (D32)

Not a seventh round, and the hub fixed the scope in its own words: *"Whoever filed the FAIL re-checks
its own five blockers at the corrected tip ... by the gate that delivered it, scoped to exactly what
it filed. There is no new surface, because the corrective REMOVES claims rather than adding them."*
The cold gate did not re-run — *"it passed; its area is untouched by the corrective; re-running a
passing gate over an unchanged area is the redundancy R-18 already devalued."* That licence holds only
if the corrective touched no code, so the chip verified the condition instead of asserting it:
`git diff --name-only 85b6c367 56598bea -- apps contracts packages` returns nothing, and the gate
re-derived the same zero independently.

The verdict is transcribed in full, unedited inside a fence:
[verdict-r6-recheck.md](verdict-r6-recheck.md).

**Outcome: FAIL.** Three of five closed; two stand.

| # | The blocker as filed | Re-check |
|---|---|---|
| 1 | The `connectorsports` frozen-SHA sentence claimed an empty result without the exclusion that makes it empty | **RESOLVED** — both forms re-run at the frozen tip and at the control; the corrected sentence is true and no longer contradicts the pack's own import table |
| 2 | The S1+S2 write set ran two commits past S2 | **RESOLVED** — the re-pinned range covers exactly S1 and S2, returns the stated 13, and the gate confirmed the two sets are equal, which is the coincidence the pack now states out loud |
| 3 | The tool silently discarded anchors it could not resolve | **DISSOLVED-BY-R-24, and the gate accepted that as closure** — its own injected name and call span both surface now, and `--strict` still exits 0 with them present |
| 4 | The mutable-axis pass called itself a ban and is not one | STOOD at `56598bea` — **closed under R-25 by deletion**, hub-verified absent |
| 5 | The prose scan called itself a SUPERSET and is not one | STOOD at `56598bea` — **closed under R-25 by deletion**, hub-verified absent |

**Both survivors were one defect: the corrective removed the totality claims from what the tool
*prints* and left them standing in what the pack and the implementation *say*.** Three holders of the
same claim, one of them reached. The gate's verdict, transcribed unedited, carries the exact strings
and coordinates; they are not restated here, because a false sentence quoted as a live description is
the thing being removed.

**R-25 is the rule this round produced, and it corrects the chip.** The first draft of this section
ANNOTATED the two false sentences as false and left them in place, reasoning from ADR-17. The hub
rejected that instrument: *honest-unknown is for gaps, not for falsehoods. You disclose what you do
not know. You DELETE what is wrong.* Annotating leaves two sentences where one belongs and makes the
reader arbitrate — and whoever meets the false one first gets no signal to keep reading. R-24 said the
false claim is removed, not excused, and **disclosure is a way of excusing**. ADR-17 applied correctly
to the tool's coverage, which is a gap; one step too far to two sentences, which were simply wrong.

**§11 did not fire, and the reason is the distinction the gate drew:** *"The two blocking survivor
lines already existed in `f1397cf7^`; they were omitted by the corrective rather than introduced by
it."* §11 fires on a remedy APPLIED and a defect RECURRED in a new disguise. This was incomplete
application, not recurrence — and the two want opposite responses: recurrence wants a better remedy,
incompleteness wants the remedy finished. The hub recorded a second, independent tell: **§11's own
remedy is more apparatus, and R-24 had ruled one ruling earlier that the answer is not more
apparatus.** A rule whose remedy contradicts the ruling governing the round does not fire in it. The
hub recorded this explicitly rather than quietly, because *"the rule is inconvenient, so I skip it"* is
the reasoning refused all chip long.

**The remedy was shaped so a fresh falsehood could not enter it.** The comment was DELETED, not
reworded — a deleted comment cannot be false, and the regex beneath it is the truth and needs no
gloss. The two sentences were re-aimed at what the report PRINTS, a target the gate had already
certified accurate (*"correctly discloses all five holes"*), so the rewrite answered to a fixed
external standard rather than to the author's judgement. And the disclosure paragraph left with the
falsehoods it disclosed: with the claims gone there was nothing to disclose, and keeping it would have
made the pack describe a state it is no longer in — the decay class, one last time.

The FOREIGN-fence SHOULD-FIX survives untouched and the gate says so rather than counting it against
the corrective: a pack author can still exempt authored text by labelling it `quoted-verdict`, because
the tool checks the tag and never the provenance. Out of the corrective's scope, still true, still
open.

### P6 DISCHARGE — the gate closes, and how each finding closed

**AGREEMENT — P6 discharged.** The cold gate returned PASS on SEMANTICS at `85b6c367` and did not
re-run, because the corrective touches no code and `git diff --name-only 85b6c367..<final tip> -- apps
contracts packages` returns nothing — verified by the chip, re-derived by the GPT gate, and verified
again by the hub at the final tip. The GPT gate's five DERIVATION blockers are all discharged, three by
the gate itself and two by hub observation after a remedy fully determined by the finding.

**This sentence does not stand alone, and that is deliberate.** A closing line that omits HOW the close
was reached is round 4's silent-omission class committed in the last line of the document. The ledger
is the sentence's other half:

| # | Finding | Closed by | Means |
|---|---|---|---|
| 1 | `connectorsports` sentence claimed an empty result without the exclusion that makes it empty | **the GPT gate**, at the corrected tip | Re-derived both grep forms at the frozen tip and at the control. Sentence true, no longer contradicts the pack's own import table |
| 2 | S1+S2 write set ran two commits past S2 | **the GPT gate**, at the corrected tip | Range re-pinned; `SETS_EQUAL` confirmed; the coincidence stated in the pack rather than the SHA quietly swapped |
| 3 | The tool silently discarded anchors it could not resolve | **the GPT gate**, which had filed it BLOCKING | DISSOLVED-BY-R-24 and accepted as closure by its own filer: injected name and call span both surface, `--strict` still exits 0 with them present, and the non-fatality read as intended rather than as an oversight |
| 4 | The mutable-axis pass called itself a ban | **HUB OBSERVATION, not a gate** | Sentence re-aimed at what the report prints — a target the gate had already certified accurate — then absence verified by string |
| 5 | The prose scan called itself a SUPERSET | **HUB OBSERVATION, not a gate** | Comment DELETED, not reworded; absence verified by string |

**The honest limit, in the hub's words and not softened:** *GPT does not re-run and does not certify
the final tip. I hold (a) its verdict that these two lines are the ENTIRE remaining defect and that the
corrective introduced none, and (b) my own observation of absence. I DO NOT HOLD A MODEL-SIDE VERDICT
ON THE FINAL TIP.* That gap is accepted deliberately: the alternative is a round whose entire content
is confirming three strings are absent, and that round would need another if it found a typo — the
regress R-24 ended.

**One string class survives the hub's grep on purpose.** The exact false sentences still appear inside
`verdict-r6-recheck.md`, within the `quoted-verdict` fence, because that is the GPT gate's text
transcribed verbatim under R-23 and the chip may not edit it. A false claim quoted as the report of the
gate that found it is not the pack asserting it. The chip's own preamble around that fence was
rewritten, since it *had* decayed into describing a state no longer true — the same class, caught on
the same day it was ruled.

Backlogged, not blocking: the FOREIGN-fence SHOULD-FIX. A pack author can still exempt authored text
by labelling it `quoted-verdict`, because the tool checks the tag and never the provenance. It is a
provenance hole in an artifact that no longer gates anything, which is the same reason the other three
dissolved.

## The rulings rounds 4, 5 and 6 produced — R-13 through R-25

R-9…R-12 were issued *before* round 4 ran. What round 4 then found forced nine more, in three
batches, round 5 forced R-22, and round 6 forced the last two. They are transcribed here because the
reasoning is the durable half. Each is the hub's, not the chip's; where the chip disagreed the hub's
text is what stands.

| Ruling | What it establishes |
|---|---|
| **R-13** | Code findings land. R-11 is ratified as doctrine — and the hub's own account of WHY it worked is the part worth keeping: the guards surfaced because reviewer attention "stopped being spent on coordinates". **Having a test and having a must-fail are different properties**, and three guards on this chip had the first without the second |
| **R-14** | **A count is a coordinate.** Verbatim: *any derived fact the pack COPIES from an artifact generated out of the pack is self-invalidating — line, count, multiplicity or list alike.* Two approved mechanisms: (1) **point instead of copy** — name the artifact, let it hold the number; (2) **scope to an axis the remaining work cannot move**. Discriminator, from the hub: `14` was legitimate because it came from a tree listing at a frozen CODE tip, which pack edits cannot touch |
| **R-15** | "Patch the numbers and freeze" — **DENIED**. Verbatim: *A freeze is discipline… A pointer is structure.* And the stopping condition: round 5 is the last of this shape; if it still finds a decayed number, **the structure failed, not the execution** |
| **R-16** | The C5 disagreement SPLITS rather than resolving to one side. GPT was right that R-10 is discharged — a witness proves EXISTENCE, and one file does that. Cold was right that pasting 4 of 14 paths is **silent truncation**, violating a separate and older rule (precedent: CHIP-MERCADO's page-1 truncation). Remedy scope: *every pasted listing in the pack is complete, or declares what it samples* |
| **R-17** | On the audit regex missing bare citations: *the defect is not its narrowness, it is its silence.* Two requirements — (1) the script reports what it could not resolve, (2) **bare `:NNN` citations are BANNED from the pack**. And the round's class is named: *all of them are silent omission inside a proof artifact… every proof artifact declares its own coverage* |
| **R-18** | The dual gate delivered **area, not redundancy**. Verbatim: *Agreement between gates looking at the same thing is worth less than disagreement between gates looking at different things.* Round 4's two gates blocked on disjoint defects; that is the gate working, not the gate quarrelling |
| **R-19** | On the corrective that decayed 40+ citations: *the remedy was never wrong. Its precondition never held.* "Final tip" was ASSERTED, never verified. The two remedies PARTITION — pack-derived facts get a **pointer**, code-derived facts get **sequencing** (land all code, freeze, then regenerate). Not competing options; different halves |
| **R-20** | **An inference must be visible as an inference.** The audit's bare-citation pass prints the file it inferred and the distance it inherited over, so a wrong inference can be READ. It is a **migration tool, not a permanent crutch**: end state is zero bare citations, and a non-zero count after the ban **is** the violation |
| **R-21** | A citation resolving to plausible-but-wrong content is **invisible to reading** — the reader compares the sentence to an expectation, never to the actual line. Therefore the audit prints what **every** citation resolved to. Found live on this chip: a bare line-448 citation inherited the wrong file and landed IN RANGE on `//  2. Load the latest audit entry…`. Written out in words here rather than in citation form, because R-17 bans the form from the pack and the audit would flag a quotation of it as a live violation |
| **R-22** | **An evidence apparatus that costs more to verify than the thing it evidences has inverted its purpose, and the answer to that is not more apparatus.** The hub's finding, and no gate could have made it: each gate sees one round, and the shape only appears across five. Mechanics: (1) prose claims become **behavioural, anchored to test and symbol names**, checkable by running the named test; (2) the coordinate table **stays** and is generated mechanically into a separate artifact; (3) the lane regenerates it and **fails on the diff**. This does NOT overturn the §11 judge, who banned cite-by-symbol as the anchor TYPE — what R-22 bans is COPYING the coordinate into prose. Landing with it: the **mutable-axis ban whole**, grep-checkable and fail-closed; **tool inversion adopted with its role demoted**; **declarations travelling with content ratified**. Round 6 authorised ONCE, gates split by area, in this order: *code first and all of it, then the frozen tip, then the table* |
| **R-23** | On the cold gate's verdict having **no durable artifact**: D30 accepted as written, and the reasoning is the part to keep — *`gate-reviewer` has no Write BY CONSTRUCTION, and that absence is exactly what makes it cold. Granting it Write to solve durability would destroy the property that makes it worth running.* One requirement added: **the cold verdict ALSO travels to the hub VERBATIM in the closing event, and the hub commits its own copy.** Two copies, two holders, two timestamps. The honest limit, in the hub's words: *this does NOT make the transcription verifiable against the original — the original dies when the transcript compacts. What it removes is UNILATERAL CONTROL.* Two general rules kept: **generate and check are separate modes; a checker does not write**, and **a new rung must be green on a clean clone, or it is not a rung, it is training to ignore the lane** |
| **R-24** | **A verification artifact keeps only the claims it can make TOTALLY. A claim it cannot make totally becomes a REPORT, not a gate.** Round 6's three blocking findings were all claims the pack's own tool made about its own coverage — total in the wording, partial in the code. The remedy deletes the claim instead of widening the recognizer, so the findings *cease to exist BY CONSTRUCTION, NOT BY PARDON*. ADR-17's honest-unknown rule turned on our own evidence apparatus. Mechanics: (1) the tool **prints what it could not resolve** — **R-17 ratified a FOURTH time**, with the CONSEQUENCE changed from an exit code to a report line; (2) **the PACK rung leaves the ladder**, so no future round can fail on it; (3) the prose **stops asserting totality anywhere**; (4) **NO WIDENING** — *do not widen a recognizer whose output no longer gates anything*. R-11's mechanism retires with its object: the table is **navigation, not evidence**, so a wrong row costs a reader a wrong jump rather than a false proof. Added to R-18: **convergence across differently-scoped gates is the strongest signal a dual gate can emit**, stronger than agreement between gates looking at the same thing, because the two lenses share no method and the finding cannot be an artifact of either. And on the severity split — cold filed SHOULD-FIX where GPT filed BLOCKING — *both are right about their own half; only the hub sees both scopes* |
| **R-25** | **HONEST-UNKNOWN IS FOR GAPS, NOT FOR FALSEHOODS. You disclose what you do not know; you DELETE what is wrong.** The chip had annotated two false sentences as false and left them standing, reasoning from ADR-17. The hub rejected the instrument: annotating leaves TWO sentences where ONE belongs and hands the reader the job of arbitrating — *and whoever lands on the false one first gets no signal to keep reading*. R-24 said the false claim is REMOVED, not excused, and **disclosure is a way of excusing**. ADR-17 was applied correctly to the tool's coverage, which IS a gap, and one step too far to two sentences, which were not. **§11 does not fire**, and not for convenience: it triggers on a remedy APPLIED and a defect RECURRED in a new disguise, whereas this remedy reached one of three holders of the same claim — *incomplete application, not recurrence*. The two want opposite responses: **recurrence wants a BETTER remedy; incompleteness wants the remedy FINISHED**, and treating the second as the first invents a new remedy for something the old one already solves, *which is precisely how apparatus gets built*. Independent tell that the reading is right: **§11's own remedy is more apparatus, and R-24 ruled one ruling earlier that the answer is not more apparatus** — a rule whose remedy contradicts the ruling governing the round does not fire in it. Recorded explicitly, because *"the rule is inconvenient, so I skip it"* is the reasoning refused all chip long. The remedy was **shaped so a fresh falsehood could not enter it**: the comment DELETED rather than reworded (a deleted comment cannot be false); the two sentences re-aimed at what the report PRINTS, a target the gate had ALREADY certified accurate, so the rewrite answered to a fixed external standard and not to the author's judgement; the disclosure paragraph leaving with the falsehoods it disclosed. And the last mechanic: **when a remedy is FULLY DETERMINED BY ITS FINDING, verification is a check of ABSENCE — one bit — and that wants a SHELL, NOT A SECOND MODEL.** The hub also verified the three strings **by string, not by line**, deliberately: *a gate's coordinate is a claim like any other*, and this one cited a line two off from where the text lived |

**R-24 is where the regress stops, and it stops structurally rather than by promise.** The hub tested
it against its own R-9, which had denied "merge on code-clean verdicts and carry the pack as debt":
there the pack defects STAYED IN PLACE as debt, here they are FIXED — the two prose errors corrected,
the false totality claims removed, nothing carried. What was declined is a further gate ROUND over
the corrected pack, not the correction. And the terminator is not a resolution to do better: the PACK
rung is off the ladder, so there is no longer a rung there for a seventh round to fail on.

The hub's own observation, which neither gate could make because each sees one half:
**neither verdict on the CODE depended on this tool.** Cold reached PASS on C1/C2/C3/C4/C6/C7/C8 by
reading source; GPT reached PASS on C2/C3/C5/C9/C10/C12 by re-deriving from git. Neither used the
table to reach a verdict about the code. All three apparatus blockers were the tool gating itself.

**R-20 and R-21 caught the remedy for R-17 in the act, which is the strongest evidence either of
them is right.** Qualifying the bare citations meant giving each one a file, and the file had to be
INFERRED from the nearest preceding one. That inference was wrong in six places, and none of them
was findable by reading:

- The ROUND-2 verification table's left column identified four rows by a bare number — the S3 note,
  the S1 note, the S2 note, the ROUND-1 CRLF claim. Those were never code lines: they were
  **`EVIDENCE.md` line numbers**, pointing at where in the pack each claim sat. Qualification turned
  four pack-line references into `link_candidate_repo.go` citations, and every one resolved to a
  real, plausible line of SQL or error handling. Fixed by dropping the numbers entirely and naming
  each claim, because a pack line number is unstable for exactly the R-14 reason and should never
  have been a citation in the first place.
- The failure-mode bullet inherited `marketplace_capability_service.go` across a sentence boundary,
  so the success-return line, the empty-code return and the engine-not-configured guard — all three
  in `generation_service.go` — were attributed to the connectors service. Line 76 exists in both
  files: in `generation_service.go` it is the four-disjunct config guard, in
  `marketplace_capability_service.go` it is an `if err != nil` inside a different method. It
  resolves, it reads plausibly, and it is the wrong function in the wrong module.
- Written out in words above rather than in citation form: R-17 bans the bare form from the pack, so
  a QUOTATION of a bare citation is itself flagged as a live violation. Naming the defect without
  reproducing it is the only way the record and the ban coexist.

**And the R-17 fix was silent in exactly R-17's way, one level down.** The first bare-citation pass
matched a WHOLE backtick span, so a span carrying several citations at once — the form the decay
table's C6/C7/C8 rows and the C2 cell all use — matched nothing. The run reported **zero** bare
citations while sixty-one were still there. The pass now scans inside each span. Recording it because
the shape is the round's subject: *the remedy for a silence had a silence of its own, and only a
tool that declares its own rule alongside its own count makes that findable.* A coverage number read
without the rule that produced it is worth nothing — which is why the artifact prints the rule.

Both classes were invisible to review and visible in one pass of the audit's printed bodies. R-21's
requirement — print what every citation resolved to — is what surfaced them; R-20's — print the
inference, not just the result — is what made them attributable to the qualifier rather than to the
author. **The remedy for a silent-omission class produced a fresh silent-omission class, and the two
rulings caught it inside the same round.** Also fixed in the same pass, from the same printed
bodies: the pass-2 range in `appendProviderDeclaredUnavailableReasons` crossed three lines into pass
3 — the cold gate had named this and the artifact had printed `continue` as the range's last line,
which belongs to the following pass — and a capability-service test range opened on a blank line.

**And one the hub stated as a general rule while ruling on a reviewer's finding:** *a reviewer's
claim arrives wearing authority.* Round 3's GPT gate handed the chip a coordinate in
`auto_link_policy_test.go`; the chip used it unverified; the script caught it in the same pass that was supposed to prove the
chip had stopped doing that. Gate output is a claim like any other — the rule the pack already
applied to worker reports, now demonstrated against a reviewer.

## Must-fail proofs (R5)

Every guard this chip touched, with the mutation that must redden it. Each was run and observed, not
reasoned about.

| Guard | Mutation applied | Observed |
|---|---|---|
| C6 unit canonicalisation | disable the mm/cm/m/pol scaling in `normalizeDimensionToken` so tokens compare lexically | RED — **chip-re-run in ROUND 3, output pasted below the table** |
| C6 contradiction still fires | — (inverse check, no mutation) | `50cm` vs `40CM` still REJECTs, so the canonicalisation did not simply disable the rule |
| C8 independent limits | restore the shared `limit*5` derivation (`linkLimit = candidateLimit * 5`) | RED on the 29-link test AND on all four independence rows — **chip-re-run in ROUND 3, output pasted below the table** |
| C8 defaults independent of `Limit` | `linkLimit = candidateLimit*100`, `auditLimit = candidateLimit*500` | new test RED on the `5` and `500` rows, `20` row still green, **and both pre-existing default tests stay GREEN** — proving the old net was blind to the class |
| F-02 display offsets | (proved by the reviewer, not the author) rebuilt the pre-corrective function standalone | returned `display="50c"`, truncated, for `"İnox 50cm"` — the guard is load-bearing |
| **S7** duplicate identity anchor (the `seen[anchor]` check in `MarketplaceCapabilityService`) | delete that check | `error = <nil>` → `TestMarketplaceCapabilityServiceRejectsDuplicateIdentityAnchor` RED |
| **R5-r4** engine-not-configured, `identityAnchors` disjunct (in `GenerateLinkCandidates`) | delete `\|\| s.identityAnchors == nil` from the 4-disjunct config guard | `error = <nil>`, `result = GeneratedCount:0` → `TestGenerateLinkCandidatesRefusesWithoutIdentityAnchorReader` RED. **The mutation is only visible because the test drives an EMPTY snapshot batch**: with an empty batch the resolve loop never runs, so an unwired engine returns a clean zero and the caller reads "nothing to link" instead of "this engine was never wired". A non-empty batch dereferences nil and announces itself — the easy fixture would have hidden the guard's purpose while still going RED for the wrong reason |
| **R5-r4** identity-anchor declaration required (the `capability.IdentityAnchors == nil` check) | delete that check | **NOT an error — an EMPTY declaration.** Observed: `IdentityAnchors() = []ports.IdentityAnchor{}, nil`. The provider then reads as supplying no anchors at all, which renders every anchor `UNAVAILABLE` — the blanket-UNAVAILABLE behaviour F-01 exists to REMOVE. So this guard is load-bearing for the FEATURE, not merely for the type, and its must-fail is the one that says the most: the mutation does not break a contract, it silently restores the bug |
| **R5-r4** unknown identity anchor (the `known[anchor]` check, keyed off `KnownIdentityAnchors`) | delete that check | `IdentityAnchors() error = nil` → `TestMarketplaceCapabilityServiceRejectsUnknownIdentityAnchor` RED |
| **S7** empty provider code (the `providerCode == ""` check in `resolveIdentityAnchors`) | delete that check | `TestGenerateLinkCandidatesFailsWhenProviderCodeIsEmpty` RED — **re-run by the chip itself**, observed verbatim: `error = PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE for provider "": identity anchor declaration is nil, want empty provider code failure`. See the honesty note below — this proof says less than the other two |
| **S7** nil declaration (the `declaration == nil` check in `resolveIdentityAnchors`) | delete that check | `error = <nil>` → `TestGenerateLinkCandidatesFailsWhenIdentityAnchorDeclarationIsNil` RED |
| **S7** clamped limit formula (`resolution_service.go`) | apply `linkLimit = max(2000, candidateLimit*4)` | `limit_5000` row RED (`link limit = 20000, want 2000`) **while `limit_5`, `limit_20` and `limit_500` all stay GREEN** — the contrast is the proof that the old 3-row table was blind to the clamped class |
| **S9** R-6a rule 2, `FOR`+`AGAINST` coexist (pass 2 of `appendProviderDeclaredUnavailableReasons`) | restore blanket per-anchor dedup on the signal side | `TestTitleMatchHardNegativeKeepsTitleForAndAgainstInSeedOrder` RED — **re-run by the chip itself**, observed verbatim: `title reasons=[]domain.LinkCandidateReason{...Anchor:"title", Direction:"FOR"...}, want FOR then AGAINST`. Not confounded: the message names the missing `AGAINST` |
| **S9** R-6a rule 1, `UNAVAILABLE` exclusive | drop `hasSignal` from the declaration skip so the declaration `UNAVAILABLE` is emitted beside a `FOR` | RED — **chip-re-run in ROUND 3, output pasted below** |
| **S9** R-6a rule 3, declaration beats seed | invert precedence so the seed `UNAVAILABLE` is kept | RED — **chip-re-run in ROUND 3, output pasted below** |
| **S9** ordering determinism | range the declarations through a map instead of the slice | RED **on 3 of 10 runs** — **chip-re-run in ROUND 3, output pasted below**. Seven runs passed under a genuinely broken ordering; see the note on `-count` |

**Honesty note on the empty-provider-code proof — it is weaker than the other two, and the
difference matters.** Deleting the duplicate-anchor and nil-declaration guards makes the error go
`<nil>`: those guards are the only thing standing between the input and a silent success, so they
are load-bearing for SAFETY. Deleting the empty-provider-code guard does NOT produce a silent
success — the request still fails closed, because the nil-declaration guard downstream catches it
(no declaration is registered under `""`). What is lost is the DIAGNOSTIC: the operator is told
"identity anchor declaration is nil" for a listing whose actual problem is that it carries no
provider code at all.

So the honest claim is: this guard is load-bearing for the error message, not for fail-closed
behaviour. The test still goes RED when it is removed only because it asserts the
`"provider code is empty"` substring in addition to the error code. Had it asserted the
`PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE` code alone — the obvious way to write it —
it would have stayed GREEN through the mutation and the proof would have been vacuous. That is the
confounded-assertion shape this mission has hit before (CHIP-PED-FILA); the chip re-ran this
specific mutation by hand rather than accepting the worker's report, precisely because the worker's
summary said the mutation "produced a downstream error" without saying whether the test reddened.

**~~One guard has NO must-fail proof~~ — SUPERSEDED by R-7.** This paragraph said the fail-closed
parse fallback could not be proven because it was unreachable through the public alternation. That
was true before S8. R-7 ordered the per-token conversion extracted into `normalizeDimensionToken` so
the fallback became directly callable, and the seam test now exists at
`TestNormalizeDimensionTokenFailsClosed`. It calls `normalizeDimensionToken("abcmm")` — a token the
alternation `hardNegativeDimensionPattern` cannot emit, since every branch requires a leading `\d` —
and asserts `key == "abcmm"`, i.e. the unparseable token is KEPT as the signature key. That is
fail-closed as a *discriminating* key: two titles with different unparseable tokens still differ,
rather than collapsing to equal. The obsolete paragraph survived the R-7 corrective; flagged by the
GPT gate in round 2 and corrected here.

### Must-fail — `cite-table.py`'s four checks, each mutated and observed

**These four were proved while the tool was a ladder rung, and R-24 has since taken it off the
ladder.** The proofs are kept and are not restated as gate evidence: they record that the checks do
what they say, which is now a property of a report rather than of a merge gate. The round-6 derivation
gate independently re-ran all six mutations below against the committed script and reproduced every
exit code, which is the first time a must-fail table in this pack has been replayed by someone other
than its author.

A check that has never failed is indistinguishable from a check that cannot fail. Every check was
disabled or violated one at a time, the exit code observed, and the mutation reverted from a byte copy
taken beforehand.

| Guard | Mutation | Observed |
|---|---|---|
| no coordinate in prose | a sentence citing a real production line, appended to the pack | `coordinates leaked into prose 1`, the offending line printed, **exit 1** |
| no mutable-axis command | a sentence naming a range command that ends at the branch tip | `mutable-axis commands in the pack 1`, **exit 1** |
| no drift in the generated table | one line appended to `coordinates.txt` by hand | `coordinate table DRIFTED`, the difference printed, **exit 1** |
| the table exists at all | `coordinates.txt` moved aside | `coordinate table MISSING`, **exit 1** |

Baseline before and after all four: **exit 0**. Both mutated files were compared byte-for-byte
against their pre-mutation copies afterwards and are identical.

A fifth proof runs the other way — a guard must also NOT fire when it should not, and this one was
about to fire everywhere. Git is configured to convert line endings on checkout in this repository,
so a fresh clone receives the generated table as CRLF while the tool writes LF; the drift guard would
have been red on every machine except the one that wrote it, for a reason having nothing to do with
drift. The comparison now normalizes line endings, and both directions were observed:

| Condition | Observed |
|---|---|
| the whole table converted to CRLF, contents otherwise identical | `coordinate table CURRENT`, **exit 0** — the guard correctly stays silent |
| one line appended by hand to that same CRLF copy | `coordinate table DRIFTED`, **exit 1** — tolerance did not cost detection |

Recorded because the failure mode is social rather than technical. A check that cries wolf is a check
someone turns off, and this pack cannot afford another guard that is technically present and
practically ignored — round 1 of this chip lost a real `gofmt` violation inside exactly this noise.

**A sixth proof, added by the R-24 corrective, replays the exact attack that used to vanish.** The
round-6 derivation gate injected a backticked name that exists nowhere in the tree and observed that
the generated table did not change and the run exited 0 — the drop was silent, which was the finding.
The same injection now, plus a call-suffix span of the kind the anchor recognizer never sees at all:

| Condition | Observed |
|---|---|
| a bare name declared nowhere, injected into the prose | appears under **UNRESOLVED** in the regenerated table; the report count rises by one |
| a `Name()` span declared nowhere, injected into the prose | appears under **NOT TREATED AS ANCHORS**; that count rises by one |
| `--strict` with both probes present and the table regenerated | **exit 0** — the report is a report; it does not gate, which is the R-24 property being proved |

Both files were compared byte-for-byte against pre-mutation copies afterwards and are identical;
baseline `--strict` before and after is exit 0. Note what the third row proves and what it does not:
the drop is now VISIBLE, and it is deliberately not FATAL. Making it fatal would rebuild the totality
claim R-24 struck, one exit code lower down.

The third row of the first table is the one worth reading twice, because the guard did not exist when
the prose describing it was first written. The tool regenerated `coordinates.txt` on every run, `--strict`
included — so a hand-edit would have been silently overwritten and the run would have reported
success. The claim "the lane regenerates and diffs, and drift fails" was true of the intent and false
of the code. Split into two modes, where the checker writes nothing, it is now true of both; and the
only reason the gap surfaced is that writing the must-fail forced the question of what, exactly,
would fail.

### Pasted must-fail outputs — the ROUND-3 re-runs

Five proofs the pack previously ASSERTED without an artifact. Each mutation was applied by the chip,
run, observed, and reverted from a byte copy taken beforehand; the mutation appears in no commit
afterwards and the full suite returns `EXIT=0`.

**C6 — disable unit scaling.** Three tests RED, and the failure text shows the exact
pre-canonicalisation behaviour F-02 exists to remove:

```quoted-output
--- FAIL: TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm (0.00s)
    generation_service_test.go:787: match_status=REJECT, want CONFIRM
--- FAIL: TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate (0.00s)
    ... Detail:"hard-negative: medida/dimensão divergente 50cm≠500MM"}, ... want ALTA/ACCEPT
--- FAIL: TestDimensionCanonicalizationUsesExactMillimetres/decimal_comma
    detectHardNegative("Produto 1,5m", "Produto 150cm")=(true, "…1,5m≠150cm"), want equivalent
--- FAIL: TestDimensionCanonicalizationUsesExactMillimetres/inches
    detectHardNegative("Produto 2pol", "Produto 50.8mm")=(true, "…2pol≠50.8mm"), want equivalent
--- FAIL: TestDimensionCanonicalizationUsesExactMillimetres/metres
    detectHardNegative("Produto 1m", "Produto 100cm")=(true, "…1m≠100cm"), want equivalent
```

**C8 — restore the shared `candidateLimit * 5` derivation.** RED on the 29-link test and on every
independence row, including 5000 — so the guard is load-bearing across the whole table, not one row:

```quoted-output
--- FAIL: TestListLinkWorkflowsUsesIndependentDefaultLimitsAndReturnsAll29Links (0.00s)
    resolution_service_test.go:391: link limit = 100, want 2000
--- FAIL: TestListLinkWorkflowsDefaultsDoNotVaryWithLimit/limit_5
    resolution_service_test.go:439: link limit = 25, want 2000
--- FAIL: .../limit_20    link limit = 100, want 2000
--- FAIL: .../limit_500   link limit = 2500, want 2000
--- FAIL: .../limit_5000  link limit = 25000, want 2000
```

**R-6a rule 1 — declaration `UNAVAILABLE` no longer suppressed by a `FOR`.** RED, and the output is
literally the contradiction the rule forbids — a `FOR` and an "anchor not supplied" on one anchor:

```quoted-output
--- FAIL: TestProviderUnavailableReasonPrecedenceKeepsObservedEvidence (0.00s)
    generation_service_test.go:201: reasons=[{Anchor:"ean", Direction:"FOR", Detail:"ean observado"},
    {Anchor:"ean", Direction:"UNAVAILABLE", Detail:"provider não fornece a âncora ean"}],
    want observed evidence [{Anchor:"ean", Direction:"FOR", Detail:"ean observado"}]
```

**R-6a rule 3 — precedence inverted so the seed beats the declaration.** RED:

```quoted-output
--- FAIL: TestProviderUnavailableReasonPrecedenceKeepsDeclarationOverSeedUnavailable (0.00s)
    generation_service_test.go:218: reasons=[{Anchor:"ean", Direction:"UNAVAILABLE",
    Detail:"ean sem correspondência"}], want declaration reason [{Anchor:"ean",
    Direction:"UNAVAILABLE", Detail:"provider não fornece a âncora ean"}]
```

**Ordering — declarations ranged through a map.** RED on **3 of 10** runs, `refforn` before `marca`:

```quoted-output
--- FAIL: TestProviderUnavailableReasonOrderingIsStable (0.00s)
    generation_service_test.go:295: reasons=[… {Anchor:"refforn", …}, {Anchor:"marca", …}],
    want stable order [… {Anchor:"marca", …}, {Anchor:"refforn", …}]
    (3 of 10 runs RED; 7 PASSED under the broken ordering)
```

**Why that last number matters, and why it is stated instead of rounded up to "RED".** Seven of ten
runs passed with output order driven by Go map iteration. A `-count=5` lane — the convention used
earlier on this chip — could have reported this guard GREEN while the ordering of PERSISTED reason
rows was genuinely non-deterministic. The proof is probabilistic, so it is labelled probabilistic.
Carried into the CLOSED payload as an upstream candidate: guards pinning the order of persisted data
run at `-count=10` minimum, and a pass at a lower count is not evidence.
