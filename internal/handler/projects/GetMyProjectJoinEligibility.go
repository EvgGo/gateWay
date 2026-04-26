package projects

import (
	"gateWay/internal/helpers"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetMyProjectJoinEligibilityHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "project_id")

		reqLog := log.With(
			"handler", "GetMyProjectJoinEligibilityHandler",
			"http_method", r.Method,
			"path", r.URL.Path,
			"raw_query", r.URL.RawQuery,
			"project_id", projectID,
		)

		if projectID == "" {
			reqLog.Warn("Пустой project_id")
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid project_id")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на проверку eligibility вступления в проект")

		resp, err := c.GetMyProjectJoinEligibility(ctx, &workspacev1.GetMyProjectJoinEligibilityRequest{
			ProjectId: projectID,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC GetMyProjectJoinEligibility", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Eligibility успешно получен",
			"can_request_join", resp.GetCanRequestJoin(),
			"is_project_open", resp.GetIsProjectOpen(),
			"already_member", resp.GetAlreadyMember(),
			"has_pending_join_request", resp.GetHasPendingJoinRequest(),
			"has_pending_invitation", resp.GetHasPendingInvitation(),
			"matched_requirements_count", resp.GetMatchedRequirementsCount(),
			"total_requirements_count", resp.GetTotalRequirementsCount(),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
