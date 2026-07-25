package httpx

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{
			"service": "marketplace-central-server-core",
			"status":  "ok",
		})
	})
	return mux
}

// CORSMiddleware adds permissive CORS headers for local development.
// All origins are allowed; in production this should be scoped to known domains.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Cache-Control is what the web client sends when the operator presses a
		// Refresh control (it asks for a bypass of every cache in the path). It is
		// not a CORS-simple header, so leaving it out of this list made the browser
		// block the actual request after a successful preflight: every Refresh
		// button failed silently in a real browser while passing in tests, whose
		// fetch is mocked. Pragma travels with it for HTTP/1.0 caches.
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, Pragma")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
