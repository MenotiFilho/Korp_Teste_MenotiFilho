package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
)

type stockRepositoryStub struct {
	decreaseFn func(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error
}

func (s stockRepositoryStub) DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
	return s.decreaseFn(ctx, items, idempotencyKey)
}

func TestStockService_DecreaseStock_WhenItemsAreValid_ShouldCallRepository(t *testing.T) {
	// Arrange
	called := false
	repo := stockRepositoryStub{decreaseFn: func(_ context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error {
		called = true
		if idempotencyKey != "idem-123" {
			t.Fatalf("expected idempotency key idem-123, got %q", idempotencyKey)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		if items[0].Codigo != "P-001" || items[0].Quantidade != 2 {
			t.Fatalf("unexpected item payload: %+v", items[0])
		}
		return nil
	}}
	svc := NewStockService(repo)

	// Act
	err := svc.DecreaseStock(context.Background(), []StockDecreaseInput{{Codigo: "P-001", Quantidade: 2}}, "idem-123")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected repository to be called")
	}
}

func TestStockService_DecreaseStock_WhenInsufficientStock_ShouldReturnRepositoryError(t *testing.T) {
	// Arrange
	repo := stockRepositoryStub{decreaseFn: func(_ context.Context, _ []domain.StockDecreaseItem, _ string) error {
		return repository.ErrProductInsufficientStock
	}}
	svc := NewStockService(repo)

	// Act
	err := svc.DecreaseStock(context.Background(), []StockDecreaseInput{{Codigo: "P-001", Quantidade: 2}}, "idem-456")

	// Assert
	if !errors.Is(err, repository.ErrProductInsufficientStock) {
		t.Fatalf("expected ErrProductInsufficientStock, got %v", err)
	}
}
