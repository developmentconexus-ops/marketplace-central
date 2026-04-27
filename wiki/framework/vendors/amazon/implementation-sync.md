# Amazon Docs Implementation Sync Runbook

Last updated: 2026-04-27  
Owner: WikiKeeper (docs sync during Amazon implementation)

## Purpose

Keep Amazon SP-API docs aligned with real implementation state in small, auditable updates.

## Update triggers

- Provider metadata changes for `amazon`.
- Auth/signing behavior changes (OAuth, LWA, SigV4, RDT).
- Connector capability changes (listings, inventory, orders, notifications, finance).
- Contract changes that affect Amazon runtime behavior.
- Security/compliance changes (secret rotation, restricted data posture).

## Mandatory sections per update

1. `What changed`
2. `Why it changed`
3. `Implementation reference`
4. `Evidence`
5. `Impact`
6. `Follow-ups`

## Evidence checklist

- Scope check:
  - `rg -n "amazon|sp-api|sellingpartner|lwa|sigv4" apps/server_core contracts packages wiki`
- Backend tests (impacted packages minimum):
  - `cd apps/server_core`
  - PowerShell: `$env:GOCACHE='.gocache'; go test ./...`
- Frontend/SDK verification (if contract/UX changed):
  - `npm run -w packages/sdk-runtime build`
  - `npm run -w apps/web test`
- Contract alignment if endpoint/schema changed:
  - `contracts/api/marketplace-central.openapi.yaml`
- Wiki cross-link check:
  - `wiki/framework/vendors/amazon/*`

