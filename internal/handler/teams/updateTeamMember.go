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

// UpdateTeamMemberHandler - обновление данных участника команды
func UpdateTeamMemberHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		userID := chi.URLParam(r, "user_id")

		if teamID == "" || userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id and user_id are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА UpdateTeamMember запрос для team %v, user %v", teamID, userID))

		var req workspacev1.UpdateTeamMemberRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("UpdateTeamMember: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		req.TeamId = teamID
		req.UserId = userID
		resp, err := c.UpdateTeamMember(ctx, &req)
		if err != nil {
			log.Warn("UpdateTeamMember failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
