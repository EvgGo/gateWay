package profile

import (
	"encoding/json"
	"fmt"
	ssov1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func parseSkillsSelection(raw json.RawMessage) (*ssov1.SkillSelection, error) {

	trimmed := strings.TrimSpace(string(raw))

	// null / пусто => очистить все skills
	if trimmed == "" || trimmed == "null" {
		return &ssov1.SkillSelection{Ids: []string{}}, nil
	}

	if strings.HasPrefix(trimmed, "[") {
		var ids []string
		if err := json.Unmarshal(raw, &ids); err != nil {
			return nil, fmt.Errorf("skills must be an object {\"ids\": [...]} or an array of strings")
		}
		norm, err := normalizeSkillIDs(ids)
		if err != nil {
			return nil, err
		}
		return &ssov1.SkillSelection{Ids: norm}, nil
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

	return &ssov1.SkillSelection{Ids: norm}, nil
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

func parseSkillIDs(q map[string][]string) []string {
	rawValues := q["skill_ids"]
	if len(rawValues) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(rawValues))
	out := make([]string, 0, len(rawValues))

	for _, raw := range rawValues {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if _, ok := seen[part]; ok {
				continue
			}
			seen[part] = struct{}{}
			out = append(out, part)
		}
	}

	return out
}

func parseOptionalInt32(v string) (int32, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}

	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "page_size must be a valid int32")
	}

	return int32(n), nil
}

func parseOptionalBool(v string) (bool, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return false, nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, status.Error(codes.InvalidArgument, "open_suggestions_only must be a boolean")
	}

	return b, nil
}

func parseUserSkillMatchMode(v string) (ssov1.UserSkillMatchMode, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "":
		return ssov1.UserSkillMatchMode_USER_SKILL_MATCH_MODE_UNSPECIFIED, nil
	case "any", "user_skill_match_mode_any":
		return ssov1.UserSkillMatchMode_USER_SKILL_MATCH_MODE_ANY, nil
	case "all", "user_skill_match_mode_all":
		return ssov1.UserSkillMatchMode_USER_SKILL_MATCH_MODE_ALL, nil
	default:
		return ssov1.UserSkillMatchMode_USER_SKILL_MATCH_MODE_UNSPECIFIED,
			status.Error(codes.InvalidArgument, "skill_match_mode must be one of: any, all")
	}
}

func parsePublicUserSortBy(v string) (ssov1.PublicUserSortBy, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "":
		return ssov1.PublicUserSortBy_PUBLIC_USER_SORT_BY_UNSPECIFIED, nil
	case "created_at", "public_user_sort_by_created_at":
		return ssov1.PublicUserSortBy_PUBLIC_USER_SORT_BY_CREATED_AT, nil
	case "profile_skill_match_percent", "public_user_sort_by_profile_skill_match_percent":
		return ssov1.PublicUserSortBy_PUBLIC_USER_SORT_BY_PROFILE_SKILL_MATCH_PERCENT, nil
	case "matched_skills_count", "public_user_sort_by_matched_skills_count":
		return ssov1.PublicUserSortBy_PUBLIC_USER_SORT_BY_MATCHED_SKILLS_COUNT, nil
	default:
		return ssov1.PublicUserSortBy_PUBLIC_USER_SORT_BY_UNSPECIFIED,
			status.Error(codes.InvalidArgument, "sort_by must be one of: created_at, profile_skill_match_percent, matched_skills_count")
	}
}

func parseSortOrder(v string) (ssov1.SortOrder, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "":
		return ssov1.SortOrder_SORT_ORDER_UNSPECIFIED, nil
	case "asc", "sort_order_asc":
		return ssov1.SortOrder_SORT_ORDER_ASC, nil
	case "desc", "sort_order_desc":
		return ssov1.SortOrder_SORT_ORDER_DESC, nil
	default:
		return ssov1.SortOrder_SORT_ORDER_UNSPECIFIED,
			status.Error(codes.InvalidArgument, "sort_order must be one of: asc, desc")
	}
}

func parseOptionalTimestamp(v string) (*timestamppb.Timestamp, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "timestamp must be RFC3339, for example 2026-04-09T12:00:00Z")
	}

	return timestamppb.New(t), nil
}
