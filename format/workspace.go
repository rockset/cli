package format

var WorkspaceFormatter = StructFormatter{
	[]Header{
		{
			FieldName: "Name",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "Description",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Collection Count",
			FieldName:   "CollectionCount",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created By",
			FieldName:   "CreatedBy",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Created At",
			FieldName:   "CreatedAt",
			FieldFn:     getFieldByName,
		},
	},
}
