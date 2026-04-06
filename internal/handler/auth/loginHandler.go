package auth

import (
	"fmt"
	DTO "gateWay/internal/DTO/Auth"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"log/slog"
	"net/http"
	"strings"
)

// LoginHandler обрабатывает POST /auth/login.
//
//	decode -> минимальная валидация -> gRPC -> JSON
//
// Важные моменты:
// - email нормализуем, password trim НЕ делаем (пароль может содержать пробелы)
func LoginHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var in DTO.LoginJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		email := strings.ToLower(strings.TrimSpace(in.Email))
		if email == "" || strings.TrimSpace(in.Password) == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "email and password are required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.Login(ctx, &authv1.LoginRequest{
			Email:    email,
			Password: in.Password,
		})
		if err != nil {
			log.Warn("Login failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		log.Debug(fmt.Sprintf("На login получили %v", resp))

		tok := resp.GetTokens()

		helpers.WriteJSON(w, http.StatusOK, map[string]any{
			"tokens": map[string]any{
				"access_token":      tok.GetAccessToken(),
				"access_expires_at": tok.GetAccessExpiresAt(),
			},
		})
	}
}
