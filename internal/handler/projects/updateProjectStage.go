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

func UpdateProjectStageHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		stageID := strings.TrimSpace(chi.URLParam(r, "stage_id"))
		if stageID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "stage_id is required")
			return
		}

		var body projects.UpdateProjectStageHTTPReq
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			log.Warn("UpdateProjectStage: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		req := &workspacev1.UpdateProjectStageRequest{
			StageId: stageID,
		}

		if body.Title != nil {
			v := strings.TrimSpace(*body.Title)
			req.Title = &v
		}

		if body.Description != nil {
			v := strings.TrimSpace(*body.Description)
			req.Description = &v
		}

		if body.WeightPercent != nil {
			v := *body.WeightPercent
			req.WeightPercent = &v
		}

		if body.Status != nil {
			status, err := projects.ParseProjectStageStatusParam(*body.Status)
			if err != nil {
				log.Warn("UpdateProjectStage: invalid status", "status", *body.Status, "err", err)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid project stage status")
				return
			}

			if status == workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_UNSPECIFIED {
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "Project stage status must not be unspecified")
				return
			}

			req.Status = &status
		}

		if body.ProgressPercent != nil {
			v := *body.ProgressPercent
			req.ProgressPercent = &v
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА UpdateProjectStage запрос", "stageID", stageID)

		resp, err := c.UpdateProjectStage(ctx, req)
		if err != nil {
			log.Warn("UpdateProjectStage failed", "stageID", stageID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
