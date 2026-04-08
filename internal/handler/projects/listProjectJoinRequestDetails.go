package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func ListProjectJoinRequestDetailsHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "ListProjectJoinRequestDetailsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			reqLog.Warn("Не передан обязательный path-параметр project_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		pageSize, err := parsePageSize(r.URL.Query().Get("page_size"), 20)
		if err != nil {
			reqLog.Warn("Некорректный page_size", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		statusValue, err := helpers.ParseJoinRequestStatusOrAll(r.URL.Query().Get("status"))
		if err != nil {
			reqLog.Warn("Некорректный status", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid status")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на список детальных заявок проекта",
			"project_id", projectID,
			"status", statusValue.String(),
			"page_size", pageSize,
			"page_token", pageToken,
		)

		resp, err := c.ListProjectJoinRequestDetails(
			ctx,
			&workspacev1.ListProjectJoinRequestDetailsRequest{
				ProjectId: projectID,
				Status:    statusValue,
				PageSize:  pageSize,
				PageToken: pageToken,
			},
		)
		if err != nil {
			reqLog.Error("Ошибка gRPC ListProjectJoinRequestDetails", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Список детальных заявок проекта успешно получен",
			"project_id", projectID,
			"requests_count", len(resp.GetRequests()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
