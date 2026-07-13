# M-09-C03 Wording Normalization Validation

- base SHA: `32b32f6de00875589468c71eb70c6eb3e5d49278`
- scope: comments/documentation only; no executable or contract behavior changed.
- normalized paths:
  - `apps/server_core/internal/modules/catalog/ports/repository.go`
  - `apps/server_core/internal/modules/catalog/application/service.go`
  - `docker/dev/README.md`
- catalog port wording now identifies the active canonical source adapter.
- enrichment wording records MPC manual enrichment precedence over retained
  canonical source values.
- Docker documentation records that M-09 removed the legacy catalog database
  variable from active composition, `MC_DATABASE_URL` remains MPC PostgreSQL,
  and governed Oracle configuration owns internal reads.
- exact active-residue scan over server entrypoint, composition, config,
  catalog port/application/transport, runtime governance, Compose, and Docker
  dev paths for `MetalShopping`, `MS_DATABASE_URL`, `MS_TENANT_ID`, and
  `platform/msdb`: **PASS**, zero matches.
- side effects: none beyond repository wording/evidence edits; no database,
  provider, Oracle, config, environment, dependency, or runtime action.
