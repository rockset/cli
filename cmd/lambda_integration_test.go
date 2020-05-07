// +build integration

package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecuteLambdaCmd(t *testing.T) {
	params := "testdata/params.json"
	buf := bytes.NewBufferString("")
	cmd := newExecuteLambdaCmd()
	cmd.SetArgs([]string{"--params", params, "commons.events"})
	cmd.SetOut(buf)

	err := cmd.Execute()

	require.Nil(t, err)
	assert.Equal(t, ``, buf.String())
}
