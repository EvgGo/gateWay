package projects

import (
	"fmt"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
)

type CreateProjectStageHTTPReq struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Position        int32  `json:"position"`
	WeightPercent   int32  `json:"weight_percent"`
	Status          string `json:"status"`
	ProgressPercent int32  `json:"progress_percent"`
}

type UpdateProjectStageHTTPReq struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	WeightPercent   *int32  `json:"weight_percent"`
	Status          *string `json:"status"`
	ProgressPercent *int32  `json:"progress_percent"`
}

type ReorderProjectStagesHTTPReq struct {
	Items []projectStageOrderHTTPItem `json:"items"`
}

type projectStageOrderHTTPItem struct {
	StageID  string `json:"stage_id"`
	Position int32  `json:"position"`
}

type EvaluateProjectStageHTTPReq struct {
	Score        int32  `json:"score"`
	ScoreComment string `json:"score_comment"`
}

// ParseProjectStageStatusParam переводит REST-значение статуса этапа в proto enum
//
// Поддерживает короткие frontend-значения:
// planned, in_progress, done, cancelled.
//
// Также поддерживает proto-значения:
// PROJECT_STAGE_STATUS_PLANNED,
// PROJECT_STAGE_STATUS_IN_PROGRESS,
// PROJECT_STAGE_STATUS_DONE,
// PROJECT_STAGE_STATUS_CANCELLED
func ParseProjectStageStatusParam(raw string) (workspacev1.ProjectStageStatus, error) {
	value := strings.TrimSpace(strings.ToLower(raw))

	switch value {
	case "", "unspecified", "project_stage_status_unspecified":
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_UNSPECIFIED, nil

	case "planned", "project_stage_status_planned":
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_PLANNED, nil

	case "in_progress", "project_stage_status_in_progress":
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_IN_PROGRESS, nil

	case "done", "project_stage_status_done":
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_DONE, nil

	case "cancelled", "canceled", "project_stage_status_cancelled":
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_CANCELLED, nil

	default:
		return workspacev1.ProjectStageStatus_PROJECT_STAGE_STATUS_UNSPECIFIED,
			fmt.Errorf("invalid project stage status: %s", raw)
	}
}
