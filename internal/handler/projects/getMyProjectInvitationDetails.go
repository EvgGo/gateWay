package projects

import (
	"encoding/json"
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

func GetMyProjectInvitationDetailsHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		invitationID := strings.TrimSpace(chi.URLParam(r, "invitation_id"))
		if invitationID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invitation_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.GetMyProjectInvitationDetails(ctx, &workspacev1.GetMyProjectInvitationDetailsRequest{
			InvitationId: invitationID,
		})
		if err != nil {
			log.Warn("GetMyProjectInvitationDetails failed", "err", err, "invitationID", invitationID)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("GetMyProjectInvitationDetails encode response failed", "err", err, "invitationID", invitationID)
		}
	}
}
