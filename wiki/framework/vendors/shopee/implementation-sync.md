# Shopee Docs Implementation Sync Runbook

Last updated: 2026-04-27  
Owner: WikiKeeper (docs sync during Shopee implementation)

## Purpose

Keep Shopee docs aligned with real implementation state in small, auditable updates.

## 1) Update Triggers

Update Shopee docs when any of the following happens:

- Provider metadata changes for `shopee` (auth strategy, capabilities, rollout/execution mode).
- Shopee adapter behavior changes (new endpoints, request signing, retries, error mapping).
- Integration/auth flow changes (install, authorize, credential rotation, token lifecycle).
- API contract changes affecting Shopee-facing behavior (`contracts/api/marketplace-central.openapi.yaml`).
- New migration or schema impact related to Shopee integration state.
- Any bugfix where root cause was stale or missing Shopee documentation.
- Any PR/commit that adds/removes Shopee operational constraints or known risks.

## 2) Mandatory Sections Per Update

Every Shopee doc update entry must include these sections (in this order):

1. `What changed`
   - 2-5 bullets with objective deltas only.
2. `Why it changed`
   - One short paragraph with context or incident/task link.
3. `Implementation reference`
   - File paths and/or endpoint identifiers touched.
4. `Evidence`
   - Commands run + pass/fail summary.
5. `Impact`
   - Runtime impact, rollout impact, and backward-compatibility note.
6. `Follow-ups`
   - Explicit next actions or `none`.

## 3) Evidence Checklist (Tests/Commands)

Use this minimal checklist for Shopee doc-sync proof:

- [ ] Scope check:
  - `rg -n "shopee|shopee_" apps/server_core contracts packages wiki`
- [ ] Backend tests (impacted packages at minimum):
  - `cd apps/server_core`
  - PowerShell: `$env:GOCACHE='.gocache'; go test ./...`
- [ ] Frontend/SDK verification (if contract or UX changed):
  - `npm run -w packages/sdk-runtime build`
  - `npm run -w apps/web test`
- [ ] Contract alignment (if endpoint/schema changed):
  - confirm `contracts/api/marketplace-central.openapi.yaml` matches behavior and docs
- [ ] Docs consistency:
  - update relevant Shopee files in `wiki/framework/vendors/shopee/`
  - verify cross-links resolve

If a command is intentionally skipped, record reason in `Evidence`.

## 4) Source-of-Truth Links (Repo)

Use these as canonical references when updating Shopee docs:

- [ARCHITECTURE.md](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/ARCHITECTURE.md)
- [AGENTS.md](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/AGENTS.md)
- [Integrations Module Wiki](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/wiki/modules/integrations.md)
- [Marketplace Framework](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/wiki/framework/marketplace-integration-framework.md)
- [Shopee Vendor README](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/wiki/framework/vendors/shopee/README.md)
- [Shopee Sources](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/wiki/framework/vendors/shopee/sources.md)
- [OpenAPI Contract](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/contracts/api/marketplace-central.openapi.yaml)
- [Integrations Module Code](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/apps/server_core/internal/modules/integrations)
- [Connectors Module Code](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/apps/server_core/internal/modules/connectors)

## 5) Changelog Table Template

Copy this table into the Shopee doc you changed (or `README.md` when broad changes happen):

| Date (UTC) | Area | Change Summary | Evidence | Author | Related |
|---|---|---|---|---|---|
| YYYY-MM-DD | auth \| catalog \| orders \| logistics \| docs | short factual summary | `go test ...` / `npm ...` / `n/a` | @owner | PR/commit/task |

## Suggested Update Snippet (LLM-Friendly)

Use this snippet for each update block:

```md
### Update: YYYY-MM-DD - <area>

#### What changed
- ...

#### Why it changed
...

#### Implementation reference
- /absolute/path/to/file
- `METHOD /route`

#### Evidence
- PASS: `<command>`
- SKIPPED: `<command>` (reason)

#### Impact
- Runtime:
- Rollout:
- Compatibility:

#### Follow-ups
- none
```

### Update: 2026-04-27 - auth/runtime

#### What changed
- Shopee provider moved to interactive runtime metadata (`install_mode=interactive`, `execution_mode=available`, `rollout_stage=v1`).
- Shopee auth strategy changed to `shopee_partner`.
- Signed partner auth flow implemented (`/api/v2/shop/auth_partner`, `/api/v2/auth/token/get`, `/api/v2/auth/access_token/get`).
- Callback provider-account mapping now accepts `shop_id`, `merchant_id`, and `selling_partner_id`.

#### Why it changed
Shopee moved from placeholder/blocked state to first production-grade auth lifecycle implementation within the integrations framework.

#### Implementation reference
- `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/shopee/signer.go`
- `apps/server_core/internal/modules/integrations/application/auth_flow_service.go`
- `apps/server_core/internal/modules/integrations/transport/auth_handler.go`
- `contracts/api/marketplace-central.openapi.yaml`

#### Evidence
- PASS: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/integrations/...`
- PASS: `npm run -w packages/sdk-runtime test -- --run`
- PASS: `npm exec -w packages/feature-marketplaces vitest -- --run`

#### Impact
- Runtime: Shopee installation can run interactive authorization and credential refresh.
- Rollout: Shopee catalog status is now available/v1.
- Compatibility: auth strategy enum expanded to include `shopee_partner` in domain/OpenAPI/SDK.

#### Follow-ups
- Enable and validate additional Shopee capabilities beyond `pricing_fee_sync` as connector implementations land.
