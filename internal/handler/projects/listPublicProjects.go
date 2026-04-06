package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"
	"strconv"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

func ListPublicProjectsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "ListPublicProjectsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		q := r.URL.Query()

		query := strings.TrimSpace(q.Get("query"))
		statusRaw := strings.TrimSpace(q.Get("status"))
		pageToken := strings.TrimSpace(q.Get("page_token"))
		skillMatchModeRaw := strings.TrimSpace(q.Get("skill_match_mode"))
		sortByRaw := strings.TrimSpace(q.Get("sort_by"))
		sortOrderRaw := strings.TrimSpace(q.Get("sort_order"))

		pageSize := int32(10) // дефолт
		if raw := strings.TrimSpace(q.Get("page_size")); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n <= 0 {
				reqLog.Warn("Некорректный query-параметр page_size", "page_size", raw, "err", err)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "page_size must be a positive integer")
				return
			}

			if n > 100 {
				reqLog.Debug("page_size превышает допустимый лимит, значение будет ограничено", "page_size_raw", n, "page_size_clamped", 100)
				n = 100
			}

			pageSize = int32(n)
		}

		convertedStatus, err := helpers.ParseProjectStatusParam(statusRaw)
		if err != nil {
			reqLog.Warn("Некорректный query-параметр status", "status", statusRaw, "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid project status")
			return
		}

		skillIDs := parseSkillIDsQuery(q["skill_ids"])

		skillMatchMode, err := parseProjectSkillMatchModeParam(skillMatchModeRaw)
		if err != nil {
			reqLog.Warn("Некорректный query-параметр skill_match_mode", "skill_match_mode", skillMatchModeRaw, "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid skill_match_mode")
			return
		}

		sortBy, err := parseProjectPublicSortByParam(sortByRaw)
		if err != nil {
			reqLog.Warn("Некорректный query-параметр sort_by", "sort_by", sortByRaw, "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid sort_by")
			return
		}

		sortOrder, err := parseSortOrderParam(sortOrderRaw)
		if err != nil {
			reqLog.Warn("Некорректный query-параметр sort_order", "sort_order", sortOrderRaw, "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid sort_order")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info(
			"Получен HTTP-запрос на получение публичного списка проектов",
			"query", query,
			"status_raw", statusRaw,
			"page_size", pageSize,
			"page_token", pageToken,
			"skill_ids", skillIDs,
			"skill_match_mode_raw", skillMatchModeRaw,
			"skill_match_mode", skillMatchMode.String(),
			"sort_by_raw", sortByRaw,
			"sort_by", sortBy.String(),
			"sort_order_raw", sortOrderRaw,
			"sort_order", sortOrder.String(),
		)

		resp, err := c.ListPublicProjects(ctx, &workspacev1.ListPublicProjectsRequest{
			Query:          query,
			Status:         convertedStatus,
			PageSize:       pageSize,
			PageToken:      pageToken,
			SkillIds:       skillIDs,
			SkillMatchMode: skillMatchMode,
			SortBy:         sortBy,
			SortOrder:      sortOrder,
		})
		if err != nil {
			reqLog.Warn("ListPublicProjects gRPC вызов завершился ошибкой", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Публичный список проектов успешно получен",
			"projects_count", len(resp.GetProjects()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func parseSkillIDsQuery(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, raw := range values {
		if raw == "" {
			continue
		}

		// поддержка обоих форматов:
		// ?skill_ids=1&skill_ids=2
		// ?skill_ids=1,2
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			id := strings.TrimSpace(part)
			if id == "" {
				continue
			}
			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func parseProjectSkillMatchModeParam(raw string) (workspacev1.ProjectSkillMatchMode, error) {
	if raw == "" {
		return workspacev1.ProjectSkillMatchMode_PROJECT_SKILL_MATCH_MODE_UNSPECIFIED, nil
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "0", "UNSPECIFIED", "PROJECT_SKILL_MATCH_MODE_UNSPECIFIED":
		return workspacev1.ProjectSkillMatchMode_PROJECT_SKILL_MATCH_MODE_UNSPECIFIED, nil
	case "1", "ANY", "PROJECT_SKILL_MATCH_MODE_ANY":
		return workspacev1.ProjectSkillMatchMode_PROJECT_SKILL_MATCH_MODE_ANY, nil
	case "2", "ALL", "PROJECT_SKILL_MATCH_MODE_ALL":
		return workspacev1.ProjectSkillMatchMode_PROJECT_SKILL_MATCH_MODE_ALL, nil
	default:
		return workspacev1.ProjectSkillMatchMode_PROJECT_SKILL_MATCH_MODE_UNSPECIFIED,
			fmt.Errorf("unknown skill_match_mode: %q", raw)
	}
}

func parseProjectPublicSortByParam(raw string) (workspacev1.ProjectPublicSortBy, error) {
	if raw == "" {
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_UNSPECIFIED, nil
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "0", "UNSPECIFIED", "PROJECT_PUBLIC_SORT_BY_UNSPECIFIED":
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_UNSPECIFIED, nil
	case "1", "CREATED_AT", "PROJECT_PUBLIC_SORT_BY_CREATED_AT":
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_CREATED_AT, nil
	case "2", "STARTED_AT", "PROJECT_PUBLIC_SORT_BY_STARTED_AT":
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_STARTED_AT, nil
	case "3", "PROFILE_SKILL_MATCH", "PROJECT_PUBLIC_SORT_BY_PROFILE_SKILL_MATCH":
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_PROFILE_SKILL_MATCH, nil
	default:
		return workspacev1.ProjectPublicSortBy_PROJECT_PUBLIC_SORT_BY_UNSPECIFIED,
			fmt.Errorf("unknown sort_by: %q", raw)
	}
}

func parseSortOrderParam(raw string) (workspacev1.SortOrder, error) {
	if raw == "" {
		return workspacev1.SortOrder_SORT_ORDER_UNSPECIFIED, nil
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "0", "UNSPECIFIED", "SORT_ORDER_UNSPECIFIED":
		return workspacev1.SortOrder_SORT_ORDER_UNSPECIFIED, nil
	case "1", "ASC", "SORT_ORDER_ASC":
		return workspacev1.SortOrder_SORT_ORDER_ASC, nil
	case "2", "DESC", "SORT_ORDER_DESC":
		return workspacev1.SortOrder_SORT_ORDER_DESC, nil
	default:
		return workspacev1.SortOrder_SORT_ORDER_UNSPECIFIED,
			fmt.Errorf("unknown sort_order: %q", raw)
	}
}
