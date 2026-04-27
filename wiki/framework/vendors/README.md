# Vendor Playbooks

This section standardizes how we document each marketplace vendor for implementation and LLM retrieval.

## Purpose

- Keep one canonical integration playbook per vendor.
- Separate vendor-specific API behavior from framework-level architecture docs.
- Make vendor onboarding predictable for future channels (Magalu, Mercado Livre, Amazon, etc.).

## Required Structure Per Vendor

Each vendor folder must include:

1. `README.md`
   - quick status, auth model, capability map, and links
2. `getting-started.md`
   - account/app setup, auth bootstrap, sandbox, and first-call flow
3. `api-best-practices.md`
   - product, price, stock, order, logistics, returns, push, and sensitive-data constraints
4. `capability-matrix.md`
   - capability to API-area mapping, module ownership, and rollout notes
5. `sources.md`
   - source links, document IDs, and last verification date

Optional files:

- `section-index.md` and `sections/*.md` for large vendor docsets (100+ pages)
- `operations.md` for runbooks and incident handling
- `gaps-and-risks.md` for known platform limitations
- `implementation-sync.md` for docs-to-code change tracking during active implementation

## Vendor Index

- [Shopee](shopee/README.md)
- [Amazon SP-API](amazon/README.md)

## Template

- [Vendor Playbook Template](_vendor-template.md)
