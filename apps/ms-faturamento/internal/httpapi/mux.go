package httpapi

import "net/http"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	return mux
}

func RegisterInvoiceRoutes(mux *http.ServeMux, handler *InvoiceHandler) {
	mux.HandleFunc("POST /api/v1/notas", handler.CreateInvoice)
	mux.HandleFunc("GET /api/v1/notas", handler.ListInvoices)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
