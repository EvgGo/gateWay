package teams

import (
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// UpdateTeamHandler - обновление информации о команде
func UpdateTeamHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		var req workspacev1.UpdateTeamRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("UpdateTeam: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		req.TeamId = teamID
		resp, err := c.UpdateTeam(ctx, &req)
		if err != nil {
			log.Warn("UpdateTeam failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
