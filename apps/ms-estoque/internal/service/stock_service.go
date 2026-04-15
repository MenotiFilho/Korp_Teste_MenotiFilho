package service

import (
	"context"
	"errors"
	"strings"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
)

type StockRepository interface {
	DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error
	IdempotencyKeyExists(ctx context.Context, key string) error
}

type StockService struct {
	repo StockRepository
}

type StockDecreaseInput struct {
	Codigo     string
	Quantidade int
}

func NewStockService(repo StockRepository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) DecreaseStock(ctx context.Context, inputs []StockDecreaseInput, idempotencyKey string) error {
	items := make([]domain.StockDecreaseItem, 0, len(inputs))
	for _, in := range inputs {
		items = append(items, domain.StockDecreaseItem{
			Codigo:     strings.TrimSpace(in.Codigo),
			Quantidade: in.Quantidade,
		})
	}

	return s.repo.DecreaseStock(ctx, items, idempotencyKey)
}

func (s *StockService) IdempotencyKeyExists(ctx context.Context, key string) error {
	return s.repo.IdempotencyKeyExists(ctx, key)
}

func IsStockDomainError(err error) bool {
	return errors.Is(err, repository.ErrInvalidDecreaseItem) ||
		errors.Is(err, repository.ErrProductNotFound) ||
		errors.Is(err, repository.ErrProductInsufficientStock)
}
