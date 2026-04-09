package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux_WhenHealthRequested_ShouldReturn200(t *testing.T) {
	// Arrange
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Act
	mux.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
