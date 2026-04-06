package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func ListProjectJoinRequestsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		status := r.URL.Query().Get("status")
		pageSize := int32(10)
		pageToken := r.URL.Query().Get("page_token")

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		convertedStatus, err := helpers.ParseJoinRequestStatusParam(status)
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid projectJoin status")
			return
		}

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА ListProjectJoinRequests запрос для проекта %v", projectID))

		resp, err := c.ListProjectJoinRequests(ctx, &workspacev1.ListProjectJoinRequestsRequest{
			ProjectId: projectID,
			Status:    convertedStatus,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn("ListProjectJoinRequests failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
