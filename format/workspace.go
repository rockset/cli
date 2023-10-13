package format

import "github.com/rockset/rockset-go-client/openapi"

var WorkspaceDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Collections", "collection_count"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Created At", "created_at"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("Collections", "collection_count"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Created At", "created_at"),
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
