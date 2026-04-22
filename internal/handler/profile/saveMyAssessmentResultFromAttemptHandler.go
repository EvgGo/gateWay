package profile

import (
	"encoding/json"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
)

type saveMyAssessmentResultFromAttemptPayload struct {
	AttemptID int64 `json:"attempt_id"`
}

func SaveMyAssessmentResultFromAttemptHandler(
	log *slog.Logger,
	c authv1.UserProfileClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "SaveMyAssessmentResultFromAttemptHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		var body saveMyAssessmentResultFromAttemptPayload

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&body); err != nil {
			reqLog.Warn("Некорректный JSON body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid json body")
			return
		}

		if body.AttemptID <= 0 {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "attempt_id must be positive")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на сохранение результата попытки в профиль",
			"attempt_id", body.AttemptID,
		)

		resp, err := c.SaveMyAssessmentResultFromAttempt(ctx, &authv1.SaveMyAssessmentResultFromAttemptRequest{
			AttemptId: body.AttemptID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC SaveMyAssessmentResultFromAttempt", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Результат попытки успешно сохранён в профиль",
			"attempt_id", resp.GetAttemptId(),
			"assessment_id", resp.GetAssessmentId(),
			"subject_id", resp.GetSubjectId(),
			"level", resp.GetLevel(),
			"mode", resp.GetMode().String(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
