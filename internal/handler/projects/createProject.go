package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

func CreateProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req workspacev1.CreateProjectRequest

		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("CreateProject: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		log.Debug("Gateway CreateProject: decoded request",
			"name", req.GetName(),
			"teamName", req.GetTeamName(),
			"teamModeString", req.GetTeamMode().String(),
			"teamModeNumber", int32(req.GetTeamMode()),
			"skillIds", req.GetSkillIds(),
		)

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА CreateProject запрос %v", req))

		resp, err := c.CreateProject(ctx, &req)
		if err != nil {
			log.Warn("CreateProject failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
