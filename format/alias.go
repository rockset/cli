package format

import "github.com/rockset/rockset-go-client/openapi"

var AliasFormatter = StructFormatter{
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
			DisplayName: "Modified At",
			FieldName:   "ModifiedAt",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
		{
			FieldName: "Collections",
			FieldFn:   getArrayFieldByName,
		},
		{
			FieldName: "State",
			FieldFn:   getFieldByName,
		},
	},
}

// just to list available fields
var _ = openapi.Alias{
	Collections:         nil,
	CreatedAt:           nil,
	CreatedByApikeyName: nil,
	CreatorEmail:        nil,
	Description:         nil,
	ModifiedAt:          nil,
	Name:                nil,
	State:               nil,
	Workspace:           nil,
}
