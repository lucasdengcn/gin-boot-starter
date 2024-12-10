package server

import (
	"gin-boot-starter/api/controller"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/persistence/repository"
	"gin-boot-starter/service"

	"github.com/google/wire"
)

// ProviderSet
var dbSet = wire.NewSet(db.GetDBCon)

// ProviderSet
var aclServiceSet = wire.NewSet(repository.NewAclRepository, service.NewAclService)

// ProviderSet
var userServiceSet = wire.NewSet(repository.NewUserRepository, service.NewUserService)

// InitializeUserController injector
func InitializeUserController() *controller.UserController {
	wire.Build(dbSet, userServiceSet, aclServiceSet, controller.NewUserController)
	return &controller.UserController{}
}

func InitializeAccountController() *controller.AccountController {
	wire.Build(dbSet, userServiceSet, aclServiceSet, controller.NewAccountController)
	return &controller.AccountController{}
}

func InitializeUserService() *service.UserService {
	wire.Build(dbSet, userServiceSet)
	return &service.UserService{}
}

// InitializeAclService injector
func InitializeAclService() *service.AclService {
	wire.Build(dbSet, aclServiceSet)
	return &service.AclService{}
}
