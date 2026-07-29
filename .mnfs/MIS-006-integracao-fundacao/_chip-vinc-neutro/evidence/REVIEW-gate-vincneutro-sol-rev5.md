> **PROVENANCE.** Captured verbatim by `codex exec -o` into `agent__r6sol.last.md` and copied byte-for-byte
> below. Nothing after this banner was written, reordered, or summarised by this chip. Seat:
> `gpt-5.6-sol` / `medium`, OS-process, stdin closed, `--sandbox read-only`. Dispatch tip `2e9e9ce3`.
> Round 6, GPT side. Blind to the Opus seat.

---

REFUTED

The code range and pack are consistent with the frozen-review invariant:

- `git log --oneline 4c000a04..7b5c18eb`: four commits.
- `git diff --stat 7b5c18eb..HEAD -- apps/`: empty.

## Findings

### F1 — the “throwing constructor” accepts generator-impossible candidates

Severity: high, verification integrity.

[wireFixtures.ts:202](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/wireFixtures.ts:202) claims:

> “A candidate the backend can actually emit. Throws if it is not one.”

But `assertProducibleScore` validates only `(confidence, confidence_band, match_status)`, plus the `NO_CANDIDATE` state/input special case. It does not validate the state/input/reasons attached to the other score paths.

The Go authority at `generation_service.go:489-495` constructs every 95/ALTA/ACCEPT candidate as:

> `newCandidate(... LinkCandidateStateExactSKU, LinkCandidateMatchInputSellerSKU ...)`

and seeds both:

> `Anchor: "seller_sku"`  
> `Anchor: "ean"`

Nevertheless the constructor accepts:

- [QueueTab.test.tsx:553](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/QueueTab.test.tsx:553): `"state: \"exact_ean\""`
- [QueueTab.test.tsx:554](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/QueueTab.test.tsx:554): `"match_status: \"ACCEPT\""`
- [QueueTab.test.tsx:555](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/QueueTab.test.tsx:555): `"match_input: \"ean\""`

The same impossible combination occurs at `QueueTab.test.tsx:701-703`. Other accepted fixtures, including `:56-66`, `:71-82`, `:85-95`, and `:593-605`, carry reason sets the named scoring path cannot finalize.

What breaks: no immediate production component crashes because this is test infrastructure. What breaks is the claimed verification barrier: tests can remain green while asserting UI behavior from rows the running server cannot emit. Such fixtures must either be rejected or explicitly use `driftCandidate`.

### F2 — undeclared providers pass as producible wire candidates

Severity: high, verification integrity.

[wireFixtures.ts:142](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/wireFixtures.ts:142) only validates provider-specific finalization when:

> `if (providerCode === "mercado_livre")`

All other provider codes pass. Go authority says generation first resolves a registered capability declaration and aborts on failure (`generation_service.go:149-169`). The tree has only the `mercado_livre` capability declaration.

Seven fixtures nevertheless pass through `wireCandidate` as if producible:

- `QueueTab.test.tsx:498`, `"provider_code: \"shopee\""`
- `QueueTab.test.tsx:509`, `"provider_code: \"amazon_marketplace\""`
- `QueueTab.test.tsx:516`, `"provider_code: \"amazon-marketplace\""`
- `QueueTab.test.tsx:925-928`, `amazon_marketplace`, `amazon__marketplace`, `_amazon`, and `amazon`

This directly contradicts the helper’s own `driftCandidate` documentation for “NO DECLARATION HERE.”

What breaks: again, no immediate production runtime path. The provider-display tests are misclassified as real-wire fixtures and silently inherit Mercado Livre’s `marca` reason. Provider-declaration or finalizer regressions can therefore be masked.

### F3 — code and pack overstate what the mechanism proves

Severity: medium, evidence correctness.

The following strings are false given F1 and F2:

- [QueueTab.test.tsx:25](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/QueueTab.test.tsx:25):

  > “Every fixture below is built through `wireCandidate`, which THROWS on a candidate the generator cannot emit”

- [wireFixtures.ts:24](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/gifted-dhawan-5049f6/apps/web/src/pages/vinculos/wireFixtures.ts:24):

  > “UNWRITEABLE rather than detectable”

- `EVIDENCE.md:871`:

  > “an impossible candidate is now UNWRITEABLE rather than detectable”

The constructor does throw for the rules it actually implements, but it does not establish the broader quoted fact.

What breaks: the evidence pack certifies a class closure that the implementation does not provide.

## SWEEP

Counts are unique matching line sites. Population is the loose token search; extraction is the semantically relevant subset judged. Every extraction pattern matched known members; no zero result is being reported as clean.

| Class | Population | Extraction | Result |
|---|---:|---:|---|
| C1 | 57 | 29 | The 29 union-sensitive comparisons/filters were judged. `BatchPreviewModal.tsx:44,46,94` enumerate `OK/FAILED`; this is a known future-union drift exposure, not a present defect because the SDK currently has exactly those two members. QueueRow/Drawer status and direction sites either ask a single policy question or have explicit unknown fallbacks. No additional present defect. |
| C2 | 187 | 22 | The 165 generic `expect(` lines were excluded with that declared reason. The 22 containment/include/match sites were judged. Guard extraction at `wireFixtures.guard.test.ts:152,176` uses the exact `GO_SEAM.extract` regex, not a wider substring sentinel. UI `.toContain` assertions intentionally inspect rendered substrings/classes. Clean. |
| C3 | 68 | 68 | All extraction, mapping, and length sites were judged. The Go-source guard pairs readable-source checks with `> 0` for constants, extracted anchors, and declared anchors (`:131,162,186`). Empty UI collections are handled as product states, not accepted as proof of extraction. Clean. |
| C4 | 99 | 36 | Loose population includes 63 batch payload/type/reference fields. Extraction found 36 candidate fixture occurrences: 29 routed through `wireCandidate`, seven through `driftCandidate`. Syntactic routing is complete, but at least the 13 sites named under F1/F2 remain generator-impossible while passing the throwing constructor. Not clean. |
| C5 | 2,418 | 471 | Population includes source plus the entire pack, including captured logs. Extraction retained 41 source-comment claim lines and 430 pack Markdown anchor lines; logs and command-output quotations explain the difference. SHA facts `bcab8269` and `45b887b3`, the `main` deletion, and the four-commit range were measured and agree. The “UNWRITEABLE” claims at `QueueTab.test.tsx:25`, `wireFixtures.ts:24,202`, and `EVIDENCE.md:871` are false as stated. |
| C6 | 130 | 16 | The remaining 114 lines were arrays, destructuring, CSS brackets, or indexing by locally constrained keys. QueueRow’s band, direction, status, and ranking lookups have explicit fallbacks. `ImportacaoSection.tsx:36-37` remains a future wire-drift exposure through raw `statusClasses[status]`/`statusLabels[status]`, but its current SDK union is exhaustively mapped and the component is deleted on the measured `main`; no present defect added to this chip. |

Reconciliation reasons are therefore explicit:

- C1: `57 = 29 union-sensitive + 28 ordinary boolean/comment sites`.
- C2: `187 = 22 string-operation sites + 165 generic expectations`.
- C3: `68 = 68`.
- C4: `99 = 36 candidate fixtures + 63 payload/reference/type fields`.
- C5: `2,418 = 471 claim-capable source/Markdown lines + 1,947 logs/output/quotations`.
- C6: `130 = 16 wire/map-index sites + 114 unrelated bracket sites`.

## Pack custody

Exact first-command output:

```text
40
```

The second command’s stdout was empty. Git emitted this environment warning on stderr:

```text
warning: unable to access 'C:\Users\leandro.theodoro/.config/git/ignore': Permission denied
warning: unable to access 'C:\Users\leandro.theodoro/.config/git/ignore': Permission denied
```

Thus the pack currently has 40 tracked files and no porcelain changes.

## Adversarial questions

### Q1

The `Vincular` label and presence do not distinguish `REVIEW` from `ACCEPT`. `QueueRow.tsx:449` derives `noCandidate` solely from `match_status === "NO_CANDIDATE"`; every other status takes the branch containing `Vincular`.

Changing the fixture from REVIEW to ACCEPT therefore preserves the test’s explicit tab, button, and `MLB1` assertions. It also preserves the unasserted high-confidence KPI value because confidence remains 95/ALTA.

It does not preserve incidental REVIEW-path coverage: “Identificado por” changes from empty/manual to `CODPROD + EAN`. The test never asserted that behavior or any KPI despite its name, so it still proves exactly its old explicit assertions, but no longer exercises what the former REVIEW fixture incidentally exercised.

### Q2

None.

`wireFixtures.guard.test.ts` checks:

- the identity-anchor vocabulary;
- Mercado Livre’s supplied-anchor declaration;
- positive, non-vacuous extraction.

It does not read `appendProviderDeclaredUnavailableReasons` or `classifyProviderIdentityAnchor`. If finalizer closure diverges while the anchor vocabulary and capability declaration remain unchanged, no guard test necessarily fails. That is a named verification gap.

## Execution routing

No unresolved claim requires executor-seat execution. The refutation follows from the committed TypeScript, Go authority, SDK types, git facts, and pack custody.