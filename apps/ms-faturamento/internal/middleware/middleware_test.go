package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/httpapi"
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

func TestMaxBodyBytes_WhenContentLengthExceedsLimit_ShouldReturn413(t *testing.T) {
	// Arrange
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := MaxBodyBytes(5, next)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", nil)
	req.ContentLength = 6
	rec := httptest.NewRecorder()

	// Act
	h.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}

func TestMaxBodyBytes_WhenWithinLimit_ShouldCallNext(t *testing.T) {
	// Arrange
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := MaxBodyBytes(10, next)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", nil)
	req.ContentLength = 5
	rec := httptest.NewRecorder()

	// Act
	h.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !called {
		t.Fatal("expected next handler to be called")
	}
}
