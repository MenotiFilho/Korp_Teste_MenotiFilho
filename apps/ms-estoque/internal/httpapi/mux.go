package httpapi

import "net/http"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	// Fallback OPTIONS handler to ensure CORS preflight requests are handled
	// by the CORS middleware (rs/cors) and return a clean 204 when no route
	// explicitly matches the request. This is a minimal, non-invasive addition.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	})
	return mux
}

func RegisterProductRoutes(mux *http.ServeMux, handler *ProductHandler) {
	mux.HandleFunc("POST /api/v1/produtos", handler.CreateProduct)
	mux.HandleFunc("GET /api/v1/produtos", handler.ListProducts)
	mux.HandleFunc("GET /api/v1/produtos/baixo-estoque", handler.ListLowStockProducts)
	mux.HandleFunc("PUT /api/v1/produtos/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /api/v1/produtos/{id}", handler.DeleteProduct)
}

func RegisterStockRoutes(mux *http.ServeMux, handler *StockHandler) {
	mux.HandleFunc("POST /api/v1/estoque/baixa", handler.DecreaseStock)
	mux.HandleFunc("GET /api/v1/estoque/baixas/{key}", func(w http.ResponseWriter, r *http.Request) {
		handler.CheckIdempotencyKey(w, r)
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
