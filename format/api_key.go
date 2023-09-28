package format

import "github.com/rockset/rockset-go-client/openapi"

var ApiKeyDefaultSelector = DefaultSelector{
	Normal: "Name:.name,Key ID:.key,Created By:.created_by,Expiry Time:.expiry_time,Last Access Time:.last_access_time,Role:.role,State:.state",
	Wide:   "",
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
