package helpers

import (
	"net/http"
)

// RequireBearerAuth - middleware, который требует наличие Authorization: Bearer <token>.
//   - НЕ валидируем JWT здесь (это задача Auth/UserProfile сервисов),
//     gateway проверяет только наличие заголовка, чтобы быстро отсекать пустые запросы
//   - Ошибка возвращается в унифицированном JSON формате через writeAPIError
func RequireBearerAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// BearerFromRequest вернет строку вида Bearer xxx или ""
		if BearerFromRequest(r) == "" {
			WriteAPIError(w, r, http.StatusUnauthorized, "missing Authorization: Bearer token")
			return
		}

		// Если токен есть - пропускаем дальше
		next.ServeHTTP(w, r)
	})
}
