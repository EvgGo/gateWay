package projects

import (
	"fmt"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	"net/http"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func RejectProjectJoinRequestHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		requestID := chi.URLParam(r, "request_id")
		if requestID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "request_id is required")
			return
		}

		var req workspacev1.RejectProjectJoinRequestRequest
		if err := utils.DecodeJSON(w, r, &req); err != nil {
			log.Warn("RejectProjectJoinRequest: Invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "Invalid input")
			return
		}

		req.RequestId = requestID

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		// Добавляем метаданные в контекст
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug(fmt.Sprintf("НА RejectProjectJoinRequest запрос для заявки %v", requestID))

		resp, err := c.RejectProjectJoinRequest(ctx, &req)
		if err != nil {
			log.Warn("RejectProjectJoinRequest failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
