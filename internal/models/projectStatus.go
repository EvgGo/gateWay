package models

import "time"

// ProjectStatus - тип статуса проекта
type ProjectStatus string

const (
	ProjectStatusNotStarted ProjectStatus = "not_started"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusDone       ProjectStatus = "done"
	ProjectStatusOnHold     ProjectStatus = "on_hold"
)

// JoinRequestStatus - тип статуса заявки на вступление в проект
type JoinRequestStatus string

const (
	JoinRequestStatusPending   JoinRequestStatus = "pending"
	JoinRequestStatusApproved  JoinRequestStatus = "approved"
	JoinRequestStatusRejected  JoinRequestStatus = "rejected"
	JoinRequestStatusCancelled JoinRequestStatus = "cancelled"
)

// ProjectRights - модель прав участника проекта
type ProjectRights struct {
	ManagerRights   bool `json:"manager_rights"`
	ManagerMember   bool `json:"manager_member"`
	ManagerProjects bool `json:"manager_projects"`
	ManagerTasks    bool `json:"manager_tasks"`
}

// Project - модель проекта
type Project struct {
	ID          string        `json:"id"`
	TeamID      string        `json:"team_id"`
	CreatorID   string        `json:"creator_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`
	IsOpen      bool          `json:"is_open"`
	StartedAt   time.Time     `json:"started_at"`
	FinishedAt  *time.Time    `json:"finished_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
