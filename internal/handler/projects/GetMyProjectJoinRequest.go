package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func GetMyProjectJoinRequestHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "GetMyProjectJoinRequestHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			reqLog.Warn("Не передан обязательный path-параметр project_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на получение моей заявки в проект",
			"project_id", projectID,
		)

		resp, err := c.GetMyProjectJoinRequest(
			ctx,
			&workspacev1.GetMyProjectJoinRequestRequest{
				ProjectId: projectID,
			},
		)
		if err != nil {
			reqLog.Error("Ошибка gRPC GetMyProjectJoinRequest", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Моя заявка в проект успешно получена",
			"project_id", projectID,
			"has_request", resp.GetRequest() != nil,
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func parsePageSize(raw string, defaultValue int32) (int32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultValue, nil
	}

	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, fmt.Errorf("page_size must be > 0")
	}
	if n > 100 {
		n = 100
	}

	return int32(n), nil
}
