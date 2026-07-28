# CHIP-ANCHORS-3 — vereditos do round 5, colados VERBATIM

Entrada dos dois assentos: `code-diff.patch` (957 linhas, sha256 `559c15f5…4da546f3`) + `GATE-P6.md`.
Nenhum dos dois recebeu o pack. Cegos entre si.
Alvo `main @ 678a6d51` · chip `chip/anchors-3 @ f8363839` (código idêntico em `617fd41f`).

**Resultado do hub:** REFUTED. Um BLOQUEANTE, achado pelos DOIS assentos de forma independente na
mesma linha, com cenários diferentes e complementares.

---

## Assento A — Claude Opus frio (`harness:gate-reviewer`, read-only por construção)

VERDICT: APROVA

No finding in the delta names a wrong observable that survives the discriminator. Five REPORTs, most severe first. Line numbers cited as `target-file` + the line in `code-diff.patch` you gave me (so every citation is checkable in both trees).

| Sev | Crit | file:line | Defect (1 sentence) | Failure scenario (inputs/state → wrong output) |
|---|---|---|---|---|
| REPORT (borderline) | CORR-2 | `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go`, join rewritten at `code-diff.patch:249` and `:280-281` | `ltrim(x,'0')` is a *lossy* canonical form: it merges CODPRODs that differ only by leading zeros, and merges `'0'`/`'000'`/`''` into the same key. | One import holds two distinct rows `codprod='101'` and `codprod='00101'` (both valid, both distinct PK rows), one link `internal_product_id=101`: post-patch **both** rows satisfy the join, `SELECT DISTINCT products.codprod` returns 2 → `vinculados=2` of `importados=2`; the *same* input returned the correct 1 pre-patch. Symmetric on `enfileirados`: one cursor entry `"101"` matches both rows. The comment at `:243-244` justifies text-to-text against an **unguarded** cast only — `products.codprod ~ '^[0-9]+$' AND products.codprod::bigint = links.internal_product_id` is exact, non-lossy and cannot raise 22P02. |
| REPORT | CORR-1 | `generation_service.go` `identityAnchorValues` (`code-diff.patch:539-543`) + `classifyProviderIdentityAnchor` tree `generation_service.go:719-721` | The ERP `seller_sku` value is now non-empty for **every** product, so the `listingValue != "" && productValue != ""` → `emit=false` branch fires universally; on the one production path that does not seed a `seller_sku` reason the anchor vanishes from `reasons[]`. | `applyCollisionScore` on the ean anchor (`generation_service.go:371`, seeds only that one anchor): listing carries a seller_sku, EAN collided on 2+ ERP products → candidate row now has **no** `seller_sku` reason at all. This is CORR-1's own defect output #2 ("a âncora some de `reasons[]`"), not removed but generalized from "only when refforn is filled" to "always". REPORT, not BLOCKING, because the identical silence already exists for `ean` on that path when both sides carry an EAN — the change aligns `seller_sku` with the module's pre-existing design, and the line that disappears was itself false (`"produto ERP sem seller_sku cadastrado"`, `:728`). |
| REPORT | CORR-1 | `applySingleAnchorScore` ExactEAN tree `generation_service.go:538-547` + `missingMatchedAnchorReason:641-643`; test pinned at `code-diff.patch:848-862` | Reason direction flips `INCOMPARABLE/erp` → `UNAVAILABLE/side=""` for "both sides have a value and they disagree", reusing the direction that means "provider não fornece a âncora" (`:704`). | EAN-only candidate whose listing carries an unmatched seller_sku → stored `link_candidates.reasons` and the /vinculos chip now read `seller_sku UNAVAILABLE` with detail `"sem CODPROD para corroborar o EAN: o seller_sku do anúncio não casa nenhum produto"`, i.e. "unavailable" for an anchor that is present on both sides. Chip names it as the open A2-R1 AGAINST branch and defers it to the operator, which is a legitimate scope call — recorded, not held. |
| REPORT | CORR-6 | `buildConcordantCandidate` (`code-diff.patch:517-520`), test (`code-diff.patch:730-764`) | The nil guard degrades into full corroboration — `seller_sku` FOR + `ean` FOR, confidence 95, ALTA, **ACCEPT** — over a zeroed `ProductCandidate`, and the new test asserts that as expected. | If the path ever became reachable, `autoApprovals` (`generation_service.go:236-246`) selects on `MatchStatus == ACCEPT` alone, so it would hand an auto-approvable candidate with `InternalProductID == nil` to the approver. Not blocking: I verified the sole production call site is `generation_service.go:303`, which passes `&product` (address of a local, never nil) — so nobody observes it today. The finding is that the replacement *pins* the bad degrade rather than degrading like `applyUnresolvedScore:620-628` (0 / NO_CANDIDATE). |
| REPORT | CORR-3 | `erp_import/domain/import.go` `IsValid` (`code-diff.patch:299-316`) | The guard is stricter than the `uuid` column it protects: Postgres also accepts brace-wrapped and dash-less 32-hex uuid text, which now gets 400. | `GET /erp/imports/11111111111111111111111111111111` used to resolve, now 400. No caller in this repo is shown to send that shape (ids come from the API's own list response), so hygiene only. |

## Criteria I consider DISCHARGED by the diff

- **CORR-1** — `identityAnchorValues` case `"seller_sku"` now takes the ERP side from `canonicalProductID(*product)` and never from `ReferenceCode` (`code-diff.patch:539-543`); I verified `canonicalProductID` reads `InternalProductID` only (`generation_service.go:464-469`) and that `InternalProductID` *is* the canonical CODPROD (`internal_read/domain/canonical_identity.go:5-8`, `internal_product.go:4-7`).
  **The named trap is handled honestly.** The five surviving table fixtures now carry `canonicalIDPtr(90x)` (`code-diff.patch:582,590,603,618,626`), and the deleted case `"exact EAN ERP seller SKU empty"` was deleted rather than re-fixtured *for a reason that holds in the code*: with a non-nil product that has a canonical id, `productValue` is never empty, so `missingMatchedAnchorReason:638-640` can never yield `side=erp` for `seller_sku` — the state is genuinely unreachable, not merely untested. The reachable `side=erp` case (nil product, unresolved path, `generation_service.go:216`) is covered by the new end-to-end test at `code-diff.patch:792`. The new `TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference` drives `GenerateLinkCandidates` through `findProducts` (so the fixture is the production shape) and its "no refforn" case genuinely fails pre-fix (`INCOMPARABLE/erp`) and passes post-fix (`UNAVAILABLE`) — a real must-fail, not a vacuous one.
- **CORR-2** — canonicalization applied on **both** sides of both joins; the demanded regression fixture exists verbatim (`codprod='00101'`, `internal_product_id=101`, expected `Vinculados=1`, `code-diff.patch:72,79,88`) and is discriminating (pre-fix that query yields 0). The sibling counter is canonicalized in the same change with its own fixture in both directions (`'00101'` vs cursor `"101"`, `'102'` vs cursor `"00102"`, `code-diff.patch:132-133,150,159`), so the two numbers on `/importacoes` cannot disagree. See the ltrim REPORT above for the residual.
- **CORR-3** — one validation site (`readImportID`, `code-diff.patch:336-343`) used by both handlers before the query; I confirmed in the tree that `/erp/imports/{id}` and `/erp/imports/{id}/chain` are the only `{id}` routes registered (`erp_import/transport/http_handler.go:43-44`). Contract moved in the same change: `400` on both paths + `invalid_import_id` in the shared enum (`code-diff.patch:873-900`), `packages/sdk-runtime/src/erpImport.ts` union widened (`:949-954`), and a test asserts both routes carry the 400 (`:922-932`).
- **CORR-4** — the universal is gone, not annotated: the claim is narrowed to "the vocabulary THIS file governs" and the `market` exception is named explicitly. I verified both halves against the tree: `knownIdentityAnchors` is read nowhere outside `marketplace_capability.go:48`, and `market/domain/identity_resolver.go:91,101,180` really does append the anchor `"refforn"` from `product.RefForn` vs the `MODEL` attribute. No new total sentence with a partial implementation.
- **CORR-6** — `jsonb_typeof(...) = 'array'` guards `jsonb_array_elements_text` inside a `CASE` (not a predicate), so a scalar/object/NULL `pending` reads as an empty queue instead of raising 22023 for the tenant (`code-diff.patch:273-277`), with a two-case fixture (`:188-189`). Nil deref removed (`:517-520`); see REPORT above for what it degrades into.
- **Scope** — `apps/web` does not appear anywhere in the 957 lines. Nothing to report there.

## What I could NOT judge from the diff alone

- **Execution.** Whether `go test ./...`, the `//go:build integration` Postgres lane, `vitest` and the OpenAPI/SDK drift check actually pass. I have no Bash by construction; that is the hub's executing seat. I did check for the likeliest silent lane breakage from the CORR-1 change and found none: `TestCase4TitleOnlyYieldsBaixaReview:1475` still holds (its snapshot has no `SellerSKU`, so the `provider` side branch is unchanged), `TestGenerateLinkCandidatesUsesProviderDeclarationForUnavailableReasons:283` asserts only `ean`, `assertProviderDeclaredUnavailableReasons:1146-1152` asserts only `marca`, and no test exercises the collision path's `seller_sku` reason. `TestCase3` was the one affected assertion and the diff updates it.
- **Real Postgres semantics of the two new SQL forms.** `ltrim`/`jsonb_typeof` behaviour is read from the text, not executed; the `-> 'pending'` guard is only fixtured for a cursor that *is* a JSON object — a cursor that is itself a scalar or array is not covered by either new case (`code-diff.patch:188-189`). I believe `->` returns NULL there rather than raising, but I did not run it: **NOT-EVIDENCED**.
- **Markers.** `P6-DUAL-GATE:`, `LIVE-VERIFIED:`/`LIVE-WAIVED-BY-OPERATOR:` and an `EXEMPLO-IO` golden case are pack-level assertions; by this round's design no pack was in my input, so I make no claim about them — I did not see them and I did not fabricate them. The provider-touching question does not arise from this delta (no connector/provider call is added or changed).

Files read for this review: `.mnfs\MIS-006-integracao-fundacao\_hub-gate-anchors-3\GATE-P6.md`, `...\_hub-gate-anchors-3\code-diff.patch`, and for context only `apps\server_core\internal\modules\product_links\application\generation_service.go`, `...\generation_service_test.go`, `...\erp_import\adapters\postgres\query_repository.go`, `...\erp_import\adapters\postgres\chain_query_repository_integration_test.go`, `...\erp_import\transport\http_handler.go`, `...\internal_read\domain\`, `...\market\domain\identity_resolver.go`.

---

## Assento B — GPT-5.6 Sol, effort medium, via codex (`--wait`)

Model dispatched: gpt-5.6-sol, effort medium, --wait (foreground).

VERDICT: REFUTED

| Severity | Criterion | File:line | Defect | Concrete failure scenario |
|---|---|---|---|---|
| BLOCKING | CORR-2 | `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:128` | The queue join strips leading zeroes from every CODPROD, including accepted nonnumeric identifiers, collapsing distinct products. | An import containing `ABC` and `00ABC` with cursor `pending:["ABC"]` returns `enfileirados=2` instead of `1`, because both normalize to `ABC`. **Discriminator:** the operator observes an inflated queue count on the import-chain screen. |
| REPORT | CORR-6 | `apps/server_core/internal/modules/product_links/application/generation_service.go:496` | The nil defense substitutes an empty product and continues to emit 95/ALTA/ACCEPT with positive SKU/EAN reasons and a null product ID. | A direct or future caller passing `comparison.product=nil` receives an auto-approvable candidate for no product. The current production caller always supplies a nonnil product, so no presently reachable user/operator observable was established. |

Discharged:

- CORR-1: uses canonical `InternalProductID`; the new generator-level regression fixture carries ID `101` and reaches the assertion through `findProducts`.
- CORR-3: both routes use the single `readImportID` validator; OpenAPI and SDK declare the new 400/error value.
- CORR-4: the false universal is narrowed and the wired market-module exception is named.
- CORR-6: the JSON array type guard and mechanical nil-dereference prevention are present, subject to the REPORT above.

Not discharged:

- CORR-2: the required `00101` ↔ `101` regression is covered on both sides, but its unconditional normalization introduces the blocking count corruption above.

Could not judge from the diff alone: none.

Paths actually read:

- `.mnfs\MIS-006-integracao-fundacao\_hub-gate-anchors-3\GATE-P6.md`
- `.mnfs\MIS-006-integracao-fundacao\_hub-gate-anchors-3\code-diff.patch`
- `apps\server_core\internal\modules\erp_import\adapters\postgres\query_repository.go`
- `apps\server_core\internal\modules\erp_import\domain\validation.go`
- `apps\server_core\internal\modules\internal_read\domain\seller_sku.go`
- `apps\server_core\internal\modules\product_links\application\generation_service.go`
- Bounded search excerpts from:
  - `apps\server_core\internal\modules\erp_import\adapters\internalread\reader.go`
  - `apps\server_core\internal\modules\erp_import\transport\http_handler_test.go`
  - `apps\server_core\internal\modules\erp_import\domain\import.go`
  - `apps\server_core\internal\modules\erp_import\adapters\postgres\import_counts_test.go`
  - `apps\server_core\internal\modules\erp_import\adapters\postgres\mirror_repository_integration_test.go`
- `C:\Users\leandro.theodoro\.codex\skills\mnfs-codex-router\SKILL.md`
- `C:\Users\leandro.theodoro\.codex\skills\mnfs-codex-router\references\command-map.md`

---

## Assento EXECUTOR do hub — lane, medida em checkout do hub

```
checkout: .claude/worktrees/hub-exec-anchors3   tip: f8363839 (detached)
BUILD_EXIT=0   VET_EXIT=0
go test ./internal/modules/product_links/... ./internal/modules/erp_import/... -count=1 -v
  199 funções · 304 com subtestes · PASS 304 · FAIL 0 · SKIP 0 · 13 pacotes ok
```

---

## Reclassificação do HUB, verificada por string

O achado do `ltrim` é **BLOQUEANTE**, e a divergência de classe entre os dois assentos não é empate:
o assento A o classificou como REPORT *borderline* e mesmo assim escreveu o cenário completo, com
número errado e tela nomeada. O discriminador responde: o operador lê `vinculados` / `enfileirados`
inflado em `/importacoes` — a MESMA tela e o MESMO número que o CORR-2 existe para consertar, com o
sinal trocado de subcontagem para supercontagem.

Fatos verificados pelo hub, por string, antes da reclassificação:

- `erp_import_products.codprod` é `TEXT NOT NULL` e parte da PK `(tenant_id, protocol_id, codprod)`
  — `migrations/0046_create_erp_import_products.sql:4,13`. **Não há CHECK numérico.** Logo `'101'` e
  `'00101'` coexistem como linhas distintas, e `'ABC'` / `'00ABC'` também.
- `IsValidCodprod` (`internal_read/domain/seller_sku.go:19-35`) rejeita não-dígito — **mas governa o
  `seller_sku` suprido pelo provider**, não a coluna do ERP: os dois sítios de chamada são
  `internal_read/adapters/oracle/reader.go:457` e `erp_import/adapters/internalread/reader.go:475`.
  O comentário do chip (*"codprod is not guaranteed numeric"*) está **certo**; a hipótese contrária
  do hub estava errada e foi descartada por leitura, não por argumento.
- Portanto os dois cenários são alcançáveis: o do assento A (`'101'` + `'00101'`, join de
  `vinculados`) sem precisar de não-numérico, e o do assento B (`'ABC'` + `'00ABC'`, join da FILA
  contra o cursor de texto) precisando.

A propriedade exigida não é "canonicalizar os dois lados" — é **comparar de forma EXATA nos dois
lados**. `ltrim` satisfaz a primeira e viola a segunda. O chip cumpriu a instrução do hub
(*"se só um lado for canonicalizado, o operador lê a diferença entre dois números como fila
travada"*) tornando os dois lados perdedores de informação, quando a saída era tornar os dois
exatos. A instrução era minha e não dizia qual das duas; esta rodada fecha isso.
