package DTO

type ChangePasswordJSON struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
