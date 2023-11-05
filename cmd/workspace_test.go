//xgo:build integration

package cmd_test

import (
	"github.com/rockset/cli/flag"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/internal/test"
)

func TestGetWorkspaceCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	c := cmd.NewRootCmd("test")
	out := test.WrapAndExecute(t, c, "get", "workspace", flag.DefaultWorkspace)

	assert.NotEmpty(t, out.String())
}
