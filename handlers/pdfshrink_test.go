package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShrinkHandler(t *testing.T) {
	handler := ShrinkHandler

	req := httptest.NewRequest("POST", "/api/pdf/shrink", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
