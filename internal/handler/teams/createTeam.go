package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

// CreateTeamHandler - создание новой команды
func CreateTeamHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		var req workspacev1.CreateTeamRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("CreateTeam: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		log.Debug(fmt.Sprintf("НА CreateTeam запрос %v", req))

		resp, err := c.CreateTeam(ctx, &req)
		if err != nil {
			log.Warn("CreateTeam failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
