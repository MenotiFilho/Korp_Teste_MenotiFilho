package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
)

type createProductInput struct {
	Codigo    string
	Descricao string
	Saldo     int
}

type productCreatorStub struct {
	createFn func(ctx context.Context, codigo, descricao string, saldo int) (domain.Product, error)
	listFn   func(ctx context.Context) ([]domain.Product, error)
}

func (s productCreatorStub) CreateProduct(ctx context.Context, codigo, descricao string, saldo int) (domain.Product, error) {
	return s.createFn(ctx, codigo, descricao, saldo)
}

func (s productCreatorStub) ListProducts(ctx context.Context) ([]domain.Product, error) {
	if s.listFn == nil {
		return []domain.Product{}, nil
	}
	return s.listFn(ctx)
}

func TestCreateProductHandler_WhenPayloadIsValid_ShouldReturn201WithProduct(t *testing.T) {
	// Arrange
	svc := productCreatorStub{createFn: func(_ context.Context, codigo, descricao string, saldo int) (domain.Product, error) {
		return domain.Product{ID: 1, Codigo: codigo, Descricao: descricao, Saldo: saldo}, nil
	}}
	h := NewProductHandler(svc)
	body := []byte(`{"codigo":"P-001","descricao":"Produto 1","saldo":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	h.CreateProduct(rec, req)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var out productResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON body, got %v", err)
	}
	if out.Codigo != "P-001" || out.Descricao != "Produto 1" || out.Saldo != 10 {
		t.Fatalf("unexpected response body: %+v", out)
	}
}

func TestCreateProductHandler_WhenPayloadIsInvalidJSON_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := productCreatorStub{createFn: func(_ context.Context, _, _ string, _ int) (domain.Product, error) {
		return domain.Product{}, nil
	}}
	h := NewProductHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader([]byte(`{"codigo":`)))
	rec := httptest.NewRecorder()

	// Act
	h.CreateProduct(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVALID_JSON")
}

func TestCreateProductHandler_WhenDomainValidationFails_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := productCreatorStub{createFn: func(_ context.Context, _, _ string, _ int) (domain.Product, error) {
		return domain.Product{}, domain.ErrCodigoRequired
	}}
	h := NewProductHandler(svc)
	body := []byte(`{"codigo":"","descricao":"Produto 1","saldo":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	h.CreateProduct(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestCreateProductHandler_WhenCodigoDuplicated_ShouldReturn409(t *testing.T) {
	// Arrange
	svc := productCreatorStub{createFn: func(_ context.Context, _, _ string, _ int) (domain.Product, error) {
		return domain.Product{}, repository.ErrProductCodigoAlreadyExists
	}}
	h := NewProductHandler(svc)
	body := []byte(`{"codigo":"P-001","descricao":"Produto 1","saldo":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	h.CreateProduct(rec, req)

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "PRODUCT_CODIGO_ALREADY_EXISTS")
}

func TestCreateProductHandler_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	// Arrange
	svc := productCreatorStub{createFn: func(_ context.Context, _, _ string, _ int) (domain.Product, error) {
		return domain.Product{}, errors.New("unexpected db error")
	}}
	h := NewProductHandler(svc)
	body := []byte(`{"codigo":"P-001","descricao":"Produto 1","saldo":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/produtos", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	h.CreateProduct(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INTERNAL_ERROR")
}

func TestListProductsHandler_WhenProductsExist_ShouldReturn200WithOrderedList(t *testing.T) {
	// Arrange
	svc := productCreatorStub{listFn: func(_ context.Context) ([]domain.Product, error) {
		return []domain.Product{
			{ID: 1, Codigo: "P-001", Descricao: "Produto 1", Saldo: 10},
			{ID: 2, Codigo: "P-002", Descricao: "Produto 2", Saldo: 5},
		}, nil
	}}
	h := NewProductHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListProducts(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var out []productResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON body, got %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 products, got %d", len(out))
	}
	if out[0].Codigo != "P-001" || out[1].Codigo != "P-002" {
		t.Fatalf("unexpected order/content: %+v", out)
	}
}

func TestListProductsHandler_WhenNoProducts_ShouldReturn200WithEmptyList(t *testing.T) {
	// Arrange
	svc := productCreatorStub{listFn: func(_ context.Context) ([]domain.Product, error) {
		return []domain.Product{}, nil
	}}
	h := NewProductHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListProducts(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var out []productResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON body, got %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty list, got %d items", len(out))
	}
}

func TestListProductsHandler_WhenServiceFails_ShouldReturn500(t *testing.T) {
	// Arrange
	svc := productCreatorStub{listFn: func(_ context.Context) ([]domain.Product, error) {
		return nil, errors.New("db unavailable")
	}}
	h := NewProductHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/produtos", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListProducts(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INTERNAL_ERROR")
}

func assertErrorCode(t *testing.T, body []byte, expected string) {
	t.Helper()

	var out ErrorResponse
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("expected valid JSON error body, got %v", err)
	}
	if out.Code != expected {
		t.Fatalf("expected error code %q, got %q", expected, out.Code)
	}
}
