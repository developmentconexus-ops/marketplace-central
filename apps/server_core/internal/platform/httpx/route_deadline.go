package httpx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
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
//
// Every caller in this repo declares the class under the bare path while the
// transport registers the ServeMux pattern with its method ("POST /erp/imports"),
// so the lookup below matches on the path and the exact pattern is only an
// optional narrowing for a route that wants one method treated differently.
func (m *RouteClassMux) RegisterRouteClass(pattern string, class RouteClass) {
	m.classes[pattern] = class
}

func (m *RouteClassMux) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, deadlineMiddleware(m.routeClass(pattern), handler))
}

// routeClass resolves the class declared for a ServeMux pattern. The exact
// pattern wins when one was declared; otherwise the method is stripped and the
// path is looked up.
func (m *RouteClassMux) routeClass(pattern string) RouteClass {
	if class, ok := m.classes[pattern]; ok {
		return class
	}
	return m.classes[routePatternPath(pattern)]
}

// routePatternPath drops the optional method from a ServeMux pattern
// ("[METHOD ][HOST]/[PATH]"), leaving the part transports declare classes under.
func routePatternPath(pattern string) string {
	if index := strings.IndexByte(pattern, ' '); index >= 0 {
		return strings.TrimLeft(pattern[index+1:], " ")
	}
	return pattern
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

// deadlineExceededBody is the fixed envelope answered when a route-class
// deadline fires. It is a constant literal, not a call into apierror: the
// dependency direction is apierror -> httpx, and httpx importing apierror
// back would be an import cycle (same reasoning as json.go's
// encodeFailureBody, which this mirrors).
const deadlineExceededBody = `{"error":{"code":"deadline_exceeded","message":"tempo limite excedido","details":{}}}`

func writeDeadlineExceeded(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusGatewayTimeout)
	_, _ = io.WriteString(w, deadlineExceededBody)
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
