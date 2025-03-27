package dto

import (
	"github.com/google/uuid"
)

type RoleRequest struct {
	Name          string      `json:"name" binding:"required"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

type RoleResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type AddRolePermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

type RemoveRolePermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

type AssignRoleToUserRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" binding:"required"`
}

type RemoveRoleFromUserRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" binding:"required"`
}
