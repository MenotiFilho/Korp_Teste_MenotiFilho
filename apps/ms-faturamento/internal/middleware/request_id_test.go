package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/requestid"
)

func TestRequestID_WhenHeaderMissing_ShouldGenerateAndExposeID(t *testing.T) {
	// Arrange
	var gotRequestID string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = requestid.FromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if gotRequestID == "" {
		t.Fatal("expected request id in context")
	}
	if rec.Header().Get(requestIDHeader) == "" {
		t.Fatal("expected request id in response header")
	}
	if rec.Header().Get(requestIDHeader) != gotRequestID {
		t.Fatalf("expected header and context ids to match, header=%q context=%q", rec.Header().Get(requestIDHeader), gotRequestID)
	}
}

func TestRequestID_WhenHeaderPresent_ShouldReuseIncomingID(t *testing.T) {
	// Arrange
	const incomingID = "req-incoming-456"
	var gotRequestID string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = requestid.FromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(requestIDHeader, incomingID)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if gotRequestID != incomingID {
		t.Fatalf("expected context id %q, got %q", incomingID, gotRequestID)
	}
	if rec.Header().Get(requestIDHeader) != incomingID {
		t.Fatalf("expected response header id %q, got %q", incomingID, rec.Header().Get(requestIDHeader))
	}
}
