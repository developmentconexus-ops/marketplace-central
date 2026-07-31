// Package apierror is the single producer of the backend's HTTP error envelope.
// Every handler that needs to answer an error writes it through Write, so the
// wire shape cannot drift between modules.
package apierror

import (
	"net/http"

	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type envelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

// Write emits the standard error envelope through httpx.WriteJSON. details is
// always present in the output, even when nil or empty, so callers never have
// to special-case "no details" separately from "some details".
func Write(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	httpx.WriteJSON(w, status, envelope{Error: errorBody{
		Code:    code,
		Message: message,
		Details: details,
	}})
}
