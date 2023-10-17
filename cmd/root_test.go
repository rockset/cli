package cmd

import (
	"bytes"
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func skipUnlessIntegrationTest(t *testing.T) {
	const keyName = "ROCKSET_APIKEY"
	if os.Getenv(keyName) == "" {
		t.Skipf("could not find %s", keyName)
	}
}

func testWrapper(cmd *cobra.Command, args []string) *bytes.Buffer {
	out := &bytes.Buffer{}
	cmd.Flags().String(FormatFlag, string(format.JSONFormat), "")
	cmd.SetOut(out)
	cmd.SetArgs(args)

	return out
}
