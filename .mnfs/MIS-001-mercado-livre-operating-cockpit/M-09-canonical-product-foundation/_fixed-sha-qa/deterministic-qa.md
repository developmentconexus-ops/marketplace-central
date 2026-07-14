# M-09 Fixed-SHA Proportional QA Evidence

- frozen_sha: `32b32f6de00875589468c71eb70c6eb3e5d49278`
- observed_at: `2026-07-13T22:41:48.6614485Z`
- mode: `proportional_qa`
- prerequisite_reviews: `passed` (incoming C01 and inventory-clock fixed-SHA reviews)
- final_status: `failed`
- sha_drift: `false`

## Executed lanes

1. Targeted Go lane from `apps/server_core` with absolute repository `GOCACHE`
   and `-count=1` — PASS. Catalog, internal-read, product-links,
   classifications, pricing, composition, server, unit, migration runner, and
   migration contract packages passed.
2. Full Go lane `go test ./... -count=1` from `apps/server_core` with the same
   `GOCACHE` — PASS. The inventory application package passed, confirming the
   authorized test-clock unblock.
3. SDK lane
   `npm test --workspace @marketplace-central/sdk-runtime -- --run` — PASS
   (40/40).
4. Governed runner contract
   `scripts/tests/live-oracle-docker-runner.tests.ps1` — PASS (14/14).
5. Exact active-residue scan over server, composition, config, catalog,
   governance, and dev-runtime paths — FAIL.

HEAD equaled the frozen SHA around every completed major lane and after the
failure.

## Stop-condition evidence

The residue scan found four active-path matches:

- `apps/server_core/internal/modules/catalog/ports/repository.go:9` — legacy
  `MetalShopping` product-reader wording.
- `apps/server_core/internal/modules/catalog/application/service.go:72-73` —
  legacy `MetalShopping` enrichment/snapshot wording.
- `docker/dev/README.md:61` — states that server boot still opens
  `MS_DATABASE_URL` for legacy catalog composition.

This is a registered command failure and contradicts the M-09-C03 requirement
that the exact active-runtime residue scan have no matches. QA stopped
immediately. The live Oracle command was not run, and the older Oracle evidence
was not refreshed or promoted to this SHA.

No integration database-write test, provider/Oracle write, source edit,
credential/raw-row/PII exposure, or SHA drift occurred.
