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

func SetProjectOpenHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		var req workspacev1.SetProjectOpenRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("SetProjectOpen: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА SetProjectOpen запрос для проекта %v", projectID))

		req.ProjectId = projectID
		resp, err := c.SetProjectOpen(ctx, &req)
		if err != nil {
			log.Warn("SetProjectOpen failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
