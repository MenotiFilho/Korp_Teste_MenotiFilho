package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
)

type InvoiceCreator interface {
	CreateInvoice(ctx context.Context, items []domain.InvoiceItem) (domain.Invoice, error)
	ListInvoices(ctx context.Context) ([]domain.Invoice, error)
	UpdateInvoice(ctx context.Context, id int64, items []domain.InvoiceItem) (domain.Invoice, error)
	DeleteInvoice(ctx context.Context, id int64) error
}

type InvoiceHandler struct {
	service InvoiceCreator
}

type createInvoiceRequest struct {
	Itens []createInvoiceItemRequest `json:"itens"`
}

type createInvoiceItemRequest struct {
	ProdutoCodigo string `json:"produto_codigo"`
	Quantidade    int    `json:"quantidade"`
}

type invoiceResponse struct {
	ID     int64                 `json:"id"`
	Numero int                   `json:"numero"`
	Status string                `json:"status"`
	Itens  []invoiceItemResponse `json:"itens"`
}

type invoiceItemResponse struct {
	ID            int64  `json:"id"`
	ProdutoCodigo string `json:"produto_codigo"`
	Quantidade    int    `json:"quantidade"`
}

func NewInvoiceHandler(service InvoiceCreator) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req createInvoiceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "payload JSON invalido", nil)
		return
	}

	items := make([]domain.InvoiceItem, 0, len(req.Itens))
	for _, item := range req.Itens {
		items = append(items, domain.InvoiceItem{
			ProdutoCodigo: item.ProdutoCodigo,
			Quantidade:    item.Quantidade,
		})
	}

	invoice, err := h.service.CreateInvoice(r.Context(), items)
	if err != nil {
		if isDomainValidationError(err) {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "dados da nota invalidos", map[string]string{"error": err.Error()})
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toInvoiceResponse(invoice))
}

func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.service.ListInvoices(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	out := make([]invoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		out = append(out, toInvoiceResponse(inv))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func toInvoiceResponse(invoice domain.Invoice) invoiceResponse {
	items := make([]invoiceItemResponse, 0, len(invoice.Itens))
	for _, item := range invoice.Itens {
		items = append(items, invoiceItemResponse{
			ID:            item.ID,
			ProdutoCodigo: item.ProdutoCodigo,
			Quantidade:    item.Quantidade,
		})
	}
	return invoiceResponse{
		ID:     invoice.ID,
		Numero: invoice.Numero,
		Status: invoice.Status,
		Itens:  items,
	}
}

func isDomainValidationError(err error) bool {
	return err == domain.ErrInvoiceItemsRequired ||
		err == domain.ErrProdutoCodigoRequired ||
		err == domain.ErrQuantidadeMustBePositive
}

type updateInvoiceRequest struct {
	Itens []createInvoiceItemRequest `json:"itens"`
}

func (h *InvoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseInt64(idStr)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "id invalido", nil)
		return
	}

	var req updateInvoiceRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			WriteError(w, r, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "payload excede limite permitido", nil)
			return
		}
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "payload JSON invalido", nil)
		return
	}

	items := make([]domain.InvoiceItem, 0, len(req.Itens))
	for _, item := range req.Itens {
		items = append(items, domain.InvoiceItem{
			ProdutoCodigo: item.ProdutoCodigo,
			Quantidade:    item.Quantidade,
		})
	}

	invoice, err := h.service.UpdateInvoice(r.Context(), id, items)
	if err != nil {
		if isDomainValidationError(err) {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "dados da nota invalidos", map[string]string{"error": err.Error()})
			return
		}
		if err == domain.ErrInvoiceNotAberta {
			WriteError(w, r, http.StatusConflict, "INVOICE_NOT_ABERTA", "nota fiscal nao esta ABERTA", nil)
			return
		}
		if isInvoiceNotFoundError(err) {
			WriteError(w, r, http.StatusNotFound, "INVOICE_NOT_FOUND", "nota fiscal nao encontrada", nil)
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toInvoiceResponse(invoice))
}

func (h *InvoiceHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseInt64(idStr)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "id invalido", nil)
		return
	}

	err = h.service.DeleteInvoice(r.Context(), id)
	if err != nil {
		if err == domain.ErrInvoiceNotAberta {
			WriteError(w, r, http.StatusConflict, "INVOICE_NOT_ABERTA", "nota fiscal nao esta ABERTA", nil)
			return
		}
		if isInvoiceNotFoundError(err) {
			WriteError(w, r, http.StatusNotFound, "INVOICE_NOT_FOUND", "nota fiscal nao encontrada", nil)
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseInt64(s string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

func isInvoiceNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "invoice not found")
}
