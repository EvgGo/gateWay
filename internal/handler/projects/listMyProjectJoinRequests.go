package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func ListMyProjectJoinRequestsHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"handler", "ListMyProjectJoinRequestsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
		)

		q := r.URL.Query()

		var (
			status    workspacev1.JoinRequestStatus
			pageSize  int32
			pageToken string
		)

		statusRaw := strings.TrimSpace(q.Get("status"))
		if statusRaw != "" {
			v, err := strconv.ParseInt(statusRaw, 10, 32)
			if err != nil {
				reqLog.Warn("Некорректный query-параметр status", "status", statusRaw, "err", err)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid status")
				return
			}

			if _, ok := workspacev1.JoinRequestStatus_name[int32(v)]; !ok {
				reqLog.Warn("Неизвестное значение status", "status", v)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "unknown status")
				return
			}

			status = workspacev1.JoinRequestStatus(v)
		}

		pageSizeRaw := strings.TrimSpace(q.Get("page_size"))
		if pageSizeRaw != "" {
			v, err := strconv.ParseInt(pageSizeRaw, 10, 32)
			if err != nil || v < 0 {
				reqLog.Warn("Некорректный query-параметр page_size", "page_size", pageSizeRaw, "err", err)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
				return
			}
			pageSize = int32(v)
		}

		pageToken = strings.TrimSpace(q.Get("page_token"))

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		grpcReq := &workspacev1.ListMyProjectJoinRequestsRequest{
			Status:    status,
			PageSize:  pageSize,
			PageToken: pageToken,
		}

		reqLog.Info(
			"Получен HTTP-запрос на список моих заявок в проекты",
			"status", grpcReq.GetStatus(),
			"page_size", grpcReq.GetPageSize(),
			"page_token", grpcReq.GetPageToken(),
		)

		resp, err := c.ListMyProjectJoinRequests(ctx, grpcReq)
		if err != nil {
			reqLog.Error("Ошибка gRPC ListMyProjectJoinRequests", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info(
			"Список моих заявок в проекты успешно получен",
			"items_count", len(resp.GetItems()),
			"next_page_token", resp.GetNextPageToken(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
