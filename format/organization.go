package format

var OrgFormatter = StructFormatter{
	[]Header{
		{
			FieldName:   "DisplayName",
			DisplayName: "Name",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Clusters",
			FieldFn:   getArrayFieldByName,
		},
	},
}
