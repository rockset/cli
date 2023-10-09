package format

import "github.com/rockset/rockset-go-client/openapi"

var RoleDefaultSelector = DefaultSelector{
	Normal: "Name:.role_name,Created By:.created_by",
	Wide:   "Name:.role_name,Created By:.created_by,Description:.description",
}

// just to list available fields
var _ = openapi.Role{
	CreatedAt:   nil,
	CreatedBy:   nil,
	Description: nil,
	OwnerEmail:  nil,
	Privileges:  nil,
	RoleName:    nil,
}
