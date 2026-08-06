# ADR-017 — what the citations assert

**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 1378

ADR-17/ADR-017 is cited ~1378 times across the repo (Go domain/application/adapters/
transport, tests, migrations, the OpenAPI contract, the FE, packages/sdk-runtime, and
`.mnfs`/docs process artifacts) but no `ADR-017*.md` document exists anywhere in the
repo. No `scripts/.runs/` paths appeared in the citation set at all, so the count above
is the true repo-wide total. The citations are highly self-consistent: the same core
rule, restated per subsystem, plus a small number of named corollaries/extensions that
citations themselves call out as "the other side" or "a third state ADR-17 does not
admit."

## Assertion A1 — Unknown operational facts never become a fabricated zero/default
- Citations: by far the largest group (majority of the ~1378; present in nearly every
  file that cites the tag)
- Verbatim (most common phrasing): "nil means unknown/unset, never a zero Money (ADR-17)"
- Anchors: apps/server_core/internal/modules/pricing/domain/decimal.go:10 ; apps/server_core/internal/modules/pricing/domain/calcprofile.go:37 ; apps/server_core/internal/modules/pricing/domain/tariff.go:22 ; apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:403-404 ; apps/server_core/internal/modules/erp_import/domain/import.go:22 ; apps/web/src/pages/precos/DecompositionPanel.tsx:49 ; contracts/api/marketplace-central.openapi.yaml:5810 ; docs/HARNESS-PROFILE.md:246

## Assertion A2 — The rule is two-sided: a genuinely-known zero must be written/rendered as literal `0`, never suppressed as "unknown"
- Citations: ~6 explicit ("other side of ADR-17") + the golden-test lines that assert an
  explicit 0 in the ST branch
- Verbatim: "Contrapartida do outro lado do ADR-17: **zero conhecido é `0`**. Suprimir um
  zero verdadeiro é o mesmo defeito na direção oposta." / "ADR-17 has two sides: a
  **known** zero must be written as `0`, and writing it as unknown is the same defect
  mirrored."
- Anchors: docs/METODO-DE-REVISAO.md:135-136 ; docs/engineering/defect-class-catalog.md:436-438 ; apps/server_core/internal/modules/pricing/domain/icms_test.go:108-110 ; apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:404 ("0 is itself a valid, distinct grupo_icms value")

## Assertion A3 — Every unresolved/unknown component must be named, never silently dropped
- Citations: ~15+ (the FCP bug is cited repeatedly as the canonical example of getting
  this wrong)
- Verbatim: "toda vez que o fato-fonte falta, \"fcp\" entra em Unknown junto de
  icms_saida/difal/pis_cofins" / "a third state ADR-17 does not admit" (nil-but-unnamed)
- Anchors: apps/server_core/internal/modules/pricing/domain/icms.go:140-151 ; apps/server_core/internal/modules/pricing/domain/icms_test.go:106-110 ; apps/server_core/internal/modules/pricing/domain/icms_erp_golden_test.go:317-321 ; apps/server_core/internal/modules/pricing/domain/decompose.go:87 ; apps/server_core/internal/modules/product_links/domain/link_candidate.go:67-68

## Assertion A4 — "Does not apply" (structural absence) is distinct from "unknown" and must not enter the blocking/Desconhecidos set
- Citations: ~6
- Verbatim: "not an ADR-17 unknown — it never enters ComponentesDesconhecidos and never
  blocks MargemValor" / "a structural does-not-apply, not a missing fact"
- Anchors: apps/server_core/internal/modules/pricing/domain/decompose.go:99-110 ; apps/server_core/internal/modules/pricing/domain/decompose_imposto_legacy_test.go:53-55 ; apps/server_core/internal/modules/pricing/domain/decompose_icms_golden_test.go:119-123

## Assertion A5 — No silent fallback between data sources; the wrong source's data must never stand in for the missing one
- Citations: ~7
- Verbatim: "there is no fallback between sources (ADR-17, fail honest)" / "serving the
  wrong source's data would misrepresent the tenant's catalog (ADR-17)"
- Anchors: apps/server_core/internal/modules/internal_read/adapters/routing/reader.go:4,20-21 ; apps/server_core/internal/modules/internal_read/adapters/routing/stock_batch_reader.go:24-25 ; apps/server_core/internal/modules/internal_read/adapters/routing/stock_batch_reader_test.go:74

## Assertion A6 — No fabricated timestamps; a value without a real source-time must leave the timestamp field nil
- Citations: ~4
- Verbatim: "config has no per-resolve timestamp and ADR-17 forbids fabricating one" /
  "must not build a MarketSignal with a fabricated evidence timestamp (ADR-17)"
- Anchors: apps/server_core/internal/modules/pricing/domain/tariff.go:28-30 ; apps/server_core/internal/modules/pricing/domain/tariff_cotacao_test.go:13-14 ; apps/server_core/internal/modules/listings/application/read_service.go:644-645

## Assertion A7 — Unknown values are excluded from derived booleans/aggregates, never coerced into a value that participates in a comparison
- Citations: ~5
- Verbatim: "custo desconhecido excluded, never counted below cost (ADR-17)" / "unknown
  never counts as below cost"
- Anchors: apps/server_core/internal/modules/listings/domain/signal.go:78-100 ; apps/server_core/internal/modules/listings/domain/signal_test.go:96-98 ; apps/server_core/internal/modules/listings/application/read_service_test.go:134-135

## Assertion A8 — No placeholder/zero-value rows; when a cell/record cannot be resolved, write no row at all
- Citations: ~8
- Verbatim: "NO row at all for this cell — never a zero-value or placeholder row
  (ADR-17: unknown never becomes 0)" / "caller must write no row at all — not a
  zero/placeholder row (ADR-17)"
- Anchors: apps/server_core/internal/modules/internal_read/domain/icms_matrix.go:42-43 ; apps/server_core/internal/modules/internal_read/domain/icms_matrix_test.go:139,233 ; apps/server_core/internal/modules/internal_read/adapters/oracle/icms_matrix.go:79-81 ; apps/server_core/internal/modules/internal_read/adapters/oracle/icms_matrix_test.go:63-65 ; apps/server_core/internal/modules/product_links/application/batch_service.go:208-210

## Assertion A9 — Wire/DTO contract: absent or masked provider data is honest absence (omitted key / null), never a fabricated placeholder value
- Citations: 11 (all in the OpenAPI contract; mirrored by ~10 more in Go DTO
  mapping/tests)
- Verbatim (exact, see full verbatim block below): "masked ≠ error, ADR-17" / "absent
  is honest absence, never a fabricated value (ADR-17)" / "never a fabricated zero
  (ADR-17)"
- Anchors: contracts/api/marketplace-central.openapi.yaml:5721,5727,5791,5797,5802,5810,5832,5843,5856,5913,5920 ; apps/server_core/internal/modules/orders/transport/http_handler_test.go:445-460,466-468

## Assertion A10 — Opaque provider-supplied strings are rendered exactly as received, never mapped/coerced/interpreted
- Citations: ~2 (OpenAPI) + adjacent Go comments about not assuming enum mappings
- Verbatim: "Opaque: never assume/map to a CPF/CNPJ enum — render as-is (ADR-17)."
- Anchors: contracts/api/marketplace-central.openapi.yaml:5839-5843

## Assertion A11 — Attribution/provenance discipline: a decision, rule name, or derived value must reflect what was actually established, never a plausible-looking value nobody verified
- Citations: ~10
- Verbatim: "an anchor matched nothing — a fact nobody established, and one the E10
  CHECK rejects outright (ADR-17 / AC-03)" / "would assert a uniqueness nobody
  established (ADR-17)"
- Anchors: apps/server_core/internal/modules/product_links/application/resolution_service.go:210-212,808-810 ; apps/server_core/internal/modules/product_links/application/batch_service.go:165-210 ; apps/server_core/internal/modules/product_links/application/batch_service_test.go:321-322,336-337,394-395,427-429 ; apps/server_core/internal/modules/product_links/application/decision_trail_test.go:114-116 ; apps/server_core/internal/modules/product_links/domain/product_link_decision.go:42-44

## Assertion A12 — A stale/cached answer presented as current is the same violation as fabricating a fact
- Citations: ~3
- Verbatim: "A stale answer presented as fresh is worse than none (ADR-17)"
- Anchors: apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:53-54 ; apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo_integration_test.go:35-37

## Assertion A13 — Import/ingestion leniency: an optional source field absent from the file is accepted as honest-unknown (warning), never rejected outright nor coerced to 0
- Citations: ~14
- Verbatim: "honest-unknown (nil/empty) rather than rejected (ADR-17)" / "instead of
  fabricating zeros (ADR-17)"
- Anchors: apps/server_core/internal/modules/erp_import/application/import_service.go:94-96 ; apps/server_core/internal/modules/erp_import/domain/validation.go:116-118,223-225 ; apps/server_core/internal/modules/erp_import/adapters/xlsx/parser.go:28-42,120-122 ; apps/server_core/internal/modules/erp_import/adapters/xlsx/parser_raw_export_test.go:96-98,125-126 ; apps/server_core/internal/modules/erp_import/adapters/postgres/import_repository.go:60-62,110-112

## Contradictions
- None found between assertions as stated — the corpus is unusually self-consistent for
  a document that never existed. The only real tension is inherent, not contradictory:
  A1 (never fabricate zero) and A2 (never suppress a real zero) point in opposite
  directions on the same fact and are explicitly documented as two faces of one rule
  (`docs/METODO-DE-REVISAO.md:135-136`, `docs/engineering/defect-class-catalog.md:436-438`).
  Several review-history comments (e.g. `apps/server_core/internal/modules/pricing/domain/icms.go:139-151`,
  `apps/server_core/internal/modules/pricing/domain/solve.go:115-119`) record cases where
  code *violated* A1/A3 and was later fixed — these are corrections, not competing
  interpretations of the rule.
- One near-miss worth flagging: `apps/server_core/internal/modules/internal_read/application/icms_matrix_job.go:47-49`
  and `docs/engineering/defect-class-catalog.md:354-357` both note that ADR-17 "working
  correctly" (honest NULL) can *mask* a totally different defect (an empty upstream
  table) as if it were a resolved-absent value — not a contradiction of the rule, but a
  documented limit on what it detects.

## Exceptions / carve-outs
- **Known-zero carve-out (A2):** a value that is genuinely zero must be persisted/
  rendered as `0`; only a value that could not be determined is nil/NULL/omitted.
  `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:404`: "0 is
  itself a valid, distinct grupo_icms value in Sankhya." `apps/server_core/internal/modules/pricing/domain/icms_test.go:108-110`:
  ICMS de saída is explicitly `0` in the ST branch (the tax is already embedded in
  custo), and that explicit zero is not itself an ADR-17 unknown.
- **Structural "does not apply" carve-out (A4):** `imposto` on a legacy/no-cell path is
  nil but is NOT counted as an ADR-17 unknown and does not block margin — it is a
  different, non-blocking kind of absence (`decompose.go:99-110`).
- **Leniency carve-out (A13):** the client-catalog xlsx import path downgrades an
  absent optional value to a warning (not a rejection) instead of either fabricating a
  0 or hard-failing the row; the strict `ValidateRows`/`RunImport` path is explicitly
  left unchanged and still rejects the same absence (`validation_test.go:139-141`,
  `import_service.go:94-96`).

## Verbatim: the 11 OpenAPI citations
(`contracts/api/marketplace-central.openapi.yaml`)

1. L5721 — "...absent when the address is not yet honestly readable (ML obfuscates the
   buyer address until payment is confirmed — masked ≠ error, ADR-17)."
2. L5727 — "...absent (masked/not yet readable) is honest absence, never a fabricated
   value (ADR-17)."
3. L5791 — "...absent when the provider reports no sub-status, never a fabricated value
   (ADR-17)."
4. L5797 — "...absent (carrier 404 / not yet dispatched) is honest absence, never a
   fabricated carrier (ADR-17)."
5. L5802 — "...absent is honest absence, never a fabricated URL (ADR-17)."
6. L5810 — "Every amount is optional/nullable: absent (nil) means the cost could not be
   honestly sourced, never a fabricated zero (ADR-17)."
7. L5832 — "...absent when the buyer has no billing data (the provider returns 404 as a
   normal \"none\" condition — masked ≠ error, ADR-17), never a fabricated identity."
8. L5843 — "Document type string EXACTLY as the provider returns it
   (identification.type). Opaque: never assume/map to a CPF/CNPJ enum — render as-is
   (ADR-17)."
9. L5856 — "Every field optional: absent (masked/blank) is honest absence, never a
   fabricated value (ADR-17)."
10. L5913 — "...null means the tax picture could not be honestly resolved for at least
    one item (unlinked product, no fiscal snapshot, or an unresolved/ambiguous ICMS
    matrix cell) — never a fabricated 0 (ADR-17)."
11. L5920 — "D-41 PIS/COFINS débito, summed across the order's items. Same
    resolution/unknown rule as icms_saida (ADR-17)."
