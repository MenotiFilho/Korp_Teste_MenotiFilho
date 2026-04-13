package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/service"
)

type StockDecreaser interface {
	DecreaseStock(ctx context.Context, items []service.StockDecreaseInput, idempotencyKey string) error
}

type StockHandler struct {
	service StockDecreaser
}

type decreaseStockRequest struct {
	Itens []decreaseStockItemRequest `json:"itens"`
}

type decreaseStockItemRequest struct {
	Codigo     string `json:"codigo"`
	Quantidade int    `json:"quantidade"`
}

func NewStockHandler(service StockDecreaser) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	var req decreaseStockRequest
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

	if len(req.Itens) == 0 {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "itens e obrigatorio e nao pode ser vazio", nil)
		return
	}

	inputs := make([]service.StockDecreaseInput, 0, len(req.Itens))
	for _, item := range req.Itens {
		inputs = append(inputs, service.StockDecreaseInput{
			Codigo:     item.Codigo,
			Quantidade: item.Quantidade,
		})
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if strings.TrimSpace(idempotencyKey) == "" {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "cabecalho Idempotency-Key e obrigatorio", map[string]string{"field": "Idempotency-Key"})
		return
	}
	if err := h.service.DecreaseStock(r.Context(), inputs, idempotencyKey); err != nil {
		h.handleDecreaseError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StockHandler) handleDecreaseError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, repository.ErrInvalidDecreaseItem) {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "itens de baixa invalidos", map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, repository.ErrIdempotencyKeyRequired) {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "cabecalho Idempotency-Key e obrigatorio", map[string]string{"field": "Idempotency-Key"})
		return
	}

	if errors.Is(err, repository.ErrProductNotFound) {
		WriteError(w, r, http.StatusNotFound, "PRODUCT_NOT_FOUND", "produto nao encontrado", nil)
		return
	}

	if errors.Is(err, repository.ErrProductInsufficientStock) {
		WriteError(w, r, http.StatusConflict, "INSUFFICIENT_STOCK", "saldo insuficiente", nil)
		return
	}

	WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
}
