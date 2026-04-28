package projects

import (
	"gateWay/internal/DTO/projects"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"log/slog"
	"net/http"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
)

func EvaluateProjectStageHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stageID := strings.TrimSpace(chi.URLParam(r, "stage_id"))
		if stageID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "stage_id is required")
			return
		}

		var body projects.EvaluateProjectStageHTTPReq
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			log.Warn("EvaluateProjectStage: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА EvaluateProjectStage запрос",
			"stageID", stageID,
			"score", body.Score,
		)

		resp, err := c.EvaluateProjectStage(ctx, &workspacev1.EvaluateProjectStageRequest{
			StageId:      stageID,
			Score:        body.Score,
			ScoreComment: strings.TrimSpace(body.ScoreComment),
		})
		if err != nil {
			log.Warn("EvaluateProjectStage failed", "stageID", stageID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
