package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func UpdateProjectMemberRightsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		userID := chi.URLParam(r, "user_id")

		if projectID == "" || userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id and user_id are required")
			return
		}

		var req workspacev1.UpdateProjectMemberRightsRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("UpdateProjectMemberRights: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА UpdateProjectMemberRights запрос для проекта %v, участника %v", projectID, userID))

		resp, err := c.UpdateProjectMemberRights(ctx, &req)
		if err != nil {
			log.Warn("UpdateProjectMemberRights failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
