package server

import (
	"gin001/apis/controllers"
	"gin001/infra/db"
	"gin001/persistence/repository"
	"gin001/services"

	"github.com/google/wire"
)

var dbSet = wire.NewSet(db.GetDBCon)

// UserControllerSet of dependencies
var UserControllerSet = wire.NewSet(dbSet, repository.NewUserRepository, services.NewUserService, controllers.NewUserController)

// InitializeUserController injector
func InitializeUserController() *controllers.UserController {
	wire.Build(UserControllerSet)
	return &controllers.UserController{}
}
