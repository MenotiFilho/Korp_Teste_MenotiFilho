package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
)

var ErrPrintStatusUpdateFailed = errors.New("print succeeded but status update failed")

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
	stockErr := s.stockClient.DecreaseStock(ctx, stockItems, idempotencyKey)
	if stockErr != nil {
		return stockErr
	}

	statusErr := s.invoiceRepo.UpdateStatus(ctx, invoice.ID, domain.StatusFechada)
	if statusErr != nil {
		if errors.Is(statusErr, repository.ErrInvoiceAlreadyClosed) {
			return nil
		}
		return ErrPrintStatusUpdateFailed
	}

	return nil
}
