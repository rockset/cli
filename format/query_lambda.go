package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryLambdaDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("State", "latest_version", "state"),
		{
			ColumnName: "Last Executed",
			Path: []PathElem{
				{FieldName: "latest_version"},
				{FieldName: "stats"},
				{FieldName: "last_executed"},
			},
			FieldFormatter: TimeSinceFormatter{}, // TODO this doesn't display correctly
		},
		{
			ColumnName:     "Last Updated",
			Path:           []PathElem{{FieldName: "last_updated"}},
			FieldFormatter: TimeSinceFormatter{}, // TODO this doesn't display correctly
		},
		NewFieldSelection("Latest Version", "latest_version", "version"),
		NewFieldSelection("Version Count", "version_count"),
		// TODO show the number of collections
		//NewFieldSelection("Collections", "collections"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("State", "latest_version", "state"),
		NewFieldSelection("Last Updated By", "last_updated_by"),
		{
			ColumnName:     "Last Updated",
			Path:           []PathElem{{FieldName: "last_updated"}},
			FieldFormatter: TimeSinceFormatter{},
		},
		NewFieldSelection("Latest Version", "latest_version", "version"),
		NewFieldSelection("Description", "latest_version", "description"),
		NewFieldSelection("Version Count", "version_count"),
		NewFieldSelection("Collections", "collections"),
	},
}

var _ = openapi.QueryLambda{
	Collections:   nil,
	LastUpdated:   nil,
	LastUpdatedBy: nil,
	LatestVersion: &openapi.QueryLambdaVersion{
		Collections:         nil,
		CreatedAt:           nil,
		CreatedBy:           nil,
		CreatedByApikeyName: nil,
		Description:         nil,
		Name:                nil,
		PublicAccessId:      nil,
		Sql: &openapi.QueryLambdaSql{
			DefaultParameters: nil,
			Query:             "",
		},
		State: nil,
		Stats: &openapi.QueryLambdaStats{
			LastExecuted:              nil,
			LastExecutedBy:            nil,
			LastExecutionError:        nil,
			LastExecutionErrorMessage: nil,
		},
		Version:   nil,
		Workspace: nil,
	},
	Name:         nil,
	VersionCount: nil,
	Workspace:    nil,
}

var QueryLambdaTagDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "version", "workspace"),
		NewFieldSelection("Tag", "tag_name"),
		NewFieldSelection("Version", "version", "version"),
		NewFieldSelection("State", "version", "state"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "version", "workspace"),
		NewFieldSelection("Tag", "tag_name"),
		NewFieldSelection("State", "version", "state"),
		NewFieldSelection("Description", "version", "description"),
		NewFieldSelection("Version", "version", "version"),
		NewFieldSelection("State", "version", "state"),
	},
}

var _ = openapi.QueryLambdaTag{
	TagName: nil,
	Version: &openapi.QueryLambdaVersion{
		Collections:         nil,
		CreatedAt:           nil,
		CreatedBy:           nil,
		CreatedByApikeyName: nil,
		Description:         nil,
		Name:                nil,
		PublicAccessId:      nil,
		Sql: &openapi.QueryLambdaSql{
			DefaultParameters: nil,
			Query:             "",
		},
		State: nil,
		Stats: &openapi.QueryLambdaStats{
			LastExecuted:              nil,
			LastExecutedBy:            nil,
			LastExecutionError:        nil,
			LastExecutionErrorMessage: nil,
		},
		Version:   nil,
		Workspace: nil,
	},
}

var QueryLambdaVersionDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Version", "version"),
		NewFieldSelection("State", "state"),
		{
			ColumnName: "Last Executed",
			Path: []PathElem{
				{FieldName: "stats"},
				{FieldName: "last_executed"},
			},
			FieldFormatter: TimeSinceFormatter{}, // TODO this doesn't display correctly
		},
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Version", "version"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("Description", "description"),
	},
}

var _ = openapi.QueryLambdaVersion{
	Collections:         nil,
	CreatedAt:           nil,
	CreatedBy:           nil,
	CreatedByApikeyName: nil,
	Description:         nil,
	Name:                nil,
	PublicAccessId:      nil,
	Sql: &openapi.QueryLambdaSql{
		DefaultParameters: nil,
		Query:             "",
	},
	State: nil,
	Stats: &openapi.QueryLambdaStats{
		LastExecuted:              nil,
		LastExecutedBy:            nil,
		LastExecutionError:        nil,
		LastExecutionErrorMessage: nil,
	},
	Version:   nil,
	Workspace: nil,
}
