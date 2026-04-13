package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/service"
)

type invoiceLoaderStub struct {
	getByIDFn func(ctx context.Context, id int64) (domain.Invoice, error)
}

func (s invoiceLoaderStub) GetInvoiceByID(ctx context.Context, id int64) (domain.Invoice, error) {
	return s.getByIDFn(ctx, id)
}

type printServiceStub struct {
	printFn func(ctx context.Context, invoice domain.Invoice) error
}

func (s printServiceStub) Print(ctx context.Context, invoice domain.Invoice) error {
	return s.printFn(ctx, invoice)
}

func TestPrintInvoiceHandler_WhenInvoiceIsAbertaAndStockSucceeds_ShouldReturn200(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, id int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     id,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: id, ProdutoCodigo: "P-001", Quantidade: 2}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error { return nil }}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestPrintInvoiceHandler_WhenInvoiceNotFound_ShouldReturn404(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{}, repository.ErrInvoiceNotFound
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error { return nil }}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/999/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "999")

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_FOUND")
}

func TestPrintInvoiceHandler_WhenInvoiceIsFechada_ShouldReturn409(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{ID: 1, Numero: 100, Status: domain.StatusFechada, Itens: []domain.InvoiceItem{}}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return domain.ErrInvoiceNotAberta
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INVOICE_NOT_ABERTA")
}

func TestPrintInvoiceHandler_WhenEstoqueUnavailable_ShouldReturn503(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     1,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 2}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return service.ErrEstoqueUnavailable
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "ESTOQUE_UNAVAILABLE")
}

func TestPrintInvoiceHandler_WhenInsufficientStock_ShouldReturn409(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     1,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 999}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return service.ErrStockInsufficientStock
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INSUFFICIENT_STOCK")
}

func TestPrintInvoiceHandler_WhenProductNotFoundInStock_ShouldReturn404(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     1,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-999", Quantidade: 1}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return service.ErrStockProductNotFound
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "PRODUCT_NOT_FOUND_IN_STOCK")
}

func TestPrintInvoiceHandler_WhenInvalidID_ShouldReturn400(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error { return nil }}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/abc/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "abc")

	// Assert
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "VALIDATION_ERROR")
}

func TestPrintInvoiceHandler_WhenStatusUpdateFailsAfterStockSucceeds_ShouldReturn500WithSpecificCode(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, id int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     id,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: id, ProdutoCodigo: "P-001", Quantidade: 2}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return service.ErrPrintStatusUpdateFailed
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "PRINT_STATUS_UPDATE_FAILED")
}

func TestPrintInvoiceHandler_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	// Arrange
	loader := invoiceLoaderStub{getByIDFn: func(_ context.Context, _ int64) (domain.Invoice, error) {
		return domain.Invoice{
			ID:     1,
			Numero: 100,
			Status: domain.StatusAberta,
			Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 1}},
		}, nil
	}}
	printer := printServiceStub{printFn: func(_ context.Context, _ domain.Invoice) error {
		return errors.New("unexpected")
	}}

	h := NewPrintInvoiceHandler(loader, printer)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notas/1/imprimir", nil)
	rec := httptest.NewRecorder()

	// Act
	h.PrintInvoice(rec, req, "1")

	// Assert
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	assertErrorCode(t, rec.Body.Bytes(), "INTERNAL_ERROR")
}
