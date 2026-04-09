package middleware

import (
	"net/http"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/httpapi"
)

func MaxBodyBytes(limit int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > limit {
			httpapi.WriteError(w, r, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "payload excede limite permitido", nil)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, limit)
		next.ServeHTTP(w, r)
	})
}
