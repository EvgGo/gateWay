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

// ForgotPasswordHandler обрабатывает POST /auth/forgot-password
// Сервис возвращает OK даже если email не найден,
// чтобы не давать атакующему возможность перебором узнавать существующие emails
// Это правило Auth сервиса; gateway просто проксирует
func ForgotPasswordHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var in DTO.ForgotPasswordJSON
		if err := utils.DecodeJSON(w, r, &in); err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid JSON body")
			return
		}

		email := strings.ToLower(strings.TrimSpace(in.Email))
		if email == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "email is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		_, err := c.ForgotPassword(ctx, &authv1.ForgotPasswordRequest{Email: email})
		if err != nil {
			log.Warn("ForgotPassword failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
