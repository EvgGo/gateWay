package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

func UpdateTeamMemberDutiesHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		userID := chi.URLParam(r, "user_id")

		if teamID == "" || userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id and user_id are required")
			return
		}

		var req workspacev1.UpdateTeamMemberDutiesRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("UpdateTeamMemberDuties: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		req.TeamId = teamID
		req.UserId = userID

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("На UpdateTeamMemberDuties запрос для команды %v, участника %v", teamID, userID))

		resp, err := c.UpdateTeamMemberDuties(ctx, &req)
		if err != nil {
			log.Warn("UpdateTeamMemberDuties failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
