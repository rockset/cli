package format

import "github.com/rockset/rockset-go-client/openapi"

var QueryLambdaDefaultSelector = DefaultSelector{
	Normal: "Workspace:.workspace,Name:.name,Last Updated By:.last_updated_by,Last Updated:.last_updated,Latest Version:.latest_version.version,Version Count:.version_count,Collections:.collections",
	Wide:   "Workspace:.workspace,Name:.name,Last Updated By:.last_updated_by,Last Updated:.last_updated,Latest Version:.latest_version.version,Description:.latest_version.description,Version Count:.version_count,Collections:.collections",
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
	Normal: "Tag:.tag_name,State:.version.state",
	Wide:   "Tag:.tag_name,Description:.version.description,State:.version.state",
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
