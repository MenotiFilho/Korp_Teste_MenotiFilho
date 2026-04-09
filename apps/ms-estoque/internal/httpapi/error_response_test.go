package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/requestid"
)

func TestWriteError_WhenCalled_ShouldReturnStandardErrorJSON(t *testing.T) {
	// Arrange
	const reqID = "req-123"
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req = req.WithContext(requestid.WithContext(context.Background(), reqID))
	rec := httptest.NewRecorder()

	// Act
	WriteError(rec, req, http.StatusBadRequest, "INVALID_INPUT", "payload invalido", map[string]string{"field": "codigo"})

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", got)
	}

	var out ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid json, got error: %v", err)
	}
	if out.Code != "INVALID_INPUT" {
		t.Fatalf("expected code INVALID_INPUT, got %q", out.Code)
	}
	if out.Message != "payload invalido" {
		t.Fatalf("expected message payload invalido, got %q", out.Message)
	}
	if out.RequestID != reqID {
		t.Fatalf("expected request_id %q, got %q", reqID, out.RequestID)
	}
}
