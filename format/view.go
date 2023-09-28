package format

import "github.com/rockset/rockset-go-client/openapi"

var ViewDefaultSelector = DefaultSelector{
	Normal: "Workspace:.workspace,Name:.name,Created By:.creator_email,State:.state,SQL:.query_sql",
	Wide:   "Workspace:.workspace,Name:.name,Description:.description,Created By:.creator_email,Created At:.created_at,State:.state,SQL:.query_sql",
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
