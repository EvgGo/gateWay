package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func ListManageableProjectJoinRequestBucketsHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusRaw := strings.TrimSpace(r.URL.Query().Get("status"))
		query := strings.TrimSpace(r.URL.Query().Get("query"))
		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		pageSize := int32(20)
		if raw := strings.TrimSpace(r.URL.Query().Get("page_size")); raw != "" {
			v, err := strconv.Atoi(raw)
			if err != nil || v <= 0 {
				log.Warn("некорректный query-параметр page_size",
					"handler", "ListManageableProjectJoinRequestBucketsHandler",
					"http_method", r.Method,
					"path", r.URL.Path,
					"page_size_raw", raw,
					"err", err,
				)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
				return
			}
			pageSize = int32(v)
		}

		reqLog := log.With(
			"handler", "ListManageableProjectJoinRequestBucketsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"status", statusRaw,
			"query", query,
			"page_size", pageSize,
			"page_token", pageToken,
		)

		reqLog.Info("Получен HTTP-запрос на список управляемых бакетов заявок в проекты")

		convertedStatus, err := helpers.ParseJoinRequestStatusParam(statusRaw)
		if err != nil {
			reqLog.Warn("некорректный query-параметр status", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid join request status")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListManageableProjectJoinRequestBuckets(ctx, &workspacev1.ListManageableProjectJoinRequestBucketsRequest{
			Status:    convertedStatus,
			Query:     query,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			reqLog.Warn("gRPC-метод ListManageableProjectJoinRequestBuckets вернул ошибку", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Список управляемых бакетов заявок успешно получен",
			"items_count", len(resp.GetItems()),
			"has_next_page", resp.GetNextPageToken() != "",
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
