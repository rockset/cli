package format

import "github.com/rockset/rockset-go-client/openapi"

var CollectionDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Status", "status"),
		{
			ColumnName:     "Size",
			Path:           []PathElem{{FieldName: "stats"}, {FieldName: "total_size"}},
			FieldFormatter: SizeFormatter{},
		},
		{
			ColumnName:     "Last Queried",
			Path:           []PathElem{{FieldName: "stats"}, {FieldName: "last_queried_ms"}},
			FieldFormatter: TimeSinceFormatter{},
		},
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("Retention", "retention_secs"),
		NewFieldSelection("Status", "status"),
		NewFieldSelection("Insert Only", "insert_only"),
		NewFieldSelection("Read Only", "read_only"),
		NewFieldSelection("Created By", "created_by"),
		NewFieldSelection("Created At", "created_at"),
	},
}

// just to list available fields
var _ = openapi.Collection{
	Aliases:             nil,
	BulkStats:           nil,
	ClusteringKey:       nil,
	CreatedAt:           nil,
	CreatedBy:           nil,
	CreatedByApikeyName: nil,
	Description:         nil,
	FieldMappingQuery:   nil,
	FieldMappings:       nil,
	InsertOnly:          nil,
	Name:                nil,
	ReadOnly:            nil,
	RetentionSecs:       nil,
	Rrn:                 nil,
	Sources:             nil,
	Stats: &openapi.CollectionStats{
		BulkBytesInserted:    nil,
		BulkBytesOverwritten: nil,
		BytesInserted:        nil,
		BytesOverwritten:     nil,
		ColumnIndexSize:      nil,
		DocCount:             nil,
		FillProgress:         nil,
		InvertedIndexSize:    nil,
		LastQueriedMs:        nil,
		LastUpdatedMs:        nil,
		PurgedDocCount:       nil,
		PurgedDocSize:        nil,
		RangeIndexSize:       nil,
		RowIndexSize:         nil,
		TotalIndexSize:       nil,
		TotalSize:            nil,
	},
	Status:                 nil,
	StorageCompressionType: nil,
	Workspace:              nil,
}
