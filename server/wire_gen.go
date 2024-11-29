// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"gin-boot-starter/apis/controllers"
	"gin-boot-starter/infra/db"
	"gin-boot-starter/persistence/repository"
	"gin-boot-starter/services"
	"github.com/google/wire"
)

// Injectors from wire_injector.go:

// InitializeUserController injector
func InitializeUserController() *controllers.UserController {
	sqlxDB := db.GetDBCon()
	userRepository := repository.NewUserRepository(sqlxDB)
	userService := services.NewUserService(userRepository)
	aclRepository := repository.NewAclRepository(sqlxDB)
	aclService := services.NewAclService(aclRepository)
	userController := controllers.NewUserController(userService, aclService)
	return userController
}

func InitializeAccountController() *controllers.AccountController {
	sqlxDB := db.GetDBCon()
	userRepository := repository.NewUserRepository(sqlxDB)
	userService := services.NewUserService(userRepository)
	aclRepository := repository.NewAclRepository(sqlxDB)
	aclService := services.NewAclService(aclRepository)
	accountController := controllers.NewAccountController(userService, aclService)
	return accountController
}

func InitializeUserService() *services.UserService {
	sqlxDB := db.GetDBCon()
	userRepository := repository.NewUserRepository(sqlxDB)
	userService := services.NewUserService(userRepository)
	return userService
}

// InitializeAclService injector
func InitializeAclService() *services.AclService {
	sqlxDB := db.GetDBCon()
	aclRepository := repository.NewAclRepository(sqlxDB)
	aclService := services.NewAclService(aclRepository)
	return aclService
}

// wire_injector.go:

// ProviderSet
var dbSet = wire.NewSet(db.GetDBCon)

// ProviderSet
var aclServiceSet = wire.NewSet(repository.NewAclRepository, services.NewAclService)

// ProviderSet
var userServiceSet = wire.NewSet(repository.NewUserRepository, services.NewUserService)
