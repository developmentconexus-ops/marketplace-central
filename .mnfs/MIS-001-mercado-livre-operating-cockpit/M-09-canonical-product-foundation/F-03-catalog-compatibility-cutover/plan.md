# F-03 Catalog Compatibility and Cutover — Plan

## Machine Work Contract

```json
{"schema_version":"1.0","feature_id":"F-03","required_sources":[".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-02-oracle-catalog-cutover/validation.md"],"knowledge_route_ids":["portfolio-core"],"allowed_paths":["apps/server_core/internal/modules/catalog","apps/server_core/internal/modules/classifications","apps/server_core/internal/modules/pricing","apps/server_core/internal/modules/product_links","apps/server_core/internal/composition/root.go","apps/server_core/cmd/server/main.go","apps/server_core/internal/platform/msdb","apps/server_core/migrations","apps/server_core/tests/unit","contracts/api/marketplace-central.openapi.yaml","packages/sdk-runtime","contracts/governance/runtime-config.json","docker-compose.yml","docker/dev/backend-entrypoint.sh",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-02-oracle-catalog-cutover",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-03-catalog-compatibility-cutover",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/checkpoint.md"],"forbidden_paths":["apps/web",".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml"],"side_effects":{"allowed":["repository-write","isolated-cache-write"],"forbidden":["provider-write","database-mutation"]},"commands":[{"id":"compat-go","command_id":"compat-go","lane_id":"unit","expected_exit_code":0},{"id":"server-go","command_id":"server-go","lane_id":"unit","expected_exit_code":0}],"criteria":[{"id":"F03-AC01","milestone_criterion_id":"M-09-C03","command_ids":["server-go"]},{"id":"F03-AC02","milestone_criterion_id":"M-09-C04","command_ids":["compat-go"]}],"stop_conditions":[{"code":"nondeterministic-mapping","condition":"A reference cannot prove positive CODPROD equality."},{"code":"legacy-runtime-residue","condition":"An active MSDB reader remains."}],"retry_budget":{"max_correction_attempts":2},"handoff_fields":["status","commit","evidence"]}
```

1. Preserve only explicit positive numeric CODPROD equality through consumer
   compatibility adapters; record `not_found` and `identity_conflict` rather
   than guessing.
2. Remove the active MSDB composition/configuration path and prove the residue
   scan is empty for active runtime files.
