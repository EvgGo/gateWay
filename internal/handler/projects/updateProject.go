package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
	"net/http"
)

func UpdateProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := helpers.ExtractProjectID(r)
		if projectID == "" {
			log.Debug(fmt.Sprintf("НА UpdateProject запрос с путем: %s", r.URL.Path))
			return
		}

		var req workspacev1.UpdateProjectRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("UpdateProject: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА UpdateProject запрос для проекта %v", projectID))

		req.ProjectId = projectID
		resp, err := c.UpdateProject(ctx, &req)
		if err != nil {
			log.Warn("UpdateProject failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
