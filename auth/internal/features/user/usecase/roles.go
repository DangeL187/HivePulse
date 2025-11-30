package usecase

import (
	"github.com/DangeL187/erax"

	"auth/internal/shared/roles"
)

type RolesUseCase struct {
	roleManager roles.RoleManager
}

func (r *RolesUseCase) GetRoles(userID uint) ([]string, error) {
	userRoles, err := r.roleManager.GetRoles(userID)
	if err != nil {
		return nil, erax.Wrap(err, "failed to get roles")
	}

	return userRoles, nil
}

func (r *RolesUseCase) GrantRole(userID uint, role string) (bool, error) {
	wasGranted, err := r.roleManager.GrantRole(userID, role)
	if err != nil {
		return false, erax.Wrap(err, "failed to grant role")
	}

	return wasGranted, nil
}

func (r *RolesUseCase) HasRole(userID uint, role string) (bool, error) {
	hasRole, err := r.roleManager.HasRole(userID, role)
	if err != nil {
		return false, erax.Wrap(err, "failed to check role")
	}

	return hasRole, nil
}

func (r *RolesUseCase) RevokeRole(userID uint, role string) error {
	err := r.roleManager.RevokeRole(userID, role)
	if err != nil {
		return erax.Wrap(err, "failed to revoke role")
	}

	return nil
}

func NewRolesUseCase(roleManager roles.RoleManager) *RolesUseCase {
	return &RolesUseCase{
		roleManager: roleManager,
	}
}
