package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func DeleteProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА DeleteProject запрос для проекта %v", projectID))

		_, err := c.DeleteProject(ctx, &workspacev1.DeleteProjectRequest{ProjectId: projectID})
		if err != nil {
			log.Warn("DeleteProject failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteSuccess(w, http.StatusNoContent)
	}
}
