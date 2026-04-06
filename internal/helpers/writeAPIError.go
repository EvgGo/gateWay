package helpers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// apiError - единый формат ошибок gateway (HTTP)
// request_id полезен клиенту/саппорту: по нему легко найти запрос в логах/трейсах
type apiError struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

// WriteAPIError пишет унифицированный JSON-ответ ошибки
// - request_id всегда добавляем, если есть
func WriteAPIError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(apiError{
		Error:     msg,
		RequestID: middleware.GetReqID(r.Context()),
	})
}
