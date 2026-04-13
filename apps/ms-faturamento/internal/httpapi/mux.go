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
	mux.HandleFunc("GET /api/v1/notas/ultimas", handler.ListLatestInvoices)
	mux.HandleFunc("PUT /api/v1/notas/{id}", handler.UpdateInvoice)
	mux.HandleFunc("DELETE /api/v1/notas/{id}", handler.DeleteInvoice)
}

func RegisterPrintInvoiceRoutes(mux *http.ServeMux, handler *PrintInvoiceHandler) {
	mux.HandleFunc("POST /api/v1/notas/{id}/imprimir", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler.PrintInvoice(w, r, id)
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
