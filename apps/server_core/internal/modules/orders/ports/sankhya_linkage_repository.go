package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type SankhyaLinkageRepository interface {
	LoadCurrent(ctx context.Context, scope domain.LinkageScope) (domain.SankhyaLinkage, bool, error)
	AppendConfirmation(ctx context.Context, linkage domain.SankhyaLinkage) (domain.SankhyaLinkage, error)
}
