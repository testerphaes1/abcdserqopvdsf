package usecase_models

type Member struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type Role string

const (
	RoleViewer Role = "viewer"
	RoleEditor Role = "editor"
	RoleOwner  Role = "owner"
)
