package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryDefaultSelector = DefaultSelector{
	Normal: "Query ID:.query_id,Status:.status,Executed By:.executed_by,Submitted At:.submitted_at",
	Wide:   "Query ID:.query_id,Status:.status,Executed By:.executed_by,Submitted At:.submitted_at,Expires At:.expires_at,SQL:.sql",
}

var _ = openapi.QueryInfo{
	ExecutedBy: nil,
	ExpiresAt:  nil,
	LastOffset: nil,
	Pagination: &openapi.Pagination{
		StartCursor: nil,
	},
	QueryErrors: nil,
	QueryId:     nil,
	Sql:         nil,
	Stats: &openapi.Stats{
		ElapsedTimeMs:          nil,
		ResultSetBytesSize:     nil,
		ResultSetDocumentCount: nil,
		ThrottledTimeMs:        nil,
	},
	Status:      nil,
	SubmittedAt: nil,
}
