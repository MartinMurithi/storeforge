package dto

type CreateRoleRequestDTO struct {
	Name          string   `json:"full_name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

type CreateRoleResponseDTO struct {
	User    UserResponseDTO `json:"user"`
	Message string          `json:"message"`
}
