package profile

import (
	"encoding/json"
	"fmt"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"net/http"
	"strings"
)

func parseCompetenceLevels(raw json.RawMessage) (map[string]string, error) {

	trimmed := strings.TrimSpace(string(raw))

	// null / пусто => очищаем карту
	if trimmed == "" || trimmed == "null" {
		return map[string]string{}, nil
	}

	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]string{}, nil
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		key := strings.TrimSpace(k)
		if key == "" {
			return nil, fmt.Errorf("competence_levels contains empty key")
		}
		out[key] = strings.TrimSpace(v)
	}

	return out, nil
}

func parseSkillsSelection(raw json.RawMessage) (*authv1.SkillSelection, error) {

	trimmed := strings.TrimSpace(string(raw))

	// null / пусто => очистить все skills
	if trimmed == "" || trimmed == "null" {
		return &authv1.SkillSelection{Ids: []string{}}, nil
	}

	// Разрешаем и bare array: ["1","2"]
	if strings.HasPrefix(trimmed, "[") {
		var ids []string
		if err := json.Unmarshal(raw, &ids); err != nil {
			return nil, fmt.Errorf("skills must be an object {\"ids\": [...]} or an array of strings")
		}
		norm, err := normalizeSkillIDs(ids)
		if err != nil {
			return nil, err
		}
		return &authv1.SkillSelection{Ids: norm}, nil
	}

	// Основной вариант: {"ids":[...]}
	var payload struct {
		Ids []string `json:"ids"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("skills must be an object with field ids")
	}

	norm, err := normalizeSkillIDs(payload.Ids)
	if err != nil {
		return nil, err
	}

	return &authv1.SkillSelection{Ids: norm}, nil
}

func normalizeSkillIDs(ids []string) ([]string, error) {

	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))

	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" {
			return nil, fmt.Errorf("skills.ids must not contain empty values")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}

	if len(out) > 30 {
		return nil, fmt.Errorf("maximum 30 skills allowed")
	}

	return out, nil
}

func parseSkillIDsFromQuery(r *http.Request) ([]string, error) {

	values := r.URL.Query()["skill_ids"]
	if len(values) == 0 {
		return nil, nil
	}

	rawIDs := make([]string, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			rawIDs = append(rawIDs, part)
		}
	}

	return normalizeSkillIDs(rawIDs)
}
