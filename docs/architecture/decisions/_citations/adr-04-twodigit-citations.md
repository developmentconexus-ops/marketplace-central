# ADR-04 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 94 (live code: 13, .mnfs/docs only: 81)
**Three-digit `ADR-004` citations:** 20 (1 is the document's own title line `004-integration-catalog-plugin-framework.md:1`; 19 real citations elsewhere, mostly `.mnfs/MIS-001-mercado-livre-operating-cockpit/**` and `docs/superpowers/plans|specs/**`)
**Existing document `004-integration-catalog-plugin-framework.md`:** "Integration Catalog Plugin Framework" — providers self-register catalog definitions/auth-adapter factories/fee syncers; composition root consumes registries instead of owning provider-specific construction (accepted 2026-04-25) — DOES NOT MATCH: none of the four live two-digit assertions (Oracle pool/timeouts, price-evidence-in-market, upsert-merge keep-absent, single-writer ingest) is about plugin self-registration. This is the confirmed example already on record: code comments cite "ADR-04: one writer per shared seam" (see A4 below) while the actual `004-*.md` document is about the provider-plugin framework.

## Assertion A1 — One writer per shared seam: entity ingest is resource-addressed and goes through exactly one write path (`IngestOrder`/`IngestListing`); enumerators/backfill/webhook never persist directly (mission: MIS-007, ratified — dominant, code-enforced)
- Citations: ~58 (live code: 9)
- Verbatim: "IC-06's single write path (ADR-04): every ingest trigger (import enumeration today; backfill/..."
- Anchors:
  - `apps/server_core/internal/modules/orders/application/ingest_service.go:65`
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go:810`
  - `apps/server_core/internal/modules/orders/ports/order_ingestor.go:5`
  - `apps/server_core/internal/composition/root.go:591`
  - `apps/server_core/internal/composition/root.go:822`
  - `apps/server_core/tests/integration/listings_refresh_test.go:36`
  - `.mnfs/MIS-007-ml-sync/mission.md:163`
  - `.mnfs/MIS-007-ml-sync/M-03-orders-shipment-persist/F-02-ingest-order-v1/spec.md:9`
  - `.mnfs/HARNESS-DEBTS.md:420` (downstream reference: "único (ADR-04). Nada no repo impede um segundo caminho de escrita de orders renascer...")

## Assertion A2 — Upsert-merge keep-absent: a product missing from the latest snapshot is flagged `absent_in_last_snapshot`/`stale_since`, never physically deleted (mission: MIS-006, ratified)
- Citations: ~19 (live code: 4)
- Verbatim: "Keep-absent (ADR-04): a rebuild never physically deletes. When a product from..."
- Anchors:
  - `apps/server_core/migrations/0076_products_mirror_active_source.sql:13`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:4`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:114`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:57`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:130`
  - `.mnfs/MIS-006-integracao-fundacao/M-04-sankhya-adapter/milestone.md:148`
  - `docs/design/STORAGE-SCHEMA.md:36`
  - `docs/design/SYSTEM-BLUEPRINT.md:99`

## Assertion A3 — Price evidence lives in the existing `market` module (Snapshot/Signal/ValidatedOffer/Aggregate); `pricing` consumes the contract, never persists its own `market_price` (mission: MIS-004, ratified — Claude and Sol candidates agreed on this one, no churn)
- Citations: 9 (live code: 0)
- Verbatim: "ADR-04 Evidência de preço persiste no módulo `market` (Snapshot/Signal/ValidatedOffer/Aggregate); pricing consome contrato, não persiste market_price próprio"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:88`
  - `.mnfs/MIS-004-mvp-demo/mission.md:156`
  - `.mnfs/MIS-004-mvp-demo/mission.md:157`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:41`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:33`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:11`

## Assertion A4 — Oracle pool sized 12; `http.Server` Read/Write/ReadHeader timeouts; per-route-class context deadlines (15s interactive / 120s batch); batch semaphore of 4 (mission: MIS-002, ratified)
- Citations: 5 (live code: 0)
- Verbatim: "ADR-04 Pool 12; http.Server Read/Write/ReadHeader timeouts; per-route-class context deadlines (15s/120s); batch semaphore 4"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:114`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:166`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:167`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/M-01-foundation-observability/F-01-adapter-hardening-pool-deadlines/feature.md:32`

## Contradictions
- **Four unrelated decisions, one per mission** — same pattern as ADR-01/02/03: single-writer ingest (MIS-007), keep-absent merge (MIS-006), price-evidence-in-market (MIS-004), pool/timeouts (MIS-002).
- **Two numbering schemes for the same live rule, acknowledged in-repo.** `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md:41` explicitly normalizes "`ADR-04` e `ADR-004` são o mesmo" while measuring the collision, and line 53 counts "Esquemas de numeração convivendo | **2** — `ADR-04` e `ADR-004`. Só o de 3 dígitos tem arquivo" — i.e. this plan already flagged that the two-digit `ADR-04` citations point at a document (`004-integration-catalog-plugin-framework.md`) whose actual content has nothing to do with any of A1–A4.
- **`RENUMBERING-REGISTRY.md` already declares the two-digit spelling dead for this number**, without resolving which of the four live meanings inherits the slot: `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md:23` — "`ADR-17` e `ADR-04` são grafias mortas."
- **`013-webhook-is-pointer-never-data.md` cites `ADR-04` lineage for the single-writer rule (A1)** without a cross-reference disclaimer: `docs/architecture/decisions/013-webhook-is-pointer-never-data.md:41,45` — a reader who follows "ADR-04" from a *ratified, currently-numbered* document lands on the wrong three-digit file (`004-integration-catalog-plugin-framework.md`).

## Amendments
None found specific to this number beyond the cross-mission renumbering already covered in Contradictions.
