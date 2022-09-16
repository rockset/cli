package format

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
		// this is a struct: openapi.QueryLambdaVersion
		//{
		//	DisplayName: "Latest Version",
		//	FieldName:   "LatestVersion",
		//	FieldFn:     getFieldByName,
		//},
		{
			DisplayName: "Version Count",
			FieldName:   "VersionCount",
			FieldFn:     getFieldByName,
		},
	},
}
