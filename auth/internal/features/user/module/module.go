package module

import (
	"auth/internal/features/user/domain"
	"auth/internal/features/user/usecase"
	"auth/internal/shared/roles"
	"auth/internal/shared/token"
)

type Module struct {
	Auth  *usecase.AuthUseCase
	Roles *usecase.RolesUseCase
	Seed  *usecase.SeedUseCase
	User  *usecase.UserUseCase
}

func NewModule(repo domain.Repository, tokenGenerator token.Generator, roleManager roles.RoleManager) *Module {
	authUseCase := usecase.NewAuthUseCase(repo, tokenGenerator)
	rolesUseCase := usecase.NewRolesUseCase(roleManager)
	userUseCase := usecase.NewUserUseCase(repo)
	seedUseCase := usecase.NewSeedUseCase(authUseCase, rolesUseCase, userUseCase)

	return &Module{
		Auth:  authUseCase,
		Roles: rolesUseCase,
		Seed:  seedUseCase,
		User:  userUseCase,
	}
}
