package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryLambdaFormatter = StructFormatter{
	[]Header{
		{
			FieldName: "Workspace",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "Name",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Last Updated By",
			FieldName:   "LastUpdatedBy",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Last Updated",
			FieldName:   "LastUpdated",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Latest Version",
			FieldName:   "LatestVersion",
			FieldFn:     getStructFieldByName("Version"),
		},
		{
			DisplayName: "Description",
			FieldName:   "LatestVersion",
			FieldFn:     getStructFieldByName("Description"),
			Wide:        true,
		},
		{
			DisplayName: "Version Count",
			FieldName:   "VersionCount",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Collections",
			FieldFn:   getArrayFieldByName,
		},
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

var QueryLambdaTagFormatter = StructFormatter{
	[]Header{
		{
			FieldName:   "TagName",
			DisplayName: "Tag",
			FieldFn:     getFieldByName,
		},
		{
			FieldName:   "Version",
			DisplayName: "Description",
			FieldFn:     getStructFieldByName("Description"),
			Wide:        true,
		},
		{
			FieldName:   "Version",
			DisplayName: "State",
			FieldFn:     getStructFieldByName("State"),
		},
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
