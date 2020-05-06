// +build integration

package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetWorkspaceCmd(t *testing.T) {
	buf := bytes.NewBufferString("")
	cmd := newGetWorkspaceCmd()
	cmd.SetArgs([]string{"commons"})
	cmd.SetOut(buf)

	err := cmd.Execute()

	require.Nil(t, err)
	assert.Equal(t, buf.String(),
		`workspace info: {CreatedAt:2020-02-04T18:53:28Z CreatedBy: Name:commons Description:default workspace CollectionCount:4}
`)
}
