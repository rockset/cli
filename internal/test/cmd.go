package test

import (
	"bytes"

	"github.com/spf13/cobra"

	"github.com/rockset/cli/format"
)

func Wrapper(c *cobra.Command, args []string) *bytes.Buffer {
	out := &bytes.Buffer{}
	// using "format" here to avoid an import cycle
	// extract flags into a separate package to fix?
	c.Flags().String("format", string(format.JSONFormat), "")
	c.SetOut(out)
	c.SetArgs(args)

	return out
}
