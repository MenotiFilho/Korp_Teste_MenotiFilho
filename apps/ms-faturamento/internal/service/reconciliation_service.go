package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
)

type StaleInvoiceLister interface {
	ListStaleOpenInvoices(ctx context.Context, olderThan time.Duration) ([]domain.Invoice, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type StockKeyChecker interface {
	IdempotencyKeyExists(ctx context.Context, key string) (bool, error)
}

type ReconciliationService struct {
	invoiceRepo StaleInvoiceLister
	stockClient StockKeyChecker
	interval    time.Duration
	staleAfter  time.Duration
}

func NewReconciliationService(
	invoiceRepo StaleInvoiceLister,
	stockClient StockKeyChecker,
	interval time.Duration,
	staleAfter time.Duration,
) *ReconciliationService {
	return &ReconciliationService{
		invoiceRepo: invoiceRepo,
		stockClient: stockClient,
		interval:    interval,
		staleAfter:  staleAfter,
	}
}

func (s *ReconciliationService) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	slog.Info("reconciliation job started",
		"interval", s.interval,
		"stale_after", s.staleAfter,
	)

	for {
		select {
		case <-ctx.Done():
			slog.Info("reconciliation job stopped")
			return
		case <-ticker.C:
			s.reconcile(ctx)
		}
	}
}

func (s *ReconciliationService) reconcile(ctx context.Context) {
	invoices, err := s.invoiceRepo.ListStaleOpenInvoices(ctx, s.staleAfter)
	if err != nil {
		slog.Error("reconciliation: failed to list stale invoices", "error", err)
		return
	}

	if len(invoices) == 0 {
		return
	}

	slog.Info("reconciliation: found stale open invoices", "count", len(invoices))

	for _, inv := range invoices {
		s.reconcileInvoice(ctx, inv)
	}
}

func (s *ReconciliationService) reconcileInvoice(ctx context.Context, inv domain.Invoice) {
	idempotencyKey := fmt.Sprintf("invoice-print-%d", inv.ID)

	exists, err := s.stockClient.IdempotencyKeyExists(ctx, idempotencyKey)
	if err != nil {
		slog.Error("reconciliation: failed to check idempotency key",
			"invoice_id", inv.ID,
			"error", err,
		)
		return
	}

	if !exists {
		slog.Debug("reconciliation: stock not yet decreased, skipping",
			"invoice_id", inv.ID,
		)
		return
	}

	err = s.invoiceRepo.UpdateStatus(ctx, inv.ID, domain.StatusFechada)
	if err != nil {
		if errors.Is(err, repository.ErrInvoiceAlreadyClosed) {
			slog.Debug("reconciliation: invoice already closed",
				"invoice_id", inv.ID,
			)
			return
		}
		slog.Error("reconciliation: failed to close invoice",
			"invoice_id", inv.ID,
			"error", err,
		)
		return
	}

	slog.Info("reconciliation: invoice closed successfully",
		"invoice_id", inv.ID,
		"numero", inv.Numero,
	)
}
