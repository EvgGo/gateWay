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

type customAssessmentSubtopicBody struct {
	SubtopicID      int64   `json:"subtopic_id"`
	Weight          float64 `json:"weight"`
	Priority        int32   `json:"priority"`
	MinItems        int32   `json:"min_items"`
	MaxItems        int32   `json:"max_items"`
	StartDifficulty int32   `json:"start_difficulty"`
	StopConfidence  float64 `json:"stop_confidence"`
}

type createCustomAssessmentBody struct {
	SubjectID       int64                          `json:"subject_id"`
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	Visibility      string                         `json:"visibility"`
	Subtopics       []customAssessmentSubtopicBody `json:"subtopics"`
	MinItems        int32                          `json:"min_items"`
	MaxItems        int32                          `json:"max_items"`
	MinDifficulty   int32                          `json:"min_difficulty"`
	MaxDifficulty   int32                          `json:"max_difficulty"`
	StartDifficulty int32                          `json:"start_difficulty"`
	StopConfidence  float64                        `json:"stop_confidence"`
	DurationSeconds int32                          `json:"duration_seconds"`
}

type updateCustomAssessmentBody struct {
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	Visibility      string                         `json:"visibility"`
	Subtopics       []customAssessmentSubtopicBody `json:"subtopics"`
	MinItems        int32                          `json:"min_items"`
	MaxItems        int32                          `json:"max_items"`
	MinDifficulty   int32                          `json:"min_difficulty"`
	MaxDifficulty   int32                          `json:"max_difficulty"`
	StartDifficulty int32                          `json:"start_difficulty"`
	StopConfidence  float64                        `json:"stop_confidence"`
	DurationSeconds int32                          `json:"duration_seconds"`
}

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
			"kind", req.GetKind().String(),
			"visibility", req.GetVisibility().String(),
			"scope", req.GetScope().String(),
			"created_by_user_id", req.GetCreatedByUserId(),
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

	kindValue, err := parseAssessmentKind(q.Get("kind"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid kind")
	}

	visibilityValue, err := parseAssessmentVisibility(q.Get("visibility"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid visibility")
	}

	scopeValue, err := parseAssessmentListScope(q.Get("scope"))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid scope")
	}

	req := &assessmentv1.ListAssessmentsRequest{
		SubjectId:       subjectID,
		Query:           strings.TrimSpace(q.Get("query")),
		Status:          statusValue,
		PageSize:        pageSize,
		PageToken:       strings.TrimSpace(q.Get("page_token")),
		Mode:            modeValue,
		Kind:            kindValue,
		Visibility:      visibilityValue,
		Scope:           scopeValue,
		CreatedByUserId: strings.TrimSpace(q.Get("created_by_user_id")),
	}

	return req, nil
}

func CreateCustomAssessmentHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "CreateCustomAssessmentHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		var body createCustomAssessmentBody
		if err := decodeJSONBody(r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		visibilityValue, err := parseAssessmentVisibility(body.Visibility)
		if err != nil || visibilityValue == assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_UNSPECIFIED {
			reqLog.Warn("Некорректный visibility",
				"visibility", body.Visibility,
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid visibility")
			return
		}

		if body.SubjectID <= 0 {
			reqLog.Warn("Некорректный subject_id", "subject_id", body.SubjectID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid subject_id")
			return
		}

		if strings.TrimSpace(body.Title) == "" {
			reqLog.Warn("Пустой title")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "title is required")
			return
		}

		if len(body.Subtopics) == 0 {
			reqLog.Warn("Пустой список subtopics")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "subtopics are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на создание custom assessment",
			"subject_id", body.SubjectID,
			"title", body.Title,
			"visibility", visibilityValue.String(),
			"subtopics_count", len(body.Subtopics),
			"duration_seconds", body.DurationSeconds,
		)

		resp, err := c.CreateCustomAssessment(ctx, &assessmentv1.CreateCustomAssessmentRequest{
			SubjectId:       body.SubjectID,
			Title:           strings.TrimSpace(body.Title),
			Description:     strings.TrimSpace(body.Description),
			Visibility:      visibilityValue,
			Subtopics:       customAssessmentSubtopicsToProto(body.Subtopics),
			MinItems:        body.MinItems,
			MaxItems:        body.MaxItems,
			MinDifficulty:   body.MinDifficulty,
			MaxDifficulty:   body.MaxDifficulty,
			StartDifficulty: body.StartDifficulty,
			StopConfidence:  body.StopConfidence,
			DurationSeconds: body.DurationSeconds,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC CreateCustomAssessment",
				"err", err,
				"grpc_code", status.Code(err).String(),
				"grpc_message", status.Convert(err).Message(),
				"subject_id", body.SubjectID,
				"title", body.Title,
				"visibility", visibilityValue.String(),
				"subtopics_count", len(body.Subtopics),
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Custom assessment успешно создан",
			"assessment_id", resp.GetId(),
			"assessment_code", resp.GetCode(),
			"assessment_title", resp.GetTitle(),
			"assessment_kind", resp.GetKind().String(),
			"assessment_visibility", resp.GetVisibility().String(),
			"subject_id", resp.GetSubjectId(),
			"subject_code", resp.GetSubjectCode(),
			"subject_title", resp.GetSubjectTitle(),
			"is_editable", resp.GetIsEditable(),
			"can_use_in_project", resp.GetCanUseInProject(),
			"subtopics_count", len(resp.GetSubtopics()),
		)

		helpers.WriteProtoJSON(w, http.StatusCreated, resp)
	}
}

func UpdateCustomAssessmentHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "UpdateCustomAssessmentHandler",
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

		var body updateCustomAssessmentBody
		if err = decodeJSONBody(r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		visibilityValue, err := parseAssessmentVisibility(body.Visibility)
		if err != nil || visibilityValue == assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_UNSPECIFIED {
			reqLog.Warn("Некорректный visibility",
				"assessment_id", assessmentID,
				"visibility", body.Visibility,
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid visibility")
			return
		}

		if strings.TrimSpace(body.Title) == "" {
			reqLog.Warn("Пустой title", "assessment_id", assessmentID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "title is required")
			return
		}

		if len(body.Subtopics) == 0 {
			reqLog.Warn("Пустой список subtopics", "assessment_id", assessmentID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "subtopics are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на обновление custom assessment",
			"assessment_id", assessmentID,
			"title", body.Title,
			"visibility", visibilityValue.String(),
			"subtopics_count", len(body.Subtopics),
			"duration_seconds", body.DurationSeconds,
		)

		resp, err := c.UpdateCustomAssessment(ctx, &assessmentv1.UpdateCustomAssessmentRequest{
			AssessmentId:    assessmentID,
			Title:           strings.TrimSpace(body.Title),
			Description:     strings.TrimSpace(body.Description),
			Visibility:      visibilityValue,
			Subtopics:       customAssessmentSubtopicsToProto(body.Subtopics),
			MinItems:        body.MinItems,
			MaxItems:        body.MaxItems,
			MinDifficulty:   body.MinDifficulty,
			MaxDifficulty:   body.MaxDifficulty,
			StartDifficulty: body.StartDifficulty,
			StopConfidence:  body.StopConfidence,
			DurationSeconds: body.DurationSeconds,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC UpdateCustomAssessment",
				"err", err,
				"grpc_code", status.Code(err).String(),
				"grpc_message", status.Convert(err).Message(),
				"assessment_id", assessmentID,
				"title", body.Title,
				"visibility", visibilityValue.String(),
				"subtopics_count", len(body.Subtopics),
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Custom assessment успешно обновлён",
			"assessment_id", resp.GetId(),
			"assessment_code", resp.GetCode(),
			"assessment_title", resp.GetTitle(),
			"assessment_kind", resp.GetKind().String(),
			"assessment_visibility", resp.GetVisibility().String(),
			"subject_id", resp.GetSubjectId(),
			"subject_code", resp.GetSubjectCode(),
			"subject_title", resp.GetSubjectTitle(),
			"is_editable", resp.GetIsEditable(),
			"can_use_in_project", resp.GetCanUseInProject(),
			"subtopics_count", len(resp.GetSubtopics()),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ArchiveCustomAssessmentHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ArchiveCustomAssessmentHandler",
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

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на архивацию custom assessment",
			"assessment_id", assessmentID,
		)

		resp, err := c.ArchiveCustomAssessment(ctx, &assessmentv1.ArchiveCustomAssessmentRequest{
			AssessmentId: assessmentID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC ArchiveCustomAssessment",
				"err", err,
				"grpc_code", status.Code(err).String(),
				"grpc_message", status.Convert(err).Message(),
				"assessment_id", assessmentID,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Custom assessment успешно архивирован",
			"assessment_id", resp.GetId(),
			"assessment_code", resp.GetCode(),
			"assessment_title", resp.GetTitle(),
			"assessment_kind", resp.GetKind().String(),
			"assessment_visibility", resp.GetVisibility().String(),
			"assessment_status", resp.GetStatus().String(),
			"is_editable", resp.GetIsEditable(),
			"can_use_in_project", resp.GetCanUseInProject(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListMyCustomAssessmentsHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "ListMyCustomAssessmentsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		q := r.URL.Query()

		pageToken := strings.TrimSpace(q.Get("page_token"))

		pageSize, err := parsePageSize(q.Get("page_size"), 20)
		if err != nil {
			reqLog.Warn("Некорректный page_size", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		statusValue, err := parseAssessmentStatus(q.Get("status"))
		if err != nil {
			reqLog.Warn("Некорректный status", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid status")
			return
		}

		query := strings.TrimSpace(q.Get("query"))

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список моих custom assessment",
			"query", query,
			"status", statusValue.String(),
			"page_size", pageSize,
			"page_token", pageToken,
		)

		resp, err := c.ListMyCustomAssessments(ctx, &assessmentv1.ListMyCustomAssessmentsRequest{
			Query:     query,
			Status:    statusValue,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC ListMyCustomAssessments",
				"err", err,
				"grpc_code", status.Code(err).String(),
				"grpc_message", status.Convert(err).Message(),
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		editableCount := 0
		canUseInProjectCount := 0

		for _, item := range resp.GetAssessments() {
			if item.GetIsEditable() {
				editableCount++
			}

			if item.GetCanUseInProject() {
				canUseInProjectCount++
			}
		}

		reqLog.Info("Список моих custom assessment успешно получен",
			"items_count", len(resp.GetAssessments()),
			"editable_count", editableCount,
			"can_use_in_project_count", canUseInProjectCount,
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
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

func parseAssessmentKind(raw string) (assessmentv1.AssessmentKind, error) {
	value := strings.TrimSpace(strings.ToLower(raw))

	switch value {
	case "", "unspecified", "assessment_kind_unspecified", "0":
		return assessmentv1.AssessmentKind_ASSESSMENT_KIND_UNSPECIFIED, nil

	case "system", "assessment_kind_system", "1":
		return assessmentv1.AssessmentKind_ASSESSMENT_KIND_SYSTEM, nil

	case "custom", "assessment_kind_custom", "2":
		return assessmentv1.AssessmentKind_ASSESSMENT_KIND_CUSTOM, nil

	default:
		return assessmentv1.AssessmentKind_ASSESSMENT_KIND_UNSPECIFIED, fmt.Errorf("unknown assessment kind")
	}
}

func parseAssessmentVisibility(raw string) (assessmentv1.AssessmentVisibility, error) {
	value := strings.TrimSpace(strings.ToLower(raw))

	switch value {
	case "", "unspecified", "assessment_visibility_unspecified", "0":
		return assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_UNSPECIFIED, nil

	case "private", "assessment_visibility_private", "1":
		return assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_PRIVATE, nil

	case "project", "assessment_visibility_project", "2":
		return assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_PROJECT, nil

	case "public", "assessment_visibility_public", "3":
		return assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_PUBLIC, nil

	default:
		return assessmentv1.AssessmentVisibility_ASSESSMENT_VISIBILITY_UNSPECIFIED, fmt.Errorf("unknown assessment visibility")
	}
}

func parseAssessmentListScope(raw string) (assessmentv1.AssessmentListScope, error) {
	value := strings.TrimSpace(strings.ToLower(raw))

	switch value {
	case "", "unspecified", "assessment_list_scope_unspecified", "0":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_UNSPECIFIED, nil

	case "available", "assessment_list_scope_available", "1":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_AVAILABLE, nil

	case "system", "assessment_list_scope_system", "2":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_SYSTEM, nil

	case "my_custom", "my-custom", "assessment_list_scope_my_custom", "3":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_MY_CUSTOM, nil

	case "project_available", "project-available", "assessment_list_scope_project_available", "4":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_PROJECT_AVAILABLE, nil

	case "public_custom", "public-custom", "assessment_list_scope_public_custom", "5":
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_PUBLIC_CUSTOM, nil

	default:
		return assessmentv1.AssessmentListScope_ASSESSMENT_LIST_SCOPE_UNSPECIFIED, fmt.Errorf("unknown assessment list scope")
	}
}

func customAssessmentSubtopicsToProto(
	items []customAssessmentSubtopicBody,
) []*assessmentv1.CustomAssessmentSubtopicInput {
	result := make([]*assessmentv1.CustomAssessmentSubtopicInput, 0, len(items))

	for _, item := range items {
		result = append(result, &assessmentv1.CustomAssessmentSubtopicInput{
			SubtopicId:      item.SubtopicID,
			Weight:          item.Weight,
			Priority:        item.Priority,
			MinItems:        item.MinItems,
			MaxItems:        item.MaxItems,
			StartDifficulty: item.StartDifficulty,
			StopConfidence:  item.StopConfidence,
		})
	}

	return result
}
