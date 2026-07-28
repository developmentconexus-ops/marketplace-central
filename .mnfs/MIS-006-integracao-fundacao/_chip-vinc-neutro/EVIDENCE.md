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
`INCOMPARABLE` matched no branch, so it was not ranked last — it was UNDISPLAYABLE. A filter
wearing a ranking's comment, against the ADR-17 invariant the file documents at `:154-156`.

> **CORREÇÃO (rodada 8, escrita depois do drive do hub em dado real).** Este parágrafo dizia antes
> que *"uma linha cujos motivos são todos `INCOMPARABLE` … renderiza um `+2` sozinho com zero
> chips"*, e apresentava isso como o defeito. Como afirmação sobre a TELA, é falsa, e o erro é meu:
> confundi **producível** com **alcançável**.
>
> `resolveIdentityAnchors` (`:149-169`) aborta a geração se a declaração não resolve, e a única
> declaração desta árvore é `mercado_livre`, que **não** supre `marca`. Logo toda linha viva carrega
> um `marca` UNAVAILABLE — direção que a expressão velha **casava**. A célula vazia exige um provider
> que supra as quatro âncoras; nenhum existe. O hub confirmou dirigindo: nenhuma linha sem chip.
>
> O meu próprio mecanismo já tinha registrado metade disso antes do drive — a fixture do teste é
> `driftCandidate` com a razão escrita *"No declaration emits it … unreachable"*. O que eu não fiz
> foi propagar essa conclusão para cá, onde a alegação de valor mora. Guard corrigido, prosa
> mantida: mesma classe do chip inteiro, um nível acima do código.
>
> **O defeito ALCANÇÁVEL, que é o que este fix compra**, é menor e real: com `COMPACT_CHIP_LIMIT`
> = 2, uma linha com um FOR, um `ean` INCOMPARABLE e o `marca` UNAVAILABLE gastava os dois slots no
> FOR e no `marca` — permanente, acionável por ninguém — enquanto o INCOMPARABLE, sobre o qual o
> operador PODE agir, não aparecia em limite nenhum. `directionRank` põe INCOMPARABLE (2) à frente
> de UNAVAILABLE (3), então essas linhas passam a mostrar a ausência acionável. O drive do hub achou
> **três** linhas exatamente dessa forma, vivas hoje.

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
next to `<button aria-label="Mostrar todos os 2 motivos" …>+2</button>`. `Tests 1 failed | 9 skipped
(10)`. Fix restored, file green.

> **ESCOPO deste must-fail, corrigido junto com o de cima.** Ele roda sobre a fixture INALCANÇÁVEL,
> logo prova que o guard DISCRIMINA — a expressão velha volta, o teste fica vermelho — e **não**
> prova que a tela de hoje mostrava a célula vazia. Não mostrava. O must-fail que fala do dado vivo
> é o vizinho, `ranks INCOMPARABLE above UNAVAILABLE without dropping either`, que dirige o mesmo
> defeito a partir de um conjunto de motivos que o backend realmente emite. Dizer "exatamente a tela
> que o contrato proíbe" era emprestar ao artefato inalcançável a autoridade do alcançável.

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

**Confirmado em dado vivo pelo hub (rodada 8), e é aqui que o valor deste chip está.** O drive achou
três `ean INCOMPARABLE` na tela, com `detail` **byte-idêntico** e sides diferentes:

```
side=both      detail=sem EAN para corroborar o CODPROD
side=both      detail=sem EAN para corroborar o CODPROD
side=provider  detail=sem EAN para corroborar o CODPROD
```

`both` = nenhum dos dois lados tem EAN; `provider` = o ERP tem, o anúncio não. Fatos operacionais
diferentes — um manda cadastrar no ERP, o outro manda corrigir o anúncio — com texto igual. Quem
separa é `reasonSideLabel`, lendo o campo. Sem ele a distinção **não existe na tela**, e a única
alternativa seria adivinhar pelo português do `detail`, que é o mesmo nos três. Este é o defeito
alcançável do chip, medido em dado real, e não é o que o critério V-1 descreve.

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

**FINDING-7 — the FE's golden fixtures are coupled to a BACKEND capability declaration, and nothing
on the FE side can detect it drifting.** Class A of `evidence/V-fixture-producibility-sweep.md`
proves each golden fixture producible by reading the generator: `KnownIdentityAnchors()` has four
members, `mercado_livre` declares three, so `marca` arrives UNAVAILABLE and every candidate carries
that reason. That proof is a **point-in-time read of Go code from a TypeScript test.** If someone
adds a fifth known anchor, or has `mercado_livre` declare `marca`, every class-A fixture silently
goes back to non-producible and the golden goes back to gating a fiction — with the FE suite fully
green, because nothing in `apps/web` imports any of it.

That is the same shape as the defect this chip exists to fix (`QueueRow.tsx:159` enumerated a union
by string literal: type-correct, and therefore invisible), one layer up: a cross-language coupling
with no compiler and no test on the seam.

Not fixed here, and deliberately: a contract test that asserts the FE's fixtures against the
generator's real output would need a Go-side fixture dump or a golden JSON emitted by the backend —
both are backend seams, outside this chip's write-set, and neither is a rename. **Recommend the hub
route it as its own slice.** The cheapest honest version is a generator-side test emitting one
canonical candidate per status to a checked-in JSON, which the FE golden then reads instead of
hand-writing — moving the coupling from "a comment claims it" to "the build breaks".

> **F-07 CLOSED in round 5 — and the paragraph above is why it should not have been left open.**
> The Opus seat raised this same coupling as a finding, which is the correct answer to a risk this
> chip filed as "accepted": *accepted residual risk* was, here, a way of deferring the work to the
> next reader. It is closed by `apps/web/src/pages/vinculos/wireFixtures.guard.test.ts`, which reads
> `connectors/ports/marketplace_capability.go` and `mercado_livre/capability_adapter.go` at test
> time and asserts set equality with the FE's constants **in both directions** — a vocabulary that
> only grew and one that only shrank are different failures and both are real (`refforn` is the one
> that actually shrank, and it left a fixture asserting a reason no provider can emit).
>
> The correction to the recommendation above: **no backend seam was needed.** The paragraph
> concluded "outside this chip's write-set" from the assumption that closing the gap required the
> generator to EMIT something. Reading the declaration is a FE-side test file and touches no Go. The
> reasoning was sound about a dump-based design and wrong about the cheapest one, which is how a
> real constraint (write-set) came to justify not doing work that was inside it.
>
> What it still does not prove, stated so a green build is not over-read: it compares against the
> declaration **in this checkout**, not against a deployed server. The generator-side JSON dump
> remains the stronger form, and remains a hub call.

**FINDING-8 — a cold read of two Go files once blew vitest's 5s `testTimeout`, and the cause was
NOT established.** The first run of `wireFixtures.guard.test.ts` hung and timed out on
`readFileSync` of `marketplace_capability.go`, while the very next test in the same file read the
same two files in 2ms. Isolating the regexes in plain node was instant.

Reported as a FINDING and not as a fix, because the honest state is unresolved: I formed a
hypothesis about which `expect` form inside a `.map()` was hanging, tested all three candidate
forms, and **all three passed** — the hypothesis was wrong. Three later runs, including a probe
built to reproduce the hang, all came back in single-digit milliseconds.

The reads now happen at MODULE scope. That is a **mitigation, not a diagnosis**: module-load time is
not charged to `testTimeout`, so a slow cold read cannot turn a passing guard into a red lane. The
comment in the file says exactly this rather than claiming a cause. **For the hub:** if another chip
sees a first-touch filesystem read outside the vitest root time out under this runner on Windows,
that is a second data point and the pattern is worth a profile §3 false-alarm entry. One occurrence
with a refuted hypothesis is not.

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
| 7 | **Round 5 delta gate, GPT side** (range `394c83c..3915f33b`) | `gpt-5.6-sol` / `medium` | OS-process, stdin closed, `--sandbox read-only`, `-o` straight into `evidence/` | prompt `evidence/PROMPT-gate-round5-delta-sol.md` · verdict `evidence/REVIEW-gate-vincneutro-sol-rev4.md` (**captured** by `-o`, non-zero) · log `evidence/agent__gate-r5-sol.log` | **REFUTED** — V10 FAIL, 4 findings, all four verified against primary source and fixed. Caught the drawer's PRIVATE second copy of `bandClasses` — the structure round 4 hardened in `QueueRow` and declared closed — and the injectivity check running on the built string rather than the painted one. |
| 8 | **Round 5 delta gate, OPUS side** (same range) | `harness:gate-reviewer` subagent, **opus**, read-only seat (Read/Grep/Glob; no Bash, no Write) | Agent tool, background, frozen prompt read from disk | prompt `evidence/PROMPT-gate-round5-delta.md` · verdict `evidence/REVIEW-gate-vincneutro-opus-rev3.md` — **transcribed**, the seat has no Write by construction, provenance declared at the artifact's head | **REFUTED** — 9 PASS, V5/V11 NOT-PROVEN (as its prompt instructed for those two items), 4 findings, all verified and fixed. It did **not** reaffirm the refuted `V10: PASS`. Its finding 1 killed this chip's "producible under a capability declaration" disposition outright, citing `resolveIdentityAnchors`'s abort — **"a dodge, not an honest degrade"**, a classification accepted as stated. |
| 9 | **Round 6 delta gate, GPT side** (code delta `4c000a04..7b5c18eb`; pack read at the dispatch tip `2e9e9ce3`) | `gpt-5.6-sol` / `medium` | OS-process, stdin closed, `--sandbox read-only`, `-o` to `agent__r6sol.last.md`, no `&` | prompt `evidence/PROMPT-gate-vincneutro-round6.md` (frozen; brief-only commits after `7b5c18eb`) · dispatch wrapper `scratchpad/prompt-r6-sol.md` · log `agent__r6sol.log` | **DISPATCHED** at `2e9e9ce3`. Two earlier attempts on this row produced NO verdict and none was claimed: the first died on `&` inside an already-backgrounded shell (exit 0 on the wrapper, log 0 bytes, sentinel absent); the second was **ABORTED before reading** by this chip when the hub's executor seat found the population count short by two — measuring and dispatching are not parallel. |
| 10 | **Round 6 delta gate, OPUS side** (same delta, same tip) | `harness:gate-reviewer` subagent, **opus**, read-only seat (Read/Grep/Glob; no Bash, no Write) | Agent tool, background, same frozen brief read from disk | prompt `evidence/PROMPT-gate-vincneutro-round6.md` | **DISPATCHED** at `2e9e9ce3`. The prior attempt was **ABORTED before reading** alongside row 9 — the seat had produced one line, `"I'll start by reading the frozen brief."`, and nothing else. No verdict existed and none was claimed. This seat cannot run the brief's Section 3; the brief now says so in writing and routes custody to the executor seat, which holds the same clause. |
| — | **Executor seat, round 6** | the HUB, own detached worktree `.claude/worktrees/hub-exec-vinc`, own `npm ci` | not this chip's dispatch | reported in-message; lanes, four must-fail arms, merge-of-proof | **RAN at `2e5331b6`** — lanes confirmed independently (12 / 0-in-vinculos, 64 / 531), merge proven clean (0 conflicts, merged tree 67 / 544 green), and **found the sentinel defect this chip had not**. Re-run at `7b5c18eb` offered by the hub. |

**Round 5 is a DELTA round, not a fifth of the same** (hub ruling). Scope is the code no reviewer
has read: the round-4 fixes (`ebb309ac`) and the fixture sweep (`3915f33b`). Two seat-specific
instructions, both in the frozen prompts:

- The Opus prompt **quotes the previous Opus round's `V10: PASS` verbatim, with the reason it was
  refuted** (`registry.go` dedupes provider codes by exact string equality, so `amazon_marketplace`
  and `amazon-marketplace` coexist), so it cannot be reaffirmed by habit.
- The Sol prompt has **V5 and V11 removed from its scope**, with the reason stated in it: both were
  discharged by measurement from a seat that ran the commands, and a reading seat that answers
  "could not run it" has found nothing and burned a round (§11).

**Range extended from the hub's `394c83c..ebb309ac` to `394c83c..3915f33b`, declared to the hub
before any verdict returned.** The hub named the first because it was HEAD at ruling time; the sweep
landed afterwards and carries a new test plus two fixture annotations. Sending reviewers to a range
that excludes the one commit nobody but the author has read would hand them the reviewed part and
hide the new part.

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

### ROUND 5 — the delta re-gate — **both sides REFUTED, both were right, and the class broke open**

Two seats on the delta `394c83c..ebb309ac`, each with its own frozen prompt
(`evidence/PROMPT-gate-round5-delta.md` for Opus, `…-sol.md` for the GPT side, which had V5 and V11
removed from its scope). Verdicts: `evidence/REVIEW-gate-vincneutro-opus-rev3.md` (transcribed —
that seat has no Write by construction, provenance declared at its head) and
`evidence/REVIEW-gate-vincneutro-sol-rev4.md` (captured by `-o`).

Eight findings between them. All eight verified against primary source before anything moved.

| # | Source | Finding | Disposition |
|---|---|---|---|
| 1 | Opus rev3 | `cand_inc` / `cand_noside` carry a comment claiming they are "producible under a capability declaration" — but `resolveIdentityAnchors` ABORTS generation unless the declaration resolves (`generation_service.go:149-169`), so those reason arrays are unreachable under EVERY declaration, not merely today's. Called **"a dodge, not an honest degrade."** | **CONFIRMED.** The classification is accepted as stated. Both are now `driftCandidate(why, …)` with the reason written in the test. |
| 2–4 | Opus rev3 / Sol rev4 | three further fixture defects of the same class | **CONFIRMED, all fixed** — see the mechanism below rather than four more point-fixes. |
| 5 | Sol rev4 | `VinculoDrawer` keeps a PRIVATE second copy of `bandClasses`, raw-indexed — the exact structure round 4 hardened in `QueueRow` and declared closed | **CONFIRMED and FIXED.** The drawer's copy is deleted; `bandClass`/`bandLabel` are exported from `QueueRow` and imported. Red captured first: `class="… font-medium undefined"` in the drawer pill. |
| 6 | Sol rev4 | `providerDisplayName`'s injectivity check runs on the string it just built, not on what the browser PAINTS (HTML collapses whitespace runs) | **CONFIRMED and FIXED.** The round-trip now runs on the painted form. New test drives four codes, including `amazon__marketplace` and `_amazon`. |
| 7 | Opus rev3 | nothing links `wireFixtures.ts`'s copy of the anchor vocabulary to the Go source it mirrors — filed by this chip as accepted residual risk F-07 | **CONFIRMED and CLOSED.** `wireFixtures.guard.test.ts` reads `marketplace_capability.go` and `mercado_livre/capability_adapter.go` at test time and asserts set equality in BOTH directions. Must-failed both ways (vocabulary grown — using `refforn`, the anchor that actually shrank out; and shrunk — dropping `marca`). |
| 8 | both | V5/V11 unscorable from a read-only seat | **NOT-PROVEN**, as instructed. Recorded, not converted into a PASS. |

**The third occurrence fired the hub's pre-written rule.** Round 5 found the third instance of the
unproducible-fixture class, and then the fourth and fifth. Per the hub: *"não é point-fix: é
mecanismo nomeado + varredura exaustiva, regra da terceira rodada."*

**Mechanism:** `apps/web/src/pages/vinculos/wireFixtures.ts`. An impossible candidate is now
UNWRITEABLE rather than detectable — `wireCandidate` THROWS in the test that builds it, naming the
generator line the fixture contradicts. A deliberately impossible shape goes through
`driftCandidate(why, …)`, which requires the reason to be written next to the fixture. Backed by
`wireFixtures.guard.test.ts` (finding 7), which is the compiler that was missing from the seam.

**Exhaustive re-sweep, by machine.** Pointing `QueueTab.test.tsx`'s factory at the builder returned
**13 failed / 5 passed (18)** on the first run. This chip's own hand sweep
(`evidence/V-fixture-producibility-sweep.md`) had reported **two** findings in that file and called
itself EXHAUSTIVE in its title. Distinct violations: prose anchors outside the vocabulary (×2),
`refforn` (the anchor D-A REMOVED), **no `marca` reason (×6)** — the very invariant that sweep was
built to enforce — a `marca` INCOMPARABLE where mercado_livre can only emit UNAVAILABLE,
`(95, ALTA, REVIEW)`, and `(72, CRITICA, REVIEW)` (×2, the legitimate drift probes).

That document has been **corrected at source** per R-24: a banner at its head and a `## CORRECTION`
section naming each false sentence, including the enumeration that "proved" it exhaustive — a grep
whose `[a-z_]*` character class silently excluded every capitalized and accented anchor and reported
ONE violation where there were FIVE. Its verdicts are kept verbatim rather than edited into
correctness, because what they got wrong IS the finding.

`VinculosDesign.golden.test.tsx` was then routed through the same builder. Its four fixtures passed
unchanged — the hand corrections from round 4 were complete — and the wiring was must-failed
(`61/MEDIA/CONFIRM` → `wireCandidate: no scoring path assigns (61, MEDIA, CONFIRM)`) so that "green"
could not mean "never called".

### ROUND 5 — hub ruling, and the correction it forced

**The guard as first written had the defect it was built to catch.** The hub's item 2: reading Go
with regexes is string-scraping, and a rename, a moved file, or a literal replaced by a generated
`const` makes every pattern match NOTHING — *and a guard that finds nothing PASSES.* Green with the
Go intact, green with the Go gone. That is the same non-discriminating observable as the 15 tsc
errors, the `.bin`, and the gate that judged a phantom tree. My version asserted set equality and
never asserted that it could still SEE its subject.

Corrected as mandated, with all four arms must-failed and captured
(`evidence/V-guard-mustfail-round5.txt`):

| arm | what was broken | result |
|---|---|---|
| (a1) | FE vocabulary GREW (`refforn` added back) | RED — set inequality named |
| (a2) | FE vocabulary SHRANK (`marca` dropped) | RED |
| (b1) | Go source moved/renamed away | **RED, 3/3 tests** — `the file this guard reads for the identity-anchor vocabulary itself is gone or moved` |
| (b2) | Go source READABLE but the symbol renamed out of it | **RED, 3/3** — `no longer contains \`var knownIdentityAnchors\` … every assertion after this one would pass vacuously`, plus two `expected 0 to be greater than 0` from the extraction tests |

Arm (b2) is the one worth naming: the file opens, the guard reads it, and the vocabulary is simply
not in there. Before the fix that was a silent pass. The zero-extraction assertions now sit in three
places, because "extracted nothing" and "there is nothing to extract" are the two states this guard
exists to tell apart.

**What it is, stated so the pack does not overstate it (hub's words, adopted).** It is NOT a
compiler. It is a divergence detector that knows how to recognize its own blindness. It cannot prove
the FE matches a DEPLOYED server and does not try; it proves the FE's copy matches the declaration
in THIS checkout, and refuses to answer at all when it can no longer read that declaration. With arm
(b) it CLOSES F-07; without arm (b) it would only have MOVED it.

### PROCESS RULE, from the hub — the one thing `wireCandidate()` does not fix

> *"Artefato de varredura não é anunciado como fechando nada até um assento que não seja você ter
> lido. Você produz, congela, dispara, e SÓ ENTÃO escreve o que ele fechou."*

The hub located the failure precisely, and not where I had located it. My error was **not** in
reading the Go. It was in **announcing the sweep as a result in the same message that produced it** —
authority prose born and published in one act, with nothing between. In the next message its "two
findings four reviewer rounds missed" had become two false claims of mine. A throwing constructor
fixes fixtures; it cannot fix that. The rule that does is cheaper than any constructor, and it is
the same shape as `P6-DUAL-GATE:` being the hub's line and not mine.

Binding from here, not retroactive.

### The grant, and what measuring it changed

`VinculosPage.test.tsx` granted without reservation — the boundary I respected stopped existing
while I was writing, CHIP-IMPORT-CHAIN having merged at `45b887b3` and been torn down. The hub then
measured the companion rather than leaving me its hypothesis, and **retracted half of its own
concern** when the measurement came back benign: one occurrence, in a comment, supporting no
assertion.

The 29th fixture is corrected through the mechanism: `95/ALTA` + `match_status: "REVIEW"` +
`reasons: []` — three impossibilities at once, and the reason it matters is where it was. It
survived four reviewer rounds AND this chip's own hand sweep by living in a file that sweep did not
own. A checklist has a scope; the fixture was outside it.

**One correction to my own first draft of that work, caught by measuring instead of trusting the
handover.** I wrote a comment saying the Importação section "is no longer on this page and the
component file is gone". That is true on `main` and **false in this worktree**: this branch is based
at `5441fe18`, `VinculosPage.tsx:159` still renders `<ImportacaoSection />` here, and the component
and its test both still exist. The hub's fact was correct and correctly scoped to the integrated
tree; I restated it as if it described mine. Corrected in place — the mock is load-bearing HERE, and
the comment is removed because it describes a seam whose removal is already decided and simply has
not reached this checkout, which is a different and honest reason.

Checked while there, because a deleted component with a surviving test is a live merge hazard:
local `main` deleted **both** `ImportacaoSection.tsx` and `ImportacaoSection.test.tsx`, so this
merge leaves no orphan test. (`origin/main` still shows both — it is stale; nothing has been
pushed.)

**FINDING (lane, for hub ratification).** One cold run of `wireFixtures.guard.test.ts` blew the 5s
`testTimeout` reading two Go files, while the very next test read the same two files in 2ms. **The
cause was not established** — three later runs, including a probe built to reproduce it, all came
back in single-digit ms, and the hypothesis about which `expect` form hung was tested and proved
WRONG. The reads were moved to module scope, which is a **mitigation, not a diagnosis**: module-load
time is not charged to `testTimeout`, so a slow cold read can no longer turn a passing guard red.

### COUNT RECONCILIATION, run against this chip's own pack — `evidence/V-count-reconciliation-self.md`

Amendment `0cb6d7e` (profile §11) came out of this chip's round-5 failure, and the hub bound it back
onto the implementer: *"você roda a mesma varredura contra o próprio pack antes de publicar."* Run
at `2e5331b6`, before the freeze. Both counts printed, as the amendment requires:

```
POPULACAO=51  EXTRACAO=36  REFS=15  soma_confere=True      <- WRONG, corrected below
POPULACAO=53  EXTRACAO=36  REFS=17  soma_confere=True      <- occurrences, not lines
```

**The first line is short by two, and the hub's executor seat found it rather than accepting it.**
`grep -c` counts LINES, not matches; the claim was a count of fixture SITES. Two sites share a line
in each of `BatchPreviewModal.test.tsx:46` and `QueueTab.test.tsx:656`, both `approvals: [{…}, {…}]`.
Both extra occurrences land in the batch-ref residual, so extraction is unchanged at 36.

`soma_confere=True` printed on BOTH runs, and that is the finding, not a footnote: consistent
arithmetic over the wrong population is still consistent. A reconciliation can close and remain
short by exactly the amount the instrument cannot see. Arm C was added for it — a raw fixture with
two occurrences on ONE line — and run against both instruments on the same mutated tree: the old
one reports 54, the corrected one 57, and the old one still says `soma_confere=True` while missing
all three.

Population anchored loosely — `candidate_id:` across **every** `*.test.ts*` in `apps/web/src`, not
only under `pages/vinculos/`, because the 29th fixture survived four rounds by living outside the
sweep's scope. Extraction = candidate objects built through the throwing constructor. The 15-site
residual is named, not slack: batch payload refs, a different union with a different producer.

**The must-fail found a defect in this sweep, too.** First run of both arms moved the population
51→53 and printed **nothing**: the classifier bucketed both injected candidates as batch refs,
because its pattern was `status:` and `match_status:` contains `status:`. Identical in form to round
5's `grep -oh 'anchor: "[a-z_]*"'`, which dropped capitals and accents and reported 1 violation where
there were 5. Word boundary applied; both arms then red — arm A (raw candidate in a file the sweep
names) and arm B (raw candidate in a file no sweep ever named). A second, smaller one: a 6-line
lookback window marked `QueueTab.test.tsx:633-634` as *built* because a real `candidate({…})` sits
four lines above. **Having two counts is what surfaced both; one number would have been reported as
fact.** Run on a scratchpad copy, diffed back against the tree afterwards — the tree was never dirtied.

**REPORT — the residual named a second union with no mechanism.** `ProductLinkBatchPreviewItem.status:
"OK" | "FAILED"` (`packages/sdk-runtime/src/index.ts:1160-1164`; `cause` beside it is free-form
`string`, so no typeset can constrain it). `BatchPreviewModal.tsx` enumerates it by string literal in
three places that disagree under drift: `:44` builds the **apply payload**, `:46` feeds the counter at
`:87`, `:94` renders the row. A third member would be **not applied**, **counted in neither** chip —
so `:84`/`:87` stop summing to `items.length` — and yet **rendered as a failure** by the ternary's
else branch at `:97`. Three inconsistent readings of one item, all silent, all type-correct. Not a
defect today: the union has exactly two members and the partition is total. It is the drift exposure
`driftCandidate` exists for, in a union the mechanism never covered. **Not fixed here** — one freeze
per round, and widening it to a component outside the round's delta is the chip choosing its own scope.

### THE SENTINEL HAD THE SHAPE IT NAMES — hub executor seat, `evidence/V-sentinel-substring-blindness.md`

The hub took the third seat (profile §11, executor) in its own detached worktree at `2e5331b6`,
confirmed the lanes independently (12 / 0-in-vinculos, 64 / 531), proved the merge clean
(`main` @ `0cb6d7e` + `2e5331b6` → 0 conflicts, merged tree 67 files / 544 tests green), and found
one defect this chip had not: the guard's sentinel asserted `.toContain("var knownIdentityAnchors")`,
and `var knownIdentityAnchorsXX` **satisfies that**. Substring containment is not symbol presence.
A rename that appends left the sentinel green while the extraction below it matched nothing.

Verified against source before acting, then fixed red-first. Count and mutation declared together,
per §11 — the pack's `RED, 3/3` was true only of mutating BOTH seams:

| mutation | before | after | sentinel |
|---|---|---|---|
| suffix rename, port only | **1 failed / 2 passed** | **2 failed / 1 passed** | was ✓ **blind** → ✗ |
| total rename, port only | 2 / 1 | 2 / 1 | ✗ |
| mercado_livre seam only | — | 2 / 1 | ✗ |
| both seams | 3 failed | 3 failed | ✗ — this is the `3/3` |

`git diff --stat apps/server_core` = 0 lines after every arm. Suffix and total rename now agree,
which is the fix: they are one class, and the guard used to answer differently for them. The
extraction pattern moved INTO `GO_SEAM` as `extract`, and the sentinel asserts with that same
object — sentinel and extraction cannot drift apart by construction. A tighter second pattern would
have gone red too, and left two patterns free to diverge again.

Reconciled, both counts printed: 7 regex literals matched, **2 are false positives** — `/../` and
`/../s`, from the path STRING on line 37. Real population 5: two seam locators in `GO_SEAM.extract`
(3 uses), three secondary extractions inside an already-located block, each guarded; zero unguarded.
The false positives are the class again inside the tool measuring the class, and they are printed
rather than quietly subtracted.

Same category, same hour: the first run of all four arms used `grep -E "^ *(Test Files|Tests)"` and
returned **empty with exit 0** — vitest emits an ANSI escape before the leading space, so `^ *` never
matches. The hub's executor seat reported hitting the identical trap independently, and said so
first. Silence read as clean, one hour after ratifying the amendment that names it.

Merge-base corrected while there: `git merge-base main HEAD` = `bcab8269`, not `5441fe18` (an
ancestor of it). The comment at `VinculosPage.test.tsx:38` cited the wrong point; fixed in place.
The hub also confirmed independently that local `main` deleted BOTH `ImportacaoSection.tsx` and its
test, so this merge leaves no orphan — the measurement this chip reported was right.

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

**Lanes, final (after round 5).** `tsc` **15 → 12, and 3 → 0 under `pages/vinculos/`**
(`evidence/L0-tsc-baseline.txt` → `evidence/L0-tsc-after-round5.txt`); vitest 62/511 →
**64/531** (`evidence/L1-vitest-after-round5.txt`). The +20 over the baseline are the regression
tests the five gate rounds forced; the new file is `wireFixtures.guard.test.ts`, the Go-source guard
that closes F-07 — its third test is the arm the hub's ruling required, the one that fails when the
guard can no longer see its own subject.

> **CORRECTION (R-24), made at the source rather than annotated.** An earlier version of this line
> read *"the same 12 as the baseline; this chip neither added nor removed one."* That was FALSE.
> `evidence/L0-tsc-baseline.txt` carries **15** errors, and the 3 under `pages/vinculos/` are inside
> it — the three `INCOMPARABLE` ones this chip was dispatched to fix. Baseline 15/3 → HEAD 12/0:
> **this chip removed 3.** The hub's executor seat caught it and also measured `main` at `4852649d`,
> which still carries all 15 including those 3.
>
> The error understated this chip's own contribution, which does not make it less false. It is
> recorded because of what it would have cost operationally: the hub's post-merge check must expect
> **15 → 12**, not 12 → 12. Had "12 = 12" been accepted, a 15 read on integrated `main` would have
> been misread as a regression introduced by this chip.

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

---

## ROUND 6 — the mechanism was refuted, and replaced. Fix commit `26c0ad10`

Both seats REFUTED, blind to each other, on the frozen input dispatched at `2e9e9ce3`
(brief blob `524d8ce2`, unchanged since dispatch). Verdicts captured in
`evidence/REVIEW-gate-vincneutro-sol-rev5.md` (blob `b34d8dbf`, `-o` capture) and
`evidence/REVIEW-gate-vincneutro-opus-rev4.md` (blob `bd360345`, transcribed — that seat has no
`Write`; the banner declares the one mechanical transform applied).

The hub reclassified under a rubric it stated mid-round: **BLOCKING is reserved to behavior,
security, data, contract; prose, count, metadata and citation are REPORT.** Three findings blocked.

### What was actually wrong: a PARTIAL guard under a TOTAL sentence

`wireCandidate` said *"A candidate the backend can actually emit. Throws if it is not one."* It
validated `(confidence, band, status)` plus one hardcoded NO_CANDIDATE shape rule, and it asked
about exactly one `provider_code`. That is not a weaker guard than advertised — it is a guard whose
advertisement was false, which is worse than none: five rounds of reviewers read the sentence and
stopped, because the sentence said the class was closed.

The two lists were the fiction. `state`/`match_input` are decided at `newCandidate`, the triple is
decided afterwards by a scorer, so a fixture assembled from a real state and a real triple can still
be impossible. `ACCEPT` is assigned at ONE non-test site — `generation_service.go:505`, inside
`buildConcordantCandidate` — and that function builds its candidate at `:491` with
`LinkCandidateStateExactSKU` / `LinkCandidateMatchInputSellerSKU`. The ACCEPT branch `:503-505`
touches confidence/band/status only. **Every ACCEPT in the product is `(exact_sku, seller_sku)`.**

### The replacement, one mechanism and not nine patches

`PRODUCIBLE_TUPLES` — 15 rows, one per producing site, each carrying
`(confidence, band, status, state, match_input, from)`, every row read out of the Go one site at a
time rather than derived from the band thresholds in the design doc.

Two consequences that were not point-fixes but deletions:

- **the NO_CANDIDATE special case is gone**, and it was WRONG in the strict direction as well as
  the loose one. It required `state === "unresolved"`. `buildConflictCandidates:340-341` builds
  with `LinkCandidateStateConflict` and then calls `applyUnresolvedScore`, so the backend emits
  `NO_CANDIDATE` with `state === "conflict"`. The old check would have REJECTED a producible
  candidate. Found by this chip's own measurement while building the table; found by neither seat.
- **the `providerCode === "mercado_livre"` gate is gone**, replaced by
  `DECLARED_PROVIDER_CAPABILITIES`. The hub's reading of that gate is the finding, kept verbatim
  because the inversion is the whole lesson: *"o guard checa exatamente o provider onde ele não pode
  errar, e libera todos os providers cuja producibilidade é desconhecida."* A provider with no row
  is not producible at all — `resolveIdentityAnchors:149-169` ABORTS generation when a declaration
  does not resolve, so an undeclared provider emits no candidate rather than a defaulted one.

Also generalized: the unsupplied-anchor rule is now `unavailableDetail(anchor)`, mirroring
`fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor)` at `:704`. The previous form was a
`marca`-only constant, so it could only ever check one anchor of one provider.

`reasonSideLabel` (`QueueRow.tsx:230`) gains `?? reason.side`. It was the only member of its family
— `bandLabels ?? band` (`:60`), `directionLabels ?? direction` (`:292`), `anchorShortLabels ?? anchor`
(`:300`), the ranking map `?? UNKNOWN_DIRECTION_RANK` (`:356`) — that discarded an unknown wire
value instead of showing it verbatim, so a `side` shipped ahead of an SDK regen vanished from the
chip entirely.

### A TENTH impossible fixture, found by the table and by nobody in six rounds

`QueueTab.test.tsx` `cand_3` overrode `state: "conflict"` and INHERITED `match_input: "title"` from
`defaultCandidate()`. No conflict site can pair those: `buildConflictCandidates:329-333` and
`buildCollisionCandidates:360-362` set `match_input` to `seller_sku` or `ean`, because the conflict
IS between those two anchors. Six rounds of reading did not see it; the table saw it on first run.
That is the argument for the table over the checklist, stated as a measurement rather than a hope.

### A rule CONSIDERED AND REJECTED

"A supplied anchor is never UNAVAILABLE" reads as an obvious tightening from `:703-704`
(`classifyProviderIdentityAnchor` returns UNAVAILABLE only under `!anchor.Supplied`). It is **false**.
`missingMatchedAnchorReason:642` sets `Direction = Unavailable` on its `default:` branch — supplied
anchor, both sides carrying a value, and it still matched nothing. Adding the rule would have been
the exact class this chip keeps finding: an instrument wider than the fact it measures. Recorded
because the near-miss is the evidence, and because the rule was rejected by reading the third
UNAVAILABLE site rather than by stopping at the first two.

### Executed evidence — counted, not tailed

| lane | result | artifact |
|---|---|---|
| MUST-FAIL, before the fix | **3 failed \| 22 skipped (25)** | `evidence/V-mustfail-round6-RED.log` (117 lines, full) |
| MUST-FAIL, after the fix | **3 passed \| 22 skipped (25)** | `evidence/V-mustfail-round6-GREEN.log` |
| vitest, whole web suite | **64 files / 534 tests passed**, exit 0 | `evidence/V-vitest-round6-GREEN.log` |
| vitest, the run that found the tenth | 1 failed / 533 passed — `no producing site emits (20, BAIXA, REVIEW, state=conflict, match_input=title)` | `evidence/V-vitest-round6-fix.log` |
| tsc | **12 errors, 0 under `pages/vinculos/`** | `evidence/V-tsc-round6.log` |
| sentinel arm (suffix rename of `knownIdentityAnchors`, port only) | **2 failed \| 3 passed (5)**, sentinel ✗ | `evidence/V-sentinel-mustfail-round6.log` |

All six measured on the tree committed as **`26c0ad10`** (§11, amendment `25716bdb`). The sentinel
arm moved from the hub executor's `2 failed / 1 passed (3)` to `2 failed / 3 passed (5)` by exactly
the two tests this round added, both of which are Go-independent — the residual is named, not slack.
The mutated Go file was restored from a copy held OUTSIDE the repo and proved identical with
`git diff --quiet HEAD -- <file>`; no `checkout`, `reset` or `stash` was used.

The three must-fails, verbatim from the guard and the tab:

1. `MUST-FAIL 1 — ACCEPT exists only as (exact_sku, seller_sku)` — the `(exact_ean, ACCEPT, ean)`
   fixture that survived five rounds now throws `no producing site emits`.
2. `MUST-FAIL 2 — a provider with no capability declaration is not producible` —
   `wireCandidate({ provider_code: "shopee" })` throws, and the message routes to `driftCandidate`.
3. `MUST-FAIL 3 — an INCOMPARABLE 'side' outside the SDK union renders VERBATIM` — the chip shows
   `? SKU (fornecedor)` instead of swallowing the value.

### Fixtures converted

Nine fixtures the mechanism now rejects were corrected rather than exempted, plus the tenth above:
two impossible ACCEPTs (`state`/`match_input` corrected; the `identificado-por` assertion is
unaffected because `statusDecidingAnchors.ACCEPT` at `QueueRow.tsx:175` ignores `match_input` and
names both anchors by definition of corroboration), and seven undeclared-provider fixtures moved to
`driftCandidate(NO_DECLARATION_HERE, …)` — case 2, which is what they always were.

### Pack corrections filed this round (R-24, at source)

- `V-count-reconciliation-self.md` — the superseded `POPULACAO=51` line now carries an INLINE
  supersession marker. A correction that only reaches a reader who started at the top is ordering,
  not correction: `grep POPULACAO=` landed on a wrong number carrying its own `soma_confere=True`.
- the same file's published commands were elided (`…`) and read the working tree, so they could not
  reproduce the numbers printed beside them — run as printed at the fix tip the old one yields 52,
  not 51. Both are now `git grep` forms anchored to `2e5331b6`, and both were run verbatim and
  reproduce `51` and `53`. Re-measured across every tip the number has been quoted at:
  `2e5331b6 → 53`, `7b5c18eb → 53`, `2e9e9ce3 → 53`, `26c0ad10 → 54` (the `+1` is
  `cand_drift_side`, MUST-FAIL 3's fixture).

### Still REPORT, not fixed here

- **F5** — `wireFixtures.ts` cited `(:700-702)` / `(:706-708)` for `classifyProviderIdentityAnchor`;
  the real sites are `:704` (UNAVAILABLE) and `:711`/`:715`/`:723`/`:726`/`:728` (the INCOMPARABLE
  branches). Corrected in the rewritten doc block, since that block was being replaced anyway.
- **the freeze violation is mine** — `EVIDENCE.md` was edited after dispatch, and the Opus seat read
  the transient clobber window. "Frozen" was a claim, not a state. The correct form is to have
  nothing to remember: no pack writes between dispatch and verdict, ledger commit afterwards.
- `BatchPreviewModal.tsx:44/46/94` — `ProductLinkBatchPreviewItem.status` enumerated by string
  literal in three places that disagree under drift. Not a defect today (the union is total at two
  members); the union the mechanism never covered and never claimed to.
- `scripts/harness/pack-measure.sh` (`c66ea7c7`) is not on this branch, so the reconciliation
  commands above are hand-written rather than generated. They are SHA-anchored and were run
  verbatim, which is the property that matters; replacing them with the script is the hub's merge.
