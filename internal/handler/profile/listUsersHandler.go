package profile

import (
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// ListUsersHandler обрабатывает GET /users?query=&page_size=&page_token=
//
// Валидация page_size:
// - если не задан => default 20
// - ограничиваем диапазон 10..100
func ListUsersHandler(log *slog.Logger, c authv1.UserProfileClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("query"))
		pageToken := strings.TrimSpace(r.URL.Query().Get("page_token"))

		pageSize := int32(20)
		if v := strings.TrimSpace(r.URL.Query().Get("page_size")); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n <= 0 {
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "page_size must be a positive integer")
				return
			}
			if n < 10 {
				n = 10
			}
			if n > 100 {
				n = 100
			}
			pageSize = int32(n)
		}

		openSuggestionsOnly := false
		if v := strings.TrimSpace(r.URL.Query().Get("open_suggestions_only")); v != "" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "open_suggestions_only must be true or false")
				return
			}
			openSuggestionsOnly = parsed
		}

		skillIDs, err := parseSkillIDsFromQuery(r)
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListUsers(ctx, &authv1.ListUsersRequest{
			Query:               q,
			PageSize:            pageSize,
			PageToken:           pageToken,
			SkillIds:            skillIDs,
			OpenSuggestionsOnly: openSuggestionsOnly,
		})
		if err != nil {
			log.Warn("ListUsers failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
