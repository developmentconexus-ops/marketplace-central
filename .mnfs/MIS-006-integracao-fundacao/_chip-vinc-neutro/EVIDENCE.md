# CHIP-VINC-NEUTRO — EVIDENCE

- **Branch:** `chip/vinc-neutro` · **worktree:** `.claude/worktrees/gifted-dhawan-5049f6`
- **BASE-SHA (floor):** `5441fe18f64171ef61cb03b51b5bf66e2922e4eb`
- **`main` at close:** `bcab8269` · **HEAD:** `7a343fea`
- **Slice commits:**
  - `fa6ca3a2` — *feat(vinculos): render INCOMPARABLE, auto-approved badge, neutral vocabulary*
  - `7a343fea` — *feat(vinculos): render "Identificado por" from the anchors that decided*
    (the V6 fix, after the round-1 gate refuted my REPORT)
- **Scope:** `/vinculos` only, `apps/web` only. Zero Go, zero migrations, zero `contracts/`, zero
  `packages/sdk-runtime/`. No server booted, no `:8080`, no `.env*` read. No push.

Artifacts referenced below live under `_chip-vinc-neutro/evidence/`.

---

## V1 — `INCOMPARABLE` renders with its own glyph and its own tokens — **PASS**

Chosen glyph **`?`**; chosen token pair **`bg-info-soft` / `text-info`**.

Three direction-indexed maps, all four keys present, all four design tokens (zero literal
Tailwind, zero `undefined`):

- `QueueRow.tsx:135` — `directionGlyphs.INCOMPARABLE = "?"`, distinct from `UNAVAILABLE`'s `–`.
- `QueueRow.tsx` `directionClasses` — `INCOMPARABLE: "bg-info-soft text-info"`, exported.
- `VinculoDrawer.tsx` — the third map **no longer exists**. It was a duplicate local
  `directionClasses`, and duplication is precisely why a new direction can be added in one map
  and missed in the other. The drawer now imports the single exported map from `QueueRow`, so
  the file is exhaustive by construction rather than by review.

Semantics honoured per the D-122 frozen block: `info` is neither `warn` (AGAINST, "blocks") nor
`surface-2/faint` (UNAVAILABLE, "the provider never supplies this"). Asserted negatively in
`QueueTab.test.tsx:212-219` — the chip's className must NOT contain `bg-warn-soft`,
`bg-accent-soft` or `bg-surface-2`. Reusing either would have been FAIL and the test now fails
if a future edit does it.

## V2 — row 100 % `INCOMPARABLE` shows at least one motivo — **PASS** *(the deciding criterion)*

**The defect.** `QueueRow.tsx:159` built the collapsed cell as

```ts
const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")].slice(0, COMPACT_CHIP_LIMIT);
```

Enumeration by string literal: type-correct, therefore silent when D-B added a fourth direction.
A row whose motivos are all `INCOMPARABLE` fell through every branch — `shown` empty, `hidden > 0`
— and the cell rendered a lone `+2` with zero chips. A filter wearing a ranking's comment, against
the ADR-17 invariant the file documents at `:154-156`.

**The fix** — ranking that is total by construction, not by enumeration:

```ts
const directionRank: Record<ProductLinkReasonDirection, number> = {
  AGAINST: 0, FOR: 1, INCOMPARABLE: 2, UNAVAILABLE: 3,
};

const shown = reasons
  .map((reason, index) => ({ reason, index }))
  .sort((a, b) => directionRank[a.reason.direction] - directionRank[b.reason.direction] || a.index - b.index)
  .slice(0, COMPACT_CHIP_LIMIT)
  .map((entry) => entry.reason);
```

Every reason enters the sort, so `shown` is empty only when `reasons` is empty. A fifth direction
now fails to COMPILE (the `Record` is exhaustive) instead of vanishing from the cell — the class
of bug is closed, not just this instance.

**Proof is on the rendered DOM**, never on `shown`: `QueueTab.test.tsx:177-220` renders a candidate
whose two reasons are both `INCOMPARABLE` (anchors `seller_sku`/`provider` and `ean`/`erp`) and
asserts `within(row).getAllByTestId("motivo-chip")` — `length > 0`, `toHaveLength(2)`, and NO
`Mostrar todos os…` button.

## V3 — must-fail of V2 — **PASS**

`evidence/V3-must-fail.txt`, captured verbatim. Reverted `:159` to the three-direction
enumeration and re-ran the file:

```
❯ src/pages/vinculos/QueueTab.test.tsx (10 tests | 1 failed | 9 skipped) 62ms
   × QueueTab > keeps a motivo on screen for a row whose reasons are ALL INCOMPARABLE (ADR-17) 61ms
     → Unable to find an element by: [data-testid="motivo-chip"]
```

The motivo cell's DOM in the failing run, quoted from the artifact — **how many chips appeared and
what the button said**: the `<ul class="flex min-w-0 items-center gap-1" />` is EMPTY (zero chips),
next to `<button aria-label="Mostrar todos os 2 motivos" …>+2</button>`. Exactly the screen the
contract forbids. `Tests 1 failed | 9 skipped (10)`. Fix restored, file green.

## V4 — the `side` of `INCOMPARABLE` reaches the operator — **PASS**

`side` is read from the **field**, never parsed out of the Portuguese `detail`:

```ts
const incomparableSideLabels: Record<ProductLinkReasonSide, string> = {
  provider: "falta no anúncio", erp: "falta no ERP", both: "falta nos dois lados",
};
export function reasonSideLabel(reason: ProductLinkReason): string | undefined {
  if (reason.direction !== "INCOMPARABLE") return undefined;
  return reason.side ? incomparableSideLabels[reason.side] : undefined;
}
```

Rendered inline in the compact chip (`QueueRow.tsx:156`) and in the drawer pill via the shared
`reasonChipLabel`. Asserted at `QueueTab.test.tsx:206-207`:
`"? SKU (falta no anúncio)"` and `"? EAN (falta no ERP)"`.

**The side-less path is named, not invented.** `generation_service.go:711` returns
`(DirectionIncomparable, "", "não foi possível comparar a âncora %s", true)` — an `INCOMPARABLE`
with an EMPTY side. The FE fabricates nothing there: `reasonSideLabel` returns `undefined` and the
chip renders the anchor alone. Covered by a dedicated test (chip reads `? Marca`, and the row's
`textContent` does not match `/falta n/`). The other four returns at `:715/:723/:726/:728` do carry
`provider` / `both` / `provider` / `erp`.

## V5 — `tsc` of the write-set closes, baseline declared — **PASS**

`npx tsc -p apps/web --noEmit`, run from the **main repo root** (a worktree-local
`npx --no-install tsc` passes vacuously).

- **Before:** 15 errors — `evidence/L0-tsc-baseline.txt`.
- **After:** 12 errors — `evidence/L0-tsc-after.txt`. **Zero in `src/pages/vinculos/`.**

The 3 that went away were mine, all the same shape (`Record<ProductLinkReasonDirection, …>`
missing the `INCOMPARABLE` key): `QueueRow.tsx:34`, `QueueRow.tsx:75`, `VinculoDrawer.tsx:118`.

The **12 baseline errors, one by one by path**, all pre-existing, all outside the write-set, all
UNTOUCHED by this chip:

| # | Path | Error |
|---|---|---|
| 1 | `src/pages/anunciosQueries.ts(17,54)` | TS2345 `ListingListOptions` not assignable to `Record<string, unknown>` |
| 2 | `src/pages/anunciosQueryState.test.ts(128,39)` | TS2345 mock missing `refreshListings`, `listIntegrationOperationRuns` |
| 3 | `src/pages/anunciosQueryState.test.ts(153,42)` | TS2345 same |
| 4 | `src/pages/AnunciosTable.test.tsx(40,7)` | TS2739 `ListingMarketSignal` missing `median`, `min_valid`, `max_valid` |
| 5 | `src/pages/ListingsRefreshControl.test.tsx(112,36)` | TS2352 `Promise<void>` → `QueryClient` |
| 6 | `src/pages/mutations/MutationPreviewModal.tsx(197,33)` | TS2741 `onRetry` missing on `ErrorStateProps` |
| 7 | `src/pages/mutations/MutationResultSummary.tsx(22,25)` | TS2741 same |
| 8 | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(39,44)` | TS2322 `"complete"` not a `CanonicalSourceFactQuality` |
| 9 | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(40,45)` | TS2322 same |
| 10 | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(41,46)` | TS2322 same |
| 11 | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(45,3)` | TS2322 `"MATCHED"` not a `MarketPriceIntelMatchStatus` |
| 12 | `src/pages/produto/ProdutoPage.test.tsx(91,7)` | TS2322 same |

`tsc` green was never the criterion. V2 is, and V2 is proven by a must-fail (V3), not by the
compiler — the compiler could not see this defect at all.

### Tree provenance of the L0 lane (addendum, after `main` moved to `81457e4f`)

`main@81457e4f` corrected this chip's own pack with a warning that lands exactly here: *"contar 15
erros de `tsc` com a composição esperada NÃO prova que você leu a árvore certa, porque os 3 erros de
`/vinculos` também existem em `main` desde `dbdcdfb1`. Ler a árvore errada dá a mesma contagem."*

The warning is right, and my two original artifacts do not answer it on their own — the baseline was
captured with `cwd` at an `apps/web`, so its paths are relative (`src/pages/…`) and name no tree.
So I re-ran the after-lane in a form that **cannot** read another tree, and persisted it:

- `evidence/L0-tsc-after-worktree-selfcontained.txt` — `npx --no-install tsc -p apps/web --noEmit`
  run from the **worktree root**, with the worktree's own `node_modules` (from `npm ci` here).
  **12 errors, 0 under `pages/vinculos/`.**

That run is non-vacuous *because* `npm ci` ran in this worktree: `node_modules/@marketplace-central/*`
resolve to `…/gifted-dhawan-5049f6/packages/*` and `…/apps/web` — verified by
`ls -l node_modules/@marketplace-central/`. Under the old junction approach the same command would
have been the vacuous pass the pack warns about.

**What is proven, and what is not.** The after-state (12, zero in `vinculos`) is now proven against
*this* tree by two independent invocations. The before-state (15) is a claim about `main`, and its
artifact carries no tree marker; it is corroborated by the pack's own statement that `main` has
carried 3 `/vinculos` errors since `dbdcdfb1`, and by V3's must-fail, which shows the three sites are
real code in this branch's parent. I am not upgrading that to "proven" — it is a corroborated claim.
Immaterial to the verdict either way: V5 asks that the write-set close and the baseline be declared,
and 0 errors under `pages/vinculos/` is a self-contained measurement of the write-set.

## V6 — "Identificado por" — **IMPLEMENTED, PASS** *(after the gate refuted my REPORT)*

**I filed this as a REPORT and I was wrong. The round-1 gate reviewer refuted it and the
refutation holds.** Recorded in full because the failure mode matters more than the outcome: I
declared a wire gap without reading the function that consumes the wire.

My premise was "the candidate contract has no field naming the anchor set that decided". False.
`decisionRuleForCandidate` (`resolution_service.go:812-835`) — the function that writes
`rule_matched` into the decision trail — is a **pure function of two fields that are already on the
wire**:

```go
switch candidate.MatchStatus {
case …Accept:   return DecisionRuleConcordantCodprodEAN
case …Confirm:  // falls through to the MatchInput switch
default:        return DecisionRuleManual
}
switch candidate.MatchInput {
case …SellerSKU: return DecisionRuleExactCodprodUnique
case …EAN:       return DecisionRuleExactEANUnique
default:         return DecisionRuleManual
}
```

`ProductLinkCandidate.match_status` and `.match_input` are both on the wire
(`packages/sdk-runtime/src/index.ts:1085` and `:1080`), and the wire enum **does** carry `CONFIRM`
(`:1058`) — which was the one thing that could have made the derivation partial. It does not. The
column is derivable today at zero backend cost, and it lands on the same rule the trail will
record, because it reads the same two inputs through the same switch.

**Implemented** in `QueueRow.tsx` as `decidingAnchors`, mirroring the backend exactly:

| `match_status` | `match_input` | backend rule | column |
|---|---|---|---|
| ACCEPT | (any) | `concordant_codprod_ean` | `CODPROD + EAN` |
| CONFIRM | `seller_sku` | `exact_codprod_unique` | `CODPROD` |
| CONFIRM | `ean` | `exact_ean_unique` | `EAN` |
| CONFIRM | `title` / `manual` / `none` | `manual` | `—` |
| REVIEW / REJECT / NO_CANDIDATE | (any) | `manual` | `—` |

**The ACCEPT row is verified, not assumed.** `MatchStatusAccept` has exactly ONE assignment site in
the whole module — `buildConcordantCandidate`, `generation_service.go:505` — and that function
builds its candidate from `seller_sku` FOR **plus** `ean` FOR corroborating the same codprod. So
ACCEPT ⟺ corroborated CODPROD+EAN; naming both anchors asserts nothing that did not happen.
(The backend's own ACCEPT branch likewise ignores `MatchInput`, so the FE ignoring it matches.)

**Both maps are keyed by the full union**, not switched on string literals — a sixth `match_status`
or a new `match_input` fails the compiler instead of falling silently through a `default`. That is
this chip's own V2 lesson applied to its own new code.

**The negative case the contract demands is present** (`QueueTab.test.tsx`): a `title_match` /
REVIEW candidate carrying `{anchor: "title", direction: "FOR"}` renders its motivo `✓ Título` and
names **no** anchor — the Identificado por cell is the honest `—`. Deriving the column from
`direction === "FOR"` would have printed "Título" as the deciding anchor for a candidate that title
explicitly did not decide (`generation_service.go:551`, "ranking-only, nunca ACCEPT"; D-121 routes
title-only to REVIEW). That is the A2-R2 failure mode, and it is now guarded by a test.

**Must-fail captured** at `evidence/V6-must-fail.txt` — against the pre-fix cell:
`× names only the anchors that decided… → Unable to find an element by:
[data-testid="identificado-por"]`, `Tests 1 failed | 10 skipped (11)`. Restored, green.

**Column placement — the one contestable call, declared.** D-122:136 says the column *"substitui
'SKU ML'/'GTIN'"*. I replaced **GTIN only** and kept the former "SKU ML" slot as **Canal**. Reason:
D-122's plural rests on the same false premise this chip already corrected — that cell never held a
SKU, it held `provider_code` (the marketplace slug). Deleting it would delete provider identity
from the table, which V10 explicitly forbids ("neutralizar o valor é FAIL"), and which D-122 cannot
have intended since it believed the cell held a redundant seller SKU. Replacing GTIN is the part of
D-122 that is unambiguously right: "✓ igual" was the reading of ONE anchor, and the set that
decided supersedes it. Column count and order unchanged. Flagged for the hub as **FINDING-4** — if
the hub reads D-122 literally, dropping Canal is a one-line change, but it then owes V10 an answer.

## V7 — auto-approved badge fires on `actor_type === "system"` — **PASS**

Predicate in `useVinculosResolved.ts`:

```ts
export function resolutionAuditEntry(item: ProductLinkWorkflowItem): ProductLinkAuditEntry | undefined {
  const resolving = (item.audit ?? [])
    .filter((entry) => entry.next_state === "resolved")
    .sort((a, b) => Date.parse(b.created_at) - Date.parse(a.created_at));
  return resolving[0];
}
export function isAutoResolved(item: ProductLinkWorkflowItem): boolean {
  return resolutionAuditEntry(item)?.actor?.actor_type === "system";
}
```

Three DOM cases, `ResolvidosTab.test.tsx`:

1. `:80` audit entry with `actor_type: "system"` → `getByTestId("auto-resolved-badge")` has text
   `Automático`, **and** the `Vinculado ✓` chip is still there (the badge adds a fact, it does not
   replace one).
2. `:101` `actor_type: "operator"` → badge absent.
3. `:121` pre-M-05 record, `audit: []` → badge absent, `row.textContent` does not contain
   `"undefined"`, and `Desfazer` is honestly **disabled** — there is no resolution to reverse.

Case 3 is the `milestone.md:213` shape. Absence of badge reads "not automatic", which is the truth
available; a fabricated badge would not be.

## V8 — the F-04 brief correction, in writing — **PASS (recorded)**

**The `milestone.md:208` predicate was NOT implemented as written.** It asks for
`rule_matched = exact_ean_unique AND actor = system`. Both halves fail against the repo:

- **(a) `rule_matched` is not on the wire.** `grep -rn "rule_matched" contracts/ packages/sdk-runtime/src/`
  → **zero hits, exit status 1**. It exists only as a DB column in
  `0082_product_link_decisions.sql`. No FE code can read it.
- **(b) that exact pair is FORBIDDEN by the schema.**
  `apps/server_core/migrations/0082_product_link_decisions.sql:53-54`:
  ```sql
  CONSTRAINT product_link_decisions_system_is_corroborated_check
      CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean'),
  ```
  `actor='system'` is legal ONLY with `concordant_codprod_ean`. `exact_ean_unique` + `system` can
  never exist. The brief predates D-121, which made corroboration the only thing the system may
  decide alone; a single anchor now goes to human CONFIRMATION.
- **(c) the predicate that DOES exist, at zero backend cost:** auto-approval writes its audit entry
  with `Actor: domain.ActorMetadata{ActorType: "system", ActorID: "auto_linker"}` —
  `resolution_service.go:280` — and `item.audit[].actor.actor_type` is already on the wire.

**Predicate used:** `actor_type === "system"` on the audit entry that moved the link to `resolved`.

Recorded here and in a docblock on `isAutoResolved` (`useVinculosResolved.ts`) plus the describe
block of `ResolvidosTab.test.tsx:70-78`, so the next reader hits the correction at the code, not
only in `.mnfs`. Without this record the next reader opens `milestone.md` and believes it.

**M-06 F-04 brief is stale on this point and should be corrected at source.** Filed as
**FINDING-1**.

## V9 — `refforn` in the vocabulary: decision declared — **PASS**

**KEPT.** `QueueRow.tsx:127` still maps `refforn: "Ref. forn."`.

D-A (D-122) removed `refforn` from the cross-side anchor vocabulary, so **no candidate generated
today carries it** — `knownIdentityAnchors` has no `refforn`. But D-A also decided that
already-persisted reasons are **NOT migrated**, so a decided candidate can still hold a `refforn`
motivo: audit data meant to survive verbatim. Dropping the entry would degrade that historical row
from `Ref. forn.` to the raw machine name, for no gain. The reasoning is recorded in the comment at
`QueueRow.tsx:116-121` so the next reader does not "clean it up".

The fallback promise still holds either way — `QueueRow.tsx:141-143`:

```ts
function anchorShortLabel(anchor: string): string {
  return anchorShortLabels[anchor] ?? anchor;
}
```

Unknown anchors fall through **verbatim**: never hidden, never renamed into something the wire did
not say.

## V10 — neutral vocabulary did not erase provider data — **PASS (both halves)**

**Half 1 — a structural label went neutral.** `QueueTab.tsx`: `Anúncio ML` → `Anúncio`,
`SKU ML` → `Canal`. `ResolvidosTab.tsx`: `Anúncio ML` → `Anúncio`. Column count and order
unchanged.

**Half 2 — a data value still says which provider the listing belongs to.** The old `SKU ML` cell
was doubly wrong: it was a mislabel AND it rendered `candidate.provider_code`, which the wire fills
with the marketplace **slug** — the candidate contract carries no seller SKU at all. It is now the
`Canal` cell rendering the provider by display name:

```ts
const providerDisplayNames: Record<string, string> = { mercado_livre: "Mercado Livre" };
export function providerDisplayName(providerCode: string): string {
  return providerDisplayNames[providerCode] ?? providerCode;
}
```

Unknown provider codes fall through verbatim rather than being blanked — an unknown provider is
still a fact. Absent `provider_code` renders `<UnknownValue />`, never a fabricated channel.

**No raw slug on screen.** `grep -n "mercado_livre" *.tsx *.ts` over the write-set, excluding
tests, returns **four hits and all four are prose or a map key**, none a rendered value:
`QueueRow.tsx:51` (comment), `QueueRow.tsx:56` (the display-name map KEY), `QueueRow.tsx:309`
(comment), `VinculoDrawer.tsx:180` (comment). Asserted in the DOM twice — `QueueTab.test.tsx`
(the whole `container.textContent` free of `"mercado_livre"`) and
`VinculosDesign.golden.test.tsx:104-105` (`getByText("Mercado Livre")` present,
`queryByText("mercado_livre")` absent). This is the trap that hit CHIP-PED-FILA on 4 surfaces.

The drawer subtitle uses the same function, so the slug is gone from the second surface too.

## V11 — no collateral damage — **PASS**

**vitest, counts before and after — both runs persisted:**

| | artifact | Files | Tests |
|---|---|---|---|
| before | `evidence/L1-vitest-baseline.txt` | 62 passed (62) | 511 passed (511) |
| after | `evidence/L1-vitest-after.txt` | **63 passed (63)** | **520 passed (520)** |

Delta exactly +1 file (`ResolvidosTab.test.tsx`) and +9 tests, all mine (8 in the first slice, 1 for
V6). Zero green→red. There is no red vitest baseline to declare — the FE test lane was fully green
before and after; the red baseline is `tsc`, declared in V5.

*The round-1 reviewer marked V11 NOT-PROVEN because only the AFTER run existed on disk — the before
count lived in my session transcript, which is not evidence. Correct call.* The baseline is now a
real artifact: I restored the write-set to `main` (`git checkout main -- apps/web/src/pages/vinculos/`,
plus moving the new `ResolvidosTab.test.tsx` aside since it does not exist at `main`), ran the full
suite, captured it, then restored the tree to HEAD and confirmed `git status` clean. Same technique
as the V3/V6 must-fails.

**Flake, declared:** one full-suite run reported
`src/app/AppRouter.test.tsx > renders the product workspace — Test timed out in 5000ms` (20821 ms
under load). Proven load-induced, not mine: the file re-run in isolation passed 22/22 in 495 ms,
and a subsequent full run came back 63/63 · 519/519 green. `AppRouter.test.tsx` is outside the
write-set and untouched.

**`VinculosDesign.golden.test.tsx` green, and CHANGED — declared.** Nothing was loosened:

- header list `["Anúncio ML", "SKU ML", …, "GTIN", …]` → `["Anúncio", "Canal", …,
  "Identificado por", …]` — forced by the F-05 renames; the golden is a design gate and these are
  the design changes it is gating.
- fixture `provider_code: "abc-1"` → `"mercado_livre"` — the old fixture used a fake code, which is
  why the golden never noticed the slug leak.
- assertion `getByText(<provider_code>)` → `getByText("Mercado Livre")` **plus**
  `queryByText("mercado_livre")` absent. Strictly stronger: fails both if the display name
  disappears and if the raw slug returns.
- `getByText("✓ igual")` → `getByTestId("identificado-por")` has text `EAN`; and in the
  NO_CANDIDATE case `queryByText("✓ igual")` absent → `queryByTestId("identificado-por")` absent.
  Equivalent strength on the new column.
- **fixture `match_status: "REVIEW"` → `"CONFIRM"`, which needs justifying rather than asserting.**
  The fixture is a single exact-EAN match on one product. Under D-121 that is precisely the
  CONFIRM queue — one anchor resolved a single product, nothing against it, nothing corroborating
  it, so it waits for a human yes (the sdk-runtime docblock at `index.ts:1049-1052` says exactly
  this). `REVIEW` means the anchors *disagree or are ambiguous*, which this fixture is not. The old
  value was wrong on its own terms, and it predates D-121. I am not claiming the change is
  costless: had I left it, the golden's EXEMPLO-IO row would legitimately show `—` in the new
  column, and the test would still pass. I changed it because a golden whose fixture misstates the
  candidate's own state is a weaker gate, not because it was failing.
- top-of-file comment updated to record all three renames and that the `SKU ML` column had always
  been rendering the slug.

**Scope.** `git diff --name-only main HEAD` against the REAL `main` (`bcab8269`, which moved past
the BASE-SHA floor while this chip ran) — 8 paths across both commits, all inside
`apps/web/src/pages/vinculos/`:

```
apps/web/src/pages/vinculos/QueueRow.tsx
apps/web/src/pages/vinculos/QueueTab.test.tsx
apps/web/src/pages/vinculos/QueueTab.tsx
apps/web/src/pages/vinculos/ResolvidosTab.test.tsx
apps/web/src/pages/vinculos/ResolvidosTab.tsx
apps/web/src/pages/vinculos/VinculoDrawer.tsx
apps/web/src/pages/vinculos/VinculosDesign.golden.test.tsx
apps/web/src/pages/vinculos/useVinculosResolved.ts
```

**Zero** `VinculosPage.tsx`. **Zero** `ImportacaoSection.tsx`. Both are CHIP-IMPORT-CHAIN's.
`package-lock.json` untouched.

Run against the BASE-SHA floor, the same diff additionally shows 8 `.mnfs/` paths — the hub's own
wave-2 commits `af61c5e8` and `bcab8269` (the three wave-2 packs, the M-06 milestone correction,
and `_hub-gate-anchors-2/EVIDENCE.md`). **Those are not mine**; they are what `main` gained while
this chip ran, which is exactly why the contract calls the base a floor.

**Re-checked at close, after `main` moved again to `81457e4f`.** The floor moved a third time. The
only honest way to state scope is to separate the two directions:

```
git diff --name-only bcab8269 HEAD   → 8 write-set paths + 11 of my own .mnfs/ evidence paths
git diff --name-only bcab8269 main   → 3 paths, all .mnfs/ pack files, none of them mine
comm -12 (mine) (main's)             → EMPTY
```

Zero intersection, so this branch merges into `81457e4f` without touching anything the hub changed
after `bcab8269`. `main@81457e4f` rewrote `_chip-vinc-neutro/chip.md` §L0 — replacing the junction
instruction with `npm ci` at the worktree root, plus the tree-provenance warning answered in the V5
addendum above. That correction ratifies what this chip had already done in the field and recorded
as FINDING-2; nothing in it changes a verdict, and no work is owed. The other two changed paths are
CHIP-IMPORT-CHAIN's pack and contract, which I do not read or touch.

---

## FINDINGS (for hub ratification — core §0)

**FINDING-1 — the M-06 F-04 brief is stale and self-contradictory against the schema.**
`milestone.md:208` specifies `rule_matched = exact_ean_unique AND actor = system`. `rule_matched`
is not on the wire (empty grep in `contracts/` and `packages/sdk-runtime/src/`) and that exact pair
is forbidden by `0082_product_link_decisions.sql:54`. The brief predates D-121. A chip that
implemented it literally would have had to invent a wire field to satisfy a predicate the database
rejects. Full detail in V8. **Recommend correcting `milestone.md:208` at source** — the correction
currently lives only here and in the code comments.

**FINDING-2 — the "junction to `node_modules`" shortcut silently validates the WRONG tree, and
profile §3 already says so.** I began this chip by junctioning the main checkout's `node_modules`
into the worktree plus a chip-local `vitest.chip.config.ts`. Inspection showed
`node_modules/@marketplace-central/{sdk-runtime,ui,web}` are **absolute** symlinks pointing back
into `Documents/marketplace-central/packages/…` — so my early lanes were type-checking and testing
MAIN's workspace packages, not my worktree's. That is verbatim the failure profile §3 names
("silently validating the wrong code"). Corrected: removed all 5 junctions, deleted
`vitest.chip.config.ts`, ran `npm ci` at the worktree root (157 packages, 49 s), verified the
symlinks now resolve INTO the worktree, and re-ran BOTH lanes with the stock configs. Every number
in this evidence is from the post-`npm ci` state. **Recommend the profile's chip bootstrap say
`npm ci` is mandatory rather than "junction is a shortcut"** — the shortcut has no honest failure
mode, it just quietly reports someone else's tree. (`npm ci` inside the worktree does not count as
a dependency change: no manifest edited, `package-lock.json` untouched.)

**FINDING-3 — profile §3 false-alarm F-ENV-9 names a mitigation binary that is not installed.**
F-ENV-9 says to dispatch codex "via the 0.145.0-alpha.18 binary explicitly, not bare `codex` from
PATH". This machine has only `0.144.0` and `0.144.4` under
`~/.codex/packages/standalone/releases/`; `where codex` resolves to the 0.144.4 shim and the
npm-global path does not exist. Dispatched via the **explicit** 0.144.4 binary path (which still
avoids the bare-PATH shim hazard). **Recommend F-ENV-9 be re-worded** to "dispatch via an explicit
versioned binary path" plus a note on how to obtain 0.145.0-alpha.18, since as written the
mitigation is not followable.

**FINDING-4 — D-122's "Identificado por substitui 'SKU ML'/'GTIN'" is half-founded on a false
premise, and I resolved it rather than escalating.** D-122:136 orders the column to replace two
cells. Replacing GTIN is right. Replacing "SKU ML" is founded on the belief that the cell held a
redundant seller SKU; it never did — it held `provider_code`, the marketplace slug (the candidate
contract has no seller SKU at all). Executing that half literally would delete provider identity
from the table, which V10 of this chip's own contract forbids in the same breath. I kept the slot
as **Canal** and replaced GTIN only, so both normative sources are satisfied and the column count
is unchanged. **The hub owns the tie-break** — if D-122 is to be read literally, dropping Canal is
a one-line change, but V10 then needs an answer for where provider identity lives on the queue
table. Recommend amending D-122:136 to name GTIN only.

**FINDING-5 — I filed a REPORT where the work was doable, and only the adversarial pass caught
it.** V6 was declared an unclosable wire gap on a premise I never tested: I read
`ProductLinkCandidate` looking for a "deciding anchors" field, did not find one, and stopped —
without reading `decisionRuleForCandidate`, the function that derives exactly that from two fields I
had already confirmed were on the wire. R-24 exists for criteria that genuinely cannot be done
totally; it is not cover for one I had not finished investigating. The generalizable rule: **before
declaring a wire gap, read the function that CONSUMES the wire on the other side** — a derivation
that the backend already performs is not a gap, it is unread code. This is the third time on this
mission that a confident claim about the repo turned out to be an unverified reading (cf. the M-06
brief's `rule_matched`, and D-122's "SKU ML"), which is the pattern worth ratifying, not the
individual slip.

**FINDING-6 — observed, NOT fixed: the drawer still prints a raw machine name.**
`VinculoDrawer.tsx:107` renders `<Fact label="Entrada de match">{candidate.match_input}</Fact>`,
i.e. the literal wire token `seller_sku` / `ean` / `title` / `none`, three feet from a table that
now carefully translates those same tokens through `anchorShortLabel`. Pre-existing — it is
identical on `main` (`git show main:…/VinculoDrawer.tsx:104`), so it is not a regression, and no
V-criterion covers it. It IS inside my write-set, so I could have fixed it; I did not, because a
correct fix needs display names for `manual` and `none` too (which `anchorShortLabels` does not
have, and inventing them is a vocabulary decision, not a rename) and because expanding the diff
after the gate had already reviewed it trades a real risk for a cosmetic gain. Flagging it so it
is a known open item rather than something the next reviewer discovers.

**Answering the "did replacing GTIN lose information?" question directly: no.** The table cell
rendered a BOOLEAN summary (`✓ igual` when `match_input === "ean" && match_value`). The drawer's
`GTIN / Ref. interna` Fact (`VinculoDrawer.tsx:100-106`) renders the actual value —
`internal_reference_code`, falling back to `match_value` on the same EAN condition. The operator
therefore still reaches the GTIN, in a strictly richer form than the cell that was replaced.

## REQUESTS (to the hub — no edit made)

**REQUEST-1 (optional, no longer blocking) — consider surfacing `rule_matched` on the wire.** V6 is
now implemented by mirroring `decisionRuleForCandidate` in the FE, which means the same rule is
derived in two places from the same two inputs. That is correct today and compiler-guarded on both
unions, but it will drift the day the backend's rule changes and the FE's copy does not. Exposing
`product_link_decisions.rule_matched` (or a `decided_by: string[]`) on `ProductLinkCandidate` would
collapse the duplication to a read. Contract + SDK + backend = one writer, not this chip. Filed as
a maintenance note, not a defect: nothing today is wrong.

**REPORT-2 — the invariant the "Identificado por" column silently depends on, filed for ratification
at the hub's instruction.**

The column renders `CODPROD + EAN` for every `ACCEPT`. That is a READING of the backend, not a
guess, and it is honest for exactly one reason:

> `LinkCandidateMatchStatusAccept` has **exactly one assignment site** in the whole server:
> `apps/server_core/internal/modules/product_links/application/generation_service.go:505`, inside
> `buildConcordantCandidate`, whose candidate is built from a `seller_sku` FOR reason and an `ean`
> FOR reason resolving the **same** codprod (`:493-497`), reachable only from `:293-303`.

Verified by grep across `apps/server_core/` — one site, not "the sites I found". Both the Sol side
and the Opus side re-derived it independently and reached the same single site.

**The invariant, stated so it can be broken loudly instead of quietly:** *any second route to
`ACCEPT` makes the column assert a corroboration that did not happen.* No test fails, no type
breaks, and no reviewer sees it — the FE keeps saying `CODPROD + EAN` about a candidate that only
one anchor decided. That is precisely the ADR-17 class of defect D-122 exists to kill, arriving
silently.

The hub is taking this to its backlog as a candidate compile-time assertion, on the M-02 doctrine:
there, an optional port erased by a decorator was only caught later because nothing pinned it. This
REPORT is the pin's paperwork — a chip cannot add the guard (it lives in `apps/server_core/`, zero
Go in this write-set), so it names the site, the invariant, and the failure mode instead.

---

## Dispatch ledger

| # | Role | Model / effort | Path | Artifact (in-repo, committed) | Outcome |
|---|---|---|---|---|---|
| 1 | Adversarial gate reviewer, round 1 (commit `fa6ca3a2`) | `gpt-5.6-sol` / `medium` | OS-process, explicit 0.144.4 binary, stdin closed, `--sandbox read-only` | prompt `evidence/PROMPT-gate-vincneutro-rev1.md` · verdict `evidence/REVIEW-gate-vincneutro-rev1.md` (3053 B) | **REFUTED** — V6 FAIL, V11 NOT-PROVEN. Both accepted; both fixed. |
| 2 | Adversarial gate reviewer, round 2 (commit `7a343fea`, the fix) | `gpt-5.6-sol` / `medium` | same | prompt `evidence/PROMPT-gate-vincneutro-rev2.md` · verdict `evidence/REVIEW-gate-vincneutro-rev2.md` (2573 B) | **REFUTED**, 10/11 PASS. Derivation confirmed faithful; one finding (column placement) ESCALATED to the hub, not self-graded. |
| 3 | Implementation + planning of the slice | **Opus, in-session — `inline (DESVIO §4.2)`** | none | — | Not dispatched. Hub ruling: do not re-do it, do not launder the ledger row, and make the Opus gate name these hunks as its priority target. Done — see `PROMPT-gate-vincneutro-opus.md` §PRIORITY TARGET. |
| 4 | **Dual gate, OPUS side** (commits `fa6ca3a2` + `7a343fea`) | `harness:gate-reviewer` subagent, **opus**, read-only seat (Read/Grep/Glob; no Bash, no Write) | Agent tool, synchronous, frozen prompt read from disk | prompt `evidence/PROMPT-gate-vincneutro-opus.md` · verdict `evidence/REVIEW-gate-vincneutro-opus.md` | **REFUTED** — 9 PASS, 2 NOT-PROVEN (both seat limits), 4 findings. Three fixed in `394c83c`; the fourth is the standing escalation. |
| 5 | **Re-gate of the fixes, GPT side** (commit `394c83c`) | `gpt-5.6-sol` / `medium` | OS-process, stdin closed, `--sandbox read-only` | prompt `evidence/PROMPT-gate-vincneutro-rev3.md` · verdict `evidence/REVIEW-gate-vincneutro-rev3.md` (`.last.md` 2218 B, non-zero, copied verbatim) | **REFUTED** — **V10 FAIL**. Caught what the Opus side explicitly cleared: `providerDisplayName` was lossy, colliding two registrable provider codes onto one name. Both findings accepted and fixed. |
| 6 | **Re-gate of the fixes, OPUS side** (commit `394c83c`) | cold Opus subagent, read-only seat (no Bash / Edit / Write) | Agent tool, background, same frozen prompt as row 5 | verdict `evidence/REVIEW-gate-vincneutro-opus-rev2.md` — **transcribed, not captured: the task `.output` artifact came back 0 bytes** | **REFUTED** — 8 PASS, 3 NOT-PROVEN (all seat limits: no Bash, no EVIDENCE.md), 4 findings, all four verified independently and fixed. Its `V10: PASS` is **wrong** — it reviewed the pre-injectivity function and explicitly cleared the collision row 5 found. |

Artifacts are copied **into the repo** under `evidence/`, not left in the session scratchpad — the
scratchpad is temp-dir and does not survive. Every `.last.md` was checked non-zero before being
filed (streaming is not persisting: the three CHIP-ANCHORS-2 reviewers left 0-byte `.output` files
and those verdicts are irrecoverable).

**Row 6 is the exception and is labelled as one.** The Opus re-gate's on-disk `.output` came back
**0 bytes** — the same failure mode as CHIP-ANCHORS-2 — so its verdict reached this session only
through the completion notification and is *transcribed*, not captured. It is filed as transcribed,
with the provenance stated at the top of the artifact, rather than presented as a captured file.
What keeps the round auditable is not the transcription: all four of its findings were re-derived
independently against the named backend sources before any fix was written, and that derivation is
in `evidence/V-round4-must-fail.txt` with file:line. The findings stand on that, not on the text.

Planning and implementation ran **in-session (Opus, this chip session)**, not dispatched: the slice
is one screen module inside an already-authored pack, and the pack named the defect and its
location. No codex planner or implement worker was dispatched — recorded as a fact, not omitted.

## Reviewer verdicts and disposition

Both verdicts are filed verbatim in `evidence/`. Neither is paraphrased away.

### Round 1 — `REVIEW-gate-vincneutro-rev1.md` — **REFUTED**

| Finding | Disposition |
|---|---|
| **V6 FAIL** — the wire CAN supply the deciding anchors; the REPORT was unfounded | **ACCEPTED, fixed** in `7a343fea`. The reviewer was right and my premise was false. See V6 above and FINDING-5. |
| **V11 NOT-PROVEN** — only the after-count existed on disk | **ACCEPTED, fixed.** Baseline re-run and persisted as `evidence/L1-vitest-baseline.txt` (62 files / 511 tests). |

The reviewer's V6 line also framed the column as replacing `Canal`/`GTIN`; that framing is what
round 2 pressed on, and it is the open item below.

### Round 2 — `REVIEW-gate-vincneutro-rev2.md` — **REFUTED** (10 of 11 PASS)

The fix itself was confirmed: *"The status/input derivation at QueueRow.tsx faithfully mirrors
resolution_service.go:812"*. V11 flipped to PASS on the persisted baseline. One finding stands:

| Finding | Disposition |
|---|---|
| **V6 FAIL** — keeping `Canal` overrides D-122:136's requirement that the new column replace **both** former slots | **NOT resolved by this chip. ESCALATED to the hub** (see below). |

### HUB RULING — the escalation is CLOSED, and it closed by correcting the source

**Decision: option 1. The code stays as shipped.** Delivered by the hub after it verified both
sides against `5441fe18` and against the wire. Recorded here verbatim in substance, with the four
facts it turned on, because a pack that leaves two contested REFUTED verdicts standing is worse
than one that is wrong and says so.

| # | Fact the hub verified itself | Where |
|---|---|---|
| 1 | `ProductLinkCandidateItem` carries **no seller-SKU field** at all — the full field list was enumerated. The chip's premise holds. | `packages/sdk-runtime/src/index.ts:1070-1089` |
| 2 | The cell under the `SKU ML` header rendered `candidate.provider_code` raw — the slug `"mercado_livre"` reaching the operator under a label promising a SKU. **The sixth surface of the `provider_code` trap**: four in CHIP-PED-FILA, one in `PedidosTable.tsx:101`, this one now. | `QueueRow.tsx:218` @ `5441fe18` |
| 3 | The old `GTIN` cell rendered `✓ igual` or `—` — a **derived boolean, never a GTIN**. Replacing it erases no datum; it swaps a poor derivation for a better one. | `QueueRow.tsx:253-260` @ `5441fe18` |
| 4 | **The column count did not change.** Nine `th` before (select · Anúncio ML · SKU ML · Produto sugerido · SKU HUB · GTIN · Confiança · Motivo · Ação), nine after (select · Anúncio · Canal · Produto sugerido · SKU HUB · Identificado por · Confiança · Motivo · Ação). Two columns renamed **in place** — nothing added, nothing retained extra. | both trees |

Fact 4 is the one that dissolves the finding. Both reviewers wrote that the chip "kept" a slot and
"overrode a two-column replacement" — language that describes a table that got **wider**, and it did
not. The count is identical on both sides.

**Why D-122:136 did not mandate what the reviewers read.** The frozen line is
`Coluna "Identificado por" (/vinculos, substitui "SKU ML"/"GTIN"):` — a **parenthetical locator**
("where on the screen this goes"), followed by three bullets that are entirely about the column's
CONTENT: the anchors that decided joined by ` + `, the distinction against Motivo, and double
corroboration as the D-121 explanation. **No bullet concerns column count.** The normative part is
satisfied whole. Two cold reviewers promoted the hub's parenthetical to a clause — a defensible
reading of a block marked frozen, which is why the hub did not fault them for it, and wrong.

And the parenthetical carried the falsehood: it was written believing `SKU ML` held a seller SKU.
It never did. **Under R-25 a false sentence in an authority artifact is CORRECTED, not annotated** —
the hub corrects D-122 on its own side, with the two `file:line` above. This chip does not touch the
decisions file.

**Why not option 3, which was my own reading.** The hub judged the relocation good engineering and
refused it on economy, not merit: the branch is closed and awaiting the Opus gate side, so changing
the `Anúncio` cell's shape would invalidate the golden's already-measured `within(row)` assertions
and leave the gate reviewing a different diff than the one that was closed; and with a single
installation `Canal` is constant on every row, so the gain is marginal until a second provider
exists — at which point relocation becomes its own chip and pays the golden cost once. The hub's
closing note on this is the rule I want on record, because it is the one I was applying:

> Você fez a chamada certa ao não fazer sem pedir. Grant de write-set não é grant de design.

**So the divergence was resolved by correcting the normative source, not by a chip override.** That
distinction is the whole point: had I shipped option 1 on my own reading, I would have been right
about the facts and wrong about the authority. The facts were mine to establish; the ruling was not.

---

**Why I did not grade this PASS myself.** Two independent cold reviewers read D-122:136 literally and
called the retained `Canal` an override of a frozen decision. My counter-argument is in V6 and
FINDING-4 and I still think it is correct on the merits — D-122's "substitui 'SKU ML'" is founded on
the belief that the cell held a seller SKU, and it never did. But a chip does not get to overrule a
frozen mission decision on its own reading, and it certainly does not get to mark its own contested
call PASS after two REFUTED verdicts. The code is left in the state that **loses no information**;
the tie-break is the hub's.

**Three options, for the hub:**

1. **Keep as shipped** — `Canal` + `Identificado por`, GTIN replaced. Satisfies V10 literally,
   deviates from D-122's plural. Zero further work; amend D-122:136 to name GTIN only.
2. **Literal D-122** — drop `Canal`, both old slots collapse into `Identificado por`. One-line
   change. Costs V10's second half: the queue table would no longer say which marketplace an
   anúncio belongs to (it survives only in the drawer subtitle).
3. **Dissolve the conflict** — move provider identity INTO the Anúncio cell as a secondary line
   (`MLB123` / `Mercado Livre`), freeing the former SKU-ML slot entirely. Then D-122's two-slot
   replacement AND V10's "provider stays as data" are both literally satisfied, and the column
   count drops by one as D-122 implies. This is the option I would take, but it is a visible design
   change to a golden-locked layout at close time, which is the hub's call, not a chip's.

### The OPUS side of the dual gate — `REVIEW-gate-vincneutro-opus.md` — **REFUTED**

Rounds 1 and 2 were both `gpt-5.6-sol`. The hub caught that: two Sol rounds are **two rounds of the
same side**, not a dual gate, and this session cannot be the Opus side of a gate on code it wrote
(§4 item 3, implementer ≠ reviewer). Correct, and it is the mirror image of what happened to
CHIP-ANCHORS-2 in the other direction. So a cold Opus reviewer ran on the frozen input, in a
**physically read-only seat**, forbidden by its prompt from reading this EVIDENCE or either Sol
verdict, and pointed at the inline-written hunks as priority target per the hub's §4.2 ruling.

It came back **REFUTED with four findings, and three of them were real code defects that two Sol
rounds had not caught.** Dispositions:

| # | Finding | Disposition |
|---|---|---|
| A | D-122:136 vs the retained `Canal` — the chip closed a conflict with a frozen decision by itself | **STANDING ESCALATION.** The reviewer independently verified my premise (`ProductLinkCandidateItem` carries no seller-SKU field, `sdk-runtime/src/index.ts:1070-1089`; the hub's own live drive recorded the cell rendering `provider_code`) and judged the deviation *defensible on the merits* — and still refused to bless it, because "the frozen premise is factually wrong" is a hub ruling. That is exactly my position. Third independent reviewer, same conclusion: not a chip's call. |
| B | `VinculosDesign.golden.test.tsx` — the corrected fixture is still not producible: `CONFIRM` with `confidence: 92`/`ALTA` | **ACCEPTED, fixed.** Verified myself first: `exact_ean` CONFIRM emits 60/MEDIA (`generation_service.go:538`), `exact_sku` CONFIRM 70/MEDIA (`:529`), and `ALTA` is emitted at exactly ONE place — 95, ACCEPT, `buildConcordantCandidate` (`:503-505`). So 92/ALTA was unreachable for any status. My REVIEW→CONFIRM fix in the previous round was half a fix. |
| C | `providerDisplayName` falls through verbatim → an unmapped provider puts a raw slug on screen | **ACCEPTED, fixed.** Latent, not live (only `mercado_livre` is mapped today), which is precisely why it would have shipped: the defect waits for the second marketplace. |
| D | `decidingAnchors` invokes an unchecked map lookup → a wire `match_status` outside the SDK union throws `undefined is not a function` **during render** | **ACCEPTED, fixed.** The most severe of the four: it does not degrade a cell, it unmounts the queue table. |

**On B, C and D I did not take the reviewer's word.** B was checked against the generator's own
`switch` before touching the fixture; C and D were reproduced as failing tests before the fix. D
reproduced verbatim:

```
TypeError: statusDecidingAnchors[candidate.match_status] is not a function
 ❯ decidingAnchors src/pages/vinculos/QueueRow.tsx:109:54
 ❯ QueueRow src/pages/vinculos/QueueRow.tsx:331:19
```

Full run in `evidence/V-opus-findings-must-fail.txt` — 2 failed / 11 passed, written before the fix,
green after. This is the same discipline as V3 and V6: the guard is proven by the failure it
catches, not by the passing run.

**The two NOT-PROVEN are seat limits, not defects, and I am not flipping them to PASS myself.**
V8 asks for a correction recorded *in this EVIDENCE*, which the reviewer was forbidden to open — it
corroborated the substance independently instead (`rule_matched` absent from `contracts/` and
`sdk-runtime`, and the `0082` CHECK reading `actor <> 'system' OR rule_matched =
'concordant_codprod_ean'`). V11's scope half needs `git`, and the seat has no Bash; the hub has
since verified the same scope independently and said so.

### What the three fixes changed

- `QueueRow.tsx` — `providerDisplayName` now TYPESETS an unmapped code instead of passing it
  through: separators to spaces, words capitalised, every character the wire sent still legible in
  order (`shopee` → `Shopee`). That is deliberately different from `anchorShortLabel`, which still
  returns an unknown anchor verbatim, and the docblock says why: an anchor is operator vocabulary
  only the ERP can name, a `provider_code` is a machine identifier that never belonged on screen in
  that form.
- `QueueRow.tsx` — `decidingAnchors` guards the lookup and returns `[]`, which renders the honest
  `—`. The `Record<Union, …>` above it is untouched and still fails the build on a sixth status:
  the compile-time guard answers "the SDK grew", the runtime guard answers "the wire grew first".
  Two different failures, and only one of them can be caught by a compiler.
- `VinculosDesign.golden.test.tsx` — the EXEMPLO-IO row is now the producible `exact_ean` CONFIRM
  (60/MEDIA) *with its second anchor as INCOMPARABLE on the provider side*, which is the state this
  whole chip exists to render, so the design golden now gates it too. The ALTA/accent token path
  did not lose coverage: the off-theme sweep's first row became the producible ACCEPT (95/ALTA,
  `buildConcordantCandidate` verbatim).

---

### ROUND 4 — the re-gate of the fixes — **both sides REFUTED, both were right**

Both sides ran on the same frozen input (`394c83c`) and the same frozen prompt
(`evidence/PROMPT-gate-vincneutro-rev3.md`). Both came back REFUTED. Six findings between them, all
six verified independently before any code moved, all six fixed. Full verbatim trail with the red
runs: `evidence/V-round4-must-fail.txt`.

**The two sides did not overlap on the thing that mattered, and that is the argument for the dual
gate as a mechanism, not as a ritual.** The GPT side's single FAIL (V10, `providerDisplayName` is
lossy) is precisely the point the Opus side had just cleared *by name* — attack 2, "the split
`/[_\-\s]+/` only collapses separators, so no registered `provider_code` typesets into another
provider's name." That sentence is false. `registry.go` `buildDefinitions` dedupes provider codes by
exact string equality only, so `amazon_marketplace` and `amazon-marketplace` are both legal and can
be registered simultaneously — and both rendered "Amazon Marketplace". Two providers wearing one
name is wrong information; the raw slug it replaced was only ugly. One reviewer asserted the
property; the other tested it.

| # | Source | Finding | Verified against | Disposition |
|---|---|---|---|---|
| 1 | GPT rev3 | `providerDisplayName` collapses distinct registrable codes onto one display name | `integrations/adapters/providers/registry.go` `buildDefinitions` — dedupe is exact string equality | **FIXED.** Typeset applied only where it is INJECTIVE: the form must round-trip to the exact wire code, else verbatim. |
| 2 | GPT rev3 | the unmapped-provider test proves one hard-coded pair and survives a constant fallback | the test itself | **FIXED.** Now renders three codes in one pass and asserts two stay distinct. Proven red against BOTH implementations the reviewer named — the lossy collapse *and* the fabricated `return "Shopee"`. |
| 3 | Opus rev2 | drift hardening covered `match_status` only; `direction` is an unchecked `Record` index **and the V2 fix made it reachable** | `QueueRow.tsx` by grep-by-STRING | **FIXED, and extended.** The reviewer's reachability argument is exactly right: the old literal enumeration silently DROPPED an unknown direction, the new total sort keeps it. The same grep found a **fourth** site the reviewer did not name — `bandLabels`/`bandClasses` on `confidence_band` — whose pill renders with a class attribute ending in the literal `undefined`. Both hardened; fall-through is the wire value verbatim, never a known member. |
| 4 | Opus rev2 | both "made honest" golden fixtures are still not producible — the `marca` motivo is missing | `KnownIdentityAnchors` (4 anchors) → `identity_anchor_adapter.go:28-35` (declares all 4) → `mercado_livre/capability_adapter.go:90` (declares 3 ⇒ `marca` unsupplied) → `classifyProviderIdentityAnchor` (:700-702, UNAVAILABLE, emit=true) | **CONFIRMED and FIXED.** Every mercado_livre candidate carries `marca` UNAVAILABLE. With 3 reasons and a 2-chip cap, production renders two chips **plus a "+1" toggle** — the golden was locking a layout the backend never emits. It now asserts the toggle. |
| 5 | Opus rev2 | the NO_CANDIDATE fixtures carry `state: "exact_ean"` with `reasons: []` — same class, same file, left untouched by the previous fix | `applyUnresolvedScore` is the only NO_CANDIDATE path and is reached via `LinkCandidateStateUnresolved` (:215, :379), always seeding two absence reasons (:620-628) | **CONFIRMED and FIXED.** `state: "unresolved"` + the two INCOMPARABLE/erp reasons + `marca`. That fixture's `provider_code: "xyz-9"` also went to `mercado_livre`: no capability declaration exists for `xyz-9`, so no reason array could be claimed producible for it. |
| 6 | Opus rev2 | after the swap the band pill has no positive assertion; only the negative off-theme sweep covers it | `git show bcab8269:…golden.test.tsx` | **FIXED, with one correction to the finding.** The finding left "did the pre-swap golden assert it?" NOT-PROVEN. It did not — pre-swap asserted `92%` and no band pill, and its only `ALTA` mention was a *negative* assertion in the NO_CANDIDATE row. So the swap did not LOSE the coverage; it never existed. The finding is still right that it should exist, and it now does on both bands. |

**What round 4 cost the chip, stated plainly.** Two of the six findings (#4, #5) are the *same
defect class the previous round's fix claimed to have closed*, one of them in the very file that fix
edited. "I made the fixture honest" was an overclaim: what was actually done was to correct the
confidence/band/status of the fixture, while its reasons array stayed a plausible invention. The
generator was read for the scoring fields and not for the finalizing step
(`appendProviderDeclaredUnavailableReasons`) that runs after every one of them.

## Verdict roll-up

| Criterion | Verdict |
|---|---|
| V1 `INCOMPARABLE` glyph + own tokens | PASS |
| V2 all-`INCOMPARABLE` row shows a motivo | **PASS** (deciding criterion) |
| V3 must-fail of V2 | PASS |
| V4 `side` reaches the operator | PASS |
| V5 `tsc` write-set clean, 12 baseline declared | PASS |
| V6 "Identificado por" | derivation **PASS** (all three gates confirm it mirrors the backend, the Opus side walking all 25 status×input pairs) · column placement **RULED — PASS**: hub decided option 1, code unchanged, D-122:136's parenthetical corrected at the source (see the ruling above) |
| V7 auto-approved badge, 3 DOM cases | PASS |
| V8 F-04 brief correction recorded | PASS |
| V9 `refforn` decision declared | PASS (KEPT) |
| V10 neutral vocabulary, provider data intact | PASS — and now for a provider the map does not know, which is what the Opus side caught |
| V11 no collateral damage | PASS |

**Lanes, final.** `tsc` 15 → **12, zero under `pages/vinculos/`**
(`evidence/L0-tsc-after-round4.txt` — the same 12 as the baseline; this chip neither added nor
removed one); vitest 62/511 → **63/524** (`evidence/L1-vitest-after-round4.txt`), the +13 over the
baseline being the regression tests the four gate rounds forced, the last two being the `direction`
and `confidence_band` wire-drift tests.

**P6 status — NOT self-declared discharged.** The earlier `AGREEMENT` in this file was written when
both gate rounds were `gpt-5.6-sol`; the hub refused it, on two grounds that were both right: there
was no Opus side, and a REFUTED round with an open finding is not an agreement. The Opus side then
ran and came back REFUTED with three real code defects, all fixed in `394c83c`.

The escalated item is now **closed by hub ruling** (option 1, code unchanged, the false parenthetical
in D-122:136 corrected at its source — full record above). So nothing is open on the contract side.

Both sides then re-gated `394c83c` against one frozen prompt
(`evidence/PROMPT-gate-vincneutro-rev3.md`) and **both returned REFUTED** — six findings, all
verified independently, all fixed in this commit (round 4 above). The prompt's first attack, whether
the runtime guard loosened the compile-time exhaustiveness this chip is built on, was answered by
construction before either reviewer reported: widening the union by one member makes the compiler
name the missing member (`evidence/V-compiletime-guard-still-bites.txt`), and both sides
independently reached the same conclusion.

**Where that leaves P6, honestly: no round of this gate has yet returned APPROVED on the code as it
now stands.** Four rounds have run, each found real defects, each was fixed — but the fixes to
rounds 3 and 4 have themselves not been reviewed by anyone but their author. The last two rounds are
the reason not to wave that away: round 4 found that two of round 3's fixes were the *same* defect
class round 3 claimed to have closed. This chip is not in a position to certify that round 4 broke
that pattern, because it is the same author.

So this file does **not** declare P6 discharged and does not carry an `AGREEMENT` marker. What it
asserts is narrower and checkable: V1–V11 hold on the current tree, the lanes are green
(`L0-tsc-after-round4.txt`, `L1-vitest-after-round4.txt`), and every finding raised by any reviewer
is either fixed with a red-first artifact or recorded as a REPORT/REQUEST.

Whether one more round is required, and whether the ledger's row-6 provenance gap (a transcribed
rather than captured Opus verdict) is acceptable, are the hub's calls. The `P6-DUAL-GATE:` line and
the `AGREEMENT` marker remain the hub's to write. Not this chip's.
