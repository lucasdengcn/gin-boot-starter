package server

import (
	"gin001/apis/controllers"
	"gin001/infra/db"
	"gin001/persistence/repository"
	"gin001/services"

	"github.com/google/wire"
)

// reuseable
var dbSet = wire.NewSet(db.GetDBCon)

// reuseable
var userServiceSet = wire.NewSet(dbSet, repository.NewUserRepository, services.NewUserService)

// InitializeUserController injector
func InitializeUserController() *controllers.UserController {
	return controllers.NewUserController(InitializeUserService())
}

func InitializeAccountController() *controllers.AccountController {
	return controllers.NewAccountController(InitializeUserService())
}

func InitializeUserService() *services.UserService {
	wire.Build(userServiceSet)
	return &services.UserService{}
}
