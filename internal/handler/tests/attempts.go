package tests

import (
	"gateWay/internal/helpers"
	"log/slog"
	"net/http"
	"strings"

	assessmentv1 "github.com/EvgGo/proto/proto/gen/go/tests"
	"github.com/go-chi/chi/v5"
)

type startAttemptBody struct {
	AssessmentID        int64 `json:"assessment_id"`
	RestartIfInProgress bool  `json:"restart_if_in_progress"`
}

type submitAnswerBody struct {
	QuestionID       int64 `json:"question_id"`
	SelectedOptionID int64 `json:"selected_option_id"`
}

type finishAttemptBody struct {
	Action string `json:"action"`
}

func StartAttemptHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "StartAttemptHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		var body startAttemptBody
		if err := decodeJSONBody(r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		if body.AssessmentID <= 0 {
			reqLog.Warn("Некорректный assessment_id", "assessment_id", body.AssessmentID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid assessment_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на старт попытки",
			"assessment_id", body.AssessmentID,
			"restart_if_in_progress", body.RestartIfInProgress,
		)

		resp, err := c.StartAttempt(ctx, &assessmentv1.StartAttemptRequest{
			AssessmentId:        body.AssessmentID,
			RestartIfInProgress: body.RestartIfInProgress,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC StartAttempt", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		summary := resp.GetAttempt().GetAssessmentSummary()

		reqLog.Info("Попытка успешно создана",
			"attempt_id", resp.GetAttempt().GetId(),
			"assessment_id", summary.GetAssessmentId(),
			"assessment_code", summary.GetAssessmentCode(),
			"assessment_title", summary.GetAssessmentTitle(),
			"subject_id", summary.GetSubjectId(),
			"subject_code", summary.GetSubjectCode(),
			"subject_title", summary.GetSubjectTitle(),
			"assessment_mode", summary.GetMode().String(),
			"has_first_question", resp.GetFirstQuestion() != nil,
		)

		helpers.WriteProtoJSON(w, http.StatusCreated, resp)
	}
}

func GetAttemptHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "GetAttemptHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		attemptID, err := parsePathInt64(chi.URLParam(r, "attempt_id"))
		if err != nil {
			reqLog.Warn("Некорректный attempt_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid attempt_id")
			return
		}

		includeSubtopicStates, err := parseOptionalBool(r.URL.Query().Get("include_subtopic_states"), false)
		if err != nil {
			reqLog.Warn("Некорректный include_subtopic_states", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid include_subtopic_states")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на получение попытки",
			"attempt_id", attemptID,
			"include_subtopic_states", includeSubtopicStates,
		)

		resp, err := c.GetAttempt(ctx, &assessmentv1.GetAttemptRequest{
			AttemptId:             attemptID,
			IncludeSubtopicStates: includeSubtopicStates,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetAttempt", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		summary := resp.GetAssessmentSummary()

		reqLog.Info("Попытка успешно получена",
			"attempt_id", attemptID,
			"assessment_id", summary.GetAssessmentId(),
			"assessment_code", summary.GetAssessmentCode(),
			"assessment_title", summary.GetAssessmentTitle(),
			"subject_id", summary.GetSubjectId(),
			"subject_code", summary.GetSubjectCode(),
			"subject_title", summary.GetSubjectTitle(),
			"assessment_mode", summary.GetMode().String(),
			"status", resp.GetStatus().String(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func GetAttemptProgressHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "GetAttemptProgressHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		attemptID, err := parsePathInt64(chi.URLParam(r, "attempt_id"))
		if err != nil {
			reqLog.Warn("Некорректный attempt_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid attempt_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на прогресс попытки", "attempt_id", attemptID)

		resp, err := c.GetAttemptProgress(ctx, &assessmentv1.GetAttemptProgressRequest{
			AttemptId: attemptID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetAttemptProgress", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Прогресс попытки успешно получен", "attempt_id", attemptID)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func GetNextQuestionHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "GetNextQuestionHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		attemptID, err := parsePathInt64(chi.URLParam(r, "attempt_id"))
		if err != nil {
			reqLog.Warn("Некорректный attempt_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid attempt_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на следующий вопрос", "attempt_id", attemptID)

		resp, err := c.GetNextQuestion(ctx, &assessmentv1.GetNextQuestionRequest{
			AttemptId: attemptID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetNextQuestion", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Следующий вопрос успешно получен",
			"attempt_id", attemptID,
			"completed", resp.GetCompleted(),
			"has_next_question", resp.GetNextQuestion() != nil,
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func SubmitAnswerHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "SubmitAnswerHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		attemptID, err := parsePathInt64(chi.URLParam(r, "attempt_id"))
		if err != nil {
			reqLog.Warn("Некорректный attempt_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid attempt_id")
			return
		}

		var body submitAnswerBody
		if err = decodeJSONBody(r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		if body.QuestionID <= 0 {
			reqLog.Warn("Некорректный question_id", "question_id", body.QuestionID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid question_id")
			return
		}

		if body.SelectedOptionID <= 0 {
			reqLog.Warn("Некорректный selected_option_id", "selected_option_id", body.SelectedOptionID)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid selected_option_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на отправку ответа",
			"attempt_id", attemptID,
			"question_id", body.QuestionID,
			"selected_option_id", body.SelectedOptionID,
		)

		resp, err := c.SubmitAnswer(ctx, &assessmentv1.SubmitAnswerRequest{
			AttemptId:        attemptID,
			QuestionId:       body.QuestionID,
			SelectedOptionId: body.SelectedOptionID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC SubmitAnswer", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Ответ успешно принят",
			"attempt_id", attemptID,
			"completed", resp.GetCompleted(),
			"has_next_question", resp.GetNextQuestion() != nil,
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func FinishAttemptHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "FinishAttemptHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		attemptID, err := parsePathInt64(chi.URLParam(r, "attempt_id"))
		if err != nil {
			reqLog.Warn("Некорректный attempt_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid attempt_id")
			return
		}

		var body finishAttemptBody
		if err := decodeJSONBody(r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		actionValue, err := parseFinishAction(body.Action)
		if err != nil {
			reqLog.Warn("Некорректный action",
				"raw_action", body.Action,
				"normalized_action", strings.ToLower(strings.TrimSpace(body.Action)),
				"attempt_id", attemptID,
				"err", err,
			)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid action")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на завершение попытки",
			"attempt_id", attemptID,
			"action", strings.ToLower(strings.TrimSpace(body.Action)),
		)

		resp, err := c.FinishAttempt(ctx, &assessmentv1.FinishAttemptRequest{
			AttemptId: attemptID,
			Action:    actionValue,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC FinishAttempt", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		summary := resp.GetAssessmentSummary()

		reqLog.Info("Попытка успешно завершена",
			"attempt_id", attemptID,
			"assessment_id", summary.GetAssessmentId(),
			"assessment_code", summary.GetAssessmentCode(),
			"assessment_title", summary.GetAssessmentTitle(),
			"subject_id", summary.GetSubjectId(),
			"subject_code", summary.GetSubjectCode(),
			"subject_title", summary.GetSubjectTitle(),
			"assessment_mode", summary.GetMode().String(),
			"status", resp.GetStatus().String(),
			"finish_reason", resp.GetFinishReason().String(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListMyAttemptsHandler(
	log *slog.Logger,
	c assessmentv1.AdaptiveTestingClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "ListMyAttemptsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		assessmentID, err := parseInt64Query(r.URL.Query().Get("assessment_id"), 0)
		if err != nil {
			reqLog.Warn("Некорректный assessment_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid assessment_id")
			return
		}

		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"), 20)
		if err != nil {
			reqLog.Warn("Некорректный page_size", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		statusValue, err := parseAttemptStatus(r.URL.Query().Get("status"))
		if err != nil {
			reqLog.Warn("Некорректный status", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid status")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список моих попыток",
			"assessment_id", assessmentID,
			"status", statusValue.String(),
			"page_size", pageSize,
			"page_token", pageToken,
		)

		resp, err := c.ListMyAttempts(ctx, &assessmentv1.ListMyAttemptsRequest{
			AssessmentId: assessmentID,
			Status:       statusValue,
			PageSize:     pageSize,
			PageToken:    pageToken,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC ListMyAttempts", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		globalCount := 0
		subtopicCount := 0

		for _, item := range resp.GetAttempts() {
			switch item.GetAssessmentSummary().GetMode() {
			case assessmentv1.AssessmentMode_ASSESSMENT_MODE_GLOBAL:
				globalCount++
			case assessmentv1.AssessmentMode_ASSESSMENT_MODE_SUBTOPIC:
				subtopicCount++
			}
		}

		reqLog.Info("Список попыток успешно получен",
			"items_count", len(resp.GetAttempts()),
			"global_count", globalCount,
			"subtopic_count", subtopicCount,
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
