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
