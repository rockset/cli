package format

import "github.com/rockset/rockset-go-client/openapi"

var MountDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Collection Path", "collection_path"),
		NewFieldSelection("State", "state"),
		{
			ColumnName:     "Last Queried",
			Path:           []PathElem{{FieldName: "stats"}, {FieldName: "last_queried_ms"}},
			FieldFormatter: TimeSinceFormatter{},
		},
	},
	Wide: []FieldSelection{
		NewFieldSelection("Collection Path", "collection_path"),
		NewFieldSelection("ID", "id"),
		NewFieldSelection("State", "state"),
		{
			ColumnName:     "Last Queried",
			Path:           []PathElem{{FieldName: "stats"}, {FieldName: "last_queried_ms"}},
			FieldFormatter: TimeSinceFormatter{},
		},
		NewFieldSelection("Virtual Instance ID", "virtual_instance_id"),
	},
}

var _ = openapi.CollectionMount{
	CollectionPath:               nil,
	CreatedAt:                    nil,
	Id:                           nil,
	LastRefreshTimeMillis:        nil,
	Rrn:                          nil,
	SnapshotExpirationTimeMillis: nil,
	State:                        nil,
	Stats: &openapi.CollectionMountStats{
		LastQueriedMs: nil,
	},
	VirtualInstanceId:  nil,
	VirtualInstanceRrn: nil,
}
