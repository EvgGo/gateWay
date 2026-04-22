package profile

import (
	"fmt"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ListPublicCandidatesHandler(
	log *slog.Logger,
	c authv1.UserProfileClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "ListPublicCandidatesHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		req, err := buildListPublicCandidatesRequest(r)
		if err != nil {
			reqLog.Warn("Некорректные query params для ListPublicCandidates", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на поиск публичных кандидатов",
			"query", req.GetQuery(),
			"page_size", req.GetPageSize(),
			"page_token", req.GetPageToken(),
			"skill_ids_count", len(req.GetSkillIds()),
			"assessment_filters_count", len(req.GetAssessmentFilters()),
			"skill_match_mode", req.GetSkillMatchMode().String(),
			"assessment_match_mode", req.GetAssessmentMatchMode().String(),
			"open_suggestions_only", req.GetOpenSuggestionsOnly(),
		)

		resp, err := c.ListPublicCandidates(ctx, req)
		if err != nil {
			reqLog.Error("Ошибка gRPC ListPublicCandidates", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Публичные кандидаты успешно получены",
			"count", len(resp.GetCandidates()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

// buildListPublicCandidatesRequest вариант под 1 предмет
func buildListPublicCandidatesRequest(r *http.Request) (*authv1.ListPublicCandidatesRequest, error) {
	q := r.URL.Query()

	pageSize, err := parseOptionalInt32(q.Get("page_size"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page_size")
	}

	openSuggestionsOnly, err := parseOptionalBool(q.Get("open_suggestions_only"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid open_suggestions_only")
	}

	createdFrom, err := parseOptionalTimestamp(q.Get("created_from"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid created_from")
	}

	createdTo, err := parseOptionalTimestamp(q.Get("created_to"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid created_to")
	}

	if createdFrom != nil && createdTo != nil && createdFrom.AsTime().After(createdTo.AsTime()) {
		return nil, status.Error(codes.InvalidArgument, "created_from must be <= created_to")
	}

	skillMatchMode, err := parseUserSkillMatchMode(q.Get("skill_match_mode"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid skill_match_mode")
	}

	sortBy, err := parsePublicUserSortBy(q.Get("sort_by"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sort_by")
	}

	sortOrder, err := parseSortOrder(q.Get("sort_order"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid sort_order")
	}

	assessmentFilters, err := parseCandidateAssessmentFilters(q)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assessment_filter")
	}

	assessmentMatchMode, err := parseCandidateAssessmentMatchMode(
		q.Get("assessment_match_mode"),
		len(assessmentFilters) > 0,
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assessment_match_mode")
	}

	req := &authv1.ListPublicCandidatesRequest{
		Query:               strings.TrimSpace(q.Get("query")),
		PageSize:            pageSize,
		PageToken:           strings.TrimSpace(q.Get("page_token")),
		SkillIds:            parseSkillIDs(q),
		SkillMatchMode:      skillMatchMode,
		CreatedFrom:         createdFrom,
		CreatedTo:           createdTo,
		OpenSuggestionsOnly: openSuggestionsOnly,
		SortBy:              sortBy,
		SortOrder:           sortOrder,
		AssessmentFilters:   assessmentFilters,
		AssessmentMatchMode: assessmentMatchMode,
	}

	return req, nil
}

func parseCandidateAssessmentMatchMode(
	raw string,
	hasAssessmentFilters bool,
) (authv1.CandidateAssessmentMatchMode, error) {

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "":
		if hasAssessmentFilters {
			return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_ALL, nil
		}
		return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_UNSPECIFIED, nil
	case "unspecified":
		return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_UNSPECIFIED, nil
	case "any":
		return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_ANY, nil
	case "all":
		return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_ALL, nil
	default:
		return authv1.CandidateAssessmentMatchMode_CANDIDATE_ASSESSMENT_MATCH_MODE_UNSPECIFIED, fmt.Errorf("unknown assessment_match_mode")
	}
}

func parseCandidateAssessmentFilters(q url.Values) ([]*authv1.CandidateAssessmentFilter, error) {
	assessmentIDRaw := strings.TrimSpace(q.Get("assessment_id"))
	minLevelRaw := strings.TrimSpace(q.Get("assessment_min_level"))

	if assessmentIDRaw != "" {
		assessmentID, err := strconv.ParseInt(assessmentIDRaw, 10, 64)
		if err != nil || assessmentID <= 0 {
			return nil, fmt.Errorf("invalid assessment_id")
		}

		minLevel := int32(1)
		if minLevelRaw != "" {
			minLevel64, err := strconv.ParseInt(minLevelRaw, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid assessment_min_level")
			}

			minLevel = int32(minLevel64)
		}

		if minLevel < 1 || minLevel > 5 {
			return nil, fmt.Errorf("assessment_min_level must be in range 1..5")
		}

		return []*authv1.CandidateAssessmentFilter{
			{
				AssessmentId: assessmentID,
				MinLevel:     minLevel,
			},
		}, nil
	}

	rawItems := make([]string, 0)
	rawItems = append(rawItems, q["assessment_filter"]...)
	rawItems = append(rawItems, q["assessment_filters"]...)

	if len(rawItems) == 0 {
		return nil, nil
	}

	seen := make(map[int64]struct{}, len(rawItems))
	out := make([]*authv1.CandidateAssessmentFilter, 0, len(rawItems))

	for _, raw := range rawItems {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		parts := strings.Split(raw, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid assessment filter format")
		}

		assessmentID, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil || assessmentID <= 0 {
			return nil, fmt.Errorf("invalid assessment_id")
		}

		minLevel64, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid min_level")
		}

		minLevel := int32(minLevel64)
		if minLevel < 1 || minLevel > 5 {
			return nil, fmt.Errorf("min_level must be in range 1..5")
		}

		if _, exists := seen[assessmentID]; exists {
			continue
		}
		seen[assessmentID] = struct{}{}

		out = append(out, &authv1.CandidateAssessmentFilter{
			AssessmentId: assessmentID,
			MinLevel:     minLevel,
		})
	}

	return out, nil
}
