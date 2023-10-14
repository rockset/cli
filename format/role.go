package format

import "github.com/rockset/rockset-go-client/openapi"

var RoleDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "role_name"),
		NewFieldSelection("Created By", "created_by"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Name", "role_name"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Description", "description"),
	},
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
