package service

import (
	"gin-boot-starter/core/exception"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/security"
	"gin-boot-starter/persistence/repository"
	"sync"

	"github.com/gin-gonic/gin"
)

var onceAclService sync.Once
var instanceAclService *AclService

type AclService struct {
	aclRepository *repository.AclRepository
}

// NewAclService with repository
func NewAclService(aclRepository *repository.AclRepository) *AclService {
	onceAclService.Do(func() {
		instanceAclService = &AclService{
			aclRepository: aclRepository,
		}
	})
	return instanceAclService
}

// SetForNewUser to user
func (acl *AclService) SetForNewUser(ctx *gin.Context, userId uint) bool {
	ok, err := acl.aclRepository.AssignRole(ctx, userId, security.RoleUser)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign role to new user failed. %v", userId)
		panic(exception.NewServiceError(ctx, "ACL_ASSIGN_ROLE_500", "Assign role to new user failed"))
	}
	ok, err = acl.aclRepository.AssignPolicy(ctx, userId, "user", security.AclActionRead)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign policy to new user failed. %v", userId)
		panic(exception.NewServiceError(ctx, "ACL_ASSIGN_POLICY_500", "Assign policy to new user failed"))
	}
	return ok
}

// AssignRole to user
func (acl *AclService) AssignRole(ctx *gin.Context, userId uint, role string) bool {
	ok, err := acl.aclRepository.AssignRole(ctx, userId, role)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign role to user failed. %v, %v", userId, role)
		panic(exception.NewRepositoryError(ctx, "ACL_ASSIGN_ROLE_500", "Assign role to user failed"))
	}
	return ok
}

// RemoveRole from user
func (acl *AclService) RemoveRole(ctx *gin.Context, userId uint, role string) bool {
	ok, err := acl.aclRepository.RemoveRole(ctx, userId, role)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Remove role from user failed. %v, %v", userId, role)
		panic(exception.NewRepositoryError(ctx, "ACL_REMOVE_ROLE_500", "Remove role from user failed"))
	}
	return ok
}

// AssignPolicy to user on object with action
func (acl *AclService) AssignPolicy(ctx *gin.Context, userId uint, object, act string) bool {
	ok, err := acl.aclRepository.AssignPolicy(ctx, userId, object, act)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Assign policy to user failed. %v, %v, %v", userId, object, act)
		panic(exception.NewRepositoryError(ctx, "ACL_ASSIGN_POLICY_500", "Assign policy to user failed"))
	}
	return ok
}

// HasPolicy check for user on resource (obj, url) with action
func (acl *AclService) HasPolicy(ctx *gin.Context, userId uint, object, act string) bool {
	ok := acl.aclRepository.HasPolicy(ctx, userId, object, act)
	if !ok {
		panic(exception.NewACLError(ctx, "ACL_HAS_POLICY_403", "Permission denied"))
	}
	return ok
}

// RemovePolicy from user on object with action
func (acl *AclService) RemovePolicy(ctx *gin.Context, userId uint, object, act string) bool {
	ok, err := acl.aclRepository.RemovePolicy(ctx, userId, object, act)
	if err != nil {
		logging.Error(ctx).Err(err).Msgf("Remove policy from user failed. %v, %v, %v", userId, object, act)
		panic(exception.NewRepositoryError(ctx, "ACL_REMOVE_POLICY_500", "Remove policy from user failed"))
	}
	return ok
}

// LoadPolicy all from DB
func (acl *AclService) LoadPolicy() error {
	return acl.aclRepository.LoadPolicy()
}
