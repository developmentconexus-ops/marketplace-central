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

---

## Dispatch ledger

| # | Role | Model / effort | Path | Artifact (in-repo, committed) | Outcome |
|---|---|---|---|---|---|
| 1 | Adversarial gate reviewer, round 1 (commit `fa6ca3a2`) | `gpt-5.6-sol` / `medium` | OS-process, explicit 0.144.4 binary, stdin closed, `--sandbox read-only` | prompt `evidence/PROMPT-gate-vincneutro-rev1.md` · verdict `evidence/REVIEW-gate-vincneutro-rev1.md` (3053 B) | **REFUTED** — V6 FAIL, V11 NOT-PROVEN. Both accepted; both fixed. |
| 2 | Adversarial gate reviewer, round 2 (commit `7a343fea`, the fix) | `gpt-5.6-sol` / `medium` | same | prompt `evidence/PROMPT-gate-vincneutro-rev2.md` · verdict `evidence/REVIEW-gate-vincneutro-rev2.md` (2573 B) | **REFUTED**, 10/11 PASS. Derivation confirmed faithful; one finding (column placement) ESCALATED to the hub, not self-graded. |

Artifacts are copied **into the repo** under `evidence/`, not left in the session scratchpad — the
scratchpad is temp-dir and does not survive. Both `.last.md` files were checked non-zero before
being filed (streaming is not persisting: the three CHIP-ANCHORS-2 reviewers left 0-byte `.output`
files and those verdicts are irrecoverable).

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

**Why I am not grading this PASS.** Two independent cold reviewers read D-122:136 literally and
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

---

## Verdict roll-up

| Criterion | Verdict |
|---|---|
| V1 `INCOMPARABLE` glyph + own tokens | PASS |
| V2 all-`INCOMPARABLE` row shows a motivo | **PASS** (deciding criterion) |
| V3 must-fail of V2 | PASS |
| V4 `side` reaches the operator | PASS |
| V5 `tsc` write-set clean, 12 baseline declared | PASS |
| V6 "Identificado por" | derivation **PASS** (both gates confirm it mirrors the backend) · column **placement ESCALATED** — hub tie-break between D-122:136 and V10 |
| V7 auto-approved badge, 3 DOM cases | PASS |
| V8 F-04 brief correction recorded | PASS |
| V9 `refforn` decision declared | PASS (KEPT) |
| V10 neutral vocabulary, provider data intact | PASS |
| V11 no collateral damage | PASS |

**AGREEMENT — P6 discharged**, with one open item escalated rather than closed: the
`Identificado por` column placement (D-122:136 literal vs V10). The chip does not claim that
decision. Everything else is PASS with artifacts on disk.

The `P6-DUAL-GATE:` line is the hub's to write, not this chip's.
