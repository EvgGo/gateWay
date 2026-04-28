package projects

import (
	"gateWay/internal/DTO/projects"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
)

func CreateProjectStageHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		var body projects.CreateProjectStageHTTPReq
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			log.Warn("CreateProjectStage: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		status, err := projects.ParseProjectStageStatusParam(body.Status)
		if err != nil {
			log.Warn("CreateProjectStage: invalid status", "status", body.Status, "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid project stage status")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА CreateProjectStage запрос",
			"projectID", projectID,
			"title", body.Title,
			"position", body.Position,
			"weight_percent", body.WeightPercent,
			"status", body.Status,
		)

		resp, err := c.CreateProjectStage(ctx, &workspacev1.CreateProjectStageRequest{
			ProjectId:       projectID,
			Title:           strings.TrimSpace(body.Title),
			Description:     strings.TrimSpace(body.Description),
			Position:        body.Position,
			WeightPercent:   body.WeightPercent,
			Status:          status,
			ProgressPercent: body.ProgressPercent,
		})
		if err != nil {
			log.Warn("CreateProjectStage failed", "projectID", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusCreated, resp)
	}
}
