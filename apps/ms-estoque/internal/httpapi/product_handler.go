package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/repository"
)

type ProductCreator interface {
	CreateProduct(ctx context.Context, codigo, descricao string, saldo int) (domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, descricao string, saldo int) (domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
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
		if strings.Contains(err.Error(), "http: request body too large") {
			WriteError(w, r, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "payload excede limite permitido", nil)
			return
		}
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

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	out := make([]productResponse, 0, len(products))
	for _, p := range products {
		out = append(out, productResponse{
			ID:        p.ID,
			Codigo:    p.Codigo,
			Descricao: p.Descricao,
			Saldo:     p.Saldo,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

// ListLowStockProducts handles GET /api/v1/produtos/baixo-estoque?threshold=N&limit=M
func (h *ProductHandler) ListLowStockProducts(w http.ResponseWriter, r *http.Request) {
	qThreshold := r.URL.Query().Get("threshold")
	qLimit := r.URL.Query().Get("limit")

	threshold := 6
	limit := 6
	if qThreshold != "" {
		t, err := strconv.Atoi(qThreshold)
		if err != nil || t <= 0 {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "threshold invalido", nil)
			return
		}
		threshold = t
	}
	if qLimit != "" {
		l, err := strconv.Atoi(qLimit)
		if err != nil || l <= 0 {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "limit invalido", nil)
			return
		}
		if l > 100 {
			l = 100
		}
		limit = l
	}

	products, err := h.service.(interface {
		ListLowStock(context.Context, int, int) ([]domain.Product, error)
	}).ListLowStock(r.Context(), threshold, limit)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	out := make([]productResponse, 0, len(products))
	for _, p := range products {
		out = append(out, productResponse{ID: p.ID, Codigo: p.Codigo, Descricao: p.Descricao, Saldo: p.Saldo})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
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

type updateProductRequest struct {
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseInt64(idStr)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "id invalido", nil)
		return
	}

	var req updateProductRequest
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

	product, err := h.service.UpdateProduct(r.Context(), id, req.Descricao, req.Saldo)
	if err != nil {
		if errors.Is(err, domain.ErrDescricaoRequired) || errors.Is(err, domain.ErrSaldoNegative) {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "dados do produto invalidos", map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, repository.ErrProductNotFound) {
			WriteError(w, r, http.StatusNotFound, "PRODUCT_NOT_FOUND", "produto nao encontrado", nil)
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(productResponse{
		ID:        product.ID,
		Codigo:    product.Codigo,
		Descricao: product.Descricao,
		Saldo:     product.Saldo,
	})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseInt64(idStr)
	if err != nil {
		WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "id invalido", nil)
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			WriteError(w, r, http.StatusNotFound, "PRODUCT_NOT_FOUND", "produto nao encontrado", nil)
			return
		}
		WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "erro interno do servidor", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
