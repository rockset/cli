package format

var CollectionFormatter = StructFormatter{
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
			FieldName: "Description",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Retention",
			FieldName:   "RetentionSecs",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Status",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Insert Only",
			FieldName:   "InsertOnly",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Read Only",
			FieldName:   "ReadOnly",
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
