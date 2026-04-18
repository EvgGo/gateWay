package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func ListTeamProjectsForAssignmentHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := chi.URLParam(r, "team_id")
		userID := chi.URLParam(r, "user_id")

		if teamID == "" || userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id and user_id are required")
			return
		}

		pageSize, err := parseInt32Query(r, "page_size")
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		req := &workspacev1.ListTeamProjectsForAssignmentRequest{
			TeamId:    teamID,
			UserId:    userID,
			Query:     strings.TrimSpace(r.URL.Query().Get("query")),
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА ListTeamProjectsForAssignment запрос для команды %v, участника %v", teamID, userID))

		resp, err := c.ListTeamProjectsForAssignment(ctx, req)
		if err != nil {
			log.Warn("ListTeamProjectsForAssignment failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
