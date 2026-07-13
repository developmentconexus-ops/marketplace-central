# F-02 Oracle Catalog Cutover — Plan

## Machine Work Contract

```json
{"schema_version":"1.0","feature_id":"F-02","required_sources":[".mnfs/MIS-001-mercado-livre-operating-cockpit/research/mnos-sankhya-read-interface-contract.md",".mnfs/MIS-001-mercado-livre-operating-cockpit/research/mvp-operator-workspace-interface-contract.md"],"knowledge_route_ids":["portfolio-core"],"allowed_paths":["apps/server_core/internal/modules/catalog","apps/server_core/internal/modules/internal_read","apps/server_core/internal/composition/root.go","apps/server_core/cmd/server/main.go","apps/server_core/internal/platform/msdb","apps/server_core/tests/unit","contracts/api/marketplace-central.openapi.yaml","packages/sdk-runtime","contracts/governance/runtime-config.json","docker-compose.yml","docker/dev/backend-entrypoint.sh",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-02-oracle-catalog-cutover",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-03-catalog-compatibility-cutover",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/checkpoint.md"],"forbidden_paths":["apps/web",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml"],"side_effects":{"allowed":["repository-write","isolated-cache-write"],"forbidden":["provider-write","database-mutation"]},"commands":[{"id":"catalog-go","command_id":"catalog-go","lane_id":"unit","expected_exit_code":0},{"id":"server-go","command_id":"server-go","lane_id":"unit","expected_exit_code":0},{"id":"sdk","command_id":"sdk","lane_id":"unit","expected_exit_code":0}],"criteria":[{"id":"F02-AC01","milestone_criterion_id":"M-09-C02","command_ids":["catalog-go"]},{"id":"F02-AC02","milestone_criterion_id":"M-09-C03","command_ids":["server-go"]},{"id":"F02-AC03","milestone_criterion_id":"M-09-C05","command_ids":["catalog-go"]}],"stop_conditions":[{"code":"oracle-unavailable","condition":"Governed Oracle read is unavailable."},{"code":"zero-default","condition":"An unknown fact would receive a numeric default."},{"code":"legacy-fallback","condition":"The catalog would retain an MSDB fallback."},{"code":"api-sdk-drift","condition":"The public contract cannot remain atomic."}],"retry_budget":{"max_correction_attempts":2},"handoff_fields":["status","commit","evidence"]}
```

1. Add a catalog-owned adapter over MPC's internal-read port. It maps only a
   positive CODPROD plus separate EAN/reference fields and source facts.
2. Remove active MSDB catalog composition and config. Oracle unavailability is
   reported as `source_unavailable`; it never falls back to MSDB.
3. Prove deterministic catalog behavior plus a sanitized governed Oracle read.

## Stop Rule

Absent canonical Oracle configuration or explicit live-read opt-in is
`externally_blocked`: do not implement against fixtures, infer CODPROD from
legacy strings, or retain an MSDB fallback.
