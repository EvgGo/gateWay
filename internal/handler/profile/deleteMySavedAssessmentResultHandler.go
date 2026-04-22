package profile

import (
	"fmt"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func DeleteMySavedAssessmentResultHandler(
	log *slog.Logger,
	c authv1.UserProfileClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "DeleteMySavedAssessmentResultHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		assessmentID, err := parsePositiveInt64Param(chi.URLParam(r, "assessment_id"))
		if err != nil {
			reqLog.Warn("Некорректный assessment_id", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid assessment_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на удаление сохранённого результата теста",
			"assessment_id", assessmentID,
		)

		resp, err := c.DeleteMySavedAssessmentResult(ctx, &authv1.DeleteMySavedAssessmentResultRequest{
			AssessmentId: assessmentID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC DeleteMySavedAssessmentResult", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Сохранённый результат теста успешно удалён",
			"assessment_id", assessmentID,
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func parsePositiveInt64Param(raw string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("must be positive int64")
	}
	return value, nil
}
