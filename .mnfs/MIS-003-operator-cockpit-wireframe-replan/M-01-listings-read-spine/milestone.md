# M-01-listings-read-spine

```yaml
id: M-01
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-16
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-003 — `../mission.md`; contracts: IC-02 (`../research/listings-read-interface-contract.md`), ADR-12/17.

## Outcome

A new read-only `listings` Go module exists: canonical tenant-scoped `listings` table populated from Mercado Livre via the existing connectors `ListListings` capability, served by `GET /listings`, `GET /listings/by-product`, `GET /listings/{id}`, `GET /listings/summary`, `POST /listings/refresh` — all per IC-02, with OpenAPI + sdk-runtime updated in the same commit. Observable: seeded integration tests return the IC-02 canonical shapes; a refresh against a connected installation upserts real rows (live-provider-read lane, read-only).

## Why This Milestone Exists

Every downstream milestone consumes listing rows or the IC-02 filter grammar (M-02 UI, M-03 selection, M-05 counters). Building it first prevents the six-screen drift that motivated the mission.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | listings-module-ingestion | Module skeleton, migration, domain model, ML ingestion + refresh endpoint |
| F-02 | listings-read-api | The four read endpoints + OpenAPI + SDK |

F-01 before F-02 (F-02 reads what F-01 persists). No parallel execution — same module, one writer.

## Dependencies

None (first milestone). Consumes existing: connectors ML capability, integrations installations, product_links (link join), governance lanes.

## Risks

- RK-02 (ingestion volume/pagination): full-pull semantics; failure lands on operation run.
- RK-06 (G1 OAuth defect — StartAuthorize clobbers connected installation — trips refresh/ingestion): live refresh lane (M01-C10) runs against an already-connected installation without re-authorizing; ingestion auth failures land on the operation run as `provider_auth`-class errors, never silently zeroed; operator reconnects via Integrações before re-running.
- Adapter mapping gaps: unmappable status → `unknown` enum value; unmappable modality → `listing_type` NULL (nullable fact, no `unknown` code) — never guessed (ADR-17, IC-02).

## Done Means

All IC-02 operations implemented and proven per `validation-contract.md` (M01-C01..C10); governance lanes green; no UI change; no write capability touched.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 with feature context pack.
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: none.

## Correction Handoff

Not applicable at planning time.
