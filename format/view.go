package format

import "github.com/rockset/rockset-go-client/openapi"

var ViewDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Created By", "creator_email"),
		NewFieldSelection("State", "state"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("Created By", "creator_email"),
		NewFieldSelection("Created At", "created_at"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("SQL", "query_sql"),
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
