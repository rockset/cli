package format

var OrgFormatter = StructFormatter{
	[]Header{
		{
			DisplayName: "DisplayName",
			FieldName:   "DisplayName",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "ID",
			FieldName:   "Id",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "User",
			FieldName:   "RocksetUser",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "External ID",
			FieldName:   "ExternalId",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created At",
			FieldName:   "CreatedAt",
			FieldFn:     getFieldByName,
		},
		// []openapi.Cluster
		//{
		//	FieldName: "Clusters",
		//	FieldFn:   getArrayFieldByName,
		//},
	},
}
