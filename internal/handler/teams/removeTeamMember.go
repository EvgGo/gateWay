package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// RemoveTeamMemberHandler - удаление участника из команды
func RemoveTeamMemberHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

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

		log.Debug(fmt.Sprintf("НА RemoveTeamMember запрос для team %v, user %v", teamID, userID))

		_, err := c.RemoveTeamMember(ctx, &workspacev1.RemoveTeamMemberRequest{
			TeamId: teamID,
			UserId: userID,
		})
		if err != nil {
			log.Warn("RemoveTeamMember failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteSuccess(w, http.StatusNoContent)
	}
}
