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
	PermissionIDs []string `json:"permission_ids" binding:"required,gt=0"` // Must have at least one ID
}

type CreateRoleResponseDTO struct {
	Role    RoleResponseDTO `json:"role"`
	Message string          `json:"message"`
}


type GetRoleByIDRequestDTO struct {
	Id          string   `json:"id" binding:"required"`
}

type GetRoleByIDResponseDTO struct {
	Role          RoleResponseDTO   `json:"role"`
}

type UpdateRoleRequestDTO struct {
	Name          string   `json:"name" binding:"omitempty"`
	Description   string   `json:"description" binding:"omitempty"`
	PermissionIDs []string `json:"permission_ids" binding:"omitempty"` // Must have at least one ID
}

type UpdateRoleResponseDTO struct {
	Role    RoleResponseDTO `json:"role"`
	Message string          `json:"message"`
}