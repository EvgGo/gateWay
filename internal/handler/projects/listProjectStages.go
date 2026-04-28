package projects

import (
	"gateWay/internal/helpers"
	"log/slog"
	"net/http"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
)

func ListProjectStagesHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА ListProjectStages запрос", "projectID", projectID)

		resp, err := c.ListProjectStages(ctx, &workspacev1.ListProjectStagesRequest{
			ProjectId: projectID,
		})
		if err != nil {
			log.Warn("ListProjectStages failed", "projectID", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
