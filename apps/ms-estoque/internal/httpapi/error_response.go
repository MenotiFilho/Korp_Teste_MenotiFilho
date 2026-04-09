package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/MenotiFilho/Korp_Teste_MenotiFilho/apps/ms-estoque/internal/requestid"
)

type ErrorResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details"`
	RequestID string      `json:"request_id"`
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string, details interface{}) {
	out := ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: requestid.FromContext(r.Context()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(out)
}
