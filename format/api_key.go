package format

import "github.com/rockset/rockset-go-client/openapi"

var ApiKeyDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Key ID", "key"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Expiry Time", "expiry_time"),
		NewFieldSelection("Last Access Time", "last_access_time"),
		NewFieldSelection("Role", "role"),
		NewFieldSelection("State", "state"),
	},
}

// just to list available fields
var _ = openapi.ApiKey{
	CreatedAt:           nil,
	CreatedBy:           nil,
	CreatedByApikeyName: nil,
	ExpiryTime:          nil,
	Key:                 "",
	LastAccessTime:      nil,
	Name:                "",
	Role:                nil,
	State:               nil,
}
