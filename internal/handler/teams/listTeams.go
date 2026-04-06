package teams

import (
	"fmt"
	"gateWay/internal/helpers"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
)

// ListTeamsHandler - получение списка команд
func ListTeamsHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		query := r.URL.Query().Get("query")
		onlyMy := r.URL.Query().Get("only_my") == "true"
		pageSize := int32(10)
		pageToken := r.URL.Query().Get("page_token")

		log.Debug(fmt.Sprintf("НА ListTeams запрос %v", query))

		resp, err := c.ListTeams(ctx, &workspacev1.ListTeamsRequest{
			Query:     query,
			OnlyMy:    onlyMy,
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Warn("ListTeams failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
