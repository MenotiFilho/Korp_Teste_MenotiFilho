package service

import (
	"context"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error)
	ListInvoices(ctx context.Context) ([]domain.Invoice, error)
	GetInvoiceByID(ctx context.Context, id int64) (domain.Invoice, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateInvoiceItems(ctx context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error)
	SoftDeleteInvoice(ctx context.Context, id int64) error
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

func (s *InvoiceService) UpdateInvoice(ctx context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error) {
	invoice, err := s.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return domain.Invoice{}, err
	}

	if err := domain.ValidateInvoiceUpdate(invoice.Status, items); err != nil {
		return domain.Invoice{}, err
	}

	return s.repo.UpdateInvoiceItems(ctx, id, items)
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, id int64) error {
	invoice, err := s.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}

	if err := domain.ValidateInvoiceDelete(invoice.Status); err != nil {
		return err
	}

	return s.repo.SoftDeleteInvoice(ctx, id)
}
