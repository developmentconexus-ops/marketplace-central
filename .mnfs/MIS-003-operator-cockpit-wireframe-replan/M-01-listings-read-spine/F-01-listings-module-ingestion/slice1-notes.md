# F-01 Slice 1 notes

## Blocking governance conflict

The required governance command was run literally from the worktree root:

```text
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance
```

Actual result:

```text
status=failed
error_code=GOV_SEMANTIC_DRIFT
id=base-sha-invalid
```

The harness requires a base SHA for its `governance` (validate + drift) mode even though the slice brief does not provide one. With the current HEAD as the base (`d0d30d68b8b58f5b42bfc2fa45f04a9869ca67b9`), the rerun reached validation and reported:

```text
error_code=GOV_COMPOSITION_MISSING
id=listings
path=apps/server_core/internal/composition/root.go
```

Slice 1 explicitly allows only the migration, migration test, listings domain entity/test, and `contracts/governance/modules.json`. It also forbids work outside this slice. Registering a composition-required module therefore creates a governance requirement for composition wiring that cannot be satisfied within the allowed paths. `apps/server_core/internal/composition/root.go` must be handled by the composition owner or the slice scope must be amended by the milestone owner.

The same governance run also reported pre-existing repository findings for catalog/internal_read, internal_read/inventory, orders/internal_read, dynamic configuration readers, and baseline exceptions. No new migration-prefix exception was added; the only migration-prefix baseline reported was `migration-prefix-0021-duplicate`.
