package server

import (
	"gin-boot-starter/apis/controllers"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/persistence/repository"
	"gin-boot-starter/services"

	"github.com/google/wire"
)

// ProviderSet
var dbSet = wire.NewSet(db.GetDBCon)

// ProviderSet
var aclServiceSet = wire.NewSet(repository.NewAclRepository, services.NewAclService)

// ProviderSet
var userServiceSet = wire.NewSet(repository.NewUserRepository, services.NewUserService)

// InitializeUserController injector
func InitializeUserController() *controllers.UserController {
	wire.Build(dbSet, userServiceSet, aclServiceSet, controllers.NewUserController)
	return &controllers.UserController{}
}

func InitializeAccountController() *controllers.AccountController {
	wire.Build(dbSet, userServiceSet, aclServiceSet, controllers.NewAccountController)
	return &controllers.AccountController{}
}

func InitializeUserService() *services.UserService {
	wire.Build(dbSet, userServiceSet)
	return &services.UserService{}
}

// InitializeAclService injector
func InitializeAclService() *services.AclService {
	wire.Build(dbSet, aclServiceSet)
	return &services.AclService{}
}
