package format

import "github.com/rockset/rockset-go-client/openapi"

var AliasDefaultSelector = DefaultSelector{
	Normal: "Workspace:.workspace,Name:.name,Collections:.collections,State:.state",
	Wide:   "Workspace:.workspace,Name:.name,Description:.description,Modified At:.modified_at,Collections:.collections,State:.state",
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
