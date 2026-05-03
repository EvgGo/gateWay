package teams

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func LeaveTeamHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		teamID := strings.TrimSpace(chi.URLParam(r, "team_id"))
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("НА LeaveTeam запрос", "teamID", teamID)

		_, err := c.LeaveTeam(ctx, &workspacev1.LeaveTeamRequest{
			TeamId: teamID,
		})
		if err != nil {
			log.Warn("LeaveTeam failed", "teamID", teamID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
