package httpapi

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/domain"
	"github.com/stretchr/testify/require"
)

// reuse productCreatorStub from product_handler_test.go by declaring a local alias
func TestListLowStockProductsHandler_WhenExists_ShouldReturn200(t *testing.T) {
	svc := productCreatorStub{listFn: func(_ context.Context) ([]domain.Product, error) {
		return []domain.Product{{ID: 1, Codigo: "P-1", Descricao: "p", Saldo: 2}}, nil
	}}
	h := NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/produtos/baixo-estoque", nil)
	rec := httptest.NewRecorder()

	h.ListLowStockProducts(rec, req)
	require.Equal(t, 200, rec.Code)

	var out []domain.Product
	err := json.Unmarshal(rec.Body.Bytes(), &out)
	require.NoError(t, err)
	require.Len(t, out, 1)
}
