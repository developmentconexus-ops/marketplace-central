# P5 Sol Decomposition Audit — r01 (VERBATIM)

```yaml
id: P5-SOL-R01
type: planning-review
status: complete
owner: GPT-5.6 Sol (medium) via codex exec OS-process
parent: MIS-004
created: 2026-07-17
lifecycle_scope: support
```

Dispatch: codex exec --model gpt-5.6-sol -c model_reasoning_effort=medium --sandbox read-only
Input manifest: planning-reviews/p5-input-r01.sha256 (top digest 801c29282f14cb8602f9346267bf658db4d9142e57967632f7f2078a9cda829b)
Log: scratchpad agent__p5-sol-r01.log · Last message below, UNEDITED.

---
## Checked paths

Frozen inputs, relative to mission root:

- `planning-reviews/p5-input-r01.sha256`
- `mission.md`
- `planning-reviews/p5-decomposition-passes-r01.md`
- `research/identity-matching-interface-contract.md`
- `research/erp-xlsx-import-interface-contract.md`
- `research/market-evidence-read-interface-contract.md`
- `research/pricing-difal-interface-contract.md`
- `research/fe-shell-seams-interface-contract.md`
- `research/ml-read-ports-interface-contract.md`
- `research/w1-merge-addendum-2026-07-17.md`
- `research/design-screens-2026-07-17.md`
- `research/p1-clarified-decisions-2026-07-17.md`
- `research/repo-baseline-2026-07-17.md`
- `M-01-erp-xlsx-identity/milestone.md`
- `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`
- `M-01-erp-xlsx-identity/F-02-erp-import-module/feature.md`
- `M-01-erp-xlsx-identity/F-03-reader-adapter-selection/feature.md`
- `M-02-price-intel-core/milestone.md`
- `M-02-price-intel-core/F-01-ml-adapter-read-ports/feature.md`
- `M-02-price-intel-core/F-02-market-persistence/feature.md`
- `M-02-price-intel-core/F-03-identity-resolver/feature.md`
- `M-02-price-intel-core/F-04-collect-verdict-api/feature.md`
- `M-03-shell-retheme/milestone.md`
- `M-03-shell-retheme/F-01-theme-tokens-fonts/feature.md`
- `M-03-shell-retheme/F-02-header-nav-routes/feature.md`
- `M-03-shell-retheme/F-03-shared-primitives/feature.md`
- `M-04-vinculos-import-ui/milestone.md`
- `M-04-vinculos-import-ui/F-01-product-links-api-gaps/feature.md`
- `M-04-vinculos-import-ui/F-02-vinculos-screen/feature.md`
- `M-05-anuncios-sinais/milestone.md`
- `M-05-anuncios-sinais/F-01-listings-signals-api/feature.md`
- `M-05-anuncios-sinais/F-02-anuncios-ui-sinais/feature.md`
- `M-07-simulador/milestone.md`
- `M-07-simulador/F-01-pricing-calc-difal-service/feature.md`
- `M-07-simulador/F-02-simulador-ui/feature.md`
- `M-08-pedidos/milestone.md`
- `M-08-pedidos/F-01-orders-projection-api/feature.md`
- `M-08-pedidos/F-02-pedidos-ui/feature.md`
- `M-06-produto-detalhe/milestone.md`
- `M-06-produto-detalhe/F-01-produto-detalhe-page/feature.md`
- `M-09-dashboard-demo/milestone.md`
- `M-09-dashboard-demo/F-01-dashboard-mpc/feature.md`
- `planning-reviews/p3-reconciliation-r01.md`

Repository code-fact verification, relative to workspace root:

- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/modules/mutations/domain/protocol.go`
- `apps/server_core/internal/modules/mutations/application/writer.go`
- `apps/server_core/internal/modules/mutations/application/preview.go`
- `apps/server_core/internal/modules/mutations/application/service.go`
- `apps/server_core/internal/modules/mutations/transport/command_handler.go`
- `apps/server_core/internal/modules/mutations/adapters/listings/selection_resolver.go`
- `apps/server_core/internal/modules/mutations/adapters/connectors/price_writer.go`
- `apps/server_core/internal/modules/mutations/adapters/stub/writer.go`
- `apps/server_core/internal/modules/mutations/ports/linkage.go`
- `apps/server_core/internal/modules/catalog/domain/product.go`
- `apps/server_core/internal/modules/catalog/domain/canonical_product.go`
- `apps/server_core/internal/modules/catalog/adapters/internalread/reader.go`
- `apps/server_core/internal/modules/internal_read/ports/reader.go`
- `apps/server_core/internal/platform/migrate/runner_test.go`
- `apps/web/src/pages/mutations/MutationPreviewModal.tsx`
- `apps/server_core/migrations/` file inventory

## Findings

1. **P5-F-01 — BLOCKING — Check 1: Wave B contradicts a declared code edge.**

   - Cited excerpt: “Wave B (após A): `M-04 ∥ M-05 ∥ M-07 ∥ M-08`” and “`M-07→M-08 (decomposição/DIFAL read IC-04)`.”
   - Defect locus: `mission.md`, `## Parallel Execution Plan`, lines 127–130; reinforced by `M-08-pedidos/milestone.md`, `## Dependencies` and `## Ownership & Concurrency`, lines 38 and 46.
   - Offending token/value: `M-07 ∥ M-08` together with `M-07→M-08`.
   - Defect: M-08 F-01 consumes the to-be-created `Decompose` and `DifalForUF` code from M-07 F-01. “Contract-first” can permit preparatory work, but it does not make the producing and consuming milestones parallel in the write-DAG.
   - **Yes-if:** Wave/feature ordering explicitly places publication of M-07 F-01’s frozen IC-04 ports before the dependent M-08 F-01 work, while retaining only genuinely independent M-08 work in parallel.

2. **P5-F-02 — BLOCKING — Check 1: Two consumed milestone artifacts lack edges.**

   - Cited excerpt: M-06 says it composes APIs from “`M-01/M-02/M-04/M-05`” and consumes “`listings.ts` (by-product + sinais — M-05)”; M-08 says “`M-01 (custo — edge de dado via Reader IC-02)`.”
   - Defect locus: `M-06-produto-detalhe/milestone.md`, `## Outcome`/`## Dependencies`, lines 21 and 35; `M-06.../F-01.../feature.md`, `## Inputs`, line 30; `M-08-pedidos/milestone.md`, `## Dependencies`, line 36; absent from `mission.md` edge list at line 129.
   - Offending token/value: missing `M-05→M-06` and missing `M-01→M-08`.
   - Defect: M-06 directly needs M-05’s SDK/API, and M-08’s demo outcome needs the xlsx-backed cost source created by M-01. The mission edge list names neither dependency.
   - **Yes-if:** both edges are recorded with their forcing artifacts and their timing is reflected in milestone dispatch/close rules.

3. **P5-F-03 — BLOCKING — Check 1: The aggregate M-09 edge creates a false M-06 dependency.**

   - Cited excerpt: “Wave C: `M-06 ∥ M-09`” followed by “`{M-01..M-08}→M-09`.”
   - Defect locus: `mission.md`, `## Parallel Execution Plan`, lines 128–129.
   - Offending token/value: `{M-01..M-08}→M-09`, specifically `M-06→M-09`.
   - Defect: M-09 names orders, listings, product_links, and erp_import as sources; it consumes no M-06 artifact. The broad edge both falsely serializes M-09 behind M-06 and contradicts the declared Wave C parallelism.
   - **Yes-if:** the aggregate edge is replaced by only the actual producing milestones; no M-06→M-09 edge remains unless an M-06-produced artifact is named.

4. **P5-F-04 — BLOCKING — Checks 2, 3, and 5: Every feature-level migration ownership path points to a nonexistent tree.**

   - Cited excerpt: examples include “`db/migrations/0045*–0047*`,” “`db/migrations/0050*–0053*`,” “`db/migrations/0055*–0058*`,” “`db/migrations/0060*–0062*`,” and “`db/migrations/0065*–0067*`.”
   - Defect locus: migration-owning feature `## Ownership` sections in M-01 F-02/F-03, M-02 F-02/F-04, M-04 F-01, M-07 F-01, and M-08 F-01.
   - Offending token/value: `db/migrations/...`.
   - Defect: the manifested baseline states the canonical path is `apps/server_core/migrations/*.sql` (`research/repo-baseline-2026-07-17.md`, line 49), and repository inspection confirms `db/migrations` does not exist. Thus the numeric blocks are allocated, but the actual write-set is not.
   - **Yes-if:** all migration ownership entries bind their assigned prefixes to `apps/server_core/migrations/`, with the existing `apps/server_core/internal/platform/migrate/runner_test.go` seam retained for hub adjudication.

5. **P5-F-05 — BLOCKING — Checks 2, 4, and 5: M-01 F-01 promises a catalog API change while forbidding the only existing SDK contract file.**

   - Cited excerpt: “expor os campos canônicos (`codprod`, `ean`, `refforn`, `marca`, `ncm`) na API de produto do catalog,” but “Forbidden paths: … `sdk-runtime/**`.”
   - Defect locus: `M-01.../F-01.../feature.md`, `## Brief` and `## Ownership`, lines 25 and 54–55.
   - Offending token/value: `sdk-runtime/**` forbidden.
   - Defect: the current catalog contract is defined in `packages/sdk-runtime/src/index.ts:154–175`; it has `manufacturer_reference` and `brand_name` and no `ncm`. No other feature is assigned this catalog SDK update. This violates ADR-12’s OpenAPI+SDK same-commit rule and leaves an orphaned API promise.
   - **Yes-if:** M-01 F-01 owns the matching catalog OpenAPI schema and existing SDK catalog types, or a named upstream feature produces those exact fields before F-01’s consumers.

6. **P5-F-06 — BLOCKING — Checks 4 and 6: F-01 reverses the ratified REFERENCIA/REFFORN semantics.**

   - Cited excerpt: F-01 says “caso contrário REFERENCIA vai para `refforn`”; IC-01 says “`ean` … TGFPRO.REFERENCIA; … inválido ⇒ tratado como ausente + warning” and “`refforn` … TGFPRO.REFFORN.”
   - Defect locus: `M-01.../F-01.../feature.md`, `## Brief`/`## Expected Output`, lines 25, 36, and 38; `research/identity-matching-interface-contract.md`, `## Resources Or Entities`, lines 28–29.
   - Offending token/value: invalid `REFERENCIA → refforn`.
   - Defect: an invalid GTIN in REFERENCIA must become absent-with-warning; it cannot overwrite the independent manufacturer reference from REFFORN.
   - **Yes-if:** invalid REFERENCIA produces `ean:null` plus the contracted warning, while `refforn` is populated only from REFFORN.

7. **P5-F-07 — BLOCKING — Checks 4 and 6: ERP import lifecycle and invalid-EAN behavior diverge from IC-02.**

   - Cited excerpt: F-02 promises “`POST /erp/imports` … ⇒ 202,” statuses `UPLOADED|VALIDATING|COMPLETED|FAILED`, and line rejection `INVALID_EAN`; IC-02 fixes `201`, synchronous processing, statuses `COMPLETED|REJECTED`, and invalid EAN as a warning treated absent.
   - Defect locus: `M-01.../F-02.../feature.md`, `## Expected Output`/`## Inputs/Outputs`, lines 36–44 and 70; `research/erp-xlsx-import-interface-contract.md`, `## File Format`, `## Operations`, and `## Persistence Expectations`, lines 34, 45, and 51.
   - Offending tokens/values: `202`, `UPLOADED`, `VALIDATING`, `FAILED`, rejected-line `INVALID_EAN`.
   - **Yes-if:** the brief uses IC-02’s synchronous `201` response, contracted persisted statuses, and warning-only invalid-EAN treatment exactly.

8. **P5-F-08 — BLOCKING — Checks 4 and 6: Unknown reserved stock is converted into a numeric availability.**

   - Cited excerpt: F-03 says absent ESTOQUE_RESERVADO gives “`disponível = físico + flag reservado_desconhecido: true`”; IC-02 says “unknown (`disponível fica unknown — propaga, não zera`).”
   - Defect locus: `M-01.../F-03.../feature.md`, `## Inputs`, line 29; `research/erp-xlsx-import-interface-contract.md`, `## File Format`, line 33.
   - Offending value: `disponível = físico`.
   - Defect: a flag beside a fabricated numeric availability does not satisfy ADR-17.
   - **Yes-if:** sellable availability itself is nullable/unknown whenever reserved stock is unknown.

9. **P5-F-09 — BLOCKING — Checks 4 and 6: The identity resolver weakens hard-negative REJECT to REVIEW.**

   - Cited excerpt: F-03 says EAN match plus contradictory brand yields “no máximo REVIEW”; IC-01 says any hard negative “⇒ REJECT do candidato mesmo com EAN igual.”
   - Defect locus: `M-02.../F-03.../feature.md`, `## Expected Output`, line 39; `research/identity-matching-interface-contract.md`, `## Matching Gate`, line 37.
   - Offending token/value: `REVIEW` instead of `REJECT`.
   - **Yes-if:** every contracted hard negative deterministically yields `REJECT`, including EAN matches.

10. **P5-F-10 — BLOCKING — Check 6: Market aggregate/verdict query keys diverge from IC-03.**

    - Cited excerpt: F-04 specifies `GET /market/aggregates?product_ids=` and `GET /market/verdicts?product_ids=`; IC-03 specifies `GET /market/aggregates?codprod=` and `GET /market/verdicts?codprod=`.
    - Defect locus: `M-02.../F-04.../feature.md`, `## Expected Output`, line 39; `research/market-evidence-read-interface-contract.md`, `## Operations`, lines 36–37.
    - Offending token/value: `product_ids` versus `codprod`.
    - **Yes-if:** producer and all consumers use IC-03’s exact `codprod` parameter and identifier semantics.

11. **P5-F-11 — BLOCKING — Checks 5 and 8: M-02 defers the collection execution model to feature execution.**

    - Cited excerpt: “execução síncrona-curta ou job local (decisão de spec), resultado consultável.”
    - Defect locus: `M-02.../F-04.../feature.md`, `## Expected Output`, line 35.
    - Offending phrase: `síncrona-curta ou job local (decisão de spec)`.
    - Defect: those alternatives have different persistence, status, and prerequisite requirements. No job status endpoint/table is verified or assigned, and the brief explicitly postpones the choice.
    - **Yes-if:** the brief fixes the IC-03-compatible execution model and names every required result/status resource at planning time.

12. **P5-F-12 — BLOCKING — Checks 4 and 5: M-04’s claimed reuse of `link_apply` is not satisfiable by the existing envelope.**

    - Cited excerpt: M-04 promises “UM protocolo … com N itens, payload por item `{action: "approve_candidate", candidate_id, ...}`,” a preview with “2 OK + 1 FAILED,” and forbids `modules/mutations/**`.
    - Defect locus: `M-04.../F-01.../feature.md`, `## Expected Output`, `## Constraints`, and `## Validation Expectations`, lines 37, 56, and 68; contradicted by:
      - `apps/server_core/internal/modules/mutations/adapters/listings/selection_resolver.go:18–71`
      - `apps/server_core/internal/modules/mutations/application/preview.go:37–79`
      - `apps/server_core/internal/modules/mutations/application/service.go:42–68`
      - `apps/server_core/internal/composition/root.go:568–588`
    - Offending promise: per-item candidate intents through the unchanged `/mutations` envelope.
    - Defect: the current envelope accepts listing IDs and one protocol-wide intent, copies that same intent to every item, and previews existing listings. It cannot create distinct candidate payloads or preview an invalid candidate as one failed item. Additionally, when `MPC_PROVIDER_WRITES_ENABLED` is false, composition replaces the whole writer router—including local linkage—with the generic stub, so approving `link_apply` does not apply local resolutions.
    - **Yes-if:** an already-approved envelope contract supports per-item linkage intents and keeps the local linkage writer active while provider writers remain disabled, with ownership assigned to a feature allowed to change that seam.

13. **P5-F-13 — BLOCKING — Checks 2 and 6: M-05 leaves its cross-module transport seam undecided.**

    - Cited excerpt: “GET signals batch ou port Go interno — spec decide transporte.”
    - Defect locus: `M-05.../F-01.../feature.md`, `## Inputs`, line 31.
    - Offending phrase: `ou ... spec decide`.
    - Defect: M-02 F-04 names an HTTP API, not a Go port; M-05’s write-set and prerequisite differ depending on the choice. The seam is therefore neither frozen nor assigned.
    - **Yes-if:** the brief selects an already-produced IC-03 public boundary, or names the feature that creates the exact Go port before M-05 consumes it.

14. **P5-F-14 — BLOCKING — Check 6: The shell brief omits the ratified Configurações→DIFAL entry.**

    - Cited excerpt: IC-05 fixes the gear menu as “Configurações → seção DIFAL read/drawer; Integrações; Catálogo; Estoque”; F-02 lists “Catálogo, Estoque, Vínculos, Integrações, Protocolos.”
    - Defect locus: `research/fe-shell-seams-interface-contract.md`, `## Nav canônica`, line 31; `M-03.../F-02.../feature.md`, `## Expected Output`, line 36.
    - Offending omission: `Configurações → DIFAL`.
    - **Yes-if:** F-02 propagates the exact IC-05 Configurações/DIFAL access while retaining the other approved reachable routes.

15. **P5-F-15 — BLOCKING — Check 7: A UI feature lacks the required State/Interaction Model section.**

    - Cited excerpt: M-03 F-03 goes directly from “`## Negative Scenarios`” to “`## Constraints`” and “`## Ownership`.”
    - Defect locus: `M-03-shell-retheme/F-03-shared-primitives/feature.md`, after `## Negative Scenarios`, lines 41–52.
    - Offending omission: `## State / Interaction Model` or equivalent.
    - Defect: this is a UI feature defining `MarginChip`, `DataTable`, and `DetailDrawer`, but does not state controlled/uncontrolled state, drawer interaction ownership, or event responsibilities.
    - **Yes-if:** the required section is present and fixes those UI state/interaction responsibilities without adding product scope.

16. **P5-F-16 — BLOCKING — Check 6: IC-03 evidence fields are not propagated to all planned market-price displays.**

    - Cited excerpt: IC-03 requires every market-price display to carry “`source`, `fetched_at`, `n_offers`/`n_sellers` when aggregated, `match_status`.” M-05’s listing shape contains only source/fetched freshness, and M-07 merely says “mostrar faixa vs preço simulado.”
    - Defect locus:
      - `research/market-evidence-read-interface-contract.md`, `## Evidence Fields`, line 41
      - `M-05.../F-01.../feature.md`, `## Expected Output`, line 36
      - `M-05.../F-02.../feature.md`, lines 36 and 39
      - `M-07.../F-02.../feature.md`, `## Expected Output`, line 39
    - Offending omissions: `match_status`, aggregate `n_offers/n_sellers`, and M-07 `source/fetched_at`.
    - **Yes-if:** every M-05 and M-07 market-price rendering receives and visibly presents all evidence fields required by IC-03/P1c, using unknown states rather than omitted or fabricated values.

17. **P5-F-17 — BLOCKING — Checks 4 and 6: M-07 materially drifts from IC-04.**

    - Cited excerpts:
      - F-01: `NO_SOLUTION`; IC-04: `UNREACHABLE_TARGET`.
      - F-01: invalid UF `422 INVALID_UF`; IC-04: `404 UF_NOT_FOUND`.
      - F-01: rate `>30%` rejected; IC-04 permits `0–35%`.
      - F-01: missing initial regime remains unknown; IC-04 fixes `SIMPLES default 4%`.
      - F-01’s `Decompose` output omits `tarifa_full`; IC-04 includes it in CalcProfile and the unique formula.
      - F-01’s `DifalForUF` returns seed/override metadata but omits the contracted version output.
    - Defect locus: `M-07.../F-01.../feature.md`, `## Expected Output`/`## Negative Scenarios`, lines 37–53; `research/pricing-difal-interface-contract.md`, `## Resources Or Entities`/`## Operations`/`## Error Matrix`, lines 25–50.
    - Offending tokens/values: `NO_SOLUTION`, `422 INVALID_UF`, `30%`, no default regime, omitted `tarifa_full`, omitted version.
    - **Yes-if:** F-01 adopts IC-04’s exact defaults, bounds, error/status vocabulary, full decomposition, and port result.

18. **P5-F-18 — BLOCKING — Checks 4 and 5: The promised 27-UF DIFAL seed lacks the internal-rate dataset needed by the contract.**

    - Cited excerpt: F-01 promises “27 UFs seedadas: 12% MG/PR/RJ/RS/SC/SP, 7% demais”; IC-04 defines those as `interestadual_pct` and requires `efetivo_pct = interna_pct − interestadual_pct`.
    - Defect locus: `M-07.../F-01.../feature.md`, `## Expected Output`, line 36; `research/pricing-difal-interface-contract.md`, `## Resources Or Entities`, line 26.
    - Offending omission: the 27 authoritative `interna_pct` seed values or a named frozen source containing them.
    - Defect: the 12%/7% values alone cannot calculate effective DIFAL. Neither the IC, the feature, nor P5 prerequisite evidence supplies the internal-rate dataset.
    - **Yes-if:** the existing approved contract names a frozen 27-UF `interna_pct` dataset/source available before execution, while retaining origin SC, 12%/7%, and the required fiscal disclaimer.

19. **P5-F-19 — BLOCKING — Check 8: Functional Cancelados scope is reintroduced despite explicit deferral.**

    - Cited excerpt: M-08 F-02 promises functional tabs “Todos/Pagos/Enviados/Entregues/Cancelados” and behavior “While tab Cancelados ativa”; M-08 F-01 says cancelled orders are “visível na tab própria.” Mission non-scope says “Cancelados/Devoluções/claims/NF-e produção (MIS-005 M-04),” and the research note says Cancelados/Devoluções must remain an explicit “em breve” stub.
    - Defect locus:
      - `M-08.../F-01.../feature.md`, `## Negative Scenarios`, line 52
      - `M-08.../F-02.../feature.md`, `## Expected Output`, lines 37 and 41
      - `mission.md`, `## Non-Scope`, line 41
      - `research/design-screens-2026-07-17.md`, `### Pedidos`, line 57
    - Offending scope: functional `Cancelados` tab.
    - **Yes-if:** Cancelados remains only the approved disabled/“em breve” stub in MIS-004, with functional handling left to MIS-005.

## Non-findings worth noting

- The frozen input was complete: 42 manifest entries, zero file-hash mismatches, and the normalized manifest-entry digest matched `801c29282f14cb8602f9346267bf658db4d9142e57967632f7f2078a9cda829b`.
- Numeric migration blocks themselves are disjoint and match the audit allocation: M-01 0045–0049, M-02 0050–0054, M-07 0055–0059, M-08 0060–0064, M-04 0065–0069, reserve 0070–0074.
- Shared root/module-governance, SDK barrel, OpenAPI root, route-swap, and migration-count seams are named as hub-adjudicated.
- Other declared feature write overlaps have serial edges: M-01 F-02→F-03, M-02 F-02→F-03→F-04, M-03 F-01→F-02/F-03, and each full-stack API feature→UI feature.
- `link_apply` and `price_update` are real enabled `ProtocolType` values. `link_apply` dispatches to a local linkage writer when that router is wired; `price_update` dispatches to the ML writer only after approval.
- Provider writes are default-off in current code through `MPC_PROVIDER_WRITES_ENABLED`; M-07’s explicit ceiling at `previewed` is additionally consistent with zero live ML writes.
- Confidence bands ≥85 / 50–84 / <50, five-seller minimum, thresholds 18/10, fixed fee 6,50 below 79, origin SC, and the 12%/7% interstate rule are otherwise propagated.
- M-09 correctly distinguishes a real zero with an explicit time window from unavailable source data represented as null plus reason.
- The parser-library uncertainty is explicitly routed through the harness dependency-request protocol rather than silently assuming a package exists.

## Verdict

**NEEDS-FOLD** — blocking findings P5-F-01 through P5-F-19 are listed above.