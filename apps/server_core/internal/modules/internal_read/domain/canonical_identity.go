package domain

import catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"

// InternalProductID is shared with catalog so every canonical CODPROD uses the
// same positive-integer value object. Legacy ProductID remains separate until
// the deterministic cutover feature removes it.
type InternalProductID = catalogdomain.InternalProductID
