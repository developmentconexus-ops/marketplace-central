# Module: Product Links

Layer: operational mapping
Path: `apps/server_core/internal/modules/product_links/` (planned)

## Main Question It Answers

"Which internal product or SKU does this Mercado Livre listing or variation represent?"

## What This Module Owns

| Entity | Purpose |
|---|---|
| `ProductLink` | Tenant-scoped link between internal product identity and provider listing identity |
| `LinkCandidate` | Suggested link from SKU/EAN/title heuristics before operator approval |
| `LinkAudit` | Trace of who/what created, changed, or rejected a link |

## Rules

- Product links are prerequisites for stock writes, order margin, and price recommendations.
- Ambiguous links must remain unresolved; downstream modules must surface data-quality flags instead of guessing.
- Provider identifiers are strings, even when the provider returns numeric-looking IDs.
- The module owns mapping state; connectors only fetch provider identifiers.
- `internal_read` missing or ambiguous product candidates map to unresolved link states, not fallback guesses.

## Initial Mercado Livre Scope

- Map internal `CODPROD`/SKU/EAN to Mercado Livre `item_id` and variation ID when present.
- Track link confidence: manual, exact SKU/EAN, title match, unresolved, conflict.
- Block write actions when link state is unresolved or conflict.
