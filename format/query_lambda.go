package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryLambdaDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Last Updated By", "last_updated_by"),
		{
			ColumnName:     "Last Updated",
			Path:           []PathElem{{FieldName: "last_updated"}},
			FieldFormatter: TimeSinceFormatter{}, // TODO this doesn't display correctly
		},
		NewFieldSelection("Latest Version", "latest_version", "version"),
		NewFieldSelection("Version Count", "version_count"),
		NewFieldSelection("Collections", "collections"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Workspace", "workspace"),
		NewFieldSelection("Name", "name"),
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
	Normal: nil, // "Tag:.tag_name,State:.version.state",
	Wide:   nil, // "Tag:.tag_name,Description:.version.description,State:.version.state",
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
