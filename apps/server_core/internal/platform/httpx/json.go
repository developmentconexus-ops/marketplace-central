package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

// encodeFailureBody is the fixed envelope answered when the encoder itself
// fails. It is a constant literal, not a call back into the encoder (which
// just failed) or into apierror (which would create an import cycle:
// apierror -> httpx is the intended direction, httpx must not import
// apierror).
const encodeFailureBody = `{"error":{"code":"internal_error","message":"failed to encode response","details":{}}}`

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		slog.Error("httpx.WriteJSON", "result", "500", "error", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(encodeFailureBody))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}
