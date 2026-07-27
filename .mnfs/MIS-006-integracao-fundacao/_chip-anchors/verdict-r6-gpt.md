# ROUND-6 GPT DERIVATION VERDICT — verbatim transcription

**What this file is.** The derivation gate ran as an OS-process `codex exec` dispatch (ledger row
D29), which does leave a durable artifact: the CLI's own `-o` last-message file. That artifact
lives in a session-scoped scratchpad outside the repository, so it dies with the session exactly
as the cold gate's transcript does. The transcription is therefore made for the same reason D30
gives for the cold side, and this file is the pack's copy.

**What it is not.** The reviewer did not write this file; the chip did. The reviewer's own artifact
was `agent__p6-sol-r6.last.md`, produced by the codex CLI at dispatch time and read once. As with
the cold verdict, this transcription is not verifiable against the original once the session
compacts. What it removes is unilateral control: the same text also travels to the hub verbatim in
the BLOCKED event, so two holders keep two copies with two timestamps.

**Its area.** Derivation and the apparatus. It does not judge F-01/F-02/F-03 semantics, and it says
so itself under `## Honest limit`. The cold gate's verdict for round 6, covering that half, is in
[verdict-r6-cold.md](verdict-r6-cold.md) and reads PASS. **The gates split.** No AGREEMENT string
was authored by the chip, and none is recorded anywhere in this pack.

Everything below the fence is the reviewer's text, unedited. The fence is five backticks because
the verdict quotes fenced blocks of its own.

`````quoted-verdict
# VERDICT: FAIL

Multiple BLOCKING apparatus defects stand. Most frozen-SHA counts now rederive correctly, including the corrected `mercado_livre` counts, diff statistics, and four complete listings. However, the tool silently omits named anchors, misses mutable commands, and its claimed coordinate SUPERSET is false. One frozen-SHA grep claim and one slice-range pin are also defective.

## Per-criterion result

| Criterion | Result | Derivation-side basis |
|---|---|---|
| C1 | NOT-CHECKED | Behaviour, port correctness, and composition semantics belong to the cold semantic gate. The coordinate apparatus itself fails to resolve some real named symbols, but I did not infer C1’s behaviour from that. |
| C2 | PASS | At `85b6c367`: `mandatoryUnavailableReasons` = 0 tracked Go hits; quoted `"marca"`/`"refforn"` = 0 production hits. |
| C3 | PASS | Complete 34-file listing matches `git ls-tree`; provider-branch grep returns zero. |
| C4 | NOT-CHECKED | Requires semantic inspection of the two detail forms and precedence. |
| C5 | PASS | 14-path contract/SDK witness is complete; frozen diff is empty; control diff is non-empty. |
| C6 | NOT-CHECKED | Test meaning and mutation adequacy belong to the semantic gate. |
| C7 | NOT-CHECKED | Hard-negative behavioural coverage was not reviewed here. |
| C8 | NOT-CHECKED | Independent-limit behaviour and tests were not semantically reviewed. |
| C9 | PASS | Migration tree contains 75 files at the cited SHA; migration diffs are empty; product-links control reports nine changed files. |
| C10 | PASS | Frozen code-path list contains exactly 15 paths, matches the paste, and contains zero `apps/web/` paths. |
| C11 | NOT-CHECKED | Go lanes and governance execution were explicitly forbidden. The committed hub artifacts contain the figures the pack copies, but raw runs were not reproducible here. |
| C12 | PASS | Whitespace-insensitive frozen diff is exactly `3 0`; existence and non-empty controls reproduce. |

## Numeric-claim ledger

| Claim | Class | Re-derivation |
|---|---|---|
| Contract authored one commit after base | git-at-frozen-SHA | `git rev-list --count 917f7bb5..f81b8975` → `1`; `f81b8975^` → full base SHA. Both chip files exist there and the validation contract is absent at `ac72eb82`. |
| Swept `product_links` files = 34 | git-at-frozen-SHA | `git ls-tree -r --name-only 85b6c367 -- …/product_links` → 34. Programmatic comparison: `pasted=34 live=34 exact=True`. |
| Contract/SDK witness = 14 | git-at-frozen-SHA | `git ls-tree … contracts/api packages/sdk-runtime` → 14; paste exact. |
| Code-path list = 15 | git-at-frozen-SHA | `git diff --name-only 917f7bb5 85b6c367 -- apps/ contracts/ packages/` → 15; paste exact; web paths = 0. |
| Governance synchrony differential = 5 paths | git-at-frozen-SHA | `git diff --name-only 8e37958a 93c90330 -- apps/server_core` → 5; paste exact. |
| Old `08308afb` `mercado_livre` count = 82 | git-at-frozen-SHA | Per-file counts sum to 82. The correction from the erroneous 86 is right. |
| Final `85b6c367` `mercado_livre` count = 84 | git-at-frozen-SHA | Per-file counts `13+12+4+2+24+2+18+9 = 84`; production count = 0. |
| Provider comparisons = 0 | git-at-frozen-SHA | Frozen `git grep -nE` exits 1 with zero rows. |
| `mandatoryUnavailableReasons` = 0 | git-at-frozen-SHA | Frozen grep exits 1 with zero rows. |
| Final server diff = 15 files, +1356/−167 | git-at-frozen-SHA | Exact `git diff --shortstat 917f7bb5 85b6c367 -- apps/server_core` output matches. |
| Historical code-tip diff = 17 files, +2817/−167 | git-at-frozen-SHA | Exact `git diff --shortstat 917f7bb5 2921d563` output matches. |
| Final integration-test diff = +17/−2 | git-at-frozen-SHA | `git diff --numstat 917f7bb5 85b6c367 -- generation_integration_test.go` → `17 2`; decayed insertion count was correctly removed from the proof row. |
| Migration witness = 75 | git-at-frozen-SHA | `git ls-tree … ac72eb82 -- migrations` → 75. |
| Product-links migration control = nine files | git-at-frozen-SHA | All nine pasted `--numstat` rows match exactly; migration diff is empty. |
| C12 controls = generation service `174 63`, root `3 0` | git-at-frozen-SHA | Both reproduce at `ac72eb82`; root remains `3 0` at `85b6c367`. |
| S1 write set = eight paths at `5bc55219` | git-at-frozen-SHA | Correct slice commit and eight-path list. |
| S1+S2 write set = 13 paths at `633bf9fa` | git-at-frozen-SHA | Count is 13, but the endpoint is wrong: the range also includes S3 commit `11e68f6f` and the later format commit. Correct S1+S2 endpoint is `b9da6d2e`, also yielding 13 paths. |
| S4 commit = two files | git-at-frozen-SHA | `030fa58c` reports exactly the two resolution-service files. |
| S5 = two files, +92/−31 | git-at-frozen-SHA | `git show --stat f92ca9c7` matches. |
| S6 = one file, +41/−0 | git-at-frozen-SHA | `git show --numstat 1df627f8` matches; production-file range is empty. |
| S9 production = +20/−13 | git-at-frozen-SHA | `2921d563` production file is `20 13`; its test file separately adds 67 lines. |
| Governance 53 violations / 175 lines / 14 exceptions | COPIED | Matches committed pointer `main@7c54bef`; raw scratchpad outputs are unavailable here and the lane was forbidden. |
| Live 34 listings / 38 generated / 29 resolved and cited percentages | COPIED | Matches committed pointer `main@40623b57`; live drive was not rerun. |
| Coordinate totals | pointer | Prose does not copy them. `coordinates.txt` alone reports 49 resolved and 4 ambiguous. R-14 discipline is respected here. |
| PACK mutations | COPIED, independently replayed | Four required failures returned exit 1; wholesale CRLF returned 0; CRLF plus edit returned 1. |
| Foreign exemptions | pointer to lane output | Current strict run prints 4 `quoted-contract` lines and 10 `quoted-verdict` lines; direct fence inspection confirms exactly 4 and 10 interior lines. |

## Findings

### BLOCKING — Named anchors are silently discarded

The generator does not fail or report when an anchor is recognized but absent from its limited top-level declaration map:

```text
cite-table.py:282
if not locs:
    continue
```

Independent inventory of the current pack:

```text
recognized_unique=123
resolved_names=53
unresolved_names=70
ambiguous_names=4
```

Several omissions are real repository symbols, not incidental prose:

```text
PACK: `IdentityAnchorSellerSKU`, `hardNegativeDimensionPattern`, `TitleMatch`
TABLE_LOOKUPS: no matches, rg exit 1
CODE:
marketplace_capability.go:25 IdentityAnchorSellerSKU
generation_service.go:684 hardNegativeDimensionPattern
link_candidate.go:15 LinkCandidateStateTitleMatch
```

An injected `` `TestThisAnchorDefinitelyDoesNotExistR22` `` also leaves the generated table unchanged and exits 0. This directly falsifies the table’s claimed coverage of anchors named by the prose.

### BLOCKING — Mutable-axis ban is narrow and bypassable

Both attacks exit 0 with zero mutable-axis findings:

```text
`git diff 917f7bb5..main -- apps/`
`git rev-parse HEAD`
```

Output:

```text
EXIT=0
mutable-axis commands in the pack  0
```

The implementation only recognizes selected `git diff|grep|ls-tree|show|log` forms and treats the presence of any SHA anywhere on a line as sufficient. Bare `HEAD`, mutable named refs, `rev-parse`, and ordinary worktree `grep` commands evade it.

### BLOCKING — Coordinate SUPERSET claim is false

The recognizer is:

```python
re.compile(r"(?<![0-9]):(\d+)...")
```

Attack `` `123:456` ``:

```text
EXIT=0
coordinates leaked into prose 0
```

The negative lookbehind contradicts the script’s own statement that the preceding character is unrestricted. Additionally, any arbitrary `quoted-*` fence suppresses coordinate and anchor inspection:

```text
```quoted-fraud
Attack coordinate real.go:456
```
```

This also exits 0, although the printed contract says the coordinate exemption is specifically `quoted-output`.

### BLOCKING — Frozen-SHA `connectorsports` claim is false

[EVIDENCE.md](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-anchors/.mnfs/MIS-006-integracao-fundacao/_chip-anchors/EVIDENCE.md:406) says the application-package grep at `85b6c367` returns nothing.

Actual command:

```text
git grep -c connectorsports 85b6c367 -- apps/server_core/internal/modules/product_links/application

85b6c367:.../generation_integration_test.go:2
exit=0
```

Only the production-excluding form returns nothing. The later correction section acknowledges those two test-file occurrences, so the pack contradicts itself.

### BLOCKING — One write-set re-pin does not end at the described slice

[EVIDENCE.md](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-anchors/.mnfs/MIS-006-integracao-fundacao/_chip-anchors/EVIDENCE.md:143) describes `917f7bb5..633bf9fa` as S1+S2.

The range actually contains:

```text
5bc55219 S1
b9da6d2e S2
11e68f6f S3
633bf9fa later composition formatting
```

The 13-path result happens to equal the correct `917f7bb5..b9da6d2e` path set because S3 modified already-listed paths. The count is coincidentally right, but the range is not the slice it claims to evidence.

### SHOULD-FIX — Foreign fences are accurate today but trust an unaudited label

Current counts are correct:

```text
quoted-contract: 4 interior lines
quoted-verdict: 10 interior lines
```

Removing those tags exposes three mutable findings: two on the quoted contract command and one `git ls-files` verdict line. Their present use is legitimate quotation. Structurally, however, the pack author can hide authored content merely by applying a FOREIGN tag; the tool verifies neither provenance nor that only foreign text is fenced.

### PASS — Required mutation behaviour and non-writing strict mode

Real-file `--strict`:

```text
coordinate table CURRENT
coordinates leaked into prose 0
mutable-axis commands in the pack 0
EXIT=0
```

The six exact virtual-file mutations produced:

```text
coordinate-in-prose      EXIT=1
mutable-axis             EXIT=1
table-drift              EXIT=1
table-absent             EXIT=1
table-wholesale-crlf     EXIT=0
crlf-plus-hand-edit      EXIT=1
```

Because the environment enforced read-only access, these were executed against the exact committed script through an in-memory filesystem shim. Real hashes before and after remained identical. Final status contains only the pre-existing `package-lock.json` modification plus `verdict-r6-cold.md`, which appeared concurrently from the cold reviewer; I created neither.

## Honest limit

This review did not judge F-01/F-02/F-03 semantics, whether named tests assert the prose claims, production behaviour, live integration, Go build/vet/test, governance execution, or UI behaviour. Those belong to the cold semantic gate or were explicitly prohibited. I also did not report the `slices.SortFunc` stability suspicion or `refforn`, as directed.
`````
