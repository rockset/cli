package format

import "github.com/rockset/rockset-go-client/openapi"

var OrganizationDefaultSelector = DefaultSelector{
	Normal: []FieldSelection{
		NewFieldSelection("Name", "display_name"),
		NewFieldSelection("ID", "id"),
		NewFieldSelection("User", "rockset_user"),
	},
	Wide: []FieldSelection{
		NewFieldSelection("Name", "display_name"),
		NewFieldSelection("ID", "id"),
		NewFieldSelection("User", "rockset_user"),
		NewFieldSelection("External ID", "external_id"),
		NewFieldSelection("Created At", "created_at"),
	},

	//	Normal: nil, //  "Display Name:.display_name,ID:.id,User:.rockset_user",
	//	Wide:   nil, // "Display Name:.display_name,ID:.id,User:.rockset_user,External ID:.external_id,Created At:.created_at",
}

// just to list available fields
var _ = openapi.Organization{
	Clusters:      nil,
	CreatedAt:     nil,
	DisplayName:   nil,
	ExternalId:    nil,
	Id:            nil,
	RocksetUser:   nil,
	SsoConnection: nil,
	SsoOnly:       nil,
}
