package format

import "github.com/rockset/rockset-go-client/openapi"

var AliasDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Collections", "collections"),
		NewFieldSelection("State", "state"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("Modified At", "modified_at"),
		NewFieldSelection("Collections", "collections"),
		NewFieldSelection("State", "state"),
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
