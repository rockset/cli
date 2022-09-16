package cmd_test

import (
	"context"
	"github.com/rockset/cli/cmd"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestStreamDocuments(t *testing.T) {
	ctx := context.TODO()
	//in := bytes.NewBufferString(`{}`)
	in, err := os.Open("testdata/test.json")
	require.NoError(t, err)

	f := &fake{}

	s := cmd.NewStreamer(f, cmd.StreamConfig{
		Workspace:  "commons",
		Collection: "writetest",
		BatchSize:  3,
	})

	cnt, err := s.Stream(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), cnt)
}

type fake struct{}

func (f *fake) AddDocuments(ctx context.Context, workspace, collection string,
	docs []interface{}) ([]openapi.DocumentStatus, error) {
	log.Printf("%d docs added", len(docs))
	res := make([]openapi.DocumentStatus, len(docs))
	added := "ADDED"
	for i := range docs {
		res[i] = openapi.DocumentStatus{
			Status: &added,
		}
	}
	return res, nil
}
