# ADR-020: Market/competitor data comes only through CollectorPort, never scraping

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision was ratified inside MIS-003
(`operator-cockpit-wireframe-replan`) as its "ADR-14", a different decision than the
MIS-007 "ADR-14" reconstructed separately as process rules in `docs/HARNESS-PROFILE.md`
(reference P-5 of the renumbering registry). No standalone document existed for either
sense. This document reconstructs the MIS-003 sense from its 3 live citations, harvested
at `docs/architecture/decisions/_citations/adr-014-citations.md`, Assertion A1, and
receives the new global number 020 per
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.

**This is the thinnest of the three reconstructed decisions in this batch: the harvest's
3 citations are all inside MIS-003's own planning documents, none in Go source or
tests.** A second pass against the codebase found that the `market` module and its
`CollectorPort` interface do exist in code (§3–§5 below) — that finding was not part of
the original 3-citation harvest and is added here as independently verified code, not as
a fourth mission citation. No production *adapter* implementing `CollectorPort` was found
or cited; the interface and its structural guards exist, the collector behind it does
not.

## Context

The MIS-003 wireframe's "Mercado" screens were designed around competitor pricing and
market-reference data, but no lawful, reliable source for that data existed at
mission time: ML's official OAuth grant for the relevant scope was still pending
(recorded as "G1 OAuth still FAILED" in the milestone brief), and scraping a competitor's
listing pages is forbidden by ML's own terms of service. Shipping the Mercado screens
against scraped or fabricated data would have produced a demo-ready but dishonest
surface — the same class of defect ADR-017 (unknown is never zero) exists to prevent,
here at the level of an entire data source rather than a single field.

## Decision

**Market and competitor data is served only through a `CollectorPort` behind a
contract-only `market` module. No scraping. Where no collector exists yet, the API
responds with an honest empty result, never a fabricated or scraped one.**

**§1 — The decision is named as a cross-cutting mission commitment, not a
milestone-local choice.** It sits in the mission's decision table alongside the other
ratified ADRs, scoped explicitly to preventing "fabricated market facts" and "scraping
drift."
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:113` — "ADR-14 Market
> contract-only behind CollectorPort | decided | fabricated market facts; scraping drift
> | 6-signal separation; honest-empty; no production adapter; G1–G7 sequence for
> successor | M-06 criteria |".

**§2 — The owning milestone commits to shipping zero production adapter.** M-06's brief
binds itself to this decision and to R-04 (a separate, G1–G7-numbered risk-gate sequence
forbidding scraping) as the reason its "market" module ships as a contract skeleton with
no adapter behind it, honestly, rather than a working-looking stub.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-06-corrigir-atributo-market-contracts/milestone.md:17` —
> "contracts: ... R-04 (G1–G7, scraping forbidden), ADR-14 (market contract-only) ...".

**§3 — The `CollectorPort` interface, as shipped.** It declares exactly two methods:
observations for a set of listings, and reference data for a set of products. No other
methods exist on it.
```go
type CollectorPort interface {
	CollectObservations(ctx context.Context, installationID string, listingIDs []string) ([]domain.MarketObservation, error)
	CollectReferences(ctx context.Context, productIDs []string) ([]domain.MarketReference, error)
}
```
> `apps/server_core/internal/modules/market/ports/collector.go:9-14` — "CollectorPort is
> the contract for a future market-data collector. This mission defines the interface
> only; no production implementation is wired." `CollectionService` (the only consumer)
> depends on it as a constructor-injected field, never a concrete adapter.
> `apps/server_core/internal/modules/market/application/collection_service.go:17,23`.

**§4 — No-scraping is enforced structurally, not just by convention.** The `Source`
type's valid values are a closed enum — `official_api`, `vendor`, `manual` — checked by
`validSource`; a value like `"scraper"` fails domain validation before it can be
persisted. A test exists specifically proving this rejection.
> `apps/server_core/internal/modules/market/domain/market.go:48-53,250-257` — `type
> Source string` with `SourceOfficialAPI`, `SourceVendor`, `SourceManual` as the only
> members; `validSource` returns `false` for anything else.
> `apps/server_core/internal/modules/market/domain/market_test.go:101` — `TestMarketEnumValidation`
> sets `Source("scraper")` and asserts the input is rejected with a `"source"` validation
> error.

**§5 — Honest-empty is enforced in the read path, not just described in the contract.**
When no stored row exists for a requested listing or product, the service returns a
struct with nil/zero fields (or the `no_price_evidence` state) instead of omitting the
item or fabricating a value.
> `apps/server_core/internal/modules/market/application/evidence_read_service.go:45-47` —
> "Signals returns the latest competitive signal per requested listing. A listing with no
> stored row gets an honest empty CompetitiveSignal (nil Money fields, zero FetchedAt)
> rather than a fabricated price."
> `apps/server_core/internal/modules/market/application/evidence_read_service.go:75-77` —
> "Aggregates returns the latest product-keyed market aggregate per requested codprod. A
> codprod with no stored row gets an honest empty MarketAggregate (NO_PRICE_EVIDENCE, n=0)
> rather than an invented sample."

## Rationale

A contract-only module that returns honest-empty until a real collector exists is the
same trade this repository makes everywhere else absence is possible (ADR-017): it is
cheaper, and more honest, to ship a seam with nothing behind it yet than to fill that
seam with data obtained by a method (scraping) the source's own terms forbid, or with
numbers invented to make a wireframe demo look complete.

## Consequences

- The Mercado screens described in the MIS-003 wireframe are absent from that mission's
  delivered scope by design — recorded in the mission's own accepted trade-offs list as
  "Mercado screens absent from MVP despite being wireframe headline; honesty over
  demo-value" (`mission.md:118`, not independently re-quoted as a clause above because it
  restates §1 rather than adding a new assertion).
- `CollectorPort` and the structural guards in §3–§5 exist and are covered by tests
  (`market/application/collection_service_test.go:207` and `market/domain/market_test.go:24`
  hold compile-time `var _ ports.CollectorPort = fakeCollector{}` / equivalent
  conformance assertions; `TestCollectorPortRoundTripPersistsMandatoryProvenance` at
  `market/application/collection_service_test.go:57` exercises the round trip). No
  production adapter behind the port was found — every implementation located is a test
  double.
- A successor mission is expected to plug in a real collector behind the same port
  without reshaping it, per the "G1–G7 sequence for successor" language in §1, but no
  citation describes what that sequence requires beyond its name.

## Unverified claims

- The citation harvest names a `market-data-interface-contract.md` research file and an
  `IC-04` designation as the presumed home of a "6-signal separation" and the
  `evidence_state` vocabulary. The `evidence_state` vocabulary is confirmed in code
  (§4/§5 anchors show `EvidenceStateObserved`, `EvidenceStateInsufficientMarket`,
  `EvidenceStateNoPriceEvidence`), but "6-signal separation" itself was not independently
  verified against that research file or against a specific set of six named signals in
  code — this document does not claim to know what the six signals are.
- No production adapter implementing `CollectorPort` was found anywhere in the module
  (only `ports/collector.go`'s interface, its consumers, and test doubles). Whether one
  exists elsewhere in the repo outside `internal/modules/market/` was not checked.
- The "G1–G7 sequence for successor" named in §1 is not defined anywhere this pass
  looked; only its name is confirmed by citation, not its contents.

## Alternatives Considered

None recorded in the citations. The 3 citations available state the decision and its
owning milestone but do not preserve a rejected-alternatives discussion.
