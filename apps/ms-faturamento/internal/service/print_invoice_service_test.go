package service

import (
	"context"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type invoiceRepoStub struct {
	updateStatusFn func(ctx context.Context, id int64, status string) error
}

func (s invoiceRepoStub) UpdateStatus(ctx context.Context, id int64, status string) error {
	return s.updateStatusFn(ctx, id, status)
}

type stockClientStub struct {
	decreaseFn func(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error
}

func (s stockClientStub) DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
	return s.decreaseFn(ctx, items, idempotencyKey)
}

func TestPrintInvoiceService_WhenInvoiceIsAbertaAndStockSucceeds_ShouldCloseInvoice(t *testing.T) {
	// Arrange
	repo := invoiceRepoStub{updateStatusFn: func(_ context.Context, id int64, status string) error {
		if id != 1 {
			t.Fatalf("expected id 1, got %d", id)
		}
		if status != domain.StatusFechada {
			t.Fatalf("expected status %q, got %q", domain.StatusFechada, status)
		}
		return nil
	}}
	stock := stockClientStub{decreaseFn: func(_ context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
		if idempotencyKey != "invoice-print-1" {
			t.Fatalf("expected idempotency key invoice-print-1, got %q", idempotencyKey)
		}
		return nil
	}}

	svc := NewPrintInvoiceService(repo, stock)

	invoice := domain.Invoice{
		ID:     1,
		Numero: 100,
		Status: domain.StatusAberta,
		Itens: []domain.InvoiceItem{
			{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 2},
			{ID: 2, NotaID: 1, ProdutoCodigo: "P-002", Quantidade: 5},
		},
	}

	// Act
	err := svc.Print(context.Background(), invoice)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPrintInvoiceService_WhenInvoiceIsFechada_ShouldReturnErrorWithoutCallingStock(t *testing.T) {
	// Arrange
	stockCalled := false
	repo := invoiceRepoStub{updateStatusFn: func(_ context.Context, _ int64, _ string) error { return nil }}
	stock := stockClientStub{decreaseFn: func(_ context.Context, _ []domain.StockDecreaseItem, _ string) error {
		stockCalled = true
		return nil
	}}

	svc := NewPrintInvoiceService(repo, stock)

	invoice := domain.Invoice{
		ID:     1,
		Numero: 100,
		Status: domain.StatusFechada,
		Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 2}},
	}

	// Act
	err := svc.Print(context.Background(), invoice)

	// Assert
	if err == nil {
		t.Fatal("expected error for non-ABERTA invoice")
	}
	if err != ErrInvoiceNotAberta {
		t.Fatalf("expected ErrInvoiceNotAberta, got %v", err)
	}
	if stockCalled {
		t.Fatal("expected stock client to NOT be called")
	}
}

func TestPrintInvoiceService_WhenEstoqueUnavailable_ShouldReturnEstoqueErrorAndNotUpdateStatus(t *testing.T) {
	// Arrange
	statusCalled := false
	repo := invoiceRepoStub{updateStatusFn: func(_ context.Context, _ int64, _ string) error {
		statusCalled = true
		return nil
	}}
	stock := stockClientStub{decreaseFn: func(_ context.Context, _ []domain.StockDecreaseItem, _ string) error {
		return ErrEstoqueUnavailable
	}}

	svc := NewPrintInvoiceService(repo, stock)

	invoice := domain.Invoice{
		ID:     1,
		Numero: 100,
		Status: domain.StatusAberta,
		Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 2}},
	}

	// Act
	err := svc.Print(context.Background(), invoice)

	// Assert
	if err == nil {
		t.Fatal("expected error when estoque unavailable")
	}
	if err != ErrEstoqueUnavailable {
		t.Fatalf("expected ErrEstoqueUnavailable, got %v", err)
	}
	if statusCalled {
		t.Fatal("expected status to NOT be updated when stock fails")
	}
}

func TestPrintInvoiceService_WhenInsufficientStock_ShouldReturnErrorAndNotUpdateStatus(t *testing.T) {
	// Arrange
	statusCalled := false
	repo := invoiceRepoStub{updateStatusFn: func(_ context.Context, _ int64, _ string) error {
		statusCalled = true
		return nil
	}}
	stock := stockClientStub{decreaseFn: func(_ context.Context, _ []domain.StockDecreaseItem, _ string) error {
		return ErrStockInsufficientStock
	}}

	svc := NewPrintInvoiceService(repo, stock)

	invoice := domain.Invoice{
		ID:     1,
		Numero: 100,
		Status: domain.StatusAberta,
		Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-001", Quantidade: 999}},
	}

	// Act
	err := svc.Print(context.Background(), invoice)

	// Assert
	if err == nil {
		t.Fatal("expected error for insufficient stock")
	}
	if err != ErrStockInsufficientStock {
		t.Fatalf("expected ErrStockInsufficientStock, got %v", err)
	}
	if statusCalled {
		t.Fatal("expected status to NOT be updated when stock fails")
	}
}

func TestPrintInvoiceService_WhenProductNotFoundInStock_ShouldReturnErrorAndNotUpdateStatus(t *testing.T) {
	// Arrange
	statusCalled := false
	repo := invoiceRepoStub{updateStatusFn: func(_ context.Context, _ int64, _ string) error {
		statusCalled = true
		return nil
	}}
	stock := stockClientStub{decreaseFn: func(_ context.Context, _ []domain.StockDecreaseItem, _ string) error {
		return ErrStockProductNotFound
	}}

	svc := NewPrintInvoiceService(repo, stock)

	invoice := domain.Invoice{
		ID:     1,
		Numero: 100,
		Status: domain.StatusAberta,
		Itens:  []domain.InvoiceItem{{ID: 1, NotaID: 1, ProdutoCodigo: "P-999", Quantidade: 1}},
	}

	// Act
	err := svc.Print(context.Background(), invoice)

	// Assert
	if err == nil {
		t.Fatal("expected error for product not found")
	}
	if err != ErrStockProductNotFound {
		t.Fatalf("expected ErrStockProductNotFound, got %v", err)
	}
	if statusCalled {
		t.Fatal("expected status to NOT be updated when stock fails")
	}
}
