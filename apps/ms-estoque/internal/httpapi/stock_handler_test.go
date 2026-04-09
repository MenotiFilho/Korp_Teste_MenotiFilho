package httpapi

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/service"
)

type stockDecreaserStub struct {
	decreaseFn func(ctx context.Context, items []service.StockDecreaseInput, idempotencyKey string) error
}

func (s stockDecreaserStub) DecreaseStock(ctx context.Context, items []service.StockDecreaseInput, idempotencyKey string) error {
	return s.decreaseFn(ctx, items, idempotencyKey)
}

func TestDecreaseStockHandler_WhenPayloadIsValid_ShouldReturn200(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, items []service.StockDecreaseInput, idempotencyKey string) error {
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		if items[0].Codigo != "P-001" || items[0].Quantidade != 2 {
			t.Fatalf("unexpected item: %+v", items[0])
		}
		if idempotencyKey != "idem-abc" {
			t.Fatalf("expected idempotency key idem-abc, got %q", idempotencyKey)
		}
		return nil
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"P-001","quantidade":2}]}`)))
	req.Header.Set("Idempotency-Key", "idem-abc")
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDecreaseStockHandler_WhenPayloadIsInvalidJSON_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error { return nil }}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":`)))
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVALID_JSON")
}

func TestDecreaseStockHandler_WhenIdempotencyKeyMissing_ShouldReturn400(t *testing.T) {
	// Arrange
	called := false
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error {
		called = true
		return nil
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"P-001","quantidade":2}]}`)))
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if called {
		t.Fatal("expected service to NOT be called when Idempotency-Key is missing")
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestDecreaseStockHandler_WhenInvalidItem_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error {
		return repository.ErrInvalidDecreaseItem
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"","quantidade":0}]}`)))
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestDecreaseStockHandler_WhenProductNotFound_ShouldReturn404(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error {
		return repository.ErrProductNotFound
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"P-999","quantidade":1}]}`)))
	req.Header.Set("Idempotency-Key", "idem-notfound")
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "PRODUCT_NOT_FOUND")
}

func TestDecreaseStockHandler_WhenInsufficientStock_ShouldReturn409(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error {
		return repository.ErrProductInsufficientStock
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"P-001","quantidade":999}]}`)))
	req.Header.Set("Idempotency-Key", "idem-insufficient")
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INSUFFICIENT_STOCK")
}

func TestDecreaseStockHandler_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	// Arrange
	svc := stockDecreaserStub{decreaseFn: func(_ context.Context, _ []service.StockDecreaseInput, _ string) error {
		return errors.New("db timeout")
	}}
	h := NewStockHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/estoque/baixa", bytes.NewReader([]byte(`{"itens":[{"codigo":"P-001","quantidade":1}]}`)))
	req.Header.Set("Idempotency-Key", "idem-unexpected")
	rec := httptest.NewRecorder()

	// Act
	h.DecreaseStock(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INTERNAL_ERROR")
}
