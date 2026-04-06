package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

func ListProjectsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		teamID := r.URL.Query().Get("team_id")
		creatorID := r.URL.Query().Get("creator_id")
		status := r.URL.Query().Get("status")
		onlyOpen := r.URL.Query().Get("only_open") == "true"
		query := r.URL.Query().Get("query")
		pageSize := int32(10) // по умолчанию 10 проектов
		pageToken := r.URL.Query().Get("page_token")

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		// Добавляем метаданные в контекст
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА ListProjects запрос для team %v, creator %v", teamID, creatorID))

		convertedStatus, err := helpers.ParseProjectStatusParam(status)
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid project status")
			return
		}

		resp, err := c.ListProjects(ctx, &workspacev1.ListProjectsRequest{
			TeamId:    teamID,
			CreatorId: creatorID,
			Status:    convertedStatus,
			OnlyOpen:  onlyOpen,
			Query:     query,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn("ListProjects failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
