package helpers

import (
	"net/http"
	"strings"
)

// BearerFromRequest достает Authorization header и строго проверяет формат:
//
//	Authorization: Bearer <token>
//
// Возвращает:
// - исходную строку заголовка целиком ("Bearer <token>") - удобно прокидывать в gRPC metadata без пересборки,
// - либо "" если заголовка нет / формат неверный.
func BearerFromRequest(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return ""
	}

	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	// EqualFold - case-insensitive ("bearer", "BEARER"), чтобы быть терпимее к клиентам
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	// Токен не должен быть пустым
	if strings.TrimSpace(parts[1]) == "" {
		return ""
	}

	// Возвращаем как есть. Это полезно, если downstream ожидает именно токен
	return h
}
