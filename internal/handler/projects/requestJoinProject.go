package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

func RequestJoinProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		projectID := helpers.ExtractProjectID(r)
		if projectID == "" {
			log.Debug(fmt.Sprintf("НА RequestJoinProjectHandler запрос с путем: %s", r.URL.Path))
			return
		}

		var in requestJoinProjectJson
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		message := strings.TrimSpace(in.Message)
		if message == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "message is required")
			return
		}

		req := &workspacev1.RequestJoinProjectRequest{
			ProjectId: projectID,
			Message:   message,
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА RequestJoinProject запрос для проекта %v", projectID))

		resp, err := c.RequestJoinProject(ctx, req)
		if err != nil {
			log.Warn("RequestJoinProject failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

type requestJoinProjectJson struct {
	Message string `json:"message"`
}
