package projects

import (
	"gateWay/internal/helpers"
	"net/http"
	"strconv"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func RemoveProjectMemberHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		userID := chi.URLParam(r, "user_id")

		if strings.TrimSpace(projectID) == "" || strings.TrimSpace(userID) == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id and user_id are required")
			return
		}

		removeFromTeam := false
		rawRemoveFromTeam := strings.TrimSpace(r.URL.Query().Get("remove_from_team"))
		if rawRemoveFromTeam != "" {
			parsedValue, err := strconv.ParseBool(rawRemoveFromTeam)
			if err != nil {
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "remove_from_team must be a boolean")
				return
			}
			removeFromTeam = parsedValue
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(
			"RemoveProjectMember request",
			"projectID", projectID,
			"userID", userID,
			"removeFromTeam", removeFromTeam,
		)

		_, err := c.RemoveProjectMember(ctx, &workspacev1.RemoveProjectMemberRequest{
			ProjectId:      projectID,
			UserId:         userID,
			RemoveFromTeam: removeFromTeam,
		})
		if err != nil {
			log.Warn(
				"RemoveProjectMember failed",
				"projectID", projectID,
				"userID", userID,
				"removeFromTeam", removeFromTeam,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
