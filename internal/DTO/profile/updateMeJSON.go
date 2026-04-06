package profile

import "encoding/json"

// UpdateMeJSON - DTO для PATCH /users/me.
// Используем *string/*bool, чтобы различать:
// - поле отсутствует в JSON => не менять,
// - поле присутствует и пустое/false => явно установить
//
// competence_levels tricky:
// - в proto это map, но map не имеет presence как optional,
// - в PATCH важно различать поле не присылали и прислали {} => очистить
//
// Поэтому мы держим json.RawMessage - если ключ отсутствовал, RawMessage будет nil
type UpdateMeJSON struct {
	FirstName             *string         `json:"first_name,omitempty"`
	LastName              *string         `json:"last_name,omitempty"`
	Phone                 *string         `json:"phone,omitempty"`
	About                 *string         `json:"about,omitempty"`
	CompetenceLevels      json.RawMessage `json:"competence_levels,omitempty"` // различаем нет поля и {}
	IsUserOpenSuggestions *bool           `json:"is_user_open_suggestions,omitempty"`
	IsProfileHidden       *bool           `json:"is_profile_hidden,omitempty"`
	Skills                json.RawMessage `json:"skill_ids,omitempty"`
	SkillsSet             bool            `json:"skills_set"`
}
