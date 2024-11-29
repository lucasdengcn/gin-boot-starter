package services

import (
	"gin-boot-starter/core"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/security"
	"gin-boot-starter/persistence/repository"

	"github.com/gin-gonic/gin"
)

type AclService struct {
	aclRepository *repository.AclRepository
}

// NewAclService with repository
func NewAclService(aclRepository *repository.AclRepository) *AclService {
	return &AclService{
		aclRepository: aclRepository,
	}
}

// SetForNewUser to user
func (acl *AclService) SetForNewUser(ctx *gin.Context, userId uint) bool {
	ok, err := acl.aclRepository.AssignRole(ctx, userId, security.RoleUser)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign role to new user failed. %v", userId)
		panic(core.NewServiceError(500, "Assign role to new user failed", "AclService.SetForNewUser"))
	}
	ok, err = acl.aclRepository.AssignPolicy(ctx, userId, "user", security.AclActionRead)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign policy to new user failed. %v", userId)
		panic(core.NewServiceError(500, "Assign policy to new user failed", "AclService.SetForNewUser"))
	}
	return ok
}

// AssignRole to user
func (acl *AclService) AssignRole(ctx *gin.Context, userId uint, role string) bool {
	ok, err := acl.aclRepository.AssignRole(ctx, userId, role)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign role to user failed. %v, %v", userId, role)
		panic(core.NewRepositoryError(500, "Assign role to user failed", "AclRepository.AssignRole"))
	}
	return ok
}

// RemoveRole from user
func (acl *AclService) RemoveRole(ctx *gin.Context, userId uint, role string) bool {
	ok, err := acl.aclRepository.RemoveRole(ctx, userId, role)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Remove role from user failed. %v, %v", userId, role)
		panic(core.NewRepositoryError(500, "Remove role from user failed", "AclRepository.RemoveRole"))
	}
	return ok
}

// AssignPolicy to user on object with action
func (acl *AclService) AssignPolicy(ctx *gin.Context, userId uint, object, act string) bool {
	ok, err := acl.aclRepository.AssignPolicy(ctx, userId, object, act)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign policy to user failed. %v, %v, %v", userId, object, act)
		panic(core.NewRepositoryError(500, "Assign policy to user failed", "AclRepository.AssignPolicy"))
	}
	return ok
}

// HasPolicy check for user on resource (obj, url) with action
func (acl *AclService) HasPolicy(ctx *gin.Context, userId uint, object, act string) bool {
	return acl.aclRepository.HasPolicy(ctx, userId, object, act)
}

// RemovePolicy from user on object with action
func (acl *AclService) RemovePolicy(ctx *gin.Context, userId uint, object, act string) bool {
	ok, err := acl.aclRepository.RemovePolicy(ctx, userId, object, act)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Remove policy from user failed. %v, %v, %v", userId, object, act)
		panic(core.NewRepositoryError(500, "Remove policy from user failed", "AclRepository.RemovePolicy"))
	}
	return ok
}

// LoadPolicy all from DB
func (acl *AclService) LoadPolicy() error {
	return acl.aclRepository.LoadPolicy()
}
