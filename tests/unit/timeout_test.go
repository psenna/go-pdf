package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/psenna/go-pdf/pkg/middleware"
)

func TestTimeoutMiddleware(t *testing.T) {
	// Create a handler that takes long
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		w.WriteHeader(http.StatusRequestTimeout)
	})

	// Create middleware with 50ms timeout
	timeout := 50 * time.Millisecond
	mw := middleware.TimeoutMiddleware(timeout)
	wrapped := mw(slowHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Slow request should timeout, got %d", w.Code)
	}
}
