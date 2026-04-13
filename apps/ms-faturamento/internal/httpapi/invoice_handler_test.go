package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type invoiceServiceStub struct {
	createFn func(ctx context.Context, items []domain.InvoiceItem) (domain.Invoice, error)
	listFn   func(ctx context.Context) ([]domain.Invoice, error)
	updateFn func(ctx context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (s invoiceServiceStub) CreateInvoice(ctx context.Context, items []domain.InvoiceItem) (domain.Invoice, error) {
	return s.createFn(ctx, items)
}

func (s invoiceServiceStub) ListInvoices(ctx context.Context) ([]domain.Invoice, error) {
	if s.listFn == nil {
		return []domain.Invoice{}, nil
	}
	return s.listFn(ctx)
}

func (s invoiceServiceStub) UpdateInvoice(ctx context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error) {
	return s.updateFn(ctx, id, items)
}

func (s invoiceServiceStub) DeleteInvoice(ctx context.Context, id int64) error {
	return s.deleteFn(ctx, id)
}

// ListLatest delegates to listFn when available (used by ListLatestInvoices handler tests)
func (s invoiceServiceStub) ListLatest(ctx context.Context, limit int) ([]domain.Invoice, error) {
	if s.listFn == nil {
		return []domain.Invoice{}, nil
	}
	return s.listFn(ctx)
}

func TestCreateInvoiceHandler_WhenPayloadIsValid_ShouldReturn201WithInvoice(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{createFn: func(_ context.Context, items []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{ID: 1, Numero: 100, Status: domain.StatusAberta, Itens: items}, nil
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-001","quantidade":2}]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	h.CreateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var out invoiceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}
	if out.Numero != 100 {
		t.Fatalf("expected numero 100, got %d", out.Numero)
	}
	if out.Status != domain.StatusAberta {
		t.Fatalf("expected status %q, got %q", domain.StatusAberta, out.Status)
	}
	if len(out.Itens) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out.Itens))
	}
}

func TestCreateInvoiceHandler_WhenPayloadIsInvalidJSON_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{createFn: func(_ context.Context, _ []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{}, nil
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", bytes.NewReader([]byte(`{"itens":`)))
	rec := httptest.NewRecorder()

	// Act
	h.CreateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVALID_JSON")
}

func TestCreateInvoiceHandler_WhenItemsEmpty_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{createFn: func(_ context.Context, _ []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{}, domain.ErrInvoiceItemsRequired
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	h.CreateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestCreateInvoiceHandler_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{createFn: func(_ context.Context, _ []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{}, errors.New("db unavailable")
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-001","quantidade":2}]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	h.CreateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INTERNAL_ERROR")
}

func TestListInvoicesHandler_WhenInvoicesExist_ShouldReturn200WithList(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{listFn: func(_ context.Context) ([]domain.Invoice, error) {
		return []domain.Invoice{
			{ID: 1, Numero: 100, Status: domain.StatusAberta, Itens: []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 2}}},
			{ID: 2, Numero: 101, Status: domain.StatusFechada, Itens: []domain.InvoiceItem{{ID: 2, NotaID: 2, ProdutoCodigo: "P-002", Quantidade: 5}}},
		}, nil
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/notas", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListInvoices(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var out []invoiceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 invoices, got %d", len(out))
	}
	if out[0].Numero != 100 || out[1].Numero != 101 {
		t.Fatalf("unexpected order: %+v", out)
	}
	if len(out[0].Itens) != 1 || len(out[1].Itens) != 1 {
		t.Fatalf("expected items in each invoice")
	}
}

func TestListInvoicesHandler_WhenNoInvoices_ShouldReturn200WithEmptyList(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{listFn: func(_ context.Context) ([]domain.Invoice, error) {
		return []domain.Invoice{}, nil
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/notas", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListInvoices(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var out []invoiceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty list, got %d items", len(out))
	}
}

func TestListInvoicesHandler_WhenServiceFails_ShouldReturn500(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{listFn: func(_ context.Context) ([]domain.Invoice, error) {
		return nil, errors.New("db error")
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/notas", nil)
	rec := httptest.NewRecorder()

	// Act
	h.ListInvoices(rec, req)

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

func TestUpdateInvoiceHandler_WhenPayloadIsValid_ShouldReturn200(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{updateFn: func(_ context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{ID: id, Numero: 100, Status: domain.StatusAberta, Itens: items}, nil
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-002","quantidade":5}]}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/notas/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	// Act
	h.UpdateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var out invoiceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("expected valid JSON, got %v", err)
	}
	if len(out.Itens) != 1 || out.Itens[0].ProdutoCodigo != "P-002" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestUpdateInvoiceHandler_WhenIDIsInvalid_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-001","quantidade":1}]}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/notas/abc", bytes.NewReader(body))
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	// Act
	h.UpdateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestUpdateInvoiceHandler_WhenPayloadIsInvalidJSON_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/notas/1", bytes.NewReader([]byte(`{"itens":`)))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	// Act
	h.UpdateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVALID_JSON")
}

func TestUpdateInvoiceHandler_WhenInvoiceNotFound_ShouldReturn404(t *testing.T) {
	// Arrange
	repoErr := errors.New("invoice not found: id=999")
	svc := invoiceServiceStub{updateFn: func(_ context.Context, _ int64, _ []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{}, repoErr
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-001","quantidade":1}]}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/notas/999", bytes.NewReader(body))
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()

	// Act
	h.UpdateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_FOUND")
}

func TestUpdateInvoiceHandler_WhenInvoiceNotAberta_ShouldReturn409(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{updateFn: func(_ context.Context, _ int64, _ []domain.InvoiceItem) (domain.Invoice, error) {
		return domain.Invoice{}, domain.ErrInvoiceNotAberta
	}}
	h := NewInvoiceHandler(svc)
	body := []byte(`{"itens":[{"produto_codigo":"P-001","quantidade":1}]}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/notas/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	// Act
	h.UpdateInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_ABERTA")
}

func TestDeleteInvoiceHandler_WhenInvoiceIsAberta_ShouldReturn204(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{deleteFn: func(_ context.Context, _ int64) error {
		return nil
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notas/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	// Act
	h.DeleteInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestDeleteInvoiceHandler_WhenIDIsInvalid_ShouldReturn400(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notas/abc", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	// Act
	h.DeleteInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestDeleteInvoiceHandler_WhenInvoiceNotFound_ShouldReturn404(t *testing.T) {
	// Arrange
	repoErr := errors.New("invoice not found: id=999")
	svc := invoiceServiceStub{deleteFn: func(_ context.Context, _ int64) error {
		return repoErr
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notas/999", nil)
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()

	// Act
	h.DeleteInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_FOUND")
}

func TestDeleteInvoiceHandler_WhenInvoiceNotAberta_ShouldReturn409(t *testing.T) {
	// Arrange
	svc := invoiceServiceStub{deleteFn: func(_ context.Context, _ int64) error {
		return domain.ErrInvoiceNotAberta
	}}
	h := NewInvoiceHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notas/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	// Act
	h.DeleteInvoice(rec, req)

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_ABERTA")
}
