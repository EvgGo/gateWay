package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// WriteSuccess записывает успешный ответ с заданным HTTP статусом (например, 204 No Content)
func WriteSuccess(w http.ResponseWriter, statusCode int) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	// Пишем пустой JSON в случае No Content
	if statusCode == http.StatusNoContent {
		return
	}

	// Пишем простой успешный JSON-ответ
	response := map[string]string{
		"status": "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to write success response", "error", err)
	}
}
