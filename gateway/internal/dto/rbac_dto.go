package dto

type RoleResponseDTO struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Slug        string                  `json:"slug"`
	Description string                  `json:"description"`
	IsSystem    bool                    `json:"is_system"`
	Permissions []PermissionResponseDTO `json:"permissions"`
	CreatedAt   string                  `json:"created_at"`
}

type PermissionResponseDTO struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CreateRoleRequestDTO struct {
	Name          string   `json:"name" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

type CreateRoleResponseDTO struct {
	Role    RoleResponseDTO `json:"role"`
	Message string          `json:"message"`
}
