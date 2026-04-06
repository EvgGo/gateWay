package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func RemoveProjectMemberHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		userID := chi.URLParam(r, "user_id")

		if projectID == "" || userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id and user_id are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА RemoveProjectMember запрос для проекта %v, участник %v", projectID, userID))

		_, err := c.RemoveProjectMember(ctx, &workspacev1.RemoveProjectMemberRequest{
			ProjectId: projectID,
			UserId:    userID,
		})
		if err != nil {
			log.Warn("RemoveProjectMember failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteSuccess(w, http.StatusNoContent)
	}
}
