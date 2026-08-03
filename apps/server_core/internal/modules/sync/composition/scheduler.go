// Package composition wires the sync module for the application root. It exists
// so composition/root.go depends on a single sync entrypoint (one import, one
// ticker line) instead of reaching into the adapter and application packages
// directly.
package composition

// InstallationScopeERP is the installation_id used for ERP-sourced entities
// (products via xlsx/sankhya) that are not tied to a marketplace installation.
// It keeps sync_state's (tenant, installation, entity) key well-formed without
// fabricating an ML installation id.
const InstallationScopeERP = "erp"

// InstallationScopeMarket is the installation_id used for the tenant-wide
// market-evidence collection job (F-A3). market_aggregates is keyed by
// (tenant_id, product_id, source) — not per marketplace installation — so this
// mirrors InstallationScopeERP's sentinel rather than a real installation id.
// Named (not inlined) so a future producer keying sync_state.entity="market"
// under a different installation_id shows up in a grep for this constant
// instead of colliding invisibly (see migration 0093).
const InstallationScopeMarket = "market"
