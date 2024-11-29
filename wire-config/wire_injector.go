package server

import (
	"gin-boot-starter/apis/controllers"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/persistence/repository"
	"gin-boot-starter/services"

	"github.com/google/wire"
)

// reuseable
var dbSet = wire.NewSet(db.GetDBCon)

var aclServiceSet = wire.NewSet(dbSet, repository.NewAclRepository, services.NewAclService)

// reuseable
var userServiceSet = wire.NewSet(dbSet, repository.NewUserRepository, services.NewUserService)

// InitializeUserController injector
func InitializeUserController() *controllers.UserController {
	return controllers.NewUserController(InitializeUserService())
}

func InitializeAccountController() *controllers.AccountController {
	return controllers.NewAccountController(InitializeUserService(), InitializeAclService())
}

func InitializeUserService() *services.UserService {
	wire.Build(userServiceSet)
	return &services.UserService{}
}

func InitializeAclService() *services.AclService {
	wire.Build(aclServiceSet)
	return &services.AclService{}
}
