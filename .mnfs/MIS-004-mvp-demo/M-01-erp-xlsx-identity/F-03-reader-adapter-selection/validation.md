# F-03 reader-adapter-selection — validation evidence

Feature = serve the internal_read Reader from EITHER Oracle OR an xlsx snapshot, selected by config,
with proven substitutability. Maps to validation-contract **C04** (Reader port semantics), **C05**
(Oracle path behavior-preserved), and the source-selection seam of **C06**.

Branch `chip/m01-erp-xlsx-identity` · base `59d0e62f`. Failing-test-first, independent review per slice.

| Slice | What | Dispatch | Commit | Review verdict |
|---|---|---|---|---|
| F03-S1 | erp_import→internal_read xlsx Reader adapter (reader.go + unit + integration tests) | D13 (sol/low) | `e314c74a` | PASS — 1🔴 (non-numeric CODPROD hard-fail) fixed pre-commit → skip-continue mirrors catalogPage |
| F03-S2 | composition source-selection + erp_import wiring + governance entry | D15 (luna/high) | `45481943` | PASS — 0🔴 0🟡; oracle path byte-preserved, xlsx cache→timing order, write path always-registered, tenant via cfg only |
| F03-S3 | xlsx/oracle substitutability contract test + 2 fixtures (example-erp 55 rows, identity-rejections 4 rows) | D16 (luna/high) | `4da666fb` | 2🔴 test-theater (oracle fake compared xlsx-to-itself) fixed pre-commit → divergence model |

Plus dual-gate hardening on F03-S3's tests:
| Gate fix | What | Commit |
|---|---|---|
| R3 | persist assertion was tautology (`persistedSnapshot` truthiness) → explicit `persistCalled bool` | `db91f385` |
| R4 | integration contention test held wrong advisory-lock key → held production key `'erp_import:'\|\|tenant` | `db91f385` |

## Contract criteria satisfied
- **C04 Reader port:** cost carries as-of source time; reserved present ⇒ available = físico − reservado;
  reserved ABSENT ⇒ available UNKNOWN (nil), never físico-as-available; unsupported queries
  (GetCurrentPrice / GetSalesHistory) return a TYPED `unsupported_query` outcome (ADR-17, never zero/empty).
  Pinned in source_contract_test.go as xlsx-vs-oracle-shaped-fake divergence/parity assertions.
- **C05 Oracle path intact:** composition source==oracle when erpSource empty; sankhya/profitability/
  oracleDB.Close paths byte-preserved through the source-selection restructure (D15R 8/8 verified).
- **Substitutability:** xlsx Reader and an independent hand-written oracle-shaped Reader (implements
  internal_read/ports.Reader 6 methods + CatalogPageReader) return the same shape for the example fixture;
  unavailable (no completed snapshot) → dual-matchable `ErrNoErpSnapshot` / `ReadErrorSourceUnavailable`.

## Notes / carried findings
- Contract test is HERMETIC (no integration tag, no DB, no server); fixtures hand-rolled via archive/zip
  (NO excelize — grep clean, DEP-GRANT-01 honored); byte-stable (deterministic clock/IDs).
- **GATE-A 🟡:** listingCostReader over nil oracleDB in xlsx mode — confirm honest-degrade-not-panic (hub/QA).
- **GATE-B Y1:** Reader port surfaces sellable stock but not physical stock independently when reserved
  absent; C04 wording "físico segue consultável como físico" (contract:82) may want a port-level physical
  read. Snapshot retains físico; port doesn't expose it. Hub triage.

**Status: COMPLETE** (3/3 slices merged + dual-gate R3/R4 fixed). Aggregate ladder green @ db91f385.
