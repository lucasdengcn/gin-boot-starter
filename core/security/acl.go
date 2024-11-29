package security

// casbin integration with sqlx

var (
	AclActionUpdate = "update"
	AclActionCreate = "create"
	AclActionDelete = "delete"
	AclActionRead   = "read"
)

var (
	RoleUser      = "User"
	RoleAdmin     = "Administrator"
	RoleOperation = "Operator"
)
