package format

var UserFormatter = StructFormatter{
	[]Header{
		{
			FieldName:   "FirstName",
			DisplayName: "First Name",
			FieldFn:     getFieldByName,
		},
		{
			FieldName:   "LastName",
			DisplayName: "Last Name",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Email",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "State",
			FieldFn:   getFieldByName,
		},
		{
			FieldName:   "CreatedAt",
			DisplayName: "Created At",
			FieldFn:     getFieldByName,
		},
		{
			FieldName: "Roles",
			Wide:      true,
			FieldFn:   getArrayFieldByName,
		},
	},
}
