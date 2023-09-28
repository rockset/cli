package format

import "github.com/rockset/rockset-go-client/openapi"

var WorkspaceDefaultSelector = DefaultSelector{
	Normal: "Name:.name,Collections:.collection_count,Created By:.created_by",
	Wide:   "Name:.name,Description:.description,Collections:.collection_count,Created By:.created_by,Created At:.created_at",
}

// just to list available fields
var _ = openapi.Workspace{
	CollectionCount: nil,
	CreatedAt:       nil,
	CreatedBy:       nil,
	Description:     nil,
	Name:            nil,
}
