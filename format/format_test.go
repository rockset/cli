package format_test

import (
	"bytes"
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/rockset/rockset-go-client/openapi"
)

func TestFormatter(t *testing.T) {
	var testCases = []struct {
		i interface{}
		s string
	}{
		{
			i: openapi.Workspace{
				CreatedAt:       openapi.PtrString("created at"),
				CreatedBy:       openapi.PtrString("created by"),
				Name:            openapi.PtrString("name"),
				Description:     openapi.PtrString("description"),
				CollectionCount: openapi.PtrInt64(100),
			},
			s: "name,description,100,created by,created at\n",
		},
		{
			i: openapi.User{
				CreatedAt: openapi.PtrString("created at"),
				Email:     "email",
				FirstName: openapi.PtrString("first name"),
				LastName:  openapi.PtrString("last name"),
				Roles:     []string{"a", "b"},
				State:     openapi.PtrString("state"),
			},
			s: "first name,last name,email,state,created at,\"a, b\"\n",
		},
		{
			i: openapi.Organization{
				RocksetUser: openapi.PtrString("user"),
				ExternalId:  openapi.PtrString("ext id"),
				Clusters: []openapi.Cluster{
					{Id: openapi.PtrString("usw2a1")},
					{Id: openapi.PtrString("euc1a1")},
				},
				Id:          openapi.PtrString("id"),
				CreatedAt:   openapi.PtrString("created at"),
				DisplayName: openapi.PtrString("name"),
			},
			s: "name,id,user,ext id,created at\n",
		},
		{
			i: openapi.QueryLambda{
				Workspace:     openapi.PtrString("ws"),
				LastUpdatedBy: openapi.PtrString("updated by"),
				LastUpdated:   openapi.PtrString("last updated"),
				Name:          openapi.PtrString("name"),
				VersionCount:  openapi.PtrInt32(10),
				Collections:   []string{"a", "b"},
				LatestVersion: &openapi.QueryLambdaVersion{
					Version:     openapi.PtrString("version"),
					Description: openapi.PtrString("desc"),
				},
			},
			s: "ws,name,updated by,last updated,version,desc,10,\"a, b\"\n",
		},
		{
			i: openapi.Collection{
				CreatedAt:         openapi.PtrString("created at"),
				CreatedBy:         openapi.PtrString("created by"),
				Name:              openapi.PtrString("name"),
				Description:       openapi.PtrString("desc"),
				Workspace:         openapi.PtrString("ws"),
				Status:            openapi.PtrString("status"),
				Sources:           nil,
				Stats:             nil,
				RetentionSecs:     openapi.PtrInt64(10),
				FieldMappings:     nil,
				FieldMappingQuery: nil,
				ClusteringKey:     nil,
				Aliases:           nil,
				InsertOnly:        openapi.PtrBool(true),
				ReadOnly:          openapi.PtrBool(true),
			},
			s: "ws,name,desc,10,status,true,true,created by,created at\n",
		},
		{
			i: openapi.Integration{
				Name:        "",
				Description: nil,
				Collections: nil,
				CreatedBy:   "",
				CreatedAt:   nil,
				S3:          nil,
				Kinesis:     nil,
				Dynamodb:    nil,
				Gcs:         nil,
				Kafka:       nil,
				Mongodb:     nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.i), func(t *testing.T) {
			buf := bytes.NewBufferString("")
			f := format.FormatterFor(buf, format.CSVFormat, false)
			f.Format(true, tc.i)
			assert.Equal(t, tc.s, buf.String())
		})
	}
}

func TestNewFormatter(t *testing.T) {
	var testCases = []struct {
		i interface{}
		s string
	}{
		{
			i: openapi.User{
				CreatedAt: openapi.PtrString("created at"),
				Email:     "email",
				FirstName: openapi.PtrString("first name"),
				LastName:  openapi.PtrString("last name"),
				Roles:     []string{"a", "b"},
				State:     openapi.PtrString("state"),
			},
			s: "first name,last name,email,state,created at,\"a, b\"\n",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%T", tc.i), func(t *testing.T) {
			buf := bytes.NewBufferString("")
			f := format.FormatterFor(buf, format.CSVFormat, false)
			f.Format(true, tc.i)
			assert.Equal(t, tc.s, buf.String())
		})
	}
}
