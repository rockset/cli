package format

import "github.com/rockset/rockset-go-client/openapi"

var APIKeyFormatter = StructFormatter{
	[]Header{
		{
			FieldName: "Name",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Key ID",
			FieldName:   "Key",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created By",
			FieldName:   "CreatedBy",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Expiry Time",
			FieldName:   "ExpiryTime",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Last Access Time",
			FieldName:   "LastAccessTime",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Role",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "State",
			FieldFn:   getFieldByName,
		},
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
