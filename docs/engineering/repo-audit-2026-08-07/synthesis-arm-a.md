# Phase 2 — Synthesis, Arm A

Method: `MetalDocs/docs/engineering/repo-audit-playbook.md`, Phase 2. Binding brief:
`SYNTHESIS-BRIEF.md`. Inputs: `PHASE-0.md` (14 established facts) plus the ten lane reports in
`lanes/`. Repo state read at `11b9e494` (one commit past the `1473e863` PHASE-0 cites).

Arm A had no visibility into Arm B and makes no prediction about it.

**Two amendments absorbed; affected sections rewritten, not annotated.**

*Amendment 1.* The zero-spend constraint is dead: GitHub Actions is available, billing is resolved,
and the failed run `29712873135` is history. Firing level 3 became reachable for anything expressible
as a command that exits non-zero, so parking such a rule at level 5 is now a defect in the
recommendation.

*Amendment 2, and it corrected a wrong fact.* The remote is **not** stale or unresolvable. That was a
keyring artifact: `gh` asks as `leandrotcawork` while git authenticates as `developmentconexus-ops`.
Measured: `origin/main` = `7df7d011` (2026-08-04); local `main` is **98 ahead, 0 behind,
fast-forwardable, nothing diverged**. The `legacy` remote (`leandrotcawork/marketplace-central`,
default branch `master`, last push 2026-07-20) is superseded. Additionally: **the repository will be
made public on a free personal account**, which puts rulesets and enforced required status checks in
reach at zero cost and kills the credential-custody fallback branch entirely; the flip is gated on a
separate secrets/PII sweep over the *history*, not the tree; and **the 98 commits stay local until
the gates exist** — the first remote execution must already be gated.

*Amendment 3, and it re-ranked the sequence.* The reachability question every lane left open is
measured: **the unauthenticated PII surface is network-reachable by design, via two independent
paths.** `deploy/Caddyfile` terminates TLS on `{$MPC_DOMAIN}` and reverse-proxies `/orders` and
`/orders/*` — for any request not asking for HTML — straight to `backend:8080`, reaching
`orders/transport/http_handler.go:608-618` with no identity check anywhere in the chain
(`composition/root.go:994` is `CORSMiddleware(apierror.Recover(mux))`). In development, the
`oauth` compose profile runs `ngrok http --url="$callback_host" frontend:5174` because Mercado
Livre's OAuth callback requires a public URL, and `apps/web/vite.config.ts` proxies `/orders`,
`/pricing`, `/marketplaces/*`, `/profitability` and more to `backend:8080` — so while that profile is
up, the same surface is on the public internet. Whether a production host is live right now is
unverified (`MPC_DOMAIN` is external), but the development path is one `--profile oauth` away and
exists because the integration demands it. Amendment 3 also supplied a second finding — three
hand-maintained route tables that must agree — which lands in A4.

Net effect on this synthesis: A1's *cause* changed and its content shrank; the gate topology got
stronger (enforced checks, not custody) and gained a public-repo threat surface it did not have; the
sequence inverted at the top — gates land before the push, not after; and **A8 moved from fourth to
first, because the mitigation I had written for it was premised on a localhost assumption that is now
measured false.** See §1 A4 and A8, §4, §5, and §7 D8.

Everything below judges. The lanes measured and were forbidden to weigh.

---

## Reading conventions

- `lanes/<file>.md:N` cites a lane's own line. `path/to/file.go:N` cites the repo.
- **[ARM A]** marks a claim no lane made, always with the measurement that would confirm or kill it.
- Where two lanes disagree I name the disagreement rather than splitting it.
- "A control that exists and does not fire is absent" is applied literally throughout. Where a
  control is inert, the axis is sized as if it were not there, because it is not.

---

# 1. Axes, not a list

Nine axes. Each named for its cause. Every finding sits in exactly one. Where a finding looked like
it belonged to two, that was a signal the cause was misnamed, and the axis was renamed until it
stopped being true.

---

## A1 — A working delivery path exists and nothing has ever been required to pass through it

**Cause.** Under the second amendment this axis got smaller and more embarrassing. The remote is
real, reachable, and 0 behind: `origin/main` = `7df7d011`, local `main` 98 ahead, fast-forwardable,
no divergence. The route is open. **What is missing is not the path but the requirement** — no change
has ever been observed by a verifier the author does not control, because no pull request has ever
existed and nothing on the remote is configured to demand one.

Every mechanism the operator wants — issues, PRs, PR review, CodeRabbit, required checks — attaches
to a pull request. Work is authored, committed, and merged inside one machine, by the same actor that
would be doing the checking. The seam exists in the tooling and has never been used.

**Findings.**

- Zero pull requests, ever. `lanes/cicd.md` F1: the only workflow is `release-images.yml`, triggering
  on `push: branches: [main]`, `tags: ["v*"]`, and `workflow_dispatch`. Nothing triggers on
  `pull_request`.
- 98 commits sit local. The count has been read three times during this audit — 96 (PHASE-0), 97
  (cicd lane at `11b9e494`), 98 (amendment 2). **Do not reconcile these; they are monotonic growth
  from the audit's own commits.** The number is not a fact about the backlog, it is a fact about how
  long the backlog has been growing.
- **`lanes/cicd.md` F2 — "no branch protection reachable" — was a measurement artifact and its
  conclusion is void.** Two GitHub identities in the keyring: `gh` asks as `leandrotcawork`, git
  authenticates as `developmentconexus-ops`. The lane's instrument was authenticated as the wrong
  principal. See §7 D8; this is not a small correction, because F2 was the single fact that made the
  whole program look blocked.
- The one recorded remote run (`gh run view 29712873135`, 2026-07-20) died on billing before any
  build step (`lanes/cicd.md` F3). History now, with one residual value: **the only workflow in this
  repo has never successfully executed even once**, so no claim about its behaviour is empirically
  grounded.
- `scripts/release.sh:1-9` declares in its own header that it bypasses CI (`lanes/cicd.md` F14).
  Currently honest, because there is nothing to bypass. It becomes a live problem the moment A2
  lands, and the fix is to change the header's premise, not the script.

**[ARM A] The hazard in pushing, and it survived the amendment intact.** `release-images.yml`
triggers on `push: branches: [main]` and pushes the `latest` tag to GHCR. **A fast-forward of 98
commits to `main` publishes production images from a tree whose `arch-gate` fails from at least five
independent causes, with no verification in front of it.** No lane connected the trigger to the
backlog. Confirmed by reading the workflow's `on:` block directly. This is now the primary technical
reason — beyond the operator's stated preference — that the gates must land *before* the push.

**Why these fix together, and what the operator's sequencing decision means.** Amendment 2 is
explicit: the gate topology lands before the push, so the first remote execution is already gated
rather than a red run against 98 commits of history this audit has already catalogued. That forces a
specific and slightly unusual shape:

1. `verify.yml`, the ratchet baselines, and the `release-images.yml` guard are developed **locally**,
   inside the existing backlog, as part of A2/A3.
2. They are driven green **locally** against the current tree before anything is pushed.
3. The first push is a fast-forward of 98 + N commits to a `main` that already contains the workflow
   and already passes it.
4. The ruleset (required checks, no bypass) is configured immediately after, and **from that moment
   every change goes through a PR.**

**The 98-commit backlog is therefore grandfathered by decision, not by verification.** That is a
deliberate, defensible trade — this audit is itself the record of what is being grandfathered — but
it must be stated plainly rather than allowed to look like the backlog passed something. The three
ratchet baselines (44 archscan, 58 governance, 234 ADR-023) plus the 22 genuinely unformatted files
**are** the written record of the grandfathered debt, which is exactly why they must be committed in
the same change as the gates.

**Dependency on the secrets/PII sweep.** The public flip is gated on a sweep of the **entire
history**, not the current tree. Nothing in this axis or any other should be executed against a
public remote until that verdict lands. Two of my recommendations depend on it directly: the public
flip itself, and §4.3's requirement that gate output be safe to publish (CI logs on a public repo are
world-readable). If the sweep finds anything in history, the flip decision — and with it the
availability of free required checks — reopens.

**Size:** 1 day of actual work, all of it at the end of A2/A3 rather than at the start. The credential
confusion is already diagnosed and costs minutes to resolve.

---

## A2 — The verifier is a set of scripts a person has to remember, not one product with one entry point

**Cause.** There is a large amount of real verification in this repo and no single thing that runs
it. Each check has its own invocation, its own prerequisites, and its own audience. What ties them
together is human memory, and human memory is level 5, which the brief correctly rules is not a
control.

**Findings.**

- **369 of 405 Go test files (91%) are reachable from no npm or harness command** (`lanes/testing.md`
  T-1: 353 `internal/**` files, ~1800 `func Test`; T-8: `migrations/*_test.go`, 15 files / 26 funcs,
  and `cmd/catalogingest/main_test.go`). `npm run harness:unit` runs `go test ./tests/unit/...` plus
  FE vitest — never `./internal/...` (fact 3, D-51).
- `grep -n "arch-gate" package.json` → 0 hits (fact 3). `arch-gate.sh` is the only thing that runs
  `./internal/...`, and it is wired to nothing.
- **`tsc --noEmit` is wired to no command. Run directly it fails with 12 errors, 3 in production
  code** (`lanes/frontend.md` F2). Frontend F3 is that fact's live proof: two shipped components
  render `<ErrorState detail="…" />` with the required `onRetry` omitted, and the type checker that
  would have refused it has never run in anger.
- 9 FE test files / 111 cases are structurally unreachable (`lanes/testing.md` T-3: `sdk-runtime`
  7/86, `feature-classifications` 1/15, `feature-inventory` 1/10). `apps/web/vitest.config.ts:14`
  name-pins `feature-products` by exact filename (T-10, B-7) — a glob would have caught the rest.
  Root `vitest.config.ts` is orphaned (T-12).
- 13 `scripts/tests/*.tests.ps1`, roughly 250 assertions, invoked by nothing: `Invoke-Pester` returns
  0 non-vendor hits (`lanes/cicd.md` F10, `lanes/testing.md` T-7). 6 of 11 are not green when run
  cold (`lanes/cicd.md` F11) — which is what always happens to a suite nobody runs.
- Governance's diff-scoped rule block is gated on `-BaseSha` with no default and no derivation
  (`lanes/cicd.md` F13, `Policy.psm1:447`). Omit the argument and an entire class of rules silently
  does not apply.
- No git hooks, no husky, no lint-staged (`lanes/cicd.md` F4); no `.claude/settings.json` hooks
  (F5).
- **[ARM A] `scripts/arch-gate.sh` cannot be wired as-is, and no lane said so.** Two structural
  defects, both read directly from the 57-line script. (a) Step 1 runs `gofmt -l "$SERVER/internal"`
  — **`cmd/`, `migrations/`, and `tests/` are gofmt-checked by nothing at all**, which is the same
  blind spot as T-1/T-8 in a different instrument. (b) Step 5 is `git status --porcelain
  --untracked-files=all` and fails on any dirty tree: *"a gate cannot certify a tree it did not
  see."* The intent is right and the placement is fatal — **it makes the script structurally
  incapable of being a pre-commit or pre-push gate on work in progress.** It can only certify an
  already-committed clean tree. This is why the natural next step ("just wire arch-gate into
  package.json") does not work and would produce a gate that fails on every real invocation.
  Confirmed by reading the file; no lane report contains the string `--untracked-files`.

**Why these fix together.** They are all the same missing artifact: a single `verify` command with a
declared manifest of everything that must pass, no undeclared prerequisites, identical locally and
on a runner. Fixing them one at a time means N separate wirings, each of which can rot
independently. Fixing them together means one manifest and one place where a check can be added.

**Sizing note.** Under control-versus-effect this axis is sized as if none of the 369 unreachable
Go test files, none of the 111 FE cases, and none of the 250 PowerShell assertions exist today —
because for the purpose of stopping a bad change, they do not. That is also the reason this axis is
cheap relative to its payoff: the tests are already written. The work is wiring, plus paying down
whatever the first honest run surfaces.

**Size:** 4–6 days. Wiring is 1–2; the rest is the first-run debt (`tsc`'s 12 errors, the 6 red
Pester scripts, whatever `go test ./internal/...` shows after 91% of it has been dark).

---

## A3 — A green result does not prove the check ran

**Cause.** Several controls here emit success indistinguishably from not having executed. This is a
distinct cause from A2: A2 is "nobody invokes it," A3 is "it is invoked and its output cannot be
read." A2's fix makes A3 *more* dangerous, because on a runner nobody is watching the terminal.

**Findings.**

- **`TestModuleBoundaryADR023` has one test function, no fixture, no positive control, and no
  must-fail proof** (`lanes/testing.md` T-2). It is the primary enforcement of ADR-023 and the only
  evidence it works is that it currently reports 234 violations. Contrast the correct shape in the
  same repo: `internal/arch/scan_test.go` and `archguard_test.go` both ship violations-and-clean
  fixtures.
- **The integration lane prints no RUN/PASS/SKIP/FAIL counts, so an all-skipped run is byte-identical
  to a green run** (`lanes/testing.md` T-6, D-26).
- 6 `t.Skip` sites behind `MPC_ORACLE_LIVE_TEST=1` that no lane and no command sets
  (`lanes/testing.md` T-11). Those tests report as passing.
- `router_registration_test.go:172-179` asserts only "not 404" (`lanes/testing.md` T-9). A route
  returning 500 satisfies it.
- **A negative-fixture test asserts `GOV_CONTEXT_UNREGISTERED`, an error code `Policy.psm1` never
  implements** (`lanes/cicd.md` F9). The fixture is green because the assertion is unreachable. This
  is the purest specimen in the audit: a test written to prove a guard exists, passing precisely
  because the guard does not.
- Two boundary invariants are proven only by tests authored in the same commit as the fix (catalog
  RLS at `7e3dcc47`; catalogingest guard at `47a76837`) (`lanes/testing.md` T-4). See §7 D4 — I
  reclassify T-4(a), and only (a).

**Why these fix together.** One rule closes all of them: **every guard ships, in the same change,
with an input that makes it fail, and every lane prints an attributable count.** That is a single
review standard and a single output convention, not six repairs.

**Size:** 3–4 days. Writing the missing negative fixtures is most of it. `TestModuleBoundaryADR023`'s
fixture is the largest single piece and it has a working template 20 lines away.

---

## A4 — The API surface is transcribed by hand into every consumer and nothing compares the copies

**Cause.** One fact — the shape of a resource, or the set of route prefixes the backend owns — is
written out by hand in every place that needs to know it, and no mechanism compares the copies. Drift
is not a risk here; it is the expected steady state, and the measurement confirms it already
happened, twice, in two different dimensions.

**Six hand-maintained transcriptions of one surface.** Four for the *shape*: Go domain type → DTO →
OpenAPI document → TypeScript SDK. Two more for the *routing*: `deploy/Caddyfile`'s matcher table and
`apps/web/vite.config.ts`'s proxy table, with the Go router as the third authority neither is
generated from.

**Findings.**

- SDK is 100% hand-written: `packages/sdk-runtime/src/index.ts`, 2595 lines, 172 interfaces (fact 5),
  **plus four more hand-written modules** — `activeSource.ts`, `dashboard.ts`, `erpImport.ts`,
  `market.ts`, 354 more lines and 2 additional client factories (`lanes/frontend.md` F13).
- `contracts/api/marketplace-central.openapi.yaml`: 8574 lines, 111 operations, 95 paths, 239
  schemas — all hand-maintained.
- **`GOV_API_SDK_SPLIT` requires only that both files change in the same commit; it never checks that
  they agree** (fact 6; `Policy.psm1:460`, `if ($apiChanged -xor $sdkChanged)`). Under
  control-versus-effect this is not a weak control, it is an absent one, and the axis is sized that
  way.
- The duplication lane found two live `ListingReadModel` divergences (`lanes/duplication.md` DUP-1).
  **[ARM A] There is a third in the same sample they took**: `read_model.go:122-152` declares
  `PublishedQuantity *int` with no `omitempty`, so the wire always carries the key and it may be
  `null`; `index.ts:380-405` declares `published_quantity: number` — non-optional, non-nullable. Any
  consumer that trusts the SDK type will do arithmetic on `null`. Confirm by diffing those two
  ranges field by field; I did, for that one struct only. **The point is not the third divergence.
  The point is that a hand sample of one struct yielded 3, which is the correct way to size this
  axis** — nobody has counted the other 171 interfaces and nothing can.
- The only thing connecting OpenAPI to the SDK is a comment at `index.ts:1747`
  (`lanes/frontend.md` F1).
- `oapi-codegen` and `openapi-typescript` were approved on 2026-08-07 and are not installed (fact 7).
  The decision is made; the work is not started.
- Delivery F14 belongs here, re-sized. See §7 D1: the claim of "10 operations with zero SDK method"
  is wrong; the real number is 3, and all 3 are admin surfaces. It stays in this axis as evidence
  that a manual reconciliation of the two documents is unreliable *even when performed carefully by
  a dedicated lane*, which is the strongest available argument that it must be mechanical.
- **Three route tables that must agree, with the requirement recorded only as a comment** (amendment
  3). `deploy/Caddyfile`'s own comment says it mirrors the dev proxy table in
  `apps/web/vite.config.ts` and asks a human to *"keep both in sync when a new backend route prefix
  appears"* — while the Go router is the actual authority and generates neither. **Named failure
  mode: a route added to the router and missed in one table works in one environment and 404s in the
  other, silently, and the environment it breaks in is whichever one the author did not run.** That
  is a strictly worse failure than the contract-shape drift above, because shape drift eventually
  produces a wrong value a user can see, whereas this produces a route that simply is not there in
  production while every local test passes.

**Why the route tables belong here and not in A2.** The tempting placement is the verification axis:
write a check that compares the three tables. That would be treating the symptom. The cause is
identical to the other five copies in this axis — **one fact transcribed by hand into every consumer,
with a code comment where a mechanism should be** — and the brief is explicit that if two axes want
the same finding, either they are one axis or the cause is named wrong. A comparison check is what
A4 falls back to *before* generation exists (§4.3), not what the finding is.

**Why these fix together.** Generation from a single authority collapses five of the six copies into
one authored artifact and one build step: `oapi-codegen` and `openapi-typescript` for the shape, and
the route-prefix table emitted from the router (or from the OpenAPI `paths`) into both the Caddyfile
and the Vite config. It is a single change of one kind; splitting it produces two half-migrated
generators and a third dialect.

**A4 also feeds A5 and A6.** Generated request/response types are the vocabulary that the boundary
kit (A5) binds to, and generated money fields are the forcing function for the money type (A6).
Doing A4 after either of them means doing that work twice.

**Size:** 4–6 days. Not the generation — that is a day. The cost is that generated types will not
match the hand-written ones at 172 sites and each mismatch is a real decision about which side was
right.

---

## A5 — Every request boundary is re-invented per module instead of provided once

**Cause.** There is no shared kit for the edge. Each module author writes their own parsing, their
own error mapping, their own envelope, their own status decision — in Go at the HTTP boundary and in
React at the query boundary. The failure mode is not ugliness; it is that a correct pattern
established once in one module has no mechanism to reach the next one.

**Findings — Go transport.**

- **12 independent `map*Error` functions dispatching on `strings.HasPrefix(err.Error(), …)`, and 0
  uses of `errors.Is`/`errors.As` anywhere in transport** (`lanes/delivery.md` F3). HTTP status is
  coupled to the prose of an error string, which no compiler protects and no rename catches. The
  codebase has 42 sentinel errors, 120 `errors.Is` and 37 `errors.As` calls (`lanes/goidiom.md`) —
  the idiom is known and used everywhere *except* the layer where it decides what the client sees.
- 51 method-prefixed routes fall through to net/http's plain-text 405; 47 bare-path routes use the
  JSON envelope (`lanes/delivery.md` F4, `server.go:2707`). Same API, two error content types, split
  by an implementation detail of registration.
- `parsePositiveInt` copied 5×, and **already drifted**: `product_links/transport/http_handler.go:514-521`
  uses `fmt.Sscanf` and accepts `"20xyz"` (`lanes/delivery.md` F5).
- 7 page-envelope structs (F6); 3 divergent response shapes including
  `market/collection_handler.go:64` with the JSON key `json:"decisões"` (F7); 12 distinct error codes
  for a malformed JSON body (F12); `MaxBytesReader` present in 2 of 20 transport packages (F10); 5
  SDK query builders plus inline concatenation (F15); 2 hand-copied error-envelope JSON literals
  (`httpx/json.go:15`, `httpx/route_deadline.go:129`) forced by an import cycle (F13).
- 19 module-local `write*Error` wrappers, 3 byte-identical (`lanes/duplication.md` DUP-7).

**Findings — React boundary.**

- **Zero React error boundaries in the entire app** (`lanes/frontend.md` F4).
- 33 of 41 `<ErrorState>` call sites discard the typed error code they were handed (F7).
- Three coexisting data-fetching patterns (F9). `ClassificationsPage.tsx` uses 7× `catch (err: any)`
  and no react-query (F8).
- The D-11 fix — the exhaustive `query.status` switch at `SyncHealthCard.tsx:127-190` — was applied
  to exactly one file; 9 remain on the vulnerable ternary chain (F10). **That is this axis in one
  finding: a correct pattern with no mechanism to propagate.**
- Four `apiBaseUrl` implementations, 3 of which bypass the token-refresh `fetchImpl` (F5).

**Why these fix together.** One kit, two languages: a Go `httpx` package that owns parsing, limits,
method handling, typed-error→status mapping, and the envelope; and a React query-state component
that owns loading/error/empty exhaustively. Each finding above disappears when its module adopts the
kit. Fixed individually they are 30 edits that re-diverge.

**Dependency.** The typed-error registry this needs is far cheaper once A4 has generated the error
schema, and the whole migration is only safe once A2's `verify` can prove 30 transport edits broke
nothing.

**Size:** 5–7 days. The kit is 2; adoption across 20 transport packages and 10 FE pages is the rest.

---

## A6 — There is no single compiled vocabulary for a value

**Cause.** Money, percentages, and enumerations each exist as several unrelated representations with
no type that makes the wrong one unrepresentable. This is the one axis where the domain itself
supplies the argument: Brazilian fiscal arithmetic is exact by regulation, and IEEE-754 binary
floating point cannot represent 0.10. That is admissible evidence independent of how this repo is
currently organised.

**Findings.**

- **`Money` is defined 5 times** — `connectors/domain/money.go:10`, `listings/domain/read_model.go:83`,
  `market/domain/market.go:13`, `pricing/domain/decimal.go:12`, and the canonical
  `kernel/exact/money.go:45`. Four are byte-identical (`lanes/duplication.md` DUP-2).
- **`kernel/exact` — the canonical one — has zero production callers** (DUP-3). The correct type
  exists and is inert.
- 12 `strconv.ParseFloat` money sites (DUP-4), including
  `orders/adapters/pricingtax/reader.go:309`, which round-trips the output of `FormatRatHalfUp(r, 2)`
  back through `float64`. Exactness is achieved and then discarded, in one expression.
- **The margin formula is implemented twice with different unit conventions**:
  `orders/domain/order_decomposition.go:118-163` uses `*float64` as a fraction;
  `pricing/domain/decompose.go:141-286` uses `big.Rat` ×100 (DUP-5). Two answers, off by 100×,
  neither labelled.
- 373 `float64`/`float32` occurrences across 78 files, concentrated in money fields in `orders`,
  `profitability`, and `marketplaces` (`lanes/goidiom.md` G-1). 15 files hand-construct `big.Rat`
  (DUP-9).
- Persistence mirrors the split: `pgtype.Float8` at 25 sites in 2 files in `orders`/`profitability`,
  versus the `::text` discipline at `pricing/matrix_reader.go:34-39,64-68` (`lanes/persistence.md`
  P-5). On the wire, money is string-exact in some responses and a bare `float64` in others
  (`lanes/delivery.md` F8).
- The money formatter is reimplemented 3× in the frontend, including once *inside* `packages/ui`
  itself at `ProductPicker.tsx:44` (`lanes/frontend.md` F6) — the shared package does not use its
  own shared function. `apps/web` computes a DIFAL sum client-side over only the loaded page (F14).
- **65 `CHECK (… IN (…))` constraints and 0 `CREATE TYPE`**, with at least 2 vocabularies duplicated
  verbatim across migrations (`lanes/persistence.md` P-6). Adding a value means finding every copy.

**The memory record is consistent with this and sharpens it.** `float64` in orders is currently safe
*by accident* — `unit_price numeric(14,2)` makes the reachable-precision-loss input unreachable from
the database. No line of code states that premise, nothing enforces it, and the in-memory ingest path
does not go through that column. A safety property that no artifact asserts is not a safety property.

**Why these fix together.** One exact money type with **no `float64` constructor and no `float64`
accessor**, adopted everywhere, plus generated DB enums. Every finding above is an instance. Done
piecemeal, each converted site needs a float bridge at its edges, and the bridges are exactly the
defect.

**Size:** 7–10 days, the largest axis. It is a mechanical change across 78 files whose only real
safety net is a test suite that today is 91% unreachable. **This axis must not start before A2
finishes.** That ordering is worth more than any amount of care.

---

## A7 — The module boundary has four instruments and none of them covers the tree the code is moving into

**Cause.** ADR-023 is enforced by four separate mechanisms that were each written against the tree as
it stood when they were written. The code is migrating from `internal/modules` (21 legacy modules)
into `internal/contexts` (fact 10), and each instrument is anchored to a literal path prefix, so
**enforcement goes silent exactly where the new code is being written.** A detector anchored to a
prefix fails precisely during a migration, which is the only time it matters.

**Findings.**

- `TestModuleBoundaryADR023` matches only the literal `…/internal/modules/` prefix
  (`module_boundary_arch_test.go:35`) and is blind to modules↔contexts edges
  (`lanes/layering.md` L-04).
- **`internal/contexts` has zero governance coverage**: `Policy.psm1:302` hardcodes `internal/modules`
  (`lanes/cicd.md` F8). And a folder with no registry entry yields *zero findings*, not a failure —
  the registry walks the registry, never the tree.
- **`internal/composition` is exempt from every detector and contains 43 vendor-token occurrences
  across 5 files** (`lanes/layering.md` L-03). `lanes/cicd.md` F7 measures archscan at 44 findings
  total, 42 of them vendor-token in composition. The DI root is both the largest violator and the
  one place nothing looks.
- `modules.json` `temporary_exceptions` is a second, disconnected ledger: `Policy.psm1:98-102` never
  opens a `.go` file, **3 of the 5 declared exception paths do not exist on disk**, and all 5 carry
  `removal_owner: M-10`, deferred since 2026-07-14 (L-01).
- Only 2 of 21 modules have a `boundary_test.go`, and the two are byte-identical
  (`channelfees`, `divergences`); `arch.ScanVendorTokens` was never pointed at `internal/modules`
  (L-02).
- `internal/platform/archguard` is a third instrument with 5 passing tests whose scope exclusion is
  asserted only in a code comment (L-07).
- The 234 violations decompose: 146 (62%) originate in `adapters/`, 42 in `application/`, 9 in
  `ports/`. **The structural leak is ports typed by the target's `domain`** —
  `connectors/ports/catalog_read.go:6,10,13-15` (L-05). 49 directed edges among 17 modules, 3
  bidirectional pairs, `erp_import↔internal_read` heaviest at 26 combined, `connectors` the heaviest
  sink at 64 inbound from 6 modules (L-06).
- ADR-023's own prose at `023-module-protocol.md:82,300,306` still says "35 violations" against a
  measured 234 (fact 9, D-55). The specification and the measurement have been out of contact
  through two corrections.

**Why these fix together.** One detector that **discovers roots by walking the tree instead of
matching a prefix**, one registry keyed on the `(kind, id)` pair, one shrink-only ratchet on the
count, and the exception ledger folded into the detector's own input so a declared exception that
does not exist on disk is itself a failure. Four instruments becoming one is the fix; tuning any one
of them leaves the other three to disagree.

**Ordering.** The 234 count only becomes meaningful once A2 runs the detector and A3 gives it a
must-fail fixture. Ratcheting a number nobody measures is theatre.

**Size:** 4–5 days for the instrument and the ratchet. Driving 234 → 0 is a separate, larger,
open-ended program and should be explicitly excluded from this axis; the axis's job is to stop the
number from growing.

---

## A8 — Nothing at runtime knows who is calling or which tenant owns the row, so every isolation control is decorative

**Cause.** The system has no identity at any layer. Not at the HTTP edge, not in the database
session. Every mechanism that would normally provide isolation is present in some form and inert
because its precondition is missing.

**Findings.**

- **Zero authentication on 94 route registrations** (`lanes/security.md` SEC-1, `lanes/delivery.md`
  F1). `composition/root.go:994` composes the entire stack as
  `CORSMiddleware(apierror.Recover(mux))` — there is no place in that expression where identity
  could be checked.
- **Buyer PII — full name, CPF/CNPJ, full address — is served by that unauthenticated API**
  (`lanes/security.md` SEC-2, `orders/transport/http_handler.go:608-618`, migration `0089`). The data
  is deliberately unmasked and *correctly so*: a Brazilian nota fiscal legally requires it, and the
  code says so. The defect is not the unmasking. **The defect is that the endpoint is open.**
- **And it is reachable from the internet by design, through two independent paths** (amendment 3;
  every lane left this open and PHASE-0 said so). Production: `deploy/Caddyfile` terminates TLS on
  `{$MPC_DOMAIN}` and matches `@orders_api { path /orders /orders/* ; not header Accept *text/html* }`
  → `reverse_proxy backend:8080`. A request carrying `Accept: application/json` lands on the PII
  handler. **The matcher's HTML exclusion means a browser gets the SPA and a script gets the
  data — the surface is, by construction, invisible to the way a person would check it.**
  Development: the `oauth` compose profile runs `ngrok http --url="$callback_host" frontend:5174`,
  and `apps/web/vite.config.ts` proxies `/orders`, `/pricing`, `/marketplaces/*`, `/profitability`
  and others to `backend:8080` — so while that profile is up the same surface is on a public URL.
  Whether a production host is live at this moment is unverified because `MPC_DOMAIN` is external.
  **That is the only part still unknown, and it does not change the sizing: the dev path is one
  `--profile oauth` away and exists because the Mercado Livre OAuth callback requires a public URL.**
  This axis is not sized as "localhost only," and my pre-amendment mitigation, which assumed it was,
  is withdrawn below.
- No per-request tenant: `pgdb/tenant.go:5-9`, `cfg.DefaultTenantID` (`lanes/delivery.md` F2).
  `pgdb/config.go:23-24` still silently defaults to `tenant_default` for `cmd/server`, with 70
  `cfg.DefaultTenantID` occurrences in `root.go`; D-39 fixed only `cmd/catalogingest`
  (`lanes/persistence.md` P-4, `lanes/security.md` SEC-3). A process-wide constant cannot distinguish
  two tenants under any schema.
- **RLS is on 4 of 64 tables and both of its preconditions fail**: the app DSN connects as the table
  owner and bypasses RLS (the `mpc_app` role from migration `0098` is `NOLOGIN`), and the
  `app.tenant_id` GUC is set at exactly one site,
  `contexts/catalog/internal/postgres/repository.go:49` (fact 12, D-44; `lanes/security.md` SEC-4;
  `lanes/persistence.md` P-1, P-3). 1 of 85 migrations enables RLS; 51 of 52 tenant-bearing
  migrations have none (SEC-5).
- **Tenant scoping is 100% hand convention across 246 raw query call sites with 0 automated checks**
  (`lanes/persistence.md` P-2). The lane verified 285/285 statements currently carry `tenant_id` —
  which is genuinely good work and is exactly why the axis is urgent: perfect adherence with no
  mechanism is one distracted commit from broken, and nothing would report it.
- CORS is `*` on all ~104 routes with `Authorization` in Allow-Headers (`httpx/router.go:20`;
  `lanes/delivery.md` F11, `lanes/security.md` SEC-6). See §7 D7 on its classification.
- Dead fail-open `pgdb.DefaultTenantID` helper, 0 callers (SEC-8) — delete it in this axis, since a
  fail-open helper sitting next to a fail-open default is the thing that gets picked up next.

**Why these fix together.** They are one missing thing: identity that flows from the request through
the context into the database session. Adding RLS without a non-bypassing role does nothing. Adding
a role without a per-request tenant does nothing. Adding auth without narrowing CORS in the *same
change* opens a credentialed cross-origin surface that is strictly worse than today's open one.

**The honest split, revised under amendment 3.** There is an edge control measured in hours, a
structural change measured in weeks, and a product decision I will not size:

- **~~Immediate: bind the container to `127.0.0.1` rather than `0.0.0.0:8080`.~~ Withdrawn.** It was
  premised on the surface being exposed only by a published container port. It is not. Caddy reaches
  `backend:8080` over the compose network, and ngrok reaches the frontend, which proxies to the same
  place. **A container-level bind would not have closed either path, and it would have produced a
  written record saying the surface was closed.** Recording the withdrawal rather than quietly
  replacing it, because a mitigation that reads as done and is inert is the exact defect class this
  audit is about — and mine was one.
- **Immediate (hours), at the edge, where the exposure actually is.** Two changes, both
  configuration:
  1. In `deploy/Caddyfile`, stop reverse-proxying the data surface to the public internet
     unconditionally. Either delete the `@orders_api` route until authentication exists, or put a
     control in front of it at the proxy — Caddy `basic_auth`, `forward_auth`, mTLS, or an IP
     allowlist. Caddy can enforce this today with no Go change, which is why it is hours.
  2. Stop the `oauth` profile from publishing the whole application. Mercado Livre needs exactly one
     public path — the OAuth callback. Everything else behind that tunnel, including the Vite proxy
     entries for `/orders`, `/pricing`, `/marketplaces/*` and `/profitability`, should not be
     reachable through it.
  **Negative fixture, and it must ship with the change:** a request to `/orders` with
  `Accept: application/json` from outside the allowlist must be refused. Assert the refusal, not the
  configuration — the `not header Accept *text/html*` matcher is a live demonstration that a surface
  can be invisible to a browser check while wide open to a script.
- **Then (5–8 days), structurally:** per-request tenant plumbing, `pgdb/config.go:23-24` made
  fail-closed to match the shape already proven by `MC_DATABASE_URL` and `MPC_ENCRYPTION_KEY`, a
  non-bypassing DB role, RLS extended with the GUC set on every path, and a boot assertion that the
  connecting role cannot bypass RLS.
- **Authentication itself is a product decision, not an engineering estimate, and I decline to put a
  number on it.** Who the users are, whether it is operator-only, and whether an IdP is in play are
  operator inputs. Sizing it here would be fluent invention.

**Interaction with the public flip.** Making the repository public does not itself expose the API,
but it publishes `deploy/Caddyfile`, `docker-compose.yml`, and `vite.config.ts` — that is, the exact
map of which paths are proxied where, and the `Accept`-header trick that reaches the data. The edge
control above should land **before** the flip, not after. This is a sequencing constraint, not a
reason to avoid going public.

---

## A9 — A failure leaves no durable trace, so the only detector is a person looking at a screen

**Cause.** When something goes wrong, the evidence is either absent or lives only in a stream nobody
is reading. Detection is therefore level 5 by construction — not by policy, but because no other
channel exists.

**Findings.**

- **`/healthz` returns 200 unconditionally** (`httpx/router.go:5-14`) **and it is the Docker
  healthcheck** (`docker-compose.yml:46`) (`lanes/observability.md` OBS-1). A container with a dead
  database reports healthy, forever.
- **5 of 6 ticker loops have no `recover()`** — one panic in a background job takes the whole
  process down. Only `sync/application/scheduler.go:201-208 safeInvoke` has one, and it covers only
  the job body (OBS-2).
- `FeeSyncScheduler` has zero logging in the file (`:52`, `:74`) (OBS-3). `runJob` skips silently on
  a cursor-read error (`:146-149`) and discards the errors from `RecordFailure` (`:154`) and
  `RecordSuccess` (`:160`) (OBS-4) — the failure-recording path itself fails silently.
- **`slog` is used in 36 files and no handler is ever configured**, so output is text, not JSON, and
  is not queryable (OBS-9). Two logging mechanisms coexist (OBS-10). There are zero `slog.Debug`
  calls anywhere, so the failure mode is always "not logged at all," never "logged at the wrong
  level."
- 13 of 15 ML adapter files never call `slog` (OBS-5). Zero metrics, zero tracing, no `/metrics`
  (OBS-7). No correlation or request ID (OBS-8; `lanes/delivery.md` F9).
- `cmd/catalogingest` failure produces no durable row (OBS-11, D-45). `/sync/health` collapses
  `installation_id` while the FE keys on `entity.entity` (OBS-6).
- ~9 discarded `json.Unmarshal` errors on read paths (`order_repo.go:297,413,1109,1112`, others)
  (`lanes/goidiom.md`) — a malformed stored payload silently becomes a zero value, which is the
  precise thing the repo's own doctrine forbids.

**Explicitly closed, do not reopen.** The ML token-refresh failure is now durable and tested
(`auth_flow_service.go:454-574`, `integrations_refresh_failure_test.go`) (`lanes/observability.md`).
This contradicts the prior memory record that it was invisible at every layer. The lane's evidence is
better and more recent. Recording it loudly here so nobody re-audits it.

**Why these fix together.** One decision — every failure writes a durable row before it writes a log
line, and health reflects dependencies — plus one JSON handler configured once at boot. Twelve
separate logging improvements produce twelve dialects of the thing OBS-10 already complains about.

**Size:** 4–6 days. The `recover()` wrappers and the real `/healthz` are a day and should not wait
for the rest.

---

# 2. Kill the noise

These are real findings. They do not earn a slot. For each: what it costs to fix versus what it costs
to leave.

**Leave forever.**

- **DUP-6 — 82 hand-written `for rows.Next()` loops, 0 `pgx.CollectRows`.** Fixing this replaces
  explicit `rows.Scan(&a, &b)` — where the compiler checks arity and types — with
  `pgx.RowToStructByName`, whose correctness depends on **struct tags matching column names, checked
  at runtime, by hand**. That is a *new* hand-synchronised seam, which is the exact defect class this
  entire audit exists to eliminate. Trading a compiler-checked verbosity for a runtime-checked
  convenience is a regression dressed as cleanup. The boilerplate has produced zero reported defects.
  **See §7 D6: this is not deferred, it is refused.**
- **DUP-10 — 83 `strings.TrimSpace(x) == ""` sites.** Idiomatic Go. A helper saves 4 characters and
  adds an import.
- **G-4 — 170 `errors.New("SCREAMING_SNAKE")` sites, 157 following the same deliberate convention.**
  92% consistency on a chosen convention is not a defect. The actual problem with error handling is
  A5's `strings.HasPrefix` dispatch, and renaming the errors would not touch it. Leave the naming.
- **DUP-11 — `Assert-True` defined 6× in test scripts.** Test-local, no runtime blast radius.
- **G-2, G-3, G-6, G-8** — inconsistent clock injection (6 files), 2 `fmt.Errorf` with no verbs, one
  `%v` where `%w` belongs at `auth_flow_service.go:478`, 200 `map[string]any` with only 4 in
  `domain/` and 0 domain-layer key access. `go vet` is clean; 355 `%w` against 12 `%v`. These are
  within normal variance for a codebase at solid-professional calibration. The one exception:
  `classifications/application/service.go:46,74` drops `.UTC()`, which is a latent timezone bug worth
  a 2-line fix — fold it into whatever change next touches that file, never its own work item.
- **G-7 — 27 of 35 adapter interfaces unexported.** Correct Go, not a finding.
- **`lanes/layering.md` L-08 — no detector looks outside `.go`, so the 77 FE vendor-token occurrences
  are uncovered.** They are not a violation of any stated rule. Writing the rule first is the
  prerequisite, and there is no evidence the rule is needed.
- **`lanes/frontend.md` F15 — accessibility baseline is OK.** Not a finding; it is in §3.
- **SEC-7 — raw `recover()` value written to a log sink.** Log-only, never client-exposed. It becomes
  a real exposure only if the logs become a queryable product surface; revisit then, and A9 will
  touch that file anyway.
- **OBS-12 — webhook stats permanently zero.** The subsystem was never built. Statistics for a
  non-existent subsystem are not a defect; building the subsystem is a product decision.
- **`lanes/frontend.md` F12 — 3 empty workspace packages plus the `apps/landing` shell.** Delete-on-
  sight if something else touches them. Not worth a change of their own.

**Fold in, do not track separately.**

- **`lanes/frontend.md` F11 — dead `pages/DashboardPage.tsx`, 160 lines, 2 live SDK calls.** Delete
  it during A4, where the SDK surface is being regenerated anyway. Tracked alone, it competes with
  real work.
- **`lanes/persistence.md` P-8 — duplicate migration numbers 0021 and 0093.** `GOV_MIGRATION_PREFIX`
  already fires on this; `lanes/cicd.md` F6 counts it among the 58 governance violations at HEAD. It
  is not its own finding, it is one row of A2's baseline debt. The runner sorts by filename
  (`migrate/runner.go:25`), so the ordering is deterministic today.
- **`lanes/persistence.md` P-7 — `stale_since` written at 2 sites, read at 0.** Either wire the
  reader or drop the column, during A9 when that area is open. It is dead weight, not a hazard.
- **`lanes/persistence.md` P-9 — 14 `fmt.Sprintf` WHERE sites, all sampled safe.** Add the check to
  the linter during A2 so it cannot regress; do not chase the existing 14.
- **`lanes/duplication.md` DUP-8 — 2 near-duplicate policy_reader adapters.** Collapses naturally
  under A5 or A6, whichever reaches it first.
- **`lanes/observability.md` OBS-10 — two logging mechanisms.** This is A9's fix, not a finding
  beside it.
- **`lanes/cicd.md` F12 — `knowledge-routes.json:12` cites a nonexistent test file.** One-line fix,
  but the systemic answer is A2 running a validator that would have caught it. Fix it when the
  validator is wired, and use it as that validator's negative fixture.
- **F15 (cicd) — a warn-only `wiki-lint` PR workflow existed on a legacy branch, failed 3/3, never
  merged.** History. Its only value is as the precedent that warn-only gates get ignored, which
  strengthens §4's "block, do not annotate."

**Explicitly demoted from the lanes' own weighting.**

- **D-52 / the 635 CRLF gofmt artifacts.** Under the amendment this stops being a gate question:
  `ubuntu-latest` checks out LF, so pure-CRLF failures cannot occur on a runner, and only the **22
  genuinely unformatted files** survive. The `.gitattributes` renormalisation decision shrinks to
  local developer experience. **It does not vanish entirely** — a Windows operator running any local
  fmt check still sees 639 — but the answer is now "scope the local check or renormalise when
  convenient," not "this blocks the program." Re-sized down, deliberately.

---

# 3. Consolidated do-not-touch list

Merged from every lane's "what is actually fine." Anything on this list that a later change breaks is
a regression, not a refactor. Several entries are also the **reference shape** the corresponding axis
should copy — marked **[TEMPLATE]**.

**Verification and architecture**

- `internal/arch/scan_test.go` — violations fixture plus clean fixture. **[TEMPLATE for A3.]** This
  is the model `TestModuleBoundaryADR023` must be rebuilt against.
- `internal/platform/archguard` and its 5 must-fail fixtures — 5/5 passing. **[TEMPLATE]**
- `GOV_FRONTEND_FETCH` — a real, enforced rule with 0 violations. Proof the governance harness works
  when a rule is fully implemented.
- `Policy.psm1:198-206` exclusion list — deliberate and correct.
- The 8 `panic` sites, exactly matching `invariants.json`. An exact match between a declared
  invariant set and the code is rare here; do not disturb it.
- `go vet ./...` and `go build ./...` clean.
- 141 of 280 interfaces in `ports/`, 0 in `domain/` — the port discipline is real where it is
  applied.

**Transport**

- `apierror.Write` as the single envelope producer — 192 call sites, 0 `http.Error` bypasses.
  CHIP-ERROR-UNIFY landed. **Do not re-litigate the envelope**; A5 extends this, never replaces it.
- `route_deadline.go` — deadline handling is correct.
- The one `errgroup` use.

**Persistence**

- The ephemeral `mpc_test_<hex>` database per integration run, and
  `Invoke-Integration`'s fail-closed `HPG_EXTERNAL_TARGET_FORBIDDEN` guard at
  `scripts/harness.ps1:69`, which refuses to run against an externally supplied DSN. **[TEMPLATE for
  a fail-closed lane guard.]**
- `migrate/runner.go` per-migration transaction.
- All 116 `timestamptz` columns; 0 `CONCURRENTLY`.
- The `::text` money read discipline at `pricing/matrix_reader.go:34-39,64-68` — this is the target
  shape for A6's persistence side.
- The ERP→Postgres mirror write path (`mirror_repository.go:71-172`). **Oracle/Sankhya stays
  read-only. No axis in this document proposes an ERP write.**

**Domain**

- `pricing/domain.Money` and `FormatRatHalfUp` — **[TEMPLATE for A6.]** The correct answer already
  exists in this repo; A6 is adoption, not invention.
- `orders/adapters/pricingtax/reader.go:266 lineTotal` guard.
- The DIFAL override path — real, wired, and audited (`lanes/duplication.md` contradicts the prior
  memory record here; the lane is right and current).

**Security and integration**

- AES-256-GCM credential encryption; the Oracle config loader. **Both fail closed — [TEMPLATE] for
  every boot guard A8 adds.**
- The fail-closed inventory: `MC_DATABASE_URL`, `MPC_ENCRYPTION_KEY`, Oracle credentials,
  `MPC_PROVIDER_WRITES_ENABLED`, `MPC_ML_CATALOG_OFFERS_ENABLED`,
  `MPC_ASSISTED_SANKHYA_LINKAGE_ENABLED`.
- PII redaction on the raw-capture path, and its tests.
- **The deliberate non-masking of the fiscal block (`orders/transport/http_handler.go:608-618`).
  Brazilian nota fiscal law requires the buyer's full identification. Do not "fix" this by masking.**
  A8 closes the door to the room; it does not redact the legally required contents.
- The ML token-refresh durable failure path and `integrations_refresh_failure_test.go` — fixed,
  tested, closed.
- `sync.Scheduler.safeInvoke` — the only correct ticker in the repo. **[TEMPLATE for A9's other 5.]**
- The Sankhya/Oracle read path and the Mercado Livre integration.

**Frontend**

- `/sync/health` plus `SyncHealthCard.tsx:127-190`'s exhaustive `query.status` switch. **[TEMPLATE
  for A5's frontend half.]**
- `packages/ui/src/money.ts` and its honest-unknown rendering — null or em-dash, never 0.
- `tsconfig.base.json` `"strict": true`.
- Zero `fetch()` bypasses of the SDK.
- The 125 local type declarations the frontend lane verified as legitimate FE-local state — **these
  are not shadow DTOs and A4 must not "consolidate" them into generated types.**
- The accessibility baseline (`lanes/frontend.md` F15).

**Infrastructure**

- The Docker dev stack; `scripts/release.sh` as a deploy mechanism (its CI-bypass header stops being
  true once A1 lands, at which point the header is what changes, not the script).
- `release-images.yml`'s build/push structure — buildx, GHCR, `type=gha` cache, digest tags. It is
  well built. The only thing wrong with it is what triggers it (A1).

---

# 4. The target gate topology

Firing hierarchy: **1 unrepresentable → 2 boot-fatal → 3 red build → 4 runtime assertion → 5
discipline. Level 5 is not a control.** Every guard below ships with the input that makes it fail, in
the same change. A guard with no negative fixture does not count as landed.

Two constraints from the brief shape everything: **human PR review is not available** (solo operator
plus agents), and **the author of nearly all code is a machine whose failure mode is fluent,
internally consistent wrongness — the wrong noun used correctly everywhere, including inside the
guard written to catch it.** The second is why nearly every gate below is paired with a fixture: an
agent will write a guard that is coherent, well-named, and inert, and only an input that must fail
distinguishes that from a working one. `GOV_CONTEXT_UNREGISTERED` (`lanes/cicd.md` F9) is that exact
specimen, already in the repo.

## 4.0 The precondition — where the gate stands

A required status check runs on a pull request. Today there are none. **A1 is not a gate; it is the
requirement all gates attach to.**

**The availability question is closed.** The repository goes public on a free personal account, so
rulesets and enforced required status checks cost nothing. The credential-custody fallback that a
private/free repo would have forced is dead and is not designed for anywhere below. Every gate in
this section assumes an **enforced** check the author cannot wave through.

Order of operations, per amendment 2 — gates before push:

1. Build `verify.yml`, the ratchet baselines, and the `release-images.yml` guard locally, inside the
   existing 98-commit backlog (this is A2 + A3's output).
2. Drive them green locally.
3. Clear the secrets/PII sweep over the full history.
4. Flip to public.
5. Fast-forward `main`. The first remote execution is a green gated tree, not a red archaeology run.
6. Configure the ruleset on `main`: required checks `verify-fast` and `verify-full`, require a pull
   request, **linear history**, block force-push and deletion, and **an empty bypass list**.
7. From that point, every change is a PR.

**Two settings do the real work in step 6 and both are easy to get wrong.** A required check that
never reports must block the merge as *expected but not received* — that is what makes deleting the
workflow a losing move rather than a winning one (§4.5). And **the bypass list must be empty**: on a
personal-account repository the owner is inherently an admin, and a ruleset with the owner on its
bypass list is decorative for exactly the actor whose credential an agent operates under. Verify both
deliberately, once, with the fixture in §4.5.

## 4.1 Level 1 — unrepresentable

The strongest and the cheapest to maintain, because nothing has to fire.

| Gate | Makes unrepresentable | Negative fixture |
|---|---|---|
| **Exact money type, no `float64` constructor or accessor** (A6) | Constructing money from, or extracting money into, a binary float. The `strconv.ParseFloat` sites and the 373 `float64` occurrences stop compiling rather than being flagged. | A `_test.go` that attempts `Money(1.5)` and `m.Float64()`; the build must fail. Compile-fail fixtures live in a `//go:build ignore` file asserted by a script that expects a non-zero `go build`. |
| **Generated request/response types** (A4) | A hand-edited SDK type diverging from OpenAPI. The copy stops existing rather than being compared. | CI regenerates and runs `git diff --exit-code` on the generated tree. Fixture: a PR hand-editing one generated field must go red. |
| **Postgres enum types replacing 65 `CHECK (… IN (…))`** (A6, P-6) | Inserting a value outside the vocabulary, and vocabularies drifting between copies. | A migration test inserting an invalid value; must error. |
| **Typed error registry with an exhaustive switch** (A5) | An error code with no HTTP status mapping. | Add a code without a mapping; the build must fail. |
| **Generated route-prefix tables** (A4, amendment 3) | A backend route prefix existing in the Go router but absent from `deploy/Caddyfile` or `apps/web/vite.config.ts`. Both tables are emitted from one authority instead of hand-maintained under a comment. | CI regenerates and runs `git diff --exit-code`. Fixture: add a route to the router without regenerating; must go red. **Interim, before generation exists: a `verify-fast` step that parses all three and fails on disagreement — a comparison check is the fallback, not the target.** |

`tsconfig.base.json` `"strict": true` already sits at this level and is the reason the frontend
lane's F3 defect is *findable at all* — the checker just never runs (A2).

## 4.2 Level 2 — boot-fatal

Refuse to start rather than run wrong. The repo already does this correctly for `MC_DATABASE_URL`,
`MPC_ENCRYPTION_KEY`, and the Oracle loader — those are the template.

| Gate | Fires | Blocks | Negative fixture |
|---|---|---|---|
| **Tenant resolution fail-closed** (A8, P-4/SEC-3) | `cmd/server` boot | Process refuses to start with no explicit tenant source. Deletes the `tenant_default` fallback at `pgdb/config.go:23-24` and the dead `pgdb.DefaultTenantID` helper (SEC-8). | Boot with the variable unset; must exit non-zero with a named reason. |
| **DB role cannot bypass RLS** (A8, fact 12/D-44) | boot | Asserts at startup that the connecting role is not the table owner and not `BYPASSRLS`. Refuses otherwise. | Boot against the owner DSN; must refuse. **This fixture is the whole point** — it is the assertion that would have caught D-44 the day RLS was written. |
| **`slog` JSON handler configured** (A9, OBS-9) | boot | Not fatal on its own; fatal if the log destination is unwritable. | Point at an unwritable destination; must exit rather than silently degrade to text. |
| **Route table has no unauthenticated route outside an explicit allowlist** (A8) | boot | Refuses to start if a registered route is absent from the auth allowlist. Applicable only after auth exists. | Register a route without an allowlist entry; boot must fail. |
| **PII-bearing route is never composed without an identity check** (A8, amendment 3) | boot | `composition/root.go:994` currently reads `CORSMiddleware(apierror.Recover(mux))`; assert at boot that every route in the PII set is wrapped by the identity middleware. This is the assertion that would have made the Caddy exposure impossible rather than merely undiscovered. | Compose the PII route without the middleware; boot must fail. |

## 4.3 Level 3 — red build (the primary enforcement layer)

One workflow, `verify.yml`, on `pull_request`, `ubuntu-latest`. **The same command runs locally**; a
gate the operator cannot reproduce on their own machine is its own defect class.

**Job `verify-fast`** — every push to a PR branch, target under 3 minutes:

| Step | Closes | Notes |
|---|---|---|
| `gofmt -l` over **the whole Go module**, not just `internal` | A2 [ARM A] | On a Linux runner the 635 CRLF artifacts cannot occur; only the 22 real ones remain. This is where D-52 dissolves. |
| `go vet ./...`, `go build ./...` | baseline | Currently clean — keep it that way. |
| `tsc --noEmit` per workspace | A2/frontend F2, F3 | 12 errors today, 3 in production code. Must be paid down in the same change that wires it. |
| archscan over **all four roots including `internal/composition`** | A7, L-03, cicd F7 | 44 findings today. Ships with a shrink-only ratchet file; the ratchet's negative fixture is a PR that adds one violation. |
| governance validate + drift, `-BaseSha` derived from `git merge-base origin/main HEAD` | A2, cicd F13 | 58 violations at HEAD with zero diff. Baselined on day one; ratchet shrink-only. |
| **Mixed-change ban** | §4.5 | See below. |
| **Count assertion on every step** | A3 | Each step prints an attributable count; a step reporting zero executed units fails. |

**Job `verify-full`** — on PR and on merge to `main`, target under 12 minutes:

| Step | Closes | Notes |
|---|---|---|
| `go test ./...` — the whole module | A2 T-1/T-8 | Brings 369 currently-dark test files into the light. Expect the first run to be red; that is the point. |
| vitest across **all** workspaces, glob-based, **no filename pins** | A2 T-3/T-10 | Removes the `apps/web/vitest.config.ts:14` name-pin and the orphaned root config. |
| Pester over `scripts/tests/*.tests.ps1` | A2 T-7, cicd F10/F11 | ~250 assertions, 6/11 red cold. `pwsh` is present on `ubuntu-latest`. |
| **Integration lane** | A2, T-6 | **[ARM A] This is CI-able and no lane said so.** `Invoke-Integration` (`scripts/harness.ps1:68-88`) needs only Docker plus `pwsh`, both present on `ubuntu-latest`; it provisions its own ephemeral Postgres and embeds its migrations. It needs **no** ERP and **no** dev stack. Confirm by running the lane on a runner once. Cost: ~259s measured locally, so budget it as the dominant term in `verify-full`. Must emit RUN/PASS/SKIP/FAIL counts before it is trusted (T-6). |
| `git diff --exit-code` on generated SDK/types | A4 | After A4 lands. |

**Public-repo threat surface — a design constraint on every workflow, not a footnote.** Going public
means fork pull requests execute workflows and every log line is world-readable. Stated per workflow:

| Workflow | Trigger | `permissions:` | Fork-reachable? | Ruling |
|---|---|---|---|---|
| `verify-fast` | `pull_request` | `contents: read` **only** | Yes | Safe. Runs with a read-only token; GitHub withholds secrets from fork PRs, and this job needs none. |
| `verify-full` | `pull_request`, `push: [main]` | `contents: read` **only** | Yes (PR half) | Safe. The integration lane provisions its own ephemeral Postgres in-runner and needs no secret. |
| `release-images.yml` | `push: [main]`, `tags: v*`, dispatch | `contents: read`, `packages: write` | **No** — `push` never fires from a fork | Correct as written. Its write scope is why it must never gain a `pull_request` trigger. |

Three rules follow, and all three are enforceable:

- **`pull_request_target` is banned outright.** It runs the base repo's workflow with a writable token
  against attacker-controlled head code — the textbook privilege-escalation path on a public repo.
  There is no gate in this document that needs it. Add a `verify-fast` step that greps `.github/`
  for the string and fails; its negative fixture is a PR introducing it.
- **`permissions:` is declared explicitly at the top of every workflow**, never inherited. Default
  token scope is a repository-wide setting an agent could widen; an explicit block in the file makes
  the intent reviewable in the diff. Fixture: a workflow with no `permissions:` block fails the same
  grep step.
- **Enable "require approval for all outside collaborators" on fork workflow runs.** Costs nothing,
  and without it a fork PR burns runner minutes on demand.

**Log safety is now a correctness requirement, and it interacts with A3.** Every count, failure
token, and diagnostic this topology adds is published to the world. Two concrete consequences:
`Invoke-Integration` generates a random Postgres password (`scripts/harness.ps1:80-82`) which must
never reach stdout, and no lane may print a DSN, a tenant identifier, or a row of buyer data in a
failure message. The A3 convention is therefore "**attributable and publishable**" — a failure token
that names a test is both; a failure message that dumps the failing row is neither. **[ARM A]** —
confirm by reading each lane's failure-path output before the public flip; this is a check the
running secrets/PII sweep does not cover, because the sweep reads history and this is about future
output.

**Ratchets, not thresholds.** Three shrink-only baselines: archscan findings (44), governance
violations (58), ADR-023 boundary violations (234). Each is a file in the repo; CI fails if the
measured number exceeds it and fails if the file was edited upward in the same PR as the code that
made it grow. **Every ratchet ships with a fixture PR that adds one violation and must go red.**

**What CI blocks versus annotates: it blocks.** `lanes/cicd.md` F15 is the local precedent — a
warn-only `wiki-lint` workflow failed 3 of 3 runs and was ignored into oblivion. Nothing in this
topology is warn-only. CodeRabbit, when it arrives, is the one exception: it annotates, because it is
judgment, and judgment is not a gate.

## 4.4 Level 4 — runtime assertion

For what cannot be known until a request exists.

| Assertion | Fires | Behaviour |
|---|---|---|
| **Query without `tenant_id` in a tenant-bearing table** (A8, P-2's 246 call sites) | per query, dev/staging fatal, prod logged + durable row | The only mechanism that can cover 246 hand-written sites. Fixture: a deliberately unscoped query in a test must panic. |
| **`/healthz` reflects dependencies** (A9, OBS-1) | per probe | Returns non-200 when the DB is unreachable. Fixture: an integration test with the DB stopped asserting non-200 — today it returns 200 and Docker calls it healthy. |
| **`recover()` on all 6 ticker loops, writing a durable row** (A9, OBS-2) | per tick | A panic degrades one job instead of killing the process, and leaves evidence. Fixture: a job that panics; assert the process survives *and* the row exists. |
| **React error boundaries** (A5, frontend F4) | per render subtree | A component throw degrades one panel. Fixture: a test component that throws; assert the boundary renders and the rest of the page survives. |

## 4.5 The hole the amendment named: the author can edit the gate

The agent writing the code has write access to `.github/`. Anything it can edit, it can weaken —
fluently, with a plausible commit message. Three layers, in order of strength:

1. **The required-check configuration lives outside the repo, in ruleset state the operator holds.**
   This is the layer that terminates the regress: **deleting `verify.yml` does not make the check
   pass — a required check that never reports blocks the merge as "expected, not received."** An
   agent cannot delete its way out. This is now unconditionally available (public repo, §4.0), so
   unlike in the pre-amendment draft it is the primary control and not a hope.
   **Negative fixture for the gate topology itself:** open a PR that deletes `verify.yml` and confirm
   the merge is blocked. Do this once, deliberately, before trusting anything else in this document.
   Pair it with a second fixture — a PR that adds the operator to the ruleset bypass list — to
   confirm the bypass list is empty and stays empty.
2. **Mixed-change ban, in `verify-fast`.** Fail any PR whose diff touches both a gate definition and
   code under that gate's jurisdiction. Jurisdiction set: `.github/**`, `scripts/harness/**`,
   `contracts/governance/**`, `internal/arch/**`, and the three ratchet files. Expressible as a
   command that exits non-zero, therefore level 3, therefore mandatory rather than optional. Negative
   fixture: a synthetic PR touching `Policy.psm1` and a module together must go red.
3. **CODEOWNERS: propose it, and be honest that it buys very little here.** `developmentconexus-ops`
   is a personal account, not an organization. There are no teams. CODEOWNERS entries must name users
   with write access, and there is exactly one such user — the operator — who is also the repository
   admin. So "require review from Code Owners" resolves to "the operator approves," which is the same
   credential custody that already exists, re-expressed as a file. It adds one genuine thing and it
   is worth the ten lines for that alone: **the PR page names the protected path that was touched**,
   so a diff that quietly edits `Policy.psm1` under a feature-shaped title is labelled rather than
   scanned for. That is annotation, not enforcement, and it should be recorded as annotation.

**Where the real separation of powers actually comes from, on a solo public repo.** Not from a second
reviewer, which does not exist. It comes from the split between two credentials: the agent holds
repository write and can open PRs; **the operator alone holds the admin credential that edits ruleset
state.** The gate binds the agent completely and binds the operator not at all — and that is the
correct arrangement, because the operator is the trusted principal and the agent is the one whose
failure mode (fluent, internally consistent wrongness) this whole topology exists to catch. It is
worth stating explicitly so nobody later mistakes the operator's ability to bypass for a hole. **The
one thing that would turn it into a hole is the agent ever operating under the admin credential.
Don't.**

## 4.6 What legitimately stays at level 5, stated rather than disguised

Per amendment point 4, and only these:

- **The Oracle/Sankhya live lane.** `MPC_ORACLE_LIVE_TEST=1`, cgo, 6 skipped tests
  (`lanes/testing.md` T-11). **No ERP credential goes into a runner.** Local only. It also carries a
  known trap: `//go:build cgo` on a host without gcc yields "no tests to run" — green without
  executing — which is A3's failure mode in the one lane that must stay at level 5. So it needs the
  count assertion *more* than anything in CI does.
- **The Docker dev-stack live drive and browser QA.** Needs the full compose stack and a real
  browser. Local, operator-driven.
- **Judgment.** Whether a decomposition is right, whether an ADR is coherent, whether a name is
  honest. CodeRabbit assists here and annotates; it does not block.

Everything else that was going to be parked at level 5 has been moved up. Under the amendment,
proposing level 5 for a rule a workflow can express is a defect in the recommendation, and this
section is deliberately short as a result.

## 4.7 Cost

Amendment point 5: availability is not the same as cost.

- **`ubuntu-latest` only. No matrix. No `windows-latest`** (2× minute multiplier, and it would
  reintroduce the CRLF problem the Linux runner dissolves). **No `macos`** (10×).
- Split fast/full exactly as above so the ~259s integration lane does not run on every intermediate
  push.
- Cache `GOCACHE`/`GOMODCACHE` and npm via `actions/setup-go` and `actions/setup-node`. Free, and it
  is most of the difference between a 12-minute and a 4-minute `verify-full`.
- Rough envelope: `verify-fast` ~2–3 min, `verify-full` ~8–12 min. At ~30 PR pushes a week that is
  roughly 1,300–1,600 minutes a month.
- **Public repository on standard runners: minutes are not billed.** The cost question that shaped
  the pre-amendment draft is gone, and the fast/full split above survives on a different
  justification — **feedback latency, not money.** A gate that takes 12 minutes on every intermediate
  push is a gate an agent learns to work around; one that answers in under 3 is one it works with.
  Keep the split for that reason and say so, rather than leaving a cost rationale in place that no
  longer applies.
- Caching stops being a cost control and becomes purely a latency control. Still do it; the
  justification changed.
- The one thing that still needs watching is not minutes but **concurrency and queue depth** on a
  free account, plus fork PRs consuming runners — hence the outside-collaborator approval setting in
  §4.3.
- `release-images.yml` already uses `type=gha` caching and is well built; leave its cost profile
  alone.

---

# 5. The sequence

Dependency order, agent-driven days, and for each: which later axis it makes cheaper.

**Two amendments moved this sequence, in opposite directions.** Amendment 2 pushed A1 *down* — gates
land before the push, so A1 closes the first phase instead of opening it, and A2/A3 are built and
driven green locally against a remote that was reachable all along. Amendment 3 pulled A8's edge half
*up* to the front, because the mitigation I had written for it was premised on a localhost assumption
that measurement destroyed.

| # | Axis | Days | Makes cheaper |
|---|---|---|---|
| 0 | **A8 — edge control** (Caddy `@orders_api`; ngrok scope) | hours | Nothing. Pure risk reduction, and it is a precondition of the public flip. |
| 1 | **A2 — one `verify` product, built and green locally** | 4–6 | Every subsequent axis. This is the multiplier. |
| 2 | **A3 — the signal proves execution** | 3–4 | A7 (a ratchet on an unproven detector is theatre), and every guard afterwards. |
| 3 | **A1 — flip public, push, enforce** | 1 | Converts A2+A3 from a local habit into an enforced requirement. Without it they are level 5. |
| 4 | **A8 — identity, structural** (5–8 days; auth itself unsized) | see §1 | Nothing technically. It is here on dependency, having already been de-risked at step 0. |
| 5 | **A4 — generate the contract and the route tables** | 4–6 | A5 (types the kit binds to) and A6 (money on the wire). |
| 6 | **A5 — the boundary kit** | 5–7 | A9 (one error path to instrument instead of 12). |
| 7 | **A6 — one money vocabulary** | 7–10 | Nothing after it. It is last among the large ones because it is the most dangerous. |
| 8 | **A7 — one boundary instrument** | 4–5 | Bounds the 234; does not chase it. |
| 9 | **A9 — durable failure** | 4–6 | — |

**Total: roughly 33–46 days of agent-driven work, plus an unsized authentication decision.** Stated
honestly: this is a program, not a sprint. It is also mostly parallelisable after A1 — A4/A5/A6 touch
transport and domain, A7 touches tooling, A9 touches adapters and boot, so 2–3 chips can run
concurrently once `verify` is enforced and can catch their collisions. **Running them in parallel
before A1 is a mistake**, because until the check is required, a collision is caught by whoever
remembers to look.

Three hard constraints on this order, all non-negotiable: the secrets/PII sweep over the full history
clears **before** step 3; `release-images.yml` is guarded **before** the fast-forward, so the first
push does not publish `latest` from a tree this audit has already documented as failing; and **step 0
lands before the public flip**, because the flip publishes `deploy/Caddyfile` and `vite.config.ts` —
the map of which public paths reach the PII handler and the `Accept`-header condition that gets past
the HTML matcher.

## Justifying the first axis

**A2 is first, and it is first for a reason that survived both amendments.** Every later axis is a
large mechanical change to code whose test suite is 91% unreachable — 369 of 405 Go test files, 111
FE cases, 250 PowerShell assertions, and a type checker that has never run. Doing A4, A5, A6, or A9
before A2 means making thousands of edits with no way to know what broke. A6 in particular — 78 files
of money conversion — is close to reckless without it. A2 is also the axis whose cost is most likely
to be underestimated, because the first honest run surfaces debt that has been accumulating in the
dark for months; that debt is not a surprise, it is the reason to run it.

**A1 is deliberately not first, and that is a change from my pre-amendment draft.** When the remote
looked unreachable, A1 was the blocking unknown. It is not: the remote is 0 behind and
fast-forwardable, and the blockage was two GitHub identities in one keyring. What remains is one day
of work whose *value* is entirely conditional on A2 and A3 existing first — pushing to a remote and
enabling required checks that run nothing would enforce nothing. **The operator's sequencing
instruction and the technical dependency agree here, which is worth noting explicitly, because they
often would not.**

**A3 sits between them because A2 without A3 produces a green wall**, and a green wall on a public
repo is worse than one in private: the badge is visible, the logs are world-readable, and the
confidence is contagious. This repo already contains a test that passes because its assertion is
unreachable (`lanes/cicd.md` F9) and a lane whose all-skipped output is byte-identical to green
(T-6). Enforcing more checks before the ability to read their output is fixed manufactures exactly
the false confidence the program exists to eliminate.

**Step 0 — the edge control — comes before everything, and amendment 3 is why.** An unauthenticated
API serving CPF/CNPJ and full addresses is the only finding in this audit with a legal dimension, and
it is now measured as internet-reachable by design through two paths, not hypothetically reachable.
The repository is also about to become public, which publishes the routing map. Hours of work at the
proxy. Do not sequence it behind anything.

**What changed in my ordering, stated explicitly.** Before amendment 3 I placed A8 fourth and
justified it as *"on risk, not on dependency,"* with a one-day container-bind mitigation up front.
Two things were wrong with that. The mitigation would not have worked — Caddy and ngrok both reach
the backend without traversing a published container port — and the deferral rested on an unmeasured
reachability assumption that every lane had flagged as open and I treated as benign. **The correct
reading of "unverified" was not "probably fine."** The structural half of A8 stays at position 4,
because that is a genuine dependency (it changes the data path of every repository in the system, and
nobody should attempt that while 91% of the tests are dark). The exposure half moved to zero.

*(Superseded reasoning, recorded so the change is auditable: my pre-amendment draft justified A1 as
first on the grounds that the remote was unreachable and the ~96 commits had nowhere to go. The
premise was false — see §7 D8 — and the conclusion moved with it.)*


**A6 is deliberately last among the large axes.** It is the biggest blast radius and the one most
likely to be broken by a fluent, confident, wrong mechanical edit — the exact failure mode the brief
names. It should be attempted only when `verify` can prove that a 78-file change broke nothing, and
that is A2's output.

**Sequencing constraint, non-negotiable:** if authentication lands (A8), **CORS must be narrowed in
the same change**, never after. `Access-Control-Allow-Origin: *` combined with `Authorization` in
Allow-Headers is meaningfully worse with credentials than without them. See §7 D7.

---

# 6. Inversion tests

One line per structural conclusion. Each states what survives if the current implementation were the
opposite in every respect.

1. **The requirement, not the route, is what makes a verifier real (A1):** survives an opposite
   repository layout, an opposite branching model, and an opposite hosting choice, because a verifier
   the author can decline to invoke has the same effect as no verifier in every topology — which is
   why an open, unused push path and a blocked one were indistinguishable in effect until one was
   made mandatory.
2. **Verification as one command with a declared manifest (A2):** survives an opposite repo layout and
   an opposite toolchain, because a control whose invocation depends on human memory is equal to no
   control regardless of what it checks.
3. **Every guard ships a negative fixture (A3):** survives an opposite test topology and an opposite
   assertion library, because a guard never observed to fail is indistinguishable from a guard that
   cannot fail, in any framework.
4. **Attributable counts on every lane (A3):** survives an opposite runner and an opposite output
   format, because a run that executed nothing and a run that passed everything are the same bytes
   unless something counts.
5. **One authored copy of the contract (A4):** survives an opposite schema design, an opposite
   protocol, and an opposite client language, because N hand-maintained copies of one fact drift at a
   rate set by human attention, not by the fact's shape.
6. **Error identity as a type, not a string prefix (A5):** survives an opposite error taxonomy and an
   opposite HTTP status map, because `strings.HasPrefix(err.Error(), …)` couples the client-visible
   status to prose that no compiler and no rename tool protects, under every taxonomy.
7. **One provided edge kit rather than per-module edges (A5):** survives an opposite module
   decomposition, because a correct pattern with no propagation mechanism reaches exactly the files
   its author touched — which is what the D-11 fix reaching 1 of 10 files demonstrates.
8. **Money with no float64 constructor (A6):** survives an opposite domain model and an opposite
   database, because IEEE-754 binary floating point cannot represent 0.10 and Brazilian fiscal
   calculation is exact by regulation — a property of the domain, not of this code.
9. **Vocabulary in the type system, not in 65 CHECK constraints (A6):** survives an opposite schema,
   because a vocabulary duplicated across N definitions requires N edits to extend and offers no
   mechanism that notices when it got N−1.
10. **Detector discovers roots by walking the tree, never by matching a prefix (A7):** survives an
    opposite directory layout, because a prefix-anchored detector goes silent exactly when the tree
    moves, and a tree moving is the definition of the migration it is supposed to police.
11. **Registry keyed on the `(kind, id)` pair (A7):** survives an opposite module naming scheme,
    because a folder with no registry entry yields zero findings rather than a failure, under any
    naming scheme.
12. **Request-scoped tenant identity (A8):** survives an opposite schema design, because a
    process-wide constant cannot distinguish two tenants in any topology, and row-level security is
    ineffective for a bypassing connection in every topology.
13. **Fail-closed configuration (A8):** survives an opposite deployment model, because a default that
    silently substitutes a plausible value converts a configuration error into a data-integrity
    error, and the second is undetectable.
14. **Durable record before log line (A9):** survives an opposite observability stack, because an
    operator not watching a terminal at 03:00 has no channel to a signal that existed only in stdout.
15. **Health reflects dependencies (A9):** survives an opposite orchestrator, because an
    unconditional 200 makes the orchestrator's restart logic inoperative under every orchestrator.
16. **Gate definitions protected from the code they gate (§4.5):** survives an opposite CI provider
    and an opposite hosting platform, because an author who can weaken their own verifier has, in
    effect, no verifier — the only durable termination is a control the author cannot write to.
17. **Route tables generated from one authority (A4, amendment 3):** survives an opposite proxy, an
    opposite bundler, and an opposite deployment topology, because a set of prefixes maintained by
    hand in N places under a comment fails in exactly the environment the author did not run, in
    every topology — and the environments differ precisely so that one of them is the one nobody ran.
18. **Authorisation asserted at composition, not at the proxy (A8, amendment 3):** survives an
    opposite edge stack, because a matcher on `Accept: */*text/html*` shows that an exposure can be
    invisible to every check a person performs in a browser while fully open to a script, in any
    proxy language — so the only placement that holds is one the request cannot route around.
19. **Reachability is measured, never assumed (A8, method):** survives an opposite deployment, because
    "not verified" and "not exposed" are different propositions in every topology, and this audit
    converted one into the other for nine lanes and one synthesis draft before measurement corrected
    it.

---

# 7. Where I disagree with the lanes

Eight, ordered by consequence. The lanes measured well; these are weighing errors, with three
measurement errors among them — and one of the three is mine.

## D1 — `delivery` F14 is wrong on 4 of its 10 items and misattributes 3 more. The real number is 3.

The lane states of 10 operations present in OpenAPI and the router that *"none of them has any
`getJson`/`postJson`/`putJson` call anywhere in `packages/sdk-runtime/src`"* (`lanes/delivery.md:77-78`).

Measured directly:

```
packages/sdk-runtime/src/market.ts:179:  return postJson<MarketPriceIntelCollectionResponse>(`/market/collections${query}`, body);
packages/sdk-runtime/src/market.ts:182:  return getJson<MarketPriceIntelSignal[]>(`/market/signals${idsQuery("listing_ids", listingIds)}`);
packages/sdk-runtime/src/market.ts:185:  return getJson<MarketPriceIntelAggregate[]>(`/market/aggregates${idsQuery("codprod", codprods)}`);
packages/sdk-runtime/src/market.ts:190:  return getJson<MarketPriceIntelVerdict[]>(`/market/verdicts${idsQuery("codprod", codprods)}`);
```

Four of the ten are implemented, in `src`, exactly as the lane's own predicate specifies.

Three more are **OAuth browser-redirect endpoints**: `/integrations/auth/callback` and
`/connectors/melhor-envio/auth/{start,callback}`
(`contracts/api/marketplace-central.openapi.yaml:1207-1213, 2829-2835, 2849-2855` — `302` responses,
`code` query parameters). An SDK method that `fetch`es a redirect target would be a bug, not a gap.

**The genuine gap is 3 operations** — `/admin/fee-schedules/{seed,sync}` and `/pricing/tariff-defaults`
GET and PUT — all operator/admin surfaces, none user-facing.

Two things make this worth stating loudly rather than quietly correcting:

- **The refuting fact was in another lane the whole time.** `lanes/frontend.md:39` (F13) names
  `market.ts` as one of four additional hand-written SDK modules. Nobody reconciled the two lanes.
- **The mechanism of the error is the audit's own thesis.** The join almost certainly ran against
  `index.ts` while the prose claimed `src`. A careful, dedicated, single-purpose manual
  cross-reference of two hand-maintained documents got 4 of 10 wrong. That is the strongest available
  argument for A4 being mechanical, and it is stronger than any of the arguments the duplication lane
  offered.

## D2 — `frontend` F3 is right about the defect and wrong about the consequence

The lane claims a *"real runtime `TypeError: onRetry is not a function` the moment the user clicks"*
which, with no error boundary (F4), *"takes the whole SPA down"* (`lanes/frontend.md:45-56`).

Measured. `packages/ui/src/ErrorState.tsx` is 17 lines; it renders `<Button onClick={onRetry}>`.
`packages/ui/src/Button.tsx` spreads `{...props}` onto a plain `<button>` with no wrapper logic.
React DOM ignores an event handler prop that is `undefined` or `null` — it throws only for a
non-function, non-null value. So `onClick={undefined}` produces **an inert button, not an exception.**

Both call sites confirmed: `apps/web/src/pages/mutations/MutationPreviewModal.tsx:210` and
`apps/web/src/pages/mutations/MutationResultSummary.tsx:22`, both omitting `onRetry`.

The defect is real and worth fixing — an error screen whose only affordance silently does nothing is
arguably worse than one with no button, because the user retries and concludes the product is
broken. But it is a **UX dead end, not an availability incident**, and F4 is not its amplifier.
This was the frontend lane's heaviest finding and it does not justify emergency sequencing; it is
ordinary A5 work.

**Confirming test:** render `MutationResultSummary` with `items.isError`, click "Tentar novamente",
assert no exception is thrown and no refetch is issued.

**What the finding actually proves, and this is the valuable part:** a required prop was omitted at
two production call sites and shipped, because `tsc --noEmit` is wired to nothing (F2). **F3 is not a
frontend finding. It is the live proof of A2.**

## D3 — `duplication` DUP-1 undercounts its own sample

The lane reports two live `ListingReadModel` divergences. There is a third in the same struct:
`listings/domain/read_model.go:122-152` declares `PublishedQuantity *int` with no `omitempty` — the
key is always present and may be `null` — while `packages/sdk-runtime/src/index.ts:380-405` declares
`published_quantity: number`, neither optional nor nullable. A consumer trusting the SDK type does
arithmetic on `null`.

I am not claiming the lane was careless. I am claiming the opposite: **a careful hand-comparison of
one struct out of 172 found 2 of 3.** That miss rate is the correct way to size A4, and it is a
better argument than the raw line counts.

## D4 — `testing` T-4(a) is misclassified, and the misclassification points at a remediation that closes nothing

The lane files the catalog RLS test under "no independent evidence," because the test was authored in
the same commit as the fix (`7e3dcc47`). Same-commit authorship is a real weakness and the lane is
right that it exists.

But under the brief's control-versus-effect rule, the correct classification is different and much
stronger: **the control is inert, so the axis must be sized as if RLS were absent** — which is exactly
what `lanes/security.md` SEC-4 and `lanes/persistence.md` P-1/P-3 independently established. The app
DSN connects as the table owner and bypasses RLS; the GUC is set at one site.

Same-commit authorship is second-order here. **Even a fully independent test would prove nothing
about the running application**, because the test's own comment concedes it cannot reach the
production connection role. Filing this as an evidence-quality problem invites the remediation "get an
independent test written," which would close the finding and change nothing about the system. The
correct remediation is the boot assertion in §4.2 — refuse to start if the connecting role can bypass
RLS.

This applies to T-4(a) only. T-4(b), the catalogingest guard at `47a76837`, is a genuine
same-commit-evidence finding and the lane's classification of it is right.

## D5 — no lane read `arch-gate.sh` closely enough, and two details change the recommendation

Neither the cicd nor the layering lane reports that `arch-gate.sh` step 1 runs
`gofmt -l "$SERVER/internal"` — **so `cmd/`, `migrations/`, and `tests/` are gofmt-checked by
nothing** — nor that step 5 fails on any dirty working tree
(`git status --porcelain --untracked-files=all`, with the message *"a gate cannot certify a tree it
did not see"*).

The second is the consequential one. It means **`arch-gate.sh` is structurally incapable of being a
pre-commit or pre-push gate**: it can only certify an already-committed clean tree. The obvious next
step everyone would take — "wire arch-gate into `package.json`" — produces a gate that fails on every
real invocation and gets disabled within a week.

**My recommendation therefore diverges: retire `arch-gate.sh` rather than wire it.** Its five steps
redistribute cleanly — gofmt (module-wide), vet, build, test, and archscan all become steps in
`verify.yml`, where the checkout is always clean by construction and the CRLF problem does not exist.
The clean-tree step is simply deleted: a CI checkout is always clean, and a local gate must never
require it. That is a crisper outcome than fixing a script whose central assumption is wrong for its
intended use.

## D6 — `duplication` DUP-6 should be refused, not deferred

82 hand-written `for rows.Next()` loops with 0 `pgx.CollectRows` is presented as duplication to
remove. It is not. `rows.Scan(&a, &b)` is checked by the compiler for arity and type.
`pgx.RowToStructByName` is checked at **runtime**, against **hand-written struct tags** that must
match column names.

Adopting it would replace a compiler-checked verbosity with a new hand-synchronised seam — which is
the identical defect class as A4's four-copy contract and A6's five money types, the two largest
axes in this document. It would be internally consistent, well-named, reviewed as an improvement, and
wrong. That is precisely the machine failure mode the brief warns about, and it is worth naming as a
near-miss.

**Verdict: leave forever.** The boilerplate has produced zero reported defects across 82 sites.

## D7 — SEC-6 (CORS `*` with `Authorization`) is misclassed as idiom and carries a sequencing rule nobody stated

`Access-Control-Allow-Origin: *` on ~104 routes (`httpx/router.go:20`) is currently harmless in the
strict sense: there are no credentials to steal because there is no authentication (SEC-1). Browsers
also refuse to send credentials to a wildcard origin. The security lane's placement reflects that.

But its risk is **conditional on A8**, and the condition inverts it. The moment authentication exists,
a wildcard origin with `Authorization` in Allow-Headers becomes a live cross-origin credential
surface. So the finding is not "an idiom to tidy" — it is a **sequencing constraint on A8**: CORS must
be narrowed **in the same change** as any authentication work, never as a follow-up, and that change
needs its own negative fixture (a request from a disallowed origin must be rejected).

I have placed it inside A8 for this reason, not in the noise pile where its current standalone
severity would put it. **The classification error is filing a conditional hazard by its severity
today rather than by its severity in the state the program is deliberately moving toward.**

## D8 — `cicd` F2 measured a credential, not a repository, and its conclusion is void

The lane reported that `origin` could not be resolved from this machine and concluded that branch
protection and ruleset state were unknowable. That conclusion propagated into PHASE-0, into my
pre-amendment draft's first axis, and very nearly into a program plan whose opening move was
"re-establish a lost remote."

Measured (operator, amendment 2): two GitHub identities share the keyring. `gh` authenticates as
`leandrotcawork`; git authenticates as `developmentconexus-ops`. **`origin/main` = `7df7d011`,
2026-08-04. Local `main` is 98 ahead, 0 behind, fast-forwardable, nothing diverged.** The remote was
never stale and never unreachable. The `legacy` remote (`leandrotcawork/marketplace-central`, default
branch `master`) is a genuinely superseded artifact and is the likely source of the confusion.

This is not a lane being careless — the tool genuinely returned "cannot resolve." It is worth calling
out because of what the failure *pattern* is: **an instrument reported the state of its own
authentication and the report was read as the state of the world.** That is the same shape as
`GOV_CONTEXT_UNREGISTERED` passing because its guard does not exist (A3), and the same shape as the
`Accept: */*text/html*` matcher in amendment 3, where a browser check reports "SPA" and a script gets
the data. Three instances, three instruments, one class. A negative control — "authenticate as the
other identity and re-ask" — would have caught it in seconds, which is exactly the discipline §4
demands of every guard.

**I include my own instance of this class rather than only the lanes'.** My step-0 mitigation
("bind to `127.0.0.1`") was withdrawn under amendment 3 because it addressed a published container
port while the actual exposure runs through Caddy and ngrok over the compose network. It would have
been written down as done and been inert. See A8.

---

## Appendix — contradictions of the memory record, stated loudly

Two lanes contradict standing memory. Both lanes are right and more recent; recording so nobody
re-audits them:

- **ML token-refresh failure is now durable and tested** (`auth_flow_service.go:454-574`,
  `integrations_refresh_failure_test.go`). The record that it was invisible at every layer is
  superseded.
- **The DIFAL override is real, wired, and audited**, and the "three copies of the tariff table"
  count is stale — three different facts with three different owners
  (`lanes/duplication.md`).

And one where a lane's number is better than a document's: **ADR-023's prose says 35 violations;
measurement says 234** (fact 9, D-55). The document is wrong, has been wrong through two corrections,
and is inside A7's scope to fix.

**Two of PHASE-0's own facts were corrected mid-synthesis by operator measurement**, both loudly:
the remote is reachable and 0 behind (§7 D8), and the unauthenticated PII surface is
internet-reachable by design through Caddy in production and ngrok in development (§1 A8). The second
is the more serious, because nine lanes and this synthesis all had "reachability unverified" in front
of them and none of us treated it as the open question it was labelled.

---

*Arm A. Nine axes, ~33–46 agent-days plus an unsized authentication decision, one refused finding,
nineteen inversion tests, eight disagreements — one of them with myself.*
