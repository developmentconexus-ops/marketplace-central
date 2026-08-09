// Package application is an ADR-023 detector fixture: ports is the one public
// surface, so this import is legal and must produce no finding.
package application

import _ "marketplace-central/apps/server_core/internal/modules/beta/ports"
