package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// GetTeamHandler - получение команды по ID
func GetTeamHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА GetTeam запрос %v", teamID))

		resp, err := c.GetTeam(ctx, &workspacev1.GetTeamRequest{TeamId: teamID})
		if err != nil {
			log.Warn("GetTeam failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		utils.PrintReadable(resp)
		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
