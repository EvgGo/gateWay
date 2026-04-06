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

// ChangePasswordHandler обрабатывает POST /auth/change-password
// Требует Bearer token (middleware)
// оба поля должны быть непустые
func ChangePasswordHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var in DTO.ChangePasswordJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if strings.TrimSpace(in.CurrentPassword) == "" || strings.TrimSpace(in.NewPassword) == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "current_password and new_password are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		_, err := c.ChangePassword(ctx, &authv1.ChangePasswordRequest{
			CurrentPassword: in.CurrentPassword,
			NewPassword:     in.NewPassword,
		})
		if err != nil {
			log.Warn("ChangePassword failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
