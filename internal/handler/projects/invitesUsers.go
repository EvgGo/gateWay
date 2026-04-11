package projects

import (
	"errors"
	"gateWay/internal/DTO/profile"
	"gateWay/internal/helpers"
	"gateWay/pkg/utils"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func InviteUserToProjectHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		var body profile.InviteUserToProjectBody
		if err := utils.DecodeJSON(w, r, &body); err != nil {
			log.Warn("InviteUserToProject: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid input")
			return
		}
		if body.UserID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "user_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.InviteUserToProject(ctx, &workspacev1.InviteUserToProjectRequest{
			ProjectId: projectID,
			UserId:    body.UserID,
			Message:   body.Message,
		})
		if err != nil {
			log.Warn("InviteUserToProject failed", "project_id", projectID, "user_id", body.UserID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListProjectInvitationsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		statusValue, err := parseProjectInvitationStatus(r.URL.Query().Get("status"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		pageSize, err := parseOptionalInt32(r.URL.Query().Get("page_size"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListProjectInvitations(ctx, &workspacev1.ListProjectInvitationsRequest{
			ProjectId: projectID,
			Status:    statusValue,
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		})
		if err != nil {
			log.Warn("ListProjectInvitations failed", "project_id", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListProjectInvitationDetailsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		statusValue, err := parseProjectInvitationStatus(r.URL.Query().Get("status"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		pageSize, err := parseOptionalInt32(r.URL.Query().Get("page_size"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListProjectInvitationDetails(ctx, &workspacev1.ListProjectInvitationDetailsRequest{
			ProjectId: projectID,
			Status:    statusValue,
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		})
		if err != nil {
			log.Warn("ListProjectInvitationDetails failed", "project_id", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func RevokeProjectInvitationHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		invitationID := chi.URLParam(r, "invitation_id")
		if invitationID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invitation_id is required")
			return
		}

		var body profile.RevokeProjectInvitationBody
		if err := helpers.DecodeOptionalJSON(r, &body); err != nil {
			log.Warn("RevokeProjectInvitation: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.RevokeProjectInvitation(ctx, &workspacev1.RevokeProjectInvitationRequest{
			InvitationId: invitationID,
			Reason:       body.Reason,
		})
		if err != nil {
			log.Warn("RevokeProjectInvitation failed", "invitation_id", invitationID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func GetMyProjectInvitationHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectID := chi.URLParam(r, "project_id")
		if projectID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "project_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.GetMyProjectInvitation(ctx, &workspacev1.GetMyProjectInvitationRequest{
			ProjectId: projectID,
		})
		if err != nil {
			log.Warn("GetMyProjectInvitation failed", "project_id", projectID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListMyProjectInvitationsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		statusValue, err := parseProjectInvitationStatus(r.URL.Query().Get("status"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		pageSize, err := parseOptionalInt32(r.URL.Query().Get("page_size"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListMyProjectInvitations(ctx, &workspacev1.ListMyProjectInvitationsRequest{
			Status:    statusValue,
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		})
		if err != nil {
			log.Warn("ListMyProjectInvitations failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func AcceptProjectInvitationHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		invitationID := chi.URLParam(r, "invitation_id")
		if invitationID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invitation_id is required")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.AcceptProjectInvitation(ctx, &workspacev1.AcceptProjectInvitationRequest{
			InvitationId: invitationID,
		})
		if err != nil {
			log.Warn("AcceptProjectInvitation failed", "invitation_id", invitationID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func RejectProjectInvitationHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		invitationID := chi.URLParam(r, "invitation_id")
		if invitationID == "" {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invitation_id is required")
			return
		}

		var body profile.RejectProjectInvitationBody
		if err := helpers.DecodeOptionalJSON(r, &body); err != nil {
			log.Warn("RejectProjectInvitation: invalid input", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid input")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.RejectProjectInvitation(ctx, &workspacev1.RejectProjectInvitationRequest{
			InvitationId: invitationID,
			Reason:       body.Reason,
		})
		if err != nil {
			log.Warn("RejectProjectInvitation failed", "invitation_id", invitationID, "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func ListMyInvitableProjectsHandler(log *slog.Logger, c workspacev1.ProjectsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pageSize, err := parseOptionalInt32(r.URL.Query().Get("page_size"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid page_size")
			return
		}

		onlyOpen, err := parseOptionalBool(r.URL.Query().Get("only_open"))
		if err != nil {
			helpers.WriteAPIError(w, r, http.StatusBadRequest, "invalid only_open")
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()
		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		resp, err := c.ListMyInvitableProjects(ctx, &workspacev1.ListMyInvitableProjectsRequest{
			Query:     r.URL.Query().Get("query"),
			OnlyOpen:  onlyOpen,
			PageSize:  pageSize,
			PageToken: r.URL.Query().Get("page_token"),
		})
		if err != nil {
			log.Warn("ListMyInvitableProjects failed", "err", err)
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func parseProjectInvitationStatus(raw string) (workspacev1.ProjectInvitationStatus, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "all", "unspecified":
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_UNSPECIFIED, nil
	case "pending":
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_PENDING, nil
	case "accepted":
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_ACCEPTED, nil
	case "rejected":
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_REJECTED, nil
	case "revoked":
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_REVOKED, nil
	default:
		return workspacev1.ProjectInvitationStatus_PROJECT_INVITATION_STATUS_UNSPECIFIED, errors.New("invalid status")
	}
}

func parseOptionalInt32(raw string) (int32, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func parseOptionalBool(raw string) (bool, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, nil
	}
	return strconv.ParseBool(raw)
}
