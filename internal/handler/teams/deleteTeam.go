package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// DeleteTeamHandler - удаление команды по ID
func DeleteTeamHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА DeleteTeam запрос %v", teamID))

		_, err := c.DeleteTeam(ctx, &workspacev1.DeleteTeamRequest{TeamId: teamID})
		if err != nil {
			log.Warn("DeleteTeam failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteSuccess(w, http.StatusNoContent)
	}
}
