package middleware

import (
	"net/http"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/requestid"
)

const requestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(requestIDHeader)
		if rid == "" {
			rid = requestid.New()
		}

		w.Header().Set(requestIDHeader, rid)
		r = r.WithContext(requestid.WithContext(r.Context(), rid))
		next.ServeHTTP(w, r)
	})
}
