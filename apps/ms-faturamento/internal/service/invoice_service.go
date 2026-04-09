package service

import (
	"context"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error)
	ListInvoices(ctx context.Context) ([]domain.Invoice, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type InvoiceService struct {
	repo InvoiceRepository
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, items []domain.InvoiceItem) (domain.Invoice, error) {
	invoice, err := domain.NewInvoice(items)
	if err != nil {
		return domain.Invoice{}, err
	}

	return s.repo.CreateInvoice(ctx, invoice)
}

func (s *InvoiceService) ListInvoices(ctx context.Context) ([]domain.Invoice, error) {
	return s.repo.ListInvoices(ctx)
}
