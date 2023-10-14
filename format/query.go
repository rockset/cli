package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Query ID", "query_id"),
		NewFieldSelection("Status", "status"),
		NewFieldSelection("Executed By", "executed_by"),
		NewFieldSelection("Submitted At", "submitted_at"),
		NewFieldSelection("User", "rockset_user"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Query ID", "query_id"),
		NewFieldSelection("Status", "status"),
		NewFieldSelection("Executed By", "executed_by"),
		NewFieldSelection("Submitted At", "submitted_at"),
		NewFieldSelection("Expires At", "expires_at"),
		NewFieldSelection("SQL", "sql"),
	},
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
