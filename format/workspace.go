package format

import "github.com/rockset/rockset-go-client/openapi"

var WorkspaceFormatter = StructFormatter{
	[]Header{
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
			DisplayName: "Collections",
			FieldName:   "CollectionCount",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created By",
			FieldName:   "CreatedBy",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created At",
			FieldName:   "CreatedAt",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
	},
}

// just to list available fields
var _ = openapi.Workspace{
	CollectionCount: nil,
	CreatedAt:       nil,
	CreatedBy:       nil,
	Description:     nil,
	Name:            nil,
}
