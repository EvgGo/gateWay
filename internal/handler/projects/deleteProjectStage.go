package projects

import (
	"gateWay/internal/helpers"
	"log/slog"
	"net/http"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
)

func DeleteProjectStageHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stageID := strings.TrimSpace(chi.URLParam(r, "stage_id"))
		if stageID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "stage_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА DeleteProjectStage запрос", "stageID", stageID)

		resp, err := c.DeleteProjectStage(ctx, &workspacev1.DeleteProjectStageRequest{
			StageId: stageID,
		})
		if err != nil {
			log.Warn("DeleteProjectStage failed", "stageID", stageID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
