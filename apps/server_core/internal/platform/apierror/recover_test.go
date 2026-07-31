package apierror

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecover_PanicBecomesEnvelope(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom-SENTINEL-9f3a")
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	want := trimJSON(`{"error":{"code":"internal_error","message":"erro interno do servidor","details":{}}}`)
	got := trimJSON(w.Body.String())
	if got != want {
		t.Fatalf("body =\n%s\nwant\n%s", got, want)
	}
	if strings.Contains(w.Body.String(), "SENTINEL") {
		t.Fatalf("body leaked panic value: %s", w.Body.String())
	}
}

func TestRecover_SaneRouteUntouched(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Body.String(); got != `{"ok":true}` {
		t.Fatalf("body = %q, want {\"ok\":true}", got)
	}
	if strings.Contains(w.Body.String(), "internal_error") {
		t.Fatalf("body has panic shape: %s", w.Body.String())
	}
}

func TestRecover_AlreadyCommittedResponseIsLeftAlone(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("partial"))
		panic("boom-after-commit")
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (unchanged)", w.Code)
	}
	if got := w.Body.String(); got != "partial" {
		t.Fatalf("body = %q, want %q (envelope must not be appended)", got, "partial")
	}
}

func TestRecover_ErrAbortHandlerIsRepanicked(t *testing.T) {
	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	panicked := false
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	}()

	if !panicked {
		t.Fatal("http.ErrAbortHandler was swallowed, want it to propagate out of Recover")
	}
}

func TestRecover_PanicDetailReachesLog(t *testing.T) {
	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(previous)

	handler := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom-SENTINEL-log-9f3a")
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if !strings.Contains(buf.String(), "SENTINEL") {
		t.Fatalf("log output missing panic sentinel:\n%s", buf.String())
	}
	if strings.Contains(w.Body.String(), "SENTINEL") {
		t.Fatalf("body leaked panic value: %s", w.Body.String())
	}
}
