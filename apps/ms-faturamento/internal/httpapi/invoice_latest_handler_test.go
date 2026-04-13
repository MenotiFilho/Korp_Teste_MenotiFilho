package httpapi

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-faturamento/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestListLatestInvoicesHandler_WhenExists_ShouldReturn200(t *testing.T) {
	// reuse invoiceServiceStub declared in invoice_handler_test.go
	svc := invoiceServiceStub{listFn: func(_ context.Context) ([]domain.Invoice, error) {
		return []domain.Invoice{{ID: 1, Numero: 10, Status: domain.StatusAberta}}, nil
	}}
	// ensure ListLatest uses same behavior via ListLatest delegating to listFn
	h := NewInvoiceHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/notas/ultimas", nil)
	rec := httptest.NewRecorder()

	h.ListLatestInvoices(rec, req)
	require.Equal(t, 200, rec.Code)

	var out []domain.Invoice
	err := json.Unmarshal(rec.Body.Bytes(), &out)
	require.NoError(t, err)
	require.Len(t, out, 1)
}
