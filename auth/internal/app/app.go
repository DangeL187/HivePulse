package app

import (
	"github.com/DangeL187/erax"

	deviceInfra "auth/internal/features/device/infra"
	deviceModule "auth/internal/features/device/module"
	userInfra "auth/internal/features/user/infra"
	userModule "auth/internal/features/user/module"
	"auth/internal/infra/database"
	"auth/internal/infra/jwt"
	"auth/internal/infra/role_manager"
	"auth/internal/shared/config"
)

type App struct {
	Config      *config.Config
	JWTManager  *jwt.Manager
	RoleManager *role_manager.RoleManager

	DeviceModule *deviceModule.Module
	UserModule   *userModule.Module
}

func NewApp() (*App, error) {
	app := &App{}

	var err error
	app.Config, err = config.NewConfig()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load config")
	}

	db, err := database.NewPostgres(app.Config)
	if err != nil {
		return nil, erax.Wrap(err, "failed to connect to DB")
	}

	app.JWTManager, err = jwt.NewJWTManager()
	if err != nil {
		return nil, erax.Wrap(err, "failed to create JWT Manager")
	}

	app.RoleManager, err = role_manager.NewRoleManager(app.Config, db)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create RoleManager")
	}

	deviceRepo := deviceInfra.NewDeviceRepo(db)
	app.DeviceModule = deviceModule.NewModule(deviceRepo, app.JWTManager)

	userRepo := userInfra.NewUserRepo(db)
	app.UserModule = userModule.NewModule(userRepo, app.JWTManager, app.RoleManager)

	err = app.UserModule.Seed.SeedAdmin()
	if err != nil {
		return nil, erax.Wrap(err, "failed to init admin")
	}

	return app, nil
}
