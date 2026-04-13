package httpapi

import "net/http"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	return mux
}

func RegisterProductRoutes(mux *http.ServeMux, handler *ProductHandler) {
	mux.HandleFunc("POST /api/v1/produtos", handler.CreateProduct)
	mux.HandleFunc("GET /api/v1/produtos", handler.ListProducts)
	mux.HandleFunc("PUT /api/v1/produtos/{id}", handler.UpdateProduct)
	mux.HandleFunc("DELETE /api/v1/produtos/{id}", handler.DeleteProduct)
}

func RegisterStockRoutes(mux *http.ServeMux, handler *StockHandler) {
	mux.HandleFunc("POST /api/v1/estoque/baixa", handler.DecreaseStock)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
