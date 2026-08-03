package composition

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	syncpg "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// NewInstallationScheduler constrói um scheduler de sync ligado a UMA
// instalação, sem job registrado — quem chama registra a entidade que lhe cabe.
//
// Existe para que módulos de fora do sync (orders, listings) parem de importar
// sync/adapters/postgres só para instanciar o repositório de estado: alvo de
// camada `adapters` é violação de governança (GOV_MODULE_LAYER), e o construtor
// aqui é o mesmo trecho que NewProductsScheduler já executava — só que
// reutilizável.
func NewInstallationScheduler(pool *pgxpool.Pool, tenantID, installationID string, interval time.Duration) *syncapp.Scheduler {
	return syncapp.NewScheduler(syncpg.NewSyncStateRepository(pool, tenantID), installationID, interval, time.Now)
}
