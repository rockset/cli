//go:build integration

package cmd_test

import (
	"testing"

	"github.com/rockset/cli/cmd"
	"github.com/rockset/cli/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestExecuteLambdaCmd(t *testing.T) {
	test.SkipUnlessIntegrationTest(t)

	c := cmd.NewRootCmd("test")
	// TODO test ql with --param
	out := test.WrapAndExecute(t, c, "execute", "ql", "events2")

	assert.NotEmpty(t, out.String())
}
