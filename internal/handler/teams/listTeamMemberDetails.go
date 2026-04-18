package teams

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func ListTeamMemberDetailsHandler(log *slog.Logger, c workspacev1.TeamsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		teamID := chi.URLParam(r, "team_id")
		if teamID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "team_id is required")
			return
		}

		pageSize, err := parseInt32Query(r, "page_size")
		if err != nil {
			http.Error(w, "invalid page_size", http.StatusBadRequest)
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListTeamMemberDetails(ctx, &workspacev1.ListTeamMemberDetailsRequest{
			TeamId:    teamID,
			Query:     strings.TrimSpace(r.URL.Query().Get("query")),
			SkillIds:  parseRepeatedQuery(r, "skill_id", "skill_ids"),
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		})
		if err != nil {
			log.Warn("ListTeamMemberDetails failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func parseInt32Query(r *http.Request, key string) (int32, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}

func parseRepeatedQuery(r *http.Request, keys ...string) []string {
	values := make([]string, 0)

	for _, key := range keys {
		for _, raw := range r.URL.Query()[key] {
			for _, part := range strings.Split(raw, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					values = append(values, part)
				}
			}
		}
	}

	return values
}
