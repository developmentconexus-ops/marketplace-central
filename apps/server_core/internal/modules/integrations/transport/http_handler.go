package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type ProviderReader interface {
	ListProviderDefinitions(ctx context.Context) ([]domain.ProviderDefinition, error)
}

type InstallationReader interface {
	List(ctx context.Context) ([]domain.Installation, error)
	Get(ctx context.Context, installationID string) (domain.Installation, bool, error)
	CreateDraft(ctx context.Context, input application.CreateInstallationInput) (domain.Installation, error)
}

type Handler struct {
	providerReader     ProviderReader
	installationReader InstallationReader
}

func NewHandler(providerReader ProviderReader, installationReader InstallationReader) Handler {
	return Handler{
		providerReader:     providerReader,
		installationReader: installationReader,
	}
}

// writeIntegrationError forwards to apierror.Write. Kept (not inlined at
// every call site) because auth_handler.go in this package calls it too and
// is outside this slice's write_set — deleting it would break that file's
// build. See slice report for the alternatives-considered note.
func writeIntegrationError(w http.ResponseWriter, status int, code, message string) {
	apierror.Write(w, status, code, message, nil)
}

func mapIntegrationError(err error) (int, string, string) {
	if err == nil {
		return http.StatusInternalServerError, "INTEGRATIONS_INTERNAL_ERROR", "internal error"
	}

	msg := err.Error()
	if strings.HasPrefix(msg, "INTEGRATIONS_") {
		return http.StatusBadRequest, msg, msg
	}

	return http.StatusInternalServerError, "INTEGRATIONS_INTERNAL_ERROR", "internal error"
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/integrations/providers", h.handleProviders)
	mux.HandleFunc("/integrations/installations", h.handleInstallations)
}

func (h Handler) handleProviders(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		slog.Info("integrations.providers", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, http.StatusMethodNotAllowed, "INTEGRATIONS_PROVIDER_METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}

	items, err := h.providerReader.ListProviderDefinitions(r.Context())
	if err != nil {
		status, code, message := mapIntegrationError(err)
		slog.Error("integrations.providers", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, status, code, message, nil)
		return
	}

	slog.Info("integrations.providers", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleInstallations(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	switch r.Method {
	case http.MethodGet:
		items, err := h.installationReader.List(r.Context())
		if err != nil {
			status, code, message := mapIntegrationError(err)
			slog.Error("integrations.installations", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			apierror.Write(w, status, code, message, nil)
			return
		}

		slog.Info("integrations.installations", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})

	case http.MethodPost:
		var req struct {
			InstallationID string `json:"installation_id"`
			ProviderCode   string `json:"provider_code"`
			Family         string `json:"family"`
			DisplayName    string `json:"display_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Info("integrations.installations", "action", "decode", "result", "400", "duration_ms", time.Since(start).Milliseconds())
			apierror.Write(w, http.StatusBadRequest, "INTEGRATIONS_INSTALLATION_INVALID", "malformed request body", nil)
			return
		}

		installation, err := h.installationReader.CreateDraft(r.Context(), application.CreateInstallationInput{
			InstallationID: req.InstallationID,
			ProviderCode:   req.ProviderCode,
			Family:         req.Family,
			DisplayName:    req.DisplayName,
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			slog.Error("integrations.installations", "action", "create_draft", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			apierror.Write(w, status, code, message, nil)
			return
		}

		slog.Info("integrations.installations", "action", "create_draft", "result", "201", "installation_id", installation.InstallationID, "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusCreated, installation)

	default:
		w.Header().Set("Allow", "GET, POST")
		slog.Info("integrations.installations", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, http.StatusMethodNotAllowed, "INTEGRATIONS_INSTALLATION_METHOD_NOT_ALLOWED", "method not allowed", nil)
	}
}
