//go:build integration

package cmd

import (
	"bytes"
	"github.com/rockset/cli/internal/cluster"
	"github.com/rockset/cli/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecuteLambdaCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	params := "testdata/params.json"
	buf := bytes.NewBufferString("")
	cmd := newExecuteQueryLambdaCmd()
	cmd.Flags().Set("region", cluster.Usw2a1)
	cmd.SetArgs([]string{"--params", params, "commons.events2"})
	cmd.SetOut(buf)

	err := cmd.Execute()

	require.Nil(t, err)
	assert.Equal(t, ``, buf.String())
}
