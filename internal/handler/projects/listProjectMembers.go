package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func ListProjectMembersHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА ListProjectMembers запрос для проекта %v", projectID))

		pageSize := int32(10) // По умолчанию 10 участников
		pageToken := r.URL.Query().Get("page_token")

		resp, err := c.ListProjectMembers(ctx, &workspacev1.ListProjectMembersRequest{
			ProjectId: projectID,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn("ListProjectMembers failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
