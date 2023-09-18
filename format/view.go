package format

import "github.com/rockset/rockset-go-client/openapi"

var ViewFormatter = StructFormatter{
	[]Header{
		{
			FieldName: "Workspace",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "Name",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "Description",
			FieldFn:   getFieldByName,
			Wide:      true,
		},
		{
			DisplayName: "Created By",
			FieldName:   "CreatorEmail",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created At",
			FieldName:   "CreatedAt",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
		{
			FieldName: "State",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "SQL",
			FieldName:   "QuerySql",
			FieldFn:     getFieldByName,
		},
	},
}

// just to list available fields
var _ = openapi.View{
	CreatedAt:           nil,
	CreatedByApikeyName: nil,
	CreatorEmail:        nil,
	Description:         nil,
	Entities:            nil,
	ModifiedAt:          nil,
	Name:                nil,
	OwnerEmail:          nil,
	Path:                nil,
	QuerySql:            nil,
	State:               nil,
	Workspace:           nil,
}
