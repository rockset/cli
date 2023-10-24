//go:build integration

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/internal/test"
)

func TestGetWorkspaceCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	buf := &bytes.Buffer{}
	c := cmd.NewGetWorkspaceCmd()
	c.SetArgs([]string{cmd.DefaultWorkspace})
	err := c.Flags().Set("region", config.Use1a1)
	require.NoError(t, err)
	c.SetOut(buf)

	err = c.Execute()

	require.Nil(t, err)
	assert.Equal(t,
		`workspace info: {CreatedAt:2020-02-04T18:53:28Z CreatedBy: Name:commons Description:default workspace CollectionCount:4},
`, buf.String())
}
