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

func ReorderProjectStagesHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := strings.TrimSpace(chi.URLParam(r, "project_id"))
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		var body projects.ReorderProjectStagesHTTPReq
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			log.Warn("ReorderProjectStages: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		items := make([]*workspacev1.ProjectStageOrderItem, 0, len(body.Items))
		for _, item := range body.Items {
			items = append(items, &workspacev1.ProjectStageOrderItem{
				StageId:  strings.TrimSpace(item.StageID),
				Position: item.Position,
			})
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА ReorderProjectStages запрос",
			"projectID", projectID,
			"items_count", len(items),
		)

		resp, err := c.ReorderProjectStages(ctx, &workspacev1.ReorderProjectStagesRequest{
			ProjectId: projectID,
			Items:     items,
		})
		if err != nil {
			log.Warn("ReorderProjectStages failed", "projectID", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
