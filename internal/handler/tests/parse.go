package tests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	assessmentv1 "github.com/EvgGo/proto/proto/gen/go/tests"
)

func parsePageSize(raw string, defaultValue int32) (int32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid page_size")
	}

	return int32(value), nil
}

func parseInt64Query(raw string, defaultValue int64) (int64, error) {

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < 0 {
		return 0, errors.New("invalid int64 query value")
	}

	return value, nil
}

func parsePathInt64(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, errors.New("empty path param")
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid path param")
	}

	return value, nil
}

func parseOptionalBool(raw string, defaultValue bool) (bool, error) {

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, errors.New("invalid bool")
	}

	return value, nil
}

func parseAssessmentStatus(raw string) (assessmentv1.AssessmentStatus, error) {

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "unspecified", "all":
		return assessmentv1.AssessmentStatus_ASSESSMENT_STATUS_UNSPECIFIED, nil
	case "active":
		return assessmentv1.AssessmentStatus_ASSESSMENT_STATUS_ACTIVE, nil
	case "archived":
		return assessmentv1.AssessmentStatus_ASSESSMENT_STATUS_ARCHIVED, nil
	default:
		return assessmentv1.AssessmentStatus_ASSESSMENT_STATUS_UNSPECIFIED, errors.New("invalid assessment status")
	}
}

func parseAttemptStatus(raw string) (assessmentv1.AttemptStatus, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "unspecified", "all":
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_UNSPECIFIED, nil
	case "in_progress":
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_IN_PROGRESS, nil
	case "completed":
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_COMPLETED, nil
	case "abandoned":
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_ABANDONED, nil
	case "expired":
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_EXPIRED, nil
	default:
		return assessmentv1.AttemptStatus_ATTEMPT_STATUS_UNSPECIFIED, errors.New("invalid attempt status")
	}
}

func parseFinishAction(raw string) (assessmentv1.FinishAction, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "complete":
		return assessmentv1.FinishAction_FINISH_ACTION_COMPLETE, nil
	case "abandon":
		return assessmentv1.FinishAction_FINISH_ACTION_ABANDON, nil
	default:
		return assessmentv1.FinishAction_FINISH_ACTION_UNSPECIFIED, errors.New("invalid finish action")
	}
}

func decodeJSONBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain a single json object")
	}

	return nil
}
