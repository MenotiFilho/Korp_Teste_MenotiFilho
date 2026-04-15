package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
)

type staleInvoiceListerStub struct {
	listFn      func(ctx context.Context, olderThan time.Duration) ([]domain.Invoice, error)
	updateFn    func(ctx context.Context, id int64, status string) error
	listCalls   int
	updateCalls int
	closedIDs   []int64
}

func (s *staleInvoiceListerStub) ListStaleOpenInvoices(ctx context.Context, olderThan time.Duration) ([]domain.Invoice, error) {
	s.listCalls++
	return s.listFn(ctx, olderThan)
}

func (s *staleInvoiceListerStub) UpdateStatus(ctx context.Context, id int64, status string) error {
	s.updateCalls++
	s.closedIDs = append(s.closedIDs, id)
	return s.updateFn(ctx, id, status)
}

type stockKeyCheckerStub struct {
	checkFn func(ctx context.Context, key string) (bool, error)
}

func (s *stockKeyCheckerStub) IdempotencyKeyExists(ctx context.Context, key string) (bool, error) {
	return s.checkFn(ctx, key)
}

func TestReconciliation_WhenStockDecreased_ShouldCloseInvoice(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{{ID: 1, Numero: 100, Status: domain.StatusAberta}}, nil
		},
		updateFn: func(_ context.Context, id int64, status string) error {
			if id != 1 {
				t.Fatalf("expected invoice id 1, got %d", id)
			}
			if status != domain.StatusFechada {
				t.Fatalf("expected status %q, got %q", domain.StatusFechada, status)
			}
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, key string) (bool, error) {
			if key != "invoice-print-1" {
				t.Fatalf("expected key invoice-print-1, got %q", key)
			}
			return true, nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 1 {
		t.Fatalf("expected 1 UpdateStatus call, got %d", invoiceRepo.updateCalls)
	}
	if len(invoiceRepo.closedIDs) != 1 || invoiceRepo.closedIDs[0] != 1 {
		t.Fatalf("expected invoice 1 to be closed, got %v", invoiceRepo.closedIDs)
	}
}

func TestReconciliation_WhenStockNotDecreased_ShouldSkipInvoice(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{{ID: 2, Numero: 101, Status: domain.StatusAberta}}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			t.Fatal("UpdateStatus should not be called")
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 0 {
		t.Fatalf("expected 0 UpdateStatus calls, got %d", invoiceRepo.updateCalls)
	}
}

func TestReconciliation_WhenNoStaleInvoices_ShouldDoNothing(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			t.Fatal("UpdateStatus should not be called")
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, _ string) (bool, error) {
			t.Fatal("IdempotencyKeyExists should not be called")
			return false, nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.listCalls != 1 {
		t.Fatalf("expected 1 ListStaleOpenInvoices call, got %d", invoiceRepo.listCalls)
	}
	if invoiceRepo.updateCalls != 0 {
		t.Fatalf("expected 0 UpdateStatus calls, got %d", invoiceRepo.updateCalls)
	}
}

func TestReconciliation_WhenStockCheckFails_ShouldNotCloseInvoice(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{{ID: 3, Numero: 102, Status: domain.StatusAberta}}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			t.Fatal("UpdateStatus should not be called")
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, _ string) (bool, error) {
			return false, ErrEstoqueUnavailable
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 0 {
		t.Fatalf("expected 0 UpdateStatus calls, got %d", invoiceRepo.updateCalls)
	}
}

func TestReconciliation_WhenInvoiceAlreadyClosed_ShouldSkip(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{{ID: 4, Numero: 103, Status: domain.StatusAberta}}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			return repository.ErrInvoiceAlreadyClosed
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 1 {
		t.Fatalf("expected 1 UpdateStatus call (attempted), got %d", invoiceRepo.updateCalls)
	}
}

func TestReconciliation_WhenListFails_ShouldNotCheckStock(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return nil, errors.New("db error")
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			t.Fatal("UpdateStatus should not be called")
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, _ string) (bool, error) {
			t.Fatal("IdempotencyKeyExists should not be called")
			return false, nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 0 {
		t.Fatalf("expected 0 UpdateStatus calls, got %d", invoiceRepo.updateCalls)
	}
}

func TestReconciliation_WhenMultipleInvoices_ShouldCloseAllWithDecreasedStock(t *testing.T) {
	// Arrange
	invoiceRepo := &staleInvoiceListerStub{
		listFn: func(_ context.Context, _ time.Duration) ([]domain.Invoice, error) {
			return []domain.Invoice{
				{ID: 1, Numero: 100, Status: domain.StatusAberta},
				{ID: 2, Numero: 101, Status: domain.StatusAberta},
				{ID: 3, Numero: 102, Status: domain.StatusAberta},
			}, nil
		},
		updateFn: func(_ context.Context, _ int64, _ string) error {
			return nil
		},
	}
	stockClient := &stockKeyCheckerStub{
		checkFn: func(_ context.Context, key string) (bool, error) {
			// Only invoice 1 has stock decreased
			return key == "invoice-print-1", nil
		},
	}
	svc := NewReconciliationService(invoiceRepo, stockClient, time.Hour, 2*time.Minute)

	// Act
	svc.reconcile(context.Background())

	// Assert
	if invoiceRepo.updateCalls != 1 {
		t.Fatalf("expected 1 UpdateStatus call, got %d", invoiceRepo.updateCalls)
	}
	if len(invoiceRepo.closedIDs) != 1 || invoiceRepo.closedIDs[0] != 1 {
		t.Fatalf("expected only invoice 1 to be closed, got %v", invoiceRepo.closedIDs)
	}
}
