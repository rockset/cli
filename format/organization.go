package format

import "github.com/rockset/rockset-go-client/openapi"

var OrganizationDefaultSelector = DefaultSelector{
	Normal: "Display Name:.display_name,ID:.id,User:.rockset_user",
	Wide:   "Display Name:.display_name,ID:.id,User:.rockset_user,External ID:.external_id,Created At:.created_at",
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
