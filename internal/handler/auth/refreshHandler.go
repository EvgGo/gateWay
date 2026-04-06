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

// RefreshHandler обрабатывает POST /auth/refresh
//
// Особенность: refresh_token может прийти
// - в JSON body,
// - или в HttpOnly cookie (refresh_token)
//
// - если decode сломался или body пустой - все равно попробуем cookie
// - но если токена нет нигде - 400
func RefreshHandler(log *slog.Logger, c authv1.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in DTO.RefreshJSON

		_ = utils.DecodeJSON(w, r, &in)

		refresh := strings.TrimSpace(in.RefreshToken)
		if refresh == "" {

			if ck, err := r.Cookie("refresh_token"); err == nil {
				refresh = strings.TrimSpace(ck.Value)
			}
		}
		if refresh == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "refresh_token is required (body or cookie)")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.Refresh(ctx, &authv1.RefreshRequest{
			RefreshToken: refresh,
			// session_id - опционально. Если клиент его не прислал, оставим пустым
			SessionId: strings.TrimSpace(in.SessionID),
		})
		if err != nil {
			log.Warn("Refresh failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
