package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type RoleService struct {
	roleRepository *repository.RoleRepository
}

func NewRoleService(roleRepository *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

func (s *RoleService) Create(req *dto.RoleRequest) (*dto.RoleResponse, error) {
	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepository.Create(role); err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepository.AddPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return s.GetByID(role.ID)
}

func (s *RoleService) GetAll() ([]dto.RoleResponse, error) {
	roles, err := s.roleRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.toResponse(&role))
	}

	return responses, nil
}

func (s *RoleService) GetByID(id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := s.roleRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(role), nil
}

func (s *RoleService) Update(id uuid.UUID, req *dto.RoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := s.roleRepository.Update(role); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *RoleService) Delete(id uuid.UUID) error {
	return s.roleRepository.Delete(id)
}

func (s *RoleService) AddPermissions(id uuid.UUID, permissionIDs []uuid.UUID) error {
	return s.roleRepository.AddPermissions(id, permissionIDs)
}

func (s *RoleService) RemovePermissions(id uuid.UUID, permissionIDs []uuid.UUID) error {
	return s.roleRepository.RemovePermissions(id, permissionIDs)
}

func (s *RoleService) AssignToUser(userID uuid.UUID, roleIDs []uuid.UUID) error {
	return s.roleRepository.AssignToUser(userID, roleIDs)
}

func (s *RoleService) RemoveFromUser(userID uuid.UUID, roleIDs []uuid.UUID) error {
	return s.roleRepository.RemoveFromUser(userID, roleIDs)
}

func (s *RoleService) toResponse(role *entity.Role) *dto.RoleResponse {
	var permissions []dto.PermissionResponse
	for _, permission := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			CreatedAt:   permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
