package services

import (
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
)

// Permission constants
const (
	PermissionHome      = "HomePermission"
	PermissionCreateJob = "CreateJobPermission"
	PermissionViewJob   = "ViewJobPermission"
	PermissionIAM       = "IamPermission"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

// CheckRolePermission checks if a role has a specific permission.
func (s *PermissionService) CheckRolePermission(roleID string, permission string) (bool, error) {
	if roleID == "" {
		return false, nil
	}

	// Fetch the role
	var role models.Role
	if err := db.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return false, err
	}

	switch permission {
	case PermissionHome:
		return role.HomePermission, nil
	case PermissionCreateJob:
		return role.CreateJobPermission, nil
	case PermissionViewJob:
		return role.ViewJobPermission, nil
	case PermissionIAM:
		return role.IamPermission, nil
	default:
		return false, nil
	}
}

// GetUserPermissions returns all permissions for a user's role
func (s *PermissionService) GetUserPermissions(roleID string) (map[string]bool, error) {
	if roleID == "" {
		return make(map[string]bool), nil
	}

	var role models.Role
	if err := db.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}

	return map[string]bool{
		PermissionHome:      role.HomePermission,
		PermissionCreateJob: role.CreateJobPermission,
		PermissionViewJob:   role.ViewJobPermission,
		PermissionIAM:       role.IamPermission,
	}, nil
}
