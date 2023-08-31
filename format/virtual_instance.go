package format

var VirtualInstanceFormatter = StructFormatter{
	[]Header{
		{
			FieldName: "Name",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "Id",
			FieldFn:   getFieldByName,
			Wide:      true,
		},
		{
			FieldName: "Description",
			FieldFn:   getFieldByName,
		},
		{
			FieldName: "State",
			FieldFn:   getFieldByName,
		},
		{
			DisplayName: "Default VI",
			FieldName:   "DefaultVi",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Current Size",
			FieldName:   "CurrentSize",
			FieldFn:     getFieldByName,
		},
		{
			DisplayName: "Desired Size",
			FieldName:   "DesiredSize",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
	},
}
