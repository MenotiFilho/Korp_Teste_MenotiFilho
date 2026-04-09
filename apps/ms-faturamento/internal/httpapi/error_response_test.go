package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/requestid"
)

func TestWriteError_WhenCalled_ShouldReturnStandardErrorJSON(t *testing.T) {
	// Arrange
	const reqID = "req-456"
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req = req.WithContext(requestid.WithContext(context.Background(), reqID))
	rec := httptest.NewRecorder()

	// Act
	WriteError(rec, req, http.StatusBadRequest, "VALIDATION_ERROR", "dados invalidos", map[string]string{"field": "numero"})

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", got)
	}

	var out ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if out.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected code VALIDATION_ERROR, got %q", out.Code)
	}
	if out.RequestID != reqID {
		t.Fatalf("expected request_id %q, got %q", reqID, out.RequestID)
	}
}
