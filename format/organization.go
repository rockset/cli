package format

import "github.com/rockset/rockset-go-client/openapi"

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
			Wide:        true,
		},
		{
			DisplayName: "External ID",
			FieldName:   "ExternalId",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
		{
			DisplayName: "Created At",
			FieldName:   "CreatedAt",
			FieldFn:     getFieldByName,
			Wide:        true,
		},
		{
			FieldName: "Clusters",
			FieldFn:   getArrayStructFieldByName("ApiserverUrl"),
			Wide:      true,
		},
	},
}

// just to list available fields
var _ = openapi.Organization{
	Clusters: []openapi.Cluster{
		{
			ApiserverUrl:   nil,
			AwsRegion:      nil,
			ClusterType:    nil,
			Domain:         nil,
			Id:             nil,
			TopLevelDomain: nil,
		},
	},
	CreatedAt:     nil,
	DisplayName:   nil,
	ExternalId:    nil,
	Id:            nil,
	RocksetUser:   nil,
	SsoConnection: nil,
	SsoOnly:       nil,
}
