package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/httpapi"
)

func TestRecover_WhenHandlerPanics_ShouldReturnStandardErrorJSON(t *testing.T) {
	// Arrange
	h := Recover(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom")
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var out httpapi.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if out.Code != "INTERNAL_ERROR" {
		t.Fatalf("expected INTERNAL_ERROR, got %q", out.Code)
	}
}
