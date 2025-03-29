package repositories

import (
	"errors"
	"fmt"
	"go-json/internal/models"
	"sync"
)

type RoleRepository interface {
	FindByRoleName(role string) (*models.Role, error)
	FindByRoleID(roleID string) (*models.Role, error)
	FindRoleByUserID(userID string) (*[]models.UserRole, error)
}

type roleRepository struct {
	roles     []models.Role
	userRoles []models.UserRole
	mu        sync.RWMutex
}

func NewRoleRepository(roles []models.Role, userRole []models.UserRole) RoleRepository {
	return &roleRepository{
		roles:     roles,
		userRoles: userRole,
		mu:        sync.RWMutex{},
	}
}

func (r *roleRepository) FindByRoleName(role string) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, r := range r.roles {
		if r.Name == role {
			roleCopy := r
			return &roleCopy, nil
		}
	}
	return nil, errors.New("role not found")
}

func (r *roleRepository) FindByRoleID(roleID string) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, r := range r.roles {
		if r.ID == roleID {
			roleCopy := r
			return &roleCopy, nil
		}
	}
	return nil, errors.New("role not found")
}

func (r *roleRepository) FindRoleByUserID(userID string) (*[]models.UserRole, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userRoles []models.UserRole
	for _, ur := range r.userRoles {
		if ur.UserID == userID {
			userRoles = append(userRoles, ur)
		}
	}

	fmt.Println("UserRoles found for userID:", userID, "=>", userRoles)

	if len(userRoles) == 0 {
		return nil, errors.New("user has no roles")
	}

	return &userRoles, nil
}
