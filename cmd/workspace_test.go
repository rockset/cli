//go:build integration

package cmd

import (
	"bytes"
	"github.com/rockset/cli/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetWorkspaceCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	buf := &bytes.Buffer{}
	cmd := newGetWorkspaceCmd()
	cmd.SetArgs([]string{DefaultWorkspace})
	cmd.Flags().Set("region", "usw2a1")
	cmd.SetOut(buf)

	err := cmd.Execute()

	require.Nil(t, err)
	assert.Equal(t,
		`workspace info: {CreatedAt:2020-02-04T18:53:28Z CreatedBy: Name:commons Description:default workspace CollectionCount:4},
`, buf.String())
}
