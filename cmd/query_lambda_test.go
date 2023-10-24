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

func TestExecuteLambdaCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	params := "testdata/params.json"
	buf := bytes.NewBufferString("")
	c := cmd.NewExecuteQueryLambdaCmd()
	err := c.Flags().Set("region", config.Usw2a1)
	require.NoError(t, err)
	c.SetArgs([]string{"--params", params, "commons.events2"})
	c.SetOut(buf)

	err = c.Execute()

	require.Nil(t, err)
	assert.Equal(t, ``, buf.String())
}
