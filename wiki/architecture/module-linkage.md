# Module Linkage: Marketplaces <-> Integrations

## Summary

Two modules, two concerns, loose coupling.

- Marketplaces owns business config (`Account`, `Policy`, `FeeSchedule`).
- Integrations owns technical operations (`Installation`, `Credential`, `OperationRun`).

## Database Relationship

One-way and nullable:

```sql
marketplace_accounts.integration_installation_id
  -> integration_installations(tenant_id, installation_id)
  ON DELETE SET NULL
```

Guarantees:
- A marketplace account can exist without integration installation.
- An integration installation can exist without marketplace account.
- There is no reverse FK from integrations tables to marketplaces tables.

## Application Glue

Fee sync is where both modules collaborate:

```text
FeeSyncService.StartSync(installation_id)
  -> load installation (integrations)
  -> resolve marketplace_code from provider_code
  -> upsert fee schedules (marketplaces)
  -> persist operation run (integrations)
```

This is service-level orchestration, not DB-level coupling.

## Why This Matters

- Failure isolation: integration auth issues do not block marketplace config CRUD.
- Optional linkage: teams can configure account/policy before wiring OAuth.
- Easier migration path to MetalShopping module boundaries later.
