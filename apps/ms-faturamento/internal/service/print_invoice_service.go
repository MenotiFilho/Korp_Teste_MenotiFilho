package service

import (
	"context"
	"fmt"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type InvoiceUpdater interface {
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type StockClientInterface interface {
	DecreaseStock(ctx context.Context, items []domain.StockDecreaseItem, idempotencyKey string) error
}

type PrintInvoiceService struct {
	invoiceRepo InvoiceUpdater
	stockClient StockClientInterface
}

func NewPrintInvoiceService(invoiceRepo InvoiceUpdater, stockClient StockClientInterface) *PrintInvoiceService {
	return &PrintInvoiceService{
		invoiceRepo: invoiceRepo,
		stockClient: stockClient,
	}
}

func (s *PrintInvoiceService) Print(ctx context.Context, invoice domain.Invoice) error {
	if invoice.Status != domain.StatusAberta {
		return domain.ErrInvoiceNotAberta
	}

	stockItems := invoice.StockDecreaseItems()
	idempotencyKey := fmt.Sprintf("invoice-print-%d", invoice.ID)
	if err := s.stockClient.DecreaseStock(ctx, stockItems, idempotencyKey); err != nil {
		return err
	}

	if err := s.invoiceRepo.UpdateStatus(ctx, invoice.ID, domain.StatusFechada); err != nil {
		return err
	}

	return nil
}
