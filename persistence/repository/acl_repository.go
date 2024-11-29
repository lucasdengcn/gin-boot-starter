package repository

import (
	"gin-boot-starter/config"
	"gin-boot-starter/core"
	"gin-boot-starter/core/logging"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/casbin/casbin/v2"
	sqlxadapter "github.com/memwey/casbin-sqlx-adapter"
)

var onceAclRepository sync.Once
var instanceAclRepository *AclRepository

type AclRepository struct {
	TransactionRepo
	enforcer *casbin.Enforcer
}

func NewAclRepository(dbCon *sqlx.DB) *AclRepository {
	// init adapter with existing connection
	onceAclRepository.Do(func() {
		opts := &sqlxadapter.AdapterOptions{
			DB:        dbCon,
			TableName: "casbin_rule",
		}
		a := sqlxadapter.NewAdapterFromOptions(opts)
		//
		path := config.GetConfig().Application.CfgPath + "/rbac_model.conf"
		e, err := casbin.NewEnforcer(path, a)
		if err != nil {
			logging.Fatal(nil).Err(err).Msg("Init casbin enforcer Error.")
		}

		instanceAclRepository = &AclRepository{
			TransactionRepo: NewTransactionRepo(dbCon),
			enforcer:        e,
		}
	})
	return instanceAclRepository
}

// AssignRole to user
func (acl *AclRepository) AssignRole(c *gin.Context, userId uint, role string) (bool, error) {
	return acl.enforcer.AddGroupingPolicy(core.StringFromUint(userId), role)
}

// RemoveRole from user
func (acl *AclRepository) RemoveRole(c *gin.Context, userId uint, role string) (bool, error) {
	return acl.enforcer.RemoveGroupingPolicy(core.StringFromUint(userId), role)
}

// AssignPolicy to user on object with action
func (acl *AclRepository) AssignPolicy(c *gin.Context, userId uint, object, act string) (bool, error) {
	return acl.enforcer.AddPolicy(core.StringFromUint(userId), object, act)
}

// HasPolicy check for user on resource (obj, url) with action
func (acl *AclRepository) HasPolicy(ctx *gin.Context, userId uint, object, act string) bool {
	return acl.enforcer.HasPolicy(core.StringFromUint(userId), object, act)
}

// RemovePolicy from user on object with action
func (acl *AclRepository) RemovePolicy(c *gin.Context, userId uint, object, act string) (bool, error) {
	return acl.enforcer.RemovePolicy(core.StringFromUint(userId), object, act)
}

// LoadPolicy all from DB
func (acl *AclRepository) LoadPolicy() error {
	return acl.enforcer.LoadPolicy()
}
