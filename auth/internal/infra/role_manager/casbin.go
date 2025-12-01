package role_manager

import (
	"gorm.io/gorm"
	"strconv"

	"github.com/DangeL187/erax"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"

	"auth/internal/shared/config"
)

type RoleManager struct {
	enforcer *casbin.Enforcer
}

func (rm *RoleManager) Enforce(rvals ...interface{}) (bool, error) {
	return rm.enforcer.Enforce(rvals...)
}

func (rm *RoleManager) GetRoles(userID uint) ([]string, error) {
	userIDString := strconv.FormatUint(uint64(userID), 10)
	roles, err := rm.enforcer.GetRolesForUser(userIDString)
	if err != nil {
		return nil, erax.Wrap(err, "failed to get roles for user")
	}

	return roles, nil
}

func (rm *RoleManager) GrantRole(userID uint, role string) (bool, error) {
	userIDString := strconv.FormatUint(uint64(userID), 10)
	wasAdded, err := rm.enforcer.AddGroupingPolicy(userIDString, role)
	if err != nil {
		return false, erax.Wrap(err, "failed to add grouping policy")
	}

	return wasAdded, nil
}

func (rm *RoleManager) HasRole(userID uint, role string) (bool, error) {
	userIDString := strconv.FormatUint(uint64(userID), 10)
	hasRole, err := rm.enforcer.HasGroupingPolicy(userIDString, role)
	if err != nil {
		return false, erax.Wrap(err, "failed to check grouping policy")
	}

	return hasRole, nil
}

func (rm *RoleManager) RevokeRole(userID uint, role string) error {
	userIDString := strconv.FormatUint(uint64(userID), 10)
	_, err := rm.enforcer.RemoveGroupingPolicy(userIDString, role)
	if err != nil {
		return erax.Wrap(err, "failed to remove grouping policy")
	}

	return nil
}

func (rm *RoleManager) initEnforcer(cfg *config.Config, db *gorm.DB) error {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return erax.Wrap(err, "failed to create gorm adapter")
	}

	rm.enforcer, err = casbin.NewEnforcer(cfg.CasbinModelConfigPath, adapter)
	if err != nil {
		return erax.Wrap(err, "failed to create enforcer")
	}

	err = rm.enforcer.LoadPolicy()
	if err != nil {
		return erax.Wrap(err, "failed to load policy")
	}

	return nil
}

func (rm *RoleManager) loadPolicy() error {
	policies := [][]any{
		{"admin", "user", "grant_role", "allow"},
		{"admin", "user", "revoke_role", "allow"},
		{"operator", "device", "register", "allow"},
		{"operator", "device", "watch", "allow"},
	}

	for _, policy := range policies {
		_, err := rm.enforcer.AddPolicy(policy...)
		if err != nil {
			return erax.Wrap(err, "failed to add policy")
		}
	}

	return nil
}

func NewRoleManager(cfg *config.Config, db *gorm.DB) (*RoleManager, error) {
	rm := &RoleManager{}

	err := rm.initEnforcer(cfg, db)
	if err != nil {
		return nil, erax.Wrap(err, "failed to init enforcer")
	}

	err = rm.loadPolicy()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load policy")
	}

	return rm, nil
}
