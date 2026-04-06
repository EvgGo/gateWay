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

// ResetPasswordHandler обрабатывает POST /auth/reset-password
//
// Валидация:
//   - три поля обязательны
//   - email нормализуем
//   - reset_token/new_password не тримим агрессивно (тримим только проверку на пустоту),
//     чтобы случайно не"сломать токен, если в нем возможны пробелы
func ResetPasswordHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		reqLog := log.With(
			"op", "ResetPasswordHandler",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"request_id", r.Header.Get("X-Request-ID"),
		)

		var in DTO.ResetPasswordJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			reqLog.Warn("invalid JSON body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		email := strings.ToLower(strings.TrimSpace(in.Email))
		resetToken := strings.TrimSpace(in.ResetToken)
		newPassword := strings.TrimSpace(in.NewPassword)

		if email == "" || resetToken == "" || newPassword == "" {
			reqLog.Warn("validation failed",
				"email_empty", email == "",
				"reset_token_empty", resetToken == "",
				"new_password_empty", newPassword == "",
			)

			helpers.WriteAPIError(w, r, http.StatusBadRequest, "email, reset_token and new_password are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		_, err := c.ResetPassword(ctx, &authv1.ResetPasswordRequest{
			Email:       email,
			ResetToken:  in.ResetToken,
			NewPassword: in.NewPassword,
		})
		if err != nil {
			reqLog.Warn("ResetPassword failed",
				"email", email,
				"err", err,
			)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
