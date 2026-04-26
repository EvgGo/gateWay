package projects

import (
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type projectAssessmentRequirementInputBody struct {
	AssessmentID int64 `json:"assessment_id"`
	MinLevel     int32 `json:"min_level"`
}

type setProjectAssessmentRequirementsBody struct {
	Requirements []projectAssessmentRequirementInputBody `json:"requirements"`
}

func SetProjectAssessmentRequirementsHandler(
	log *slog.Logger,
	c workspacev1.ProjectsClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "project_id")

		reqLog := log.With(
			"handler", "SetProjectAssessmentRequirementsHandler",
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

		var body setProjectAssessmentRequirementsBody
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			reqLog.Warn("Некорректный body", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid body")
			return
		}

		requirements := make([]*workspacev1.ProjectAssessmentRequirementInput, 0, len(body.Requirements))
		for i, item := range body.Requirements {
			if item.AssessmentID <= 0 {
				reqLog.Warn("Некорректный assessment_id", "index", i, "assessment_id", item.AssessmentID)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid assessment_id")
				return
			}
			if item.MinLevel < 1 || item.MinLevel > 5 {
				reqLog.Warn("Некорректный min_level", "index", i, "min_level", item.MinLevel)
				helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid min_level")
				return
			}

			requirements = append(requirements, &workspacev1.ProjectAssessmentRequirementInput{
				AssessmentId: item.AssessmentID,
				MinLevel:     item.MinLevel,
			})
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		reqLog.Info("Получен HTTP-запрос на полную замену требований проекта",
			"requirements_count", len(requirements),
		)

		resp, err := c.SetProjectAssessmentRequirements(ctx, &workspacev1.SetProjectAssessmentRequirementsRequest{
			ProjectId:    projectID,
			Requirements: requirements,
		})
		if err != nil {
			reqLog.Error("Ошибка gRPC SetProjectAssessmentRequirements", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		reqLog.Info("Требования проекта успешно обновлены",
			"project_id", resp.GetId(),
			"requirements_count", len(resp.GetAssessmentRequirements()),
		)

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}
