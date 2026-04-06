package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// ListTeamMembersHandler - получение списка участников команды
func ListTeamMembersHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		// Добавляем метаданные в контекст
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА ListTeamMembers запрос для команды %v", teamID))

		pageSize := int32(10) // по умолчанию 10 участников
		pageToken := r.URL.Query().Get("page_token")

		resp, err := c.ListTeamMembers(ctx, &workspacev1.ListTeamMembersRequest{
			TeamId:    teamID,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn("ListTeamMembers failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
