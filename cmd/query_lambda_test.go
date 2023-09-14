package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecuteLambdaCmd(t *testing.T) {
	t.Skipf("skip for now")
	skipUnlessIntegrationTest(t)

	params := "testdata/params.json"
	buf := bytes.NewBufferString("")
	cmd := newExecuteQueryLambdaCmd()
	cmd.Flags().Set("region", "usw2a1")
	cmd.SetArgs([]string{"--params", params, "commons.events2"})
	cmd.SetOut(buf)

	err := cmd.Execute()

	require.Nil(t, err)
	assert.Equal(t, ``, buf.String())
}
