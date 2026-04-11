package profile

type InviteUserToProjectBody struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type RevokeProjectInvitationBody struct {
	Reason string `json:"reason"`
}

type RejectProjectInvitationBody struct {
	Reason string `json:"reason"`
}
