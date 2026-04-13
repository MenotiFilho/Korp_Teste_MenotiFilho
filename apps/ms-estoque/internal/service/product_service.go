package service

import (
	"context"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	ListLowStockProducts(ctx context.Context, threshold int, limit int) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, descricao string, saldo int) (domain.Product, error)
	SoftDeleteProduct(ctx context.Context, id int64) error
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, codigo, descricao string, saldo int) (domain.Product, error) {
	p, err := domain.NewProduct(codigo, descricao, saldo)
	if err != nil {
		return domain.Product{}, err
	}

	return s.repo.CreateProduct(ctx, p)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.ListProducts(ctx)
}

// ListLowStock returns products with low stock using defaults if needed
func (s *ProductService) ListLowStock(ctx context.Context, threshold, limit int) ([]domain.Product, error) {
	if threshold <= 0 {
		threshold = 10
	}
	if limit <= 0 {
		limit = 6
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.ListLowStockProducts(ctx, threshold, limit)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, descricao string, saldo int) (domain.Product, error) {
	if err := domain.ValidateProductUpdate(descricao, saldo); err != nil {
		return domain.Product{}, err
	}

	return s.repo.UpdateProduct(ctx, id, descricao, saldo)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.SoftDeleteProduct(ctx, id)
}
