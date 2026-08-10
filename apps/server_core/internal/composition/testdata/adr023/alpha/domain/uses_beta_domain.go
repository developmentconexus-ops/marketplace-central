// Package domain is an ADR-023 detector fixture: a domain layer importing a
// sibling module's domain is THE violation. Never compiled (testdata).
package domain

import _ "marketplace-central/apps/server_core/internal/modules/beta/domain"
