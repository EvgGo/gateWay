package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func ListProjectMemberDetailsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "project_id")
		if strings.TrimSpace(projectID) == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		var pageSize int32
		rawPageSize := strings.TrimSpace(r.URL.Query().Get("page_size"))
		if rawPageSize != "" {
			parsedPageSize, err := strconv.Atoi(rawPageSize)
			if err != nil {
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "page_size must be an integer")
				return
			}
			pageSize = int32(parsedPageSize)
		}

		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(
			"ListProjectMemberDetails request",
			"projectID", projectID,
			"pageSize", pageSize,
			"pageToken", pageToken,
		)

		resp, err := c.ListProjectMemberDetails(ctx, &workspacev1.ListProjectMemberDetailsRequest{
			ProjectId: projectID,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn(
				"ListProjectMemberDetails failed",
				"projectID", projectID,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteJSON(w, http.StatusOK, resp)
	}
}
