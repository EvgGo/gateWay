package auth

import (
	DTO "gateWay/internal/DTO/Auth"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
	"strings"
)

// LogoutHandler обрабатывает POST /auth/logout.
//
// Требуем Bearer token через middleware requireBearerAuth.
func LogoutHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var in DTO.LogoutJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if strings.TrimSpace(in.SessionID) == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "session_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		_, err := c.Logout(ctx, &authv1.LogoutRequest{
			SessionId: strings.TrimSpace(in.SessionID),
		})
		if err != nil {
			log.Warn("Logout failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		// 204 No Content - стандартный ответ - для операции успешно выполнена, тела нет.
		w.WriteHeader(http.StatusNoContent)
	}
}
