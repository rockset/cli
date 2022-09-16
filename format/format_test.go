package format_test

import (
	"bytes"
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/stretchr/testify/assert"
	"testing"
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
				Roles:     &[]string{"a", "b"},
				State:     openapi.PtrString("state"),
			},
			s: "first name,last name,email,state,created at\n",
		},
		{
			i: openapi.Organization{
				DeletionScheduledAt: openapi.PtrString("deletion at"),
				RocksetUser:         openapi.PtrString("user"),
				ExternalId:          openapi.PtrString("ext id"),
				Clusters:            nil,
				Id:                  openapi.PtrString("id"),
				CreatedAt:           openapi.PtrString("created at"),
				DisplayName:         openapi.PtrString("name"),
			},
			s: "name,user,id,ext id,created at,deletion at\n",
		},
		{
			i: openapi.QueryLambda{
				Workspace:     openapi.PtrString("ws"),
				LastUpdatedBy: openapi.PtrString("updated by"),
				LastUpdated:   openapi.PtrString("last updated"),
				Name:          openapi.PtrString("name"),
				VersionCount:  openapi.PtrInt32(10),
				Collections:   nil,
				LatestVersion: nil,
			},
			s: "ws,name,updated by,last updated,10\n",
		},
		{
			i: openapi.Collection{
				CreatedAt:                         openapi.PtrString("created at"),
				CreatedBy:                         openapi.PtrString("created by"),
				Name:                              openapi.PtrString("name"),
				Description:                       openapi.PtrString("desc"),
				Workspace:                         openapi.PtrString("ws"),
				Status:                            openapi.PtrString("status"),
				Sources:                           nil,
				Stats:                             nil,
				RetentionSecs:                     openapi.PtrInt64(10),
				FieldMappings:                     nil,
				FieldMappingQuery:                 nil,
				ClusteringKey:                     nil,
				Aliases:                           nil,
				FieldSchemas:                      nil,
				InvertedIndexGroupEncodingOptions: nil,
				FieldPartitions:                   nil,
				InsertOnly:                        openapi.PtrBool(true),
				EnableExactlyOnceWrites:           openapi.PtrBool(true),
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
				Segment:     nil,
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
				Roles:     &[]string{"a", "b"},
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
