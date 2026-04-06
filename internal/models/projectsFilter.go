package models

import (
	"github.com/google/uuid"
	"time"
)

type ProjectsFilter struct {
	TeamID    string
	CreatorID string
	Status    string
	UserID    string
	OnlyOpen  bool
	Query     string

	// Пагинация
	PageSize int
	// Курсор: значение created_at и id последнего элемента предыдущей страницы
	// Кодируется в строку pageToken
	Cursor *ProjectCursor
}

// ProjectCursor представляет декодированный курсор
type ProjectCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID // или string, в зависимости от типа id
}
