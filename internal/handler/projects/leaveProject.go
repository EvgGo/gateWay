package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func LeaveProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА LeaveProject запрос", "projectID", projectID)

		_, err := c.LeaveProject(ctx, &workspacev1.LeaveProjectRequest{
			ProjectId: projectID,
		})
		if err != nil {
			log.Warn("LeaveProject failed", "projectID", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
