package format

var UserDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("First Name", "first_name"),
		NewFieldSelection("Last Name", "last_name"),
		NewFieldSelection("Email", "email"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("Created At", "created_at"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("First Name", "first_name"),
		NewFieldSelection("Last Name", "last_name"),
		NewFieldSelection("Email", "email"),
		NewFieldSelection("State", "state"),
		NewFieldSelection("Created At", "created_at"),
		NewFieldSelection("Roles", "roles"),
	},
}
