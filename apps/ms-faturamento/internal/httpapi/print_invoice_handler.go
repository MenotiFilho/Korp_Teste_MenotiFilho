package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/repository"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/service"
)

type InvoiceLoader interface {
	GetInvoiceByID(ctx context.Context, id int64) (domain.Invoice, error)
}

type InvoicePrinter interface {
	Print(ctx context.Context, invoice domain.Invoice) error
}

type PrintInvoiceHandler struct {
	loader  InvoiceLoader
	printer InvoicePrinter
}

func NewPrintInvoiceHandler(loader InvoiceLoader, printer InvoicePrinter) *PrintInvoiceHandler {
	return &PrintInvoiceHandler{
		loader:  loader,
		printer: printer,
	}
}

func (h *PrintInvoiceHandler) PrintInvoice(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "id da nota invalido", nil)
		return
	}

	invoice, err := h.loader.GetInvoiceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrInvoiceNotFound) {
			WriteError(w, r, http.StatusNotFound, "INVOICE_NOT_FOUND", "nota fiscal nao encontrada", nil)
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	if err := h.printer.Print(r.Context(), invoice); err != nil {
		h.handlePrintError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PrintInvoiceHandler) handlePrintError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, domain.ErrInvoiceNotAberta) {
		WriteError(w, r, http.StatusConflict, "INVOICE_NOT_ABERTA", "nota fiscal nao esta em status ABERTA", nil)
		return
	}

	if errors.Is(err, service.ErrEstoqueUnavailable) {
		WriteError(w, r, http.StatusServiceUnavailable, "ESTOQUE_UNAVAILABLE", "servico de estoque indisponivel", nil)
		return
	}

	if errors.Is(err, service.ErrStockInsufficientStock) {
		WriteError(w, r, http.StatusConflict, "INSUFFICIENT_STOCK", "saldo insuficiente no estoque", nil)
		return
	}

	if errors.Is(err, service.ErrStockProductNotFound) {
		WriteError(w, r, http.StatusNotFound, "PRODUCT_NOT_FOUND_IN_STOCK", "produto nao encontrado no estoque", nil)
		return
	}

	if errors.Is(err, service.ErrPrintStatusUpdateFailed) {
		WriteError(w, r, http.StatusInternalServerError, "PRINT_STATUS_UPDATE_FAILED", "impressao realizada mas falha ao atualizar status, tente novamente", nil)
		return
	}

	WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
}
