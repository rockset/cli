package test

import (
	"bytes"
	"github.com/rockset/cli/config"
	"github.com/stretchr/testify/require"
	"log"
	"testing"

	"github.com/spf13/cobra"

	"github.com/rockset/cli/format"
)

// Wrapper wraps a cobra.Command and adds the --region and --format flags
func Wrapper(t *testing.T, c *cobra.Command, args ...string) *bytes.Buffer {
	t.Helper()

	out := &bytes.Buffer{}
	// using string values for flags here to avoid an import cycle, extract flags into a separate package to fix?
	require.NoError(t, c.PersistentFlags().Set("format", string(format.JSONFormat)))
	require.NoError(t, c.PersistentFlags().Set("cluster", config.Use1a1))

	c.SetOut(out)
	c.SetArgs(args)

	return out
}

// WrapAndExecute wraps the command with Wrapper and then executes it
func WrapAndExecute(t *testing.T, c *cobra.Command, args ...string) *bytes.Buffer {
	t.Helper()
	log.Printf("args: %+v", args)
	out := Wrapper(t, c, args...)
	require.NoError(t, c.Execute())

	return out
}
