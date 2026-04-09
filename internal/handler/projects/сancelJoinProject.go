package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func CancelJoinProjectHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "CancelJoinProjectHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		requestID := strings.TrimSpace(chi.URLParam(r, "request_id"))
		if requestID == "" {
			reqLog.Warn("Не передан обязательный path-параметр request_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "request_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на отмену своей заявки в проект",
			"request_id", requestID,
		)

		resp, err := c.CancelJoinProject(
			ctx,
			&workspacev1.CancelJoinProjectRequest{
				RequestId: requestID,
			},
		)
		if err != nil {
			reqLog.Error("Ошибка gRPC CancelJoinProject", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Заявка успешно отменена",
			"request_id", resp.GetId(),
			"status", resp.GetStatus(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
