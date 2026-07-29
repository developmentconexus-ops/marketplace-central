# Veredito GPT-5.6 Sol effort=xhigh — análise global da harness (VERBATIM)

`proveniência: codex exec --model gpt-5.6-sol --sandbox read-only -c model_reasoning_effort=xhigh · 2026-07-29 · insumos: .mnfs/HARNESS-DEBTS.md + docs/HARNESS-PROFILE.md + Documents/mnfs-harness · prompt: scratchpad sol-harness-global-prompt.md · stream bruto 310128 bytes, bloco final extraído linhas 3389-3560`

## 1. Failure-pattern taxonomy (the core deliverable)

Ranked by observed aggregate rework, review rounds, latency, and risk of false acceptance.

1. **Observable-blind gates — the review system judges representations of reality without being provisioned to observe reality.**

   - **Evidence:** Seven dual-gate rounds reasoned about a generator and preserved a false UI-reachability premise; one live drive disproved it immediately (Profile amendment 2026-07-28, “PRODUCIBLE… vs REACHABLE…”). CHIP-ANCHORS-3 and CHIP-IMPORT-CHAIN reviewers could not run Go, PostgreSQL, or the requested lanes, leaving execution certified only by the implementer (HARNESS-DEBTS C-4; Profile amendment 2026-07-28, “dual gate has a THIRD seat”). CHIP-ANCHORS-3 and CHIP-VINC-NEUTRO then spent eight and nine rounds respectively, much of it on wording rather than shipped behavior (Profile amendment 2026-07-28, “R-24 IS SCOPED…”).
   - **Root cause at the ENGINEERING level:** Role independence was designed as context/model independence, but not as information-source independence. Two reviewers reading the same diff, generator, or pack share the same blind spot regardless of model diversity. Execution capability was removed from reviewers without first assigning each executable criterion to another accountable instrument.
   - **Structural remedy candidate:** Compile validation contracts into a typed verification graph: every criterion declares `static`, `executable`, or `live`, its required capabilities, its observable, and its independent executor. A capability scheduler produces immutable raw receipts at a named SHA; reading seats receive only the code diff, contract, and those receipts. This subsumes the Profile’s manual “third seat,” “diff not pack,” live-evidence carve-out, browser-drive custody, and producible/reachable prose, and retires C-4 plus the gate half of C-2.

2. **Non-discriminating green — lanes report success without proving that the intended code or test population executed.**

   - **Evidence:** Missing DB environment yielded RUN 27 / PASS 1 / SKIP 26 while the package printed `ok`; a mutation still printed `ok`; and a mojibake sweep returned zero because its regex matched nothing (HARNESS-DEBTS B-1). A worker reported 130 passes where the output contained 190 (B-2). Other greens came from a wrong test name, omitted integration tags, all-skipped integration tests, and output filters blinded by ANSI escapes (Profile amendments 2026-07-28, “VACUOUS GREEN” and “fifth instance”). A new migration produced 13 SQLSTATE 42703 failures because the base was never migrated (B-4; Profile amendment 2026-07-29, “A DB LANE ONLY PROVES…”).
   - **Root cause at the ENGINEERING level:** The harness equates process exit status and human summaries with semantic coverage. Commands have no declared expected population, mandatory preconditions, or machine-readable execution counts. The same actor runs, interprets, and summarizes the lane.
   - **Structural remedy candidate:** Introduce a standard lane runner that emits a signed JSON receipt containing checkout, SHA, cwd, command identity, prerequisite status, discovered/executed/pass/fail/skip counts, named expected tests, migration head, and artifact hashes. A lane cannot be green if the expected population is absent or unexpectedly skipped; parsers and filters require known-answer and negative-control fixtures. This retires B-1, B-2, B-4 and collapses the Profile’s vacuous-green, counting, ANSI-filter, DB-env, migration, and cwd-specific evidence rules.

3. **Unversioned prose as control plane — briefs, rulings, snapshots, and ownership exist as manually synchronized text.**

   - **Evidence:** M-06 carried three false assertions; S8 received a broken lane command; S10 inherited an already-refuted A-17 premise (HARNESS-DEBTS A-2). Six rulings, A-17 through A-22, required manual plan/pack amendments before dispatch (A-3). S10 lacked an accepted-round byte snapshot, making a narrow follow-up impossible to prove and forcing weaker re-review (A-5). The hub cut a patch without notifying the chip, while cross-session rulings crossed in delivery queues (Profile amendment 2026-07-28, “THE HUB ANNOUNCES THE CUT”; D-1). A data-dependent acceptance predicate moved from 3,822 to 2,923 after three measurable ERP facts were discovered only during execution (Profile amendment 2026-07-29, “A DESIGN DECISION…”).
   - **Root cause at the ENGINEERING level:** Markdown and messages are being used as a database without keys, versions, referential integrity, invalidation, or transactions. The system relies on orchestrator memory to propagate a ruling to every dependent brief and to preserve the exact parent state of each correction round.
   - **Structural remedy candidate:** Replace mutable cards and free-form rulings with immutable dispatch manifests and an event-sourced state store. Each dispatch references fact receipts, parent round/SHA, accepted snapshot, write-set, commands, model/effort, ownership version, and correlation IDs. A ruling transaction automatically invalidates dependent manifests; ACK is emitted only after the new state version is durable. This retires A-2 through A-5, D-1, D-2, and the Profile’s manual CUT, snapshot, ruling-amendment, and owner-registry ceremonies.

4. **Non-hermetic execution — command meaning varies with cwd, partial installs, caches, migrations, CLI versions, and checkout topology.**

   - **Evidence:** The Profile requires manual defenses for relative GOCACHE, half-installed `node_modules`, incorrect vitest cwd, CRLF-only gofmt noise, codex 0.144.4/0.145.0 cache incompatibility, and cold build timeouts (Profile §3; amendments 2026-07-17, 2026-07-18, and 2026-07-28). The Profile itself printed the forbidden relative GOCACHE form three times, producing an 83-byte failure with zero tests executed (amendment 2026-07-28, GOCACHE correction). Migration state remained external to source state and produced 13 schema errors (B-4).
   - **Root cause at the ENGINEERING level:** A worktree is treated as a reproducible build unit even though essential state lives outside it. Environment construction is an advisory sequence, not an atomic prerequisite of lane execution. Incident signatures are documented after each failure instead of eliminating the state space that caused them.
   - **Structural remedy candidate:** Provide hermetic lane environments with pinned toolchains, fixed working directories, completed dependency installation, isolated caches, and automatic schema convergence before execution. The runner creates or verifies that environment atomically and refuses partially initialized states. This subsumes most of Profile §3, the repeated GOCACHE/npm/vitest/CLI/migration instructions, B-4, and the worktree-provisioning portion of D-5.

5. **Capability-role mismatch — work is dispatched to seats that are physically unable to satisfy their contracts.**

   - **Evidence:** The workspace-write sandbox prevents FE workers from traversing to required config, so vitest aborts before collection even though Go lanes run in the same sandbox (B-3). Read-only gate seats were asked to execute Go and PostgreSQL criteria they could not run (C-4; Profile amendment 2026-07-28, “third seat”). Browser screenshot tooling was broken, and the resulting workaround has the hub temporarily check chip files into its own checkout to drive the UI (Profile §3 screenshot signature; §6 browser-drive seam).
   - **Root cause at the ENGINEERING level:** Role contracts describe responsibilities but not required capabilities. Scheduling considers model and organizational independence, while filesystem reach, temp-space access, database access, browser access, and artifact-write guarantees are discovered only after dispatch.
   - **Structural remedy candidate:** Publish machine-verifiable capability descriptors for every worker and seat, probe them before dispatch, and schedule only compatible criteria. Provision writers with repo-root read plus scoped write, and provide executors with controlled temp, DB, and browser facilities. This retires B-3, C-4, the sandbox-blindness mitigation, the manual third-seat routing, and most of the browser-file-borrow workaround.

6. **Doctrine/runtime divergence — rules are ratified faster than the executable harness can be released and activated.**

   - **Evidence:** The source/runtime harness is explicitly drifted and synchronization is operator-gated (D-4). The stop hook issued a third false accusation after the Profile had already ratified path+SHA attribution and `unknown` degradation; the runtime behavior did not change (D-7; Profile amendment 2026-07-28, “AN AUTOMATED GATE MUST NAME THE TREE…”). The installed 0.4.0 `dispatch-lint.sh` and `stop-gate.sh` differ from the harness source, and the runtime stop hook still matches any `CLOSED` substring and resolves evidence against session cwd. Stale skills also survive in caches and require a dispatch-level denylist (Profile §10; amendment 2026-07-16).
   - **Root cause at the ENGINEERING level:** Source, marketplace copy, plugin cache, repo-vendored scripts, and active runtime are separate release surfaces without an atomic version handshake. Ratification creates normative truth, but no deployment transaction proves that the enforcing code now implements it.
   - **Structural remedy candidate:** Ship hooks, skills, schemas, and scripts as one signed release bundle. Hub boot attests active hashes/version against the mission’s required version and blocks execution on mismatch. Activation is atomic and gated by known-answer hook tests, with automatic rollback. This retires D-3, D-4, D-7, most of Profile §10, and runtime-specific defensive prose.

7. **Instance-scoped correction loops — the system repeatedly repairs loci without retaining a machine-actionable defect class.**

   - **Evidence:** Five instances of duplicated HTTP error behavior were found before the harness classified the shared mechanism (C-1). CHIP-M05 spent six rounds point-fixing the same claim-from-the-wrong-layer mechanism across a comment, scope prose, and OpenAPI text (Profile amendment 2026-07-25). A delta brief named one repaired site and caused one reviewer to miss the same false universal 83 lines below (amendment 2026-07-28, “delta brief…”). An “exhaustive” regex sweep counted a population of 23 but silently excluded four accented/capitalized members (amendment 2026-07-28, “A SWEEP IS ONLY…”).
   - **Root cause at the ENGINEERING level:** Findings are keyed to round and location rather than mechanism. Delta review minimizes immediate token cost, so neither the correction state machine nor reviewer incentives require class-wide closure. Search populations and defect-class identities are free-form prose.
   - **Structural remedy candidate:** Give every finding a stable defect-class ID, occurrence counter, semantic search population, and detector. Repetition triggers a class-level correction unit and detector promotion, while a round budget forces redesign rather than another patch. This subsumes C-1 and the Profile’s third-round rule, delta-sweep instructions, regex reconciliation tutorial, and 2026-07-29 stop-the-line prose.

8. **Self-referential verification — producer and oracle share the same assumptions, representation, or positional structure.**

   - **Evidence:** Marshaling and unmarshaling through the same struct made JSON-tag transposition invisible; symmetric fixtures also failed to distinguish the fields (C-5). Positional test windows silently expanded as code was appended, with one guard reading six 500 responses where it intended four (C-3). Must-fail mutations could be red for multiple causes and could never detect an over-strict guard without a corresponding must-pass case (Profile amendment 2026-07-28, “A MUST-FAIL ARM…”). S6 verified that a value arrived but not that it came from the correct producer, accepting hardcoded `IncludeAll` (C-2).
   - **Root cause at the ENGINEERING level:** Verification maps describe commands but not oracle provenance. Tests reuse the production representation or author’s model, so apparent independence exists at the process level while semantic independence is absent.
   - **Structural remedy candidate:** Extend verification schemas with producer, consumer, oracle source, and independence classification. Transport tests compare raw wire forms against external golden fixtures; wiring tests trace producer-to-consumer provenance; guards require isolated must-fail and producible must-pass cases. This retires C-2, C-3, C-5 and the Profile’s positional-window, symmetric-fixture, and mutation-theory prose.

9. **Static resource scheduling — model effort, quota, and tool health are policy constants rather than admission-controlled resources.**

   - **Evidence:** A slice classified complex was still dispatched with Sol/low and paid a corrective review round; actual worker model/effort was absent from the ledger (A-1, A-4). Codex quota is an invisible wall detected only after dispatch (D-6). A temporary Claude-only contingency had to be written into the Profile after quota exhaustion (Profile amendment 2026-07-18).
   - **Root cause at the ENGINEERING level:** The model matrix is static and dispatch receipts do not close the loop between predicted complexity, actual model/effort, observed failures, quota, and latency. Availability is discovered by failure rather than admission control.
   - **Structural remedy candidate:** Add a dispatch admission controller that probes quota/tool health, validates actual model and effort, records the receipt, and selects among predeclared fallbacks using measured slice risk. This retires A-1, A-4, D-6, the quota contingency block, and repeated model/effort prose in core and Profile.

## 2. The Profile as symptom

The Profile is no longer primarily a repository binding. It is 1,032 lines with 67 dated amendment rows in roughly two weeks; 30 amendments landed on 2026-07-28 alone (Profile amendment log). That concentration is not evidence of unusually rich repository doctrine. It is evidence that one field session converted every newly discovered instrument limitation into binding prose.

The healthy repo-specific core is comparatively small:

- Identity, stack, shell, and default branch (§1).
- Exact lane entrypoints and stack commands (§2).
- DB isolation and lifecycle (§4).
- Concrete collision axes and owned seams (§§5–6).
- Product integrity invariants, truth order, and human gates (§§7–9).

Those match the intended Profile contract in HARNESS-CORE §0 and `PROFILE-TEMPLATE.md`.

The workaround families are much larger:

- Most of §3 is an incident runbook for non-hermetic tooling: cache paths, partial installs, cwd-sensitive vitest, CLI incompatibility, CRLF, sandbox behavior, and rasterizer failure.
- §10 compensates for plugin/cache distribution failure by teaching every dispatch which stale skills to ignore.
- Most of §11 is generic review and testing method: executing seats, class sweeps, anti-vacuity, mutation semantics, evidence custody, blocking/report classification, and reviewer inputs.
- Much of §12 is generic dispatch transport, freeze protocol, model routing, and hook attribution.
- The browser-drive seam in §6 is a controlled but dangerous workaround for the lack of an isolated live environment capable of running a chip build (Profile amendment 2026-07-28, “BROWSER-DRIVE SEAM…”).

The boundary is therefore unhealthy. HARNESS-CORE §0 says a method-level finding belongs upstream, yet the Profile holds generic rules about evidence discrimination, review topology, event custody, encoding, hook attribution, and mutation testing. There is also schema drift: `PROFILE-TEMPLATE.md` defines §11 as enforcement bindings, while the actual Profile uses §11 for code-review method; HARNESS-CORE §9 still assumes Profile §11 carries the required settings block. A supposedly binding schema has no validator capable of detecting that the section contract itself diverged.

Several large rule families would disappear if their instruments were fixed:

- A hermetic lane runner collapses the GOCACHE, npm, vitest-cwd, CLI-version, migration, CRLF, and partial-install clauses.
- Capability-aware scheduling collapses sandbox blindness, read-only execution findings, the manual third seat, and much of browser-drive custody.
- A typed event/artifact store collapses CUT announcements, round snapshots, verdict-paste custody, owner registries, correlation conventions, and manual ruling propagation.
- Atomic runtime attestation collapses stale-skill denylists, hook-drift warnings, and runtime-specific false-alarm signatures.
- Structured lane receipts collapse counting instructions, skip interpretation, output-filter recipes, and many sweep controls.

Consolidation should leave three artifacts:

1. A method package containing schemas, state machines, capability contracts, and executable validators.
2. A short repo Profile containing configuration values and genuine product invariants.
3. A non-binding, expiry-managed incident catalogue for tool signatures that cannot yet be eliminated.

Each existing Profile paragraph should be classified as `repo binding`, `method defect`, or `incident signature`. Method defects migrate upstream only with executable enforcement; incident signatures receive an owner and expiry; duplicated prose is then deleted. The current “living document grows from every field finding” rule in HARNESS-CORE §0 should be narrowed, because it currently rewards recording environmental entropy instead of removing it.

## 3. Process economics

The most disproportionate burn came from reviewing claims rather than resolving observables:

- CHIP-ANCHORS-3 reached eight rounds and CHIP-VINC-NEUTRO nine; CHIP-M05 reached six. Seven dual-gate rounds failed to resolve a reachability question that one live drive answered immediately (Profile amendments 2026-07-25 and 2026-07-28).
- Evidence packs reached 9,730 / 11,502 / 21,788 lines against 605 / 1,811 / 1,318 code lines: approximately 16.1×, 6.4×, and 16.5×. The “10–20×” framing is directionally correct but not universally accurate (Profile amendment 2026-07-28, “A finding BLOCKS…”).
- Cold remedy reviews cost 80–121k tokens before warm delta reuse was introduced (REVIEW-STANDARD §9). Agent-tool gates ran 281–679 seconds without visibility (HARNESS-CORE §1). A one-line scribe continuation consumed roughly 64k tokens, and monolithic chips carried about 130k context tokens per call while making 430–640 tool calls and injecting 0.9–1.4M characters (HARNESS-CORE §§1 and 8).
- Coordination overhead was also material: six manual ruling amendments in a week (A-3), four restart round-trips in M-01 (Profile §6), four held slices plus four hub round-trips from predictable seam gaps (HARNESS-CORE §4), and repeated cross-session/ownership waits (D-1, D-2).

The main economic feedback loops are:

1. **Untrusted lane → more evidence prose → reviewers inspect prose → findings about prose → larger evidence packs.** The Profile records the terminal state directly: gates were reviewing their own paperwork while comment-only findings opened new rounds (amendment 2026-07-28, “A finding BLOCKS…”).

2. **Restricted reviewer → implementer self-certifies execution → add another reviewer → all reviewers retain the same missing observable.** The system bought model diversity while withholding DB/browser/execution access; seven reachability rounds demonstrate that correlated inputs dominate reviewer count (C-4; Profile “PRODUCIBLE vs REACHABLE”).

3. **Mutable brief → stale claim → corrective ruling → manual propagation → another stale descendant.** A-17 propagated into S10, six amendments depended on hub memory, and the ERP predicate changed after dispatch because measurable facts were not represented as prerequisites (A-2, A-3; Profile 2026-07-29).

4. **False hook alarm → prose amendment → runtime unchanged → another false alarm.** D-7 is the third occurrence after the Profile had already specified the correct behavior. The marginal amendment cost bought no enforcement value.

5. **Delta review saves tokens locally but hides classes globally.** A site-specific delta brief caused one reviewer to stop exactly where instructed, while a second occurrence remained 83 lines away; the resulting additional round erased the intended saving (Profile amendment 2026-07-28, “delta brief…”).

6. **Static model economics externalize dispatch mistakes into review.** Missing model/effort receipts and late quota detection make the corrective gate pay for decisions the dispatcher should have rejected before launch (A-1, A-4, D-6).

## 4. Refutations

1. **A-1 does not prove that Sol/low caused the defects.** It establishes correlation: one complex slice ran at low effort and later had a renamed-not-removed function and a magic value. The core itself deliberately maps complex work to Sol/low (HARNESS-CORE §1). Without matched slices or historical failure rates, “complex must medium/high” is an unsupported remedy. The valid finding is that dispatch policy is neither instrumented nor enforced (A-1, A-4).

2. **C-4 is stale as written.** It says there is still no structural fix, but the Profile already ratified a hub-owned independent executing seat and removed execution criteria from reading-seat prompts (Profile amendment 2026-07-28, “dual gate has a THIRD seat”). That fix remains manual and incomplete, but it is plainly structural relative to implementer self-certification.

3. **D-3’s candidate fix is insufficient because the hook already attempts that distinction.** Both source and installed `dispatch-lint.sh` say they apply only to chip dispatches, but implement the distinction by searching arbitrary prompt text for `CHIP-`. A backlog-seeding prompt can mention that token and be misclassified. The missing mechanism is a typed operation kind supplied outside prompt prose, not a better substring convention (D-3; harness hook source).

4. **B-4 should not move to core in its concrete form.** “Run `cmd/migrate`” is a repository binding. Core should require a schema-convergence prerequisite and a receipt; the Profile should supply the exact command. Moving the command itself upstream would reproduce the boundary violation already visible throughout §11 (B-4; HARNESS-CORE §5’s semantics-versus-command split).

5. **C-3 and C-5 are not primarily missing doctrine.** The Profile demonstrates that prose rules can be ratified while the runtime remains unchanged (D-7), and the same faculty can reproduce its assumption in fixture, sweep, and proof of sweep (Profile amendment 2026-07-28, “A SWEEP IS ONLY…”). The remedy must encode oracle independence and population boundaries in test tooling or schemas, not add another transport-testing essay to core.

6. **A universal “stop on the second occurrence” threshold is not justified.** The Profile already contains a different threshold: third defect of the same shape or third correction round (amendment 2026-07-25), followed by a second-occurrence stop rule on 2026-07-29. Neither rule includes severity, expected loss, detector cost, or class prevalence. A class registry and risk budget are defensible; a global numeric trigger is another point rule.

7. **The prompt’s “24 field-collected debts” count is misleading for the current file.** The current A–D list contains 21 debt IDs, including the late-added D-7. The remaining three bullets are under §E and explicitly labeled product legacy, “context, not harness debt.” Counting those yields 24 items, not 24 harness debts.

8. **Not every Profile amendment is a post-failure point-fix.** The 67 rows include initial extraction, upstream method adoptions, planned execution-model changes, and a temporary quota contingency as well as incident amendments (Profile log entries 2026-07-15 through 2026-07-18). The accretion critique remains valid, especially for the 30 rows on 2026-07-28, but treating every row as identical would distort the evidence.

9. **“Evidence packs are 10–20× the code” is not exact across all cited packs.** The measured ratios are approximately 16.1×, 6.4×, and 16.5×. The economic pathology is real, but one pack falls materially below the stated range (Profile amendment 2026-07-28, “A finding BLOCKS…”).

## 5. Top 5 interventions

1. **Replace P6 with an observable-first gate compiler and class-aware state machine.**

   Validation criteria become typed executable objects with required capabilities, observable, executor, reachability source, and defect-class ID. Reviewers receive current-tip code diff plus raw receipts, never full evidence packs. Repeated classes trigger detector-backed class closure, while prose-only findings cannot open rounds. This attacks the largest measured burn: eight/nine-round gates and seven rounds defeated by one live drive (Profile amendments 2026-07-28).

   **Expected kill-list:** C-1, C-2, C-4; Profile “third seat,” “diff not pack,” BLOCKING/REPORT, producible/reachable, live-evidence carve-out, delta class-sweep, third-round, and stop-the-line rules.

2. **Build one standard lane/evidence runner with anti-vacuity and oracle-independence contracts.**

   It must control cwd and prerequisites, enumerate expected tests, emit structured counts, distinguish skips, record migration state, run known-answer controls, and require independent wire/golden oracles where producer and consumer representations overlap. This directly prevents false green rather than teaching agents how to interpret it (B-1, B-2, C-3, C-5; Profile 2026-07-28/29).

   **Expected kill-list:** B-1, B-2, B-4, C-3, C-5; Profile vacuous-green, sweep-count, ANSI-filter, DB-env, migration, positional-guard, symmetric-fixture, and must-fail/must-pass prose.

3. **Introduce a typed, versioned dispatch and event control plane.**

   Dispatch manifests should carry verified fact receipts, model/effort, ownership snapshot, base/current/parent SHAs, round identity, commands, write-set, and correlation IDs. Rulings become transactions that invalidate dependent manifests; ACK, CUT, and completion artifacts are atomic state transitions. This removes orchestrator memory from correctness (A-2, A-3, A-5, D-1).

   **Expected kill-list:** A-2 through A-5, D-1, D-2; manual pack amendments, CUT announcements, accepted-byte snapshots, verdict-paste custody, removal-owner lookup, and substantial dispatch-ledger prose.

4. **Provision hermetic, capability-declared execution environments and perform admission control before dispatch.**

   Pin toolchains and cwd, complete dependency/bootstrap operations atomically, expose repo-root reads with scoped writes, provision controlled temp/DB/browser capabilities, and probe quota/model/tool health. Dispatch fails before tokens are spent when the requested role cannot execute its verification plan (A-1, B-3, C-4, D-6).

   **Expected kill-list:** A-1, B-3, C-4, D-5 worktree provisioning, D-6; most Profile §3 false-alarm signatures, GOCACHE/npm/vitest/CLI/rasterizer clauses, quota contingency, browser-file borrowing, and manual model-effort routing.

5. **Make the harness an atomically released, attested product with conformance tests.**

   Hooks, skills, scripts, schemas, and role tables ship as one signed version. Hub boot verifies active runtime hashes, refuses source/cache drift, and runs known-answer tests for hook scope, checkout attribution, false-positive degradation, encoding, and transport. No doctrine amendment is considered enforced until the matching release is active (D-4, D-7).

   **Expected kill-list:** D-3, D-4, D-7 and remaining D-5 distribution issues; Profile §10 stale-protocol denylist, hook-path/SHA workaround rules, runtime-drift warnings, manual `harness-sync` dependence, and repeated transport/encoding caveats.
