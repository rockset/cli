package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryInfoFormatter = StructFormatter{
	[]Header{
		{
			DisplayName: "Query ID",
			FieldName:   "QueryId",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Status",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Executed By",
			FieldName:   "ExecutedBy",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Submitted At",
			FieldName:   "SubmittedAt",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Expires At",
			FieldName:   "ExpiresAt",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
		{
			DisplayName: "SQL",
			FieldName:   "Sql",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
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
