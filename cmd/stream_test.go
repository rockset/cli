package cmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/rockset/rockset-go-client/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/flag"
)

func TestStreamDocuments(t *testing.T) {
	ctx := context.TODO()
	//in := bytes.NewBufferString(`{}`)
	in, err := os.Open("testdata/test.json")
	require.NoError(t, err)

	f := &fake{t}

	s := cmd.NewStreamer(f, cmd.StreamConfig{
		Workspace:  flag.DefaultWorkspace,
		Collection: "writetest",
		BatchSize:  3,
	})

	cnt, err := s.Stream(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), cnt)
}

type fake struct {
	t *testing.T
}

func (f *fake) AddDocuments(ctx context.Context, workspace, collection string,
	docs []interface{}) ([]openapi.DocumentStatus, error) {
	f.t.Logf("%d docs added", len(docs))
	res := make([]openapi.DocumentStatus, len(docs))
	added := "ADDED"
	for i := range docs {
		res[i] = openapi.DocumentStatus{
			Status: &added,
		}
	}
	return res, nil
}
