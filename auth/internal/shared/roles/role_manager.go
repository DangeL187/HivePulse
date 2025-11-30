package roles

type RoleManager interface {
	GrantRole(userID uint, role string) (bool, error)
	GetRoles(userID uint) ([]string, error)
	HasRole(userID uint, role string) (bool, error)
	RevokeRole(userID uint, role string) error
}
