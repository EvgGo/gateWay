package profile

import (
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

// GetProfileHandler обрабатывает GET /users/{user_id}
//
// Это публичная ручка, Bearer token не требуется
// Возвращает PublicUser (без скрытых/приватных полей)
func GetProfileHandler(log *slog.Logger, c authv1.UserProfileClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := strings.TrimSpace(chi.URLParam(r, "user_id"))
		if userID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "user_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.GetProfile(ctx, &authv1.GetProfileRequest{UserId: userID})
		if err != nil {
			log.Warn("GetProfile failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
