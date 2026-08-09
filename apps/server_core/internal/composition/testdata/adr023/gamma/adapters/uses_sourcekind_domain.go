// Package adapters is an ADR-023 detector fixture: sourcekind is the shared
// core carve-out and must produce no finding at any layer.
package adapters

import _ "marketplace-central/apps/server_core/internal/modules/sourcekind/domain"
