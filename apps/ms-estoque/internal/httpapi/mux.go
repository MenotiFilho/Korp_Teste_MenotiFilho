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
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
