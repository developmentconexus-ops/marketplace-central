// Package transport wires the sync module's read-only aggregate endpoints
// onto HTTP. /sync/health (F-01, M-09) is the first route in this package;
// /sync/runs stays owned by the integrations module and is untouched here.
package transport

import (
	"context"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// HealthReadService is the narrow surface HealthHandler depends on — the
// concrete *application.HealthService satisfies it, and a fake can stand in
// for it in tests without importing the whole application package.
type HealthReadService interface {
	GetHealth(ctx context.Context) (application.HealthResult, error)
}

type HealthHandler struct {
	service HealthReadService
}

type healthRouteClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func NewHealthHandler(service HealthReadService) HealthHandler {
	return HealthHandler{service: service}
}

// Register mounts GET /sync/health as an INTERACTIVE route (the 15s budget,
// not the 120s batch one) — this is a single tenant-scoped read against
// sync_state, never a bulk operation. Explicitly outside registerBatchRoutes.
func (h HealthHandler) Register(mux httpx.RouteRegistrar) {
	if registrar, ok := mux.(healthRouteClassRegistrar); ok {
		registrar.RegisterRouteClass("/sync/health", httpx.InteractiveRouteClass)
	}
	mux.HandleFunc("GET /sync/health", h.HandleGet)
}

func (h HealthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		apierror.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", nil)
		return
	}

	result, err := h.service.GetHealth(r.Context())
	if err != nil {
		apierror.Write(w, http.StatusInternalServerError, "internal_error", "internal error", nil)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toHealthResponseDTO(result))
}

type healthResponseDTO struct {
	Entities []healthEntityDTO `json:"entities"`
	Webhook  healthWebhookDTO  `json:"webhook"`
}

type healthEntityDTO struct {
	Entity              string            `json:"entity"`
	LastSuccessAt       *time.Time        `json:"last_success_at"`
	LastIncrementalAt   *time.Time        `json:"last_incremental_at"`
	ConsecutiveFailures int               `json:"consecutive_failures"`
	Phase               *string           `json:"phase"`
	LastError           *domain.SyncError `json:"last_error"`
}

type healthWebhookDTO struct {
	LastNotificationAt *time.Time `json:"last_notification_at"`
	Pending            int        `json:"pending"`
	Dropped24h         int        `json:"dropped_24h"`
}

func toHealthResponseDTO(result application.HealthResult) healthResponseDTO {
	entities := make([]healthEntityDTO, 0, len(result.Entities))
	for _, e := range result.Entities {
		entities = append(entities, healthEntityDTO{
			Entity:              e.Entity,
			LastSuccessAt:       e.LastSuccessAt,
			LastIncrementalAt:   e.LastIncrementalAt,
			ConsecutiveFailures: e.ConsecutiveFailures,
			Phase:               e.Phase,
			LastError:           e.LastError,
		})
	}
	return healthResponseDTO{
		Entities: entities,
		Webhook: healthWebhookDTO{
			LastNotificationAt: result.Webhook.LastNotificationAt,
			Pending:            result.Webhook.Pending,
			Dropped24h:         result.Webhook.Dropped24h,
		},
	}
}
