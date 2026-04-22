package tests

import (
	"fmt"
	"gateWay/internal/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"

	assessmentv1 "github.com/EvgGo/proto/proto/gen/go/tests"
	"github.com/go-chi/chi/v5"
)

func ListSubjectsHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ListSubjectsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		query := strings.TrimSpace(r.URL.Query().Get("query"))
		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"), 20)
		if err != nil {
			reqLog.Warn("Некорректный page_size", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список предметов",
			"query", query,
			"page_size", pageSize,
			"page_token", pageToken,
		)

		resp, err := c.ListSubjects(ctx, &assessmentv1.ListSubjectsRequest{
			Query:     query,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC ListSubjects", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Список предметов успешно получен",
			"items_count", len(resp.GetSubjects()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListSubtopicsHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ListSubtopicsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		query := strings.TrimSpace(r.URL.Query().Get("query"))
		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		subjectID, err := parseInt64Query(r.URL.Query().Get("subject_id"), 0)
		if err != nil {
			reqLog.Warn("Некорректный subject_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid subject_id")
			return
		}

		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"), 20)
		if err != nil {
			reqLog.Warn("Некорректный page_size", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список подтем",
			"subject_id", subjectID,
			"query", query,
			"page_size", pageSize,
			"page_token", pageToken,
		)

		resp, err := c.ListSubtopics(ctx, &assessmentv1.ListSubtopicsRequest{
			SubjectId: subjectID,
			Query:     query,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC ListSubtopics", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Список подтем успешно получен",
			"items_count", len(resp.GetSubtopics()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListAssessmentsHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ListAssessmentsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		req, err := buildListAssessmentsRequest(r)
		if err != nil {
			reqLog.Warn("Некорректные query params для ListAssessments", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список тестов",
			"subject_id", req.GetSubjectId(),
			"status", req.GetStatus().String(),
			"mode", req.GetMode().String(),
			"query", req.GetQuery(),
			"page_size", req.GetPageSize(),
			"page_token", req.GetPageToken(),
		)

		resp, err := c.ListAssessments(ctx, req)
		if err != nil {
			reqLog.Error("Ошибка gRPC ListAssessments", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		globalCount := 0
		subtopicCount := 0

		for _, item := range resp.GetAssessments() {
			switch item.GetMode() {
			case assessmentv1.AssessmentMode_ASSESSMENT_MODE_GLOBAL:
				globalCount++
			case assessmentv1.AssessmentMode_ASSESSMENT_MODE_SUBTOPIC:
				subtopicCount++
			}
		}

		reqLog.Info("Список тестов успешно получен",
			"items_count", len(resp.GetAssessments()),
			"global_count", globalCount,
			"subtopic_count", subtopicCount,
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func buildListAssessmentsRequest(r *http.Request) (*assessmentv1.ListAssessmentsRequest, error) {
	q := r.URL.Query()

	subjectID, err := parseInt64Query(q.Get("subject_id"), 0)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subject_id")
	}

	pageSize, err := parsePageSize(q.Get("page_size"), 20)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page_size")
	}

	statusValue, err := parseAssessmentStatus(q.Get("status"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	modeValue, err := parseAssessmentMode(q.Get("mode"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid mode")
	}

	req := &assessmentv1.ListAssessmentsRequest{
		SubjectId: subjectID,
		Query:     strings.TrimSpace(q.Get("query")),
		Status:    statusValue,
		PageSize:  pageSize,
		PageToken: strings.TrimSpace(q.Get("page_token")),
		Mode:      modeValue,
	}

	return req, nil
}

func parseAssessmentMode(raw string) (assessmentv1.AssessmentMode, error) {
	value := strings.TrimSpace(raw)

	switch strings.ToLower(value) {
	case "", "unspecified", "assessment_mode_unspecified", "0":
		return assessmentv1.AssessmentMode_ASSESSMENT_MODE_UNSPECIFIED, nil

	case "subtopic", "assessment_mode_subtopic", "1":
		return assessmentv1.AssessmentMode_ASSESSMENT_MODE_SUBTOPIC, nil

	case "global", "assessment_mode_global", "2":
		return assessmentv1.AssessmentMode_ASSESSMENT_MODE_GLOBAL, nil

	default:
		return assessmentv1.AssessmentMode_ASSESSMENT_MODE_UNSPECIFIED, fmt.Errorf("unknown assessment mode")
	}
}

func GetAssessmentHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "GetAssessmentHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		assessmentID, err := parsePathInt64(chi.URLParam(r, "assessment_id"))
		if err != nil {
			reqLog.Warn("Некорректный assessment_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid assessment_id")
			return
		}

		includeSubtopics, err := parseOptionalBool(r.URL.Query().Get("include_subtopics"), true)
		if err != nil {
			reqLog.Warn("Некорректный include_subtopics", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid include_subtopics")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на получение теста",
			"assessment_id", assessmentID,
			"include_subtopics", includeSubtopics,
		)

		resp, err := c.GetAssessment(ctx, &assessmentv1.GetAssessmentRequest{
			AssessmentId:     assessmentID,
			IncludeSubtopics: includeSubtopics,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetAssessment", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Тест успешно получен",
			"assessment_id", assessmentID,
			"assessment_mode", resp.GetMode().String(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
