package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/httpapi"
)

func TestMaxBodyBytes_WhenContentLengthExceedsLimit_ShouldReturn413(t *testing.T) {
	// Arrange
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := MaxBodyBytes(5, next)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", strings.NewReader("123456"))
	req.ContentLength = 6
	rec := httptest.NewRecorder()

	// Act
	h.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
	var out httpapi.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if out.Code != "PAYLOAD_TOO_LARGE" {
		t.Fatalf("expected PAYLOAD_TOO_LARGE, got %q", out.Code)
	}
}

func TestMaxBodyBytes_WhenWithinLimit_ShouldCallNext(t *testing.T) {
	// Arrange
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})
	h := MaxBodyBytes(10, next)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", strings.NewReader("12345"))
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
