package profile

import (
	"encoding/json"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
	"strings"
)

func GetProfilesByIdsHandler(log *slog.Logger, c authv1.UserProfileClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		reqLog := log.With(
			"handler", "GetProfilesByIdsHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
		)

		reqLog.Info("Получен HTTP-запрос на получение публичных профилей по списку user_id")

		userIDs := collectUserIDs(r)
		if len(userIDs) == 0 {
			reqLog.Warn("Не передан ни один user_id/user_ids")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "at least one user_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.GetProfilesByIds(ctx, &authv1.GetProfilesByIdsRequest{
			UserIds: userIDs,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetProfilesByIds", "err", err)

			helpers.WriteGRPCError(w, r, err)
		}

		reqLog.Info("Публичные профили успешно получены", "count", len(resp.GetUsers()))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"users": resp.GetUsers(),
		})
	}
}

func collectUserIDs(r *http.Request) []string {
	q := r.URL.Query()

	rawValues := make([]string, 0)

	// поддержка:
	// ?user_ids=1&user_ids=2
	rawValues = append(rawValues, q["user_ids"]...)

	// поддержка:
	// ?user_ids=1,2,3
	if joined := strings.TrimSpace(q.Get("user_ids")); joined != "" {
		rawValues = append(rawValues, strings.Split(joined, ",")...)
	}

	// поддержка:
	// ?user_id=1&user_id=2
	rawValues = append(rawValues, q["user_id"]...)

	seen := make(map[string]struct{}, len(rawValues))
	result := make([]string, 0, len(rawValues))

	for _, v := range rawValues {
		id := strings.TrimSpace(v)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}
