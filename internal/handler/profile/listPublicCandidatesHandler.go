package profile

import (
	"gateWay/internal/helpers"
	ssov1 "github.com/EvgGo/proto/proto/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

func ListPublicCandidatesHandler(log *slog.Logger, c ssov1.UserProfileClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req, err := buildListPublicCandidatesRequest(r)
		if err != nil {
			log.Warn("ListPublicCandidates: invalid query params", "err", err)
			helpers.WriteAPIError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := helpers.CtxWithOutgoingMeta(r)
		defer cancel()

		ctx = helpers.AppendCommonGRPCMetadata(ctx, r)

		log.Debug("ListPublicCandidates request",
			"query", req.GetQuery(),
			"page_size", req.GetPageSize(),
			"page_token", req.GetPageToken(),
			"skill_ids_count", len(req.GetSkillIds()),
			"skill_match_mode", req.GetSkillMatchMode().String(),
			"open_suggestions_only", req.GetOpenSuggestionsOnly(),
			"sort_by", req.GetSortBy().String(),
			"sort_order", req.GetSortOrder().String(),
		)

		resp, err := c.ListPublicCandidates(ctx, req)
		if err != nil {
			helpers.WriteGRPCError(w, r, err)
			return
		}

		helpers.WriteProtoJSON(w, http.StatusOK, resp)
	}
}

func buildListPublicCandidatesRequest(r *http.Request) (*ssov1.ListPublicCandidatesRequest, error) {
	q := r.URL.Query()

	pageSize, err := parseOptionalInt32(q.Get("page_size"))
	if err != nil {
		return nil, err
	}

	openSuggestionsOnly, err := parseOptionalBool(q.Get("open_suggestions_only"))
	if err != nil {
		return nil, err
	}

	createdFrom, err := parseOptionalTimestamp(q.Get("created_from"))
	if err != nil {
		return nil, err
	}

	createdTo, err := parseOptionalTimestamp(q.Get("created_to"))
	if err != nil {
		return nil, err
	}

	if createdFrom != nil && createdTo != nil && createdFrom.AsTime().After(createdTo.AsTime()) {
		return nil, status.Error(codes.InvalidArgument, "created_from must be <= created_to")
	}

	skillMatchMode, err := parseUserSkillMatchMode(q.Get("skill_match_mode"))
	if err != nil {
		return nil, err
	}

	sortBy, err := parsePublicUserSortBy(q.Get("sort_by"))
	if err != nil {
		return nil, err
	}

	sortOrder, err := parseSortOrder(q.Get("sort_order"))
	if err != nil {
		return nil, err
	}

	req := &ssov1.ListPublicCandidatesRequest{
		Query:               strings.TrimSpace(q.Get("query")),
		PageSize:            pageSize,
		PageToken:           strings.TrimSpace(q.Get("page_token")),
		SkillIds:            parseSkillIDs(q),
		SkillMatchMode:      skillMatchMode,
		CreatedFrom:         createdFrom,
		CreatedTo:           createdTo,
		OpenSuggestionsOnly: openSuggestionsOnly,
		SortBy:              sortBy,
		SortOrder:           sortOrder,
	}

	return req, nil
}
