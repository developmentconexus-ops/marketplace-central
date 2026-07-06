package transport

import (
	"net/http"
	"reflect"

	connectorports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type Handler struct {
	meAuth connectorports.MEAuthPort // nil if ME_CLIENT_ID not set
}

func NewHandler(meAuth connectorports.MEAuthPort) *Handler {
	meAuth = normalizeMEAuth(meAuth)
	return &Handler{meAuth: meAuth}
}

func normalizeMEAuth(meAuth connectorports.MEAuthPort) connectorports.MEAuthPort {
	if meAuth == nil {
		return nil
	}

	v := reflect.ValueOf(meAuth)
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return nil
		}
	}

	return meAuth
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/connectors/melhor-envio/auth/start", h.handleMEAuthStart)
	mux.HandleFunc("/connectors/melhor-envio/auth/callback", h.handleMEAuthCallback)
	mux.HandleFunc("/connectors/melhor-envio/status", h.handleMEStatus)
}

func writeConnectorsError(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func (h *Handler) handleMEAuthStart(w http.ResponseWriter, r *http.Request) {
	if h.meAuth == nil {
		writeConnectorsError(w, http.StatusServiceUnavailable, "CONNECTORS_ME_NOT_CONFIGURED", "Melhor Envio is not configured (ME_CLIENT_ID missing)")
		return
	}
	h.meAuth.HandleStart(w, r)
}

func (h *Handler) handleMEAuthCallback(w http.ResponseWriter, r *http.Request) {
	if h.meAuth == nil {
		writeConnectorsError(w, http.StatusServiceUnavailable, "CONNECTORS_ME_NOT_CONFIGURED", "Melhor Envio is not configured")
		return
	}
	h.meAuth.HandleCallback(w, r)
}

func (h *Handler) handleMEStatus(w http.ResponseWriter, r *http.Request) {
	if h.meAuth == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"connected": false})
		return
	}
	h.meAuth.HandleStatus(w, r)
}
