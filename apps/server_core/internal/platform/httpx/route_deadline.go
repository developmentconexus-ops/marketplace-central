package httpx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// RouteRegistrar is the small registration surface shared by the standard
// ServeMux and the composition mux that applies route-class deadlines.
type RouteRegistrar interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type RouteClass string

const (
	InteractiveRouteClass RouteClass = "interactive"
	BatchRouteClass       RouteClass = "batch"
	interactiveDeadline              = 15 * time.Second
	batchDeadline                    = 120 * time.Second
)

// RouteClassMux keeps route-class policy at transport registration while
// retaining the standard library's method-aware ServeMux matching.
type RouteClassMux struct {
	mux     *http.ServeMux
	classes map[string]RouteClass
}

func NewRouteClassMux() *RouteClassMux {
	return &RouteClassMux{
		mux:     http.NewServeMux(),
		classes: make(map[string]RouteClass),
	}
}

// RegisterRouteClass must be called before the matching transport route is
// registered. Routes without an explicit class are interactive by default.
func (m *RouteClassMux) RegisterRouteClass(pattern string, class RouteClass) {
	m.classes[pattern] = class
}

func (m *RouteClassMux) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, deadlineMiddleware(m.classes[pattern], handler))
}

func (m *RouteClassMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(handler))
}

func (m *RouteClassMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func routeClassDeadline(class RouteClass) time.Duration {
	if class == BatchRouteClass {
		return batchDeadline
	}
	return interactiveDeadline
}

func deadlineMiddleware(class RouteClass, next http.Handler) http.Handler {
	return deadlineMiddlewareWithTimeout(routeClassDeadline(class), next)
}

func deadlineMiddlewareWithTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		buffered := newDeadlineResponseWriter()
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(buffered, r.WithContext(ctx))
		}()

		select {
		case <-done:
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				writeDeadlineExceeded(w)
				return
			}
			buffered.replay(w)
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				writeDeadlineExceeded(w)
			}
		}
	})
}

func writeDeadlineExceeded(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusGatewayTimeout)
	_, _ = io.WriteString(w, `{"error":"deadline_exceeded"}`)
}

type deadlineResponseWriter struct {
	mu           sync.Mutex
	header       http.Header
	body         bytes.Buffer
	status       int
	wroteHeaders bool
}

func newDeadlineResponseWriter() *deadlineResponseWriter {
	return &deadlineResponseWriter{header: make(http.Header)}
}

func (w *deadlineResponseWriter) Header() http.Header {
	return w.header
}

func (w *deadlineResponseWriter) WriteHeader(status int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.wroteHeaders {
		return
	}
	w.status = status
	w.wroteHeaders = true
}

func (w *deadlineResponseWriter) Write(payload []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.wroteHeaders {
		w.status = http.StatusOK
		w.wroteHeaders = true
	}
	return w.body.Write(payload)
}

func (w *deadlineResponseWriter) replay(dst http.ResponseWriter) {
	w.mu.Lock()
	status := w.status
	if status == 0 {
		status = http.StatusOK
	}
	header := w.header.Clone()
	body := append([]byte(nil), w.body.Bytes()...)
	w.mu.Unlock()

	for key, values := range header {
		dst.Header()[key] = values
	}
	dst.WriteHeader(status)
	_, _ = dst.Write(body)
}
