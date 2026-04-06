package helpers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

// ExtractProjectID получает projectID из параметра маршрута chi или из последнего сегмента пути URL
// Возвращает пустую строку, если идентификатор не найден
func ExtractProjectID(r *http.Request) string {

	// Пытаемся получить параметр из chi (например, /projects/{project_id})
	projectID := chi.URLParam(r, "project_id")
	if projectID != "" {
		return projectID
	}

	// Если параметра нет, извлекаем последний сегмент пути
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
