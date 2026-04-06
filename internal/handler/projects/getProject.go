package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func GetProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			projectID = parts[len(parts)-1]
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА GetProject запрос для проекта %v", projectID))

		resp, err := c.GetProject(ctx, &workspacev1.GetProjectRequest{ProjectId: projectID})
		if err != nil {
			log.Warn("GetProject failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
