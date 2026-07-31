package apierror

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// recoverMessage is the fixed message answered when a handler panics. The
// recovered value itself never reaches the client — only the log gets it,
// following the same split as httpx.WriteJSON's encode-failure fallback.
const recoverMessage = "erro interno do servidor"

// Recover wraps next so a panicking handler answers the standard 500
// envelope instead of the connection being closed with no body, matching
// what net/http's default server does today. It must live in apierror (not
// httpx) because it answers through Write, and the dependency direction is
// apierror -> httpx: httpx importing apierror back would be a cycle.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &recoverResponseWriter{ResponseWriter: w}
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}
			if rw.written {
				slog.Error("apierror.Recover", "result", "response already committed", "panic", rec)
				return
			}
			slog.Error("apierror.Recover", "result", "500", "panic", rec, "stack", string(debug.Stack()))
			Write(w, http.StatusInternalServerError, "internal_error", recoverMessage, nil)
		}()
		next.ServeHTTP(rw, r)
	})
}

// recoverResponseWriter records whether a response has already been
// committed (headers or body written), so a panic after that point does not
// attempt a second WriteHeader. Modeled on
// internal/platform/httpx/route_deadline.go:156's deadlineResponseWriter.
type recoverResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *recoverResponseWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *recoverResponseWriter) Write(payload []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(payload)
}
