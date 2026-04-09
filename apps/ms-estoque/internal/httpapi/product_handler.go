package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
)

type ProductCreator interface {
	CreateProduct(ctx context.Context, codigo, descricao string, saldo int) (domain.Product, error)
}

type ProductHandler struct {
	service ProductCreator
}

type createProductRequest struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

type productResponse struct {
	ID        int64  `json:"id"`
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

func NewProductHandler(service ProductCreator) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "payload JSON invalido", nil)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), req.Codigo, req.Descricao, req.Saldo)
	if err != nil {
		h.handleCreateError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(productResponse{
		ID:        product.ID,
		Codigo:    product.Codigo,
		Descricao: product.Descricao,
		Saldo:     product.Saldo,
	})
}

func (h *ProductHandler) handleCreateError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, domain.ErrCodigoRequired) || errors.Is(err, domain.ErrDescricaoRequired) || errors.Is(err, domain.ErrSaldoNegative) {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "dados do produto invalidos", map[string]string{"error": err.Error()})
		return
	}

	if errors.Is(err, repository.ErrProductCodigoAlreadyExists) {
		WriteError(w, r, http.StatusConflict, "PRODUCT_CODIGO_ALREADY_EXISTS", "codigo de produto ja cadastrado", nil)
		return
	}

	WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
}
