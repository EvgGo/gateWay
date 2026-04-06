package profile

import (
	"fmt"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"net/http"
)

// GetMeHandler обрабатывает GET /users/me
// Требует Bearer token (middleware)
// Возвращает полный User, включая приватные поля (email)
func GetMeHandler(log *slog.Logger, c authv1.UserProfileClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА GetMeHandler запрос %v", r))

		resp, err := c.GetMe(ctx, &emptypb.Empty{})
		if err != nil {
			log.Warn("GetMe failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
