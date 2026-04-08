package helpers

import (
	"fmt"
	"gateWay/internal/models"
	"strconv"
	"strings"

	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
)

func ParseProjectStatusParam(s string) (workspacev1.ProjectStatus, error) {

	s = strings.TrimSpace(s)
	if s == "" {
		return workspacev1.ProjectStatus_PROJECT_STATUS_UNSPECIFIED, nil
	}

	if n, err := strconv.Atoi(s); err == nil {
		v := workspacev1.ProjectStatus(n)
		switch v {
		case workspacev1.ProjectStatus_PROJECT_STATUS_UNSPECIFIED,
			workspacev1.ProjectStatus_PROJECT_STATUS_NOT_STARTED,
			workspacev1.ProjectStatus_PROJECT_STATUS_IN_PROGRESS,
			workspacev1.ProjectStatus_PROJECT_STATUS_DONE,
			workspacev1.ProjectStatus_PROJECT_STATUS_ON_HOLD:
			return v, nil
		default:
			return workspacev1.ProjectStatus_PROJECT_STATUS_UNSPECIFIED, fmt.Errorf("invalid project status: %q", s)
		}
	}

	// можно string values: not_started / in_progress / done / on_hold
	switch strings.ToLower(s) {
	case string(models.ProjectStatusNotStarted):
		return workspacev1.ProjectStatus_PROJECT_STATUS_NOT_STARTED, nil
	case string(models.ProjectStatusInProgress):
		return workspacev1.ProjectStatus_PROJECT_STATUS_IN_PROGRESS, nil
	case string(models.ProjectStatusDone):
		return workspacev1.ProjectStatus_PROJECT_STATUS_DONE, nil
	case string(models.ProjectStatusOnHold):
		return workspacev1.ProjectStatus_PROJECT_STATUS_ON_HOLD, nil
	case "unspecified":
		return workspacev1.ProjectStatus_PROJECT_STATUS_UNSPECIFIED, nil
	default:
		return workspacev1.ProjectStatus_PROJECT_STATUS_UNSPECIFIED, fmt.Errorf("invalid project status: %q", s)
	}
}

func ParseJoinRequestStatusParam(s string) (workspacev1.JoinRequestStatus, error) {

	s = strings.TrimSpace(s)
	if s == "" {
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, nil
	}

	// можно numeric enum values in query: status=1
	if n, err := strconv.Atoi(s); err == nil {
		v := workspacev1.JoinRequestStatus(n)
		switch v {
		case workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED,
			workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_PENDING,
			workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_APPROVED,
			workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_REJECTED,
			workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_CANCELLED:
			return v, nil
		default:
			return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, fmt.Errorf("invalid join request status: %q", s)
		}
	}

	// можно string values: pending / approved / rejected / cancelled
	switch strings.ToLower(s) {
	case string(models.JoinRequestStatusPending):
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_PENDING, nil
	case string(models.JoinRequestStatusApproved):
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_APPROVED, nil
	case string(models.JoinRequestStatusRejected):
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_REJECTED, nil
	case string(models.JoinRequestStatusCancelled):
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_CANCELLED, nil
	case "unspecified":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, nil
	default:
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, fmt.Errorf("invalid join request status: %q", s)
	}
}

func ParseManageableBucketStatus(raw string) (workspacev1.JoinRequestStatus, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "pending":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_PENDING, nil
	case "approved":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_APPROVED, nil
	case "rejected":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_REJECTED, nil
	case "cancelled":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_CANCELLED, nil
	case "unspecified", "all":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, nil
	default:
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, fmt.Errorf("unknown join request status: %q", raw)
	}
}

func ParseJoinRequestStatusOrAll(raw string) (workspacev1.JoinRequestStatus, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "all", "unspecified":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, nil
	case "pending":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_PENDING, nil
	case "approved":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_APPROVED, nil
	case "rejected":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_REJECTED, nil
	case "cancelled":
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_CANCELLED, nil
	default:
		return workspacev1.JoinRequestStatus_JOIN_REQUEST_STATUS_UNSPECIFIED, fmt.Errorf("unknown join request status: %q", raw)
	}
}
