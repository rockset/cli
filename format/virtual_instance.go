package format

var VirtualInstanceDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("Default VI", "default_vi"),
		NewFieldSelection("Current Size", "current_size"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Name", "name"),
		NewFieldSelection("ID", "id"),
		NewFieldSelection("Description", "description"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("Default VI", "default_vi"),
		NewFieldSelection("Current Size", "current_size"),
		NewFieldSelection("Desired Size", "desired_size"),
	},
}
